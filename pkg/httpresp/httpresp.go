package httpresp

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Standard error code constants used across all API endpoints.
const (
	CodeBadRequest       = "bad_request"
	CodeMethodNotAllowed = "method_not_allowed"
	CodeUnauthorized     = "unauthorized"
	CodeForbidden        = "forbidden"
	CodeNotFound         = "not_found"
	CodeConflict         = "conflict"
	CodeRateLimited      = "rate_limited"
	CodeUnavailable      = "service_unavailable"
	CodeInternal         = "internal_error"
	CodeTimeout          = "timeout"
	CodeCancelled        = "cancelled"
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
	WriteError(w, http.StatusMethodNotAllowed, CodeMethodNotAllowed, "Method not allowed")
	return false
}

// WriteMethodNotAllowed writes a 405 JSON error with an Allow header.
func WriteMethodNotAllowed(w http.ResponseWriter, allowed ...string) {
	if len(allowed) > 0 {
		w.Header().Set("Allow", strings.Join(allowed, ", "))
	}
	WriteError(w, http.StatusMethodNotAllowed, CodeMethodNotAllowed, "Method not allowed")
}

// WriteBadRequest writes a 400 JSON error with the given code and message.
func WriteBadRequest(w http.ResponseWriter, code, msg string) {
	if code == "" {
		code = CodeBadRequest
	}
	WriteError(w, http.StatusBadRequest, code, msg)
}

// WriteUnauthorized writes a 401 JSON error.
func WriteUnauthorized(w http.ResponseWriter, msg string) {
	WriteError(w, http.StatusUnauthorized, CodeUnauthorized, msg)
}

// WriteForbidden writes a 403 JSON error.
func WriteForbidden(w http.ResponseWriter, msg string) {
	WriteError(w, http.StatusForbidden, CodeForbidden, msg)
}

// WriteNotFound writes a 404 JSON error.
func WriteNotFound(w http.ResponseWriter, msg string) {
	WriteError(w, http.StatusNotFound, CodeNotFound, msg)
}

// WriteRateLimited writes a 429 JSON error.
func WriteRateLimited(w http.ResponseWriter, msg string) {
	WriteError(w, http.StatusTooManyRequests, CodeRateLimited, msg)
}

// WriteUnavailable writes a 503 JSON error.
func WriteUnavailable(w http.ResponseWriter, msg string) {
	WriteError(w, http.StatusServiceUnavailable, CodeUnavailable, msg)
}

// WriteInternalError writes a 500 JSON error.
func WriteInternalError(w http.ResponseWriter, msg string) {
	WriteError(w, http.StatusInternalServerError, CodeInternal, msg)
}

// WriteTimeout writes a 408 JSON error.
func WriteTimeout(w http.ResponseWriter, msg string) {
	WriteError(w, http.StatusRequestTimeout, CodeTimeout, msg)
}
