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

// writeJSONError sends a JSON error body with an automatically-derived error code.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	httpresp.WriteError(w, status, codeFromStatus(status), message)
}

// writeJSON sets Content-Type to application/json, writes the status code, and encodes payload as JSON.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	httpresp.WriteJSON(w, status, payload)
}

// codeFromStatus maps common HTTP status codes to the canonical httpresp error code strings.
func codeFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return httpresp.CodeBadRequest
	case http.StatusMethodNotAllowed:
		return httpresp.CodeMethodNotAllowed
	case http.StatusUnauthorized:
		return httpresp.CodeUnauthorized
	case http.StatusForbidden:
		return httpresp.CodeForbidden
	case http.StatusNotFound:
		return httpresp.CodeNotFound
	case http.StatusConflict:
		return httpresp.CodeConflict
	case http.StatusTooManyRequests:
		return httpresp.CodeRateLimited
	case http.StatusServiceUnavailable:
		return httpresp.CodeUnavailable
	case http.StatusInternalServerError:
		return httpresp.CodeInternal
	case http.StatusRequestTimeout:
		return httpresp.CodeTimeout
	default:
		return "error"
	}
}
