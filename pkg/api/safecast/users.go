package safecast

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) adapterUsers(w http.ResponseWriter, r *http.Request) {
	if !wantsJSON(r) {
		http.Error(w, "Accept application/json", http.StatusNotAcceptable)
		return
	}
	addCORS(w, r)
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodGet {
		h.adapterUsersList(w, r)
		return
	}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// adapterUsersList returns the list of users (stub: empty array).
//
//	@Summary		List users
//	@Description	Returns users. Stub implementation returns empty array. Also at /users and /api/v2/users.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	UserRails	"List of users (stub: empty)"
//	@Failure		406	"Accept must be application/json"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v1/users [get]
func (h *Handler) adapterUsersList(w http.ResponseWriter, r *http.Request) {
	out := []UserRails{}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// adapterUserByID returns a user by ID (not implemented).
//
//	@Summary		Get user by ID
//	@Description	Returns one user by ID. Returns 501 Not Implemented.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Failure		501	"Not implemented"
//	@Failure		406	"Accept must be application/json"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v1/users/{id} [get]
func (h *Handler) adapterUserByID(w http.ResponseWriter, r *http.Request) {
	if !wantsJSON(r) {
		http.Error(w, "Accept application/json", http.StatusNotAcceptable)
		return
	}
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
