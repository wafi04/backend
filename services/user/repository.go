package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/backend/pkg/logger"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
)

type Database struct {
	db  *sqlx.DB
	log logger.Logger
}
type UserRepository interface {
	CreateUserDetails(ctx context.Context, req *request.ReqCreateUserDetails) (*types.UserDetails, error)
	UpdateUserDetails(ctx context.Context, req *request.ReqUpdateUserDetails) (*types.UserDetails, error)
	GetUserDetails(ctx context.Context, userID string) (*types.UserDetails, error)
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &Database{db: db}
}

func (d *Database) CreateUserDetails(ctx context.Context, req *request.ReqCreateUserDetails) (*types.UserDetails, error) {
	d.log.Log(logger.InfoLevel, "incoming request : %v", req)

	var userDetails types.UserDetails
	query := `
    INSERT INTO user_details 
    (
        user_id,
        place_birth,
        date_birth,
        gender,
        phone_number,
        bio,
        preferences,
        created_at,
        updated_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING 
        user_id,
        place_birth,
        date_birth,
        bio,
        gender,
        phone_number,
        preferences,
        created_at,
        updated_at
    `

	// Temporary variable to hold raw preferences data
	var rawPreferences []byte

	err := d.db.QueryRowContext(ctx, query,
		req.UserID,
		req.PlaceBirth,
		req.DateBirth,
		req.Gender,
		req.PhoneNumber,
		req.Bio,
		req.Preferences,
		time.Now(),
		time.Now(),
	).Scan(
		&userDetails.UserID,
		&userDetails.PlaceBirth,
		&userDetails.DateBirth,
		&userDetails.Bio,
		&userDetails.Gender,
		&userDetails.PhoneNumber,
		&rawPreferences,
		&userDetails.CreatedAt,
		&userDetails.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user details : %v", err)
	}

	if err := json.Unmarshal(rawPreferences, &userDetails.Preferences); err != nil {
		return nil, fmt.Errorf("failed to parse preferences : %v", err)
	}

	return &userDetails, nil
}
func (d *Database) UpdateUserDetails(ctx context.Context, req *request.ReqUpdateUserDetails) (*types.UserDetails, error) {
	var userDetails types.UserDetails

	placeBirth := sql.NullString{String: "", Valid: false}
	if req.PlaceBirth != nil {
		placeBirth = sql.NullString{String: *req.PlaceBirth, Valid: true}
	}

	dateBirth := sql.NullTime{Time: time.Time{}, Valid: false}
	if req.DateBirth != nil {
		dateBirth = sql.NullTime{Time: *req.DateBirth, Valid: true}
	}

	gender := sql.NullString{String: "", Valid: false}
	if req.Gender != nil {
		gender = sql.NullString{String: *req.Gender, Valid: true}
	}

	phoneNumber := sql.NullString{String: "", Valid: false}
	if req.PhoneNumber != nil {
		phoneNumber = sql.NullString{String: *req.PhoneNumber, Valid: true}
	}

	bio := sql.NullString{String: "", Valid: false}
	if req.Bio != nil {
		bio = sql.NullString{String: *req.Bio, Valid: true}
	}

	preferences := sql.NullString{String: "", Valid: false}
	if req.Preferences != nil {
		jsonData, err := json.Marshal(req.Preferences)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal preferences: %v", err)
		}
		preferences = sql.NullString{String: string(jsonData), Valid: true}
	}

	query := `
        UPDATE user_details
        SET
            place_birth = COALESCE($1, place_birth),
            date_birth = COALESCE($2, date_birth),
            gender = COALESCE($3, gender),
            phone_number = COALESCE($4, phone_number),
            bio = COALESCE($5, bio),
            preferences = COALESCE($6::jsonb, preferences),
            updated_at = $7
        WHERE user_id = $8
        RETURNING
            user_id,
            place_birth,
            date_birth,
            gender,
            phone_number,
            bio,
            preferences,
            created_at,
            updated_at
    `
	var rawPreferences []byte

	err := d.db.QueryRowContext(ctx, query,
		placeBirth,
		dateBirth,
		gender,
		phoneNumber,
		bio,
		preferences,
		time.Now(),
		req.UserID,
	).Scan(
		&userDetails.UserID,
		&userDetails.PlaceBirth,
		&userDetails.DateBirth,
		&userDetails.Gender,
		&userDetails.PhoneNumber,
		&userDetails.Bio,
		&rawPreferences,
		&userDetails.CreatedAt,
		&userDetails.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user details: %v", err)
	}

	if rawPreferences != nil {
		if err := json.Unmarshal(rawPreferences, &userDetails.Preferences); err != nil {
			return nil, fmt.Errorf("failed to parse preferences: %v", err)
		}
	} else {
		userDetails.Preferences = types.Preferences{}
	}

	return &userDetails, nil
}
func (d *Database) GetUserDetails(ctx context.Context, userID string) (*types.UserDetails, error) {
	var userDetails types.UserDetails

	// Define the query
	query := `
    SELECT 
        user_id,
        place_birth,
        date_birth,
        gender,
        phone_number,
        bio,
        preferences,
        created_at,
        updated_at
    FROM user_details
    WHERE user_id = $1
    `

	var rawPreferences []byte

	err := d.db.QueryRowContext(ctx, query, userID).Scan(
		&userDetails.UserID,
		&userDetails.PlaceBirth,
		&userDetails.DateBirth,
		&userDetails.Gender,
		&userDetails.PhoneNumber,
		&userDetails.Bio,
		&rawPreferences,
		&userDetails.CreatedAt,
		&userDetails.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user details: %v", err)
	}

	if rawPreferences != nil {
		if err := json.Unmarshal(rawPreferences, &userDetails.Preferences); err != nil {
			return nil, fmt.Errorf("failed to parse preferences: %v", err)
		}
	} else {
		userDetails.Preferences = types.Preferences{}
	}

	return &userDetails, nil
}
