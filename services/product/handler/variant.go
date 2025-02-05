package producthandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	httpresponse "github.com/wafi04/backend/pkg/response"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
	"github.com/wafi04/backend/pkg/utils"
)

func (h *ProductHandler) HandleCreateVariants(c *gin.Context) {
	var req struct {
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "error : %s", err.Error())
		return
	}

	var sku string
	if sku == "" {
		sku = utils.GenerateSku(req.Color)
	} else if !utils.IsSkuValid(sku) {
		return
	}

	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	variants, err := h.productService.CreateProductVariant(c, &request.CreateProductVariantRequest{
		ProductID: id,
		Color:     req.Color,
		SKU:       sku,
	})

	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to create variants", err.Error())
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusCreated, "Created variants successfully", variants)

}

func (p *ProductHandler) HandleUpdateVariants(c *gin.Context) {

	var req struct {
		Color     string `json:"color"`
		ProductID string `json:"product_id"`
		Sku       string `json:"sku"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "error : %s", err.Error())
		return
	}

	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	if req.Sku == "" {
		req.Sku = utils.GenerateSku(req.Color)
	} else if !utils.IsSkuValid(req.Sku) {
		return
	}

	update, err := p.productService.UpdateProductVariant(c, &request.UpdateProductVariantRequest{
		Variant: &types.ProductVariant{
			ID:        id,
			Color:     req.Color,
			SKU:       req.Sku,
			ProductID: req.ProductID,
		},
	})

	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to update variants", err.Error())
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusOK, "Update variants successfully", update)
}

func (p *ProductHandler) HandleDeleteVariants(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	delete, err := p.productService.DeleteProductVariant(c, &request.DeleteProductVariantRequest{
		ID: id,
	})
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to update variants", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Update variants successfully", delete)

}

func (p *ProductHandler) HandleGetProductVariant(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	log.Printf("Variant id : %s", id)

	variant, err := p.productService.GetProductVariant(c, &request.GetProductVariantRequest{
		VariantID: id,
	})
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to Get variant", err.Error())
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get Product variant successfully", variant)
}
func (p *ProductHandler) HandleGetProductVariants(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}
	log.Printf("Product id : %s", id)

	// Fetch product details
	product, err := p.productService.GetProduct(c, &request.GetProductRequest{
		ID: id,
	})
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to get product", err.Error())
		return
	}

	response := map[string]interface{}{
		"product": product,
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get product and variants successfully", response)
}
