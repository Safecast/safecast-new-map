package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// handleSensors handles GET /api/sensors
//
// @Summary     Discover active fixed sensors
// @Description Lists active fixed radiation sensors (Pointcast, Solarcast, bGeigieZen, etc.) with their location, type, and last reading timestamp. Requires database connection.
// @Tags        realtime
// @Produce     json
// @Param       type    query  string  false "Filter by sensor type (e.g. Pointcast, Solarcast, bGeigieZen)"
// @Param       min_lat query  number  false "Southern boundary latitude" default(-90)
// @Param       max_lat query  number  false "Northern boundary latitude" default(90)
// @Param       min_lon query  number  false "Western boundary longitude" default(-180)
// @Param       max_lon query  number  false "Eastern boundary longitude" default(180)
// @Param       limit   query  integer false "Maximum number of sensors (1 to 1000)" default(50)
// @Success     200 {object} map[string]interface{} "Sensor list with locations and last reading times"
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /sensors [get]
func (h *RESTHandler) handleSensors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !dbAvailable() {
		writeError(w, http.StatusServiceUnavailable, "database connection required for sensor data")
		return
	}

	q := r.URL.Query()

	sensorType := q.Get("type")

	minLat := -90.0
	if s := q.Get("min_lat"); s != "" {
		var err error
		minLat, err = strconv.ParseFloat(s, 64)
		if err != nil || minLat < -90 || minLat > 90 {
			writeError(w, http.StatusBadRequest, "min_lat must be between -90 and 90")
			return
		}
	}
	maxLat := 90.0
	if s := q.Get("max_lat"); s != "" {
		var err error
		maxLat, err = strconv.ParseFloat(s, 64)
		if err != nil || maxLat < -90 || maxLat > 90 {
			writeError(w, http.StatusBadRequest, "max_lat must be between -90 and 90")
			return
		}
	}
	minLon := -180.0
	if s := q.Get("min_lon"); s != "" {
		var err error
		minLon, err = strconv.ParseFloat(s, 64)
		if err != nil || minLon < -180 || minLon > 180 {
			writeError(w, http.StatusBadRequest, "min_lon must be between -180 and 180")
			return
		}
	}
	maxLon := 180.0
	if s := q.Get("max_lon"); s != "" {
		var err error
		maxLon, err = strconv.ParseFloat(s, 64)
		if err != nil || maxLon < -180 || maxLon > 180 {
			writeError(w, http.StatusBadRequest, "max_lon must be between -180 and 180")
			return
		}
	}

	limit := 50
	if s := q.Get("limit"); s != "" {
		var err error
		limit, err = strconv.Atoi(s)
		if err != nil || limit < 1 || limit > 1000 {
			writeError(w, http.StatusBadRequest, "limit must be between 1 and 1000")
			return
		}
	}

	result, err := listSensorsDB(r.Context(), sensorType, minLat, maxLat, minLon, maxLon, limit)
	serveMCPResult(w, result, err)
}

// handleSensor routes /api/sensor/{id}/current and /api/sensor/{id}/history
//
// @Summary     Get readings from a specific sensor
// @Description Routes to current or history endpoint based on the path suffix. Use /current for the latest reading, /history for time-series data.
// @Tags        realtime
// @Produce     json
// @Param       id         path    string  true  "Device identifier"
// @Param       start_date query   string  false "Start date for history (YYYY-MM-DD) — required for /history"
// @Param       end_date   query   string  false "End date for history (YYYY-MM-DD, default: today)"
// @Param       limit      query   integer false "Maximum number of results (1 to 1000)" default(25)
// @Success     200 {object} map[string]interface{} "Sensor readings"
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /sensor/{id}/current [get]
// @Router      /sensor/{id}/history [get]
func (h *RESTHandler) handleSensor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !dbAvailable() {
		writeError(w, http.StatusServiceUnavailable, "database connection required for sensor data")
		return
	}

	// Parse /api/sensor/{id}/current or /api/sensor/{id}/history
	path := strings.TrimPrefix(r.URL.Path, "/api/sensor/")
	var deviceID, action string
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		deviceID = path[:idx]
		action = path[idx+1:]
	} else {
		deviceID = path
	}
	if deviceID == "" {
		writeError(w, http.StatusBadRequest, "device id is required in path: /api/sensor/{id}/current or /history")
		return
	}

	q := r.URL.Query()

	switch action {
	case "current", "":
		limit := 25
		if s := q.Get("limit"); s != "" {
			var err error
			limit, err = strconv.Atoi(s)
			if err != nil || limit < 1 || limit > 1000 {
				writeError(w, http.StatusBadRequest, "limit must be between 1 and 1000")
				return
			}
		}
		result, err := sensorCurrentDB(r.Context(), deviceID, -90, 90, -180, 180, limit)
		serveMCPResult(w, result, err)

	case "history":
		startDateStr := q.Get("start_date")
		if startDateStr == "" {
			writeError(w, http.StatusBadRequest, "start_date is required for sensor history (YYYY-MM-DD)")
			return
		}
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "start_date must be in YYYY-MM-DD format")
			return
		}

		endDateStr := q.Get("end_date")
		if endDateStr == "" {
			endDateStr = time.Now().Format("2006-01-02")
		}
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "end_date must be in YYYY-MM-DD format")
			return
		}
		if endDate.Before(startDate) {
			writeError(w, http.StatusBadRequest, "end_date must be after start_date")
			return
		}

		limit := 200
		if s := q.Get("limit"); s != "" {
			limit, err = strconv.Atoi(s)
			if err != nil || limit < 1 || limit > 10000 {
				writeError(w, http.StatusBadRequest, "limit must be between 1 and 10000")
				return
			}
		}

		result, err := sensorHistoryDB(r.Context(), deviceID, startDate, endDate, limit)
		serveMCPResult(w, result, err)

	default:
		writeError(w, http.StatusNotFound, "unknown sensor endpoint: use /current or /history")
	}
}
