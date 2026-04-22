package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	modeladapter "safecast-new-map/cmd/unified-server/model-adapter"
)

// =============================================================================
// Request/response payloads
// =============================================================================

// aiHintPayload is the structured JSON exchanged with the admin UI. It mirrors
// the on-disk JSON files and the ai_hints row shape so a single payload can be
// used for both create and import.
type aiHintPayload struct {
	Model                 string                                  `json:"model,omitempty"`
	DisplayName           string                                  `json:"display_name"`
	Capabilities          []string                                `json:"capabilities"`
	SystemPrompt          string                                  `json:"system_prompt"`
	Tools                 map[string]*modeladapter.ToolHintJSON   `json:"tools"`
	GlobalFormattingRules map[string]string                       `json:"global_formatting_rules"`
}

type aiHintListItem struct {
	Model        string     `json:"model"`
	DisplayName  string     `json:"display_name"`
	Capabilities []string   `json:"capabilities"`
	UpdatedAt    *time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type aiHintFullRow struct {
	ID                    int64                                 `json:"id"`
	Model                 string                                `json:"model"`
	DisplayName           string                                `json:"display_name"`
	Capabilities          []string                              `json:"capabilities"`
	SystemPrompt          string                                `json:"system_prompt"`
	Tools                 map[string]*modeladapter.ToolHintJSON `json:"tools"`
	GlobalFormattingRules map[string]string                     `json:"global_formatting_rules"`
	CreatedAt             *time.Time                            `json:"created_at"`
	UpdatedAt             *time.Time                            `json:"updated_at"`
	DeletedAt             *time.Time                            `json:"deleted_at,omitempty"`
}

type aiHintHistoryItem struct {
	ID         int64      `json:"id"`
	Model      string     `json:"model"`
	ChangeKind string     `json:"change_kind"`
	ChangedAt  *time.Time `json:"changed_at"`
	Snapshot   json.RawMessage `json:"snapshot,omitempty"`
}

// =============================================================================
// Helpers
// =============================================================================

// modelSlugPattern: lowercase alphanumerics and dashes only; no leading/trailing dash.
func validateModelSlug(s string) error {
	if s == "" {
		return fmt.Errorf("model slug is required")
	}
	if len(s) > 64 {
		return fmt.Errorf("model slug too long (max 64 chars)")
	}
	for i, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			// ok
		case r == '-':
			if i == 0 || i == len(s)-1 {
				return fmt.Errorf("model slug cannot start or end with dash")
			}
		default:
			return fmt.Errorf("model slug may only contain lowercase letters, digits, and dashes")
		}
	}
	return nil
}

// modelFromPath extracts the {model} path parameter from URLs like
// /api/admin/ai-hints/claude or /api/admin/ai-hints/claude/history.
func modelFromPath(path string, basePath string, trailingSegments int) (string, []string, error) {
	stripped := strings.TrimPrefix(path, basePath)
	stripped = strings.Trim(stripped, "/")
	parts := strings.Split(stripped, "/")
	if len(parts) < 1 || parts[0] == "" {
		return "", nil, fmt.Errorf("missing model slug in path")
	}
	model := parts[0]
	if err := validateModelSlug(model); err != nil {
		return "", nil, err
	}
	rest := parts[1:]
	if len(rest) != trailingSegments {
		return "", nil, fmt.Errorf("unexpected path segments")
	}
	return model, rest, nil
}

// writeAIHintJSON emits a JSON response with sensible defaults.
func writeAIHintJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

// marshalHintPayload converts an aiHintPayload into the JSONB column values.
func marshalHintPayload(p *aiHintPayload) (capsJSON, toolsJSON, rulesJSON string) {
	caps, _ := json.Marshal(valueOrEmptyStringSlice(p.Capabilities))
	tools, _ := json.Marshal(valueOrEmptyToolMap(p.Tools))
	rules, _ := json.Marshal(valueOrEmptyStringMap(p.GlobalFormattingRules))
	return string(caps), string(tools), string(rules)
}

// fetchFullHintRow reads a single row (including soft-deleted) by model slug.
func fetchFullHintRow(model string) (*aiHintFullRow, error) {
	row := db.DB.QueryRow(`
		SELECT id, model, display_name, capabilities, system_prompt, tools,
		       global_formatting_rules, created_at, updated_at, deleted_at
		FROM ai_hints WHERE model = $1`, model)

	var (
		id                                  int64
		modelOut, displayName, systemPrompt string
		capsJSON, toolsJSON, rulesJSON      sql.NullString
		createdAt, updatedAt, deletedAt     sql.NullTime
	)
	if err := row.Scan(&id, &modelOut, &displayName, &capsJSON, &systemPrompt, &toolsJSON, &rulesJSON, &createdAt, &updatedAt, &deletedAt); err != nil {
		return nil, err
	}

	out := &aiHintFullRow{
		ID:                    id,
		Model:                 modelOut,
		DisplayName:           displayName,
		SystemPrompt:          systemPrompt,
		Capabilities:          []string{},
		Tools:                 map[string]*modeladapter.ToolHintJSON{},
		GlobalFormattingRules: map[string]string{},
	}
	if capsJSON.Valid && capsJSON.String != "" {
		_ = json.Unmarshal([]byte(capsJSON.String), &out.Capabilities)
	}
	if toolsJSON.Valid && toolsJSON.String != "" {
		_ = json.Unmarshal([]byte(toolsJSON.String), &out.Tools)
	}
	if rulesJSON.Valid && rulesJSON.String != "" {
		_ = json.Unmarshal([]byte(rulesJSON.String), &out.GlobalFormattingRules)
	}
	if createdAt.Valid {
		out.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		out.UpdatedAt = &updatedAt.Time
	}
	if deletedAt.Valid {
		out.DeletedAt = &deletedAt.Time
	}
	return out, nil
}

// snapshotHint writes the current full row for model into ai_hints_history.
// changeKind is a short tag like "update", "delete", "restore", "manual".
func snapshotHint(model string, changeKind string) error {
	current, err := fetchFullHintRow(model)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // nothing to snapshot (create path)
		}
		return err
	}
	snap, err := json.Marshal(current)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec(`
		INSERT INTO ai_hints_history (model, snapshot, change_kind)
		VALUES ($1, $2, $3)`, model, string(snap), changeKind)
	return err
}

// =============================================================================
// Handlers
// =============================================================================

// adminAIHintsListHandler handles GET /api/admin/ai-hints.
// Query params: include_deleted=true to include soft-deleted rows.
func adminAIHintsListHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}

	includeDeleted := r.URL.Query().Get("include_deleted") == "true"
	q := `SELECT model, display_name, capabilities, updated_at, deleted_at FROM ai_hints`
	if !includeDeleted {
		q += ` WHERE deleted_at IS NULL`
	}
	q += ` ORDER BY display_name ASC`

	rows, err := db.DB.Query(q)
	if err != nil {
		log.Printf("admin ai-hints list error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
		return
	}
	defer rows.Close()

	out := []aiHintListItem{}
	for rows.Next() {
		var (
			model, display             string
			capsJSON                   sql.NullString
			updatedAt, deletedAt       sql.NullTime
		)
		if err := rows.Scan(&model, &display, &capsJSON, &updatedAt, &deletedAt); err != nil {
			log.Printf("admin ai-hints list scan error: %v", err)
			continue
		}
		item := aiHintListItem{Model: model, DisplayName: display, Capabilities: []string{}}
		if capsJSON.Valid && capsJSON.String != "" {
			_ = json.Unmarshal([]byte(capsJSON.String), &item.Capabilities)
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			item.UpdatedAt = &t
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			item.DeletedAt = &t
		}
		out = append(out, item)
	}

	writeAIHintJSON(w, http.StatusOK, map[string]interface{}{"data": out})
}

// adminAIHintGetHandler handles GET /api/admin/ai-hints/{model}.
func adminAIHintGetHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	model, _, err := modelFromPath(r.URL.Path, "/api/admin/ai-hints/", 0)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	row, err := fetchFullHintRow(model)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if err != nil {
		log.Printf("admin ai-hint get error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
		return
	}
	writeAIHintJSON(w, http.StatusOK, row)
}

// adminAIHintCreateHandler handles POST /api/admin/ai-hints.
// Body: aiHintPayload. model slug is derived from display_name unless provided
// explicitly; the caller may supply a slug to break slugify ties.
func adminAIHintCreateHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}

	var p aiHintPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if strings.TrimSpace(p.DisplayName) == "" {
		writeError(w, http.StatusBadRequest, "display_name is required")
		return
	}
	slug := strings.TrimSpace(p.Model)
	if slug == "" {
		slug = slugifyModel(p.DisplayName)
	}
	if err := validateModelSlug(slug); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	capsJSON, toolsJSON, rulesJSON := marshalHintPayload(&p)

	_, err := db.DB.Exec(`
		INSERT INTO ai_hints (model, display_name, capabilities, system_prompt, tools, global_formatting_rules)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		slug, p.DisplayName, capsJSON, p.SystemPrompt, toolsJSON, rulesJSON,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			writeError(w, http.StatusConflict, "A bot with this model slug already exists")
			return
		}
		log.Printf("admin ai-hint create error: %v", err)
		writeError(w, http.StatusInternalServerError, "Create failed")
		return
	}

	writeAIHintJSON(w, http.StatusCreated, map[string]string{"status": "ok", "model": slug})
}

// adminAIHintUpdateHandler handles PUT /api/admin/ai-hints/{model}.
func adminAIHintUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	model, _, err := modelFromPath(r.URL.Path, "/api/admin/ai-hints/", 0)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var p aiHintPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if strings.TrimSpace(p.DisplayName) == "" {
		writeError(w, http.StatusBadRequest, "display_name is required")
		return
	}

	if err := snapshotHint(model, "update"); err != nil {
		log.Printf("admin ai-hint snapshot error: %v", err)
	}

	capsJSON, toolsJSON, rulesJSON := marshalHintPayload(&p)
	res, err := db.DB.Exec(`
		UPDATE ai_hints
		SET display_name = $2,
		    capabilities = $3,
		    system_prompt = $4,
		    tools = $5,
		    global_formatting_rules = $6,
		    updated_at = NOW()
		WHERE model = $1`,
		model, p.DisplayName, capsJSON, p.SystemPrompt, toolsJSON, rulesJSON,
	)
	if err != nil {
		log.Printf("admin ai-hint update error: %v", err)
		writeError(w, http.StatusInternalServerError, "Update failed")
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	writeAIHintJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// adminAIHintDeleteHandler handles DELETE /api/admin/ai-hints/{model}.
// Soft-deletes by setting deleted_at.
func adminAIHintDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	model, _, err := modelFromPath(r.URL.Path, "/api/admin/ai-hints/", 0)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := snapshotHint(model, "delete"); err != nil {
		log.Printf("admin ai-hint snapshot(delete) error: %v", err)
	}

	res, err := db.DB.Exec(`
		UPDATE ai_hints SET deleted_at = NOW(), updated_at = NOW()
		WHERE model = $1 AND deleted_at IS NULL`, model)
	if err != nil {
		log.Printf("admin ai-hint delete error: %v", err)
		writeError(w, http.StatusInternalServerError, "Delete failed")
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "Not found or already deleted")
		return
	}
	writeAIHintJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// adminAIHintRestoreHandler handles POST /api/admin/ai-hints/{model}/restore.
func adminAIHintRestoreHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	model, _, err := modelFromPath(r.URL.Path, "/api/admin/ai-hints/", 1)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := snapshotHint(model, "restore"); err != nil {
		log.Printf("admin ai-hint snapshot(restore) error: %v", err)
	}

	res, err := db.DB.Exec(`
		UPDATE ai_hints SET deleted_at = NULL, updated_at = NOW()
		WHERE model = $1 AND deleted_at IS NOT NULL`, model)
	if err != nil {
		log.Printf("admin ai-hint restore error: %v", err)
		writeError(w, http.StatusInternalServerError, "Restore failed")
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "Not found or not deleted")
		return
	}
	writeAIHintJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// adminAIHintSnapshotHandler handles POST /api/admin/ai-hints/{model}/snapshot.
// Manually records the current state into history.
func adminAIHintSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	model, _, err := modelFromPath(r.URL.Path, "/api/admin/ai-hints/", 1)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := fetchFullHintRow(model); err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "Not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Query failed")
		return
	}
	if err := snapshotHint(model, "manual"); err != nil {
		log.Printf("admin ai-hint manual snapshot error: %v", err)
		writeError(w, http.StatusInternalServerError, "Snapshot failed")
		return
	}
	writeAIHintJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// adminAIHintHistoryListHandler handles GET /api/admin/ai-hints/{model}/history.
func adminAIHintHistoryListHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	model, _, err := modelFromPath(r.URL.Path, "/api/admin/ai-hints/", 1)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rows, err := db.DB.Query(`
		SELECT id, model, change_kind, changed_at
		FROM ai_hints_history
		WHERE model = $1
		ORDER BY changed_at DESC
		LIMIT 200`, model)
	if err != nil {
		log.Printf("admin ai-hint history list error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
		return
	}
	defer rows.Close()

	out := []aiHintHistoryItem{}
	for rows.Next() {
		var (
			id         int64
			m, kind    string
			changedAt  sql.NullTime
		)
		if err := rows.Scan(&id, &m, &kind, &changedAt); err != nil {
			continue
		}
		item := aiHintHistoryItem{ID: id, Model: m, ChangeKind: kind}
		if changedAt.Valid {
			t := changedAt.Time
			item.ChangedAt = &t
		}
		out = append(out, item)
	}
	writeAIHintJSON(w, http.StatusOK, map[string]interface{}{"data": out})
}

// adminAIHintHistoryRestoreHandler handles
// POST /api/admin/ai-hints/{model}/history/{id}/restore.
func adminAIHintHistoryRestoreHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	// Path looks like /api/admin/ai-hints/{model}/history/{id}/restore
	stripped := strings.TrimPrefix(r.URL.Path, "/api/admin/ai-hints/")
	stripped = strings.Trim(stripped, "/")
	parts := strings.Split(stripped, "/")
	if len(parts) != 4 || parts[1] != "history" || parts[3] != "restore" {
		writeError(w, http.StatusBadRequest, "Invalid history restore path")
		return
	}
	model := parts[0]
	if err := validateModelSlug(model); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	historyID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid history id")
		return
	}

	var snapJSON string
	err = db.DB.QueryRow(`
		SELECT snapshot FROM ai_hints_history WHERE id = $1 AND model = $2`,
		historyID, model,
	).Scan(&snapJSON)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "History entry not found")
		return
	}
	if err != nil {
		log.Printf("admin ai-hint history restore fetch error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
		return
	}

	var snap aiHintFullRow
	if err := json.Unmarshal([]byte(snapJSON), &snap); err != nil {
		writeError(w, http.StatusInternalServerError, "Corrupt snapshot")
		return
	}

	// Snapshot the pre-restore state first so the current version is itself
	// recoverable from history.
	if err := snapshotHint(model, "restore"); err != nil {
		log.Printf("admin ai-hint snapshot(restore-pre) error: %v", err)
	}

	caps, _ := json.Marshal(valueOrEmptyStringSlice(snap.Capabilities))
	tools, _ := json.Marshal(valueOrEmptyToolMap(snap.Tools))
	rules, _ := json.Marshal(valueOrEmptyStringMap(snap.GlobalFormattingRules))

	res, err := db.DB.Exec(`
		UPDATE ai_hints
		SET display_name = $2,
		    capabilities = $3,
		    system_prompt = $4,
		    tools = $5,
		    global_formatting_rules = $6,
		    deleted_at = NULL,
		    updated_at = NOW()
		WHERE model = $1`,
		model, snap.DisplayName, string(caps), snap.SystemPrompt, string(tools), string(rules),
	)
	if err != nil {
		log.Printf("admin ai-hint history restore update error: %v", err)
		writeError(w, http.StatusInternalServerError, "Restore failed")
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		// Row was hard-deleted externally; re-insert from snapshot.
		_, err = db.DB.Exec(`
			INSERT INTO ai_hints (model, display_name, capabilities, system_prompt, tools, global_formatting_rules)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			model, snap.DisplayName, string(caps), snap.SystemPrompt, string(tools), string(rules),
		)
		if err != nil {
			log.Printf("admin ai-hint history restore insert error: %v", err)
			writeError(w, http.StatusInternalServerError, "Restore failed")
			return
		}
	}

	writeAIHintJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// adminAIHintReloadHandler handles POST /api/admin/ai-hints/reload.
// Rebuilds the in-memory HintsLoader from DB state.
func adminAIHintReloadHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	if mcpHintsLoader == nil {
		writeError(w, http.StatusServiceUnavailable, "Hints loader not initialized")
		return
	}
	count, err := loadAIHintsFromDB(mcpHintsLoader)
	if err != nil {
		log.Printf("admin ai-hint reload error: %v", err)
		writeError(w, http.StatusInternalServerError, "Reload failed")
		return
	}
	writeAIHintJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "loaded": count})
}

// adminAIHintExportHandler handles GET /api/admin/ai-hints/{model}/export.
// Returns the hint formatted exactly like the on-disk JSON files so the
// response can be committed back to the repo if desired.
func adminAIHintExportHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	model, _, err := modelFromPath(r.URL.Path, "/api/admin/ai-hints/", 1)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	row, err := fetchFullHintRow(model)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "Not found")
		return
	}
	if err != nil {
		log.Printf("admin ai-hint export error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
		return
	}

	exportShape := modeladapter.ModelHint{
		Model:                 row.Model,
		DisplayName:           row.DisplayName,
		Capabilities:          row.Capabilities,
		SystemPrompt:          row.SystemPrompt,
		Tools:                 row.Tools,
		GlobalFormattingRules: row.GlobalFormattingRules,
	}
	body, err := json.MarshalIndent(exportShape, "", "  ")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Encode failed")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.json"`, model))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// adminAIHintImportHandler handles POST /api/admin/ai-hints/import.
// Body is a single JSON file contents matching the on-disk hint format.
// If the model exists, it is updated; otherwise a new row is created.
func adminAIHintImportHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	data, err := io.ReadAll(io.LimitReader(r.Body, 2*1024*1024)) // 2MB hard cap
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to read body")
		return
	}
	var p aiHintPayload
	if err := json.Unmarshal(data, &p); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	slug := strings.TrimSpace(p.Model)
	if slug == "" {
		slug = slugifyModel(p.DisplayName)
	}
	if err := validateModelSlug(slug); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(p.DisplayName) == "" {
		writeError(w, http.StatusBadRequest, "display_name is required")
		return
	}

	// If the row exists, snapshot before overwriting.
	if _, err := fetchFullHintRow(slug); err == nil {
		if err := snapshotHint(slug, "import"); err != nil {
			log.Printf("admin ai-hint import snapshot error: %v", err)
		}
	}

	capsJSON, toolsJSON, rulesJSON := marshalHintPayload(&p)
	_, err = db.DB.Exec(`
		INSERT INTO ai_hints (model, display_name, capabilities, system_prompt, tools, global_formatting_rules)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (model) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			capabilities = EXCLUDED.capabilities,
			system_prompt = EXCLUDED.system_prompt,
			tools = EXCLUDED.tools,
			global_formatting_rules = EXCLUDED.global_formatting_rules,
			deleted_at = NULL,
			updated_at = NOW()`,
		slug, p.DisplayName, capsJSON, p.SystemPrompt, toolsJSON, rulesJSON,
	)
	if err != nil {
		log.Printf("admin ai-hint import error: %v", err)
		writeError(w, http.StatusInternalServerError, "Import failed")
		return
	}
	writeAIHintJSON(w, http.StatusOK, map[string]string{"status": "ok", "model": slug})
}
