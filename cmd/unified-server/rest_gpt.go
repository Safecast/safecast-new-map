package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

// GPT-friendly compact types — short field names to minimize response size.
type gptItem struct {
	USvH float64 `json:"usvh"`
	At   string  `json:"at"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Det  string  `json:"det,omitempty"`
	Dist float64 `json:"dist_m,omitempty"`
}

type gptResp struct {
	Count int       `json:"n"`
	Total int       `json:"total,omitempty"`
	Src   string    `json:"src"`
	Items []gptItem `json:"items"`
}

// RegisterGPT wires /api/gpt/* routes — compact endpoints for ChatGPT Custom GPT Actions.
// All routes are hard-capped at 5 results and return non-indented JSON.
func (h *RESTHandler) RegisterGPT(mux *http.ServeMux) {
	mux.HandleFunc("/api/gpt/radiation", h.handleGPTRadiation)
	mux.HandleFunc("/api/gpt/area", h.handleGPTArea)
	mux.HandleFunc("/api/gpt/stats", h.handleGPTStats)
}

func (h *RESTHandler) handleGPTRadiation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET only", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	lat, err := strconv.ParseFloat(q.Get("lat"), 64)
	if err != nil || lat < -90 || lat > 90 {
		writeError(w, http.StatusBadRequest, "invalid lat")
		return
	}
	lon, err := strconv.ParseFloat(q.Get("lon"), 64)
	if err != nil || lon < -180 || lon > 180 {
		writeError(w, http.StatusBadRequest, "invalid lon")
		return
	}
	radiusM := 1500.0
	if s := q.Get("radius_m"); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			radiusM = v
		}
	}

	var result *mcp.CallToolResult
	if dbAvailable() {
		result, _ = queryRadiationDB(r.Context(), lat, lon, radiusM, 5)
	} else {
		result, _ = queryRadiationAPI(r.Context(), lat, lon, radiusM, 5)
	}

	writeGPT(w, result)
}

func (h *RESTHandler) handleGPTArea(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET only", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	minLat, e1 := strconv.ParseFloat(q.Get("min_lat"), 64)
	maxLat, e2 := strconv.ParseFloat(q.Get("max_lat"), 64)
	minLon, e3 := strconv.ParseFloat(q.Get("min_lon"), 64)
	maxLon, e4 := strconv.ParseFloat(q.Get("max_lon"), 64)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		writeError(w, http.StatusBadRequest, "min_lat, max_lat, min_lon, max_lon required")
		return
	}

	var result *mcp.CallToolResult
	if dbAvailable() {
		result, _ = searchAreaDB(r.Context(), minLat, maxLat, minLon, maxLon, 5)
	} else {
		result, _ = searchAreaAPI(r.Context(), minLat, maxLat, minLon, maxLon, 5)
	}

	writeGPT(w, result)
}

func (h *RESTHandler) handleGPTStats(w http.ResponseWriter, r *http.Request) {
	// Stats responses are already small; just proxy through.
	h.handleStats(w, r)
}

// writeGPT extracts JSON from an MCP tool result, maps it to compact gptResp, and writes it.
func writeGPT(w http.ResponseWriter, result *mcp.CallToolResult) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if result == nil || len(result.Content) == 0 {
		_ = json.NewEncoder(w).Encode(gptResp{Src: "empty"})
		return
	}

	// Extract JSON text from MCP result
	raw := ""
	for _, c := range result.Content {
		if tc, ok := mcp.AsTextContent(c); ok && tc.Text != "" {
			raw = tc.Text
			break
		}
	}
	if raw == "" {
		_ = json.NewEncoder(w).Encode(gptResp{Src: "empty"})
		return
	}

	// Parse the full verbose result
	type loc struct {
		Lat float64 `json:"latitude"`
		Lon float64 `json:"longitude"`
	}
	type fullItem struct {
		Value      *float64 `json:"value"`
		CapturedAt string   `json:"captured_at"`
		Location   *loc     `json:"location"`
		DistanceM  float64  `json:"distance_m"`
		Detector   string   `json:"detector"`
	}
	type fullResult struct {
		Count        int        `json:"count"`
		TotalAvail   int        `json:"total_available"`
		Source       string     `json:"source"`
		Measurements []fullItem `json:"measurements"`
	}

	var full fullResult
	if err := json.Unmarshal([]byte(raw), &full); err != nil {
		// Fall back to raw if we can't parse
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(raw))
		return
	}

	items := make([]gptItem, 0, len(full.Measurements))
	for _, m := range full.Measurements {
		item := gptItem{
			At:   m.CapturedAt,
			Det:  m.Detector,
			Dist: m.DistanceM,
		}
		if m.Value != nil {
			item.USvH = *m.Value
		}
		if m.Location != nil {
			item.Lat = m.Location.Lat
			item.Lon = m.Location.Lon
		}
		items = append(items, item)
	}

	_ = json.NewEncoder(w).Encode(gptResp{
		Count: full.Count,
		Total: full.TotalAvail,
		Src:   full.Source,
		Items: items,
	})
}
