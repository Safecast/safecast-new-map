package main

import (
	"net/http"
	"strconv"
)

// handleArea handles GET /api/area
//
// @Summary     Find radiation measurements within a bounding box
// @Description Returns historical bGeigie measurements inside the specified geographic rectangle. Falls back to Safecast REST API if no database is configured.
// @Tags        historical
// @Produce     json
// @Param       min_lat query  number  true  "Southern boundary latitude (-90 to 90)"
// @Param       max_lat query  number  true  "Northern boundary latitude (-90 to 90)"
// @Param       min_lon query  number  true  "Western boundary longitude (-180 to 180)"
// @Param       max_lon query  number  true  "Eastern boundary longitude (-180 to 180)"
// @Param       limit   query  integer false "Maximum number of results (1 to 10000)" default(100)
// @Success     200 {object} map[string]interface{} "Measurements with count, bbox, and source"
// @Failure     400 {object} map[string]string "Invalid or missing parameters"
// @Router      /area [get]
func (h *RESTHandler) handleArea(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	q := r.URL.Query()

	parseRequired := func(key string, min, max float64) (float64, bool) {
		s := q.Get(key)
		if s == "" {
			writeError(w, http.StatusBadRequest, key+" is required")
			return 0, false
		}
		v, err := strconv.ParseFloat(s, 64)
		if err != nil || v < min || v > max {
			writeError(w, http.StatusBadRequest, key+" must be between "+strconv.FormatFloat(min, 'f', 0, 64)+" and "+strconv.FormatFloat(max, 'f', 0, 64))
			return 0, false
		}
		return v, true
	}

	minLat, ok := parseRequired("min_lat", -90, 90)
	if !ok {
		return
	}
	maxLat, ok := parseRequired("max_lat", -90, 90)
	if !ok {
		return
	}
	minLon, ok := parseRequired("min_lon", -180, 180)
	if !ok {
		return
	}
	maxLon, ok := parseRequired("max_lon", -180, 180)
	if !ok {
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

	limit := 5
	if s := q.Get("limit"); s != "" {
		var err error
		limit, err = strconv.Atoi(s)
		if err != nil || limit < 1 || limit > 10000 {
			writeError(w, http.StatusBadRequest, "limit must be between 1 and 10000")
			return
		}
	}
	if limit > 10 {
		limit = 10
	}

	if dbAvailable() {
		result, err := searchAreaDB(r.Context(), minLat, maxLat, minLon, maxLon, limit)
		serveMCPResult(w, result, err)
	} else {
		result, err := searchAreaAPI(r.Context(), minLat, maxLat, minLon, maxLon, limit)
		serveMCPResult(w, result, err)
	}
}
