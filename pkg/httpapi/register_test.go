// register_test.go tests that Register installs the expected routes and that auth/admin guards behave correctly.
package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"safecast-new-map/pkg/auth"
)

// testRegistrar is a routeRegistrar that registers stub handlers for a fixed list of paths.
type testRegistrar struct {
	routes []string
}

func (t testRegistrar) Register(mux *http.ServeMux) {
	for _, route := range t.routes {
		mux.HandleFunc(route, func(http.ResponseWriter, *http.Request) {})
	}
}

func TestRegisterRouteInventory(t *testing.T) {
	mux := http.NewServeMux()

	authManager := &auth.Manager{}

	cfg := RegisterConfig{
		WebServer: testRegistrar{routes: []string{
			"/api/markers/spectra",
			"/api/update-coordinates",
			"/api/tracks/bounds",
		}},
		APIHandler: testRegistrar{routes: []string{
			"/api",
			"/api/latest",
			"/api/track/",
		}},
		AuthManager:   authManager,
		AdminPassword: "secret",

		AdminUploadsHandler: func(http.ResponseWriter, *http.Request) {},
		AdminTracksHandler:  func(http.ResponseWriter, *http.Request) {},
		AdminMCPDataHandler:            func(http.ResponseWriter, *http.Request) {},
		AdminMCPExportHandler:          func(http.ResponseWriter, *http.Request) {},
		AdminMCPDeleteHandler:          func(http.ResponseWriter, *http.Request) {},
		AdminRealtimeDataHandler:       func(http.ResponseWriter, *http.Request) {},
		AdminRealtimeExportHandler:     func(http.ResponseWriter, *http.Request) {},
		AdminRealtimeDeleteHandler:     func(http.ResponseWriter, *http.Request) {},
		AdminTranslationsReloadHandler: func(http.ResponseWriter, *http.Request) {},
		AdminTranslationByIDHandler:    func(http.ResponseWriter, *http.Request) {},
		AdminTranslationsHandler:       func(http.ResponseWriter, *http.Request) {},
	}
	Register(mux, cfg)

	requiredRoutes := []string{
		"/api",
		"/api/latest",
		"/api/track/",
		"/api/markers/spectra",
		"/api/update-coordinates",
		"/api/tracks/bounds",
		"/api/auth/login",
		"/api/user/profile",
		"/api/user/uploads",
		"/api/admin/users",
		"/api/admin/uploads",
		"/api/admin/tracks",
		"/api/admin/mcp/data",
		"/api/admin/mcp/export",
		"/api/admin/mcp/delete",
		"/api/admin/realtime/data",
		"/api/admin/realtime/export",
		"/api/admin/realtime/delete",
		"/api/admin/translations/reload",
		"/api/admin/translations/",
		"/api/admin/translations",
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

func TestRegisterTrackRoute(t *testing.T) {
	mux := http.NewServeMux()
	cfg := RegisterConfig{
		APIHandler: NewHandler(nil, "sqlite", nil, nil, nil, ""),
	}
	Register(mux, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/track/abc.json", nil)
	_, pattern := mux.Handler(req)
	if pattern == "" {
		t.Fatal("route /api/track/ was not registered")
	}
}

func TestRegisterAuthAdminGuards(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, RegisterConfig{
		AuthManager:   &auth.Manager{},
		AdminPassword: "secret",
	})

	t.Run("admin users requires auth or password", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("got status %d, want 401", rec.Code)
		}
	})

	t.Run("user uploads requires session auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/uploads", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("got status %d, want 401", rec.Code)
		}
	})
}
