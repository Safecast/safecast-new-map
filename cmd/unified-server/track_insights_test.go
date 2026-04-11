package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrackInsightsHandler_MissingTrackID(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/track//insights", nil)
	trackInsightsHandler(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestTrackInsightsHandler_DBUnavailable(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/track/{id}/insights", trackInsightsHandler)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/track/abc123/insights", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}
