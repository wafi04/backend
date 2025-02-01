package producthandler

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/backend/internal/handler/dto/request"
	"github.com/wafi04/backend/pkg/logger"
	httpresponse "github.com/wafi04/backend/pkg/response"
)

func (h *ProductHandler) HandleAddProductImage(c *gin.Context) {
	maxSize := int64(10 << 20)
	if err := c.Request.ParseMultipartForm(maxSize); err != nil {
		h.log.Log(logger.ErrorLevel, "Error parsing multipart form: %v", err)
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to parse form data : %s", err.Error())
		return
	}
	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}

	file, _, err := c.Request.FormFile("image")
	var imageUrl string
	if err == nil {
		defer file.Close()

		PublicID := fmt.Sprintf("%06d", rand.Intn(1000000))
		uploadRequest := &request.FileUploadRequest{
			FileData: file,
			Folder:   "products",
			PublicID: PublicID,
		}

		uploadResponse, err := h.filesclient.UploadFile(c, uploadRequest)
		if err != nil {
			httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Failed to upload image")
			return
		}

		imageUrl = uploadResponse.URL
	}

	var imageUrlPtr *string
	if imageUrl != "" {
		imageUrlPtr = &imageUrl
	}

	productImage, err := h.productService.AddProductImage(c, &request.AddProductImageRequest{
		VariantID: id,
		URL:       *imageUrlPtr,
		IsMain:    true,
	})
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Failed to created product image")
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusCreated, "Created Product Image succesfully", productImage)
}

func (p *ProductHandler) HandleDeleteProductImage(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}
	productImage, err := p.productService.DeleteProductImage(c, &request.DeleteProductImageRequest{
		ID: id,
	})

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Failed to Delete product image")
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusOK, "Delete Product Image succesfully", productImage)
}
