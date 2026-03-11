package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

var queryRadiationToolDef = mcp.NewTool("query_radiation",
	mcp.WithDescription("Find radiation measurements near a geographic location. Returns measurements within a specified radius of the given coordinates. For villages and rural areas use a radius of at least 25000-50000m to account for geocoding imprecision. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool. CRITICAL: Present all findings in an objective, scientific manner without using personal pronouns (I, we, I'll, you) or conversational language (Perfect!, Great!). Format as factual statements only."),
	mcp.WithNumber("lat",
		mcp.Description("Latitude (-90 to 90)"),
		mcp.Min(-90), mcp.Max(90),
		mcp.Required(),
	),
	mcp.WithNumber("lon",
		mcp.Description("Longitude (-180 to 180)"),
		mcp.Min(-180), mcp.Max(180),
		mcp.Required(),
	),
	mcp.WithNumber("radius_m",
		mcp.Description("Search radius in meters (default: 1500, max: 50000)"),
		mcp.Min(25), mcp.Max(50000),
		mcp.DefaultNumber(1500),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of results to return (default: 25, max: 10000)"),
		mcp.Min(1), mcp.Max(10000),
		mcp.DefaultNumber(25),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleQueryRadiation(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	lat, err := req.RequireFloat("lat")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	lon, err := req.RequireFloat("lon")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	radiusM := req.GetFloat("radius_m", 1500)
	limit := req.GetInt("limit", 25)

	if lat < -90 || lat > 90 {
		return mcp.NewToolResultError("Latitude must be between -90 and 90"), nil
	}
	if lon < -180 || lon > 180 {
		return mcp.NewToolResultError("Longitude must be between -180 and 180"), nil
	}
	if radiusM < 25 || radiusM > 50000 {
		return mcp.NewToolResultError("Radius must be between 25 and 50000 meters"), nil
	}
	if limit < 1 || limit > 10000 {
		return mcp.NewToolResultError("Limit must be between 1 and 10000"), nil
	}

	if dbAvailable() {
		return queryRadiationDB(ctx, lat, lon, radiusM, limit)
	}
	return queryRadiationAPI(ctx, lat, lon, radiusM, limit)
}

func queryRadiationDB(ctx context.Context, lat, lon, radiusM float64, limit int) (*mcp.CallToolResult, error) {
	// Use a bounding box pre-filter (&&) to hit the geometry spatial index first,
	// then refine with ST_DWithin on geography for precise meter-based distance.
	// Without the bbox filter, the geography cast bypasses the index → full table scan → timeout.
	//
	// PERFORMANCE: Use a subquery to filter and sort BEFORE joining to uploads/users.
	// This limits the join to only N rows instead of joining 90k+ rows then sorting.
	query := `
		WITH top_markers AS (
			SELECT m.id, m.doserate, m.date, m.lat, m.lon,
				m.device_id, m.altitude, m.detector, m.trackid, m.has_spectrum, m.geom
			FROM markers m
			WHERE m.geom && ST_Expand(ST_SetSRID(ST_MakePoint($2, $1), 4326), $3 / 111000.0)
			  AND ST_DWithin(m.geom::geography, ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography, $3)
			ORDER BY m.date DESC
			LIMIT $4
		)
		SELECT m.id, m.doserate AS value, 'µSv/h' AS unit,
			to_timestamp(m.date) AS captured_at,
			m.lat AS latitude, m.lon AS longitude,
			m.device_id, m.altitude AS height, m.detector,
			m.trackid, m.has_spectrum,
			ST_Distance(m.geom::geography, ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography) AS distance_m,
			u.internal_user_id, usr.username AS uploader_username, usr.email AS uploader_email
		FROM top_markers m
		LEFT JOIN uploads u ON u.track_id = m.trackid
		LEFT JOIN users usr ON u.internal_user_id = usr.id::text
		ORDER BY m.date DESC`

	rows, err := queryRows(ctx, query, lat, lon, radiusM, limit)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get total count (with same bbox pre-filter for performance)
	countRow, _ := queryRow(ctx, `
		SELECT count(*) AS total
		FROM markers m
		WHERE m.geom && ST_Expand(ST_SetSRID(ST_MakePoint($2, $1), 4326), $3 / 111000.0)
		  AND ST_DWithin(m.geom::geography, ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography, $3)`,
		lat, lon, radiusM)
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
			"device_id":    r["device_id"],
			"height":       r["height"],
			"detector":     r["detector"],
			"track_id":     r["trackid"],
			"has_spectrum":  r["has_spectrum"],
			"distance_m":   r["distance_m"],
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
		"query": map[string]any{
			"lat":      lat,
			"lon":      lon,
			"radius_m": radiusM,
		},
		"measurements": measurements,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The .unit. field indicates measurement units - CPM means .counts per minute. NOT .counts per second.. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I.ll, I.m, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: .Latest reading: X CPM at location Y. NOT .I found a reading of X CPM. or .Perfect! The sensor shows..... State only objective facts and measurements. (3) **FORMATTING — REQUIRED**: Always present results in a markdown table. Every location (lat/lon pair) MUST be a clickable map link: [lat°N, lon°E](https://simplemap.safecast.org/?lat=LAT&lon=LON&zoom=15). Never show plain coordinates without a link.",
		"_next_step": map[string]any{
			"instruction": "REQUIRED: This result contains ONLY historical mobile survey data. You MUST immediately call sensor_current using the bounding box below to check for real-time fixed sensors. Do NOT report 'no real-time data' until sensor_current has been called.",
			"call_sensor_current_with": map[string]any{
				"min_lat": lat - 0.5,
				"max_lat": lat + 0.5,
				"min_lon": lon - 0.5,
				"max_lon": lon + 0.5,
			},
		},
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}

func queryRadiationAPI(ctx context.Context, lat, lon, radiusM float64, limit int) (*mcp.CallToolResult, error) {
	resp, err := client.GetLatestNearby(ctx, lat, lon, radiusM, limit)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	markers, _ := resp["markers"].([]any)
	normalized := make([]map[string]any, 0, len(markers))
	for _, raw := range markers {
		if m, ok := raw.(map[string]any); ok {
			normalized = append(normalized, normalizeLatestMarker(m))
		}
	}

	result := map[string]any{
		"count":  len(normalized),
		"source": "api",
		"query": map[string]any{
			"lat":      lat,
			"lon":      lon,
			"radius_m": radiusM,
		},
		"measurements": normalized,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The .unit. field indicates measurement units - CPM means .counts per minute. NOT .counts per second.. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I.ll, I.m, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: .Latest reading: X CPM at location Y. NOT .I found a reading of X CPM. or .Perfect! The sensor shows..... State only objective facts and measurements. (3) **FORMATTING — REQUIRED**: Always present results in a markdown table. Every location (lat/lon pair) MUST be a clickable map link: [lat°N, lon°E](https://simplemap.safecast.org/?lat=LAT&lon=LON&zoom=15). Never show plain coordinates without a link.",
		"_next_step": map[string]any{
			"instruction": "REQUIRED: This result contains ONLY historical mobile survey data. You MUST immediately call sensor_current using the bounding box below to check for real-time fixed sensors. Do NOT report 'no real-time data' until sensor_current has been called.",
			"call_sensor_current_with": map[string]any{
				"min_lat": lat - 0.5,
				"max_lat": lat + 0.5,
				"min_lon": lon - 0.5,
				"max_lon": lon + 0.5,
			},
		},
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}
