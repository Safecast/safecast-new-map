package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"safecast-new-map/pkg/auth"
)

// One vote per client per answer, so a thumbs-up count means distinct people.
// In-memory and bounded; resets on restart, which is acceptable for a
// spam brake (the global chat limiter and admin approval remain the backstops).
const maxVoteKeys = 200000

var (
	voteMu   sync.Mutex
	voteSeen = map[string]struct{}{}
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
	} else if th := autoApproveScore(); th > 0 {
		_, err = duckDB.Exec(`UPDATE qa_embeddings SET status='active' WHERE chat_id=? AND status='pending' AND feedback_score >= ?`, chatID, th)
	}
	if err != nil {
		log.Printf("auto-approve (non-fatal): %v", err)
	}
}
