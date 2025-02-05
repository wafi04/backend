package request

import (
	"time"

	"github.com/wafi04/backend/pkg/types"
)

type ReqCreateUserDetails struct {
	UserID      string            `json:"user_id" db:"user_id"`
	PlaceBirth  string            `json:"place_birth" db:"place_birth"`
	DateBirth   *time.Time        `json:"date_birth,omitempty" db:"date_birth"`
	Gender      string            `json:"gender" db:"gender"`
	PhoneNumber string            `json:"phone_number" db:"phone_number"`
	Bio         string            `json:"bio" db:"bio"`
	Preferences types.Preferences `json:"preferences" db:"preferences,type:JSONB"`
}

type ReqUpdateUserDetails struct {
	UserID      string             `json:"user_id" db:"user_id"`
	PlaceBirth  *string            `json:"place_birth,omitempty" db:"place_birth"`
	DateBirth   *time.Time         `json:"date_birth,omitempty" db:"date_birth"`
	Gender      *string            `json:"gender,omitempty" db:"gender"`
	PhoneNumber *string            `json:"phone_number,omitempty" db:"phone_number"`
	Bio         *string            `json:"bio,omitempty" db:"bio"`
	Preferences *types.Preferences `json:"preferences,omitempty" db:"preferences,type:JSONB"`
}

type CreateAddressReq struct {
	UserID         string  `json:"user_id"`
	AddressID      string  `json:"address_id"`
	RecipientName  string  `json:"recipient_name"`
	Recipientphone string  `json:"recipient_phone"`
	FullAddress    string  `json:"full_address"`
	City           string  `json:"city"`
	Province       string  `json:"province"`
	PostalCode     string  `json:"postal_code"`
	Country        string  `json:"country"`
	Label          *string `json:"label,omitempty"`
	IsDefault      bool    `json:"is_default"`
}

type UpdateAddressReq struct {
	UserID         string  `json:"user_id"`
	AddressID      *string `json:"address_id,omitempty"`
	RecipientName  *string `json:"recipient_name,omitempty"`
	Recipientphone *string `json:"recipient_phone,omitempty"`
	FullAddress    *string `json:"full_address,omitempty"`
	City           *string `json:"city,omitempty"`
	Province       *string `json:"province,omitempty"`
	PostalCode     *string `json:"postal_code,omitempty"`
	Country        *string `json:"country,omitempty"`
	Label          *string `json:"label,omitempty"`
	IsDefault      *bool   `json:"is_default,omitempty"`
}
