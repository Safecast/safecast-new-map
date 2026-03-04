package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	log.Println("DEBUG: safecast MCP server binary version 2026-02-18-1")
	// Create MCP server
	mcpServer := server.NewMCPServer(
		"safecast-mcp",
		"1.0.0",
	)

	// Initialize database connection
	if os.Getenv("DATABASE_URL") != "" {
		if err := initDB(); err != nil {
			log.Printf("Warning: database connection failed: %v (using REST API fallback)", err)
		} else {
			log.Println("Connected to PostgreSQL database")
		}
	} else {
		log.Println("No DATABASE_URL set, using REST API only")
	}

	// Initialize DuckDB Analytics
	if err := initDuckDB(); err != nil {
		log.Printf("Warning: failed to initialize DuckDB: %v (analytics features disabled)", err)
	} else {
		log.Println("Initialized DuckDB analytics engine")
	}

	// Register tools
	mcpServer.AddTool(
		mcp.NewTool("ping",
			mcp.WithDescription("Health check tool"),
		),
		instrument("ping", pingHandler),
	)

	mcpServer.AddTool(queryRadiationToolDef, instrument("query_radiation", handleQueryRadiation))
	mcpServer.AddTool(searchAreaToolDef, instrument("search_area", handleSearchArea))
	mcpServer.AddTool(listTracksToolDef, instrument("list_tracks", handleListTracks))
	mcpServer.AddTool(getTrackToolDef, instrument("get_track", handleGetTrack))
	mcpServer.AddTool(deviceHistoryToolDef, instrument("device_history", handleDeviceHistory))
	mcpServer.AddTool(getSpectrumToolDef, instrument("get_spectrum", handleGetSpectrum))
	mcpServer.AddTool(listSpectraToolDef, instrument("list_spectra", handleListSpectra))
	mcpServer.AddTool(radiationInfoToolDef, instrument("radiation_info", handleRadiationInfo))
	mcpServer.AddTool(dbInfoToolDef, instrument("db_info", handleDBInfo))
	mcpServer.AddTool(listSensorsToolDef, instrument("list_sensors", handleListSensors))
	mcpServer.AddTool(sensorCurrentToolDef, instrument("sensor_current", handleSensorCurrent))
	mcpServer.AddTool(sensorHistoryToolDef, instrument("sensor_history", handleSensorHistory))
	mcpServer.AddTool(queryAnalyticsToolDef, instrument("query_analytics", handleQueryAnalytics))
	mcpServer.AddTool(radiationStatsToolDef, instrument("radiation_stats", handleRadiationStats))
	mcpServer.AddTool(queryDuckDBLogsToolDef, instrument("query_duckdb_logs", handleQueryDuckDBLogs))
	mcpServer.AddTool(queryExtremeReadingsToolDef, instrument("query_extreme_readings", handleQueryExtremeReadings))
	mcpServer.AddTool(topUploadersToolDef, instrument("top_uploaders", handleTopUploaders))
	mcpServer.AddTool(searchTracksLocationToolDef, instrument("search_tracks_by_location", handleSearchTracksByLocation))

	// 🚨 TRANSPORT SWITCH
	if os.Getenv("MCP_TRANSPORT") == "stdio" {

		log.Println("Starting MCP server in stdio mode (Claude Desktop)")

		stdioServer := server.NewStdioServer(mcpServer)

		err := stdioServer.Listen(
			context.Background(),
			os.Stdin,
			os.Stdout,
		)

		if err != nil {
			log.Fatal(err)
		}

		return
	}

	// Default: HTTP mode (production)

	baseURL := os.Getenv("MCP_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3333"
	}

	sseServer := server.NewSSEServer(mcpServer,
		server.WithBaseURL(baseURL),
		server.WithStaticBasePath("/mcp"),
	)

	httpServer := server.NewStreamableHTTPServer(mcpServer,
		server.WithEndpointPath("/mcp-http"),
	)

	mux := http.NewServeMux()
	mux.Handle("/mcp-http", httpServer)
	mux.Handle("/mcp/", sseServer) // SSE server handles /mcp/sse and /mcp/message

	rest := &RESTHandler{}
	rest.Register(mux)

	port := os.Getenv("MCP_PORT")
	if port == "" {
		port = "3333"
	}

	listenAddr := ":" + port

	log.Printf("Starting MCP server on %s", listenAddr)
	log.Println("  SSE endpoint: /mcp/sse")
	log.Println("  Streamable HTTP endpoint: /mcp-http")

	log.Println("  REST API: /api/...")
	log.Println("  Swagger UI: /docs/")

	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatal(err)
	}
	}

// pingHandler is the health check tool implementation.
func pingHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("pong"), nil
}

func instrument(
	name string,
	h func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error),
) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		start := time.Now()

		// Extract user info from MCP request arguments
		userID := ""
		userEmail := ""

		if req.Params.Arguments != nil {

			if args, ok := req.Params.Arguments.(map[string]any); ok {

				if v, ok := args["user_id"].(string); ok {
					userID = v
				}

				if v, ok := args["user_email"].(string); ok {
					userEmail = v
				}
			}
		}

		// Execute tool
		res, err := h(ctx, req)

		duration := time.Since(start)

		// Existing DuckDB analytics log
		resultCount := 0
		if res != nil {
			resultCount = len(res.Content)
		}

		args := map[string]any{}
		if req.Params.Arguments != nil {
			if argsMap, ok := req.Params.Arguments.(map[string]any); ok {
				args = argsMap
			}
		}

		LogQueryAsync(name, args, resultCount, duration, "claude-client")

		logAISessionWithUser(
			name,
			"",
			duration.Milliseconds(),
			err,
			userID,
			userEmail,
		)

		return res, err
	}
}