package main

import (
	"net/http"
	"strconv"
	"strings"
)

// handleDevice routes /api/device/{id}/history
//
// @Summary     Get historical measurements from a device
// @Description Returns time-series radiation data from a bGeigie mobile device or fixed sensor. Queries both the markers table (bGeigie imports) and realtime_measurements table (fixed sensors).
// @Tags        historical
// @Produce     json
// @Param       id    path    string  true  "Device identifier"
// @Param       days  query   integer false "Days of history to retrieve (1 to 365)" default(30)
// @Param       limit query   integer false "Maximum number of results (1 to 10000)" default(200)
// @Success     200 {object} map[string]interface{} "Device measurements with period metadata"
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Router      /device/{id}/history [get]
func (h *RESTHandler) handleDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract device ID from path: /api/device/{id}/history
	path := strings.TrimPrefix(r.URL.Path, "/api/device/")
	// Strip trailing /history
	path = strings.TrimSuffix(path, "/history")
	deviceID := strings.TrimSuffix(path, "/")
	if deviceID == "" {
		writeError(w, http.StatusBadRequest, "device id is required in path: /api/device/{id}/history")
		return
	}

	q := r.URL.Query()

	days := 30
	if s := q.Get("days"); s != "" {
		var err error
		days, err = strconv.Atoi(s)
		if err != nil || days < 1 || days > 365 {
			writeError(w, http.StatusBadRequest, "days must be between 1 and 365")
			return
		}
	}

	limit := 200
	if s := q.Get("limit"); s != "" {
		var err error
		limit, err = strconv.Atoi(s)
		if err != nil || limit < 1 || limit > 10000 {
			writeError(w, http.StatusBadRequest, "limit must be between 1 and 10000")
			return
		}
	}

	if dbAvailable() {
		result, err := deviceHistoryDB(r.Context(), deviceID, days, limit)
		serveMCPResult(w, result, err)
	} else {
		result, err := deviceHistoryAPI(r.Context(), deviceID, days, limit)
		serveMCPResult(w, result, err)
	}
}
