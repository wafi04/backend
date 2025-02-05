package types

import "time"

type Cart struct {
	CartID    string     `db:"cart_id" json:"cart_id"`
	UserID    string     `db:"user_id" json:"user_id"`
	Total     float64    `db:"total" json:"total"`
	Item      []CartItem `json:"cart_items"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type CartItem struct {
	CartItemID  string  `db:"cart_item_id" json:"cart_item_id"`
	CartID      string  `db:"cart_id" json:"cart_id"`
	VariantID   string  `db:"variant_id" json:"variant_id"`
	Size        string  `db:"size" json:"size"`
	Quantity    int64   `db:"quantity" json:"quantity"`
	SubTotal    float64 `db:"sub_total" json:"sub_total"`
	CreatedAt   string  `db:"created_at" json:"created_at"`
	UpdatedAt   string  `db:"updated_at" json:"updated_at"`
	ImageURL    *string `json:"image_url,omitempty"`
	Color       *string `json:"color,omitempty"`
	SKU         *string `json:"sku,omitempty"`
	ProductName *string `json:"product_name"`
}
