package productRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wafi04/backend/internal/handler/dto/request"
	"github.com/wafi04/backend/internal/handler/dto/response"
	"github.com/wafi04/backend/internal/handler/dto/types"
	"github.com/wafi04/backend/pkg/logger"
)

func (pr *Database) AddProductImage(ctx context.Context, req *request.AddProductImageRequest) (*types.ProductImage, error) {
	imageID := uuid.New().String()

	querySelect := `
        SELECT id FROM product_variants WHERE id = $1
    `
	var variantID string
	err := pr.DB.QueryRowContext(ctx, querySelect, req.VariantID).Scan(&variantID)
	if err != nil {
		pr.log.Log(logger.ErrorLevel, "Variant not found: %v", err)
		return nil, fmt.Errorf("variant not found: %v", err)
	}

	queryInsert := `
        INSERT INTO product_images (id, url, variant_id, is_main)
        VALUES ($1, $2, $3, $4)
        RETURNING id, url, variant_id, is_main
    `

	var image types.ProductImage
	err = pr.DB.QueryRowContext(ctx, queryInsert, imageID, req.URL, req.VariantID, req.IsMain).Scan(
		&image.ID,
		&image.URL,
		&image.VariantID,
		&image.IsMain,
	)
	if err != nil {
		pr.log.Log(logger.ErrorLevel, "Failed to create product image: %v", err)
		return nil, fmt.Errorf("failed to create product image: %v", err)
	}

	return &image, nil
}

func (pr *Database) UpdateProductImage(ctx context.Context, req *request.UpdateProductImageRequest) (*types.ProductImage, error) {
	querySelect := `
        SELECT id FROM product_images WHERE id = $1
    `
	var imageID string
	err := pr.DB.QueryRowContext(ctx, querySelect, req.Image.ID).Scan(&imageID)
	if err != nil {
		pr.log.Log(logger.ErrorLevel, "Image not found: %v", err)
		return nil, fmt.Errorf("image not found: %v", err)
	}

	queryUpdate := `
        UPDATE product_images
        SET url = $1, is_main = $2
        WHERE id = $3
        RETURNING id, url, variant_id, is_main
    `

	var image types.ProductImage
	err = pr.DB.QueryRowContext(ctx, queryUpdate, req.Image.URL, req.Image.IsMain, req.Image.ID).Scan(
		&image.ID,
		&image.URL,
		&image.VariantID,
		&image.IsMain,
	)
	if err != nil {
		pr.log.Log(logger.ErrorLevel, "Failed to update product image: %v", err)
		return nil, fmt.Errorf("failed to update product image: %v", err)
	}

	return &image, nil
}

func (pr *Database) DeleteProductImage(ctx context.Context, req *request.DeleteProductImageRequest) (*response.DeleteProductResponse, error) {
	queryDelete := `
        DELETE FROM product_images
        WHERE id = $1
    `

	result, err := pr.DB.ExecContext(ctx, queryDelete, req.ID)
	if err != nil {
		pr.log.Log(logger.ErrorLevel, "Failed to delete product image: %v", err)
		return nil, fmt.Errorf("failed to delete product image: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		pr.log.Log(logger.ErrorLevel, "Failed to get rows affected: %v", err)
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return &response.DeleteProductResponse{
			Success: false,
		}, nil
	}

	return &response.DeleteProductResponse{
		Success: true,
	}, nil
}
