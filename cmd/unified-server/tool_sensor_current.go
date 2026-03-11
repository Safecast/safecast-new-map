package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

var sensorCurrentToolDef = mcp.NewTool("sensor_current",
	mcp.WithDescription("Get the MOST RECENT readings from REAL-TIME fixed sensors. USE THIS TOOL for: bGeigieZen (type=geigiecast-zen, IDs like geigiecast-zen:65002), Pointcast (type=pointcast, IDs like pointcast:10042), Solarcast (type=solarcast), Notehub/Radnote (type=notehub, IDs like note:dev:867648049123019), nGeigie (type=ngeigie, IDs like ngeigie:101), device-tcp (IDs like safecast:3474557222). Use when users ask about 'current', 'latest', 'live', 'now', or 'real-time' data, OR to look up a specific fixed sensor by device ID. When searching by location, always call this tool AND query_radiation together to cover both fixed and mobile sources. Use a LARGE bounding box (at least ±0.5 degrees, ~50km) for villages and rural areas to account for geocoding imprecision. DO NOT use device_history for any of these fixed sensor types — device_history is ONLY for mobile bGeigie (type=geigiecast). DO NOT use query_radiation for current data. CPM = counts per minute (convert to µSv/h using ~0.0069 for LND 7318). Always report the captured_at timestamp. Present data objectively without personal pronouns."),
	mcp.WithString("device_id",
		mcp.Description("Specific device ID to get latest reading from"),
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
		mcp.Description("Maximum number of readings to return (default: 25, max: 1000)"),
		mcp.Min(1), mcp.Max(1000),
		mcp.DefaultNumber(25),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleSensorCurrent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	deviceID := req.GetString("device_id", "")
	minLat := req.GetFloat("min_lat", -90)
	maxLat := req.GetFloat("max_lat", 90)
	minLon := req.GetFloat("min_lon", -180)
	maxLon := req.GetFloat("max_lon", 180)
	limit := req.GetInt("limit", 25)

	if limit < 1 || limit > 1000 {
		return mcp.NewToolResultError("Limit must be between 1 and 1000"), nil
	}

	if dbAvailable() {
		return sensorCurrentDB(ctx, deviceID, minLat, maxLat, minLon, maxLon, limit)
	}
	
	// Fallback to API if database not available
	return mcp.NewToolResultError("Database connection required for sensor_current tool. Please ensure DATABASE_URL is set to access real-time sensor data."), nil
}

func sensorCurrentDB(ctx context.Context, deviceID string, minLat, maxLat, minLon, maxLon float64, limit int) (*mcp.CallToolResult, error) {
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
	
	var query string
	var args []interface{}

	if deviceID != "" {
		// Get latest reading from specific device
		query = fmt.Sprintf(`
			SELECT
				id,
				device_id,
				COALESCE(device_name, device_id) AS device_name,
				value,
				COALESCE(unit, 'µSv/h') AS unit,
				to_timestamp(measured_at) AS captured_at,
				lat AS latitude,
				lon AS longitude,
				COALESCE(transport, '') AS transport
			FROM %s
			WHERE device_id = $1
			  AND to_timestamp(measured_at) <= NOW()
			ORDER BY measured_at DESC
			LIMIT 1`, realtimeTable)

		args = []interface{}{deviceID}
	} else {
		// Get latest readings from all sensors in geographic area
		query = fmt.Sprintf(`
			SELECT
				rm.id,
				rm.device_id,
				COALESCE(rm.device_name, rm.device_id) AS device_name,
				rm.value,
				COALESCE(rm.unit, 'µSv/h') AS unit,
				to_timestamp(rm.measured_at) AS captured_at,
				rm.lat AS latitude,
				rm.lon AS longitude,
				COALESCE(rm.transport, '') AS transport
			FROM %s rm
			INNER JOIN (
				SELECT device_id, MAX(measured_at) as max_measured_at
				FROM %s
				WHERE lat >= $1 AND lat <= $2 AND lon >= $3 AND lon <= $4
				  AND to_timestamp(measured_at) <= NOW()
				GROUP BY device_id
			) latest ON rm.device_id = latest.device_id AND rm.measured_at = latest.max_measured_at
			WHERE rm.lat >= $1 AND rm.lat <= $2 AND rm.lon >= $3 AND rm.lon <= $4
			  AND to_timestamp(rm.measured_at) <= NOW()
			ORDER BY rm.measured_at DESC
			LIMIT $5`, realtimeTable, realtimeTable)

		args = []interface{}{minLat, maxLat, minLon, maxLon, limit}
	}

	rows, err := queryRows(ctx, query, args...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error querying %s table: %v", realtimeTable, err)), nil
	}

	readings := make([]map[string]any, len(rows))
	for i, r := range rows {
		// Fix incorrect unit: Geiger counters report in CPM (counts per minute), not CPS
		unit := r["unit"]
		if unitStr, ok := unit.(string); ok {
			unit = strings.ReplaceAll(strings.ReplaceAll(unitStr, "cps", "cpm"), "CPS", "CPM")
		}

		readings[i] = map[string]any{
			"id":          r["id"],
			"device_id":   r["device_id"],
			"device_name": r["device_name"],
			"value":       r["value"],
			"unit":        unit,
			"captured_at": r["captured_at"],
			"location": map[string]any{
				"latitude":  r["latitude"],
				"longitude": r["longitude"],
			},
			"type": r["transport"],
		}
	}

	result := map[string]any{
		"count":    len(readings),
		"source":   "database",
		"readings": readings,
		"table_used": realtimeTable,
		"available_tables": availableTables,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) **REAL-TIME DATA**: This tool returns the MOST RECENT readings from fixed sensors. Readings with future timestamps (sensor clock errors) are automatically filtered out. Always check the 'captured_at' timestamp and report it to the user - if the data is more than 24 hours old, mention this to the user and suggest checking if the sensor is still active. (2) **UNITS**: CPM means 'counts per minute' NOT 'counts per second'. Always convert to µSv/h using detector-specific factors (LND 7318: ~0.0069 µSv/h per CPM). (3) **TOOL SELECTION**: For latest sensor data, use 'sensor_current'. For historical trends, use 'sensor_history'. For mobile measurements, use 'device_history'. Do NOT use 'query_radiation' for current sensor data as it searches the historical markers table. (4) **PRESENTATION**: State objective facts only - no personal pronouns (I, we, you), exclamations, or conversational phrases. (5) **FORMATTING — REQUIRED**: Always present results in a markdown table. Every device_id MUST be a clickable map link using the format [device_id](https://simplemap.safecast.org/?lat=LATITUDE&lon=LONGITUDE&zoom=15) substituting the actual latitude and longitude from the location field. Example: [geigiecast-zen:65002](https://simplemap.safecast.org/?lat=34.48265&lon=136.16314&zoom=15). Never show plain device IDs without a link. Timestamps MUST be shown in UTC.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}