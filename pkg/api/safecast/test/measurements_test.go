// measurements_test.go tests the measurements API: list (GET), get-by-id (GET :id),
// count (GET count), and create (POST). With a nil database the list and count
// endpoints return 500; we assert that status and that validation/error shapes
// (e.g. 422 with errors.value or errors.base) match the migration plan contract.

package safecast_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestMeasurements_GET_List_NilDB checks that GET measurements on adapter and v1
// paths returns 500 when no database is configured. The handler tries to query
// the DB and returns Internal Server Error when the DB is nil.
func TestMeasurements_GET_List_NilDB(t *testing.T) {
	h := newTestHandler(nil, nil)
	tests := []struct {
		path string
	}{
		{"/measurements"},
		{"/measurements.json"},
		{"/api/v1/measurements"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Accept", "application/json")
			rec := serve(req, h)
			assertStatus(t, rec, http.StatusInternalServerError)
		})
	}
}

// TestMeasurements_GET_ByID_NotFoundPath ensures that a path with no id segment
// (e.g. GET /measurements/) returns 404 Not Found.
func TestMeasurements_GET_ByID_NotFoundPath(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/measurements/", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusNotFound)
}

// TestMeasurements_GET_ByID_InvalidID ensures that a non-numeric id in the path
// (e.g. /measurements/abc) returns 404, as the id cannot be parsed.
func TestMeasurements_GET_ByID_InvalidID(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/measurements/abc", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusNotFound)
}

// TestMeasurements_GET_Count_NilDB checks that the count endpoint (multiple paths:
// /count, /measurements/count, etc.) returns 500 when no database is configured.
func TestMeasurements_GET_Count_NilDB(t *testing.T) {
	h := newTestHandler(nil, nil)
	tests := []struct {
		path string
	}{
		{"/count"},
		{"/measurements/count"},
		{"/measurements/count.json"},
		{"/api/v1/measurements/count"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Accept", "application/json")
			rec := serve(req, h)
			assertStatus(t, rec, http.StatusInternalServerError)
		})
	}
}

// TestMeasurements_POST_InvalidJSON ensures that POST /measurements with a body
// that is not valid JSON returns 422 Unprocessable Entity and a JSON error body
// with an "errors" map containing "base" (e.g. "invalid JSON") per the Rails contract.
func TestMeasurements_POST_InvalidJSON(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/measurements.json", bytes.NewBufferString("not json"))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rec := serve(req, h)

	assertStatus(t, rec, http.StatusUnprocessableEntity)
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
	}
	var errBody struct {
		Errors map[string][]string `json:"errors"`
	}
	if err := decodeJSON(rec, &errBody); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if errBody.Errors == nil {
		t.Fatal("expected errors map")
	}
	if base, ok := errBody.Errors["base"]; !ok || len(base) == 0 {
		t.Errorf("expected errors.base, got %v", errBody.Errors)
	}
}

// TestMeasurements_POST_EmptyMeasurement ensures that POST with an empty measurement
// object (no value, no unit) returns 422 and errors.value (e.g. "can't be blank").
func TestMeasurements_POST_EmptyMeasurement(t *testing.T) {
	h := newTestHandler(nil, nil)
	body := `{"measurement": {}}`
	req := httptest.NewRequest(http.MethodPost, "/measurements.json", bytes.NewBufferString(body))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rec := serve(req, h)

	assertStatus(t, rec, http.StatusUnprocessableEntity)
	var errBody struct {
		Errors map[string][]string `json:"errors"`
	}
	if err := decodeJSON(rec, &errBody); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if errBody.Errors == nil {
		t.Fatal("expected errors map")
	}
	if val, ok := errBody.Errors["value"]; !ok || len(val) == 0 {
		t.Errorf("expected errors.value, got %v", errBody.Errors)
	}
}

// TestMeasurements_POST_ValidBody_NilDB checks that a valid POST body (value, unit,
// lat/lon) is accepted by validation but returns 500 when the database is nil,
// since creating the measurement requires a DB.
func TestMeasurements_POST_ValidBody_NilDB(t *testing.T) {
	h := newTestHandler(nil, nil)
	body := `{"measurement": {"value": 100, "unit": "cpm", "latitude": 35.0, "longitude": 139.0}}`
	req := httptest.NewRequest(http.MethodPost, "/measurements.json", bytes.NewBufferString(body))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusInternalServerError)
}

// TestMeasurements_List_RequiresJSON ensures that GET /measurements without
// Accept: application/json and without .json in the path returns 406 Not Acceptable.
func TestMeasurements_List_RequiresJSON(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/measurements", nil)
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusNotAcceptable)
}

// TestMeasurements_Count_MethodNotAllowed ensures that POST /count returns 405
// and that the Allow header lists GET as the permitted method.
func TestMeasurements_Count_MethodNotAllowed(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/count", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusMethodNotAllowed)
	if allow := rec.Header().Get("Allow"); allow != "GET" {
		t.Errorf("Allow = %q, want GET", allow)
	}
}
