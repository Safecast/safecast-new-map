-- Migration: Add translations table
-- Run: psql -h 127.0.0.1 -U postgres -d safecast -f migrations/add_translations_table.sql

CREATE TABLE IF NOT EXISTS translations (
  id              BIGSERIAL PRIMARY KEY,
  language_code   VARCHAR(10) NOT NULL,
  key             VARCHAR(255) NOT NULL,
  value           TEXT NOT NULL,
  updated_at      TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT translations_unique UNIQUE (language_code, key)
);

CREATE INDEX IF NOT EXISTS idx_translations_lang ON translations(language_code);
CREATE INDEX IF NOT EXISTS idx_translations_key ON translations(key);
