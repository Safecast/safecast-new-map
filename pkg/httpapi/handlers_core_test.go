// handlers_core_test.go tests the core API handlers (overview, latest, tracks, track data, countries, shorten) with a temp SQLite DB.
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"safecast-new-map/pkg/database"
	_ "safecast-new-map/pkg/database/drivers"
)

const (
	coreTestTrackID  = "test-api-track-001"
	coreTestMarkerID = 42
)

// newTestHandler builds a Handler backed by a temporary SQLite database with seeded tracks/markers.
func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	cfg := database.Config{DBType: "sqlite", DBPath: dbPath, Port: 1}
	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}
	if err := db.InitSchema(cfg); err != nil {
		t.Fatalf("InitSchema: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	seedCoreTestDB(t, db)
	return NewHandler(db, "sqlite", nil, nil, nil, "")
}

// seedCoreTestDB inserts one track and one marker for core API tests.
func seedCoreTestDB(t *testing.T, db *database.Database) {
	t.Helper()
	ctx := context.Background()
	_, err := db.DB.ExecContext(ctx, `INSERT INTO tracks (trackID) VALUES (?)`, coreTestTrackID)
	if err != nil {
		t.Fatalf("seed tracks: %v", err)
	}
	_, err = db.DB.ExecContext(ctx,
		`INSERT INTO markers (id, doseRate, date, lon, lat, countRate, zoom, speed, trackID)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		coreTestMarkerID, 0.12, time.Now().Unix(), 139.7, 35.6, 8.0, 10, 0.0, coreTestTrackID)
	if err != nil {
		t.Fatalf("seed markers: %v", err)
	}
}

func TestHandleOverview(t *testing.T) {
	h := newTestHandler(t)
	mux := http.NewServeMux()
	h.Register(mux)
	t.Run("GET returns JSON overview", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("got status %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if _, ok := out["endpoints"]; !ok {
			t.Error("response missing endpoints")
		}
		if _, ok := out["totalTracks"]; !ok {
			t.Error("response missing totalTracks")
		}
	})
}

func TestHandleLatestNearby(t *testing.T) {
	h := newTestHandler(t)
	mux := http.NewServeMux()
	h.Register(mux)
	t.Run("missing lat/lon returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/latest", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})
	t.Run("invalid lat returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/latest?lat=abc&lon=139.7", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})
	t.Run("lat out of range returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/latest?lat=95&lon=0", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})
	t.Run("valid request returns markers array", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/latest?lat=35.6&lon=139.7&radius_m=50000", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if _, ok := out["markers"]; !ok {
			t.Error("response missing markers")
		}
		if _, ok := out["center"]; !ok {
			t.Error("response missing center")
		}
	})
	t.Run("POST returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/latest", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got %d, want 405", rec.Code)
		}
	})
}

func TestHandleTracksList(t *testing.T) {
	h := newTestHandler(t)
	mux := http.NewServeMux()
	h.Register(mux)
	t.Run("returns track list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/tracks", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if _, ok := out["tracks"]; !ok {
			t.Error("response missing tracks")
		}
	})
}

func TestHandleTrackData(t *testing.T) {
	h := newTestHandler(t)
	mux := http.NewServeMux()
	h.Register(mux)
	t.Run("valid track ID returns data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/track/"+coreTestTrackID+".json", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
	})
	t.Run("missing track ID returns 400 or 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/track/", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest && rec.Code != http.StatusNotFound {
			t.Errorf("got %d, want 400 or 404", rec.Code)
		}
	})
}

func TestHandleCountries(t *testing.T) {
	h := newTestHandler(t)
	mux := http.NewServeMux()
	h.Register(mux)
	t.Run("returns countries", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/countries", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if _, ok := out["countries"]; !ok {
			t.Error("response missing countries")
		}
	})
}

func TestHandleShorten(t *testing.T) {
	h := newTestHandler(t)
	mux := http.NewServeMux()
	h.Register(mux)
	t.Run("GET returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/shorten", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got %d, want 405", rec.Code)
		}
	})
	t.Run("invalid JSON returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader("not json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})
	t.Run("missing url returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":""}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})
	t.Run("preview returns short code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten",
			strings.NewReader(`{"url":"https://example.com/map?lat=35.6&lon=139.7","commit":false}`))
		req.Header.Set("Content-Type", "application/json")
		req.Host = "example.com"
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if _, ok := out["code"]; !ok {
			t.Error("response missing code")
		}
		if _, ok := out["short"]; !ok {
			t.Error("response missing short")
		}
	})
	t.Run("commit stores short link", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten",
			strings.NewReader(`{"url":"https://example.com/map?lat=35.6&lon=139.7","commit":true}`))
		req.Header.Set("Content-Type", "application/json")
		req.Host = "example.com"
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if stored, _ := out["stored"].(bool); !stored {
			t.Error("expected stored=true on commit")
		}
	})
	t.Run("foreign host rejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten",
			strings.NewReader(`{"url":"https://evil.com/path","commit":false}`))
		req.Header.Set("Content-Type", "application/json")
		req.Host = "example.com"
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})
	t.Run("relative path accepted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten",
			strings.NewReader(`{"url":"/map?lat=35.6&lon=139.7","commit":false}`))
		req.Header.Set("Content-Type", "application/json")
		req.Host = "example.com"
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
	})
}

func TestHandleTracksByYear(t *testing.T) {
	h := newTestHandler(t)
	mux := http.NewServeMux()
	h.Register(mux)
	t.Run("returns tracks for year", func(t *testing.T) {
		year := time.Now().UTC().Format("2006")
		req := httptest.NewRequest(http.MethodGet, "/api/tracks/years/"+year, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
	})
}
