package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

var validTopics = []string{"units", "dose_rates", "safety_levels", "detectors", "background_levels", "isotopes"}

var radiationInfoToolDef = mcp.NewTool("radiation_info",
	mcp.WithDescription("Get educational reference information about radiation units, safety levels, detectors, and related topics. Returns static reference content. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool."),
	mcp.WithString("topic",
		mcp.Description("Topic to retrieve information about"),
		mcp.Enum(validTopics...),
		mcp.Required(),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleRadiationInfo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topic, err := req.RequireString("topic")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	normalized := strings.ToLower(strings.ReplaceAll(topic, "-", "_"))

	content, ok := referenceData[normalized]
	if !ok {
		return mcp.NewToolResultError(fmt.Sprintf(
			"Invalid topic: %q. Valid topics: %s", topic, strings.Join(validTopics, ", "),
		)), nil
	}

	result := map[string]any{
		"topic":   normalized,
		"content": content,
		"_ai_hint": "CRITICAL INSTRUCTIONS: (1) The .unit. field indicates measurement units - CPM means .counts per minute. NOT .counts per second.. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I.ll, I.m, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: .Latest reading: X CPM at location Y. NOT .I found a reading of X CPM. or .Perfect! The sensor shows..... State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}
