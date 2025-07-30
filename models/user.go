package models

import (
	"time"
)

// User status constants
const (
	USER_STATUS_INACTIVE = 0 // inactive
	USER_STATUS_ACTIVE   = 1 // active
	USER_STATUS_BLOCKED  = 2 // terblokir
)

type User struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	Name           string          `json:"name" gorm:"not null"`
	Phone          string          `json:"phone" gorm:"unique;not null"`
	MotherName     string          `json:"mother_name" gorm:"not null"`
	PinAtm         string          `json:"-" gorm:"not null"` // Hidden from JSON
	Balance        int64           `json:"balance" gorm:"default:0"`
	Status         int             `json:"status" gorm:"default:1"` // 0=inactive, 1=active, 2=suspended, 3=terblokir
	Avatar         string          `json:"avatar" gorm:"size:500"`
	BankAccounts   []BankAccount   `json:"bank_accounts,omitempty" gorm:"foreignKey:UserID"`
	DeviceSessions []DeviceSession `json:"device_sessions,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// Action-based Request Structure
type UserActionRequest struct {
	Action string `json:"action" binding:"required"`
	Name   string `json:"name,omitempty"`
	Phone  string `json:"phone,omitempty"`
	ID     uint   `json:"id,omitempty"`
}

type UserResponse struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Phone        string        `json:"phone"`
	MotherName   string        `json:"mother_name"`
	Balance      int64         `json:"balance"`
	Status       int           `json:"status"`
	Avatar       string        `json:"avatar"`
	BankAccounts []BankAccount `json:"bank_accounts,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type UsersListResponse struct {
	Users   []UserResponse `json:"users"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Name:         u.Name,
		Phone:        u.Phone,
		MotherName:   u.MotherName,
		Balance:      u.Balance,
		Status:       u.Status,
		Avatar:       u.Avatar,
		BankAccounts: u.BankAccounts,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
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

// Status validation helper functions
func (u *User) IsActive() bool {
	return u.Status == USER_STATUS_ACTIVE
}

func (u *User) IsInactive() bool {
	return u.Status == USER_STATUS_INACTIVE
}

func (u *User) IsBlocked() bool {
	return u.Status == USER_STATUS_BLOCKED
}

func ValidateStatus(status int) bool {
	return status == USER_STATUS_INACTIVE || status == USER_STATUS_ACTIVE || status == USER_STATUS_BLOCKED
}
