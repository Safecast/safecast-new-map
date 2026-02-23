// adapter_contract_test.go tests the stub and adapter endpoints that do not
// require a database: users, devices, radiation_index, ingest, device_stories.
// These return 200 with JSON (empty arrays or stub data). We also check that
// by-id endpoints (users/1, devices/1) return 501 Not Implemented and that
// method-not-allowed and Accept-header behavior are correct.

package safecast_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAdapter_GET_Users_ReturnsEmptyArray checks that GET /users (and .json and
// /api/v1/users) returns 200, application/json, and an empty JSON array [],
// since the users list is currently a stub with no DB.
func TestAdapter_GET_Users_ReturnsEmptyArray(t *testing.T) {
	h := newTestHandler(nil, nil)
	tests := []struct {
		path string
	}{
		{"/users"},
		{"/users.json"},
		{"/api/v1/users"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Accept", "application/json")
			rec := serve(req, h)
			assertStatus(t, rec, http.StatusOK)
			if rec.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
			}
			var out []interface{}
			if err := decodeJSON(rec, &out); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if out == nil || len(out) != 0 {
				t.Errorf("expected empty array, got %v", out)
			}
		})
	}
}

// TestAdapter_GET_Users_RequiresJSON ensures that GET /users without
// Accept: application/json returns 406 Not Acceptable.
func TestAdapter_GET_Users_RequiresJSON(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusNotAcceptable)
}

// TestAdapter_GET_Devices_ReturnsEmptyArray checks that GET /devices (and .json
// and /api/v1/devices) returns 200 and an empty JSON array [].
func TestAdapter_GET_Devices_ReturnsEmptyArray(t *testing.T) {
	h := newTestHandler(nil, nil)
	tests := []struct {
		path string
	}{
		{"/devices"},
		{"/devices.json"},
		{"/api/v1/devices"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Accept", "application/json")
			rec := serve(req, h)
			assertStatus(t, rec, http.StatusOK)
			var out []interface{}
			if err := decodeJSON(rec, &out); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if out == nil || len(out) != 0 {
				t.Errorf("expected empty array, got %v", out)
			}
		})
	}
}

// TestAdapter_GET_RadiationIndex_ReturnsJSON checks that GET /radiation_index
// with Accept: application/json returns 200 and a JSON array (stub data).
func TestAdapter_GET_RadiationIndex_ReturnsJSON(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/radiation_index", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusOK)
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
	}
	var out []map[string]float64
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out == nil {
		t.Error("expected non-nil array")
	}
}

// TestAdapter_GET_Ingest_ReturnsJSON checks that GET /ingest returns 200 and
// application/json (stub empty array; the real implementation would use Elasticsearch).
func TestAdapter_GET_Ingest_ReturnsJSON(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/ingest", nil)
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusOK)
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
	}
	var out []interface{}
	if err := decodeJSON(rec, &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out == nil {
		t.Error("expected non-nil array")
	}
}

// TestAdapter_GET_DeviceStories_ReturnsJSON checks that GET /device_stories
// returns 200 and a JSON array (stub).
func TestAdapter_GET_DeviceStories_ReturnsJSON(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/device_stories", nil)
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusOK)
	var out []interface{}
	if err := decodeJSON(rec, &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out == nil {
		t.Error("expected non-nil array")
	}
}

// TestAdapter_GET_UserByID_NotImplemented ensures that GET /users/1 returns
// 501 Not Implemented, as the single-user fetch is not yet implemented.
func TestAdapter_GET_UserByID_NotImplemented(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusNotImplemented)
}

// TestAdapter_GET_DeviceByID_NotImplemented ensures that GET /devices/1 returns
// 501 Not Implemented.
func TestAdapter_GET_DeviceByID_NotImplemented(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/devices/1", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusNotImplemented)
}

// TestAdapter_RadiationIndex_MethodNotAllowed ensures that POST /radiation_index
// returns 405 and Allow: GET.
func TestAdapter_RadiationIndex_MethodNotAllowed(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/radiation_index", nil)
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusMethodNotAllowed)
	if allow := rec.Header().Get("Allow"); allow != "GET" {
		t.Errorf("Allow = %q, want GET", allow)
	}
}

// TestAdapter_Users_MethodNotAllowed ensures that PUT /users returns 405 and
// that Allow lists GET and POST.
func TestAdapter_Users_MethodNotAllowed(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodPut, "/users", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusMethodNotAllowed)
	if allow := rec.Header().Get("Allow"); allow != "GET, POST" {
		t.Errorf("Allow = %q, want GET, POST", allow)
	}
}
