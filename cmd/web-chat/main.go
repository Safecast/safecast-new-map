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

	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

//go:embed index.html
var indexHTML []byte

//go:embed safecast-square-ct.png
var logoPNG []byte

const systemPrompt = `Safecast radiation monitoring assistant with REAL-TIME sensor data and historical archives.

**Tool Selection**
- Current/live data: sensor_current (returns actual CPM readings), list_sensors (metadata/discovery only)
- Time-series from fixed sensors: sensor_history
- Extreme readings with locations: query_extreme_readings
- Statistics: radiation_stats
- Historical surveys: query_radiation, search_area, list_tracks
- NEVER use query_radiation for current data (historical only)
- NEVER use radiation_stats for specific extreme location queries
- NEVER use list_sensors when the user wants radiation readings — use sensor_current instead

**Device type names** (exact values in the database):
- bGeigieZen → "geigiecast-zen" (IDs like geigiecast-zen:65002)
- bGeigie → "geigiecast" (IDs like geigiecast:62007) — MOBILE only
- Pointcast → "pointcast" (IDs like pointcast:10042)
- Solarcast → "solarcast"
- Notehub/Radnote/Blues → "notehub" (IDs like note:dev:867648049123019)
- nGeigie → "ngeigie" (IDs like ngeigie:101)
- Direct TCP → "device-tcp" (IDs like safecast:3474557222)

**Data Types**
- Real-time fixed stations: geigiecast-zen, pointcast, solarcast, notehub, ngeigie, device-tcp → sensor_current or sensor_history
- Mobile surveys only: geigiecast → query_radiation, list_tracks, device_history
- CPM → µSv/h: multiply by ~0.0069 (LND 7318)
- NEVER use device_history for any fixed sensor type — use sensor_current instead

**"Latest readings" at a location: ALWAYS do BOTH steps**
1. Call query_radiation (historical mobile data)
2. Call sensor_current with geographic bounds (fixed real-time sensors)
Report both results. Only say "no real-time data" after sensor_current returns empty results.

**Looking up a specific fixed sensor by device ID** (any non-geigiecast type):
- Use sensor_current with device_id parameter, NOT device_history
- Note: notehub device IDs contain colons, e.g. "note:dev:867648049123019" — pass the full string as device_id
- device_history is ONLY for mobile bGeigie (geigiecast) devices

**Radius Selection** (query_radiation, sensor_current):
Address: 1000-2000m | District: 5000-10000m | Village/Town: 25-50km | City: 50km | Metro: 75-100km
When in doubt, use a LARGER radius — it is better to return too many results than to miss nearby sensors due to geocoding imprecision. Always state radius used.

**Formatting**
- Hide "_ai_generated_note" field (internal use only)
- **CRITICAL: ALL devices/coords MUST be clickable map links:**
  * Devices: [pointcast:10042](https://simplemap.safecast.org/?lat=LAT&lon=LON&zoom=15)
  * Tracks: [track_id](https://simplemap.safecast.org/?lat=LAT&lon=LON&zoom=12)
  * Coords: [37.72°N, 140.48°E](https://simplemap.safecast.org/?lat=37.72&lon=140.48&zoom=15)
  * NEVER plain device names or "Visit: https://..." text
- Sensor/track data: ALWAYS use markdown tables (not lists)
- Table columns: Device ID, Type, Location, Reading, Timestamp
- Timestamp: ALWAYS display in UTC — convert from any timezone, format as "2026-03-03 22:14 UTC"
- Concise coords: "37.48°N, 140.48°E"

Be concise. Ask for clarification if location unclear.`

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
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	Error string `json:"error,omitempty"`
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
		mc, err := mcpclient.NewStreamableHttpClient(mcpURL)
		if err != nil {
			writeChunkBuffered(w, chunk{Type: "error", Error: fmt.Sprintf("MCP connect: %v", err)}, &buffer, isCloudfFront)
			if isCloudfFront {
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
			writeChunkBuffered(w, chunk{Type: "error", Error: fmt.Sprintf("MCP init: %v", err)}, &buffer, isCloudfFront)
			if isCloudfFront {
				flushBuffer(w, buffer)
			}
			return
		}

		toolsResult, err := mc.ListTools(ctx, mcp.ListToolsRequest{})
		if err != nil {
			writeChunkBuffered(w, chunk{Type: "error", Error: fmt.Sprintf("list tools: %v", err)}, &buffer, isCloudfFront)
			if isCloudfFront {
				flushBuffer(w, buffer)
			}
			return
		}
		tools := mcpToolsToAnthropic(toolsResult.Tools)

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
			for _, block := range resp.Content {
				switch block.Type {
				case "text":
					// Stream each text block as it arrives (or buffer if CloudFront)
					writeChunkBuffered(w, chunk{Type: "text", Text: block.Text}, &buffer, isCloudfFront)
				case "tool_use":
					toolUses = append(toolUses, block)
				}
			}

			if resp.StopReason == "end_turn" || len(toolUses) == 0 {
				break
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

// ── Main ───────────────────────────────────────────────────────────────────

func main() {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY is required")
	}
	model := os.Getenv("CLAUDE_MODEL")
	if model == "" {
		model = "claude-sonnet-4-5"
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
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	log.Printf("Safecast web-chat on :%s  MCP→%s  model=%s", port, mcpURL, model)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
