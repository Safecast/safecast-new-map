// Package safecast_test contains contract and parity tests for the Safecast API
// (pkg/api/safecast). The API replaces api.safecast.org and serves measurements,
// bgeigie imports, users, and devices over HTTP. These tests call the API handlers
// without starting a real server: we build an HTTP request, pass it to the handler,
// and inspect the recorded response (status code, headers, JSON body). Tests use
// the standard Go testing package and net/http/httptest; no external test framework
// is required. Run all tests with: go test ./pkg/api/safecast/...
package safecast_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"safecast-new-map/pkg/api/safecast"
	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
)

// newTestHandler builds an API handler for tests. Pass nil for db and authMgr when
// the test does not need a database or authentication (e.g. root endpoint, stubs).
// The handler is configured with a test base URL and no-op logging.
func newTestHandler(db *database.Database, authMgr *auth.Manager) *safecast.Handler {
	return safecast.NewHandler(db, "sqlite", authMgr, "https://api.test", nil)
}

// serve sends the given HTTP request to the handler and returns the response.
// Internally it creates a router (ServeMux), registers the Safecast routes on it,
// and uses httptest.ResponseRecorder to capture the response without opening a
// real network port. This is the standard Go approach for testing HTTP handlers.
func serve(req *http.Request, h *safecast.Handler) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	h.Register(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

// decodeJSON parses the response body as JSON into the value pointed at by v.
// Use this after checking the status code (rec.Code) to inspect success or error
// payloads. v is typically a struct or map with `json:"field"` tags.
func decodeJSON(rec *httptest.ResponseRecorder, v interface{}) error {
	return json.NewDecoder(rec.Body).Decode(v)
}

// assertStatus fails the test if the response status code is not the expected value.
// It also prints the response body on failure to help with debugging. The t.Helper()
// call ensures that failure reports point to the calling test line, not this function.
func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Errorf("status = %d, want %d; body: %s", rec.Code, want, rec.Body.String())
	}
}
