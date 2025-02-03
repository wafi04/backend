package inventoryService

import (
	"context"

	"github.com/wafi04/backend/internal/handler/dto/request"
	inventoryrepo "github.com/wafi04/backend/internal/repository/inventory"
	"github.com/wafi04/backend/pkg/logger"
)

type InventoryService struct {
	inventoryRepo inventoryrepo.InventoryRepository
	log           logger.Logger
}

func NewInventoryService(inventoryRepo inventoryrepo.InventoryRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
	}
}

func (s *InventoryService) GetInventoryByVariant(ctx context.Context, req *request.GetInventoryByVariantRequest) (*request.GetInventoryByVariantResponse, error) {
	s.log.Log(logger.DebugLevel, "Incoming Request From : %s", req.VariantID)
	return s.inventoryRepo.GetInventoryByVariant(ctx, req)
}
