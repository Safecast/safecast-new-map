package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

var searchAreaToolDef = mcp.NewTool("search_area",
	mcp.WithDescription("Find radiation measurements within a geographic bounding box. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool. CRITICAL: Present all findings in an objective, scientific manner without using personal pronouns (I, we, I'll, you) or conversational language (Perfect!, Great!). Format as factual statements only."),
	mcp.WithNumber("min_lat",
		mcp.Description("Southern boundary latitude"),
		mcp.Min(-90), mcp.Max(90),
		mcp.Required(),
	),
	mcp.WithNumber("max_lat",
		mcp.Description("Northern boundary latitude"),
		mcp.Min(-90), mcp.Max(90),
		mcp.Required(),
	),
	mcp.WithNumber("min_lon",
		mcp.Description("Western boundary longitude"),
		mcp.Min(-180), mcp.Max(180),
		mcp.Required(),
	),
	mcp.WithNumber("max_lon",
		mcp.Description("Eastern boundary longitude"),
		mcp.Min(-180), mcp.Max(180),
		mcp.Required(),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of results to return (default: 100, max: 10000)"),
		mcp.Min(1), mcp.Max(10000),
		mcp.DefaultNumber(100),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleSearchArea(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	minLat, err := req.RequireFloat("min_lat")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	maxLat, err := req.RequireFloat("max_lat")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	minLon, err := req.RequireFloat("min_lon")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	maxLon, err := req.RequireFloat("max_lon")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	limit := req.GetInt("limit", 100)

	if minLat < -90 || minLat > 90 || maxLat < -90 || maxLat > 90 {
		return mcp.NewToolResultError("Latitude must be between -90 and 90"), nil
	}
	if minLon < -180 || minLon > 180 || maxLon < -180 || maxLon > 180 {
		return mcp.NewToolResultError("Longitude must be between -180 and 180"), nil
	}
	if minLat >= maxLat {
		return mcp.NewToolResultError("min_lat must be less than max_lat"), nil
	}
	if minLon >= maxLon {
		return mcp.NewToolResultError("min_lon must be less than max_lon"), nil
	}
	if limit < 1 || limit > 10000 {
		return mcp.NewToolResultError("Limit must be between 1 and 10000"), nil
	}

	if dbAvailable() {
		return searchAreaDB(ctx, minLat, maxLat, minLon, maxLon, limit)
	}
	return searchAreaAPI(ctx, minLat, maxLat, minLon, maxLon, limit)
}

func searchAreaDB(ctx context.Context, minLat, maxLat, minLon, maxLon float64, limit int) (*mcp.CallToolResult, error) {
	query := `
		SELECT m.id, m.doserate AS value, 'µSv/h' AS unit,
			to_timestamp(m.date) AS captured_at,
			m.lat AS latitude, m.lon AS longitude,
			m.device_id, m.altitude AS height, m.detector,
			m.trackid, m.has_spectrum,
			u.internal_user_id, usr.username AS uploader_username, usr.email AS uploader_email
		FROM markers m
		LEFT JOIN uploads u ON u.track_id = m.trackid
		LEFT JOIN users usr ON u.internal_user_id = usr.id::text
		WHERE m.geom && ST_MakeEnvelope($1, $2, $3, $4, 4326)
		ORDER BY m.date DESC
		LIMIT $5`

	rows, err := queryRows(ctx, query, minLon, minLat, maxLon, maxLat, limit)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	countRow, _ := queryRow(ctx, `
		SELECT count(*) AS total
		FROM markers m
		WHERE m.geom && ST_MakeEnvelope($1, $2, $3, $4, 4326)`,
		minLon, minLat, maxLon, maxLat)
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
	for i, r := range rows {
		measurement := map[string]any{
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
			"track_id":    r["trackid"],
			"has_spectrum": r["has_spectrum"],
		}

		// Add uploader information if available
		if uploaderUsername, ok := r["uploader_username"]; ok && uploaderUsername != nil && uploaderUsername != "" {
			measurement["uploader"] = map[string]any{
				"username": uploaderUsername,
				"email":    r["uploader_email"],
			}
		}

		measurements[i] = measurement
	}

	result := map[string]any{
		"count":           len(measurements),
		"total_available": total,
		"source":          "database",
		"bbox": map[string]any{
			"min_lat": minLat,
			"max_lat": maxLat,
			"min_lon": minLon,
			"max_lon": maxLon,
		},
		"measurements": measurements,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The .unit. field indicates measurement units - CPM means .counts per minute. NOT .counts per second.. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I.ll, I.m, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: .Latest reading: X CPM at location Y. NOT .I found a reading of X CPM. or .Perfect! The sensor shows..... State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}

func searchAreaAPI(ctx context.Context, minLat, maxLat, minLon, maxLon float64, limit int) (*mcp.CallToolResult, error) {
	markers, err := client.GetMarkers(ctx, minLat, minLon, maxLat, maxLon)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if limit > len(markers) {
		limit = len(markers)
	}
	limited := markers[:limit]

	normalized := make([]map[string]any, len(limited))
	for i, m := range limited {
		normalized[i] = normalizeGetMarker(m)
	}

	result := map[string]any{
		"count":         len(normalized),
		"total_in_bbox": len(markers),
		"source":        "api",
		"bbox": map[string]any{
			"min_lat": minLat,
			"max_lat": maxLat,
			"min_lon": minLon,
			"max_lon": maxLon,
		},
		"measurements": normalized,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The .unit. field indicates measurement units - CPM means .counts per minute. NOT .counts per second.. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I.ll, I.m, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: .Latest reading: X CPM at location Y. NOT .I found a reading of X CPM. or .Perfect! The sensor shows..... State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	default:
		return 0, false
	}
}
