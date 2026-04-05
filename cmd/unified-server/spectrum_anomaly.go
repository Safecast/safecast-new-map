// spectrum_anomaly.go — post-upload anomaly detection for spectrum tracks.
//
// After a spectrum track upload completes, triggerSpectrumAnomalyAnalysis
// launches a background goroutine that:
//  1. Loads all spectra for the track from PostgreSQL (joined via markers).
//  2. Normalises each spectrum to counts/second (divides by live_time_sec).
//  3. Computes the median normalised spectrum (median count rate per channel).
//  4. Computes cosine similarity of each spectrum against the median.
//  5. Flags spectra whose similarity falls below the configured threshold.
//  6. Stores the summary in the DuckLake spectrum_anomaly_summaries table.
//
// The summary is exposed via GET /api/track/{id}/insights (track_insights.go)
// and the /admin/anomaly admin page.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"time"
)

// SpectrumAnomalySummary holds the spectral anomaly analysis for one track.
type SpectrumAnomalySummary struct {
	TrackID            string                `json:"track_id"`
	ComputedAt         time.Time             `json:"computed_at"`
	SpectrumCount      int                   `json:"spectrum_count"`
	AnomalyCount       int                   `json:"anomaly_count"`
	AnomalySpectrumIDs []int64               `json:"anomaly_spectrum_ids,omitempty"`
	MedianTotalCounts  float64               `json:"median_total_counts"`
	CosineThreshold    float64               `json:"cosine_threshold"`
	AnomalyDetails     []spectrumAnomalyItem `json:"anomaly_details,omitempty"`
}

// spectrumAnomalyItem describes one anomalous spectrum.
type spectrumAnomalyItem struct {
	SpectrumID      int64   `json:"spectrum_id"`
	MarkerID        int64   `json:"marker_id"`
	CosineSimilarity float64 `json:"cosine_similarity"`
	TotalCountsPerSec float64 `json:"total_counts_per_sec"`
}

// triggerSpectrumAnomalyAnalysis launches background spectrum analysis.
// Safe to call from the upload goroutine — returns immediately.
func triggerSpectrumAnomalyAnalysis(trackID string) {
	if trackID == "" {
		return
	}
	go func() {
		time.Sleep(2 * time.Second)
		summary, err := computeSpectrumAnomalies(trackID)
		if err != nil {
			log.Printf("[spectrum-anomaly] %s: compute failed: %v", trackID, err)
			return
		}
		if err := storeSpectrumAnomalySummary(summary); err != nil {
			log.Printf("[spectrum-anomaly] %s: store failed: %v", trackID, err)
			return
		}
		log.Printf("[spectrum-anomaly] %s: %d spectra, %d anomalies (cosine threshold %.2f)",
			trackID, summary.SpectrumCount, summary.AnomalyCount, summary.CosineThreshold)
	}()
}

// spectrumRow holds raw data loaded from PostgreSQL for one spectrum.
type spectrumRow struct {
	id          int64
	markerID    int64
	channelsJSON string
	liveTimeSec float64
}

// computeSpectrumAnomalies performs the full analysis for a track.
func computeSpectrumAnomalies(trackID string) (*SpectrumAnomalySummary, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}
	cfg := getAnomalyConfig()

	rows, err := db.DB.Query(`
		SELECT s.id, s.marker_id, s.channels, COALESCE(s.live_time_sec, 0)
		FROM spectra s
		JOIN markers m ON s.marker_id = m.id
		WHERE m.trackid = $1
		  AND s.channels IS NOT NULL
		  AND s.channels != ''
		  AND s.channels != '[]'
		ORDER BY s.id`,
		trackID,
	)
	if err != nil {
		return nil, fmt.Errorf("query spectra: %w", err)
	}
	defer rows.Close()

	var spectra []spectrumRow
	for rows.Next() {
		var sr spectrumRow
		if err := rows.Scan(&sr.id, &sr.markerID, &sr.channelsJSON, &sr.liveTimeSec); err != nil {
			continue
		}
		if sr.liveTimeSec < cfg.SpectrumMinLiveTime {
			continue // too short — unreliable
		}
		spectra = append(spectra, sr)
	}

	summary := &SpectrumAnomalySummary{
		TrackID:         trackID,
		ComputedAt:      time.Now().UTC(),
		CosineThreshold: cfg.SpectrumCosineThresh,
		SpectrumCount:   len(spectra),
	}

	if len(spectra) < 2 {
		// Need at least 2 spectra to compare.
		return summary, nil
	}

	// Parse and normalise spectra to counts/second.
	type normSpectrum struct {
		id          int64
		markerID    int64
		normalised  []float64 // counts/sec per channel
		totalCPS    float64
	}

	var normalised []normSpectrum
	for _, sr := range spectra {
		var channels []int
		if err := json.Unmarshal([]byte(sr.channelsJSON), &channels); err != nil {
			continue
		}
		if len(channels) == 0 {
			continue
		}
		ns := normSpectrum{id: sr.id, markerID: sr.markerID}
		ns.normalised = make([]float64, len(channels))
		for i, c := range channels {
			cps := float64(c) / sr.liveTimeSec
			ns.normalised[i] = cps
			ns.totalCPS += cps
		}
		normalised = append(normalised, ns)
	}
	summary.SpectrumCount = len(normalised)

	if len(normalised) < 2 {
		return summary, nil
	}

	// Compute median total CPS across all spectra.
	totalCPSValues := make([]float64, len(normalised))
	for i, ns := range normalised {
		totalCPSValues[i] = ns.totalCPS
	}
	sort.Float64s(totalCPSValues)
	n := len(totalCPSValues)
	if n%2 == 0 {
		summary.MedianTotalCounts = (totalCPSValues[n/2-1] + totalCPSValues[n/2]) / 2
	} else {
		summary.MedianTotalCounts = totalCPSValues[n/2]
	}

	// Compute the median spectrum (median CPS per channel).
	// All spectra must have the same channel count for this to be meaningful;
	// use the most common channel count and skip outliers.
	channelCounts := make(map[int]int)
	for _, ns := range normalised {
		channelCounts[len(ns.normalised)]++
	}
	dominantLen := 0
	dominantCount := 0
	for l, c := range channelCounts {
		if c > dominantCount {
			dominantLen, dominantCount = l, c
		}
	}

	// Filter to spectra with the dominant channel count.
	var filtered []normSpectrum
	for _, ns := range normalised {
		if len(ns.normalised) == dominantLen {
			filtered = append(filtered, ns)
		}
	}
	if len(filtered) < 2 {
		return summary, nil
	}
	summary.SpectrumCount = len(filtered)

	// Build median spectrum channel-by-channel.
	medianSpectrum := make([]float64, dominantLen)
	channelVals := make([]float64, len(filtered))
	for ch := 0; ch < dominantLen; ch++ {
		for i, ns := range filtered {
			channelVals[i] = ns.normalised[ch]
		}
		sort.Float64s(channelVals)
		m := len(channelVals)
		if m%2 == 0 {
			medianSpectrum[ch] = (channelVals[m/2-1] + channelVals[m/2]) / 2
		} else {
			medianSpectrum[ch] = channelVals[m/2]
		}
	}

	// Cosine similarity of each spectrum against the median.
	medianNorm := vectorNorm(medianSpectrum)
	if medianNorm == 0 {
		return summary, nil // degenerate median — all zeros
	}

	var anomalyIDs []int64
	var anomalyDetails []spectrumAnomalyItem

	for _, ns := range filtered {
		sim := cosineSimilarityVec(ns.normalised, medianSpectrum, medianNorm)
		if sim < cfg.SpectrumCosineThresh {
			anomalyIDs = append(anomalyIDs, ns.id)
			anomalyDetails = append(anomalyDetails, spectrumAnomalyItem{
				SpectrumID:        ns.id,
				MarkerID:          ns.markerID,
				CosineSimilarity:  sim,
				TotalCountsPerSec: ns.totalCPS,
			})
		}
	}

	// Cap stored IDs per config.
	if len(anomalyIDs) > cfg.MaxAnomalyIDs {
		anomalyIDs = anomalyIDs[:cfg.MaxAnomalyIDs]
		anomalyDetails = anomalyDetails[:cfg.MaxAnomalyIDs]
	}

	summary.AnomalyCount = len(anomalyIDs)
	summary.AnomalySpectrumIDs = anomalyIDs
	summary.AnomalyDetails = anomalyDetails
	return summary, nil
}

// vectorNorm returns the L2 norm of a float64 slice.
func vectorNorm(v []float64) float64 {
	var sum float64
	for _, x := range v {
		sum += x * x
	}
	return math.Sqrt(sum)
}

// cosineSimilarityVec computes cosine similarity between a and b.
// bNorm is pre-computed to avoid recalculating for every spectrum.
func cosineSimilarityVec(a, b []float64, bNorm float64) float64 {
	if len(a) != len(b) || bNorm == 0 {
		return 0
	}
	var dot float64
	for i := range a {
		dot += a[i] * b[i]
	}
	aNorm := vectorNorm(a)
	if aNorm == 0 {
		return 0
	}
	return dot / (aNorm * bNorm)
}

// storeSpectrumAnomalySummary persists the summary using DELETE + INSERT.
func storeSpectrumAnomalySummary(s *SpectrumAnomalySummary) error {
	if !duckDBAvailable() {
		return fmt.Errorf("duckdb not available")
	}

	idsJSON, _ := json.Marshal(s.AnomalySpectrumIDs)
	detailsJSON, _ := json.Marshal(s.AnomalyDetails)

	if _, err := duckDB.Exec(
		`DELETE FROM spectrum_anomaly_summaries WHERE track_id = ?`, s.TrackID,
	); err != nil {
		return fmt.Errorf("delete old summary: %w", err)
	}

	if _, err := duckDB.Exec(
		`INSERT INTO spectrum_anomaly_summaries
		 (track_id, computed_at, spectrum_count, anomaly_count,
		  anomaly_spectrum_ids, median_total_counts, cosine_threshold, anomaly_details)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		s.TrackID, s.ComputedAt, s.SpectrumCount, s.AnomalyCount,
		string(idsJSON), s.MedianTotalCounts, s.CosineThreshold, string(detailsJSON),
	); err != nil {
		return fmt.Errorf("insert summary: %w", err)
	}
	return nil
}

// loadSpectrumAnomalySummary fetches the stored summary for a track.
// Returns nil, nil if not yet computed.
func loadSpectrumAnomalySummary(trackID string) (*SpectrumAnomalySummary, error) {
	if !duckDBAvailable() {
		return nil, nil
	}
	row := duckDB.QueryRow(`
		SELECT track_id, computed_at, spectrum_count, anomaly_count,
		       anomaly_spectrum_ids, median_total_counts, cosine_threshold, anomaly_details
		FROM spectrum_anomaly_summaries
		WHERE track_id = ?
		LIMIT 1`,
		trackID,
	)
	var s SpectrumAnomalySummary
	var idsJSON, detailsJSON string
	if err := row.Scan(
		&s.TrackID, &s.ComputedAt, &s.SpectrumCount, &s.AnomalyCount,
		&idsJSON, &s.MedianTotalCounts, &s.CosineThreshold, &detailsJSON,
	); err != nil {
		return nil, nil
	}
	_ = json.Unmarshal([]byte(idsJSON), &s.AnomalySpectrumIDs)
	_ = json.Unmarshal([]byte(detailsJSON), &s.AnomalyDetails)
	return &s, nil
}
