package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// AdminListUsersHandler returns a list of users with pagination support (admin only).
// Query parameters: limit, offset, search
// Note: Authentication is handled by the route handler via password parameter check.
func (m *Manager) AdminListUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse pagination parameters
	limit := 50
	offset := 0
	search := r.URL.Query().Get("search")

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get users with pagination
	users, err := GetUsersPaginated(r.Context(), m.DB, m.DBDriver, limit, offset, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Failed to fetch users"}, http.StatusInternalServerError)
		return
	}

	// Get total count for pagination
	total, err := CountUsers(r.Context(), m.DB, m.DBDriver, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Failed to count users"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}, http.StatusOK)
}

// AdminCreateUserHandler creates a new user (admin only).
// Note: Authentication is handled by the route handler via password parameter check.
func (m *Manager) AdminCreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req struct {
		Email                string `json:"email"`
		Username             string `json:"username"`
		Password             string `json:"password"`
		SendWelcomeEmail     bool   `json:"send_welcome_email"`
		RequiresPasswordSetup bool   `json:"requires_password_setup"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	// Validate input
	if err := ValidateEmail(req.Email); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if req.Username != "" {
		if err := ValidateUsername(req.Username); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
	}

	// Generate password if not provided
	password := req.Password
	if password == "" {
		generatedPassword, err := GenerateRandomPassword(16)
		if err != nil {
			writeJSON(w, map[string]string{"error": "Failed to generate password"}, http.StatusInternalServerError)
			return
		}
		password = generatedPassword
		req.RequiresPasswordSetup = true
	} else {
		if err := ValidatePassword(password); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
	}

	// Create user
	newUser := &User{
		Email:                 req.Email,
		Username:              req.Username,
		EmailVerified:         !req.RequiresPasswordSetup,
		IsActive:              true,
		RequiresPasswordSetup: req.RequiresPasswordSetup,
	}

	userID, err := CreateUser(r.Context(), m.DB, m.DBDriver, newUser, password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key") {
			writeJSON(w, map[string]string{"error": "User with this email already exists"}, http.StatusConflict)
		} else {
			writeJSON(w, map[string]string{"error": "Failed to create user"}, http.StatusInternalServerError)
		}
		return
	}

	// Send welcome or password setup email if requested
	if m.EmailSender != nil && (req.SendWelcomeEmail || req.RequiresPasswordSetup) {
		if req.RequiresPasswordSetup {
			// Generate password reset token
			token, err := GenerateToken()
			if err == nil {
				expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()
				resetToken := &PasswordResetToken{
					UserID:    userID,
					Token:     token,
					ExpiresAt: expiresAt,
					Used:      false,
				}
				if err := CreatePasswordResetToken(r.Context(), m.DB, m.DBDriver, resetToken); err == nil {
					setupURL := fmt.Sprintf("%s/reset-password?token=%s", m.BaseURL, token)
					m.EmailSender.SendPasswordSetupEmail(req.Email, req.Username, setupURL)
				}
			}
		} else {
			// Send regular welcome email with verification
			token, err := GenerateToken()
			if err == nil {
				expiresAt := time.Now().Add(24 * time.Hour).Unix()
				verifyToken := &EmailVerificationToken{
					UserID:    userID,
					Token:     token,
					ExpiresAt: expiresAt,
					Used:      false,
				}
				if err := CreateEmailVerificationToken(r.Context(), m.DB, m.DBDriver, verifyToken); err == nil {
					verifyURL := fmt.Sprintf("%s/api/auth/verify-email?token=%s", m.BaseURL, token)
					m.EmailSender.SendWelcomeEmail(req.Email, req.Username, verifyURL)
				}
			}
		}
	}

	writeJSON(w, map[string]interface{}{
		"id":      userID,
		"email":   req.Email,
		"message": "User created successfully",
	}, http.StatusCreated)
}

// AdminUpdateUserHandler updates an existing user (admin only).
// Note: Authentication is handled by the route handler via password parameter check.
func (m *Manager) AdminUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeJSON(w, map[string]string{"error": "User ID required"}, http.StatusBadRequest)
		return
	}
	userIDStr := pathParts[3]
	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Invalid user ID"}, http.StatusBadRequest)
		return
	}

	// Parse request
	var req struct {
		Email         *string `json:"email"`
		Username      *string `json:"username"`
		IsActive      *bool   `json:"is_active"`
		IsAdmin       *bool   `json:"is_admin"`
		EmailVerified *bool   `json:"email_verified"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	// Get existing user
	targetUser, err := GetUserByID(r.Context(), m.DB, m.DBDriver, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, map[string]string{"error": "User not found"}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": "Failed to fetch user"}, http.StatusInternalServerError)
		}
		return
	}

	// Update fields if provided
	if req.Email != nil {
		if err := ValidateEmail(*req.Email); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
		targetUser.Email = *req.Email
	}

	if req.Username != nil {
		if *req.Username != "" {
			if err := ValidateUsername(*req.Username); err != nil {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
				return
			}
		}
		targetUser.Username = *req.Username
	}

	if req.IsActive != nil {
		targetUser.IsActive = *req.IsActive
	}

	if req.IsAdmin != nil {
		targetUser.IsAdmin = *req.IsAdmin
	}

	if req.EmailVerified != nil {
		targetUser.EmailVerified = *req.EmailVerified
	}

	// Update user in database
	if err := UpdateUser(r.Context(), m.DB, m.DBDriver, targetUser); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key") {
			writeJSON(w, map[string]string{"error": "Email or username already in use"}, http.StatusConflict)
		} else {
			writeJSON(w, map[string]string{"error": "Failed to update user"}, http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]interface{}{
		"message": "User updated successfully",
		"user":    targetUser,
	}, http.StatusOK)
}

// AdminDeleteUserHandler deletes a user (admin only).
// Note: Authentication is handled by the route handler via password parameter check.
func (m *Manager) AdminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeJSON(w, map[string]string{"error": "User ID required"}, http.StatusBadRequest)
		return
	}
	userIDStr := pathParts[3]
	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Invalid user ID"}, http.StatusBadRequest)
		return
	}

	// Delete user
	if err := DeleteUser(r.Context(), m.DB, m.DBDriver, targetUserID); err != nil {
		writeJSON(w, map[string]string{"error": "Failed to delete user"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "User deleted successfully"}, http.StatusOK)
}

// AdminResetUserPasswordHandler sends a password reset email to a user (admin only).
// Note: Authentication is handled by the route handler via password parameter check.
func (m *Manager) AdminResetUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeJSON(w, map[string]string{"error": "User ID required"}, http.StatusBadRequest)
		return
	}
	userIDStr := pathParts[3]
	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Invalid user ID"}, http.StatusBadRequest)
		return
	}

	// Get target user
	targetUser, err := GetUserByID(r.Context(), m.DB, m.DBDriver, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, map[string]string{"error": "User not found"}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": "Failed to fetch user"}, http.StatusInternalServerError)
		}
		return
	}

	// Generate password reset token
	token, err := GenerateToken()
	if err != nil {
		writeJSON(w, map[string]string{"error": "Failed to generate token"}, http.StatusInternalServerError)
		return
	}

	// Token expires in 1 hour for admin-initiated resets
	expiresAt := time.Now().Add(1 * time.Hour).Unix()
	resetToken := &PasswordResetToken{
		UserID:    targetUserID,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	if err := CreatePasswordResetToken(r.Context(), m.DB, m.DBDriver, resetToken); err != nil {
		writeJSON(w, map[string]string{"error": "Failed to create reset token"}, http.StatusInternalServerError)
		return
	}

	// Send email if email sender is configured
	if m.EmailSender != nil {
		resetURL := fmt.Sprintf("%s/reset-password?token=%s", m.BaseURL, token)
		if err := m.EmailSender.SendPasswordResetEmail(targetUser.Email, resetURL); err != nil {
			// Log the actual error and return it to the user
			errMsg := fmt.Sprintf("Failed to send email: %v", err)
			writeJSON(w, map[string]string{"error": errMsg}, http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]string{
		"message": "Password reset email sent successfully",
		"token":   token, // Include token for admin to share manually if needed
	}, http.StatusOK)
}

// AdminRegenerateAPIKeyHandler regenerates an API key for a user (admin only).
// POST /api/admin/users/{id}/regenerate-api-key
// Note: Authentication is handled by the route handler via password parameter check.
func (m *Manager) AdminRegenerateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeJSON(w, map[string]string{"error": "User ID required"}, http.StatusBadRequest)
		return
	}
	userIDStr := pathParts[3]
	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Invalid user ID"}, http.StatusBadRequest)
		return
	}

	// Get target user
	targetUser, err := GetUserByID(r.Context(), m.DB, m.DBDriver, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, map[string]string{"error": "User not found"}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": "Failed to fetch user"}, http.StatusInternalServerError)
		}
		return
	}

	// Regenerate API key
	newAPIKey, err := RegenerateAPIKey(r.Context(), m.DB, m.DBDriver, targetUserID)
	if err != nil {
		log.Printf("ERROR: Failed to regenerate API key for user %d: %v", targetUserID, err)
		writeJSON(w, map[string]string{"error": "Failed to regenerate API key"}, http.StatusInternalServerError)
		return
	}

	// Send email notification if email sender is configured
	if m.EmailSender != nil {
		if err := m.EmailSender.SendAPIKeyRegeneratedEmail(targetUser.Email, targetUser.Username, newAPIKey); err != nil {
			log.Printf("WARNING: Failed to send API key regeneration email to %s: %v", targetUser.Email, err)
			// Don't fail the request - API key was successfully regenerated
		}
	}

	writeJSON(w, map[string]interface{}{
		"message": "API key regenerated successfully",
		"api_key": newAPIKey,
		"user_id": targetUserID,
	}, http.StatusOK)
}

// IsAdmin checks if a user has admin privileges.
func IsAdmin(user *User) bool {
	return user != nil && user.IsAdmin
}
