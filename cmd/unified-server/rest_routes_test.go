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
