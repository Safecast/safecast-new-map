// Interactive NVIDIA API Model Tester
// Usage: export NVIDIA_API_KEY=nvapi-... && go run test-nvidia-models.go

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
)

// NVIDIA API models (update as new models become available)
var availableModels = []string{
	"qwen/qwen3.5-122b-a10b",
	"qwen/qwen-2.5-coder-32b",
	"meta/llama-3.1-405b-instruct",
	"meta/llama-3.1-70b-instruct",
	"meta/llama-3.2-90b-vision-instruct",
	"google/gemma-2-27b-it",
	"google/gemma-2-9b-it",
	"mistralai/mistral-large-2-instruct",
	"mistralai/mixtral-8x22b-instruct-v0.1",
	"nvidia/nemotron-4-340b-instruct",
}

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
		fmt.Println()
		fmt.Println("Get your API key from: https://build.nvidia.com/explore/discover")
		fmt.Println("Then run: export NVIDIA_API_KEY=nvapi-...")
		fmt.Println()
		
		// Offer to set it interactively
		fmt.Print("Or paste your API key here: ")
		reader := bufio.NewReader(os.Stdin)
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)
		if key != "" {
			apiKey = key
			fmt.Println()
			fmt.Println("💡 Tip: To avoid entering this each time, run:")
			fmt.Printf("   export NVIDIA_API_KEY=%s\n", key)
		} else {
			fmt.Println("No API key provided. Exiting.")
			os.Exit(1)
		}
	}

	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║     NVIDIA API Interactive Model Tester                   ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Show available models
	fmt.Println("📦 Available Models:")
	fmt.Println("────────────────────────────────────────────────────────────")
	for i, model := range availableModels {
		fmt.Printf("  %2d. %s\n", i+1, model)
	}
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println()

	// Let user select model
	fmt.Print("Select model (1-10) or enter custom model name: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var selectedModel string
	if num, err := strconv.Atoi(input); err == nil && num > 0 && num <= len(availableModels) {
		selectedModel = availableModels[num-1]
	} else if input != "" {
		selectedModel = input
	} else {
		selectedModel = availableModels[0] // Default to first
	}

	fmt.Printf("\n✅ Selected: %s\n", selectedModel)
	fmt.Println()

	// Chat loop
	messages := []Message{
		{Role: "system", Content: "You are a helpful AI assistant."},
	}

	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║  Enter your prompts (type 'quit' to exit, 'clear' to     ║")
	fmt.Println("║  reset conversation, 'model' to change model)             ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

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

		if strings.ToLower(prompt) == "clear" {
			messages = []Message{{Role: "system", Content: "You are a helpful AI assistant."}}
			fmt.Println("🗑️  Conversation cleared.\n")
			continue
		}

		if strings.ToLower(prompt) == "model" {
			fmt.Println("\n📦 Available Models:")
			for i, m := range availableModels {
				fmt.Printf("  %2d. %s\n", i+1, m)
			}
			fmt.Print("\nSelect model: ")
			modelInput, _ := reader.ReadString('\n')
			modelInput = strings.TrimSpace(modelInput)
			if num, err := strconv.Atoi(modelInput); err == nil && num > 0 && num <= len(availableModels) {
				selectedModel = availableModels[num-1]
			} else if modelInput != "" {
				selectedModel = modelInput
			}
			fmt.Printf("✅ Changed to: %s\n\n", selectedModel)
			continue
		}

		// Add user message
		messages = append(messages, Message{Role: "user", Content: prompt})

		// Call API
		fmt.Print("🤖 AI: ")
		
		response, err := callNVIDIA(apiKey, selectedModel, messages)
		if err != nil {
			fmt.Printf("\n❌ Error: %v\n\n", err)
			continue
		}

		fmt.Println(response)
		fmt.Println()

		// Add AI response to history
		messages = append(messages, Message{Role: "assistant", Content: response})
	}
}

func callNVIDIA(apiKey, model string, messages []Message) (string, error) {
	reqBody := Request{
		Model:     model,
		Messages:  messages,
		MaxTokens: 2048,
	}

	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://integrate.api.nvidia.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(raw))
	}

	var result Response
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("%s: %s", result.Error.Type, result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from model")
	}

	return result.Choices[0].Message.Content, nil
}
