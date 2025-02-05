package cart

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/wafi04/backend/pkg/logger"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
	response "github.com/wafi04/backend/pkg/types/res"
)

type Database struct {
	db     *sqlx.DB
	logger logger.Logger
}

type CartRepository interface {
	AddCart(ctx context.Context, req *request.CartRequest) (*response.CartResponse, error)
	RemoveFromCart(ctx context.Context, req *request.ReqRemoveCartByID) (*response.ResRemoveCartItem, error)
	ClearCart(ctx context.Context, req *request.ClearCart) (*response.ResRemoveCartItem, error)
	UpdateQuantity(ctx context.Context, req *request.UpdateQuantity) (*types.CartItem, error)
	GetCart(ctx context.Context, userID string) (*types.Cart, error)
	GetCartItemCount(ctx context.Context, cartID string) (int, error)
}

func NewCartRepository(db *sqlx.DB) CartRepository {
	return &Database{db: db}
}
func (d *Database) RemoveFromCart(ctx context.Context, req *request.ReqRemoveCartByID) (*response.ResRemoveCartItem, error) {
	d.logger.Log(logger.InfoLevel, "Removing cart item: %s for user: %s", req.CartItemID, req.UserID)

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var cartID string
	err = tx.QueryRowContext(ctx, `
        SELECT cart_id 
        FROM cart_items 
        WHERE cart_item_id = $1
    `, req.CartItemID).Scan(&cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart_id: %w", err)
	}

	queryDelete := `
        DELETE FROM cart_items 
        WHERE cart_item_id = $1
    `
	result, err := tx.ExecContext(ctx, queryDelete, req.CartItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete cart item: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	d.logger.Log(logger.InfoLevel, "Deleted %d cart item(s)", rowsAffected)

	err = d.UpdateCartTotal(ctx, tx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart total: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &response.ResRemoveCartItem{
		Success: true,
	}, nil
}

func (d *Database) ClearCart(ctx context.Context, req *request.ClearCart) (*response.ResRemoveCartItem, error) {
	query := `
		DELETE FROM carts
		WHERE user_id = $1 
	`
	_, err := d.db.ExecContext(ctx, query, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear cart : %v", err)
	}

	return &response.ResRemoveCartItem{
		Success: true,
	}, nil
}

func (d *Database) UpdateQuantity(ctx context.Context, req *request.UpdateQuantity) (*types.CartItem, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var cart types.CartItem
	queryGet := `
		SELECT 
			cart_item_id,
			cart_id,
			product_variant_id,
			quantity,
			size,
			sub_total,
			created_at,
			updated_at
		FROM cart_items
		WHERE cart_item_id = $1
	`
	err = tx.QueryRowContext(ctx, queryGet, req.CartItemID).Scan(
		&cart.CartItemID,
		&cart.CartID,
		&cart.VariantID,
		&cart.Quantity,
		&cart.Size,
		&cart.SubTotal,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cart item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}

	price, err := d.getProductPrice(ctx, tx, cart.VariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get price from product : %v", err)
	}

	newSubTotal := float64(req.Quantity) * price

	queryUpdate := `
		UPDATE cart_items
		SET 
			quantity = $1,
			sub_total = $2,
			size = $3,
			updated_at = $4
		WHERE cart_item_id = $5
		RETURNING 
			cart_item_id,
			cart_id,
			product_variant_id,
			quantity,
			size,
			sub_total,
			created_at,
			updated_at
	`
	now := time.Now()
	err = tx.QueryRowContext(ctx, queryUpdate,
		req.Quantity,
		newSubTotal,
		req.Size,
		now,
		req.CartItemID,
	).Scan(
		&cart.CartItemID,
		&cart.CartID,
		&cart.VariantID,
		&cart.Quantity,
		&cart.Size,
		&cart.SubTotal,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	err = d.UpdateCartTotal(ctx, tx, cart.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart total: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &cart, nil
}
func (d *Database) GetCart(ctx context.Context, userID string) (*types.Cart, error) {
	var cart types.Cart

	queryCart := `
    SELECT 
        cart_id,
        user_id,
        total,
        created_at,
        updated_at
    FROM carts
    WHERE user_id = $1
    `
	err := d.db.QueryRowContext(ctx, queryCart, userID).Scan(
		&cart.CartID,
		&cart.UserID,
		&cart.Total,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart not found for user: %s", userID)
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	queryItems := `
    SELECT 
        ci.cart_item_id,
        ci.cart_id,
        ci.product_variant_id,
        ci.size,
        ci.quantity,
        ci.sub_total,
        ci.created_at,
        ci.updated_at,
        pi.url AS image_url,
        pv.color,           
        pv.sku,            
        p.name AS product_name 
    FROM cart_items ci
    LEFT JOIN product_variants pv ON ci.product_variant_id = pv.id
    LEFT JOIN products p ON pv.product_id = p.id
    LEFT JOIN product_images pi ON pv.id = pi.variant_id AND pi.is_main = TRUE
    WHERE ci.cart_id = $1
    `
	rows, err := d.db.QueryContext(ctx, queryItems, cart.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart items: %w", err)
	}
	defer rows.Close()

	var items []types.CartItem
	for rows.Next() {
		var item types.CartItem
		err := rows.Scan(
			&item.CartItemID,
			&item.CartID,
			&item.VariantID,
			&item.Size,
			&item.Quantity,
			&item.SubTotal,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ImageURL,
			&item.Color,
			&item.SKU,
			&item.ProductName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cart items: %w", err)
	}

	cart.Item = items
	return &cart, nil
}
