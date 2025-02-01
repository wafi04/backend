package request

import "github.com/wafi04/backend/internal/handler/dto/types"

type CreateProductRequest struct {
	Name        string  `json:"name"`
	SubTitle    string  `json:"sub_title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`
}

type GetProductRequest struct {
	ID string `json:"id"`
}

type UpdateProductRequest struct {
	Product *types.Product `json:"product,omitempty"`
}

type DeleteProductRequest struct {
	ID string `json:"id,omitempty"`
}

type ListProductsRequest struct {
	PageSize  int32  `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}

type CreateProductVariantRequest struct {
	ProductID string `json:"product_id,omitempty"`
	Color     string `json:"color,omitempty"`
	SKU       string `json:"sku,omitempty"`
}
type GetProductVariantRequest struct {
	VariantID string `json:"variant_id,omitempty"`
}
type GetProductVariantsRequest struct {
	ProductID string `json:"product_id,omitempty"`
}

type UpdateProductVariantRequest struct {
	Variant *types.ProductVariant `json:"variant,omitempty"`
}

type DeleteProductVariantRequest struct {
	ID string `json:"id,omitempty"`
}
type AddProductImageRequest struct {
	VariantID string `json:"variant_id,omitempty"`
	URL       string `json:"url,omitempty"`
	IsMain    bool   `json:"is_main,omitempty"`
}

type UpdateProductImageRequest struct {
	Image *types.ProductImage `json:"image,omitempty"`
}

type DeleteProductImageRequest struct {
	ID string `json:"id,omitempty"`
}
