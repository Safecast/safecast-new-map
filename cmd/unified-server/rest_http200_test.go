package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
)

// radiationInfoCanonicalTopics matches referenceData / validTopics for /api/info/{topic}.
var radiationInfoCanonicalTopics = []string{
	"units",
	"dose_rates",
	"safety_levels",
	"detectors",
	"background_levels",
	"isotopes",
}

func topicPathVariants(canonical string) []string {
	seen := make(map[string]struct{})
	add := func(s string) {
		if s == "" {
			return
		}
		if _, ok := seen[s]; ok {
			return
		}
		seen[s] = struct{}{}
	}
	kebab := strings.ReplaceAll(canonical, "_", "-")
	add(canonical)
	add(kebab)
	add(strings.ToUpper(canonical))
	add(strings.ToUpper(kebab))
	if len(canonical) > 0 {
		add(strings.ToUpper(canonical[:1]) + canonical[1:])
	}
	if len(kebab) > 0 {
		add(strings.ToUpper(kebab[:1]) + kebab[1:])
	}
	// Mixed token casing for multi-part slugs (ASCII-only topics).
	if strings.Contains(canonical, "_") {
		parts := strings.Split(canonical, "_")
		mixed := make([]string, len(parts))
		for i, p := range parts {
			if i%2 == 0 {
				mixed[i] = strings.ToUpper(p)
			} else {
				mixed[i] = p
			}
		}
		add(strings.Join(mixed, "_"))
		add(strings.Join(mixed, "-"))
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func buildRadiationInfoHTTP200Cases() []struct {
	name string
	url  string
} {
	var pathSegs []string
	for _, topic := range radiationInfoCanonicalTopics {
		pathSegs = append(pathSegs, topicPathVariants(topic)...)
	}
	queries := []string{""}
	for i := 0; i < 64; i++ {
		queries = append(queries, fmt.Sprintf("?cache_bust=%d", i))
	}
	var cases []struct {
		name string
		url  string
	}
	n := 0
outer:
	for _, seg := range pathSegs {
		for _, q := range queries {
			if n >= 200 {
				break outer
			}
			url := "/api/info/" + seg + q
			safe := strings.Map(func(r rune) rune {
				switch {
				case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
					return r
				case r == '-' || r == '_':
					return r
				default:
					return '_'
				}
			}, url)
			cases = append(cases, struct {
				name string
				url  string
			}{name: fmt.Sprintf("%03d_%s", n, safe), url: url})
			n++
		}
	}
	return cases
}

func TestBuildRadiationInfoHTTP200Cases_Count(t *testing.T) {
	if got := len(buildRadiationInfoHTTP200Cases()); got != 200 {
		t.Fatalf("expected exactly 200 radiation info OK cases, got %d", got)
	}
}

// TestRESTRadiationInfo_SuiteHTTP200 exercises GET /api/info/{topic} success paths:
// multiple spellings (snake, kebab, case) and ignored query strings. Static content only; no DB.
func TestRESTRadiationInfo_SuiteHTTP200(t *testing.T) {
	mux := newRESTMux(t)
	for _, tc := range buildRadiationInfoHTTP200Cases() {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
			}
			ct := rec.Header().Get("Content-Type")
			if !strings.Contains(ct, "application/json") {
				t.Fatalf("expected JSON content-type, got %q", ct)
			}
			var body map[string]any
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("decode JSON: %v", err)
			}
			topic, ok := body["topic"].(string)
			if !ok || topic == "" {
				t.Fatalf("expected string topic in body, got %#v", body["topic"])
			}
			content, ok := body["content"].(string)
			if !ok || content == "" {
				t.Fatalf("expected non-empty string content, got %#v", body["content"])
			}
		})
	}
}
