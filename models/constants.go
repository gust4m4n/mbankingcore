package models

// Response Codes
const (
	// ===========================================
	// GENERAL SUCCESS (200)
	// ===========================================
	CODE_SUCCESS = 200

	// ===========================================
	// GENERAL/SYSTEM ERRORS (250-299)
	// ===========================================
	CODE_INTERNAL_SERVER   = 250
	CODE_DATABASE_FAILED   = 251
	CODE_VALIDATION_FAILED = 252
	CODE_INVALID_REQUEST   = 253
	CODE_NOT_FOUND         = 254

	// ===========================================
	// AUTHENTICATION ERRORS (300-349) - Public endpoints
	// ===========================================
	CODE_UNAUTHORIZED     = 300
	CODE_INVALID_TOKEN    = 301
	CODE_TOKEN_EXPIRED    = 302
	CODE_INVALID_PASSWORD = 303
	CODE_MISSING_TOKEN    = 304
	CODE_INVALID_EMAIL    = 305
	CODE_EMAIL_EXISTS     = 306
	CODE_PHONE_EXISTS     = 307
	CODE_REGISTER_FAILED  = 308
	CODE_LOGIN_FAILED     = 309
	CODE_REFRESH_FAILED   = 310

	// ===========================================
	// PUBLIC CONTENT ERRORS (350-399) - Terms, Privacy Policy
	// ===========================================
	CODE_TERMS_CONDITIONS_NOT_FOUND       = 360
	CODE_TERMS_CONDITIONS_CREATE_FAILED   = 361
	CODE_TERMS_CONDITIONS_UPDATE_FAILED   = 362
	CODE_TERMS_CONDITIONS_DELETE_FAILED   = 363
	CODE_TERMS_CONDITIONS_RETRIEVE_FAILED = 364

	CODE_PRIVACY_POLICY_NOT_FOUND       = 370
	CODE_PRIVACY_POLICY_CREATE_FAILED   = 371
	CODE_PRIVACY_POLICY_UPDATE_FAILED   = 372
	CODE_PRIVACY_POLICY_DELETE_FAILED   = 373
	CODE_PRIVACY_POLICY_RETRIEVE_FAILED = 374

	// ===========================================
	// USER MANAGEMENT ERRORS (400-449) - Protected endpoints
	// ===========================================
	CODE_USER_NOT_FOUND       = 400
	CODE_INVALID_USER_ID      = 401
	CODE_USER_CREATE_FAILED   = 402
	CODE_USER_UPDATE_FAILED   = 403
	CODE_USER_DELETE_FAILED   = 404
	CODE_USER_RETRIEVE_FAILED = 405
	CODE_USER_LIST_FAILED     = 406

	// ===========================================
	// PROFILE MANAGEMENT ERRORS (450-499) - Protected endpoints
	// ===========================================
	CODE_PROFILE_NOT_FOUND        = 450
	CODE_PROFILE_UPDATE_FAILED    = 451
	CODE_PROFILE_RETRIEVE_FAILED  = 452
	CODE_PASSWORD_CHANGE_FAILED   = 453
	CODE_CURRENT_PASSWORD_INVALID = 454

	// ===========================================
	// ARTICLE MANAGEMENT ERRORS (500-549) - Protected endpoints
	// ===========================================
	// Reserved for future article management features

	// ===========================================
	// PHOTO MANAGEMENT ERRORS (550-599) - Protected endpoints
	// ===========================================
	// Reserved for future photo management features

	// ===========================================
	// CONFIG MANAGEMENT ERRORS (600-649) - Protected/Admin endpoints
	// ===========================================
	CODE_CONFIG_NOT_FOUND       = 600
	CODE_CONFIG_CREATE_FAILED   = 601
	CODE_CONFIG_UPDATE_FAILED   = 602
	CODE_CONFIG_DELETE_FAILED   = 603
	CODE_CONFIG_RETRIEVE_FAILED = 604
	CODE_CONFIG_KEY_INVALID     = 605
	CODE_CONFIG_VALUE_INVALID   = 606

	// ===========================================
	// ADMIN USER MANAGEMENT ERRORS (650-699) - Admin only
	// ===========================================
	// Reserved for future admin user management features

	// ===========================================
	// ADMIN CONTENT MANAGEMENT ERRORS (700-749) - Admin only
	// ===========================================
	// Reserved for future admin content management features

	// ===========================================
	// PERMISSION/ACCESS ERRORS (750-799)
	// ===========================================
	CODE_FORBIDDEN                = 750
	CODE_INSUFFICIENT_PERMISSIONS = 751
	CODE_ADMIN_REQUIRED           = 752
	CODE_OWNER_REQUIRED           = 753
)

// Success Messages
const (
	MSG_USER_CREATED               = "User created successfully"
	MSG_USER_RETRIEVED             = "User retrieved successfully"
	MSG_USER_UPDATED               = "User updated successfully"
	MSG_USER_DELETED               = "User deleted successfully"
	MSG_USERS_LISTED               = "Users retrieved successfully"
	MSG_PROFILE_UPDATED            = "Profile updated successfully"
	MSG_PROFILE_RETRIEVED          = "Profile retrieved successfully"
	MSG_PASSWORD_CHANGED           = "Password changed successfully"
	MSG_ARTICLE_CREATED            = "Article created successfully"
	MSG_ARTICLE_RETRIEVED          = "Article retrieved successfully"
	MSG_ARTICLE_UPDATED            = "Article updated successfully"
	MSG_ARTICLE_DELETED            = "Article deleted successfully"
	MSG_ARTICLES_LISTED            = "Articles retrieved successfully"
	MSG_PHOTO_CREATED              = "Photo created successfully"
	MSG_PHOTO_RETRIEVED            = "Photo retrieved successfully"
	MSG_PHOTO_UPDATED              = "Photo updated successfully"
	MSG_PHOTO_DELETED              = "Photo deleted successfully"
	MSG_PHOTOS_LISTED              = "Photos retrieved successfully"
	MSG_CONFIG_CREATED             = "Configuration created successfully"
	MSG_CONFIG_RETRIEVED           = "Configuration retrieved successfully"
	MSG_CONFIG_UPDATED             = "Configuration updated successfully"
	MSG_CONFIG_DELETED             = "Configuration deleted successfully"
	MSG_CONFIGS_LISTED             = "Configurations retrieved successfully"
	MSG_ONBOARDING_CREATED         = "Onboarding created successfully"
	MSG_ONBOARDING_RETRIEVED       = "Onboarding retrieved successfully"
	MSG_ONBOARDING_UPDATED         = "Onboarding updated successfully"
	MSG_ONBOARDING_DELETED         = "Onboarding deleted successfully"
	MSG_ONBOARDINGS_LISTED         = "Onboardings retrieved successfully"
	MSG_TERMS_CONDITIONS_CREATED   = "Terms and conditions created successfully"
	MSG_TERMS_CONDITIONS_RETRIEVED = "Terms and conditions retrieved successfully"
	MSG_TERMS_CONDITIONS_UPDATED   = "Terms and conditions updated successfully"
	MSG_TERMS_CONDITIONS_DELETED   = "Terms and conditions deleted successfully"
	MSG_PRIVACY_POLICY_CREATED     = "Privacy policy created successfully"
	MSG_PRIVACY_POLICY_RETRIEVED   = "Privacy policy retrieved successfully"
	MSG_PRIVACY_POLICY_UPDATED     = "Privacy policy updated successfully"
	MSG_PRIVACY_POLICY_DELETED     = "Privacy policy deleted successfully"
	MSG_LOGIN_SUCCESS              = "Login successful"
	MSG_REGISTER_SUCCESS           = "Registration successful"
	MSG_REFRESH_SUCCESS            = "Token refreshed successfully"
)

// Error Messages
const (
	// General/System Error Messages
	MSG_INTERNAL_SERVER   = "Internal server error"
	MSG_DATABASE_FAILED   = "Database operation failed"
	MSG_VALIDATION_FAILED = "Validation failed"
	MSG_INVALID_REQUEST   = "Invalid request data"
	MSG_NOT_FOUND         = "Resource not found"

	// Authentication Error Messages
	MSG_UNAUTHORIZED     = "Unauthorized access"
	MSG_INVALID_TOKEN    = "Invalid or malformed token"
	MSG_TOKEN_EXPIRED    = "Token has expired"
	MSG_INVALID_PASSWORD = "Invalid phone or password"
	MSG_MISSING_TOKEN    = "Authorization token required"
	MSG_INVALID_EMAIL    = "Invalid email format"
	MSG_EMAIL_EXISTS     = "Email already exists"
	MSG_PHONE_EXISTS     = "Phone already exists"
	MSG_REGISTER_FAILED  = "Registration failed"
	MSG_LOGIN_FAILED     = "Login failed"
	MSG_REFRESH_FAILED   = "Token refresh failed"

	// Public Content Error Messages
	MSG_TERMS_CONDITIONS_NOT_FOUND       = "Terms and conditions not found"
	MSG_TERMS_CONDITIONS_CREATE_FAILED   = "Failed to create terms and conditions"
	MSG_TERMS_CONDITIONS_UPDATE_FAILED   = "Failed to update terms and conditions"
	MSG_TERMS_CONDITIONS_DELETE_FAILED   = "Failed to delete terms and conditions"
	MSG_TERMS_CONDITIONS_RETRIEVE_FAILED = "Failed to retrieve terms and conditions"

	MSG_PRIVACY_POLICY_NOT_FOUND       = "Privacy policy not found"
	MSG_PRIVACY_POLICY_CREATE_FAILED   = "Failed to create privacy policy"
	MSG_PRIVACY_POLICY_UPDATE_FAILED   = "Failed to update privacy policy"
	MSG_PRIVACY_POLICY_DELETE_FAILED   = "Failed to delete privacy policy"
	MSG_PRIVACY_POLICY_RETRIEVE_FAILED = "Failed to retrieve privacy policy"

	// User Management Error Messages
	MSG_USER_NOT_FOUND       = "User not found"
	MSG_INVALID_USER_ID      = "Invalid user ID"
	MSG_USER_CREATE_FAILED   = "Failed to create user"
	MSG_USER_UPDATE_FAILED   = "Failed to update user"
	MSG_USER_DELETE_FAILED   = "Failed to delete user"
	MSG_USER_RETRIEVE_FAILED = "Failed to retrieve user"
	MSG_USER_LIST_FAILED     = "Failed to retrieve users list"

	// Profile Management Error Messages
	MSG_PROFILE_NOT_FOUND        = "Profile not found"
	MSG_PROFILE_UPDATE_FAILED    = "Failed to update profile"
	MSG_PROFILE_RETRIEVE_FAILED  = "Failed to retrieve profile"
	MSG_PASSWORD_CHANGE_FAILED   = "Failed to change password"
	MSG_CURRENT_PASSWORD_INVALID = "Current password is invalid"

	// Config Management Error Messages
	MSG_CONFIG_NOT_FOUND       = "Configuration not found"
	MSG_CONFIG_CREATE_FAILED   = "Failed to create configuration"
	MSG_CONFIG_UPDATE_FAILED   = "Failed to update configuration"
	MSG_CONFIG_DELETE_FAILED   = "Failed to delete configuration"
	MSG_CONFIG_RETRIEVE_FAILED = "Failed to retrieve configuration"
	MSG_CONFIG_KEY_INVALID     = "Invalid configuration key"
	MSG_CONFIG_VALUE_INVALID   = "Invalid configuration value"

	// Permission/Access Error Messages
	MSG_FORBIDDEN                = "Access forbidden - insufficient permissions"
	MSG_INSUFFICIENT_PERMISSIONS = "Insufficient permissions to perform this action"
	MSG_ADMIN_REQUIRED           = "Admin privileges required"
	MSG_OWNER_REQUIRED           = "Owner privileges required"
)
