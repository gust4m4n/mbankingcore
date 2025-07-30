package handlers

import (
	"net/http"
	"strconv"

	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionHandler struct {
	DB *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
}

type TopupRequest struct {
	Amount int64 `json:"amount" binding:"required,min=1"`
}

type WithdrawRequest struct {
	Amount int64 `json:"amount" binding:"required,min=1"`
}

type TransferRequest struct {
	ToAccountNumber string `json:"to_account_number" binding:"required"`
	Amount          int64  `json:"amount" binding:"required,min=1"`
	Description     string `json:"description"`
}

// Topup - Add balance to user account
func (h *TransactionHandler) Topup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "User not authenticated",
		})
		return
	}

	var req TopupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// Start transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get user with lock
	var user models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&user, userID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "User not found",
		})
		return
	}

	// Record balance before transaction
	balanceBefore := user.Balance

	// Update user balance
	user.Balance += req.Amount
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update balance",
		})
		return
	}

	// Create transaction record
	transaction := models.Transaction{
		UserID:        user.ID,
		Type:          "topup",
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
		Status:        "completed",
		Description:   "Balance top-up",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create transaction record",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to commit transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Top-up successful",
		"data": gin.H{
			"transaction_id": transaction.ID,
			"amount":         req.Amount,
			"balance_before": balanceBefore,
			"balance_after":  user.Balance,
			"transaction_at": transaction.CreatedAt,
		},
	})
}

// Withdraw - Deduct balance from user account
func (h *TransactionHandler) Withdraw(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "User not authenticated",
		})
		return
	}

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// Start transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get user with lock
	var user models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&user, userID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "User not found",
		})
		return
	}

	// Check if user has sufficient balance
	if user.Balance < req.Amount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Insufficient balance",
		})
		return
	}

	// Record balance before transaction
	balanceBefore := user.Balance

	// Update user balance
	user.Balance -= req.Amount
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update balance",
		})
		return
	}

	// Create transaction record
	transaction := models.Transaction{
		UserID:        user.ID,
		Type:          "withdraw",
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
		Status:        "completed",
		Description:   "Balance withdrawal",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create transaction record",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to commit transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Withdrawal successful",
		"data": gin.H{
			"transaction_id": transaction.ID,
			"amount":         req.Amount,
			"balance_before": balanceBefore,
			"balance_after":  user.Balance,
			"transaction_at": transaction.CreatedAt,
		},
	})
}

// GetUserTransactions - Get transaction history for authenticated user
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "User not authenticated",
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var transactions []models.Transaction
	var total int64

	// Count total transactions
	h.DB.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&total)

	// Get transactions with pagination
	if err := h.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch transactions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Transactions retrieved successfully",
		"data": gin.H{
			"transactions": transactions,
			"pagination": gin.H{
				"current_page": page,
				"per_page":     limit,
				"total":        total,
				"total_pages":  (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetAllTransactions - Get all transactions for admin monitoring
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	userID := c.Query("user_id")
	transactionType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	query := h.DB.Model(&models.Transaction{}).Preload("User")

	// Apply filters
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}

	var transactions []models.Transaction
	var total int64

	// Count total transactions
	query.Count(&total)

	// Get transactions with pagination
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch transactions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All transactions retrieved successfully",
		"data": gin.H{
			"transactions": transactions,
			"pagination": gin.H{
				"current_page": page,
				"per_page":     limit,
				"total":        total,
				"total_pages":  (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// Transfer - Transfer balance between users using account number
func (h *TransactionHandler) Transfer(c *gin.Context) {
	senderUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "User not authenticated",
		})
		return
	}

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// Start transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get sender user with lock
	var senderUser models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&senderUser, senderUserID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Sender user not found",
		})
		return
	}

	// Find receiver user by account number through bank_accounts table
	var receiverBankAccount models.BankAccount
	if err := tx.Preload("User").Where("account_number = ? AND is_active = ?", req.ToAccountNumber, true).
		First(&receiverBankAccount).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Recipient account number not found or inactive",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to find recipient account",
			})
		}
		return
	}

	// Check if sender is trying to transfer to themselves
	if senderUser.ID == receiverBankAccount.User.ID {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Cannot transfer to your own account",
		})
		return
	}

	// Get receiver user with lock
	var receiverUser models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&receiverUser, receiverBankAccount.User.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Recipient user not found",
		})
		return
	}

	// Check if sender has sufficient balance
	if senderUser.Balance < req.Amount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Insufficient balance",
		})
		return
	}

	// Record balances before transaction
	senderBalanceBefore := senderUser.Balance
	receiverBalanceBefore := receiverUser.Balance

	// Update sender balance (subtract)
	senderUser.Balance -= req.Amount
	if err := tx.Save(&senderUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update sender balance",
		})
		return
	}

	// Update receiver balance (add)
	receiverUser.Balance += req.Amount
	if err := tx.Save(&receiverUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update receiver balance",
		})
		return
	}

	// Prepare description
	transferDesc := req.Description
	if transferDesc == "" {
		transferDesc = "Transfer to " + receiverBankAccount.AccountNumber
	}

	// Create transaction record for sender (debit)
	senderTransaction := models.Transaction{
		UserID:        senderUser.ID,
		Type:          "transfer_out",
		Amount:        req.Amount,
		BalanceBefore: senderBalanceBefore,
		BalanceAfter:  senderUser.Balance,
		Status:        "completed",
		Description:   transferDesc,
	}

	if err := tx.Create(&senderTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create sender transaction record",
		})
		return
	}

	// Create transaction record for receiver (credit)
	receiverTransaction := models.Transaction{
		UserID:        receiverUser.ID,
		Type:          "transfer_in",
		Amount:        req.Amount,
		BalanceBefore: receiverBalanceBefore,
		BalanceAfter:  receiverUser.Balance,
		Status:        "completed",
		Description:   "Transfer from " + senderUser.Phone,
	}

	if err := tx.Create(&receiverTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create receiver transaction record",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to commit transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Transfer successful",
		"data": gin.H{
			"transaction_id":        senderTransaction.ID,
			"to_account_number":     req.ToAccountNumber,
			"to_account_name":       receiverBankAccount.AccountName,
			"amount":                req.Amount,
			"sender_balance_before": senderBalanceBefore,
			"sender_balance_after":  senderUser.Balance,
			"description":           transferDesc,
			"transaction_at":        senderTransaction.CreatedAt,
		},
	})
}
