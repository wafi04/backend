package inventory

import (
	"context"

	"github.com/wafi04/backend/pkg/logger"
	request "github.com/wafi04/backend/pkg/types/req"
)

type InventoryService struct {
	inventoryRepo InventoryRepository
	log           logger.Logger
}

func NewInventoryService(inventoryRepo InventoryRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
	}
}

func (s *InventoryService) GetInventoryByVariant(ctx context.Context, req *request.GetInventoryByVariantRequest) (*request.GetInventoryByVariantResponse, error) {
	s.log.Log(logger.DebugLevel, "Incoming Request From : %s", req.VariantID)
	return s.inventoryRepo.GetInventoryByVariant(ctx, req)
}
