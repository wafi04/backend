package server

import (
	"github.com/gin-gonic/gin"
	config "github.com/wafi04/backend/config/development"
	authhandler "github.com/wafi04/backend/internal/handler/auth"
	categoryhandler "github.com/wafi04/backend/internal/handler/category"
	producthandler "github.com/wafi04/backend/internal/handler/product"
	"github.com/wafi04/backend/pkg/middleware"
	httpresponse "github.com/wafi04/backend/pkg/response"
	"github.com/wafi04/backend/pkg/utils"
)

func Allroutes(
	authHandler *authhandler.AuthHandler,
	categoryHandler *categoryhandler.CategoryHandler,
	producthandler *producthandler.ProductHandler,
) *gin.Engine {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	config.SetupCORS(r)

	r.Use(httpresponse.ResponseTimeMiddleware())

	r.GET("/health", utils.ConnectionHealthy)

	public := r.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.CreateUser)
			auth.POST("/login", authHandler.Login)

		}
	}
	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		user := protected.Group("/user")
		{
			user.GET("/profile", authHandler.GetUser)
			user.GET("/verify-email", authHandler.VerifyEmail)
			user.POST("/resend-verification", authHandler.ResendVerification)
			user.POST("/logout", authHandler.Logout)
			user.POST("/revoke-session", authHandler.RevokeSession)
			user.POST("/refresh-token", authHandler.RefreshToken)
			user.GET("/sessions", authHandler.ListSessions)
		}
		category := protected.Group("/category")
		{
			category.GET("/list-categories", categoryHandler.HandleGetCategory)
			category.POST("", categoryHandler.HandleCreateCategory)
			category.PUT("/update/:id", categoryHandler.HandleUpdateCategory)
			category.DELETE("/:id", categoryHandler.HandleDeleteCategory)
		}
		product := protected.Group("/product")
		{
			product.POST("", producthandler.HandleCreateProduct)
			product.GET("/:id", producthandler.HandleGetProduct)
			product.GET("/all", producthandler.HandleListProducts)
			product.PUT("/:id", producthandler.HandleUpdateProduct)
			product.DELETE("{id}", producthandler.HandleDeleteProduct)

			// variants
			product.POST("/:id/variant", producthandler.HandleCreateVariants)
			product.PUT("/:id/variant", producthandler.HandleUpdateVariants)
			product.GET("/:id/variant", producthandler.HandleGetProductVariant)
			product.GET("/:id/variants", producthandler.HandleGetProductVariants)
			product.DELETE("/:id/variant", producthandler.HandleDeleteVariants)

			// images
			product.POST("/:id/variant/images", producthandler.HandleAddProductImage)
			product.DELETE("/:id/variant/images", producthandler.HandleDeleteProductImage)
		}
	}

	return r
}
