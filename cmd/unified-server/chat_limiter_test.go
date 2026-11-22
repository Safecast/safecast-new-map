package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time          { return c.t }
func (c *fakeClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func newTestLimiter(perMinute, perDay, global int) (*chatLimiter, *fakeClock) {
	clock := &fakeClock{t: time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC)}
	l := newChatLimiter(perMinute, perDay, global)
	l.now = clock.now
	return l, clock
}

func TestChatLimiterPerMinute(t *testing.T) {
	l, clock := newTestLimiter(3, 0, 0)
	for i := 0; i < 3; i++ {
		if ok, _, _ := l.allow("a"); !ok {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	ok, retry, reason := l.allow("a")
	if ok || reason != "minute" {
		t.Fatalf("4th request in a minute: ok=%v reason=%q, want blocked by minute", ok, reason)
	}
	if retry <= 0 || retry > time.Minute {
		t.Fatalf("retry-after %v out of range", retry)
	}
	if ok, _, _ := l.allow("b"); !ok {
		t.Fatal("a different client should not be affected")
	}
	clock.advance(time.Minute)
	if ok, _, _ := l.allow("a"); !ok {
		t.Fatal("limit should reset after a minute")
	}
}

func TestChatLimiterPerDay(t *testing.T) {
	l, clock := newTestLimiter(0, 2, 0)
	l.allow("a")
	l.allow("a")
	if ok, _, reason := l.allow("a"); ok || reason != "day" {
		t.Fatalf("3rd request in a day: ok=%v reason=%q", ok, reason)
	}
	clock.advance(23 * time.Hour)
	if ok, _, _ := l.allow("a"); ok {
		t.Fatal("still inside the day window")
	}
	clock.advance(time.Hour)
	if ok, _, _ := l.allow("a"); !ok {
		t.Fatal("limit should reset after 24 hours")
	}
}

func TestChatLimiterGlobalCapIgnoresClientKey(t *testing.T) {
	l, clock := newTestLimiter(0, 0, 3)
	for _, k := range []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"} {
		if ok, _, _ := l.allow(k); !ok {
			t.Fatalf("%s should be allowed", k)
		}
	}
	if ok, _, reason := l.allow("4.4.4.4"); ok || reason != "global" {
		t.Fatalf("new client past the global cap: ok=%v reason=%q", ok, reason)
	}
	clock.advance(24 * time.Hour)
	if ok, _, _ := l.allow("4.4.4.4"); !ok {
		t.Fatal("global cap should reset after 24 hours")
	}
}

func TestChatLimiterBlockedRequestsDoNotCount(t *testing.T) {
	l, clock := newTestLimiter(1, 2, 0)
	l.allow("a")
	for i := 0; i < 5; i++ {
		l.allow("a") // blocked by the minute limit
	}
	clock.advance(time.Minute)
	if ok, _, _ := l.allow("a"); !ok {
		t.Fatal("blocked requests should not use up the daily allowance")
	}
}

func TestChatLimiterZeroDisables(t *testing.T) {
	l, _ := newTestLimiter(0, 0, 0)
	for i := 0; i < 1000; i++ {
		if ok, _, _ := l.allow("a"); !ok {
			t.Fatalf("request %d blocked with all limits off", i+1)
		}
	}
}

func TestChatLimiterSweepsOldClients(t *testing.T) {
	l, clock := newTestLimiter(10, 10, 0)
	l.allow("old")
	clock.advance(25 * time.Hour)
	l.allow("new")
	if _, ok := l.clients["old"]; ok {
		t.Fatal("expired client window should be removed")
	}
}

func TestChatClientKeyPrefersCloudFrontViewerAddress(t *testing.T) {
	cases := []struct {
		viewer, xff, want string
	}{
		{"203.0.113.7:51234", "198.51.100.1", "203.0.113.7"},
		{"2001:db8::1:46532", "198.51.100.1", "2001:db8::1"},
		{"", "198.51.100.1, 203.0.113.9", "198.51.100.1"},
	}
	for _, c := range cases {
		r := httptest.NewRequest(http.MethodPost, "/chat", nil)
		if c.viewer != "" {
			r.Header.Set("CloudFront-Viewer-Address", c.viewer)
		}
		r.Header.Set("X-Forwarded-For", c.xff)
		if got := chatClientKey(r); got != c.want {
			t.Errorf("viewer=%q xff=%q: got %q want %q", c.viewer, c.xff, got, c.want)
		}
	}
}

func TestWithChatLimitsReturns429WithErrorChunk(t *testing.T) {
	l, _ := newTestLimiter(1, 0, 0)
	calls := 0
	h := withChatLimits(l, 1024, func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	})
	send := func() *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(`{"message":"hi"}`))
		r.Header.Set("X-Forwarded-For", "198.51.100.1")
		rec := httptest.NewRecorder()
		h(rec, r)
		return rec
	}
	if rec := send(); rec.Code != http.StatusOK {
		t.Fatalf("first request: got %d", rec.Code)
	}
	rec := send()
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("second request: got %d, want 429", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Error("missing Retry-After header")
	}
	if !strings.Contains(rec.Body.String(), `"type":"error"`) {
		t.Errorf("body should be an error chunk, got %q", rec.Body.String())
	}
	if calls != 1 {
		t.Errorf("handler called %d times, want 1", calls)
	}
}

func TestWithChatLimitsPassesPreflight(t *testing.T) {
	l, _ := newTestLimiter(1, 1, 1)
	l.allow("x") // use up the global cap
	called := false
	h := withChatLimits(l, 1024, func(w http.ResponseWriter, r *http.Request) { called = true })
	h(httptest.NewRecorder(), httptest.NewRequest(http.MethodOptions, "/chat", nil))
	if !called {
		t.Fatal("OPTIONS preflight should reach the handler without using a slot")
	}
}

func TestWithChatLimitsCapsBodySize(t *testing.T) {
	l, _ := newTestLimiter(0, 0, 0)
	var readErr error
	h := withChatLimits(l, 16, func(w http.ResponseWriter, r *http.Request) {
		_, readErr = io.ReadAll(r.Body)
	})
	h(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(strings.Repeat("x", 100))))
	if readErr == nil {
		t.Fatal("reading a body over the cap should fail")
	}
}
