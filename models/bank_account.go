package models

import (
	"time"
)

// BankAccount represents a bank account belonging to a user
type BankAccount struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null;index"`
	User          User      `json:"user" gorm:"foreignKey:UserID"`
	AccountNumber string    `json:"account_number" gorm:"not null;uniqueIndex:idx_user_account;size:50"`
	AccountName   string    `json:"account_name" gorm:"not null;size:100"` // Name as it appears on the account
	BankName      string    `json:"bank_name" gorm:"size:100"`             // Bank institution name
	BankCode      string    `json:"bank_code" gorm:"size:10"`              // Bank code (e.g., "014" for BCA)
	AccountType   string    `json:"account_type" gorm:"size:20"`           // e.g., "saving", "checking", "current"
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	IsPrimary     bool      `json:"is_primary" gorm:"default:false"` // Primary account for the user
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BankAccountRequest for creating/updating bank account
type BankAccountRequest struct {
	AccountNumber string `json:"account_number" binding:"required,min=8,max=20"`
	AccountName   string `json:"account_name" binding:"required,min=3,max=100"`
	BankName      string `json:"bank_name" binding:"omitempty,max=100"`
	BankCode      string `json:"bank_code" binding:"omitempty,max=10"`
	AccountType   string `json:"account_type" binding:"omitempty,max=20"`
	IsPrimary     bool   `json:"is_primary"`
}

// BankAccountResponse for API responses
type BankAccountResponse struct {
	ID            uint      `json:"id"`
	AccountNumber string    `json:"account_number"`
	AccountName   string    `json:"account_name"`
	BankName      string    `json:"bank_name"`
	BankCode      string    `json:"bank_code"`
	AccountType   string    `json:"account_type"`
	IsActive      bool      `json:"is_active"`
	IsPrimary     bool      `json:"is_primary"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
