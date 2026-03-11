package main

import (
	"context"
	"fmt"

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
		mcp.Description("Maximum number of sensors to return (default: 50, max: 1000)"),
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
	
	// Fallback to API if database not available
	return mcp.NewToolResultError("Database connection required for list_sensors tool. Please ensure DATABASE_URL is set to access real-time sensor data."), nil
}

func listSensorsDB(ctx context.Context, sensorType string, minLat, maxLat, minLon, maxLon float64, limit int) (*mcp.CallToolResult, error) {
	// Check what tables are available in the database
	tablesQuery := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
		ORDER BY table_name
	`
	
	tableRows, err := queryRows(ctx, tablesQuery)
	if err != nil {
		return mcp.NewToolResultError("Could not query database schema: " + err.Error()), nil
	}
	
	// Look for tables that might contain real-time sensor data
	availableTables := make([]string, len(tableRows))
	realtimeTable := ""
	for i, row := range tableRows {
		if tableName, ok := row["table_name"].(string); ok {
			availableTables[i] = tableName
			// Check for possible real-time sensor data tables
			if tableName == "realtime_measurements" || 
			   tableName == "measurements_realtime" || 
			   tableName == "sensors" ||
			   tableName == "devices" {
				realtimeTable = tableName
			}
		}
	}
	
	if realtimeTable == "" {
		// If no real-time table found, return available tables for debugging
		result := map[string]any{
			"message": "No known real-time sensor data tables found in database.",
			"available_tables": availableTables,
			"suggestion": "Real-time sensor data may not be available through this database connection.",
		}
		return jsonResult(result)
	}
	
	// Query the appropriate real-time table to find unique devices/sensors
	var query string
	var args []interface{}

	if sensorType != "" {
		// Filter by sensor type
		// FIXED: Get the actual latest reading per device, not grouped by lat/lon
		// which causes stale data when sensors move or have multiple positions
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
		// No filter by type
		// FIXED: Get the actual latest reading per device, not grouped by lat/lon
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
		return mcp.NewToolResultError(fmt.Sprintf("Error querying %s table: %v", realtimeTable, err)), nil
	}

	sensors := make([]map[string]any, len(rows))
	for i, r := range rows {
		sensors[i] = map[string]any{
			"device_id":       r["device_id"],
			"device_name":     r["device_name"],
			"type":            r["transport"],
			"location": map[string]any{
				"latitude":  r["latitude"],
				"longitude": r["longitude"],
			},
			"last_reading_at": r["last_reading_at"],
		}
	}

	result := map[string]any{
		"count":   len(sensors),
		"source":  "database",
		"sensors": sensors,
		"table_used": realtimeTable,
		"available_tables": availableTables,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The 'unit' field indicates measurement units - CPM means 'counts per minute' NOT 'counts per second'. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: 'Latest reading: X CPM at location Y' NOT 'I found a reading of X CPM' or 'Perfect! The sensor shows...'. State only objective facts and measurements. (3) **FORMATTING — REQUIRED**: Always present results in a markdown table. Every device_id MUST be a clickable map link: [device_id](https://simplemap.safecast.org/?lat=LATITUDE&lon=LONGITUDE&zoom=15) using the actual lat/lon from the location field. Never show plain device IDs without a link.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}