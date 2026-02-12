package safecast

import (
	"net/http"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
)

// Handler wires the Safecast API (api.safecast.org replacement) to the database and auth.
// It registers both the adapter (unversioned + v1) and v2 routes; both call the same core logic.
type Handler struct {
	DB      *database.Database
	DBType  string
	Auth    *auth.Manager
	BaseURL string
	Logf    func(string, ...any)
}

// NewHandler constructs a Handler with the given dependencies.
func NewHandler(db *database.Database, dbType string, authMgr *auth.Manager, baseURL string, logf func(string, ...any)) *Handler {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	return &Handler{
		DB:      db,
		DBType:  dbType,
		Auth:    authMgr,
		BaseURL: baseURL,
		Logf:    logf,
	}
}

// Register attaches Safecast API routes to the provided mux.
// Adapter routes: unversioned + v1 (Rails contract).
// v2 routes: direct to core (same logic, different HTTP format).
func (h *Handler) Register(mux *http.ServeMux) {
	registerAdapter(mux, h)
	registerV2(mux, h)
}
