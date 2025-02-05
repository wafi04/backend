package main

import (
	"github.com/cloudinary/cloudinary-go/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	config "github.com/wafi04/backend/config/development"
	authhandler "github.com/wafi04/backend/services/auth/handler"
	authrepo "github.com/wafi04/backend/services/auth/repository"
	authservice "github.com/wafi04/backend/services/auth/service"
	"github.com/wafi04/backend/services/cart"
	categoryhandler "github.com/wafi04/backend/services/category/handler"
	category "github.com/wafi04/backend/services/category/repository"
	"github.com/wafi04/backend/services/category/service"
	"github.com/wafi04/backend/services/files"
	"github.com/wafi04/backend/services/inventory"
	producthandler "github.com/wafi04/backend/services/product/handler"
	productRepository "github.com/wafi04/backend/services/product/repository"
	productservice "github.com/wafi04/backend/services/product/service"
	"github.com/wafi04/backend/services/user"

	"github.com/wafi04/backend/pkg/logger"

	"github.com/wafi04/backend/pkg/server"
)

func main() {
	log := logger.NewLogger()
	if err := config.LoadConfig("development"); err != nil {
		log.Log(logger.ErrorLevel, "Error loading config: %v", err)
	}

	log.Info("gaiaiaia")
	db, err := config.NewDB()
	if err != nil {
		log.Log(logger.ErrorLevel, "Database connection failed: %v", err)
	}
	defer db.Close()

	cld, err := cloudinary.NewFromParams(
		config.LoadEnv("CLOUDINARY_CLOUD_NAME"),
		config.LoadEnv("CLOUDINARY_API_KEY"),
		config.LoadEnv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Log(logger.ErrorLevel, "Failed to initialize Cloudinary: %v", err)
		return
	}

	// Check database health
	health := db.Health()
	log.Log(logger.InfoLevel, "Database health: %v", health["status"])

	userRepo := authrepo.NewDB(db.DB)
	userService := authservice.NewAuthService(userRepo)
	categoryRepo := category.NewCategoryRepository(db.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	productrepo := productRepository.NewProductRepository(db.DB)
	productservice := productservice.NewProductService(productrepo)
	inventoryrepo := inventory.NewInventoryRepository(db.DB)
	inventoryService := inventory.NewInventoryService(inventoryrepo)
	cartrepo := cart.NewCartRepository(db.DB)
	cartService := cart.NewCartService(cartrepo)
	userrepos := user.NewUserRepository(db.DB)
	shipAddrrepo := user.NewShippingAddressRepo(db.DB)

	userHandler := user.NewUserHandler(userrepos)
	filesService := files.NewCloudinaryService(cld)
	authHandler := authhandler.NewAuthHandler(userService)
	categoryhandler := categoryhandler.NewCategoryHandler(categoryService, filesService)
	producthandler := producthandler.NewProductHandler(productservice, filesService)
	inventoryHandler := inventory.NewInventoryHandler(inventoryService)
	cartHandler := cart.NewCartHandler(cartService)
	shiphnadler := user.NewShippingHandler(shipAddrrepo)

	router := server.Allroutes(authHandler, userHandler, categoryhandler, producthandler, inventoryHandler, cartHandler, shiphnadler)

	log.Info("Starting server on : %s", config.LoadEnv("PORT"))
	if err := router.Run(":8080"); err != nil {
		log.Log(logger.ErrorLevel, "Failed to start server: %s", err)
	}
}
