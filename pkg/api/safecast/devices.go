package safecast

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) adapterDevices(w http.ResponseWriter, r *http.Request) {
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
		h.adapterDevicesList(w, r)
		return
	}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// adapterDevicesList returns the list of devices (stub: empty array).
//
//	@Summary		List devices
//	@Description	Returns devices. Stub implementation returns empty array. Also at /devices and /api/v2/devices.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	DeviceRails	"List of devices (stub: empty)"
//	@Failure		406	"Accept must be application/json"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v1/devices [get]
func (h *Handler) adapterDevicesList(w http.ResponseWriter, r *http.Request) {
	out := []DeviceRails{}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// adapterDeviceByID returns a device by ID (not implemented).
//
//	@Summary		Get device by ID
//	@Description	Returns one device by ID. Returns 501 Not Implemented.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Device ID"
//	@Failure		501	"Not implemented"
//	@Failure		406	"Accept must be application/json"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v1/devices/{id} [get]
func (h *Handler) adapterDeviceByID(w http.ResponseWriter, r *http.Request) {
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
