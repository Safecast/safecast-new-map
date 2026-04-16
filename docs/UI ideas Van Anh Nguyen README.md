# Safecast Map — UI Fix & Tutorial Layer
## Feature Requirements for Engineering
**Scope:** Frontend only. No infrastructure changes. No backend dependencies beyond what currently exists.

---

## 1. QUICK FIXES

### 1.1 Filter Panel — Label All Icons
**Current state:** Left panel shows emoji/icon-only checkboxes with no text labels.
**Required:** Add a visible text label next to each icon. Labels should be concise (1–2 words).
Suggested labels based on assumed icon meanings:
- 🧍 → "Walk"
- ✈️ → "Air"
- 🚜 / vehicle → "Vehicle"
- 📡 / antenna → "Station"
- 📊 / bar chart → "Aggregate"

**Note to engineer:** Confirm actual filter categories with the data team and update labels accordingly. The goal is zero ambiguity.

---

### 1.2 Data Point Tooltip — Add Unit and Context
**Current state:** Hovering or clicking a data point shows a floating value (e.g., "8.57") with no unit, no timestamp, no hardware type.
**Required:**
- Always display unit: `μR/h` or `nSv/h` — whichever Safecast uses as standard — inline with the value
- Add a one-line reading classification beneath the value. Example:
  - `8.57 μR/h — Within normal background range`
  - `213 μR/h — Elevated — above typical background`
- Add collection metadata: date collected, hardware type (mobile / stationary), contributor ID if public
- Classification thresholds should be hardcoded as static constants on the frontend — no backend call required. Confirm thresholds with data team.

---

### 1.3 Color Scale Legend — Add Baseline Annotation
**Current state:** Right-side legend shows a gradient from dark blue (3) to yellow (6554) with no interpretive labels.
**Required:**
- Add two annotation markers on the legend:
  - "Typical background" bracket — spanning the normal range (confirm range with data team, likely ~5–20 μR/h)
  - "Elevated" label at the threshold where readings become notable
- Add a single line of microcopy beneath the legend: `Global average background radiation: ~10 μR/h`
- Unit label (`μR/h`) should be fully visible, not clipped at the bottom edge (currently truncated in the screenshot)

---

### 1.4 Geolocation Button — Fix or Disable
**Current state:** "See your location" button is non-functional for some users.
**Required:**
- If geolocation permission is denied or unavailable, show a clear inline message: `Location access unavailable. Please search for a city above.`
- Do not silently fail. Either it works or it tells the user why it didn't.

---

### 1.5 Station vs. Mobile Data — Visual Distinction
**Current state:** Stationary sensors and mobile readings are rendered identically as dots.
**Required:**
- Differentiate visually. Suggested: stationary sensors as a small ring/hollow circle, mobile readings as filled dots. Color behavior remains the same.
- Add a two-line legend key beneath the color scale identifying the two types.
- This does not require infrastructure changes if the data type field already exists in the current data response — confirm with data team.

---

## 2. TUTORIAL LAYER

### Overview
A contextual, dismissible tutorial that activates on first visit and can be re-opened at any time. It does not replace the map — it overlays it. The user should feel oriented within 30 seconds and in control of dismissing it at any point.

**Implementation approach:** Store tutorial completion state in `localStorage`. No login required, no backend call.

---

### 2.1 Entry Point — Welcome Tooltip on Load
**Trigger:** First visit (no `safecast_tutorial_seen` key in localStorage) OR user clicks a persistent "?" help button (see 2.5).
**Behavior:**
- On map load, before the user interacts, show a single centered modal overlay (not full-screen, approximately 480px wide).
- Content:
  - **Headline:** `You're looking at real radiation data, collected by volunteers around the world.`
  - **Body (2 sentences max):** `Each dot represents a measurement taken by a Safecast device — either carried by someone on the move, or transmitted continuously from a fixed sensor. The color tells you the intensity.`
  - **Two buttons:** `Show me how to read this` (proceeds to step 2.2) and `Explore on my own` (dismisses, sets localStorage flag)
- Clicking outside the modal dismisses it and sets the localStorage flag.

---

### 2.2 Step-by-Step Highlight Tour (4 Steps)
Activated by "Show me how to read this." Each step highlights one UI element with a spotlight overlay (darken rest of map) and a positioned tooltip. User advances with a `Next →` button. Can exit at any step.

**Step 1 — The Color Scale**
- Highlight: right-side legend
- Tooltip text: `Color shows radiation intensity. Blue is typical background. As readings rise, colors shift through cyan, green, yellow — readings in the yellow range are significantly elevated above normal.`
- Add: a static annotation directly on the legend showing where "normal" sits (this can be CSS-injected during the tutorial step, no data change required)

**Step 2 — The Filters**
- Highlight: left filter panel
- Tooltip text: `Use these filters to show data by how it was collected — on foot, by vehicle, by air, or from fixed stations. You can combine them.`
- Each filter icon should display its text label during this step (labels should already exist per fix 1.1)

**Step 3 — Clicking a Data Point**
- Highlight: a pre-selected representative dot near the center of the default map view (hardcode coordinates for a dot in Tokyo that reliably exists in the dataset — confirm with data team)
- Tooltip text: `Click any dot to see the exact reading, when it was collected, and what type of device recorded it.`
- Trigger the actual tooltip for that dot during this step so the user sees a live example

**Step 4 — The Time Filter (if applicable)**
- Highlight: time/playback control (confirm location with engineer — not clearly visible in screenshot)
- Tooltip text: `Use the time filter to step through historical data. You can see how readings in an area have changed over months or years.`
- If the time feature is not present on the current build, skip this step.

**End of tour:**
- Show a small completion card: `You're ready to explore. If you have questions about the data, open the chat below.`
- Single button: `Start exploring`
- Set `safecast_tutorial_seen = true` in localStorage

---

### 2.3 Contextual Tooltips on First Interaction (Post-Tour)
For users who dismissed the tour or have already completed it, add lightweight first-interaction hints — these fire once per element, per session, then never again.

- **First time opening filters:** Small tooltip appears above panel: `Filter by collection method. Hover each icon for details.`
- **First time clicking a data point:** If the tooltip shows a raw value without context (pre-fix 1.2), show a one-time nudge: `Readings below ~20 μR/h are within normal background range globally.`
- **First time using search:** No tooltip needed — search bar behavior is self-evident.

All contextual tooltips are dismissed on click-away and tracked in localStorage per key (e.g., `safecast_hint_filter_seen`, `safecast_hint_datapoint_seen`).

---

### 2.4 Inline Explainer Card — "What is Safecast?"
**Location:** Persistent collapsible panel, accessible from a small info icon `ⓘ` near the Safecast logo in the bottom-left corner.
**Content (static HTML, no backend):**
- **What you're seeing:** One paragraph explaining the map, the dataset scale, and the two hardware types.
- **How the data is collected:** Two sentences — mobile devices carried by volunteers, stationary sensors transmitting continuously.
- **How to trust it:** One sentence on open data, CC0 license, and the fact that all raw data is publicly downloadable.
- **Link:** `→ Learn more about Safecast's methodology` (links to existing about/methodology page)

This panel is not part of the tutorial flow. It is always accessible independently.

---

### 2.5 Persistent Help Button
**Location:** Fixed position, bottom-right corner, above the chat button if one exists.
**Appearance:** Small circular `?` button, unobtrusive.
**Behavior:** Clicking it restarts the tutorial from Step 2.1 (the welcome modal), regardless of localStorage state.
**Purpose:** Gives any user a way back into orientation without hunting for it.

---

## 3. STORY WALL — "The People Behind the Data"

### Overview
This is the most engaging part of the tutorial experience and the most human entry point into understanding what Safecast is. Rather than explaining citizen science abstractly, the story wall shows it through the people who actually do it — a volunteer who mounted a device on their bike and rode across Amsterdam, someone who measured radiation near their home in Fukushima for years, a researcher who drove hundreds of kilometers through rural Ukraine.

The engineering task now is to **build the containers**. Content will be supplied by leadership separately. No story content is required to ship this feature — placeholder cards ship first, real stories replace them when ready.

**Implementation approach:** Static JSON file (`stories.json`) as the data source. No backend. Leadership edits the JSON file directly to add or update stories. Engineer documents the schema clearly (see 3.3).

---

### 3.1 Entry Point — Story Wall Trigger in Tutorial
**Where it lives:** Added as a new final step in the tutorial flow, after the existing Step 4 (time filter). It is the emotional close of the tutorial before the user begins exploring.

**Behavior:**
- After Step 4, instead of jumping directly to the completion card, show a full-width panel overlay (not a small tooltip — this gets more space).
- Headline: `The data you're looking at was collected by people like these.`
- Subhead: `Every dot on this map has a person behind it.`
- Below the headline: render 3 story cards horizontally (see 3.2 for card spec)
- Two buttons at the bottom:
  - `Read more stories →` — opens the full Story Wall page (see 3.4)
  - `Start exploring` — dismisses the tutorial, sets localStorage flag, proceeds to map

---

### 3.2 Story Card Spec
Each card is a self-contained unit. Build it to be reusable — it should render identically in the tutorial panel, the full story wall page, and any future embed.

**Card dimensions:** Approximately 280px wide, flexible height. Three cards render side by side on desktop, stacked on mobile.

**Card anatomy (top to bottom):**
1. **Photo container** — fixed aspect ratio 3:2, rounded corners. Accepts an image URL from the JSON. If no image is provided, renders a neutral placeholder with the Safecast logo centered. No broken image states.
2. **Name** — bold, 16px
3. **Location tag** — small, muted text. Format: `City, Country` or `Region, Country`
4. **Hardware badge** — small pill/tag indicating which device they used. Values: `bGeigie Nano`, `Pointcast`, `bGeigie`, or `Other`. Styled consistently with the map's existing filter labels.
5. **Story excerpt** — 2–3 sentences max. Plain text, no markdown. If the excerpt exceeds 180 characters, truncate with a `Read more` inline link.
6. **"View full story" link** — bottom of card, links to the individual story page or an expanded modal (see 3.5)

**Placeholder card (for pre-launch):**
- Photo container: light grey fill with centered Safecast logo
- Name: `Volunteer Name`
- Location: `City, Country`
- Hardware badge: `Device`
- Excerpt: `This story is coming soon. We're collecting stories from our volunteers around the world.`
- No "View full story" link on placeholder cards

---

### 3.3 Data Schema — stories.json
Engineer creates this file at project root or in `/data/`. Leadership fills it in.

```json
[
  {
    "id": "story-001",
    "name": "Volunteer Name",
    "location": "City, Country",
    "hardware": "bGeigie Nano",
    "photo_url": "",
    "excerpt": "2-3 sentence summary of their story and why they collect data.",
    "full_story_url": "",
    "is_placeholder": true
  }
]
```

**Field rules for leadership (document these in a separate README):**
- `photo_url` — direct image URL or relative path to `/assets/stories/`. Leave empty string for placeholder.
- `excerpt` — 2–3 sentences, plain text only, no links or formatting.
- `full_story_url` — can point to an external blog post, a Safecast page, or leave empty. If empty, no "View full story" link renders.
- `is_placeholder` — set to `false` when real content is ready. Placeholder cards render differently (see 3.2).
- Ship with 3 placeholder entries so the layout renders correctly from day one.

---

### 3.4 Full Story Wall Page
A standalone page, linked from the tutorial panel and accessible from the main nav.

**URL:** `/stories` or `/people` — confirm with team

**Layout:**
- Page headline: `The People Behind the Data`
- Subhead (1 sentence, static): `Safecast's dataset is built by volunteers around the world — researchers, cyclists, teachers, and curious neighbors — who carry or install devices and contribute readings to the open record.`
- Grid of story cards (see 3.2) — 3 columns desktop, 2 tablet, 1 mobile
- Cards sourced from `stories.json` — filter out `is_placeholder: true` entries on this page. Placeholder cards only appear in the tutorial panel, not the public story wall.
- If zero real stories exist yet, the page shows: `Stories coming soon. We're gathering voices from our global community of volunteers.`

---

### 3.5 Expanded Story Modal (Optional — Build if Capacity Allows)
If `full_story_url` is empty but leadership still wants a longer story displayed, the card's "View full story" link opens an in-page modal rather than navigating away.

**Modal content:**
- Full-size photo (if available)
- Name, location, hardware badge
- Full story text (new field `full_story_text` in JSON — plain text, paragraph breaks supported via `\n\n`)
- Close button top-right

**If this is out of capacity for this sprint:** skip the modal entirely. Cards without a `full_story_url` simply do not render the "View full story" link. This is acceptable for launch.

---

## 4. OUT OF SCOPE (DO NOT BUILD IN THIS CYCLE)
The following were raised in research but require infrastructure, content decisions, or design work beyond this sprint:

- AI chat integration or natural language querying
- User accounts or contributor profiles
- Story or anecdote layer on individual data points
- API access UI
- Partnership or institutional onboarding flows
- Mobile-responsive redesign

---

## 4. OPEN QUESTIONS FOR DATA TEAM (Before Engineering Starts)
1. What are the confirmed threshold values for normal / elevated / high radiation in μR/h or nSv/h?
2. What is the standard display unit across the platform — μR/h or nSv/h?
3. What fields are currently returned per data point in the map response? (Needed for tooltip fix 1.2 and station type fix 1.5)
4. Is there a reliable, always-present data point near the Tokyo default map center that can be used as the tutorial example in Step 3?
5. Is the time filter feature currently live on the map build being shipped?