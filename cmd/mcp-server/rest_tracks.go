package main

import (
	"net/http"
	"strconv"
	"strings"
)

// handleTracks handles GET /api/tracks
//
// @Summary     Browse bGeigie measurement tracks
// @Description Lists bGeigie Import tracks (bulk radiation measurement drives). Each track represents measurements from a single bGeigie session. Can filter by year, month, and detector/device name.
// @Tags        historical
// @Produce     json
// @Param       year     query  integer false "Filter by year (2000–2100)"
// @Param       month    query  integer false "Filter by month (1–12, requires year)"
// @Param       detector query  string  false "Filter by detector/device name (e.g., 'bGeigieZen', 'bGeigie', 'Pointcast'). Partial match supported."
// @Param       limit    query  integer false "Maximum number of results (1 to 50000)" default(50)
// @Success     200 {object} map[string]interface{} "Track list with count and filter metadata"
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Router      /tracks [get]
func (h *RESTHandler) handleTracks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	q := r.URL.Query()

	year := 0
	if s := q.Get("year"); s != "" {
		var err error
		year, err = strconv.Atoi(s)
		if err != nil || year < 2000 || year > 2100 {
			writeError(w, http.StatusBadRequest, "year must be between 2000 and 2100")
			return
		}
	}

	month := 0
	if s := q.Get("month"); s != "" {
		var err error
		month, err = strconv.Atoi(s)
		if err != nil || month < 1 || month > 12 {
			writeError(w, http.StatusBadRequest, "month must be between 1 and 12")
			return
		}
		if year == 0 {
			writeError(w, http.StatusBadRequest, "month filter requires year parameter")
			return
		}
	}

	limit := 50
	if s := q.Get("limit"); s != "" {
		var err error
		limit, err = strconv.Atoi(s)
		if err != nil || limit < 1 || limit > 50000 {
			writeError(w, http.StatusBadRequest, "limit must be between 1 and 50000")
			return
		}
	}

	detector := q.Get("detector")

	// DB is always preferred — calling listTracksAPI would call simplemap.safecast.org/api/tracks
	// which is this server itself, causing infinite recursion.
	if dbAvailable() {
		result, err := listTracksDB(r.Context(), year, month, detector, "", limit)
		serveMCPResult(w, result, err)
		return
	}

	// No DB: detector filter impossible without DB
	if detector != "" {
		writeError(w, http.StatusServiceUnavailable, "Detector filtering requires database access")
		return
	}

	// No DB: fall back to API (only valid for years before the MCP server existed as the backend)
	result, err := listTracksAPI(r.Context(), year, month, limit)
	serveMCPResult(w, result, err)
}

// handleTrack handles GET /api/track/{id}
//
// @Summary     Get all measurements from a specific track
// @Description Retrieves radiation measurements recorded during a specific bGeigie drive. Use GET /api/tracks to find track IDs first.
// @Tags        historical
// @Produce     json
// @Param       id    path    string  true  "Track identifier (e.g. 8eh5m1)"
// @Param       from  query   integer false "Start marker ID for filtering"
// @Param       to    query   integer false "End marker ID for filtering"
// @Param       limit query   integer false "Maximum number of results (1 to 10000)" default(200)
// @Success     200 {object} map[string]interface{} "Measurements for the track"
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Router      /track/{id} [get]
func (h *RESTHandler) handleTrack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract track ID from path: /api/track/{id}
	trackID := strings.TrimPrefix(r.URL.Path, "/api/track/")
	if trackID == "" {
		writeError(w, http.StatusBadRequest, "track id is required in path: /api/track/{id}")
		return
	}

	q := r.URL.Query()

	fromID := 0
	if s := q.Get("from"); s != "" {
		var err error
		fromID, err = strconv.Atoi(s)
		if err != nil || fromID < 0 {
			writeError(w, http.StatusBadRequest, "from must be a non-negative integer")
			return
		}
	}

	toID := 0
	if s := q.Get("to"); s != "" {
		var err error
		toID, err = strconv.Atoi(s)
		if err != nil || toID < 0 {
			writeError(w, http.StatusBadRequest, "to must be a non-negative integer")
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
		result, err := getTrackDB(r.Context(), trackID, fromID, toID, limit)
		serveMCPResult(w, result, err)
	} else {
		result, err := getTrackAPI(r.Context(), trackID, fromID, toID, limit)
		serveMCPResult(w, result, err)
	}
}
