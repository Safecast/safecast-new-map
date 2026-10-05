package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

var geocodePlaceToolDef = mcp.NewTool("geocode_place",
	mcp.WithDescription("Convert a place name (village, town, city, region) to latitude/longitude and a bounding box. ALWAYS call this FIRST when the user names a place, then pass the returned min_lat/max_lat/min_lon/max_lon to sensor_current and query_radiation. Never guess coordinates yourself."),
	mcp.WithString("query",
		mcp.Required(),
		mcp.Description("Place name, ideally with region/country, e.g. 'Mitsue, Nara, Japan'"),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

type geocodeResult struct {
	Name   string  `json:"name"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	MinLat float64 `json:"min_lat"`
	MaxLat float64 `json:"max_lat"`
	MinLon float64 `json:"min_lon"`
	MaxLon float64 `json:"max_lon"`
}

var (
	geocodeMu    sync.Mutex
	geocodeCache = map[string][]geocodeResult{}
	geocodeLast  time.Time
)

func handleGeocodePlace(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := strings.TrimSpace(req.GetString("query", ""))
	if query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}
	key := strings.ToLower(query)

	geocodeMu.Lock()
	defer geocodeMu.Unlock()
	results, ok := geocodeCache[key]
	if !ok {
		// Nominatim usage policy: max 1 request/second, identifying User-Agent.
		if wait := time.Second - time.Since(geocodeLast); wait > 0 {
			time.Sleep(wait)
		}
		geocodeLast = time.Now()
		var err error
		if results, err = nominatimSearch(ctx, query); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("geocoding failed: %v", err)), nil
		}
		geocodeCache[key] = results
	}
	if len(results) == 0 {
		return mcp.NewToolResultError("no place found; try adding the region or country"), nil
	}
	return jsonResult(map[string]any{"query": query, "results": results})
}

func nominatimSearch(ctx context.Context, query string) ([]geocodeResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	u := "https://nominatim.openstreetmap.org/search?format=json&limit=3&q=" + url.QueryEscape(query)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "safecast-new-map/1.0 (https://simplemap.safecast.org)")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nominatim HTTP %d", resp.StatusCode)
	}
	var raw []struct {
		DisplayName string   `json:"display_name"`
		Lat         string   `json:"lat"`
		Lon         string   `json:"lon"`
		BoundingBox []string `json:"boundingbox"` // south, north, west, east
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	var out []geocodeResult
	for _, p := range raw {
		lat, _ := strconv.ParseFloat(p.Lat, 64)
		lon, _ := strconv.ParseFloat(p.Lon, 64)
		g := geocodeResult{Name: p.DisplayName, Lat: lat, Lon: lon, MinLat: lat - 0.5, MaxLat: lat + 0.5, MinLon: lon - 0.5, MaxLon: lon + 0.5}
		if len(p.BoundingBox) == 4 {
			g.MinLat, _ = strconv.ParseFloat(p.BoundingBox[0], 64)
			g.MaxLat, _ = strconv.ParseFloat(p.BoundingBox[1], 64)
			g.MinLon, _ = strconv.ParseFloat(p.BoundingBox[2], 64)
			g.MaxLon, _ = strconv.ParseFloat(p.BoundingBox[3], 64)
		}
		out = append(out, g)
	}
	return out, nil
}
