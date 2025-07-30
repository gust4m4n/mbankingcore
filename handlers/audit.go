package handlers

import (
	"net/http"
	"strconv"
	"time"

	"mbankingcore/config"
	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditHandler handles audit trail operations
type AuditHandler struct {
	DB *gorm.DB
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler() *AuditHandler {
	return &AuditHandler{
		DB: config.DB,
	}
}

// GetAuditLogs retrieves audit logs with filtering and pagination
// @Summary Get audit logs
// @Description Retrieve audit logs with filtering and pagination (Admin only)
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param entity_type query string false "Filter by entity type"
// @Param entity_id query int false "Filter by entity ID"
// @Param user_id query int false "Filter by user ID"
// @Param admin_id query int false "Filter by admin ID"
// @Param action query string false "Filter by action"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param ip_address query string false "Filter by IP address"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} models.APIResponse{data=models.AuditResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/admin/audit/logs [get]
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	// Parse query parameters
	req := &models.AuditRequest{}

	if entityType := c.Query("entity_type"); entityType != "" {
		req.EntityType = entityType
	}
	if entityID := c.Query("entity_id"); entityID != "" {
		if id, err := strconv.ParseUint(entityID, 10, 32); err == nil {
			req.EntityID = uint(id)
		}
	}
	if userID := c.Query("user_id"); userID != "" {
		if id, err := strconv.ParseUint(userID, 10, 32); err == nil {
			req.UserID = uint(id)
		}
	}
	if adminID := c.Query("admin_id"); adminID != "" {
		if id, err := strconv.ParseUint(adminID, 10, 32); err == nil {
			req.AdminID = uint(id)
		}
	}
	if action := c.Query("action"); action != "" {
		req.Action = action
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			req.StartDate = date
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			req.EndDate = date.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}
	if ipAddress := c.Query("ip_address"); ipAddress != "" {
		req.IPAddress = ipAddress
	}

	// Parse pagination
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	req.Page = page

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	req.Limit = limit

	// Get audit logs
	response, err := models.GetAuditLogs(h.DB, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "Failed to retrieve audit logs"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(200, "Audit logs retrieved successfully", response))
}

// GetLoginAuditLogs retrieves login audit logs with filtering and pagination
// @Summary Get login audit logs
// @Description Retrieve login audit logs with filtering and pagination (Admin only)
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query int false "Filter by user ID"
// @Param admin_id query int false "Filter by admin ID"
// @Param login_type query string false "Filter by login type"
// @Param status query string false "Filter by status"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param ip_address query string false "Filter by IP address"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} models.APIResponse{data=models.LoginAuditResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/admin/audit/login [get]
func (h *AuditHandler) GetLoginAuditLogs(c *gin.Context) {
	// Parse query parameters
	req := &models.LoginAuditRequest{}

	if userID := c.Query("user_id"); userID != "" {
		if id, err := strconv.ParseUint(userID, 10, 32); err == nil {
			req.UserID = uint(id)
		}
	}
	if adminID := c.Query("admin_id"); adminID != "" {
		if id, err := strconv.ParseUint(adminID, 10, 32); err == nil {
			req.AdminID = uint(id)
		}
	}
	if loginType := c.Query("login_type"); loginType != "" {
		req.LoginType = loginType
	}
	if status := c.Query("status"); status != "" {
		req.Status = status
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			req.StartDate = date
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			req.EndDate = date.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}
	if ipAddress := c.Query("ip_address"); ipAddress != "" {
		req.IPAddress = ipAddress
	}

	// Parse pagination
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	req.Page = page

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	req.Limit = limit

	// Get login audit logs
	response, err := models.GetLoginAuditLogs(h.DB, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(500, "Failed to retrieve login audit logs"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(200, "Login audit logs retrieved successfully", response))
}

// LogUserActivity creates an audit log entry for user activities
func LogUserActivity(userID uint, action, entityType string, entityID uint, ipAddress, userAgent, endpoint, method string, statusCode int) {
	auditLog := models.AuditLog{
		UserID:        &userID,
		Action:        action,
		EntityType:    entityType,
		EntityID:      entityID,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		APIEndpoint:   endpoint,
		RequestMethod: method,
		StatusCode:    statusCode,
	}

	config.DB.Create(&auditLog)
}

// LogAdminActivity creates an audit log entry for admin activities
func LogAdminActivity(adminID uint, action, entityType string, entityID uint, ipAddress, userAgent, endpoint, method string, statusCode int) {
	auditLog := models.AuditLog{
		AdminID:       &adminID,
		Action:        action,
		EntityType:    entityType,
		EntityID:      entityID,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		APIEndpoint:   endpoint,
		RequestMethod: method,
		StatusCode:    statusCode,
	}

	config.DB.Create(&auditLog)
}

// LogLoginActivity creates a login audit entry
func LogLoginActivity(userID *uint, adminID *uint, loginType, status, ipAddress, userAgent, failureReason string) {
	loginAudit := models.LoginAudit{
		UserID:        userID,
		AdminID:       adminID,
		LoginType:     loginType,
		Status:        status,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		FailureReason: failureReason,
	}

	config.DB.Create(&loginAudit)
}
