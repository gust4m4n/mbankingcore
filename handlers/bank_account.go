package handlers

import (
	"net/http"
	"strconv"

	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BankAccountHandler struct {
	DB *gorm.DB
}

func NewBankAccountHandler(db *gorm.DB) *BankAccountHandler {
	return &BankAccountHandler{DB: db}
}

// GetBankAccounts returns all bank accounts for the authenticated user
func (h *BankAccountHandler) GetBankAccounts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	var bankAccounts []models.BankAccount
	err := h.DB.Where("user_id = ? AND is_active = ?", userID, true).
		Order("is_primary DESC, created_at ASC").
		Find(&bankAccounts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	var responses []models.BankAccountResponse
	for _, account := range bankAccounts {
		responses = append(responses, models.BankAccountResponse{
			ID:            account.ID,
			AccountNumber: account.AccountNumber,
			AccountName:   account.AccountName,
			BankName:      account.BankName,
			BankCode:      account.BankCode,
			AccountType:   account.AccountType,
			IsActive:      account.IsActive,
			IsPrimary:     account.IsPrimary,
			CreatedAt:     account.CreatedAt,
			UpdatedAt:     account.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Bank accounts retrieved successfully",
		Data: gin.H{
			"accounts": responses,
			"total":    len(responses),
		},
	})
}

// CreateBankAccount creates a new bank account for the authenticated user
func (h *BankAccountHandler) CreateBankAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	var req models.BankAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Check if account number already exists for this user
	var existingAccount models.BankAccount
	err := h.DB.Where("user_id = ? AND account_number = ?", userID, req.AccountNumber).First(&existingAccount).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "Account number already exists for this user",
		})
		return
	}

	// If this is set as primary, make other accounts non-primary
	if req.IsPrimary {
		h.DB.Model(&models.BankAccount{}).Where("user_id = ?", userID).Update("is_primary", false)
	}

	// Create new bank account
	bankAccount := models.BankAccount{
		UserID:        userID.(uint),
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
		BankName:      req.BankName,
		BankCode:      req.BankCode,
		AccountType:   req.AccountType,
		IsActive:      true,
		IsPrimary:     req.IsPrimary,
	}

	if err := h.DB.Create(&bankAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.CreateFailedResponse())
		return
	}

	response := models.BankAccountResponse{
		ID:            bankAccount.ID,
		AccountNumber: bankAccount.AccountNumber,
		AccountName:   bankAccount.AccountName,
		BankName:      bankAccount.BankName,
		BankCode:      bankAccount.BankCode,
		AccountType:   bankAccount.AccountType,
		IsActive:      bankAccount.IsActive,
		IsPrimary:     bankAccount.IsPrimary,
		CreatedAt:     bankAccount.CreatedAt,
		UpdatedAt:     bankAccount.UpdatedAt,
	}

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Bank account created successfully",
		Data:    response,
	})
}

// UpdateBankAccount updates an existing bank account
func (h *BankAccountHandler) UpdateBankAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	accountID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid account ID",
		})
		return
	}

	var req models.BankAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Find the bank account
	var bankAccount models.BankAccount
	err = h.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&bankAccount).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Bank account not found",
		})
		return
	}

	// If this is set as primary, make other accounts non-primary
	if req.IsPrimary && !bankAccount.IsPrimary {
		h.DB.Model(&models.BankAccount{}).Where("user_id = ? AND id != ?", userID, accountID).Update("is_primary", false)
	}

	// Update bank account
	bankAccount.AccountName = req.AccountName
	bankAccount.BankName = req.BankName
	bankAccount.BankCode = req.BankCode
	bankAccount.AccountType = req.AccountType
	bankAccount.IsPrimary = req.IsPrimary

	if err := h.DB.Save(&bankAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.UpdateFailedResponse())
		return
	}

	response := models.BankAccountResponse{
		ID:            bankAccount.ID,
		AccountNumber: bankAccount.AccountNumber,
		AccountName:   bankAccount.AccountName,
		BankName:      bankAccount.BankName,
		BankCode:      bankAccount.BankCode,
		AccountType:   bankAccount.AccountType,
		IsActive:      bankAccount.IsActive,
		IsPrimary:     bankAccount.IsPrimary,
		CreatedAt:     bankAccount.CreatedAt,
		UpdatedAt:     bankAccount.UpdatedAt,
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Bank account updated successfully",
		Data:    response,
	})
}

// DeleteBankAccount soft deletes a bank account
func (h *BankAccountHandler) DeleteBankAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	accountID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid account ID",
		})
		return
	}

	// Find the bank account
	var bankAccount models.BankAccount
	err = h.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&bankAccount).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Bank account not found",
		})
		return
	}

	// Check if this is the only active account
	var activeCount int64
	h.DB.Model(&models.BankAccount{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&activeCount)

	if activeCount <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Cannot delete the last active bank account",
		})
		return
	}

	// Soft delete the account
	bankAccount.IsActive = false
	if err := h.DB.Save(&bankAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.DeleteFailedResponse())
		return
	}

	// If this was the primary account, make another account primary
	if bankAccount.IsPrimary {
		var newPrimary models.BankAccount
		err = h.DB.Where("user_id = ? AND is_active = ? AND id != ?", userID, true, accountID).
			First(&newPrimary).Error
		if err == nil {
			newPrimary.IsPrimary = true
			h.DB.Save(&newPrimary)
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Bank account deleted successfully",
		Data:    nil,
	})
}

// SetPrimaryAccount sets a bank account as primary
func (h *BankAccountHandler) SetPrimaryAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	accountID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid account ID",
		})
		return
	}

	// Find the bank account
	var bankAccount models.BankAccount
	err = h.DB.Where("id = ? AND user_id = ? AND is_active = ?", accountID, userID, true).First(&bankAccount).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Bank account not found",
		})
		return
	}

	// Make other accounts non-primary
	h.DB.Model(&models.BankAccount{}).Where("user_id = ?", userID).Update("is_primary", false)

	// Set this account as primary
	bankAccount.IsPrimary = true
	if err := h.DB.Save(&bankAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.UpdateFailedResponse())
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Primary bank account updated successfully",
		Data:    nil,
	})
}
