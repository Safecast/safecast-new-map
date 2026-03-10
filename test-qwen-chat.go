// Minimal test client for Qwen3.5 via NVIDIA API
// Usage: export NVIDIA_API_KEY=nvapi-... && go run test-qwen-chat.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Response struct {
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
	apiKey := os.Getenv("NVIDIA_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ NVIDIA_API_KEY not set")
		fmt.Println("   Get it from: https://build.nvidia.com/qwen/qwen3.5-122b-a10b")
		fmt.Println("   Then: export NVIDIA_API_KEY=nvapi-...")
		os.Exit(1)
	}

	fmt.Println("=== Testing Qwen3.5-122B via NVIDIA API ===")
	fmt.Println()

	// Test query
	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant for Safecast radiation data."},
		{Role: "user", Content: "What is the typical background radiation level in µSv/h?"},
	}

	reqBody := Request{
		Model:     "qwen/qwen3.5-122b-a10b",
		Messages:  messages,
		MaxTokens: 500,
	}

	body, _ := json.Marshal(reqBody)

	fmt.Printf("Request: POST https://integrate.api.nvidia.com/v1/chat/completions\n")
	fmt.Printf("Model: qwen/qwen3.5-122b-a10b\n")
	fmt.Println()

	// Create request with proper headers
	req, err := http.NewRequest("POST", "https://integrate.api.nvidia.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("❌ Failed to create request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Printf("HTTP Status: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Println()

	raw, _ := io.ReadAll(resp.Body)

	// Print raw
	fmt.Println("Raw response:")
	fmt.Println(string(raw))
	fmt.Println()

	// Try to parse
	var result Response
	err = json.Unmarshal(raw, &result)
	if err != nil {
		fmt.Printf("❌ Failed to parse JSON: %v\n", err)
		os.Exit(1)
	}

	if result.Error != nil {
		fmt.Printf("❌ API Error: %s - %s\n", result.Error.Type, result.Error.Message)
		os.Exit(1)
	}

	if len(result.Choices) > 0 {
		fmt.Println("✅ Success!")
		fmt.Println()
		fmt.Println("Response:")
		fmt.Println(result.Choices[0].Message.Content)
		fmt.Println()
		fmt.Printf("Finish reason = %s\n", result.Choices[0].FinishReason)
	} else {
		fmt.Println("❌ No choices in response")
		os.Exit(1)
	}
}
