package models

import (
	"time"
)

type Onboarding struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Image       string    `gorm:"size:500;not null" json:"image"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Onboarding request structures
type CreateOnboardingRequest struct {
	Image       string `json:"image" binding:"required"`
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"required,min=1"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type UpdateOnboardingRequest struct {
	Image       string `json:"image,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// Onboarding response structures
type OnboardingResponse struct {
	ID          uint      `json:"id"`
	Image       string    `json:"image"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OnboardingListResponse struct {
	Onboardings []OnboardingResponse `json:"onboardings"`
	Total       int                  `json:"total"`
}

// Simple list response without pagination
type OnboardingSimpleListResponse struct {
	Onboardings []OnboardingResponse `json:"onboardings"`
}

// Response helper functions
func OnboardingCreatedResponse(onboarding Onboarding) Response {
	return Response{
		Code:    201,
		Message: "Onboarding created successfully",
		Data: OnboardingResponse{
			ID:          onboarding.ID,
			Image:       onboarding.Image,
			Title:       onboarding.Title,
			Description: onboarding.Description,
			IsActive:    onboarding.IsActive,
			CreatedAt:   onboarding.CreatedAt,
			UpdatedAt:   onboarding.UpdatedAt,
		},
	}
}

func OnboardingRetrievedResponse(onboarding Onboarding) Response {
	return Response{
		Code:    200,
		Message: "Onboarding retrieved successfully",
		Data: OnboardingResponse{
			ID:          onboarding.ID,
			Image:       onboarding.Image,
			Title:       onboarding.Title,
			Description: onboarding.Description,
			IsActive:    onboarding.IsActive,
			CreatedAt:   onboarding.CreatedAt,
			UpdatedAt:   onboarding.UpdatedAt,
		},
	}
}

func OnboardingListRetrievedResponse(onboardings []Onboarding, total, page, perPage int) Response {
	var onboardingResponses []OnboardingResponse
	for _, onboarding := range onboardings {
		onboardingResponses = append(onboardingResponses, OnboardingResponse{
			ID:          onboarding.ID,
			Image:       onboarding.Image,
			Title:       onboarding.Title,
			Description: onboarding.Description,
			IsActive:    onboarding.IsActive,
			CreatedAt:   onboarding.CreatedAt,
			UpdatedAt:   onboarding.UpdatedAt,
		})
	}

	return Response{
		Code:    200,
		Message: "Onboardings retrieved successfully",
		Data: OnboardingListResponse{
			Onboardings: onboardingResponses,
			Total:       total,
		},
	}
}

func OnboardingSimpleListRetrievedResponse(onboardings []Onboarding) Response {
	var onboardingResponses []OnboardingResponse
	for _, onboarding := range onboardings {
		onboardingResponses = append(onboardingResponses, OnboardingResponse{
			ID:          onboarding.ID,
			Image:       onboarding.Image,
			Title:       onboarding.Title,
			Description: onboarding.Description,
			IsActive:    onboarding.IsActive,
			CreatedAt:   onboarding.CreatedAt,
			UpdatedAt:   onboarding.UpdatedAt,
		})
	}

	return Response{
		Code:    200,
		Message: "Onboardings retrieved successfully",
		Data: OnboardingSimpleListResponse{
			Onboardings: onboardingResponses,
		},
	}
}

func OnboardingUpdatedResponse(onboarding Onboarding) Response {
	return Response{
		Code:    200,
		Message: "Onboarding updated successfully",
		Data: OnboardingResponse{
			ID:          onboarding.ID,
			Image:       onboarding.Image,
			Title:       onboarding.Title,
			Description: onboarding.Description,
			IsActive:    onboarding.IsActive,
			CreatedAt:   onboarding.CreatedAt,
			UpdatedAt:   onboarding.UpdatedAt,
		},
	}
}

func OnboardingDeletedResponse() Response {
	return Response{
		Code:    200,
		Message: "Onboarding deleted successfully",
		Data:    nil,
	}
}

func OnboardingNotFoundResponse() Response {
	return Response{
		Code:    404,
		Message: "Onboarding not found",
		Data:    nil,
	}
}
