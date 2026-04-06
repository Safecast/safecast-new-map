# Admin Panel

Guide to administering the Safecast New Map platform.

[[Home|← Back to Home]]

---

## Overview

The admin panel provides comprehensive management tools for:
- User accounts and permissions
- Upload management and moderation
- MCP analytics and monitoring
- Real-time sensor management
- Translation management (29 languages)
- Cache management
- Track management

---

## Enable Admin Panel

### Method 1: Admin Password Flag

```bash
./safecast-new-map -admin-password your-secure-password
```

Access at: `http://localhost:8765/admin/users?password=your-secure-password`

### Method 2: Admin User Login

1. Create an admin user:
```bash
psql -d safecast -c "UPDATE users SET is_admin = true WHERE email = 'admin@example.com';"
```

2. Log in as the admin user
3. Access admin panel via username → "Admin"

---

## Admin Pages

### Users Management

**URL:** `/admin/users`

**Features:**
- View all registered users
- Search users by email or username
- Edit user details (username, email)
- Grant/revoke admin privileges
- Activate/deactivate accounts
- Regenerate API keys
- Reset user passwords
- Delete users

#### User List

| Column | Description |
|--------|-------------|
| ID | User ID |
| Email | User email address |
| Username | Display name |
| API Key | User's API key |
| Verified | Email verification status |
| Active | Account status |
| Admin | Admin privileges |
| Created | Registration date |
| Last Login | Last login timestamp |

#### User Actions

**Edit User:**
- Change username
- Update email
- Toggle admin status
- Activate/deactivate account

**Regenerate API Key:**
```bash
# API call
POST /api/admin/users/{userId}/regenerate-api-key?password=admin-password
```

**Reset Password:**
```bash
# API call
POST /api/admin/users/{userId}/reset-password?password=admin-password
```

**Delete User:**
```bash
# API call
DELETE /api/admin/users/{userId}?password=admin-password
```

#### Create New User

```bash
POST /api/admin/users/create?password=admin-password
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "secure-password",
  "username": "newuser",
  "is_admin": false
}
```

---

### Uploads Management

**URL:** `/admin/uploads`

**Features:**
- View all uploaded files
- Search uploads by filename, username, or date
- Update upload metadata
- Import from Safecast API
- Delete tracks and uploads
- Bulk operations

#### Upload List

| Column | Description |
|--------|-------------|
| ID | Upload ID |
| Filename | Original filename |
| Source | Upload source (web, API, Safecast) |
| Uploaded | Upload timestamp |
| Recording Date | Measurement date |
| Detector | Device type |
| Username | Uploader |
| Track Count | Number of measurements |

#### Update Upload Metadata

```bash
POST /api/admin/uploads/update?password=admin-password
Content-Type: application/json

{
  "uploadId": 123,
  "recordingDate": "2024-01-15T10:00:00Z",
  "detector": "bGeigie",
  "username": "john_doe",
  "name": "Tokyo Urban Survey",
  "notes": "Measured in central Tokyo",
  "comment": "High readings near station"
}
```

#### Import from Safecast API

**Import by Safecast ID:**
```bash
POST /api/admin/import-from-safecast?password=admin-password
Content-Type: application/json

{
  "safecastId": 12345
}
```

**Import by ID Range:**
```bash
POST /api/admin/import-by-id?password=admin-password
Content-Type: application/json

{
  "startId": 100000,
  "endId": 110000
}
```

**Backfill by Country:**
```bash
POST /api/admin/backfill-countries?password=admin-password
Content-Type: application/json

{
  "countries": ["Japan", "United States", "Germany"]
}
```

#### Delete Tracks

**Single Track:**
```bash
POST /api/admin/delete?password=admin-password
Content-Type: application/json

{
  "trackId": "abc123"
}
```

**Multiple Tracks:**
```bash
POST /api/admin/delete-multiple?password=admin-password
Content-Type: application/json

{
  "trackIds": ["abc123", "def456", "ghi789"]
}
```

#### Backfill Metadata from Old Safecast API

```bash
POST /api/admin/tracks/import-safecast?password=admin-password
Content-Type: application/json

{
  "trackId": "abc123"
}
```

---

### MCP Analytics

**URL:** `/admin/mcp`

**Features:**
- Monitor MCP tool usage
- View AI query logs
- Track tool execution statistics
- Analyze chat questions
- Monitor performance metrics

#### Tool Usage Statistics

| Metric | Description |
|--------|-------------|
| Tool Name | MCP tool identifier |
| Executions | Total number of calls |
| Avg Response Time | Average execution time |
| Cache Hit Rate | Percentage of cached responses |

#### AI Query Log

View all AI queries with:
- Query text
- Timestamp
- User session
- Tools invoked
- Response time
- Cache status

**Example Queries:**
```
"What's the radiation level in Tokyo?"
"Show me the highest readings in Japan"
"Find tracks near Chernobyl"
"List all sensors currently online"
```

#### Performance Monitoring

Track:
- Query execution times
- Cache hit/miss ratios
- Tool usage patterns
- Peak usage times

---

### Real-Time Sensors

**URL:** `/admin/realtime`

**Features:**
- View all active sensors
- Monitor sensor status
- View current readings
- Manage sensor data
- Configure polling settings

#### Sensor List

| Column | Description |
|--------|-------------|
| ID | Sensor identifier |
| Name | Sensor display name |
| Type | Sensor type (Pointcast, Solarcast, bGeigieZen) |
| Location | GPS coordinates |
| Current Reading | Latest dose rate |
| Last Update | Last data received |
| Status | Online/Offline |

#### Sensor Actions

**Refresh Sensor Data:**
- Force immediate poll
- Update current readings
- Refresh map display

**Configure Polling:**
```bash
./safecast-new-map \
  -safecast-realtime \
  -safecast-fetcher-interval 5m
```

**Delete Sensor Data:**
- Remove stale readings
- Clean up offline sensors
- Reset sensor history

---

### Translations Management

**URL:** `/admin/translations`

**Features:**
- Edit UI translations for 29 languages
- Search translations by key or value
- Add new translation keys
- Delete translation keys
- Reload translations into memory

#### Supported Languages

| Code | Language |
|------|----------|
| ar | Arabic |
| bg | Bulgarian |
| cs | Czech |
| da | Danish |
| de | German |
| el | Greek |
| en | English |
| es | Spanish |
| fa | Persian |
| fi | Finnish |
| fr | French |
| he | Hebrew |
| hi | Hindi |
| hu | Hungarian |
| id | Indonesian |
| it | Italian |
| ja | Japanese |
| ko | Korean |
| ms | Malay |
| nl | Dutch |
| no | Norwegian |
| pl | Polish |
| pt | Portuguese |
| ru | Russian |
| sv | Swedish |
| th | Thai |
| tr | Turkish |
| uk | Ukrainian |
| vi | Vietnamese |
| zh | Chinese |

#### Translation Workflow

1. **Select language** from dropdown
2. **Search** for specific keys or values
3. **Edit** translations inline
4. **Save** changes to database
5. **Reload into Memory** to apply changes

#### Add New Translation Key

```bash
POST /api/admin/translations?password=admin-password
Content-Type: application/json

{
  "key": "new.translation.key",
  "translations": {
    "en": "English text",
    "ja": "日本語テキスト",
    "fr": "Texte français"
  }
}
```

#### Reload Translations

Click "Reload into Memory" button to apply changes without server restart.

**Brand Rule:** "Safecast" must remain untranslated in all languages.

---

### Cache Management

**URL:** Accessible from any admin page

**Features:**
- Clear LRU cache
- Force data refresh
- View cache statistics

#### Clear Cache

```bash
POST /api/admin/cache?password=admin-password
```

**Effects:**
- Clears all cached map tiles
- Forces fresh database queries
- Invalidates track info cache
- Resets translation cache

#### Cache Statistics

View:
- Cache size
- Hit rate
- Memory usage
- Oldest entries

---

## Admin API Reference

### Authentication

All admin endpoints require authentication via:
- Admin password in query parameter: `?password=admin-password`
- Admin user session cookie

### Endpoints

See [API Documentation](API-Documentation#admin-endpoints) for complete admin API reference.

---

## Database Schema

### Admin-Related Tables

**Users:**
```sql
SELECT id, email, username, is_admin, is_active, created_at, last_login_at
FROM users
WHERE is_admin = true;
```

**Uploads:**
```sql
SELECT id, filename, source, uploaded_at, recording_date, detector, username
FROM uploads
ORDER BY uploaded_at DESC;
```

**Translations:**
```sql
SELECT language_code, key, value
FROM translations
WHERE language_code = 'en'
ORDER BY key;
```

**Real-time Measurements:**
```sql
SELECT sensor_id, doseRate, timestamp, lat, lon
FROM realtime_measurements
ORDER BY timestamp DESC;
```

---

## Security Best Practices

### Admin Access

1. **Use strong admin password** (minimum 16 characters)
2. **Enable HTTPS** for admin access
3. **Restrict admin access by IP** (via firewall or proxy)
4. **Monitor admin actions** in logs
5. **Use admin user accounts** instead of password flag in production

### Audit Trail

Log all admin actions:
- User creation/deletion
- Upload modifications
- Translation changes
- Cache clears
- Import operations

### Backup Before Changes

Before bulk operations:
```bash
# Backup database
pg_dump -h localhost -U safecast_user safecast > backup_$(date +%Y%m%d).sql

# Or use tool
./tools/backup_database.sh
```

---

## Common Admin Tasks

### Link Historical Uploads to Users

```bash
# Run migration
psql -d safecast -f migrations/link_historical_uploads_to_users.sql

# Or use tool
go run ./cmd/tools/add-internal-user-id
```

### Populate Usernames from Safecast API

```bash
./tools/populate_usernames.sh
```

### Fix Recording Dates

```bash
# PostgreSQL
psql -d safecast -f migrations/fix_recording_dates.sql

# SQLite
sqlite3 data.db < migrations/fix_recording_dates_sqlite.sql
```

### Refresh Track Statistics

```bash
./tools/refresh_track_stats.sh
```

### Clean Up Test Users

```bash
go run ./cmd/tools/cleanup-test-users
```

---

## Troubleshooting

### Can't Access Admin Panel

**Check admin password:**
```bash
# Verify flag is set
./safecast-new-map -admin-password your-password
```

**Check admin user:**
```sql
SELECT id, email, is_admin FROM users WHERE is_admin = true;
```

### Translation Changes Not Applied

**Reload translations:**
1. Go to `/admin/translations`
2. Click "Reload into Memory"

**Or restart server:**
```bash
systemctl restart safecast
```

### MCP Analytics Not Showing

**Check DuckLake configuration:**
```bash
echo $DUCKLAKE_PG_URL
echo $DUCKLAKE_DATA_PATH
```

**Verify logs table:**
```sql
SELECT COUNT(*) FROM mcp_query_log;
```

---

## See Also

- [User Authentication](User-Authentication) - User management
- [API Documentation](API-Documentation) - Admin API endpoints
- [Internationalization](Internationalization) - Translation system
- [MCP Server & AI Integration](MCP-Server-AI-Integration) - MCP analytics
- [Database Maintenance](Database-Maintenance) - Database utilities
