package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// tourStepColumns lists the sortable columns for the admin list endpoint.
var tourStepColumns = []string{
	"id", "step_key", "sort_order", "selector", "enabled", "updated_at",
}

// tourStep is the full row shape used by the admin API.
type tourStep struct {
	ID             int64          `json:"id"`
	StepKey        string         `json:"step_key"`
	SortOrder      int            `json:"sort_order"`
	Selector       string         `json:"selector"`
	Center         bool           `json:"center"`
	Enabled        bool           `json:"enabled"`
	RequireLogin   sql.NullBool   `json:"-"`
	RequireAdmin   sql.NullBool   `json:"-"`
	ShowIfFeature  sql.NullString `json:"-"`
	Viewport       string         `json:"viewport"`
	FirstTimeOnly  bool           `json:"first_time_only"`
	UpdatedAt      interface{}    `json:"updated_at"`
}

// MarshalJSON emits nullable columns as JSON null (not zero values) so the
// admin UI can distinguish "don't care" from "must be false".
func (t tourStep) toMap() map[string]interface{} {
	m := map[string]interface{}{
		"id":              t.ID,
		"step_key":        t.StepKey,
		"sort_order":      t.SortOrder,
		"selector":        t.Selector,
		"center":          t.Center,
		"enabled":         t.Enabled,
		"viewport":        t.Viewport,
		"first_time_only": t.FirstTimeOnly,
		"updated_at":      t.UpdatedAt,
	}
	if t.RequireLogin.Valid {
		m["require_login"] = t.RequireLogin.Bool
	} else {
		m["require_login"] = nil
	}
	if t.RequireAdmin.Valid {
		m["require_admin"] = t.RequireAdmin.Bool
	} else {
		m["require_admin"] = nil
	}
	if t.ShowIfFeature.Valid {
		m["show_if_feature"] = t.ShowIfFeature.String
	} else {
		m["show_if_feature"] = nil
	}
	return m
}

// tourStepInput is the writable subset used by Create/Update.
type tourStepInput struct {
	StepKey        string  `json:"step_key"`
	Selector       string  `json:"selector"`
	Center         bool    `json:"center"`
	Enabled        bool    `json:"enabled"`
	RequireLogin   *bool   `json:"require_login"`
	RequireAdmin   *bool   `json:"require_admin"`
	ShowIfFeature  *string `json:"show_if_feature"`
	Viewport       string  `json:"viewport"`
	FirstTimeOnly  bool    `json:"first_time_only"`
}

// tourPublicStep is the shape returned by the public /api/tour/steps endpoint.
// Text is resolved for the requested language; conditions are a nested object
// so the client can filter cheaply.
type tourPublicStep struct {
	StepKey    string                 `json:"step_key"`
	Selector   string                 `json:"selector"`
	Center     bool                   `json:"center"`
	Text       string                 `json:"text"`
	Conditions map[string]interface{} `json:"conditions"`
}

type tourService struct{}

func newTourService() *tourService { return &tourService{} }

// List returns all tour steps ordered by sort_order (admin use).
func (s *tourService) List(sortCol, order string) ([]map[string]interface{}, error) {
	if !isValidColumn(sortCol, tourStepColumns) {
		sortCol = "sort_order"
	}
	order = strings.ToUpper(order)
	if order != "DESC" {
		order = "ASC"
	}
	q := fmt.Sprintf(`
		SELECT id, step_key, sort_order, selector, center, enabled,
		       require_login, require_admin, show_if_feature,
		       viewport, first_time_only, updated_at
		  FROM tour_steps ORDER BY %s %s`, sortCol, order)
	rows, err := db.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]map[string]interface{}, 0, 16)
	for rows.Next() {
		var t tourStep
		if err := rows.Scan(&t.ID, &t.StepKey, &t.SortOrder, &t.Selector,
			&t.Center, &t.Enabled, &t.RequireLogin, &t.RequireAdmin,
			&t.ShowIfFeature, &t.Viewport, &t.FirstTimeOnly, &t.UpdatedAt); err != nil {
			log.Printf("tourService.List scan: %v", err)
			continue
		}
		out = append(out, t.toMap())
	}
	return out, nil
}

// ListPublic returns enabled steps with text resolved for the requested
// language (falling back to English). Used by the map's tour loader.
func (s *tourService) ListPublic(lang string) ([]tourPublicStep, error) {
	if lang == "" {
		lang = "en"
	}
	rows, err := db.DB.Query(`
		SELECT step_key, selector, center,
		       require_login, require_admin, show_if_feature,
		       viewport, first_time_only
		  FROM tour_steps
		 WHERE enabled = TRUE
		 ORDER BY sort_order ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	translationsMu.RLock()
	tr := translations
	translationsMu.RUnlock()

	resolve := func(key string) string {
		if m, ok := tr[lang]; ok {
			if v, ok := m[key]; ok && v != "" {
				return v
			}
		}
		if m, ok := tr["en"]; ok {
			if v, ok := m[key]; ok {
				return v
			}
		}
		return ""
	}

	out := make([]tourPublicStep, 0, 16)
	for rows.Next() {
		var (
			stepKey, selector, viewport string
			center, firstTimeOnly       bool
			reqLogin, reqAdmin          sql.NullBool
			showIf                      sql.NullString
		)
		if err := rows.Scan(&stepKey, &selector, &center,
			&reqLogin, &reqAdmin, &showIf, &viewport, &firstTimeOnly); err != nil {
			log.Printf("tourService.ListPublic scan: %v", err)
			continue
		}

		conds := map[string]interface{}{
			"viewport":        viewport,
			"first_time_only": firstTimeOnly,
		}
		if reqLogin.Valid {
			conds["require_login"] = reqLogin.Bool
		} else {
			conds["require_login"] = nil
		}
		if reqAdmin.Valid {
			conds["require_admin"] = reqAdmin.Bool
		} else {
			conds["require_admin"] = nil
		}
		if showIf.Valid {
			conds["show_if_feature"] = showIf.String
		} else {
			conds["show_if_feature"] = nil
		}

		out = append(out, tourPublicStep{
			StepKey:    stepKey,
			Selector:   selector,
			Center:     center,
			Text:       resolve("tour." + stepKey + ".text"),
			Conditions: conds,
		})
	}
	return out, nil
}

// Create inserts a new step at the end of the current order.
func (s *tourService) Create(in tourStepInput) (int64, error) {
	if in.StepKey == "" || in.Selector == "" {
		return 0, fmt.Errorf("step_key and selector are required")
	}
	if in.Viewport == "" {
		in.Viewport = "any"
	}
	var nextOrder int
	if err := db.DB.QueryRow(
		`SELECT COALESCE(MAX(sort_order), 0) + 1 FROM tour_steps`,
	).Scan(&nextOrder); err != nil {
		return 0, err
	}

	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO tour_steps (
			step_key, sort_order, selector, center, enabled,
			require_login, require_admin, show_if_feature,
			viewport, first_time_only
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id`,
		in.StepKey, nextOrder, in.Selector, in.Center, in.Enabled,
		nullBool(in.RequireLogin), nullBool(in.RequireAdmin),
		nullString(in.ShowIfFeature), in.Viewport, in.FirstTimeOnly,
	).Scan(&id)
	return id, err
}

// Update rewrites all mutable columns on a row. Full replace is simpler than
// building a dynamic UPDATE and matches how the admin UI sends the record
// back after editing.
func (s *tourService) Update(id int64, in tourStepInput) error {
	if in.Viewport == "" {
		in.Viewport = "any"
	}
	_, err := db.DB.Exec(`
		UPDATE tour_steps
		   SET step_key = $1, selector = $2, center = $3, enabled = $4,
		       require_login = $5, require_admin = $6, show_if_feature = $7,
		       viewport = $8, first_time_only = $9, updated_at = NOW()
		 WHERE id = $10`,
		in.StepKey, in.Selector, in.Center, in.Enabled,
		nullBool(in.RequireLogin), nullBool(in.RequireAdmin),
		nullString(in.ShowIfFeature), in.Viewport, in.FirstTimeOnly, id,
	)
	return err
}

// Delete drops the step row and its tour.<step_key>.* translation rows so the
// translations table does not accumulate orphans.
func (s *tourService) Delete(id int64) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var stepKey string
	if err := tx.QueryRow(`SELECT step_key FROM tour_steps WHERE id = $1`, id).Scan(&stepKey); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM tour_steps WHERE id = $1`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(
		`DELETE FROM translations WHERE key LIKE 'tour.' || $1 || '.%'`, stepKey,
	); err != nil {
		return err
	}
	return tx.Commit()
}

// Reorder rewrites sort_order to match the given slice of IDs (index 0 ⇒
// order 1). Uses a single transaction and a temporary offset to avoid
// violating any potential future unique-ordering constraint.
func (s *tourService) Reorder(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Phase 1: move every target into a high, non-colliding range so the
	// final assignment can never conflict with a row not yet updated.
	for i, id := range ids {
		if _, err := tx.Exec(
			`UPDATE tour_steps SET sort_order = $1 WHERE id = $2`,
			10000+i, id,
		); err != nil {
			return err
		}
	}
	// Phase 2: assign the real, compact order.
	for i, id := range ids {
		if _, err := tx.Exec(
			`UPDATE tour_steps SET sort_order = $1, updated_at = NOW() WHERE id = $2`,
			i+1, id,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func nullBool(p *bool) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

func nullString(p *string) interface{} {
	if p == nil || *p == "" {
		return nil
	}
	return *p
}
