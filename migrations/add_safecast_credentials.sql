-- Add api.safecast.org credential columns to users table.
-- safecast_api_key:  the user's own api.safecast.org API key, entered on /profile
-- safecast_user_id:  numeric Safecast user id, resolved via GET /users/me.json
--                     when safecast_api_key is set; not user-editable
--
-- These columns are also added automatically at unified-server startup
-- (see pkg/database/database.go: ensureUsersSafecastColumns) — this file is the
-- standalone record for ops/manual application.

ALTER TABLE users ADD COLUMN IF NOT EXISTS safecast_api_key TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS safecast_user_id TEXT;
