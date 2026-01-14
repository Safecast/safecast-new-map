package auth

// User represents a user account in the system.
type User struct {
	ID                    int64  `json:"id"`
	Email                 string `json:"email"`
	PasswordHash          string `json:"-"` // Never expose password hash in JSON
	Username              string `json:"username,omitempty"`
	EmailVerified         bool   `json:"email_verified"`
	CreatedAt             int64  `json:"created_at"`
	UpdatedAt             int64  `json:"updated_at"`
	LastLoginAt           int64  `json:"last_login_at,omitempty"`
	IsActive              bool   `json:"is_active"`
	ExternalID            string `json:"external_id,omitempty"`
	ExternalSource        string `json:"external_source,omitempty"`
	RequiresPasswordSetup bool   `json:"requires_password_setup"`
	APIKey                string `json:"api_key,omitempty"`
}

// Session represents a user session stored in the database.
type Session struct {
	ID             string `json:"id"`
	UserID         int64  `json:"userId"`
	CreatedAt      int64  `json:"createdAt"`
	ExpiresAt      int64  `json:"expiresAt"`
	IPAddress      string `json:"ipAddress,omitempty"`
	UserAgent      string `json:"userAgent,omitempty"`
	LastActivityAt int64  `json:"lastActivityAt"`
}

// PasswordResetToken represents a token for password reset requests.
type PasswordResetToken struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userId"`
	Token     string `json:"token"`
	CreatedAt int64  `json:"createdAt"`
	ExpiresAt int64  `json:"expiresAt"`
	Used      bool   `json:"used"`
	UsedAt    int64  `json:"usedAt,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
}

// EmailVerificationToken represents a token for email verification.
type EmailVerificationToken struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userId"`
	Token     string `json:"token"`
	CreatedAt int64  `json:"createdAt"`
	ExpiresAt int64  `json:"expiresAt"`
	Used      bool   `json:"used"`
	UsedAt    int64  `json:"usedAt,omitempty"`
}
