package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// CreateUser creates a new user in the database with a hashed password.
// Returns the new user ID.
func CreateUser(ctx context.Context, db *sql.DB, dbDriver string, user *User, password string) (int64, error) {
	// Hash the password
	passwordHash, err := HashPassword(password)
	if err != nil {
		return 0, err
	}

	now := time.Now().Unix()
	if user.CreatedAt == 0 {
		user.CreatedAt = now
	}
	if user.UpdatedAt == 0 {
		user.UpdatedAt = now
	}

	var query string
	var args []interface{}
	var id int64

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			INSERT INTO users (email, password_hash, username, email_verified, created_at, updated_at,
			                   is_active, external_id, external_source, requires_password_setup)
			VALUES ($1, $2, $3, $4, to_timestamp($5), to_timestamp($6), $7, $8, $9, $10)
			RETURNING id
		`
		args = []interface{}{
			user.Email, passwordHash, user.Username, user.EmailVerified,
			user.CreatedAt, user.UpdatedAt, user.IsActive,
			user.ExternalID, user.ExternalSource, user.RequiresPasswordSetup,
		}
		err = db.QueryRowContext(ctx, query, args...).Scan(&id)
		return id, err

	case "sqlite", "chai":
		query = `
			INSERT INTO users (email, password_hash, username, email_verified, created_at, updated_at,
			                   is_active, external_id, external_source, requires_password_setup)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		emailVerified := 0
		if user.EmailVerified {
			emailVerified = 1
		}
		isActive := 1
		if !user.IsActive {
			isActive = 0
		}
		requiresPasswordSetup := 0
		if user.RequiresPasswordSetup {
			requiresPasswordSetup = 1
		}

		args = []interface{}{
			user.Email, passwordHash, user.Username, emailVerified,
			user.CreatedAt, user.UpdatedAt, isActive,
			user.ExternalID, user.ExternalSource, requiresPasswordSetup,
		}
		result, err := db.ExecContext(ctx, query, args...)
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()

	default:
		return 0, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// GetUserByEmail retrieves a user by their email address.
func GetUserByEmail(ctx context.Context, db *sql.DB, dbDriver string, email string) (*User, error) {
	var user User
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			SELECT id, email, password_hash, username, email_verified,
			       EXTRACT(EPOCH FROM created_at)::BIGINT,
			       EXTRACT(EPOCH FROM updated_at)::BIGINT,
			       COALESCE(EXTRACT(EPOCH FROM last_login_at)::BIGINT, 0),
			       is_active, COALESCE(external_id, ''), COALESCE(external_source, ''),
			       requires_password_setup
			FROM users
			WHERE email = $1
		`
	case "sqlite", "chai":
		query = `
			SELECT id, email, password_hash, username, email_verified,
			       created_at, updated_at, COALESCE(last_login_at, 0),
			       is_active, COALESCE(external_id, ''), COALESCE(external_source, ''),
			       requires_password_setup
			FROM users
			WHERE email = ?
		`
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	var emailVerifiedInt int
	var isActiveInt int
	var requiresPasswordSetupInt int

	if dbDriver == "sqlite" || dbDriver == "chai" {
		err := db.QueryRowContext(ctx, query, email).Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Username,
			&emailVerifiedInt, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
			&isActiveInt, &user.ExternalID, &user.ExternalSource, &requiresPasswordSetupInt,
		)
		if err != nil {
			return nil, err
		}
		user.EmailVerified = emailVerifiedInt != 0
		user.IsActive = isActiveInt != 0
		user.RequiresPasswordSetup = requiresPasswordSetupInt != 0
	} else {
		err := db.QueryRowContext(ctx, query, email).Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Username,
			&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
			&user.IsActive, &user.ExternalID, &user.ExternalSource, &user.RequiresPasswordSetup,
		)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(ctx context.Context, db *sql.DB, dbDriver string, userID int64) (*User, error) {
	var user User
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			SELECT id, email, password_hash, username, email_verified,
			       EXTRACT(EPOCH FROM created_at)::BIGINT,
			       EXTRACT(EPOCH FROM updated_at)::BIGINT,
			       COALESCE(EXTRACT(EPOCH FROM last_login_at)::BIGINT, 0),
			       is_active, COALESCE(external_id, ''), COALESCE(external_source, ''),
			       requires_password_setup
			FROM users
			WHERE id = $1
		`
	case "sqlite", "chai":
		query = `
			SELECT id, email, password_hash, username, email_verified,
			       created_at, updated_at, COALESCE(last_login_at, 0),
			       is_active, COALESCE(external_id, ''), COALESCE(external_source, ''),
			       requires_password_setup
			FROM users
			WHERE id = ?
		`
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	var emailVerifiedInt int
	var isActiveInt int
	var requiresPasswordSetupInt int

	if dbDriver == "sqlite" || dbDriver == "chai" {
		err := db.QueryRowContext(ctx, query, userID).Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Username,
			&emailVerifiedInt, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
			&isActiveInt, &user.ExternalID, &user.ExternalSource, &requiresPasswordSetupInt,
		)
		if err != nil {
			return nil, err
		}
		user.EmailVerified = emailVerifiedInt != 0
		user.IsActive = isActiveInt != 0
		user.RequiresPasswordSetup = requiresPasswordSetupInt != 0
	} else {
		err := db.QueryRowContext(ctx, query, userID).Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Username,
			&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
			&user.IsActive, &user.ExternalID, &user.ExternalSource, &user.RequiresPasswordSetup,
		)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// UpdateUserPassword updates a user's password.
func UpdateUserPassword(ctx context.Context, db *sql.DB, dbDriver string, userID int64, newPassword string) error {
	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			UPDATE users
			SET password_hash = $1, updated_at = to_timestamp($2), requires_password_setup = FALSE
			WHERE id = $3
		`
		_, err = db.ExecContext(ctx, query, passwordHash, now, userID)
	case "sqlite", "chai":
		query = `
			UPDATE users
			SET password_hash = ?, updated_at = ?, requires_password_setup = 0
			WHERE id = ?
		`
		_, err = db.ExecContext(ctx, query, passwordHash, now, userID)
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	return err
}

// UpdateUserLastLogin updates the user's last login timestamp.
func UpdateUserLastLogin(ctx context.Context, db *sql.DB, dbDriver string, userID int64) error {
	now := time.Now().Unix()
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `UPDATE users SET last_login_at = to_timestamp($1) WHERE id = $2`
		_, err := db.ExecContext(ctx, query, now, userID)
		return err
	case "sqlite", "chai":
		query = `UPDATE users SET last_login_at = ? WHERE id = ?`
		_, err := db.ExecContext(ctx, query, now, userID)
		return err
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// VerifyUserEmail marks a user's email as verified.
func VerifyUserEmail(ctx context.Context, db *sql.DB, dbDriver string, userID int64) error {
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `UPDATE users SET email_verified = TRUE WHERE id = $1`
		_, err := db.ExecContext(ctx, query, userID)
		return err
	case "sqlite", "chai":
		query = `UPDATE users SET email_verified = 1 WHERE id = ?`
		_, err := db.ExecContext(ctx, query, userID)
		return err
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// CreateSession creates a new session in the database.
func CreateSession(ctx context.Context, db *sql.DB, dbDriver string, session *Session) error {
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			INSERT INTO sessions (id, user_id, created_at, expires_at, ip_address, user_agent, last_activity_at)
			VALUES ($1, $2, to_timestamp($3), to_timestamp($4), $5, $6, to_timestamp($7))
		`
		_, err := db.ExecContext(ctx, query,
			session.ID, session.UserID, session.CreatedAt, session.ExpiresAt,
			session.IPAddress, session.UserAgent, session.LastActivityAt,
		)
		return err

	case "sqlite", "chai":
		query = `
			INSERT INTO sessions (id, user_id, created_at, expires_at, ip_address, user_agent, last_activity_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err := db.ExecContext(ctx, query,
			session.ID, session.UserID, session.CreatedAt, session.ExpiresAt,
			session.IPAddress, session.UserAgent, session.LastActivityAt,
		)
		return err

	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// GetSession retrieves a session by its ID.
func GetSession(ctx context.Context, db *sql.DB, dbDriver string, sessionID string) (*Session, error) {
	var session Session
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			SELECT id, user_id,
			       EXTRACT(EPOCH FROM created_at)::BIGINT,
			       EXTRACT(EPOCH FROM expires_at)::BIGINT,
			       COALESCE(ip_address, ''), COALESCE(user_agent, ''),
			       EXTRACT(EPOCH FROM last_activity_at)::BIGINT
			FROM sessions
			WHERE id = $1
		`
	case "sqlite", "chai":
		query = `
			SELECT id, user_id, created_at, expires_at,
			       COALESCE(ip_address, ''), COALESCE(user_agent, ''),
			       last_activity_at
			FROM sessions
			WHERE id = ?
		`
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	err := db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt,
		&session.IPAddress, &session.UserAgent, &session.LastActivityAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// UpdateSessionActivity updates the last activity timestamp for a session.
func UpdateSessionActivity(ctx context.Context, db *sql.DB, dbDriver string, sessionID string) error {
	now := time.Now().Unix()
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `UPDATE sessions SET last_activity_at = to_timestamp($1) WHERE id = $2`
		_, err := db.ExecContext(ctx, query, now, sessionID)
		return err
	case "sqlite", "chai":
		query = `UPDATE sessions SET last_activity_at = ? WHERE id = ?`
		_, err := db.ExecContext(ctx, query, now, sessionID)
		return err
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// DeleteSession removes a session from the database.
func DeleteSession(ctx context.Context, db *sql.DB, dbDriver string, sessionID string) error {
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `DELETE FROM sessions WHERE id = $1`
	case "sqlite", "chai":
		query = `DELETE FROM sessions WHERE id = ?`
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	_, err := db.ExecContext(ctx, query, sessionID)
	return err
}

// DeleteExpiredSessions removes all expired sessions from the database.
func DeleteExpiredSessions(ctx context.Context, db *sql.DB, dbDriver string) error {
	now := time.Now().Unix()
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `DELETE FROM sessions WHERE expires_at < to_timestamp($1)`
		_, err := db.ExecContext(ctx, query, now)
		return err
	case "sqlite", "chai":
		query = `DELETE FROM sessions WHERE expires_at < ?`
		_, err := db.ExecContext(ctx, query, now)
		return err
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// CreatePasswordResetToken creates a new password reset token.
func CreatePasswordResetToken(ctx context.Context, db *sql.DB, dbDriver string, token *PasswordResetToken) error {
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			INSERT INTO password_reset_tokens (user_id, token, created_at, expires_at, used, ip_address)
			VALUES ($1, $2, to_timestamp($3), to_timestamp($4), $5, $6)
		`
		_, err := db.ExecContext(ctx, query,
			token.UserID, token.Token, token.CreatedAt, token.ExpiresAt, token.Used, token.IPAddress,
		)
		return err

	case "sqlite", "chai":
		used := 0
		if token.Used {
			used = 1
		}
		query = `
			INSERT INTO password_reset_tokens (user_id, token, created_at, expires_at, used, ip_address)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err := db.ExecContext(ctx, query,
			token.UserID, token.Token, token.CreatedAt, token.ExpiresAt, used, token.IPAddress,
		)
		return err

	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// GetPasswordResetToken retrieves a password reset token by its token string.
func GetPasswordResetToken(ctx context.Context, db *sql.DB, dbDriver string, tokenStr string) (*PasswordResetToken, error) {
	var token PasswordResetToken
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			SELECT id, user_id, token,
			       EXTRACT(EPOCH FROM created_at)::BIGINT,
			       EXTRACT(EPOCH FROM expires_at)::BIGINT,
			       used,
			       COALESCE(EXTRACT(EPOCH FROM used_at)::BIGINT, 0),
			       COALESCE(ip_address, '')
			FROM password_reset_tokens
			WHERE token = $1
		`
	case "sqlite", "chai":
		query = `
			SELECT id, user_id, token, created_at, expires_at, used,
			       COALESCE(used_at, 0), COALESCE(ip_address, '')
			FROM password_reset_tokens
			WHERE token = ?
		`
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	if dbDriver == "sqlite" || dbDriver == "chai" {
		var usedInt int
		err := db.QueryRowContext(ctx, query, tokenStr).Scan(
			&token.ID, &token.UserID, &token.Token, &token.CreatedAt, &token.ExpiresAt,
			&usedInt, &token.UsedAt, &token.IPAddress,
		)
		if err != nil {
			return nil, err
		}
		token.Used = usedInt != 0
	} else {
		err := db.QueryRowContext(ctx, query, tokenStr).Scan(
			&token.ID, &token.UserID, &token.Token, &token.CreatedAt, &token.ExpiresAt,
			&token.Used, &token.UsedAt, &token.IPAddress,
		)
		if err != nil {
			return nil, err
		}
	}

	return &token, nil
}

// MarkPasswordResetTokenUsed marks a password reset token as used.
func MarkPasswordResetTokenUsed(ctx context.Context, db *sql.DB, dbDriver string, tokenID int64) error {
	now := time.Now().Unix()
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `UPDATE password_reset_tokens SET used = TRUE, used_at = to_timestamp($1) WHERE id = $2`
		_, err := db.ExecContext(ctx, query, now, tokenID)
		return err
	case "sqlite", "chai":
		query = `UPDATE password_reset_tokens SET used = 1, used_at = ? WHERE id = ?`
		_, err := db.ExecContext(ctx, query, now, tokenID)
		return err
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// CreateEmailVerificationToken creates a new email verification token.
func CreateEmailVerificationToken(ctx context.Context, db *sql.DB, dbDriver string, token *EmailVerificationToken) error {
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			INSERT INTO email_verification_tokens (user_id, token, created_at, expires_at, used)
			VALUES ($1, $2, to_timestamp($3), to_timestamp($4), $5)
		`
		_, err := db.ExecContext(ctx, query,
			token.UserID, token.Token, token.CreatedAt, token.ExpiresAt, token.Used,
		)
		return err

	case "sqlite", "chai":
		used := 0
		if token.Used {
			used = 1
		}
		query = `
			INSERT INTO email_verification_tokens (user_id, token, created_at, expires_at, used)
			VALUES (?, ?, ?, ?, ?)
		`
		_, err := db.ExecContext(ctx, query,
			token.UserID, token.Token, token.CreatedAt, token.ExpiresAt, used,
		)
		return err

	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// GetEmailVerificationToken retrieves an email verification token by its token string.
func GetEmailVerificationToken(ctx context.Context, db *sql.DB, dbDriver string, tokenStr string) (*EmailVerificationToken, error) {
	var token EmailVerificationToken
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			SELECT id, user_id, token,
			       EXTRACT(EPOCH FROM created_at)::BIGINT,
			       EXTRACT(EPOCH FROM expires_at)::BIGINT,
			       used,
			       COALESCE(EXTRACT(EPOCH FROM used_at)::BIGINT, 0)
			FROM email_verification_tokens
			WHERE token = $1
		`
	case "sqlite", "chai":
		query = `
			SELECT id, user_id, token, created_at, expires_at, used, COALESCE(used_at, 0)
			FROM email_verification_tokens
			WHERE token = ?
		`
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	if dbDriver == "sqlite" || dbDriver == "chai" {
		var usedInt int
		err := db.QueryRowContext(ctx, query, tokenStr).Scan(
			&token.ID, &token.UserID, &token.Token, &token.CreatedAt, &token.ExpiresAt,
			&usedInt, &token.UsedAt,
		)
		if err != nil {
			return nil, err
		}
		token.Used = usedInt != 0
	} else {
		err := db.QueryRowContext(ctx, query, tokenStr).Scan(
			&token.ID, &token.UserID, &token.Token, &token.CreatedAt, &token.ExpiresAt,
			&token.Used, &token.UsedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return &token, nil
}

// MarkEmailVerificationTokenUsed marks an email verification token as used.
func MarkEmailVerificationTokenUsed(ctx context.Context, db *sql.DB, dbDriver string, tokenID int64) error {
	now := time.Now().Unix()
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `UPDATE email_verification_tokens SET used = TRUE, used_at = to_timestamp($1) WHERE id = $2`
		_, err := db.ExecContext(ctx, query, now, tokenID)
		return err
	case "sqlite", "chai":
		query = `UPDATE email_verification_tokens SET used = 1, used_at = ? WHERE id = ?`
		_, err := db.ExecContext(ctx, query, now, tokenID)
		return err
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// LinkUserIDToUploads updates existing uploads with a user_id to link them to a newly created user.
// This is used during user migration to link historical uploads to the migrated user account.
func LinkUserIDToUploads(ctx context.Context, db *sql.DB, dbDriver string, externalUserID string, newUserID string) error {
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `UPDATE uploads SET user_id = $1 WHERE user_id = $2`
		_, err := db.ExecContext(ctx, query, newUserID, externalUserID)
		return err
	case "sqlite", "chai":
		query = `UPDATE uploads SET user_id = ? WHERE user_id = ?`
		_, err := db.ExecContext(ctx, query, newUserID, externalUserID)
		return err
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}

// GetAllUsers retrieves all users from the database.
func GetAllUsers(ctx context.Context, db *sql.DB, dbDriver string) ([]*User, error) {
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			SELECT id, email, password_hash, username, email_verified,
			       EXTRACT(EPOCH FROM created_at)::BIGINT as created_at,
			       EXTRACT(EPOCH FROM updated_at)::BIGINT as updated_at,
			       COALESCE(EXTRACT(EPOCH FROM last_login_at)::BIGINT, 0) as last_login_at,
			       is_active, COALESCE(external_id, ''), COALESCE(external_source, ''),
			       requires_password_setup
			FROM users
			ORDER BY created_at DESC
		`
	case "sqlite", "chai":
		query = `
			SELECT id, email, password_hash, username, email_verified, created_at, updated_at,
			       COALESCE(last_login_at, 0) as last_login_at,
			       is_active, COALESCE(external_id, ''), COALESCE(external_source, ''),
			       requires_password_setup
			FROM users
			ORDER BY created_at DESC
		`
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Username,
			&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
			&user.IsActive, &user.ExternalID, &user.ExternalSource,
			&user.RequiresPasswordSetup,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// UpdateUser updates an existing user in the database.
func UpdateUser(ctx context.Context, db *sql.DB, dbDriver string, user *User) error {
	user.UpdatedAt = time.Now().Unix()

	var query string
	var args []interface{}

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			UPDATE users
			SET email = $1, username = $2, email_verified = $3, updated_at = to_timestamp($4),
			    is_active = $5, external_id = $6, external_source = $7, requires_password_setup = $8
			WHERE id = $9
		`
		args = []interface{}{
			user.Email, user.Username, user.EmailVerified, user.UpdatedAt,
			user.IsActive, user.ExternalID, user.ExternalSource, user.RequiresPasswordSetup,
			user.ID,
		}
	case "sqlite", "chai":
		query = `
			UPDATE users
			SET email = ?, username = ?, email_verified = ?, updated_at = ?,
			    is_active = ?, external_id = ?, external_source = ?, requires_password_setup = ?
			WHERE id = ?
		`
		args = []interface{}{
			user.Email, user.Username, user.EmailVerified, user.UpdatedAt,
			user.IsActive, user.ExternalID, user.ExternalSource, user.RequiresPasswordSetup,
			user.ID,
		}
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteUser deletes a user from the database.
// This will cascade delete all associated sessions, tokens, etc.
func DeleteUser(ctx context.Context, db *sql.DB, dbDriver string, userID int64) error {
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `DELETE FROM users WHERE id = $1`
		_, err := db.ExecContext(ctx, query, userID)
		return err
	case "sqlite", "chai":
		query = `DELETE FROM users WHERE id = ?`
		_, err := db.ExecContext(ctx, query, userID)
		return err
	default:
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}
}
