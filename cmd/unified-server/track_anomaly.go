// track_anomaly.go — post-upload anomaly detection for radiation tracks.
//
// After a track upload completes, triggerAnomalyAnalysis launches a background
// goroutine that:
//  1. Loads all doserate values for the track from PostgreSQL.
//  2. Computes mean, median, and stddev (Welford single-pass).
//  3. Flags readings where doserate > mean + 3*stddev.
//  4. Stores the summary in the DuckLake track_anomaly_summaries table.
//
// The summary is exposed via GET /api/track/{id}/insights (track_insights.go).

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"time"
)

// AnomalySummary holds the statistical anomaly analysis for a single track.
type AnomalySummary struct {
	TrackID        string    `json:"track_id"`
	ComputedAt     time.Time `json:"computed_at"`
	ReadingCount   int       `json:"reading_count"`
	MeanDoserate   float64   `json:"mean_doserate"`
	MedianDoserate float64   `json:"median_doserate"`
	StddevDoserate float64   `json:"stddev_doserate"`
	Threshold      float64   `json:"threshold"`
	AnomalyCount   int       `json:"anomaly_count"`
	AnomalyIDs     []int64   `json:"anomaly_ids,omitempty"`
	AnomalyMax     float64   `json:"anomaly_max,omitempty"`
	AnomalyMin     float64   `json:"anomaly_min,omitempty"`
}

// triggerAnomalyAnalysis launches a background goroutine to analyse the track.
// It waits 2 seconds to let the database commit settle before querying.
func triggerAnomalyAnalysis(trackID string) {
	if trackID == "" {
		return
	}
	go func() {
		time.Sleep(2 * time.Second)
		summary, err := computeAnomalies(trackID)
		if err != nil {
			log.Printf("[anomaly] %s: compute failed: %v", trackID, err)
			return
		}
		if err := storeAnomalySummary(summary); err != nil {
			log.Printf("[anomaly] %s: store failed: %v", trackID, err)
			return
		}
		log.Printf("[anomaly] %s: %d readings, %d anomalies (threshold %.4f µSv/h)",
			trackID, summary.ReadingCount, summary.AnomalyCount, summary.Threshold)
	}()
}

// computeAnomalies queries PostgreSQL for all doserate values on the track and
// computes the statistical summary.
func computeAnomalies(trackID string) (*AnomalySummary, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	cfg := getAnomalyConfig()
	rows, err := db.DB.Query(
		`SELECT id, doserate FROM markers
		 WHERE trackid = $1 AND doserate > $2 AND doserate < $3
		 ORDER BY id`,
		trackID, cfg.MinDoserate, cfg.MaxDoserate,
	)
	if err != nil {
		return nil, fmt.Errorf("query markers: %w", err)
	}
	defer rows.Close()

	var ids []int64
	var values []float64
	for rows.Next() {
		var id int64
		var v float64
		if err := rows.Scan(&id, &v); err != nil {
			continue
		}
		ids = append(ids, id)
		values = append(values, v)
	}

	summary := &AnomalySummary{
		TrackID:    trackID,
		ComputedAt: time.Now().UTC(),
	}

	n := len(values)
	summary.ReadingCount = n

	if n == 0 {
		return summary, nil
	}

	// Mean and variance via Welford's online algorithm.
	var mean, m2 float64
	for i, v := range values {
		delta := v - mean
		mean += delta / float64(i+1)
		m2 += delta * (v - mean)
	}
	stddev := 0.0
	if n > 1 {
		stddev = math.Sqrt(m2 / float64(n-1))
	}

	// Median via sorted copy.
	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)
	var median float64
	if n%2 == 0 {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	} else {
		median = sorted[n/2]
	}

	threshold := mean + cfg.SigmaThreshold*stddev

	summary.MeanDoserate = mean
	summary.MedianDoserate = median
	summary.StddevDoserate = stddev
	summary.Threshold = threshold

	// Flag anomalous readings.
	var anomalyIDs []int64
	anomalyMax := -math.MaxFloat64
	anomalyMin := math.MaxFloat64

	for i, v := range values {
		if v > threshold {
			anomalyIDs = append(anomalyIDs, ids[i])
			if v > anomalyMax {
				anomalyMax = v
			}
			if v < anomalyMin {
				anomalyMin = v
			}
		}
	}

	// Cap stored IDs per config to keep the JSON compact.
	if len(anomalyIDs) > cfg.MaxAnomalyIDs {
		anomalyIDs = anomalyIDs[:cfg.MaxAnomalyIDs]
	}

	summary.AnomalyCount = len(anomalyIDs)
	summary.AnomalyIDs = anomalyIDs
	if summary.AnomalyCount > 0 {
		summary.AnomalyMax = anomalyMax
		summary.AnomalyMin = anomalyMin
	}

	return summary, nil
}

// storeAnomalySummary persists the summary to DuckLake using DELETE + INSERT
// (DuckLake does not support ON CONFLICT / upsert).
func storeAnomalySummary(s *AnomalySummary) error {
	if !duckDBAvailable() {
		return fmt.Errorf("duckdb not available")
	}

	idsJSON, err := json.Marshal(s.AnomalyIDs)
	if err != nil {
		idsJSON = []byte("[]")
	}

	if _, err := duckDB.Exec(
		`DELETE FROM track_anomaly_summaries WHERE track_id = ?`, s.TrackID,
	); err != nil {
		return fmt.Errorf("delete old summary: %w", err)
	}

	if _, err := duckDB.Exec(
		`INSERT INTO track_anomaly_summaries
		 (track_id, computed_at, reading_count,
		  mean_doserate, median_doserate, stddev_doserate,
		  threshold, anomaly_count, anomaly_ids,
		  anomaly_max, anomaly_min)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.TrackID, s.ComputedAt, s.ReadingCount,
		s.MeanDoserate, s.MedianDoserate, s.StddevDoserate,
		s.Threshold, s.AnomalyCount, string(idsJSON),
		s.AnomalyMax, s.AnomalyMin,
	); err != nil {
		return fmt.Errorf("insert summary: %w", err)
	}

	return nil
}

// loadAnomalySummary fetches the stored anomaly summary for a track from DuckLake.
// Returns nil, nil if no summary has been computed yet.
func loadAnomalySummary(trackID string) (*AnomalySummary, error) {
	if !duckDBAvailable() {
		return nil, nil
	}

	row := duckDB.QueryRow(
		`SELECT track_id, computed_at, reading_count,
		        mean_doserate, median_doserate, stddev_doserate,
		        threshold, anomaly_count, anomaly_ids,
		        anomaly_max, anomaly_min
		 FROM track_anomaly_summaries
		 WHERE track_id = ?
		 LIMIT 1`,
		trackID,
	)

	var s AnomalySummary
	var idsJSON string
	if err := row.Scan(
		&s.TrackID, &s.ComputedAt, &s.ReadingCount,
		&s.MeanDoserate, &s.MedianDoserate, &s.StddevDoserate,
		&s.Threshold, &s.AnomalyCount, &idsJSON,
		&s.AnomalyMax, &s.AnomalyMin,
	); err != nil {
		// No row — not yet computed.
		return nil, nil
	}

	if idsJSON != "" && idsJSON != "null" {
		_ = json.Unmarshal([]byte(idsJSON), &s.AnomalyIDs)
	}

	return &s, nil
}
