package categoryhandler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/backend/internal/handler/dto/request"
	service "github.com/wafi04/backend/internal/service/category"
	"github.com/wafi04/backend/internal/service/files"
	"github.com/wafi04/backend/pkg/logger"
	httpresponse "github.com/wafi04/backend/pkg/response"
	"github.com/wafi04/backend/pkg/utils"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
	log             logger.Logger
	filesclient     *files.Cloudinary
}

func NewCategoryHandler(service *service.CategoryService, files *files.Cloudinary) *CategoryHandler {
	return &CategoryHandler{
		categoryService: service,
		filesclient:     files,
	}
}

func (h *CategoryHandler) HandleCreateCategory(c *gin.Context) {
	h.log.Log(logger.InfoLevel, "Content-Type: %s", c.Request.Header.Get("Content-Type"))

	maxSize := int64(10 << 20)
	if err := c.Request.ParseMultipartForm(maxSize); err != nil {
		h.log.Log(logger.ErrorLevel, "Error parsing multipart form: %v", err)
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to parse form data : %s", err.Error())
		return
	}

	// Log semua form values
	h.log.Log(logger.InfoLevel, "All form values: %+v", c.Request)

	name := c.PostForm("name")
	description := c.PostForm("description")
	parentID := c.PostForm("parent_id")

	if name == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Name is Required")
		return
	}

	if description == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Description is Required")
		return
	}

	file, _, err := c.Request.FormFile("file")
	var imageUrl string
	if err == nil {
		defer file.Close()
		PublicID := utils.GenerateRandomId("CAT")
		uploadRequest := &request.FileUploadRequest{
			FileData: file,
			Folder:   "categories",
			PublicID: PublicID,
		}

		uploadResponse, err := h.filesclient.UploadFile(c, uploadRequest)
		if err != nil {
			httpresponse.SendErrorResponse(c, http.StatusInternalServerError, "Failed to read file")
			return
		}

		imageUrl = uploadResponse.URL
	}

	var parentIDPtr *string
	if parentID != "" {
		parentIDPtr = &parentID
	}

	var imageUrlPtr *string
	if imageUrl != "" {
		imageUrlPtr = &imageUrl
	}

	resp, err := h.categoryService.CreateCategory(c, &request.CreateCategoryRequest{
		Name:        name,
		Description: description,
		Image:       imageUrlPtr,
		ParentID:    parentIDPtr,
	})

	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadGateway, "Failed To Create Category : %s", err.Error())
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusCreated, "Succes to create category", resp)
}

func (h *CategoryHandler) HandleUpdateCategory(c *gin.Context) {
	h.log.Log(logger.InfoLevel, "Content-Type: %s", c.Request.Header.Get("Content-Type"))

	id := c.Param("id")

	// Validate ID
	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}
	maxSize := int64(10 << 20)
	if err := c.Request.ParseMultipartForm(maxSize); err != nil {
		h.log.Log(logger.ErrorLevel, "Error parsing multipart form: %v", err)
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to parse form data : %s", err.Error())
		return
	}

	h.log.Log(logger.InfoLevel, "All form values: %+v", c.Request)

	name := c.PostForm("name")
	description := c.PostForm("description")
	parentID := c.PostForm("parent_id")

	if name == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Name is Required")
		return
	}

	if description == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Description is Required")
		return
	}

	file, _, err := c.Request.FormFile("file")
	var imageUrl string
	if err == nil {
		defer file.Close()
		PublicID := utils.GenerateRandomId("CAT")
		uploadRequest := &request.FileUploadRequest{
			FileData: file,
			Folder:   "categories",
			PublicID: PublicID,
		}

		uploadResponse, err := h.filesclient.UploadFile(c, uploadRequest)
		if err != nil {
			httpresponse.SendErrorResponse(c, http.StatusInternalServerError, "Failed to read file")
			return
		}

		imageUrl = uploadResponse.URL
	}

	var parentIDPtr *string
	if parentID != "" {
		parentIDPtr = &parentID
	}

	var imageUrlPtr *string
	if imageUrl != "" {
		imageUrlPtr = &imageUrl
	}

	resp, err := h.categoryService.UppdateCategory(c, &request.UpdateCategoryRequest{
		ID:          id,
		Name:        &name,
		Description: &description,
		Image:       imageUrlPtr,
		ParentID:    parentIDPtr,
	})

	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadGateway, "Failed To Create Category : %s", err.Error())
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusCreated, "Succes to create category", resp)
}

func (h *CategoryHandler) HandleGetCategory(c *gin.Context) {
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	parentID := c.Request.URL.Query().Get("parent_id")
	includeChildren := c.Request.URL.Query().Get("include_children") == "true"

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	req := &request.ListCategoriesRequest{
		Page:            int32(page),
		Limit:           int32(limit),
		IncludeChildren: includeChildren,
	}

	if parentID != "" {
		req.ParentID = &parentID
	}

	resp, err := h.categoryService.GetCategories(c, req)
	if err != nil {
		h.log.Log(logger.ErrorLevel, "Error calling ListCategories: %v", err)
		httpresponse.SendErrorResponse(c, http.StatusInternalServerError, "Error retrieving categories")
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Category  Retrieved Successfully", resp)
}

func (h *CategoryHandler) HandleDeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// Alternative ways to get parameters:
	// Query parameter (?id=123): c.Query("id")
	// Optional query parameter with default: c.DefaultQuery("id", "default_value")
	// Check if query exists: id, exists := c.GetQuery("id")

	// Validate ID
	if id == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}
	updateReq := &request.DeleteCategoryRequest{
		ID:             id,
		DeleteChildren: false,
	}

	category, err := h.categoryService.DeleteCategory(c, updateReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			httpresponse.SendErrorResponse(c, http.StatusNotFound, "Category Not Found")
		case strings.Contains(err.Error(), "invalid"):
			httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Invalid Request")
		default:
			httpresponse.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")

		}
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Delete Category Succesfully", category)
}
