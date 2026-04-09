# RAG & Semantic Cache Architecture

The Safecast AI assistant uses a lightweight Retrieval-Augmented Generation (RAG) system
with a semantic cache layer built entirely in Go — no external embedding API required.

> **Mermaid diagram:** [rag-semantic-cache.mmd](rag-semantic-cache.mmd)
> Render with any Mermaid-compatible viewer (VS Code extension, GitHub, mermaid.live).

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          User sends a question                              │
│                    (map widget or /assistant/ page)                         │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 1 — EMBEDDING                                                         │
│                                                                             │
│  getEmbedding(question)                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Tokenise → unigrams + bigrams                                      │   │
│  │  FNV-1a 32-bit hash each token → index into 512-dim vector          │   │
│  │  L2-normalise → []float32 (512 dims)                                │   │
│  │                                                                     │   │
│  │  Pure Go, deterministic, no external API.                           │   │
│  │  Swap point: replace with Voyage AI (voyage-3-lite) if needed.      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 2 — SEMANTIC CACHE CHECK                                              │
│                                                                             │
│  checkSemanticCache(embedding)                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Load all qa_embeddings WHERE feedback_score > 0                    │   │
│  │  Cosine similarity = dot(a, b)  [both L2-normalised]                │   │
│  │  Best score ≥ 0.85 → CACHE HIT                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│              │                              │                               │
│         CACHE HIT                      CACHE MISS                          │
│              │                              │                               │
│              ▼                              ▼                               │
│  Return cached answer            Continue to Step 3                        │
│  Send done {chat_id, cached:true}                                          │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │ (cache miss only)
                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 3 — RAG CONTEXT INJECTION                                             │
│                                                                             │
│  buildRAGContext(embedding)            getLocationKnowledge()               │
│  ┌──────────────────────────────┐    ┌──────────────────────────────────┐  │
│  │ Top-3 past Q&A with          │    │ Up to 20 curated location notes  │  │
│  │ similarity ≥ 0.50            │    │ from location_knowledge table    │  │
│  │ Answers truncated to 500 ch  │    │ (auto-built from 👍 feedback)    │  │
│  └──────────────────────────────┘    └──────────────────────────────────┘  │
│                 │                                   │                       │
│                 └──────────────┬────────────────────┘                       │
│                                ▼                                            │
│                   enrichSystemPrompt(base + RAG + locations)                │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 4 — AGENTIC LOOP (Claude + MCP tools)                                 │
│                                                                             │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │  Call Claude API (claude-sonnet-4-5) with enriched system prompt     │  │
│  │       │                                                              │  │
│  │       ├─ Claude calls MCP tools (port 3333):                        │  │
│  │       │   query_radiation, sensor_current, get_track, search_area…  │  │
│  │       │                                                              │  │
│  │       └─ Tool results fed back → Claude iterates until end_turn      │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                 │                                           │
│                            Final answer                                     │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 5 — STORE & LOG                                                       │
│                                                                             │
│  logChatQuestionWithAnswer()           storeQAEmbeddingAsync()              │
│  ┌──────────────────────────────┐    ┌──────────────────────────────────┐  │
│  │ chat_questions table         │    │ qa_embeddings table (async)      │  │
│  │  id = embeddingChatID        │    │  chat_id = embeddingChatID       │  │
│  │  question, answer, source,   │    │  question, answer, embedding     │  │
│  │  ip, browser, country, model │    │  feedback_score = 0              │  │
│  └──────────────────────────────┘    └──────────────────────────────────┘  │
│                                                                             │
│  Send done event { chat_id: embeddingChatID } → browser stores it          │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 6 — USER FEEDBACK (optional)                                          │
│                                                                             │
│  User clicks 👍 or 👎                                                       │
│  POST /api/feedback { chat_id, score: +1|-1 }                              │
│                                                                             │
│  RecordFeedback()                                                           │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │  INSERT chat_feedback (chat_id, score)   ← admin dashboard reads    │  │
│  │  UPDATE qa_embeddings SET feedback_score = feedback_score + score    │  │
│  │                                                                      │  │
│  │  If 👍: extractLocationKnowledge() [async]                          │  │
│  │    → regex search answer for lat/lon coordinates                    │  │
│  │    → INSERT location_knowledge (lat, lon, radius_m=1000, note)      │  │
│  │      (feeds back into Step 3 for future questions)                  │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

## DuckLake Tables

All tables are persisted via DuckLake (PostgreSQL catalog + Parquet files at
`/var/lib/safecast/ducklake/`). They are created automatically on startup by
`createDuckDBSchema()` in `cmd/unified-server/duckdb_analytics.go`.

| Table | Purpose | Key columns |
|-------|---------|-------------|
| `qa_embeddings` | Semantic cache store | `chat_id`, `question`, `answer`, `embedding` (JSON float32[512]), `feedback_score` |
| `chat_questions` | Full question/answer log with metadata | `id` (= chat_id), `question`, `answer`, `source`, `country`, `model` |
| `chat_feedback` | User thumbs up/down votes | `chat_id`, `score` (+1 or -1) |
| `location_knowledge` | Curated geographic context (auto-extracted from 👍 answers) | `lat`, `lon`, `radius_m`, `note` |

All four tables share `chat_id` / `id` = `time.Now().UnixMilli()` as the common key.

## Embedding: Feature Hashing

Embeddings are generated in `cmd/unified-server/embeddings.go` using **FNV-1a hashing**:

1. Tokenise: lowercase, split on non-alphanumeric, generate unigrams + bigrams
2. Hash each token: `fnv32a(token) % 512` → increment that dimension
3. L2-normalise the 512-dim vector

This is fully deterministic, requires no external API, and produces consistent
similarity scores for domain-specific radiation/sensor vocabulary.

**Swap point:** To upgrade to neural embeddings, replace `getEmbedding()` in
`embeddings.go` with a call to Voyage AI (`voyage-3-lite`). The return type
`[]float32` and all downstream code remain unchanged.

## Cache Thresholds

| Threshold | Value | Effect |
|-----------|-------|--------|
| `cacheHitThreshold` | 0.85 | Return cached answer, skip Claude entirely |
| `ragContextThreshold` | 0.50 | Include past Q&A in system prompt as context |
| `ragTopK` | 3 | Max past Q&A pairs injected per request |

Cache entries are only eligible if `feedback_score > 0` — at least one thumbs-up
vote is required before a cached answer can be served to other users.

## Feedback Loop

```
👍 vote  →  feedback_score +1  →  answer eligible for cache hits
              +  lat/lon extracted from answer  →  location_knowledge populated
              +  future questions near that location get curated context injected

👎 vote  →  feedback_score -1  →  answer eventually filtered out of cache
```

## Key Files

| File | Responsibility |
|------|---------------|
| `cmd/unified-server/embeddings.go` | Feature hashing, cosine similarity |
| `cmd/unified-server/semantic_cache.go` | Cache check, RAG context, feedback recording, location extraction |
| `cmd/unified-server/mcp_register.go` | Chat handler, agentic loop, streaming, embedding lifecycle |
| `cmd/unified-server/chat_logging.go` | Insert into `chat_questions` |
| `cmd/unified-server/duckdb_analytics.go` | DuckLake init, table schemas |
| `cmd/unified-server/admin_mcp.go` | Admin API for analytics (including feedback counts) |

## Local Development

See [ducklake-local-setup.md](ducklake-local-setup.md) for DuckLake connection setup
(`DUCKLAKE_PG_URL`, `DUCKLAKE_DATA_PATH`). Without DuckLake, the server starts fine
but all analytics and caching are silently disabled.
