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
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
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
					c.Set("email", claims.Email)
				}
			}
		}
		c.Next()
	}
}

// AdminMiddleware validates that the user has admin or owner role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check if user is authenticated
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"code":    models.CODE_UNAUTHORIZED,
				"message": models.MSG_MISSING_TOKEN,
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

		// Check if user role is admin or owner
		if claims.Role != models.ROLE_ADMIN && claims.Role != models.ROLE_OWNER {
			c.JSON(403, gin.H{
				"code":    403,
				"message": models.MSG_FORBIDDEN,
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// OwnerMiddleware validates that the user has owner role only
func OwnerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check if user is authenticated
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"code":    models.CODE_UNAUTHORIZED,
				"message": models.MSG_MISSING_TOKEN,
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

		// Check if user role is owner
		if claims.Role != models.ROLE_OWNER {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "Only owner can perform this action",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}
