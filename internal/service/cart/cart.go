package cartservice

import (
	"context"

	"github.com/wafi04/backend/internal/handler/dto/request"
	"github.com/wafi04/backend/internal/handler/dto/response"
	"github.com/wafi04/backend/internal/handler/dto/types"
	cartrepository "github.com/wafi04/backend/internal/repository/cart"
	"github.com/wafi04/backend/pkg/logger"
)

type CartService struct {
	cartrepo cartrepository.CartRepository
	log      logger.Logger
}

func NewCartService(cartRepo cartrepository.CartRepository) *CartService {
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
