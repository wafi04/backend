package productRepository

import (
	"context"

	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
	response "github.com/wafi04/backend/pkg/types/res"
)

type ProductRepository interface {
	// product
	CreateProduct(ctx context.Context, req *types.Product) (*types.Product, error)
	GetProduct(ctx context.Context, req *request.GetProductRequest) (*types.Product, error)
	ListProducts(ctx context.Context, req *request.ListProductsRequest) (*response.ListProductsResponse, error)
	UpdateProduct(ctx context.Context, req *request.UpdateProductRequest) (*types.Product, error)
	DeleteProduct(ctx context.Context, req *request.DeleteProductRequest) (*response.DeleteProductResponse, error)

	// variant
	CreateProductVariant(ctx context.Context, req *request.CreateProductVariantRequest) (*types.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, req *request.UpdateProductVariantRequest) (*types.ProductVariant, error)
	GetProductsVariant(ctx context.Context, req *request.GetProductVariantRequest) (*types.ProductVariant, error)
	DeleteProductVariant(ctx context.Context, req *request.DeleteProductVariantRequest) (*response.DeleteProductResponse, error)
	GetProductVariants(ctx context.Context, req *request.GetProductVariantsRequest) (*response.GetProductVariantsResponse, error)
	// images
	AddProductImage(ctx context.Context, req *request.AddProductImageRequest) (*types.ProductImage, error)
	UpdateProductImage(ctx context.Context, req *request.UpdateProductImageRequest) (*types.ProductImage, error)
	DeleteProductImage(ctx context.Context, req *request.DeleteProductImageRequest) (*response.DeleteProductResponse, error)
}
