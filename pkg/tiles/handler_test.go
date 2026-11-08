package tiles

import (
	"context"
	"encoding/json"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	_ "modernc.org/sqlite"
	"safecast-new-map/pkg/database"
)

func TestHandlerAggregatePNGAndProbe(t *testing.T) {
	db := openTileDB(t)
	inside := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC).Unix()
	outside := time.Date(2024, 1, 4, 12, 0, 0, 0, time.UTC).Unix()
	insertPoint(t, db, 1, 0.1, inside, 35.6, 139.7, "road", false)
	insertPoint(t, db, 2, 0.5, inside, 35.6, 139.7, "road", false)
	insertPoint(t, db, 3, 0.2, inside, 35.6, 139.7, "other", true)
	insertPoint(t, db, 4, 9, outside, 35.6, 139.7, "spike", false)

	svc, err := New(db, "sqlite", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	svc.Register(mux)

	id, ok := TileAt(35.6, 139.7, 12)
	if !ok {
		t.Fatal("TileAt")
	}
	filtered := "/tiles/aggregate/" + strconv.Itoa(id.Z) + "/" + strconv.Itoa(id.X) + "/" + strconv.Itoa(id.Y) + ".png?from=" + strconv.FormatInt(inside, 10) + "&to=" + strconv.FormatInt(inside, 10)
	rec := httptest.NewRequest(http.MethodGet, filtered, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, rec)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	if rr.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("content type %q", rr.Header().Get("Content-Type"))
	}
	if rr.Header().Get("Cache-Control") != "public, max-age=60" {
		t.Fatalf("cache header %q", rr.Header().Get("Cache-Control"))
	}
	img, err := png.Decode(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	x, y, ok := Pixel(35.6, 139.7, id)
	if !ok {
		t.Fatal("pixel")
	}
	assertPixel(t, img, x, y, 0.5, quantizeAlpha(alphaForCount(3)))
	originX := x / CellPixels * CellPixels
	originY := y / CellPixels * CellPixels
	assertPixel(t, img, originX+CellPixels-1, originY, 0.5, quantizeAlpha(alphaForCount(3)))
	nextX := originX + CellPixels
	if nextX >= TileSize {
		nextX = originX - 1
	}
	_, _, _, a := img.At(nextX, originY).RGBA()
	if a != 0 {
		t.Fatal("pixel outside the cell is not transparent")
	}

	probeURL := "/tiles/aggregate/probe?lat=35.6&lon=139.7&z=12&from=" + strconv.FormatInt(inside, 10) + "&to=" + strconv.FormatInt(inside, 10)
	pr := httptest.NewRecorder()
	mux.ServeHTTP(pr, httptest.NewRequest(http.MethodGet, probeURL, nil))
	if pr.Code != http.StatusOK {
		t.Fatalf("probe %d %s", pr.Code, pr.Body.String())
	}
	var hit ProbeHit
	if err := json.Unmarshal(pr.Body.Bytes(), &hit); err != nil {
		t.Fatal(err)
	}
	if !hit.Found || hit.ID != 2 || hit.DoseRate != 0.5 || hit.TrackID != "road" || hit.HasSpectrum {
		t.Fatalf("probe %+v", hit)
	}
	if hit.Count != 3 || hit.Tracks != 2 {
		t.Fatalf("probe support %+v", hit)
	}

	miss := httptest.NewRecorder()
	mux.ServeHTTP(miss, httptest.NewRequest(http.MethodGet, "/tiles/nope/1/0/0.png", nil))
	if miss.Code != http.StatusNotFound {
		t.Fatalf("unknown processor status %d", miss.Code)
	}
}

func TestWarmTileIsServedFromCache(t *testing.T) {
	db := openTileDB(t)
	insertPoint(t, db, 1, 0.4, 1_700_000_000, 35.6, 139.7, "road", false)
	svc, err := New(db, "sqlite", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	svc.RegisterProcessor(&countingAggregate{Aggregate: *svc.agg, calls: &calls})
	id, _ := TileAt(35.6, 139.7, 12)
	req := Request{Tile: id}
	if _, err := svc.Render(context.Background(), "aggregate", req); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Render(context.Background(), "aggregate", req); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("renders %d, want 1 warm cache hit", calls)
	}
}

func TestZoom13TileIsServed(t *testing.T) {
	if MaxOverviewZoom != 13 {
		t.Fatalf("overview stops at zoom %d, want 13", MaxOverviewZoom)
	}
	db := openTileDB(t)
	insertPoint(t, db, 1, 0.4, 1_700_000_000, 35.6, 139.7, "road", false)
	svc, err := New(db, "sqlite", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	svc.Register(mux)
	id, ok := TileAt(35.6, 139.7, 13)
	if !ok {
		t.Fatal("TileAt")
	}
	url := "/tiles/aggregate/" + strconv.Itoa(id.Z) + "/" + strconv.Itoa(id.X) + "/" + strconv.Itoa(id.Y) + ".png"
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, url, nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("zoom 13 status %d %s", rr.Code, rr.Body.String())
	}
	if rr.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("content type %q", rr.Header().Get("Content-Type"))
	}
}

type countingAggregate struct {
	Aggregate
	calls *int
}

func (c *countingAggregate) Name() string        { return "aggregate" }
func (c *countingAggregate) Version() string     { return "1" }
func (c *countingAggregate) ContentType() string { return "image/png" }
func (c *countingAggregate) Render(ctx context.Context, req Request) (Result, error) {
	*c.calls++
	return c.Aggregate.Render(ctx, req)
}

func assertPixel(t *testing.T, img image.Image, x, y int, dose float64, alpha uint8) {
	t.Helper()
	pr, pg, pb, pa := straight(img.At(x, y))
	wantR, wantG, wantB := ColorForDose(dose)
	if pr != wantR || pg != wantG || pb != wantB || pa != alpha {
		t.Fatalf("pixel rgba(%d,%d,%d,%d), want rgb(%d,%d,%d) a=%d", pr, pg, pb, pa, wantR, wantG, wantB, alpha)
	}
}

func TestSpeedFilterDropsAirplanes(t *testing.T) {
	db := openTileDB(t)
	const lat, lon = 35.6, 139.7
	insertSpeed := func(id int, dose, speed float64, track string, nullSpeed bool) {
		t.Helper()
		var sp any = speed
		if nullSpeed {
			sp = nil
		}
		_, err := db.DB.ExecContext(context.Background(),
			`INSERT INTO markers (id, doseRate, date, lon, lat, countRate, zoom, speed, trackID, has_spectrum) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, dose, int64(1_700_000_000), lon, lat, 1.0, 0, sp, track, 0)
		if err != nil {
			t.Fatal(err)
		}
	}
	insertSpeed(1, 0.10, 0.5, "walk", false)
	insertSpeed(2, 0.40, 10, "car", false)
	insertSpeed(3, 5.00, 80, "plane", false)
	insertSpeed(4, 0.15, 0, "blank", true)
	insertSpeed(5, 9.00, -1, "live", false)

	svc, err := New(db, "sqlite", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	svc.Register(mux)
	id, ok := TileAt(lat, lon, 12)
	if !ok {
		t.Fatal("TileAt")
	}
	fetch := func(speed string) image.Image {
		t.Helper()
		url := "/tiles/aggregate/" + strconv.Itoa(id.Z) + "/" + strconv.Itoa(id.X) + "/" + strconv.Itoa(id.Y) + ".png?speed=" + speed
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, url, nil))
		if rr.Code != http.StatusOK {
			t.Fatalf("speed %s status %d %s", speed, rr.Code, rr.Body.String())
		}
		img, err := png.Decode(rr.Body)
		if err != nil {
			t.Fatal(err)
		}
		return img
	}
	x, y, ok := Pixel(lat, lon, id)
	if !ok {
		t.Fatal("pixel")
	}
	// Default map filters: cars and pedestrians. The plane's 5 µSv/h must not win.
	assertPixel(t, fetch("cw"), x, y, 0.40, quantizeAlpha(alphaForCount(3)))
	assertPixel(t, fetch("a"), x, y, 5.00, quantizeAlpha(alphaForCount(1)))
	// A missing speed is walking, so it beats the 0.10 walking sample.
	assertPixel(t, fetch("w"), x, y, 0.15, quantizeAlpha(alphaForCount(2)))
	_, _, _, clear := straight(fetch("none").At(x, y))
	if clear != 0 {
		t.Fatal("speed=none still painted a cell")
	}

	pr := httptest.NewRecorder()
	mux.ServeHTTP(pr, httptest.NewRequest(http.MethodGet, "/tiles/aggregate/probe?lat=35.6&lon=139.7&z=12&speed=cw", nil))
	var hit ProbeHit
	if err := json.Unmarshal(pr.Body.Bytes(), &hit); err != nil {
		t.Fatal(err)
	}
	if !hit.Found || hit.ID != 2 || hit.DoseRate != 0.40 || hit.Count != 3 {
		t.Fatalf("probe %+v", hit)
	}
}

func openTileDB(t *testing.T) *database.Database {
	t.Helper()
	dir := t.TempDir()
	cfg := database.Config{DBType: "sqlite", DBPath: filepath.Join(dir, "tiles.db"), Port: 1}
	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.InitSchema(cfg); err != nil {
		t.Fatal(err)
	}
	return db
}

func insertPoint(t *testing.T, db *database.Database, id int, dose float64, date int64, lat, lon float64, track string, spectrum bool) {
	t.Helper()
	spec := 0
	if spectrum {
		spec = 1
	}
	_, err := db.DB.ExecContext(context.Background(),
		`INSERT INTO markers (id, doseRate, date, lon, lat, countRate, zoom, speed, trackID, has_spectrum) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, dose, date, lon, lat, 1.0, 0, 0.0, track, spec)
	if err != nil {
		t.Fatal(err)
	}
}
