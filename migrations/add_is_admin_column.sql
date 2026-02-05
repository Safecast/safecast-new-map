-- Add is_admin column to users table
-- PostgreSQL version
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN DEFAULT FALSE;

-- For SQLite (run separately if needed):
-- ALTER TABLE users ADD COLUMN is_admin INTEGER DEFAULT 0;
