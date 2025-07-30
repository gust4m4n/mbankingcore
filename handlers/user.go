package handlers

import (
	"strconv"

	"mbankingcore/config"
	"mbankingcore/models"
	"mbankingcore/utils"

	"github.com/gin-gonic/gin"
)

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var request models.CreateUserRequest

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Create new user
	user := models.User{
		Name:       request.Name,
		Phone:      request.Phone,
		MotherName: request.MotherName,
	}

	// Hash the PIN
	hashedPin, err := utils.HashPassword(request.PinAtm)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_INTERNAL_SERVER,
			"message": "Failed to hash PIN",
		})
		return
	}
	user.PinAtm = hashedPin

	// Save to database
	result := config.DB.Create(&user)
	if result.Error != nil {
		// Check for duplicate phone error
		if result.Error.Error() == `ERROR: duplicate key value violates unique constraint "uni_users_phone" (SQLSTATE 23505)` {
			c.JSON(409, gin.H{
				"code":    models.CODE_PHONE_EXISTS,
				"message": "Phone already exists",
			})
			return
		}

		c.JSON(500, gin.H{
			"code":    models.CODE_USER_CREATE_FAILED,
			"message": "Failed to create user",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    models.CODE_SUCCESS,
		"message": "User created successfully",
		"data":    user.ToResponse(),
	})
}

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

// UpdateUser updates an existing user by ID
func UpdateUser(c *gin.Context) {
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

	var request models.UpdateUserRequest

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	var user models.User

	// Find user first to ensure it exists
	result := config.DB.First(&user, uint(id))
	if result.Error != nil {
		c.JSON(404, gin.H{
			"code":    models.CODE_USER_NOT_FOUND,
			"message": "User not found",
		})
		return
	}

	// Prepare update data (only update non-empty fields)
	updateData := make(map[string]interface{})

	if request.Name != "" {
		updateData["name"] = request.Name
	}
	if request.Phone != "" {
		updateData["phone"] = request.Phone
	}
	if request.MotherName != "" {
		updateData["mother_name"] = request.MotherName
	}
	if request.Balance != nil {
		updateData["balance"] = *request.Balance
	}
	if request.Status != nil {
		if models.ValidateStatus(*request.Status) {
			updateData["status"] = *request.Status
		} else {
			c.JSON(400, gin.H{
				"code":    models.CODE_INVALID_REQUEST,
				"message": "Invalid status. Must be 0 (inactive), 1 (active), or 2 (blocked)",
			})
			return
		}
	}

	// If no fields to update
	if len(updateData) == 0 {
		c.JSON(400, gin.H{
			"code":    models.CODE_INVALID_REQUEST,
			"message": "No fields to update",
		})
		return
	}

	// Update user in database
	result = config.DB.Model(&user).Updates(updateData)
	if result.Error != nil {
		// Check for duplicate phone error
		if result.Error.Error() == `ERROR: duplicate key value violates unique constraint "uni_users_phone" (SQLSTATE 23505)` {
			c.JSON(409, gin.H{
				"code":    models.CODE_PHONE_EXISTS,
				"message": "Phone already exists",
			})
			return
		}

		c.JSON(500, gin.H{
			"code":    models.CODE_USER_UPDATE_FAILED,
			"message": "Failed to update user",
			"error":   result.Error.Error(),
		})
		return
	}

	// Reload user to get updated data
	config.DB.First(&user, uint(id))

	c.JSON(200, gin.H{
		"code":    models.CODE_SUCCESS,
		"message": "User updated successfully",
		"data":    user.ToResponse(),
	})
}
