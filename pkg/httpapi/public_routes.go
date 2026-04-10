package httpapi

// Shared public API route patterns used by both the map listener registrar and
// MCP REST mirror wiring. Keeping these literals centralized reduces drift when
// handlers move between composition layers.
const (
	RouteAPIAuthRegister       = "/api/auth/register"
	RouteAPIAuthLogin          = "/api/auth/login"
	RouteAPIAuthLogout         = "/api/auth/logout"
	RouteAPIAuthForgotPassword = "/api/auth/forgot-password"
	RouteAPIAuthResetPassword  = "/api/auth/reset-password"
	RouteAPIAuthVerifyEmail    = "/api/auth/verify-email"

	RouteAPIUserProfile        = "/api/user/profile"
	RouteAPIUserChangePassword = "/api/user/change-password"
	RouteAPIUserUploads        = "/api/user/uploads"

	RouteAPIAdminUsers       = "/api/admin/users"
	RouteAPIAdminUsersCreate = "/api/admin/users/create"
	RouteAPIAdminUsersByID   = "/api/admin/users/"

	RouteAPIAdminUploads              = "/api/admin/uploads"
	RouteAPIAdminTracks               = "/api/admin/tracks"
	RouteAPIAdminBackfill             = "/api/admin/backfill"
	RouteAPIAdminBackfillCountries    = "/api/admin/backfill-countries"
	RouteAPIAdminDelete               = "/api/admin/delete"
	RouteAPIAdminDeleteMultiple       = "/api/admin/delete-multiple"
	RouteAPIAdminImportFromSafecast   = "/api/admin/import-from-safecast"
	RouteAPIAdminImportByID           = "/api/admin/import-by-id"
	RouteAPIAdminTracksUpdate         = "/api/admin/tracks/update"
	RouteAPIAdminUploadsUpdate        = "/api/admin/uploads/update"
	RouteAPIAdminTracksImportSafecast = "/api/admin/tracks/import-safecast"
	RouteAPIAdminCache                = "/api/admin/cache"
	RouteAPIAdminMCPData              = "/api/admin/mcp/data"
	RouteAPIAdminMCPExport            = "/api/admin/mcp/export"
	RouteAPIAdminMCPDelete            = "/api/admin/mcp/delete"
	RouteAPIAdminMCPUpdate            = "/api/admin/mcp/update"
	RouteAPIAdminRealtimeData         = "/api/admin/realtime/data"
	RouteAPIAdminRealtimeExport       = "/api/admin/realtime/export"
	RouteAPIAdminRealtimeDelete       = "/api/admin/realtime/delete"
	RouteAPIAdminTranslationsReload   = "/api/admin/translations/reload"
	RouteAPIAdminTranslationsByID     = "/api/admin/translations/"
	RouteAPIAdminTranslations         = "/api/admin/translations"

	RouteAPISensors              = "/api/sensors"
	RouteAPISensorsExport        = "/api/sensors/export"
	RouteAPISensorByID           = "/api/sensor/"
	RouteAPIFeedback             = "/api/feedback"
	RoutePatternAPITrackInsights = "GET /api/track/{id}/insights"
)
