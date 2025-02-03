package carthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/backend/internal/handler/dto/request"
	cartservice "github.com/wafi04/backend/internal/service/cart"
	"github.com/wafi04/backend/pkg/middleware"
	httpresponse "github.com/wafi04/backend/pkg/response"
)

type CartHandler struct {
	cartservice *cartservice.CartService
}

func NewCartHandler(service *cartservice.CartService) *CartHandler {
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
	}
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	cart, err := h.cartservice.AddToCart(c, &request.CartRequest{
		VariantID: req.VariantID,
		Size:      req.Size,
		Quantity:  req.Quantity,
		UserID:    user.UserID,
		Total:     req.Total,
	})

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
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

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get Cart Successfully", cart)
}
