// register.go wires all HTTP routes: web/API handlers plus auth and admin endpoints.
package httpapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
	"safecast-new-map/pkg/httpresp"
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

	// Optional public API handlers that do not belong to auth/admin groups.
	APISensorsHandler       http.HandlerFunc
	APISensorsExportHandler http.HandlerFunc
	APISensorByIDHandler    http.HandlerFunc
	APIFeedbackHandler      http.HandlerFunc
	APITrackInsightsHandler http.HandlerFunc

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
	AdminMCPUpdateHandler            http.HandlerFunc
	AdminRealtimeDataHandler         http.HandlerFunc
	AdminRealtimeExportHandler       http.HandlerFunc
	AdminRealtimeDeleteHandler       http.HandlerFunc
	AdminTranslationsReloadHandler   http.HandlerFunc
	AdminTranslationByIDHandler      http.HandlerFunc
	AdminTranslationsHandler         http.HandlerFunc
	APITourStepsHandler              http.HandlerFunc
	AdminTourStepsReorderHandler     http.HandlerFunc
	AdminTourStepsByIDHandler        http.HandlerFunc
	AdminTourStepsHandler            http.HandlerFunc

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
		mux.HandleFunc(RouteAPIAuthRegister, cfg.AuthManager.RegisterHandler)
		mux.HandleFunc(RouteAPIAuthLogin, cfg.AuthManager.LoginHandler)
		mux.HandleFunc(RouteAPIAuthLogout, cfg.AuthManager.LogoutHandler)
		mux.HandleFunc(RouteAPIAuthForgotPassword, cfg.AuthManager.ForgotPasswordHandler)
		mux.HandleFunc(RouteAPIAuthResetPassword, cfg.AuthManager.ResetPasswordHandler)
		mux.HandleFunc(RouteAPIAuthVerifyEmail, cfg.AuthManager.VerifyEmailHandler)
		mux.HandleFunc(RouteAPIUserProfile, cfg.AuthManager.RequireAuth(cfg.AuthManager.ProfileHandler))
		mux.HandleFunc(RouteAPIUserChangePassword, cfg.AuthManager.RequireAuth(cfg.AuthManager.ChangePasswordHandler))
		mux.HandleFunc(RouteAPIUserUploads, cfg.AuthManager.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			handleUserUploads(w, r, cfg)
		}))

		mux.HandleFunc(RouteAPIAdminUsers, cfg.AuthManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			if !checkAdminAccess(w, r, cfg.AdminPassword) {
				return
			}
			cfg.AuthManager.AdminListUsersHandler(w, r)
		}))
		mux.HandleFunc(RouteAPIAdminUsersCreate, cfg.AuthManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			if !checkAdminAccess(w, r, cfg.AdminPassword) {
				return
			}
			cfg.AuthManager.AdminCreateUserHandler(w, r)
		}))
		mux.HandleFunc(RouteAPIAdminUsersByID, cfg.AuthManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
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
				writeJSONError(w, http.StatusNotFound, "Not found")
			}
		default:
			httpresp.WriteMethodNotAllowed(w)
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

	registerOptional(RouteAPIAdminUploads, cfg.AdminUploadsHandler)
	registerOptional(RouteAPIAdminTracks, cfg.AdminTracksHandler)
	registerOptional(RouteAPIAdminBackfill, cfg.AdminBackfillHandler)
	registerOptional(RouteAPIAdminBackfillCountries, cfg.AdminBackfillCountriesHandler)
	registerOptional(RouteAPIAdminDelete, cfg.AdminDeleteTrackHandler)
	registerOptional(RouteAPIAdminDeleteMultiple, cfg.AdminDeleteMultipleTracksHandler)
	registerOptional(RouteAPIAdminImportFromSafecast, cfg.AdminImportFromSafecastHandler)
	registerOptional(RouteAPIAdminImportByID, cfg.AdminImportByIDHandler)
	registerOptional(RouteAPIAdminTracksUpdate, cfg.AdminUpdateTrackHandler)
	registerOptional(RouteAPIAdminUploadsUpdate, cfg.AdminUpdateUploadHandler)
	registerOptional(RouteAPIAdminTracksImportSafecast, cfg.AdminImportSafecastMetaHandler)
	registerOptional(RouteAPIAdminCache, cfg.AdminCacheHandler)

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

	registerOptionalAdmin(RouteAPIAdminMCPData, cfg.AdminMCPDataHandler)
	registerOptionalAdmin(RouteAPIAdminMCPExport, cfg.AdminMCPExportHandler)
	registerOptionalAdmin(RouteAPIAdminMCPDelete, cfg.AdminMCPDeleteHandler)
	registerOptionalAdmin(RouteAPIAdminMCPUpdate, cfg.AdminMCPUpdateHandler)
	registerOptionalAdmin(RouteAPIAdminRealtimeData, cfg.AdminRealtimeDataHandler)
	registerOptionalAdmin(RouteAPIAdminRealtimeExport, cfg.AdminRealtimeExportHandler)
	registerOptionalAdmin(RouteAPIAdminRealtimeDelete, cfg.AdminRealtimeDeleteHandler)
	registerOptionalAdmin(RouteAPIAdminTranslationsReload, cfg.AdminTranslationsReloadHandler)
	registerOptionalAdmin(RouteAPIAdminTranslationsByID, cfg.AdminTranslationByIDHandler)
	registerOptionalAdmin(RouteAPIAdminTranslations, cfg.AdminTranslationsHandler)
	// Tour step admin endpoints: reorder must be registered before the
	// {id} pattern because mux longest-match prefers the more specific path.
	registerOptionalAdmin(RouteAPIAdminTourStepsReorder, cfg.AdminTourStepsReorderHandler)
	registerOptionalAdmin(RouteAPIAdminTourStepsByID, cfg.AdminTourStepsByIDHandler)
	registerOptionalAdmin(RouteAPIAdminTourSteps, cfg.AdminTourStepsHandler)

	registerOptionalPublic := func(path string, handler http.HandlerFunc) {
		if handler == nil {
			return
		}
		mux.HandleFunc(path, handler)
	}
	registerOptionalPublic(RouteAPISensors, cfg.APISensorsHandler)
	registerOptionalPublic(RouteAPISensorsExport, cfg.APISensorsExportHandler)
	registerOptionalPublic(RouteAPISensorByID, cfg.APISensorByIDHandler)
	registerOptionalPublic(RouteAPIFeedback, cfg.APIFeedbackHandler)
	registerOptionalPublic(RoutePatternAPITrackInsights, cfg.APITrackInsightsHandler)
	registerOptionalPublic(RouteAPITourSteps, cfg.APITourStepsHandler)
}

// checkAdminAccess returns true if the request is from an admin user or carries the admin password.
func checkAdminAccess(w http.ResponseWriter, r *http.Request, adminPassword string) bool {
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		return true
	}
	if adminPassword != "" && r.URL.Query().Get("password") == adminPassword {
		return true
	}
	httpresp.WriteUnauthorized(w, "Unauthorized - Please login as admin or provide password")
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

	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		httpresp.WriteUnauthorized(w, "Unauthorized")
		return
	}
	if cfg.DB == nil || cfg.DB.DB == nil {
		httpresp.WriteUnavailable(w, "Database not available")
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
		httpresp.WriteInternalError(w, "Failed to fetch uploads")
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
