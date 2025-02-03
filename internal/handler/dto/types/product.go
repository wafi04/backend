package types

type Product struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	SubTitle    string            `json:"sub_title"`
	Description string            `json:"description"`
	SKU         string            `json:"sku"`
	Price       float64           `json:"price"`
	Variants    []*ProductVariant `json:"variants,omitempty"`
	CategoryID  string            `json:"category_id"`
	CreatedAt   int64             `json:"created_at,omitempty"`
	UpdatedAt   int64             `json:"updated_at,omitempty"`
}
type Inventory struct {
	VariantID      string `json:"variant_id" db:"variant_id"`
	ID             string `json:"id" db:"id"`
	Size           string `json:"size" db:"size"`
	Stock          int    `json:"stock" db:"stock"`
	ReservedStock  int    `json:"reserved_stock" db:"reserved_stock"`
	AvailableStock int    `json:"available_stock" db:"available_stock"`
	CreatedAt      int64  `json:"created_at" db:"created_at"`
	UpdatedAt      int64  `json:"updated_at" db:"updated_at"`
}
type ProductVariant struct {
	ID        string          `json:"id,omitempty"`
	Color     string          `json:"color,omitempty"`
	SKU       string          `json:"sku,omitempty"`
	Images    []*ProductImage `json:"images,omitempty"`
	Inventory []*Inventory    `json:"inventory,omitempty"`
	ProductID string          `json:"product_id,omitempty"`
}

type ProductImage struct {
	ID        string `json:"id,omitempty"`
	URL       string `json:"url,omitempty"`
	VariantID string `json:"variant_id,omitempty"`
	IsMain    bool   `json:"is_main,omitempty"`
}
