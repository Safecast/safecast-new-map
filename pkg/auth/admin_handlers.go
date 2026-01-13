package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// AdminListUsersHandler returns a list of all users (admin only).
func (m *Manager) AdminListUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user from context
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Check if user is admin (for now, check if they have a specific email or flag)
	// TODO: Add proper role-based access control
	if !isAdmin(user) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Admin access required"})
		return
	}

	// Get all users from database
	users, err := GetAllUsers(r.Context(), m.DB, m.DBDriver)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"users": users})
}

// AdminCreateUserHandler creates a new user (admin only).
func (m *Manager) AdminCreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user from context
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Check if user is admin
	if !isAdmin(user) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Admin access required"})
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
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate input
	if err := ValidateEmail(req.Email); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if req.Username != "" {
		if err := ValidateUsername(req.Username); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
	}

	// Generate password if not provided
	password := req.Password
	if password == "" {
		generatedPassword, err := GenerateRandomPassword(16)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate password"})
			return
		}
		password = generatedPassword
		req.RequiresPasswordSetup = true
	} else {
		if err := ValidatePassword(password); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
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
			writeJSON(w, http.StatusConflict, map[string]string{"error": "User with this email already exists"})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
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

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":      userID,
		"email":   req.Email,
		"message": "User created successfully",
	})
}

// AdminUpdateUserHandler updates an existing user (admin only).
func (m *Manager) AdminUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user from context
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Check if user is admin
	if !isAdmin(user) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Admin access required"})
		return
	}

	// Get user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID required"})
		return
	}
	userIDStr := pathParts[3]
	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Parse request
	var req struct {
		Email         *string `json:"email"`
		Username      *string `json:"username"`
		IsActive      *bool   `json:"is_active"`
		EmailVerified *bool   `json:"email_verified"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Get existing user
	targetUser, err := GetUserByID(r.Context(), m.DB, m.DBDriver, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
		}
		return
	}

	// Update fields if provided
	if req.Email != nil {
		if err := ValidateEmail(*req.Email); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		targetUser.Email = *req.Email
	}

	if req.Username != nil {
		if *req.Username != "" {
			if err := ValidateUsername(*req.Username); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
		}
		targetUser.Username = *req.Username
	}

	if req.IsActive != nil {
		targetUser.IsActive = *req.IsActive
	}

	if req.EmailVerified != nil {
		targetUser.EmailVerified = *req.EmailVerified
	}

	// Update user in database
	if err := UpdateUser(r.Context(), m.DB, m.DBDriver, targetUser); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key") {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "Email or username already in use"})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User updated successfully",
		"user":    targetUser,
	})
}

// AdminDeleteUserHandler deletes a user (admin only).
func (m *Manager) AdminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user from context
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Check if user is admin
	if !isAdmin(user) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Admin access required"})
		return
	}

	// Get user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID required"})
		return
	}
	userIDStr := pathParts[3]
	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Prevent self-deletion
	if targetUserID == user.ID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot delete your own account"})
		return
	}

	// Delete user
	if err := DeleteUser(r.Context(), m.DB, m.DBDriver, targetUserID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// AdminResetUserPasswordHandler sends a password reset email to a user (admin only).
func (m *Manager) AdminResetUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user from context
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Check if user is admin
	if !isAdmin(user) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Admin access required"})
		return
	}

	// Get user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "User ID required"})
		return
	}
	userIDStr := pathParts[3]
	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Get target user
	targetUser, err := GetUserByID(r.Context(), m.DB, m.DBDriver, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
		}
		return
	}

	// Generate password reset token
	token, err := GenerateToken()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
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
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create reset token"})
		return
	}

	// Send email if email sender is configured
	if m.EmailSender != nil {
		resetURL := fmt.Sprintf("%s/reset-password?token=%s", m.BaseURL, token)
		if err := m.EmailSender.SendPasswordResetEmail(targetUser.Email, resetURL); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Password reset email sent successfully",
		"token":   token, // Include token for admin to share manually if needed
	})
}

// isAdmin checks if a user has admin privileges.
// TODO: Implement proper role-based access control.
// For now, this is a placeholder that always returns true for authenticated users.
// In production, you should check against a role field or admin flag.
func isAdmin(user *User) bool {
	// Placeholder: Check if user email contains "admin" or has a specific domain
	// In production, add an IsAdmin or Role field to the User model
	return strings.Contains(strings.ToLower(user.Email), "admin") ||
		strings.Contains(strings.ToLower(user.Email), "@safecast.org")
}
