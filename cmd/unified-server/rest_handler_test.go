package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRESTMux(t *testing.T) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	h := &RESTHandler{}
	h.registerAPIRoutes(mux)
	return mux
}

func TestHandleRadiation_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, "/api/radiation?lat=35.0&lon=139.0", nil)
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleRadiation_BadRequest(t *testing.T) {
	mux := newRESTMux(t)
	cases := []struct {
		name string
		url  string
	}{
		{"missing lat and lon", "/api/radiation"},
		{"missing lon", "/api/radiation?lat=35.0"},
		{"missing lat", "/api/radiation?lon=139.0"},
		{"non-numeric lat", "/api/radiation?lat=invalid&lon=139.0"},
		{"lat above range", "/api/radiation?lat=91&lon=139.0"},
		{"lat below range", "/api/radiation?lat=-91&lon=139.0"},
		{"non-numeric radius_m", "/api/radiation?lat=35.0&lon=139.0&radius_m=bad"},
		{"radius_m too small", "/api/radiation?lat=35.0&lon=139.0&radius_m=1"},
		{"radius_m too large", "/api/radiation?lat=35.0&lon=139.0&radius_m=99999"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleArea_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/area", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleArea_BadRequest(t *testing.T) {
	mux := newRESTMux(t)
	cases := []struct {
		name string
		url  string
	}{
		{"missing all params", "/api/area"},
		{"missing max_lat", "/api/area?min_lat=35.0&min_lon=139.0&max_lon=140.0"},
		{"non-numeric min_lat", "/api/area?min_lat=bad&max_lat=36.0&min_lon=139.0&max_lon=140.0"},
		{"min_lat out of range", "/api/area?min_lat=91&max_lat=92&min_lon=139.0&max_lon=140.0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleTracks_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tracks", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleTracks_BadRequest(t *testing.T) {
	mux := newRESTMux(t)
	cases := []struct {
		name string
		url  string
	}{
		{"non-numeric year", "/api/tracks?year=notanumber"},
		{"year too early", "/api/tracks?year=1999"},
		{"year too late", "/api/tracks?year=2101"},
		{"month without year", "/api/tracks?month=6"},
		{"month out of range", "/api/tracks?year=2024&month=13"},
		{"month zero", "/api/tracks?year=2024&month=0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleTrack_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/track/abc123", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleTrack_MissingID(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/track/", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleTrack_BadRequest(t *testing.T) {
	mux := newRESTMux(t)
	cases := []struct {
		name string
		url  string
	}{
		{"negative from", "/api/track/x?from=-1"},
		{"non-numeric from", "/api/track/x?from=bad"},
		{"negative to", "/api/track/x?to=-1"},
		{"limit zero", "/api/track/x?limit=0"},
		{"limit too large", "/api/track/x?limit=10001"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleDevice_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/device/some-id/history", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleDevice_MissingID(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/device/", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleStats_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/stats", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleStats_InvalidInterval(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/stats?interval=weekly", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleInfo_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/info/units", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleInfo_MissingTopic(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/info/", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleSensors_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sensors", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleSensors_DBUnavailable(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/sensors", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleSensorsExport_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sensors/export", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleSensorsExport_DBUnavailable(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/sensors/export", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleSensor_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sensor/dev-1/current", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleSensor_DBUnavailable(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/sensor/dev-1/current", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleSpectra_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/spectra", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleSpectra_DBUnavailable(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/spectra", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleSpectrum_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/spectrum/1", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleSpectrum_BadRequest(t *testing.T) {
	mux := newRESTMux(t)
	cases := []struct {
		name string
		url  string
	}{
		{"missing marker id", "/api/spectrum/"},
		{"non-numeric marker", "/api/spectrum/abc"},
		{"zero marker", "/api/spectrum/0"},
		{"negative marker", "/api/spectrum/-1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleGPTRadiation_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/gpt/radiation?lat=35&lon=139", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleGPTRadiation_BadRequest(t *testing.T) {
	mux := newRESTMux(t)
	cases := []struct {
		name string
		url  string
	}{
		{"missing lat", "/api/gpt/radiation?lon=139"},
		{"invalid lat", "/api/gpt/radiation?lat=bad&lon=139"},
		{"lat out of range", "/api/gpt/radiation?lat=99&lon=139"},
		{"invalid lon", "/api/gpt/radiation?lat=35&lon=bad"},
		{"lon out of range", "/api/gpt/radiation?lat=35&lon=200"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleGPTArea_MethodNotAllowed(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/gpt/area?min_lat=1&max_lat=2&min_lon=3&max_lon=4", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleGPTArea_BadRequest(t *testing.T) {
	mux := newRESTMux(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/gpt/area?min_lat=1&max_lat=2", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}
