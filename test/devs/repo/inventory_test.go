package inventoryrepo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/wafi04/backend/internal/handler/dto/request"
	inventoryrepo "github.com/wafi04/backend/internal/repository/inventory"
)

func TestCreateInventory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := &inventoryrepo.Database{
		DB: sqlxDB,
	}

	tests := []struct {
		name          string
		req           *request.CreateInventoryRequest
		mockBehavior  func()
		expectedError bool
	}{
		{
			name: "Successful Inventory Creation",
			req: &request.CreateInventoryRequest{
				Stock: 10,
				Size:  "M",
			},
			mockBehavior: func() {
				mock.ExpectQuery(`INSERT INTO inventory \(id, stock, size, available_stock, reserved_stock, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) RETURNING id, stock, size, available_stock, reserved_stock, created_at, updated_at`).
					WithArgs(
						sqlmock.AnyArg(),
						10,
						"M",
						10,
						0,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnRows(sqlmock.NewRows([]string{"id", "stock", "size", "available_stock", "reserved_stock", "created_at", "updated_at"}).
						AddRow("INV123", 10, "M", 10, 0, time.Now(), time.Now()))
			},
			expectedError: false,
		},
		{
			name: "Invalid Stock Value",
			req: &request.CreateInventoryRequest{
				Stock: -5,
				Size:  "M",
			},
			mockBehavior:  func() {},
			expectedError: true,
		},
		{
			name: "Invalid Size Value",
			req: &request.CreateInventoryRequest{
				Stock: 10,
				Size:  "",
			},
			mockBehavior:  func() {},
			expectedError: true,
		},
		{
			name: "Database Error",
			req: &request.CreateInventoryRequest{
				Stock: 10,
				Size:  "M",
			},
			mockBehavior: func() {
				mock.ExpectQuery(`INSERT INTO inventory \(id, stock, size, available_stock, reserved_stock, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) RETURNING id, stock, size, available_stock, reserved_stock, created_at, updated_at`).
					WithArgs(
						sqlmock.AnyArg(),
						10,
						"M",
						10,
						0,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnError(errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			result, err := repo.CreateInventory(context.Background(), tt.req)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Stock, result.Stock)
				assert.Equal(t, tt.req.Size, result.Size)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
func TestUpdateInventory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &inventoryrepo.Database{DB: sqlxDB}

	tests := []struct {
		name          string
		req           *request.UpdateInventoryRequest
		mockBehavior  func()
		expectedError bool
	}{
		{
			name: "Successful Update",
			req: &request.UpdateInventoryRequest{
				ID:             "INV123",
				Size:           "L",
				Stock:          50,
				ReservedStock:  10,
				AvailableStock: 40,
			},
			mockBehavior: func() {
				mock.ExpectQuery(`^UPDATE inventory SET size = \$2, stock = \$3, reserved_stock = \$4, available_stock = \$5, updated_at = NOW\(\) WHERE id = \$1 RETURNING id, stock, size, available_stock, reserved_stock, created_at, updated_at$`).
					WithArgs("INV123", "L", 50, 10, 40).
					WillReturnRows(sqlmock.NewRows([]string{"id", "stock", "size", "available_stock", "reserved_stock", "created_at", "updated_at"}).
						AddRow("INV123", 50, "L", 40, 10, time.Now(), time.Now()))
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			result, err := repo.UpdateInventory(context.Background(), tt.req)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.ID, result.ID)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
