package handlers

import (
	"net/http"
	"strconv"

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

// Login handles login from various platforms and devices (renamed from MultiPlatformLogin)
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.MultiPlatformLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Validate required fields based on provider
	if req.Provider == models.LoginProviderEmail {
		if req.Email == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, models.ValidationFailedResponse())
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

	var user models.User
	var err error

	switch req.Provider {
	case models.LoginProviderEmail:
		user, err = h.authenticateEmail(req.Email, req.Password)
	case models.LoginProviderGoogle:
		user, err = h.authenticateGoogle(req.ProviderID, req.Email)
	case models.LoginProviderApple:
		user, err = h.authenticateApple(req.ProviderID, req.Email)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Unsupported login provider",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.InvalidPasswordResponse())
		return
	}

	// Check if there's already an active session with the same device info
	// Must check device_type, device_id, AND device_name to prevent duplicate sessions
	var existingSession models.DeviceSession
	err = h.DB.Where("user_id = ? AND device_type = ? AND device_id = ? AND device_name = ? AND is_active = ?",
		user.ID, req.DeviceInfo.DeviceType, req.DeviceInfo.DeviceID, req.DeviceInfo.DeviceName, true).
		First(&existingSession).Error

	if err == nil {
		// Session with same device info already exists
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "Device is already logged in. Please logout from this device first or use a different device.",
			"data": gin.H{
				"existing_session": gin.H{
					"device_type":   existingSession.DeviceType,
					"device_id":     existingSession.DeviceID,
					"device_name":   existingSession.DeviceName,
					"last_activity": existingSession.LastActivity,
				},
			},
		})
		return
	}

	// Get client IP
	ipAddress := utils.GetClientIP(c)

	// Create new device session
	session, err := h.SessionManager.CreateSession(user.ID, req, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Remove password from response
	user.Password = ""

	response := models.MultiPlatformLoginResponse{
		User:         user,
		AccessToken:  session.SessionToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    24 * 60 * 60, // 24 hours in seconds
		SessionID:    session.ID,
		DeviceInfo:   req.DeviceInfo,
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Login successful",
		Data:    response,
	})
}

// Register handles registration with platform info (renamed from RegisterMultiPlatform)
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Name       string               `json:"name" binding:"required"`
		Email      string               `json:"email" binding:"required,email"`
		Phone      string               `json:"phone"`
		Password   string               `json:"password" binding:"required,min=6"`
		Role       string               `json:"role"`
		Provider   models.LoginProvider `json:"provider"`
		ProviderID string               `json:"provider_id,omitempty"`
		DeviceInfo models.DeviceInfo    `json:"device_info"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.InvalidRequestResponse())
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, models.EmailExistsResponse())
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Validate and set role
	role := req.Role
	if !models.ValidateRole(role) {
		role = models.ROLE_USER // Default to user if invalid role provided
	}

	// Create user with platform info
	user := models.User{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Password:      hashedPassword,
		Role:          role,
		EmailVerified: false,
	}

	// Set provider-specific fields
	switch req.Provider {
	case models.LoginProviderGoogle:
		user.GoogleID = req.ProviderID
	case models.LoginProviderApple:
		user.AppleID = req.ProviderID
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.CreateFailedResponse())
		return
	}

	// Create MultiPlatformLoginRequest for session creation
	loginReq := models.MultiPlatformLoginRequest{
		Email:      req.Email,
		Provider:   req.Provider,
		ProviderID: req.ProviderID,
		DeviceInfo: req.DeviceInfo,
	}

	// Create initial device session
	ipAddress := utils.GetClientIP(c)
	session, err := h.SessionManager.CreateSession(user.ID, loginReq, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Remove password from response
	user.Password = ""

	response := models.MultiPlatformLoginResponse{
		User:         user,
		AccessToken:  session.SessionToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    24 * 60 * 60,
		SessionID:    session.ID,
		DeviceInfo:   req.DeviceInfo,
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "User registered successfully",
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

	// Remove password from response
	user.Password = ""

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

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, models.UserUpdatedResponse(user))
}

// ChangePassword changes user password and invalidates all sessions
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse())
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
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

	// Check current password
	if err := utils.CheckPassword(user.Password, req.CurrentPassword); err != nil {
		c.JSON(http.StatusUnauthorized, models.InvalidPasswordResponse())
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerResponse())
		return
	}

	// Update password
	user.Password = hashedPassword
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.UpdateFailedResponse())
		return
	}

	// Invalidate all sessions for security
	h.SessionManager.LogoutAllSessions(userID.(uint))

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Password changed successfully. All sessions have been invalidated for security.",
		Data:    nil,
	})
}

// authenticateEmail handles email/password authentication
func (h *AuthHandler) authenticateEmail(email, password string) (models.User, error) {
	var user models.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	// Check password using existing utils
	if err := utils.CheckPassword(user.Password, password); err != nil {
		return user, err
	}

	return user, nil
}

// authenticateGoogle handles Google SSO authentication
func (h *AuthHandler) authenticateGoogle(googleID, email string) (models.User, error) {
	var user models.User

	// First try to find user by Google ID
	err := h.DB.Where("google_id = ?", googleID).First(&user).Error
	if err == nil {
		return user, nil
	}

	// If not found by Google ID, try by email
	err = h.DB.Where("email = ?", email).First(&user).Error
	if err == nil {
		// Link Google ID to existing user
		user.GoogleID = googleID
		h.DB.Save(&user)
		return user, nil
	}

	// Return error if user not found
	return user, gorm.ErrRecordNotFound
}

// authenticateApple handles Apple SSO authentication
func (h *AuthHandler) authenticateApple(appleID, email string) (models.User, error) {
	var user models.User

	// First try to find user by Apple ID
	err := h.DB.Where("apple_id = ?", appleID).First(&user).Error
	if err == nil {
		return user, nil
	}

	// If not found by Apple ID, try by email
	err = h.DB.Where("email = ?", email).First(&user).Error
	if err == nil {
		// Link Apple ID to existing user
		user.AppleID = appleID
		h.DB.Save(&user)
		return user, nil
	}

	// Return error if user not found
	return user, gorm.ErrRecordNotFound
}
