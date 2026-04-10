package httpapi

import (
	"net/http"

	"safecast-new-map/pkg/auth"
)

// LegacyRoutesConfig holds handlers for non-/api legacy endpoints that remain
// part of the public web surface for backward compatibility.
type LegacyRoutesConfig struct {
	AuthManager *auth.Manager
	RequireAuth bool

	UploadHandler          http.HandlerFunc
	UploadProgressHandler  http.HandlerFunc
	GetMarkersHandler      http.HandlerFunc
	StreamMarkersHandler   http.HandlerFunc
	RealtimeHistoryHandler http.HandlerFunc
	TrackByIDHandler       http.HandlerFunc
	TracksByPrefixHandler  http.HandlerFunc
}

// RegisterLegacyRoutes attaches legacy non-/api endpoints to mux.
// Upload auth behavior intentionally matches the previous main() logic:
// protect /upload only when both RequireAuth and AuthManager are set.
func RegisterLegacyRoutes(mux *http.ServeMux, cfg LegacyRoutesConfig) {
	if mux == nil {
		return
	}

	registerOptional := func(path string, handler http.HandlerFunc) {
		if handler == nil {
			return
		}
		mux.HandleFunc(path, handler)
	}

	if cfg.UploadHandler != nil {
		uploadHandler := cfg.UploadHandler
		if cfg.RequireAuth && cfg.AuthManager != nil {
			uploadHandler = cfg.AuthManager.RequireAuth(uploadHandler)
		}
		mux.HandleFunc("/upload", uploadHandler)
	}

	registerOptional("/upload/progress", cfg.UploadProgressHandler)
	registerOptional("/get_markers", cfg.GetMarkersHandler)
	registerOptional("/stream_markers", cfg.StreamMarkersHandler)
	registerOptional("/realtime_history", cfg.RealtimeHistoryHandler)
	registerOptional("/trackid/", cfg.TrackByIDHandler)
	registerOptional("/tracks/", cfg.TracksByPrefixHandler)
}
