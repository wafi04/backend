package cart

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/wafi04/backend/pkg/logger"
	request "github.com/wafi04/backend/pkg/types/req"
	response "github.com/wafi04/backend/pkg/types/res"
)

func (d *Database) getCartQuery(ctx context.Context, tx *sql.Tx, userID string) (string, error) {
	var cartID string
	getCartQuery := `
        SELECT cart_id 
        FROM carts 
        WHERE user_id = $1
    `

	err := tx.QueryRowContext(ctx, getCartQuery, userID).Scan(&cartID)

	return cartID, err

}
func (d *Database) getProductPrice(ctx context.Context, tx *sql.Tx, variantID string) (float64, error) {
	var price float64
	d.logger.Log(logger.InfoLevel, "VariantId : %s", variantID)
	getProductPriceQuery := `
		SELECT price FROM product_variants WHERE id = $1`

	err := tx.QueryRowContext(ctx, getProductPriceQuery, variantID).Scan(&price)
	if err != nil {
		return 0, fmt.Errorf("failed to get product price: %w", err)
	}
	return price, nil
}
func (d *Database) AddCart(ctx context.Context, req *request.CartRequest) (*response.CartResponse, error) {
	d.logger.Log(logger.InfoLevel, "Adding/updating cart for user: %s, variant: %s", req.UserID, req.VariantID)

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	cartID, err := d.getCartQuery(ctx, tx, req.UserID)
	if err == sql.ErrNoRows {
		cartID = uuid.New().String()
		_, err = tx.ExecContext(ctx, `
            INSERT INTO carts (cart_id, user_id, total)
            VALUES ($1, $2, 0)
        `, cartID, req.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to create new cart: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to query cart: %w", err)
	}

	subTotal := float64(req.Quantity) * req.Price
	var existingItemID string
	var existingQuantity int

	checkItemQuery := `
        SELECT cart_item_id, quantity 
        FROM cart_items 
        WHERE cart_id = $1 
        AND product_variant_id = $2 
        AND size = $3
    `
	err = tx.QueryRowContext(ctx, checkItemQuery, cartID, req.VariantID, req.Size).Scan(&existingItemID, &existingQuantity)

	if err == sql.ErrNoRows {
		newItemID := uuid.New().String()
		_, err = tx.ExecContext(ctx, `
            INSERT INTO cart_items (
                cart_item_id, cart_id, product_variant_id,
                size, quantity, sub_total
            )
            VALUES ($1, $2, $3, $4, $5, $6)
        `, newItemID, cartID, req.VariantID, req.Size, req.Quantity, subTotal)
		if err != nil {
			return nil, fmt.Errorf("failed to insert cart item: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check existing item: %w", err)
	} else {
		newQuantity := existingQuantity + int(req.Quantity)
		newSubTotal := float64(newQuantity) * req.Price

		_, err = tx.ExecContext(ctx, `
            UPDATE cart_items 
            SET quantity = $1, 
                sub_total = $2,
                updated_at = CURRENT_TIMESTAMP
            WHERE cart_item_id = $3
        `, newQuantity, newSubTotal, existingItemID)
		if err != nil {
			return nil, fmt.Errorf("failed to update cart item: %w", err)
		}
	}

	err = d.UpdateCartTotal(ctx, tx, cartID)
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

func (d *Database) UpdateCartTotal(ctx context.Context, tx *sql.Tx, cartID string) error {
	d.logger.Log(logger.InfoLevel, "Updating cart total for cartID=%s", cartID)

	query := `
        UPDATE carts
        SET total = (
            SELECT COALESCE(SUM(sub_total), 0)
            FROM cart_items
            WHERE cart_id = $1
        ),
        updated_at = CURRENT_TIMESTAMP
        WHERE cart_id = $1
    `
	_, err := tx.ExecContext(ctx, query, cartID)
	if err != nil {
		d.logger.Log(logger.ErrorLevel, "Failed to update cart total: %v", err)
		return fmt.Errorf("failed to update cart total: %w", err)
	}

	d.logger.Log(logger.InfoLevel, "Successfully updated cart total for cartID=%s", cartID)
	return nil
}
