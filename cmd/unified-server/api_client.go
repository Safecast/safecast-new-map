package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

const defaultBaseURL = "https://simplemap.safecast.org"

var client = NewSafecastClient()

type SafecastClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewSafecastClient() *SafecastClient {
	baseURL := os.Getenv("SIMPLEMAP_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &SafecastClient{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		baseURL:    baseURL,
	}
}

// GetLatestNearby queries /api/latest for measurements near a location.
func (c *SafecastClient) GetLatestNearby(ctx context.Context, lat, lon, radiusM float64, limit int) (map[string]any, error) {
	v := url.Values{}
	v.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	v.Set("lon", strconv.FormatFloat(lon, 'f', -1, 64))
	v.Set("radius_m", strconv.FormatFloat(radiusM, 'f', -1, 64))
	v.Set("limit", strconv.Itoa(limit))
	return c.getObject(ctx, "/api/latest", v)
}

// GetMarkers queries /get_markers for measurements within a bounding box.
func (c *SafecastClient) GetMarkers(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]map[string]any, error) {
	v := url.Values{}
	v.Set("minLat", strconv.FormatFloat(minLat, 'f', -1, 64))
	v.Set("minLon", strconv.FormatFloat(minLon, 'f', -1, 64))
	v.Set("maxLat", strconv.FormatFloat(maxLat, 'f', -1, 64))
	v.Set("maxLon", strconv.FormatFloat(maxLon, 'f', -1, 64))
	v.Set("zoom", "0")
	return c.getList(ctx, "/get_markers", v)
}

// GetTracks queries /api/tracks for all track summaries.
func (c *SafecastClient) GetTracks(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "/api/tracks", nil)
}

// GetTracksByYear queries /api/tracks/years/{year}.
// Returns an empty tracks map (not an error) if the year has no data.
func (c *SafecastClient) GetTracksByYear(ctx context.Context, year int) (map[string]any, error) {
	path := fmt.Sprintf("/api/tracks/years/%d", year)
	result, err := c.getObject(ctx, path, nil)
	if err != nil {
		// Upstream returns 404 when no tracks exist for the year — treat as empty.
		if isNotFound(err) {
			return map[string]any{"tracks": []any{}}, nil
		}
		return nil, err
	}
	return result, nil
}

// GetTracksByMonth queries /api/tracks/months/{year}/{month}.
// Returns an empty tracks map (not an error) if the month has no data.
func (c *SafecastClient) GetTracksByMonth(ctx context.Context, year, month int) (map[string]any, error) {
	path := fmt.Sprintf("/api/tracks/months/%d/%d", year, month)
	result, err := c.getObject(ctx, path, nil)
	if err != nil {
		// Upstream returns 404 when no tracks exist for the month — treat as empty.
		if isNotFound(err) {
			return map[string]any{"tracks": []any{}}, nil
		}
		return nil, err
	}
	return result, nil
}

// GetTrackData queries /api/track/{trackID}.json for full track markers.
func (c *SafecastClient) GetTrackData(ctx context.Context, trackID string, from, to int) (map[string]any, error) {
	path := fmt.Sprintf("/api/track/%s.json", url.PathEscape(trackID))
	v := url.Values{}
	if from != 0 {
		v.Set("from", strconv.Itoa(from))
	}
	if to != 0 {
		v.Set("to", strconv.Itoa(to))
	}
	return c.getObject(ctx, path, v)
}

// GetRealtimeHistory queries /realtime_history for device measurement history.
func (c *SafecastClient) GetRealtimeHistory(ctx context.Context, deviceID string) (map[string]any, error) {
	v := url.Values{}
	v.Set("device", deviceID)
	return c.getObject(ctx, "/realtime_history", v)
}

// GetSpectrum queries /api/spectrum/{markerID} for gamma spectroscopy data.
func (c *SafecastClient) GetSpectrum(ctx context.Context, markerID int) (map[string]any, error) {
	path := fmt.Sprintf("/api/spectrum/%d", markerID)
	return c.getObject(ctx, path, nil)
}

func (c *SafecastClient) getObject(ctx context.Context, path string, params url.Values) (map[string]any, error) {
	body, err := c.doGet(ctx, path, params)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *SafecastClient) getList(ctx context.Context, path string, params url.Values) ([]map[string]any, error) {
	body, err := c.doGet(ctx, path, params)
	if err != nil {
		return nil, err
	}
	var result []map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *SafecastClient) doGet(ctx context.Context, path string, params url.Values) ([]byte, error) {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("no response from simplemap API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("simplemap API error (404): %s", resp.Status)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("simplemap API error (%d): %s", resp.StatusCode, resp.Status)
	}

	return body, nil
}

// isNotFound returns true if the error is a 404 from the upstream API.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	return len(err.Error()) >= 30 && err.Error()[:30] == "simplemap API error (404): 404"
}

// normalizeLatestMarker converts a marker from /api/latest to MCP output format.
func normalizeLatestMarker(m map[string]any) map[string]any {
	return map[string]any{
		"id":          m["id"],
		"value":       m["doseRateMicroSvH"],
		"unit":        "µSv/h",
		"captured_at": m["timeUTC"],
		"location": map[string]any{
			"latitude":  m["lat"],
			"longitude": m["lon"],
		},
		"height":     m["altitudeM"],
		"detector":   m["detectorType"],
		"speed_ms":   m["speedMS"],
		"count_rate": m["countRateCPS"],
	}
}

// normalizeGetMarker converts a marker from /get_markers to MCP output format.
func normalizeGetMarker(m map[string]any) map[string]any {
	result := map[string]any{
		"id":    m["id"],
		"value": m["doseRate"],
		"unit":  "µSv/h",
		"location": map[string]any{
			"latitude":  m["lat"],
			"longitude": m["lon"],
		},
		"track_id":     m["trackID"],
		"height":       m["altitude"],
		"detector":     m["detector"],
		"has_spectrum": m["hasSpectrum"],
	}
	if date, ok := toFloat(m["date"]); ok && date > 0 {
		t := time.Unix(int64(date), 0).UTC()
		result["captured_at"] = t.Format(time.RFC3339)
	}
	return result
}

// jsonResult serializes v to indented JSON and returns it as a tool result.
func jsonResult(v any) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError("failed to serialize response"), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}
