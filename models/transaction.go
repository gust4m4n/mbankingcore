package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	UserID         uint           `json:"user_id" gorm:"not null;index"`
	Type           string         `json:"type" gorm:"not null"`              // "topup", "withdraw", "transfer_out", "transfer_in", "reversal"
	Amount         int64          `json:"amount" gorm:"not null"`            // Amount dalam format int64
	BalanceBefore  int64          `json:"balance_before" gorm:"not null"`    // Balance sebelum transaksi
	BalanceAfter   int64          `json:"balance_after" gorm:"not null"`     // Balance setelah transaksi
	Description    string         `json:"description"`                       // Deskripsi transaksi
	Status         string         `json:"status" gorm:"default:'completed'"` // "completed", "failed", "pending"
	OriginalTxnID  *uint          `json:"original_txn_id,omitempty"`         // ID transaksi asli (untuk reversal)
	ReversedTxnID  *uint          `json:"reversed_txn_id,omitempty"`         // ID transaksi reversal (untuk transaksi yang di-reverse)
	IsReversed     bool           `json:"is_reversed" gorm:"default:false"`  // Apakah transaksi sudah di-reverse
	ReversalReason string         `json:"reversal_reason,omitempty"`         // Alasan reversal
	ReversedAt     *time.Time     `json:"reversed_at,omitempty"`             // Waktu reversal
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	User        User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OriginalTxn *Transaction `json:"original_txn,omitempty" gorm:"foreignKey:OriginalTxnID"`
	ReversedTxn *Transaction `json:"reversed_txn,omitempty" gorm:"foreignKey:ReversedTxnID"`
}

// Request structures for transaction operations
type TopupRequest struct {
	Amount      int64  `json:"amount" binding:"required,min=1"`
	Description string `json:"description"`
}

type WithdrawRequest struct {
	Amount      int64  `json:"amount" binding:"required,min=1"`
	Description string `json:"description"`
}

type TransferRequest struct {
	ToAccountNumber string `json:"to_account_number" binding:"required"`
	Amount          int64  `json:"amount" binding:"required,min=1"`
	Description     string `json:"description"`
}

type ReversalRequest struct {
	TransactionID  uint   `json:"transaction_id" binding:"required"`
	ReversalReason string `json:"reversal_reason" binding:"required,min=10,max=500"`
	AdminComments  string `json:"admin_comments,omitempty"`
}
