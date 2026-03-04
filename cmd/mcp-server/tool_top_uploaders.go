package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

var topUploadersToolDef = mcp.NewTool("top_uploaders",
	mcp.WithDescription("Get statistics about which users or devices uploaded the most radiation measurement data to Safecast. Supports grouping by user (default) or device. Returns aggregated upload counts, individual marker counts, file sizes, and associated devices/users. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool."),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of top uploaders to return (default: 20, max: 100)"),
		mcp.Min(1), mcp.Max(100),
		mcp.DefaultNumber(20),
	),
	mcp.WithString("sort_by",
		mcp.Description("Sort by 'upload_count' (number of tracks), 'marker_count' (number of individual measurements), or 'total_size' (total data in bytes). Default: marker_count"),
		mcp.Enum("upload_count", "marker_count", "total_size"),
		mcp.DefaultString("marker_count"),
	),
	mcp.WithString("group_by",
		mcp.Description("Group results by 'user' (uploader) or 'device' (individual device). Default: user"),
		mcp.Enum("user", "device"),
		mcp.DefaultString("user"),
	),
	mcp.WithNumber("year",
		mcp.Description("Filter by year (e.g., 2024, 2026). Optional."),
		mcp.Min(2000), mcp.Max(2100),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleTopUploaders(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !dbAvailable() {
		return mcp.NewToolResultError("Database connection required for top uploaders query"), nil
	}

	limit := req.GetInt("limit", 20)
	sortBy := req.GetString("sort_by", "marker_count")
	groupBy := req.GetString("group_by", "user")
	year := req.GetInt("year", 0)

	if limit < 1 || limit > 100 {
		return mcp.NewToolResultError("Limit must be between 1 and 100"), nil
	}
	if year != 0 && (year < 2000 || year > 2100) {
		return mcp.NewToolResultError("Year must be between 2000 and 2100"), nil
	}

	// Build the query - joins uploads with markers to count individual measurement points
	var query string
	args := []any{}
	argIdx := 1

	if groupBy == "device" {
		// Group by device
		query = `
			WITH device_stats AS (
				SELECT
					u.detector AS device_name,
					COUNT(DISTINCT u.id) AS upload_count,
					SUM(COALESCE(u.file_size, 0)) AS total_size,
					COUNT(DISTINCT m.id) AS marker_count,
					array_agg(DISTINCT COALESCE(usr.username, u.username, 'Unknown') ORDER BY COALESCE(usr.username, u.username, 'Unknown')) FILTER (WHERE COALESCE(usr.username, u.username) IS NOT NULL) AS users
				FROM uploads u
				LEFT JOIN users usr ON u.internal_user_id = usr.id::text
				LEFT JOIN markers m ON u.track_id = m.track_id
				WHERE u.detector IS NOT NULL`

		// Add year filter if specified
		if year != 0 {
			query += fmt.Sprintf(" AND EXTRACT(YEAR FROM recording_date) = $%d", argIdx)
			args = append(args, year)
			argIdx++
		}

		query += `
				GROUP BY u.detector
			)
			SELECT
				device_name,
				upload_count,
				marker_count,
				total_size,
				users,
				CASE
					WHEN array_length(users, 1) = 1 THEN users[1]
					WHEN array_length(users, 1) > 1 THEN 'Multiple'
					ELSE NULL
				END AS primary_user
			FROM device_stats`
	} else {
		// Group by user (default)
		query = `
			WITH uploader_stats AS (
				SELECT
					COALESCE(usr.username, u.username, 'Unknown') AS username,
					u.internal_user_id,
					COUNT(DISTINCT u.id) AS upload_count,
					SUM(COALESCE(u.file_size, 0)) AS total_size,
					COUNT(DISTINCT m.id) AS marker_count,
					array_agg(DISTINCT u.detector ORDER BY u.detector) FILTER (WHERE u.detector IS NOT NULL) AS devices
				FROM uploads u
				LEFT JOIN users usr ON u.internal_user_id = usr.id::text
				LEFT JOIN markers m ON u.track_id = m.track_id
				WHERE 1=1`

		// Add year filter if specified
		if year != 0 {
			query += fmt.Sprintf(" AND EXTRACT(YEAR FROM recording_date) = $%d", argIdx)
			args = append(args, year)
			argIdx++
		}

		query += `
				GROUP BY COALESCE(usr.username, u.username, 'Unknown'), u.internal_user_id
			)
			SELECT
				username,
				upload_count,
				marker_count,
				total_size,
				devices,
				CASE
					WHEN array_length(devices, 1) = 1 THEN devices[1]
					WHEN array_length(devices, 1) > 1 THEN 'Multiple'
					ELSE NULL
				END AS primary_device
			FROM uploader_stats`
	}

	// Add ORDER BY clause based on sort_by parameter
	if sortBy == "total_size" {
		query += " ORDER BY total_size DESC, marker_count DESC"
	} else if sortBy == "marker_count" {
		query += " ORDER BY marker_count DESC, upload_count DESC"
	} else {
		query += " ORDER BY upload_count DESC, marker_count DESC"
	}

	query += fmt.Sprintf(" LIMIT $%d", argIdx)
	args = append(args, limit)

	rows, err := queryRows(ctx, query, args...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Query failed: %v", err)), nil
	}

	results := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		var item map[string]any
		if groupBy == "device" {
			item = map[string]any{
				"device_name":    r["device_name"],
				"upload_count":   r["upload_count"],
				"marker_count":   r["marker_count"],
				"total_size_mb":  bytesToMB(r["total_size"]),
				"users":          r["users"],
				"primary_user":   r["primary_user"],
			}
		} else {
			item = map[string]any{
				"username":       r["username"],
				"upload_count":   r["upload_count"],
				"marker_count":   r["marker_count"],
				"total_size_mb":  bytesToMB(r["total_size"]),
				"devices":        r["devices"],
				"primary_device": r["primary_device"],
			}
		}
		results = append(results, item)
	}

	var dataKey string
	if groupBy == "device" {
		dataKey = "devices"
	} else {
		dataKey = "uploaders"
	}

	result := map[string]any{
		"count":    len(results),
		"sort_by":  sortBy,
		"group_by": groupBy,
		dataKey:    results,
		"filters": map[string]any{
			"year": nilIfZero(year),
		},
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The 'unit' field indicates measurement units - CPM means 'counts per minute' NOT 'counts per second'. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: 'Latest reading: X CPM at location Y' NOT 'I found a reading of X CPM' or 'Perfect! The sensor shows...'. State only objective facts and measurements. (3) When presenting statistics: 'upload_count' = number of track files (survey routes) uploaded, 'marker_count' = total number of individual measurement points across all tracks, 'total_size_mb' = total data volume in megabytes. (4) Data is grouped by 'group_by' parameter: when grouped by 'user', shows per-user stats with their devices; when grouped by 'device', shows per-device stats with their users.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}

// bytesToMB converts bytes to megabytes with 1 decimal place precision
func bytesToMB(v any) float64 {
	var bytes float64
	switch val := v.(type) {
	case int64:
		bytes = float64(val)
	case float64:
		bytes = val
	case int:
		bytes = float64(val)
	default:
		return 0
	}
	mb := bytes / (1024 * 1024)
	return float64(int(mb*10)) / 10 // Round to 1 decimal place
}
