package cart

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wafi04/backend/pkg/logger"
	"github.com/wafi04/backend/pkg/middleware"
	httpresponse "github.com/wafi04/backend/pkg/response"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
)

type CartHandler struct {
	cartservice *CartService
}

func NewCartHandler(service *CartService) *CartHandler {
	return &CartHandler{
		cartservice: service,
	}
}

func (h *CartHandler) HandleAddToCart(c *gin.Context) {
	var req struct {
		VariantID string  `json:"variant_id"`
		Size      string  `json:"size"`
		Quantity  int64   `json:"quantity"`
		Total     float64 `json:"total"`
		Price     float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "error : %s", err.Error())
		return
	}
	user, err := middleware.GetUserFromGinContext(c)
	h.cartservice.log.Log(logger.InfoLevel, "user : %s", user.UserID)
	if err != nil {
		h.cartservice.log.Log(logger.InfoLevel, "user : %v", err)
		httpresponse.SendErrorResponseWithDetails(c, http.StatusUnauthorized, "Unauthorized : %w", err.Error())
		return
	}
	cart, err := h.cartservice.AddToCart(c, &request.CartRequest{
		VariantID: req.VariantID,
		Size:      req.Size,
		Quantity:  req.Quantity,
		UserID:    user.UserID,
		Total:     req.Total,
		Price:     req.Price,
	})

	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to create cart :", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusCreated, "Created Cart Success", cart)
}

func (h *CartHandler) HandleGetCart(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cart, err := h.cartservice.GetCart(c, user.UserID)
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to get Cart", err.Error())
		return
	}
	h.cartservice.log.Log(logger.InfoLevel, "total : %f", cart.Total)
	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get Cart Successfully", cart)
}

func (h *CartHandler) HandleGetCountCart(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cart, err := h.cartservice.GetCart(c, user.UserID)
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to get Cart", err.Error())
		return
	}

	cartItems, err := h.cartservice.GetCountCartItem(c, cart.CartID)
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to get Cart", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get Product Successfully", cartItems)
}

func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cartItemID := c.Param("id")
	if cartItemID == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "cart item id is required")
		return
	}

	delete, err := h.cartservice.RemoveCart(c, &request.ReqRemoveCartByID{
		CartItemID: cartItemID,
		UserID:     user.UserID,
	})

	if err != nil {
		h.cartservice.log.Log(logger.InfoLevel, "err : %v", err)
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Failed to remove cart Items ")
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Remove Cart item successfully", delete)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}
	clear, err := h.cartservice.ClearCart(c, &request.ClearCart{
		UserID: user.UserID,
	})

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "failed to clear cart")
		return
	}
	httpresponse.SendSuccessResponse(c, http.StatusOK, "Clear Cart Succesfully ", clear)
}

func (h *CartHandler) UpdateQuantity(c *gin.Context) {
	cartItemID := c.Param("id")
	if cartItemID == "" {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "cart item id is required")
		return
	}

	var req struct {
		Size     string `json:"size"`
		Quantity int64  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "error : %s", err.Error())
		return
	}

	update, err := h.cartservice.UpdateQuantity(c, &request.UpdateQuantity{
		CartItemID: cartItemID,
		Size:       req.Size,
		Quantity:   req.Quantity,
	})

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusBadRequest, "Failed To update quantity")
		return
	}

	types.Broadcast <- fmt.Sprintf("Updated quantity for cart item %s: %d", cartItemID, req.Quantity)

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Updated Quantity suceessfully", update)
}
