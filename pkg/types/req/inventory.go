package request

import "github.com/wafi04/backend/pkg/types"

type GetInventoryByVariantRequest struct {
	VariantID string `json:"variant_id"`
}

type GetInventoryByVariantResponse struct {
	Inventory []types.Inventory `json:"inventory"`
}

type CreateInventoryRequest struct {
	VariantID string `json:"variant_id"`
	Size      string `json:"size"`
	Stock     int    `json:"stock"`
}

type UpdateInventoryRequest struct {
	ID             string `json:"id"`
	Size           string `json:"size"`
	Stock          int    `json:"stock"`
	ReservedStock  int    `json:"reserved_stock"`
	AvailableStock int    `json:"available_stock"`
}
