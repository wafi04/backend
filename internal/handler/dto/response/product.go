package response

import "github.com/wafi04/backend/internal/handler/dto/types"

type DeleteProductResponse struct {
	Success bool `json:"success,omitempty"`
}

type ListProductsResponse struct {
	Products      []*types.Product `json:"products,omitempty"`
	NextPageToken string           `json:"next_page_token,omitempty"`
}

type GetProductVariantsResponse struct {
	Variants []*types.ProductVariant `json:"variants,omitempty"`
}
