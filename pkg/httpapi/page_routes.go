package httpapi

import (
	"net/http"

	"safecast-new-map/pkg/auth"
)

// PageRoutesConfig holds non-API page routes served by the main listener.
// These routes are mostly HTML/UI surfaces and are intentionally kept outside
// the /api registrars.
type PageRoutesConfig struct {
	AuthManager   *auth.Manager
	AdminPassword string

	HomeHandler        http.HandlerFunc
	StoriesPageHandler http.HandlerFunc
	StoriesDataHandler http.HandlerFunc
	MapHandler         http.HandlerFunc

	ProfileHandler       http.HandlerFunc
	ResetPasswordHandler http.HandlerFunc

	AdminUsersPageHandler        http.HandlerFunc
	AdminUploadsPageHandler      http.HandlerFunc
	AdminMCPPageHandler          http.HandlerFunc
	AdminRealtimePageHandler     http.HandlerFunc
	AdminTranslationsPageHandler http.HandlerFunc
	AdminTourPageHandler         http.HandlerFunc
}

// RegisterPageRoutes attaches UI/admin page routes to mux.
// Auth-sensitive routes are only exposed when AuthManager is configured.
func RegisterPageRoutes(mux *http.ServeMux, cfg PageRoutesConfig) {
	if mux == nil {
		return
	}

	registerOptional := func(path string, handler http.HandlerFunc) {
		if handler == nil {
			return
		}
		mux.HandleFunc(path, handler)
	}

	registerOptional("/home", cfg.HomeHandler)
	registerOptional("/stories.html", cfg.StoriesPageHandler)
	registerOptional("/data/stories.json", cfg.StoriesDataHandler)
	registerOptional("/", cfg.MapHandler)

	if cfg.AuthManager == nil {
		return
	}

	if cfg.ProfileHandler != nil {
		mux.HandleFunc("/profile", cfg.AuthManager.OptionalAuth(cfg.ProfileHandler))
	}
	registerOptional("/reset-password", cfg.ResetPasswordHandler)

	checkAdminAccess := func(w http.ResponseWriter, r *http.Request) bool {
		if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
			return true
		}
		if cfg.AdminPassword != "" && r.URL.Query().Get("password") == cfg.AdminPassword {
			return true
		}
		http.Error(w, "Unauthorized - Please login as admin or provide password", http.StatusUnauthorized)
		return false
	}

	registerOptionalAdmin := func(path string, handler http.HandlerFunc) {
		if handler == nil {
			return
		}
		mux.HandleFunc(path, cfg.AuthManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			if !checkAdminAccess(w, r) {
				return
			}
			handler(w, r)
		}))
	}

	registerOptionalAdmin("/admin/users", cfg.AdminUsersPageHandler)
	registerOptionalAdmin("/admin/uploads", cfg.AdminUploadsPageHandler)
	registerOptionalAdmin("/admin/mcp", cfg.AdminMCPPageHandler)
	registerOptionalAdmin("/admin/realtime", cfg.AdminRealtimePageHandler)
	registerOptionalAdmin("/admin/translations", cfg.AdminTranslationsPageHandler)
	registerOptionalAdmin("/admin/tour", cfg.AdminTourPageHandler)
}
