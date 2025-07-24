package handlers

import (
	"mbankingcore/config"
	"mbankingcore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateOnboarding - Create a new onboarding
func CreateOnboarding(c *gin.Context) {
	var req models.CreateOnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	onboarding := models.Onboarding{
		Image:       req.Image,
		Title:       req.Title,
		Description: req.Description,
	}

	// Set is_active field - default to true if not provided
	if req.IsActive != nil {
		onboarding.IsActive = *req.IsActive
	} else {
		onboarding.IsActive = true
	}

	if err := config.DB.Create(&onboarding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusCreated, models.OnboardingCreatedResponse(onboarding))
}

// GetOnboardings - Get list of onboardings without pagination
func GetOnboardings(c *gin.Context) {
	var onboardings []models.Onboarding

	// Get all onboardings ordered by created_at
	if err := config.DB.Order("created_at ASC").Find(&onboardings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.OnboardingSimpleListRetrievedResponse(onboardings))
}

// GetOnboarding - Get a single onboarding by ID
func GetOnboarding(c *gin.Context) {
	id := c.Param("id")
	var onboarding models.Onboarding

	if err := config.DB.First(&onboarding, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.OnboardingNotFoundResponse())
		return
	}

	c.JSON(http.StatusOK, models.OnboardingRetrievedResponse(onboarding))
}

// UpdateOnboarding - Update an onboarding
func UpdateOnboarding(c *gin.Context) {
	id := c.Param("id")
	var onboarding models.Onboarding

	// Check if onboarding exists
	if err := config.DB.First(&onboarding, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.OnboardingNotFoundResponse())
		return
	}

	var req models.UpdateOnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	
	if req.Image != "" {
		updates["image"] = req.Image
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "No fields to update",
			Data:    nil,
		})
		return
	}

	if err := config.DB.Model(&onboarding).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Reload the onboarding to get updated data
	config.DB.First(&onboarding, id)

	c.JSON(http.StatusOK, models.OnboardingUpdatedResponse(onboarding))
}

// DeleteOnboarding - Delete an onboarding (hard delete)
func DeleteOnboarding(c *gin.Context) {
	id := c.Param("id")
	var onboarding models.Onboarding

	// Check if onboarding exists
	if err := config.DB.First(&onboarding, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.OnboardingNotFoundResponse())
		return
	}

	// Hard delete
	if err := config.DB.Delete(&onboarding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.OnboardingDeletedResponse())
}
