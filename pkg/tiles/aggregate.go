package tiles

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"strings"

	"safecast-new-map/pkg/database"
)

// Aggregate reduces the measurements in each Web Mercator cell to one value
// and paints that cell. The color is the reduced dose. The opacity is how
// many measurements fell in the cell, on a logarithmic scale, so a single
// reading does not look like a survey of the whole cell. Opacity is density,
// not a confidence interval. A cell with no measurement stays transparent.
type Aggregate struct {
	DB     *database.Database
	DBType string
}

// Name is the URL path segment for this processor.
func (a *Aggregate) Name() string { return "aggregate" }

// ContentType is the media type Render returns.
func (a *Aggregate) ContentType() string { return "image/png" }

// Version participates in the cache key. Bump it when the painted image changes
// so a cached stamp from an older renderer is not reused.
func (a *Aggregate) Version() string { return "4" }

// Render reads the raw points inside the tile and encodes a paletted PNG.
func (a *Aggregate) Render(ctx context.Context, req Request) (Result, error) {
	if a == nil || a.DB == nil {
		return Result{}, fmt.Errorf("tiles: aggregate database is not configured")
	}
	bounds, ok := TileBounds(req.Tile)
	if !ok {
		return Result{}, fmt.Errorf("tiles: tile %d/%d/%d is outside the XYZ grid", req.Tile.Z, req.Tile.X, req.Tile.Y)
	}
	points, err := a.Load(ctx, bounds, req.From, req.To, req.Speed)
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

// CellPixels is the screen width of one aggregation cell.
// It divides TileSize, so a cell lies on the global Web Mercator pixel grid
// and is never split by a tile edge. The geographic size is CellPixels
// world-pixels at that zoom, which halves each time the user zooms in.
const CellPixels = 8

const cellsPerAxis = TileSize / CellPixels

// Cell is one aggregation cell. Dose is the reduced value (the maximum for
// the default reducer). Count is the number of measurements and is what
// the painter turns into opacity. Tracks and Winner feed the click probe.
type Cell struct {
	Occupied bool
	Dose     float64
	Winner   Point
	Count    int
	Tracks   int
	Latest   int64
	reduce   Reducer
	trackIDs map[string]struct{}
}

// Reducer combines the points in one cell into the dose that is drawn.
// Count stays outside the reducer so a later mean or median can replace
// the maximum without changing the density scale.
type Reducer interface {
	Observe(Point)
	Dose() float64
	Winner() Point
}

type maxReducer struct {
	set    bool
	dose   float64
	winner Point
}

func newReducer() Reducer { return &maxReducer{} }

func (r *maxReducer) Observe(p Point) {
	if !r.set || hotter(p, r.winner) {
		r.set = true
		r.dose = p.Dose
		r.winner = p
	}
}

func (r *maxReducer) Dose() float64 { return r.dose }

func (r *maxReducer) Winner() Point { return r.winner }

func (c *Cell) observe(p Point) {
	if c.reduce == nil {
		c.reduce = newReducer()
	}
	c.reduce.Observe(p)
	c.Dose = c.reduce.Dose()
	c.Winner = c.reduce.Winner()
	c.Occupied = true
	c.Count++
	if c.trackIDs == nil {
		c.trackIDs = map[string]struct{}{}
	}
	if p.TrackID != "" {
		c.trackIDs[p.TrackID] = struct{}{}
	}
	c.Tracks = len(c.trackIDs)
	if p.Date > c.Latest {
		c.Latest = p.Date
	}
}

// Grid is the aggregation of one tile. Cells are indexed in global Web
// Mercator order, not by the viewport that requested the tile.
type Grid struct {
	Zoom    int
	OriginX int
	OriginY int
	Cells   [cellsPerAxis * cellsPerAxis]Cell
}

// CellIndex is the global cell that contains the point at zoom z.
// The same point returns the same index for every tile and every pan.
func CellIndex(lat, lon float64, z int) (x, y int, ok bool) {
	wx, wy, ok := worldPixel(lat, lon, z)
	if !ok {
		return 0, 0, false
	}
	return int(wx) / CellPixels, int(wy) / CellPixels, true
}

func worldPixel(lat, lon float64, z int) (x, y float64, ok bool) {
	if z < 0 || z > 22 || lat < -webMercatorLimit || lat > webMercatorLimit || lon < -180 || lon > 180 {
		return 0, 0, false
	}
	n := float64(uint(TileSize) << uint(z))
	x = (lon + 180) / 360 * n
	y = webY(lat) * n
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= n {
		x = n - 1
	}
	if y >= n {
		y = n - 1
	}
	return x, y, true
}

// At returns the cell in this tile for a global cell index.
func (g Grid) At(cx, cy int) (Cell, bool) {
	lx := cx - g.OriginX
	ly := cy - g.OriginY
	if lx < 0 || ly < 0 || lx >= cellsPerAxis || ly >= cellsPerAxis {
		return Cell{}, false
	}
	return g.Cells[ly*cellsPerAxis+lx], true
}

// Bin assigns each point to the global cell that contains it and reduces
// that cell. Points that belong to a neighboring tile are skipped, so a
// shared edge is owned by exactly one tile.
func Bin(points []Point, id TileID) Grid {
	g := Grid{
		Zoom:    id.Z,
		OriginX: id.X * cellsPerAxis,
		OriginY: id.Y * cellsPerAxis,
	}
	for _, p := range points {
		cx, cy, ok := CellIndex(p.Lat, p.Lon, id.Z)
		if !ok {
			continue
		}
		lx := cx - g.OriginX
		ly := cy - g.OriginY
		if lx < 0 || ly < 0 || lx >= cellsPerAxis || ly >= cellsPerAxis {
			continue
		}
		g.Cells[ly*cellsPerAxis+lx].observe(p)
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

// Probe returns the winning raw point in the aggregation cell that contains lat, lon.
// speed uses the same bands as the map checkboxes; see Request.Speed.
func (a *Aggregate) Probe(ctx context.Context, lat, lon float64, z int, from, to int64, speed string) (ProbeHit, error) {
	id, ok := TileAt(lat, lon, z)
	if !ok {
		return ProbeHit{}, fmt.Errorf("tiles: probe location is outside the map")
	}
	bounds, ok := TileBounds(id)
	if !ok {
		return ProbeHit{}, fmt.Errorf("tiles: probe tile is invalid")
	}
	points, err := a.Load(ctx, bounds, from, to, speed)
	if err != nil {
		return ProbeHit{}, err
	}
	grid := Bin(points, id)
	cx, cy, ok := CellIndex(lat, lon, id.Z)
	if !ok {
		return ProbeHit{}, nil
	}
	cell, ok := grid.At(cx, cy)
	if !ok || !cell.Occupied {
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
// (zoom, lat, lon) btree can bound the scan. Date bounds and the speed
// band are optional. speed "all" does not constrain speed.
func (a *Aggregate) Load(ctx context.Context, bounds Bounds, from, to int64, speed string) ([]Point, error) {
	query, args := pointQuery(a.DBType, bounds, from, to, speed)
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

func pointQuery(dbType string, b Bounds, from, to int64, speed string) (string, []any) {
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
	q = appendSpeed(q, cleanSpeed(speed), ph, &args)
	return q, args
}

// Speeds match shouldDisplayBySpeed in map.html. The column is metres per second.
const (
	speedWalkMax  = 1.94  // 7 km/h
	speedPlaneMin = 55.56 // 200 km/h
)

// appendSpeed keeps the rows the map's speed checkboxes would draw.
// A missing speed counts as walking, the same way a marker with no speed does.
// Negative speeds are realtime sensors and stay out of these bands.
func appendSpeed(q, token string, ph func() string, args *[]any) string {
	switch token {
	case "all":
		return q
	case "none":
		return q + " AND 0 = 1"
	}
	var parts []string
	if strings.Contains(token, "w") {
		p := ph()
		(*args)[len(*args)-1] = speedWalkMax
		parts = append(parts, fmt.Sprintf("(COALESCE(speed, 0) >= 0 AND COALESCE(speed, 0) < %s)", p))
	}
	if strings.Contains(token, "c") {
		lo, hi := ph(), ph()
		(*args)[len(*args)-2] = speedWalkMax
		(*args)[len(*args)-1] = speedPlaneMin
		parts = append(parts, fmt.Sprintf("(speed >= %s AND speed < %s)", lo, hi))
	}
	if strings.Contains(token, "a") {
		p := ph()
		(*args)[len(*args)-1] = speedPlaneMin
		parts = append(parts, fmt.Sprintf("(speed >= %s)", p))
	}
	if len(parts) == 0 {
		return q + " AND 0 = 1"
	}
	return q + " AND (" + strings.Join(parts, " OR ") + ")"
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

// alphaCountCap is the observation count at which a cell becomes fully opaque.
// Counts beyond it do not get darker, so a dense road cannot use up the scale.
const alphaCountCap = 32

// alphaSingle is the opacity of a cell that contains one measurement.
// It stays well below opaque: one reading is not a survey of the cell.
const alphaSingle = 64

// alphaForCount maps an observation count to an opacity in 0..255.
// Zero stays invisible. From one observation up to alphaCountCap the scale
// is logarithmic. The result is density, not a confidence interval.
func alphaForCount(n int) uint8 {
	if n <= 0 {
		return 0
	}
	if n >= alphaCountCap {
		return 255
	}
	t := math.Log(float64(n)) / math.Log(float64(alphaCountCap))
	a := float64(alphaSingle) + (255-float64(alphaSingle))*t
	return uint8(math.Round(a))
}

// quantizeAlpha collapses opacity onto 16-step palette entries so a tile
// with many doses still fits in a 256-color PNG.
func quantizeAlpha(a uint8) uint8 {
	if a == 0 {
		return 0
	}
	const step = 16
	q := (int(a) + step/2) / step * step
	if q > 255 {
		q = 255
	}
	if q == 0 {
		q = step
	}
	return uint8(q)
}

// EncodePNG paints each occupied cell as one solid square. Empty cells stay
// fully transparent. Color comes from the reduced dose and opacity from the
// observation count; the two are not mixed.
func EncodePNG(g Grid) ([]byte, error) {
	pal := color.Palette{color.NRGBA{0, 0, 0, 0}}
	indexOf := map[[4]uint8]uint8{}
	img := image.NewPaletted(image.Rect(0, 0, TileSize, TileSize), pal)
	for ly := 0; ly < cellsPerAxis; ly++ {
		for lx := 0; lx < cellsPerAxis; lx++ {
			c := g.Cells[ly*cellsPerAxis+lx]
			if !c.Occupied {
				continue
			}
			cr, cg, cb := ColorForDose(c.Dose)
			ca := quantizeAlpha(alphaForCount(c.Count))
			key := [4]uint8{cr, cg, cb, ca}
			idx, ok := indexOf[key]
			if !ok {
				if len(img.Palette) >= 256 {
					idx = nearest(img.Palette, cr, cg, cb, ca)
				} else {
					idx = uint8(len(img.Palette))
					img.Palette = append(img.Palette, color.NRGBA{R: cr, G: cg, B: cb, A: ca})
					indexOf[key] = idx
				}
			}
			x0 := lx * CellPixels
			y0 := ly * CellPixels
			for y := y0; y < y0+CellPixels; y++ {
				for x := x0; x < x0+CellPixels; x++ {
					img.SetColorIndex(x, y, idx)
				}
			}
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("tiles: encode png: %w", err)
	}
	return buf.Bytes(), nil
}

func nearest(pal color.Palette, r, g, b, a uint8) uint8 {
	best := uint8(1)
	bestD := 1 << 30
	for i := 1; i < len(pal); i++ {
		cr, cg, cb, ca := pal[i].RGBA()
		dr := int(cr>>8) - int(r)
		dg := int(cg>>8) - int(g)
		db := int(cb>>8) - int(b)
		da := int(ca>>8) - int(a)
		d := dr*dr + dg*dg + db*db + da*da
		if d < bestD {
			bestD = d
			best = uint8(i)
		}
	}
	return best
}
