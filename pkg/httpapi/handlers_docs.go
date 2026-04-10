// handlers_docs.go — GET /licenses/{mit|cc0}.
package httpapi

import (
	"bytes"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

// license serves plain-text license files.
//
// @Summary     Download project license text
// @Description Returns embedded license content for supported license codes.
// @Tags        web
// @Param       code path string true "License code: mit or cc0"
// @Success     200 {string} string "License text"
// @Failure     404 {string} string "License not found"
// @Failure     405 {string} string "Method not allowed"
// @Router      /licenses/{code} [get]
func (s *Server) license(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	rawPath := strings.TrimPrefix(r.URL.Path, "/licenses/")
	if rawPath == r.URL.Path {
		http.NotFound(w, r)
		return
	}
	code := strings.ToLower(strings.Trim(rawPath, "/"))
	var file string
	switch code {
	case "mit":
		file = "LICENSE"
	case "cc0":
		file = "LICENSE.CC0"
	default:
		http.NotFound(w, r)
		return
	}
	data, err := fs.ReadFile(s.Content, file)
	if err != nil {
		s.Logf("license handler: %s read error: %v", file, err)
		writeJSONError(w, http.StatusInternalServerError, "unable to load license")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.ServeContent(w, r, file, time.Time{}, bytes.NewReader(data))
}
