package inventory

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/backend/pkg/logger"
	httpresponse "github.com/wafi04/backend/pkg/response"
	request "github.com/wafi04/backend/pkg/types/req"
)

type InventoryHandler struct {
	inventoryService *InventoryService
	log              logger.Logger
}

func NewInventoryHandler(service *InventoryService) *InventoryHandler {
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
