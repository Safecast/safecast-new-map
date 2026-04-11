package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleFeedback_MethodNotAllowed(t *testing.T) {
	handler := handleFeedback()
	for _, method := range []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	} {
		t.Run(method, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, "/api/feedback", nil)
			handler(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleFeedback_OptionsPreflight(t *testing.T) {
	handler := handleFeedback()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/feedback", nil)
	handler(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleFeedback_BadRequest(t *testing.T) {
	handler := handleFeedback()
	cases := []struct {
		name string
		body string
	}{
		{"malformed JSON", `{not valid json`},
		{"missing chat_id field", `{"score": 1}`},
		{"chat_id explicitly zero", `{"chat_id": 0, "score": 1}`},
		{"empty body", ``},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			handler(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}
