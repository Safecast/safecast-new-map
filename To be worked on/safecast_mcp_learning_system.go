// safecast_mcp_learning_system.go
// Pseudocode for adding a semantic cache + feedback loop to the Safecast MCP server
// This plugs into the existing web-chat handler in go/cmd/web-chat/

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"math"
	"time"
)

// ============================================================
// 1. DATA STRUCTURES
// ============================================================

// QARecord stores a question-answer pair with its embedding and feedback score
type QARecord struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Embedding []float64 `json:"embedding"` // 1536-dim from embedding model
	Score     int       `json:"score"`      // cumulative thumbs up (+1) / down (-1)
	UsedCount int       `json:"used_count"` // how often this was served from cache
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LocationKnowledge stores curated geo-annotations
// These are manually or semi-automatically created explanations
// for areas with notable radiation readings
type LocationKnowledge struct {
	ID          string    `json:"id"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	RadiusM     float64   `json:"radius_m"`
	Explanation string    `json:"explanation"` // e.g. "National Museum of Nuclear Science & History"
	Tags        []string  `json:"tags"`        // e.g. ["museum", "nuclear", "artifacts"]
	Source      string    `json:"source"`       // "user_feedback" | "manual" | "auto_detected"
	CreatedAt   time.Time `json:"created_at"`
}

// FeedbackEvent from the web chat UI
type FeedbackEvent struct {
	QuestionID string `json:"question_id"`
	Rating     int    `json:"rating"` // +1 or -1
	Comment    string `json:"comment,omitempty"`
}

// ============================================================
// 2. DUCKDB SCHEMA SETUP
// ============================================================

func initLearningTables(db *sql.DB) error {
	// DuckDB supports fixed-size arrays for embeddings
	// and has built-in list_cosine_similarity() for vector search
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS qa_cache (
			id            VARCHAR PRIMARY KEY,
			question      TEXT NOT NULL,
			answer        TEXT NOT NULL,
			embedding     DOUBLE[1536],
			score         INTEGER DEFAULT 0,
			used_count    INTEGER DEFAULT 0,
			created_at    TIMESTAMP DEFAULT current_timestamp,
			updated_at    TIMESTAMP DEFAULT current_timestamp
		);

		CREATE TABLE IF NOT EXISTS location_knowledge (
			id            VARCHAR PRIMARY KEY,
			lat           DOUBLE NOT NULL,
			lon           DOUBLE NOT NULL,
			radius_m      DOUBLE DEFAULT 1000,
			explanation   TEXT NOT NULL,
			tags          VARCHAR[],
			source        VARCHAR DEFAULT 'manual',
			created_at    TIMESTAMP DEFAULT current_timestamp
		);

		CREATE TABLE IF NOT EXISTS feedback_log (
			id            VARCHAR PRIMARY KEY,
			question_id   VARCHAR NOT NULL,
			rating        INTEGER NOT NULL,
			comment       TEXT,
			created_at    TIMESTAMP DEFAULT current_timestamp
		);
	`)
	return err
}

// ============================================================
// 3. EMBEDDING SERVICE
// ============================================================

// EmbeddingService wraps an embedding model API
// Could use OpenAI text-embedding-3-small, or a local model
// like sentence-transformers via a sidecar container
type EmbeddingService struct {
	apiURL string
	apiKey string
}

func (e *EmbeddingService) Embed(ctx context.Context, text string) ([]float64, error) {
	// Call embedding API
	// For cost efficiency with Safecast's volume:
	//   - OpenAI text-embedding-3-small: $0.02/1M tokens, 1536 dims
	//   - Or self-hosted: all-MiniLM-L6-v2 via Go (384 dims, free)
	//
	// Returns normalized vector for cosine similarity
	return nil, nil // placeholder
}

// CosineSimilarity computes similarity between two normalized vectors
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dot float64
	for i := range a {
		dot += a[i] * b[i]
	}
	return dot
}

// ============================================================
// 4. SEMANTIC CACHE
// ============================================================

type SemanticCache struct {
	db       *sql.DB
	embedder *EmbeddingService
	// Similarity threshold: 0.92 = very similar question
	// Tune this based on real-world performance
	threshold float64
}

func NewSemanticCache(db *sql.DB, embedder *EmbeddingService) *SemanticCache {
	return &SemanticCache{
		db:        db,
		embedder:  embedder,
		threshold: 0.92,
	}
}

// Lookup checks if a semantically similar question has been well-answered before
func (sc *SemanticCache) Lookup(ctx context.Context, question string) (*QARecord, error) {
	// 1. Embed the incoming question
	embedding, err := sc.embedder.Embed(ctx, question)
	if err != nil {
		return nil, err
	}

	// 2. Search DuckDB for similar questions with positive feedback
	// DuckDB has list_cosine_similarity() built in
	rows, err := sc.db.QueryContext(ctx, `
		SELECT id, question, answer, score, used_count,
			   list_cosine_similarity(embedding, $1::DOUBLE[1536]) as similarity
		FROM qa_cache
		WHERE score >= 1
		ORDER BY similarity DESC
		LIMIT 5
	`, embedding)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 3. Return best match above threshold
	for rows.Next() {
		var record QARecord
		var similarity float64
		if err := rows.Scan(&record.ID, &record.Question, &record.Answer,
			&record.Score, &record.UsedCount, &similarity); err != nil {
			continue
		}
		if similarity >= sc.threshold {
			// Update usage count
			sc.db.ExecContext(ctx, `
				UPDATE qa_cache SET used_count = used_count + 1, 
				updated_at = current_timestamp WHERE id = $1
			`, record.ID)
			return &record, nil
		}
	}
	return nil, nil // no cache hit
}

// TopKSimilar returns the k most similar past Q&A pairs for RAG context
// These don't need to be exact matches — they provide context
func (sc *SemanticCache) TopKSimilar(ctx context.Context, embedding []float64, k int) ([]QARecord, error) {
	rows, err := sc.db.QueryContext(ctx, `
		SELECT id, question, answer, score,
			   list_cosine_similarity(embedding, $1::DOUBLE[1536]) as similarity
		FROM qa_cache
		WHERE score >= 0
		ORDER BY similarity DESC
		LIMIT $2
	`, embedding, k)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []QARecord
	for rows.Next() {
		var r QARecord
		var sim float64
		if err := rows.Scan(&r.ID, &r.Question, &r.Answer, &r.Score, &sim); err != nil {
			continue
		}
		results = append(results, r)
	}
	return results, nil
}

// Store saves a new Q&A pair to the cache
func (sc *SemanticCache) Store(ctx context.Context, question, answer string, embedding []float64) error {
	id := generateID() // nanoid or UUID
	_, err := sc.db.ExecContext(ctx, `
		INSERT INTO qa_cache (id, question, answer, embedding)
		VALUES ($1, $2, $3, $4::DOUBLE[1536])
	`, id, question, answer, embedding)
	return err
}

// ============================================================
// 5. LOCATION KNOWLEDGE STORE
// ============================================================

type LocationStore struct {
	db *sql.DB
}

// FindNearby returns location knowledge entries near given coordinates
// This is where curated explanations like "Nuclear Museum" live
func (ls *LocationStore) FindNearby(ctx context.Context, lat, lon, radiusM float64) ([]LocationKnowledge, error) {
	// Simple bounding box filter + haversine for accuracy
	// PostGIS would be better but DuckDB spatial extension works too
	degRadius := radiusM / 111000.0 // rough meters to degrees

	rows, err := ls.db.QueryContext(ctx, `
		SELECT id, lat, lon, radius_m, explanation, tags, source
		FROM location_knowledge
		WHERE lat BETWEEN $1 AND $2
		  AND lon BETWEEN $3 AND $4
	`, lat-degRadius, lat+degRadius, lon-degRadius, lon+degRadius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []LocationKnowledge
	for rows.Next() {
		var lk LocationKnowledge
		if err := rows.Scan(&lk.ID, &lk.Lat, &lk.Lon, &lk.RadiusM,
			&lk.Explanation, &lk.Tags, &lk.Source); err != nil {
			continue
		}
		// Verify with haversine
		dist := haversine(lat, lon, lk.Lat, lk.Lon)
		if dist <= math.Max(lk.RadiusM, radiusM) {
			results = append(results, lk)
		}
	}
	return results, nil
}

// AddFromFeedback creates a location knowledge entry from user feedback
// This is how the system "learns" — when a user confirms an explanation
// about a location, it gets stored for future queries
func (ls *LocationStore) AddFromFeedback(ctx context.Context, lat, lon float64, explanation string) error {
	id := generateID()
	_, err := ls.db.ExecContext(ctx, `
		INSERT INTO location_knowledge (id, lat, lon, explanation, source)
		VALUES ($1, $2, $3, $4, 'user_feedback')
	`, id, lat, lon, explanation)
	return err
}

// ============================================================
// 6. ENHANCED CHAT HANDLER
// ============================================================

// This replaces or wraps the existing handleChat in go/cmd/web-chat/main.go

type LearningChatHandler struct {
	cache         *SemanticCache
	locations     *LocationStore
	embedder      *EmbeddingService
	claudeClient  interface{} // existing Anthropic client
}

func (h *LearningChatHandler) HandleChat(ctx context.Context, userQuestion string) (string, error) {
	// STEP 1: Embed the question
	embedding, err := h.embedder.Embed(ctx, userQuestion)
	if err != nil {
		// Fall through to normal LLM call if embedding fails
		return h.callClaude(ctx, userQuestion, nil, nil)
	}

	// STEP 2: Check semantic cache for high-confidence match
	cached, err := h.cache.Lookup(ctx, userQuestion)
	if err == nil && cached != nil {
		// Cache hit! Return the cached answer
		// Optionally prefix with a note that this is from prior knowledge
		return cached.Answer, nil
	}

	// STEP 3: No cache hit — gather RAG context
	// 3a. Get top-k similar past Q&A pairs
	similarQA, _ := h.cache.TopKSimilar(ctx, embedding, 3)

	// 3b. Extract any location coordinates from the question
	//     (use the LLM or regex to parse lat/lon from the question)
	lat, lon, hasLocation := extractLocation(userQuestion)

	var locationContext []LocationKnowledge
	if hasLocation {
		// 3c. Check location knowledge store
		locationContext, _ = h.locations.FindNearby(ctx, lat, lon, 2000)
	}

	// STEP 4: Build enriched system prompt with RAG context
	answer, err := h.callClaudeWithContext(ctx, userQuestion, similarQA, locationContext)
	if err != nil {
		return "", err
	}

	// STEP 5: Store the Q&A pair (score starts at 0, needs feedback to go positive)
	h.cache.Store(ctx, userQuestion, answer, embedding)

	return answer, nil
}

// callClaudeWithContext builds the system prompt with injected RAG context
func (h *LearningChatHandler) callClaudeWithContext(
	ctx context.Context,
	question string,
	similarQA []QARecord,
	locations []LocationKnowledge,
) (string, error) {

	// Build context injection block
	contextBlock := ""

	if len(similarQA) > 0 {
		contextBlock += "\n## Previously answered similar questions:\n"
		for _, qa := range similarQA {
			contextBlock += "Q: " + qa.Question + "\n"
			contextBlock += "A: " + qa.Answer + "\n"
			contextBlock += "---\n"
		}
	}

	if len(locations) > 0 {
		contextBlock += "\n## Known location information:\n"
		for _, loc := range locations {
			contextBlock += loc.Explanation + "\n"
		}
	}

	// Prepend to the existing system prompt
	systemPrompt := `You are the Safecast radiation data assistant.
` + contextBlock + `
Use the above context when relevant to provide better answers.
If a similar question was answered before, you can build on that answer.`

	// Call Claude API (existing pattern from web-chat)
	return h.callClaude(ctx, question, nil, nil)
}

// ============================================================
// 7. FEEDBACK ENDPOINT
// ============================================================

// HandleFeedback processes thumbs up/down from the web chat UI
func (h *LearningChatHandler) HandleFeedback(ctx context.Context, event FeedbackEvent) error {
	// 1. Update the Q&A record score
	_, err := h.cache.db.ExecContext(ctx, `
		UPDATE qa_cache
		SET score = score + $1, updated_at = current_timestamp
		WHERE id = $2
	`, event.Rating, event.QuestionID)
	if err != nil {
		return err
	}

	// 2. Log the feedback event
	_, err = h.cache.db.ExecContext(ctx, `
		INSERT INTO feedback_log (id, question_id, rating, comment)
		VALUES ($1, $2, $3, $4)
	`, generateID(), event.QuestionID, event.Rating, event.Comment)

	// 3. If positive feedback and the answer mentions a location,
	//    auto-create a location knowledge entry
	if event.Rating > 0 {
		// Retrieve the Q&A record
		var question, answer string
		h.cache.db.QueryRowContext(ctx, `
			SELECT question, answer FROM qa_cache WHERE id = $1
		`, event.QuestionID).Scan(&question, &answer)

		// Extract location if present and store as knowledge
		if lat, lon, ok := extractLocation(question); ok {
			h.locations.AddFromFeedback(ctx, lat, lon, answer)
		}
	}

	return err
}

// ============================================================
// 8. INTEGRATION WITH EXISTING MCP SERVER
// ============================================================

// The existing MCP server in go/cmd/web-chat/ needs these additions:
//
// 1. Add a POST /api/feedback endpoint for the thumbs up/down UI
//
// 2. Wrap the existing chat handler with LearningChatHandler
//
// 3. Add DuckDB tables on startup (initLearningTables)
//
// 4. Add an embedding service (start with OpenAI, can switch to local later)
//
// 5. Add thumbs up/down buttons to the web chat HTML template:
//
//    <div class="feedback">
//      <button onclick="sendFeedback(msgId, +1)">👍</button>
//      <button onclick="sendFeedback(msgId, -1)">👎</button>
//    </div>
//
// 6. Optional: Add an admin endpoint to manually add location knowledge
//    POST /api/locations { lat, lon, radius_m, explanation, tags }
//
// 7. Optional: Periodic job to prune low-scoring Q&A pairs:
//    DELETE FROM qa_cache WHERE score < -3 AND updated_at < now() - interval '30 days'

// ============================================================
// HELPER STUBS
// ============================================================

func generateID() string                                                     { return "" }
func haversine(lat1, lon1, lat2, lon2 float64) float64                       { return 0 }
func extractLocation(text string) (lat, lon float64, ok bool)                { return 0, 0, false }
func (h *LearningChatHandler) callClaude(ctx context.Context, q string, a interface{}, b interface{}) (string, error) {
	return "", nil
}
