# User Authentication

Guide to user registration, login, API keys, and security features.

[[Home|← Back to Home]]

---

## Overview

The platform provides comprehensive user authentication with:
- Email-based registration and verification
- Session-based authentication (30-day sessions)
- API key authentication for programmatic access
- Password reset via email
- Role-based access (admin, regular user)
- Comprehensive authentication logging

---

## Enable Authentication

### Basic Setup

```bash
./safecast-new-map \
  -allow-registration \
  -session-secret "your-random-secret-key"
```

### With Email Support

```bash
./safecast-new-map \
  -allow-registration \
  -require-auth \
  -smtp-host smtp.gmail.com \
  -smtp-port 587 \
  -smtp-username your-email@gmail.com \
  -smtp-password your-app-password \
  -smtp-from your-email@gmail.com \
  -session-secret "your-random-secret-key" \
  -base-url "https://your-domain.com"
```

**Flags:**
- `-allow-registration` - Enable user registration
- `-require-auth` - Require authentication for uploads
- `-session-secret` - Secret key for session encryption (required)
- `-smtp-*` - Email configuration for verification and password reset
- `-base-url` - Base URL for email links

---

## User Registration

### Web Registration

1. Navigate to the login page
2. Click "Register"
3. Enter email, username, and password
4. Check email for verification link
5. Click verification link to activate account

### API Registration

```bash
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "username": "johndoe"
}
```

**Response:**
```json
{
  "message": "Registration successful. Please check your email to verify your account."
}
```

### Email Verification

After registration, users receive a verification email:

```
Click this link to verify your email:
https://your-domain.com/api/auth/verify-email?token=verification-token
```

**Verify via API:**
```bash
GET /api/auth/verify-email?token={token}
```

---

## Login Methods

### Password Login

**Web Interface:**
1. Enter email and password
2. Click "Login"
3. Session cookie is set automatically

**API:**
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your-password"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "user": {
    "id": 42,
    "email": "user@example.com",
    "username": "johndoe",
    "api_key": "abcdefghijklmnopqrst"
  }
}
```

### API Key Login

**Web Interface:**
1. Enter email
2. Enter API key instead of password
3. Click "Login"

**API:**
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "api_key": "your-20-char-api-key"
}
```

### Using API Keys

**Header method (recommended):**
```bash
curl -H "X-API-Key: your-api-key" \
  http://localhost:8765/api/protected-endpoint
```

**Query parameter method:**
```bash
curl "http://localhost:8765/api/protected-endpoint?api_key=your-api-key"
```

---

## API Keys

### Finding Your API Key

1. Log in to your account
2. Click your username → "Profile"
3. Your API key is displayed in the "API Key" section
4. Click "Copy" to copy to clipboard

### API Key Uses

- **Web login** - Alternative to password authentication
- **Programmatic uploads** - Automated data submission
- **API access** - Protected endpoint authentication
- **CLI tools** - Command-line authentication

### Security Best Practices

✅ **Do:**
- Treat API key like a password
- Use header method for authentication
- Rotate keys periodically
- Use environment variables in scripts

❌ **Don't:**
- Share your API key
- Commit API keys to version control
- Use API keys in client-side code
- Send API keys in URLs (use headers instead)

### Regenerate API Key

**Self-service:**
1. Go to Profile page
2. Click "Regenerate API Key"
3. Confirm the action
4. Copy new key immediately (old key becomes invalid)

**Admin regenerate:**
```bash
POST /api/admin/users/{userId}/regenerate-api-key?password=admin-password
```

**Response:**
```json
{
  "message": "API key regenerated successfully",
  "api_key": "new-20-char-key",
  "user_id": 42
}
```

The user receives an email notification with their new API key.

---

## Password Management

### Password Reset

**Web Interface:**
1. Click "Forgot Password" on login page
2. Enter email address
3. Check email for reset link
4. Click link and enter new password

**API:**
```bash
# Request reset
POST /api/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}

# Reset password
POST /api/auth/reset-password
Content-Type: application/json

{
  "token": "reset-token-from-email",
  "password": "new-secure-password"
}
```

### Change Password

**Web Interface:**
1. Log in
2. Go to Profile
3. Click "Change Password"
4. Enter current and new password

**API:**
```bash
POST /api/user/change-password
Content-Type: application/json

{
  "currentPassword": "old-password",
  "newPassword": "new-password"
}
```

---

## User Profile

### Access Profile

**Web:**
1. Log in
2. Click username → "Profile"

**API:**
```bash
GET /api/user/profile
```

**Response:**
```json
{
  "id": 42,
  "email": "user@example.com",
  "username": "johndoe",
  "api_key": "abcdefghijklmnopqrst",
  "email_verified": true,
  "is_active": true,
  "is_admin": false,
  "created_at": "2024-01-01T00:00:00Z",
  "last_login_at": "2024-01-15T10:30:00Z"
}
```

### Upload History

View your uploaded files:

**Web:**
1. Go to Profile page
2. Scroll to "Upload History" section

**API:**
```bash
GET /api/user/uploads?page=1&limit=20
```

**Response:**
```json
{
  "uploads": [
    {
      "id": 123,
      "filename": "tokyo_survey.kml",
      "uploaded_at": "2024-01-15T10:00:00Z",
      "recording_date": "2024-01-15T09:00:00Z",
      "detector": "bGeigie",
      "track_count": 1250
    }
  ],
  "total": 15,
  "page": 1,
  "limit": 20
}
```

---

## Session Management

### Session Duration

- Default: 30 days
- Sessions are automatically renewed on activity
- Expired sessions require re-login

### Session Storage

Sessions are stored in the database:
- Session ID (random token)
- User ID
- Expiration time
- Last activity timestamp

### Logout

**Web:**
1. Click username → "Logout"

**API:**
```bash
POST /api/auth/logout
```

### Clear All Sessions

Admin can clear all sessions for a user:
```bash
# Via admin panel
# Or delete from database
DELETE FROM sessions WHERE user_id = 42;
```

---

## Authentication Logging

All authentication events are logged with `AUTH:` prefix:

### Log Events

| Event | Log Format |
|-------|-----------|
| **Successful login** | `AUTH: Successful login - user_id=42 email=user@example.com method=password ip=192.168.1.100` |
| **Failed login** | `AUTH: Failed login attempt - email=user@example.com method=password reason=invalid_password ip=192.168.1.100` |
| **Registration** | `AUTH: New user registered - user_id=43 email=new@example.com username=john_doe ip=192.168.1.100` |
| **Email verification** | `AUTH: Email verified - user_id=43 email=new@example.com ip=192.168.1.100` |
| **Password change** | `AUTH: Password changed - user_id=42 email=user@example.com ip=192.168.1.100` |
| **Logout** | `AUTH: User logged out - user_id=42 email=user@example.com ip=192.168.1.100` |
| **API key regeneration** | `AUTH: API key regenerated - user_id=42 email=user@example.com admin_id=1 ip=192.168.1.100` |

### Filter Logs

```bash
# View all auth events
grep "AUTH:" /var/log/safecast.log

# View failed logins
grep "AUTH: Failed login" /var/log/safecast.log

# Monitor for brute force
grep "AUTH: Failed login.*ip=192.168.1.100" /var/log/safecast.log | wc -l

# View successful logins
grep "AUTH: Successful login" /var/log/safecast.log
```

### Security Monitoring

**Detect brute force:**
```bash
# Count failed logins per IP
grep "AUTH: Failed login" /var/log/safecast.log | \
  grep -oP 'ip=\K[0-9.]+' | sort | uniq -c | sort -rn
```

**Monitor suspicious activity:**
```bash
# Multiple failed logins in short time
grep "AUTH: Failed login" /var/log/safecast.log | \
  grep "2024-01-15T10:" | wc -l
```

---

## Admin User Management

### Create User

```bash
POST /api/admin/users/create?password=admin-password
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "password",
  "username": "newuser",
  "is_admin": false
}
```

### List Users

```bash
GET /api/admin/users?password=admin-password
```

### Update User

```bash
PUT /api/admin/users/{userId}?password=admin-password
Content-Type: application/json

{
  "username": "newname",
  "is_admin": true,
  "is_active": true
}
```

### Delete User

```bash
DELETE /api/admin/users/{userId}?password=admin-password
```

### Admin Password Reset

```bash
POST /api/admin/users/{userId}/reset-password?password=admin-password
Content-Type: application/json

{
  "new_password": "temporary-password"
}
```

---

## Database Schema

### Users Table

```sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  api_key TEXT UNIQUE NOT NULL,
  username TEXT,
  email_verified BOOLEAN DEFAULT false,
  is_active BOOLEAN DEFAULT true,
  is_admin BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  last_login_at TIMESTAMPTZ
);
```

### Sessions Table

```sql
CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  user_id BIGINT REFERENCES users(id),
  expires_at TIMESTAMPTZ,
  last_activity TIMESTAMPTZ
);
```

### Password Reset Tokens

```sql
CREATE TABLE password_reset_tokens (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id),
  token TEXT UNIQUE,
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ
);
```

### Email Verification Tokens

```sql
CREATE TABLE email_verification_tokens (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id),
  token TEXT UNIQUE,
  used BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ
);
```

---

## Migration

### Add Authentication to Existing Database

```bash
# PostgreSQL
psql -d safecast -f migrations/create_users_table.sql
psql -d safecast -f migrations/link_historical_uploads_to_users.sql

# SQLite
sqlite3 data.db < migrations/create_users_table_sqlite.sql

# DuckDB
duckdb data.duckdb < migrations/create_users_table_duckdb.sql
```

### Link Historical Uploads

```sql
-- Link uploads to users by email
UPDATE uploads u
SET internal_user_id = usr.id
FROM users usr
WHERE u.username = usr.username
  AND u.internal_user_id IS NULL;
```

### Import API Keys from CSV

```bash
go run ./cmd/tools/import-api-keys \
  -file api_keys.csv \
  -db-conn "postgres://user:pass@localhost/dbname"
```

**CSV format:**
```csv
user_id,api_key
1,abcdefghijklmnopqrst
2,uvwxyz1234567890abcd
```

---

## Best Practices

### Security

1. **Use HTTPS** in production
2. **Set strong session secret** (minimum 32 characters)
3. **Enable email verification** for new accounts
4. **Monitor authentication logs** for suspicious activity
5. **Rotate API keys** periodically
6. **Use strong passwords** (minimum 8 characters)

### User Experience

1. **Provide clear error messages** for login failures
2. **Send helpful emails** for password reset
3. **Display API key** prominently in profile
4. **Show upload history** for user engagement

### Administration

1. **Review failed login attempts** regularly
2. **Clean up inactive users** periodically
3. **Backup user database** regularly
4. **Document admin procedures** for user management

---

## Troubleshooting

### Registration Fails

**Check logs:**
```bash
grep "AUTH: New user registered" /var/log/safecast.log
```

**Check email configuration:**
```bash
# Test SMTP connection
telnet smtp.gmail.com 587
```

### Login Fails

**Verify credentials:**
```sql
-- Check if user exists
SELECT id, email, is_active FROM users WHERE email = 'user@example.com';

-- Check if email verified
SELECT email_verified FROM users WHERE id = 42;
```

### Session Expires Too Quickly

**Check session configuration:**
```bash
# Ensure session secret is set
echo $SESSION_SECRET
```

### API Key Not Working

**Verify API key:**
```sql
SELECT id, email, api_key FROM users WHERE api_key = 'your-api-key';
```

**Regenerate API key:**
```bash
# Via admin API
curl -X POST "http://localhost:8765/api/admin/users/42/regenerate-api-key?password=admin-password"
```

---

## See Also

- [Admin Panel](Admin-Panel) - User management via web interface
- [API Documentation](API-Documentation) - Authentication endpoints
- [Database Setup](Database-Setup) - Users table schema
- [Deployment](Deployment) - Production security configuration
