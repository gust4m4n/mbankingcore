package models

import (
	"time"
)

type User struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	Name           string          `json:"name" gorm:"not null"`
	Email          string          `json:"email" gorm:"unique;not null"`
	Phone          string          `json:"phone"`
	Password       string          `json:"-" gorm:""`
	Role           string          `json:"role" gorm:"size:20;default:'user'"`
	GoogleID       string          `json:"google_id" gorm:"size:255"`
	AppleID        string          `json:"apple_id" gorm:"size:255"`
	FacebookID     string          `json:"facebook_id" gorm:"size:255"`
	EmailVerified  bool            `json:"email_verified" gorm:"default:false"`
	Avatar         string          `json:"avatar" gorm:"size:500"`
	DeviceSessions []DeviceSession `json:"device_sessions,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone"`
	Role  string `json:"role,omitempty"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Role  string `json:"role,omitempty"`
}

// Action-based Request Structure
type UserActionRequest struct {
	Action string `json:"action" binding:"required"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Role   string `json:"role,omitempty"`
	ID     uint   `json:"id,omitempty"`
}

type UserResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	Avatar        string    `json:"avatar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UsersListResponse struct {
	Users   []UserResponse `json:"users"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		Phone:         u.Phone,
		Role:          u.Role,
		EmailVerified: u.EmailVerified,
		Avatar:        u.Avatar,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// Response helper functions for users
func UsersListRetrievedResponse(users []User, total, page, perPage int) Response {
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return Response{
		Code:    200,
		Message: "Users retrieved successfully",
		Data: UsersListResponse{
			Users:   userResponses,
			Total:   total,
			Page:    page,
			PerPage: perPage,
		},
	}
}

// Role validation helper functions
func (u *User) IsAdmin() bool {
	return u.Role == ROLE_ADMIN
}

func (u *User) IsUser() bool {
	return u.Role == ROLE_USER
}

func (u *User) IsOwner() bool {
	return u.Role == ROLE_OWNER
}

func (u *User) CanManageRoles() bool {
	return u.Role == ROLE_OWNER
}

func ValidateRole(role string) bool {
	return role == ROLE_USER || role == ROLE_ADMIN || role == ROLE_OWNER
}
