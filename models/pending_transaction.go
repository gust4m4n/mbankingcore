package models

import (
	"time"

	"gorm.io/gorm"
)

// PendingTransaction represents transactions that require approval (checker-maker system)
type PendingTransaction struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	UserID             uint           `json:"user_id" gorm:"not null;index"`           // Target user for the transaction
	MakerAdminID       uint           `json:"maker_admin_id" gorm:"not null;index"`    // Admin who initiated (maker)
	CheckerAdminID     *uint          `json:"checker_admin_id,omitempty" gorm:"index"` // Admin who approved/rejected (checker)
	TransactionType    string         `json:"transaction_type" gorm:"not null"`        // "topup", "withdraw", "transfer", "balance_adjustment", "balance_set"
	Amount             int64          `json:"amount" gorm:"not null"`                  // Transaction amount
	CurrentBalance     int64          `json:"current_balance" gorm:"not null"`         // User's current balance when request was made
	ExpectedBalance    int64          `json:"expected_balance" gorm:"not null"`        // Expected balance after transaction
	Description        string         `json:"description"`                             // Transaction description
	Reason             string         `json:"reason"`                                  // Reason for the transaction (for adjustments)
	Status             string         `json:"status" gorm:"default:'pending'"`         // "pending", "approved", "rejected", "expired"
	Priority           string         `json:"priority" gorm:"default:'normal'"`        // "low", "normal", "high", "critical"
	RequiresApproval   bool           `json:"requires_approval" gorm:"default:true"`   // Whether this transaction requires approval
	ApprovalThreshold  int64          `json:"approval_threshold" gorm:"not null"`      // Threshold amount that triggered approval requirement
	RequestData        string         `json:"request_data" gorm:"type:text"`           // JSON of original request data
	ApprovalComments   string         `json:"approval_comments"`                       // Comments from checker
	RejectionReason    string         `json:"rejection_reason"`                        // Reason for rejection
	ExpiresAt          *time.Time     `json:"expires_at"`                              // When the pending transaction expires
	ApprovedAt         *time.Time     `json:"approved_at"`                             // When approved
	RejectedAt         *time.Time     `json:"rejected_at"`                             // When rejected
	ProcessedAt        *time.Time     `json:"processed_at"`                            // When actually processed (transaction created)
	FinalTransactionID *uint          `json:"final_transaction_id,omitempty"`          // ID of the final created transaction
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User             User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	MakerAdmin       Admin        `json:"maker_admin,omitempty" gorm:"foreignKey:MakerAdminID"`
	CheckerAdmin     *Admin       `json:"checker_admin,omitempty" gorm:"foreignKey:CheckerAdminID"`
	FinalTransaction *Transaction `json:"final_transaction,omitempty" gorm:"foreignKey:FinalTransactionID"`
}

// ApprovalThreshold represents the approval requirements for different transaction types
type ApprovalThreshold struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	TransactionType       string         `json:"transaction_type" gorm:"not null;uniqueIndex"` // "topup", "withdraw", "transfer", "balance_adjustment", "balance_set"
	AmountThreshold       int64          `json:"amount_threshold" gorm:"not null"`             // Amount above which approval is required
	RequiresDualApproval  bool           `json:"requires_dual_approval" gorm:"default:false"`  // Whether it requires two approvers for very high amounts
	DualApprovalThreshold int64          `json:"dual_approval_threshold"`                      // Amount above which dual approval is required
	AutoExpireHours       int            `json:"auto_expire_hours" gorm:"default:24"`          // Hours after which pending transaction expires
	IsActive              bool           `json:"is_active" gorm:"default:true"`                // Whether this threshold is active
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request structures for checker-maker operations
type PendingTransactionRequest struct {
	UserID          uint   `json:"user_id" binding:"required"`
	TransactionType string `json:"transaction_type" binding:"required,oneof=topup withdraw transfer balance_adjustment balance_set"`
	Amount          int64  `json:"amount" binding:"required"`
	Description     string `json:"description"`
	Reason          string `json:"reason"` // Required for balance adjustments
	Priority        string `json:"priority" binding:"omitempty,oneof=low normal high critical"`
	RequestData     string `json:"request_data"` // JSON of original request
}

type ApprovalRequest struct {
	Action          string `json:"action" binding:"required,oneof=approve reject"`
	Comments        string `json:"comments"`
	RejectionReason string `json:"rejection_reason"` // Required if action is reject
}

type ApprovalThresholdRequest struct {
	TransactionType       string `json:"transaction_type" binding:"required,oneof=topup withdraw transfer balance_adjustment balance_set"`
	AmountThreshold       int64  `json:"amount_threshold" binding:"required,min=0"`
	RequiresDualApproval  bool   `json:"requires_dual_approval"`
	DualApprovalThreshold int64  `json:"dual_approval_threshold" binding:"min=0"`
	AutoExpireHours       int    `json:"auto_expire_hours" binding:"required,min=1,max=168"` // 1 hour to 1 week
}

// Response structures
type PendingTransactionResponse struct {
	ID                uint       `json:"id"`
	UserID            uint       `json:"user_id"`
	UserName          string     `json:"user_name"`
	MakerAdminID      uint       `json:"maker_admin_id"`
	MakerAdminName    string     `json:"maker_admin_name"`
	CheckerAdminID    *uint      `json:"checker_admin_id,omitempty"`
	CheckerAdminName  *string    `json:"checker_admin_name,omitempty"`
	TransactionType   string     `json:"transaction_type"`
	Amount            int64      `json:"amount"`
	CurrentBalance    int64      `json:"current_balance"`
	ExpectedBalance   int64      `json:"expected_balance"`
	Description       string     `json:"description"`
	Reason            string     `json:"reason,omitempty"`
	Status            string     `json:"status"`
	Priority          string     `json:"priority"`
	ApprovalThreshold int64      `json:"approval_threshold"`
	ApprovalComments  string     `json:"approval_comments,omitempty"`
	RejectionReason   string     `json:"rejection_reason,omitempty"`
	ExpiresAt         *time.Time `json:"expires_at"`
	CreatedAt         time.Time  `json:"created_at"`
	DaysToExpire      int        `json:"days_to_expire"`
	HoursToExpire     int        `json:"hours_to_expire"`
}

// ApprovalStats represents approval statistics for dashboard
type ApprovalStats struct {
	PendingCount          int `json:"pending_count"`
	ApprovedToday         int `json:"approved_today"`
	RejectedToday         int `json:"rejected_today"`
	ExpiredCount          int `json:"expired_count"`
	HighPriorityCount     int `json:"high_priority_count"`
	CriticalPriorityCount int `json:"critical_priority_count"`
}
