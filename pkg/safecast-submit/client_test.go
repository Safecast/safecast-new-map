package safecastsubmit

import (
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResolveUserID(t *testing.T) {
	tests := []struct {
		name       string
		respStatus int
		respBody   string
		wantID     string
		wantErr    bool
	}{
		{"success", http.StatusOK, `{"id":42,"name":"tester"}`, "42", false},
		{"unauthorized", http.StatusUnauthorized, ``, "", true},
		{"empty id", http.StatusOK, `{"id":0}`, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/users/me.json" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				if r.URL.Query().Get("api_key") != "test-key" {
					t.Errorf("missing/incorrect api_key param: %q", r.URL.Query().Get("api_key"))
				}
				w.WriteHeader(tt.respStatus)
				_, _ = w.Write([]byte(tt.respBody))
			}))
			defer srv.Close()

			c := NewClientWithBaseURL(srv.URL)
			id, err := c.ResolveUserID(context.Background(), "test-key")
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if id != tt.wantID {
				t.Errorf("got id %q, want %q", id, tt.wantID)
			}
		})
	}
}

func TestCheckExists(t *testing.T) {
	tests := []struct {
		name       string
		respBody   string
		filename   string
		wantExists bool
	}{
		{"match found", `[{"id":1,"source":{"url":"https://s3.example.com/uploads/bgeigie_import/source/1/2026-01-01_1200.log"}}]`, "2026-01-01_1200.log", true},
		{"no match", `[{"id":1,"source":{"url":"https://s3.example.com/uploads/bgeigie_import/source/1/other.log"}}]`, "2026-01-01_1200.log", false},
		{"empty list", `[]`, "2026-01-01_1200.log", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bgeigie_imports.json" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				q := r.URL.Query()
				if q.Get("by_user_id") != "7" {
					t.Errorf("missing/incorrect by_user_id param: %q", q.Get("by_user_id"))
				}
				if q.Get("q") != tt.filename {
					t.Errorf("missing/incorrect q param: %q", q.Get("q"))
				}
				if q.Get("api_key") != "test-key" {
					t.Errorf("missing/incorrect api_key param: %q", q.Get("api_key"))
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.respBody))
			}))
			defer srv.Close()

			c := NewClientWithBaseURL(srv.URL)
			exists, err := c.CheckExists(context.Background(), "test-key", "7", tt.filename)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if exists != tt.wantExists {
				t.Errorf("got exists=%v, want %v", exists, tt.wantExists)
			}
		})
	}
}

func TestCheckExists_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL)
	_, err := c.CheckExists(context.Background(), "test-key", "7", "f.log")
	if err == nil {
		t.Fatal("expected error on 500 response, got nil")
	}
}

func TestSubmit(t *testing.T) {
	const wantContent = "$BNRDD,0000000000,A,2026-01-01T12:00:00Z,100,60,..."

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/bgeigie_imports.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("api_key") != "test-key" {
			t.Errorf("missing/incorrect api_key param: %q", r.URL.Query().Get("api_key"))
		}

		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || !strings.HasPrefix(mediaType, "multipart/") {
			t.Fatalf("expected multipart content-type, got %q (err=%v)", r.Header.Get("Content-Type"), err)
		}

		mr := multipart.NewReader(r.Body, params["boundary"])
		var gotDescription, gotFilename, gotContent string
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("read multipart part: %v", err)
			}
			data, _ := io.ReadAll(part)
			switch part.FormName() {
			case "bgeigie_import[description]":
				gotDescription = string(data)
			case "bgeigie_import[source]":
				gotFilename = part.FileName()
				gotContent = string(data)
			}
		}

		if gotDescription == "" {
			t.Error("expected non-empty description field")
		}
		if gotFilename != "2026-01-01_1200.log" {
			t.Errorf("got filename %q, want %q", gotFilename, "2026-01-01_1200.log")
		}
		if gotContent != wantContent {
			t.Errorf("got content %q, want %q", gotContent, wantContent)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":123}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL)
	importID, err := c.Submit(context.Background(), "test-key", "2026-01-01_1200.log", []byte(wantContent))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if importID != "123" {
		t.Errorf("got importID %q, want %q", importID, "123")
	}
}

func TestSubmit_RejectedStatus(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"bad request", http.StatusBadRequest},
		{"unauthorized", http.StatusUnauthorized},
		{"forbidden", http.StatusForbidden},
		{"server error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			}))
			defer srv.Close()

			c := NewClientWithBaseURL(srv.URL)
			_, err := c.Submit(context.Background(), "test-key", "f.log", []byte("data"))
			if err == nil {
				t.Fatalf("expected error for status %d, got nil", tt.status)
			}
		})
	}
}
