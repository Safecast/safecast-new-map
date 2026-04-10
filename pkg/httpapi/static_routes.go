package httpapi

import (
	"io/fs"
	"net/http"
)

// StaticRoutesConfig holds static file serving settings.
type StaticRoutesConfig struct {
	StaticFS fs.FS
	JSDir    string
}

// RegisterStaticRoutes attaches static asset routes to mux.
// /static serves files from embedded assets and /js serves files from a
// filesystem directory for runtime JS worker compatibility.
func RegisterStaticRoutes(mux *http.ServeMux, cfg StaticRoutesConfig) {
	if mux == nil {
		return
	}
	if cfg.StaticFS != nil {
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(cfg.StaticFS))))
	}
	if cfg.JSDir != "" {
		mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(cfg.JSDir))))
	}
}
