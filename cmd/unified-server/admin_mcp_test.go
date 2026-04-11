package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAdminMCPUpdateHandler_MethodNotAllowed(t *testing.T) {
	for _, method := range []string{
		http.MethodGet,
		http.MethodDelete,
		http.MethodPatch,
	} {
		t.Run(method, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, "/api/admin/mcp/update?table=chat_questions&id=1", nil)
			adminMCPUpdateHandler(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestAdminMCPUpdateHandler_AnalyticsUnavailable(t *testing.T) {
	for _, method := range []string{http.MethodPut, http.MethodPost} {
		t.Run(method, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method,
				"/api/admin/mcp/update?table=chat_questions&id=1",
				strings.NewReader(`{"question":"updated question"}`))
			req.Header.Set("Content-Type", "application/json")
			adminMCPUpdateHandler(rec, req)
			if rec.Code != http.StatusServiceUnavailable {
				t.Errorf("expected 503, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}
