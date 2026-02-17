-- Migration: Link historical uploads to internal users
-- This matches uploads.user_id (external Safecast user ID) to users.external_id
-- and sets the internal_user_id field for proper user profile tracking

-- Update uploads that have a user_id matching a user's external_id
UPDATE uploads u
SET internal_user_id = usr.id::text
FROM users usr
WHERE u.internal_user_id IS NULL
  AND u.user_id IS NOT NULL
  AND u.user_id != ''
  AND usr.external_id = u.user_id;

-- Show summary of what was updated
SELECT
    COUNT(*) as total_uploads_linked,
    COUNT(DISTINCT u.internal_user_id) as unique_users
FROM uploads u
WHERE u.internal_user_id IS NOT NULL;

-- Show breakdown by user
SELECT
    usr.email,
    usr.username,
    COUNT(*) as upload_count
FROM uploads u
JOIN users usr ON u.internal_user_id = usr.id::text
GROUP BY usr.email, usr.username
ORDER BY upload_count DESC
LIMIT 20;
