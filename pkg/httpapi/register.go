// register.go wires all HTTP routes: web/API handlers plus auth and admin endpoints.
package httpapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
)

// routeRegistrar is implemented by both Server (web) and Handler (core API).
type routeRegistrar interface {
	Register(*http.ServeMux)
}

// RegisterConfig holds everything needed to register all API and web routes.
// Fill WebServer (from NewWebServer) and APIHandler (from NewHandler); optionally
// set AuthManager, AdminPassword, and the admin handler funcs for protected routes.
type RegisterConfig struct {
	WebServer  routeRegistrar
	APIHandler routeRegistrar

	AuthManager   *auth.Manager
	DB            *database.Database
	AdminPassword string

	// Admin handlers; each is guarded by optional auth or admin password.
	AdminUploadsHandler              http.HandlerFunc
	AdminTracksHandler               http.HandlerFunc
	AdminBackfillHandler             http.HandlerFunc
	AdminBackfillCountriesHandler    http.HandlerFunc
	AdminDeleteTrackHandler          http.HandlerFunc
	AdminDeleteMultipleTracksHandler http.HandlerFunc
	AdminImportFromSafecastHandler   http.HandlerFunc
	AdminImportByIDHandler           http.HandlerFunc
	AdminUpdateTrackHandler          http.HandlerFunc
	AdminUpdateUploadHandler         http.HandlerFunc
	AdminImportSafecastMetaHandler   http.HandlerFunc
	AdminCacheHandler                http.HandlerFunc
	AdminMCPDataHandler              http.HandlerFunc
	AdminMCPExportHandler            http.HandlerFunc
	AdminMCPDeleteHandler            http.HandlerFunc
	AdminRealtimeDataHandler         http.HandlerFunc
	AdminRealtimeExportHandler       http.HandlerFunc
	AdminRealtimeDeleteHandler       http.HandlerFunc
	AdminTranslationsReloadHandler   http.HandlerFunc
	AdminTranslationByIDHandler      http.HandlerFunc
	AdminTranslationsHandler         http.HandlerFunc

	Logf func(string, ...any)
}

// Register attaches all routes to mux: first web and core API routes from the
// registrars, then auth and admin routes from RegisterConfig.
func Register(mux *http.ServeMux, cfg RegisterConfig) {
	if cfg.WebServer != nil {
		cfg.WebServer.Register(mux)
	}
	if cfg.APIHandler != nil {
		cfg.APIHandler.Register(mux)
	}
	registerAuthAndAdminRoutes(mux, cfg)
}

// registerAuthAndAdminRoutes adds /api/auth/*, /api/user/*, and /api/admin/* routes.
func registerAuthAndAdminRoutes(mux *http.ServeMux, cfg RegisterConfig) {
	if mux == nil {
		return
	}
	if cfg.AuthManager != nil {
		mux.HandleFunc("/api/auth/register", cfg.AuthManager.RegisterHandler)
		mux.HandleFunc("/api/auth/login", cfg.AuthManager.LoginHandler)
		mux.HandleFunc("/api/auth/logout", cfg.AuthManager.LogoutHandler)
		mux.HandleFunc("/api/auth/forgot-password", cfg.AuthManager.ForgotPasswordHandler)
		mux.HandleFunc("/api/auth/reset-password", cfg.AuthManager.ResetPasswordHandler)
		mux.HandleFunc("/api/auth/verify-email", cfg.AuthManager.VerifyEmailHandler)
		mux.HandleFunc("/api/user/profile", cfg.AuthManager.RequireAuth(cfg.AuthManager.ProfileHandler))
		mux.HandleFunc("/api/user/change-password", cfg.AuthManager.RequireAuth(cfg.AuthManager.ChangePasswordHandler))
		mux.HandleFunc("/api/user/uploads", cfg.AuthManager.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			handleUserUploads(w, r, cfg)
		}))

		mux.HandleFunc("/api/admin/users", cfg.AuthManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			if !checkAdminAccess(w, r, cfg.AdminPassword) {
				return
			}
			cfg.AuthManager.AdminListUsersHandler(w, r)
		}))
		mux.HandleFunc("/api/admin/users/create", cfg.AuthManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			if !checkAdminAccess(w, r, cfg.AdminPassword) {
				return
			}
			cfg.AuthManager.AdminCreateUserHandler(w, r)
		}))
		mux.HandleFunc("/api/admin/users/", cfg.AuthManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			if !checkAdminAccess(w, r, cfg.AdminPassword) {
				return
			}
			switch r.Method {
			case http.MethodPut, http.MethodPatch:
				cfg.AuthManager.AdminUpdateUserHandler(w, r)
			case http.MethodDelete:
				cfg.AuthManager.AdminDeleteUserHandler(w, r)
			case http.MethodPost:
				if strings.HasSuffix(r.URL.Path, "/reset-password") {
					cfg.AuthManager.AdminResetUserPasswordHandler(w, r)
				} else if strings.HasSuffix(r.URL.Path, "/regenerate-api-key") {
					cfg.AuthManager.AdminRegenerateAPIKeyHandler(w, r)
				} else {
					http.Error(w, "Not found", http.StatusNotFound)
				}
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}))
	}

	guard := func(next http.HandlerFunc) http.HandlerFunc {
		if next == nil {
			return nil
		}
		if cfg.AuthManager != nil {
			return cfg.AuthManager.OptionalAuth(next)
		}
		return next
	}

	registerOptional := func(path string, handler http.HandlerFunc) {
		if handler == nil {
			return
		}
		mux.HandleFunc(path, guard(handler))
	}

	registerOptional("/api/admin/uploads", cfg.AdminUploadsHandler)
	registerOptional("/api/admin/tracks", cfg.AdminTracksHandler)
	registerOptional("/api/admin/backfill", cfg.AdminBackfillHandler)
	registerOptional("/api/admin/backfill-countries", cfg.AdminBackfillCountriesHandler)
	registerOptional("/api/admin/delete", cfg.AdminDeleteTrackHandler)
	registerOptional("/api/admin/delete-multiple", cfg.AdminDeleteMultipleTracksHandler)
	registerOptional("/api/admin/import-from-safecast", cfg.AdminImportFromSafecastHandler)
	registerOptional("/api/admin/import-by-id", cfg.AdminImportByIDHandler)
	registerOptional("/api/admin/tracks/update", cfg.AdminUpdateTrackHandler)
	registerOptional("/api/admin/uploads/update", cfg.AdminUpdateUploadHandler)
	registerOptional("/api/admin/tracks/import-safecast", cfg.AdminImportSafecastMetaHandler)
	registerOptional("/api/admin/cache", cfg.AdminCacheHandler)

	registerOptionalAdmin := func(path string, handler http.HandlerFunc) {
		if handler == nil {
			return
		}
		registerOptional(path, func(w http.ResponseWriter, r *http.Request) {
			if !checkAdminAccess(w, r, cfg.AdminPassword) {
				return
			}
			handler(w, r)
		})
	}

	registerOptionalAdmin("/api/admin/mcp/data", cfg.AdminMCPDataHandler)
	registerOptionalAdmin("/api/admin/mcp/export", cfg.AdminMCPExportHandler)
	registerOptionalAdmin("/api/admin/mcp/delete", cfg.AdminMCPDeleteHandler)
	registerOptionalAdmin("/api/admin/realtime/data", cfg.AdminRealtimeDataHandler)
	registerOptionalAdmin("/api/admin/realtime/export", cfg.AdminRealtimeExportHandler)
	registerOptionalAdmin("/api/admin/realtime/delete", cfg.AdminRealtimeDeleteHandler)
	registerOptionalAdmin("/api/admin/translations/reload", cfg.AdminTranslationsReloadHandler)
	registerOptionalAdmin("/api/admin/translations/", cfg.AdminTranslationByIDHandler)
	registerOptionalAdmin("/api/admin/translations", cfg.AdminTranslationsHandler)
}

// checkAdminAccess returns true if the request is from an admin user or carries the admin password.
func checkAdminAccess(w http.ResponseWriter, r *http.Request, adminPassword string) bool {
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		return true
	}
	if adminPassword != "" && r.URL.Query().Get("password") == adminPassword {
		return true
	}
	http.Error(w, "Unauthorized - Please login as admin or provide password", http.StatusUnauthorized)
	return false
}

// handleUserUploads returns the authenticated user's uploads (paginated) as JSON.
//
// @Summary     List current user uploads
// @Description Returns paginated uploads for the authenticated user.
// @Tags        auth
// @Produce     json
// @Param       limit query int false "Page size"
// @Param       offset query int false "Offset"
// @Success     200 {object} map[string]interface{} "Upload rows"
// @Failure     401 {string} string "Unauthorized"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/user/uploads [get]
func handleUserUploads(w http.ResponseWriter, r *http.Request, cfg RegisterConfig) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if r.Method != http.MethodGet {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
	}
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if cfg.DB == nil || cfg.DB.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			if parsed > 10000 {
				limit = 10000
			} else {
				limit = parsed
			}
		}
	}
	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	internalUserID := fmt.Sprintf("%d", user.ID)
	ctx := r.Context()
	uploads, err := cfg.DB.GetUploadsPaginated(ctx, limit, offset, internalUserID, "")
	if err != nil {
		if cfg.Logf != nil {
			cfg.Logf("Error fetching user uploads: %v", err)
		}
		http.Error(w, "Failed to fetch uploads", http.StatusInternalServerError)
		return
	}
	totalCount, _ := cfg.DB.CountUploads(ctx, internalUserID, "")

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"uploads": uploads,
		"total":   totalCount,
		"limit":   limit,
		"offset":  offset,
	})
}
