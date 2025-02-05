package server

import (
	"github.com/gin-gonic/gin"
	config "github.com/wafi04/backend/config/development"
	authhandler "github.com/wafi04/backend/services/auth/handler"
	"github.com/wafi04/backend/services/cart"
	categoryhandler "github.com/wafi04/backend/services/category/handler"
	"github.com/wafi04/backend/services/inventory"
	producthandler "github.com/wafi04/backend/services/product/handler"
	"github.com/wafi04/backend/services/user"

	"github.com/wafi04/backend/pkg/middleware"
	httpresponse "github.com/wafi04/backend/pkg/response"
	"github.com/wafi04/backend/pkg/types"
	"github.com/wafi04/backend/pkg/utils"
)

func Allroutes(
	authHandler *authhandler.AuthHandler,
	userhandler *user.UserHandler,
	categoryHandler *categoryhandler.CategoryHandler,
	producthandler *producthandler.ProductHandler,
	inventoryhandler *inventory.InventoryHandler,
	carthandler *cart.CartHandler,
	shippingHandler *user.ShippingHandler,
) *gin.Engine {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	config.SetupCORS(r)
	r.Use(httpresponse.ResponseTimeMiddleware())
	r.GET("/ws", WebSocketHandler)

	go BroadcastMessages()
	types.Broadcast <- "Hello, WebSocket clients!"

	r.GET("/health", utils.ConnectionHealthy)

	public := r.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.CreateUser)
			auth.POST("/login", authHandler.Login)

		}
	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		user := protected.Group("/user")
		{
			user.GET("/profile", authHandler.GetUser)
			user.GET("/verify", authHandler.Verify)
			user.GET("/verify-email", authHandler.VerifyEmail)
			user.POST("/resend-verification", authHandler.ResendVerification)
			user.POST("/logout", authHandler.Logout)
			user.POST("/revoke-session", authHandler.RevokeSession)
			user.POST("/refresh-token", authHandler.RefreshToken)
			user.GET("/sessions", authHandler.ListSessions)
			user.GET("/details", userhandler.GetUserDetails)
			user.POST("/details", userhandler.HandleCreateUserDetails)
			user.PATCH("/details/:id", userhandler.HandleUpdateProfiles)
			user.POST("/address", shippingHandler.CreateAddressReq)
			user.GET("/address", shippingHandler.GetAll)
			user.PATCH("/address/:id", shippingHandler.UpdateShipping)
			user.DELETE("/address/:id", shippingHandler.DeleteShipping)

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
			public.GET("/product/all", producthandler.HandleListProducts)
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
		inv := protected.Group("/stock")
		{
			inv.GET("/:id", inventoryhandler.HandleGetInvetory)
		}

		cart := protected.Group("/cart")
		{
			cart.POST("", carthandler.HandleAddToCart)
			cart.GET("", carthandler.HandleGetCart)
			cart.DELETE("/clear", carthandler.ClearCart)
			cart.PATCH("/items/:id", carthandler.UpdateQuantity)
			cart.DELETE("/items/:id", carthandler.RemoveFromCart)
		}

	}

	return r
}
