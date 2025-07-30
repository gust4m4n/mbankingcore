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

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	var users []models.User
	var total int64

	// Count total users
	if err := config.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_RETRIEVE_FAILED,
			"message": "Failed to retrieve users",
			"error":   err.Error(),
		})
		return
	}

	// Get users with pagination
	if err := config.DB.Order("created_at DESC").
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
	idParam := c.Param("id")
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

// DeleteUser deletes a user by ID
func DeleteUser(c *gin.Context) {
	// Get ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_USER_ID,
			"message": "Invalid user ID",
		})
		return
	}

	var user models.User

	// Find user first to get the name for response
	result := config.DB.First(&user, uint(id))
	if result.Error != nil {
		c.JSON(404, gin.H{
			"code":    models.CODE_USER_NOT_FOUND,
			"message": "User not found",
		})
		return
	}

	// Delete the user
	result = config.DB.Delete(&user, uint(id))
	if result.Error != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_DELETE_FAILED,
			"message": "Failed to delete user",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    models.CODE_SUCCESS,
		"message": "User deleted successfully",
		"data": gin.H{
			"id":   user.ID,
			"name": user.Name,
		},
	})
}
