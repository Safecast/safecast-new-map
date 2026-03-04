-- Add index on uploads.detector for efficient device name filtering
CREATE INDEX IF NOT EXISTS idx_uploads_detector ON uploads(detector) WHERE detector IS NOT NULL;

-- Create a case-insensitive index for ILIKE queries
CREATE INDEX IF NOT EXISTS idx_uploads_detector_lower ON uploads(LOWER(detector)) WHERE detector IS NOT NULL;
