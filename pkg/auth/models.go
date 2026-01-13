package auth

// User represents a user account in the system.
type User struct {
	ID                   int64  `json:"id"`
	Email                string `json:"email"`
	PasswordHash         string `json:"-"` // Never expose password hash in JSON
	Username             string `json:"username,omitempty"`
	EmailVerified        bool   `json:"emailVerified"`
	CreatedAt            int64  `json:"createdAt"`
	UpdatedAt            int64  `json:"updatedAt"`
	LastLoginAt          int64  `json:"lastLoginAt,omitempty"`
	IsActive             bool   `json:"isActive"`
	ExternalID           string `json:"externalId,omitempty"`
	ExternalSource       string `json:"externalSource,omitempty"`
	RequiresPasswordSetup bool  `json:"requiresPasswordSetup"`
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
