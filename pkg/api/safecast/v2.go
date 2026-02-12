package safecast

import (
	"net/http"
)

// registerV2 registers /api/v2/* routes that call the same core directly (no adapter).
func registerV2(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/api/v2", h.v2Root)
	mux.HandleFunc("/api/v2/", h.v2Root)
	mux.HandleFunc("/api/v2/measurements", h.v2Measurements)
	mux.HandleFunc("/api/v2/measurements/", h.v2MeasurementByID)
	mux.HandleFunc("/api/v2/measurements/count", h.v2Count)
	mux.HandleFunc("/api/v2/bgeigie_imports", h.v2BgeigieImports)
	mux.HandleFunc("/api/v2/bgeigie_imports/", h.v2BgeigieImportByID)
	mux.HandleFunc("/api/v2/users", h.v2Users)
	mux.HandleFunc("/api/v2/users/", h.v2UserByID)
	mux.HandleFunc("/api/v2/devices", h.v2Devices)
	mux.HandleFunc("/api/v2/devices/", h.v2DeviceByID)
	mux.HandleFunc("/api/v2/radiation_index", h.v2RadiationIndex)
	mux.HandleFunc("/api/v2/ingest", h.v2Ingest)
	mux.HandleFunc("/api/v2/device_stories", h.v2DeviceStories)
	mux.HandleFunc("/api/v2/device_stories/", h.v2DeviceStoryByID)
}

func (h *Handler) v2Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v2" && r.URL.Path != "/api/v2/" {
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

func (h *Handler) v2Measurements(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodGet {
		h.adapterMeasurementsList(w, r)
		return
	}
	h.adapterMeasurementsCreate(w, r)
}

func (h *Handler) v2MeasurementByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.adapterMeasurementByID(w, r)
}

func (h *Handler) v2Count(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.adapterCount(w, r)
}

func (h *Handler) v2BgeigieImports(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodGet {
		h.adapterBgeigieImportsList(w, r)
		return
	}
	h.adapterBgeigieImportsCreate(w, r)
}

func (h *Handler) v2BgeigieImportByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.adapterBgeigieImportByID(w, r)
}

func (h *Handler) v2Users(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterUsers(w, r)
}

func (h *Handler) v2UserByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterUserByID(w, r)
}

func (h *Handler) v2Devices(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterDevices(w, r)
}

func (h *Handler) v2DeviceByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterDeviceByID(w, r)
}

func (h *Handler) v2RadiationIndex(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterRadiationIndex(w, r)
}

func (h *Handler) v2Ingest(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterIngest(w, r)
}

func (h *Handler) v2DeviceStories(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterDeviceStories(w, r)
}

func (h *Handler) v2DeviceStoryByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterDeviceStoryByID(w, r)
}
