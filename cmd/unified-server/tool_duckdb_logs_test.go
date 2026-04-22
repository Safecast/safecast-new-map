package main

import "testing"

func TestValidateReadOnlyQuery(t *testing.T) {
	cases := []struct {
		name        string
		in          string
		wantErr     bool
		wantCleaned string
	}{
		{"plain select", "SELECT 1", false, "SELECT 1"},
		{"lowercase select", "select 1", false, "select 1"},
		{"trailing semicolon", "SELECT 1;", false, "SELECT 1"},
		{"multiple trailing semicolons and whitespace", "SELECT 1 ;; \n", false, "SELECT 1"},
		{"leading/trailing whitespace", "  SELECT now()  ", false, "SELECT now()"},
		{"with cte", "WITH t AS (SELECT 1) SELECT * FROM t", false, "WITH t AS (SELECT 1) SELECT * FROM t"},
		{"lowercase with cte", "with t as (select 1) select * from t;", false, "with t as (select 1) select * from t"},

		{"empty", "", true, ""},
		{"only semicolons", ";;;", true, ""},
		{"only whitespace", "   \n\t", true, ""},
		{"non-select verb", "DELETE FROM mcp_query_log", true, ""},
		{"drop table", "DROP TABLE foo", true, ""},
		{"update", "UPDATE foo SET x=1", true, ""},
		{"insert", "INSERT INTO foo VALUES (1)", true, ""},

		{"stacked select then delete", "SELECT 1; DELETE FROM foo", true, ""},
		{"stacked with stray whitespace", "SELECT 1 ; DELETE FROM foo ;", true, ""},
		{"inner semicolon in comment", "SELECT 1 /* ; */", true, ""},
		{"inner semicolon in string literal", "SELECT ';'", true, ""}, // false-positive reject is safer than allowing stacked DDL

		{"attach statement", "ATTACH 'evil.db' AS x", true, ""},
		{"pragma", "PRAGMA disable_verification", true, ""},
		{"copy", "COPY foo FROM 'file'", true, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cleaned, errMsg := validateReadOnlyQuery(tc.in)
			if tc.wantErr {
				if errMsg == "" {
					t.Fatalf("expected error for %q, got cleaned=%q", tc.in, cleaned)
				}
				return
			}
			if errMsg != "" {
				t.Fatalf("unexpected error for %q: %s", tc.in, errMsg)
			}
			if cleaned != tc.wantCleaned {
				t.Fatalf("cleaned mismatch for %q: got %q want %q", tc.in, cleaned, tc.wantCleaned)
			}
		})
	}
}
