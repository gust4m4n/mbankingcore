package models

import (
	"time"
)

// DeviceType represents the type of device/platform
type DeviceType string

const (
	DeviceTypeAndroid     DeviceType = "android"
	DeviceTypeIOS         DeviceType = "ios"
	DeviceTypeWeb         DeviceType = "web"
	DeviceTypeDesktop     DeviceType = "desktop"
	DeviceTypeGoogleSSO   DeviceType = "google_sso"
	DeviceTypeAppleSSO    DeviceType = "apple_sso"
	DeviceTypeFacebookSSO DeviceType = "facebook_sso"
)

// LoginProvider represents the authentication provider
type LoginProvider string

const (
	LoginProviderEmail    LoginProvider = "email"
	LoginProviderGoogle   LoginProvider = "google"
	LoginProviderApple    LoginProvider = "apple"
	LoginProviderFacebook LoginProvider = "facebook"
)

// DeviceSession represents an active user session on a specific device
type DeviceSession struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	User         User           `json:"user" gorm:"foreignKey:UserID"`
	SessionToken string         `json:"session_token" gorm:"unique;not null;size:255"`
	RefreshToken string         `json:"refresh_token" gorm:"unique;not null;size:255"`
	DeviceType   DeviceType     `json:"device_type" gorm:"not null;size:50"`
	DeviceID     string        `json:"device_id" gorm:"size:255;index"`
	DeviceName   string        `json:"device_name" gorm:"size:255"`
	Provider     LoginProvider `json:"provider" gorm:"not null;size:50"`
	ProviderID   string        `json:"provider_id" gorm:"size:255"`
	IPAddress    string        `json:"ip_address" gorm:"size:45"`
	UserAgent    string        `json:"user_agent" gorm:"size:500"`
	IsActive     bool          `json:"is_active" gorm:"default:true"`
	LastActivity time.Time     `json:"last_activity" gorm:"autoUpdateTime"`
	ExpiresAt    time.Time     `json:"expires_at" gorm:"not null"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// LoginRequest for multi-platform authentication
type MultiPlatformLoginRequest struct {
	Email      string        `json:"email"`
	Password   string        `json:"password"` // SHA256 hash
	Provider   LoginProvider `json:"provider"`
	ProviderID string        `json:"provider_id,omitempty"`
	DeviceInfo DeviceInfo    `json:"device_info"`
}

// DeviceInfo contains information about the device
type DeviceInfo struct {
	DeviceType DeviceType `json:"device_type"`
	DeviceID   string     `json:"device_id"`
	DeviceName string     `json:"device_name"`
	UserAgent  string     `json:"user_agent,omitempty"`
}

// MultiPlatformLoginResponse for successful authentication
type MultiPlatformLoginResponse struct {
	User         User       `json:"user"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in"`
	SessionID    uint       `json:"session_id"`
	DeviceInfo   DeviceInfo `json:"device_info"`
}

// UserSessionsResponse for listing active sessions
type UserSessionsResponse struct {
	Sessions []DeviceSessionInfo `json:"sessions"`
	Total    int                 `json:"total"`
}

// DeviceSessionInfo for session listing
type DeviceSessionInfo struct {
	ID           uint          `json:"id"`
	DeviceType   DeviceType    `json:"device_type"`
	DeviceName   string        `json:"device_name"`
	Provider     LoginProvider `json:"provider"`
	IPAddress    string        `json:"ip_address"`
	LastActivity time.Time     `json:"last_activity"`
	IsActive     bool          `json:"is_active"`
	IsCurrent    bool          `json:"is_current"`
	CreatedAt    time.Time     `json:"created_at"`
}

// RefreshTokenRequest for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutRequest for specific device logout
type LogoutRequest struct {
	SessionID  *uint `json:"session_id,omitempty"`
	AllDevices bool  `json:"all_devices,omitempty"`
}
