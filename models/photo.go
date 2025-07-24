package models

import (
	"time"
)

type Photo struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Image     string    `gorm:"size:500;not null" json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Photo request structures
type CreatePhotoRequest struct {
	Image string `json:"image" binding:"required"`
}

type UpdatePhotoRequest struct {
	Image string `json:"image,omitempty"`
}

// Photo response structures
type PhotoResponse struct {
	ID        uint      `json:"id"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type PhotosListResponse struct {
	Photos  []PhotoResponse `json:"photos"`
	Total   int             `json:"total"`
	Page    int             `json:"page"`
	PerPage int             `json:"per_page"`
}

// Response helper functions
func PhotoCreatedResponse(photo Photo) Response {
	return Response{
		Code:    201,
		Message: "Photo created successfully",
		Data: PhotoResponse{
			ID:        photo.ID,
			Image:     photo.Image,
			CreatedAt: photo.CreatedAt,
		},
	}
}

func PhotoRetrievedResponse(photo Photo) Response {
	return Response{
		Code:    200,
		Message: "Photo retrieved successfully",
		Data: PhotoResponse{
			ID:        photo.ID,
			Image:     photo.Image,
			CreatedAt: photo.CreatedAt,
		},
	}
}

func PhotosListRetrievedResponse(photos []Photo, total, page, perPage int) Response {
	var photoResponses []PhotoResponse
	for _, photo := range photos {
		photoResponses = append(photoResponses, PhotoResponse{
			ID:        photo.ID,
			Image:     photo.Image,
			CreatedAt: photo.CreatedAt,
		})
	}

	return Response{
		Code:    200,
		Message: "Photos retrieved successfully",
		Data: PhotosListResponse{
			Photos:  photoResponses,
			Total:   total,
			Page:    page,
			PerPage: perPage,
		},
	}
}

func PhotoUpdatedResponse(photo Photo) Response {
	return Response{
		Code:    200,
		Message: "Photo updated successfully",
		Data: PhotoResponse{
			ID:        photo.ID,
			Image:     photo.Image,
			CreatedAt: photo.CreatedAt,
		},
	}
}

func PhotoDeletedResponse() Response {
	return Response{
		Code:    200,
		Message: "Photo deleted successfully",
		Data:    nil,
	}
}

func PhotoNotFoundResponse() Response {
	return Response{
		Code:    404,
		Message: "Photo not found",
		Data:    nil,
	}
}
