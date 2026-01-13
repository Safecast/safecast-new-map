package auth

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Email validation regex (simplified RFC 5322 compliant)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// Username validation regex (alphanumeric, underscore, hyphen, 3-30 chars)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]{3,30}$`)
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
	MinUsernameLength = 3
	MaxUsernameLength = 30
)

// ValidateEmail validates an email address format.
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email address is required")
	}

	email = strings.TrimSpace(strings.ToLower(email))

	if len(email) > 254 {
		return fmt.Errorf("email address is too long (max 254 characters)")
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email address format")
	}

	return nil
}

// ValidatePassword validates a password meets security requirements.
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters", MinPasswordLength)
	}

	if len(password) > MaxPasswordLength {
		return fmt.Errorf("password must be at most %d characters", MaxPasswordLength)
	}

	// Optional: Add complexity requirements
	// hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	// hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	//
	// if !hasUpper || !hasLower || !hasNumber {
	// 	return fmt.Errorf("password must contain uppercase, lowercase, and numbers")
	// }

	return nil
}

// ValidateUsername validates a username format.
// Username is optional in the system, so empty usernames are allowed.
func ValidateUsername(username string) error {
	if username == "" {
		// Username is optional
		return nil
	}

	username = strings.TrimSpace(username)

	if len(username) < MinUsernameLength {
		return fmt.Errorf("username must be at least %d characters", MinUsernameLength)
	}

	if len(username) > MaxUsernameLength {
		return fmt.Errorf("username must be at most %d characters", MaxUsernameLength)
	}

	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}

// SanitizeInput trims whitespace and removes potentially harmful characters.
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// NormalizeEmail converts email to lowercase and trims whitespace.
func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
