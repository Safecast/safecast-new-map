package httpresp

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ErrorBody is the standard JSON error envelope used by API endpoints.
type ErrorBody struct {
	Status string      `json:"status,omitempty"`
	Error  string      `json:"error"`
	Code   string      `json:"code,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
}

// WriteJSON writes a JSON response with the provided status code.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteError writes a JSON error response envelope.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	body := ErrorBody{
		Status: "error",
		Error:  message,
		Code:   strings.TrimSpace(code),
	}
	WriteJSON(w, status, body)
}

// RequireMethodJSON verifies that the request method is allowed.
// On mismatch it writes a JSON error body and returns false.
func RequireMethodJSON(w http.ResponseWriter, r *http.Request, allowed ...string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, method := range allowed {
		if r.Method == method {
			return true
		}
	}
	w.Header().Set("Allow", strings.Join(allowed, ", "))
	WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	return false
}
