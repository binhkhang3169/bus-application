package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler handles health check requests for the API gateway.
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "API Gateway is healthy and operational 🚀",
		"data": gin.H{
			"timestamp": time.Now().Format(time.RFC3339Nano),
		},
	})
}
