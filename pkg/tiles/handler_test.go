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
	assertColor(t, img, x, y, 0.5)
	nx := x + 1
	if nx >= TileSize {
		nx = x - 1
	}
	_, _, _, a := img.At(nx, y).RGBA()
	if a != 0 {
		t.Fatal("neighbor pixel is not transparent")
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

func assertColor(t *testing.T, img image.Image, x, y int, dose float64) {
	t.Helper()
	pr, pg, pb, pa := img.At(x, y).RGBA()
	wantR, wantG, wantB := ColorForDose(dose)
	if pa>>8 != 255 || uint8(pr>>8) != wantR || uint8(pg>>8) != wantG || uint8(pb>>8) != wantB {
		t.Fatalf("pixel rgba(%d,%d,%d,%d), want rgb(%d,%d,%d) opaque", pr>>8, pg>>8, pb>>8, pa>>8, wantR, wantG, wantB)
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
