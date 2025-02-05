package inventory

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/backend/pkg/logger"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
	"github.com/wafi04/backend/pkg/utils"
)

type Database struct {
	DB  *sqlx.DB
	log logger.Logger
}

type InventoryRepository interface {
	GetInventoryByVariant(ctx context.Context, req *request.GetInventoryByVariantRequest) (*request.GetInventoryByVariantResponse, error)
	CreateInventory(ctx context.Context, req *request.CreateInventoryRequest) (*types.Inventory, error)
	UpdateInventory(ctx context.Context, req *request.UpdateInventoryRequest) (*types.Inventory, error)
	CheckAvailability(ctx context.Context, req *Req) (*Res, error)
}

func NewInventoryRepository(DB *sqlx.DB) InventoryRepository {
	return &Database{DB: DB}
}

func (r *Database) GetInventoryByVariant(ctx context.Context, req *request.GetInventoryByVariantRequest) (*request.GetInventoryByVariantResponse, error) {
	query := `
        SELECT 
            id,
            variant_id,
            size,
            stock,
			available_stock,
			reserved_stock
        FROM inventory 
        WHERE variant_id = $1
    `
	rows, err := r.DB.QueryContext(ctx, query, req.VariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var inventorys []types.Inventory

	for rows.Next() {
		var inv types.Inventory
		err := rows.Scan(
			&inv.ID,
			&inv.VariantID,
			&inv.Size,
			&inv.Stock,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		inventorys = append(inventorys, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return &request.GetInventoryByVariantResponse{
		Inventory: inventorys,
	}, nil
}

func (r *Database) CreateInventory(ctx context.Context, req *request.CreateInventoryRequest) (*types.Inventory, error) {
	if req.Stock < 0 {
		return nil, fmt.Errorf("invalid stock value: %d", req.Stock)
	}
	if req.Size <= "" {
		return nil, fmt.Errorf("invalid size value: %s", req.Size)
	}

	InventoryID := utils.GenerateRandomId("INV")

	var created_at, updated_at time.Time
	var inv types.Inventory

	query := `
        INSERT INTO inventory (id, stock, size, available_stock, reserved_stock, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, stock, size, available_stock, reserved_stock, created_at, updated_at
    `

	err := r.DB.QueryRowContext(ctx, query, InventoryID, req.Stock, req.Size, req.Stock, 0, time.Now(), time.Now()).Scan(
		&inv.ID,
		&inv.Stock,
		&inv.Size,
		&inv.AvailableStock,
		&inv.ReservedStock,
		&created_at,
		&updated_at,
	)
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to create inventory: %v", err)
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}

	inv.CreatedAt = created_at.Unix()
	inv.UpdatedAt = updated_at.Unix()

	return &inv, nil
}

func (r *Database) UpdateInventory(ctx context.Context, req *request.UpdateInventoryRequest) (*types.Inventory, error) {
	var created_at, updated_at time.Time
	var inv types.Inventory

	query := `
        UPDATE inventory 
        SET 
            size = $2,
            stock = $3,
            reserved_stock = $4,
            available_stock = $5,
            updated_at = NOW() 
        WHERE 
            id = $1
        RETURNING id, stock, size, available_stock, reserved_stock, created_at, updated_at
    `

	err := r.DB.QueryRowContext(ctx, query, req.ID, req.Size, req.Stock, req.ReservedStock, req.AvailableStock).Scan(
		&inv.ID,
		&inv.Stock,
		&inv.Size,
		&inv.AvailableStock,
		&inv.ReservedStock,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	inv.CreatedAt = created_at.Unix()
	inv.UpdatedAt = updated_at.Unix()

	return &inv, nil
}

type Req struct {
	VariantId string `json:"variant_id"`
	Quantity  int64  `json:"quantity"`
}

type Res struct {
	Available      bool  `json:"available"`
	AvailableStock int64 `json:"available_stock"`
}

func (r *Database) CheckAvailability(ctx context.Context, req *Req) (*Res, error) {
	query := `
    SELECT 
        stock - reserved_stock AS available_stock
    FROM 
        inventory
    WHERE 
        variant_id = $1;
    `

	var availableStock int64
	err := r.DB.QueryRowContext(ctx, query, req.VariantId).Scan(&availableStock)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Res{
				Available:      false,
				AvailableStock: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}

	available := availableStock >= req.Quantity

	return &Res{
		Available:      available,
		AvailableStock: availableStock,
	}, nil
}
