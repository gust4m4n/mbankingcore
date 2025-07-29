package handlers

import (
	"strconv"

	"mbankingcore/config"
	"mbankingcore/models"

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
		Name:  request.Name,
		Email: request.Email,
		Phone: request.Phone,
	}

	// Set role - validate and default to user if not provided or invalid
	if request.Role != "" {
		// Get current user's role from context if trying to create admin or owner
		if request.Role != models.ROLE_USER {
			currentUserRole, exists := c.Get("role")
			if !exists {
				c.JSON(401, gin.H{
					"code":    models.CODE_UNAUTHORIZED,
					"message": "User role not found in context",
				})
				return
			}

			// Only owner can create admin or owner users
			if currentUserRole != models.ROLE_OWNER {
				c.JSON(403, gin.H{
					"code":    403,
					"message": "Only owner can create admin or owner users",
				})
				return
			}
		}

		if models.ValidateRole(request.Role) {
			user.Role = request.Role
		} else {
			c.JSON(400, gin.H{
				"code":    models.CODE_INVALID_REQUEST,
				"message": "Invalid role. Must be 'user', 'admin', or 'owner'",
			})
			return
		}
	} else {
		user.Role = models.ROLE_USER // Default to user role
	}

	// Save to database
	result := config.DB.Create(&user)
	if result.Error != nil {
		// Check for duplicate email error
		if result.Error.Error() == `ERROR: duplicate key value violates unique constraint "uni_users_email" (SQLSTATE 23505)` {
			c.JSON(409, gin.H{
				"code":    models.CODE_EMAIL_EXISTS,
				"message": "Email already exists",
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

// ListAdminUsers retrieves all admin and owner users with pagination (Admin/Owner only)
func ListAdminUsers(c *gin.Context) {
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

	// Count total admin and owner users
	if err := config.DB.Model(&models.User{}).
		Where("role IN (?)", []string{models.ROLE_ADMIN, models.ROLE_OWNER}).
		Count(&total).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_RETRIEVE_FAILED,
			"message": "Failed to retrieve admin users",
			"error":   err.Error(),
		})
		return
	}

	// Get admin and owner users with pagination
	if err := config.DB.Where("role IN (?)", []string{models.ROLE_ADMIN, models.ROLE_OWNER}).
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&users).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    models.CODE_USER_RETRIEVE_FAILED,
			"message": "Failed to retrieve admin users",
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
	if request.Email != "" {
		updateData["email"] = request.Email
	}
	if request.Phone != "" {
		updateData["phone"] = request.Phone
	}
	if request.Role != "" {
		// Get current user's role from context
		currentUserRole, exists := c.Get("role")
		if !exists {
			c.JSON(401, gin.H{
				"code":    models.CODE_UNAUTHORIZED,
				"message": "User role not found in context",
			})
			return
		}

		// Only owner can change user roles
		if currentUserRole != models.ROLE_OWNER {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "Only owner can change user roles",
			})
			return
		}

		// Validate the new role
		if models.ValidateRole(request.Role) {
			updateData["role"] = request.Role
		} else {
			c.JSON(400, gin.H{
				"code":    models.CODE_INVALID_REQUEST,
				"message": "Invalid role. Must be 'user', 'admin', or 'owner'",
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
		// Check for duplicate email error
		if result.Error.Error() == `ERROR: duplicate key value violates unique constraint "uni_users_email" (SQLSTATE 23505)` {
			c.JSON(409, gin.H{
				"code":    models.CODE_EMAIL_EXISTS,
				"message": "Email already exists",
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
