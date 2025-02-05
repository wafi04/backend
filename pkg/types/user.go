package types

import "time"

type UserDetails struct {
	UserID      string      `json:"user_id" db:"user_id"`
	PlaceBirth  string      `json:"place_birth" db:"place_birth"`
	DateBirth   *time.Time  `json:"date_birth,omitempty" db:"date_birth"`
	Gender      string      `json:"gender" db:"gender"`
	PhoneNumber string      `json:"phone_number" db:"phone_number"`
	Bio         string      `json:"bio" db:"bio"`
	Preferences Preferences `json:"preferences" db:"preferences,type:JSONB"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

type Preferences struct {
	Theme         string `json:"theme"`
	Notifications bool   `json:"notifications"`
}

type ShippingAddress struct {
	UserID         string    `json:"user_id"`
	AddressID      string    `json:"address_id"`
	RecipientName  string    `json:"recipient_name"`
	Recipientphone string    `json:"recipient_phone"`
	FullAddress    string    `json:"full_address"`
	City           string    `json:"city"`
	Province       string    `json:"province"`
	PostalCode     string    `json:"postal_code"`
	Country        string    `json:"country"`
	Label          *string   `json:"label,omitempty"`
	IsDefault      bool      `json:"is_default"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ListShippingAddress struct {
	Address []*ShippingAddress `json:"address"`
}
