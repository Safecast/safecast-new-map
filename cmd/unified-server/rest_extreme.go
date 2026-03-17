package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// @Summary Find extreme radiation readings
// @Description Find the highest or lowest radiation readings in the database with full location details. Supports excluding anomalous sources.
// @Tags analytics
// @Produce json
// @Param direction query string false "Direction: 'highest' or 'lowest'" default(highest)
// @Param limit query int false "Number of results (1-100)" default(10)
// @Param min_lat query number false "Southern boundary for geographic filter" default(-90)
// @Param max_lat query number false "Northern boundary for geographic filter" default(90)
// @Param min_lon query number false "Western boundary for geographic filter" default(-180)
// @Param max_lon query number false "Eastern boundary for geographic filter" default(180)
// @Param exclude_devices query string false "Comma-separated device IDs to exclude (e.g., 'bGeigie-2113,bGeigie-456')"
// @Param exclude_areas query string false "JSON array of bounding boxes to exclude (e.g., '[{\"min_lat\":51.8,\"max_lat\":52.0,\"min_lon\":-8.6,\"max_lon\":-8.3}]')"
// @Success 200 {object} map[string]interface{} "Extreme readings with location details"
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/extreme [get]
func handleRESTExtremeReadings(w http.ResponseWriter, r *http.Request) {
	direction := r.URL.Query().Get("direction")
	if direction == "" {
		direction = "highest"
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	minLat := -90.0
	if v := r.URL.Query().Get("min_lat"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			minLat = parsed
		}
	}

	maxLat := 90.0
	if v := r.URL.Query().Get("max_lat"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			maxLat = parsed
		}
	}

	minLon := -180.0
	if v := r.URL.Query().Get("min_lon"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			minLon = parsed
		}
	}

	maxLon := 180.0
	if v := r.URL.Query().Get("max_lon"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			maxLon = parsed
		}
	}

	// Parse exclusion parameters
	var excludeDevices []string
	if devicesStr := r.URL.Query().Get("exclude_devices"); devicesStr != "" {
		excludeDevices = strings.Split(devicesStr, ",")
		// Trim whitespace from each device ID
		for i, dev := range excludeDevices {
			excludeDevices[i] = strings.TrimSpace(dev)
		}
	}

	excludeAreas := r.URL.Query().Get("exclude_areas")

	// Create MCP request
	req := mcp.CallToolRequest{}
	req.Params.Name = "query_extreme_readings"
	args := map[string]any{
		"direction": direction,
		"limit":     float64(limit),
		"min_lat":   minLat,
		"max_lat":   maxLat,
		"min_lon":   minLon,
		"max_lon":   maxLon,
	}

	if len(excludeDevices) > 0 {
		args["exclude_devices"] = excludeDevices
	}
	if excludeAreas != "" {
		args["exclude_areas"] = excludeAreas
	}

	req.Params.Arguments = args

	result, err := handleQueryExtremeReadings(r.Context(), req)
	serveMCPResult(w, result, err)
}
