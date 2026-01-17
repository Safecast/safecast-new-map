package auth

import (
	"context"
	"database/sql"
	"net/http"
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
	SendPasswordSetupEmail(to, username, setupURL string) error
	SendPasswordResetEmail(to, resetURL string) error
	SendPasswordChangedEmail(to string) error
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

// setSessionCookie sets the session cookie in the response.
func (m *Manager) setSessionCookie(w http.ResponseWriter, sessionID string, expiresAt int64) {
	cookie := &http.Cookie{
		Name:     m.SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Unix(expiresAt, 0),
		HttpOnly: true, // Prevent JavaScript access (XSS protection)
		Secure:   true, // Only send over HTTPS (should be true in production)
		SameSite: http.SameSiteLaxMode, // CSRF protection
	}
	http.SetCookie(w, cookie)
}

// clearSessionCookie clears the session cookie in the response.
func (m *Manager) clearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     m.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}
