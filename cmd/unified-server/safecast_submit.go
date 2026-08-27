package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
	"safecast-new-map/pkg/httpresp"
	safecastsubmit "safecast-new-map/pkg/safecast-submit"
)

// safecastSubmitter is the subset of safecastsubmit.Client used here, so tests
// can inject a fake without hitting the network.
type safecastSubmitter interface {
	ResolveUserID(ctx context.Context, apiKey string) (string, error)
	CheckExists(ctx context.Context, apiKey, safecastUserID, filename string) (bool, error)
	Submit(ctx context.Context, apiKey, filename string, content []byte) (string, error)
}

// safecastSubmitClient is the real client used in production; tests replace it.
var safecastSubmitClient safecastSubmitter = safecastsubmit.NewClient()

// submitToSafecastIfNeeded best-effort submits a bGeigie log to api.safecast.org
// on behalf of the uploading user, using their own Safecast API key. It never
// blocks or fails the local upload: it runs in its own goroutine (call it with
// "go"), every outcome is logged and recorded on the upload row (keyed by
// uploadID, since track_id can be shared across files or reused by duplicate
// detection), and errors are swallowed here.
func submitToSafecastIfNeeded(trackID string, uploadID int64, filename string, content []byte, safecastAPIKey, safecastUserID string, db *database.Database) {
	if safecastAPIKey == "" || safecastUserID == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	exists, err := safecastSubmitClient.CheckExists(ctx, safecastAPIKey, safecastUserID, filename)
	if err != nil {
		log.Printf("[safecast-submit] %s: check exists failed for %s: %v", trackID, filename, err)
		_ = db.UpdateUploadSafecastStatus(context.Background(), uploadID, "", err.Error())
		return
	}
	if exists {
		log.Printf("[safecast-submit] %s: %s already exists on api.safecast.org, skipping", trackID, filename)
		return
	}

	importID, err := safecastSubmitClient.Submit(ctx, safecastAPIKey, filename, content)
	if err != nil {
		log.Printf("[safecast-submit] %s: submit failed for %s: %v", trackID, filename, err)
		_ = db.UpdateUploadSafecastStatus(context.Background(), uploadID, "", err.Error())
		return
	}

	log.Printf("[safecast-submit] %s: submitted %s to api.safecast.org, import id %s", trackID, filename, importID)
	if err := db.UpdateUploadSafecastStatus(context.Background(), uploadID, importID, ""); err != nil {
		log.Printf("[safecast-submit] %s: failed to record submit status: %v", trackID, err)
	}
}

// safecastCredentialsRequest is the body for POST /api/user/safecast-credentials.
// An empty SafecastAPIKey clears the stored credentials (opt-out).
type safecastCredentialsRequest struct {
	SafecastAPIKey string `json:"safecast_api_key"`
}

// updateSafecastCredentialsHandler lets a logged-in user set (or clear) their
// api.safecast.org API key. On set, it resolves and caches the corresponding
// Safecast user id via GET /users/me.json so later uploads don't need to look
// it up on every submit.
func updateSafecastCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	if !httpresp.RequireMethodJSON(w, r, http.MethodPost) {
		return
	}

	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		httpresp.WriteUnauthorized(w, "Unauthorized")
		return
	}

	var req safecastCredentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresp.WriteBadRequest(w, "invalid_body", "Invalid request body")
		return
	}

	if req.SafecastAPIKey == "" {
		if err := auth.UpdateUserSafecastCredentials(r.Context(), db.DB, *dbType, user.ID, "", ""); err != nil {
			httpresp.WriteInternalError(w, "Failed to clear Safecast credentials")
			return
		}
		httpresp.WriteJSON(w, http.StatusOK, map[string]any{"safecast_api_key": "", "safecast_user_id": ""})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	safecastUserID, err := safecastSubmitClient.ResolveUserID(ctx, req.SafecastAPIKey)
	if err != nil {
		log.Printf("[safecast-submit] resolve user id failed for user %d: %v", user.ID, err)
		httpresp.WriteBadRequest(w, "invalid_api_key", "Could not verify this API key with api.safecast.org")
		return
	}

	if err := auth.UpdateUserSafecastCredentials(r.Context(), db.DB, *dbType, user.ID, req.SafecastAPIKey, safecastUserID); err != nil {
		httpresp.WriteInternalError(w, "Failed to save Safecast credentials")
		return
	}

	httpresp.WriteJSON(w, http.StatusOK, map[string]any{
		"safecast_api_key": req.SafecastAPIKey,
		"safecast_user_id": safecastUserID,
	})
}
