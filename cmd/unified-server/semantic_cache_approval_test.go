//go:build duckdb

package main

import (
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// setupTestDuckDB opens a throwaway local DuckDB with the analytics schema.
func setupTestDuckDB(t *testing.T) {
	t.Helper()
	t.Setenv("DUCKLAKE_PG_URL", "")
	t.Setenv("DUCKDB_LOCAL_PATH", filepath.Join(t.TempDir(), "analytics.duckdb"))
	if err := initDuckDBAnalytics(); err != nil {
		t.Fatalf("init duckdb: %v", err)
	}
	t.Cleanup(func() {
		duckDB.Close()
		duckDB = nil
	})
}

// waitForQARow waits for the async insert from storeQAEmbeddingAsync.
func waitForQARow(t *testing.T, chatID int64) string {
	t.Helper()
	for i := 0; i < 100; i++ {
		var status string
		if duckDB.QueryRow(`SELECT status FROM qa_embeddings WHERE chat_id = ?`, chatID).Scan(&status) == nil {
			return status
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("qa_embeddings row for chat %d never appeared", chatID)
	return ""
}

func TestNewAnswersAreNotReusedUntilApproved(t *testing.T) {
	setupTestDuckDB(t)

	const chatID = int64(4242)
	question := "what is the radiation level near the station in testville"
	answer := "Planted answer. Location: 37.4210, 141.0330"
	emb := []float32{1, 0, 0, 0}

	storeQAEmbeddingAsync(nil, chatID, question, answer, emb, "en")
	if status := waitForQARow(t, chatID); status != "pending" {
		t.Fatalf("new answer status = %q, want pending", status)
	}

	// A visitor's thumbs-up must not make it reusable.
	if err := RecordFeedback(chatID, 1); err != nil {
		t.Fatalf("feedback: %v", err)
	}
	time.Sleep(200 * time.Millisecond) // extractLocationKnowledge runs async

	if got, _ := checkSemanticCache(emb, question, "", "en"); got != "" {
		t.Fatalf("pending answer served from cache: %q", got)
	}
	if ctx := buildRAGContext(emb); strings.Contains(ctx, "Planted") {
		t.Fatalf("pending answer injected as RAG context: %q", ctx)
	}
	if notes := getLocationKnowledge(); strings.Contains(notes, "Planted") {
		t.Fatalf("location note from pending answer injected: %q", notes)
	}

	// Admin approves (Promote sets status = 'active').
	if _, err := duckDB.Exec(`UPDATE qa_embeddings SET feedback_score = feedback_score + 1, status = 'active' WHERE chat_id = ?`, chatID); err != nil {
		t.Fatalf("approve: %v", err)
	}
	if got, _ := checkSemanticCache(emb, question, "", "en"); got != answer {
		t.Fatalf("approved answer not served from cache, got %q", got)
	}
	if ctx := buildRAGContext(emb); !strings.Contains(ctx, "Planted") {
		t.Fatal("approved answer missing from RAG context")
	}
	if notes := getLocationKnowledge(); !strings.Contains(notes, "Planted") {
		t.Fatal("location note from approved answer missing")
	}
}

func TestAdminAddedLocationNotesAreUsed(t *testing.T) {
	setupTestDuckDB(t)
	if _, err := duckDB.Exec(`INSERT INTO location_knowledge (id, lat, lon, radius_m, note, source_chat_id) VALUES (1, 37.4, 141.0, 1000, 'Admin note', NULL)`); err != nil {
		t.Fatalf("insert: %v", err)
	}
	if notes := getLocationKnowledge(); !strings.Contains(notes, "Admin note") {
		t.Fatalf("admin-added note missing: %q", notes)
	}
}

func TestNewQAIDIsUniqueAndJSSafe(t *testing.T) {
	seen := make(map[int64]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 2000; i++ {
				id := newQAID()
				mu.Lock()
				if seen[id] {
					mu.Unlock()
					t.Errorf("duplicate id %d", id)
					return
				}
				seen[id] = true
				mu.Unlock()
				if id > 1<<53 {
					t.Errorf("id %d exceeds JS MAX_SAFE_INTEGER", id)
					return
				}
			}
		}()
	}
	wg.Wait()
}
