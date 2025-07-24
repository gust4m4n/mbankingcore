package models

import (
	"time"
)

type Article struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Image     string    `gorm:"size:500" json:"image"`
	Content   string    `gorm:"type:text" json:"content"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Article request structures
type CreateArticleRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=255"`
	Image    string `json:"image"`
	Content  string `json:"content" binding:"required,min=1"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type UpdateArticleRequest struct {
	Title    string `json:"title,omitempty"`
	Image    string `json:"image,omitempty"`
	Content  string `json:"content,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// Article response structures
type ArticleResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	Content   string    `json:"content"`
	IsActive  bool      `json:"is_active"`
	UserID    uint      `json:"user_id"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ArticlesListResponse struct {
	Articles []ArticleResponse `json:"articles"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PerPage  int               `json:"per_page"`
}

// Response helper functions
func ArticleCreatedResponse(article Article) Response {
	return Response{
		Code:    201,
		Message: "Article created successfully",
		Data: ArticleResponse{
			ID:        article.ID,
			Title:     article.Title,
			Image:     article.Image,
			Content:   article.Content,
			IsActive:  article.IsActive,
			UserID:    article.UserID,
			Author:    article.User.Name,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		},
	}
}

func ArticleRetrievedResponse(article Article) Response {
	return Response{
		Code:    200,
		Message: "Article retrieved successfully",
		Data: ArticleResponse{
			ID:        article.ID,
			Title:     article.Title,
			Image:     article.Image,
			Content:   article.Content,
			IsActive:  article.IsActive,
			UserID:    article.UserID,
			Author:    article.User.Name,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		},
	}
}

func ArticlesListRetrievedResponse(articles []Article, total, page, perPage int) Response {
	var articleResponses []ArticleResponse
	for _, article := range articles {
		articleResponses = append(articleResponses, ArticleResponse{
			ID:        article.ID,
			Title:     article.Title,
			Image:     article.Image,
			Content:   article.Content,
			IsActive:  article.IsActive,
			UserID:    article.UserID,
			Author:    article.User.Name,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		})
	}

	return Response{
		Code:    200,
		Message: "Articles retrieved successfully",
		Data: ArticlesListResponse{
			Articles: articleResponses,
			Total:    total,
			Page:     page,
			PerPage:  perPage,
		},
	}
}

func ArticleUpdatedResponse(article Article) Response {
	return Response{
		Code:    200,
		Message: "Article updated successfully",
		Data: ArticleResponse{
			ID:        article.ID,
			Title:     article.Title,
			Image:     article.Image,
			Content:   article.Content,
			IsActive:  article.IsActive,
			UserID:    article.UserID,
			Author:    article.User.Name,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		},
	}
}

func ArticleDeletedResponse() Response {
	return Response{
		Code:    200,
		Message: "Article deleted successfully",
		Data:    nil,
	}
}

func ArticleNotFoundResponse() Response {
	return Response{
		Code:    404,
		Message: "Article not found",
		Data:    nil,
	}
}
