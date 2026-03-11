package auth

import "net/http"

// RequireStaticAdminBasic guards an endpoint with optional static BasicAuth.
// It is intended for transitional admin routes that still use configured
// admin password auth instead of session-based policies.
func RequireStaticAdminBasic(adminPassword string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if next == nil {
			http.NotFound(w, r)
			return
		}
		// Preserve handler-level method semantics; only enforce for write calls.
		if r.Method != http.MethodPost || adminPassword == "" {
			next.ServeHTTP(w, r)
			return
		}

		if IsStaticAdminAuthorized(r, adminPassword) {
			next.ServeHTTP(w, r)
			return
		}

		ChallengeStaticAdminBasic(w)
	}
}

// IsStaticAdminAuthorized checks BasicAuth credentials for static admin routes.
func IsStaticAdminAuthorized(r *http.Request, adminPassword string) bool {
	if adminPassword == "" {
		return true
	}
	user, pass, ok := r.BasicAuth()
	return ok && user == "admin" && pass == adminPassword
}

// ChallengeStaticAdminBasic writes a standard BasicAuth challenge response.
func ChallengeStaticAdminBasic(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="safecast"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
