package models

import (
	"time"

	"gorm.io/gorm"
)

// Admin status constants
const (
	ADMIN_STATUS_INACTIVE = 0 // inactive
	ADMIN_STATUS_ACTIVE   = 1 // active
	ADMIN_STATUS_BLOCKED  = 2 // terblokir
)

// Admin role constants
const (
	ADMIN_ROLE_ADMIN = "admin" // admin
	ADMIN_ROLE_SUPER = "super" // super admin
)

type Admin struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"` // Hidden from JSON
	Role      string         `json:"role" gorm:"default:admin"`
	Status    int            `json:"status" gorm:"default:1"` // 0=inactive, 1=active, 2=blocked
	Avatar    string         `json:"avatar" gorm:"size:500"`
	LastLogin *time.Time     `json:"last_login,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateAdminRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin super"`
}

type UpdateAdminRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
	Status   *int   `json:"status,omitempty"`
}

type AdminResponse struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Status    int        `json:"status"`
	Avatar    string     `json:"avatar"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type AdminLoginResponse struct {
	Admin       AdminResponse `json:"admin"`
	AccessToken string        `json:"access_token"`
	ExpiresIn   int64         `json:"expires_in"`
}

type AdminListResponse struct {
	Admins  []AdminResponse `json:"admins"`
	Total   int             `json:"total"`
	Page    int             `json:"page"`
	PerPage int             `json:"per_page"`
}

func (a *Admin) ToResponse() AdminResponse {
	return AdminResponse{
		ID:        a.ID,
		Name:      a.Name,
		Email:     a.Email,
		Role:      a.Role,
		Status:    a.Status,
		Avatar:    a.Avatar,
		LastLogin: a.LastLogin,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

// Status validation helper functions
func (a *Admin) IsActive() bool {
	return a.Status == ADMIN_STATUS_ACTIVE
}

func (a *Admin) IsInactive() bool {
	return a.Status == ADMIN_STATUS_INACTIVE
}

func (a *Admin) IsBlocked() bool {
	return a.Status == ADMIN_STATUS_BLOCKED
}

func (a *Admin) IsSuper() bool {
	return a.Role == ADMIN_ROLE_SUPER
}

func ValidateAdminStatus(status int) bool {
	return status == ADMIN_STATUS_INACTIVE || status == ADMIN_STATUS_ACTIVE || status == ADMIN_STATUS_BLOCKED
}

func ValidateAdminRole(role string) bool {
	return role == ADMIN_ROLE_ADMIN || role == ADMIN_ROLE_SUPER
}

// Response helper functions for admins
func AdminListRetrievedResponse(admins []Admin, total, page, perPage int) Response {
	var adminResponses []AdminResponse
	for _, admin := range admins {
		adminResponses = append(adminResponses, admin.ToResponse())
	}

	return Response{
		Code:    200,
		Message: "Admins retrieved successfully",
		Data: AdminListResponse{
			Admins:  adminResponses,
			Total:   total,
			Page:    page,
			PerPage: perPage,
		},
	}
}

func AdminCreatedResponse(admin Admin) Response {
	return Response{
		Code:    201,
		Message: "Admin created successfully",
		Data:    admin.ToResponse(),
	}
}

func AdminUpdatedResponse(admin Admin) Response {
	return Response{
		Code:    200,
		Message: "Admin updated successfully",
		Data:    admin.ToResponse(),
	}
}

func AdminDeletedResponse() Response {
	return Response{
		Code:    200,
		Message: "Admin deleted successfully",
		Data:    nil,
	}
}

func AdminLoginSuccessResponse(admin Admin, token string, expiresIn int64) Response {
	return Response{
		Code:    200,
		Message: "Admin login successful",
		Data: AdminLoginResponse{
			Admin:       admin.ToResponse(),
			AccessToken: token,
			ExpiresIn:   expiresIn,
		},
	}
}

func AdminLogoutSuccessResponse() Response {
	return Response{
		Code:    200,
		Message: "Admin logout successful",
		Data:    nil,
	}
}

// Admin Soft Delete Response
type AdminSoftDeletedResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	DeletedAt time.Time `json:"deleted_at"`
}

func AdminSoftDeletedSuccessResponse(admin *Admin) Response {
	return Response{
		Code:    200,
		Message: "Admin soft deleted successfully",
		Data: AdminSoftDeletedResponse{
			ID:        admin.ID,
			Name:      admin.Name,
			Email:     admin.Email,
			Role:      admin.Role,
			DeletedAt: admin.DeletedAt.Time,
		},
	}
}

// Admin Restore Response
type AdminRestoredResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status int    `json:"status"`
}

func AdminRestoredSuccessResponse(admin *Admin) Response {
	return Response{
		Code:    200,
		Message: "Admin restored successfully",
		Data: AdminRestoredResponse{
			ID:     admin.ID,
			Name:   admin.Name,
			Email:  admin.Email,
			Role:   admin.Role,
			Status: admin.Status,
		},
	}
}
