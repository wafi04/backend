package authrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/wafi04/backend/pkg/logger"
	"github.com/wafi04/backend/pkg/middleware"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
	response "github.com/wafi04/backend/pkg/types/res"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	DB     *sqlx.DB
	logger logger.Logger
}

func NewDB(DB *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		DB: DB,
	}
}

func (D *AuthRepository) CreateUser(ctx context.Context, req *request.CreateUserRequest) (response.CreateUserResponse, error) {

	role := "user"
	if req.Email == "wafiq610@gmail.com" {
		role = "admin"
	}
	userID := uuid.New().String()
	now := time.Now()
	D.logger.Log(logger.InfoLevel, "Data : %v", req)
	query := `
        INSERT INTO users (
            user_id, name, email, password_hash, role,
            is_active, is_email_verified, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := D.DB.ExecContext(
		ctx, query,
		userID, req.Name, req.Email, req.Password, role,
		true, false, now, now,
	)

	if err != nil {
		return response.CreateUserResponse{}, fmt.Errorf("failed to create verification token: %w", err)
	}

	accces_token, err := middleware.GenerateToken(&types.UserInfo{
		UserID:          userID,
		Name:            req.Name,
		Email:           req.Email,
		Role:            role,
		IsEmailVerified: false,
	}, 24)
	if err != nil {
		return response.CreateUserResponse{}, fmt.Errorf("failed to generate tokens: %w", err)
	}
	refresh_token, err := middleware.GenerateToken(&types.UserInfo{
		UserID:          userID,
		Name:            req.Name,
		Email:           req.Email,
		Role:            role,
		IsEmailVerified: false,
	}, 168)
	if err != nil {
		return response.CreateUserResponse{}, fmt.Errorf("failed to generate tokens: %w", err)
	}

	session := types.Session{
		SessionID:      uuid.New().String(),
		UserID:         userID,
		AccessToken:    accces_token,
		RefreshToken:   refresh_token,
		IPAddress:      req.IPAddress,
		DeviceInfo:     req.DeviceInfo,
		CreatedAt:      time.Now().Unix(),
		LastActivityAt: time.Now().Unix(),
		IsActive:       true,
		ExpiresAt:      time.Now().Unix(),
	}

	err = D.CreateSession(ctx, &session)
	if err != nil {
		return response.CreateUserResponse{}, fmt.Errorf("failed to create session: %w", err)
	}

	return response.CreateUserResponse{
		UserID:  userID,
		Name:    req.Name,
		Email:   req.Email,
		Role:    role,
		Picture: req.Picture,
		SessionInfo: &types.Session{
			SessionID:  session.SessionID,
			DeviceInfo: session.DeviceInfo,
			IPAddress:  session.IPAddress,
		},
	}, nil

}

type dbUser struct {
	UserID          string
	Name            string
	Email           string
	Role            string
	Password        string
	Picture         string
	IsEmailVerified bool
	CreatedAt       int64
	UpdatedAt       int64
	LastLoginAt     int64
	IsActive        bool
}

func (r *AuthRepository) Login(ctx context.Context, login *request.LoginRequest) (*response.LoginResponse, error) {
	query := `
    SELECT
        user_id,
        name,
        email,
        role,
        password_hash,
        COALESCE(picture, ''),
        COALESCE(is_email_verified, false)::boolean,  
        EXTRACT(EPOCH FROM created_at)::bigint,
        EXTRACT(EPOCH FROM updated_at)::bigint,
        EXTRACT(EPOCH FROM COALESCE(last_login_at, created_at))::bigint,
        is_active::boolean
    FROM users
    WHERE email = $1
`
	var dbuser dbUser
	err := r.DB.QueryRowContext(ctx, query, login.Email).Scan(
		&dbuser.UserID,
		&dbuser.Name,
		&dbuser.Email,
		&dbuser.Role,
		&dbuser.Password,
		&dbuser.Picture,
		&dbuser.IsEmailVerified,
		&dbuser.CreatedAt,
		&dbuser.UpdatedAt,
		&dbuser.LastLoginAt,
		&dbuser.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Convert Unix timestamps to time.Time
	createdAt := time.Unix(dbuser.CreatedAt, 0)
	updatedAt := time.Unix(dbuser.UpdatedAt, 0)
	lastLoginAt := time.Unix(dbuser.LastLoginAt, 0)

	userInfo := &types.UserInfo{
		UserID:          dbuser.UserID,
		Name:            dbuser.Name,
		Email:           dbuser.Email,
		Role:            dbuser.Role,
		IsEmailVerified: dbuser.IsEmailVerified,
		CreatedAt:       createdAt.Unix(),
		UpdatedAt:       updatedAt.Unix(),
		LastLoginAt:     lastLoginAt.Unix(),
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(dbuser.Password), []byte(login.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	access_token, err := middleware.GenerateToken(userInfo, 24)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}
	refresh_token, err := middleware.GenerateToken(userInfo, 168)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Check for existing session
	query = `
        SELECT 
            session_id, 
            ip_address,
            device_info, 
            EXTRACT(EPOCH FROM created_at)::bigint, 
            EXTRACT(EPOCH FROM last_activity_at)::bigint
        FROM sessions 
        WHERE user_id = $1 AND is_active = true AND device_info = $2
    `
	var existingSession types.Session
	err = r.DB.QueryRowContext(ctx, query, userInfo.UserID, login.DeviceInfo).Scan(
		&existingSession.SessionID,
		&existingSession.IPAddress,
		&existingSession.DeviceInfo,
		&existingSession.CreatedAt,
		&existingSession.LastActivityAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking existing session: %w", err)
	}

	// Update or create session
	if err == sql.ErrNoRows {
		existingSession = types.Session{
			SessionID:      uuid.New().String(),
			UserID:         userInfo.UserID,
			AccessToken:    access_token,
			RefreshToken:   refresh_token,
			IPAddress:      login.IPAddress,
			DeviceInfo:     login.DeviceInfo,
			CreatedAt:      time.Now().Unix(),
			LastActivityAt: time.Now().Unix(),
			IsActive:       true,
			ExpiresAt:      time.Now().Add(7 * 24 * time.Hour).Unix(),
		}
		err = r.CreateSession(ctx, &existingSession)
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %w", err)
		}
	} else {
		_, err = r.UpdateSessionActivity(ctx, existingSession.SessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to update session activity: %w", err)
		}
	}

	_, err = r.DB.ExecContext(
		ctx,
		"UPDATE users SET last_login_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1",
		userInfo.UserID,
	)
	if err != nil {
		r.logger.Log(logger.ErrorLevel, "Failed to update last login: %v", err)
	}

	return &response.LoginResponse{
		AccessToken:   access_token,
		Refresh_token: refresh_token,
		UserID:        userInfo.UserID,
		SessionInfo: &types.SessionInfo{
			SessionID:      existingSession.SessionID,
			DeviceInfo:     existingSession.DeviceInfo,
			IPAddress:      existingSession.IPAddress,
			CreatedAt:      existingSession.CreatedAt,
			LastActivityAt: existingSession.LastActivityAt,
		},
	}, nil
}

func (sr *AuthRepository) GetUser(ctx context.Context, req *request.GetUserRequest) (*response.GetUserResponse, error) {
	query := `
        SELECT 
            user_id, 
            name, 
            email,
            picture, 
            role, 
            is_active, 
            is_email_verified,
            created_at, 
            updated_at, 
            last_login_at
        FROM users
        WHERE user_id = $1
    `

	user := &response.GetUserResponse{
		User: types.UserInfo{},
	}

	var (
		isActive                          bool
		createdAt, updatedAt, lastLoginAt time.Time
		picture                           sql.NullString
	)
	err := sr.DB.QueryRowContext(ctx, query, req.UserID).Scan(
		&user.User.UserID,
		&user.User.Name,
		&user.User.Email,
		&picture,
		&user.User.Role,
		&isActive,
		&user.User.IsEmailVerified,
		&createdAt,
		&updatedAt,
		&lastLoginAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		sr.logger.Log(logger.ErrorLevel, "Error fetching user: %v", err)
		return nil, fmt.Errorf("database error")
	}

	if picture.Valid {
		user.User.Picture = picture.String
	}

	user.User.CreatedAt = createdAt.Unix()
	user.User.UpdatedAt = updatedAt.Unix()
	user.User.LastLoginAt = lastLoginAt.Unix()
	return user, nil
}

func (sr *AuthRepository) Logout(ctx context.Context, req *request.LogoutRequest) (*response.LogoutResponse, error) {
	query := `
	DELETE FROM sessions
    WHERE access_token = $1 AND user_id = $2
	`
	_, err := sr.DB.ExecContext(ctx, query, req.AccessToken, req.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		sr.logger.Log(logger.ErrorLevel, "Error fetching user: %v", err)
		return nil, fmt.Errorf("database error")
	}

	return &response.LogoutResponse{
		Success: true,
	}, nil
}
