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
// Both /static/ and /js/ serve from the embedded StaticFS so they work
// regardless of the process working directory in production.
func RegisterStaticRoutes(mux *http.ServeMux, cfg StaticRoutesConfig) {
	if mux == nil {
		return
	}
	if cfg.StaticFS != nil {
		fileServer := http.FileServer(http.FS(cfg.StaticFS))
		mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
		// /js/foo.js → js/foo.js inside StaticFS (i.e. public_html/js/foo.js).
		mux.Handle("/js/", fileServer)
	} else if cfg.JSDir != "" {
		mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(cfg.JSDir))))
	}
}
