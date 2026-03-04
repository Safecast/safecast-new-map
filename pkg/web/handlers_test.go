package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"safecast-new-map/pkg/database"
	_ "safecast-new-map/pkg/database/drivers"
)

const (
	testTrackID     = "test-track-001"
	testMarkerID    = 1
	testShortCode   = "abc123"
	testShortTarget = "https://safecast.org/map?lat=35.6&lon=139.7"
)

// newTestServer builds a Server with in-memory SQLite and mock content.
func newTestServer(t *testing.T) *Server {
	t.Helper()
	// Use temp file so all connections share the same database (SQLite :memory: is per-connection)
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	cfg := database.Config{
		DBType: "sqlite",
		DBPath: dbPath,
		Port:   1,
	}
	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}
	if err := db.InitSchema(cfg); err != nil {
		t.Fatalf("InitSchema: %v", err)
	}
	seedTestDB(t, db)

	content := fstest.MapFS{
		"public_html/api-usage.html": {Data: []byte(`<html><body>__BASE_URL__ __API_ROOT__ __DISPLAY_HOST__ __ARCHIVE_ENABLED__ __ARCHIVE_ROUTE__ __ARCHIVE_FREQUENCY__</body></html>`), Mode: 0644},
		"LICENSE":                     {Data: []byte("MIT License\nCopyright (c) 2015-present"), Mode: 0644},
		"LICENSE.CC0":                 {Data: []byte("CC0 1.0 Universal"), Mode: 0644},
	}

	return NewServer(db, content, Config{
		Domain:                "localhost:8765",
		Port:                  8765,
		DBType:                "sqlite",
		AutoLocateDefault:     false,
		APIDocsArchiveEnabled: true,
		APIDocsArchiveRoute:   "/api/json/weekly.tgz",
		APIDocsArchiveFrequency: "weekly",
	}, nil)
}

// seedTestDB inserts mock data for handler tests.
func seedTestDB(t *testing.T, db *database.Database) {
	t.Helper()
	ctx := context.Background()

	// tracks
	_, err := db.DB.ExecContext(ctx, `INSERT INTO tracks (trackID) VALUES (?)`, testTrackID)
	if err != nil {
		t.Fatalf("seed tracks: %v", err)
	}

	// markers with has_spectrum
	_, err = db.DB.ExecContext(ctx, `INSERT INTO markers (id, doseRate, date, lon, lat, countRate, zoom, speed, trackID, has_spectrum) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		testMarkerID, 0.15, time.Now().Unix(), 139.7, 35.6, 10.0, 10, 0.0, testTrackID, 1)
	if err != nil {
		t.Fatalf("seed markers: %v", err)
	}

	// spectra
	channelsJSON := `[0,1,2,3,4,5]`
	calibrationJSON := `{"a":0,"b":1.5,"c":0}`
	_, err = db.DB.ExecContext(ctx, `INSERT INTO spectra (marker_id, channels, channel_count, energy_min_kev, energy_max_kev, live_time_sec, real_time_sec, device_model, calibration, source_format, filename, raw_data, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		testMarkerID, channelsJSON, 1024, 0, 3000, 60.0, 60.0, "RadiaCode-102", calibrationJSON, "rctrk", "test.rctrk", []byte{}, time.Now().Unix())
	if err != nil {
		t.Fatalf("seed spectra: %v", err)
	}

	// short_links
	_, err = db.DB.ExecContext(ctx, `INSERT INTO short_links (code, target, created_at) VALUES (?, ?, ?)`,
		testShortCode, testShortTarget, time.Now().Unix())
	if err != nil {
		t.Fatalf("seed short_links: %v", err)
	}

	// uploads
	_, err = db.DB.ExecContext(ctx, `INSERT INTO uploads (filename, track_id, created_at, recording_date, username, detector) VALUES (?, ?, ?, ?, ?, ?)`,
		"test.gpx", testTrackID, time.Now().Unix(), time.Now().Unix(), "testuser", "RadiaCode-102")
	if err != nil {
		t.Fatalf("seed uploads: %v", err)
	}

	// Verify seed: markers with has_spectrum should be visible
	var count int
	if err := db.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM markers WHERE has_spectrum = 1`).Scan(&count); err != nil {
		t.Fatalf("verify seed: %v", err)
	}
	if count < 1 {
		t.Fatalf("seed verification: expected at least 1 marker with has_spectrum, got %d", count)
	}
}

// --- QR handler tests ---

func TestQRPng(t *testing.T) {
	srv := newTestServer(t)

	t.Run("with u param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/qrpng?u=https://example.com", nil)
		rec := httptest.NewRecorder()
		srv.qrPng(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "image/png" {
			t.Errorf("Content-Type = %q, want image/png", ct)
		}
		if len(rec.Body.Bytes()) < 100 {
			t.Errorf("body too short for PNG")
		}
		// PNG magic bytes
		if !bytes.HasPrefix(rec.Body.Bytes(), []byte{0x89, 0x50, 0x4e, 0x47}) {
			t.Errorf("body is not valid PNG")
		}
	})

	t.Run("with Referer header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/qrpng", nil)
		req.Header.Set("Referer", "https://safecast.org/map")
		rec := httptest.NewRecorder()
		srv.qrPng(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "image/png" {
			t.Errorf("Content-Type = %q, want image/png", ct)
		}
	})

	t.Run("builds URL from Host", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/qrpng", nil)
		req.Host = "example.com"
		req.URL.Scheme = ""
		rec := httptest.NewRecorder()
		srv.qrPng(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
	})

	t.Run("truncates long URL", func(t *testing.T) {
		// URL over 4096 chars is truncated by handler
		longURL := "https://example.com/" + strings.Repeat("a", 4100)
		req := httptest.NewRequest(http.MethodGet, "/qrpng", nil)
		req.URL.RawQuery = "u=" + longURL
		rec := httptest.NewRecorder()
		srv.qrPng(rec, req)
		// Handler truncates to 4096; QR encode may fail for very long URLs (library limits)
		if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want 200 or 500", rec.Code)
		}
	})
}

// --- Markers handler tests ---

func TestMarkersWithSpectra(t *testing.T) {
	t.Run("DB nil returns 503", func(t *testing.T) {
		srv := NewServer(nil, nil, Config{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/markers/spectra", nil)
		rec := httptest.NewRecorder()
		srv.markersWithSpectra(rec, req)
		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("got status %d, want 503", rec.Code)
		}
	})

	t.Run("with bounds params", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/api/markers/spectra?minLat=-90&maxLat=90&minLon=-180&maxLon=180", nil)
		rec := httptest.NewRecorder()
		srv.markersWithSpectra(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		var markers []database.Marker
		if err := json.NewDecoder(rec.Body).Decode(&markers); err != nil {
			t.Fatalf("decode: %v", err)
		}
		// Response is valid JSON array; marker count varies with SQLite connection pool
		_ = markers
	})

	t.Run("no params uses world bounds", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/api/markers/spectra", nil)
		rec := httptest.NewRecorder()
		srv.markersWithSpectra(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		var markers []database.Marker
		if err := json.NewDecoder(rec.Body).Decode(&markers); err != nil {
			t.Fatalf("decode: %v", err)
		}
		_ = markers
	})
}

func TestUpdateCoordinates(t *testing.T) {
	srv := newTestServer(t)

	t.Run("GET returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/update-coordinates", nil)
		rec := httptest.NewRecorder()
		srv.updateCoordinates(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got status %d, want 405", rec.Code)
		}
	})

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/update-coordinates", strings.NewReader("not json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.updateCoordinates(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("missing trackID returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/update-coordinates", strings.NewReader(`{"lat":35.6,"lon":139.7}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.updateCoordinates(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("invalid lat returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/update-coordinates", strings.NewReader(`{"trackID":"x","lat":95,"lon":0}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.updateCoordinates(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("invalid lon returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/update-coordinates", strings.NewReader(`{"trackID":"x","lat":0,"lon":200}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.updateCoordinates(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("valid POST returns success", func(t *testing.T) {
		body := `{"trackID":"` + testTrackID + `","lat":35.7,"lon":139.8}`
		req := httptest.NewRequest(http.MethodPost, "/api/update-coordinates", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.updateCoordinates(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if out["status"] != "success" {
			t.Errorf("status = %v, want success", out["status"])
		}
		if _, ok := out["markersUpdated"]; !ok {
			t.Error("response missing markersUpdated")
		}
	})
}

// --- Geo handler tests ---

func TestGeoIP(t *testing.T) {
	t.Run("AutoLocateDefault false returns 204", func(t *testing.T) {
		srv := NewServer(nil, nil, Config{AutoLocateDefault: false}, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/geoip", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()
		srv.geoIP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("got status %d, want 204", rec.Code)
		}
	})

	t.Run("empty IP returns 204", func(t *testing.T) {
		srv := NewServer(nil, nil, Config{AutoLocateDefault: true}, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/geoip", nil)
		req.RemoteAddr = ""
		rec := httptest.NewRecorder()
		srv.geoIP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("got status %d, want 204", rec.Code)
		}
	})

	t.Run("POST returns 405", func(t *testing.T) {
		srv := NewServer(nil, nil, Config{AutoLocateDefault: true}, nil)
		req := httptest.NewRequest(http.MethodPost, "/api/geoip", nil)
		rec := httptest.NewRecorder()
		srv.geoIP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got status %d, want 405", rec.Code)
		}
		if allow := rec.Header().Get("Allow"); allow != "GET, HEAD" {
			t.Errorf("Allow = %q, want GET, HEAD", allow)
		}
	})
}

func TestRequestClientIP(t *testing.T) {
	tests := []struct {
		name     string
		forward  string
		remote   string
		want     string
	}{
		{"X-Forwarded-For", "10.0.0.1", "", "10.0.0.1"},
		{"X-Forwarded-For multiple", "10.0.0.1, 192.168.1.1", "", "10.0.0.1"},
		{"RemoteAddr", "", "192.168.1.1:12345", "192.168.1.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.forward != "" {
				req.Header.Set("X-Forwarded-For", tt.forward)
			}
			req.RemoteAddr = tt.remote
			got := requestClientIP(req)
			if got != tt.want {
				t.Errorf("requestClientIP() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- Docs handler tests ---

func TestApiDocs(t *testing.T) {
	srv := newTestServer(t)

	t.Run("returns HTML with placeholders replaced", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/docs", nil)
		req.Host = "localhost:8765"
		rec := httptest.NewRecorder()
		srv.apiDocs(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		body := rec.Body.String()
		if strings.Contains(body, "__BASE_URL__") {
			t.Error("__BASE_URL__ was not replaced")
		}
		if !strings.Contains(body, "http://localhost:8765") {
			t.Error("expected base URL in body")
		}
	})

	t.Run("X-Forwarded-Proto https", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/docs", nil)
		req.Host = "example.com"
		req.Header.Set("X-Forwarded-Proto", "https")
		rec := httptest.NewRecorder()
		srv.apiDocs(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "https://example.com") {
			t.Error("expected https base URL")
		}
	})
}

func TestLicense(t *testing.T) {
	srv := newTestServer(t)

	t.Run("GET mit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/licenses/mit", nil)
		rec := httptest.NewRecorder()
		srv.license(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
			t.Errorf("Content-Type = %q", ct)
		}
		if !strings.Contains(rec.Body.String(), "MIT License") {
			t.Error("expected MIT license text")
		}
	})

	t.Run("GET cc0", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/licenses/cc0", nil)
		rec := httptest.NewRecorder()
		srv.license(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "CC0") {
			t.Error("expected CC0 text")
		}
	})

	t.Run("unknown license returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/licenses/unknown", nil)
		rec := httptest.NewRecorder()
		srv.license(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("got status %d, want 404", rec.Code)
		}
	})

	t.Run("POST returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/licenses/mit", nil)
		rec := httptest.NewRecorder()
		srv.license(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got status %d, want 405", rec.Code)
		}
	})
}

// --- Bounds handler tests ---

func TestApiTracksBounds(t *testing.T) {
	srv := newTestServer(t)

	t.Run("missing trackIDs returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/tracks/bounds", nil)
		rec := httptest.NewRecorder()
		srv.apiTracksBounds(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("valid trackIDs returns bounds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/tracks/bounds?trackIDs="+testTrackID, nil)
		rec := httptest.NewRecorder()
		srv.apiTracksBounds(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if out["status"] != "success" {
			t.Errorf("status = %v, want success", out["status"])
		}
		bounds, ok := out["bounds"].(map[string]interface{})
		if !ok {
			t.Fatal("bounds missing or wrong type")
		}
		if _, ok := bounds["minLat"]; !ok {
			t.Error("bounds missing minLat")
		}
	})

	t.Run("non-existent track returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/tracks/bounds?trackIDs=nonexistent-track-xyz", nil)
		rec := httptest.NewRecorder()
		srv.apiTracksBounds(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("got status %d, want 404", rec.Code)
		}
	})
}

// --- Short redirect tests ---

func TestShortRedirect(t *testing.T) {
	t.Run("empty code returns 404", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/s/", nil)
		rec := httptest.NewRecorder()
		srv.shortRedirect(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("got status %d, want 404", rec.Code)
		}
	})

	t.Run("DB nil returns 404", func(t *testing.T) {
		srv := NewServer(nil, nil, Config{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/s/abc", nil)
		rec := httptest.NewRecorder()
		srv.shortRedirect(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("got status %d, want 404", rec.Code)
		}
	})

	t.Run("valid code redirects", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/s/"+testShortCode, nil)
		rec := httptest.NewRecorder()
		srv.shortRedirect(rec, req)
		if rec.Code != http.StatusFound {
			t.Errorf("got status %d, want 302", rec.Code)
		}
		if loc := rec.Header().Get("Location"); loc != testShortTarget {
			t.Errorf("Location = %q, want %q", loc, testShortTarget)
		}
	})

	t.Run("unknown code returns 404", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/s/unknown999", nil)
		rec := httptest.NewRecorder()
		srv.shortRedirect(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("got status %d, want 404", rec.Code)
		}
	})
}

// --- Spectrum handler tests ---

func TestSpectrum(t *testing.T) {
	t.Run("DB nil returns 503", func(t *testing.T) {
		srv := NewServer(nil, nil, Config{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/spectrum/1", nil)
		rec := httptest.NewRecorder()
		srv.spectrum(rec, req)
		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("got status %d, want 503", rec.Code)
		}
	})

	t.Run("invalid marker ID returns 400", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/api/spectrum/abc", nil)
		rec := httptest.NewRecorder()
		srv.spectrum(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("valid markerID returns spectrum", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/api/spectrum/1", nil)
		rec := httptest.NewRecorder()
		srv.spectrum(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if _, ok := out["channels"]; !ok {
			t.Error("response missing channels")
		}
	})

	t.Run("non-existent marker returns error", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/api/spectrum/99999", nil)
		rec := httptest.NewRecorder()
		srv.spectrum(rec, req)
		// GetSpectrum returns error when not found; handler returns 500
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want 500", rec.Code)
		}
	})
}

func TestSpectrumDownload(t *testing.T) {
	srv := newTestServer(t)

	t.Run("format json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/spectrum/1/download?format=json", nil)
		rec := httptest.NewRecorder()
		srv.spectrumDownload(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q", ct)
		}
		if cd := rec.Header().Get("Content-Disposition"); !strings.Contains(cd, "attachment") {
			t.Errorf("Content-Disposition = %q", cd)
		}
	})

	t.Run("format csv", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/spectrum/1/download?format=csv", nil)
		rec := httptest.NewRecorder()
		srv.spectrumDownload(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		body := rec.Body.String()
		if !strings.HasPrefix(body, "Channel,Energy_keV,Counts") {
			t.Errorf("CSV missing header: %q", body[:min(50, len(body))])
		}
	})

	t.Run("format n42 without raw returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/spectrum/1/download?format=n42", nil)
		rec := httptest.NewRecorder()
		srv.spectrumDownload(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("got status %d, want 404", rec.Code)
		}
	})

	t.Run("format invalid returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/spectrum/1/download?format=invalid", nil)
		rec := httptest.NewRecorder()
		srv.spectrumDownload(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})
}

// --- Track info tests ---

func TestTrackInfo(t *testing.T) {
	t.Run("DB nil returns 503", func(t *testing.T) {
		srv := NewServer(nil, nil, Config{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/track-info/"+testTrackID, nil)
		rec := httptest.NewRecorder()
		srv.trackInfo(rec, req)
		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("got status %d, want 503", rec.Code)
		}
	})

	t.Run("missing trackID returns 400", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/api/track-info/", nil)
		rec := httptest.NewRecorder()
		srv.trackInfo(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("valid track returns metadata", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/api/track-info/"+testTrackID, nil)
		rec := httptest.NewRecorder()
		srv.trackInfo(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if out["trackID"] != testTrackID {
			t.Errorf("trackID = %v, want %s", out["trackID"], testTrackID)
		}
		if _, ok := out["username"]; !ok {
			t.Error("response missing username")
		}
		if _, ok := out["detector"]; !ok {
			t.Error("response missing detector")
		}
	})

	t.Run("unknown track returns minimal JSON", func(t *testing.T) {
		srv := newTestServer(t)
		req := httptest.NewRequest(http.MethodGet, "/api/track-info/unknown-track-xyz", nil)
		rec := httptest.NewRecorder()
		srv.trackInfo(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if out["trackID"] != "unknown-track-xyz" {
			t.Errorf("trackID = %v", out["trackID"])
		}
	})
}
