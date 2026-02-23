package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	// SessionDuration is the default duration for user sessions (30 days)
	SessionDuration = 30 * 24 * time.Hour

	// PasswordResetTokenDuration is how long password reset tokens are valid (1 hour)
	PasswordResetTokenDuration = 1 * time.Hour

	// EmailVerificationTokenDuration is how long email verification tokens are valid (24 hours)
	EmailVerificationTokenDuration = 24 * time.Hour
)

// RegisterHandler handles user registration requests.
// POST /api/auth/register
func (m *Manager) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !m.AllowRegistration {
		writeJSONError(w, "Registration is currently disabled", http.StatusForbidden)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate inputs
	req.Email = NormalizeEmail(req.Email)
	if err := ValidateEmail(req.Email); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Username = SanitizeInput(req.Username)
	if err := ValidateUsername(req.Username); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ValidatePassword(req.Password); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	existingUser, err := GetUserByEmail(r.Context(), m.DB, m.DBDriver, req.Email)
	if err != nil && err != sql.ErrNoRows {
		writeJSONError(w, "Database error", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		writeJSONError(w, "Email address already registered", http.StatusConflict)
		return
	}

	// Create user
	user := &User{
		Email:         req.Email,
		Username:      req.Username,
		EmailVerified: false,
		IsActive:      true,
	}

	clientIP := getClientIP(r)
	userID, err := CreateUser(r.Context(), m.DB, m.DBDriver, user, req.Password)
	if err != nil {
		log.Printf("AUTH: Failed registration - email=%s username=%s reason=create_failed ip=%s error=%v", req.Email, req.Username, clientIP, err)
		writeJSONError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Printf("AUTH: New user registered - user_id=%d email=%s username=%s ip=%s", userID, req.Email, req.Username, clientIP)

	// Create email verification token
	token, err := GenerateToken()
	if err != nil {
		writeJSONError(w, "Failed to generate verification token", http.StatusInternalServerError)
		return
	}

	now := time.Now().Unix()
	verificationToken := &EmailVerificationToken{
		UserID:    userID,
		Token:     token,
		CreatedAt: now,
		ExpiresAt: now + int64(EmailVerificationTokenDuration.Seconds()),
		Used:      false,
	}

	if err := CreateEmailVerificationToken(r.Context(), m.DB, m.DBDriver, verificationToken); err != nil {
		writeJSONError(w, "Failed to create verification token", http.StatusInternalServerError)
		return
	}

	// Get the newly created user to retrieve the API key
	newUser, err := GetUserByID(r.Context(), m.DB, m.DBDriver, userID)
	if err != nil {
		log.Printf("WARNING: Failed to retrieve newly created user for API key: %v", err)
	}

	// Send verification email with API key
	if m.EmailSender != nil {
		verificationURL := fmt.Sprintf("%s/api/auth/verify-email?token=%s", m.BaseURL, token)

		// Include API key in the welcome email if user was retrieved successfully
		if newUser != nil && newUser.APIKey != "" {
			if err := m.EmailSender.SendWelcomeEmailWithAPIKey(req.Email, req.Username, verificationURL, newUser.APIKey); err != nil {
				log.Printf("ERROR: Failed to send welcome email with API key to %s: %v", req.Email, err)
				// Don't fail registration - user can request a new verification email later
			}
		} else {
			// Fallback to regular welcome email without API key
			if err := m.EmailSender.SendWelcomeEmail(req.Email, req.Username, verificationURL); err != nil {
				log.Printf("ERROR: Failed to send welcome email to %s: %v", req.Email, err)
			}
		}
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
		"message": "Registration successful. Please check your email to verify your account.",
		"userId":  userID,
	}, http.StatusCreated)
}

// LoginHandler handles user login requests.
// POST /api/auth/login
// Accepts either password or API key for authentication
func (m *Manager) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		APIKey   string `json:"api_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.Email = NormalizeEmail(req.Email)
	clientIP := getClientIP(r)

	var user *User
	var err error
	var authMethod string

	// If API key is provided, authenticate with API key
	if req.APIKey != "" {
		authMethod = "api_key"
		user, err = GetUserByAPIKey(r.Context(), m.DB, m.DBDriver, req.APIKey)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("AUTH: Failed login attempt - email=%s method=api_key reason=invalid_key ip=%s", req.Email, clientIP)
				writeJSONError(w, "Invalid API key", http.StatusUnauthorized)
				return
			}
			log.Printf("AUTH: Failed login attempt - email=%s method=api_key reason=database_error ip=%s error=%v", req.Email, clientIP, err)
			writeJSONError(w, "Database error", http.StatusInternalServerError)
			return
		}
		// Verify the email matches the API key's user
		if user.Email != req.Email {
			log.Printf("AUTH: Failed login attempt - email=%s method=api_key reason=email_mismatch ip=%s", req.Email, clientIP)
			writeJSONError(w, "Invalid email or API key", http.StatusUnauthorized)
			return
		}
	} else if req.Password != "" {
		authMethod = "password"
		// Otherwise, authenticate with password
		user, err = GetUserByEmail(r.Context(), m.DB, m.DBDriver, req.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("AUTH: Failed login attempt - email=%s method=password reason=user_not_found ip=%s", req.Email, clientIP)
				writeJSONError(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}
			log.Printf("AUTH: Failed login attempt - email=%s method=password reason=database_error ip=%s error=%v", req.Email, clientIP, err)
			writeJSONError(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Verify password
		if !VerifyPassword(user.PasswordHash, req.Password) {
			log.Printf("AUTH: Failed login attempt - email=%s method=password reason=invalid_password ip=%s", req.Email, clientIP)
			writeJSONError(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
	} else {
		log.Printf("AUTH: Failed login attempt - email=%s method=none reason=no_credentials ip=%s", req.Email, clientIP)
		writeJSONError(w, "Either password or API key is required", http.StatusBadRequest)
		return
	}

	// Check if user is active
	if !user.IsActive {
		log.Printf("AUTH: Failed login attempt - email=%s user_id=%d method=%s reason=account_disabled ip=%s", user.Email, user.ID, authMethod, clientIP)
		writeJSONError(w, "Account is disabled", http.StatusForbidden)
		return
	}

	// Create session
	sessionID, err := GenerateSessionID()
	if err != nil {
		writeJSONError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	now := time.Now().Unix()
	expiresAt := now + int64(SessionDuration.Seconds())

	session := &Session{
		ID:             sessionID,
		UserID:         user.ID,
		CreatedAt:      now,
		ExpiresAt:      expiresAt,
		IPAddress:      getClientIP(r),
		UserAgent:      r.UserAgent(),
		LastActivityAt: now,
	}

	if err := CreateSession(r.Context(), m.DB, m.DBDriver, session); err != nil {
		writeJSONError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Update last login timestamp
	_ = UpdateUserLastLogin(r.Context(), m.DB, m.DBDriver, user.ID)

	// Log successful login
	log.Printf("AUTH: Successful login - user_id=%d email=%s method=%s ip=%s user_agent=%s",
		user.ID, user.Email, authMethod, clientIP, r.UserAgent())

	// Set session cookie
	m.setSessionCookie(w, sessionID, expiresAt)

	// Return user info (without password hash)
	user.PasswordHash = ""
	writeJSON(w, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"user":    user,
	}, http.StatusOK)
}

// LogoutHandler handles user logout requests.
// POST /api/auth/logout
func (m *Manager) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientIP := getClientIP(r)
	var userEmail string
	var userID int64

	// Get session ID from cookie and try to get user info for logging
	cookie, err := r.Cookie(m.SessionCookieName)
	if err == nil && cookie.Value != "" {
		// Try to get session info for logging before deleting
		if session, err := GetSession(r.Context(), m.DB, m.DBDriver, cookie.Value); err == nil {
			if user, err := GetUserByID(r.Context(), m.DB, m.DBDriver, session.UserID); err == nil {
				userEmail = user.Email
				userID = user.ID
			}
		}
		// Delete session from database
		_ = DeleteSession(r.Context(), m.DB, m.DBDriver, cookie.Value)
	}

	// Log logout
	if userEmail != "" {
		log.Printf("AUTH: User logout - user_id=%d email=%s ip=%s", userID, userEmail, clientIP)
	} else {
		log.Printf("AUTH: Logout attempt (no active session) - ip=%s", clientIP)
	}

	// Clear session cookie
	m.clearSessionCookie(w)

	writeJSON(w, map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	}, http.StatusOK)
}

// ForgotPasswordHandler handles password reset requests.
// POST /api/auth/forgot-password
func (m *Manager) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.Email = NormalizeEmail(req.Email)

	// Get user from database
	user, err := GetUserByEmail(r.Context(), m.DB, m.DBDriver, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't reveal if email exists or not (security)
			writeJSON(w, map[string]interface{}{
				"success": true,
				"message": "If that email address is registered, a password reset link has been sent.",
			}, http.StatusOK)
			return
		}
		writeJSONError(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Create password reset token
	token, err := GenerateToken()
	if err != nil {
		writeJSONError(w, "Failed to generate reset token", http.StatusInternalServerError)
		return
	}

	now := time.Now().Unix()
	resetToken := &PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		CreatedAt: now,
		ExpiresAt: now + int64(PasswordResetTokenDuration.Seconds()),
		Used:      false,
		IPAddress: getClientIP(r),
	}

	if err := CreatePasswordResetToken(r.Context(), m.DB, m.DBDriver, resetToken); err != nil {
		writeJSONError(w, "Failed to create reset token", http.StatusInternalServerError)
		return
	}

	// Send password reset email
	if m.EmailSender != nil {
		resetURL := fmt.Sprintf("%s/reset-password?token=%s", m.BaseURL, token)
		if err := m.EmailSender.SendPasswordResetEmail(req.Email, resetURL); err != nil {
			log.Printf("ERROR: Failed to send password reset email to %s: %v", req.Email, err)
		}
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
		"message": "If that email address is registered, a password reset link has been sent.",
	}, http.StatusOK)
}

// ResetPasswordHandler handles password reset with token.
// POST /api/auth/reset-password
func (m *Manager) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := ValidatePassword(req.Password); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get reset token from database
	token, err := GetPasswordResetToken(r.Context(), m.DB, m.DBDriver, req.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, "Invalid or expired reset token", http.StatusBadRequest)
			return
		}
		writeJSONError(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check if token is expired
	now := time.Now().Unix()
	if token.ExpiresAt < now {
		writeJSONError(w, "Reset token has expired", http.StatusBadRequest)
		return
	}

	// Check if token has already been used
	if token.Used {
		writeJSONError(w, "Reset token has already been used", http.StatusBadRequest)
		return
	}

	// Update user password
	if err := UpdateUserPassword(r.Context(), m.DB, m.DBDriver, token.UserID, req.Password); err != nil {
		writeJSONError(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Mark token as used
	_ = MarkPasswordResetTokenUsed(r.Context(), m.DB, m.DBDriver, token.ID)

	// Get user email for notification
	user, err := GetUserByID(r.Context(), m.DB, m.DBDriver, token.UserID)
	if err == nil && m.EmailSender != nil {
		_ = m.EmailSender.SendPasswordChangedEmail(user.Email)
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
		"message": "Password reset successful",
	}, http.StatusOK)
}

// VerifyEmailHandler handles email verification with token.
// GET /api/auth/verify-email?token=xxx
func (m *Manager) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		writeJSONError(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Get verification token from database
	token, err := GetEmailVerificationToken(r.Context(), m.DB, m.DBDriver, tokenStr)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, "Invalid or expired verification token", http.StatusBadRequest)
			return
		}
		writeJSONError(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check if token is expired
	now := time.Now().Unix()
	if token.ExpiresAt < now {
		writeJSONError(w, "Verification token has expired", http.StatusBadRequest)
		return
	}

	// Check if token has already been used
	if token.Used {
		writeJSONError(w, "Email already verified", http.StatusOK)
		return
	}

	// Get user info for logging
	user, err := GetUserByID(r.Context(), m.DB, m.DBDriver, token.UserID)
	if err != nil {
		writeJSONError(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	// Verify user email
	if err := VerifyUserEmail(r.Context(), m.DB, m.DBDriver, token.UserID); err != nil {
		writeJSONError(w, "Failed to verify email", http.StatusInternalServerError)
		return
	}

	// Mark token as used
	_ = MarkEmailVerificationTokenUsed(r.Context(), m.DB, m.DBDriver, token.ID)

	log.Printf("AUTH: Email verified - user_id=%d email=%s ip=%s", user.ID, user.Email, getClientIP(r))

	// Redirect to map page with success message
	http.Redirect(w, r, "/?emailVerified=true", http.StatusSeeOther)
}

// ProfileHandler returns the current user's profile.
// GET /api/user/profile
func (m *Manager) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent CloudFront from caching user-specific data
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Don't expose password hash
	user.PasswordHash = ""

	writeJSON(w, user, http.StatusOK)
}

// ChangePasswordHandler allows a logged-in user to change their password.
// POST /api/user/change-password
func (m *Manager) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		writeJSONError(w, "Current password and new password are required", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 8 {
		writeJSONError(w, "New password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	// Verify current password
	if !VerifyPassword(user.PasswordHash, req.CurrentPassword) {
		writeJSONError(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// Update password
	if err := UpdateUserPassword(r.Context(), m.DB, m.DBDriver, user.ID, req.NewPassword); err != nil {
		writeJSONError(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	log.Printf("AUTH: Password changed - user_id=%d email=%s ip=%s", user.ID, user.Email, getClientIP(r))

	writeJSON(w, map[string]interface{}{
		"success": true,
		"message": "Password changed successfully",
	}, http.StatusOK)
}

// Helper functions

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, map[string]interface{}{
		"error": message,
	}, status)
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}
