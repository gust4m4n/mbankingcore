package handlers

import (
	"net/http"

	"mbankingcore/config"
	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPrivacyPolicy retrieves the currently active privacy policy from config
func GetPrivacyPolicy(c *gin.Context) {
	db := config.GetDB()
	
	var configPrivacy models.Config
	err := db.Where("key = ?", "privacy-policy").First(&configPrivacy).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(http.StatusNotFound, "Privacy policy not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve privacy policy"))
		return
	}

	// Return privacy policy content from config
	c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Privacy policy retrieved successfully", gin.H{
		"content": configPrivacy.Value,
		"updated_at": configPrivacy.UpdatedAt,
	}))
}

// SetPrivacyPolicy sets the privacy policy content using config API (admin/owner only)
func SetPrivacyPolicy(c *gin.Context) {
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
	err := db.Where("key = ?", "privacy-policy").First(&existingConfig).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to check existing privacy policy"))
		return
	}

	if err == gorm.ErrRecordNotFound {
		// Create new privacy policy config
		configPrivacy := models.Config{
			Key:   "privacy-policy",
			Value: request.Content,
		}

		if err := db.Create(&configPrivacy).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to create privacy policy"))
			return
		}

		c.JSON(http.StatusCreated, models.NewSuccessResponse(http.StatusCreated, "Privacy policy created successfully", gin.H{
			"content": configPrivacy.Value,
			"updated_at": configPrivacy.UpdatedAt,
		}))
	} else {
		// Update existing privacy policy config
		existingConfig.Value = request.Content

		if err := db.Save(&existingConfig).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to update privacy policy"))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(http.StatusOK, "Privacy policy updated successfully", gin.H{
			"content": existingConfig.Value,
			"updated_at": existingConfig.UpdatedAt,
		}))
	}
}
