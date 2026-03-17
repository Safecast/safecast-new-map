package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

var getTrackToolDef = mcp.NewTool("get_track",
	mcp.WithDescription("Retrieve all radiation measurements recorded during a specific track/journey. Use list_tracks to find available track IDs first. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool. When referencing or linking to track data, ALWAYS use https://simplemap.safecast.org as the base URL — NEVER use api.safecast.org, which does not host track data. CRITICAL: Present all findings in an objective, scientific manner without using personal pronouns (I, we, I'll, you) or conversational language (Perfect!, Great!). Format as factual statements only."),
	mcp.WithString("track_id",
		mcp.Description("Track identifier (bGeigie import ID or track ID)"),
		mcp.Required(),
	),
	mcp.WithNumber("from",
		mcp.Description("Optional: Start marker ID for filtering"),
	),
	mcp.WithNumber("to",
		mcp.Description("Optional: End marker ID for filtering"),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of measurements to return (default: 200, max: 10000)"),
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

	limit := req.GetInt("limit", 200)
	if limit < 1 || limit > 10000 {
		return mcp.NewToolResultError("Limit must be between 1 and 10000"), nil
	}

	fromID := req.GetInt("from", 0)
	toID := req.GetInt("to", 0)

	if dbAvailable() {
		return getTrackDB(ctx, trackIDStr, fromID, toID, limit)
	}
	return getTrackAPI(ctx, trackIDStr, fromID, toID, limit)
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
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The .unit. field indicates measurement units - CPM means .counts per minute. NOT .counts per second.. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I.ll, I.m, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: .Latest reading: X CPM at location Y. NOT .I found a reading of X CPM. or .Perfect! The sensor shows..... State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
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
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The .unit. field indicates measurement units - CPM means .counts per minute. NOT .counts per second.. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I.ll, I.m, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: .Latest reading: X CPM at location Y. NOT .I found a reading of X CPM. or .Perfect! The sensor shows..... State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}
