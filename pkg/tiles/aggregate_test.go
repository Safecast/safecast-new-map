package tiles

import (
	"bytes"
	"image/png"
	"testing"
)

func TestBinKeepsMaxAndLeavesNeighborsEmpty(t *testing.T) {
	id, ok := TileAt(35.6, 139.7, 12)
	if !ok {
		t.Fatal("TileAt")
	}
	x, y, ok := Pixel(35.6, 139.7, id)
	if !ok {
		t.Fatal("Pixel")
	}
	points := []Point{
		{ID: 1, Lat: 35.6, Lon: 139.7, Dose: 0.1, Date: 10, TrackID: "a"},
		{ID: 2, Lat: 35.6, Lon: 139.7, Dose: 0.5, Date: 11, TrackID: "a"},
		{ID: 3, Lat: 35.6, Lon: 139.7, Dose: 0.2, Date: 12, TrackID: "b"},
	}
	grid := Bin(points, id)
	cell := grid.Cells[y*TileSize+x]
	if !cell.Occupied || cell.Dose != 0.5 || cell.Winner.ID != 2 {
		t.Fatalf("winner = %+v", cell)
	}
	if cell.Count != 3 || cell.Tracks != 2 {
		t.Fatalf("count=%d tracks=%d", cell.Count, cell.Tracks)
	}
	nx := x + 1
	if nx >= TileSize {
		nx = x - 1
	}
	if grid.Cells[y*TileSize+nx].Occupied {
		t.Fatal("neighbor pixel was painted")
	}

	body, err := EncodePNG(grid)
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	pr, pg, pb, pa := img.At(x, y).RGBA()
	wantR, wantG, wantB := ColorForDose(0.5)
	if pa>>8 != 255 || uint8(pr>>8) != wantR || uint8(pg>>8) != wantG || uint8(pb>>8) != wantB {
		t.Fatalf("pixel rgba(%d,%d,%d,%d), want rgb(%d,%d,%d) opaque", pr>>8, pg>>8, pb>>8, pa>>8, wantR, wantG, wantB)
	}
	_, _, _, na := img.At(nx, y).RGBA()
	if na != 0 {
		t.Fatal("empty pixel is not transparent")
	}
}

func TestRoadPaintsOnlyItsPixels(t *testing.T) {
	id, ok := TileAt(35.6, 139.7, 12)
	if !ok {
		t.Fatal("TileAt")
	}
	var points []Point
	for i := 0; i < 8; i++ {
		lat, lon := pointInPixel(t, id, 20+i, 40)
		points = append(points, Point{ID: int64(i + 1), Lat: lat, Lon: lon, Dose: 0.2, Date: int64(i), TrackID: "road"})
	}
	grid := Bin(points, id)
	painted := 0
	for _, c := range grid.Cells {
		if c.Occupied {
			painted++
		}
	}
	if painted != 8 {
		t.Fatalf("painted %d pixels, want 8", painted)
	}
}

func pointInPixel(t *testing.T, id TileID, x, y int) (lat, lon float64) {
	t.Helper()
	b, ok := TileBounds(id)
	if !ok {
		t.Fatal("bounds")
	}
	lon = b.MinLon + (float64(x)+0.5)/TileSize*(b.MaxLon-b.MinLon)
	north := webY(b.MaxLat)
	south := webY(b.MinLat)
	wy := north + (float64(y)+0.5)/TileSize*(south-north)
	lat = latFromWebY(wy)
	px, py, ok := Pixel(lat, lon, id)
	if !ok || px != x || py != y {
		t.Fatalf("reconstructed pixel %d,%d, want %d,%d", px, py, x, y)
	}
	return lat, lon
}
