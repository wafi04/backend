package producthandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/backend/pkg/logger"
	httpresponse "github.com/wafi04/backend/pkg/response"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
	"github.com/wafi04/backend/services/files"
	productservice "github.com/wafi04/backend/services/product/service"
)

type ProductHandler struct {
	productService *productservice.ProductService
	log            logger.Logger
	filesclient    *files.Cloudinary
}

func NewProductHandler(service *productservice.ProductService, files *files.Cloudinary) *ProductHandler {
	return &ProductHandler{
		productService: service,
		filesclient:    files,
	}
}

func (h *ProductHandler) HandleCreateProduct(c *gin.Context) {
	h.log.Log(logger.InfoLevel, "Content-Type: %s", c.Request.Header.Get("Content-Type"))

	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "error : %s", err.Error())
		return
	}

	resp, err := h.productService.CreateProduct(c, &req)

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Failed to create user")
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusCreated, "Created user successfully", resp)

}

func (h *ProductHandler) HandleGetProduct(c *gin.Context) {
	log.Printf("Received get product request: %s %s", c.Request.Method, c.Request.URL.Path)

	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	res, err := h.productService.GetProduct(c, &request.GetProductRequest{
		ID: id,
	})
	if err != nil {
		log.Printf("Failed to get Product: %v", err)
		httpresponse.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to get products", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get product success", res)
}

func (h *ProductHandler) HandleListProducts(c *gin.Context) {
	log.Printf("Received get product request: %s %s", c.Request.Method, c.Request.URL.Path)

	_ = c.Query("page")
	// limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	req := &request.ListProductsRequest{
		PageSize:  10,
		PageToken: "0",
	}

	res, err := h.productService.ListProducts(c, req)

	if err != nil {
		log.Printf("Failed to get Product: %v", err)
		httpresponse.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to get products", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get product successfuly ", res)
}

func (h *ProductHandler) HandleUpdateProduct(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusNotFound, "Product id not foung")
		return
	}

	log.Printf("Product id : %s", id)

	var req struct {
		Name        string  `json:"name"`
		SubTitle    string  `json:"sub_title"`
		Description string  `json:"description"`
		CategoryID  string  `json:"category_id"`
		Price       float64 `json:"price"`
		Sku         string  `json:"sku"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "error : %s", err.Error())
		return
	}

	update, err := h.productService.UpdateProduct(c, &request.UpdateProductRequest{
		Product: &types.Product{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
			SubTitle:    req.SubTitle,
			Price:       req.Price,
			CategoryID:  req.CategoryID,
		},
	})

	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to update product ", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Updated Product Success", update)

}

func (h *ProductHandler) HandleDeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusNotFound, "Product id not foung")
		return
	}

	delete, err := h.productService.DeleteProduct(c, &request.DeleteProductRequest{
		ID: id,
	})

	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to delete product ", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Delete Product Success", delete)

}
