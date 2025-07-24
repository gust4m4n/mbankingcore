package handlers

import (
	"net/http"
	"strconv"

	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	DB *gorm.DB
}

func NewArticleHandler(db *gorm.DB) *ArticleHandler {
	return &ArticleHandler{
		DB: db,
	}
}

// CreateArticle creates a new article
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	var req models.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Create article
	article := models.Article{
		Title:   req.Title,
		Image:   req.Image,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	// Set is_active field - default to true if not provided
	if req.IsActive != nil {
		article.IsActive = *req.IsActive
	} else {
		article.IsActive = true
	}

	if err := h.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Load user information for response
	if err := h.DB.Preload("User").First(&article, article.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusCreated, models.ArticleCreatedResponse(article))
}

// GetArticles retrieves all articles with pagination
func (h *ArticleHandler) GetArticles(c *gin.Context) {
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

	var articles []models.Article
	var total int64

	// Count total articles
	if err := h.DB.Model(&models.Article{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Get articles with pagination and preload user
	if err := h.DB.Preload("User").
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.ArticlesListRetrievedResponse(articles, int(total), page, perPage))
}

// GetArticleByID retrieves a single article by ID
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid article ID",
		})
		return
	}

	var article models.Article
	if err := h.DB.Preload("User").First(&article, uint(articleID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ArticleNotFoundResponse())
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.ArticleRetrievedResponse(article))
}

// UpdateArticle updates an existing article
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid article ID",
		})
		return
	}

	var req models.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Find the article
	var article models.Article
	if err := h.DB.First(&article, uint(articleID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ArticleNotFoundResponse())
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Check if user owns the article
	if article.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "You don't have permission to update this article",
		})
		return
	}

	// Update fields if provided
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Image != "" {
		article.Image = req.Image
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.IsActive != nil {
		article.IsActive = *req.IsActive
	}

	if err := h.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Load user information for response
	if err := h.DB.Preload("User").First(&article, article.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.ArticleUpdatedResponse(article))
}

// DeleteArticle deletes an article
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid article ID",
		})
		return
	}

	// Find the article
	var article models.Article
	if err := h.DB.First(&article, uint(articleID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ArticleNotFoundResponse())
			return
		}
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Check if user owns the article
	if article.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "You don't have permission to delete this article",
		})
		return
	}

	// Delete the article (hard delete)
	if err := h.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.ArticleDeletedResponse())
}

// GetMyArticles retrieves articles created by the current user
func (h *ArticleHandler) GetMyArticles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

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

	var articles []models.Article
	var total int64

	// Count total articles for this user
	if err := h.DB.Model(&models.Article{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Get user's articles with pagination
	if err := h.DB.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.ArticlesListRetrievedResponse(articles, int(total), page, perPage))
}
