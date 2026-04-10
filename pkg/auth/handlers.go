package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"safecast-new-map/pkg/httpresp"
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
//
// @Summary     Register a new user
// @Description Creates a new user account and sends an email verification token when email is configured.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body object true "Registration request"
// @Success     201 {object} map[string]interface{} "Registration created"
// @Failure     400 {object} map[string]interface{} "Invalid request"
// @Failure     403 {object} map[string]interface{} "Registration disabled"
// @Failure     409 {object} map[string]interface{} "Email already exists"
// @Failure     500 {object} map[string]interface{} "Server error"
// @Router      /api/auth/register [post]
func (m *Manager) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !httpresp.RequireMethodJSON(w, r, http.MethodPost) {
		return
	}

	if !m.AllowRegistration {
		httpresp.WriteForbidden(w, "Registration is currently disabled")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Invalid request body")
		return
	}

	// Validate inputs
	req.Email = NormalizeEmail(req.Email)
	if err := ValidateEmail(req.Email); err != nil {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, err.Error())
		return
	}

	req.Username = SanitizeInput(req.Username)
	if err := ValidateUsername(req.Username); err != nil {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, err.Error())
		return
	}

	if err := ValidatePassword(req.Password); err != nil {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, err.Error())
		return
	}

	// Check if user already exists
	existingUser, err := GetUserByEmail(r.Context(), m.DB, m.DBDriver, req.Email)
	if err != nil && err != sql.ErrNoRows {
		httpresp.WriteInternalError(w, "Database error")
		return
	}
	if existingUser != nil {
		httpresp.WriteError(w, http.StatusConflict, httpresp.CodeConflict, "Email address already registered")
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
		httpresp.WriteInternalError(w, "Failed to create user")
		return
	}

	log.Printf("AUTH: New user registered - user_id=%d email=%s username=%s ip=%s", userID, req.Email, req.Username, clientIP)

	// Create email verification token
	token, err := GenerateToken()
	if err != nil {
		httpresp.WriteInternalError(w, "Failed to generate verification token")
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
		httpresp.WriteInternalError(w, "Failed to create verification token")
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

	httpresp.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Registration successful. Please check your email to verify your account.",
		"userId":  userID,
	})
}

// LoginHandler handles user login requests.
//
// @Summary     Login user
// @Description Authenticates by email+password or email+API key and creates a session cookie.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body object true "Login request"
// @Success     200 {object} map[string]interface{} "Login success"
// @Failure     400 {object} map[string]interface{} "Invalid request"
// @Failure     401 {object} map[string]interface{} "Unauthorized"
// @Failure     403 {object} map[string]interface{} "Account disabled"
// @Failure     500 {object} map[string]interface{} "Server error"
// @Router      /api/auth/login [post]
func (m *Manager) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !httpresp.RequireMethodJSON(w, r, http.MethodPost) {
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		APIKey   string `json:"api_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Invalid request body")
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
				httpresp.WriteUnauthorized(w, "Invalid API key")
				return
			}
			log.Printf("AUTH: Failed login attempt - email=%s method=api_key reason=database_error ip=%s error=%v", req.Email, clientIP, err)
			httpresp.WriteInternalError(w, "Database error")
			return
		}
		// Verify the email matches the API key's user
		if user.Email != req.Email {
			log.Printf("AUTH: Failed login attempt - email=%s method=api_key reason=email_mismatch ip=%s", req.Email, clientIP)
			httpresp.WriteUnauthorized(w, "Invalid email or API key")
			return
		}
	} else if req.Password != "" {
		authMethod = "password"
		// Otherwise, authenticate with password
		user, err = GetUserByEmail(r.Context(), m.DB, m.DBDriver, req.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("AUTH: Failed login attempt - email=%s method=password reason=user_not_found ip=%s", req.Email, clientIP)
				httpresp.WriteUnauthorized(w, "Invalid email or password")
				return
			}
			log.Printf("AUTH: Failed login attempt - email=%s method=password reason=database_error ip=%s error=%v", req.Email, clientIP, err)
			httpresp.WriteInternalError(w, "Database error")
			return
		}

		// Verify password
		if !VerifyPassword(user.PasswordHash, req.Password) {
			log.Printf("AUTH: Failed login attempt - email=%s method=password reason=invalid_password ip=%s", req.Email, clientIP)
			httpresp.WriteUnauthorized(w, "Invalid email or password")
			return
		}
	} else {
		log.Printf("AUTH: Failed login attempt - email=%s method=none reason=no_credentials ip=%s", req.Email, clientIP)
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Either password or API key is required")
		return
	}

	// Check if user is active
	if !user.IsActive {
		log.Printf("AUTH: Failed login attempt - email=%s user_id=%d method=%s reason=account_disabled ip=%s", user.Email, user.ID, authMethod, clientIP)
		httpresp.WriteForbidden(w, "Account is disabled")
		return
	}

	// Create session
	sessionID, err := GenerateSessionID()
	if err != nil {
		httpresp.WriteInternalError(w, "Failed to create session")
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
		httpresp.WriteInternalError(w, "Failed to create session")
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
	httpresp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"user":    user,
	})
}

// LogoutHandler handles user logout requests.
//
// @Summary     Logout user
// @Description Deletes active session and clears the session cookie.
// @Tags        auth
// @Produce     json
// @Success     200 {object} map[string]interface{} "Logout success"
// @Failure     405 {object} map[string]interface{} "Method not allowed"
// @Router      /api/auth/logout [post]
func (m *Manager) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if !httpresp.RequireMethodJSON(w, r, http.MethodPost) {
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

	httpresp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

// ForgotPasswordHandler handles password reset requests.
//
// @Summary     Request password reset
// @Description Generates a password reset token and sends reset email when account exists.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body object true "Forgot-password request"
// @Success     200 {object} map[string]interface{} "Request accepted"
// @Failure     400 {object} map[string]interface{} "Invalid request"
// @Failure     500 {object} map[string]interface{} "Server error"
// @Router      /api/auth/forgot-password [post]
func (m *Manager) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if !httpresp.RequireMethodJSON(w, r, http.MethodPost) {
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Invalid request body")
		return
	}

	req.Email = NormalizeEmail(req.Email)

	// Get user from database
	user, err := GetUserByEmail(r.Context(), m.DB, m.DBDriver, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't reveal if email exists or not (security)
			httpresp.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"message": "If that email address is registered, a password reset link has been sent.",
			})
			return
		}
		httpresp.WriteInternalError(w, "Database error")
		return
	}

	// Create password reset token
	token, err := GenerateToken()
	if err != nil {
		httpresp.WriteInternalError(w, "Failed to generate reset token")
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
		httpresp.WriteInternalError(w, "Failed to create reset token")
		return
	}

	// Send password reset email
	if m.EmailSender != nil {
		resetURL := fmt.Sprintf("%s/reset-password?token=%s", m.BaseURL, token)
		if err := m.EmailSender.SendPasswordResetEmail(req.Email, resetURL); err != nil {
			log.Printf("ERROR: Failed to send password reset email to %s: %v", req.Email, err)
		}
	}

	httpresp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "If that email address is registered, a password reset link has been sent.",
	})
}

// ResetPasswordHandler handles password reset with token.
//
// @Summary     Reset password with token
// @Description Resets account password using a valid reset token.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body object true "Reset-password request"
// @Success     200 {object} map[string]interface{} "Password reset success"
// @Failure     400 {object} map[string]interface{} "Invalid token or request"
// @Failure     500 {object} map[string]interface{} "Server error"
// @Router      /api/auth/reset-password [post]
func (m *Manager) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if !httpresp.RequireMethodJSON(w, r, http.MethodPost) {
		return
	}

	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Invalid request body")
		return
	}

	if err := ValidatePassword(req.Password); err != nil {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, err.Error())
		return
	}

	// Get reset token from database
	token, err := GetPasswordResetToken(r.Context(), m.DB, m.DBDriver, req.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Invalid or expired reset token")
			return
		}
		httpresp.WriteInternalError(w, "Database error")
		return
	}

	// Check if token is expired
	now := time.Now().Unix()
	if token.ExpiresAt < now {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Reset token has expired")
		return
	}

	// Check if token has already been used
	if token.Used {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Reset token has already been used")
		return
	}

	// Update user password
	if err := UpdateUserPassword(r.Context(), m.DB, m.DBDriver, token.UserID, req.Password); err != nil {
		httpresp.WriteInternalError(w, "Failed to update password")
		return
	}

	// Mark token as used
	_ = MarkPasswordResetTokenUsed(r.Context(), m.DB, m.DBDriver, token.ID)

	// Get user email for notification
	user, err := GetUserByID(r.Context(), m.DB, m.DBDriver, token.UserID)
	if err == nil && m.EmailSender != nil {
		_ = m.EmailSender.SendPasswordChangedEmail(user.Email)
	}

	httpresp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password reset successful",
	})
}

// VerifyEmailHandler handles email verification with token.
//
// @Summary     Verify email address
// @Description Verifies user email using a token and redirects to the map UI on success.
// @Tags        auth
// @Produce     json
// @Param       token query string true "Verification token"
// @Success     303 {string} string "Redirect to app"
// @Failure     400 {object} map[string]interface{} "Invalid token"
// @Failure     500 {object} map[string]interface{} "Server error"
// @Router      /api/auth/verify-email [get]
func (m *Manager) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	if !httpresp.RequireMethodJSON(w, r, http.MethodGet) {
		return
	}

	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Token is required")
		return
	}

	// Get verification token from database
	token, err := GetEmailVerificationToken(r.Context(), m.DB, m.DBDriver, tokenStr)
	if err != nil {
		if err == sql.ErrNoRows {
			httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Invalid or expired verification token")
			return
		}
		httpresp.WriteInternalError(w, "Database error")
		return
	}

	// Check if token is expired
	now := time.Now().Unix()
	if token.ExpiresAt < now {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Verification token has expired")
		return
	}

	// Check if token has already been used
	if token.Used {
		httpresp.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Email already verified",
		})
		return
	}

	// Get user info for logging
	user, err := GetUserByID(r.Context(), m.DB, m.DBDriver, token.UserID)
	if err != nil {
		httpresp.WriteInternalError(w, "Failed to fetch user")
		return
	}

	// Verify user email
	if err := VerifyUserEmail(r.Context(), m.DB, m.DBDriver, token.UserID); err != nil {
		httpresp.WriteInternalError(w, "Failed to verify email")
		return
	}

	// Mark token as used
	_ = MarkEmailVerificationTokenUsed(r.Context(), m.DB, m.DBDriver, token.ID)

	log.Printf("AUTH: Email verified - user_id=%d email=%s ip=%s", user.ID, user.Email, getClientIP(r))

	// Redirect to map page with success message
	http.Redirect(w, r, "/?emailVerified=true", http.StatusSeeOther)
}

// ProfileHandler returns the current user's profile.
//
// @Summary     Get current user profile
// @Description Returns profile details for the authenticated user.
// @Tags        auth
// @Produce     json
// @Success     200 {object} User "User profile"
// @Failure     401 {object} map[string]interface{} "Unauthorized"
// @Router      /api/user/profile [get]
func (m *Manager) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent CloudFront from caching user-specific data
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if !httpresp.RequireMethodJSON(w, r, http.MethodGet) {
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		httpresp.WriteUnauthorized(w, "Unauthorized")
		return
	}

	// Don't expose password hash
	user.PasswordHash = ""

	httpresp.WriteJSON(w, http.StatusOK, user)
}

// ChangePasswordHandler allows a logged-in user to change their password.
//
// @Summary     Change current user password
// @Description Changes password for the authenticated user after validating current password.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body object true "Change-password request"
// @Success     200 {object} map[string]interface{} "Password changed"
// @Failure     400 {object} map[string]interface{} "Invalid request"
// @Failure     401 {object} map[string]interface{} "Unauthorized or wrong password"
// @Failure     500 {object} map[string]interface{} "Server error"
// @Router      /api/user/change-password [post]
func (m *Manager) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if !httpresp.RequireMethodJSON(w, r, http.MethodPost) {
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		httpresp.WriteUnauthorized(w, "Unauthorized")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Invalid request body")
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "Current password and new password are required")
		return
	}

	if len(req.NewPassword) < 8 {
		httpresp.WriteBadRequest(w, httpresp.CodeBadRequest, "New password must be at least 8 characters")
		return
	}

	// Verify current password
	if !VerifyPassword(user.PasswordHash, req.CurrentPassword) {
		httpresp.WriteUnauthorized(w, "Current password is incorrect")
		return
	}

	// Update password
	if err := UpdateUserPassword(r.Context(), m.DB, m.DBDriver, user.ID, req.NewPassword); err != nil {
		httpresp.WriteInternalError(w, "Failed to update password")
		return
	}

	log.Printf("AUTH: Password changed - user_id=%d email=%s ip=%s", user.ID, user.Email, getClientIP(r))

	httpresp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password changed successfully",
	})
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
