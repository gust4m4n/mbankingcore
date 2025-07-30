package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null;index"`
	Type          string         `json:"type" gorm:"not null"`              // "topup" atau "withdraw"
	Amount        int64          `json:"amount" gorm:"not null"`            // Amount dalam format int64
	BalanceBefore int64          `json:"balance_before" gorm:"not null"`    // Balance sebelum transaksi
	BalanceAfter  int64          `json:"balance_after" gorm:"not null"`     // Balance setelah transaksi
	Description   string         `json:"description"`                       // Deskripsi transaksi
	Status        string         `json:"status" gorm:"default:'completed'"` // "completed", "failed", "pending"
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
