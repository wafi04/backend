package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/wafi04/backend/pkg/utils"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(connStr string) (*Database, error) {
	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

type ProductImage struct {
	ID        string `db:"id"`
	URL       string `db:"url"`
	VariantID string `db:"variant_id"`
	IsMain    bool   `db:"is_main"`
}

func (d *Database) AddProductImages(variantID string, imageURLs []string) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, imageURL := range imageURLs {
		image := ProductImage{
			ID:        utils.GenerateRandomId("IMG"),
			URL:       imageURL,
			VariantID: variantID,
			IsMain:    i == 0,
		}

		_, err := tx.NamedExec(`
            INSERT INTO product_images (id, url, variant_id, is_main)
            VALUES (:id, :url, :variant_id, :is_main)
        `, image)

		if err != nil {
			return fmt.Errorf("failed to insert image %s: %w", imageURL, err)
		}

		log.Printf("Inserted image %s for variant %s\n", imageURL, variantID)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Successfully seeded all images!")
	return nil
}
func (d *Database) AddInventoryVariants() error {
	// Buka file JSON
	file, err := os.Open("inventory.json")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Baca isi file JSON
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON ke slice struct
	var inventory []struct {
		VariantID string `json:"variant_id"`
		Size      string `json:"size"`
		Stock     int8   `json:"stock"`
	}
	if err := json.Unmarshal(byteValue, &inventory); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Mulai transaksi
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Loop melalui setiap item inventaris
	for _, inv := range inventory {
		query := `
            INSERT INTO inventory (
                id, variant_id, size, stock, available_stock, reserved_stock, created_at, updated_at
            ) VALUES (
                $1, $2, $3, $4, $5, $6, NOW(), NOW()
            )
        `
		// Generate ID unik untuk setiap entri
		id := utils.GenerateRandomId("INV")

		// Eksekusi query
		_, err := tx.Exec(query,
			id,
			inv.VariantID,
			inv.Size,
			inv.Stock,
			inv.Stock, // available_stock = stock
			0,         // reserved_stock = 0 (default)
		)
		if err != nil {
			return fmt.Errorf("failed to insert inventory for variant %s, size %s: %w", inv.VariantID, inv.Size, err)
		}

		log.Printf("Inserted inventory for variant %s, size %s\n", inv.VariantID, inv.Size)
	}

	// Commit transaksi
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Successfully seeded all inventory!")
	return nil
}

func main() {
	connStr := "postgres://postgres:postgres@192.168.100.9:5432/backend?sslmode=disable&search_path=public"
	db, err := NewDatabase(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	image := []string{
		"https://static.nike.com/a/images/c_limit,w_592,f_auto/t_product_v1/34dc1f16-c254-4022-8813-2953827643a1/W+NIKE+REVOLUTION+7.png",
		"https://static.nike.com/a/images/t_default/34dc1f16-c254-4022-8813-2953827643a1/W+NIKE+REVOLUTION+7.png",
		"https://static.nike.com/a/images/t_default/f35822e9-d199-4ca6-96c9-4c88b6ff1443/W+NIKE+REVOLUTION+7.png",
		"https://static.nike.com/a/images/t_default/60c18f19-1324-4bdb-a6b4-fc3afe7efd0b/W+NIKE+REVOLUTION+7.png",
		"https://static.nike.com/a/images/t_default/cd797719-ff30-4eda-a97f-ac76eb2ee646/W+NIKE+REVOLUTION+7.png",
		"https://static.nike.com/a/images/t_default/5d4d4520-f897-4e97-9d69-25732d927c23/W+NIKE+REVOLUTION+7.png",
		"https://static.nike.com/a/images/t_default/295d9957-2b74-4fa7-af3d-9e366a82461c/W+NIKE+REVOLUTION+7.png",
		"https://static.nike.com/a/images/t_default/5b6a828d-6a3b-488b-bf5b-d2f11bb7ebbc/W+NIKE+REVOLUTION+7.png",
		"https://static.nike.com/a/images/t_default/28d20a58-53b4-46e6-859a-4ef7e7e23d16/W+NIKE+REVOLUTION+7.png",
		"https://static.nike.com/a/images/t_default/c1512ee0-3c91-43a3-a834-3b5ccde08ddd/W+NIKE+REVOLUTION+7.png",
	}
	variantID := "VAR-602961174156"

	if err := db.AddProductImages(variantID, image); err != nil {
		log.Fatalf("Failed to add inventory variants: %v", err)
	}
}
