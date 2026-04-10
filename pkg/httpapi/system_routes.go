package httpapi

import "net/http"

// SystemRoutesConfig holds infrastructure/lifecycle handlers that do not
// belong to public API, legacy UI, or static route groups.
type SystemRoutesConfig struct {
	SelfUpgradeHandler http.Handler
}

// RegisterSystemRoutes attaches infrastructure routes to mux.
func RegisterSystemRoutes(mux *http.ServeMux, cfg SystemRoutesConfig) {
	if mux == nil {
		return
	}
	if cfg.SelfUpgradeHandler != nil {
		mux.Handle("/selfupgrade/", cfg.SelfUpgradeHandler)
	}
}
