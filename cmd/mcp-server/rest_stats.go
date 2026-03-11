package main

import (
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

// handleStats handles GET /api/stats
//
// @Summary     Get aggregate radiation statistics
// @Description Returns aggregate radiation dose-rate statistics from the Safecast database grouped by time interval. Powered by DuckDB + PostgreSQL.
// @Tags        reference
// @Produce     json
// @Param       interval query string false "Aggregation interval: year, month, or overall" Enums(year, month, overall) default(year)
// @Success     200 {object} map[string]interface{} "Statistics data with interval and source metadata"
// @Failure     400 {object} map[string]string "Invalid interval value"
// @Failure     503 {object} map[string]string "Analytics engine unavailable"
// @Router      /stats [get]
func (h *RESTHandler) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = "year"
	}
	switch interval {
	case "year", "month", "overall":
		// valid
	default:
		writeError(w, http.StatusBadRequest, "interval must be one of: year, month, overall")
		return
	}

	// Construct a minimal MCP request and reuse the existing handler.
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"interval": interval}
	result, err := handleRadiationStats(r.Context(), req)
	serveMCPResult(w, result, err)
}
