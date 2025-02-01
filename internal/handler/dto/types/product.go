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
type ProductVariant struct {
	ID        string          `json:"id,omitempty"`
	Color     string          `json:"color,omitempty"`
	SKU       string          `json:"sku,omitempty"`
	Images    []*ProductImage `json:"images,omitempty"`
	ProductID string          `json:"product_id,omitempty"`
}

type ProductImage struct {
	ID        string `json:"id,omitempty"`
	URL       string `json:"url,omitempty"`
	VariantID string `json:"variant_id,omitempty"`
	IsMain    bool   `json:"is_main,omitempty"`
}
