package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Product struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	SubTitle    string  `json:"sub_title"`
	Price       float64 `json:"price"`
}

// var (
// 	db      *sqlx.DB
// 	connStr = "postgres://postgres:postgres@192.168.100.9:5432/backend?sslmode=disable&search_path=public"
// )

// func main() {
// 	var err error
// 	db, err = sqlx.Connect("pgx", connStr)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}
// 	defer db.Close()

// 	file, err := os.Open("data.json")
// 	if err != nil {
// 		log.Fatal("Failed to open file:", err)
// 	}
// 	defer file.Close()

// 	byteValue, _ := ioutil.ReadAll(file)

// 	var products []Product
// 	err = json.Unmarshal(byteValue, &products)
// 	if err != nil {
// 		log.Fatal("Failed to parse JSON:", err)
// 	}

// 	for _, product := range products {
// 		id := utils.GenerateRandomId("PROD")
// 		sku := utils.GenerateSku(product.Name)
// 		_, err := db.Exec(
// 			"INSERT INTO products (id,name, description,sub_title,sku,price,category_id,created_at,updated_at) VALUES ($1, $2, $3,$4,$5,$6,$7,$8,$9)",
// 			id, product.Name, product.Description, product.SubTitle, sku, product.Price, product.Category, time.Now(), time.Now(),
// 		)
// 		if err != nil {
// 			log.Println("Failed to insert product:", err)
// 		} else {
// 			fmt.Println("Inserted:", product.Name)
// 		}
// 	}

// 	fmt.Println("Seeding from file completed!")
// }
