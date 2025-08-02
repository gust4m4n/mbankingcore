package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CheckerMakerHandler struct {
	DB *gorm.DB
}

func NewCheckerMakerHandler(db *gorm.DB) *CheckerMakerHandler {
	return &CheckerMakerHandler{DB: db}
}

// CreatePendingTransaction - Create a new pending transaction for approval
func (h *CheckerMakerHandler) CreatePendingTransaction(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Admin authentication required",
		})
		return
	}

	var req models.PendingTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// Check if user exists
	var user models.User
	if err := h.DB.First(&user, req.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to validate user",
		})
		return
	}

	// Get approval threshold for this transaction type
	var threshold models.ApprovalThreshold
	if err := h.DB.Where("transaction_type = ? AND is_active = ?", req.TransactionType, true).
		First(&threshold).Error; err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("No approval threshold configured for transaction type: %s", req.TransactionType),
		})
		return
	}

	// Check if amount requires approval
	if req.Amount < threshold.AmountThreshold {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Transaction amount %d is below approval threshold %d", req.Amount, threshold.AmountThreshold),
		})
		return
	}

	// Calculate expected balance
	var expectedBalance int64
	switch req.TransactionType {
	case "topup", "balance_adjustment":
		expectedBalance = user.Balance + req.Amount
	case "withdraw":
		expectedBalance = user.Balance - req.Amount
		if expectedBalance < 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Insufficient balance for withdrawal",
			})
			return
		}
	case "balance_set":
		expectedBalance = req.Amount
	default:
		expectedBalance = user.Balance
	}

	// Set default priority if not provided
	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(threshold.AutoExpireHours) * time.Hour)

	// Create pending transaction
	pendingTxn := models.PendingTransaction{
		UserID:            req.UserID,
		MakerAdminID:      adminID.(uint),
		TransactionType:   req.TransactionType,
		Amount:            req.Amount,
		CurrentBalance:    user.Balance,
		ExpectedBalance:   expectedBalance,
		Description:       req.Description,
		Reason:            req.Reason,
		Status:            "pending",
		Priority:          priority,
		RequiresApproval:  true,
		ApprovalThreshold: threshold.AmountThreshold,
		RequestData:       req.RequestData,
		ExpiresAt:         &expiresAt,
	}

	if err := h.DB.Create(&pendingTxn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create pending transaction",
		})
		return
	}

	// Load relationships for response
	h.DB.Preload("User").Preload("MakerAdmin").First(&pendingTxn, pendingTxn.ID)

	// Create audit log
	auditDetails := map[string]interface{}{
		"pending_transaction_id": pendingTxn.ID,
		"user_id":                req.UserID,
		"transaction_type":       req.TransactionType,
		"amount":                 req.Amount,
		"priority":               priority,
	}
	auditDetailsJSON, _ := json.Marshal(auditDetails)
	auditDetailsRaw := json.RawMessage(auditDetailsJSON)

	auditLog := models.AuditLog{
		EntityType: "pending_transaction",
		EntityID:   pendingTxn.ID,
		Action:     "CREATE",
		AdminID:    &[]uint{adminID.(uint)}[0],
		IPAddress:  c.ClientIP(),
		NewValues:  &auditDetailsRaw,
	}

	if err := h.DB.Create(&auditLog).Error; err != nil {
		// Log error but continue (audit shouldn't break the main operation)
		fmt.Printf("Failed to create audit log: %v\n", err)
	}

	response := models.PendingTransactionResponse{
		ID:                pendingTxn.ID,
		UserID:            pendingTxn.UserID,
		UserName:          pendingTxn.User.Name,
		MakerAdminID:      pendingTxn.MakerAdminID,
		MakerAdminName:    pendingTxn.MakerAdmin.Name,
		TransactionType:   pendingTxn.TransactionType,
		Amount:            pendingTxn.Amount,
		CurrentBalance:    pendingTxn.CurrentBalance,
		ExpectedBalance:   pendingTxn.ExpectedBalance,
		Description:       pendingTxn.Description,
		Reason:            pendingTxn.Reason,
		Status:            pendingTxn.Status,
		Priority:          pendingTxn.Priority,
		ApprovalThreshold: pendingTxn.ApprovalThreshold,
		ExpiresAt:         pendingTxn.ExpiresAt,
		CreatedAt:         pendingTxn.CreatedAt,
	}

	// Calculate time to expire
	if pendingTxn.ExpiresAt != nil {
		timeToExpire := time.Until(*pendingTxn.ExpiresAt)
		response.DaysToExpire = int(timeToExpire.Hours() / 24)
		response.HoursToExpire = int(timeToExpire.Hours()) % 24
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Code:    http.StatusCreated,
		Message: "Pending transaction created successfully",
		Data:    response,
	})
}

// GetPendingTransactions - Get list of pending transactions for approval
func (h *CheckerMakerHandler) GetPendingTransactions(c *gin.Context) {
	// Parse query parameters
	page := 1
	limit := 20
	status := c.Query("status")
	transactionType := c.Query("transaction_type")
	priority := c.Query("priority")
	userID := c.Query("user_id")
	makerAdminID := c.Query("maker_admin_id")

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	// Build query
	query := h.DB.Model(&models.PendingTransaction{}).
		Preload("User").
		Preload("MakerAdmin").
		Preload("CheckerAdmin")

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if transactionType != "" {
		query = query.Where("transaction_type = ?", transactionType)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if userID != "" {
		if parsedUserID, err := strconv.ParseUint(userID, 10, 32); err == nil {
			query = query.Where("user_id = ?", uint(parsedUserID))
		}
	}
	if makerAdminID != "" {
		if parsedMakerID, err := strconv.ParseUint(makerAdminID, 10, 32); err == nil {
			query = query.Where("maker_admin_id = ?", uint(parsedMakerID))
		}
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get transactions
	var pendingTxns []models.PendingTransaction
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&pendingTxns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch pending transactions",
		})
		return
	}

	// Convert to response format
	var responses []models.PendingTransactionResponse
	for _, txn := range pendingTxns {
		response := models.PendingTransactionResponse{
			ID:                txn.ID,
			UserID:            txn.UserID,
			UserName:          txn.User.Name,
			MakerAdminID:      txn.MakerAdminID,
			MakerAdminName:    txn.MakerAdmin.Name,
			TransactionType:   txn.TransactionType,
			Amount:            txn.Amount,
			CurrentBalance:    txn.CurrentBalance,
			ExpectedBalance:   txn.ExpectedBalance,
			Description:       txn.Description,
			Reason:            txn.Reason,
			Status:            txn.Status,
			Priority:          txn.Priority,
			ApprovalThreshold: txn.ApprovalThreshold,
			ApprovalComments:  txn.ApprovalComments,
			RejectionReason:   txn.RejectionReason,
			ExpiresAt:         txn.ExpiresAt,
			CreatedAt:         txn.CreatedAt,
		}

		// Set checker admin info if exists
		if txn.CheckerAdmin != nil {
			checkerID := txn.CheckerAdmin.ID
			checkerName := txn.CheckerAdmin.Name
			response.CheckerAdminID = &checkerID
			response.CheckerAdminName = &checkerName
		}

		// Calculate time to expire
		if txn.ExpiresAt != nil {
			timeToExpire := time.Until(*txn.ExpiresAt)
			response.DaysToExpire = int(timeToExpire.Hours() / 24)
			response.HoursToExpire = int(timeToExpire.Hours()) % 24
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    http.StatusOK,
		Message: "Pending transactions retrieved successfully",
		Data: map[string]interface{}{
			"pending_transactions": responses,
			"pagination": map[string]interface{}{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// ApproveOrRejectTransaction - Approve or reject a pending transaction
func (h *CheckerMakerHandler) ApproveOrRejectTransaction(c *gin.Context) {
	checkerAdminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Admin authentication required",
		})
		return
	}

	pendingIDStr := c.Param("id")
	pendingID, err := strconv.ParseUint(pendingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid pending transaction ID",
		})
		return
	}

	var req models.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// Validate rejection reason if rejecting
	if req.Action == "reject" && req.RejectionReason == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Rejection reason is required when rejecting transaction",
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

	// Get pending transaction with lock
	var pendingTxn models.PendingTransaction
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("User").
		Preload("MakerAdmin").
		First(&pendingTxn, uint(pendingID)).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Pending transaction not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch pending transaction",
		})
		return
	}

	// Check if transaction is still pending
	if pendingTxn.Status != "pending" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Transaction is already %s", pendingTxn.Status),
		})
		return
	}

	// Check if transaction has expired
	if pendingTxn.ExpiresAt != nil && time.Now().After(*pendingTxn.ExpiresAt) {
		// Mark as expired
		tx.Model(&pendingTxn).Updates(map[string]interface{}{
			"status":     "expired",
			"updated_at": time.Now(),
		})
		tx.Commit()

		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Transaction has expired",
		})
		return
	}

	// Check that checker is not the same as maker (segregation of duties)
	if pendingTxn.MakerAdminID == checkerAdminID.(uint) {
		tx.Rollback()
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "Admin cannot approve their own transaction (segregation of duties)",
		})
		return
	}

	now := time.Now()
	checkerID := checkerAdminID.(uint)

	if req.Action == "approve" {
		// Approve and process the transaction
		if err := h.processApprovedTransaction(tx, &pendingTxn, checkerID, req.Comments); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("Failed to process approved transaction: %v", err),
			})
			return
		}

		// Update pending transaction status
		updates := map[string]interface{}{
			"status":            "approved",
			"checker_admin_id":  checkerID,
			"approval_comments": req.Comments,
			"approved_at":       &now,
			"processed_at":      &now,
			"updated_at":        now,
		}
		if err := tx.Model(&pendingTxn).Updates(updates).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update pending transaction",
			})
			return
		}

	} else if req.Action == "reject" {
		// Reject the transaction
		updates := map[string]interface{}{
			"status":            "rejected",
			"checker_admin_id":  checkerID,
			"rejection_reason":  req.RejectionReason,
			"approval_comments": req.Comments,
			"rejected_at":       &now,
			"updated_at":        now,
		}
		if err := tx.Model(&pendingTxn).Updates(updates).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update pending transaction",
			})
			return
		}
	}

	tx.Commit()

	// Create audit log
	auditDetails := map[string]interface{}{
		"pending_transaction_id": pendingTxn.ID,
		"action":                 req.Action,
		"comments":               req.Comments,
		"rejection_reason":       req.RejectionReason,
	}
	auditDetailsJSON, _ := json.Marshal(auditDetails)
	auditDetailsRaw := json.RawMessage(auditDetailsJSON)

	auditLog := models.AuditLog{
		EntityType: "pending_transaction",
		EntityID:   pendingTxn.ID,
		Action:     strings.ToUpper(req.Action),
		AdminID:    &checkerID,
		IPAddress:  c.ClientIP(),
		NewValues:  &auditDetailsRaw,
	}

	if err := h.DB.Create(&auditLog).Error; err != nil {
		// Log error but continue (audit shouldn't break the main operation)
		fmt.Printf("Failed to create audit log: %v\n", err)
	}

	// Get updated pending transaction for response
	h.DB.Preload("User").Preload("MakerAdmin").Preload("CheckerAdmin").First(&pendingTxn, pendingTxn.ID)

	// Get checker admin info
	var checkerAdmin models.Admin
	h.DB.First(&checkerAdmin, checkerID)

	response := models.PendingTransactionResponse{
		ID:                pendingTxn.ID,
		UserID:            pendingTxn.UserID,
		UserName:          pendingTxn.User.Name,
		MakerAdminID:      pendingTxn.MakerAdminID,
		MakerAdminName:    pendingTxn.MakerAdmin.Name,
		CheckerAdminID:    &checkerID,
		CheckerAdminName:  &checkerAdmin.Name,
		TransactionType:   pendingTxn.TransactionType,
		Amount:            pendingTxn.Amount,
		CurrentBalance:    pendingTxn.CurrentBalance,
		ExpectedBalance:   pendingTxn.ExpectedBalance,
		Description:       pendingTxn.Description,
		Reason:            pendingTxn.Reason,
		Status:            pendingTxn.Status,
		Priority:          pendingTxn.Priority,
		ApprovalThreshold: pendingTxn.ApprovalThreshold,
		ApprovalComments:  pendingTxn.ApprovalComments,
		RejectionReason:   pendingTxn.RejectionReason,
		ExpiresAt:         pendingTxn.ExpiresAt,
		CreatedAt:         pendingTxn.CreatedAt,
	}

	message := fmt.Sprintf("Transaction %s successfully", req.Action)
	if req.Action == "approve" {
		message += " and processed"
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    response,
	})
}

// processApprovedTransaction creates the actual transaction when approved
func (h *CheckerMakerHandler) processApprovedTransaction(tx *gorm.DB, pendingTxn *models.PendingTransaction, checkerAdminID uint, comments string) error {
	// Get fresh user data with lock
	var user models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&user, pendingTxn.UserID).Error; err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	// Verify current balance matches expected (in case it changed since request)
	if user.Balance != pendingTxn.CurrentBalance {
		// Update current balance in pending transaction for audit
		tx.Model(pendingTxn).Update("current_balance", user.Balance)

		// Recalculate expected balance
		var newExpectedBalance int64
		switch pendingTxn.TransactionType {
		case "topup", "balance_adjustment":
			newExpectedBalance = user.Balance + pendingTxn.Amount
		case "withdraw":
			newExpectedBalance = user.Balance - pendingTxn.Amount
			if newExpectedBalance < 0 {
				return fmt.Errorf("insufficient balance: current=%d, required=%d", user.Balance, pendingTxn.Amount)
			}
		case "balance_set":
			newExpectedBalance = pendingTxn.Amount
		default:
			newExpectedBalance = user.Balance
		}

		// Update expected balance
		tx.Model(pendingTxn).Update("expected_balance", newExpectedBalance)
		pendingTxn.ExpectedBalance = newExpectedBalance
	}

	// Create the actual transaction
	description := pendingTxn.Description
	if comments != "" {
		description = fmt.Sprintf("%s (Approved by admin: %s)", description, comments)
	}

	transaction := models.Transaction{
		UserID:        pendingTxn.UserID,
		Type:          pendingTxn.TransactionType,
		Amount:        pendingTxn.Amount,
		BalanceBefore: user.Balance,
		BalanceAfter:  pendingTxn.ExpectedBalance,
		Description:   description,
		Status:        "completed",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	// Update user balance
	if err := tx.Model(&user).Update("balance", pendingTxn.ExpectedBalance).Error; err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	// Link the final transaction to pending transaction
	if err := tx.Model(pendingTxn).Update("final_transaction_id", transaction.ID).Error; err != nil {
		return fmt.Errorf("failed to link final transaction: %v", err)
	}

	return nil
}

// GetApprovalStats - Get approval statistics for dashboard
func (h *CheckerMakerHandler) GetApprovalStats(c *gin.Context) {
	var stats models.ApprovalStats

	today := time.Now().Truncate(24 * time.Hour)

	// Pending count
	h.DB.Model(&models.PendingTransaction{}).
		Where("status = ?", "pending").
		Count(&[]int64{int64(stats.PendingCount)}[0])

	// Approved today
	h.DB.Model(&models.PendingTransaction{}).
		Where("status = ? AND approved_at >= ?", "approved", today).
		Count(&[]int64{int64(stats.ApprovedToday)}[0])

	// Rejected today
	h.DB.Model(&models.PendingTransaction{}).
		Where("status = ? AND rejected_at >= ?", "rejected", today).
		Count(&[]int64{int64(stats.RejectedToday)}[0])

	// Expired count
	h.DB.Model(&models.PendingTransaction{}).
		Where("status = ? OR (status = 'pending' AND expires_at < ?)", "expired", time.Now()).
		Count(&[]int64{int64(stats.ExpiredCount)}[0])

	// High priority count
	h.DB.Model(&models.PendingTransaction{}).
		Where("status = 'pending' AND priority = ?", "high").
		Count(&[]int64{int64(stats.HighPriorityCount)}[0])

	// Critical priority count
	h.DB.Model(&models.PendingTransaction{}).
		Where("status = 'pending' AND priority = ?", "critical").
		Count(&[]int64{int64(stats.CriticalPriorityCount)}[0])

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    http.StatusOK,
		Message: "Approval statistics retrieved successfully",
		Data:    stats,
	})
}
