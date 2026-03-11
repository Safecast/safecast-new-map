package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// Tool Definitions

var queryAnalyticsToolDef = mcp.NewTool("query_analytics",
	mcp.WithDescription("Get usage statistics for MCP tools (call counts, duration). Powered by DuckDB local logs."),
)

var radiationStatsToolDef = mcp.NewTool("radiation_stats",
	mcp.WithDescription("Get aggregate radiation statistics from the Safecast database (e.g., average dose rate by year/month). Powered by DuckDB+Postgres."),
	mcp.WithString("interval",
		mcp.Description("Aggregation interval: 'year', 'month', or 'overall'"),
		mcp.Enum("year", "month", "overall"),
		mcp.DefaultString("year"),
	),
)

// Handlers

func handleQueryAnalytics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if duckDB == nil {
		return mcp.NewToolResultError("DuckDB analytics engine is not initialized"), nil
	}

	// Execute query
	rows, err := duckDB.Query(`
		SELECT tool_name, COUNT(*) as count,
               AVG(duration_ms) as avg_ms,
               MAX(duration_ms) as max_ms
		FROM mcp_query_log
		GROUP BY tool_name
		ORDER BY count DESC
	`)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Query failed: %v", err)), nil
	}
	defer rows.Close()

	var stats []map[string]any
	for rows.Next() {
		var tool string
		var count int64
		var avgMs, maxMs float64
		if err := rows.Scan(&tool, &count, &avgMs, &maxMs); err != nil {
			continue
		}
		stats = append(stats, map[string]any{
			"tool":   tool,
			"calls":  count,
			"avg_ms": avgMs,
			"max_ms": maxMs,
		})
	}

	return jsonResult(map[string]any{
		"stats":              stats,
		"source":             "duckdb_local_log",
		"_ai_hint":           "CRITICAL INSTRUCTIONS: (1) The 'unit' field indicates measurement units - CPM means 'counts per minute' NOT 'counts per second'. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: 'Latest reading: X CPM at location Y' NOT 'I found a reading of X CPM' or 'Perfect! The sensor shows...'. State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	})
}

func handleRadiationStats(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if duckDB == nil {
		return mcp.NewToolResultError("DuckDB analytics engine is not initialized"), nil
	}

	interval := req.GetString("interval", "year")

	var query string
	switch interval {
	case "year":
		// Query attached Postgres DB
		// Note: 'postgres_db' is the name we attached it as in duckdb_client.go
		query = `
			SELECT
				EXTRACT(YEAR FROM to_timestamp(date)::TIMESTAMP) AS year,
				COUNT(*) AS count,
				AVG(doserate) AS avg_value,
				MAX(doserate) AS max_value
			FROM postgres_db.public.markers
			WHERE doserate > 0 AND doserate < 1000
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 20
		`
	case "month":
		query = `
			SELECT
				DATE_TRUNC('month', to_timestamp(date)::TIMESTAMP) AS month,
				COUNT(*) AS count,
				AVG(doserate) AS avg_value
			FROM postgres_db.public.markers
			WHERE doserate > 0 AND doserate < 1000
			  AND date > CAST(EXTRACT(EPOCH FROM (now() - INTERVAL '1 year')) AS BIGINT)
			GROUP BY 1
			ORDER BY 1 DESC
		`
	default: // overall
		query = `
			SELECT
				COUNT(*) AS count,
				AVG(doserate) AS avg_value,
				MAX(doserate) AS max_value
			FROM postgres_db.public.markers
			WHERE doserate > 0 AND doserate < 1000
		`
	}

	// Execute against DuckDB which proxies to Postgres
	rows, err := duckDB.Query(query)
	if err != nil {
		// Provide helpful error if table doesn't exist (e.g. schema mismatch)
		return mcp.NewToolResultError(fmt.Sprintf("Analytics query failed (check if postgres is attached): %v", err)), nil
	}
	defer rows.Close()

	// Generic scanner for results
	cols, _ := rows.Columns()
	var results []map[string]any

	for rows.Next() {
		// Create generic pointers
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		row := make(map[string]any)
		for i, colName := range cols {
			val := columns[i]
			// Handle byte arrays (often strings in db drivers)
			if b, ok := val.([]byte); ok {
				row[colName] = string(b)
			} else {
				row[colName] = val
			}
		}
		results = append(results, row)
	}

	return jsonResult(map[string]any{
		"interval":           interval,
		"data":               results,
		"source":             "duckdb_postgres_attach",
		"_ai_hint":           "CRITICAL INSTRUCTIONS: (1) The 'unit' field indicates measurement units - CPM means 'counts per minute' NOT 'counts per second'. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: 'Latest reading: X CPM at location Y' NOT 'I found a reading of X CPM' or 'Perfect! The sensor shows...'. State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	})
}
