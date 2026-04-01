package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

var getTrackToolDef = mcp.NewTool("get_track",
	mcp.WithDescription("Retrieve radiation measurements for a specific track by its track ID. The track_id can be used directly — you do NOT need to call list_tracks first if you already know the ID (e.g. from a URL like simplemap.safecast.org/trackid/8fVdyw or from a prior query). Use stats_only=true to get a statistical summary (min/max/mean/percentiles/count) without returning all individual measurements — this is strongly preferred for large tracks (>500 measurements) or when the user asks about peak values, distribution, or anomalies. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool. When referencing or linking to track data, ALWAYS use https://simplemap.safecast.org as the base URL — NEVER use api.safecast.org, which does not host track data. CRITICAL: Present all findings in an objective, scientific manner without using personal pronouns (I, we, I'll, you) or conversational language (Perfect!, Great!). Format as factual statements only."),
	mcp.WithString("track_id",
		mcp.Description("Track identifier (bGeigie import ID or track ID). Can be used directly without calling list_tracks first."),
		mcp.Required(),
	),
	mcp.WithBoolean("stats_only",
		mcp.Description("If true, return only statistical aggregates (count, min, max, mean, median, p95, p99, stddev) plus the top-10 peak measurements. No individual measurement list is returned. Use this for large tracks or when the user wants a summary, peak values, or anomaly detection."),
		mcp.DefaultBool(false),
	),
	mcp.WithNumber("from",
		mcp.Description("Optional: Start marker ID for filtering"),
	),
	mcp.WithNumber("to",
		mcp.Description("Optional: End marker ID for filtering"),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of measurements to return (default: 200, max: 10000). Ignored when stats_only=true."),
		mcp.Min(1), mcp.Max(10000),
		mcp.DefaultNumber(200),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleGetTrack(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	trackIDStr, err := req.RequireString("track_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	statsOnly := req.GetBool("stats_only", false)
	limit := req.GetInt("limit", 200)
	if limit < 1 || limit > 10000 {
		return mcp.NewToolResultError("Limit must be between 1 and 10000"), nil
	}

	fromID := req.GetInt("from", 0)
	toID := req.GetInt("to", 0)

	if statsOnly {
		if dbAvailable() {
			return getTrackStatsDB(ctx, trackIDStr, fromID, toID)
		}
		// API fallback: fetch all and compute stats client-side
		return getTrackStatsAPI(ctx, trackIDStr, fromID, toID)
	}

	if dbAvailable() {
		return getTrackDB(ctx, trackIDStr, fromID, toID, limit)
	}
	return getTrackAPI(ctx, trackIDStr, fromID, toID, limit)
}

func getTrackStatsDB(ctx context.Context, trackID string, fromID, toID int) (*mcp.CallToolResult, error) {
	// Build optional WHERE clauses for from/to marker filtering
	filterClauses := ""
	args := []any{trackID}
	argIdx := 2
	if fromID != 0 {
		filterClauses += fmt.Sprintf(" AND id >= $%d", argIdx)
		args = append(args, fromID)
		argIdx++
	}
	if toID != 0 {
		filterClauses += fmt.Sprintf(" AND id <= $%d", argIdx)
		args = append(args, toID)
		argIdx++
	}

	statsQuery := fmt.Sprintf(`
		SELECT
			count(*)                                        AS total,
			min(doserate)                                   AS min_value,
			max(doserate)                                   AS max_value,
			avg(doserate)                                   AS mean_value,
			percentile_cont(0.50) WITHIN GROUP (ORDER BY doserate) AS median,
			percentile_cont(0.95) WITHIN GROUP (ORDER BY doserate) AS p95,
			percentile_cont(0.99) WITHIN GROUP (ORDER BY doserate) AS p99,
			stddev(doserate)                                AS stddev,
			min(lat) AS min_lat, max(lat) AS max_lat,
			min(lon) AS min_lon, max(lon) AS max_lon,
			to_timestamp(min(date)) AS first_captured_at,
			to_timestamp(max(date)) AS last_captured_at
		FROM markers
		WHERE trackid = $1%s`, filterClauses)

	statsRow, err := queryRow(ctx, statsQuery, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Top-10 peak measurements
	peakArgs := []any{trackID}
	peakArgIdx := 2
	peakFilter := ""
	if fromID != 0 {
		peakFilter += fmt.Sprintf(" AND id >= $%d", peakArgIdx)
		peakArgs = append(peakArgs, fromID)
		peakArgIdx++
	}
	if toID != 0 {
		peakFilter += fmt.Sprintf(" AND id <= $%d", peakArgIdx)
		peakArgs = append(peakArgs, toID)
	}
	peakQuery := fmt.Sprintf(`
		SELECT id, doserate AS value, 'µSv/h' AS unit,
			to_timestamp(date) AS captured_at,
			lat AS latitude, lon AS longitude
		FROM markers
		WHERE trackid = $1%s
		ORDER BY doserate DESC
		LIMIT 10`, peakFilter)

	peakRows, err := queryRows(ctx, peakQuery, peakArgs...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	peaks := make([]map[string]any, len(peakRows))
	for i, r := range peakRows {
		peaks[i] = map[string]any{
			"id":          r["id"],
			"value":       r["value"],
			"unit":        r["unit"],
			"captured_at": r["captured_at"],
			"location": map[string]any{
				"latitude":  r["latitude"],
				"longitude": r["longitude"],
			},
		}
	}

	result := map[string]any{
		"track_id":    trackID,
		"map_url":     "https://simplemap.safecast.org/trackid/" + trackID,
		"source":      "database",
		"stats_only":  true,
		"from_marker": nilIfZero(fromID),
		"to_marker":   nilIfZero(toID),
		"statistics":  statsRow,
		"top10_peaks": peaks,
		"_ai_hint": "STATS INTERPRETATION (all values in µSv/h): " +
			"min_value = lowest single reading on the track. " +
			"max_value = highest single reading — use this to answer 'what was the peak?' questions. " +
			"mean_value = arithmetic average. median = middle value (robust to outliers). " +
			"p95/p99 = 95th/99th percentile — values this high or higher occurred in only 5%/1% of readings. " +
			"stddev = spread; high stddev relative to mean indicates localised hotspots. " +
			"top10_peaks = the 10 highest individual readings with coordinates — link each location as " +
			"https://simplemap.safecast.org/?lat=LAT&lon=LON&zoom=16 so the user can inspect it on the map. " +
			"If more detail is needed (e.g. all readings above a threshold), call get_track again without stats_only and use limit up to 10000. " +
			"Present findings as objective factual statements only — no personal pronouns, no exclamations.",
		"_ai_generated_note": aiGeneratedNote,
	}
	return jsonResult(result)
}

func getTrackStatsAPI(ctx context.Context, trackIDStr string, fromID, toID int) (*mcp.CallToolResult, error) {
	resp, err := client.GetTrackData(ctx, trackIDStr, fromID, toID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	markers, _ := resp["markers"].([]any)

	type statAccum struct {
		count                      int
		sum, min, max              float64
		vals                       []float64
	}
	acc := statAccum{min: 1e18, max: -1e18}

	type peak struct {
		value float64
		raw   map[string]any
	}
	var peaks []peak

	for _, raw := range markers {
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		nm := normalizeLatestMarker(m)
		v, ok := nm["value"].(float64)
		if !ok {
			continue
		}
		acc.count++
		acc.sum += v
		acc.vals = append(acc.vals, v)
		if v < acc.min {
			acc.min = v
		}
		if v > acc.max {
			acc.max = v
		}
		peaks = append(peaks, peak{value: v, raw: nm})
	}

	var mean, stddev, median, p95, p99 float64
	if acc.count > 0 {
		mean = acc.sum / float64(acc.count)
		// sort for percentiles
		sorted := make([]float64, len(acc.vals))
		copy(sorted, acc.vals)
		for i := 1; i < len(sorted); i++ {
			for j := i; j > 0 && sorted[j] < sorted[j-1]; j-- {
				sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
			}
		}
		median = percentileVal(sorted, 0.50)
		p95 = percentileVal(sorted, 0.95)
		p99 = percentileVal(sorted, 0.99)
		for _, v := range acc.vals {
			d := v - mean
			stddev += d * d
		}
		if acc.count > 1 {
			stddev = stddev / float64(acc.count-1)
		}
		// sqrt approximation
		stddev = stddevSqrt(stddev)
	}

	// sort peaks descending and take top 10
	for i := 1; i < len(peaks); i++ {
		for j := i; j > 0 && peaks[j].value > peaks[j-1].value; j-- {
			peaks[j], peaks[j-1] = peaks[j-1], peaks[j]
		}
	}
	if len(peaks) > 10 {
		peaks = peaks[:10]
	}
	top10 := make([]map[string]any, len(peaks))
	for i, p := range peaks {
		top10[i] = p.raw
	}

	result := map[string]any{
		"track_id":   trackIDStr,
		"map_url":    "https://simplemap.safecast.org/trackid/" + trackIDStr,
		"source":     "api",
		"stats_only": true,
		"statistics": map[string]any{
			"total":      acc.count,
			"min_value":  acc.min,
			"max_value":  acc.max,
			"mean_value": mean,
			"median":     median,
			"p95":        p95,
			"p99":        p99,
			"stddev":     stddev,
		},
		"top10_peaks":        top10,
		"_ai_hint": "STATS INTERPRETATION (all values in µSv/h): " +
			"min_value = lowest reading. max_value = peak reading on the entire track. " +
			"mean_value = average. median = middle value (robust to outliers). " +
			"p95/p99 = 95th/99th percentile — values this high or higher in only 5%/1% of readings. " +
			"top10_peaks = 10 highest individual readings with coordinates — link each as " +
			"https://simplemap.safecast.org/?lat=LAT&lon=LON&zoom=16 for map inspection. " +
			"For all individual readings (up to 10000), call get_track again without stats_only. " +
			"Present as objective factual statements — no personal pronouns, no exclamations.",
		"_ai_generated_note": aiGeneratedNote,
	}
	return jsonResult(result)
}

// percentileVal returns the value at the given quantile (0–1) from a sorted slice.
func percentileVal(sorted []float64, q float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := q * float64(len(sorted)-1)
	lo := int(idx)
	hi := lo + 1
	if hi >= len(sorted) {
		return sorted[lo]
	}
	frac := idx - float64(lo)
	return sorted[lo]*(1-frac) + sorted[hi]*frac
}

// stddevSqrt is a simple Newton's method sqrt to avoid importing math.
func stddevSqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 20; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

func getTrackDB(ctx context.Context, trackID string, fromID, toID, limit int) (*mcp.CallToolResult, error) {
	query := `
		SELECT m.id, m.doserate AS value, 'µSv/h' AS unit,
			to_timestamp(m.date) AS captured_at,
			m.lat AS latitude, m.lon AS longitude,
			m.device_id, m.altitude AS height, m.detector,
			m.has_spectrum,
			u.internal_user_id, usr.username AS uploader_username, usr.email AS uploader_email
		FROM markers m
		LEFT JOIN uploads u ON u.track_id = m.trackid
		LEFT JOIN users usr ON u.internal_user_id = usr.id::text
		WHERE m.trackid = $1`

	args := []any{trackID}
	argIdx := 2

	if fromID != 0 {
		query += fmt.Sprintf(" AND id >= $%d", argIdx)
		args = append(args, fromID)
		argIdx++
	}
	if toID != 0 {
		query += fmt.Sprintf(" AND id <= $%d", argIdx)
		args = append(args, toID)
		argIdx++
	}

	query += " ORDER BY date ASC"
	query += fmt.Sprintf(" LIMIT $%d", argIdx)
	args = append(args, limit)

	rows, err := queryRows(ctx, query, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get total count for this track
	countRow, _ := queryRow(ctx, `SELECT count(*) AS total FROM markers WHERE trackid = $1`, trackID)
	total := 0
	if countRow != nil {
		if t, ok := countRow["total"]; ok {
			switch v := t.(type) {
			case int64:
				total = int(v)
			case float64:
				total = int(v)
			}
		}
	}

	measurements := make([]map[string]any, len(rows))
	var uploaderUsername, uploaderEmail any
	for i, r := range rows {
		measurements[i] = map[string]any{
			"id":    r["id"],
			"value": r["value"],
			"unit":  r["unit"],
			"captured_at": r["captured_at"],
			"location": map[string]any{
				"latitude":  r["latitude"],
				"longitude": r["longitude"],
			},
			"device_id":   r["device_id"],
			"height":      r["height"],
			"detector":    r["detector"],
			"has_spectrum": r["has_spectrum"],
		}

		// Store uploader info from first row (all rows for same track have same uploader)
		if i == 0 {
			uploaderUsername = r["uploader_username"]
			uploaderEmail = r["uploader_email"]
		}
	}

	result := map[string]any{
		"track_id":        trackID,
		"map_url":         "https://simplemap.safecast.org/trackid/" + trackID,
		"count":           len(measurements),
		"total_available": total,
		"source":          "database",
		"from_marker":     nilIfZero(fromID),
		"to_marker":       nilIfZero(toID),
		"measurements":    measurements,
		"_ai_hint": "INSTRUCTIONS: (1) All doserate values are in µSv/h. " +
			"(2) 'count' is the number of measurements returned; 'total_available' is the full track size. " +
			"If total_available > count, only a partial view was returned — call get_track with stats_only=true " +
			"to get min/max/mean/percentiles/top-10 peaks for the ENTIRE track without hitting the limit. " +
			"To retrieve all individual readings use limit=10000. " +
			"(3) Link every location as https://simplemap.safecast.org/?lat=LAT&lon=LON&zoom=16. " +
			"(4) Present as objective factual statements — no personal pronouns, no exclamations.",
		"_ai_generated_note": aiGeneratedNote,
	}

	// Add uploader information if available
	if uploaderUsername != nil && uploaderUsername != "" {
		result["uploader"] = map[string]any{
			"username": uploaderUsername,
			"email":    uploaderEmail,
		}
	}

	return jsonResult(result)
}

func getTrackAPI(ctx context.Context, trackIDStr string, fromID, toID, limit int) (*mcp.CallToolResult, error) {
	resp, err := client.GetTrackData(ctx, trackIDStr, fromID, toID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	markers, _ := resp["markers"].([]any)
	totalAvailable := len(markers)

	if limit > len(markers) {
		limit = len(markers)
	}
	limited := markers[:limit]

	normalized := make([]map[string]any, 0, len(limited))
	for _, raw := range limited {
		if m, ok := raw.(map[string]any); ok {
			normalized = append(normalized, normalizeLatestMarker(m))
		}
	}

	result := map[string]any{
		"track": map[string]any{
			"track_id":     resp["trackID"],
			"marker_count": resp["markerCount"],
			"track_index":  resp["trackIndex"],
			"map_url":      "https://simplemap.safecast.org/trackid/" + resp["trackID"].(string),
		},
		"count":           len(normalized),
		"total_available": totalAvailable,
		"source":          "api",
		"from_marker":     nilIfZero(fromID),
		"to_marker":       nilIfZero(toID),
		"measurements":    normalized,
		"_ai_hint": "INSTRUCTIONS: (1) All doserate values are in µSv/h. " +
			"(2) 'count' is measurements returned; 'total_available' is the full track size. " +
			"If total_available > count, the track is larger than the limit — call get_track with " +
			"stats_only=true to get min/max/mean/percentiles/top-10 peaks for the entire track, " +
			"or use limit=10000 to retrieve up to 10000 individual readings. " +
			"(3) Link every location as https://simplemap.safecast.org/?lat=LAT&lon=LON&zoom=16. " +
			"(4) Present as objective factual statements — no personal pronouns, no exclamations.",
		"_ai_generated_note": aiGeneratedNote,
	}

	return jsonResult(result)
}
