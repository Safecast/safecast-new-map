package main

import (
	"net/http"
	"strconv"
	"strings"
)

// handleSpectra handles GET /api/spectra
//
// @Summary     Browse gamma spectroscopy records
// @Description Returns spectroscopy metadata (filename, device model, energy range, location) without channel data. Use GET /api/spectrum/{marker_id} to fetch full channel data. Requires database connection.
// @Tags        spectroscopy
// @Produce     json
// @Param       min_lat      query  number  false "Southern boundary latitude (requires all 4 bbox params)"
// @Param       max_lat      query  number  false "Northern boundary latitude (requires all 4 bbox params)"
// @Param       min_lon      query  number  false "Western boundary longitude (requires all 4 bbox params)"
// @Param       max_lon      query  number  false "Eastern boundary longitude (requires all 4 bbox params)"
// @Param       source_format query string  false "Filter by file format (e.g. spe, csv)"
// @Param       device_model  query string  false "Filter by detector model name (partial match)"
// @Param       track_id      query string  false "Filter by track identifier"
// @Param       limit         query integer false "Maximum number of results (1 to 500)" default(50)
// @Success     200 {object} map[string]interface{} "Spectroscopy records with metadata"
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /spectra [get]
func (h *RESTHandler) handleSpectra(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !dbAvailable() {
		writeError(w, http.StatusServiceUnavailable, "database connection required for spectroscopy data")
		return
	}

	q := r.URL.Query()

	// Optional bounding box — all four must be provided together.
	_, hasMinLat := q["min_lat"]
	_, hasMaxLat := q["max_lat"]
	_, hasMinLon := q["min_lon"]
	_, hasMaxLon := q["max_lon"]
	hasBBox := hasMinLat || hasMaxLat || hasMinLon || hasMaxLon

	var minLat, maxLat, minLon, maxLon float64
	if hasBBox {
		if !(hasMinLat && hasMaxLat && hasMinLon && hasMaxLon) {
			writeError(w, http.StatusBadRequest, "all four bbox parameters (min_lat, max_lat, min_lon, max_lon) must be provided together")
			return
		}
		var err error
		minLat, err = strconv.ParseFloat(q.Get("min_lat"), 64)
		if err != nil || minLat < -90 || minLat > 90 {
			writeError(w, http.StatusBadRequest, "min_lat must be between -90 and 90")
			return
		}
		maxLat, err = strconv.ParseFloat(q.Get("max_lat"), 64)
		if err != nil || maxLat < -90 || maxLat > 90 {
			writeError(w, http.StatusBadRequest, "max_lat must be between -90 and 90")
			return
		}
		minLon, err = strconv.ParseFloat(q.Get("min_lon"), 64)
		if err != nil || minLon < -180 || minLon > 180 {
			writeError(w, http.StatusBadRequest, "min_lon must be between -180 and 180")
			return
		}
		maxLon, err = strconv.ParseFloat(q.Get("max_lon"), 64)
		if err != nil || maxLon < -180 || maxLon > 180 {
			writeError(w, http.StatusBadRequest, "max_lon must be between -180 and 180")
			return
		}
		if minLat >= maxLat {
			writeError(w, http.StatusBadRequest, "min_lat must be less than max_lat")
			return
		}
		if minLon >= maxLon {
			writeError(w, http.StatusBadRequest, "min_lon must be less than max_lon")
			return
		}
	}

	sourceFormat := q.Get("source_format")
	deviceModel := q.Get("device_model")
	trackID := q.Get("track_id")

	limit := 50
	if s := q.Get("limit"); s != "" {
		var err error
		limit, err = strconv.Atoi(s)
		if err != nil || limit < 1 || limit > 500 {
			writeError(w, http.StatusBadRequest, "limit must be between 1 and 500")
			return
		}
	}

	result, err := listSpectraDB(r.Context(), hasBBox, minLat, maxLat, minLon, maxLon, sourceFormat, deviceModel, trackID, limit)
	serveMCPResult(w, result, err)
}

// handleSpectrum handles GET /api/spectrum/{marker_id}
//
// @Summary     Get full spectroscopy channel data
// @Description Returns complete gamma spectroscopy data including all channel counts for the specified measurement point. Use GET /api/spectra to find marker IDs.
// @Tags        spectroscopy
// @Produce     json
// @Param       marker_id path integer true "Marker/measurement identifier"
// @Success     200 {object} map[string]interface{} "Full spectroscopy channel data"
// @Failure     400 {object} map[string]string "Invalid marker_id"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /spectrum/{marker_id} [get]
func (h *RESTHandler) handleSpectrum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract marker_id from path: /api/spectrum/{marker_id}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/spectrum/")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "marker_id is required in path: /api/spectrum/{marker_id}")
		return
	}
	markerID, err := strconv.Atoi(idStr)
	if err != nil || markerID < 1 {
		writeError(w, http.StatusBadRequest, "marker_id must be a positive integer")
		return
	}

	if dbAvailable() {
		result, err := getSpectrumDB(r.Context(), markerID)
		serveMCPResult(w, result, err)
	} else {
		result, err := getSpectrumAPI(r.Context(), markerID)
		serveMCPResult(w, result, err)
	}
}
