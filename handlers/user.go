package handlers

import (
	"strconv"

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
