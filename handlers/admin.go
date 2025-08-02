package handlers

import (
	"encoding/json"
	"mbankingcore/models"
	"mbankingcore/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

// AdminLogin handles admin authentication
func (h *AdminHandler) AdminLogin(c *gin.Context) {
	var request models.AdminLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid request data",
			Data:    err.Error(),
		})
		return
	}

	// Find admin by email
	var admin models.Admin
	if err := h.DB.Where("email = ?", request.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_LOGIN_FAILED,
			Message: "Invalid email or password",
			Data:    nil,
		})
		return
	}

	// Check if admin is active
	if !admin.IsActive() {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_LOGIN_FAILED,
			Message: "Admin account is not active",
			Data:    nil,
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_LOGIN_FAILED,
			Message: "Invalid email or password",
			Data:    nil,
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateAdminJWT(admin.ID, admin.Email, admin.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to generate token",
			Data:    nil,
		})
		return
	}

	// Update last login
	now := time.Now()
	admin.LastLogin = &now
	h.DB.Save(&admin)

	c.JSON(http.StatusOK, models.AdminLoginSuccessResponse(admin, token, 24*60*60)) // 24 hours
}

// AdminLogout handles admin logout
func (h *AdminHandler) AdminLogout(c *gin.Context) {
	// For JWT, logout is handled client-side by removing the token
	// You can implement token blacklisting here if needed
	c.JSON(http.StatusOK, models.AdminLogoutSuccessResponse())
}

// CreateAdmin creates a new admin user
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var request models.CreateAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid request data",
			Data:    err.Error(),
		})
		return
	}

	// Check if email already exists
	var existingAdmin models.Admin
	if err := h.DB.Where("email = ?", request.Email).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusConflict, models.Response{
			Code:    models.CODE_EMAIL_EXISTS,
			Message: "Email already exists",
			Data:    nil,
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to hash password",
			Data:    nil,
		})
		return
	}

	// Create admin
	admin := models.Admin{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role:     request.Role,
		Status:   models.ADMIN_STATUS_ACTIVE,
	}

	if err := h.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_USER_CREATE_FAILED,
			Message: "Failed to create admin",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, models.AdminCreatedResponse(admin))
}

// UpdateAdmin updates an existing admin
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	adminIDStr := c.Param("admin_id")
	adminID, err := strconv.ParseUint(adminIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_INVALID_USER_ID,
			Message: "Invalid admin ID",
			Data:    nil,
		})
		return
	}

	var request models.UpdateAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid request data",
			Data:    err.Error(),
		})
		return
	}

	// Find admin
	var admin models.Admin
	if err := h.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    models.CODE_USER_NOT_FOUND,
			Message: "Admin not found",
			Data:    nil,
		})
		return
	}

	// Update fields
	if request.Name != "" {
		admin.Name = request.Name
	}
	if request.Email != "" {
		// Check if new email already exists
		var existingAdmin models.Admin
		if err := h.DB.Where("email = ? AND id != ?", request.Email, adminID).First(&existingAdmin).Error; err == nil {
			c.JSON(http.StatusConflict, models.Response{
				Code:    models.CODE_EMAIL_EXISTS,
				Message: "Email already exists",
				Data:    nil,
			})
			return
		}
		admin.Email = request.Email
	}
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    models.CODE_INTERNAL_SERVER,
				Message: "Failed to hash password",
				Data:    nil,
			})
			return
		}
		admin.Password = string(hashedPassword)
	}
	if request.Role != "" {
		if !models.ValidateAdminRole(request.Role) {
			c.JSON(http.StatusBadRequest, models.Response{
				Code:    models.CODE_VALIDATION_FAILED,
				Message: "Invalid role",
				Data:    nil,
			})
			return
		}
		admin.Role = request.Role
	}
	if request.Status != nil {
		if !models.ValidateAdminStatus(*request.Status) {
			c.JSON(http.StatusBadRequest, models.Response{
				Code:    models.CODE_VALIDATION_FAILED,
				Message: "Invalid status",
				Data:    nil,
			})
			return
		}
		admin.Status = *request.Status
	}

	// Save changes
	if err := h.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_USER_UPDATE_FAILED,
			Message: "Failed to update admin",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.AdminUpdatedResponse(admin))
}

// DeleteAdmin deletes an admin
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	adminIDStr := c.Param("admin_id")
	adminID, err := strconv.ParseUint(adminIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_INVALID_USER_ID,
			Message: "Invalid admin ID",
			Data:    nil,
		})
		return
	}

	// Get current admin ID from context
	currentAdminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_UNAUTHORIZED,
			Message: "Admin authentication required",
			Data:    nil,
		})
		return
	}

	// Prevent admin from deleting themselves
	if uint(adminID) == currentAdminID.(uint) {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Cannot delete your own account",
			Data:    nil,
		})
		return
	}

	// Find admin
	var admin models.Admin
	if err := h.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    models.CODE_USER_NOT_FOUND,
			Message: "Admin not found",
			Data:    nil,
		})
		return
	}

	// Soft delete admin
	if err := h.DB.Delete(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_USER_DELETE_FAILED,
			Message: "Failed to soft delete admin",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.AdminSoftDeletedSuccessResponse(&admin))
}

// RestoreAdmin restores a soft deleted admin by ID
func (h *AdminHandler) RestoreAdmin(c *gin.Context) {
	adminID := c.Param("admin_id")
	id, err := strconv.ParseUint(adminID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_INVALID_REQUEST,
			Message: "Invalid admin ID",
			Data:    nil,
		})
		return
	}

	var admin models.Admin
	// Find soft deleted admin
	if err := h.DB.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", uint(id)).First(&admin).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    models.CODE_USER_NOT_FOUND,
				Message: "Deleted admin not found",
				Data:    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    models.CODE_USER_DELETE_FAILED,
				Message: "Database error",
				Data:    nil,
			})
		}
		return
	}

	// Restore the admin by setting deleted_at to NULL
	if err := h.DB.Unscoped().Model(&admin).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_USER_DELETE_FAILED,
			Message: "Failed to restore admin",
			Data:    nil,
		})
		return
	}

	// Refresh admin data
	h.DB.First(&admin, uint(id))

	c.JSON(http.StatusOK, models.AdminRestoredSuccessResponse(&admin))
}

// GetDeletedAdmins retrieves all soft deleted admins
func (h *AdminHandler) GetDeletedAdmins(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	var admins []models.Admin
	var total int64

	// Count total soft deleted admins
	h.DB.Unscoped().Where("deleted_at IS NOT NULL").Model(&models.Admin{}).Count(&total)

	// Get soft deleted admins
	if err := h.DB.Unscoped().Where("deleted_at IS NOT NULL").
		Offset(offset).Limit(perPage).Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_USER_LIST_FAILED,
			Message: "Failed to retrieve deleted admins",
			Data:    nil,
		})
		return
	}

	// Convert to response format
	adminResponses := make([]models.AdminResponse, len(admins))
	for i, admin := range admins {
		adminResponses[i] = admin.ToResponse()
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    models.CODE_SUCCESS,
		Message: "Deleted admins retrieved successfully",
		Data: gin.H{
			"admins":   adminResponses,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

// PermanentDeleteAdmin permanently deletes a soft deleted admin (unrecoverable)
func (h *AdminHandler) PermanentDeleteAdmin(c *gin.Context) {
	adminID := c.Param("admin_id")
	id, err := strconv.ParseUint(adminID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_INVALID_REQUEST,
			Message: "Invalid admin ID",
			Data:    nil,
		})
		return
	}

	var admin models.Admin
	// Find soft deleted admin first
	if err := h.DB.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", uint(id)).First(&admin).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    models.CODE_USER_NOT_FOUND,
				Message: "Deleted admin not found",
				Data:    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    models.CODE_USER_DELETE_FAILED,
				Message: "Database error",
				Data:    nil,
			})
		}
		return
	}

	// Get current admin info for self-deletion check
	currentAdminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_UNAUTHORIZED,
			Message: "Admin authentication required",
			Data:    nil,
		})
		return
	}

	// Prevent self permanent deletion
	if admin.ID == currentAdminID.(uint) {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_INVALID_REQUEST,
			Message: "Cannot permanently delete your own account",
			Data:    nil,
		})
		return
	}

	// Permanently delete the admin
	if err := h.DB.Unscoped().Delete(&admin, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_USER_DELETE_FAILED,
			Message: "Failed to permanently delete admin",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    models.CODE_SUCCESS,
		Message: "Admin permanently deleted successfully",
		Data: gin.H{
			"id":    admin.ID,
			"name":  admin.Name,
			"email": admin.Email,
			"role":  admin.Role,
		},
	})
}

// GetAdmins retrieves all admins with pagination
func (h *AdminHandler) GetAdmins(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Get search filter parameters
	search := c.Query("search")
	name := c.Query("name")
	email := c.Query("email")
	role := c.Query("role")
	status := c.Query("status")

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Build query with filters
	query := h.DB.Model(&models.Admin{})

	// Apply search filter (searches across name and email)
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply specific field filters
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		if status == "active" {
			query = query.Where("status = ?", true)
		} else if status == "inactive" {
			query = query.Where("status = ?", false)
		}
	}

	// Get total count with filters
	var total int64
	query.Count(&total)

	// Get admins with filters and pagination
	var admins []models.Admin
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_USER_LIST_FAILED,
			Message: "Failed to retrieve admins",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.AdminListRetrievedResponse(admins, int(total), page, perPage))
}

// GetAdminByID retrieves a specific admin by ID
func (h *AdminHandler) GetAdminByID(c *gin.Context) {
	adminIDStr := c.Param("admin_id")
	adminID, err := strconv.ParseUint(adminIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_INVALID_USER_ID,
			Message: "Invalid admin ID",
			Data:    nil,
		})
		return
	}

	// Find admin
	var admin models.Admin
	if err := h.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    models.CODE_USER_NOT_FOUND,
			Message: "Admin not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Admin retrieved successfully",
		Data:    admin.ToResponse(),
	})
}

// GetDashboard retrieves dashboard statistics
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	var dashboard models.DashboardStats

	// Get total users count
	if err := h.DB.Model(&models.User{}).Count(&dashboard.TotalUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.DashboardRetrieveFailedResponse())
		return
	}

	// Get total admins count
	if err := h.DB.Model(&models.Admin{}).Count(&dashboard.TotalAdmins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.DashboardRetrieveFailedResponse())
		return
	}

	// Get current time and calculate time ranges
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Calculate start of week (Monday)
	weekday := now.Weekday()
	daysFromMonday := int(weekday) - 1
	if weekday == time.Sunday {
		daysFromMonday = 6
	}
	startOfWeek := startOfToday.AddDate(0, 0, -daysFromMonday)

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	// Get total transactions for all periods
	// Today's transactions
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ? AND created_at < ?", startOfToday, startOfToday.AddDate(0, 0, 1)).
		Count(&dashboard.TotalTransactions.Today)

	// This week's transactions
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ? AND created_at < ?", startOfWeek, startOfWeek.AddDate(0, 0, 7)).
		Count(&dashboard.TotalTransactions.ThisWeek)

	// This month's transactions
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ? AND created_at < ?", startOfMonth, startOfMonth.AddDate(0, 1, 0)).
		Count(&dashboard.TotalTransactions.ThisMonth)

	// This year's transactions
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ? AND created_at < ?", startOfYear, startOfYear.AddDate(1, 0, 0)).
		Count(&dashboard.TotalTransactions.ThisYear)

	// All time transactions
	h.DB.Model(&models.Transaction{}).
		Count(&dashboard.TotalTransactions.AllTime)

	// Get total transaction amounts for all periods
	// Today's total transaction amount
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ? AND created_at < ?", startOfToday, startOfToday.AddDate(0, 0, 1)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TotalTransactions.TodayAmount)

	// This week's total transaction amount
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ? AND created_at < ?", startOfWeek, startOfWeek.AddDate(0, 0, 7)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TotalTransactions.ThisWeekAmount)

	// This month's total transaction amount
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ? AND created_at < ?", startOfMonth, startOfMonth.AddDate(0, 1, 0)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TotalTransactions.ThisMonthAmount)

	// This year's total transaction amount
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ? AND created_at < ?", startOfYear, startOfYear.AddDate(1, 0, 0)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TotalTransactions.ThisYearAmount)

	// All time total transaction amount
	h.DB.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TotalTransactions.AllTimeAmount)

	// Get topup transactions for all periods
	// Today's topup transactions
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "topup", startOfToday, startOfToday.AddDate(0, 0, 1)).
		Count(&dashboard.TopupTransactions.Today)

	// This week's topup transactions
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "topup", startOfWeek, startOfWeek.AddDate(0, 0, 7)).
		Count(&dashboard.TopupTransactions.ThisWeek)

	// This month's topup transactions
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "topup", startOfMonth, startOfMonth.AddDate(0, 1, 0)).
		Count(&dashboard.TopupTransactions.ThisMonth)

	// This year's topup transactions
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "topup", startOfYear, startOfYear.AddDate(1, 0, 0)).
		Count(&dashboard.TopupTransactions.ThisYear)

	// All time topup transactions
	h.DB.Model(&models.Transaction{}).
		Where("type = ?", "topup").
		Count(&dashboard.TopupTransactions.AllTime)

	// Get total topup amounts for all periods
	// Today's topup amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "topup", startOfToday, startOfToday.AddDate(0, 0, 1)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TopupTransactions.TodayAmount)

	// This week's topup amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "topup", startOfWeek, startOfWeek.AddDate(0, 0, 7)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TopupTransactions.ThisWeekAmount)

	// This month's topup amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "topup", startOfMonth, startOfMonth.AddDate(0, 1, 0)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TopupTransactions.ThisMonthAmount)

	// This year's topup amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "topup", startOfYear, startOfYear.AddDate(1, 0, 0)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TopupTransactions.ThisYearAmount)

	// All time topup amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ?", "topup").
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TopupTransactions.AllTimeAmount)

	// Get withdraw transactions for all periods
	// Today's withdraw transactions
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "withdraw", startOfToday, startOfToday.AddDate(0, 0, 1)).
		Count(&dashboard.WithdrawTransactions.Today)

	// This week's withdraw transactions
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "withdraw", startOfWeek, startOfWeek.AddDate(0, 0, 7)).
		Count(&dashboard.WithdrawTransactions.ThisWeek)

	// This month's withdraw transactions
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "withdraw", startOfMonth, startOfMonth.AddDate(0, 1, 0)).
		Count(&dashboard.WithdrawTransactions.ThisMonth)

	// This year's withdraw transactions
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "withdraw", startOfYear, startOfYear.AddDate(1, 0, 0)).
		Count(&dashboard.WithdrawTransactions.ThisYear)

	// All time withdraw transactions
	h.DB.Model(&models.Transaction{}).
		Where("type = ?", "withdraw").
		Count(&dashboard.WithdrawTransactions.AllTime)

	// Get total withdraw amounts for all periods
	// Today's withdraw amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "withdraw", startOfToday, startOfToday.AddDate(0, 0, 1)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.WithdrawTransactions.TodayAmount)

	// This week's withdraw amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "withdraw", startOfWeek, startOfWeek.AddDate(0, 0, 7)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.WithdrawTransactions.ThisWeekAmount)

	// This month's withdraw amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "withdraw", startOfMonth, startOfMonth.AddDate(0, 1, 0)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.WithdrawTransactions.ThisMonthAmount)

	// This year's withdraw amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "withdraw", startOfYear, startOfYear.AddDate(1, 0, 0)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.WithdrawTransactions.ThisYearAmount)

	// All time withdraw amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ?", "withdraw").
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.WithdrawTransactions.AllTimeAmount)

	// Get transfer transactions for all periods (both transfer_out and transfer_in)
	// Today's transfer transactions
	h.DB.Model(&models.Transaction{}).
		Where("(type = ? OR type = ?) AND created_at >= ? AND created_at < ?", "transfer_out", "transfer_in", startOfToday, startOfToday.AddDate(0, 0, 1)).
		Count(&dashboard.TransferTransactions.Today)

	// This week's transfer transactions
	h.DB.Model(&models.Transaction{}).
		Where("(type = ? OR type = ?) AND created_at >= ? AND created_at < ?", "transfer_out", "transfer_in", startOfWeek, startOfWeek.AddDate(0, 0, 7)).
		Count(&dashboard.TransferTransactions.ThisWeek)

	// This month's transfer transactions
	h.DB.Model(&models.Transaction{}).
		Where("(type = ? OR type = ?) AND created_at >= ? AND created_at < ?", "transfer_out", "transfer_in", startOfMonth, startOfMonth.AddDate(0, 1, 0)).
		Count(&dashboard.TransferTransactions.ThisMonth)

	// This year's transfer transactions
	h.DB.Model(&models.Transaction{}).
		Where("(type = ? OR type = ?) AND created_at >= ? AND created_at < ?", "transfer_out", "transfer_in", startOfYear, startOfYear.AddDate(1, 0, 0)).
		Count(&dashboard.TransferTransactions.ThisYear)

	// All time transfer transactions
	h.DB.Model(&models.Transaction{}).
		Where("(type = ? OR type = ?)", "transfer_out", "transfer_in").
		Count(&dashboard.TransferTransactions.AllTime)

	// Get total transfer amounts for all periods (only transfer_out to avoid double counting)
	// Today's transfer amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "transfer_out", startOfToday, startOfToday.AddDate(0, 0, 1)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TransferTransactions.TodayAmount)

	// This week's transfer amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "transfer_out", startOfWeek, startOfWeek.AddDate(0, 0, 7)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TransferTransactions.ThisWeekAmount)

	// This month's transfer amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "transfer_out", startOfMonth, startOfMonth.AddDate(0, 1, 0)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TransferTransactions.ThisMonthAmount)

	// This year's transfer amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "transfer_out", startOfYear, startOfYear.AddDate(1, 0, 0)).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TransferTransactions.ThisYearAmount)

	// All time transfer amount
	h.DB.Model(&models.Transaction{}).
		Where("type = ?", "transfer_out").
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TransferTransactions.AllTimeAmount)

	// Calculate last 7 days and last 30 days for all transaction types
	last7DaysStart := startOfToday.AddDate(0, 0, -6)   // Last 7 days including today
	last30DaysStart := startOfToday.AddDate(0, 0, -29) // Last 30 days including today

	// Total transactions - Last 7 days
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ?", last7DaysStart).
		Count(&dashboard.TotalTransactions.Last7Days)
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ?", last7DaysStart).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TotalTransactions.Last7DaysAmount)

	// Total transactions - Last 30 days
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ?", last30DaysStart).
		Count(&dashboard.TotalTransactions.Last30Days)
	h.DB.Model(&models.Transaction{}).
		Where("created_at >= ?", last30DaysStart).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TotalTransactions.Last30DaysAmount)

	// Topup transactions - Last 7 days
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ?", "topup", last7DaysStart).
		Count(&dashboard.TopupTransactions.Last7Days)
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ?", "topup", last7DaysStart).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TopupTransactions.Last7DaysAmount)

	// Topup transactions - Last 30 days
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ?", "topup", last30DaysStart).
		Count(&dashboard.TopupTransactions.Last30Days)
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ?", "topup", last30DaysStart).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TopupTransactions.Last30DaysAmount)

	// Withdraw transactions - Last 7 days
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ?", "withdraw", last7DaysStart).
		Count(&dashboard.WithdrawTransactions.Last7Days)
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ?", "withdraw", last7DaysStart).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.WithdrawTransactions.Last7DaysAmount)

	// Withdraw transactions - Last 30 days
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ?", "withdraw", last30DaysStart).
		Count(&dashboard.WithdrawTransactions.Last30Days)
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ?", "withdraw", last30DaysStart).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.WithdrawTransactions.Last30DaysAmount)

	// Transfer transactions - Last 7 days (count both in and out, amount only out to avoid double counting)
	h.DB.Model(&models.Transaction{}).
		Where("(type = ? OR type = ?) AND created_at >= ?", "transfer_out", "transfer_in", last7DaysStart).
		Count(&dashboard.TransferTransactions.Last7Days)
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ?", "transfer_out", last7DaysStart).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TransferTransactions.Last7DaysAmount)

	// Transfer transactions - Last 30 days (count both in and out, amount only out to avoid double counting)
	h.DB.Model(&models.Transaction{}).
		Where("(type = ? OR type = ?) AND created_at >= ?", "transfer_out", "transfer_in", last30DaysStart).
		Count(&dashboard.TransferTransactions.Last30Days)
	h.DB.Model(&models.Transaction{}).
		Where("type = ? AND created_at >= ?", "transfer_out", last30DaysStart).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&dashboard.TransferTransactions.Last30DaysAmount)

	// Get performance data for charts
	// Populate Last7Days chart data (daily)
	for i := 6; i >= 0; i-- {
		dayStart := startOfToday.AddDate(0, 0, -i)
		dayEnd := dayStart.AddDate(0, 0, 1)
		dayLabel := dayStart.Format("Jan 02")

		var dayCount int64
		var dayAmount int64

		h.DB.Model(&models.Transaction{}).
			Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
			Count(&dayCount)

		h.DB.Model(&models.Transaction{}).
			Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
			Select("COALESCE(SUM(amount), 0)").Row().Scan(&dayAmount)

		dashboard.Performance.Last7Days = append(dashboard.Performance.Last7Days, models.PerformanceDataPoint{
			Period: dayLabel,
			Count:  dayCount,
			Amount: dayAmount,
		})
	}

	// Populate Last30Days chart data (daily)
	for i := 29; i >= 0; i-- {
		dayStart := startOfToday.AddDate(0, 0, -i)
		dayEnd := dayStart.AddDate(0, 0, 1)
		dayLabel := dayStart.Format("Jan 02")

		var dayCount int64
		var dayAmount int64

		h.DB.Model(&models.Transaction{}).
			Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
			Count(&dayCount)

		h.DB.Model(&models.Transaction{}).
			Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
			Select("COALESCE(SUM(amount), 0)").Row().Scan(&dayAmount)

		dashboard.Performance.Last30Days = append(dashboard.Performance.Last30Days, models.PerformanceDataPoint{
			Period: dayLabel,
			Count:  dayCount,
			Amount: dayAmount,
		})
	}

	c.JSON(http.StatusOK, models.DashboardSuccessResponse(dashboard))
}

// GetAdminTransactionByID - Get transaction detail by ID for admin
func (h *AdminHandler) GetAdminTransactionByID(c *gin.Context) {
	// Get transaction ID from URL parameter
	txnIDStr := c.Param("id")
	txnID, err := strconv.ParseUint(txnIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid transaction ID",
			Data:    nil,
		})
		return
	}

	// Find transaction by ID with all related data preloaded
	var transaction models.Transaction
	err = h.DB.Preload("User").
		Preload("OriginalTxn").
		Preload("OriginalTxn.User").
		Preload("ReversedTxn").
		Preload("ReversedTxn.User").
		First(&transaction, uint(txnID)).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    models.CODE_NOT_FOUND,
				Message: "Transaction not found",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to retrieve transaction",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    models.CODE_SUCCESS,
		Message: "Transaction retrieved successfully",
		Data:    transaction,
	})
}

// AdminTopupUserBalance - Admin can topup user balance
func (h *AdminHandler) AdminTopupUserBalance(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid user ID",
			Data:    nil,
		})
		return
	}

	// Parse request body
	var request models.TopupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid request data",
			Data:    err.Error(),
		})
		return
	}

	// Validate amount
	if request.Amount <= 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Amount must be greater than zero",
			Data:    nil,
		})
		return
	}

	// Check if user exists and is active
	var user models.User
	if err := h.DB.First(&user, uint(userID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    models.CODE_NOT_FOUND,
				Message: "User not found",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to find user",
			Data:    err.Error(),
		})
		return
	}

	// Check if user is active
	if user.Status != models.USER_STATUS_ACTIVE {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "User account is not active",
			Data:    nil,
		})
		return
	}

	// Get admin information from context
	adminInterface, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_UNAUTHORIZED,
			Message: "Admin context not found",
			Data:    nil,
		})
		return
	}

	admin, ok := adminInterface.(models.Admin)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_UNAUTHORIZED,
			Message: "Invalid admin context",
			Data:    nil,
		})
		return
	}

	// Start database transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lock user for update to prevent race conditions
	var lockedUser models.User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&lockedUser, uint(userID)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to lock user record",
			Data:    err.Error(),
		})
		return
	}

	// Calculate new balance
	balanceBefore := lockedUser.Balance
	balanceAfter := balanceBefore + request.Amount

	// Create transaction record
	description := request.Description
	if description == "" {
		description = "Admin top-up balance"
	}

	transaction := models.Transaction{
		UserID:        uint(userID),
		Type:          "topup",
		Amount:        request.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   description,
		Status:        "completed",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to create transaction record",
			Data:    err.Error(),
		})
		return
	}

	// Update user balance
	if err := tx.Model(&lockedUser).Update("balance", balanceAfter).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to update user balance",
			Data:    err.Error(),
		})
		return
	}

	// Create audit log for admin action
	auditDetails := map[string]interface{}{
		"amount":           request.Amount,
		"balance_before":   balanceBefore,
		"balance_after":    balanceAfter,
		"description":      description,
		"transaction_id":   transaction.ID,
		"target_user_id":   userID,
		"target_user_name": user.Name,
	}

	auditDetailsJSON, _ := json.Marshal(auditDetails)
	auditDetailsRaw := json.RawMessage(auditDetailsJSON)

	auditLog := models.AuditLog{
		EntityType: "user_balance",
		EntityID:   uint(userID),
		Action:     "ADMIN_TOPUP",
		AdminID:    &admin.ID,
		IPAddress:  c.ClientIP(),
		NewValues:  &auditDetailsRaw,
	}

	if err := tx.Create(&auditLog).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to create audit log",
			Data:    err.Error(),
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to commit transaction",
			Data:    err.Error(),
		})
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"transaction_id": transaction.ID,
		"user_id":        userID,
		"user_name":      user.Name,
		"amount":         request.Amount,
		"balance_before": balanceBefore,
		"balance_after":  balanceAfter,
		"description":    description,
		"admin_id":       admin.ID,
		"admin_name":     admin.Name,
		"created_at":     transaction.CreatedAt,
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    models.CODE_SUCCESS,
		Message: "User balance topped up successfully",
		Data:    responseData,
	})
}

// AdminAdjustUserBalance - Admin can adjust user balance with credit/debit
func (h *AdminHandler) AdminAdjustUserBalance(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid user ID",
			Data:    nil,
		})
		return
	}

	// Parse request body
	var request models.BalanceAdjustmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid request data",
			Data:    err.Error(),
		})
		return
	}

	// Validate amount (cannot be zero)
	if request.Amount == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Amount cannot be zero",
			Data:    nil,
		})
		return
	}

	// Check if user exists and is active
	var user models.User
	if err := h.DB.First(&user, uint(userID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    models.CODE_NOT_FOUND,
				Message: "User not found",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to find user",
			Data:    err.Error(),
		})
		return
	}

	// Check if user is active
	if user.Status != models.USER_STATUS_ACTIVE {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "User account is not active",
			Data:    nil,
		})
		return
	}

	// Get admin information from context
	adminInterface, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_UNAUTHORIZED,
			Message: "Admin context not found",
			Data:    nil,
		})
		return
	}

	admin, ok := adminInterface.(models.Admin)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_UNAUTHORIZED,
			Message: "Invalid admin context",
			Data:    nil,
		})
		return
	}

	// Start database transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lock user for update to prevent race conditions
	var lockedUser models.User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&lockedUser, uint(userID)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to lock user record",
			Data:    err.Error(),
		})
		return
	}

	// Calculate new balance
	balanceBefore := lockedUser.Balance
	balanceAfter := balanceBefore + request.Amount

	// Validate that balance doesn't go negative
	if balanceAfter < 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Insufficient balance for debit adjustment",
			Data: map[string]interface{}{
				"current_balance":   balanceBefore,
				"adjustment_amount": request.Amount,
				"resulting_balance": balanceAfter,
			},
		})
		return
	}

	// Determine transaction type based on amount
	var transactionType string
	if request.Amount > 0 {
		transactionType = "adjustment_credit"
	} else {
		transactionType = "adjustment_debit"
	}

	// Create transaction record
	description := request.Description
	if description == "" {
		if request.Amount > 0 {
			description = "Admin balance credit adjustment"
		} else {
			description = "Admin balance debit adjustment"
		}
	}

	transaction := models.Transaction{
		UserID:        uint(userID),
		Type:          transactionType,
		Amount:        request.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   description,
		Status:        "completed",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to create transaction record",
			Data:    err.Error(),
		})
		return
	}

	// Update user balance
	if err := tx.Model(&lockedUser).Update("balance", balanceAfter).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to update user balance",
			Data:    err.Error(),
		})
		return
	}

	// Create audit log for admin action
	auditDetails := map[string]interface{}{
		"amount":           request.Amount,
		"balance_before":   balanceBefore,
		"balance_after":    balanceAfter,
		"reason":           request.Reason,
		"adjustment_type":  request.Type,
		"description":      description,
		"transaction_id":   transaction.ID,
		"target_user_id":   userID,
		"target_user_name": user.Name,
	}

	auditDetailsJSON, _ := json.Marshal(auditDetails)
	auditDetailsRaw := json.RawMessage(auditDetailsJSON)

	auditLog := models.AuditLog{
		EntityType: "user_balance",
		EntityID:   uint(userID),
		Action:     "ADMIN_BALANCE_ADJUSTMENT",
		AdminID:    &admin.ID,
		IPAddress:  c.ClientIP(),
		NewValues:  &auditDetailsRaw,
	}

	if err := tx.Create(&auditLog).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to create audit log",
			Data:    err.Error(),
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to commit transaction",
			Data:    err.Error(),
		})
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"transaction_id":    transaction.ID,
		"user_id":           userID,
		"user_name":         user.Name,
		"adjustment_amount": request.Amount,
		"adjustment_type":   request.Type,
		"balance_before":    balanceBefore,
		"balance_after":     balanceAfter,
		"reason":            request.Reason,
		"description":       description,
		"admin_id":          admin.ID,
		"admin_name":        admin.Name,
		"created_at":        transaction.CreatedAt,
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    models.CODE_SUCCESS,
		Message: "User balance adjusted successfully",
		Data:    responseData,
	})
}

// AdminSetUserBalance - Admin can set exact user balance
func (h *AdminHandler) AdminSetUserBalance(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid user ID",
			Data:    nil,
		})
		return
	}

	// Parse request body
	var request models.BalanceSetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid request data",
			Data:    err.Error(),
		})
		return
	}

	// Check if user exists and is active
	var user models.User
	if err := h.DB.First(&user, uint(userID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    models.CODE_NOT_FOUND,
				Message: "User not found",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to find user",
			Data:    err.Error(),
		})
		return
	}

	// Check if user is active
	if user.Status != models.USER_STATUS_ACTIVE {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "User account is not active",
			Data:    nil,
		})
		return
	}

	// Get admin information from context
	adminInterface, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_UNAUTHORIZED,
			Message: "Admin context not found",
			Data:    nil,
		})
		return
	}

	admin, ok := adminInterface.(models.Admin)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    models.CODE_UNAUTHORIZED,
			Message: "Invalid admin context",
			Data:    nil,
		})
		return
	}

	// Start database transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lock user for update to prevent race conditions
	var lockedUser models.User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&lockedUser, uint(userID)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to lock user record",
			Data:    err.Error(),
		})
		return
	}

	// Calculate adjustment amount
	balanceBefore := lockedUser.Balance
	balanceAfter := request.Balance
	adjustmentAmount := balanceAfter - balanceBefore

	// Skip if balance is already at the desired amount
	if adjustmentAmount == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "User balance is already at the specified amount",
			Data: map[string]interface{}{
				"current_balance":   balanceBefore,
				"requested_balance": request.Balance,
			},
		})
		return
	}

	// Determine transaction type based on adjustment
	var transactionType string
	if adjustmentAmount > 0 {
		transactionType = "balance_set_credit"
	} else {
		transactionType = "balance_set_debit"
	}

	// Create transaction record
	description := request.Description
	if description == "" {
		description = "Admin balance set operation"
	}

	transaction := models.Transaction{
		UserID:        uint(userID),
		Type:          transactionType,
		Amount:        adjustmentAmount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   description,
		Status:        "completed",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to create transaction record",
			Data:    err.Error(),
		})
		return
	}

	// Update user balance
	if err := tx.Model(&lockedUser).Update("balance", balanceAfter).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to update user balance",
			Data:    err.Error(),
		})
		return
	}

	// Create audit log for admin action
	auditDetails := map[string]interface{}{
		"balance_before":    balanceBefore,
		"balance_after":     balanceAfter,
		"adjustment_amount": adjustmentAmount,
		"reason":            request.Reason,
		"description":       description,
		"transaction_id":    transaction.ID,
		"target_user_id":    userID,
		"target_user_name":  user.Name,
	}

	auditDetailsJSON, _ := json.Marshal(auditDetails)
	auditDetailsRaw := json.RawMessage(auditDetailsJSON)

	auditLog := models.AuditLog{
		EntityType: "user_balance",
		EntityID:   uint(userID),
		Action:     "ADMIN_BALANCE_SET",
		AdminID:    &admin.ID,
		IPAddress:  c.ClientIP(),
		NewValues:  &auditDetailsRaw,
	}

	if err := tx.Create(&auditLog).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to create audit log",
			Data:    err.Error(),
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to commit transaction",
			Data:    err.Error(),
		})
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"transaction_id":    transaction.ID,
		"user_id":           userID,
		"user_name":         user.Name,
		"balance_before":    balanceBefore,
		"balance_after":     balanceAfter,
		"adjustment_amount": adjustmentAmount,
		"reason":            request.Reason,
		"description":       description,
		"admin_id":          admin.ID,
		"admin_name":        admin.Name,
		"created_at":        transaction.CreatedAt,
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    models.CODE_SUCCESS,
		Message: "User balance set successfully",
		Data:    responseData,
	})
}

// AdminGetUserBalanceHistory - Get user balance change history
func (h *AdminHandler) AdminGetUserBalanceHistory(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    models.CODE_VALIDATION_FAILED,
			Message: "Invalid user ID",
			Data:    nil,
		})
		return
	}

	// Check if user exists
	var user models.User
	if err := h.DB.First(&user, uint(userID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    models.CODE_NOT_FOUND,
				Message: "User not found",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to find user",
			Data:    err.Error(),
		})
		return
	}

	// Parse query parameters for pagination and filtering
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	transactionType := c.Query("type") // Filter by transaction type

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Build query for balance-affecting transactions
	query := h.DB.Where("user_id = ?", userID).
		Where("type IN (?)", []string{
			"topup", "withdraw", "transfer_in", "transfer_out",
			"adjustment_credit", "adjustment_debit",
			"balance_set_credit", "balance_set_debit",
			"reversal_credit", "reversal_debit",
		})

	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}

	// Get total count
	var total int64
	if err := query.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to count transactions",
			Data:    err.Error(),
		})
		return
	}

	// Get transactions with pagination
	var transactions []models.Transaction
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_INTERNAL_SERVER,
			Message: "Failed to retrieve transactions",
			Data:    err.Error(),
		})
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":              user.ID,
			"name":            user.Name,
			"phone":           user.Phone,
			"current_balance": user.Balance,
			"status":          user.Status,
		},
		"balance_history": transactions,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    models.CODE_SUCCESS,
		Message: "User balance history retrieved successfully",
		Data:    responseData,
	})
}
