package middleware

import (
	"strings"

	"mbankingcore/models"
	"mbankingcore/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and sets user info in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"code":    models.CODE_UNAUTHORIZED,
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{
				"code":    models.CODE_UNAUTHORIZED,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(401, gin.H{
				"code":    models.CODE_UNAUTHORIZED,
				"message": "Token required",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    models.CODE_UNAUTHORIZED,
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("phone", claims.Phone)
		c.Next()
	}
}

// OptionalAuth middleware that doesn't abort if no token is provided
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString != "" {
				claims, err := utils.ValidateJWT(tokenString)
				if err == nil {
					c.Set("userID", claims.UserID)
					c.Set("phone", claims.Phone)
				}
			}
		}
		c.Next()
	}
}
