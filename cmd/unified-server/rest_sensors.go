package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
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

// handleSensorsExport handles GET /api/sensors/export
//
// @Summary     Export all active sensors
// @Description Downloads all active fixed sensors matching the given filters in CSV, JSON, or Excel format. No row limit — returns up to 10 000 devices.
// @Tags        realtime
// @Produce     text/csv application/json application/vnd.ms-excel
// @Param       format  query  string  false "Output format: csv (default), json, xlsx"
// @Param       type    query  string  false "Filter by sensor type"
// @Param       min_lat query  number  false "Southern boundary" default(-90)
// @Param       max_lat query  number  false "Northern boundary" default(90)
// @Param       min_lon query  number  false "Western boundary" default(-180)
// @Param       max_lon query  number  false "Eastern boundary" default(180)
// @Success     200 {string} string "Sensor data file"
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /sensors/export [get]
func (h *RESTHandler) handleSensorsExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !dbAvailable() {
		writeError(w, http.StatusServiceUnavailable, "database connection required for sensor export")
		return
	}

	q := r.URL.Query()

	format := q.Get("format")
	if format == "" {
		format = "csv"
	}
	if format != "csv" && format != "json" && format != "xlsx" && format != "excel" {
		writeError(w, http.StatusBadRequest, "format must be csv, json, or xlsx")
		return
	}
	if format == "excel" {
		format = "xlsx"
	}

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

	realtimeTable, _, err := findRealtimeTable(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if realtimeTable == "" {
		writeError(w, http.StatusServiceUnavailable, "real-time sensor table not found in database")
		return
	}

	sensors, _, err := listSensorsQuery(r.Context(), realtimeTable, sensorType, minLat, maxLat, minLon, maxLon, 10000)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	date := time.Now().UTC().Format("2006-01-02")

	switch format {
	case "csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="safecast_sensors_%s.csv"`, date))
		cw := csv.NewWriter(w)
		_ = cw.Write([]string{"Device_ID", "Type", "Latitude", "Longitude", "Last_Reading"})
		for _, s := range sensors {
			loc, _ := s["location"].(map[string]any)
			lat := fmt.Sprintf("%v", loc["latitude"])
			lon := fmt.Sprintf("%v", loc["longitude"])
			_ = cw.Write([]string{
				fmt.Sprintf("%v", s["device_id"]),
				fmt.Sprintf("%v", s["type"]),
				lat,
				lon,
				fmt.Sprintf("%v", s["last_reading_at"]),
			})
		}
		cw.Flush()

	case "json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="safecast_sensors_%s.json"`, date))
		_ = json.NewEncoder(w).Encode(sensors)

	case "xlsx":
		w.Header().Set("Content-Type", "application/vnd.ms-excel; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="safecast_sensors_%s.xls"`, date))
		xmlEsc := func(s string) string {
			s = strings.ReplaceAll(s, "&", "&amp;")
			s = strings.ReplaceAll(s, "<", "&lt;")
			s = strings.ReplaceAll(s, ">", "&gt;")
			return s
		}
		cell := func(v string) string {
			return `<Cell><Data ss:Type="String">` + xmlEsc(v) + `</Data></Cell>`
		}
		row := func(cells []string) string {
			out := "<Row>"
			for _, c := range cells {
				out += cell(c)
			}
			return out + "</Row>"
		}
		fmt.Fprint(w, `<?xml version="1.0"?><?mso-application progid="Excel.Sheet"?>`)
		fmt.Fprint(w, `<Workbook xmlns="urn:schemas-microsoft-com:office:spreadsheet" xmlns:ss="urn:schemas-microsoft-com:office:spreadsheet">`)
		fmt.Fprint(w, `<Worksheet ss:Name="Safecast Sensors"><Table>`)
		fmt.Fprint(w, row([]string{"Device_ID", "Type", "Latitude", "Longitude", "Last_Reading"}))
		for _, s := range sensors {
			loc, _ := s["location"].(map[string]any)
			fmt.Fprint(w, row([]string{
				fmt.Sprintf("%v", s["device_id"]),
				fmt.Sprintf("%v", s["type"]),
				fmt.Sprintf("%v", loc["latitude"]),
				fmt.Sprintf("%v", loc["longitude"]),
				fmt.Sprintf("%v", s["last_reading_at"]),
			}))
		}
		fmt.Fprint(w, `</Table></Worksheet></Workbook>`)
	}
}
