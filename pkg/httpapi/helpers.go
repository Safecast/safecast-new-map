// helpers.go provides shared HTTP response helpers for API handlers.
package httpapi

import (
	"net/http"
	"safecast-new-map/pkg/httpresp"
)

// requireMethod checks that the request method matches; if not, it writes 405
// with an Allow header and returns false. Returns true when the method matches.
func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	return httpresp.RequireMethodJSON(w, r, method)
}

// writeJSONError sends a JSON body {"status":"error","error":message} with the given status code.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	httpresp.WriteError(w, status, "", message)
}

// writeJSON sets Content-Type to application/json, writes the status code, and encodes payload as JSON.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	httpresp.WriteJSON(w, status, payload)
}
