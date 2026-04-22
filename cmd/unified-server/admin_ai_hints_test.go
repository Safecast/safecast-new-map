package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSlugifyModel(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"Claude 3.5 Sonnet", "claude-3-5-sonnet"},
		{"GPT-4o", "gpt-4o"},
		{"Kimi  K2", "kimi-k2"},
		{"  Qwen/2.5 ", "qwen-2-5"},
		{"日本語 Bot", "bot"},
		{"---weird---name---", "weird-name"},
	}
	for _, c := range cases {
		if got := slugifyModel(c.in); got != c.out {
			t.Errorf("slugifyModel(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}

func TestValidateModelSlug(t *testing.T) {
	ok := []string{"claude", "gpt-4o", "kimi-k2", "a1", "my-custom-bot-2"}
	bad := []string{"", "Claude", "gpt_4", "-lead", "trail-", "with space", strings.Repeat("a", 65)}
	for _, s := range ok {
		if err := validateModelSlug(s); err != nil {
			t.Errorf("validateModelSlug(%q) should be valid, got %v", s, err)
		}
	}
	for _, s := range bad {
		if err := validateModelSlug(s); err == nil {
			t.Errorf("validateModelSlug(%q) should be invalid", s)
		}
	}
}

func TestModelFromPath(t *testing.T) {
	model, rest, err := modelFromPath("/api/admin/ai-hints/claude", "/api/admin/ai-hints/", 0)
	if err != nil || model != "claude" || len(rest) != 0 {
		t.Errorf("simple model: got model=%q rest=%v err=%v", model, rest, err)
	}
	model, rest, err = modelFromPath("/api/admin/ai-hints/gpt-4o/history", "/api/admin/ai-hints/", 1)
	if err != nil || model != "gpt-4o" || len(rest) != 1 || rest[0] != "history" {
		t.Errorf("model/history: got model=%q rest=%v err=%v", model, rest, err)
	}
	if _, _, err := modelFromPath("/api/admin/ai-hints/", "/api/admin/ai-hints/", 0); err == nil {
		t.Errorf("empty model should fail")
	}
	if _, _, err := modelFromPath("/api/admin/ai-hints/BadSlug", "/api/admin/ai-hints/", 0); err == nil {
		t.Errorf("invalid slug should fail")
	}
}

// The following tests assume the global `db` is nil at package scope (no main()
// initialization during `go test`). They check that handlers fail closed with
// 503 when the DB is unavailable, matching the pattern in admin_mcp_test.go.

func TestAdminAIHintsHandlers_DatabaseUnavailable(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		path    string
		body    string
		handler http.HandlerFunc
	}{
		{"list", http.MethodGet, "/api/admin/ai-hints", "", adminAIHintsListHandler},
		{"get", http.MethodGet, "/api/admin/ai-hints/claude", "", adminAIHintGetHandler},
		{"create", http.MethodPost, "/api/admin/ai-hints", `{"display_name":"Test"}`, adminAIHintCreateHandler},
		{"update", http.MethodPut, "/api/admin/ai-hints/claude", `{"display_name":"X"}`, adminAIHintUpdateHandler},
		{"delete", http.MethodDelete, "/api/admin/ai-hints/claude", "", adminAIHintDeleteHandler},
		{"restore", http.MethodPost, "/api/admin/ai-hints/claude/restore", "", adminAIHintRestoreHandler},
		{"snapshot", http.MethodPost, "/api/admin/ai-hints/claude/snapshot", "", adminAIHintSnapshotHandler},
		{"history-list", http.MethodGet, "/api/admin/ai-hints/claude/history", "", adminAIHintHistoryListHandler},
		{"history-restore", http.MethodPost, "/api/admin/ai-hints/claude/history/1/restore", "", adminAIHintHistoryRestoreHandler},
		{"export", http.MethodGet, "/api/admin/ai-hints/claude/export", "", adminAIHintExportHandler},
		{"import", http.MethodPost, "/api/admin/ai-hints/import", `{"display_name":"Imported"}`, adminAIHintImportHandler},
		{"reload", http.MethodPost, "/api/admin/ai-hints/reload", "", adminAIHintReloadHandler},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			var body *strings.Reader
			if tc.body != "" {
				body = strings.NewReader(tc.body)
			}
			var req *http.Request
			if body != nil {
				req = httptest.NewRequest(tc.method, tc.path, body)
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			}
			tc.handler(rec, req)
			if rec.Code != http.StatusServiceUnavailable {
				t.Errorf("expected 503 with db=nil, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestAdminAIHintCreateHandler_ValidationBeforeDB(t *testing.T) {
	// With db=nil we'd normally get 503 first, so invalid JSON cases can't be
	// distinguished here without a DB stub. This test documents the guard order
	// by asserting that the 503 short-circuit happens even for malformed input.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/ai-hints", strings.NewReader("{not json"))
	req.Header.Set("Content-Type", "application/json")
	adminAIHintCreateHandler(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 (db guard runs first), got %d: %s", rec.Code, rec.Body.String())
	}
}
