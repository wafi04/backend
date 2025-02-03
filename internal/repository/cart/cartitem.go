package cartrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wafi04/backend/internal/handler/dto/request"
	"github.com/wafi04/backend/internal/handler/dto/response"
	"github.com/wafi04/backend/internal/handler/dto/types"
	"github.com/wafi04/backend/pkg/logger"
)

func (d *Database) getOrCreateCart(ctx context.Context, tx *sql.Tx, req *request.CartRequest) (*types.Cart, error) {
	var cart types.Cart
	getCartQuery := `
		SELECT cart_id, user_id, total, created_at, updated_at
		FROM carts
		WHERE user_id = $1`

	err := tx.QueryRowContext(ctx, getCartQuery, req.UserID).Scan(
		&cart.CartID,
		&cart.UserID,
		&cart.Total,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	insertCartQuery := `
		INSERT INTO carts (cart_id, user_id, total)
		VALUES ($1, $2, $3)`

	if err == sql.ErrNoRows {
		cart.CartID = uuid.New().String()
		cart.UserID = req.UserID
		cart.Total = req.Total

		_, err = tx.ExecContext(ctx, insertCartQuery,
			cart.CartID, cart.UserID, cart.Total)
		if err != nil {
			return nil, fmt.Errorf("failed to create cart: %w", err)
		}
		return &cart, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &cart, nil
}

func (d *Database) getCartItem(ctx context.Context, tx *sql.Tx, cartID, variantID, size string) (*types.CartItem, error) {
	var item types.CartItem
	getCartItemQuery := `
		SELECT cart_item_id, cart_id, product_variant_id, size, quantity, sub_total
		FROM cart_items
		WHERE cart_id = $1 AND product_variant_id = $2 AND size = $3`
	err := tx.QueryRowContext(ctx, getCartItemQuery, cartID, variantID, size).Scan(
		&item.CartItemID,
		&item.CartID,
		&item.VariantID,
		&item.Size,
		&item.Quantity,
		&item.SubTotal,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}

	return &item, nil
}

func (d *Database) getProductPrice(ctx context.Context, tx *sql.Tx, variantID string) (float64, error) {
	var price float64

	getProductPriceQuery := `
		SELECT price FROM product_variants WHERE id = $1`

	err := tx.QueryRowContext(ctx, getProductPriceQuery, variantID).Scan(&price)
	if err != nil {
		return 0, fmt.Errorf("failed to get product price: %w", err)
	}
	return price, nil
}

func (d *Database) AddCart(ctx context.Context, req *request.CartRequest) (*response.CartResponse, error) {
	d.logger.Log(logger.InfoLevel, "Incoming request data: %s", req.UserID)

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	cart, err := d.getOrCreateCart(ctx, tx, req)
	if err != nil {
		return nil, err
	}

	item, err := d.getCartItem(ctx, tx, cart.CartID, req.VariantID, req.Size)
	if err != nil {
		return nil, err
	}

	price, err := d.getProductPrice(ctx, tx, req.VariantID)
	if err != nil {
		return nil, err
	}

	subTotal := float64(req.Quantity) * price
	insertCartItemQuery := `
		INSERT INTO cart_items (
			cart_item_id, cart_id, product_variant_id,
			size, quantity, sub_total
		)
		VALUES ($1, $2, $3, $4, $5, $6)`

	if item == nil {
		newItem := types.CartItem{
			CartItemID: uuid.New().String(),
			CartID:     cart.CartID,
			VariantID:  req.VariantID,
			Size:       req.Size,
			Quantity:   req.Quantity,
			SubTotal:   subTotal,
		}

		_, err = tx.ExecContext(ctx, insertCartItemQuery,
			newItem.CartItemID, newItem.CartID, newItem.VariantID,
			newItem.Size, newItem.Quantity, newItem.SubTotal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert cart item: %w", err)
		}
	} else {

		updateCartItemQuery := `
			UPDATE cart_items
			SET quantity = $1, sub_total = $2, updated_at = $3
			WHERE cart_item_id = $4
		`
		_, err = tx.ExecContext(ctx, updateCartItemQuery,
			req.Quantity, subTotal, time.Now(), item.CartItemID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update cart item: %w", err)
		}
	}

	err = d.UpdateCartTotal(ctx, cart.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart total: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &response.CartResponse{
		VariantID: req.VariantID,
		Size:      req.Size,
		Quantity:  req.Quantity,
	}, nil
}

func (d *Database) UpdateCartTotal(ctx context.Context, cartID string) error {
	updateCartTotalQuery := `
		UPDATE carts
		SET total = (
			SELECT COALESCE(SUM(sub_total), 0)
			FROM cart_items
			WHERE cart_id = $1
		),
		updated_at = $2
		WHERE cart_id = $1
	`

	_, err := d.db.ExecContext(ctx, updateCartTotalQuery, cartID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update cart total: %w", err)
	}
	return nil
}
