package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterAPIRoutesInventory(t *testing.T) {
	mux := http.NewServeMux()
	h := &RESTHandler{}
	h.registerAPIRoutes(mux)

	required := []string{
		"/api/radiation",
		"/api/area",
		"/api/tracks",
		"/api/track/",
		"/api/device/",
		"/api/sensors",
		"/api/sensors/export",
		"/api/sensor/",
		"/api/spectra",
		"/api/spectrum/",
		"/api/stats",
		"/api/extreme",
		"/api/info/",
		"/api/gpt/radiation",
		"/api/gpt/area",
		"/api/gpt/stats",
	}

	for _, route := range required {
		t.Run(route, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, route, nil)
			_, pattern := mux.Handler(req)
			if pattern == "" {
				t.Fatalf("route %s was not registered", route)
			}
		})
	}
}

func TestRegisterAPIRoutesSensorsCoexistence(t *testing.T) {
	mux := http.NewServeMux()
	h := &RESTHandler{}
	h.registerAPIRoutes(mux)

	reqList := httptest.NewRequest(http.MethodGet, "/api/sensors", nil)
	_, patternList := mux.Handler(reqList)
	if patternList != "/api/sensors" {
		t.Fatalf("expected /api/sensors pattern, got %q", patternList)
	}

	reqExport := httptest.NewRequest(http.MethodGet, "/api/sensors/export", nil)
	_, patternExport := mux.Handler(reqExport)
	if patternExport != "/api/sensors/export" {
		t.Fatalf("expected /api/sensors/export pattern, got %q", patternExport)
	}
}

func TestRegisterCompanionRoutesInventory(t *testing.T) {
	mainMux := http.NewServeMux()
	mcpMux := http.NewServeMux()

	registerCompanionRoutes(companionRoutesConfig{
		MainMux:             mainMux,
		MCPMux:              mcpMux,
		RESTHandler:         &RESTHandler{},
		FeedbackHandler:     func(http.ResponseWriter, *http.Request) {},
		TrackInsightsHandle: func(http.ResponseWriter, *http.Request) {},
		ChatHandler:         func(http.ResponseWriter, *http.Request) {},
	})

	mainRequired := []string{
		"/api/feedback",
		"/api/sensors",
		"/api/sensors/export",
		"/api/sensor/",
		"/chat",
	}
	for _, route := range mainRequired {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		_, pattern := mainMux.Handler(req)
		if pattern == "" {
			t.Fatalf("main mux route %s was not registered", route)
		}
	}

	reqInsights := httptest.NewRequest(http.MethodGet, "/api/track/abc/insights", nil)
	_, insightsPattern := mainMux.Handler(reqInsights)
	if insightsPattern != "GET /api/track/{id}/insights" {
		t.Fatalf("expected track insights pattern, got %q", insightsPattern)
	}

	mcpRequired := []string{
		"/api/feedback",
		"/chat",
	}
	for _, route := range mcpRequired {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		_, pattern := mcpMux.Handler(req)
		if pattern == "" {
			t.Fatalf("mcp mux route %s was not registered", route)
		}
	}
}

func TestRegisterCompanionRoutesTrackInsightsPrecedence(t *testing.T) {
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/api/track/", func(http.ResponseWriter, *http.Request) {})

	registerCompanionRoutes(companionRoutesConfig{
		MainMux:             mainMux,
		TrackInsightsHandle: func(http.ResponseWriter, *http.Request) {},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/track/abc/insights", nil)
	_, pattern := mainMux.Handler(req)
	if pattern != "GET /api/track/{id}/insights" {
		t.Fatalf("expected insights route precedence, got %q", pattern)
	}
}
