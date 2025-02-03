package cartrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/backend/internal/handler/dto/request"
	"github.com/wafi04/backend/internal/handler/dto/response"
	"github.com/wafi04/backend/internal/handler/dto/types"
	"github.com/wafi04/backend/pkg/logger"
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
}

func NewCartRepository(db *sqlx.DB) CartRepository {
	return &Database{db: db}
}

func (d *Database) RemoveFromCart(ctx context.Context, req *request.ReqRemoveCartByID) (*response.ResRemoveCartItem, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var cartID string
	queryVerify := `
		SELECT ci.cart_id
		FROM cart_items ci
		JOIN carts c ON ci.cart_id = c.cart_id
		WHERE ci.cart_item_id = $1 AND c.user_id = $2
	`
	err = tx.QueryRowContext(ctx, queryVerify, req.CartItemID, req.UserID).Scan(&cartID)
	if err != nil {
		return nil, err
	}
	queryDelete := `
		DELETE FROM cart_items 
		WHERE cart_item_id = $1
	`
	_, err = tx.ExecContext(ctx, queryDelete, req.CartItemID)
	if err != nil {
		return nil, err
	}

	queryUpdateTotal := `
		UPDATE carts 
		SET total = (
			SELECT COALESCE(SUM(sub_total), 0) 
			FROM cart_items 
			WHERE cart_id = $1
		),
		updated_at = CURRENT_TIMESTAMP
		WHERE cart_id = $1
	`
	_, err = tx.ExecContext(ctx, queryUpdateTotal, cartID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &response.ResRemoveCartItem{
		Success: true,
	}, nil
}

func (d *Database) ClearCart(ctx context.Context, req *request.ClearCart) (*response.ResRemoveCartItem, error) {
	query := `
		DELETE FROM carts
		WHERE carts_id = $1 
	`

	_, err := d.db.ExecContext(ctx, query, req.CartID)
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

	err = d.UpdateCartTotal(ctx, cart.CartID)
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
		return nil, err
	}

	queryItems := `
    SELECT 
        cart_item_id,
        cart_id,
        variant_id,
        size,
        quantity,
        sub_total,
        created_at,
        updated_at
    FROM cart_items
    WHERE cart_id = $1
    `
	rows, err := d.db.QueryContext(ctx, queryItems, cart.CartID)
	if err != nil {
		return nil, err
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
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	cart.Item = items
	return &cart, nil
}
