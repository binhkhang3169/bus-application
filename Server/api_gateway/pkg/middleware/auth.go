// file: pkg/middleware/auth.go
package middleware

import (
	"api_gateway/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// === HÀM KẾT HỢP MIDDLEWARE (ĐÃ THÊM) ===

// AuthMiddleware là một hàm tiện ích kết hợp cả hai bước xác thực và phân quyền.
// Nó trả về một chuỗi các gin.HandlerFunc để áp dụng cho các route group.
// 1. OptionalAuthMiddleware: Luôn xác định danh tính (guest hoặc user đã login).
// 2. RBACMiddleware: Dựa trên danh tính đó để kiểm tra quyền truy cập.
func AuthMiddleware(authService *utils.Auth, routePermissions map[string][]string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		OptionalAuthMiddleware(authService),
		RBACMiddleware(routePermissions),
	}
}

// === MIDDLEWARE 1: XÁC ĐỊNH DANH TÍNH (AUTHENTICATION) ===

// OptionalAuthMiddleware xác định danh tính người dùng một cách "mềm".
// - Nếu có token hợp lệ, nó sẽ gán thông tin user vào context.
// - Nếu không có token, nó sẽ gán vai trò ROLE_GUEST và userID = 0.
// - Nó chỉ từ chối request nếu token được cung cấp nhưng KHÔNG HỢP LỆ (hết hạn, sai chữ ký).
func OptionalAuthMiddleware(authService *utils.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// --- Trường hợp 1: Khách vãng lai (không có token) ---
		if authHeader == "" {
			c.Set("userID", 0) // ID 0 cho khách
			c.Set("userRole", "ROLE_GUEST")
			log.Printf("Identity: No token. Identified as ROLE_GUEST for path '%s'", c.Request.URL.Path)
			c.Next()
			return
		}

		// --- Trường hợp 2: Có token, tiến hành xác thực ---
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			// Token có nhưng sai định dạng -> Vẫn coi là khách và cho đi tiếp
			c.Set("userID", 0)
			c.Set("userRole", "ROLE_GUEST")
			log.Printf("Identity: Malformed token. Identified as ROLE_GUEST for path '%s'", c.Request.URL.Path)
			c.Next()
			return
		}
		tokenString := parts[1]

		validatedTokenInfo, err := authService.ValidateToken(tokenString)
		if err != nil {
			// Token có nhưng KHÔNG HỢP LỆ -> Đây là lỗi, phải từ chối request.
			log.Printf("Identity: Token validation error for path '%s': %v", c.Request.URL.Path, err)
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "Invalid or expired token", "details": err.Error()})
			c.Abort()
			return
		}

		// --- Trường hợp 3: Token hợp lệ ---
		c.Set("userID", validatedTokenInfo.UserID)
		c.Set("userRole", validatedTokenInfo.Role)
		log.Printf("Identity: Token validated. UserID: %d, Role: '%s', Path: '%s'", validatedTokenInfo.UserID, validatedTokenInfo.Role, c.Request.URL.Path)
		c.Next()
	}
}

// === MIDDLEWARE 2: KIỂM TRA QUYỀN HẠN (AUTHORIZATION - RBAC) ===

func isRoleAllowed(userRole string, allowedRoles []string) bool {
	for _, r := range allowedRoles {
		if r == userRole {
			return true
		}
	}
	return false
}

// RBACMiddleware kiểm tra xem vai trò của người dùng (đã được OptionalAuthMiddleware xác định)
// có được phép truy cập vào route hiện tại hay không.
func RBACMiddleware(routePermissions map[string][]string) gin.HandlerFunc {
	var sortedRulePrefixes []string
	for prefix := range routePermissions {
		sortedRulePrefixes = append(sortedRulePrefixes, prefix)
	}
	sort.Slice(sortedRulePrefixes, func(i, j int) bool {
		return len(sortedRulePrefixes[i]) > len(sortedRulePrefixes[j])
	})

	log.Printf("RBAC: Middleware configured with rule prefixes: %v", sortedRulePrefixes)

	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path

		userRoleVal, _ := c.Get("userRole")
		userRole := userRoleVal.(string)

		// Tìm rule phù hợp nhất cho path hiện tại
		var bestMatchRulePrefix string
		var rolesForBestMatch []string

		for _, rulePrefix := range sortedRulePrefixes {
			if strings.HasPrefix(requestPath, rulePrefix) {
				bestMatchRulePrefix = rulePrefix
				rolesForBestMatch = routePermissions[rulePrefix]
				break
			}
		}

		// Nếu có rule cho path này, kiểm tra quyền. Nếu không có rule, mặc định cho qua.
		if bestMatchRulePrefix != "" {
			if !isRoleAllowed(userRole, rolesForBestMatch) {
				log.Printf("RBAC: Role '%s' FORBIDDEN for path '%s' (rule '%s' matched, role not in %v).", userRole, requestPath, bestMatchRulePrefix, rolesForBestMatch)
				c.JSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "message": fmt.Sprintf("Access denied. Role ('%s') not authorized for this resource.", userRole)})
				c.Abort()
				return
			}
			log.Printf("RBAC: Role '%s' ALLOWED for path '%s' (rule '%s' matched).", userRole, requestPath, bestMatchRulePrefix)
		} else {
			log.Printf("RBAC: No specific rule for path '%s'. Access GRANTED by default for role '%s'.", requestPath, userRole)
		}

		c.Next()
	}
}
