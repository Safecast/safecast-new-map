-- Migration: add cache-management columns to qa_embeddings
--
-- These columns power the /admin/qa-embeddings curation tab. They are also
-- applied automatically at unified-server startup
-- (see cmd/unified-server/duckdb_analytics.go) — this file is the standalone
-- record for ops/manual application.
--
-- Target: DuckLake `analytics` catalog (PostgreSQL + Parquet).
-- Apply with the DuckDB CLI attached to the catalog, e.g.:
--   ATTACH 'ducklake:postgres:dbname=ducklake_catalog host=localhost user=ducklake_rw'
--          AS analytics (DATA_PATH '/var/lib/safecast/ducklake/');
--   USE analytics;
--   .read migrations/add_qa_embeddings_columns.sql
--
-- Columns:
--   used_count    — bumped on every cache hit; surfaces popular Q&A in admin
--   last_used_at  — timestamp of the most recent cache hit; drives staleness UI
--   status        — 'active' | 'demoted' | 'archived'; only 'active' is eligible
--                   for cache hits and RAG retrieval
--   lang          — IETF tag (e.g. 'en', 'ja'); scopes lookups so a Japanese
--                   question never matches an English cached answer

ALTER TABLE qa_embeddings ADD COLUMN IF NOT EXISTS used_count   INTEGER     DEFAULT 0;
ALTER TABLE qa_embeddings ADD COLUMN IF NOT EXISTS last_used_at TIMESTAMPTZ;
ALTER TABLE qa_embeddings ADD COLUMN IF NOT EXISTS status       VARCHAR     DEFAULT 'active';
ALTER TABLE qa_embeddings ADD COLUMN IF NOT EXISTS lang         VARCHAR     DEFAULT '';
