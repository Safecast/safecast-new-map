package safecast

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Root(t *testing.T) {
	h := NewHandler(nil, "sqlite", nil, "https://api.example.com", nil)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET /api/v1: status %d, want 200", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
	}
	if rec.Body.Len() == 0 {
		t.Error("empty response body")
	}
}

func TestHandler_V2Root(t *testing.T) {
	h := NewHandler(nil, "sqlite", nil, "https://api.example.com", nil)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v2", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET /api/v2: status %d, want 200", rec.Code)
	}
}
