-- Add comment field to uploads table to store the user-supplied description
-- from the Safecast API bgeigie_imports JSON (the "comment" field).

ALTER TABLE uploads ADD COLUMN IF NOT EXISTS comment TEXT;
