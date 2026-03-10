// Test MCP Server + AI + Hints with Real REST API
// Usage: 
//   1. Start MCP server: ./bin/mcp-server-test
//   2. Run: export NVIDIA_API_KEY=nvapi-... && go run test-mcp-ai.go

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var availableModels = []string{
	"qwen/qwen3.5-122b-a10b",
	"meta/llama-3.1-70b-instruct",
	"mistralai/mistral-large-2-instruct",
}

// MCP REST API tool definitions (matching your server)
var mcpTools = []struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
	Endpoint    string
	Method      string
}{
	{
		Name:        "query_radiation",
		Description: "Find radiation measurements near a location. Returns historical bGeigie data.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"lat": map[string]interface{}{"type": "number", "description": "Latitude"},
				"lon": map[string]interface{}{"type": "number", "description": "Longitude"},
				"radius_m": map[string]interface{}{"type": "integer", "description": "Radius in meters (25-50000)"},
				"limit": map[string]interface{}{"type": "integer", "description": "Max results (1-10000)"},
			},
			"required": []string{"lat", "lon"},
		},
		Endpoint: "/api/radiation",
		Method:   "GET",
	},
	{
		Name:        "sensor_current",
		Description: "Get current real-time readings from fixed sensors. CRITICAL: Always call this after query_radiation to check for real-time data.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"min_lat": map[string]interface{}{"type": "number"},
				"max_lat": map[string]interface{}{"type": "number"},
				"min_lon": map[string]interface{}{"type": "number"},
				"max_lon": map[string]interface{}{"type": "number"},
			},
		},
		Endpoint: "/api/sensors",
		Method:   "GET",
	},
	{
		Name:        "list_spectra",
		Description: "List available gamma spectroscopy records. Returns marker IDs for full spectrum data.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"limit": map[string]interface{}{"type": "integer", "description": "Max results (1-500)"},
			},
		},
		Endpoint: "/api/spectra",
		Method:   "GET",
	},
	{
		Name:        "get_spectrum",
		Description: "Get full gamma spectrum data for a marker. Includes channel counts for isotope identification.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"marker_id": map[string]interface{}{"type": "integer", "description": "Marker ID from list_spectra"},
			},
			"required": []string{"marker_id"},
		},
		Endpoint: "/api/spectrum/{marker_id}",
		Method:   "GET",
	},
	{
		Name:        "search_area",
		Description: "Find radiation measurements in a geographic bounding box.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"min_lat": map[string]interface{}{"type": "number"},
				"max_lat": map[string]interface{}{"type": "number"},
				"min_lon": map[string]interface{}{"type": "number"},
				"max_lon": map[string]interface{}{"type": "number"},
				"limit": map[string]interface{}{"type": "integer"},
			},
			"required": []string{"min_lat", "max_lat", "min_lon", "max_lon"},
		},
		Endpoint: "/api/area",
		Method:   "GET",
	},
	{
		Name:        "radiation_stats",
		Description: "Get overall database statistics: total measurements, coverage, contributors.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
		},
		Endpoint: "/api/stats",
		Method:   "GET",
	},
}

type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type AIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	Tools     []Tool    `json:"tools,omitempty"`
	MaxTokens int       `json:"max_tokens"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type AIResponse struct {
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║     MCP Server + AI + Hints Integration Test              ║")
	fmt.Println("║     Tests: query_radiation, sensor_current, spectra       ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	mcpURL := os.Getenv("MCP_URL")
	if mcpURL == "" {
		mcpURL = "http://localhost:3333"
	}

	fmt.Printf("🔧 MCP Server: %s\n", mcpURL)
	fmt.Print("📡 Testing connection... ")
	
	// Test with stats endpoint
	resp, err := http.Get(mcpURL + "/api/stats")
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		fmt.Println("\n⚠️  Start MCP server first:")
		fmt.Println("   ./bin/mcp-server-test")
		os.Exit(1)
	}
	resp.Body.Close()
	fmt.Println("✅ Connected")
	fmt.Println()

	apiKey := os.Getenv("NVIDIA_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ NVIDIA_API_KEY not set")
		fmt.Println("   export NVIDIA_API_KEY=nvapi-...")
		os.Exit(1)
	}
	fmt.Println("✅ NVIDIA API key found")
	fmt.Println()

	fmt.Println("📦 AI Models:")
	for i, model := range availableModels {
		fmt.Printf("  %2d. %s\n", i+1, model)
	}
	fmt.Print("\nSelect model (1-3): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	modelIdx := 0
	if num, err := strconv.Atoi(strings.TrimSpace(input)); err == nil && num > 0 && num <= len(availableModels) {
		modelIdx = num - 1
	}
	selectedModel := availableModels[modelIdx]
	fmt.Printf("\n✅ Using: %s\n\n", selectedModel)

	fmt.Println("🛠️  Available MCP Tools:")
	for _, tool := range mcpTools {
		fmt.Printf("  • %-20s %s\n", tool.Name, truncate(tool.Description, 50))
	}
	fmt.Println()

	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║  Try these queries:                                       ║")
	fmt.Println("║  • What's radiation near Tokyo?                           ║")
	fmt.Println("║  • Show me gamma spectra data                             ║")
	fmt.Println("║  • Current sensor readings in Fukushima                   ║")
	fmt.Println("║  • Database statistics                                    ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	messages := []Message{
		{Role: "system", Content: `You are an AI assistant for Safecast radiation data. You have access to these tools:

- query_radiation: Historical bGeigie measurements near a location
- sensor_current: Real-time fixed sensor readings (ALWAYS call this after query_radiation)
- list_spectra: List gamma spectroscopy records
- get_spectrum: Full gamma spectrum with channel data
- search_area: Radiation in a bounding box
- radiation_stats: Database statistics

IMPORTANT:
1. For location queries, ALWAYS call sensor_current after query_radiation
2. For spectra questions, call list_spectra first, then get_spectrum
3. Present data in tables with map links: [lat,lon](https://simplemap.safecast.org/?lat=X&lon=Y&zoom=10)
4. Include units (µSv/h, CPM) and compare to background (0.05-0.20 µSv/h)`},
	}

	for {
		fmt.Print("🧑 You: ")
		prompt, _ := reader.ReadString('\n')
		prompt = strings.TrimSpace(prompt)

		if prompt == "" {
			continue
		}
		if strings.ToLower(prompt) == "quit" || strings.ToLower(prompt) == "exit" {
			fmt.Println("\n👋 Goodbye!")
			break
		}

		messages = append(messages, Message{Role: "user", Content: prompt})

		fmt.Print("🤖 AI: ")

		response, toolCalls, err := callAIWithTools(apiKey, selectedModel, messages, mcpTools)
		if err != nil {
			fmt.Printf("\n❌ AI Error: %v\n\n", err)
			continue
		}

		if len(toolCalls) > 0 {
			fmt.Printf("\n🔧 Calling %d MCP tool(s)...\n", len(toolCalls))

			for _, tc := range toolCalls {
				fmt.Printf("   → %s(...)\n", tc.Function.Name)

				result, err := callMCPTool(mcpURL, tc.Function.Name, tc.Function.Arguments)
				if err != nil {
					fmt.Printf("   ❌ Error: %v\n", err)
					continue
				}

				preview := result
				if len(preview) > 200 {
					preview = preview[:200] + "..."
				}
				fmt.Printf("   ✅ Got: %s\n", preview)

				messages = append(messages, Message{
					Role:    "tool",
					Content: string(result),
				})
			}

			fmt.Print("\n🤖 AI (with data): ")
			response, _, err = callAIWithTools(apiKey, selectedModel, messages, mcpTools)
			if err != nil {
				fmt.Printf("\n❌ AI Error: %v\n\n", err)
				continue
			}
		}

		fmt.Printf("\n%s\n\n", response)
		messages = append(messages, Message{Role: "assistant", Content: response})
	}
}

func callAIWithTools(apiKey, model string, messages []Message, tools []struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
	Endpoint    string
	Method      string
}) (string, []ToolCall, error) {
	aiTools := make([]Tool, len(tools))
	for i, t := range tools {
		aiTools[i] = Tool{
			Type: "function",
			Function: ToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			},
		}
	}

	reqBody := AIRequest{
		Model:     model,
		Messages:  messages,
		Tools:     aiTools,
		MaxTokens: 2048,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://integrate.api.nvidia.com/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	var result AIResponse
	json.Unmarshal(raw, &result)

	if result.Error != nil {
		return "", nil, fmt.Errorf("%s: %s", result.Error.Type, result.Error.Message)
	}
	if len(result.Choices) == 0 {
		return "", nil, fmt.Errorf("no response")
	}

	msg := result.Choices[0].Message
	return msg.Content, msg.ToolCalls, nil
}

func callMCPTool(mcpURL, toolName, argsJSON string) (string, error) {
	var args map[string]interface{}
	json.Unmarshal([]byte(argsJSON), &args)

	// Find tool endpoint
	var endpoint, method string
	for _, tool := range mcpTools {
		if tool.Name == toolName {
			endpoint = tool.Endpoint
			method = tool.Method
			break
		}
	}

	// Build URL with query params
	url := endpoint
	if method == "GET" && len(args) > 0 {
		params := []string{}
		for k, v := range args {
			if k == "marker_id" {
				url = strings.Replace(url, "{marker_id}", fmt.Sprintf("%v", v), 1)
			} else {
				params = append(params, fmt.Sprintf("%s=%v", k, v))
			}
		}
		if len(params) > 0 {
			url += "?" + strings.Join(params, "&")
		}
	}

	resp, err := http.Get(mcpURL + url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)
	return string(result), nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
