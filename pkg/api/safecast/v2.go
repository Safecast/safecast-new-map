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

// v2Root returns the API root descriptor (v2 path).
//
//	@Summary		API root (v2)
//	@Description	Returns the Safecast API root with name, URI, and subresource links. Same response as /api/v1.
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	RootRails	"Root descriptor with subresource_uris"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2 [get]
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

// v2Measurements lists or creates measurements (v2 path). Same behavior and response shape as /api/v1/measurements.
//
//	@Summary		List measurements (v2)
//	@Description	Returns measurements with optional filters. Also supports POST to create (same body as v1).
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Param			latitude		query	string				false	"Center latitude for distance filter"
//	@Param			longitude		query	string				false	"Center longitude for distance filter"
//	@Param			distance		query	string				false	"Max distance (km) from lat/lon"
//	@Param			captured_after	query	string				false	"Filter by captured_at >="
//	@Param			captured_before	query	string				false	"Filter by captured_at <"
//	@Param			user_id			query	string				false	"Filter by user ID"
//	@Param			device_id		query	string				false	"Filter by device ID"
//	@Param			page			query	string				false	"Page number (default 1)"
//	@Param			per_page		query	string				false	"Per page (default 100)"
//	@Success		200				{array}	MeasurementRails	"List of measurements"
//	@Failure		405				"Method not allowed"
//	@Router			/api/v2/measurements [get]
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

// v2MeasurementByID returns one measurement by ID (v2 path).
//
//	@Summary		Get measurement by ID (v2)
//	@Description	Returns one measurement by numeric ID. Same response as /api/v1/measurements/{id}.
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"Measurement ID"
//	@Success		200	{object}	MeasurementRails	"Measurement"
//	@Failure		404	"Not found"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/measurements/{id} [get]
func (h *Handler) v2MeasurementByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.adapterMeasurementByID(w, r)
}

// v2Count returns total count of measurements (v2 path).
//
//	@Summary		Count measurements (v2)
//	@Description	Returns total count with same query params as list. Same response as /api/v1/count.
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	CountRails	"Count"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/measurements/count [get]
func (h *Handler) v2Count(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.adapterCount(w, r)
}

// v2BgeigieImports lists or creates bGeigie imports (v2 path). Same behavior and response shape as /api/v1/bgeigie_imports.
//
//	@Summary		List bGeigie imports (v2)
//	@Description	Returns bGeigie import records. Also supports POST (same as v1; returns 501 Not Implemented).
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Param			status		query	string				false	"Filter by status"
//	@Param			page		query	string				false	"Page number (default 1)"
//	@Param			per_page	query	string				false	"Per page (default 25)"
//	@Success		200			{array}	BgeigieImportRails	"List of bGeigie imports"
//	@Failure		405			"Method not allowed"
//	@Router			/api/v2/bgeigie_imports [get]
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

// v2BgeigieImportByID returns one bGeigie import by ID (v2 path).
//
//	@Summary		Get bGeigie import by ID (v2)
//	@Description	Returns one bGeigie import by numeric ID. Same response as /api/v1/bgeigie_imports/{id}.
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"Bgeigie import ID"
//	@Success		200	{object}	BgeigieImportRails	"Bgeigie import"
//	@Failure		404	"Not found"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/bgeigie_imports/{id} [get]
func (h *Handler) v2BgeigieImportByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.adapterBgeigieImportByID(w, r)
}

// v2Users lists users (v2 path). Same response as /api/v1/users.
//
//	@Summary		List users (v2)
//	@Description	Returns users. Stub implementation returns empty array. Same as /api/v1/users.
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	UserRails	"List of users (stub: empty)"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/users [get]
func (h *Handler) v2Users(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterUsers(w, r)
}

// v2UserByID returns one user by ID (v2 path). Same behavior as /api/v1/users/{id}.
//
//	@Summary		Get user by ID (v2)
//	@Description	Returns one user by ID. Returns 501 Not Implemented (stub).
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Failure		501	"Not implemented"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/users/{id} [get]
func (h *Handler) v2UserByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterUserByID(w, r)
}

// v2Devices lists devices (v2 path). Same response as /api/v1/devices.
//
//	@Summary		List devices (v2)
//	@Description	Returns devices. Stub implementation returns empty array. Same as /api/v1/devices.
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	DeviceRails	"List of devices (stub: empty)"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/devices [get]
func (h *Handler) v2Devices(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterDevices(w, r)
}

// v2DeviceByID returns one device by ID (v2 path). Same behavior as /api/v1/devices/{id}.
//
//	@Summary		Get device by ID (v2)
//	@Description	Returns one device by ID. Returns 501 Not Implemented (stub).
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Device ID"
//	@Failure		501	"Not implemented"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/devices/{id} [get]
func (h *Handler) v2DeviceByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterDeviceByID(w, r)
}

// v2RadiationIndex returns radiation index data (v2 path). Stub returns empty array.
//
//	@Summary		Radiation index (v2)
//	@Description	Returns radiation index. Stub returns empty array.
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	object	"Stub: empty array"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/radiation_index [get]
func (h *Handler) v2RadiationIndex(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterRadiationIndex(w, r)
}

// v2Ingest returns ingest data (v2 path). Stub returns empty array.
//
//	@Summary		Ingest (v2)
//	@Description	Returns ingest data. Stub returns empty array.
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	object	"Stub: empty array"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/ingest [get]
func (h *Handler) v2Ingest(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterIngest(w, r)
}

// v2DeviceStories returns device stories (v2 path). Stub returns empty array.
//
//	@Summary		List device stories (v2)
//	@Description	Returns device stories. Stub returns empty array.
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	object	"Stub: empty array"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/device_stories [get]
func (h *Handler) v2DeviceStories(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterDeviceStories(w, r)
}

// v2DeviceStoryByID returns one device story by ID (v2 path). Returns 501 Not Implemented.
//
//	@Summary		Get device story by ID (v2)
//	@Description	Returns one device story. Returns 501 Not Implemented.
//	@Tags			API v2
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Device story ID"
//	@Failure		501	"Not implemented"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v2/device_stories/{id} [get]
func (h *Handler) v2DeviceStoryByID(w http.ResponseWriter, r *http.Request) {
	addCORS(w, r)
	h.adapterDeviceStoryByID(w, r)
}
