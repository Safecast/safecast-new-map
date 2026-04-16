# Tour Steps Admin — Database-backed, Multilingual

**Status:** Draft for team discussion
**Date:** 2026-04-16
**Author:** Rob (with Claude)

## Background

The map tour is currently a hardcoded JavaScript array in
[`cmd/unified-server/public_html/map.html:11611-11630`](../../../cmd/unified-server/public_html/map.html#L11611-L11630):

```js
var tourSteps = [
  { selector: '#map', center: true, text: 'Welcome to the Safecast map — …' },
  { selector: '.upload-btn-container', text: 'Upload your own radiation measurements…' },
  …
];
```

Changing step text, reordering, adding/removing steps, or translating into any
of our 30 supported languages requires a code change, rebuild and deploy. There
is no way for a non-developer admin to maintain the tour.

## Goal

Make the tour fully editable at runtime through an admin UI, multilingual from
day one, with support for conditional display (logged in / admin / feature
gates / viewport / first-time only), and drag-and-drop reordering.

Later, the same pattern will be applied to the Stories wall — out of scope for
this spec but kept in mind so the data-model split generalises.

## Non-goals

- Per-language enable/disable of individual steps (language filtering not
  requested).
- A full in-page WYSIWYG overlay editor. We considered it and picked a lighter
  hybrid; see "Alternatives considered".
- Versioning / history of tour edits. Admin edits are DB state; migrations are
  tracked in git.
- Bulk import/export UI.

## Design

### 1. Data model

A new PostgreSQL table `tour_steps` holds structural data; text lives in the
existing `translations` table.

```sql
CREATE TABLE tour_steps (
  id              BIGSERIAL PRIMARY KEY,
  step_key        VARCHAR(64) NOT NULL UNIQUE,  -- "welcome", "upload_btn"…
  sort_order      INTEGER NOT NULL,
  selector        TEXT NOT NULL,                -- CSS selector of target
  center          BOOLEAN NOT NULL DEFAULT FALSE,
  enabled         BOOLEAN NOT NULL DEFAULT TRUE,
  require_login   BOOLEAN,                      -- NULL = don't care
  require_admin   BOOLEAN,
  show_if_feature VARCHAR(64),                  -- 'ai_enabled' | 'realtime' | NULL
  viewport        VARCHAR(16) DEFAULT 'any',    -- 'any' | 'desktop' | 'mobile'
  first_time_only BOOLEAN NOT NULL DEFAULT FALSE,
  updated_at      TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_tour_steps_order ON tour_steps(sort_order);
```

Step text lives in the existing `translations` table under keys:

- `tour.<step_key>.text` — the step body
- `tour.<step_key>.title` — optional heading (reserved for future use)

This gives us multilingual editing for free, reusing the translations admin
page, the seed / cache / reload flow, and the existing language filter.

**Seeding.** A new `seedTourStepsDB()` function runs at startup after
`seedTranslationsDB()`. It inserts the 9 current steps with
`INSERT … ON CONFLICT DO NOTHING` so admin edits are never overwritten. It
also seeds English text into `translations` for each step using the same
ON CONFLICT approach.

Migration file: `migrations/add_tour_steps_table.sql`.

### 2. Backend (Go)

Mirrors the pattern established by `admin_translations*.go`.

New files under `cmd/unified-server/`:

- `tour_steps.go` — `TourStep` struct, `tourService` with `List()`, `Get()`,
  `Create()`, `Update()`, `Delete()`, `Reorder([]int64)`, `seedTourStepsDB()`.
- `admin_tour_steps.go` — HTTP handlers.
- `tour_steps_test.go` — service + handler tests.

#### Public API (used by `map.html`)

```
GET /api/tour/steps?lang=xx
```

Returns only enabled steps, ordered, with text resolved for the requested
language (falling back to English, then to the first available language, as
`translationsForLang` already does):

```json
[
  {
    "step_key": "welcome",
    "selector": "#map",
    "center": true,
    "text": "Welcome to the Safecast map — …",
    "conditions": {
      "require_login": null,
      "require_admin": null,
      "show_if_feature": null,
      "viewport": "any",
      "first_time_only": false
    }
  }
]
```

#### Admin API (requires admin auth, same as translations admin)

| Method | Path                                  | Purpose                                              |
|--------|---------------------------------------|------------------------------------------------------|
| GET    | `/api/admin/tour/steps`               | List all (enabled + disabled)                        |
| POST   | `/api/admin/tour/steps`               | Create; auto-assigns `sort_order = MAX+1`            |
| PUT    | `/api/admin/tour/steps/{id}`          | Update any column                                    |
| DELETE | `/api/admin/tour/steps/{id}`          | Delete row and its `tour.<step_key>.*` translations  |
| POST   | `/api/admin/tour/steps/reorder`       | Body `{"ids":[3,1,2,...]}`, rewrites `sort_order`    |

### 3. Admin UI — `/admin/tour`

New page `cmd/unified-server/public_html/admin-tour.html`, visual style
identical to `admin-translations.html`.

Layout:

- Language dropdown (same 30 langs) — controls which language's text you edit.
- **New step** button.
- Table with columns: *drag-handle | order | step_key | selector | text
  (current lang) | conditions badge | enabled toggle | actions*.
- Drag-and-drop reordering via **SortableJS** (one CDN link, ~10 KB) → on drop,
  the page sends `POST /api/admin/tour/steps/reorder`.
- Inline edit-in-place for `text` (textarea) and `selector` (text input).
- "Conditions" column shows compact badges (e.g. `login-only`, `admin-only`,
  `feature:ai_enabled`, `desktop-only`, `first-time`). Click opens a modal
  with the checkboxes/dropdowns.
- **Pick element** icon button next to the selector field. Opens `/?pick=1`
  in a new tab; user clicks an element on the map; the tab posts the computed
  selector back via `window.opener.postMessage`, then closes.
- **Preview tour** button. Opens `/?preview_tour=1` in a new tab and runs the
  tour end-to-end using live DB data.

### 4. Changes to `map.html`

Replace the hardcoded `tourSteps` array with a fetch at tour start:

```js
async function startTour() {
  const lang = currentLang || 'en';
  const res  = await fetch('/api/tour/steps?lang=' + lang);
  tourSteps  = (await res.json()).filter(stepMatchesConditions);
  if (tourSteps.length === 0) return;
  tourIndex = 0;
  showTourStep(0);
}
```

`stepMatchesConditions(step)` evaluates each condition against runtime state:

- `require_login` → `window.currentUser` present?
- `require_admin` → `window.currentUser?.is_admin`?
- `show_if_feature` → `window.featureFlags[...]`?
- `viewport` → `window.matchMedia('(max-width: 768px)')`
- `first_time_only` → `localStorage['safecast_tutorial_seen']` not set?

The existing "hidden element → skip" safeguard in `showTourStep` stays
unchanged; condition filtering is an additional pre-flight filter.

#### `?pick=1` mode (element picker)

Loaded only when the query parameter is present (~150 lines, behind a guard):

- Outlines every hoverable element with a coloured border.
- On click, computes a stable CSS selector. Preference order:
  1. `#id` if the element has an id and it's unique on the page.
  2. `.unique-class` if exactly one element has that class.
  3. A tag + `:nth-child` path from the nearest ancestor with an id.
- Posts `{type:'tour-selector-picked', selector}` to `window.opener`, then
  closes the tab.

No new dependencies.

#### `?preview_tour=1` mode

Forces `startTour()` on load and ignores the `safecast_tutorial_seen` flag.

### 5. Rollout

1. Create and run `migrations/add_tour_steps_table.sql` against local Postgres.
2. Ship backend + seeder. Verify `GET /api/tour/steps?lang=en` returns the
   same nine steps and the existing tour (still reading the hardcoded array)
   still works.
3. Ship the admin page. Smoke-test add / edit / delete / reorder / preview.
4. Ship `map.html` client changes. The hardcoded array is removed in this
   commit so any regression is easy to bisect.
5. Single push to `main`; GitHub Actions pipeline deploys all three binaries.

### 6. Swagger / API docs

All new handlers must carry the same Swagger annotations used elsewhere in the
codebase (`@Summary`, `@Description`, `@Tags`, `@Param`, `@Success`,
`@Failure`, `@Router`) — see [`admin_translations.go`](../../../cmd/unified-server/admin_translations.go)
for the established pattern. Tag the public handler `tour` and the admin
handlers `admin`.

After adding the handlers, regenerate the swagger artifacts:

```bash
cd cmd/unified-server && swag init -g doc.go -o docs/api \
  --parseDependency --parseInternal --parseDependencyLevel 2 \
  --instanceName unifiedapi
```

The combined Map API + MCP API docs page at
`https://simplemap.safecast.org/docs/` will pick up the new endpoints on next
deploy.

### 7. Testing

Automated:

- `tour_steps_test.go` — service CRUD, `Reorder` transaction behaviour,
  `seedTourStepsDB` idempotency.
- `admin_tour_steps_test.go` — handler auth gates + request/response shapes.

Manual:

- Run the tour in English, Japanese, and one RTL language (Arabic or Farsi);
  confirm text renders correctly.
- Toggle each condition type in the admin and verify the step shows/hides
  appropriately in a fresh browser session.
- Exercise the element picker on four distinct targets: a button with an id,
  a div with a unique class, a nav item with no id, and an SVG element.

## Alternatives considered

### A. Standard admin table only

Simpler (no picker, no preview button), but writing correct CSS selectors by
hand is the actual pain point — this approach leaves that pain in place.

### B. Full in-page overlay editor

Admin mode on the live map: click each spotlight to edit text, click an
element to retarget, drag to reorder. WYSIWYG but heavy (~2-3 weeks). Editing
multilingual text inside an overlay is awkward, and you'd still want an
underlying list view for bulk operations. Rejected for cost/benefit.

### C. Hybrid (chosen)

Standard admin table for all editing, plus a small "Pick element" button that
solves the one genuinely hard part (writing a correct selector). Estimated
~1.5 weeks. Keeps the existing admin pattern and adds minimal new UI code.

## Risks / open questions

- **`first_time_only` semantics.** Today the whole tutorial modal is
  first-visit-only via `safecast_tutorial_seen`. Per-step first-time-only
  would need a per-step localStorage key. Proposed default: interpret
  `first_time_only` as "hide this step on repeat tour runs" using the
  existing `safecast_tutorial_seen` flag — simpler, no new storage.
- **Feature-gate list.** Starts as a frozen enum (`ai_enabled`,
  `realtime_enabled`, `logged_in_uploads`) shown as a dropdown in the admin.
  Frozen, not free-text, so typos can't silently break the tour.
- **Step count growth.** Fine for the current nine steps; no pagination
  needed. Revisit if the tour grows past ~50 steps.

## Future work (not in this spec)

- **Stories admin.** Apply the same split (structural table +
  translations-keyed text) to the Stories wall currently backed by
  `/data/stories.json`. The `step_key` → `slug` mapping transfers directly;
  most of the backend/admin scaffolding from this spec can be reused.
