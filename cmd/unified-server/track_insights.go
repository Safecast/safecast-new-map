// track_insights.go — GET /api/track/{id}/insights
//
// Returns cached RAG responses (Q&A pairs with positive feedback) that are
// spatially or semantically relevant to the given track, plus curated
// location_knowledge notes that fall within the track's bounding box.
//
// Relevance ranking:
//  1. Hard filter: only qa_embeddings with feedback_score > 0.
//  2. Semantic similarity: embed "radiation measurements near <centroid>"
//     and cosine-compare against stored embeddings (threshold = ragContextThreshold).
//  3. Return top 5 by similarity score.
//  4. Location notes: location_knowledge rows whose (lat, lon) falls inside bbox.
//
// Registration: http.HandleFunc("GET /api/track/{id}/insights", trackInsightsHandler)
// in mcp_register.go, registered on http.DefaultServeMux (port 8765).

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"safecast-new-map/pkg/database"
)

// trackInsight is a single cached Q&A relevant to the track.
type trackInsight struct {
	ChatID          int64   `json:"chat_id"`
	Question        string  `json:"question"`
	Answer          string  `json:"answer"`
	SimilarityScore float32 `json:"similarity_score"`
	FeedbackScore   int     `json:"feedback_score"`
}

// trackLocationNote is a curated geographic note near the track.
type trackLocationNote struct {
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Note          string  `json:"note"`
	SourceChatID  int64   `json:"source_chat_id,omitempty"`
	FeedbackScore int     `json:"feedback_score"`
}

// trackInsightsResponse is the JSON payload for GET /api/track/{id}/insights.
type trackInsightsResponse struct {
	TrackID                string                  `json:"track_id"`
	Insights               []trackInsight          `json:"insights"`
	LocationNotes          []trackLocationNote     `json:"location_notes"`
	AnomalySummary         *AnomalySummary         `json:"anomaly_summary,omitempty"`
	SpectrumAnomalySummary *SpectrumAnomalySummary `json:"spectrum_anomaly_summary,omitempty"`
}

// trackInsightsHandler serves GET /api/track/{id}/insights.
// Registered on http.DefaultServeMux via mcp_register.go using Go 1.22 pattern routing.
func trackInsightsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	trackID := r.PathValue("id")
	if trackID == "" {
		writeError(w, http.StatusBadRequest, "track id required")
		return
	}

	if !dbAvailable() {
		writeError(w, http.StatusServiceUnavailable, "database not available")
		return
	}

	// 1. Load bounding box from PostgreSQL.
	bounds, found, err := db.GetTrackBounds(ctx, trackID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load track bounds")
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "track not found")
		return
	}

	anomaly, _ := loadAnomalySummary(trackID)
	spectrumAnomaly, _ := loadSpectrumAnomalySummary(trackID)
	resp := trackInsightsResponse{
		TrackID:                trackID,
		Insights:               insightsByEmbedding(ctx, bounds, trackID),
		LocationNotes:          locationNotesInBbox(bounds, trackID),
		AnomalySummary:         anomaly,
		SpectrumAnomalySummary: spectrumAnomaly,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp) //nolint:errcheck
}

// insightsByEmbedding returns the top-5 positively-rated Q&A pairs whose
// embeddings are most similar to a query generated from the track centroid.
func insightsByEmbedding(ctx context.Context, bounds database.Bounds, trackID string) []trackInsight {
	out := []trackInsight{}
	if !duckDBAvailable() {
		return out
	}

	centerLat := (bounds.MinLat + bounds.MaxLat) / 2
	centerLon := (bounds.MinLon + bounds.MaxLon) / 2
	autoQuery := fmt.Sprintf("radiation measurements near %.3f %.3f", centerLat, centerLon)

	// getEmbedding ignores ctx (local feature-hash implementation).
	embedding, err := getEmbedding(ctx, autoQuery)
	if err != nil || len(embedding) == 0 {
		return out
	}

	// Load only Q&A that reference this specific track ID so the sidebar
	// doesn't show insights from unrelated tracks.
	entries, err := loadQAEmbeddingsForTrack(trackID)
	if err != nil {
		return out
	}

	type scored struct {
		e     qaEntry
		score float32
	}
	// Include ALL positively-rated entries; the centroid auto-query is too
	// generic to reach ragContextThreshold (0.50) against natural-language
	// questions. We still rank by similarity so the most topically relevant
	// results float to the top.
	var candidates []scored
	for _, e := range entries {
		candidates = append(candidates, scored{e, cosineSimilarity(embedding, e.Embedding)})
	}

	// Sort: primary = feedback_score descending, tie-break = cosine similarity descending.
	sort.Slice(candidates, func(i, j int) bool {
		fi, fj := candidates[i].e.FeedbackScore, candidates[j].e.FeedbackScore
		if fi != fj {
			return fi > fj
		}
		return candidates[i].score > candidates[j].score
	})

	limit := 5
	if len(candidates) < limit {
		limit = len(candidates)
	}

	for _, c := range candidates[:limit] {
		out = append(out, trackInsight{
			ChatID:          c.e.ChatID,
			Question:        c.e.Question,
			Answer:          c.e.Answer,
			SimilarityScore: c.score,
			FeedbackScore:   c.e.FeedbackScore,
		})
	}
	return out
}

// locationNotesInBbox returns location_knowledge rows whose point falls
// inside the track bounding box.
func locationNotesInBbox(bounds database.Bounds, trackID string) []trackLocationNote {
	out := []trackLocationNote{}
	if !duckDBAvailable() {
		return out
	}

	rows, err := duckDB.Query(
		`SELECT lk.lat, lk.lon, lk.note, lk.source_chat_id,
		        COALESCE(qe.feedback_score, 0) AS feedback_score
		 FROM location_knowledge lk
		 LEFT JOIN qa_embeddings qe ON qe.chat_id = lk.source_chat_id
		 WHERE lk.lat BETWEEN ? AND ? AND lk.lon BETWEEN ? AND ?
		 ORDER BY feedback_score DESC, lk.created_at DESC LIMIT 20`,
		bounds.MinLat, bounds.MaxLat, bounds.MinLon, bounds.MaxLon,
	)
	if err != nil {
		return out
	}
	defer rows.Close()

	for rows.Next() {
		var n trackLocationNote
		var srcID sql.NullInt64
		if rows.Scan(&n.Lat, &n.Lon, &n.Note, &srcID, &n.FeedbackScore) == nil {
			if srcID.Valid {
				n.SourceChatID = srcID.Int64
			}
			// Skip notes that mention a different track ID.
			if noteTrack := extractTrackID(n.Note); noteTrack != "" && noteTrack != trackID {
				continue
			}
			out = append(out, n)
		}
	}
	return out
}
