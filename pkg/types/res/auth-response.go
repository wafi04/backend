package response

import (
	"time"

	"github.com/wafi04/backend/pkg/types"
)

type GetUserResponse struct {
	User types.UserInfo `json:"user,omitempty"`
}
type CreateUserResponse struct {
	UserID      string         `json:"user_id"`
	AccessToken string         `json:"access_token"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	Role        string         `json:"role"`
	CreatedAt   int64          `json:"created_at"`
	Picture     string         `json:"picture"`
	SessionInfo *types.Session `json:"session_info,omitempty"`
}

type LoginResponse struct {
	Refresh_token string             `json:"refresh_token"`
	AccessToken   string             `json:"access_token"`
	UserID        string             `json:"user_id"`
	UserInfo      *types.UserInfo    `json:"user_info"`
	SessionInfo   *types.SessionInfo `json:"session_info"`
	ExpiresAt     int64              `json:"expires_at"`
}

type LogoutResponse struct {
	Success bool `json:"success"`
}

type ValidateTokenResponse struct {
	Valid     bool   `json:"valid"`
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expires_at"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type VerifyEmailResponse struct {
	Success bool   `json:"success"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}
type RevokeSessionResponse struct {
	Success bool `json:"succcess"`
}

type ListSessionResponse struct {
	Sessions []*types.SessionInfo `json:"sessions"`
}

type ResendVerificationResponse struct {
	Success           bool      `json:"success"`
	VerificationToken string    `json:"verification_token"`
	ExpiresAt         time.Time `json:"expires_at"`
	VerifyCode        string    `json:"verify_code"`
}

type ResetPasswordResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	UpdatedAt int64  `json:"updated_at"`
}
