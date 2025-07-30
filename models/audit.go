package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Pagination represents pagination information for API responses
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// AuditLog represents an audit log entry for tracking all system changes
type AuditLog struct {
	ID            uint             `json:"id" gorm:"primaryKey"`
	UserID        *uint            `json:"user_id,omitempty"`
	AdminID       *uint            `json:"admin_id,omitempty"`
	EntityType    string           `json:"entity_type" gorm:"not null;size:50"` // 'user', 'transaction', 'admin', etc.
	EntityID      uint             `json:"entity_id" gorm:"not null"`
	Action        string           `json:"action" gorm:"not null;size:20"`         // 'CREATE', 'UPDATE', 'DELETE', 'LOGIN', etc.
	OldValues     *json.RawMessage `json:"old_values,omitempty" gorm:"type:jsonb"` // Data before change
	NewValues     *json.RawMessage `json:"new_values,omitempty" gorm:"type:jsonb"` // Data after change
	IPAddress     string           `json:"ip_address" gorm:"type:inet"`
	UserAgent     string           `json:"user_agent"`
	APIEndpoint   string           `json:"api_endpoint" gorm:"size:255"`
	RequestMethod string           `json:"request_method" gorm:"size:10"`
	StatusCode    int              `json:"status_code"`
	CreatedAt     time.Time        `json:"created_at"`
}

// LoginAudit represents login/logout audit trail
type LoginAudit struct {
	ID            uint             `json:"id" gorm:"primaryKey"`
	UserID        *uint            `json:"user_id,omitempty"`
	AdminID       *uint            `json:"admin_id,omitempty"`
	LoginType     string           `json:"login_type" gorm:"not null;size:20"` // 'user_login', 'admin_login', 'logout'
	Status        string           `json:"status" gorm:"not null;size:20"`     // 'success', 'failed', 'blocked'
	IPAddress     string           `json:"ip_address" gorm:"type:inet"`
	UserAgent     string           `json:"user_agent"`
	DeviceInfo    *json.RawMessage `json:"device_info,omitempty" gorm:"type:jsonb"`
	FailureReason string           `json:"failure_reason,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
}

// AuditRequest represents request structure for audit queries
type AuditRequest struct {
	EntityType string    `form:"entity_type"` // Filter by entity type
	EntityID   uint      `form:"entity_id"`   // Filter by entity ID
	UserID     uint      `form:"user_id"`     // Filter by user ID
	AdminID    uint      `form:"admin_id"`    // Filter by admin ID
	Action     string    `form:"action"`      // Filter by action
	StartDate  time.Time `form:"start_date"`  // Date range start
	EndDate    time.Time `form:"end_date"`    // Date range end
	IPAddress  string    `form:"ip_address"`  // Filter by IP
	Page       int       `form:"page"`        // Pagination
	Limit      int       `form:"limit"`       // Items per page
}

// LoginAuditRequest represents request structure for login audit queries
type LoginAuditRequest struct {
	UserID    uint      `form:"user_id"`    // Filter by user ID
	AdminID   uint      `form:"admin_id"`   // Filter by admin ID
	LoginType string    `form:"login_type"` // Filter by login type
	Status    string    `form:"status"`     // Filter by status
	StartDate time.Time `form:"start_date"` // Date range start
	EndDate   time.Time `form:"end_date"`   // Date range end
	IPAddress string    `form:"ip_address"` // Filter by IP
	Page      int       `form:"page"`       // Pagination
	Limit     int       `form:"limit"`      // Items per page
}

// AuditResponse represents paginated audit response
type AuditResponse struct {
	Logs       []AuditLog `json:"logs"`
	Pagination Pagination `json:"pagination"`
}

// LoginAuditResponse represents paginated login audit response
type LoginAuditResponse struct {
	Logs       []LoginAudit `json:"logs"`
	Pagination Pagination   `json:"pagination"`
}

// CreateAuditLog creates a new audit log entry
func CreateAuditLog(db *gorm.DB, log *AuditLog) error {
	return db.Create(log).Error
}

// CreateLoginAudit creates a new login audit entry
func CreateLoginAudit(db *gorm.DB, log *LoginAudit) error {
	return db.Create(log).Error
}

// GetAuditLogs retrieves audit logs with filtering and pagination
func GetAuditLogs(db *gorm.DB, req *AuditRequest) (*AuditResponse, error) {
	var logs []AuditLog
	var total int64

	// Build query with filters
	query := db.Model(&AuditLog{})

	if req.EntityType != "" {
		query = query.Where("entity_type = ?", req.EntityType)
	}
	if req.EntityID > 0 {
		query = query.Where("entity_id = ?", req.EntityID)
	}
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.AdminID > 0 {
		query = query.Where("admin_id = ?", req.AdminID)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if !req.StartDate.IsZero() {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if !req.EndDate.IsZero() {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	if req.IPAddress != "" {
		query = query.Where("ip_address = ?", req.IPAddress)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Get paginated results
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &AuditResponse{
		Logs: logs,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

// GetLoginAuditLogs retrieves login audit logs with filtering and pagination
func GetLoginAuditLogs(db *gorm.DB, req *LoginAuditRequest) (*LoginAuditResponse, error) {
	var logs []LoginAudit
	var total int64

	// Build query with filters
	query := db.Model(&LoginAudit{})

	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.AdminID > 0 {
		query = query.Where("admin_id = ?", req.AdminID)
	}
	if req.LoginType != "" {
		query = query.Where("login_type = ?", req.LoginType)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if !req.StartDate.IsZero() {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if !req.EndDate.IsZero() {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	if req.IPAddress != "" {
		query = query.Where("ip_address = ?", req.IPAddress)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Get paginated results
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &LoginAuditResponse{
		Logs: logs,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}
