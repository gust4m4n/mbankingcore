package handlers

import (
	"net/http"
	"strconv"
	"time"

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

	var req models.TopupRequest
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

	var req models.WithdrawRequest
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
	if err := h.DB.Preload("User").Where("user_id = ?", userID).
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

	var req models.TransferRequest
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

// Reversal - Reverse a completed transaction (Admin only)
func (h *TransactionHandler) Reversal(c *gin.Context) {
	var req models.ReversalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request data",
		})
		return
	}

	// Start database transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lock and get original transaction
	var originalTxn models.Transaction
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND deleted_at IS NULL", req.TransactionID).
		First(&originalTxn).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Transaction not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get transaction",
			})
		}
		return
	}

	// Check if transaction is already reversed
	if originalTxn.IsReversed {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Transaction has already been reversed",
		})
		return
	}

	// Check if transaction can be reversed (only completed transactions)
	if originalTxn.Status != "completed" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Only completed transactions can be reversed",
		})
		return
	}

	// Check if transaction is not too old (business rule: can only reverse within 30 days)
	// You can adjust this timeframe based on business requirements
	// For now, we'll allow all completed transactions to be reversed

	// Lock and get user for balance update
	var user models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", originalTxn.UserID).
		First(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user",
		})
		return
	}

	// Calculate reversal operation
	var reversalType string
	var newBalance int64
	var reversalAmount int64

	currentBalance := user.Balance

	switch originalTxn.Type {
	case "topup":
		// Reverse topup: deduct the amount
		reversalType = "reversal"
		reversalAmount = originalTxn.Amount
		newBalance = currentBalance - originalTxn.Amount
		if newBalance < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Insufficient balance for reversal. Cannot reverse topup transaction.",
			})
			return
		}
	case "withdraw":
		// Reverse withdraw: add the amount back
		reversalType = "reversal"
		reversalAmount = originalTxn.Amount
		newBalance = currentBalance + originalTxn.Amount
	case "transfer_out":
		// Reverse outgoing transfer: add the amount back to sender
		reversalType = "reversal"
		reversalAmount = originalTxn.Amount
		newBalance = currentBalance + originalTxn.Amount

		// For transfer reversals, we also need to reverse the corresponding transfer_in
		// Find the corresponding transfer_in transaction
		var correspondingTxn models.Transaction
		if err := tx.Where("user_id != ? AND type = 'transfer_in' AND amount = ? AND created_at BETWEEN ? AND ?",
			originalTxn.UserID, originalTxn.Amount,
			originalTxn.CreatedAt.Add(-time.Minute), originalTxn.CreatedAt.Add(time.Minute)).
			First(&correspondingTxn).Error; err == nil {

			// Lock and get receiver user
			var receiverUser models.User
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", correspondingTxn.UserID).
				First(&receiverUser).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Failed to get receiver user for transfer reversal",
				})
				return
			}

			// Check if receiver has sufficient balance for reversal
			if receiverUser.Balance < correspondingTxn.Amount {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Receiver has insufficient balance for transfer reversal",
				})
				return
			}

			// Create reversal transaction for receiver (deduct the amount)
			receiverReversalBalance := receiverUser.Balance - correspondingTxn.Amount
			receiverReversalTxn := models.Transaction{
				UserID:         correspondingTxn.UserID,
				Type:           "reversal",
				Amount:         correspondingTxn.Amount,
				BalanceBefore:  receiverUser.Balance,
				BalanceAfter:   receiverReversalBalance,
				Description:    "Transfer reversal - " + req.ReversalReason,
				Status:         "completed",
				OriginalTxnID:  &correspondingTxn.ID,
				ReversalReason: req.ReversalReason,
			}

			if err := tx.Create(&receiverReversalTxn).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Failed to create receiver reversal transaction",
				})
				return
			}

			// Update receiver balance
			if err := tx.Model(&receiverUser).Update("balance", receiverReversalBalance).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Failed to update receiver balance",
				})
				return
			}

			// Mark corresponding transaction as reversed
			now := time.Now()
			if err := tx.Model(&correspondingTxn).Updates(map[string]interface{}{
				"is_reversed":     true,
				"reversed_txn_id": receiverReversalTxn.ID,
				"reversed_at":     &now,
			}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Failed to mark corresponding transaction as reversed",
				})
				return
			}
		}
	case "transfer_in":
		// Reverse incoming transfer: deduct the amount
		reversalType = "reversal"
		reversalAmount = originalTxn.Amount
		newBalance = currentBalance - originalTxn.Amount
		if newBalance < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Insufficient balance for transfer_in reversal",
			})
			return
		}
	default:
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Transaction type cannot be reversed",
		})
		return
	}

	// Create reversal transaction
	reversalTxn := models.Transaction{
		UserID:         originalTxn.UserID,
		Type:           reversalType,
		Amount:         reversalAmount,
		BalanceBefore:  currentBalance,
		BalanceAfter:   newBalance,
		Description:    "Reversal of transaction #" + strconv.Itoa(int(originalTxn.ID)) + " - " + req.ReversalReason,
		Status:         "completed",
		OriginalTxnID:  &originalTxn.ID,
		ReversalReason: req.ReversalReason,
	}

	if err := tx.Create(&reversalTxn).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create reversal transaction",
		})
		return
	}

	// Update user balance
	if err := tx.Model(&user).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user balance",
		})
		return
	}

	// Mark original transaction as reversed
	now := time.Now()
	if err := tx.Model(&originalTxn).Updates(map[string]interface{}{
		"is_reversed":     true,
		"reversed_txn_id": reversalTxn.ID,
		"reversed_at":     &now,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to mark original transaction as reversed",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to commit reversal transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Transaction reversed successfully",
		"data": gin.H{
			"reversal_transaction_id": reversalTxn.ID,
			"original_transaction_id": originalTxn.ID,
			"reversed_amount":         reversalAmount,
			"balance_before":          currentBalance,
			"balance_after":           newBalance,
			"reversal_reason":         req.ReversalReason,
			"reversed_at":             now,
		},
	})
}

// GetTransactionByID - Get transaction detail by ID
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	// Get transaction ID from URL parameter
	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Transaction ID is required",
		})
		return
	}

	// Convert to uint
	id, err := strconv.ParseUint(transactionID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid transaction ID format",
		})
		return
	}

	var transaction models.Transaction

	// Get transaction with user information
	if err := h.DB.Preload("User").
		Preload("OriginalTxn").
		Preload("ReversedTxn").
		Where("id = ?", uint(id)).
		First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Transaction not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Transaction retrieved successfully",
		"data": gin.H{
			"transaction": transaction,
		},
	})
}
