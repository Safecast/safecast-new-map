// MCP Server Integration for Unified Safecast Server
// Starts MCP server on separate port (default 3333)
// Requires: PostgreSQL (main DB), DuckDB (optional for analytics), ANTHROPIC_API_KEY (optional, for web chat)

package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	httpSwagger "github.com/swaggo/http-swagger"
	"safecast-new-map/cmd/unified-server/model-adapter"
)

var (
	mcpModelAdapter *modeladapter.Adapter
	mcpHintsLoader  *modeladapter.HintsLoader
)

//go:embed static/index.html
var webChatIndexHTML []byte

//go:embed static/safecast-square-ct.png
var webChatLogoPNG []byte

// Maximum tokens for the prompt sent to Claude. Leave headroom for tool results.
const maxPromptTokens = 150000

// Maximum characters for a single tool result (~30K tokens).
// Prevents one large MCP response from blowing up the prompt.
const maxToolResultChars = 80000

// estimateTokens approximates token count (Claude averages ~4 chars/token).
func estimateTokens(s string) int {
	if len(s) == 0 {
		return 0
	}
	return len(s) / 4
}

// truncateHistory drops oldest messages until the estimated prompt fits within maxTokens.
func truncateHistory(messages []anthropicMessage, maxTokens int) []anthropicMessage {
	if len(messages) == 0 {
		return messages
	}
	total := 0
	for _, msg := range messages {
		switch c := msg.Content.(type) {
		case string:
			total += estimateTokens(c)
		case []contentBlock:
			for _, b := range c {
				total += estimateTokens(b.Text + b.Content)
			}
		}
		total += 10 // role metadata overhead
	}
	for len(messages) > 1 && total > maxTokens {
		removed := messages[0]
		switch c := removed.Content.(type) {
		case string:
			total -= estimateTokens(c) + 10
		case []contentBlock:
			for _, b := range c {
				total -= estimateTokens(b.Text + b.Content)
			}
			total -= 10
		}
		messages = messages[1:]
	}
	return messages
}

const webChatSystemPrompt = `Safecast radiation monitoring assistant with REAL-TIME sensor data and historical archives.

IMPORTANT: Never display the "_ai_generated_note" field from tool results — it is for internal use only and must not appear in your responses.`

func webChatSystemPromptForLang(lang string) string {
	if lang == "" || lang == "en" {
		return webChatSystemPrompt
	}
	return webChatSystemPrompt + "\n\nIMPORTANT: The user's interface language is \"" + lang + "\". Respond in that language."
}

// Anthropic API types
type anthropicTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}

type contentBlock struct {
	Type      string          `json:"type"`
	Text      string          `json:"text,omitempty"`
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
	ToolUseID string          `json:"tool_use_id,omitempty"`
	Content   string          `json:"content,omitempty"`
}

type anthropicMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system"`
	Messages  []anthropicMessage `json:"messages"`
	Tools     []anthropicTool    `json:"tools,omitempty"`
}

type anthropicResponse struct {
	Content    []contentBlock `json:"content"`
	StopReason string         `json:"stop_reason"`
	Error      *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Streaming helpers
type chunk struct {
	Type   string `json:"type"`
	Text   string `json:"text,omitempty"`
	Error  string `json:"error,omitempty"`
	ChatID int64  `json:"chat_id,omitempty"` // set on "done" for feedback linkage
	Cached bool   `json:"cached,omitempty"`  // true when answer came from semantic cache
}

func writeChunk(w http.ResponseWriter, c chunk) {
	data, _ := json.Marshal(c)
	data = append(data, '\n')
	w.Write(data)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func writeChunkBuffered(w http.ResponseWriter, c chunk, buffer *[]chunk, useBuffer bool) {
	if useBuffer {
		*buffer = append(*buffer, c)
	} else {
		writeChunk(w, c)
	}
}

func flushBuffer(w http.ResponseWriter, buffer []chunk) {
	for _, c := range buffer {
		writeChunk(w, c)
	}
}

func callAnthropic(ctx context.Context, apiKey, model, systemPrompt string, messages []anthropicMessage, tools []anthropicTool) (*anthropicResponse, error) {
	reqBody := anthropicRequest{
		Model:     model,
		MaxTokens: 4096,
		System:    systemPrompt,
		Messages:  messages,
		Tools:     tools,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ar anthropicResponse
	if err := json.Unmarshal(raw, &ar); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if ar.Error != nil {
		return nil, fmt.Errorf("anthropic %s: %s", ar.Error.Type, ar.Error.Message)
	}
	return &ar, nil
}

func mcpToolsToAnthropic(tools []mcp.Tool) []anthropicTool {
	var out []anthropicTool
	for _, t := range tools {
		schema, _ := json.Marshal(t.InputSchema)
		out = append(out, anthropicTool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: json.RawMessage(schema),
		})
	}
	return out
}

func handleWebChat(mcpURL, apiKey, model string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		isCloudFront := r.Header.Get("CloudFront-Viewer-Country") != "" ||
			r.Header.Get("CloudFront-Forwarded-Proto") != "" ||
			r.Header.Get("X-Amz-Cf-Id") != ""

		w.Header().Set("Content-Type", "application/x-ndjson")
		if !isCloudFront {
			w.Header().Set("Transfer-Encoding", "chunked")
			w.Header().Set("X-Accel-Buffering", "no")
		}
		w.Header().Set("Cache-Control", "no-cache, no-store")

		var buffer []chunk
		ctx := r.Context()

		var chatReq struct {
			Message         string             `json:"message"`
			History         []anthropicMessage `json:"history,omitempty"`
			Source          string             `json:"source,omitempty"`
			Lang            string             `json:"lang,omitempty"`
			ClientTimestamp string             `json:"client_timestamp,omitempty"`
			TrackID         string             `json:"track_id,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&chatReq); err != nil || chatReq.Message == "" {
			w.WriteHeader(http.StatusBadRequest)
			writeChunkBuffered(w, chunk{Type: "error", Error: "invalid request: message required"}, &buffer, isCloudFront)
			if isCloudFront {
				flushBuffer(w, buffer)
			}
			return
		}

		// Default source if not provided by frontend
		source := chatReq.Source
		if source == "" {
			source = "web-chat"
		}
		// Capture request metadata now; the full row (with answer) is inserted after the AI responds
		chatReqRef := r
		chatQuestion := chatReq.Message
		chatSource := source
		chatModel := model
		chatHistory := len(chatReq.History)
		chatClientTS := chatReq.ClientTimestamp
		var answerText strings.Builder

		// Assign a stable ID for this exchange (used by both chat_questions and qa_embeddings).
		embeddingChatID := time.Now().UnixMilli() // UnixNano exceeds JS MAX_SAFE_INTEGER

		// --- Semantic cache + RAG layer (requires OPENAI_API_KEY) ---
		embedding, embErr := getEmbedding(ctx, chatReq.Message)
		if embErr != nil {
			log.Printf("embedding error (continuing without cache): %v", embErr)
		}

		if len(embedding) > 0 {
			// 1. Check semantic cache: high-similarity + positive feedback → return instantly.
			if cachedAnswer, _ := checkSemanticCache(embedding, chatReq.Message, chatReq.TrackID); cachedAnswer != "" {
				writeChunkBuffered(w, chunk{Type: "text", Text: cachedAnswer}, &buffer, isCloudFront)
				writeChunkBuffered(w, chunk{Type: "done", ChatID: embeddingChatID, Cached: true}, &buffer, isCloudFront)
				if isCloudFront {
					flushBuffer(w, buffer)
				}
				logChatQuestionWithAnswer(chatReqRef, chatQuestion, chatSource, chatModel, "", chatHistory, chatClientTS, cachedAnswer, embeddingChatID)
				return
			}
		}

		mc, err := mcpclient.NewStreamableHttpClient(mcpURL)
		if err != nil {
			writeChunkBuffered(w, chunk{Type: "error", Error: fmt.Sprintf("MCP connect: %v", err)}, &buffer, isCloudFront)
			if isCloudFront {
				flushBuffer(w, buffer)
			}
			return
		}
		defer mc.Close()

		if _, err := mc.Initialize(ctx, mcp.InitializeRequest{
			Params: mcp.InitializeParams{
				ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
				ClientInfo:      mcp.Implementation{Name: "safecast-web-chat", Version: "1.0.0"},
			},
		}); err != nil {
			writeChunkBuffered(w, chunk{Type: "error", Error: fmt.Sprintf("MCP init: %v", err)}, &buffer, isCloudFront)
			if isCloudFront {
				flushBuffer(w, buffer)
			}
			return
		}

		toolsResult, err := mc.ListTools(ctx, mcp.ListToolsRequest{})
		if err != nil {
			writeChunkBuffered(w, chunk{Type: "error", Error: fmt.Sprintf("list tools: %v", err)}, &buffer, isCloudFront)
			if isCloudFront {
				flushBuffer(w, buffer)
			}
			return
		}
		tools := mcpToolsToAnthropic(toolsResult.Tools)

		messages := chatReq.History
		if messages == nil {
			messages = []anthropicMessage{}
		}
		messages = append(messages, anthropicMessage{Role: "user", Content: chatReq.Message})
		messages = truncateHistory(messages, maxPromptTokens)

		// 2. Build RAG context from similar past Q&A + location knowledge.
		sysPrompt := webChatSystemPromptForLang(chatReq.Lang)
		if chatReq.TrackID != "" {
			sysPrompt = "The user is currently viewing track " + chatReq.TrackID + " on the Safecast map. When the user refers to 'this track', 'the track', or similar, they mean track " + chatReq.TrackID + ".\n\n" + sysPrompt
		}
		if len(embedding) > 0 {
			ragCtx := buildRAGContext(embedding)
			locKnowledge := getLocationKnowledge()
			sysPrompt = enrichSystemPrompt(sysPrompt, ragCtx, locKnowledge)
		}

		for {
			resp, err := callAnthropic(ctx, apiKey, model, sysPrompt, messages, tools)
			if err != nil {
				writeChunkBuffered(w, chunk{Type: "error", Error: err.Error()}, &buffer, isCloudFront)
				if isCloudFront {
					flushBuffer(w, buffer)
				}
				return
			}

			messages = append(messages, anthropicMessage{
				Role:    "assistant",
				Content: resp.Content,
			})

			var toolUses []contentBlock
			for _, block := range resp.Content {
				switch block.Type {
				case "text":
					answerText.WriteString(block.Text)
					writeChunkBuffered(w, chunk{Type: "text", Text: block.Text}, &buffer, isCloudFront)
				case "tool_use":
					toolUses = append(toolUses, block)
				}
			}

			if resp.StopReason == "end_turn" || len(toolUses) == 0 {
				break
			}

			var toolResults []contentBlock
			for _, tu := range toolUses {
				var args map[string]any
				_ = json.Unmarshal(tu.Input, &args)

				callReq := mcp.CallToolRequest{}
				callReq.Params.Name = tu.Name
				callReq.Params.Arguments = args

				var resultText string
				toolResult, err := mc.CallTool(ctx, callReq)
				if err != nil {
					resultText = fmt.Sprintf("tool error: %v", err)
				} else {
					for _, c := range toolResult.Content {
						if tc, ok := c.(mcp.TextContent); ok {
							resultText += tc.Text
						}
					}
				}
				if len(resultText) > maxToolResultChars {
					resultText = resultText[:maxToolResultChars] + "\n\n... [truncated — result too large. Ask the user to narrow their query or use a smaller limit.]"
				}

				toolResults = append(toolResults, contentBlock{
					Type:      "tool_result",
					ToolUseID: tu.ID,
					Content:   resultText,
				})
			}

			messages = append(messages, anthropicMessage{
				Role:    "user",
				Content: toolResults,
			})
		}

		finalAnswer := strings.TrimSpace(answerText.String())
		logChatQuestionWithAnswer(chatReqRef, chatQuestion, chatSource, chatModel, "", chatHistory, chatClientTS, finalAnswer, embeddingChatID)

		// 3. Async: store Q&A + embedding in semantic cache for future lookups.
		if len(embedding) > 0 && finalAnswer != "" {
			storeQAEmbeddingAsync(ctx, embeddingChatID, chatQuestion, finalAnswer, embedding)
		}

		writeChunkBuffered(w, chunk{Type: "done", ChatID: embeddingChatID}, &buffer, isCloudFront)
		if isCloudFront {
			flushBuffer(w, buffer)
		}
	}
}

// handleFeedback accepts a thumbs-up (+1) or thumbs-down (-1) for a chat response.
// The frontend sends: POST /api/feedback {"chat_id": <int>, "score": 1|-1}
func handleFeedback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			ChatID int64 `json:"chat_id"`
			Score  int   `json:"score"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ChatID == 0 {
			http.Error(w, "invalid request: chat_id required", http.StatusBadRequest)
			return
		}
		if err := RecordFeedback(req.ChatID, req.Score); err != nil {
			log.Printf("feedback error: %v", err)
			http.Error(w, "failed to record feedback", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}
}

// RegisterMCP starts the MCP server on a separate port (default 3333).
func RegisterMCP() {
	log.Println("DEBUG: safecast unified server with MCP integration")

	// Initialize DuckDB for analytics
	if err := initDuckDBAnalytics(); err != nil {
		log.Printf("Warning: DuckDB initialization failed: %v (analytics features disabled)", err)
	}

	// Initialize hints loader
	hintsDir := os.Getenv("MCP_HINTS_DIR")
	if hintsDir == "" {
		execPath, _ := os.Executable()
		hintsDir = filepath.Join(filepath.Dir(execPath), "hints")
	}

	mcpHintsLoader = modeladapter.NewHintsLoader(hintsDir)
	if err := mcpHintsLoader.Load(); err != nil {
		log.Printf("Warning: failed to load hints: %v (using default hints)", err)
	} else {
		log.Printf("Loaded hints for models: %v", mcpHintsLoader.GetAllModels())
	}

	mcpModelAdapter = modeladapter.NewAdapter()
	mcpModelAdapter.SetHintsLoader(mcpHintsLoader)

	// Create MCP server
	mcpServer := server.NewMCPServer("safecast-mcp", "1.0.0")

	if db != nil {
		log.Println("Using existing PostgreSQL connection for MCP")
	}
	if duckDB != nil {
		log.Println("Using existing DuckDB connection for MCP analytics")
	}

	// Register MCP tools
	mcpServer.AddTool(
		mcp.NewTool("ping", mcp.WithDescription("Health check tool")),
		instrumentMCP("ping", func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultText("pong"), nil
		}),
	)

	mcpServer.AddTool(queryRadiationToolDef, instrumentMCP("query_radiation", handleQueryRadiation))
	mcpServer.AddTool(searchAreaToolDef, instrumentMCP("search_area", handleSearchArea))
	mcpServer.AddTool(listTracksToolDef, instrumentMCP("list_tracks", handleListTracks))
	mcpServer.AddTool(getTrackToolDef, instrumentMCP("get_track", handleGetTrack))
	mcpServer.AddTool(deviceHistoryToolDef, instrumentMCP("device_history", handleDeviceHistory))
	mcpServer.AddTool(getSpectrumToolDef, instrumentMCP("get_spectrum", handleGetSpectrum))
	mcpServer.AddTool(listSpectraToolDef, instrumentMCP("list_spectra", handleListSpectra))
	mcpServer.AddTool(radiationInfoToolDef, instrumentMCP("radiation_info", handleRadiationInfo))
	mcpServer.AddTool(dbInfoToolDef, instrumentMCP("db_info", handleDBInfo))
	mcpServer.AddTool(listSensorsToolDef, instrumentMCP("list_sensors", handleListSensors))
	mcpServer.AddTool(sensorCurrentToolDef, instrumentMCP("sensor_current", handleSensorCurrent))
	mcpServer.AddTool(sensorHistoryToolDef, instrumentMCP("sensor_history", handleSensorHistory))
	mcpServer.AddTool(queryAnalyticsToolDef, instrumentMCP("query_analytics", handleQueryAnalytics))
	mcpServer.AddTool(radiationStatsToolDef, instrumentMCP("radiation_stats", handleRadiationStats))
	mcpServer.AddTool(queryDuckDBLogsToolDef, instrumentMCP("query_duckdb_logs", handleQueryDuckDBLogs))
	mcpServer.AddTool(queryExtremeReadingsToolDef, instrumentMCP("query_extreme_readings", handleQueryExtremeReadings))
	mcpServer.AddTool(topUploadersToolDef, instrumentMCP("top_uploaders", handleTopUploaders))
	mcpServer.AddTool(searchTracksLocationToolDef, instrumentMCP("search_tracks_by_location", handleSearchTracksByLocation))

	log.Println("MCP tools registered")

	// MCP port configuration
	mcpPort := os.Getenv("MCP_PORT")
	if mcpPort == "" {
		mcpPort = "3333"
	}

	baseURL := os.Getenv("MCP_BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%s", mcpPort)
	}

	sseServer := server.NewSSEServer(mcpServer,
		server.WithBaseURL(baseURL),
		server.WithStaticBasePath("/mcp"),
	)

	httpServer := server.NewStreamableHTTPServer(mcpServer,
		server.WithEndpointPath("/mcp-http"),
	)

	mux := http.NewServeMux()
	mux.Handle("/mcp-http", modeladapter.ModelDetectionMiddleware(httpServer))
	mux.Handle("/mcp/", modeladapter.ModelDetectionMiddleware(sseServer))

	// Register Swagger docs
	registerSwaggerDocs(mux)

	// Web Chat routes
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	model := os.Getenv("CLAUDE_MODEL")
	if model == "" {
		model = "claude-sonnet-4-5"
	}
	mcpURL := fmt.Sprintf("http://localhost:%s/mcp-http", mcpPort)

	feedbackHandler := handleFeedback()
	// Register feedback on both mux and main mux regardless of apiKey,
	// so the endpoint is always reachable even if chat is reconfigured.
	mux.HandleFunc("/api/feedback", feedbackHandler)
	http.HandleFunc("/api/feedback", feedbackHandler)

	// Track insights: register on main mux (port 8765) using Go 1.22 pattern routing.
	// The specific pattern "GET /api/track/{id}/insights" takes precedence over the
	// pkg/httpapi catch-all "/api/track/" handler.
	http.HandleFunc("GET /api/track/{id}/insights", trackInsightsHandler)

	if apiKey != "" {
		chatHandler := handleWebChat(mcpURL, apiKey, model)

		// Register on MCP mux (port 3333)
		mux.HandleFunc("/assistant/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(webChatIndexHTML)
		})
		mux.HandleFunc("/safecast-square-ct.png", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/png")
			w.Header().Set("Cache-Control", "public, max-age=86400")
			w.Write(webChatLogoPNG)
		})
		mux.HandleFunc("/chat", chatHandler)

		// Also register /chat on main map server (port 8765) so the
		// embedded widget can use a relative "/chat" URL without
		// cross-origin or CloudFront routing issues.
		http.HandleFunc("/chat", chatHandler)

		log.Printf("Web chat enabled at http://localhost:%s/assistant/ (model=%s)", mcpPort, model)
	} else {
		log.Println("Web chat disabled: ANTHROPIC_API_KEY not set")
	}

	log.Printf("MCP Server starting on port %s", mcpPort)
	log.Println("  SSE endpoint: /mcp/sse")
	log.Println("  Streamable HTTP endpoint: /mcp-http")
	log.Printf("  Hints directory: %s", hintsDir)
	log.Println("  REST API: /api/...")
	log.Println("  Swagger UI: /mcp-api/")
	if apiKey != "" {
		log.Printf("  Web Chat: http://localhost:%s/assistant/", mcpPort)
	}

	// Start MCP server on separate port
	go func() {
		listenAddr := ":" + mcpPort
		log.Printf("MCP goroutine: starting listener on %s", listenAddr)
		if err := http.ListenAndServe(listenAddr, mux); err != nil {
			log.Printf("ERROR: MCP server on port %s failed: %v", mcpPort, err)
		}
	}()
	log.Printf("MCP goroutine launched for port %s", mcpPort)
}

func instrumentMCP(
	name string,
	h func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error),
) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()

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

		res, err := h(ctx, req)

		if mcpModelAdapter != nil && res != nil {
			res = mcpModelAdapter.EnrichResult(ctx, res)
		}

		duration := time.Since(start)

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

		LogQueryAsync(name, args, resultCount, duration, "unified-server")

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

// registerSwaggerDocs registers the Swagger UI at /mcp-api/
func registerSwaggerDocs(mux *http.ServeMux) {
	mapBaseForDocs := strings.TrimSpace(os.Getenv("MAP_BASE_URL"))
	if mapBaseForDocs == "" {
		if strings.TrimSpace(*domain) != "" {
			mapBaseForDocs = "https://" + strings.TrimSpace(*domain)
		} else {
			mapBaseForDocs = fmt.Sprintf("http://localhost:%d", *port)
		}
	}
	mapDocsURL := strings.TrimRight(mapBaseForDocs, "/") + "/map-api/"
	mcpAPINavScript := fmt.Sprintf(`function() {
				document.title = 'Safecast MCP API Docs';
				const mapDocsURL = %q;

				// ── Favicons ──
				const link16 = document.createElement('link');
				link16.rel = 'icon'; link16.type = 'image/png'; link16.sizes = '16x16';
				link16.href = '/mcp-api/favicon-16x16.png';
				document.head.appendChild(link16);
				const link32 = document.createElement('link');
				link32.rel = 'icon'; link32.type = 'image/png'; link32.sizes = '32x32';
				link32.href = '/mcp-api/favicon-32x32.png';
				document.head.appendChild(link32);
				const linkICO = document.createElement('link');
				linkICO.rel = 'shortcut icon'; linkICO.href = '/mcp-api/favicon.ico';
				document.head.appendChild(linkICO);

				// ── Theme CSS ──
				const style = document.createElement('link');
				style.rel = 'stylesheet'; style.href = '/mcp-api/swagger-theme.css';
				document.head.appendChild(style);

				// ── Remove Swagger logo ──
				const swaggerLogo = document.querySelector('.topbar-wrapper .link');
				if (swaggerLogo) swaggerLogo.remove();
				document.querySelectorAll('.topbar-wrapper img').forEach(img => img.remove());

				// ── Inject "Switch to Map API" button into topbar ──
				const topbar = document.querySelector('.swagger-ui .topbar');
				if (topbar) {
					const existingBtn = document.getElementById('safecast-switch-btn');
					if (existingBtn) existingBtn.remove();
					const switchBtn = document.createElement('a');
					switchBtn.id = 'safecast-switch-btn';
					switchBtn.href = mapDocsURL;
					switchBtn.textContent = '\u2190 Switch to Map API';
					switchBtn.style.cssText = 'display:inline-block;margin-left:auto;margin-right:16px;padding:6px 14px;background:#085e58;color:#fff;border-radius:6px;font:600 13px/1.4 -apple-system,BlinkMacSystemFont,sans-serif;text-decoration:none;border:1px solid rgba(255,255,255,0.25);white-space:nowrap;';
					switchBtn.onmouseover = function() { this.style.background = '#064a45'; };
					switchBtn.onmouseout  = function() { this.style.background = '#085e58'; };
					const wrapper = topbar.querySelector('.topbar-wrapper');
					if (wrapper) {
						wrapper.style.display = 'flex';
						wrapper.style.alignItems = 'center';
						wrapper.style.width = '100%%';
						wrapper.appendChild(switchBtn);
					} else {
						topbar.appendChild(switchBtn);
					}
				}

				// ── Dark mode toggle ──
				const btn = document.createElement('button');
				btn.id = 'dark-mode-toggle';
				btn.textContent = '\u{1F319} Dark Mode';
				const isDark = localStorage.getItem('safecastMCPDarkMode') === 'true';
				if (isDark) { document.body.classList.add('dark-mode'); btn.textContent = '\u2600\uFE0F Light Mode'; }
				btn.onclick = function() {
					document.body.classList.toggle('dark-mode');
					const nowDark = document.body.classList.contains('dark-mode');
					btn.textContent = nowDark ? '\u2600\uFE0F Light Mode' : '\u{1F319} Dark Mode';
					localStorage.setItem('safecastMCPDarkMode', nowDark);
				};
				document.body.appendChild(btn);

				// ── Preamble section (inserted before #swagger-ui) ──
				const existing = document.getElementById('safecast-preamble');
				if (existing) existing.remove();
				const preamble = document.createElement('div');
				preamble.id = 'safecast-preamble';
				preamble.style.cssText = 'font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:#fff;border-bottom:3px solid #0d9488;padding:28px 32px 24px;line-height:1.6;';
				preamble.innerHTML = '<div style="max-width:900px;margin:0 auto;">' +
					'<h2 style="margin:0 0 6px;font-size:22px;color:#0f3d38;">Safecast MCP API</h2>' +
					'<p style="margin:0 0 16px;font-size:15px;color:#555;">' +
						'This API is designed for AI assistants and automated tools. It exposes the Safecast radiation dataset ' +
						'as a set of callable functions (MCP tools) and standard REST endpoints, so that language models and ' +
						'applications can query radiation data programmatically. If you want to ask questions about the data ' +
						'in plain English instead, try the <a href="/assistant/" style="color:#0d9488;">AI web chat</a>.' +
					'</p>' +
					'<details style="margin-bottom:16px;">' +
						'<summary style="cursor:pointer;font-weight:600;font-size:14px;color:#0f3d38;user-select:none;">For developers &mdash; transports &amp; tools overview</summary>' +
						'<div style="margin-top:10px;display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:10px;">' +
							'<div style="background:#f0fdfa;border:1px solid #99f6e4;border-radius:8px;padding:12px;">' +
								'<strong style="color:#0f3d38;">MCP Streamable HTTP</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">Primary MCP transport at <code>/mcp-http</code>. Compatible with Claude, GPT, and most MCP clients.</p>' +
							'</div>' +
							'<div style="background:#f0fdfa;border:1px solid #99f6e4;border-radius:8px;padding:12px;">' +
								'<strong style="color:#0f3d38;">MCP SSE</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">Server-Sent Events transport at <code>/mcp/sse</code> for streaming-capable clients.</p>' +
							'</div>' +
							'<div style="background:#f0fdfa;border:1px solid #99f6e4;border-radius:8px;padding:12px;">' +
								'<strong style="color:#0f3d38;">REST Endpoints</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">All MCP tools are also accessible as plain HTTP GET/POST calls under <code>/api/</code>.</p>' +
							'</div>' +
							'<div style="background:#f0fdfa;border:1px solid #99f6e4;border-radius:8px;padding:12px;">' +
								'<strong style="color:#0f3d38;">AI Web Chat</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">Human-friendly interface at <code>/assistant/</code> for querying data with natural language.</p>' +
							'</div>' +
						'</div>' +
						'<p style="margin:12px 0 0;font-size:13px;color:#777;">No authentication required. All data is CC0 licensed.</p>' +
					'</details>' +
					'<div style="display:flex;gap:10px;flex-wrap:wrap;align-items:center;">' +
						'<a href="/assistant/" style="display:inline-block;padding:7px 16px;background:#0f3d38;color:#fff;border-radius:6px;font:600 13px/1.4 sans-serif;text-decoration:none;">Open AI Chat</a>' +
						'<a href="' + mapDocsURL + '" style="display:inline-block;padding:7px 16px;background:#1a3a5c;color:#fff;border-radius:6px;font:600 13px/1.4 sans-serif;text-decoration:none;">\u2190 Switch to Map API</a>' +
					'</div>' +
				'</div>';
				const swaggerUIEl = document.getElementById('swagger-ui');
				if (swaggerUIEl) {
					document.body.insertBefore(preamble, swaggerUIEl);
				} else {
					document.body.prepend(preamble);
				}
			}`, mapDocsURL)
	mux.HandleFunc("/mcp-api/favicon.ico", serveFavicon)
	mux.HandleFunc("/mcp-api/favicon-16x16.png", serveFavicon16)
	mux.HandleFunc("/mcp-api/favicon-32x32.png", serveFavicon32)
	mux.HandleFunc("/mcp-api/swagger-theme.css", serveSwaggerTheme)
	mux.Handle("/mcp-api/", httpSwagger.Handler(
		httpSwagger.URL("/mcp-api/doc.json"),
		httpSwagger.UIConfig(map[string]string{
			"onComplete": mcpAPINavScript,
		}),
	))
}
