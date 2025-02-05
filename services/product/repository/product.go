package productRepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/wafi04/backend/pkg/logger"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
	response "github.com/wafi04/backend/pkg/types/res"
)

type Database struct {
	DB  *sqlx.DB
	log logger.Logger
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &Database{DB: db}
}

func (s *Database) CreateProduct(ctx context.Context, req *types.Product) (*types.Product, error) {
	now := time.Now()
	query := `
    INSERT INTO products  
    (id, name, sub_title, description, sku, price, category_id, created_at, updated_at)
    VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id, name, sub_title, description, sku, price, category_id, created_at, updated_at
    `

	var product types.Product
	var createdAt, updatedAt time.Time

	err := s.DB.QueryRowContext(ctx, query,
		req.ID,
		req.Name,
		req.SubTitle,
		req.Description,
		req.SKU,
		req.Price,
		req.CategoryID,
		now,
		now,
	).Scan(
		&product.ID,
		&product.Name,
		&product.SubTitle,
		&product.Description,
		&product.SKU,
		&product.Price,
		&product.CategoryID,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		s.log.Log(logger.ErrorLevel, "Failed to insert product: %v", err)
		return nil, fmt.Errorf("failed to insert Product: %v", err)
	}

	product.CreatedAt = time.Now().Unix()
	product.UpdatedAt = time.Now().Unix()

	return &product, nil
}
func (r *Database) GetProduct(ctx context.Context, req *request.GetProductRequest) (*types.Product, error) {
	product, err := r.getProductBase(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	variants, err := r.getProductVariants(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if len(variants) > 0 {
		variantIDs := extractVariantIDs(variants)
		if err := r.enrichVariantsWithImages(ctx, variants, variantIDs); err != nil {
			return nil, err
		}
		if err := r.EnrichVariantsWithInventory(ctx, variants, variantIDs); err != nil {
			return nil, err
		}
	}

	product.Variants = variants
	return product, nil
}

func (r *Database) getProductBase(ctx context.Context, productID string) (*types.Product, error) {
	const query = `
        SELECT 
            id, name, sub_title, description, 
            price, sku, category_id, 
            created_at, updated_at
        FROM products
        WHERE id = $1
    `

	product := &types.Product{}
	var subTitle sql.NullString
	var createdAt, updatedAt time.Time

	err := r.DB.QueryRowContext(ctx, query, productID).Scan(
		&product.ID, &product.Name, &subTitle, &product.Description,
		&product.Price, &product.SKU, &product.CategoryID,
		&createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Log(logger.ErrorLevel, "Product not found: %v", err)
			return nil, fmt.Errorf("product not found")
		}
		r.log.Log(logger.ErrorLevel, "Failed to get product: %v", err)
		return nil, fmt.Errorf("failed to get product")
	}

	if subTitle.Valid {
		product.SubTitle = subTitle.String
	}
	product.CreatedAt = createdAt.Unix()
	product.UpdatedAt = updatedAt.Unix()

	return product, nil
}

func (r *Database) EnrichVariantsWithInventory(ctx context.Context, variants []*types.ProductVariant, variantIDs []string) error {
	const query = `
        SELECT 
            id, variant_id, size, stock,
            reserved_stock, available_stock,
            created_at, updated_at
        FROM inventory
        WHERE variant_id = ANY($1)
        ORDER BY size
    `

	rows, err := r.DB.QueryContext(ctx, query, pq.Array(variantIDs))
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to get inventory: %v", err)
		return fmt.Errorf("failed to get inventory")
	}
	defer rows.Close()

	variantMap := createVariantMap(variants)

	for rows.Next() {
		var inv types.Inventory
		var variantID string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&inv.ID, &variantID, &inv.Size, &inv.Stock,
			&inv.ReservedStock, &inv.AvailableStock,
			&createdAt, &updatedAt,
		); err != nil {
			r.log.Log(logger.ErrorLevel, "Failed to scan inventory row: %v", err)
			return fmt.Errorf("failed to scan inventory row")
		}

		if variant, exists := variantMap[variantID]; exists {
			inv.VariantID = variantID
			inv.CreatedAt = createdAt.Unix()
			inv.UpdatedAt = updatedAt.Unix()
			variant.Inventory = append(variant.Inventory, &inv)
		}
	}

	return nil
}

// Helper functions
func extractVariantIDs(variants []*types.ProductVariant) []string {
	ids := make([]string, len(variants))
	for i, v := range variants {
		ids[i] = v.ID
	}
	return ids
}

func createVariantMap(variants []*types.ProductVariant) map[string]*types.ProductVariant {
	variantMap := make(map[string]*types.ProductVariant)
	for _, v := range variants {
		variantMap[v.ID] = v
	}
	return variantMap
}

func (s *Database) ListProducts(ctx context.Context, req *request.ListProductsRequest) (*response.ListProductsResponse, error) {
	if req.PageToken == "" {
		req.PageToken = "0"
	}

	baseQuery := `
        SELECT 
            p.id,
            p.name,
            p.sub_title,
            p.description,
            p.price,
            p.sku,
            p.category_id,
            p.created_at,
            p.updated_at,
            (
                SELECT COALESCE(JSON_AGG(
                    json_build_object(
                        'id', v.id,
                        'color', v.color,
                        'sku', v.sku,
                        'product_id', v.product_id,
                        'images', (
                            SELECT COALESCE(JSON_AGG(
                                json_build_object(
                                    'id', i.id,
                                    'url', i.url,
                                    'variant_id', i.variant_id,
                                    'is_main', i.is_main
                                )
                            ), '[]'::json)
                            FROM product_images i
                            WHERE i.variant_id = v.id
                        )
                    )
                ), '[]'::json)
                FROM product_variants v
                WHERE v.product_id = p.id
            ) AS variants
        FROM 
            products p
        WHERE 1=1
        ORDER BY p.created_at DESC
        LIMIT $1
        OFFSET ($1 * COALESCE(NULLIF($2, ''), '0')::integer)
    `

	params := []interface{}{
		req.PageSize,
		req.PageToken,
	}

	rows, err := s.DB.QueryxContext(ctx, baseQuery, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %v", err)
	}
	defer rows.Close()

	var products []*types.Product
	for rows.Next() {
		var product struct {
			ID          string          `db:"id"`
			Name        string          `db:"name"`
			SubTitle    sql.NullString  `db:"sub_title"`
			Description string          `db:"description"`
			Price       float64         `db:"price"`
			SKU         string          `db:"sku"`
			CategoryID  string          `db:"category_id"`
			CreatedAt   time.Time       `db:"created_at"`
			UpdatedAt   time.Time       `db:"updated_at"`
			Variants    json.RawMessage `db:"variants"`
		}

		if err := rows.StructScan(&product); err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}

		var variants []*types.ProductVariant
		if err := json.Unmarshal(product.Variants, &variants); err != nil {
			return nil, fmt.Errorf("failed to parse variants: %v", err)
		}

		pbProduct := &types.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			SKU:         product.SKU,
			CategoryID:  product.CategoryID,
			CreatedAt:   product.CreatedAt.Unix(),
			UpdatedAt:   product.UpdatedAt.Unix(),
			Variants:    variants,
		}
		if product.SubTitle.Valid {
			pbProduct.SubTitle = product.SubTitle.String
		}

		products = append(products, pbProduct)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %v", err)
	}

	// Calculate next page token
	nextPageToken := ""
	if len(products) == int(req.PageSize) {
		currentPage, _ := strconv.Atoi(req.PageToken)
		nextPageToken = strconv.Itoa(currentPage + 1)
	}

	return &response.ListProductsResponse{
		Products:      products,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Database) UpdateProduct(ctx context.Context, req *request.UpdateProductRequest) (*types.Product, error) {
	var product types.Product
	query := `
	UPDATE products
	SET 
		name = $1,
		sub_title  = $2,
		description = $3,
		price = $4,
		sku = $5,
		category_id  = $6
	WHERE id = $7
	RETURNING 
		id,
		name,
		sub_title,
		description,
		price,
		sku,
		category_id,
		created_at,
		updated_at
	`
	var createdAt, updatedAt time.Time

	err := s.DB.QueryRowContext(ctx, query,
		req.Product.Name,
		req.Product.SubTitle,
		req.Product.Description,
		req.Product.Price,
		req.Product.SKU,
		req.Product.CategoryID,
		req.Product.ID,
	).Scan(
		&product.ID,
		&product.Name,
		&product.SubTitle,
		&product.Description,
		&product.Price,
		&product.SKU,
		&product.CategoryID,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to delete product: %v", err)
	}

	product.CreatedAt = createdAt.Unix()
	product.UpdatedAt = updatedAt.Unix()
	return &product, nil
}

func (s *Database) DeleteProduct(ctx context.Context, req *request.DeleteProductRequest) (*response.DeleteProductResponse, error) {
	query := `
        DELETE FROM products WHERE id = $1;
    `

	_, err := s.DB.ExecContext(ctx, query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete product: %v", err)
	}

	return &response.DeleteProductResponse{
		Success: true,
	}, nil
}
