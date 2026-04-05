// admin_anomaly.go — admin handlers for anomaly detection config and results.
//
// Routes (registered in main.go):
//   GET  /admin/anomaly                     — serve HTML page
//   GET  /api/admin/anomaly/config          — read current config
//   POST /api/admin/anomaly/config          — save config
//   GET  /api/admin/anomaly/data            — paginated anomaly summaries from DuckLake
//   POST /api/admin/anomaly/rerun/{trackID} — re-run analysis for one track

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// anomalySettings holds the runtime anomaly detection configuration
// for both dose-rate tracks and spectrum tracks.
type anomalySettings struct {
	SigmaThreshold       float64 `json:"sigma_threshold"`
	MinDoserate          float64 `json:"min_doserate"`
	MaxDoserate          float64 `json:"max_doserate"`
	MaxAnomalyIDs        int     `json:"max_anomaly_ids"`
	SpectrumCosineThresh float64 `json:"spectrum_cosine_threshold"`
	SpectrumMinLiveTime  float64 `json:"spectrum_min_live_time"`
}

var (
	anomalyConfig   anomalySettings
	anomalyConfigMu sync.RWMutex
)

func init() {
	anomalyConfig = anomalySettings{
		SigmaThreshold:       3.0,
		MinDoserate:          0.0,
		MaxDoserate:          10000.0,
		MaxAnomalyIDs:        50,
		SpectrumCosineThresh: 0.85,
		SpectrumMinLiveTime:  1.0,
	}
}

// getAnomalyConfig returns a safe copy of the current config.
func getAnomalyConfig() anomalySettings {
	anomalyConfigMu.RLock()
	defer anomalyConfigMu.RUnlock()
	return anomalyConfig
}

// initAnomalyConfig creates the anomaly_config table (if absent) and loads
// the stored values into the in-memory config. Called once at server startup.
func initAnomalyConfig() {
	if db == nil {
		return
	}
	// Create table with defaults if it doesn't exist.
	_, err := db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS anomaly_config (
			id                       INTEGER PRIMARY KEY DEFAULT 1,
			sigma_threshold          DOUBLE PRECISION    DEFAULT 3.0,
			min_doserate             DOUBLE PRECISION    DEFAULT 0.0,
			max_doserate             DOUBLE PRECISION    DEFAULT 10000.0,
			max_anomaly_ids          INTEGER             DEFAULT 50,
			spectrum_cosine_threshold DOUBLE PRECISION   DEFAULT 0.85,
			spectrum_min_live_time   DOUBLE PRECISION    DEFAULT 1.0
		)`)
	if err != nil {
		log.Printf("[anomaly-config] create table: %v", err)
		return
	}
	// Ensure exactly one row exists.
	_, err = db.DB.Exec(`
		INSERT INTO anomaly_config (id, sigma_threshold, min_doserate, max_doserate, max_anomaly_ids, spectrum_cosine_threshold, spectrum_min_live_time)
		VALUES (1, 3.0, 0.0, 10000.0, 50, 0.85, 1.0)
		ON CONFLICT (id) DO NOTHING`)

	// Add new columns to existing installations (idempotent).
	for _, col := range []string{
		"ALTER TABLE anomaly_config ADD COLUMN IF NOT EXISTS spectrum_cosine_threshold DOUBLE PRECISION DEFAULT 0.85",
		"ALTER TABLE anomaly_config ADD COLUMN IF NOT EXISTS spectrum_min_live_time DOUBLE PRECISION DEFAULT 1.0",
	} {
		if _, e := db.DB.Exec(col); e != nil {
			log.Printf("[anomaly-config] migrate: %v", e)
		}
	}
	if err != nil {
		log.Printf("[anomaly-config] seed row: %v", err)
		return
	}
	loadAnomalyConfig()
}

// loadAnomalyConfig reads from PostgreSQL and updates the in-memory config.
func loadAnomalyConfig() {
	if db == nil {
		return
	}
	var cfg anomalySettings
	err := db.DB.QueryRow(`
		SELECT sigma_threshold, min_doserate, max_doserate, max_anomaly_ids,
		       spectrum_cosine_threshold, spectrum_min_live_time
		FROM anomaly_config WHERE id = 1`).
		Scan(&cfg.SigmaThreshold, &cfg.MinDoserate, &cfg.MaxDoserate, &cfg.MaxAnomalyIDs,
			&cfg.SpectrumCosineThresh, &cfg.SpectrumMinLiveTime)
	if err != nil {
		log.Printf("[anomaly-config] load: %v", err)
		return
	}
	anomalyConfigMu.Lock()
	anomalyConfig = cfg
	anomalyConfigMu.Unlock()
	log.Printf("[anomaly-config] loaded: sigma=%.1f min=%.2f max=%.0f maxIDs=%d cosine=%.2f minLive=%.1fs",
		cfg.SigmaThreshold, cfg.MinDoserate, cfg.MaxDoserate, cfg.MaxAnomalyIDs,
		cfg.SpectrumCosineThresh, cfg.SpectrumMinLiveTime)
}

// adminAnomalyConfigGetHandler — GET /api/admin/anomaly/config
func adminAnomalyConfigGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getAnomalyConfig()) //nolint:errcheck
}

// adminAnomalyConfigSaveHandler — POST /api/admin/anomaly/config
func adminAnomalyConfigSaveHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}
	var payload anomalySettings
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Sanity-check inputs.
	if payload.SigmaThreshold < 1.0 || payload.SigmaThreshold > 10.0 {
		http.Error(w, "sigma_threshold must be between 1.0 and 10.0", http.StatusBadRequest)
		return
	}
	if payload.MinDoserate < 0 || payload.MinDoserate >= payload.MaxDoserate {
		http.Error(w, "min_doserate must be >= 0 and less than max_doserate", http.StatusBadRequest)
		return
	}
	if payload.MaxAnomalyIDs < 1 || payload.MaxAnomalyIDs > 500 {
		http.Error(w, "max_anomaly_ids must be between 1 and 500", http.StatusBadRequest)
		return
	}
	if payload.SpectrumCosineThresh < 0.0 || payload.SpectrumCosineThresh > 1.0 {
		http.Error(w, "spectrum_cosine_threshold must be between 0.0 and 1.0", http.StatusBadRequest)
		return
	}
	if payload.SpectrumMinLiveTime < 0 {
		http.Error(w, "spectrum_min_live_time must be >= 0", http.StatusBadRequest)
		return
	}
	_, err := db.DB.Exec(`
		UPDATE anomaly_config
		SET sigma_threshold=$1, min_doserate=$2, max_doserate=$3, max_anomaly_ids=$4,
		    spectrum_cosine_threshold=$5, spectrum_min_live_time=$6
		WHERE id=1`,
		payload.SigmaThreshold, payload.MinDoserate, payload.MaxDoserate, payload.MaxAnomalyIDs,
		payload.SpectrumCosineThresh, payload.SpectrumMinLiveTime)
	if err != nil {
		log.Printf("[anomaly-config] save: %v", err)
		http.Error(w, "Save failed", http.StatusInternalServerError)
		return
	}
	// Reload into memory immediately.
	loadAnomalyConfig()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"}) //nolint:errcheck
}

// adminAnomalyDataHandler — GET /api/admin/anomaly/data
// Returns paginated rows from track_anomaly_summaries in DuckLake.
func adminAnomalyDataHandler(w http.ResponseWriter, r *http.Request) {
	if !duckDBAvailable() {
		http.Error(w, "Analytics not available", http.StatusServiceUnavailable)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	validCols := []string{"track_id", "computed_at", "reading_count", "anomaly_count", "mean_doserate", "threshold"}
	sortCol := r.URL.Query().Get("sort")
	if !isValidColumn(sortCol, validCols) {
		sortCol = "computed_at"
	}
	order := strings.ToUpper(r.URL.Query().Get("order"))
	if order != "ASC" {
		order = "DESC"
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))

	whereSQL := ""
	var args []interface{}
	if search != "" {
		whereSQL = fmt.Sprintf("WHERE track_id LIKE ?")
		args = append(args, "%"+search+"%")
	}

	// Total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM track_anomaly_summaries %s", whereSQL)
	var total int
	if err := duckDB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("[anomaly-admin] count: %v", err)
		total = 0
	}

	// Data
	dataQuery := fmt.Sprintf(`
		SELECT track_id, computed_at, reading_count,
		       mean_doserate, median_doserate, stddev_doserate,
		       threshold, anomaly_count, anomaly_max, anomaly_min
		FROM track_anomaly_summaries %s
		ORDER BY %s %s
		LIMIT %d OFFSET %d`,
		whereSQL, sortCol, order, limit, offset)

	queryArgs := append(args, limit, offset)
	// DuckLake doesn't use $n placeholders for LIMIT/OFFSET in this query form,
	// so just build those inline — they're safe integers from strconv.Atoi.
	dataQuery = fmt.Sprintf(`
		SELECT track_id, computed_at, reading_count,
		       mean_doserate, median_doserate, stddev_doserate,
		       threshold, anomaly_count, anomaly_max, anomaly_min
		FROM track_anomaly_summaries %s
		ORDER BY %s %s
		LIMIT %d OFFSET %d`,
		whereSQL, sortCol, order, limit, offset)
	queryArgs = args

	rows, err := duckDB.Query(dataQuery, queryArgs...)
	if err != nil {
		log.Printf("[anomaly-admin] query: %v", err)
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var (
			trackID                                      string
			computedAt                                   interface{}
			readingCount, anomalyCount                   int
			meanDoserate, medianDoserate, stddevDoserate float64
			threshold, anomalyMax, anomalyMin            float64
		)
		if err := rows.Scan(
			&trackID, &computedAt, &readingCount,
			&meanDoserate, &medianDoserate, &stddevDoserate,
			&threshold, &anomalyCount, &anomalyMax, &anomalyMin,
		); err != nil {
			continue
		}
		results = append(results, map[string]interface{}{
			"track_id":        trackID,
			"computed_at":     computedAt,
			"reading_count":   readingCount,
			"mean_doserate":   meanDoserate,
			"median_doserate": medianDoserate,
			"stddev_doserate": stddevDoserate,
			"threshold":       threshold,
			"anomaly_count":   anomalyCount,
			"anomaly_max":     anomalyMax,
			"anomaly_min":     anomalyMin,
		})
	}
	if results == nil {
		results = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
		"data":  results,
		"total": total,
	})
}

// adminSpectrumAnomalyDataHandler — GET /api/admin/anomaly/spectrum-data
// Returns paginated rows from spectrum_anomaly_summaries in DuckLake.
func adminSpectrumAnomalyDataHandler(w http.ResponseWriter, r *http.Request) {
	if !duckDBAvailable() {
		http.Error(w, "Analytics not available", http.StatusServiceUnavailable)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	validCols := []string{"track_id", "computed_at", "spectrum_count", "anomaly_count", "median_total_counts"}
	sortCol := r.URL.Query().Get("sort")
	if !isValidColumn(sortCol, validCols) {
		sortCol = "computed_at"
	}
	order := strings.ToUpper(r.URL.Query().Get("order"))
	if order != "ASC" {
		order = "DESC"
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	whereSQL := ""
	var args []interface{}
	if search != "" {
		whereSQL = "WHERE track_id LIKE ?"
		args = append(args, "%"+search+"%")
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM spectrum_anomaly_summaries %s", whereSQL)
	if err := duckDB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		total = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT track_id, computed_at, spectrum_count, anomaly_count,
		       median_total_counts, cosine_threshold
		FROM spectrum_anomaly_summaries %s
		ORDER BY %s %s
		LIMIT %d OFFSET %d`,
		whereSQL, sortCol, order, limit, offset)

	rows, err := duckDB.Query(dataQuery, args...)
	if err != nil {
		log.Printf("[spectrum-anomaly-admin] query: %v", err)
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var (
			trackID                    string
			computedAt                 interface{}
			spectrumCount, anomalyCount int
			medianTotal, cosineThresh  float64
		)
		if err := rows.Scan(&trackID, &computedAt, &spectrumCount, &anomalyCount,
			&medianTotal, &cosineThresh); err != nil {
			continue
		}
		results = append(results, map[string]interface{}{
			"track_id":           trackID,
			"computed_at":        computedAt,
			"spectrum_count":     spectrumCount,
			"anomaly_count":      anomalyCount,
			"median_total_counts": medianTotal,
			"cosine_threshold":   cosineThresh,
		})
	}
	if results == nil {
		results = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
		"data":  results,
		"total": total,
	})
}

// adminAnomalyRerunHandler — POST /api/admin/anomaly/rerun/{trackID}
// Triggers re-computation of dose-rate anomaly analysis for a single track.
func adminAnomalyRerunHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	trackID := parts[len(parts)-1]
	if trackID == "" || trackID == "rerun" {
		http.Error(w, "track_id required", http.StatusBadRequest)
		return
	}
	go func() {
		summary, err := computeAnomalies(trackID)
		if err != nil {
			log.Printf("[anomaly] rerun %s: %v", trackID, err)
			return
		}
		if err := storeAnomalySummary(summary); err != nil {
			log.Printf("[anomaly] rerun store %s: %v", trackID, err)
			return
		}
		log.Printf("[anomaly] rerun %s: %d readings, %d anomalies", trackID, summary.ReadingCount, summary.AnomalyCount)
	}()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "queued", "track_id": trackID}) //nolint:errcheck
}

// adminSpectrumAnomalyRerunHandler — POST /api/admin/anomaly/spectrum-rerun/{trackID}
// Triggers re-computation of spectrum anomaly analysis for a single track.
func adminSpectrumAnomalyRerunHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	trackID := parts[len(parts)-1]
	if trackID == "" || trackID == "spectrum-rerun" {
		http.Error(w, "track_id required", http.StatusBadRequest)
		return
	}
	go func() {
		summary, err := computeSpectrumAnomalies(trackID)
		if err != nil {
			log.Printf("[spectrum-anomaly] rerun %s: %v", trackID, err)
			return
		}
		if err := storeSpectrumAnomalySummary(summary); err != nil {
			log.Printf("[spectrum-anomaly] rerun store %s: %v", trackID, err)
			return
		}
		log.Printf("[spectrum-anomaly] rerun %s: %d spectra, %d anomalies", trackID, summary.SpectrumCount, summary.AnomalyCount)
	}()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "queued", "track_id": trackID}) //nolint:errcheck
}
