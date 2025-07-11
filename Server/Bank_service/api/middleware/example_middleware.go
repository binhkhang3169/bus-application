package middleware

import (
	// "log"
	// "net/http"
	// "strings"
	// "time"

	"github.com/gin-gonic/gin"
	// "bank/token" // Giả sử bạn có package token
)

// LoggingMiddleware ghi log mỗi request.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// start := time.Now()
		// path := c.Request.URL.Path
		// raw := c.Request.URL.RawQuery

		c.Next() // Xử lý request

		// latency := time.Since(start)
		// clientIP := c.ClientIP()
		// method := c.Request.Method
		// statusCode := c.Writer.Status()
		// errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// bodySize := c.Writer.Size()

		// if raw != "" {
		// 	path = path + "?" + raw
		// }

		// log.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %#v %s\n",
		// 	latency,
		// 	statusCode,
		// 	bodySize,
		// 	clientIP,
		// 	method,
		// 	path,
		// 	errorMessage,
		// )
	}
}

// AuthMiddleware kiểm tra token xác thực.
// func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
// 			return
// 		}

// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
// 			return
// 		}

// 		tokenString := parts[1]
// 		payload, err := token.VerifyToken(tokenString, jwtSecret) // Giả sử bạn có hàm VerifyToken
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			return
// 		}

// 		// Lưu payload vào context để các handler sau có thể sử dụng
// 		c.Set("userID", payload.UserID) // Ví dụ: payload.UserID
// 		c.Set("username", payload.Username)

// 		c.Next()
// 	}
// }
