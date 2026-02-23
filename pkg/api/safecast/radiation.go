package safecast

import (
	"encoding/json"
	"net/http"
)

// adapterRadiationIndex returns radiation index data (stub: empty array).
//
//	@Summary		Radiation index
//	@Description	Returns radiation index. Stub returns empty array. Also at /radiation_index.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	object	"Stub: empty array"
//	@Failure		405	"Method not allowed"
//	@Router			/radiation_index [get]
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

// adapterIngest returns ingest data (stub: empty array).
//
//	@Summary		Ingest
//	@Description	Returns ingest data. Stub returns empty array.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	object	"Stub: empty array"
//	@Failure		405	"Method not allowed"
//	@Router			/ingest [get]
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

// adapterDeviceStories returns device stories (stub: empty array).
//
//	@Summary		List device stories
//	@Description	Returns device stories. Stub returns empty array.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	object	"Stub: empty array"
//	@Failure		405	"Method not allowed"
//	@Router			/device_stories [get]
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

// adapterDeviceStoryByID returns a device story by ID (not implemented).
//
//	@Summary		Get device story by ID
//	@Description	Returns one device story. Returns 501 Not Implemented.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Device story ID"
//	@Failure		501	"Not implemented"
//	@Failure		405	"Method not allowed"
//	@Router			/device_stories/{id} [get]
func (h *Handler) adapterDeviceStoryByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// adapterAirnote returns airnote data (not implemented).
//
//	@Summary		Airnote
//	@Description	Returns airnote data. Returns 501 Not Implemented.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Failure		501	"Not implemented"
//	@Failure		405	"Method not allowed"
//	@Router			/airnote [get]
func (h *Handler) adapterAirnote(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
