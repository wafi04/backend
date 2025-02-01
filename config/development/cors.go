package config

import (
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func SetupCORS(r *gin.Engine) {
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://192.168.100.6:3000"}
    config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{
        "Origin",
        "Content-Type",
        "Accept",
        "Authorization",
        "Cookie",
        "X-Requested-With",
        "Access-Control-Allow-Credentials", 
    }
    config.AllowCredentials = true
    config.ExposeHeaders = []string{
        "Content-Length",
        "Access-Control-Allow-Origin",
        "Access-Control-Allow-Headers",
        "Access-Control-Allow-Credentials",
    }
    r.Use(cors.New(config))
}