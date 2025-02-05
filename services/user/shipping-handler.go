package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/backend/pkg/middleware"
	httpresponse "github.com/wafi04/backend/pkg/response"
	request "github.com/wafi04/backend/pkg/types/req"
	"github.com/wafi04/backend/pkg/utils"
)

type ShippingHandler struct {
	shippingRepo ShippingAddressRepo
}

func NewShippingHandler(shippingrepo ShippingAddressRepo) *ShippingHandler {
	return &ShippingHandler{
		shippingRepo: shippingrepo,
	}
}

func (h *ShippingHandler) CreateAddressReq(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		RecipientName  string  `json:"recipient_name" binding:"required"`
		Recipientphone string  `json:"recipient_phone" binding:"required"`
		FullAddress    string  `json:"full_address" binding:"required"`
		City           string  `json:"city" binding:"required"`
		Province       string  `json:"province" binding:"required"`
		PostalCode     string  `json:"postal_code" binding:"required"`
		Country        string  `json:"country" binding:"required"`
		Label          *string `json:"label,omitempty"`
		IsDefault      bool    `json:"is_default"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Invalid request: %s", err.Error())
		return
	}

	if len(req.Recipientphone) < 10 || len(req.Recipientphone) > 15 {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}
	addressID := utils.GenerateRandomId("SHIP-ADDR")

	createReq := &request.CreateAddressReq{
		UserID:         user.UserID,
		AddressID:      addressID,
		RecipientName:  req.RecipientName,
		Recipientphone: req.Recipientphone,
		FullAddress:    req.FullAddress,
		City:           req.City,
		Province:       req.Province,
		PostalCode:     req.PostalCode,
		Country:        req.Country,
		Label:          req.Label,
		IsDefault:      req.IsDefault,
	}

	data, err := h.shippingRepo.CreateShippingAddres(c, createReq)
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to create address", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusCreated, "Address created successfully", data)
}

func (h *ShippingHandler) GetAll(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	data, err := h.shippingRepo.GetShippingAddress(c, user.UserID)
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to get address: %s", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get Address successfully", data)
}

func (h *ShippingHandler) UpdateShipping(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	AddressID := c.Param("id")

	if AddressID == "" {
		httpresponse.SendErrorResponse(c, http.StatusNotFound, "Adddress id is required")
		return
	}

	var req struct {
		RecipientName  *string `json:"recipient_name,omitempty"`
		Recipientphone *string `json:"recipient_phone,omitempty"`
		FullAddress    *string `json:"full_address,omitempty"`
		City           *string `json:"city,omitempty"`
		Province       *string `json:"province,omitempty"`
		PostalCode     *string `json:"postal_code,omitempty"`
		Country        *string `json:"country,omitempty"`
		Label          *string `json:"label,omitempty"`
		IsDefault      *bool   `json:"is_default,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Invalid request: %s", err.Error())
		return
	}

	if len(*req.Recipientphone) < 10 || len(*req.Recipientphone) > 15 {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}
	updateReq := &request.UpdateAddressReq{
		UserID:         user.UserID,
		AddressID:      &AddressID,
		RecipientName:  req.RecipientName,
		Recipientphone: req.Recipientphone,
		FullAddress:    req.FullAddress,
		City:           req.City,
		Province:       req.Province,
		PostalCode:     req.PostalCode,
		Country:        req.Country,
		Label:          req.Label,
		IsDefault:      req.IsDefault,
	}

	update, err := h.shippingRepo.UpdateShippingAddress(c, updateReq)

	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to get address: %s", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Updated Address successfully", update)

}

func (h *ShippingHandler) UpdateDefault(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		IsDefault bool   `json:"is_boolean"`
		AddressID string `json:"address_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Invalid request: %s", err.Error())
		return
	}
	success, err := h.shippingRepo.SetDefaultAddress(c, user.UserID, req.AddressID)
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to update : %s", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get Shipping Address successfully ", success)
}

func (h *ShippingHandler) DeleteShipping(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ShippingID := c.Param("id")
	if ShippingID == "" {
		httpresponse.SendErrorResponse(c, http.StatusNotFound, "Shipping Id  Not Fround")
		return
	}

	delete, err := h.shippingRepo.DeleteShippingAddress(c, user.UserID, ShippingID)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Shipping Delete Successfully")
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusOK, "Delete Successfully ", delete)
}
