package cart

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *Database) GetCartItemCount(ctx context.Context, cartID string) (int, error) {
	var itemCount int

	query := `
        SELECT COUNT(*) AS item_count
        FROM cart_items
        WHERE cart_id = $1
    `

	err := r.db.QueryRowContext(ctx, query, cartID).Scan(&itemCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get cart item count: %w", err)
	}

	return itemCount, nil
}
