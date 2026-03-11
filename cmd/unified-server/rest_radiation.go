package main

import (
	"net/http"
	"strconv"
)

// handleRadiation handles GET /api/radiation
//
// @Summary     Find radiation measurements near a location
// @Description Returns historical bGeigie measurements within a radius of the given coordinates, sorted by most recent. Uses PostGIS spatial index for fast queries. Falls back to Safecast REST API if no database is configured.
// @Tags        historical
// @Produce     json
// @Param       lat      query  number  true  "Latitude in decimal degrees (-90 to 90)"
// @Param       lon      query  number  true  "Longitude in decimal degrees (-180 to 180)"
// @Param       radius_m query  number  false "Search radius in meters (25 to 50000)" default(1500)
// @Param       limit    query  integer false "Maximum number of results (1 to 10000)" default(25)
// @Success     200 {object} map[string]interface{} "Radiation measurements with count, source, and query metadata"
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Router      /radiation [get]
func (h *RESTHandler) handleRadiation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	q := r.URL.Query()

	latStr := q.Get("lat")
	lonStr := q.Get("lon")
	if latStr == "" || lonStr == "" {
		writeError(w, http.StatusBadRequest, "lat and lon are required")
		return
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil || lat < -90 || lat > 90 {
		writeError(w, http.StatusBadRequest, "lat must be a number between -90 and 90")
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil || lon < -180 || lon > 180 {
		writeError(w, http.StatusBadRequest, "lon must be a number between -180 and 180")
		return
	}

	radiusM := 1500.0
	if s := q.Get("radius_m"); s != "" {
		radiusM, err = strconv.ParseFloat(s, 64)
		if err != nil || radiusM < 25 || radiusM > 50000 {
			writeError(w, http.StatusBadRequest, "radius_m must be between 25 and 50000")
			return
		}
	}

	limit := 5
	if s := q.Get("limit"); s != "" {
		limit, err = strconv.Atoi(s)
		if err != nil || limit < 1 || limit > 10000 {
			writeError(w, http.StatusBadRequest, "limit must be between 1 and 10000")
			return
		}
	}
	// Hard cap for API consumers (e.g. Custom GPT) that can't handle large responses
	if limit > 10 {
		limit = 10
	}

	if dbAvailable() {
		result, err := queryRadiationDB(r.Context(), lat, lon, radiusM, limit)
		serveMCPResult(w, result, err)
	} else {
		result, err := queryRadiationAPI(r.Context(), lat, lon, radiusM, limit)
		serveMCPResult(w, result, err)
	}
}
