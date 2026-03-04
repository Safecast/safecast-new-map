// handlers_short.go — short URL redirects
//
// Resolves short codes (e.g. /s/abc123) to full URLs stored in the database
// and redirects the client. Used for shareable, compact links.
package web

import (
	"net/http"
	"strings"
	"time"
)

// shortRedirect looks up the short code from the path (e.g. /s/xyz), fetches
// the target URL from the database, and issues a 302 redirect. Returns 404
// if the code is unknown or empty.
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
		http.Error(w, "short link lookup failed", http.StatusInternalServerError)
		return
	}
	if strings.TrimSpace(target) == "" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}
