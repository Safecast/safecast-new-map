// helpers.go provides shared HTTP response helpers for API handlers.
package httpapi

import (
	"encoding/json"
	"net/http"
)

// requireMethod checks that the request method matches; if not, it writes 405
// with an Allow header and returns false. Returns true when the method matches.
func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return false
}

// writeJSONError sends a JSON body {"status":"error","error":message} with the given status code.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"status": "error", "error": message})
}

// writeJSON sets Content-Type to application/json, writes the status code, and encodes payload as JSON.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
