package inventoryhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/backend/internal/handler/dto/request"
	inventoryService "github.com/wafi04/backend/internal/service/inventory"
	"github.com/wafi04/backend/pkg/logger"
	httpresponse "github.com/wafi04/backend/pkg/response"
)

type InventoryHandler struct {
	inventoryService *inventoryService.InventoryService
	log              logger.Logger
}

func NewInventoryHandler(service *inventoryService.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: service,
	}
}

func (h *InventoryHandler) HandleGetInvetory(c *gin.Context) {
	h.log.Log(logger.InfoLevel, "Content-Type: %s", c.Request.Header.Get("Content-Type"))

	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusNotFound, "Product id not foung")
		return
	}

	inv, err := h.inventoryService.GetInventoryByVariant(c, &request.GetInventoryByVariantRequest{
		VariantID: id,
	})

	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to get inventory ", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get Inventory Successfulyy", inv)
}
