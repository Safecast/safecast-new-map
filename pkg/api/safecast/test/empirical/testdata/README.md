# Test recordings (sanitized)

All user and measurement data in these JSON files has been **sanitized**. Real names, IDs, and any third-party tokens have been replaced with placeholders so nothing sensitive lives in the repo.

- **User lists:** Names appear as `[SANITIZED] User 1`, `[SANITIZED] User 2`, etc.; IDs are small integers (1, 2, …).
- **Measurements:** IDs and `user_id` are placeholder values; `location_name` may be `[SANITIZED]` where applicable.
- **HTML responses** (e.g. no-Accept cases): Bodies are minimal HTML with no tracking or map keys.

Tests only assert response status, Content-Type, and JSON **shape** (keys and value types), not real data.
