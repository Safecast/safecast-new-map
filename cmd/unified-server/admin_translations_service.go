package main

import (
	"fmt"
	"strings"
)

type translationRecord struct {
	ID           int64       `json:"id"`
	LanguageCode string      `json:"language_code"`
	Key          string      `json:"key"`
	Value        string      `json:"value"`
	UpdatedAt    interface{} `json:"updated_at"`
}

type translationListResult struct {
	Rows      []translationRecord `json:"data"`
	Total     int                 `json:"total"`
	Languages []string            `json:"languages"`
}

type translationService struct{}

func newTranslationService() *translationService {
	return &translationService{}
}

func (s *translationService) List(limit, offset int, sortCol, order, search, langFilter string) (*translationListResult, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if langFilter != "" {
		conditions = append(conditions, fmt.Sprintf("language_code = $%d", argIdx))
		args = append(args, langFilter)
		argIdx++
	}
	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(key ILIKE $%d OR value ILIKE $%d)", argIdx, argIdx+1))
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIdx += 2
	}

	whereSQL := ""
	if len(conditions) > 0 {
		whereSQL = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM translations %s", whereSQL)
	var total int
	if err := db.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	dataQuery := fmt.Sprintf("SELECT id, language_code, key, value, updated_at FROM translations %s ORDER BY %s %s LIMIT %d OFFSET %d",
		whereSQL, sortCol, order, limit, offset)
	rows, err := db.DB.Query(dataQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]translationRecord, 0, limit)
	for rows.Next() {
		var row translationRecord
		if err := rows.Scan(&row.ID, &row.LanguageCode, &row.Key, &row.Value, &row.UpdatedAt); err != nil {
			continue
		}
		out = append(out, row)
	}

	langRows, err := db.DB.Query("SELECT DISTINCT language_code FROM translations ORDER BY language_code")
	languages := make([]string, 0)
	if err == nil {
		defer langRows.Close()
		for langRows.Next() {
			var lang string
			if langRows.Scan(&lang) == nil {
				languages = append(languages, lang)
			}
		}
	}

	return &translationListResult{
		Rows:      out,
		Total:     total,
		Languages: languages,
	}, nil
}

func (s *translationService) Update(id int64, value string) error {
	_, err := db.DB.Exec("UPDATE translations SET value = $1, updated_at = NOW() WHERE id = $2", value, id)
	return err
}

func (s *translationService) Create(languageCode, key, value string) error {
	_, err := db.DB.Exec(
		"INSERT INTO translations (language_code, key, value) VALUES ($1, $2, $3) ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()",
		languageCode, key, value,
	)
	return err
}

func (s *translationService) Delete(id int64) error {
	_, err := db.DB.Exec("DELETE FROM translations WHERE id = $1", id)
	return err
}
