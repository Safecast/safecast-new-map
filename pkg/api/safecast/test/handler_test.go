// handler_test.go tests the API root endpoints (GET /api/v1 and GET /api/v2),
// Accept-header and method behavior, CORS headers, and parity between v1 and v2
// root responses. These endpoints return a JSON descriptor with API name, URI,
// and subresource links; they do not require a database.

package safecast_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

// TestHandler_Root checks that GET /api/v1 returns 200, application/json, and a
// body containing the API name, base URI, and non-empty subresource_uris list
// (per the migration plan root contract).
func TestHandler_Root(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)

	assertStatus(t, rec, http.StatusOK)
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
	}
	if rec.Body.Len() == 0 {
		t.Fatal("empty response body")
	}
	var root struct {
		Name             string   `json:"name"`
		URI              string   `json:"uri"`
		SubresourceURIs   []string `json:"subresource_uris"`
	}
	if err := decodeJSON(rec, &root); err != nil {
		t.Fatalf("decode root: %v", err)
	}
	if root.Name != "Safecast API" {
		t.Errorf("name = %q, want Safecast API", root.Name)
	}
	if root.URI == "" {
		t.Error("uri is empty")
	}
	if len(root.SubresourceURIs) == 0 {
		t.Error("subresource_uris is empty")
	}
}

// TestHandler_V2Root checks that GET /api/v2 returns 200 and application/json.
// v2 root uses the same descriptor as v1 but does not require Accept: application/json.
func TestHandler_V2Root(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v2", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)

	assertStatus(t, rec, http.StatusOK)
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
	}
}

// TestHandler_Root_NoAccept ensures that GET /api/v1 without Accept: application/json
// returns 406 Not Acceptable, as the API contract requires JSON for this endpoint.
func TestHandler_Root_NoAccept(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
	rec := serve(req, h)

	assertStatus(t, rec, http.StatusNotAcceptable)
}

// TestHandler_Root_MethodNotAllowed ensures that POST /api/v1 returns 405 Method Not
// Allowed and that the Allow header lists the permitted method (GET).
func TestHandler_Root_MethodNotAllowed(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)

	assertStatus(t, rec, http.StatusMethodNotAllowed)
	if allow := rec.Header().Get("Allow"); allow != "GET" {
		t.Errorf("Allow = %q, want GET", allow)
	}
}

// TestHandler_RootParity verifies that GET /api/v1 and GET /api/v2 return equivalent
// data: same API name and same set of subresource_uris (order may differ, so we
// sort before comparing). This parity check ensures both versions use the same core.
func TestHandler_RootParity(t *testing.T) {
	h := newTestHandler(nil, nil)
	rootStruct := func(path string) (name string, uris []string) {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Accept", "application/json")
		rec := serve(req, h)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status %d", path, rec.Code)
		}
		var r struct {
			Name           string   `json:"name"`
			SubresourceURIs []string `json:"subresource_uris"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&r); err != nil {
			t.Fatalf("%s decode: %v", path, err)
		}
		return r.Name, r.SubresourceURIs
	}
	name1, uris1 := rootStruct("/api/v1")
	name2, uris2 := rootStruct("/api/v2")
	if name1 != name2 {
		t.Errorf("name mismatch: v1 %q, v2 %q", name1, name2)
	}
	sort.Strings(uris1)
	sort.Strings(uris2)
	if len(uris1) != len(uris2) {
		t.Errorf("subresource_uris length: v1 %d, v2 %d", len(uris1), len(uris2))
	}
	for i := range uris1 {
		if uris1[i] != uris2[i] {
			t.Errorf("subresource_uris[%d]: v1 %q, v2 %q", i, uris1[i], uris2[i])
		}
	}
}

// TestHandler_CORS checks that when the request includes an Origin header, the
// response includes Access-Control-Allow-Origin set to that origin (per the
// migration plan CORS contract). This allows browser-based clients to call the API.
func TestHandler_CORS(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v2", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := serve(req, h)

	assertStatus(t, rec, http.StatusOK)
	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "https://example.com" {
		t.Errorf("Access-Control-Allow-Origin = %q, want https://example.com", origin)
	}
}
