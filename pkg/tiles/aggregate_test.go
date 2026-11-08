package tiles

import (
	"bytes"
	"image/color"
	"image/png"
	"testing"
)

func TestBinKeepsMaxInOneCell(t *testing.T) {
	id, ok := TileAt(35.6, 139.7, 12)
	if !ok {
		t.Fatal("TileAt")
	}
	points := []Point{
		{ID: 1, Lat: 35.6, Lon: 139.7, Dose: 0.1, Date: 10, TrackID: "a"},
		{ID: 2, Lat: 35.6, Lon: 139.7, Dose: 0.5, Date: 11, TrackID: "a"},
		{ID: 3, Lat: 35.6, Lon: 139.7, Dose: 0.2, Date: 12, TrackID: "b"},
	}
	grid := Bin(points, id)
	cx, cy, ok := CellIndex(35.6, 139.7, id.Z)
	if !ok {
		t.Fatal("CellIndex")
	}
	cell, ok := grid.At(cx, cy)
	if !ok || !cell.Occupied || cell.Dose != 0.5 || cell.Winner.ID != 2 {
		t.Fatalf("winner = %+v ok=%v", cell, ok)
	}
	if cell.Count != 3 || cell.Tracks != 2 {
		t.Fatalf("count=%d tracks=%d", cell.Count, cell.Tracks)
	}
	occupied := 0
	for _, c := range grid.Cells {
		if c.Occupied {
			occupied++
		}
	}
	if occupied != 1 {
		t.Fatalf("occupied cells %d, want 1", occupied)
	}

	body, err := EncodePNG(grid)
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	px, py, ok := Pixel(35.6, 139.7, id)
	if !ok {
		t.Fatal("Pixel")
	}
	pr, pg, pb, pa := straight(img.At(px, py))
	wantR, wantG, wantB := ColorForDose(0.5)
	wantA := quantizeAlpha(alphaForCount(3))
	if pr != wantR || pg != wantG || pb != wantB || pa != wantA {
		t.Fatalf("pixel rgba(%d,%d,%d,%d), want rgb(%d,%d,%d) a=%d", pr, pg, pb, pa, wantR, wantG, wantB, wantA)
	}
	// The whole cell is one value. The next cell is empty.
	originX := px / CellPixels * CellPixels
	originY := py / CellPixels * CellPixels
	er, eg, eb, ea := straight(img.At(originX+CellPixels-1, originY+CellPixels-1))
	if er != pr || eg != pg || eb != pb || ea != pa {
		t.Fatal("cell is not a uniform square")
	}
	nextX := originX + CellPixels
	if nextX >= TileSize {
		nextX = originX - 1
	}
	_, _, _, na := img.At(nextX, originY).RGBA()
	if na != 0 {
		t.Fatal("neighbor cell is not transparent")
	}
}

func TestCellIndexIsStableAcrossTilesAndPans(t *testing.T) {
	const lat, lon, z = 35.6, 139.7, 10
	cx, cy, ok := CellIndex(lat, lon, z)
	if !ok {
		t.Fatal("CellIndex")
	}
	againX, againY, _ := CellIndex(lat, lon, z)
	if cx != againX || cy != againY {
		t.Fatal("panning changed the cell")
	}
	id, ok := TileAt(lat, lon, z)
	if !ok {
		t.Fatal("TileAt")
	}
	point := Point{ID: 1, Lat: lat, Lon: lon, Dose: 0.4, TrackID: "road"}
	own := Bin([]Point{point}, id)
	cell, ok := own.At(cx, cy)
	if !ok || cell.Count != 1 || cell.Dose != 0.4 {
		t.Fatalf("owning tile cell %+v ok=%v", cell, ok)
	}
	east := Bin([]Point{point}, TileID{Z: z, X: id.X + 1, Y: id.Y})
	if _, ok := east.At(cx, cy); ok {
		t.Fatal("east tile claimed a cell it does not contain")
	}
	for _, c := range east.Cells {
		if c.Occupied {
			t.Fatal("point was aggregated into the neighboring tile")
		}
	}
}

func TestZoomChangesCellSize(t *testing.T) {
	lat, lon := 35.6, 139.7
	cx, cy, ok := CellIndex(lat, lon, 8)
	if !ok {
		t.Fatal("CellIndex")
	}
	latA, lonA := lonLatInCell(cx, cy, 8, 0.2, 0.5)
	latB, lonB := lonLatInCell(cx, cy, 8, 0.6, 0.5)
	ax, ay, _ := CellIndex(latA, lonA, 8)
	bx, by, _ := CellIndex(latB, lonB, 8)
	if ax != cx || ay != cy || bx != cx || by != cy {
		t.Fatalf("zoom 8 split a cell: (%d,%d) (%d,%d) origin (%d,%d)", ax, ay, bx, by, cx, cy)
	}
	hx, hy, _ := CellIndex(latA, lonA, 12)
	hx2, hy2, _ := CellIndex(latB, lonB, 12)
	if hx == hx2 && hy == hy2 {
		t.Fatal("zoom 12 still merged points that share a zoom 8 cell")
	}
	if cellLonDegrees(8) <= cellLonDegrees(12)*8 {
		t.Fatalf("cell longitude zoom8=%v zoom12=%v", cellLonDegrees(8), cellLonDegrees(12))
	}
}

func TestOpacityFollowsCountNotDose(t *testing.T) {
	if alphaForCount(0) != 0 {
		t.Fatal("empty cell is not transparent")
	}
	if alphaForCount(1) >= alphaForCount(8) || alphaForCount(8) >= alphaForCount(alphaCountCap) {
		t.Fatalf("opacity did not increase with count: 1=%d 8=%d cap=%d", alphaForCount(1), alphaForCount(8), alphaForCount(alphaCountCap))
	}
	if alphaForCount(alphaCountCap) != 255 || alphaForCount(alphaCountCap*10) != 255 {
		t.Fatal("dense cells did not saturate")
	}
	if alphaForCount(1) > 96 {
		t.Fatalf("a single observation is too opaque: %d", alphaForCount(1))
	}

	id, _ := TileAt(35.6, 139.7, 11)
	sparseLat, sparseLon := pointInCell(t, id, 2, 2, 0.2, 0.2)
	denseLat, denseLon := pointInCell(t, id, 4, 4, 0.5, 0.5)
	points := []Point{{ID: 1, Lat: sparseLat, Lon: sparseLon, Dose: 0.5}}
	for i := 0; i < 20; i++ {
		lat, lon := pointInCell(t, id, 4, 4, 0.1+float64(i)*0.03, 0.1)
		points = append(points, Point{ID: int64(i + 2), Lat: lat, Lon: lon, Dose: 0.5})
	}
	grid := Bin(points, id)
	sx, sy, _ := CellIndex(sparseLat, sparseLon, id.Z)
	dx, dy, _ := CellIndex(denseLat, denseLon, id.Z)
	sparse, _ := grid.At(sx, sy)
	dense, _ := grid.At(dx, dy)
	if sparse.Count != 1 || dense.Count != 20 {
		t.Fatalf("sparse=%d dense=%d", sparse.Count, dense.Count)
	}
	if sparse.Dose != dense.Dose {
		t.Fatalf("doses differ: sparse %v dense %v", sparse.Dose, dense.Dose)
	}
	body, err := EncodePNG(grid)
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	spx, spy, _ := Pixel(sparseLat, sparseLon, id)
	dpx, dpy, _ := Pixel(denseLat, denseLon, id)
	sr, sg, sb, sa := straight(img.At(spx, spy))
	dr, dg, db, da := straight(img.At(dpx, dpy))
	if sr != dr || sg != dg || sb != db {
		t.Fatal("count changed the dose color")
	}
	if sa >= da || sa == 0 {
		t.Fatalf("sparse alpha %d dense alpha %d", sa, da)
	}
}

func TestDenseRoadCollapsesToCells(t *testing.T) {
	id, ok := TileAt(35.65, 139.75, 11)
	if !ok {
		t.Fatal("TileAt")
	}
	var points []Point
	// Closely spaced readings along one row of cells, plus one hotspot.
	for i := 0; i < 800; i++ {
		lx := (i / 40) % cellsPerAxis
		lat, lon := pointInCell(t, id, lx, 3, 0.05+float64(i%40)*0.02, 0.5)
		points = append(points, Point{ID: int64(i + 1), Lat: lat, Lon: lon, Dose: 0.15, TrackID: "road"})
	}
	hotLat, hotLon := pointInCell(t, id, 1, 1, 0.5, 0.5)
	for i := 0; i < 40; i++ {
		lat, lon := pointInCell(t, id, 1, 1, 0.05+float64(i)*0.02, 0.2)
		points = append(points, Point{ID: int64(1000 + i), Lat: lat, Lon: lon, Dose: 0.1 + float64(i)*0.01, TrackID: "hot"})
	}
	grid := Bin(points, id)
	occupied := 0
	for _, c := range grid.Cells {
		if c.Occupied {
			occupied++
		}
	}
	if occupied >= len(points)/10 {
		t.Fatalf("occupied %d cells for %d points", occupied, len(points))
	}
	hx, hy, _ := CellIndex(hotLat, hotLon, id.Z)
	hot, ok := grid.At(hx, hy)
	if !ok || hot.Count < 40 || hot.Dose < 0.4 {
		t.Fatalf("hotspot %+v ok=%v", hot, ok)
	}
	body, err := EncodePNG(grid)
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	// Every painted pixel belongs to a uniform CellPixels square.
	for y := 0; y < TileSize; y += CellPixels {
		for x := 0; x < TileSize; x += CellPixels {
			r0, g0, b0, a0 := straight(img.At(x, y))
			r1, g1, b1, a1 := straight(img.At(x+CellPixels-1, y+CellPixels-1))
			if r0 != r1 || g0 != g1 || b0 != b1 || a0 != a1 {
				t.Fatalf("cell %d,%d is not uniform", x, y)
			}
		}
	}
}

func straight(c color.Color) (r, g, b, a uint8) {
	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	return n.R, n.G, n.B, n.A
}

func cellLonDegrees(z int) float64 {
	n := float64(uint(TileSize) << uint(z))
	return float64(CellPixels) / n * 360
}

func lonLatInCell(cx, cy, z int, fx, fy float64) (lat, lon float64) {
	n := float64(uint(TileSize) << uint(z))
	wx := (float64(cx) + fx) * CellPixels
	wy := (float64(cy) + fy) * CellPixels
	lon = wx/n*360 - 180
	lat = latFromWebY(wy / n)
	return lat, lon
}

// pointInCell returns a lat/lon inside the tile-local cell (lx, ly),
// fx and fy in 0..1 within that cell.
func pointInCell(t *testing.T, id TileID, lx, ly int, fx, fy float64) (lat, lon float64) {
	t.Helper()
	px := lx*CellPixels + int(fx*float64(CellPixels-1))
	py := ly*CellPixels + int(fy*float64(CellPixels-1))
	return pointInPixel(t, id, px, py)
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
