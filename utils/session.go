package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"mbankingcore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SessionManager handles device session management
type SessionManager struct {
	DB *gorm.DB
}

// NewSessionManager creates a new session manager
func NewSessionManager(db *gorm.DB) *SessionManager {
	return &SessionManager{DB: db}
}

// GenerateTokens generates access and refresh tokens
func (sm *SessionManager) GenerateTokens() (accessToken, refreshToken string, err error) {
	// Generate random tokens
	accessBytes := make([]byte, 32)
	refreshBytes := make([]byte, 32)

	if _, err := rand.Read(accessBytes); err != nil {
		return "", "", err
	}

	if _, err := rand.Read(refreshBytes); err != nil {
		return "", "", err
	}

	accessToken = hex.EncodeToString(accessBytes)
	refreshToken = hex.EncodeToString(refreshBytes)

	return accessToken, refreshToken, nil
}

// CreateSession creates a new device session
func (sm *SessionManager) CreateSession(userID uint, req models.MultiPlatformLoginRequest, ipAddress string) (*models.DeviceSession, error) {
	_, refreshToken, err := sm.GenerateTokens()
	if err != nil {
		return nil, err
	}

	// Get user for JWT creation
	var user models.User
	if err := sm.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Create JWT token without role
	jwtToken, err := GenerateJWT(userID, req.Phone)
	if err != nil {
		return nil, err
	}

	session := &models.DeviceSession{
		UserID:       userID,
		SessionToken: jwtToken, // Keep JWT for compatibility
		RefreshToken: refreshToken,
		DeviceType:   req.DeviceInfo.DeviceType,
		DeviceID:     req.DeviceInfo.DeviceID,
		DeviceName:   req.DeviceInfo.DeviceName,
		Provider:     req.Provider,
		ProviderID:   req.ProviderID,
		IPAddress:    ipAddress,
		IsActive:     true,
		LastActivity: time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour), // 24 hour expiry
	}

	if err := sm.DB.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

// ValidateSession validates a session token
func (sm *SessionManager) ValidateSession(sessionToken string) (*models.DeviceSession, error) {
	var session models.DeviceSession
	err := sm.DB.Where("session_token = ? AND is_active = ? AND expires_at > ?",
		sessionToken, true, time.Now()).
		Preload("User").
		First(&session).Error

	if err != nil {
		return nil, err
	}

	// Update last activity
	session.LastActivity = time.Now()
	sm.DB.Save(&session)

	return &session, nil
}

// RefreshSession refreshes a session using refresh token
func (sm *SessionManager) RefreshSession(refreshToken string) (*models.DeviceSession, string, error) {
	var session models.DeviceSession
	err := sm.DB.Where("refresh_token = ? AND is_active = ? AND expires_at > ?",
		refreshToken, true, time.Now()).
		Preload("User").
		First(&session).Error

	if err != nil {
		return nil, "", err
	}

	// Generate new JWT token
	newJWTToken, err := GenerateJWT(session.UserID, session.User.Phone)
	if err != nil {
		return nil, "", err
	}

	// Update session with new token and extend expiry
	session.SessionToken = newJWTToken
	session.LastActivity = time.Now()
	session.ExpiresAt = time.Now().Add(24 * time.Hour)

	if err := sm.DB.Save(&session).Error; err != nil {
		return nil, "", err
	}

	return &session, newJWTToken, nil
}

// GetUserSessions gets all active sessions for a user
func (sm *SessionManager) GetUserSessions(userID uint) ([]models.DeviceSessionInfo, error) {
	var sessions []models.DeviceSession
	err := sm.DB.Where("user_id = ? AND is_active = ?", userID, true).
		Order("last_activity DESC").
		Find(&sessions).Error

	if err != nil {
		return nil, err
	}

	var sessionInfos []models.DeviceSessionInfo
	for _, session := range sessions {
		sessionInfos = append(sessionInfos, models.DeviceSessionInfo{
			ID:           session.ID,
			DeviceType:   session.DeviceType,
			DeviceName:   session.DeviceName,
			Provider:     session.Provider,
			IPAddress:    session.IPAddress,
			LastActivity: session.LastActivity,
			IsActive:     session.IsActive,
			CreatedAt:    session.CreatedAt,
		})
	}

	return sessionInfos, nil
}

// LogoutSession logs out a specific session
func (sm *SessionManager) LogoutSession(sessionID uint, userID uint) error {
	return sm.DB.Model(&models.DeviceSession{}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		Update("is_active", false).Error
}

// LogoutAllSessions logs out all sessions for a user
func (sm *SessionManager) LogoutAllSessions(userID uint) error {
	return sm.DB.Model(&models.DeviceSession{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error
}

// LogoutAllOtherSessions logs out all other sessions except current
func (sm *SessionManager) LogoutAllOtherSessions(userID uint, currentSessionID uint) error {
	return sm.DB.Model(&models.DeviceSession{}).
		Where("user_id = ? AND id != ?", userID, currentSessionID).
		Update("is_active", false).Error
}

// CleanupExpiredSessions removes expired sessions
func (sm *SessionManager) CleanupExpiredSessions() error {
	return sm.DB.Model(&models.DeviceSession{}).
		Where("expires_at < ?", time.Now()).
		Update("is_active", false).Error
}

// GetClientIP extracts client IP from Gin context
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return c.ClientIP()
}
