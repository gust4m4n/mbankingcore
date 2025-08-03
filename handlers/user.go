package handlers

import (
	"strconv"
	"time"

	"mbankingcore/config"
	"mbankingcore/models"

	"github.com/gin-gonic/gin"
)

// ListUsers retrieves all users with pagination
func ListUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Get search filter parameters
	search := c.Query("search")
	name := c.Query("name")
	phone := c.Query("phone")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Build query with filters
	query := config.DB.Model(&models.User{})

	// Apply search filter (searches across name and phone)
	if search != "" {
		query = query.Where("name ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply specific field filters
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if phone != "" {
		query = query.Where("phone ILIKE ?", "%"+phone+"%")
	}
	if status != "" {
		if status == "active" {
			query = query.Where("status = ?", true)
		} else if status == "inactive" {
			query = query.Where("status = ?", false)
		}
	}

	var users []models.User
	var total int64

	// Count total users with filters
	if err := query.Count(&total).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_RETRIEVE_FAILED,
			"message": "Failed to retrieve users",
			"error":   err.Error(),
		})
		return
	}

	// Get users with filters and pagination
	if err := query.Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&users).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_RETRIEVE_FAILED,
			"message": "Failed to retrieve users",
			"error":   err.Error(),
		})
		return
	}

	// Use the new response helper function
	response := models.UsersListRetrievedResponse(users, int(total), page, perPage)
	c.JSON(response.Code, response)
}

// GetUserByID retrieves a specific user by ID
func GetUserByID(c *gin.Context) {
	// Get ID from URL parameter
	idParam := c.Param("user_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_USER_ID,
			"message": "Invalid user ID",
		})
		return
	}

	var user models.User

	// Find user by ID
	result := config.DB.First(&user, uint(id))
	if result.Error != nil {
		c.JSON(404, gin.H{
			"code":    models.CODE_USER_NOT_FOUND,
			"message": "User not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    models.CODE_SUCCESS,
		"message": "User retrieved successfully",
		"data":    user.ToResponse(),
	})
}

// DeleteUser soft deletes a user by ID
func DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid user ID",
		})
		return
	}

	var user models.User
	result := config.DB.First(&user, uint(id))
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(404, gin.H{
				"code":    models.CODE_USER_NOT_FOUND,
				"message": "User not found",
			})
		} else {
			c.JSON(500, gin.H{
				"code":    models.CODE_USER_DELETE_FAILED,
				"message": "Database error",
				"error":   result.Error.Error(),
			})
		}
		return
	}

	// Soft delete the user
	result = config.DB.Delete(&user, uint(id))
	if result.Error != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_DELETE_FAILED,
			"message": "Failed to soft delete user",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(200, models.UserSoftDeletedSuccessResponse(&user))
}

// RestoreUser restores a soft deleted user by ID
func RestoreUser(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid user ID",
		})
		return
	}

	var user models.User
	// Find soft deleted user
	result := config.DB.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", uint(id)).First(&user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(404, gin.H{
				"code":    models.CODE_USER_NOT_FOUND,
				"message": "Deleted user not found",
			})
		} else {
			c.JSON(500, gin.H{
				"code":    models.CODE_USER_DELETE_FAILED,
				"message": "Database error",
				"error":   result.Error.Error(),
			})
		}
		return
	}

	// Restore the user by setting deleted_at to NULL
	result = config.DB.Unscoped().Model(&user).Update("deleted_at", nil)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_DELETE_FAILED,
			"message": "Failed to restore user",
			"error":   result.Error.Error(),
		})
		return
	}

	// Refresh user data
	config.DB.First(&user, uint(id))

	c.JSON(200, models.UserRestoredSuccessResponse(&user))
}

// GetDeletedUsers retrieves all soft deleted users
func GetDeletedUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	var users []models.User
	var total int64

	// Count total soft deleted users
	config.DB.Unscoped().Where("deleted_at IS NOT NULL").Model(&models.User{}).Count(&total)

	// Get soft deleted users
	result := config.DB.Unscoped().Where("deleted_at IS NOT NULL").
		Offset(offset).Limit(perPage).Find(&users)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_LIST_FAILED,
			"message": "Failed to retrieve deleted users",
			"error":   result.Error.Error(),
		})
		return
	}

	// Convert to response format
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = models.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Phone:      user.Phone,
			MotherName: user.MotherName,
			Balance:    user.Balance,
			Status:     user.Status,
			Avatar:     user.Avatar,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		}
	}

	c.JSON(200, gin.H{
		"code":    models.CODE_SUCCESS,
		"message": "Deleted users retrieved successfully",
		"data": gin.H{
			"users":    userResponses,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

// PermanentDeleteUser permanently deletes a soft deleted user (unrecoverable)
func PermanentDeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid user ID",
		})
		return
	}

	var user models.User
	// Find soft deleted user first
	result := config.DB.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", uint(id)).First(&user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(404, gin.H{
				"code":    models.CODE_USER_NOT_FOUND,
				"message": "Deleted user not found",
			})
		} else {
			c.JSON(500, gin.H{
				"code":    models.CODE_USER_DELETE_FAILED,
				"message": "Database error",
				"error":   result.Error.Error(),
			})
		}
		return
	}

	// Permanently delete the user
	result = config.DB.Unscoped().Delete(&user, uint(id))
	if result.Error != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_DELETE_FAILED,
			"message": "Failed to permanently delete user",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    models.CODE_SUCCESS,
		"message": "User permanently deleted successfully",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"phone": user.Phone,
		},
	})
}

// UpdateUserStatus updates user status (direct admin action)
func UpdateUserStatus(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid user ID",
		})
		return
	}

	var req models.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Validate status
	if !models.ValidateStatus(req.Status) {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid status value",
		})
		return
	}

	// Find user
	var user models.User
	result := config.DB.First(&user, uint(id))
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(404, gin.H{
				"code":    models.CODE_USER_NOT_FOUND,
				"message": "User not found",
			})
		} else {
			c.JSON(500, gin.H{
				"code":    models.CODE_USER_UPDATE_FAILED,
				"message": "Database error",
				"error":   result.Error.Error(),
			})
		}
		return
	}

	// Get admin info for audit
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(401, gin.H{
			"code":    models.CODE_UNAUTHORIZED,
			"message": "Admin not authenticated",
		})
		return
	}

	var admin models.Admin
	config.DB.First(&admin, adminID)

	// Store previous status
	previousStatus := user.Status

	// Update status
	user.Status = req.Status
	result = config.DB.Save(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_UPDATE_FAILED,
			"message": "Failed to update user status",
			"error":   result.Error.Error(),
		})
		return
	}

	// Create audit log
	auditLog := models.AuditLog{
		AdminID:       &admin.ID,
		EntityType:    "user",
		EntityID:      user.ID,
		Action:        "UPDATE_STATUS",
		IPAddress:     c.ClientIP(),
		UserAgent:     c.GetHeader("User-Agent"),
		APIEndpoint:   c.Request.URL.Path,
		RequestMethod: c.Request.Method,
		StatusCode:    200,
	}
	config.DB.Create(&auditLog)

	c.JSON(200, gin.H{
		"code":    models.CODE_SUCCESS,
		"message": "User status updated successfully",
		"data": models.UpdateUserStatusResponse{
			ID:                   user.ID,
			Name:                 user.Name,
			Phone:                user.Phone,
			PreviousStatus:       previousStatus,
			NewStatus:            req.Status,
			PreviousStatusString: models.GetUserStatusString(previousStatus),
			NewStatusString:      models.GetUserStatusString(req.Status),
			Reason:               req.Reason,
			UpdatedBy:            admin.Name,
			UpdatedAt:            user.UpdatedAt,
		},
	})
}

// CreatePendingUserStatusChange creates a pending user status change request (maker-checker)
func CreatePendingUserStatusChange(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid user ID",
		})
		return
	}

	var req models.CreatePendingUserStatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Validate status
	if !models.ValidateStatus(req.Status) {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid status value",
		})
		return
	}

	// Validate priority if provided
	if req.Priority != "" && req.Priority != "low" && req.Priority != "normal" && req.Priority != "high" && req.Priority != "critical" {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid priority value. Must be: low, normal, high, or critical",
		})
		return
	}

	// Set default priority
	if req.Priority == "" {
		req.Priority = "normal"
	}

	// Find user
	var user models.User
	result := config.DB.First(&user, uint(id))
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(404, gin.H{
				"code":    models.CODE_USER_NOT_FOUND,
				"message": "User not found",
			})
		} else {
			c.JSON(500, gin.H{
				"code":    models.CODE_USER_UPDATE_FAILED,
				"message": "Database error",
				"error":   result.Error.Error(),
			})
		}
		return
	}

	// Check if status is already the same
	if user.Status == req.Status {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "User already has the requested status",
		})
		return
	}

	// Get admin info
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(401, gin.H{
			"code":    models.CODE_UNAUTHORIZED,
			"message": "Admin not authenticated",
		})
		return
	}

	// Check if there's already a pending request for this user
	var existingPending models.PendingUserStatusChange
	existingResult := config.DB.Where("user_id = ? AND status = 'pending'", user.ID).First(&existingPending)
	if existingResult.Error == nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "There is already a pending status change request for this user",
		})
		return
	}

	// Calculate expiration time (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Create pending request
	pending := models.PendingUserStatusChange{
		UserID:          user.ID,
		MakerAdminID:    adminID.(uint),
		CurrentStatus:   user.Status,
		RequestedStatus: req.Status,
		Reason:          req.Reason,
		Priority:        req.Priority,
		ExpiresAt:       &expiresAt,
	}

	result = config.DB.Create(&pending)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_UPDATE_FAILED,
			"message": "Failed to create pending status change request",
			"error":   result.Error.Error(),
		})
		return
	}

	// Load relationships for response
	config.DB.Preload("User").Preload("MakerAdmin").First(&pending, pending.ID)

	// Create audit log
	auditLog := models.AuditLog{
		AdminID:       &pending.MakerAdminID,
		EntityType:    "user_status_change",
		EntityID:      pending.ID,
		Action:        "CREATE_PENDING",
		IPAddress:     c.ClientIP(),
		UserAgent:     c.GetHeader("User-Agent"),
		APIEndpoint:   c.Request.URL.Path,
		RequestMethod: c.Request.Method,
		StatusCode:    201,
	}
	config.DB.Create(&auditLog)

	c.JSON(201, gin.H{
		"code":    models.CODE_SUCCESS,
		"message": "Pending user status change request created successfully",
		"data": models.PendingUserStatusChangeResponse{
			ID:                    pending.ID,
			UserID:                pending.UserID,
			UserName:              pending.User.Name,
			UserPhone:             pending.User.Phone,
			MakerAdminID:          pending.MakerAdminID,
			MakerAdminName:        pending.MakerAdmin.Name,
			CurrentStatus:         pending.CurrentStatus,
			CurrentStatusString:   models.GetUserStatusString(pending.CurrentStatus),
			RequestedStatus:       pending.RequestedStatus,
			RequestedStatusString: models.GetUserStatusString(pending.RequestedStatus),
			Reason:                pending.Reason,
			Status:                pending.Status,
			Priority:              pending.Priority,
			ExpiresAt:             pending.ExpiresAt,
			CreatedAt:             pending.CreatedAt,
			UpdatedAt:             pending.UpdatedAt,
		},
	})
}

// GetPendingUserStatusChanges retrieves all pending user status change requests
func GetPendingUserStatusChanges(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Query with filters
	query := config.DB.Model(&models.PendingUserStatusChange{}).Where("status = 'pending'")

	// Get priority filter
	priority := c.Query("priority")
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Get records with relationships
	var pending []models.PendingUserStatusChange
	result := query.Preload("User").Preload("MakerAdmin").
		Order("priority DESC, created_at ASC").
		Limit(perPage).Offset(offset).Find(&pending)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_LIST_FAILED,
			"message": "Failed to retrieve pending status changes",
			"error":   result.Error.Error(),
		})
		return
	}

	// Convert to response format
	var responses []models.PendingUserStatusChangeResponse
	for _, p := range pending {
		response := models.PendingUserStatusChangeResponse{
			ID:                    p.ID,
			UserID:                p.UserID,
			UserName:              p.User.Name,
			UserPhone:             p.User.Phone,
			MakerAdminID:          p.MakerAdminID,
			MakerAdminName:        p.MakerAdmin.Name,
			CurrentStatus:         p.CurrentStatus,
			CurrentStatusString:   models.GetUserStatusString(p.CurrentStatus),
			RequestedStatus:       p.RequestedStatus,
			RequestedStatusString: models.GetUserStatusString(p.RequestedStatus),
			Reason:                p.Reason,
			Status:                p.Status,
			Priority:              p.Priority,
			ExpiresAt:             p.ExpiresAt,
			CreatedAt:             p.CreatedAt,
			UpdatedAt:             p.UpdatedAt,
		}
		responses = append(responses, response)
	}

	totalPages := (int(total) + perPage - 1) / perPage

	c.JSON(200, gin.H{
		"code":    models.CODE_SUCCESS,
		"message": "Pending user status changes retrieved successfully",
		"data": models.PendingUserStatusChangeListResponse{
			PendingChanges: responses,
			Pagination: models.Pagination{
				Page:       page,
				Limit:      perPage,
				Total:      int(total),
				TotalPages: totalPages,
			},
		},
	})
}

// ReviewPendingUserStatusChange approves or rejects a pending status change request
func ReviewPendingUserStatusChange(c *gin.Context) {
	pendingID := c.Param("pending_id")
	id, err := strconv.Atoi(pendingID)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid pending request ID",
		})
		return
	}

	var req models.ReviewPendingUserStatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Validate action
	if req.Action != "approve" && req.Action != "reject" {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid action. Must be 'approve' or 'reject'",
		})
		return
	}

	// Get checker admin info
	checkerAdminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(401, gin.H{
			"code":    models.CODE_UNAUTHORIZED,
			"message": "Admin not authenticated",
		})
		return
	}

	// Find pending request
	var pending models.PendingUserStatusChange
	result := config.DB.Preload("User").Preload("MakerAdmin").First(&pending, uint(id))
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(404, gin.H{
				"code":    models.CODE_USER_NOT_FOUND,
				"message": "Pending request not found",
			})
		} else {
			c.JSON(500, gin.H{
				"code":    models.CODE_USER_UPDATE_FAILED,
				"message": "Database error",
				"error":   result.Error.Error(),
			})
		}
		return
	}

	// Check if already processed
	if pending.Status != "pending" {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Request has already been processed",
		})
		return
	}

	// Check if expired
	if pending.ExpiresAt != nil && time.Now().After(*pending.ExpiresAt) {
		pending.Status = "expired"
		config.DB.Save(&pending)
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Request has expired",
		})
		return
	}

	// Check if same admin (can't check own request)
	if pending.MakerAdminID == checkerAdminID.(uint) {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Cannot approve/reject your own request",
		})
		return
	}

	now := time.Now()
	checkerID := checkerAdminID.(uint)
	pending.CheckerAdminID = &checkerID
	var message string

	if req.Action == "approve" {
		// Update user status
		var user models.User
		userResult := config.DB.First(&user, pending.UserID)
		if userResult.Error != nil {
			c.JSON(500, gin.H{
				"code":    models.CODE_USER_UPDATE_FAILED,
				"message": "Failed to find user",
				"error":   userResult.Error.Error(),
			})
			return
		}

		user.Status = pending.RequestedStatus
		userResult = config.DB.Save(&user)
		if userResult.Error != nil {
			c.JSON(500, gin.H{
				"code":    models.CODE_USER_UPDATE_FAILED,
				"message": "Failed to update user status",
				"error":   userResult.Error.Error(),
			})
			return
		}

		// Update pending request
		pending.Status = "approved"
		pending.ApprovedAt = &now
		pending.ApprovalComments = req.Comments

		// Create audit log for status change
		auditLog := models.AuditLog{
			AdminID:       pending.CheckerAdminID,
			EntityType:    "user",
			EntityID:      user.ID,
			Action:        "APPROVE_STATUS_CHANGE",
			IPAddress:     c.ClientIP(),
			UserAgent:     c.GetHeader("User-Agent"),
			APIEndpoint:   c.Request.URL.Path,
			RequestMethod: c.Request.Method,
			StatusCode:    200,
		}
		config.DB.Create(&auditLog)

		message = "User status change approved and applied successfully"
	} else {
		// Reject request
		pending.Status = "rejected"
		pending.RejectedAt = &now
		pending.RejectionReason = req.Comments

		message = "User status change request rejected successfully"
	}

	// Save pending request
	result = config.DB.Save(&pending)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_UPDATE_FAILED,
			"message": "Failed to update pending request",
			"error":   result.Error.Error(),
		})
		return
	}

	// Load checker admin info for response
	config.DB.Preload("CheckerAdmin").First(&pending, pending.ID)

	// Build response
	response := models.PendingUserStatusChangeResponse{
		ID:                    pending.ID,
		UserID:                pending.UserID,
		UserName:              pending.User.Name,
		UserPhone:             pending.User.Phone,
		MakerAdminID:          pending.MakerAdminID,
		MakerAdminName:        pending.MakerAdmin.Name,
		CurrentStatus:         pending.CurrentStatus,
		CurrentStatusString:   models.GetUserStatusString(pending.CurrentStatus),
		RequestedStatus:       pending.RequestedStatus,
		RequestedStatusString: models.GetUserStatusString(pending.RequestedStatus),
		Reason:                pending.Reason,
		Status:                pending.Status,
		Priority:              pending.Priority,
		ApprovalComments:      pending.ApprovalComments,
		RejectionReason:       pending.RejectionReason,
		ExpiresAt:             pending.ExpiresAt,
		ApprovedAt:            pending.ApprovedAt,
		RejectedAt:            pending.RejectedAt,
		CreatedAt:             pending.CreatedAt,
		UpdatedAt:             pending.UpdatedAt,
	}

	if pending.CheckerAdminID != nil && pending.CheckerAdmin != nil {
		response.CheckerAdminID = pending.CheckerAdminID
		checkerName := pending.CheckerAdmin.Name
		response.CheckerAdminName = &checkerName
	}

	c.JSON(200, gin.H{
		"code":    models.CODE_SUCCESS,
		"message": message,
		"data":    response,
	})
}
