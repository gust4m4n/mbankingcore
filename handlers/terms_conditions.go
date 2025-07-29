package handlers

import (
	"net/http"

	"mbankingcore/config"
	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTermsConditions retrieves the currently active terms and conditions from config
func GetTermsConditions(c *gin.Context) {
	db := config.GetDB()

	var configTnc models.Config
	err := db.Where("key = ?", "tnc").First(&configTnc).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(http.StatusNotFound, "Terms and conditions not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve terms and conditions"))
		return
	}

	// Return terms content from config
	c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Terms and conditions retrieved successfully", gin.H{
		"content":    configTnc.Value,
		"updated_at": configTnc.UpdatedAt,
	}))
}

// SetTermsConditions sets the terms and conditions content using config API (admin/owner only)
func SetTermsConditions(c *gin.Context) {
	var request struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	db := config.GetDB()

	// Check if config key already exists
	var existingConfig models.Config
	err := db.Where("key = ?", "tnc").First(&existingConfig).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to check existing terms and conditions"))
		return
	}

	if err == gorm.ErrRecordNotFound {
		// Create new terms and conditions config
		configTnc := models.Config{
			Key:   "tnc",
			Value: request.Content,
		}

		if err := db.Create(&configTnc).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to create terms and conditions"))
			return
		}

		c.JSON(http.StatusCreated, models.NewSuccessResponse(http.StatusCreated, "Terms and conditions created successfully", gin.H{
			"content":    configTnc.Value,
			"updated_at": configTnc.UpdatedAt,
		}))
	} else {
		// Update existing terms and conditions config
		existingConfig.Value = request.Content

		if err := db.Save(&existingConfig).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to update terms and conditions"))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Terms and conditions updated successfully", gin.H{
			"content":    existingConfig.Value,
			"updated_at": existingConfig.UpdatedAt,
		}))
	}
}
