package middleware

import (
	"bytes"
	"io"
	"strconv"
	"time"

	"mbankingcore/handlers"

	"github.com/gin-gonic/gin"
)

// AuditLogMiddleware logs all API requests for audit trail
func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Read request body for logging
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// Process request
		c.Next()

		// Log after request processing
		latency := time.Since(start)

		// Get user/admin info from context
		userID, userExists := c.Get("user_id")
		adminID, adminExists := c.Get("admin_id")

		// Get client info
		clientIP := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		method := c.Request.Method
		endpoint := c.Request.URL.Path
		statusCode := c.Writer.Status()

		// Determine entity type based on endpoint
		entityType := determineEntityType(endpoint)

		// Get entity ID from URL parameters if available
		var entityIDValue uint = 0
		if idParam := c.Param("id"); idParam != "" {
			if id, err := strconv.ParseUint(idParam, 10, 32); err == nil {
				entityIDValue = uint(id)
			}
		}

		// Determine action based on HTTP method
		action := determineAction(method)

		// Log the activity
		if userExists {
			if userIDValue, ok := userID.(uint); ok {
				handlers.LogUserActivity(
					userIDValue,
					action,
					entityType,
					entityIDValue,
					clientIP,
					userAgent,
					endpoint,
					method,
					statusCode,
				)
			}
		} else if adminExists {
			if adminIDValue, ok := adminID.(uint); ok {
				handlers.LogAdminActivity(
					adminIDValue,
					action,
					entityType,
					entityIDValue,
					clientIP,
					userAgent,
					endpoint,
					method,
					statusCode,
				)
			}
		}

		// Add audit info to response headers for debugging (optional)
		if gin.Mode() == gin.DebugMode {
			c.Header("X-Audit-Latency", latency.String())
			c.Header("X-Audit-Entity-Type", entityType)
			c.Header("X-Audit-Action", action)
		}
	}
}

// determineEntityType extracts entity type from API endpoint
func determineEntityType(endpoint string) string {
	switch {
	case contains(endpoint, "/users"):
		return "user"
	case contains(endpoint, "/transactions"):
		return "transaction"
	case contains(endpoint, "/bank-accounts"):
		return "bank_account"
	case contains(endpoint, "/articles"):
		return "article"
	case contains(endpoint, "/photos"):
		return "photo"
	case contains(endpoint, "/admins"):
		return "admin"
	case contains(endpoint, "/config"):
		return "config"
	case contains(endpoint, "/onboardings"):
		return "onboarding"
	case contains(endpoint, "/auth"):
		return "auth"
	case contains(endpoint, "/profile"):
		return "profile"
	case contains(endpoint, "/sessions"):
		return "session"
	default:
		return "system"
	}
}

// determineAction maps HTTP method to audit action
func determineAction(method string) string {
	switch method {
	case "GET":
		return "READ"
	case "POST":
		return "CREATE"
	case "PUT", "PATCH":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexOf(s, substr) != -1)))
}

// indexOf returns the index of substr in s, or -1 if not found
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// AuditLoginMiddleware specifically logs login/logout activities
func AuditLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store request info before processing
		clientIP := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		endpoint := c.Request.URL.Path

		// Process request
		c.Next()

		// Log login activity after processing
		statusCode := c.Writer.Status()

		// Determine login type and status
		var loginType, status string
		var userID, adminID *uint

		if contains(endpoint, "/admin/login") {
			loginType = "admin_login"
			if statusCode == 200 {
				status = "success"
				if adminIDValue, exists := c.Get("admin_id"); exists {
					if id, ok := adminIDValue.(uint); ok {
						adminID = &id
					}
				}
			} else {
				status = "failed"
			}
		} else if contains(endpoint, "/login") {
			loginType = "user_login"
			if statusCode == 200 {
				status = "success"
				if userIDValue, exists := c.Get("user_id"); exists {
					if id, ok := userIDValue.(uint); ok {
						userID = &id
					}
				}
			} else {
				status = "failed"
			}
		} else if contains(endpoint, "/logout") {
			if contains(endpoint, "/admin") {
				loginType = "admin_logout"
			} else {
				loginType = "user_logout"
			}
			status = "success"

			// Get user/admin ID from context for logout
			if userIDValue, exists := c.Get("user_id"); exists {
				if id, ok := userIDValue.(uint); ok {
					userID = &id
				}
			}
			if adminIDValue, exists := c.Get("admin_id"); exists {
				if id, ok := adminIDValue.(uint); ok {
					adminID = &id
				}
			}
		} else {
			return // Not a login/logout endpoint
		}

		// Get failure reason if failed
		failureReason := ""
		if status == "failed" {
			if msg, exists := c.Get("failure_reason"); exists {
				if reason, ok := msg.(string); ok {
					failureReason = reason
				}
			}
		}

		// Log the login activity
		handlers.LogLoginActivity(userID, adminID, loginType, status, clientIP, userAgent, failureReason)
	}
}
