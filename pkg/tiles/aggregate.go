package tiles

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"

	"safecast-new-map/pkg/database"
)

// Aggregate paints the highest dose in each pixel of a tile.
// A pixel with no zoom=0 measurement stays transparent. Opacity is never
// used to suggest coverage: a faded color would look like a milder dose.
type Aggregate struct {
	DB     *database.Database
	DBType string
}

// Name is the URL path segment for this processor.
func (a *Aggregate) Name() string { return "aggregate" }

// ContentType is the media type Render returns.
func (a *Aggregate) ContentType() string { return "image/png" }

// Version participates in the cache key.
func (a *Aggregate) Version() string { return "1" }

// Render reads the raw points inside the tile and encodes a paletted PNG.
func (a *Aggregate) Render(ctx context.Context, req Request) (Result, error) {
	if a == nil || a.DB == nil {
		return Result{}, fmt.Errorf("tiles: aggregate database is not configured")
	}
	bounds, ok := TileBounds(req.Tile)
	if !ok {
		return Result{}, fmt.Errorf("tiles: tile %d/%d/%d is outside the XYZ grid", req.Tile.Z, req.Tile.X, req.Tile.Y)
	}
	points, err := a.Load(ctx, bounds, req.From, req.To)
	if err != nil {
		return Result{}, err
	}
	body, err := EncodePNG(Bin(points, req.Tile))
	if err != nil {
		return Result{}, err
	}
	return Result{Body: body, ContentType: a.ContentType()}, nil
}

// Point is the narrow row the aggregate needs. The full marker payload stays
// in the database and is still available through the existing marker APIs.
type Point struct {
	ID          int64
	Lat         float64
	Lon         float64
	Dose        float64
	Date        int64
	TrackID     string
	HasSpectrum bool
}

// Cell is the summary of one pixel. Count and Tracks are carried for the
// click probe. They are not drawn, and they are not a confidence interval.
type Cell struct {
	Occupied bool
	Dose     float64
	Winner   Point
	Count    int
	Tracks   int
	Latest   int64
	trackIDs map[string]struct{}
}

// Grid is the 256 by 256 aggregation of one tile, row-major, north to south.
type Grid struct {
	Cells [TileSize * TileSize]Cell
}

// Bin keeps the hottest point in each pixel. Equal doses keep the later
// timestamp, then the smaller id, so the probe is stable.
func Bin(points []Point, id TileID) Grid {
	var g Grid
	for _, p := range points {
		x, y, ok := Pixel(p.Lat, p.Lon, id)
		if !ok {
			continue
		}
		c := &g.Cells[y*TileSize+x]
		if c.trackIDs == nil {
			c.trackIDs = map[string]struct{}{}
		}
		c.Count++
		if p.TrackID != "" {
			c.trackIDs[p.TrackID] = struct{}{}
		}
		c.Tracks = len(c.trackIDs)
		if p.Date > c.Latest {
			c.Latest = p.Date
		}
		if !c.Occupied || hotter(p, c.Winner) {
			c.Occupied = true
			c.Dose = p.Dose
			c.Winner = p
		}
	}
	return g
}

func hotter(p, prev Point) bool {
	if p.Dose != prev.Dose {
		return p.Dose > prev.Dose
	}
	if p.Date != prev.Date {
		return p.Date > prev.Date
	}
	return p.ID < prev.ID
}

// ProbeHit is the raw measurement a click resolves to.
type ProbeHit struct {
	Found       bool    `json:"found"`
	ID          int64   `json:"id,omitempty"`
	DoseRate    float64 `json:"doseRate,omitempty"`
	Date        int64   `json:"date,omitempty"`
	TrackID     string  `json:"trackID,omitempty"`
	HasSpectrum bool    `json:"hasSpectrum,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	Count       int     `json:"count,omitempty"`
	Tracks      int     `json:"tracks,omitempty"`
	Latest      int64   `json:"latest,omitempty"`
}

// Probe returns the winning raw point in the pixel that contains lat, lon.
func (a *Aggregate) Probe(ctx context.Context, lat, lon float64, z int, from, to int64) (ProbeHit, error) {
	id, ok := TileAt(lat, lon, z)
	if !ok {
		return ProbeHit{}, fmt.Errorf("tiles: probe location is outside the map")
	}
	bounds, ok := TileBounds(id)
	if !ok {
		return ProbeHit{}, fmt.Errorf("tiles: probe tile is invalid")
	}
	points, err := a.Load(ctx, bounds, from, to)
	if err != nil {
		return ProbeHit{}, err
	}
	x, y, ok := Pixel(lat, lon, id)
	if !ok {
		return ProbeHit{}, nil
	}
	cell := Bin(points, id).Cells[y*TileSize+x]
	if !cell.Occupied {
		return ProbeHit{}, nil
	}
	return ProbeHit{
		Found:       true,
		ID:          cell.Winner.ID,
		DoseRate:    cell.Winner.Dose,
		Date:        cell.Winner.Date,
		TrackID:     cell.Winner.TrackID,
		HasSpectrum: cell.Winner.HasSpectrum,
		Lat:         cell.Winner.Lat,
		Lon:         cell.Winner.Lon,
		Count:       cell.Count,
		Tracks:      cell.Tracks,
		Latest:      cell.Latest,
	}, nil
}

// Load reads zoom=0 points inside bounds through the web-read lane.
// The predicate is zoom plus a latitude/longitude range so the existing
// (zoom, lat, lon) btree can bound the scan. Date bounds are optional.
func (a *Aggregate) Load(ctx context.Context, bounds Bounds, from, to int64) ([]Point, error) {
	query, args := pointQuery(a.DBType, bounds, from, to)
	var points []Point
	err := a.DB.WithWebRead(ctx, func(jobCtx context.Context, conn *sql.DB) error {
		rows, err := conn.QueryContext(jobCtx, query, args...)
		if err != nil {
			return fmt.Errorf("tiles: query points: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var p Point
			var has any
			if err := rows.Scan(&p.ID, &p.Lat, &p.Lon, &p.Dose, &p.Date, &p.TrackID, &has); err != nil {
				return fmt.Errorf("tiles: scan point: %w", err)
			}
			p.HasSpectrum = asBool(has)
			points = append(points, p)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("tiles: iterate points: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return points, nil
}

func pointQuery(dbType string, b Bounds, from, to int64) (string, []any) {
	dose, track, spectrum := "doseRate", "trackID", "COALESCE(has_spectrum, 0)"
	if strings.EqualFold(dbType, "pgx") {
		dose, track, spectrum = "doserate", "trackid", "COALESCE(has_spectrum, FALSE)"
	} else if strings.EqualFold(dbType, "duckdb") {
		spectrum = "COALESCE(has_spectrum, FALSE)"
	}
	var args []any
	ph := func() string {
		args = append(args, nil)
		if strings.EqualFold(dbType, "pgx") {
			return fmt.Sprintf("$%d", len(args))
		}
		return "?"
	}
	// Fill placeholders in order after the SQL text is chosen.
	zoom := ph()
	minLat, maxLat, minLon, maxLon := ph(), ph(), ph(), ph()
	args[0] = 0
	args[1] = b.MinLat
	args[2] = b.MaxLat
	args[3] = b.MinLon
	args[4] = b.MaxLon
	q := fmt.Sprintf(`SELECT id, lat, lon, %s, date, COALESCE(%s, ''), %s
FROM markers
WHERE zoom = %s AND lat BETWEEN %s AND %s AND lon BETWEEN %s AND %s`,
		dose, track, spectrum, zoom, minLat, maxLat, minLon, maxLon)
	if from > 0 {
		p := ph()
		args[len(args)-1] = from
		q += " AND date >= " + p
	}
	if to > 0 {
		p := ph()
		args[len(args)-1] = to
		q += " AND date <= " + p
	}
	return q, args
}

func asBool(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case int64:
		return t != 0
	case int:
		return t != 0
	case []byte:
		s := string(t)
		return s == "1" || s == "t" || s == "true"
	case string:
		return t == "1" || t == "t" || t == "true"
	default:
		return false
	}
}

// EncodePNG writes a paletted image. Index 0 is fully transparent, so empty
// pixels do not blend the dose color with the base map.
func EncodePNG(g Grid) ([]byte, error) {
	pal := color.Palette{color.NRGBA{0, 0, 0, 0}}
	indexOf := map[[3]uint8]uint8{ /* transparent is 0 and is not looked up */ }
	img := image.NewPaletted(image.Rect(0, 0, TileSize, TileSize), pal)
	for i := range g.Cells {
		c := g.Cells[i]
		if !c.Occupied {
			continue
		}
		r, gch, b := ColorForDose(c.Dose)
		key := [3]uint8{r, gch, b}
		idx, ok := indexOf[key]
		if !ok {
			if len(img.Palette) >= 256 {
				idx = nearest(img.Palette, r, gch, b)
			} else {
				idx = uint8(len(img.Palette))
				img.Palette = append(img.Palette, color.NRGBA{R: r, G: gch, B: b, A: 255})
				indexOf[key] = idx
			}
		}
		img.SetColorIndex(i%TileSize, i/TileSize, idx)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("tiles: encode png: %w", err)
	}
	return buf.Bytes(), nil
}

func nearest(pal color.Palette, r, g, b uint8) uint8 {
	best := uint8(1)
	bestD := 1 << 30
	for i := 1; i < len(pal); i++ {
		cr, cg, cb, _ := pal[i].RGBA()
		dr := int(cr>>8) - int(r)
		dg := int(cg>>8) - int(g)
		db := int(cb>>8) - int(b)
		d := dr*dr + dg*dg + db*db
		if d < bestD {
			bestD = d
			best = uint8(i)
		}
	}
	return best
}
