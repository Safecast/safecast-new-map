package safecast

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handler) adapterRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1" && r.URL.Path != "/api/v1/" {
		return
	}
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
	writeRootJSON(w, h.BaseURL)
}

// writeRootJSON writes the root API descriptor (shared by v1 and v2).
func writeRootJSON(w http.ResponseWriter, baseURL string) {
	base := baseURL
	if base == "" {
		base = "https://api.safecast.org"
	}
	base = strings.TrimSuffix(base, "/")
	out := RootRails{
		Name:           "Safecast API",
		URI:            base + "/",
		SubresourceURIs: []string{
			base + "/users.json",
			base + "/measurements.json",
			base + "/bgeigie_imports.json",
			base + "/devices.json",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}
