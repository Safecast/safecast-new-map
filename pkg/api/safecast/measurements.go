package safecast

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
)

// --- Core service (version-neutral) ---

// ListMeasurementsCore returns measurements from the database using filters.
func (h *Handler) ListMeasurementsCore(ctx context.Context, filters database.MeasurementFilters) ([]database.MeasurementRow, error) {
	return h.DB.QueryMarkersAsMeasurements(ctx, filters, h.DBType)
}

// CountMeasurementsCore returns total marker count.
func (h *Handler) CountMeasurementsCore(ctx context.Context, filters database.MeasurementFilters) (int64, error) {
	return h.DB.CountMarkers(ctx, filters, h.DBType)
}

// GetMeasurementByIDCore returns a single measurement by ID.
func (h *Handler) GetMeasurementByIDCore(ctx context.Context, id int64) (*database.MeasurementRow, error) {
	filters := database.MeasurementFilters{ID: &id, PerPage: 1, Page: 1}
	rows, err := h.DB.QueryMarkersAsMeasurements(ctx, filters, h.DBType)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

// --- Helper: parse filters from request ---

func parseMeasurementFilters(r *http.Request) database.MeasurementFilters {
	f := database.MeasurementFilters{PerPage: 100, Page: 1}
	if v := r.URL.Query().Get("latitude"); v != "" {
		if lat, err := strconv.ParseFloat(v, 64); err == nil {
			f.Latitude = &lat
		}
	}
	if v := r.URL.Query().Get("longitude"); v != "" {
		if lon, err := strconv.ParseFloat(v, 64); err == nil {
			f.Longitude = &lon
		}
	}
	if v := r.URL.Query().Get("distance"); v != "" {
		if d, err := strconv.Atoi(v); err == nil {
			f.Distance = &d
		}
	}
	if v := r.URL.Query().Get("captured_after"); v != "" {
		if t, err := parseTime(v); err == nil {
			f.CapturedAfter = &t
		}
	}
	if v := r.URL.Query().Get("captured_before"); v != "" {
		if t, err := parseTime(v); err == nil {
			f.CapturedBefore = &t
		}
	}
	if v := r.URL.Query().Get("user_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.UserID = &id
		}
	}
	if v := r.URL.Query().Get("device_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.DeviceID = &id
		}
	}
	if v := r.URL.Query().Get("measurement_import_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.MeasurementImportID = &id
		}
	}
	if v := r.URL.Query().Get("original_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.OriginalID = &id
		}
	}
	if v := r.URL.Query().Get("since"); v != "" {
		if t, err := parseTime(v); err == nil {
			f.Since = &t
		}
	}
	if v := r.URL.Query().Get("until"); v != "" {
		if t, err := parseTime(v); err == nil {
			f.Until = &t
		}
	}
	f.Unit = r.URL.Query().Get("unit")
	f.Order = r.URL.Query().Get("order")
	if v := r.URL.Query().Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p >= 1 {
			f.Page = p
		}
	}
	if v := r.URL.Query().Get("per_page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			f.PerPage = p
		}
	}
	return f
}

func parseTime(s string) (int64, error) {
	formats := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05 -0700",
		"2006-01-02",
	}
	for _, fmt := range formats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t.Unix(), nil
		}
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, nil
	}
	return 0, strconv.ErrSyntax
}

// measurementRowToRails converts a MeasurementRow to Rails-compatible JSON shape.
func measurementRowToRails(m database.MeasurementRow) MeasurementRails {
	r := MeasurementRails{
		ID:           m.ID,
		Value:        m.Value,
		Unit:         m.Unit,
		LocationName: m.LocationName,
		Latitude:     m.Latitude,
		Longitude:    m.Longitude,
		OriginalID:   m.OriginalID,
	}
	r.CapturedAt = time.Unix(m.CapturedAt, 0).Format("2006/01/02 15:04:05 -0700")
	r.UserID = m.UserID
	r.DeviceID = m.DeviceID
	r.MeasurementImportID = m.MeasurementImportID
	return r
}

// --- Adapter handlers ---

func (h *Handler) adapterMeasurements(w http.ResponseWriter, r *http.Request) {
	if !wantsJSON(r) {
		http.Error(w, "Accept application/json", http.StatusNotAcceptable)
		return
	}
	addCORS(w, r)

	switch r.Method {
	case http.MethodGet:
		h.adapterMeasurementsList(w, r)
	case http.MethodPost:
		h.adapterMeasurementsCreate(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// adapterMeasurementsList returns measurements matching the query filters.
//
//	@Summary		List measurements
//	@Description	Returns measurements with optional filters. Also at /measurements and /api/v2/measurements.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Param			latitude				query	string				false	"Center latitude for distance filter"
//	@Param			longitude				query	string				false	"Center longitude for distance filter"
//	@Param			distance				query	string				false	"Max distance (km) from lat/lon"
//	@Param			captured_after			query	string				false	"Filter by captured_at >= (ISO8601 or timestamp)"
//	@Param			captured_before			query	string				false	"Filter by captured_at <"
//	@Param			user_id					query	string				false	"Filter by user ID"
//	@Param			device_id				query	string				false	"Filter by device ID"
//	@Param			measurement_import_id	query	string				false	"Filter by measurement_import_id"
//	@Param			original_id				query	string				false	"Filter by original_id"
//	@Param			since					query	string				false	"Alias for captured_after"
//	@Param			until					query	string				false	"Alias for captured_before"
//	@Param			unit					query	string				false	"Filter by unit (e.g. cpm)"
//	@Param			order					query	string				false	"Sort order"
//	@Param			page					query	string				false	"Page number (default 1)"
//	@Param			per_page				query	string				false	"Per page (default 100)"
//	@Success		200						{array}	MeasurementRails	"List of measurements"
//	@Failure		406						"Accept must be application/json"
//	@Failure		405						"Method not allowed"
//	@Router			/api/v1/measurements [get]
func (h *Handler) adapterMeasurementsList(w http.ResponseWriter, r *http.Request) {
	filters := parseMeasurementFilters(r)
	rows, err := h.ListMeasurementsCore(r.Context(), filters)
	if err != nil {
		h.Logf("safecast measurements list: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Convert to Rails format
	out := make([]MeasurementRails, len(rows))
	for i := range rows {
		out[i] = measurementRowToRails(rows[i])
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// adapterMeasurementsCreate creates a single measurement (api_key or X-API-Key required).
//
//	@Summary		Create measurement
//	@Description	Creates a measurement. Requires api_key query or X-API-Key header. Also at /measurements and /api/v2/measurements.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Param			api_key	query		string				false	"API key (or use X-API-Key header)"
//	@Param			body	body		object				true	"JSON: { \"measurement\": { \"value\", \"unit\", \"latitude\", \"longitude\", \"captured_at\", \"location_name\", \"device_id\", \"height\" } }"
//	@Success		200		{object}	MeasurementRails	"Created measurement"
//	@Failure		401		"Unauthorized (missing or invalid API key)"
//	@Failure		422		{object}	ErrorsRails	"Validation errors (invalid JSON or value can't be blank)"
//	@Failure		406		"Accept must be application/json"
//	@Failure		405		"Method not allowed"
//	@Router			/api/v1/measurements [post]
func (h *Handler) adapterMeasurementsCreate(w http.ResponseWriter, r *http.Request) {
	apiKey := r.URL.Query().Get("api_key")
	if apiKey == "" {
		apiKey = r.Header.Get("X-API-Key")
	}
	if apiKey == "" && h.Auth != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var body struct {
		Measurement struct {
			Value        float64  `json:"value"`
			Unit         string   `json:"unit"`
			Latitude     float64  `json:"latitude"`
			Longitude    float64  `json:"longitude"`
			CapturedAt   string   `json:"captured_at"`
			LocationName string   `json:"location_name"`
			DeviceID     *int64   `json:"device_id"`
			Height       *float64 `json:"height"`
		} `json:"measurement"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorsRails{Errors: map[string][]string{"base": {"invalid JSON"}}})
		return
	}
	m := body.Measurement
	if m.Value == 0 && m.Unit == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorsRails{Errors: map[string][]string{"value": {"can't be blank"}}})
		return
	}
	if m.Unit == "" {
		m.Unit = "cpm"
	}
	capturedAt := time.Now().Unix()
	if m.CapturedAt != "" {
		if t, err := parseTime(m.CapturedAt); err == nil {
			capturedAt = t
		}
	}
	userID := ""
	if h.Auth != nil && apiKey != "" {
		if u, err := auth.GetUserByAPIKey(r.Context(), h.Auth.DB, h.Auth.DBDriver, apiKey); err == nil && u != nil {
			userID = strconv.FormatInt(u.ID, 10)
		}
	}
	markerID, err := h.DB.CreateSingleMeasurement(r.Context(), m.Value, m.Unit, m.Latitude, m.Longitude, capturedAt, userID, h.DBType)
	if err != nil {
		h.Logf("safecast POST measurement: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filters := database.MeasurementFilters{ID: &markerID, PerPage: 1, Page: 1}
	rows, err := h.DB.QueryMarkersAsMeasurements(r.Context(), filters, h.DBType)
	if err != nil || len(rows) == 0 {
		h.Logf("safecast fetch created measurement: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	out := measurementRowToRails(rows[0])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// adapterMeasurementByID returns a single measurement by ID.
//
//	@Summary		Get measurement by ID
//	@Description	Returns one measurement by numeric ID. Also at /measurements/:id and /api/v2/measurements/:id.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"Measurement ID"
//	@Success		200	{object}	MeasurementRails	"Measurement"
//	@Failure		404	"Not found"
//	@Failure		406	"Accept must be application/json"
//	@Failure		405	"Method not allowed"
//	@Router			/api/v1/measurements/{id} [get]
func (h *Handler) adapterMeasurementByID(w http.ResponseWriter, r *http.Request) {
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
	path := strings.Trim(r.URL.Path, "/")
	path = strings.TrimSuffix(path, ".json")
	parts := strings.Split(path, "/")
	idStr := ""
	for i, p := range parts {
		if p == "measurements" && i+1 < len(parts) {
			idStr = parts[i+1]
			break
		}
	}
	if idStr == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	row, err := h.GetMeasurementByIDCore(r.Context(), id)
	if err != nil {
		h.Logf("safecast measurement get: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if row == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	out := measurementRowToRails(*row)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// adapterCount returns the total count of measurements matching the query filters.
//
//	@Summary		Count measurements
//	@Description	Returns total count with same query params as list. Also at /count and /api/v2/count.
//	@Tags			API v1
//	@Accept			json
//	@Produce		json
//	@Param			latitude				query		string		false	"Center latitude"
//	@Param			longitude				query		string		false	"Center longitude"
//	@Param			distance				query		string		false	"Max distance (km)"
//	@Param			captured_after			query		string		false	"Filter captured_at >="
//	@Param			captured_before			query		string		false	"Filter captured_at <"
//	@Param			user_id					query		string		false	"Filter by user ID"
//	@Param			device_id				query		string		false	"Filter by device ID"
//	@Param			measurement_import_id	query		string		false	"Filter by measurement_import_id"
//	@Param			original_id				query		string		false	"Filter by original_id"
//	@Param			since					query		string		false	"Alias for captured_after"
//	@Param			until					query		string		false	"Alias for captured_before"
//	@Param			unit					query		string		false	"Filter by unit"
//	@Param			order					query		string		false	"Sort order"
//	@Param			page					query		string		false	"Page (for consistency; count ignores)"
//	@Param			per_page				query		string		false	"Per page (for consistency; count ignores)"
//	@Success		200						{object}	CountRails	"Count"
//	@Failure		406						"Accept must be application/json"
//	@Failure		405						"Method not allowed"
//	@Router			/api/v1/count [get]
func (h *Handler) adapterCount(w http.ResponseWriter, r *http.Request) {
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
	filters := parseMeasurementFilters(r)
	count, err := h.CountMeasurementsCore(r.Context(), filters)
	if err != nil {
		h.Logf("safecast count: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CountRails{Count: count})
}
