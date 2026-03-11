package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

var sensorHistoryToolDef = mcp.NewTool("sensor_history",
	mcp.WithDescription("Pull time-series data from REAL-TIME fixed sensors (Pointcast, Solarcast, bGeigieZen, etc.) over a date range. Use this tool for historical time-series from fixed sensors. NOT for mobile bGeigie devices - use device_history for those. The 'unit' field indicates the measurement unit - CPM means 'counts per minute' (NOT counts per second). Always present radiation values in µSv/h by converting from CPM using detector-specific factors. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool. CRITICAL: Present all findings in an objective, scientific manner without using personal pronouns (I, we, I'll, you) or conversational language (Perfect!, Great!). Format as factual statements only."),
	mcp.WithString("device_id",
		mcp.Description("Device identifier to get historical data from"),
		mcp.Required(),
	),
	mcp.WithString("start_date",
		mcp.Description("Start date in YYYY-MM-DD format"),
		mcp.Required(),
	),
	mcp.WithString("end_date",
		mcp.Description("End date in YYYY-MM-DD format (default: today)"),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of measurements to return (default: 200, max: 10000)"),
		mcp.Min(1), mcp.Max(10000),
		mcp.DefaultNumber(200),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleSensorHistory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	deviceID, err := req.RequireString("device_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	startDateStr, err := req.RequireString("start_date")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	endDateStr := req.GetString("end_date", "")
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02")
	}

	limit := req.GetInt("limit", 200)

	if limit < 1 || limit > 10000 {
		return mcp.NewToolResultError("Limit must be between 1 and 10000"), nil
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return mcp.NewToolResultError("start_date must be in YYYY-MM-DD format"), nil
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return mcp.NewToolResultError("end_date must be in YYYY-MM-DD format"), nil
	}

	if endDate.Before(startDate) {
		return mcp.NewToolResultError("end_date must be after start_date"), nil
	}

	if dbAvailable() {
		return sensorHistoryDB(ctx, deviceID, startDate, endDate, limit)
	}
	
	// Fallback to API if database not available
	return mcp.NewToolResultError("Database connection required for sensor_history tool. Please ensure DATABASE_URL is set to access real-time sensor data."), nil
}

func sensorHistoryDB(ctx context.Context, deviceID string, startDate, endDate time.Time, limit int) (*mcp.CallToolResult, error) {
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
	
	// Query the appropriate real-time table for time-series data
	query := fmt.Sprintf(`
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
			AND measured_at >= $2
			AND measured_at <= $3
			AND to_timestamp(measured_at) <= NOW()
		ORDER BY measured_at ASC
		LIMIT $4`, realtimeTable)

	startUnix := startDate.Unix()
	endUnix := endDate.Unix()

	rows, err := queryRows(ctx, query, deviceID, startUnix, endUnix, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error querying %s table: %v", realtimeTable, err)), nil
	}

	measurements := make([]map[string]any, len(rows))
	for i, r := range rows {
		// Fix incorrect unit: Geiger counters report in CPM (counts per minute), not CPS
		unit := r["unit"]
		if unitStr, ok := unit.(string); ok {
			unit = strings.ReplaceAll(strings.ReplaceAll(unitStr, "cps", "cpm"), "CPS", "CPM")
		}

		measurements[i] = map[string]any{
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
			"type":   r["transport"],
		}
	}

	capturedAfter := startDate.Format("2006-01-02") + " 00:00"
	capturedBefore := endDate.Format("2006-01-02") + " 23:59"

	result := map[string]any{
		"device": map[string]any{
			"id": deviceID,
		},
		"period": map[string]any{
			"start_date": capturedAfter,
			"end_date":   capturedBefore,
		},
		"count":        len(measurements),
		"source":       "database",
		"measurements": measurements,
		"table_used": realtimeTable,
		"available_tables": availableTables,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The 'unit' field indicates measurement units - CPM means 'counts per minute' NOT 'counts per second'. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: 'Latest reading: X CPM at location Y' NOT 'I found a reading of X CPM' or 'Perfect! The sensor shows...'. State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}