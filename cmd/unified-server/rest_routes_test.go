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
		"/api/feedback",
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
		MainMux:     mainMux,
		MCPMux:      mcpMux,
		ChatHandler: func(http.ResponseWriter, *http.Request) {},
	})

	mainRequired := []string{
		"/chat",
	}
	for _, route := range mainRequired {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		_, pattern := mainMux.Handler(req)
		if pattern == "" {
			t.Fatalf("main mux route %s was not registered", route)
		}
	}

	mcpRequired := []string{"/chat"}
	for _, route := range mcpRequired {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		_, pattern := mcpMux.Handler(req)
		if pattern == "" {
			t.Fatalf("mcp mux route %s was not registered", route)
		}
	}
}
