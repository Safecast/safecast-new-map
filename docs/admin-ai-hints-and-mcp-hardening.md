tunu# Admin AI Hints Editor + MCP Read-Only Hardening

Date: 2026-04-22
PRs: [#121](https://github.com/Safecast/safecast-new-map/pull/121), [#122](https://github.com/Safecast/safecast-new-map/pull/122), [#123](https://github.com/Safecast/safecast-new-map/pull/123)

Three related pieces of work landed together:

1. A new admin page for editing the per-model MCP hints (previously only editable by hand-editing on-disk JSON files).
2. Hardening the one MCP tool that accepted free-form SQL from the AI client.
3. Parameterizing the one MCP tool that built SQL via `fmt.Sprintf` + quote-escaping, and adding missing `read-only` annotations so all 19 MCP tools advertise themselves correctly.

## 1. Admin AI Hints editor — PR #121

### What

A new page at **`/admin/ai-hints`** for adding, editing, soft-deleting, and restoring per-model hints (Claude, Kimi, Qwen, GPT, …) that the MCP server uses to customize its tool descriptions, examples, warnings, and system prompts per AI client.

The page follows the same shape as the Translations admin (`/admin/translations`): a structured editor backed by PostgreSQL, with an **Export JSON** button so edits can be round-tripped back into the repo if desired.

### Why

Before this change, the hints were editable only by:

1. Editing `cmd/unified-server/hints/<model>.json`.
2. Rebuilding the binary.
3. Redeploying.

That is:

- Not reachable for non-engineers.
- Slow feedback loop — a typo in a hint meant a full rebuild.
- High-friction to experiment with prompt changes across bots.

Making hints editable from the admin panel (with history) makes it safe to iterate on prompts for each bot without redeploying, and keeps a durable audit trail.

### How

**Schema** ([`migrations/add_ai_hints_tables.sql`](../migrations/add_ai_hints_tables.sql))

- `ai_hints` — one row per bot slug (`claude`, `kimi`, `qwen`, `gpt`, `unknown`). JSONB columns for `capabilities`, `tools`, `global_formatting_rules`. `deleted_at` for soft delete so "delete" is recoverable.
- `ai_hints_history` — append-only snapshot table. Every mutation (`create`, `update`, `delete`, `restore`, `import`) writes the pre-change state as a JSONB `snapshot`. This powers the **History** tab's one-click "Restore this version".

**Startup flow** ([`cmd/unified-server/mcp_register.go`](../cmd/unified-server/mcp_register.go))

1. The server opens a filesystem for the hint templates. If `MCP_HINTS_DIR` is set, it points to a real directory (useful for devs editing JSON without rebuilds). Otherwise the server uses `//go:embed hints/*.json` — the templates are baked into the binary. This matters on production, where `/usr/local/bin/safecast-new-map` has no neighbouring `hints/` directory.
2. `seedAIHintsDBFromFS` copies any missing rows from the template FS into the `ai_hints` table with `ON CONFLICT (model) DO NOTHING`. Existing DB rows are never overwritten — admin edits survive restarts.
3. `loadAIHintsFromDB` reads every non-deleted row and hands the map to `HintsLoader.SetHints`. The DB is now the runtime source of truth; the JSON templates in the repo are just bootstrap seeds.

The `HintsLoader` was refactored to accept an `fs.FS` instead of a raw directory path so the same code path serves both real-disk and embedded inputs.

**UI** ([`cmd/unified-server/public_html/admin-ai-hints.html`](../cmd/unified-server/public_html/admin-ai-hints.html))

Two-column layout:

- **Sidebar**: bot list (with a "Show deleted" toggle) + **+ Add bot** button. New bots auto-derive their slug from the display name (`"Claude 3.5 Sonnet"` → `claude-3-5-sonnet`); the slug is read-only once saved.
- **Editor**: sub-tabs for General / Global rules / Tools / History.
  - **General**: display name, slug, capabilities (hybrid dropdown + free-text), system prompt.
  - **Global rules**: KV editor.
  - **Tools**: per-tool accordion with description, output hints, warnings list, formatting rules KV, and an example editor that is structured (user query + tool name + arguments KV) rather than a raw JSON textarea.
  - **History**: list of prior snapshots with a one-click "Restore this version" button. Restore itself creates a `change_kind='restore'` snapshot before overwriting, so the action is reversible.

**Action bar**: Save, Reload into Memory, Snapshot, Export JSON, Import JSON, Delete / Restore. "Reload into Memory" calls the atomic `SetHints` swap so MCP picks up changes without a restart.

**Auth**: same as Translations — admin session OR `?password=` query param, via `checkAdminAccess`.

### Backward compatibility

- The JSON templates in `cmd/unified-server/hints/` remain in the repo. They are now seed templates only; they are not re-read at runtime once the DB has a row for a given model. `hints/README.md` was updated to spell this out.
- `MCP_HINTS_DIR` still overrides the embedded templates with a real disk path, for developers who want to iterate on JSON without rebuilds.

---

## 2. Tighten SELECT-only guard on `query_duckdb_logs` — PR #122

### What

The [`query_duckdb_logs`](../cmd/unified-server/tool_duckdb_logs.go) tool accepts free-form SQL from the AI client for introspecting MCP analytics tables (`mcp_query_log`, `mcp_ai_query_log`, `chat_questions`). The previous SELECT-only check was too weak; this change replaces it with a stricter validator.

### Why

The old check was:

```go
if !strings.HasPrefix(strings.ToUpper(query), "SELECT") {
    return "Only SELECT queries are allowed", nil
}
```

That is trivially bypassable by stacking statements:

```sql
SELECT 1; DELETE FROM mcp_query_log;
```

`HasPrefix("SELECT")` returns true, and DuckDB's Go driver can execute multi-statement payloads via `Query()`. That gave the AI chatbot a potential write path even though the tool is documented as read-only — directly contradicting the project invariant that the AI must not modify data.

### How

New helper [`validateReadOnlyQuery()`](../cmd/unified-server/tool_duckdb_logs.go):

1. Strip any number of trailing `;` and whitespace (`SELECT 1;;` and `SELECT 1 ;\n` both OK).
2. Reject empty inputs.
3. Require the cleaned query to start with `SELECT` or `WITH` (CTEs are a common read-only shape).
4. Reject if **any** `;` remains in the cleaned text. This catches stacked statements and the `SELECT 1 /* ; */; DROP …` comment-embedding trick.

Trade-off: a query that legitimately contains `;` inside a string literal (e.g. `SELECT ';'`) is rejected. For analytics-log queries this is acceptable, and the prompt for this tool hints that the AI should avoid semicolons.

Also swapped `Query()` → `QueryContext()` so request cancellations propagate, and added `WithReadOnlyHintAnnotation(true)` so MCP clients see the correct capability hint.

21-case unit test ([`tool_duckdb_logs_test.go`](../cmd/unified-server/tool_duckdb_logs_test.go)) covers the accept/reject matrix.

---

## 3. Parameterize `query_extreme_readings` + annotate remaining tools — PR #123

### What

[`query_extreme_readings`](../cmd/unified-server/tool_extreme_readings.go) built its `WHERE` clause with `fmt.Sprintf` and a naive `'` → `''` escape for user-supplied `device_id` values. Refactored to use `?` placeholders for every user input. Also added `WithReadOnlyHintAnnotation(true)` to the 3 tools that were missing it.

### Why

The old pattern:

```go
deviceList[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(dev, "'", "''"))
whereConditions = append(whereConditions, fmt.Sprintf(
    "device_id NOT IN (%s)", strings.Join(deviceList, ", "),
))
```

The `'` → `''` escape is technically correct for standard SQL string literals, but:

- It is easy to get wrong under driver quirks (e.g. `E'\\''` escape-string extensions, null bytes, character-set tricks).
- `fmt.Sprintf` into a SQL string is the exact shape of every SQL-injection incident.
- Every other tool in the codebase already uses parameterized queries — this one was inconsistent.

Using the driver's `?` placeholders removes the escape-correctness question entirely: user strings never land in the SQL text.

### How

```go
whereConditions := []string{"doserate > 0 AND doserate < 10000"}
var args []any

if hasGeoFilter {
    whereConditions = append(whereConditions, "lat BETWEEN ? AND ? AND lon BETWEEN ? AND ?")
    args = append(args, minLat, maxLat, minLon, maxLon)
}

if len(excludeDevices) > 0 {
    placeholders := make([]string, len(excludeDevices))
    for i, dev := range excludeDevices {
        placeholders[i] = "?"
        args = append(args, dev)
    }
    whereConditions = append(whereConditions, fmt.Sprintf(
        "device_id NOT IN (%s)", strings.Join(placeholders, ", "),
    ))
}
// … similar for excludeAreas

rows, err := duckDB.QueryContext(ctx, query, args...)
```

`orderDir` and `limit` remain interpolated because `?` isn't allowed in `ORDER BY` or `LIMIT` positions. Both are validated earlier:

- `orderDir` is derived from the `direction` enum (`highest` → `DESC`, `lowest` → `ASC`).
- `limit` is clamped to `1..100`.

### Annotation cleanup

The MCP protocol allows tools to carry a `readOnlyHint` annotation telling the client "invoking me has no side effects." 16 of 19 tools already had it. The 3 that didn't (`query_analytics`, `radiation_stats`, `query_extreme_readings`) were read-only in behavior but silent to the client. Now all 19 advertise correctly — relevant because some MCP clients (future or today's) may use this hint to decide whether a tool needs user confirmation.

---

## Overall outcome

Before this work:

- Hints were editable only by engineers, required rebuild+redeploy, and had no audit trail.
- Two MCP tools had structural paths through which an attacker-crafted prompt could potentially trigger writes (multi-statement SQL, Sprintf-with-escape).
- Three tools did not advertise their read-only status to MCP clients.

After:

- Hints are editable live via the admin UI, persist in PostgreSQL, and every change is snapshotted.
- The embedded `//go:embed hints/*.json` fallback means the file loader works in production regardless of binary location.
- All 19 MCP tools are parameterized or validator-guarded, with no user-controlled strings reaching SQL text.
- All 19 MCP tools advertise `readOnlyHint: true`.

## Follow-up ideas (not done)

- Full admin-side audit log (who changed which hint, not just the snapshot) — right now snapshots live in `ai_hints_history` but aren't tied to an admin user ID.
- A test that actually attempts a stacked-statement payload through the real DuckDB driver to confirm the guard (unit test covers the validator logic; an integration test would cover the end-to-end).
- If the AI chatbot grows tools that *should* write (e.g. letting a user submit feedback), those should go through a dedicated write path with the requesting user's session attached — not through the MCP tool surface, which is intentionally AI-callable.
