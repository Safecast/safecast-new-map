-- Add admin-editable metadata fields to uploads table.
-- name:  human-readable display name (imported from Safecast API or set by admin)
-- notes: free-text admin annotation

ALTER TABLE uploads ADD COLUMN IF NOT EXISTS name  VARCHAR;
ALTER TABLE uploads ADD COLUMN IF NOT EXISTS notes VARCHAR;

-- Back-fill name from filename for existing rows so the column is never blank.
UPDATE uploads SET name = filename WHERE name IS NULL OR name = '';
