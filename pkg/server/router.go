package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	config "github.com/wafi04/backend/config/development"
	authhandler "github.com/wafi04/backend/internal/handler/auth"
	categoryhandler "github.com/wafi04/backend/internal/handler/category"
	inventoryhandler "github.com/wafi04/backend/internal/handler/inventory"
	producthandler "github.com/wafi04/backend/internal/handler/product"
	"github.com/wafi04/backend/pkg/middleware"
	httpresponse "github.com/wafi04/backend/pkg/response"
	"github.com/wafi04/backend/pkg/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections (you can customize this for security)
		return true
	},
}

func Allroutes(
	authHandler *authhandler.AuthHandler,
	categoryHandler *categoryhandler.CategoryHandler,
	producthandler *producthandler.ProductHandler,
	inventoryhandler *inventoryhandler.InventoryHandler,
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
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			log.Printf("Received: %s", message)

			if err := conn.WriteMessage(messageType, message); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	})
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
		inv := protected.Group("/stock")
		{
			inv.GET("/:id", inventoryhandler.HandleGetInvetory)
		}
	}

	return r
}
