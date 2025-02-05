package request

type CreateUserRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Password   string `json:"password"`
	Picture    string `json:"picture"`
	IPAddress  string `json:"ip_address"`
	DeviceInfo string `json:"device_info"`
}
type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	DeviceInfo string `json:"device_info"`
	IPAddress  string `json:"ip_address"`
}

type UpdateUserRequest struct {
	UserID   string  `json:"user_id"`
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	Role     *string `json:"role,omitempty"`
	Picture  *string `json:"picture,omitempty"`
}
type GetUserRequest struct {
	UserID string `json:"user_id"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email"`
}

type RequestPasswordResetResponse struct {
	Success    bool   `json:"success"`
	ResetToken string `json:"reset_token"`
	ExpiresAt  int64  `json:"expires_at"`
}

type ResetPasswordRequest struct {
	ResetToken  string `json:"reset_token"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type LogoutRequest struct {
	AccessToken string `json:"access_token"`
	SessionID   string `json:"session_id"`
	UserID      string `json:"user_id"`
}

type ValidateTokenRequest struct {
	AccessToken string `json:"access_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	SessionID    string `json:"session_id"`
}

type UpdateUserResponse struct {
	UserID    string `json:"user_id"`
	UpdatedAt int64  `json:"updated_at"`
}

type VerifyEmailRequest struct {
	VerificationToken string `json:"verification_token"`
	VerifyCode        string `json:"verify_code"`
}
type RevokeSessionRequest struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
}
type ListSessionRequest struct {
	UserID string `json:"user_id"`
}

// VerifyEmailResponse represents the response for email verification
type VerifyEmailResponse struct {
	Success bool   `json:"success"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// ResendVerificationRequest represents the request to resend verification
type ResendVerificationRequest struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Token  string `json:"token"`
}
