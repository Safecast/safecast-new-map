package main

import (
	"context"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

var deviceHistoryToolDef = mcp.NewTool("device_history",
	mcp.WithDescription("Get historical measurements from MOBILE bGeigie survey devices (type=geigiecast, IDs like geigiecast:62007). Use this tool ONLY for mobile bGeigie devices. DO NOT use for fixed sensors — for bGeigieZen (geigiecast-zen), Pointcast, Solarcast, Notehub/Radnote (note:dev:...), nGeigie, or device-tcp, use sensor_current instead. Radiation values are in CPM (counts per minute, NOT counts per second). Always present radiation values in µSv/h by converting from CPM using detector-specific factors. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool. CRITICAL: Present all findings in an objective, scientific manner without using personal pronouns (I, we, I'll, you) or conversational language (Perfect!, Great!). Format as factual statements only."),
	mcp.WithString("device_id",
		mcp.Description("Device identifier"),
		mcp.Required(),
	),
	mcp.WithNumber("days",
		mcp.Description("Number of days of history to retrieve (default: 30, max: 365)"),
		mcp.Min(1), mcp.Max(365),
		mcp.DefaultNumber(30),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of measurements to return (default: 200, max: 10000)"),
		mcp.Min(1), mcp.Max(10000),
		mcp.DefaultNumber(200),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleDeviceHistory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	deviceIDStr, err := req.RequireString("device_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	days := req.GetInt("days", 30)
	limit := req.GetInt("limit", 200)

	if days < 1 || days > 365 {
		return mcp.NewToolResultError("days must be between 1 and 365"), nil
	}
	if limit < 1 || limit > 10000 {
		return mcp.NewToolResultError("Limit must be between 1 and 10000"), nil
	}

	if dbAvailable() {
		return deviceHistoryDB(ctx, deviceIDStr, days, limit)
	}
	return deviceHistoryAPI(ctx, deviceIDStr, days, limit)
}

func deviceHistoryDB(ctx context.Context, deviceID string, days, limit int) (*mcp.CallToolResult, error) {
	now := time.Now().UTC()
	startDate := now.AddDate(0, 0, -days)

	// Query both markers table and realtime_measurements table
	// First, try markers table (bGeigie imports)
	markersQuery := `
		SELECT m.id, m.doserate AS value, 'µSv/h' AS unit,
			to_timestamp(m.date) AS captured_at,
			m.lat AS latitude, m.lon AS longitude,
			m.altitude AS height, m.detector, m.trackid::text AS track_id,
			u.internal_user_id, usr.username AS uploader_username, usr.email AS uploader_email
		FROM markers m
		LEFT JOIN uploads u ON u.track_id = m.trackid
		LEFT JOIN users usr ON u.internal_user_id = usr.id::text
		WHERE m.device_id = $1 AND m.date >= $2 AND m.date <= $3
		ORDER BY m.date DESC
		LIMIT $4`

	markersRows, err := queryRows(ctx, markersQuery, deviceID, startDate.Unix(), now.Unix(), limit)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Then, query realtime_measurements table (fixed sensors)
	// First, check if the table exists and what columns it has
	columnsQuery := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = 'realtime_measurements'
		ORDER BY column_name
	`
	
	var realtimeRows []map[string]any
	
	columnRows, err := queryRows(ctx, columnsQuery)
	if err != nil || len(columnRows) == 0 {
		// If we can't query the schema or table doesn't exist, try the basic query
		realtimeQuery := `
			SELECT id, value, unit,
				to_timestamp(measured_at) AS captured_at,
				lat AS latitude, lon AS longitude,
				device_name, transport, device_id
			FROM realtime_measurements
			WHERE device_id = $1 AND measured_at >= $2 AND measured_at <= $3
			ORDER BY measured_at DESC
			LIMIT $4`

		realtimeRows, err = queryRows(ctx, realtimeQuery, deviceID, startDate.Unix(), now.Unix(), limit)
		if err != nil {
			return mcp.NewToolResultError("Error querying realtime_measurements table: " + err.Error()), nil
		}
	} else {
		// Build the query based on available columns
		hasHeight := false
		for _, row := range columnRows {
			if colName, ok := row["column_name"].(string); ok && colName == "height" {
				hasHeight = true
				break
			}
		}
		
		var realtimeQuery string
		if hasHeight {
			realtimeQuery = `
				SELECT id, value, unit,
					to_timestamp(measured_at) AS captured_at,
					lat AS latitude, lon AS longitude,
					device_name, transport, device_id, height
				FROM realtime_measurements
				WHERE device_id = $1 AND measured_at >= $2 AND measured_at <= $3
				ORDER BY measured_at DESC
				LIMIT $4`
		} else {
			realtimeQuery = `
				SELECT id, value, unit,
					to_timestamp(measured_at) AS captured_at,
					lat AS latitude, lon AS longitude,
					device_name, transport, device_id
				FROM realtime_measurements
				WHERE device_id = $1 AND measured_at >= $2 AND measured_at <= $3
				ORDER BY measured_at DESC
				LIMIT $4`
		}

		realtimeRows, err = queryRows(ctx, realtimeQuery, deviceID, startDate.Unix(), now.Unix(), limit)
		if err != nil {
			return mcp.NewToolResultError("Error querying realtime_measurements table: " + err.Error()), nil
		}
	}

	// Combine results and sort by timestamp (most recent first)
	allMeasurements := make([]map[string]any, 0, len(markersRows)+len(realtimeRows))
	
	// Process markers results
	for _, r := range markersRows {
		measurement := map[string]any{
			"id":    r["id"],
			"value": r["value"],
			"unit":  r["unit"],
			"captured_at": r["captured_at"],
			"location": map[string]any{
				"latitude":  r["latitude"],
				"longitude": r["longitude"],
			},
			"height":   r["height"],
			"detector": r["detector"],
			"source":   "bgeigie_import",
		}
		if r["track_id"] != nil {
			measurement["track_id"] = r["track_id"]
		}

		// Add uploader information if available
		if uploaderUsername, ok := r["uploader_username"]; ok && uploaderUsername != nil && uploaderUsername != "" {
			measurement["uploader"] = map[string]any{
				"username": uploaderUsername,
				"email":    r["uploader_email"],
			}
		}

		allMeasurements = append(allMeasurements, measurement)
	}
	
	// Process realtime results
	for _, r := range realtimeRows {
		// Fix incorrect unit: Geiger counters report in CPM (counts per minute), not CPS
		unit := r["unit"]
		if unitStr, ok := unit.(string); ok {
			unit = strings.ReplaceAll(strings.ReplaceAll(unitStr, "cps", "cpm"), "CPS", "CPM")
		}

		measurement := map[string]any{
			"id":    r["id"],
			"value": r["value"],
			"unit":  unit,
			"captured_at": r["captured_at"],
			"location": map[string]any{
				"latitude":  r["latitude"],
				"longitude": r["longitude"],
			},
			"height":   r["height"],
			"device_name": r["device_name"],
			"type":     r["transport"],
			"source":   "realtime_sensor",
		}
		allMeasurements = append(allMeasurements, measurement)
	}

	// Sort all measurements by captured_at timestamp (most recent first)
	// Since queryRows returns timestamps as interface{}, we need to handle both string and time.Time types
	for i := 0; i < len(allMeasurements)-1; i++ {
		for j := 0; j < len(allMeasurements)-i-1; j++ {
			// Get the timestamp values - they could be strings or time.Time objects
			time1Val := allMeasurements[j]["captured_at"]
			time2Val := allMeasurements[j+1]["captured_at"]
			
			// Compare timestamps - swap if j-th element is older than (j+1)-th element
			shouldSwap := false
			
			// Handle different possible types for timestamps
			switch v1 := time1Val.(type) {
			case string:
				if v2, ok := time2Val.(string); ok {
					if v1 < v2 {
						shouldSwap = true
					}
				}
			case time.Time:
				if v2, ok := time2Val.(time.Time); ok {
					if v1.Before(v2) {
						shouldSwap = true
					}
				}
			}
			
			if shouldSwap {
				allMeasurements[j], allMeasurements[j+1] = allMeasurements[j+1], allMeasurements[j]
			}
		}
	}

	// Apply limit
	measurements := allMeasurements
	if len(measurements) > limit {
		measurements = measurements[:limit]
	}

	capturedAfter := startDate.Format("2006-01-02") + " 00:00"
	capturedBefore := now.Format("2006-01-02") + " 23:59"

	result := map[string]any{
		"device": map[string]any{
			"id": deviceID,
		},
		"period": map[string]any{
			"days":       days,
			"start_date": capturedAfter,
			"end_date":   capturedBefore,
		},
		"count":        len(measurements),
		"source":       "database",
		"measurements": measurements,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The 'unit' field indicates measurement units - CPM means 'counts per minute' NOT 'counts per second'. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: 'Latest reading: X CPM at location Y' NOT 'I found a reading of X CPM' or 'Perfect! The sensor shows...'. State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}

func deviceHistoryAPI(ctx context.Context, deviceIDStr string, days, limit int) (*mcp.CallToolResult, error) {
	resp, err := client.GetRealtimeHistory(ctx, deviceIDStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	now := time.Now().UTC()
	startDate := now.AddDate(0, 0, -days)
	capturedAfter := startDate.Format("2006-01-02") + " 00:00"
	capturedBefore := now.Format("2006-01-02") + " 23:59"
	startUnix := startDate.Unix()

	// Extract dose rate time series and filter by date range
	var measurements []map[string]any
	if series, ok := resp["series"].(map[string]any); ok {
		if doseRate, ok := series["doseRate"].([]any); ok {
			for _, raw := range doseRate {
				pt, ok := raw.(map[string]any)
				if !ok {
					continue
				}
				ts, ok := toFloat(pt["time"])
				if !ok || int64(ts) < startUnix {
					continue
				}
				t := time.Unix(int64(ts), 0).UTC()
				measurements = append(measurements, map[string]any{
					"value":       pt["value"],
					"unit":        "µSv/h",
					"captured_at": t.Format(time.RFC3339),
				})
			}
		}
	}

	totalAvailable := len(measurements)
	if limit > len(measurements) {
		limit = len(measurements)
	}
	measurements = measurements[:limit]

	deviceInfo := map[string]any{
		"id": deviceIDStr,
	}
	if name, ok := resp["deviceName"].(string); ok && name != "" {
		deviceInfo["name"] = name
	}
	if tube, ok := resp["tube"].(string); ok && tube != "" {
		deviceInfo["sensor"] = tube
	}

	result := map[string]any{
		"device": deviceInfo,
		"period": map[string]any{
			"days":       days,
			"start_date": capturedAfter,
			"end_date":   capturedBefore,
		},
		"count":           len(measurements),
		"total_available": totalAvailable,
		"source":          "api",
		"measurements":    measurements,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The 'unit' field indicates measurement units - CPM means 'counts per minute' NOT 'counts per second'. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: 'Latest reading: X CPM at location Y' NOT 'I found a reading of X CPM' or 'Perfect! The sensor shows...'. State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	if ranges, ok := resp["ranges"].(map[string]any); ok {
		if dr, ok := ranges["doseRate"].(map[string]any); ok {
			result["statistics"] = map[string]any{
				"min_usvh": dr["min"],
				"max_usvh": dr["max"],
				"avg_usvh": dr["avg"],
			}
		}
	}

	return jsonResult(result)
}
