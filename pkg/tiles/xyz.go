package tiles

import "math"

const (
	// TileSize is the pixel width of one XYZ tile. Leaflet's default is the same.
	TileSize = 256
	// PointZoom is the first zoom at which the map draws individual markers.
	// Overview tiles are served only below this zoom.
	PointZoom = 13
	// MaxOverviewZoom is the finest cached overview tile.
	MaxOverviewZoom = PointZoom - 1
)

// webMercatorLimit is the latitude Web Mercator can represent.
const webMercatorLimit = 85.0511287798066

// TileID is one XYZ tile. X increases east and Y increases south.
type TileID struct {
	Z int
	X int
	Y int
}

// Bounds is a geographic rectangle. Longitudes are in degrees, west to east,
// except when the rectangle crosses the antimeridian and MinLon > MaxLon.
type Bounds struct {
	MinLat float64
	MinLon float64
	MaxLat float64
	MaxLon float64
}

// Valid reports whether the rectangle has area and lies on Earth.
// A rectangle that crosses the antimeridian is valid when MinLon > MaxLon.
func (b Bounds) Valid() bool {
	if b.MinLat >= b.MaxLat || b.MinLat < -90 || b.MaxLat > 90 {
		return false
	}
	if b.MinLon < -180 || b.MaxLon > 180 || b.MinLon > 180 || b.MaxLon < -180 {
		return false
	}
	return b.MinLon != b.MaxLon
}

// TileAt returns the tile that contains the point at zoom z.
func TileAt(lat, lon float64, z int) (TileID, bool) {
	if z < 0 || z > 22 || lat < -webMercatorLimit || lat > webMercatorLimit || lon < -180 || lon > 180 {
		return TileID{}, false
	}
	n := float64(uint(1) << uint(z))
	x := int(math.Floor((lon + 180) / 360 * n))
	y := int(math.Floor(webY(lat) * n))
	max := int(n) - 1
	if x > max {
		x = max
	}
	if y > max {
		y = max
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	return TileID{Z: z, X: x, Y: y}, true
}

// TileBounds returns the geographic rectangle covered by id.
func TileBounds(id TileID) (Bounds, bool) {
	if id.Z < 0 || id.Z > 22 {
		return Bounds{}, false
	}
	n := 1 << uint(id.Z)
	if id.X < 0 || id.Y < 0 || id.X >= n || id.Y >= n {
		return Bounds{}, false
	}
	nf := float64(n)
	west := float64(id.X)/nf*360 - 180
	east := float64(id.X+1)/nf*360 - 180
	north := latFromWebY(float64(id.Y) / nf)
	south := latFromWebY(float64(id.Y+1) / nf)
	return Bounds{MinLat: south, MaxLat: north, MinLon: west, MaxLon: east}, true
}

// Pixel returns the tile-local pixel that contains the point.
// The east and south edges belong to the next tile, so a point is in one pixel.
func Pixel(lat, lon float64, id TileID) (x, y int, ok bool) {
	b, good := TileBounds(id)
	if !good {
		return 0, 0, false
	}
	lastX := (1 << uint(id.Z)) - 1
	if lon < b.MinLon || lat > b.MaxLat || lat < b.MinLat || lon > b.MaxLon {
		return 0, 0, false
	}
	if lon == b.MaxLon && id.X < lastX {
		return 0, 0, false
	}
	if lat == b.MinLat && id.Y < lastX {
		return 0, 0, false
	}
	x = int(((lon - b.MinLon) / (b.MaxLon - b.MinLon)) * TileSize)
	north := webY(b.MaxLat)
	south := webY(b.MinLat)
	span := south - north
	if span == 0 {
		return 0, 0, false
	}
	y = int(((webY(lat) - north) / span) * TileSize)
	if x >= TileSize {
		x = TileSize - 1
	}
	if y >= TileSize {
		y = TileSize - 1
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	return x, y, true
}

// TilesIntersecting lists the tiles at zoom z that meet b.
// A rectangle that crosses the antimeridian (MinLon > MaxLon) is split.
func TilesIntersecting(b Bounds, z int) []TileID {
	if z < 0 || z > 22 || b.MinLat > b.MaxLat {
		return nil
	}
	if b.MinLon > b.MaxLon {
		left := b
		left.MaxLon = 180
		right := b
		right.MinLon = -180
		return append(TilesIntersecting(left, z), TilesIntersecting(right, z)...)
	}
	north := math.Min(b.MaxLat, webMercatorLimit)
	south := math.Max(b.MinLat, -webMercatorLimit)
	if south > north {
		return nil
	}
	west := math.Max(b.MinLon, -180)
	east := math.Min(b.MaxLon, 180)
	nw, ok1 := TileAt(north, west, z)
	se, ok2 := TileAt(south, east, z)
	if !ok1 || !ok2 {
		return nil
	}
	out := make([]TileID, 0, (se.X-nw.X+1)*(se.Y-nw.Y+1))
	for x := nw.X; x <= se.X; x++ {
		for y := nw.Y; y <= se.Y; y++ {
			out = append(out, TileID{Z: z, X: x, Y: y})
		}
	}
	return out
}

// Overlaps reports whether the two rectangles share any area.
func Overlaps(a, b Bounds) bool {
	if a.MinLon > a.MaxLon {
		left := a
		left.MaxLon = 180
		right := a
		right.MinLon = -180
		return Overlaps(left, b) || Overlaps(right, b)
	}
	if b.MinLon > b.MaxLon {
		return Overlaps(b, a)
	}
	return a.MinLat <= b.MaxLat && a.MaxLat >= b.MinLat && a.MinLon <= b.MaxLon && a.MaxLon >= b.MinLon
}

func webY(lat float64) float64 {
	if lat > webMercatorLimit {
		lat = webMercatorLimit
	}
	if lat < -webMercatorLimit {
		lat = -webMercatorLimit
	}
	// artanh(sin φ) is the mercator projection. The 1/2 in artanh is required;
	// dropping it shifts mid-latitude cities by thousands of kilometres.
	s := math.Sin(lat * math.Pi / 180)
	merc := 0.5 * math.Log((1+s)/(1-s))
	return (1 - merc/math.Pi) / 2
}

func latFromWebY(y float64) float64 {
	return math.Atan(math.Sinh(math.Pi-2*math.Pi*y)) * 180 / math.Pi
}
