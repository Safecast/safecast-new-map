package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

var dbInfoToolDef = mcp.NewTool("db_info",
	mcp.WithDescription("Get database connection information and replication status (diagnostic tool). IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool."),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleDBInfo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !dbAvailable() {
		return mcp.NewToolResultText("No database connection (using REST API fallback)"), nil
	}

	// Get basic connection info
	info := make(map[string]any)
	
	// Query PostgreSQL version
	versionRow, err := queryRow(ctx, "SELECT version() AS version")
	if err == nil && versionRow != nil {
		info["version"] = versionRow["version"]
	}

	// Query current database name
	dbNameRow, err := queryRow(ctx, "SELECT current_database() AS db_name")
	if err == nil && dbNameRow != nil {
		info["database"] = dbNameRow["db_name"]
	}

	// Query current user
	userRow, err := queryRow(ctx, "SELECT current_user AS username")
	if err == nil && userRow != nil {
		info["user"] = userRow["username"]
	}

	// Query server address (useful for identifying which server we're connected to)
	addrRow, err := queryRow(ctx, "SELECT inet_server_addr() AS server_addr, inet_server_port() AS server_port")
	if err == nil && addrRow != nil {
		info["server_address"] = addrRow["server_addr"]
		info["server_port"] = addrRow["server_port"]
	}

	// Check if this is a replica (read-only mode)
	isReplica := false
	replicationRow, err := queryRow(ctx, "SELECT pg_is_in_recovery() AS in_recovery")
	if err == nil && replicationRow != nil {
		if val, ok := replicationRow["in_recovery"].(bool); ok {
			isReplica = val
		}
		info["is_replica"] = isReplica
		if isReplica {
			info["mode"] = "read replica (replication lag possible)"
		} else {
			info["mode"] = "primary/master"
		}
	}

	// If this is a replica, try to get replication lag
	if isReplica {
		lagRow, err := queryRow(ctx, `
			SELECT 
				EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp())) AS lag_seconds,
				pg_last_xact_replay_timestamp() AS last_replay_time
		`)
		if err == nil && lagRow != nil {
			info["replication_lag_seconds"] = lagRow["lag_seconds"]
			info["last_replay_time"] = lagRow["last_replay_time"]
		}
	}

	// Get table counts for context
	countsRow, err := queryRow(ctx, "SELECT count(*) AS total FROM uploads")
	if err == nil && countsRow != nil {
		if total, ok := countsRow["total"]; ok {
			switch v := total.(type) {
			case int64:
				info["uploads_count"] = v
			case float64:
				info["uploads_count"] = int64(v)
			}
		}
	}

	// Get highest upload ID
	maxRow, err := queryRow(ctx, "SELECT MAX(id) AS max_id FROM uploads")
	if err == nil && maxRow != nil {
		if maxID, ok := maxRow["max_id"]; ok {
			info["max_upload_id"] = maxID
		}
	}

	return jsonResult(map[string]any{
		"status":     "connected",
		"connection": info,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The .unit. field indicates measurement units - CPM means .counts per minute. NOT .counts per second.. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I.ll, I.m, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: .Latest reading: X CPM at location Y. NOT .I found a reading of X CPM. or .Perfect! The sensor shows..... State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	})
}
