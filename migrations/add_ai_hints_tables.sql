-- Migration: Add ai_hints + ai_hints_history tables
-- Run: psql -h 127.0.0.1 -U postgres -d safecast -f migrations/add_ai_hints_tables.sql

CREATE TABLE IF NOT EXISTS ai_hints (
  id                      BIGSERIAL PRIMARY KEY,
  model                   VARCHAR(64) NOT NULL UNIQUE,
  display_name            TEXT NOT NULL,
  capabilities            JSONB NOT NULL DEFAULT '[]'::jsonb,
  system_prompt           TEXT NOT NULL DEFAULT '',
  tools                   JSONB NOT NULL DEFAULT '{}'::jsonb,
  global_formatting_rules JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at              TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_ai_hints_model       ON ai_hints(model);
CREATE INDEX IF NOT EXISTS idx_ai_hints_deleted_at  ON ai_hints(deleted_at);

CREATE TABLE IF NOT EXISTS ai_hints_history (
  id          BIGSERIAL PRIMARY KEY,
  model       VARCHAR(64) NOT NULL,
  snapshot    JSONB NOT NULL,
  change_kind VARCHAR(32) NOT NULL,
  changed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_hints_history_model      ON ai_hints_history(model);
CREATE INDEX IF NOT EXISTS idx_ai_hints_history_changed_at ON ai_hints_history(changed_at DESC);
