package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"safecast-new-map/pkg/auth"
)

func TestRegisterLegacyRoutesInventory(t *testing.T) {
	mux := http.NewServeMux()
	RegisterLegacyRoutes(mux, LegacyRoutesConfig{
		UploadHandler:          func(http.ResponseWriter, *http.Request) {},
		UploadProgressHandler:  func(http.ResponseWriter, *http.Request) {},
		GetMarkersHandler:      func(http.ResponseWriter, *http.Request) {},
		StreamMarkersHandler:   func(http.ResponseWriter, *http.Request) {},
		RealtimeHistoryHandler: func(http.ResponseWriter, *http.Request) {},
		TrackByIDHandler:       func(http.ResponseWriter, *http.Request) {},
		TracksByPrefixHandler:  func(http.ResponseWriter, *http.Request) {},
	})

	requiredRoutes := []string{
		"/upload",
		"/upload/progress",
		"/get_markers",
		"/stream_markers",
		"/realtime_history",
		"/trackid/abc",
		"/tracks/abc",
	}

	for _, route := range requiredRoutes {
		t.Run(route, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, route, nil)
			_, pattern := mux.Handler(req)
			if pattern == "" {
				t.Fatalf("route %s was not registered", route)
			}
		})
	}
}

func TestRegisterLegacyRoutesUploadAuthGuard(t *testing.T) {
	t.Run("guarded when auth is required and manager exists", func(t *testing.T) {
		mux := http.NewServeMux()
		RegisterLegacyRoutes(mux, LegacyRoutesConfig{
			RequireAuth:   true,
			AuthManager:   &auth.Manager{},
			UploadHandler: func(http.ResponseWriter, *http.Request) {},
		})

		req := httptest.NewRequest(http.MethodPost, "/upload", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("left unguarded when manager is missing", func(t *testing.T) {
		mux := http.NewServeMux()
		RegisterLegacyRoutes(mux, LegacyRoutesConfig{
			RequireAuth:   true,
			AuthManager:   nil,
			UploadHandler: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusAccepted) },
		})

		req := httptest.NewRequest(http.MethodPost, "/upload", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusAccepted {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusAccepted)
		}
	})
}
