package productservice

import (
	"context"
	"time"

	"github.com/wafi04/backend/internal/handler/dto/request"
	"github.com/wafi04/backend/internal/handler/dto/response"
	"github.com/wafi04/backend/internal/handler/dto/types"
	productRepository "github.com/wafi04/backend/internal/repository/product"
	"github.com/wafi04/backend/pkg/logger"
	"github.com/wafi04/backend/pkg/utils"
)

type ProductService struct {
	productrepo productRepository.ProductRepository
	log         logger.Logger
}

func NewProductService(productrepo productRepository.ProductRepository) *ProductService {
	return &ProductService{
		productrepo: productrepo,
	}
}

func (h *ProductService) CreateProduct(ctx context.Context, req *request.CreateProductRequest) (*types.Product, error) {
	id := utils.GenerateRandomId("PROD")
	sku := utils.GenerateSku(req.Name)

	h.log.Log(logger.InfoLevel, "incoming request ")
	return h.productrepo.CreateProduct(ctx, &types.Product{
		ID:          id,
		Name:        req.Name,
		SubTitle:    req.SubTitle,
		Description: req.Description,
		SKU:         sku,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	})
}

func (h *ProductService) GetProduct(ctx context.Context, req *request.GetProductRequest) (*types.Product, error) {
	h.log.Log(logger.InfoLevel, "incoming request ")
	return h.productrepo.GetProduct(ctx, req)

}
func (h *ProductService) ListProducts(ctx context.Context, req *request.ListProductsRequest) (*response.ListProductsResponse, error) {
	h.log.Log(logger.InfoLevel, "incoming request list")
	return h.productrepo.ListProducts(ctx, req)
}
func (h *ProductService) UpdateProduct(ctx context.Context, req *request.UpdateProductRequest) (*types.Product, error) {
	h.log.Log(logger.InfoLevel, "incoming request list")
	return h.productrepo.UpdateProduct(ctx, req)
}
func (h *ProductService) DeleteProduct(ctx context.Context, req *request.DeleteProductRequest) (*response.DeleteProductResponse, error) {
	h.log.Log(logger.InfoLevel, "incoming request list")
	return h.productrepo.DeleteProduct(ctx, req)
}

func (h *ProductService) CreateProductVariant(ctx context.Context, req *request.CreateProductVariantRequest) (*types.ProductVariant, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Create Varinat")
	return h.productrepo.CreateProductVariant(ctx, req)
}

func (h *ProductService) UpdateProductVariant(ctx context.Context, req *request.UpdateProductVariantRequest) (*types.ProductVariant, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Update Varinat")
	return h.productrepo.UpdateProductVariant(ctx, req)
}

func (h *ProductService) DeleteProductVariant(ctx context.Context, req *request.DeleteProductVariantRequest) (*response.DeleteProductResponse, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Delete Varinat")
	return h.productrepo.DeleteProductVariant(ctx, req)
}
func (h *ProductService) GetProductVariant(ctx context.Context, req *request.GetProductVariantRequest) (*types.ProductVariant, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Varinat  by id")
	return h.productrepo.GetProductsVariant(ctx, req)
}
func (h *ProductService) GetProductVariants(ctx context.Context, req *request.GetProductVariantsRequest) (*response.GetProductVariantsResponse, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Get all  Varinat")
	return h.productrepo.GetProductVariants(ctx, req)
}

func (h *ProductService) AddProductImage(ctx context.Context, req *request.AddProductImageRequest) (*types.ProductImage, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Create Image")
	return h.productrepo.AddProductImage(ctx, req)
}

func (h *ProductService) UpdateProductImage(ctx context.Context, req *request.UpdateProductImageRequest) (*types.ProductImage, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Update image")
	return h.productrepo.UpdateProductImage(ctx, req)
}

func (h *ProductService) DeleteProductImage(ctx context.Context, req *request.DeleteProductImageRequest) (*response.DeleteProductResponse, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Delete image")
	return h.productrepo.DeleteProductImage(ctx, req)
}
