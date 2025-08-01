package handlers

import (
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
	adminIDStr := c.Param("id")
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
	adminIDStr := c.Param("id")
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

	// Delete admin
	if err := h.DB.Delete(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    models.CODE_USER_DELETE_FAILED,
			Message: "Failed to delete admin",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.AdminDeletedResponse())
}

// GetAdmins retrieves all admins with pagination
func (h *AdminHandler) GetAdmins(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Get total count
	var total int64
	h.DB.Model(&models.Admin{}).Count(&total)

	// Get admins
	var admins []models.Admin
	if err := h.DB.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&admins).Error; err != nil {
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
	adminIDStr := c.Param("id")
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
