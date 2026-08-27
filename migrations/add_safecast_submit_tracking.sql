-- Add api.safecast.org submission-tracking columns to uploads table.
-- safecast_import_id:       bgeigie_imports id returned by a successful submit
-- submitted_to_safecast_at: when the submit succeeded
-- safecast_submit_error:    last failure reason, if any (best-effort, non-blocking)
--
-- These columns are also added automatically at unified-server startup
-- (see pkg/database/database.go: ensureUploadsMetadataColumns) — this file is the
-- standalone record for ops/manual application.

ALTER TABLE uploads ADD COLUMN IF NOT EXISTS safecast_import_id TEXT;
ALTER TABLE uploads ADD COLUMN IF NOT EXISTS submitted_to_safecast_at TIMESTAMPTZ;
ALTER TABLE uploads ADD COLUMN IF NOT EXISTS safecast_submit_error TEXT;
