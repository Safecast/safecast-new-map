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
		"Query MCP AI logs stored in DuckDB. Supports simple SQL SELECT queries. Available tables: mcp_ai_query_log (tool execution logs), mcp_query_log (tool usage stats), chat_questions (user questions and AI answers from web-chat and map widget with metadata: timestamp, question, answer, source, ip_address, user_agent, is_mobile, os, browser, country, accept_language, referer, session_id, history_length, model, cloudfront).",
	),
	mcp.WithString(
		"query",
		mcp.Required(),
		mcp.Description("SQL SELECT (or WITH ... SELECT) query to execute against mcp_ai_query_log. Must be a single statement; semicolons other than a trailing terminator are rejected."),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

// validateReadOnlyQuery enforces that a user-supplied DuckDB query is a single
// read-only statement. Returns the cleaned query (trailing semicolons stripped)
// and an error message suitable for returning to the MCP client. On success the
// error string is empty.
//
// The previous implementation only checked HasPrefix("SELECT"), which is
// trivially bypassed by "SELECT 1; DELETE FROM foo;" — DuckDB's Go driver can
// execute multi-statement payloads through Query(). This guard also accepts
// leading WITH (CTEs), since they're a common read-only shape.
func validateReadOnlyQuery(q string) (string, string) {
	cleaned := strings.TrimSpace(q)
	// Allow one or more trailing semicolons (common from copy/paste) then
	// re-trim whitespace that may have sat between them.
	cleaned = strings.TrimRight(cleaned, "; \t\r\n")
	if cleaned == "" {
		return "", "Empty query"
	}
	upper := strings.ToUpper(cleaned)
	if !strings.HasPrefix(upper, "SELECT") && !strings.HasPrefix(upper, "WITH") {
		return "", "Only SELECT or WITH … SELECT queries are allowed"
	}
	// After stripping the trailing terminator, no further semicolons are
	// permitted. This rejects stacked statements like "SELECT 1; DROP TABLE x"
	// as well as semicolons hidden in block comments.
	if strings.Contains(cleaned, ";") {
		return "", "Multi-statement queries are not allowed; remove inner ';' characters"
	}
	return cleaned, ""
}

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

	query, vErr := validateReadOnlyQuery(q)
	if vErr != "" {
		return mcp.NewToolResultText(vErr), nil
	}

	rows, err := duckDB.QueryContext(ctx, query)
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
