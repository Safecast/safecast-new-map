//go:build duckdb

package main

import (
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
)

func TestQAColumnFilters(t *testing.T) {
	d, err := sql.Open("duckdb", "")
	if err != nil {
		t.Fatal(err)
	}
	d.SetMaxOpenConns(1)
	old := duckDB
	duckDB = d
	defer func() { duckDB = old }()
	for _, q := range []string{
		`CREATE TABLE qa_embeddings (id BIGINT, chat_id BIGINT, question VARCHAR, answer VARCHAR, embedding VARCHAR, feedback_score INTEGER DEFAULT 0, used_count INTEGER DEFAULT 0, last_used_at TIMESTAMPTZ, status VARCHAR DEFAULT 'active', lang VARCHAR DEFAULT '', created_at TIMESTAMPTZ DEFAULT now(), corrected BOOLEAN DEFAULT false)`,
		`INSERT INTO qa_embeddings (id, chat_id, question, answer, status) VALUES (1775312585310528000, 1, 'Tell me about track ABC', 'Dose is low', 'active'), (2, 2, 'Where is Fukushima', 'North Japan', 'pending')`,
	} {
		if _, err := d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	for url, want := range map[string]int{
		"/x": 2, "/x?f_question=track": 1, "/x?f_answer=japan": 1, "/x?f_status=pend": 1,
		"/x?f_id=17753": 1, "/x?status=!archived": 2, "/x?f_question=track&f_status=pend": 0, "/x?sort=answer&order=asc": 2,
	} {
		w := httptest.NewRecorder()
		adminQAEmbeddingsListHandler(w, httptest.NewRequest("GET", url, nil))
		var out struct {
			Total int `json:"total"`
			Data  []map[string]interface{}
		}
		if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &out) != nil || out.Total != want {
			t.Fatalf("%s: code=%d total=%d want=%d body=%s", url, w.Code, out.Total, want, w.Body.String())
		}
	}
}
