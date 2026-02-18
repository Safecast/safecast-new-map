package auth

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const userContextKey contextKey = "auth_user"

// Manager holds the authentication system configuration and dependencies.
type Manager struct {
	DB                *sql.DB
	DBDriver          string
	SessionCookieName string
	AllowRegistration bool
	EmailSender       EmailSender
	BaseURL           string
}

// EmailSender is an interface for sending emails.
// This allows for dependency injection and easier testing.
type EmailSender interface {
	SendWelcomeEmail(to, username, verificationURL string) error
	SendWelcomeEmailWithAPIKey(to, username, verificationURL, apiKey string) error
	SendPasswordSetupEmail(to, username, setupURL string) error
	SendPasswordResetEmail(to, resetURL string) error
	SendPasswordChangedEmail(to string) error
	SendAPIKeyRegeneratedEmail(to, username, newAPIKey string) error
}

// WithAuthContext adds a user to the request context.
func WithAuthContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUserFromContext retrieves the authenticated user from the request context.
// Returns the user and true if authenticated, nil and false otherwise.
func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextKey).(*User)
	return user, ok
}

// getUserFromSession attempts to retrieve and validate a user session from cookies.
// Returns the user if authenticated, nil otherwise.
func (m *Manager) getUserFromSession(r *http.Request) (*User, error) {
	cookie, err := r.Cookie(m.SessionCookieName)
	if err != nil {
		// No cookie found, not authenticated
		return nil, nil
	}

	sessionID := cookie.Value
	if sessionID == "" {
		return nil, nil
	}

	// Get session from database
	session, err := GetSession(r.Context(), m.DB, m.DBDriver, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Session not found in database
			return nil, nil
		}
		return nil, err
	}

	// Check if session is expired
	now := time.Now().Unix()
	if session.ExpiresAt < now {
		// Session expired, delete it
		_ = DeleteSession(r.Context(), m.DB, m.DBDriver, sessionID)
		return nil, nil
	}

	// Get user from database
	user, err := GetUserByID(r.Context(), m.DB, m.DBDriver, session.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, nil
	}

	// Update session activity (don't block on this)
	go UpdateSessionActivity(context.Background(), m.DB, m.DBDriver, sessionID)

	return user, nil
}

// RequireAuth is middleware that requires authentication.
// If the user is not authenticated, it returns 401 Unauthorized.
func (m *Manager) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := m.getUserFromSession(r)
		if err != nil {
			http.Error(w, "Authentication error", http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "Unauthorized - Please login", http.StatusUnauthorized)
			return
		}

		// Add user to request context
		ctx := WithAuthContext(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// OptionalAuth is middleware that attempts to authenticate the user but doesn't require it.
// The handler can check for the user in the context to see if they're authenticated.
func (m *Manager) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := m.getUserFromSession(r)
		if err != nil {
			// Log error but don't block the request
			// log.Printf("Optional auth error: %v", err)
		}

		if user != nil {
			// Add user to request context
			ctx := WithAuthContext(r.Context(), user)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	}
}

// RequireAdmin is middleware that requires admin authentication.
// If the user is not authenticated or is not an admin, it returns an error.
func (m *Manager) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := m.getUserFromSession(r)
		if err != nil {
			http.Error(w, "Authentication error", http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "Unauthorized - Please login", http.StatusUnauthorized)
			return
		}

		if !user.IsAdmin {
			http.Error(w, "Forbidden - Admin access required", http.StatusForbidden)
			return
		}

		// Add user to request context
		ctx := WithAuthContext(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserByAPIKey retrieves a user by their API key.
func GetUserByAPIKey(ctx context.Context, db *sql.DB, dbDriver string, apiKey string) (*User, error) {
	var user User
	var query string

	switch dbDriver {
	case "pgx", "duckdb":
		query = `
			SELECT id, email, password_hash, username, email_verified,
			       EXTRACT(EPOCH FROM created_at)::BIGINT,
			       EXTRACT(EPOCH FROM updated_at)::BIGINT,
			       COALESCE(EXTRACT(EPOCH FROM last_login_at)::BIGINT, 0),
			       is_active, is_admin, COALESCE(external_id, ''), COALESCE(external_source, ''),
			       requires_password_setup, COALESCE(api_key, '')
			FROM users
			WHERE api_key = $1 AND is_active = TRUE
		`
	case "sqlite", "chai":
		query = `
			SELECT id, email, password_hash, username, email_verified,
			       created_at, updated_at, COALESCE(last_login_at, 0),
			       is_active, COALESCE(is_admin, 0), COALESCE(external_id, ''), COALESCE(external_source, ''),
			       requires_password_setup, COALESCE(api_key, '')
			FROM users
			WHERE api_key = ? AND is_active = 1
		`
	default:
		return nil, sql.ErrNoRows
	}

	var emailVerifiedInt int
	var isActiveInt int
	var isAdminInt int
	var requiresPasswordSetupInt int

	if dbDriver == "sqlite" || dbDriver == "chai" {
		err := db.QueryRowContext(ctx, query, apiKey).Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Username,
			&emailVerifiedInt, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
			&isActiveInt, &isAdminInt, &user.ExternalID, &user.ExternalSource, &requiresPasswordSetupInt,
			&user.APIKey,
		)
		if err != nil {
			return nil, err
		}
		user.EmailVerified = emailVerifiedInt != 0
		user.IsActive = isActiveInt != 0
		user.IsAdmin = isAdminInt != 0
		user.RequiresPasswordSetup = requiresPasswordSetupInt != 0
	} else {
		err := db.QueryRowContext(ctx, query, apiKey).Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Username,
			&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
			&user.IsActive, &user.IsAdmin, &user.ExternalID, &user.ExternalSource, &user.RequiresPasswordSetup,
			&user.APIKey,
		)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// RequireAuthOrAPIKey is middleware that requires either session authentication or API key.
// This is useful for endpoints that need to support both browser and programmatic access.
func (m *Manager) RequireAuthOrAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// First try session authentication
		user, err := m.getUserFromSession(r)
		if err == nil && user != nil {
			ctx := WithAuthContext(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Then try API key authentication
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			// Also check query parameter for convenience
			apiKey = r.URL.Query().Get("api_key")
		}

		if apiKey != "" {
			user, err = GetUserByAPIKey(r.Context(), m.DB, m.DBDriver, apiKey)
			if err == nil && user != nil {
				ctx := WithAuthContext(r.Context(), user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		http.Error(w, "Unauthorized - Please login or provide API key", http.StatusUnauthorized)
	}
}

// setSessionCookie sets the session cookie in the response.
func (m *Manager) setSessionCookie(w http.ResponseWriter, sessionID string, expiresAt int64) {
	secure := !strings.HasPrefix(m.BaseURL, "http://")
	cookie := &http.Cookie{
		Name:     m.SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Unix(expiresAt, 0),
		HttpOnly: true,                    // Prevent JavaScript access (XSS protection)
		Secure:   secure,                  // Only send over HTTPS in production
		SameSite: http.SameSiteLaxMode,    // CSRF protection
	}
	http.SetCookie(w, cookie)
}

// clearSessionCookie clears the session cookie in the response.
func (m *Manager) clearSessionCookie(w http.ResponseWriter) {
	secure := !strings.HasPrefix(m.BaseURL, "http://")
	cookie := &http.Cookie{
		Name:     m.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}
