package models

import (
	"time"

	"gorm.io/gorm"
)

// User status constants
const (
	USER_STATUS_INACTIVE           = 0 // inactive
	USER_STATUS_ACTIVE             = 1 // active
	USER_STATUS_BLOCKED            = 2 // blocked
	USER_STATUS_DORMANT            = 3 // dormant
	USER_STATUS_SUSPENDED          = 4 // suspended
	USER_STATUS_CLOSED             = 5 // closed
	USER_STATUS_PENDING_ACTIVATION = 6 // pending activation
	USER_STATUS_FROZEN             = 7 // frozen
	USER_STATUS_LOCKED             = 8 // locked
	USER_STATUS_BLACKLISTED        = 9 // blacklisted
)

type User struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	Name           string          `json:"name" gorm:"not null"`
	Phone          string          `json:"phone" gorm:"unique;not null"`
	MotherName     string          `json:"mother_name" gorm:"not null"`
	PinAtm         string          `json:"-" gorm:"not null"` // Hidden from JSON
	Balance        int64           `json:"balance" gorm:"default:0"`
	Status         int             `json:"status" gorm:"default:1"` // 0=inactive, 1=active, 2=blocked, 3=dormant, 4=suspended, 5=closed, 6=pending_activation, 7=frozen, 8=locked, 9=blacklisted
	Avatar         string          `json:"avatar" gorm:"size:500"`
	BankAccounts   []BankAccount   `json:"bank_accounts,omitempty" gorm:"foreignKey:UserID"`
	DeviceSessions []DeviceSession `json:"device_sessions,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`
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
	return status >= USER_STATUS_INACTIVE && status <= USER_STATUS_BLACKLISTED
}

// GetStatusString returns string representation of user status
func GetUserStatusString(status int) string {
	switch status {
	case USER_STATUS_INACTIVE:
		return "Inactive"
	case USER_STATUS_ACTIVE:
		return "Active"
	case USER_STATUS_BLOCKED:
		return "Blocked"
	case USER_STATUS_DORMANT:
		return "Dormant"
	case USER_STATUS_SUSPENDED:
		return "Suspended"
	case USER_STATUS_CLOSED:
		return "Closed"
	case USER_STATUS_PENDING_ACTIVATION:
		return "Pending Activation"
	case USER_STATUS_FROZEN:
		return "Frozen"
	case USER_STATUS_LOCKED:
		return "Locked"
	case USER_STATUS_BLACKLISTED:
		return "Blacklisted"
	default:
		return "Unknown"
	}
}

// Update Status Request
type UpdateUserStatusRequest struct {
	Status int    `json:"status" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

// Update Status Response
type UpdateUserStatusResponse struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Phone                string    `json:"phone"`
	PreviousStatus       int       `json:"previous_status"`
	NewStatus            int       `json:"new_status"`
	PreviousStatusString string    `json:"previous_status_string"`
	NewStatusString      string    `json:"new_status_string"`
	Reason               string    `json:"reason"`
	UpdatedBy            string    `json:"updated_by"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Soft Delete Response
type UserSoftDeletedResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	DeletedAt time.Time `json:"deleted_at"`
}

func UserSoftDeletedSuccessResponse(user *User) Response {
	return Response{
		Code:    200,
		Message: "User soft deleted successfully",
		Data: UserSoftDeletedResponse{
			ID:        user.ID,
			Name:      user.Name,
			Phone:     user.Phone,
			DeletedAt: user.DeletedAt.Time,
		},
	}
}

// User Restore Response
type UserRestoredResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Status int    `json:"status"`
}

func UserRestoredSuccessResponse(user *User) Response {
	return Response{
		Code:    200,
		Message: "User restored successfully",
		Data: UserRestoredResponse{
			ID:     user.ID,
			Name:   user.Name,
			Phone:  user.Phone,
			Status: user.Status,
		},
	}
}
