package safecast

import (
	"net/http"
	"strings"
)

// registerAdapter registers unversioned and v1 routes that adapt to/from Rails format.
func registerAdapter(mux *http.ServeMux, h *Handler) {
	prefix := "/"
	prefixV1 := "/api/v1"

	// Root (skip prefix "/" to avoid overriding map handler; only register /api/v1)
	mux.HandleFunc(prefixV1, h.adapterRoot)
	mux.HandleFunc(prefixV1+"/", h.adapterRoot)

	// Measurements
	mux.HandleFunc(prefix+"measurements", h.adapterMeasurements)
	mux.HandleFunc(prefix+"measurements.json", h.adapterMeasurements)
	mux.HandleFunc(prefix+"measurements/", h.adapterMeasurementByID)
	mux.HandleFunc(prefixV1+"/measurements", h.adapterMeasurements)
	mux.HandleFunc(prefixV1+"/measurements/", h.adapterMeasurementByID)

	// Count
	mux.HandleFunc(prefix+"count", h.adapterCount)
	mux.HandleFunc(prefix+"measurements/count", h.adapterCount)
	mux.HandleFunc(prefix+"measurements/count.json", h.adapterCount)
	mux.HandleFunc(prefixV1+"/measurements/count", h.adapterCount)

	// Bgeigie imports
	mux.HandleFunc(prefix+"bgeigie_imports", h.adapterBgeigieImports)
	mux.HandleFunc(prefix+"bgeigie_imports.json", h.adapterBgeigieImports)
	mux.HandleFunc(prefix+"bgeigie_imports/", h.adapterBgeigieImportByID)
	mux.HandleFunc(prefixV1+"/bgeigie_imports", h.adapterBgeigieImports)
	mux.HandleFunc(prefixV1+"/bgeigie_imports/", h.adapterBgeigieImportByID)

	// Users
	mux.HandleFunc(prefix+"users", h.adapterUsers)
	mux.HandleFunc(prefix+"users.json", h.adapterUsers)
	mux.HandleFunc(prefix+"users/", h.adapterUserByID)
	mux.HandleFunc(prefixV1+"/users", h.adapterUsers)
	mux.HandleFunc(prefixV1+"/users/", h.adapterUserByID)

	// Devices
	mux.HandleFunc(prefix+"devices", h.adapterDevices)
	mux.HandleFunc(prefix+"devices.json", h.adapterDevices)
	mux.HandleFunc(prefix+"devices/", h.adapterDeviceByID)
	mux.HandleFunc(prefixV1+"/devices", h.adapterDevices)
	mux.HandleFunc(prefixV1+"/devices/", h.adapterDeviceByID)

	// Radiation index, ingest, device stories
	mux.HandleFunc(prefix+"radiation_index", h.adapterRadiationIndex)
	mux.HandleFunc(prefixV1+"/radiation_index", h.adapterRadiationIndex)
	mux.HandleFunc(prefix+"ingest", h.adapterIngest)
	mux.HandleFunc(prefixV1+"/ingest", h.adapterIngest)
	mux.HandleFunc(prefix+"device_stories", h.adapterDeviceStories)
	mux.HandleFunc(prefixV1+"/device_stories", h.adapterDeviceStories)
	mux.HandleFunc(prefix+"device_stories/", h.adapterDeviceStoryByID)
	mux.HandleFunc(prefixV1+"/device_stories/", h.adapterDeviceStoryByID)
	mux.HandleFunc(prefix+"airnote/", h.adapterAirnote)
	mux.HandleFunc(prefixV1+"/airnote/", h.adapterAirnote)
}

// wantsJSON returns true if the request accepts JSON.
func wantsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return true
	}
	if strings.HasSuffix(r.URL.Path, ".json") {
		return true
	}
	return false
}
