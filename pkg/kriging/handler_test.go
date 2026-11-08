package kriging

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGridHasSupportedAndBlankCells(t *testing.T) {
	// Points sit in the southwest of a small box. The northeast is farther than
	// the practical range, so those cells stay blank.
	when := int64(1_700_000_000)
	var obs []Observation
	for i := 0; i < 12; i++ {
		obs = append(obs, Observation{
			Lat: 0.001, Lon: 0.001 + float64(i)*0.0002, TimeUnix: when, DoseRate: 0.15, TrackID: "t",
		})
	}
	h := &Handler{real: NewPredictor(memorySource{obs: obs}, DefaultModel())}
	url := "/api/kriging/grid?cols=6&rows=6&time=1700000000&minLat=0&minLon=0&maxLat=0.02&maxLon=0.02"
	rec := httptest.NewRecorder()
	h.ServeGrid(rec, httptest.NewRequest(http.MethodGet, url, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Model          string `json:"model"`
		ValueUnit      string `json:"valueUnit"`
		CandidateCount int    `json:"candidateCount"`
		Note           string `json:"note"`
		Cells          []Cell `json:"cells"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Model != ModelName || body.ValueUnit != "uSv/h" || body.CandidateCount < 10 {
		t.Fatalf("unexpected meta: %+v", body)
	}
	if body.Note == "" {
		t.Fatal("missing note")
	}
	drawn, blank := 0, 0
	for _, cell := range body.Cells {
		if cell.Suppressed || cell.MedianUSvH == nil {
			blank++
			continue
		}
		if *cell.MedianUSvH <= 0 || math.IsNaN(cell.LogVariance) {
			t.Fatalf("bad drawn cell %+v", cell)
		}
		drawn++
	}
	if drawn == 0 || blank == 0 {
		t.Fatalf("drawn %d blank %d, want both", drawn, blank)
	}
}

func TestGridRefusesEmptyAndHugeBox(t *testing.T) {
	h := &Handler{real: NewPredictor(memorySource{}, DefaultModel())}
	rec := httptest.NewRecorder()
	h.ServeGrid(rec, httptest.NewRequest(http.MethodGet, "/api/kriging/grid?minLat=1&minLon=0&maxLat=1&maxLon=1", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("empty box status %d body %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	h.ServeGrid(rec, httptest.NewRequest(http.MethodGet, "/api/kriging/grid?minLat=0&minLon=0&maxLat=10&maxLon=10", nil))
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("huge box status %d body %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["refusal"] != "viewport_too_wide" {
		t.Fatalf("refusal %v", body["refusal"])
	}
}

func TestGridRefusesTooMany(t *testing.T) {
	obs := make([]Observation, DefaultMaxObservations+1)
	for i := range obs {
		obs[i] = Observation{Lat: 0.001, Lon: 0.001, TimeUnix: 50, DoseRate: 0.2, TrackID: "t"}
	}
	h := &Handler{real: NewPredictor(memorySource{obs: obs}, DefaultModel())}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/kriging/grid?minLat=0&minLon=0&maxLat=0.01&maxLon=0.01&time=50&cols=2&rows=2", nil)
	h.ServeGrid(rec, req)
	if rec.Code != http.StatusUnprocessableEntity || !strings.Contains(rec.Body.String(), "too_many_observations") {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
}

func TestOverlayScript(t *testing.T) {
	h := &Handler{}
	rec := httptest.NewRecorder()
	h.ServeOverlay(rec, httptest.NewRequest(http.MethodGet, "/kriging/overlay.js", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "SafecastKriging") {
		t.Fatal("overlay script missing attach point")
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "javascript") {
		t.Fatalf("content type %s", rec.Header().Get("Content-Type"))
	}
}
