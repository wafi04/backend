package types

type User struct {
	UserID          string   `json:"user_id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	Role            string   `json:"role"`
	PasswordHash    string   `json:"password_hash"`
	IsEmailVerified bool     `json:"is_email_verified"`
	IsActive        bool     `json:"is_active"`
	CreatedAt       int64    `json:"created_at"`
	UpdatedAt       int64    `json:"updated_at"`
	LastLoginAt     int64    `json:"last_login_at"`
	ActiveSessions  []string `json:"active_sessions"`
	Picture         string   `json:"picture"`
}

type UserInfo struct {
	UserID          string `json:"user_id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	IsEmailVerified bool   `json:"is_email_verified"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
	LastLoginAt     int64  `json:"last_login_at"`
	Picture         string `json:"picture"`
}

type Session struct {
	SessionID      string `json:"session_id"`
	UserID         string `json:"user_id"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	DeviceInfo     string `json:"device_info"`
	IPAddress      string `json:"ip_address"`
	CreatedAt      int64  `json:"created_at"`
	ExpiresAt      int64  `json:"expires_at"`
	LastActivityAt int64  `json:"last_activity_at"`
	IsActive       bool   `json:"is_active"`
}

type VerificationToken struct {
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"`
	CreatedAt int64  `json:"created_at"`
	ExpiresAt int64  `json:"expires_at"`
	IsUsed    bool   `json:"is_used"`
}

type SessionInfo struct {
	SessionID      string `json:"session_id"`
	DeviceInfo     string `json:"device_info"`
	IPAddress      string `json:"ip_address"`
	CreatedAt      int64  `json:"created_at"`
	LastActivityAt int64  `json:"last_activity_at"`
}
