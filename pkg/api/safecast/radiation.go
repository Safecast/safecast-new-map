package safecast

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) adapterRadiationIndex(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if wantsJSON(r) {
		out := []map[string]float64{}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
		return
	}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *Handler) adapterIngest(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Rails uses Elasticsearch; stub with empty data
	out := []interface{}{}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) adapterDeviceStories(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	out := []interface{}{}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) adapterDeviceStoryByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *Handler) adapterAirnote(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
