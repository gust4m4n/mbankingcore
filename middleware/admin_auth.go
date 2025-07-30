package middleware

import (
	"mbankingcore/models"
	"mbankingcore/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware validates admin JWT tokens
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    models.CODE_UNAUTHORIZED,
				Message: "Authorization header required",
				Data:    nil,
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    models.CODE_INVALID_TOKEN,
				Message: "Invalid authorization header format",
				Data:    nil,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate JWT token
		claims, err := utils.ValidateAdminJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    models.CODE_INVALID_TOKEN,
				Message: "Invalid or expired token",
				Data:    nil,
			})
			c.Abort()
			return
		}

		// Store admin information in context
		c.Set("admin_id", claims.AdminID)
		c.Set("admin_email", claims.Email)
		c.Set("admin_role", claims.Role)

		c.Next()
	}
}

// SuperAdminMiddleware ensures only super admins can access certain endpoints
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("admin_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    models.CODE_UNAUTHORIZED,
				Message: "Admin authentication required",
				Data:    nil,
			})
			c.Abort()
			return
		}

		if role != models.ADMIN_ROLE_SUPER {
			c.JSON(http.StatusForbidden, models.Response{
				Code:    403,
				Message: "Super admin privileges required",
				Data:    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
