package tiles

import (
	"math"
	"testing"
)

func TestTileRoundTrip(t *testing.T) {
	points := []struct {
		name string
		lat  float64
		lon  float64
		z    int
	}{
		{"zoom 0", 0, 0, 0},
		{"zoom 8 tokyo", 35.68, 139.76, 8},
		{"zoom 12 tokyo", 35.68, 139.76, 12},
		{"west of antimeridian", 0, 179.9, 8},
		{"east of antimeridian", 0, -179.9, 8},
	}
	for _, tt := range points {
		t.Run(tt.name, func(t *testing.T) {
			id, ok := TileAt(tt.lat, tt.lon, tt.z)
			if !ok {
				t.Fatal("TileAt rejected the point")
			}
			b, ok := TileBounds(id)
			if !ok {
				t.Fatal("TileBounds rejected the tile")
			}
			if tt.lat < b.MinLat || tt.lat > b.MaxLat || tt.lon < b.MinLon || tt.lon > b.MaxLon {
				t.Fatalf("point %f,%f outside tile bounds %+v", tt.lat, tt.lon, b)
			}
			back, ok := TileAt((b.MinLat+b.MaxLat)/2, (b.MinLon+b.MaxLon)/2, tt.z)
			if !ok || back != id {
				t.Fatalf("center of %v mapped to %v", id, back)
			}
		})
	}
}

func TestAntimeridianTilesDiffer(t *testing.T) {
	west, ok1 := TileAt(0, 179.9, 8)
	east, ok2 := TileAt(0, -179.9, 8)
	if !ok1 || !ok2 {
		t.Fatal("TileAt failed")
	}
	if west.X == east.X {
		t.Fatalf("tiles on opposite sides of the antimeridian share x %d", west.X)
	}
}

func TestTilesIntersectingSplitsAntimeridian(t *testing.T) {
	b := Bounds{MinLat: -1, MaxLat: 1, MinLon: 179, MaxLon: -179}
	got := TilesIntersecting(b, 3)
	if len(got) < 2 {
		t.Fatalf("crossing rectangle produced %d tiles", len(got))
	}
	var sawWest, sawEast bool
	for _, id := range got {
		if id.X == 0 {
			sawEast = true
		}
		if id.X == (1<<3)-1 {
			sawWest = true
		}
	}
	if !sawWest || !sawEast {
		t.Fatalf("tiles = %+v, want both edges of the world", got)
	}
}

func TestPixelIsStable(t *testing.T) {
	id, ok := TileAt(35.6, 139.7, 12)
	if !ok {
		t.Fatal("TileAt")
	}
	x, y, ok := Pixel(35.6, 139.7, id)
	if !ok || x < 0 || x >= TileSize || y < 0 || y >= TileSize {
		t.Fatalf("pixel %d,%d ok=%v", x, y, ok)
	}
	againX, againY, ok := Pixel(35.6, 139.7, id)
	if !ok || againX != x || againY != y {
		t.Fatalf("pixel moved from %d,%d to %d,%d", x, y, againX, againY)
	}
}

func TestQuantizeDay(t *testing.T) {
	// 2024-01-02 15:00 UTC and 18:00 UTC are one cache day.
	afternoon := int64(1704207600)
	evening := afternoon + 3*60*60
	if QuantizeFrom(afternoon) != QuantizeFrom(evening) {
		t.Fatal("same UTC day produced two lower bounds")
	}
	if QuantizeTo(afternoon) != QuantizeTo(evening) {
		t.Fatal("same UTC day produced two upper bounds")
	}
	if QuantizeFrom(0) != 0 || QuantizeTo(0) != 0 {
		t.Fatal("an open bound must stay open")
	}
	next := QuantizeFrom(afternoon) + 24*60*60
	a := cacheKey("aggregate", "1", normalizeRequest(Request{Tile: TileID{Z: 1, X: 1, Y: 1}, From: afternoon, Speed: "all"}))
	b := cacheKey("aggregate", "1", normalizeRequest(Request{Tile: TileID{Z: 1, X: 1, Y: 1}, From: next, Speed: "all"}))
	c := cacheKey("aggregate", "1", normalizeRequest(Request{Tile: TileID{Z: 1, X: 1, Y: 1}, From: afternoon, Speed: "walk"}))
	d := cacheKey("aggregate", "2", normalizeRequest(Request{Tile: TileID{Z: 1, X: 1, Y: 1}, From: afternoon, Speed: "all"}))
	if a == b || a == c || a == d {
		t.Fatalf("cache keys collided: %s %s %s %s", a, b, c, d)
	}
	if math.Abs(float64(QuantizeTo(afternoon)-QuantizeFrom(afternoon))) > 24*60*60 {
		t.Fatal("quantized window is longer than a day")
	}
}
