package controllers

import (
	"net/http"
	"ticket-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// TokenTestController handles routes for testing token functionality
type TokenTestController struct {
	auth *utils.Auth
}

// NewTokenTestController creates a new instance of TokenTestController
func NewTokenTestController(auth *utils.Auth) *TokenTestController {
	return &TokenTestController{
		auth: auth,
	}
}

// ValidateTokenHandler is a route specifically for testing token validation
func (t *TokenTestController) ValidateTokenHandler(c *gin.Context) {
	// Extract token from Authorization header
	token := c.GetHeader("Authorization")
	if token == "" || len(token) < 7 || token[:7] != "Bearer " {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid or missing token",
			"data":    nil,
		})
		return
	}

	// Remove "Bearer " prefix
	tokenString := token[7:]

	// Attempt to get customer ID
	customerID, err := t.auth.GetCustomerIDFromJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Token validation failed: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Check token expiration
	isExpired := t.auth.IsTokenExpired(tokenString)

	// Prepare response with token details
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Token is valid",
		"data": gin.H{
			"customer_id": customerID,
			"is_expired":  isExpired,
		},
	})
}
