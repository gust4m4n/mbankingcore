package models

import (
	"time"

	"gorm.io/gorm"
)

// PendingUserStatusChange represents user status change requests that require approval
type PendingUserStatusChange struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	UserID           uint           `json:"user_id" gorm:"not null;index"`           // Target user
	MakerAdminID     uint           `json:"maker_admin_id" gorm:"not null;index"`    // Admin who initiated (maker)
	CheckerAdminID   *uint          `json:"checker_admin_id,omitempty" gorm:"index"` // Admin who approved/rejected (checker)
	CurrentStatus    int            `json:"current_status" gorm:"not null"`          // Current user status
	RequestedStatus  int            `json:"requested_status" gorm:"not null"`        // Requested new status
	Reason           string         `json:"reason" gorm:"not null"`                  // Reason for status change
	Status           string         `json:"status" gorm:"default:'pending'"`         // "pending", "approved", "rejected", "expired"
	Priority         string         `json:"priority" gorm:"default:'normal'"`        // "low", "normal", "high", "critical"
	ApprovalComments string         `json:"approval_comments"`                       // Comments from checker
	RejectionReason  string         `json:"rejection_reason"`                        // Reason for rejection
	ExpiresAt        *time.Time     `json:"expires_at"`                              // When the pending request expires
	ApprovedAt       *time.Time     `json:"approved_at"`                             // When approved
	RejectedAt       *time.Time     `json:"rejected_at"`                             // When rejected
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	User         User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	MakerAdmin   Admin  `json:"maker_admin,omitempty" gorm:"foreignKey:MakerAdminID"`
	CheckerAdmin *Admin `json:"checker_admin,omitempty" gorm:"foreignKey:CheckerAdminID"`
}

// Pending Status Change Request
type CreatePendingUserStatusChangeRequest struct {
	Status   int    `json:"status" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
	Priority string `json:"priority,omitempty"` // optional: low, normal, high, critical
}

// Approval/Rejection Request
type ReviewPendingUserStatusChangeRequest struct {
	Action   string `json:"action" binding:"required"` // "approve" or "reject"
	Comments string `json:"comments,omitempty"`        // Optional comments
}

// Response structures
type PendingUserStatusChangeResponse struct {
	ID                    uint       `json:"id"`
	UserID                uint       `json:"user_id"`
	UserName              string     `json:"user_name"`
	UserPhone             string     `json:"user_phone"`
	MakerAdminID          uint       `json:"maker_admin_id"`
	MakerAdminName        string     `json:"maker_admin_name"`
	CheckerAdminID        *uint      `json:"checker_admin_id,omitempty"`
	CheckerAdminName      *string    `json:"checker_admin_name,omitempty"`
	CurrentStatus         int        `json:"current_status"`
	CurrentStatusString   string     `json:"current_status_string"`
	RequestedStatus       int        `json:"requested_status"`
	RequestedStatusString string     `json:"requested_status_string"`
	Reason                string     `json:"reason"`
	Status                string     `json:"status"`
	Priority              string     `json:"priority"`
	ApprovalComments      string     `json:"approval_comments,omitempty"`
	RejectionReason       string     `json:"rejection_reason,omitempty"`
	ExpiresAt             *time.Time `json:"expires_at,omitempty"`
	ApprovedAt            *time.Time `json:"approved_at,omitempty"`
	RejectedAt            *time.Time `json:"rejected_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type PendingUserStatusChangeListResponse struct {
	PendingChanges []PendingUserStatusChangeResponse `json:"pending_changes"`
	Pagination     Pagination                        `json:"pagination"`
}

func PendingUserStatusChangeSuccessResponse(pending *PendingUserStatusChange) Response {
	response := PendingUserStatusChangeResponse{
		ID:                    pending.ID,
		UserID:                pending.UserID,
		UserName:              pending.User.Name,
		UserPhone:             pending.User.Phone,
		MakerAdminID:          pending.MakerAdminID,
		MakerAdminName:        pending.MakerAdmin.Name,
		CurrentStatus:         pending.CurrentStatus,
		CurrentStatusString:   GetUserStatusString(pending.CurrentStatus),
		RequestedStatus:       pending.RequestedStatus,
		RequestedStatusString: GetUserStatusString(pending.RequestedStatus),
		Reason:                pending.Reason,
		Status:                pending.Status,
		Priority:              pending.Priority,
		ApprovalComments:      pending.ApprovalComments,
		RejectionReason:       pending.RejectionReason,
		ExpiresAt:             pending.ExpiresAt,
		ApprovedAt:            pending.ApprovedAt,
		RejectedAt:            pending.RejectedAt,
		CreatedAt:             pending.CreatedAt,
		UpdatedAt:             pending.UpdatedAt,
	}

	if pending.CheckerAdminID != nil && pending.CheckerAdmin != nil {
		response.CheckerAdminID = pending.CheckerAdminID
		checkerName := pending.CheckerAdmin.Name
		response.CheckerAdminName = &checkerName
	}

	return Response{
		Code:    200,
		Message: "Pending user status change retrieved successfully",
		Data:    response,
	}
}

func PendingUserStatusChangeListSuccessResponse(pending []PendingUserStatusChangeResponse, total int, page, perPage int) Response {
	totalPages := (total + perPage - 1) / perPage

	return Response{
		Code:    200,
		Message: "Pending user status changes retrieved successfully",
		Data: PendingUserStatusChangeListResponse{
			PendingChanges: pending,
			Pagination: Pagination{
				Page:       page,
				Limit:      perPage,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	}
}
