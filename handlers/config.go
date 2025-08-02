package handlers

import (
	"net/http"

	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetConfig sets a configuration value
func SetConfig(c *gin.Context) {
	var request models.ConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Check if config key already exists
	var existingConfig models.Config
	err := db.Where("key = ?", request.Key).First(&existingConfig).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to check existing config"))
		return
	}

	if err == gorm.ErrRecordNotFound {
		// Create new config
		config := models.Config{
			Key:   request.Key,
			Value: request.Value,
		}

		if err := db.Create(&config).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to create config"))
			return
		}

		c.JSON(http.StatusCreated, models.NewSuccessResponse(http.StatusCreated, "Config created successfully", config.ToResponse()))
	} else {
		// Update existing config
		existingConfig.Value = request.Value

		if err := db.Save(&existingConfig).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to update config"))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Config updated successfully", existingConfig.ToResponse()))
	}
}

// GetConfig retrieves a configuration value by key
func GetConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Config key is required"))
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var config models.Config
	if err := db.Where("key = ?", key).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(http.StatusNotFound, "Config not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve config"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Config retrieved successfully", config.ToResponse()))
}

// GetAllConfigs retrieves all configuration values (admin only)
func GetAllConfigs(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var configs []models.Config
	if err := db.Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve configs"))
		return
	}

	var responses []models.ConfigResponse
	for _, config := range configs {
		responses = append(responses, config.ToResponse())
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Configs retrieved successfully", responses))
}

// DeleteConfig deletes a configuration by key (admin only)
func DeleteConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Config key is required"))
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var config models.Config
	if err := db.Where("key = ?", key).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(http.StatusNotFound, "Config not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to find config"))
		return
	}

	if err := db.Delete(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to delete config"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Config deleted successfully", nil))
}

// GetAdminTermsConditions retrieves admin terms and conditions from config
func GetAdminTermsConditions(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var config models.Config
	if err := db.Where("key = ?", "admin_terms_conditions").First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(http.StatusNotFound, "Admin terms and conditions not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve admin terms and conditions"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Admin terms and conditions retrieved successfully", gin.H{
		"content": config.Value,
	}))
}

// SetAdminTermsConditions sets admin terms and conditions via config
func SetAdminTermsConditions(c *gin.Context) {
	var request struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Check if config key already exists
	var existingConfig models.Config
	err := db.Where("key = ?", "admin_terms_conditions").First(&existingConfig).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to check existing config"))
		return
	}

	if err == gorm.ErrRecordNotFound {
		// Create new config
		config := models.Config{
			Key:   "admin_terms_conditions",
			Value: request.Content,
		}

		if err := db.Create(&config).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to create admin terms and conditions"))
			return
		}

		c.JSON(http.StatusCreated, models.NewSuccessResponse(http.StatusCreated, "Admin terms and conditions created successfully", gin.H{
			"content": config.Value,
		}))
	} else {
		// Update existing config
		existingConfig.Value = request.Content

		if err := db.Save(&existingConfig).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to update admin terms and conditions"))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Admin terms and conditions updated successfully", gin.H{
			"content": existingConfig.Value,
		}))
	}
}

// GetAdminPrivacyPolicy retrieves admin privacy policy from config
func GetAdminPrivacyPolicy(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var config models.Config
	if err := db.Where("key = ?", "admin_privacy_policy").First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(http.StatusNotFound, "Admin privacy policy not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve admin privacy policy"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Admin privacy policy retrieved successfully", gin.H{
		"content": config.Value,
	}))
}

// SetAdminPrivacyPolicy sets admin privacy policy via config
func SetAdminPrivacyPolicy(c *gin.Context) {
	var request struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Check if config key already exists
	var existingConfig models.Config
	err := db.Where("key = ?", "admin_privacy_policy").First(&existingConfig).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to check existing config"))
		return
	}

	if err == gorm.ErrRecordNotFound {
		// Create new config
		config := models.Config{
			Key:   "admin_privacy_policy",
			Value: request.Content,
		}

		if err := db.Create(&config).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to create admin privacy policy"))
			return
		}

		c.JSON(http.StatusCreated, models.NewSuccessResponse(http.StatusCreated, "Admin privacy policy created successfully", gin.H{
			"content": config.Value,
		}))
	} else {
		// Update existing config
		existingConfig.Value = request.Content

		if err := db.Save(&existingConfig).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to update admin privacy policy"))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Admin privacy policy updated successfully", gin.H{
			"content": existingConfig.Value,
		}))
	}
}
