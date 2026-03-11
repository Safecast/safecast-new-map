package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

var listSpectraToolDef = mcp.NewTool("list_spectra",
	mcp.WithDescription("Browse gamma spectroscopy records. Returns metadata without channel data. Use get_spectrum with a marker_id from results to fetch full channel data. Can filter by track ID, geographic bounds, file format, or device model. Call with no filters to get all spectra. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool."),
	mcp.WithNumber("min_lat",
		mcp.Description("Southern boundary latitude (requires all 4 bbox params)"),
		mcp.Min(-90), mcp.Max(90),
	),
	mcp.WithNumber("max_lat",
		mcp.Description("Northern boundary latitude (requires all 4 bbox params)"),
		mcp.Min(-90), mcp.Max(90),
	),
	mcp.WithNumber("min_lon",
		mcp.Description("Western boundary longitude (requires all 4 bbox params)"),
		mcp.Min(-180), mcp.Max(180),
	),
	mcp.WithNumber("max_lon",
		mcp.Description("Eastern boundary longitude (requires all 4 bbox params)"),
		mcp.Min(-180), mcp.Max(180),
	),
	mcp.WithString("source_format",
		mcp.Description("Filter by spectrum file format (e.g., 'spe', 'csv')"),
	),
	mcp.WithString("device_model",
		mcp.Description("Filter by detector/device model name"),
	),
	mcp.WithString("track_id",
		mcp.Description("Filter by track identifier (e.g., '8eh5m1', '8ZnI7f')"),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of results to return (default: 50, max: 500)"),
		mcp.Min(1), mcp.Max(500),
		mcp.DefaultNumber(50),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleListSpectra(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !dbAvailable() {
		return mcp.NewToolResultError("list_spectra requires a database connection (no REST API fallback available)"), nil
	}

	// Check which bbox params were provided
	argsMap, _ := req.Params.Arguments.(map[string]any)
	if argsMap == nil {
		argsMap = map[string]any{}
	}
	_, hasMinLat := argsMap["min_lat"]
	_, hasMaxLat := argsMap["max_lat"]
	_, hasMinLon := argsMap["min_lon"]
	_, hasMaxLon := argsMap["max_lon"]
	hasBBox := hasMinLat || hasMaxLat || hasMinLon || hasMaxLon

	var minLat, maxLat, minLon, maxLon float64
	if hasBBox {
		if !(hasMinLat && hasMaxLat && hasMinLon && hasMaxLon) {
			return mcp.NewToolResultError("All four bounding box parameters (min_lat, max_lat, min_lon, max_lon) must be provided together"), nil
		}
		minLat = req.GetFloat("min_lat", 0)
		maxLat = req.GetFloat("max_lat", 0)
		minLon = req.GetFloat("min_lon", 0)
		maxLon = req.GetFloat("max_lon", 0)
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
	}

	sourceFormat := req.GetString("source_format", "")
	deviceModel := req.GetString("device_model", "")
	trackID := req.GetString("track_id", "")
	limit := req.GetInt("limit", 50)
	if limit < 1 || limit > 500 {
		return mcp.NewToolResultError("Limit must be between 1 and 500"), nil
	}

	return listSpectraDB(ctx, hasBBox, minLat, maxLat, minLon, maxLon, sourceFormat, deviceModel, trackID, limit)
}

func listSpectraDB(ctx context.Context, hasBBox bool, minLat, maxLat, minLon, maxLon float64, sourceFormat, deviceModel, trackID string, limit int) (*mcp.CallToolResult, error) {
	// Exclude s.channels to avoid huge payloads
	baseSelect := `SELECT s.id, s.marker_id, s.channel_count, s.energy_min_kev, s.energy_max_kev,
			s.live_time_sec, s.real_time_sec, s.device_model, s.calibration,
			s.source_format, s.filename, s.created_at,
			m.doserate, m.lat, m.lon, to_timestamp(m.date) AS captured_at,
			m.trackid AS track_id,
			u.internal_user_id, usr.username AS uploader_username, usr.email AS uploader_email
		FROM spectra s
		JOIN markers m ON m.id = s.marker_id
		LEFT JOIN uploads u ON u.track_id = m.trackid
		LEFT JOIN users usr ON u.internal_user_id = usr.id::text
		WHERE 1=1`

	countBase := `SELECT count(*) AS total
		FROM spectra s
		JOIN markers m ON m.id = s.marker_id
		WHERE 1=1`

	args := []any{}
	countArgs := []any{}
	argIdx := 1
	countArgIdx := 1

	if hasBBox {
		clause := fmt.Sprintf(" AND m.geom && ST_MakeEnvelope($%d, $%d, $%d, $%d, 4326)",
			argIdx, argIdx+1, argIdx+2, argIdx+3)
		baseSelect += clause
		args = append(args, minLon, minLat, maxLon, maxLat)
		argIdx += 4

		countClause := fmt.Sprintf(" AND m.geom && ST_MakeEnvelope($%d, $%d, $%d, $%d, 4326)",
			countArgIdx, countArgIdx+1, countArgIdx+2, countArgIdx+3)
		countBase += countClause
		countArgs = append(countArgs, minLon, minLat, maxLon, maxLat)
		countArgIdx += 4
	}

	if sourceFormat != "" {
		baseSelect += fmt.Sprintf(" AND s.source_format = $%d", argIdx)
		args = append(args, sourceFormat)
		argIdx++

		countBase += fmt.Sprintf(" AND s.source_format = $%d", countArgIdx)
		countArgs = append(countArgs, sourceFormat)
		countArgIdx++
	}

	if deviceModel != "" {
		baseSelect += fmt.Sprintf(" AND s.device_model ILIKE $%d", argIdx)
		args = append(args, "%"+deviceModel+"%")
		argIdx++

		countBase += fmt.Sprintf(" AND s.device_model ILIKE $%d", countArgIdx)
		countArgs = append(countArgs, "%"+deviceModel+"%")
		countArgIdx++
	}

	if trackID != "" {
		baseSelect += fmt.Sprintf(" AND m.trackid = $%d", argIdx)
		args = append(args, trackID)
		argIdx++

		countBase += fmt.Sprintf(" AND m.trackid = $%d", countArgIdx)
		countArgs = append(countArgs, trackID)
		countArgIdx++
	}

	baseSelect += " ORDER BY s.created_at DESC"
	baseSelect += fmt.Sprintf(" LIMIT $%d", argIdx)
	args = append(args, limit)

	rows, err := queryRows(ctx, baseSelect, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	countRow, _ := queryRow(ctx, countBase, countArgs...)
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

	spectra := make([]map[string]any, len(rows))
	for i, r := range rows {
		spec := map[string]any{
			"spectrum_id":    r["id"],
			"marker_id":     r["marker_id"],
			"filename":      r["filename"],
			"source_format": r["source_format"],
			"device_model":  r["device_model"],
			"channel_count": r["channel_count"],
			"energy_range": map[string]any{
				"min_kev": r["energy_min_kev"],
				"max_kev": r["energy_max_kev"],
			},
			"live_time_sec": r["live_time_sec"],
			"real_time_sec": r["real_time_sec"],
			"calibration":   r["calibration"],
			"created_at":    r["created_at"],
			"marker": map[string]any{
				"doserate":    r["doserate"],
				"latitude":    r["lat"],
				"longitude":   r["lon"],
				"captured_at": r["captured_at"],
				"track_id":    r["track_id"],
			},
		}

		// Add uploader information if available
		if uploaderUsername, ok := r["uploader_username"]; ok && uploaderUsername != nil {
			spec["uploader"] = map[string]any{
				"username": uploaderUsername,
				"email":    r["uploader_email"],
			}
		}

		spectra[i] = spec
	}

	filters := map[string]any{}
	if hasBBox {
		filters["bbox"] = map[string]any{
			"min_lat": minLat, "max_lat": maxLat,
			"min_lon": minLon, "max_lon": maxLon,
		}
	}
	if sourceFormat != "" {
		filters["source_format"] = sourceFormat
	}
	if deviceModel != "" {
		filters["device_model"] = deviceModel
	}
	if trackID != "" {
		filters["track_id"] = trackID
	}

	result := map[string]any{
		"count":           len(spectra),
		"total_available": total,
		"source":          "database",
		"filters":         filters,
		"spectra":         spectra,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The .unit. field indicates measurement units - CPM means .counts per minute. NOT .counts per second.. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I.ll, I.m, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: .Latest reading: X CPM at location Y. NOT .I found a reading of X CPM. or .Perfect! The sensor shows..... State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}
