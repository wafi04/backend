package main

import (
	"github.com/cloudinary/cloudinary-go/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	config "github.com/wafi04/backend/config/development"
	authhandler "github.com/wafi04/backend/internal/handler/auth"
	categoryhandler "github.com/wafi04/backend/internal/handler/category"
	producthandler "github.com/wafi04/backend/internal/handler/product"
	authrepository "github.com/wafi04/backend/internal/repository/auth"
	"github.com/wafi04/backend/internal/repository/category"
	productRepository "github.com/wafi04/backend/internal/repository/product"
	authservice "github.com/wafi04/backend/internal/service/auth"
	service "github.com/wafi04/backend/internal/service/category"
	"github.com/wafi04/backend/internal/service/files"
	productservice "github.com/wafi04/backend/internal/service/product"
	"github.com/wafi04/backend/pkg/logger"

	"github.com/wafi04/backend/pkg/server"
)

func main() {
	log := logger.NewLogger()
	if err := config.LoadConfig("development"); err != nil {
		log.Log(logger.ErrorLevel, "Error loading config: %v", err)
	}

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

	userRepo := authrepository.NewDB(db.DB)
	userService := authservice.NewAuthService(userRepo)
	categoryRepo := category.NewCategoryRepository(db.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	productrepo := productRepository.NewCategoryRepository(db.DB)
	productservice := productservice.NewProductService(productrepo)

	filesService := files.NewCloudinaryService(cld)
	authHandler := authhandler.NewAuthHandler(userService)
	categoryhandler := categoryhandler.NewCategoryHandler(categoryService, filesService)
	producthandler := producthandler.NewProductHandler(productservice, filesService)
	router := server.Allroutes(authHandler, categoryhandler, producthandler)

	log.Info("Starting server on : %s", config.LoadEnv("PORT"))
	if err := router.Run(config.LoadEnv("PORT")); err != nil {
		log.Log(logger.ErrorLevel, "Failed to start server: %s", err)
	}
}
