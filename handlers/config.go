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
