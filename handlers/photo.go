package handlers

import (
	"net/http"
	"strconv"

	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PhotoHandler struct {
	DB *gorm.DB
}

func NewPhotoHandler(db *gorm.DB) *PhotoHandler {
	return &PhotoHandler{
		DB: db,
	}
}

// CreatePhoto creates a new photo
func (h *PhotoHandler) CreatePhoto(c *gin.Context) {
	var req models.CreatePhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Create photo
	photo := models.Photo{
		Image: req.Image,
	}

	if err := h.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusCreated, models.PhotoCreatedResponse(photo))
}

// GetPhotos retrieves all photos with pagination
func (h *PhotoHandler) GetPhotos(c *gin.Context) {
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

	var photos []models.Photo
	var total int64

	// Count total photos
	if err := h.DB.Model(&models.Photo{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Get photos with pagination
	if err := h.DB.Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.PhotosListRetrievedResponse(photos, int(total), page, perPage))
}

// GetPhotoByID retrieves a single photo by ID
func (h *PhotoHandler) GetPhotoByID(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid photo ID",
		})
		return
	}

	var photo models.Photo
	if err := h.DB.First(&photo, uint(photoID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.PhotoNotFoundResponse())
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.PhotoRetrievedResponse(photo))
}

// UpdatePhoto updates an existing photo
func (h *PhotoHandler) UpdatePhoto(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid photo ID",
		})
		return
	}

	var req models.UpdatePhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Find the photo
	var photo models.Photo
	if err := h.DB.First(&photo, uint(photoID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.PhotoNotFoundResponse())
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Update fields if provided
	if req.Image != "" {
		photo.Image = req.Image
	}

	if err := h.DB.Save(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.PhotoUpdatedResponse(photo))
}

// DeletePhoto deletes a photo
func (h *PhotoHandler) DeletePhoto(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid photo ID",
		})
		return
	}

	// Find the photo
	var photo models.Photo
	if err := h.DB.First(&photo, uint(photoID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.PhotoNotFoundResponse())
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Delete the photo (hard delete)
	if err := h.DB.Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.PhotoDeletedResponse())
}
