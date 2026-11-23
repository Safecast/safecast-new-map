package main

import (
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"

	"safecast-new-map/pkg/auth"
)

// One vote per client per answer, so a thumbs-up count means distinct people.
// In-memory and bounded; resets on restart, which is acceptable for a
// spam brake (the feedback limiter and admin approval remain the backstops).
const maxVoteKeys = 200000

// Feedback requests are tiny; anything bigger is not a real vote.
const maxFeedbackBodyBytes = 4096

var (
	voteMu   sync.Mutex
	voteSeen = map[string]struct{}{}

	feedbackLimiterOnce sync.Once
	feedbackLimiterInst *chatLimiter
)

// firstVote reports whether this client has not yet voted on chatID.
func firstVote(chatID int64, clientKey string) bool {
	k := strconv.FormatInt(chatID, 10) + "|" + clientKey
	voteMu.Lock()
	defer voteMu.Unlock()
	if _, ok := voteSeen[k]; ok {
		return false
	}
	if len(voteSeen) >= maxVoteKeys {
		voteSeen = map[string]struct{}{}
	}
	voteSeen[k] = struct{}{}
	return true
}

// feedbackLimiter is shared by every route that serves handleFeedback, so
// registering the handler twice does not double the allowance.
// FEEDBACK_LIMIT_PER_MINUTE, FEEDBACK_LIMIT_PER_DAY, FEEDBACK_LIMIT_GLOBAL_PER_DAY; 0 disables each.
func feedbackLimiter() *chatLimiter {
	feedbackLimiterOnce.Do(func() {
		feedbackLimiterInst = newChatLimiter(
			getEnvIntOrDefault("FEEDBACK_LIMIT_PER_MINUTE", 10),
			getEnvIntOrDefault("FEEDBACK_LIMIT_PER_DAY", 100),
			getEnvIntOrDefault("FEEDBACK_LIMIT_GLOBAL_PER_DAY", 10000),
		)
		log.Printf("[feedback] limits: %d/min and %d/day per client, %d/day total (0 = off)",
			feedbackLimiterInst.perMinute, feedbackLimiterInst.perDay, feedbackLimiterInst.globalPerDay)
	})
	return feedbackLimiterInst
}

// allowFeedback applies the rate limit and body cap. It writes a 429 and
// returns false when the client is over its limit.
func allowFeedback(l *chatLimiter, w http.ResponseWriter, r *http.Request, clientKey string) bool {
	ok, retryAfter, reason := l.allow(clientKey)
	if !ok {
		if reason == "global" {
			log.Printf("[feedback] global daily limit reached (%d)", l.globalPerDay)
		}
		secs := int(math.Ceil(retryAfter.Seconds()))
		if secs < 1 {
			secs = 1
		}
		w.Header().Set("Retry-After", strconv.Itoa(secs))
		writeError(w, http.StatusTooManyRequests, "too many feedback requests")
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxFeedbackBodyBytes)
	return true
}

// autoApproveScore is the thumbs-up total at which a pending answer becomes
// active without an admin. CHAT_AUTO_APPROVE_SCORE; 0 disables (default 3).
func autoApproveScore() int {
	if v := os.Getenv("CHAT_AUTO_APPROVE_SCORE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return n
		}
	}
	return 3
}

// autoApproveThreshold is the score that auto-approves an answer, or 0 when
// auto-approve is off. Without CHAT_ORIGIN_SECRET the client key comes from
// X-Forwarded-For, which any visitor can set, so "distinct voters" cannot be
// told apart and vote-based approval stays off.
func autoApproveThreshold() int {
	if os.Getenv("CHAT_ORIGIN_SECRET") == "" {
		return 0
	}
	return autoApproveScore()
}

// logAutoApproveMode reports at startup whether vote-based approval is active.
func logAutoApproveMode() {
	switch {
	case autoApproveScore() == 0:
		log.Println("[feedback] auto-approve off (CHAT_AUTO_APPROVE_SCORE=0); admin approval only")
	case os.Getenv("CHAT_ORIGIN_SECRET") == "":
		log.Println("[feedback] auto-approve off: CHAT_ORIGIN_SECRET not set, so voters can't be told apart; admin approval only")
	default:
		log.Printf("[feedback] auto-approve at score %d", autoApproveScore())
	}
}

func isAdminRequest(r *http.Request) bool {
	u, ok := auth.GetUserFromContext(r.Context())
	return ok && u.IsAdmin
}

// maybeAutoApprove activates a pending answer after an admin thumbs-up, or
// once its score reaches the threshold.
func maybeAutoApprove(chatID int64, byAdmin bool) {
	if !duckDBAvailable() {
		return
	}
	var err error
	if byAdmin {
		_, err = duckDB.Exec(`UPDATE qa_embeddings SET status='active' WHERE chat_id=? AND status='pending'`, chatID)
	} else if th := autoApproveThreshold(); th > 0 {
		_, err = duckDB.Exec(`UPDATE qa_embeddings SET status='active' WHERE chat_id=? AND status='pending' AND feedback_score >= ?`, chatID, th)
	}
	if err != nil {
		log.Printf("auto-approve (non-fatal): %v", err)
	}
}
