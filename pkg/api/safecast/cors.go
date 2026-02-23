package safecast

import (
	"net/http"
)

// addCORS sets CORS headers per Rails application_controller.rb.
func addCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "https://safecast.org"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*, X-Requested-With")
	w.Header().Set("Access-Control-Max-Age", "100000")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
	}
}
