package cart

import (
	"context"

	"github.com/wafi04/backend/pkg/logger"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
	response "github.com/wafi04/backend/pkg/types/res"
)

type CartService struct {
	cartrepo CartRepository
	log      logger.Logger
}

func NewCartService(cartRepo CartRepository) *CartService {
	return &CartService{
		cartrepo: cartRepo,
	}
}
func (s *CartService) AddToCart(ctx context.Context, req *request.CartRequest) (*response.CartResponse, error) {
	s.log.Log(logger.DebugLevel, "Incoming request on service")
	return s.cartrepo.AddCart(ctx, req)
}

func (s *CartService) RemoveCart(ctx context.Context, req *request.ReqRemoveCartByID) (*response.ResRemoveCartItem, error) {
	s.log.Log(logger.DebugLevel, "Incoming request on service")
	return s.cartrepo.RemoveFromCart(ctx, req)
}

func (s *CartService) ClearCart(ctx context.Context, req *request.ClearCart) (*response.ResRemoveCartItem, error) {
	return s.cartrepo.ClearCart(ctx, req)
}

func (s *CartService) UpdateQuantity(ctx context.Context, req *request.UpdateQuantity) (*types.CartItem, error) {
	return s.cartrepo.UpdateQuantity(ctx, req)
}

func (s *CartService) GetCart(ctx context.Context, req string) (*types.Cart, error) {
	return s.cartrepo.GetCart(ctx, req)
}

func (s *CartService) GetCountCartItem(ctx context.Context, req string) (int, error) {
	return s.cartrepo.GetCartItemCount(ctx, req)
}
