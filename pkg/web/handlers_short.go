package web

import (
	"net/http"
	"strings"
	"time"
)

// shortRedirect resolves a short code and redirects to the stored long URL.
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
