package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)
func ConnectionHealthy(c *gin.Context) {
	serviceReady := true
	timestamp := time.Now().Format(time.RFC3339)
	data := struct {
		Ready     bool   `json:"ready"`
		Timestamp string `json:"time"`
	}{
		Ready:     serviceReady,
		Timestamp: timestamp,
	}

	if serviceReady {
		c.JSON(http.StatusOK, gin.H{
			"status":  "oke",
			"message": "Service Ready",
			"data":    data,
	
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Service Not Ready",
		})
	}
}
