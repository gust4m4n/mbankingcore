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
