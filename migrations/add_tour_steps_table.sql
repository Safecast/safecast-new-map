-- Migration: Add tour_steps table for DB-backed, multilingual map tour
-- Run: psql -h 127.0.0.1 -U postgres -d safecast -f migrations/add_tour_steps_table.sql
--
-- Structural data (selector, order, conditions) lives here; per-language text
-- lives in the existing `translations` table under keys of the form
-- `tour.<step_key>.text`.

CREATE TABLE IF NOT EXISTS tour_steps (
  id              BIGSERIAL PRIMARY KEY,
  step_key        VARCHAR(64) NOT NULL UNIQUE,
  sort_order      INTEGER     NOT NULL,
  selector        TEXT        NOT NULL,
  center          BOOLEAN     NOT NULL DEFAULT FALSE,
  enabled         BOOLEAN     NOT NULL DEFAULT TRUE,
  require_login   BOOLEAN,
  require_admin   BOOLEAN,
  show_if_feature VARCHAR(64),
  viewport        VARCHAR(16) NOT NULL DEFAULT 'any',
  first_time_only BOOLEAN     NOT NULL DEFAULT FALSE,
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT tour_steps_viewport_check CHECK (viewport IN ('any', 'desktop', 'mobile'))
);

CREATE INDEX IF NOT EXISTS idx_tour_steps_order ON tour_steps(sort_order);
