package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

var queryDuckDBLogsToolDef = mcp.NewTool(
	"query_duckdb_logs",
	mcp.WithDescription(
		"Query MCP AI logs stored in DuckDB. Supports simple SQL SELECT queries.",
	),
	mcp.WithString(
		"query",
		mcp.Required(),
		mcp.Description("SQL SELECT query to execute against mcp_ai_query_log"),
	),
)

func handleQueryDuckDBLogs(
	ctx context.Context,
	req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	if duckDB == nil {
		return mcp.NewToolResultText("DuckDB not initialized"), nil
	}

	args, ok := req.Params.Arguments.(map[string]any)
	if !ok || args == nil {
		return mcp.NewToolResultText("Missing or invalid arguments"), nil
	}

	q, ok := args["query"].(string)
	if !ok || strings.TrimSpace(q) == "" {
		return mcp.NewToolResultText("Missing or invalid 'query' argument"), nil
	}

	query := strings.TrimSpace(q)
	if !strings.HasPrefix(strings.ToUpper(query), "SELECT") {
		return mcp.NewToolResultText("Only SELECT queries are allowed"), nil
	}

	rows, err := duckDB.Query(query)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Query error: %v", err)), nil
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Columns error: %v", err)), nil
	}

	var results strings.Builder
	for rows.Next() {
		values := make([]any, len(cols))
		pointers := make([]any, len(cols))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			results.WriteString(fmt.Sprintf("Scan error: %v\n", err))
			continue
		}
		for i, col := range cols {
			if i > 0 {
				results.WriteString("  ")
			}
			results.WriteString(fmt.Sprintf("%s: %v", col, values[i]))
		}
		results.WriteString("\n")
	}

	if err := rows.Err(); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Rows error: %v", err)), nil
	}

	out := results.String()
	if out == "" {
		out = "No results"
	}
	return mcp.NewToolResultText(out), nil
}
