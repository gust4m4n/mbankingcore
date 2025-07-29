package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"mbankingcore/models"
	"mbankingcore/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB             *gorm.DB
	SessionManager *utils.SessionManager
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		DB:             db,
		SessionManager: utils.NewSessionManager(db),
	}
}

// BankingLogin handles first step of banking authentication - sends OTP
func (h *AuthHandler) BankingLogin(c *gin.Context) {
	var req models.BankingLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request data",
			"errors":  err.Error(),
		})
		return
	}

	// Additional custom validation
	if len(req.Name) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Name must be at least 8 characters long",
		})
		return
	}

	if len(req.AccountNumber) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Account number must be at least 8 characters long",
		})
		return
	}

	if len(req.Phone) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Phone number must be at least 8 characters long",
		})
		return
	}

	if len(req.MotherName) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Mother name must be at least 8 characters long",
		})
		return
	}

	if len(req.PinAtm) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "PIN ATM must be exactly 6 digits",
		})
		return
	}

	// Validate PIN is numeric
	for _, char := range req.PinAtm {
		if char < '0' || char > '9' {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "PIN ATM must contain only numeric digits",
			})
			return
		}
	}

	// Validate device info
	if req.DeviceInfo.DeviceType == "" || req.DeviceInfo.DeviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Device information is required",
		})
		return
	}

	// Check if phone number is already registered
	var existingUser models.User
	phoneExists := h.DB.Preload("BankAccounts").Where("phone = ?", req.Phone).First(&existingUser).Error == nil

	if phoneExists {
		// Phone is registered - validate user data matches
		if existingUser.Name != req.Name || existingUser.MotherName != req.MotherName {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "User information does not match our records",
			})
			return
		}

		// Check if account number exists for this user
		var bankAccount models.BankAccount
		accountExists := h.DB.Where("user_id = ? AND account_number = ?", existingUser.ID, req.AccountNumber).First(&bankAccount).Error == nil

		if !accountExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Account number not found for this user",
			})
			return
		}

		// Verify PIN (compare plain PIN with hashed PIN in database)
		if err := utils.CheckPassword(existingUser.PinAtm, req.PinAtm); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid PIN ATM",
			})
			return
		}
	} else {
		// Phone is not registered - will auto-register during OTP verification
		log.Printf("New phone number %s will be registered after OTP verification", req.Phone)
	}

	// Generate 6-digit OTP and unique login token
	otpCode := utils.GenerateOTP()
	loginToken := utils.GenerateLoginToken()

	// Create OTP session (expires in 5 minutes)
	// Store PIN in plain text temporarily for verification, will be hashed when saving to user
	otpSession := models.OTPSession{
		LoginToken:    loginToken,
		Phone:         req.Phone,
		OtpCode:       otpCode,
		Name:          req.Name,
		AccountNumber: req.AccountNumber,
		MotherName:    req.MotherName,
		PinAtm:        req.PinAtm, // Store plain PIN temporarily for verification
		DeviceType:    string(req.DeviceInfo.DeviceType),
		DeviceID:      req.DeviceInfo.DeviceID,
		DeviceName:    req.DeviceInfo.DeviceName,
		ExpiresAt:     time.Now().Add(5 * time.Minute),
		IsUsed:        false,
	}

	if err := h.DB.Create(&otpSession).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// TODO: Send OTP via SMS to req.Phone
	// For now, we'll log it (in production, integrate with SMS service)
	log.Printf("OTP for phone %s: %s", req.Phone, otpCode)

	var message string
	if phoneExists {
		message = "OTP sent to your registered phone number"
	} else {
		message = "Phone number will be registered. OTP sent for verification"
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: message,
		Data: gin.H{
			"login_token": loginToken,
			"expires_in":  300, // 5 minutes
			"is_new_user": !phoneExists,
		},
	})
}

// BankingLoginVerify handles second step of banking authentication - simplified for development
// Currently only validates login_token and always returns success
func (h *AuthHandler) BankingLoginVerify(c *gin.Context) {
	var req models.OTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Only validate login_token exists and is not used
	var otpSession models.OTPSession
	err := h.DB.Where("login_token = ? AND is_used = ?", req.LoginToken, false).
		First(&otpSession).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid login token or session not found",
		})
		return
	}

	// Mark OTP as used
	otpSession.IsUsed = true
	h.DB.Save(&otpSession)

	// Check if user exists by phone
	var user models.User
	userExists := h.DB.Where("phone = ?", otpSession.Phone).First(&user).Error == nil

	if !userExists {
		// Hash PIN for new user
		hashedPin, err := utils.HashPassword(otpSession.PinAtm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
			return
		}

		// Auto-register user
		user = models.User{
			Name:       otpSession.Name,
			Phone:      otpSession.Phone,
			MotherName: otpSession.MotherName,
			PinAtm:     hashedPin,
			Role:       models.ROLE_USER,
		}

		if err := h.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.CreateFailedResponse())
			return
		}

		// Create bank account for new user
		bankAccount := models.BankAccount{
			UserID:        user.ID,
			AccountNumber: otpSession.AccountNumber,
			AccountName:   user.Name,
			BankName:      "Unknown Bank",
			AccountType:   "saving",
			IsActive:      true,
			IsPrimary:     true,
		}

		if err := h.DB.Create(&bankAccount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.CreateFailedResponse())
			return
		}
	}

	// Create device session without validation checks (simplified for development)
	loginReq := models.MultiPlatformLoginRequest{
		Phone:      otpSession.Phone,
		Provider:   models.LoginProviderEmail,
		DeviceInfo: req.DeviceInfo,
	}

	ipAddress := utils.GetClientIP(c)
	session, err := h.SessionManager.CreateSession(user.ID, loginReq, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Remove sensitive data from response
	user.PinAtm = ""

	response := models.MultiPlatformLoginResponse{
		User:         user,
		AccessToken:  session.SessionToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    24 * 60 * 60,
		SessionID:    session.ID,
		DeviceInfo:   req.DeviceInfo,
	}

	message := "Login successful"
	if !userExists {
		message = "Account created and login successful"
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: message,
		Data:    response,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	session, newToken, err := h.SessionManager.RefreshSession(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid or expired refresh token",
		})
		return
	}

	response := gin.H{
		"access_token": newToken,
		"expires_in":   24 * 60 * 60, // 24 hours
		"session_id":   session.ID,
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Token refreshed successfully",
		Data:    response,
	})
}

// GetActiveSessions returns all active sessions for current user
func (h *AuthHandler) GetActiveSessions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	currentSessionToken := c.GetHeader("Authorization")
	if currentSessionToken != "" {
		currentSessionToken = currentSessionToken[7:] // Remove "Bearer " prefix
	}

	sessions, err := h.SessionManager.GetUserSessions(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Mark current session
	for i := range sessions {
		// You would need to implement session identification logic
		sessions[i].IsCurrent = false // Placeholder
	}

	response := models.UserSessionsResponse{
		Sessions: sessions,
		Total:    len(sessions),
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Active sessions retrieved successfully",
		Data:    response,
	})
}

// Logout handles session logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req models.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	var err error
	if req.AllDevices {
		// Logout from all devices
		err = h.SessionManager.LogoutAllSessions(userID.(uint))
	} else if req.SessionID != nil {
		// Logout specific session
		err = h.SessionManager.LogoutSession(*req.SessionID, userID.(uint))
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Either session_id or all_devices must be specified",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	message := "Logged out successfully"
	if req.AllDevices {
		message = "Logged out from all devices successfully"
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: message,
		Data:    nil,
	})
}

// LogoutOtherSessions logs out all other sessions except current
func (h *AuthHandler) LogoutOtherSessions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	// Get current session ID from token or header
	currentSessionIDStr := c.GetHeader("X-Session-ID")
	if currentSessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Current session ID is required",
		})
		return
	}

	currentSessionID, err := strconv.ParseUint(currentSessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid session ID",
		})
		return
	}

	err = h.SessionManager.LogoutAllOtherSessions(userID.(uint), uint(currentSessionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Logged out from other devices successfully",
		Data:    nil,
	})
}

// Get current user profile (protected endpoint)
func (h *AuthHandler) Profile(c *gin.Context) {
	// Get user ID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	// Find user
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.UserNotFoundResponse())
		return
	}

	// Remove sensitive data from response
	user.PinAtm = ""

	c.JSON(http.StatusOK, models.UserRetrievedResponse(user))
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Find user
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.UserNotFoundResponse())
		return
	}

	// Update user fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.UpdateFailedResponse())
		return
	}

	// Remove sensitive data from response
	user.PinAtm = ""

	c.JSON(http.StatusOK, models.UserUpdatedResponse(user))
}

// ChangePIN changes user PIN ATM and invalidates all sessions
func (h *AuthHandler) ChangePIN(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	var req struct {
		CurrentPIN string `json:"current_pin" binding:"required"`
		NewPIN     string `json:"new_pin" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Find user
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.UserNotFoundResponse())
		return
	}

	// Check current PIN
	if err := utils.CheckPassword(user.PinAtm, req.CurrentPIN); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid current PIN",
		})
		return
	}

	// Hash new PIN
	hashedPIN, err := utils.HashPassword(req.NewPIN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Update PIN
	user.PinAtm = hashedPIN
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.UpdateFailedResponse())
		return
	}

	// Invalidate all sessions for security
	h.SessionManager.LogoutAllSessions(userID.(uint))

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "PIN changed successfully. All sessions have been invalidated for security.",
		Data:    nil,
	})
}
