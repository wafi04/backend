package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
)

type Success struct {
	Success bool `json:"succcess"`
}
type ShippingAddressRepo interface {
	CreateShippingAddres(ctx context.Context, req *request.CreateAddressReq) (*types.ShippingAddress, error)
	GetShippingAddress(ctx context.Context, userID string) (*types.ListShippingAddress, error)
	UpdateShippingAddress(ctx context.Context, req *request.UpdateAddressReq) (*types.ShippingAddress, error)
	SetDefaultAddress(ctx context.Context, userID, addressID string) (*Success, error)
	DeleteShippingAddress(ctx context.Context, userID, addressID string) (*Success, error)
}

func NewShippingAddressRepo(db *sqlx.DB) ShippingAddressRepo {
	return &Database{db: db}
}

func (d *Database) CreateShippingAddres(ctx context.Context, req *request.CreateAddressReq) (*types.ShippingAddress, error) {
	var shipAddr types.ShippingAddress
	query := `
	INSERT INTO shipping_addresses (
		address_id,
		user_id,
		recipient_name,
		recipient_phone,
		full_address,
		city,
		province,
		postal_code,
		country,
		label,
		is_default,
		created_at,
		updated_at
	) VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13 
	) RETURNING 
		address_id,
		user_id,
		recipient_name,
		recipient_phone,
		full_address,
		city,
		province,
		postal_code,
		country,
		label,
		is_default,
		created_at,
		updated_at
	`
	err := d.db.DB.QueryRowContext(ctx, query,
		req.AddressID,
		req.UserID,
		req.RecipientName,
		req.Recipientphone,
		req.FullAddress,
		req.City,
		req.Province,
		req.PostalCode,
		req.Country,
		req.Label,
		false,
		time.Now(),
		time.Now(),
	).Scan(
		&shipAddr.AddressID,
		&shipAddr.UserID,
		&shipAddr.RecipientName,
		&shipAddr.Recipientphone,
		&shipAddr.FullAddress,
		&shipAddr.City,
		&shipAddr.Province,
		&shipAddr.PostalCode,
		&shipAddr.Country,
		&shipAddr.Label,
		&shipAddr.IsDefault,
		&shipAddr.CreatedAt,
		&shipAddr.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create shipping address : %v", err)
	}

	return &shipAddr, nil
}

func (d *Database) GetShippingAddress(ctx context.Context, userID string) (*types.ListShippingAddress, error) {
	query := `
		SELECT 
			address_id,
			user_id,
			recipient_name,
			recipient_phone,
			full_address,
			city,
			province,
			postal_code,
			country,
			label,
			is_default,
			created_at,
			updated_at
		FROM shipping_addresses
		WHERE user_id  = $1
	`
	rows, err := d.db.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get data : %v", err)
	}

	var shipAddrs []*types.ShippingAddress
	for rows.Next() {
		shipAddr := &types.ShippingAddress{}
		err := rows.Scan(
			&shipAddr.AddressID,
			&shipAddr.UserID,
			&shipAddr.RecipientName,
			&shipAddr.Recipientphone,
			&shipAddr.FullAddress,
			&shipAddr.City,
			&shipAddr.Province,
			&shipAddr.PostalCode,
			&shipAddr.Country,
			&shipAddr.Label,
			&shipAddr.IsDefault,
			&shipAddr.CreatedAt,
			&shipAddr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		shipAddrs = append(shipAddrs, shipAddr)
	}

	if err = rows.Err(); err != nil {
		return nil, nil
	}

	return &types.ListShippingAddress{
		Address: shipAddrs,
	}, nil
}
func (d *Database) UpdateShippingAddress(ctx context.Context, req *request.UpdateAddressReq) (*types.ShippingAddress, error) {
	var shipAddr types.ShippingAddress
	query := `
        UPDATE shipping_addresses
        SET
            recipient_name = COALESCE($1, recipient_name),
            recipient_phone = COALESCE($2, recipient_phone),
            full_address = COALESCE($3, full_address),
            city = COALESCE($4, city),
            province = COALESCE($5, province),
            postal_code = COALESCE($6, postal_code),
            country = COALESCE($7, country),
            label = COALESCE($8, label),
            is_default = COALESCE($9, is_default),
            updated_at = $10
        WHERE user_id = $11 AND address_id = $12
        RETURNING
            address_id,
            user_id,
            recipient_name,
            recipient_phone,
            full_address,
            city,
            province,
            postal_code,
            country,
            label,
            is_default,
            created_at,
            updated_at
    `

	err := d.db.QueryRowContext(ctx, query,
		req.RecipientName,
		req.Recipientphone,
		req.FullAddress,
		req.City,
		req.Province,
		req.PostalCode,
		req.Country,
		req.Label,
		req.IsDefault,
		time.Now(),
		req.UserID,
		req.AddressID,
	).Scan(
		&shipAddr.AddressID,
		&shipAddr.UserID,
		&shipAddr.RecipientName,
		&shipAddr.Recipientphone,
		&shipAddr.FullAddress,
		&shipAddr.City,
		&shipAddr.Province,
		&shipAddr.PostalCode,
		&shipAddr.Country,
		&shipAddr.Label,
		&shipAddr.IsDefault,
		&shipAddr.CreatedAt,
		&shipAddr.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update shipping address: %v", err)
	}

	return &shipAddr, nil
}

func (d *Database) SetDefaultAddress(ctx context.Context, userID, addressID string) (*Success, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// First, unset all default addresses for this user
	query1 := `
        UPDATE shipping_addresses
        SET is_default = false
        WHERE user_id = $1
    `
	_, err = tx.ExecContext(ctx, query1, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to unset default addresses: %v", err)
	}

	// Then, set the specified address as default
	query2 := `
        UPDATE shipping_addresses
        SET 
            is_default = true,
            updated_at = $1
        WHERE user_id = $2 AND address_id = $3
    `
	result, err := tx.ExecContext(ctx, query2, time.Now(), userID, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to set default address: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("address not found")
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &Success{
		Success: true,
	}, nil
}

func (d *Database) DeleteShippingAddress(ctx context.Context, userID, addressID string) (*Success, error) {
	query := `
        DELETE FROM shipping_addresses
        WHERE user_id = $1 AND address_id = $2
    `

	result, err := d.db.ExecContext(ctx, query, userID, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete shipping address: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("address not found")
	}

	return &Success{
		Success: true,
	}, nil
}
