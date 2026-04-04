// Semantic cache and RAG context layer for the web-chat handler.
//
// Flow:
//  1. Embed the incoming question (OpenAI text-embedding-3-small).
//  2. If a past Q&A has cosine similarity >= cacheHitThreshold AND positive
//     user feedback → return the cached answer immediately (no LLM call).
//  3. Otherwise inject the top-k most similar past Q&A pairs plus any curated
//     location knowledge into the system prompt as RAG context.
//  4. After the LLM responds, store the new Q&A + embedding asynchronously.
//  5. When the user gives thumbs-up, increment feedback_score; if coordinates
//     are detected in the answer, also populate location_knowledge.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// cacheHitThreshold: similarity above this + positive feedback → skip LLM.
	// Lower than neural-embedding threshold because feature-hash cosine scores
	// are sparser; 0.85 corresponds roughly to very similar phrasing.
	cacheHitThreshold = float32(0.85)
	// ragContextThreshold: similarity above this → include in RAG context.
	ragContextThreshold = float32(0.50)
	// ragTopK: maximum number of similar Q&A pairs to inject.
	ragTopK = 3
)

// qaEntry is a row from qa_embeddings.
type qaEntry struct {
	ID            int64 // internal row id (time.Now().UnixNano())
	ChatID        int64 // chat_id sent to the frontend for feedback linkage
	Question      string
	Answer        string
	Embedding     []float32
	FeedbackScore int
}

// trackIDRegexp matches 5–10 character alphanumeric track IDs as whole words.
var trackIDRegexp = regexp.MustCompile(`\b([A-Za-z0-9]{5,10})\b`)

// extractTrackID returns the first token in s that looks like a Safecast track
// ID (5–10 alphanumeric chars, mixed case). Returns "" if none found.
func extractTrackID(s string) string {
	for _, m := range trackIDRegexp.FindAllString(s, -1) {
		// Skip common English words that would match the pattern.
		lower := strings.ToLower(m)
		if lower == "track" || lower == "about" || lower == "where" ||
			lower == "what" || lower == "which" || lower == "japan" ||
			lower == "korea" || lower == "china" || lower == "india" {
			continue
		}
		// Track IDs contain both letters and digits.
		hasLetter := strings.ContainsAny(m, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		hasDigit := strings.ContainsAny(m, "0123456789")
		if hasLetter && hasDigit {
			return m
		}
	}
	return ""
}

// checkSemanticCache looks for a positively-rated cached answer with high
// cosine similarity to the given embedding.
// trackID is the explicit track the user is viewing (may be ""); question is
// used as a fallback to extract a track ID when trackID is empty. Cache hits
// from a different track are never returned. Returns ("", 0) on miss.
func checkSemanticCache(embedding []float32, question, trackID string) (answer string, chatID int64) {
	if !duckDBAvailable() || len(embedding) == 0 {
		return "", 0
	}
	entries, err := loadQAEmbeddings(true)
	if err != nil {
		log.Printf("semantic cache load: %v", err)
		return "", 0
	}

	// Prefer the explicit track ID; fall back to extracting from question text.
	queryTrackID := trackID
	if queryTrackID == "" {
		queryTrackID = extractTrackID(question)
	}

	var bestScore float32
	var best *qaEntry
	for i := range entries {
		if queryTrackID != "" {
			if !strings.Contains(entries[i].Question, queryTrackID) &&
				!strings.Contains(entries[i].Answer, queryTrackID) {
				continue // cached entry is about a different track
			}
		}
		if s := cosineSimilarity(embedding, entries[i].Embedding); s > bestScore {
			bestScore = s
			best = &entries[i]
		}
	}
	if bestScore >= cacheHitThreshold && best != nil {
		return best.Answer, best.ChatID
	}
	return "", 0
}

// buildRAGContext returns a formatted string of the top-k most similar past
// Q&A pairs (similarity >= ragContextThreshold) to prepend to the system prompt.
func buildRAGContext(embedding []float32) string {
	if !duckDBAvailable() || len(embedding) == 0 {
		return ""
	}
	entries, err := loadQAEmbeddings(false)
	if err != nil {
		log.Printf("RAG context load: %v", err)
		return ""
	}

	type scored struct {
		e     *qaEntry
		score float32
	}
	var candidates []scored
	for i := range entries {
		if s := cosineSimilarity(embedding, entries[i].Embedding); s >= ragContextThreshold {
			candidates = append(candidates, scored{&entries[i], s})
		}
	}
	if len(candidates) == 0 {
		return ""
	}

	// Partial sort: bring top-k to the front.
	limit := ragTopK
	if len(candidates) < limit {
		limit = len(candidates)
	}
	for i := 0; i < limit; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].score > candidates[i].score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("RELEVANT PAST Q&A (use as context, do not quote verbatim):\n")
	for _, c := range candidates[:limit] {
		ans := c.e.Answer
		if len(ans) > 500 {
			ans = ans[:500] + "…"
		}
		sb.WriteString(fmt.Sprintf("Q: %s\nA: %s\n\n", c.e.Question, ans))
	}
	return sb.String()
}

// getLocationKnowledge returns curated geographic notes from location_knowledge.
// Currently returns all notes (capped at 20); future: filter by detected coordinates.
func getLocationKnowledge() string {
	if !duckDBAvailable() {
		return ""
	}
	rows, err := duckDB.Query(`SELECT note FROM location_knowledge ORDER BY created_at DESC LIMIT 20`)
	if err != nil {
		return ""
	}
	defer rows.Close()

	var notes []string
	for rows.Next() {
		var note string
		if rows.Scan(&note) == nil && note != "" {
			if len(note) > 300 {
				note = note[:300] + "…"
			}
			notes = append(notes, note)
		}
	}
	if len(notes) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("LOCATION KNOWLEDGE BASE (curated context):\n")
	for _, n := range notes {
		sb.WriteString("- " + n + "\n")
	}
	return sb.String()
}

// storeQAEmbeddingAsync saves a new Q&A + embedding in the background.
// embeddingChatID is the ID to use as chat_id (for later feedback linkage).
func storeQAEmbeddingAsync(ctx context.Context, embeddingChatID int64, question, answer string, embedding []float32) {
	go func() {
		if !duckDBAvailable() || len(embedding) == 0 {
			return
		}
		embJSON, err := json.Marshal(embedding)
		if err != nil {
			return
		}
		id := time.Now().UnixNano()
		if _, err := duckDB.Exec(
			`INSERT INTO qa_embeddings (id, chat_id, question, answer, embedding, feedback_score) VALUES (?, ?, ?, ?, ?, 0)`,
			id, embeddingChatID, question, answer, string(embJSON),
		); err != nil {
			log.Printf("store qa_embedding: %v", err)
		}
	}()
}

// RecordFeedback updates the feedback score for a Q&A entry identified by chat_id.
// On positive feedback it also attempts to extract location knowledge from the answer.
func RecordFeedback(chatID int64, score int) error {
	if !duckDBAvailable() {
		return fmt.Errorf("analytics not available")
	}
	if score != 1 && score != -1 {
		return fmt.Errorf("score must be +1 or -1")
	}
	// Always record in chat_feedback — this is what the admin page reads.
	if _, err := duckDB.Exec(
		`INSERT INTO chat_feedback (chat_id, score) VALUES (?, ?)`,
		chatID, score,
	); err != nil {
		return fmt.Errorf("chat_feedback insert: %w", err)
	}
	// Best-effort update of semantic cache score (may have no matching row).
	if _, err := duckDB.Exec(
		`UPDATE qa_embeddings SET feedback_score = feedback_score + ? WHERE chat_id = ?`,
		score, chatID,
	); err != nil {
		log.Printf("qa_embeddings feedback update (non-fatal): %v", err)
	}
	if score > 0 {
		go extractLocationKnowledge(chatID)
	}
	return nil
}

// coordRegexp matches pairs of decimal numbers that look like lat/lon.
// It handles both compact forms ("45.4248, -75.094") and the verbose
// "Latitude: 45.4248..., Longitude: -75.094..." format the LLM often produces.
var coordRegexp = regexp.MustCompile(`(-?\d{1,3}\.\d{3,})[^\d\-\n]{0,20}?(-?\d{1,3}\.\d{3,})`)

// extractLocationKnowledge looks for lat/lon coordinates in the answer and
// stores a note in location_knowledge so future questions about that area
// automatically receive the context.
func extractLocationKnowledge(chatID int64) {
	if !duckDBAvailable() {
		return
	}
	var answer string
	row := duckDB.QueryRow(`SELECT answer FROM qa_embeddings WHERE chat_id = ? LIMIT 1`, chatID)
	if err := row.Scan(&answer); err != nil {
		return
	}
	matches := coordRegexp.FindStringSubmatch(answer)
	if len(matches) < 3 {
		return
	}
	lat, err1 := strconv.ParseFloat(matches[1], 64)
	lon, err2 := strconv.ParseFloat(matches[2], 64)
	if err1 != nil || err2 != nil {
		return
	}
	if math.Abs(lat) > 90 || math.Abs(lon) > 180 {
		return
	}
	id := time.Now().UnixNano()
	if _, err := duckDB.Exec(
		`INSERT INTO location_knowledge (id, lat, lon, radius_m, note, source_chat_id) VALUES (?, ?, ?, 1000, ?, ?)`,
		id, lat, lon, answer, chatID,
	); err != nil {
		log.Printf("location_knowledge insert: %v", err)
	}
}

// loadQAEmbeddingsForTrack fetches positively-rated Q&A rows that mention
// the given track ID in the question or answer text.
func loadQAEmbeddingsForTrack(trackID string) ([]qaEntry, error) {
	like := "%" + trackID + "%"
	rows, err := duckDB.Query(
		`SELECT id, chat_id, question, answer, embedding, feedback_score FROM qa_embeddings
		 WHERE feedback_score > 0
		 AND (question LIKE ? OR answer LIKE ?)`,
		like, like,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []qaEntry
	for rows.Next() {
		var e qaEntry
		var embJSON string
		if err := rows.Scan(&e.ID, &e.ChatID, &e.Question, &e.Answer, &embJSON, &e.FeedbackScore); err != nil {
			continue
		}
		if err := json.Unmarshal([]byte(embJSON), &e.Embedding); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// loadQAEmbeddings fetches Q&A rows from DuckLake.
// If positiveOnly is true, only rows with feedback_score > 0 are returned.
func loadQAEmbeddings(positiveOnly bool) ([]qaEntry, error) {
	query := `SELECT id, chat_id, question, answer, embedding, feedback_score FROM qa_embeddings`
	if positiveOnly {
		query += ` WHERE feedback_score > 0`
	}
	rows, err := duckDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []qaEntry
	for rows.Next() {
		var e qaEntry
		var embJSON string
		if err := rows.Scan(&e.ID, &e.ChatID, &e.Question, &e.Answer, &embJSON, &e.FeedbackScore); err != nil {
			continue
		}
		if err := json.Unmarshal([]byte(embJSON), &e.Embedding); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// enrichSystemPrompt prepends RAG context and location knowledge to the base
// system prompt. Returns the original prompt unchanged if nothing is available.
func enrichSystemPrompt(base string, ragContext, locationKnowledge string) string {
	var prefix strings.Builder
	if ragContext != "" {
		prefix.WriteString(ragContext)
		prefix.WriteString("\n")
	}
	if locationKnowledge != "" {
		prefix.WriteString(locationKnowledge)
		prefix.WriteString("\n")
	}
	if prefix.Len() == 0 {
		return base
	}
	return prefix.String() + base
}
