package main

import (
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// handleInfo handles GET /api/info/{topic}
//
// @Summary     Get radiation reference information
// @Description Returns static educational reference content about radiation units, safety levels, detector types, background levels, and isotopes.
// @Tags        reference
// @Produce     json
// @Param       topic path string true "Topic to retrieve" Enums(units, dose_rates, safety_levels, detectors, background_levels, isotopes)
// @Success     200 {object} map[string]interface{} "Reference content for the requested topic"
// @Failure     400 {object} map[string]string "Invalid or missing topic"
// @Router      /info/{topic} [get]
func (h *RESTHandler) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract topic from path: /api/info/{topic}
	topic := strings.TrimPrefix(r.URL.Path, "/api/info/")
	if topic == "" {
		writeError(w, http.StatusBadRequest, "topic is required. Valid topics: units, dose_rates, safety_levels, detectors, background_levels, isotopes")
		return
	}

	// Construct a minimal MCP request and reuse the existing handler.
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"topic": topic}
	result, err := handleRadiationInfo(r.Context(), req)
	serveMCPResult(w, result, err)
}
