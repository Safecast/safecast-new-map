package main

import (
	"database/sql"
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	modeladapter "safecast-new-map/cmd/unified-server/model-adapter"
)

// aiHintRow mirrors the ai_hints table row. JSONB columns are decoded into
// their concrete Go types on read.
type aiHintRow struct {
	Model                 string
	DisplayName           string
	Capabilities          []string
	SystemPrompt          string
	Tools                 map[string]*modeladapter.ToolHintJSON
	GlobalFormattingRules map[string]string
}

// toModelHint converts a DB row into the in-memory ModelHint shape the MCP
// server consumes.
func (r *aiHintRow) toModelHint() *modeladapter.ModelHint {
	return &modeladapter.ModelHint{
		Model:                 r.Model,
		DisplayName:           r.DisplayName,
		Capabilities:          r.Capabilities,
		SystemPrompt:          r.SystemPrompt,
		Tools:                 r.Tools,
		GlobalFormattingRules: r.GlobalFormattingRules,
	}
}

// seedAIHintsDB inserts the on-disk hint JSON files into ai_hints for any
// models that are not already present. Delegates to seedAIHintsDBFromFS.
func seedAIHintsDB(hintsDir string) {
	seedAIHintsDBFromFS(os.DirFS(hintsDir), hintsDir)
}

// seedAIHintsDBFromFS is the fs.FS-aware seeder. Accepts either a real-disk
// fs (os.DirFS) or an embedded tree (fs.Sub(embeddedHintsFS, "hints")).
// `label` is used only for log output. Matches the translations seeding
// pattern — never overwrites admin edits.
func seedAIHintsDBFromFS(fsys fs.FS, label string) {
	if db == nil || db.DB == nil {
		return
	}
	if fsys == nil {
		return
	}

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		log.Printf("[ai_hints] Cannot read hints source %s for seeding: %v", label, err)
		return
	}

	inserted := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := fs.ReadFile(fsys, entry.Name())
		if err != nil {
			log.Printf("[ai_hints] Cannot read %s/%s: %v", label, entry.Name(), err)
			continue
		}

		var hint modeladapter.ModelHint
		if err := json.Unmarshal(data, &hint); err != nil {
			log.Printf("[ai_hints] Cannot parse %s/%s: %v", label, entry.Name(), err)
			continue
		}

		model := hint.Model
		if entry.Name() == "default.json" {
			model = "unknown"
		}
		if model == "" {
			continue
		}

		capsJSON, _ := json.Marshal(valueOrEmptyStringSlice(hint.Capabilities))
		toolsJSON, _ := json.Marshal(valueOrEmptyToolMap(hint.Tools))
		rulesJSON, _ := json.Marshal(valueOrEmptyStringMap(hint.GlobalFormattingRules))

		_, err = db.DB.Exec(`
			INSERT INTO ai_hints (model, display_name, capabilities, system_prompt, tools, global_formatting_rules)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (model) DO NOTHING`,
			model, hint.DisplayName, string(capsJSON), hint.SystemPrompt, string(toolsJSON), string(rulesJSON),
		)
		if err != nil {
			log.Printf("[ai_hints] Error seeding %s: %v", model, err)
			continue
		}
		inserted++
	}
	if inserted > 0 {
		log.Printf("[ai_hints] Seeded up to %d hint rows from %s (skipping existing)", inserted, label)
	}
}

// loadAIHintsFromDB reads all non-deleted ai_hints rows and atomically swaps
// them into the shared HintsLoader. Returns the number loaded.
func loadAIHintsFromDB(loader *modeladapter.HintsLoader) (int, error) {
	if db == nil || db.DB == nil || loader == nil {
		return 0, nil
	}

	rows, err := db.DB.Query(`
		SELECT model, display_name, capabilities, system_prompt, tools, global_formatting_rules
		FROM ai_hints
		WHERE deleted_at IS NULL`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	hints := make(map[string]*modeladapter.ModelHint)
	for rows.Next() {
		row, err := scanAIHintRow(rows)
		if err != nil {
			log.Printf("[ai_hints] Scan error: %v", err)
			continue
		}
		hints[row.Model] = row.toModelHint()
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	loader.SetHints(hints)
	return len(hints), nil
}

// scanAIHintRow decodes a single row from a SELECT that returns the six
// canonical ai_hints columns (model, display_name, capabilities, system_prompt,
// tools, global_formatting_rules).
func scanAIHintRow(scanner interface {
	Scan(dest ...interface{}) error
}) (*aiHintRow, error) {
	var (
		model, displayName, systemPrompt string
		capsJSON, toolsJSON, rulesJSON   sql.NullString
	)
	if err := scanner.Scan(&model, &displayName, &capsJSON, &systemPrompt, &toolsJSON, &rulesJSON); err != nil {
		return nil, err
	}

	row := &aiHintRow{
		Model:        model,
		DisplayName:  displayName,
		SystemPrompt: systemPrompt,
		Capabilities: []string{},
		Tools:        map[string]*modeladapter.ToolHintJSON{},
		GlobalFormattingRules: map[string]string{},
	}
	if capsJSON.Valid && capsJSON.String != "" {
		_ = json.Unmarshal([]byte(capsJSON.String), &row.Capabilities)
	}
	if toolsJSON.Valid && toolsJSON.String != "" {
		_ = json.Unmarshal([]byte(toolsJSON.String), &row.Tools)
	}
	if rulesJSON.Valid && rulesJSON.String != "" {
		_ = json.Unmarshal([]byte(rulesJSON.String), &row.GlobalFormattingRules)
	}
	return row, nil
}

// slugifyModel turns a display name into a safe MCP routing slug. "Claude 3.5
// Sonnet" becomes "claude-3-5-sonnet". Only used when creating a new bot
// group; the slug is readonly after that point.
func slugifyModel(displayName string) string {
	s := strings.ToLower(strings.TrimSpace(displayName))
	var b strings.Builder
	lastDash := true
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				b.WriteRune('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func valueOrEmptyStringSlice(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func valueOrEmptyToolMap(m map[string]*modeladapter.ToolHintJSON) map[string]*modeladapter.ToolHintJSON {
	if m == nil {
		return map[string]*modeladapter.ToolHintJSON{}
	}
	return m
}

func valueOrEmptyStringMap(m map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return m
}
