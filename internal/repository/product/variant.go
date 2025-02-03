package productRepository

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/wafi04/backend/internal/handler/dto/request"
	"github.com/wafi04/backend/internal/handler/dto/response"
	"github.com/wafi04/backend/internal/handler/dto/types"
	"github.com/wafi04/backend/pkg/logger"
)

func (pr *Database) CreateProductVariant(ctx context.Context, req *request.CreateProductVariantRequest) (*types.ProductVariant, error) {
	variantsID := uuid.New().String()
	var variants types.ProductVariant
	query := `
		INSERT INTO product_variants (id,color,sku,product_id)
		VALUES ($1,$2,$3,$4)
		RETURNING id, color, sku, product_id
	`

	err := pr.DB.QueryRowContext(ctx, query, variantsID, req.Color, req.SKU, req.ProductID).Scan(
		&variants.ID,
		&variants.Color,
		&variants.SKU,
		&variants.ProductID,
	)

	if err != nil {
		pr.log.Error("Failed to Create Variants : %v ", err)
	}

	return &types.ProductVariant{
		ID:        variants.ID,
		Color:     variants.Color,
		SKU:       variants.SKU,
		ProductID: variants.ProductID,
	}, nil
}

func (pr *Database) UpdateProductVariant(ctx context.Context, req *request.UpdateProductVariantRequest) (*types.ProductVariant, error) {
	query := `
        UPDATE product_variants
        SET color = $1, sku = $2
        WHERE id = $3
        RETURNING id, color, sku, product_id
    `

	var variant types.ProductVariant
	err := pr.DB.QueryRowContext(ctx, query, req.Variant.Color, req.Variant.SKU, req.Variant.ID).Scan(
		&variant.ID,
		&variant.Color,
		&variant.SKU,
		&variant.ProductID,
	)

	if err != nil {
		pr.log.Error("Failed to update variant: %v", err)
		return nil, err
	}

	return &variant, nil
}

func (pr *Database) DeleteProductVariant(ctx context.Context, req *request.DeleteProductVariantRequest) (*response.DeleteProductResponse, error) {
	query := `
	DELETE FROM product_variants WHERE id = $1
	`

	_, err := pr.DB.ExecContext(ctx, query, req.ID)

	if err != nil {
		pr.log.Error("Failed to Delete Variants : %v", err)
		return nil, err
	}

	return &response.DeleteProductResponse{
		Success: true,
	}, nil

}

func (pr *Database) GetProductsVariant(ctx context.Context, req *request.GetProductVariantRequest) (*types.ProductVariant, error) {
	log.Printf("Request from : %s", req.VariantID)
	var variants types.ProductVariant
	query := `
		SELECT id,color,sku,product_id
		FROM product_variants
		WHERE id = $1
	`
	err := pr.DB.QueryRowContext(ctx, query, req.VariantID).Scan(
		&variants.ID,
		&variants.Color,
		&variants.SKU,
		&variants.ProductID,
	)

	if err != nil {
		pr.log.Error("Failed to get Variants : %v", err)
		return nil, nil
	}
	return &variants, nil
}

func (pr *Database) GetProductVariants(ctx context.Context, req *request.GetProductVariantsRequest) (*response.GetProductVariantsResponse, error) {
	query := `
    SELECT 
        id, color, sku, product_id
    FROM product_variants 
    WHERE product_id = $1
    `
	rows, err := pr.DB.QueryContext(ctx, query, req.ProductID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []*types.ProductVariant

	for rows.Next() {
		variant := &types.ProductVariant{}
		err := rows.Scan(
			&variant.ID,
			&variant.Color,
			&variant.SKU,
			&variant.ProductID,
		)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Product Variants : %v", err.Error())
		return nil, err
	}

	return &response.GetProductVariantsResponse{
		Variants: variants}, nil
}

func (r *Database) getProductVariants(ctx context.Context, productID string) ([]*types.ProductVariant, error) {
	const query = `
        SELECT id, color, sku
        FROM product_variants
        WHERE product_id = $1
    `

	rows, err := r.DB.QueryContext(ctx, query, productID)
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to get variants: %v", err)
		return nil, fmt.Errorf("failed to get product variants")
	}
	defer rows.Close()

	var variants []*types.ProductVariant
	for rows.Next() {
		var variant types.ProductVariant
		err := rows.Scan(&variant.ID, &variant.Color, &variant.SKU)
		if err != nil {
			r.log.Log(logger.ErrorLevel, "Failed to scan variant row: %v", err)
			return nil, fmt.Errorf("failed to scan variant row")
		}

		variant.ProductID = productID
		variant.Images = make([]*types.ProductImage, 0)
		variant.Inventory = make([]*types.Inventory, 0)
		variants = append(variants, &variant)
	}

	return variants, nil
}
