package authrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wafi04/backend/internal/handler/dto/request"
	"github.com/wafi04/backend/internal/handler/dto/response"
	"github.com/wafi04/backend/internal/handler/dto/types"
	"github.com/wafi04/backend/pkg/logger"
	"github.com/wafi04/backend/pkg/middleware"
)

func (D *AuthRepository) RevokeSession(ctx context.Context, req *request.RevokeSessionRequest) (*response.RevokeSessionResponse, error) {
	D.logger.Log(logger.InfoLevel, "Recieved  Session Request ")

	query := `
	DELETE FROM sessions
    WHERE session_id = $1 AND user_id = $2
	`
	_, err := D.DB.ExecContext(ctx, query, req.SessionID, req.UserID)

	if err != nil {
		D.logger.Log(logger.ErrorLevel, "Failed to Delete Session : %v", err)
		return nil, nil
	}

	return &response.RevokeSessionResponse{
		Success: true}, nil
}
func (D *AuthRepository) CreateSession(ctx context.Context, session *types.Session) error {
	insertQuery := `
       INSERT INTO sessions (
           session_id, 
           user_id, 
           access_token, 
           refresh_token, 
           ip_address, 
           device_info, 
           is_active, 
           expires_at, 
           last_activity_at, 
           created_at
       ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
   `

	updateQuery := `
       UPDATE sessions 
       SET 
           access_token = $1, 
           refresh_token = $2, 
           ip_address = $3, 
           last_activity_at = $4
       WHERE user_id = $5 AND device_info = $6
   `

	if session.SessionID == "" {
		session.SessionID = uuid.New().String()
	}

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	_, err := D.DB.ExecContext(
		ctx,
		insertQuery,
		session.SessionID,
		session.UserID,
		session.AccessToken,
		session.RefreshToken,
		session.IPAddress,
		session.DeviceInfo,
		true,
		expiresAt,
		now,
		now,
	)

	if err != nil {
		_, err = D.DB.ExecContext(
			ctx,
			updateQuery,
			session.AccessToken,
			session.RefreshToken,
			session.IPAddress,
			now,
			session.UserID,
			session.DeviceInfo,
		)
	}

	if err != nil {
		D.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    session.UserID,
			"session_id": session.SessionID,
		}).Error("Failed to create/update session")
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil

}
func (D *AuthRepository) RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.RefreshTokenResponse, error) {
	D.logger.Log(logger.InfoLevel, "Refresh Token Incoming")

	var storedRefreshToken string
	query := `
        SELECT refresh_token FROM sessions WHERE session_id = $1
    `
	err := D.DB.QueryRowContext(ctx, query, req.SessionID).Scan(&storedRefreshToken)
	if err != nil {
		D.logger.Log(logger.ErrorLevel, "Failed to retrieve refresh token: %v", err)
		return nil, err
	}

	if storedRefreshToken != req.RefreshToken {
		D.logger.Log(logger.ErrorLevel, "Refresh token mismatch")
		return nil, errors.New("invalid refresh token")
	}

	query = `
        SELECT u.user_id, u.email, u.role, u.is_email_verified
        FROM sessions s
        JOIN users u ON s.user_id = u.user_id
        WHERE s.session_id = $1
    `
	var user types.UserInfo
	err = D.DB.QueryRowContext(ctx, query, req.SessionID).Scan(
		&user.UserID,
		&user.Email,
		&user.Role,
		&user.IsEmailVerified,
	)
	if err != nil {
		D.logger.Log(logger.ErrorLevel, "Failed to retrieve user: %v", err)
		return nil, err
	}

	accces_token, err := middleware.GenerateToken(&user, 24)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	updateQuery := `
        UPDATE sessions SET access_token = $1 WHERE session_id = $2
    `
	_, err = D.DB.ExecContext(ctx, updateQuery, accces_token, req.SessionID)
	if err != nil {
		D.logger.Log(logger.ErrorLevel, "Failed to update session: %v", err)
		return nil, err
	}

	return &response.RefreshTokenResponse{
		AccessToken:  accces_token,
		RefreshToken: req.RefreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func (D *AuthRepository) ListSessions(ctx context.Context, req *request.ListSessionRequest) (*response.ListSessionResponse, error) {
	query := `
        SELECT 
            session_id,
            device_info,
            ip_address,
            EXTRACT(EPOCH FROM created_at)::bigint AS created_at,
            EXTRACT(EPOCH FROM last_activity_at)::bigint AS last_activity_at
        FROM sessions
        WHERE user_id = $1
    `

	rows, err := D.DB.QueryContext(ctx, query, req.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*types.SessionInfo
	for rows.Next() {
		session := &types.SessionInfo{}
		err := rows.Scan(
			&session.SessionID,
			&session.DeviceInfo,
			&session.IPAddress,
			&session.CreatedAt,
			&session.LastActivityAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &response.ListSessionResponse{
		Sessions: sessions,
	}, nil
}

type ReqSession struct {
	SessionID string `json:"session_id"`
}

func (d *AuthRepository) GetCurrentSession(ctx context.Context, req *ReqSession) (*types.Session, error) {
	d.logger.Log(logger.InfoLevel, "session id : %s ", req.SessionID)

	var last_activity_at, expires_at, created_at time.Time
	var sessions types.Session
	query := `
	SELECT session_id, 
           user_id, 
           access_token, 
           refresh_token, 
           ip_address, 
           device_info, 
           is_active, 
           expires_at, 
           last_activity_at, 
           created_at
	FROM sessions
	WHERE session_id = $1
	`

	err := d.DB.QueryRowContext(ctx, query, req.SessionID).Scan(
		&sessions.SessionID,
		&sessions.UserID,
		&sessions.AccessToken,
		&sessions.RefreshToken,
		&sessions.IPAddress,
		&sessions.DeviceInfo,
		&sessions.IsActive,
		&expires_at,
		&last_activity_at,
		&created_at,
	)

	if err != nil {
		d.logger.Log(logger.ErrorLevel, "session id not found :%s", err)
		return nil, err
	}

	sessions.CreatedAt = created_at.Unix()
	sessions.LastActivityAt = created_at.Unix()
	sessions.ExpiresAt = expires_at.Unix()

	return &sessions, nil
}

// func (d *AuthRepository) UpdateSessionActivity(ctx context.Context, sessionID string) (bool, error) {
// 	query := `
//         UPDATE sessions
//         SET last_activity_at = CURRENT_TIMESTAMP
//         WHERE session_id = $1
//     `

// 	result, err := d.DB.ExecContext(ctx, query, sessionID)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to update session activity: %w", err)
// 	}

// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return false, err
// 	}

// 	return rows > 0, nil
// }
