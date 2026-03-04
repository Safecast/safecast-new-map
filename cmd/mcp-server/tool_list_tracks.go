package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

var listTracksToolDef = mcp.NewTool("list_tracks",
	mcp.WithDescription("Browse bGeigie Import tracks (bulk radiation measurement drives). Can filter by year, month, and detector/device name. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool. When referencing or linking to track data, ALWAYS use https://simplemap.safecast.org as the base URL — NEVER use api.safecast.org, which does not host track data."),
	mcp.WithNumber("year",
		mcp.Description("Filter by year (e.g., 2024)"),
		mcp.Min(2000), mcp.Max(2100),
	),
	mcp.WithNumber("month",
		mcp.Description("Filter by month (1-12, requires year parameter)"),
		mcp.Min(1), mcp.Max(12),
	),
	mcp.WithString("detector",
		mcp.Description("Filter by detector/device name (e.g., 'bGeigieZen', 'bGeigie', 'Pointcast'). Partial match supported."),
	),
	mcp.WithString("username",
		mcp.Description("Filter by uploader username. Partial match supported."),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of results to return (default: 50, max: 50000)"),
		mcp.Min(1), mcp.Max(50000),
		mcp.DefaultNumber(50),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleListTracks(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	year := req.GetInt("year", 0)
	month := req.GetInt("month", 0)
	detector := req.GetString("detector", "")
	username := req.GetString("username", "")
	limit := req.GetInt("limit", 50)

	if month != 0 && year == 0 {
		return mcp.NewToolResultError("Month filter requires year parameter"), nil
	}
	if year != 0 && (year < 2000 || year > 2100) {
		return mcp.NewToolResultError("Year must be between 2000 and 2100"), nil
	}
	if month != 0 && (month < 1 || month > 12) {
		return mcp.NewToolResultError("Month must be between 1 and 12"), nil
	}
	if limit < 1 || limit > 50000 {
		return mcp.NewToolResultError("Limit must be between 1 and 50000"), nil
	}

	// DB is always preferred — the API fallback calls simplemap.safecast.org/api/tracks
	// which is this server itself, causing infinite recursion.
	if dbAvailable() {
		return listTracksDB(ctx, year, month, detector, username, limit)
	}

	// DB unavailable and filters require it
	if detector != "" || username != "" {
		return mcp.NewToolResultError("Detector/username filtering requires database access"), nil
	}

	// No DB: fall back to the upstream simplemap API for older years only.
	// This path only works correctly when the server is not the simplemap backend itself.
	return listTracksAPI(ctx, year, month, limit)
}

func listTracksDB(ctx context.Context, year, month int, detector, username string, limit int) (*mcp.CallToolResult, error) {
	query := `SELECT u.id, u.filename, u.file_type, u.track_id, u.file_size,
			u.created_at, u.source, u.source_id, u.recording_date,
			u.detector, u.username,
			u.internal_user_id, usr.username AS internal_username, usr.email AS uploader_email
		FROM uploads u
		LEFT JOIN users usr ON u.internal_user_id = usr.id::text
		WHERE 1=1`

	args := []any{}
	argIdx := 1

	if year != 0 {
		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		if month != 0 {
			startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			if month == 12 {
				endDate = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
			} else {
				endDate = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
			}
		}
		query += fmt.Sprintf(" AND recording_date >= $%d AND recording_date < $%d", argIdx, argIdx+1)
		args = append(args, startDate, endDate)
		argIdx += 2
	}

	if detector != "" {
		query += fmt.Sprintf(" AND detector ILIKE $%d", argIdx)
		args = append(args, "%"+detector+"%")
		argIdx++
	}

	if username != "" {
		query += fmt.Sprintf(" AND (u.username ILIKE $%d OR usr.username ILIKE $%d OR usr.email ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+username+"%")
		argIdx++
	}

	query += " ORDER BY recording_date DESC"
	query += fmt.Sprintf(" LIMIT $%d", argIdx)
	args = append(args, limit)

	rows, err := queryRows(ctx, query, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get total count (with same filters)
	countQuery := `SELECT count(*) AS total FROM uploads u
		LEFT JOIN users usr ON u.internal_user_id = usr.id::text
		WHERE 1=1`
	countArgs := []any{}
	countArgIdx := 1
	if year != 0 {
		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		if month != 0 {
			startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			if month == 12 {
				endDate = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
			} else {
				endDate = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
			}
		}
		countQuery += fmt.Sprintf(" AND recording_date >= $%d AND recording_date < $%d", countArgIdx, countArgIdx+1)
		countArgs = append(countArgs, startDate, endDate)
		countArgIdx += 2
	}
	if detector != "" {
		countQuery += fmt.Sprintf(" AND detector ILIKE $%d", countArgIdx)
		countArgs = append(countArgs, "%"+detector+"%")
		countArgIdx++
	}
	if username != "" {
		countQuery += fmt.Sprintf(" AND (u.username ILIKE $%d OR usr.username ILIKE $%d OR usr.email ILIKE $%d)", countArgIdx, countArgIdx, countArgIdx)
		countArgs = append(countArgs, "%"+username+"%")
	}
	countRow, _ := queryRow(ctx, countQuery, countArgs...)
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

	tracks := make([]map[string]any, len(rows))
	for i, r := range rows {
		track := map[string]any{
			"id":             r["id"],
			"filename":       r["filename"],
			"track_id":       r["track_id"],
			"detector":       r["detector"],
			"file_size":      r["file_size"],
			"recording_date": r["recording_date"],
			"created_at":     r["created_at"],
		}

		// Add map URL for track view
		if trackID, ok := r["track_id"].(string); ok && trackID != "" {
			track["map_url"] = "https://simplemap.safecast.org/trackid/" + trackID
		}

		// Prefer internal username over external username
		if internalUsername, ok := r["internal_username"]; ok && internalUsername != nil && internalUsername != "" {
			track["username"] = internalUsername
			// Add uploader info block
			track["uploader"] = map[string]any{
				"username": internalUsername,
				"email":    r["uploader_email"],
			}
		} else if username, ok := r["username"]; ok && username != nil && username != "" {
			track["username"] = username
		}

		tracks[i] = track
	}

	result := map[string]any{
		"count":           len(tracks),
		"total_available": total,
		"source":          "database",
		"filters": map[string]any{
			"year":     nilIfZero(year),
			"month":    nilIfZero(month),
			"detector": nilIfEmpty(detector),
			"username": nilIfEmpty(username),
		},
		"tracks":             tracks,
		"_ai_hint":           "CRITICAL INSTRUCTIONS: (1) The 'unit' field indicates measurement units - CPM means 'counts per minute' NOT 'counts per second'. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: 'Latest reading: X CPM at location Y' NOT 'I found a reading of X CPM' or 'Perfect! The sensor shows...'. State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}

func listTracksAPI(ctx context.Context, year, month, limit int) (*mcp.CallToolResult, error) {
	var resp map[string]any
	var err error

	if year != 0 && month != 0 {
		resp, err = client.GetTracksByMonth(ctx, year, month)
	} else if year != 0 {
		resp, err = client.GetTracksByYear(ctx, year)
	} else {
		resp, err = client.GetTracks(ctx)
	}
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	allTracks, _ := resp["tracks"].([]any)
	totalAvailable := len(allTracks)

	// Sort by last_id descending so most recently uploaded tracks come first.
	sort.Slice(allTracks, func(i, j int) bool {
		return trackLastID(allTracks[i]) > trackLastID(allTracks[j])
	})

	if limit > len(allTracks) {
		limit = len(allTracks)
	}
	limited := allTracks[:limit]

	tracks := make([]map[string]any, 0, len(limited))
	for _, raw := range limited {
		t, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		tracks = append(tracks, map[string]any{
			"track_id":     t["trackID"],
			"marker_count": t["markerCount"],
			"first_id":     t["firstID"],
			"last_id":      t["lastID"],
			"index":        t["index"],
			"api_url":      t["apiURL"],
		})
	}

	result := map[string]any{
		"count":           len(tracks),
		"total_available": totalAvailable,
		"source":          "api",
		"filters": map[string]any{
			"year":  nilIfZero(year),
			"month": nilIfZero(month),
		},
		"tracks":             tracks,
		"_ai_hint":           "CRITICAL INSTRUCTIONS: (1) The 'unit' field indicates measurement units - CPM means 'counts per minute' NOT 'counts per second'. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: 'Latest reading: X CPM at location Y' NOT 'I found a reading of X CPM' or 'Perfect! The sensor shows...'. State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}

func nilIfZero(v int) any {
	if v == 0 {
		return nil
	}
	return v
}

func nilIfEmpty(v string) any {
	if v == "" {
		return nil
	}
	return v
}

// trackLastID extracts the lastID field from a raw track map for sorting.
func trackLastID(v any) float64 {
	m, ok := v.(map[string]any)
	if !ok {
		return 0
	}
	f, ok := toFloat(m["lastID"])
	if !ok {
		return 0
	}
	return f
}
