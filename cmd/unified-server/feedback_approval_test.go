package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFirstVote(t *testing.T) {
	if !firstVote(9001, "1.2.3.4") {
		t.Fatal("first vote should count")
	}
	if firstVote(9001, "1.2.3.4") {
		t.Fatal("repeat vote should not count")
	}
	if !firstVote(9001, "5.6.7.8") || !firstVote(9002, "1.2.3.4") {
		t.Fatal("different client or answer should count")
	}
}

func TestAutoApproveScoreDefault(t *testing.T) {
	t.Setenv("CHAT_AUTO_APPROVE_SCORE", "")
	if autoApproveScore() != 3 {
		t.Fatal("default should be 3")
	}
	t.Setenv("CHAT_AUTO_APPROVE_SCORE", "0")
	if autoApproveScore() != 0 {
		t.Fatal("0 disables")
	}
}

func TestAutoApproveNeedsOriginSecret(t *testing.T) {
	t.Setenv("CHAT_AUTO_APPROVE_SCORE", "")
	t.Setenv("CHAT_ORIGIN_SECRET", "")
	if th := autoApproveThreshold(); th != 0 {
		t.Fatalf("without an origin secret voters can be spoofed; threshold %d, want 0 (off)", th)
	}
	t.Setenv("CHAT_ORIGIN_SECRET", "s3cret")
	if th := autoApproveThreshold(); th != 3 {
		t.Fatalf("with an origin secret: threshold %d, want 3", th)
	}
	t.Setenv("CHAT_AUTO_APPROVE_SCORE", "0")
	if th := autoApproveThreshold(); th != 0 {
		t.Fatalf("CHAT_AUTO_APPROVE_SCORE=0 must stay off; got %d", th)
	}
}

func TestAllowFeedbackRateLimits(t *testing.T) {
	l, _ := newTestLimiter(2, 0, 0)
	send := func(key string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{}`))
		allowFeedback(l, rec, r, key)
		return rec
	}
	send("a")
	send("a")
	rec := send("a")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("3rd vote in a minute: got %d, want 429", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Error("missing Retry-After header")
	}
	if rec := send("b"); rec.Code == http.StatusTooManyRequests {
		t.Fatal("another client should not be limited")
	}
}

func TestAllowFeedbackCapsBody(t *testing.T) {
	l, _ := newTestLimiter(0, 0, 0)
	r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(strings.Repeat("x", maxFeedbackBodyBytes+1)))
	if !allowFeedback(l, httptest.NewRecorder(), r, "a") {
		t.Fatal("should be allowed")
	}
	if _, err := io.ReadAll(r.Body); err == nil {
		t.Fatal("reading a body over the cap should fail")
	}
}

func TestHandleFeedbackRejectsBadScore(t *testing.T) {
	h := handleFeedback()
	for _, body := range []string{`{"chat_id": 7, "score": 5}`, `{"chat_id": 7, "score": 0}`, `{"chat_id": 7}`} {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(body))
		r.RemoteAddr = "192.0.2.77:1000"
		h(rec, r)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%s: got %d, want 400", body, rec.Code)
		}
	}
	// A rejected score must not use up the client's one vote on that answer.
	if !firstVote(7, "192.0.2.77") {
		t.Error("invalid score consumed the vote slot")
	}
}
