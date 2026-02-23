// bgeigie_test.go tests the bgeigie_imports API: list (GET), get-by-id (GET :id),
// and create (POST). With nil DB the list returns 500; get-by-id with invalid or
// missing id returns 404; POST is currently stubbed with 501 Not Implemented.
// We also check Accept header and method-not-allowed behavior.

package safecast_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestBgeigie_GET_List_NilDB checks that GET bgeigie_imports on adapter and v1
// paths returns 500 when no database is configured.
func TestBgeigie_GET_List_NilDB(t *testing.T) {
	h := newTestHandler(nil, nil)
	tests := []struct {
		path string
	}{
		{"/bgeigie_imports"},
		{"/bgeigie_imports.json"},
		{"/api/v1/bgeigie_imports"},
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

// TestBgeigie_GET_ByID_InvalidID ensures that a non-numeric id returns 404.
func TestBgeigie_GET_ByID_InvalidID(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/bgeigie_imports/notanid", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusNotFound)
}

// TestBgeigie_GET_ByID_NoID ensures that GET /bgeigie_imports/ with no id segment
// after the trailing slash returns 404 (no resource id to look up).
func TestBgeigie_GET_ByID_NoID(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/bgeigie_imports/", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusNotFound)
}

// TestBgeigie_POST_NotImplemented checks that POST /bgeigie_imports returns 501
// Not Implemented, as the upload flow is not yet implemented.
func TestBgeigie_POST_NotImplemented(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/bgeigie_imports", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusNotImplemented)
}

// TestBgeigie_List_RequiresJSON ensures that GET /bgeigie_imports without
// Accept: application/json returns 406 Not Acceptable.
func TestBgeigie_List_RequiresJSON(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/bgeigie_imports", nil)
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusNotAcceptable)
}

// TestBgeigie_MethodNotAllowed ensures that PUT /bgeigie_imports returns 405 and
// that the Allow header lists GET and POST.
func TestBgeigie_MethodNotAllowed(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodPut, "/bgeigie_imports", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusMethodNotAllowed)
	if allow := rec.Header().Get("Allow"); allow != "GET, POST" {
		t.Errorf("Allow = %q, want GET, POST", allow)
	}
}

// TestBgeigie_GetByID_MethodNotAllowed ensures that POST /bgeigie_imports/1 returns
// 405 and that Allow lists GET only for the by-id endpoint.
func TestBgeigie_GetByID_MethodNotAllowed(t *testing.T) {
	h := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/bgeigie_imports/1", nil)
	req.Header.Set("Accept", "application/json")
	rec := serve(req, h)
	assertStatus(t, rec, http.StatusMethodNotAllowed)
	if allow := rec.Header().Get("Allow"); allow != "GET" {
		t.Errorf("Allow = %q, want GET", allow)
	}
}
