package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
)

var listSensorsToolDef = mcp.NewTool("list_sensors",
	mcp.WithDescription("Discover active fixed sensors (Pointcast, Solarcast, bGeigieZen, Notehub/Radnote, nGeigie, etc.) by location or type, returning device IDs, locations, status, and last reading timestamp. Use for sensor discovery and metadata only — this tool does NOT return radiation readings. When the user wants actual radiation values, use sensor_current instead. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool. CRITICAL: Present all findings in an objective, scientific manner without using personal pronouns (I, we, I'll, you) or conversational language (Perfect!, Great!). Format as factual statements only."),
	mcp.WithString("type",
		mcp.Description("Filter by sensor type (e.g., 'Pointcast', 'Solarcast', 'bGeigieZen', etc.)"),
	),
	mcp.WithNumber("min_lat",
		mcp.Description("Southern boundary for geographic filter"),
		mcp.Min(-90), mcp.Max(90),
	),
	mcp.WithNumber("max_lat",
		mcp.Description("Northern boundary for geographic filter"),
		mcp.Min(-90), mcp.Max(90),
	),
	mcp.WithNumber("min_lon",
		mcp.Description("Western boundary for geographic filter"),
		mcp.Min(-180), mcp.Max(180),
	),
	mcp.WithNumber("max_lon",
		mcp.Description("Eastern boundary for geographic filter"),
		mcp.Min(-180), mcp.Max(180),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of sensors to return for display (default: 50, max: 1000). The total_count in the result always reflects ALL matching sensors regardless of this limit."),
		mcp.Min(1), mcp.Max(1000),
		mcp.DefaultNumber(50),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleListSensors(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sensorType := req.GetString("type", "")
	minLat := req.GetFloat("min_lat", -90)
	maxLat := req.GetFloat("max_lat", 90)
	minLon := req.GetFloat("min_lon", -180)
	maxLon := req.GetFloat("max_lon", 180)
	limit := req.GetInt("limit", 50)

	if limit < 1 || limit > 1000 {
		return mcp.NewToolResultError("Limit must be between 1 and 1000"), nil
	}

	if dbAvailable() {
		return listSensorsDB(ctx, sensorType, minLat, maxLat, minLon, maxLon, limit)
	}

	return mcp.NewToolResultError("Database connection required for list_sensors tool. Please ensure DATABASE_URL is set to access real-time sensor data."), nil
}

// findRealtimeTable discovers the realtime sensor table in the database.
// Returns (tableName, availableTables, error). tableName is "" if not found.
func findRealtimeTable(ctx context.Context) (string, []string, error) {
	tablesQuery := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`
	tableRows, err := queryRows(ctx, tablesQuery)
	if err != nil {
		return "", nil, fmt.Errorf("could not query database schema: %w", err)
	}

	availableTables := make([]string, 0, len(tableRows))
	realtimeTable := ""
	for _, row := range tableRows {
		if tableName, ok := row["table_name"].(string); ok {
			availableTables = append(availableTables, tableName)
			if tableName == "realtime_measurements" ||
				tableName == "measurements_realtime" ||
				tableName == "sensors" ||
				tableName == "devices" {
				realtimeTable = tableName
			}
		}
	}
	return realtimeTable, availableTables, nil
}

// listSensorsQuery executes the sensor discovery query against realtimeTable.
// Returns (sensors, totalCount, error). totalCount is the COUNT(*) for the
// given filters independent of limit.
func listSensorsQuery(ctx context.Context, realtimeTable, sensorType string, minLat, maxLat, minLon, maxLon float64, limit int) ([]map[string]any, int, error) {
	// --- total count query ---
	var countQuery string
	var countArgs []interface{}
	if sensorType != "" {
		countQuery = fmt.Sprintf(`
			SELECT COUNT(DISTINCT device_id) AS total
			FROM %s
			WHERE lat >= $1 AND lat <= $2 AND lon >= $3 AND lon <= $4
			  AND (COALESCE(transport, '') ILIKE $5 OR COALESCE(device_name, '') ILIKE $5)`,
			realtimeTable)
		countArgs = []interface{}{minLat, maxLat, minLon, maxLon, "%" + sensorType + "%"}
	} else {
		countQuery = fmt.Sprintf(`
			SELECT COUNT(DISTINCT device_id) AS total
			FROM %s
			WHERE lat >= $1 AND lat <= $2 AND lon >= $3 AND lon <= $4`,
			realtimeTable)
		countArgs = []interface{}{minLat, maxLat, minLon, maxLon}
	}

	totalCount := 0
	if countRow, err := queryRow(ctx, countQuery, countArgs...); err == nil {
		switch v := countRow["total"].(type) {
		case int64:
			totalCount = int(v)
		case int32:
			totalCount = int(v)
		case int:
			totalCount = v
		}
	}

	// --- sensor rows query ---
	var query string
	var args []interface{}

	if sensorType != "" {
		query = fmt.Sprintf(`
			SELECT
				rm.device_id,
				COALESCE(rm.device_name, rm.device_id) AS device_name,
				COALESCE(rm.transport, '') AS transport,
				rm.lat AS latitude,
				rm.lon AS longitude,
				to_timestamp(rm.measured_at) AS last_reading_at
			FROM %s rm
			INNER JOIN (
				SELECT device_id, MAX(measured_at) as max_measured_at
				FROM %s
				WHERE lat >= $1 AND lat <= $2 AND lon >= $3 AND lon <= $4
					AND (COALESCE(transport, '') ILIKE $5 OR COALESCE(device_name, '') ILIKE $5)
				GROUP BY device_id
			) latest ON rm.device_id = latest.device_id AND rm.measured_at = latest.max_measured_at
			WHERE rm.lat >= $1 AND rm.lat <= $2 AND rm.lon >= $3 AND rm.lon <= $4
			ORDER BY rm.measured_at DESC
			LIMIT $6`, realtimeTable, realtimeTable)
		args = []interface{}{minLat, maxLat, minLon, maxLon, "%" + sensorType + "%", limit}
	} else {
		query = fmt.Sprintf(`
			SELECT
				rm.device_id,
				COALESCE(rm.device_name, rm.device_id) AS device_name,
				COALESCE(rm.transport, '') AS transport,
				rm.lat AS latitude,
				rm.lon AS longitude,
				to_timestamp(rm.measured_at) AS last_reading_at
			FROM %s rm
			INNER JOIN (
				SELECT device_id, MAX(measured_at) as max_measured_at
				FROM %s
				WHERE lat >= $1 AND lat <= $2 AND lon >= $3 AND lon <= $4
				GROUP BY device_id
			) latest ON rm.device_id = latest.device_id AND rm.measured_at = latest.max_measured_at
			WHERE rm.lat >= $1 AND rm.lat <= $2 AND rm.lon >= $3 AND rm.lon <= $4
			ORDER BY rm.measured_at DESC
			LIMIT $5`, realtimeTable, realtimeTable)
		args = []interface{}{minLat, maxLat, minLon, maxLon, limit}
	}

	rows, err := queryRows(ctx, query, args...)
	if err != nil {
		return nil, totalCount, fmt.Errorf("error querying %s: %w", realtimeTable, err)
	}

	sensors := make([]map[string]any, len(rows))
	for i, r := range rows {
		sensors[i] = map[string]any{
			"device_id":   r["device_id"],
			"device_name": r["device_name"],
			"type":        r["transport"],
			"location": map[string]any{
				"latitude":  r["latitude"],
				"longitude": r["longitude"],
			},
			"last_reading_at": r["last_reading_at"],
		}
	}
	return sensors, totalCount, nil
}

// buildExportURL constructs the /api/sensors/export URL with the given filters.
func buildExportURL(sensorType string, minLat, maxLat, minLon, maxLon float64) string {
	params := url.Values{}
	params.Set("min_lat", fmt.Sprintf("%.6f", minLat))
	params.Set("max_lat", fmt.Sprintf("%.6f", maxLat))
	params.Set("min_lon", fmt.Sprintf("%.6f", minLon))
	params.Set("max_lon", fmt.Sprintf("%.6f", maxLon))
	if sensorType != "" {
		params.Set("type", sensorType)
	}
	return "/api/sensors/export?" + params.Encode()
}

func listSensorsDB(ctx context.Context, sensorType string, minLat, maxLat, minLon, maxLon float64, limit int) (*mcp.CallToolResult, error) {
	realtimeTable, availableTables, err := findRealtimeTable(ctx)
	if err != nil {
		return mcp.NewToolResultError("Could not query database schema: " + err.Error()), nil
	}

	if realtimeTable == "" {
		result := map[string]any{
			"message":          "No known real-time sensor data tables found in database.",
			"available_tables": availableTables,
			"suggestion":       "Real-time sensor data may not be available through this database connection.",
		}
		return jsonResult(result)
	}

	sensors, totalCount, err := listSensorsQuery(ctx, realtimeTable, sensorType, minLat, maxLat, minLon, maxLon, limit)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	exportURL := buildExportURL(sensorType, minLat, maxLat, minLon, maxLon)

	result := map[string]any{
		"count":       len(sensors),
		"total_count": totalCount,
		"source":      "database",
		"sensors":     sensors,
		"table_used":  realtimeTable,
		"_export_url": exportURL,
		"_ai_hint": `CRITICAL INSTRUCTIONS:
(1) CPM = counts per minute (NOT per second). Always report CPM values as counts per minute.
(2) Present data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. State only objective facts.
(3) FORMATTING — REQUIRED: Always present results in a markdown table. Every device_id MUST be a clickable map link: [device_id](https://simplemap.safecast.org/?lat=LATITUDE&lon=LONGITUDE&zoom=15) using the actual lat/lon from the location field. Never show plain device IDs without a link.
(4) COUNTS — When stating how many sensors exist, ALWAYS use total_count (the true database total), NOT the count field (which is just the query limit). Example: if total_count=131 and count=100, state "131 active sensors" — never mention the query limit.`,
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}
