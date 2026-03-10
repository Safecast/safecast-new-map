// Test MCP Server + AI + Hints with Real REST API
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
				"lat": map[string]interface{}{"type": "number"},
				"lon": map[string]interface{}{"type": "number"},
				"radius_m": map[string]interface{}{"type": "integer"},
				"limit": map[string]interface{}{"type": "integer"},
			},
			"required": []string{"lat", "lon"},
		},
		Endpoint: "/api/radiation",
		Method:   "GET",
	},
	{
		Name:        "sensor_current",
		Description: "Get current real-time readings from fixed sensors. ALWAYS call this after query_radiation.",
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
		Description: "List available gamma spectroscopy records.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"limit": map[string]interface{}{"type": "integer"},
			},
		},
		Endpoint: "/api/spectra",
		Method:   "GET",
	},
	{
		Name:        "get_spectrum",
		Description: "Get full gamma spectrum data for a marker.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"marker_id": map[string]interface{}{"type": "integer"},
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
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	mcpURL := os.Getenv("MCP_URL")
	if mcpURL == "" {
		mcpURL = "http://localhost:3333"
	}

	fmt.Printf("🔧 MCP Server: %s\n", mcpURL)
	fmt.Print("📡 Testing connection... ")
	
	resp, err := http.Get(mcpURL + "/api/sensors")
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		os.Exit(1)
	}
	resp.Body.Close()
	fmt.Println("✅ Connected")
	fmt.Println()

	apiKey := os.Getenv("NVIDIA_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ NVIDIA_API_KEY not set")
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
	fmt.Println("║  Try these queries (or type your own):                    ║")
	fmt.Println("║  1. List gamma spectra                                    ║")
	fmt.Println("║  2. Current sensors in Japan                              ║")
	fmt.Println("║  3. Radiation near Tokyo                                  ║")
	fmt.Println("║  4. Search area around Fukushima                          ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	messages := []Message{
		{Role: "system", Content: `You are a helpful AI assistant for Safecast radiation monitoring data.

AVAILABLE TOOLS:
- query_radiation: Find radiation measurements near coordinates (historical data)
- sensor_current: Get real-time sensor readings (ALWAYS use after query_radiation)
- list_spectra: List gamma spectroscopy records
- get_spectrum: Get full gamma spectrum for a marker (needs marker_id from list_spectra)
- search_area: Find radiation in a bounding box

IMPORTANT RULES:
1. ALWAYS call sensor_current after query_radiation to get complete picture
2. For spectra: first call list_spectra, then get_spectrum with a marker_id
3. Use coordinates: Japan (35.6762, 139.6503), Fukushima (37.75, 140.5)

RESPONSE FORMAT:
- Present data in tables
- Include units (µSv/h, CPM)
- Add map links: [lat,lon](https://simplemap.safecast.org/?lat=X&lon=Y&zoom=10)`},
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
		if strings.ToLower(prompt) == "1" {
			prompt = "List available gamma spectra"
			fmt.Println(prompt)
		}
		if strings.ToLower(prompt) == "2" {
			prompt = "Show current sensor readings in Japan (lat 35-38, lon 139-142)"
			fmt.Println(prompt)
		}
		if strings.ToLower(prompt) == "3" {
			prompt = "What's the radiation level near Tokyo, Japan? (lat 35.6762, lon 139.6503)"
			fmt.Println(prompt)
		}
		if strings.ToLower(prompt) == "4" {
			prompt = "Search for radiation data around Fukushima (min_lat 37, max_lat 38, min_lon 140, max_lon 141)"
			fmt.Println(prompt)
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
				fmt.Printf("   → %s(%s)\n", tc.Function.Name, tc.Function.Arguments)

				result, err := callMCPTool(mcpURL, tc.Function.Name, tc.Function.Arguments)
				if err != nil {
					fmt.Printf("   ❌ Error: %v\n", err)
					continue
				}

				// Show more data for small responses
				preview := result
				if len(preview) > 300 {
					preview = preview[:300] + "..."
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

	client := &http.Client{Timeout: 90 * time.Second}
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

	var endpoint, method string
	for _, tool := range mcpTools {
		if tool.Name == toolName {
			endpoint = tool.Endpoint
			method = tool.Method
			break
		}
	}

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
