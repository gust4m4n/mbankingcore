package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApprovalThresholdHandler struct {
	DB *gorm.DB
}

func NewApprovalThresholdHandler(db *gorm.DB) *ApprovalThresholdHandler {
	return &ApprovalThresholdHandler{DB: db}
}

// GetApprovalThresholds - Get all approval thresholds
func (h *ApprovalThresholdHandler) GetApprovalThresholds(c *gin.Context) {
	var thresholds []models.ApprovalThreshold
	if err := h.DB.Where("is_active = ?", true).
		Order("transaction_type ASC").
		Find(&thresholds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve approval thresholds",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    http.StatusOK,
		Message: "Approval thresholds retrieved successfully",
		Data:    thresholds,
	})
}

// GetApprovalThresholdByType - Get approval threshold by transaction type
func (h *ApprovalThresholdHandler) GetApprovalThresholdByType(c *gin.Context) {
	transactionType := c.Param("type")

	var threshold models.ApprovalThreshold
	if err := h.DB.Where("transaction_type = ? AND is_active = ?", transactionType, true).
		First(&threshold).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Approval threshold not found for this transaction type",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve approval threshold",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    http.StatusOK,
		Message: "Approval threshold retrieved successfully",
		Data:    threshold,
	})
}

// CreateOrUpdateApprovalThreshold - Create or update approval threshold
func (h *ApprovalThresholdHandler) CreateOrUpdateApprovalThreshold(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Admin authentication required",
		})
		return
	}

	var req models.ApprovalThresholdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// Validate that dual approval threshold is higher than regular threshold
	if req.RequiresDualApproval && req.DualApprovalThreshold <= req.AmountThreshold {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Dual approval threshold must be higher than regular threshold",
		})
		return
	}

	// Check if threshold already exists
	var existingThreshold models.ApprovalThreshold
	err := h.DB.Where("transaction_type = ?", req.TransactionType).First(&existingThreshold).Error

	if err == gorm.ErrRecordNotFound {
		// Create new threshold
		threshold := models.ApprovalThreshold{
			TransactionType:       req.TransactionType,
			AmountThreshold:       req.AmountThreshold,
			RequiresDualApproval:  req.RequiresDualApproval,
			DualApprovalThreshold: req.DualApprovalThreshold,
			AutoExpireHours:       req.AutoExpireHours,
			IsActive:              true,
		}

		if err := h.DB.Create(&threshold).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create approval threshold",
			})
			return
		}

		// Create audit log
		auditDetails := map[string]interface{}{
			"transaction_type":        req.TransactionType,
			"amount_threshold":        req.AmountThreshold,
			"requires_dual_approval":  req.RequiresDualApproval,
			"dual_approval_threshold": req.DualApprovalThreshold,
			"auto_expire_hours":       req.AutoExpireHours,
		}
		auditDetailsJSON, _ := json.Marshal(auditDetails)
		auditDetailsRaw := json.RawMessage(auditDetailsJSON)

		auditLog := models.AuditLog{
			EntityType: "approval_threshold",
			EntityID:   threshold.ID,
			Action:     "CREATE",
			AdminID:    &[]uint{adminID.(uint)}[0],
			IPAddress:  c.ClientIP(),
			NewValues:  &auditDetailsRaw,
		}
		h.DB.Create(&auditLog)

		c.JSON(http.StatusCreated, models.APIResponse{
			Code:    http.StatusCreated,
			Message: "Approval threshold created successfully",
			Data:    threshold,
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check existing threshold",
		})
		return
	}

	// Update existing threshold
	oldValues := map[string]interface{}{
		"amount_threshold":        existingThreshold.AmountThreshold,
		"requires_dual_approval":  existingThreshold.RequiresDualApproval,
		"dual_approval_threshold": existingThreshold.DualApprovalThreshold,
		"auto_expire_hours":       existingThreshold.AutoExpireHours,
		"is_active":               existingThreshold.IsActive,
	}

	existingThreshold.AmountThreshold = req.AmountThreshold
	existingThreshold.RequiresDualApproval = req.RequiresDualApproval
	existingThreshold.DualApprovalThreshold = req.DualApprovalThreshold
	existingThreshold.AutoExpireHours = req.AutoExpireHours
	existingThreshold.IsActive = true

	if err := h.DB.Save(&existingThreshold).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update approval threshold",
		})
		return
	}

	// Create audit log
	oldValuesJSON, _ := json.Marshal(oldValues)
	oldValuesRaw := json.RawMessage(oldValuesJSON)

	newValues := map[string]interface{}{
		"amount_threshold":        req.AmountThreshold,
		"requires_dual_approval":  req.RequiresDualApproval,
		"dual_approval_threshold": req.DualApprovalThreshold,
		"auto_expire_hours":       req.AutoExpireHours,
		"is_active":               true,
	}
	newValuesJSON, _ := json.Marshal(newValues)
	newValuesRaw := json.RawMessage(newValuesJSON)

	auditLog := models.AuditLog{
		EntityType: "approval_threshold",
		EntityID:   existingThreshold.ID,
		Action:     "UPDATE",
		AdminID:    &[]uint{adminID.(uint)}[0],
		IPAddress:  c.ClientIP(),
		OldValues:  &oldValuesRaw,
		NewValues:  &newValuesRaw,
	}
	h.DB.Create(&auditLog)

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    http.StatusOK,
		Message: "Approval threshold updated successfully",
		Data:    existingThreshold,
	})
}

// DeactivateApprovalThreshold - Deactivate approval threshold
func (h *ApprovalThresholdHandler) DeactivateApprovalThreshold(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Admin authentication required",
		})
		return
	}

	thresholdIDStr := c.Param("id")
	thresholdID, err := strconv.ParseUint(thresholdIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid threshold ID",
		})
		return
	}

	var threshold models.ApprovalThreshold
	if err := h.DB.First(&threshold, uint(thresholdID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Approval threshold not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve approval threshold",
		})
		return
	}

	// Update to inactive
	oldValues := map[string]interface{}{
		"is_active": threshold.IsActive,
	}

	threshold.IsActive = false
	if err := h.DB.Save(&threshold).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to deactivate approval threshold",
		})
		return
	}

	// Create audit log
	oldValuesJSON, _ := json.Marshal(oldValues)
	oldValuesRaw := json.RawMessage(oldValuesJSON)

	newValues := map[string]interface{}{
		"is_active": false,
	}
	newValuesJSON, _ := json.Marshal(newValues)
	newValuesRaw := json.RawMessage(newValuesJSON)

	auditLog := models.AuditLog{
		EntityType: "approval_threshold",
		EntityID:   threshold.ID,
		Action:     "DEACTIVATE",
		AdminID:    &[]uint{adminID.(uint)}[0],
		IPAddress:  c.ClientIP(),
		OldValues:  &oldValuesRaw,
		NewValues:  &newValuesRaw,
	}
	h.DB.Create(&auditLog)

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    http.StatusOK,
		Message: "Approval threshold deactivated successfully",
		Data:    threshold,
	})
}
