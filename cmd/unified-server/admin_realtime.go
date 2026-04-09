package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// realtimeColumns defines the valid columns for the realtime_measurements table.
var realtimeColumns = []string{
	"id", "device_id", "device_name", "tube", "transport", "country",
	"value", "unit", "lat", "lon", "measured_at", "fetched_at", "extra",
}

// latestPerDeviceCTE returns a CTE that selects only the most recent row per device_id.
const latestPerDeviceCTE = `WITH latest AS (
  SELECT DISTINCT ON (device_id) ` + "id, device_id, device_name, tube, transport, country, value, unit, lat, lon, measured_at, fetched_at, extra" + `
  FROM realtime_measurements
  ORDER BY device_id, fetched_at DESC
)`

// adminRealtimeDataHandler returns JSON data for the latest reading per device.
// GET /api/admin/realtime/data?limit=50&offset=0&sort=fetched_at&order=desc&search=...
//
// @Summary     Admin realtime device data
// @Description Returns paginated latest realtime device readings.
// @Tags        admin
// @Produce     json
// @Param       limit query int false "Page size"
// @Param       offset query int false "Offset"
// @Param       sort query string false "Sort column"
// @Param       order query string false "Sort order: asc or desc"
// @Param       search query string false "Search term"
// @Success     200 {object} map[string]interface{} "Realtime rows"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/realtime/data [get]
func adminRealtimeDataHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 1000 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	sortCol := r.URL.Query().Get("sort")
	if !isValidColumn(sortCol, realtimeColumns) {
		sortCol = "fetched_at"
	}
	order := strings.ToUpper(r.URL.Query().Get("order"))
	if order != "ASC" {
		order = "DESC"
	}

	search := r.URL.Query().Get("search")
	whereSQL, args := realtimeSearchWhere(search)

	// Count total (latest per device, then filter)
	countQuery := fmt.Sprintf("%s SELECT COUNT(*) FROM latest %s", latestPerDeviceCTE, whereSQL)
	var total int
	if err := db.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("admin realtime count error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": fmt.Sprintf("Query failed: %v", err),
			"data":  []interface{}{},
			"total": 0,
		})
		return
	}

	// Fetch data
	colList := strings.Join(realtimeColumns, ", ")
	dataQuery := fmt.Sprintf("%s SELECT %s FROM latest %s ORDER BY %s %s LIMIT %d OFFSET %d",
		latestPerDeviceCTE, colList, whereSQL, sortCol, order, limit, offset)

	rows, err := db.DB.Query(dataQuery, args...)
	if err != nil {
		log.Printf("admin realtime query error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(realtimeColumns))
		ptrs := make([]interface{}, len(realtimeColumns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			log.Printf("admin realtime scan error: %v", err)
			continue
		}
		row := make(map[string]interface{})
		for i, col := range realtimeColumns {
			if b, ok := values[i].([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = values[i]
			}
		}
		results = append(results, row)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":    results,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"columns": realtimeColumns,
	})
}

// adminRealtimeExportHandler exports the latest reading per device as CSV.
// GET /api/admin/realtime/export?search=...
//
// @Summary     Admin realtime export
// @Description Exports latest realtime device readings as CSV.
// @Tags        admin
// @Produce     text/csv
// @Param       search query string false "Search term"
// @Success     200 {file} file "CSV export"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/realtime/export [get]
func adminRealtimeExportHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}

	search := r.URL.Query().Get("search")
	whereSQL, args := realtimeSearchWhere(search)

	colList := strings.Join(realtimeColumns, ", ")
	query := fmt.Sprintf("%s SELECT %s FROM latest %s ORDER BY fetched_at DESC",
		latestPerDeviceCTE, colList, whereSQL)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Printf("admin realtime export error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=realtime_devices.csv")

	writer := csv.NewWriter(w)
	writer.Write(realtimeColumns)

	for rows.Next() {
		values := make([]interface{}, len(realtimeColumns))
		ptrs := make([]interface{}, len(realtimeColumns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			continue
		}
		record := make([]string, len(realtimeColumns))
		for i, v := range values {
			if v == nil {
				record[i] = ""
			} else if b, ok := v.([]byte); ok {
				record[i] = string(b)
			} else {
				record[i] = fmt.Sprintf("%v", v)
			}
		}
		writer.Write(record)
	}
	writer.Flush()
}

// adminRealtimeDeleteHandler deletes rows from realtime_measurements.
// DELETE /api/admin/realtime/delete?ids=123,456,789
// DELETE /api/admin/realtime/delete?all=true&search=...
//
// @Summary     Admin realtime delete
// @Description Deletes selected or filtered realtime measurements.
// @Tags        admin
// @Produce     json
// @Param       ids query string false "Comma-separated row IDs"
// @Param       all query boolean false "Delete all filtered rows"
// @Param       search query string false "Search term when all=true"
// @Success     200 {object} map[string]interface{} "Delete result"
// @Failure     400 {string} string "Invalid request"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/realtime/delete [delete]
// @Router      /api/admin/realtime/delete [post]
func adminRealtimeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}

	if r.URL.Query().Get("all") == "true" {
		search := r.URL.Query().Get("search")
		whereSQL, args := realtimeSearchWhere(search)

		query := fmt.Sprintf("DELETE FROM realtime_measurements %s", whereSQL)
		result, err := db.DB.Exec(query, args...)
		if err != nil {
			log.Printf("admin realtime delete all error: %v", err)
			writeError(w, http.StatusInternalServerError, "Delete failed")
			return
		}
		affected, _ := result.RowsAffected()
		writeJSON(w, http.StatusOK, map[string]interface{}{"deleted": affected})
		return
	}

	ids := r.URL.Query().Get("ids")
	if ids == "" {
		writeError(w, http.StatusBadRequest, "Missing ids parameter")
		return
	}

	idList := strings.Split(ids, ",")
	placeholders := make([]string, len(idList))
	args := make([]interface{}, len(idList))
	for i, id := range idList {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = strings.TrimSpace(id)
	}

	query := fmt.Sprintf("DELETE FROM realtime_measurements WHERE id IN (%s)", strings.Join(placeholders, ","))
	result, err := db.DB.Exec(query, args...)
	if err != nil {
		log.Printf("admin realtime delete error: %v", err)
		writeError(w, http.StatusInternalServerError, "Delete failed")
		return
	}
	affected, _ := result.RowsAffected()
	writeJSON(w, http.StatusOK, map[string]interface{}{"deleted": affected})
}

// realtimeSearchWhere builds a WHERE clause with parameterized search for PostgreSQL.
func realtimeSearchWhere(search string) (string, []interface{}) {
	if search == "" {
		return "", nil
	}
	var clauses []string
	for _, col := range realtimeColumns {
		clauses = append(clauses, fmt.Sprintf("CAST(%s AS TEXT) ILIKE $1", col))
	}
	whereSQL := "WHERE " + strings.Join(clauses, " OR ")
	return whereSQL, []interface{}{"%" + search + "%"}
}
