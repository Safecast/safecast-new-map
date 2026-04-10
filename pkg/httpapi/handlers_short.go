// handlers_short.go — GET /s/{code} redirects to the stored URL for the short code.
package httpapi

import (
	"net/http"
	"strings"
	"time"
)

// shortRedirect resolves and redirects a short code.
//
// @Summary     Resolve short URL code
// @Description Looks up a short code and redirects to its stored target URL.
// @Tags        web
// @Success     302 {string} string "Redirect"
// @Failure     404 {string} string "Code not found"
// @Failure     502 {string} string "Invalid redirect target"
// @Router      /s/{code} [get]
func (s *Server) shortRedirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/s/"))
	if code == "" {
		http.NotFound(w, r)
		return
	}
	if s.DB == nil || s.DB.DB == nil {
		http.NotFound(w, r)
		return
	}
	ctx, cancel := WithMinimumDeadline(r.Context(), 30*time.Second)
	defer cancel()
	target, err := s.DB.ResolveShortLink(ctx, code)
	if err != nil {
		s.Logf("short link lookup for %q failed: %v", code, err)
		writeJSONError(w, http.StatusInternalServerError, "short link lookup failed")
		return
	}
	if strings.TrimSpace(target) == "" {
		http.NotFound(w, r)
		return
	}
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		writeJSONError(w, http.StatusBadGateway, "invalid redirect target")
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}
