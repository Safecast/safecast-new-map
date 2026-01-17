package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost is the cost factor for bcrypt hashing (higher = more secure but slower)
	BcryptCost = 12

	// TokenLength is the length of generated tokens in bytes (32 bytes = 256 bits)
	TokenLength = 32
)

// HashPassword hashes a password using bcrypt with the configured cost factor.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// VerifyPassword compares a plaintext password with a bcrypt hash.
// Returns true if the password matches the hash.
func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a cryptographically secure random token.
// The token is base64-encoded for easy transmission in URLs and JSON.
func GenerateToken() (string, error) {
	bytes := make([]byte, TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateSessionID generates a unique session ID.
// This is an alias for GenerateToken for clarity in code.
func GenerateSessionID() (string, error) {
	return GenerateToken()
}

// GenerateRandomPassword generates a random password of the specified length.
// Used for creating temporary passwords for migrated users.
func GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random password: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
