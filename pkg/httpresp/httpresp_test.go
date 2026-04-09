package httpresp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteErrorEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, http.StatusBadRequest, "bad_input", "Invalid request")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("got content-type %q, want application/json", ct)
	}

	var body ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error body: %v", err)
	}
	if body.Status != "error" {
		t.Fatalf("got status %q, want error", body.Status)
	}
	if body.Code != "bad_input" {
		t.Fatalf("got code %q, want bad_input", body.Code)
	}
	if body.Error != "Invalid request" {
		t.Fatalf("got error %q, want Invalid request", body.Error)
	}
}

func TestRequireMethodJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()

	ok := RequireMethodJSON(rec, req, http.MethodPost)
	if ok {
		t.Fatal("RequireMethodJSON returned true, want false")
	}
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got := rec.Header().Get("Allow"); got != http.MethodPost {
		t.Fatalf("got Allow %q, want %q", got, http.MethodPost)
	}
}
