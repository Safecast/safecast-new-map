package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/xuri/excelize/v2"
)

//go:embed index.html
var indexHTML []byte

//go:embed safecast-square-ct.png
var logoPNG []byte

// Maximum tokens for the prompt (input to Claude). Leave room for tool results.
const maxPromptTokens = 150000

// Maximum characters for a single tool result (~30K tokens).
// Prevents a single large MCP response from blowing up the prompt.
const maxToolResultChars = 120000

// estimateTokens roughly estimates the number of tokens in a string.
// Claude uses ~4 characters per token on average for English text.
func estimateTokens(s string) int {
	if len(s) == 0 {
		return 0
	}
	return len(s) / 4
}

// truncateHistory removes older messages from history to stay under token limit.
// It keeps the most recent messages and always includes the system prompt.
func truncateHistory(messages []anthropicMessage, maxTokens int) []anthropicMessage {
	if len(messages) == 0 {
		return messages
	}

	// Estimate current token count
	totalTokens := 0
	for _, msg := range messages {
		switch content := msg.Content.(type) {
		case string:
			totalTokens += estimateTokens(content)
		case []contentBlock:
			for _, block := range content {
				totalTokens += estimateTokens(block.Text)
			}
		}
		// Add ~10 tokens per message for role metadata
		totalTokens += 10
	}

	// If under limit, return as-is
	if totalTokens <= maxTokens {
		return messages
	}

	// Remove messages from the beginning until we're under the limit
	// Keep at least the last message (current user query)
	for len(messages) > 1 && totalTokens > maxTokens {
		// Remove the oldest message
		removed := messages[0]
		switch content := removed.Content.(type) {
		case string:
			totalTokens -= estimateTokens(content) + 10
		case []contentBlock:
			for _, block := range content {
				totalTokens -= estimateTokens(block.Text)
			}
			totalTokens -= 10
		}
		messages = messages[1:]
	}

	// If still over limit, truncate the content of the first message
	if len(messages) > 0 && totalTokens > maxTokens {
		switch content := messages[0].Content.(type) {
		case string:
			// Truncate string content
			maxLen := (maxTokens * 4)
			if len(content) > maxLen {
				messages[0].Content = content[:maxLen] + "... [truncated due to length]"
			}
		case []contentBlock:
			// Truncate text blocks
			for i := range content {
				if totalTokens <= maxTokens {
					break
				}
				maxLen := (maxTokens * 4)
				if len(content[i].Text) > maxLen {
					content[i].Text = content[i].Text[:maxLen] + "... [truncated due to length]"
					totalTokens = estimateTokens(content[i].Text)
				}
			}
			messages[0].Content = content
		}
	}

	return messages
}

const systemPrompt = `You are a data-aware AI assistant embedded inside a radiation map interface (Safecast-like system). Your role is not just to answer questions, but to help users explore, filter, and export structured datasets efficiently.

You must follow these rules strictly:

-------------------------------------
1. UNDERSTAND USER INTENT
-------------------------------------

Classify every query into one of three types:

A) Informational (e.g., “Is radiation high in Tokyo?”)
→ Respond normally with explanation

B) Data preview (e.g., “show radiation readings in Tokyo last 24h”)
→ Fetch structured data and return:
   - a preview (max 20–50 rows)
   - dataset metadata (total rows, time range, region)

C) Data export (e.g., “download csv”, “export full dataset”)
→ Skip preview and trigger export mode

If unclear, default to B (preview mode).

-------------------------------------
2. NEVER DUMP FULL DATA IN CHAT
-------------------------------------

You MUST NOT output large datasets directly in chat. Always limit preview to 50 rows max and indicate total dataset size.

-------------------------------------
3. ALWAYS OFFER EXPORT WHEN RELEVANT
-------------------------------------

If:
- dataset > 50 rows OR
- user asks for data OR
- query involves filters (time, location, sensors)

Then include: EXPORT_AVAILABLE = true

-------------------------------------
4. EXPORT CONFIGURATION MODEL
-------------------------------------

When export is triggered, construct a structured export request:

{
  "format": "csv | xlsx | json",
  "limit": "preview | sample | full | custom",
  "time_range": "1h | 24h | 7d | custom",
  "bbox": [minLat, minLon, maxLat, maxLon],
  "region": "...",
  "columns": ["device_id", "cpm", "usvh", "timestamp", "lat", "lon"],
  "aggregation": "none | hourly | daily"
}

Defaults:
- format: csv
- limit: full
- columns: all
- region: current map view

-------------------------------------
5. HANDLE LARGE DATASETS INTELLIGENTLY
-------------------------------------

Estimate dataset size before exporting.
- rows < 5,000 → allow instant export
- 5,000–50,000 → allow but warn briefly
- > 50,000 → mark as ASYNC export
- > 200,000 → require user confirmation or suggest filters

For large datasets, respond: “This dataset is large. I’ll prepare it for download.”

-------------------------------------
6. RESPONSE FORMAT (STRICT)
-------------------------------------

When returning data preview, use this JSON format (start of block):
{
  "type": "data_preview",
  "summary": {
    "rows_total": 12453,
    "rows_shown": 50,
    "time_range": "last 24h",
    "region": "Tokyo"
  },
  "table": [...50 rows max...],
  "export_available": true,
  "suggested_export": {
    "format": "csv",
    "limit": "full"
  }
}

When triggering export, use this JSON format (start of block):
{
  "type": "export_request",
  "export_config": {...},
  "estimated_rows": 12453,
  "mode": "sync | async"
}

-------------------------------------
7. BEHAVIORAL RULES
-------------------------------------

- Prefer structured JSON over plain text when handling data
- Always respect map bounds if available
- Never hallucinate data — rely on backend queries
- Be concise, functional, and tool-like
- Optimize for data scientists, not casual users

-------------------------------------
8. TOOL SELECTION (LEGACY)
-------------------------------------
- Current/live data: sensor_current
- Time-series: sensor_history
- Extreme readings: query_extreme_readings
- Statistics: radiation_stats
- Historical surveys: query_radiation, search_area, list_tracks
- CPM -> µSv/h: multiply by ~0.0069 (LND 7318)
- Timestamp: ALWAYS display in UTC — convert from any timezone, format as "2026-03-03 22:14 UTC"

Concise and tool-like. Ask for clarification if location unclear.`


// ── Anthropic API types ────────────────────────────────────────────────────

type anthropicTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}

// contentBlock covers all content block variants we care about.
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
	Content interface{} `json:"content"` // string or []contentBlock
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

// ── Streaming helpers (chunked HTTP / NDJSON) ──────────────────────────────

type chunk struct {
	Type        string          `json:"type"`
	Text        string          `json:"text,omitempty"`
	Error       string          `json:"error,omitempty"`
	DataPreview json.RawMessage `json:"data_preview,omitempty"`
}

// extractDataPreview looks for a data_preview JSON object in the AI response text.
// Returns the raw JSON bytes if found, nil otherwise.
func extractDataPreview(text string) json.RawMessage {
	if !strings.Contains(text, `"data_preview"`) {
		return nil
	}
	marker := strings.Index(text, `"type":"data_preview"`)
	if marker == -1 {
		marker = strings.Index(text, `"type": "data_preview"`)
	}
	if marker == -1 {
		return nil
	}

	jsonObjectEnd := func(start int) int {
		depth := 0
		inStr := false
		escape := false
		for i := start; i < len(text); i++ {
			c := text[i]
			if escape {
				escape = false
				continue
			}
			if c == '\\' && inStr {
				escape = true
				continue
			}
			if c == '"' {
				inStr = !inStr
				continue
			}
			if inStr {
				continue
			}
			if c == '{' {
				depth++
			} else if c == '}' {
				depth--
				if depth == 0 {
					return i
				}
			}
		}
		return -1
	}

	for i := 0; i < marker; i++ {
		if text[i] != '{' {
			continue
		}
		objEnd := jsonObjectEnd(i)
		if objEnd == -1 || objEnd < marker {
			continue
		}
		candidate := text[i : objEnd+1]
		var obj map[string]json.RawMessage
		if err := json.Unmarshal([]byte(candidate), &obj); err != nil {
			continue
		}
		var typeVal string
		if err := json.Unmarshal(obj["type"], &typeVal); err != nil || typeVal != "data_preview" {
			continue
		}
		return json.RawMessage(candidate)
	}
	return nil
}

func writeChunk(w http.ResponseWriter, c chunk) {
	data, _ := json.Marshal(c)
	data = append(data, '\n')
	w.Write(data)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// writeChunkBuffered either buffers the chunk or writes immediately
func writeChunkBuffered(w http.ResponseWriter, c chunk, buffer *[]chunk, useBuffer bool) {
	if useBuffer {
		*buffer = append(*buffer, c)
	} else {
		writeChunk(w, c)
	}
}

// flushBuffer writes all buffered chunks at once
func flushBuffer(w http.ResponseWriter, buffer []chunk) {
	for _, c := range buffer {
		writeChunk(w, c)
	}
}

// ── Anthropic call ─────────────────────────────────────────────────────────

func callAnthropic(ctx context.Context, apiKey, model string, messages []anthropicMessage, tools []anthropicTool) (*anthropicResponse, error) {
	// Truncate history if it exceeds the token limit
	messages = truncateHistory(messages, maxPromptTokens)

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

// ── MCP tool conversion ────────────────────────────────────────────────────

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

// ── Chat handler ───────────────────────────────────────────────────────────

func handleChat(mcpURL, apiKey, model string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS preflight
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Detect if request comes through CloudFront
		// CloudFront adds these headers: CloudFront-Viewer-Country, CloudFront-Forwarded-Proto, etc.
		isCloudfFront := r.Header.Get("CloudFront-Viewer-Country") != "" ||
			r.Header.Get("CloudFront-Forwarded-Proto") != "" ||
			r.Header.Get("X-Amz-Cf-Id") != ""

		// Debug logging
		log.Printf("Chat request: CloudFront=%v, Headers: CF-Country=%q, CF-Proto=%q, X-Amz-Cf-Id=%q",
			isCloudfFront,
			r.Header.Get("CloudFront-Viewer-Country"),
			r.Header.Get("CloudFront-Forwarded-Proto"),
			r.Header.Get("X-Amz-Cf-Id"))

		// Chunked HTTP streaming — NDJSON, one JSON object per line, flushed immediately.
		// CloudFront buffers responses, so we collect chunks and send all at once.
		w.Header().Set("Content-Type", "application/x-ndjson")
		if !isCloudfFront {
			w.Header().Set("Transfer-Encoding", "chunked")
			w.Header().Set("X-Accel-Buffering", "no") // nginx: don't buffer
		}
		w.Header().Set("Cache-Control", "no-cache, no-store")

		// Buffer for CloudFront requests
		var buffer []chunk

		ctx := r.Context()

		var chatReq struct {
			Message string              `json:"message"`
			History []anthropicMessage `json:"history,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&chatReq); err != nil || chatReq.Message == "" {
			w.WriteHeader(http.StatusBadRequest)
			writeChunkBuffered(w, chunk{Type: "error", Error: "invalid request: message required"}, &buffer, isCloudfFront)
			if isCloudfFront {
				flushBuffer(w, buffer)
			}
			return
		}

		// ── Connect to MCP server ──────────────────────────────────────────
		var tools []anthropicTool
		mc, err := mcpclient.NewStreamableHttpClient(mcpURL)
		if err == nil {
			defer mc.Close()
			if _, errInit := mc.Initialize(ctx, mcp.InitializeRequest{
				Params: mcp.InitializeParams{
					ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
					ClientInfo:      mcp.Implementation{Name: "safecast-web-chat", Version: "1.0.0"},
				},
			}); errInit == nil {
				toolsResult, errList := mc.ListTools(ctx, mcp.ListToolsRequest{})
				if errList == nil {
					tools = mcpToolsToAnthropic(toolsResult.Tools)
				}
			}
		}

		// ── Agentic loop ───────────────────────────────────────────────────
		// Start with conversation history (if provided) and append new user message
		messages := chatReq.History
		if messages == nil {
			messages = []anthropicMessage{}
		}
		messages = append(messages, anthropicMessage{Role: "user", Content: chatReq.Message})

		for {
			resp, err := callAnthropic(ctx, apiKey, model, messages, tools)
			if err != nil {
				writeChunkBuffered(w, chunk{Type: "error", Error: err.Error()}, &buffer, isCloudfFront)
				if isCloudfFront {
					flushBuffer(w, buffer)
				}
				return
			}

			messages = append(messages, anthropicMessage{
				Role:    "assistant",
				Content: resp.Content,
			})

			var toolUses []contentBlock
			var textBlocks []string
			for _, block := range resp.Content {
				switch block.Type {
				case "text":
					textBlocks = append(textBlocks, block.Text)
				case "tool_use":
					toolUses = append(toolUses, block)
				}
			}

			isFinal := resp.StopReason == "end_turn" || len(toolUses) == 0
			if isFinal {
				// On the final turn, check if the combined text contains a data_preview
				// JSON blob. If so, send a structured data_preview event so the client
				// renders a table. Otherwise stream the text blocks as-is.
				fullText := strings.Join(textBlocks, "")
				if dp := extractDataPreview(fullText); dp != nil {
					writeChunkBuffered(w, chunk{Type: "data_preview", DataPreview: dp}, &buffer, isCloudfFront)
				} else {
					for _, t := range textBlocks {
						writeChunkBuffered(w, chunk{Type: "text", Text: t}, &buffer, isCloudfFront)
					}
				}
				break
			}
			// Non-final iteration (tool calls pending): stream text blocks normally
			for _, t := range textBlocks {
				writeChunkBuffered(w, chunk{Type: "text", Text: t}, &buffer, isCloudfFront)
			}

			// ── Execute tool calls via MCP ─────────────────────────────────
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


				// Truncate oversized tool results to prevent exceeding API limits
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

		// Send final "done" chunk
		writeChunkBuffered(w, chunk{Type: "done"}, &buffer, isCloudfFront)

		// For CloudFront requests, flush all buffered chunks at once
		if isCloudfFront {
			flushBuffer(w, buffer)
		}
	}
}

// ── Export handler ─────────────────────────────────────────────────────────

func handleExport(mcpURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var exportReq struct {
			Format        string    `json:"format"`
			Limit         string    `json:"limit"`
			TimeRange     string    `json:"time_range"`
			Bbox          []float64 `json:"bbox"`
			Region        string    `json:"region"`
			Columns       []string  `json:"columns"`
			Aggregation   string    `json:"aggregation"`
			SuggestedTool string    `json:"suggested_tool"`
			Lat           float64   `json:"lat"`
			Lon           float64   `json:"lon"`
			RadiusM       float64   `json:"radius_m"`
		}

		if err := json.NewDecoder(r.Body).Decode(&exportReq); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		var data map[string]any

		if exportReq.Limit == "testmode" {
			data = map[string]any{
				"measurements": []any{
					map[string]any{"Location": "Osaka", "Radiation_uSv": 0.05, "device_id": "80242", "lat": 34.6937, "lon": 135.5023},
					map[string]any{"Location": "Tokyo", "Radiation_uSv": 0.04, "device_id": "1002", "lat": 35.6895, "lon": 139.6917},
					map[string]any{"Location": "Kyoto", "Radiation_uSv": 0.06, "device_id": "5003", "lat": 35.0116, "lon": 135.7681},
				},
			}
		} else {
			ctx := r.Context()
			mc, err := mcpclient.NewStreamableHttpClient(mcpURL)
			if err != nil {
				// Fallback to test data if database isn't running locally
				data = map[string]any{
					"measurements": []any{
						map[string]any{"Location": "Shibuya", "Radiation_uSv": 0.197, "lat": 35.6762, "lon": 139.6503},
						map[string]any{"Location": "Chiyoda", "Radiation_uSv": 0.215, "lat": 35.6895, "lon": 139.6917},
						map[string]any{"Location": "Minato", "Radiation_uSv": 0.185, "lat": 35.7314, "lon": 139.767},
					},
				}
			} else {
				defer mc.Close()

				if _, err := mc.Initialize(ctx, mcp.InitializeRequest{
				Params: mcp.InitializeParams{
					ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
					ClientInfo:      mcp.Implementation{Name: "safecast-export", Version: "1.0.0"},
				},
			}); err != nil {
				http.Error(w, "MCP init error", http.StatusInternalServerError)
				return
			}

			toolName := "query_radiation"
			if exportReq.SuggestedTool != "" {
				toolName = exportReq.SuggestedTool
			}

			limit := 10000
			if exportReq.Limit == "full" {
				limit = 100000
			}

			args := map[string]any{"limit": limit}
			if len(exportReq.Bbox) == 4 {
				args["min_lat"] = exportReq.Bbox[0]
				args["min_lon"] = exportReq.Bbox[1]
				args["max_lat"] = exportReq.Bbox[2]
				args["max_lon"] = exportReq.Bbox[3]
				args["lat"] = (exportReq.Bbox[0] + exportReq.Bbox[2]) / 2
				args["lon"] = (exportReq.Bbox[1] + exportReq.Bbox[3]) / 2
				args["radius_m"] = 50000
			} else if exportReq.Lat != 0 || exportReq.Lon != 0 {
				args["lat"] = exportReq.Lat
				args["lon"] = exportReq.Lon
				radius := exportReq.RadiusM
				if radius == 0 {
					radius = 50000
				}
				args["radius_m"] = radius
			}

			callReq := mcp.CallToolRequest{}
			callReq.Params.Name = toolName
			callReq.Params.Arguments = args

			toolResult, err := mc.CallTool(ctx, callReq)
			if err != nil {
				http.Error(w, "Tool call failed: "+err.Error(), http.StatusInternalServerError)
				return
			}

			var resultText string
			for _, c := range toolResult.Content {
				if tc, ok := c.(mcp.TextContent); ok {
					resultText += tc.Text
				}
			}

			if err := json.Unmarshal([]byte(resultText), &data); err != nil {
				http.Error(w, "Parse tool result failed", http.StatusInternalServerError)
				return
			}
			}
		}

		switch exportReq.Format {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Disposition", "attachment; filename=safecast_export.json")
			json.NewEncoder(w).Encode(data)
		case "excel", "xlsx":
			w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			w.Header().Set("Content-Disposition", "attachment; filename=safecast_export.xlsx")

			var rows []any
			if m, ok := data["measurements"].([]any); ok {
				rows = m
			} else if r, ok := data["readings"].([]any); ok {
				rows = r
			}

			if len(rows) > 0 {
				f := excelize.NewFile()
				defer func() {
					if err := f.Close(); err != nil {
						fmt.Println(err)
					}
				}()
				
				f.SetSheetName("Sheet1", "Data")
				
				firstRow, _ := rows[0].(map[string]any)
				var headers []string
				for k := range firstRow {
					headers = append(headers, k)
				}
				sort.Strings(headers)
				
				// Write Headers
				for i, h := range headers {
					cell, _ := excelize.CoordinatesToCellName(i+1, 1)
					f.SetCellValue("Data", cell, h)
				}

				// Write Rows
				for rowIdx, r := range rows {
					rowMap, _ := r.(map[string]any)
					for colIdx, h := range headers {
						val := rowMap[h]
						cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
						if m, ok := val.(map[string]any); ok {
							b, _ := json.Marshal(m)
							f.SetCellValue("Data", cell, string(b))
						} else {
							f.SetCellValue("Data", cell, val)
						}
					}
				}
				f.Write(w)
			} else {
				// Empty sheet
				f := excelize.NewFile()
				f.Write(w)
			}
		case "csv":
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", "attachment; filename=safecast_export.csv")
			
			var rows []any
			if m, ok := data["measurements"].([]any); ok {
				rows = m
			} else if r, ok := data["readings"].([]any); ok {
				rows = r
			}

			if len(rows) > 0 {
				writer := csv.NewWriter(w)
				firstRow, _ := rows[0].(map[string]any)
				var headers []string
				for k := range firstRow {
					headers = append(headers, k)
				}
				sort.Strings(headers)
				writer.Write(headers)

				for _, r := range rows {
					row, _ := r.(map[string]any)
					var line []string
					for _, h := range headers {
						val := row[h]
						if m, ok := val.(map[string]any); ok {
							b, _ := json.Marshal(m)
							line = append(line, string(b))
						} else {
							line = append(line, fmt.Sprintf("%v", val))
						}
					}
					writer.Write(line)
				}
				writer.Flush()
			} else {
				fmt.Fprint(w, "No data available")
			}
		default:
			http.Error(w, "Unsupported format", http.StatusBadRequest)
		}
	}
}

// ── Main ───────────────────────────────────────────────────────────────────


func main() {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: ANTHROPIC_API_KEY is not set. Chat features will error.")
	}
	model := os.Getenv("CLAUDE_MODEL")
	if model == "" {
		model = "claude-haiku-4-5-20251001"
	}
	mcpURL := os.Getenv("MCP_URL")
	if mcpURL == "" {
		mcpURL = "http://localhost:3333/mcp-http"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3334"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})
	http.HandleFunc("/safecast-square-ct.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(logoPNG)
	})
	http.HandleFunc("/chat", handleChat(mcpURL, apiKey, model))
	http.HandleFunc("/export", handleExport(mcpURL))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	log.Printf("Safecast web-chat on :%s  MCP→%s  model=%s", port, mcpURL, model)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
