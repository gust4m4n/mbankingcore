package models

// Standard API Response structure
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Alias for APIResponse for backward compatibility
type Response = APIResponse

// Error Response structure
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Helper function untuk success response
func NewSuccessResponse(code int, message string, data interface{}) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Helper function untuk error response
func NewErrorResponse(code int, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// ===========================================
// SUCCESS RESPONSE FUNCTIONS
// ===========================================

// Authentication Success Responses
func LoginSuccessResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_LOGIN_SUCCESS, data)
}

func RegisterSuccessResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_REGISTER_SUCCESS, data)
}

func RefreshSuccessResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_REFRESH_SUCCESS, data)
}

// User Management Success Responses
func UserCreatedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_USER_CREATED, data)
}

func UserRetrievedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_USER_RETRIEVED, data)
}

func UserUpdatedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_USER_UPDATED, data)
}

func UserDeletedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_USER_DELETED, data)
}

func UsersListedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_USERS_LISTED, data)
}

// Profile Management Success Responses
func ProfileRetrievedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_PROFILE_RETRIEVED, data)
}

func ProfileUpdatedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_PROFILE_UPDATED, data)
}

func PasswordChangedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_PASSWORD_CHANGED, data)
}

// Config Management Success Responses
func ConfigCreatedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_CONFIG_CREATED, data)
}

func ConfigRetrievedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_CONFIG_RETRIEVED, data)
}

func ConfigUpdatedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_CONFIG_UPDATED, data)
}

func ConfigDeletedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_CONFIG_DELETED, data)
}

func ConfigsListedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_CONFIGS_LISTED, data)
}

// Terms & Conditions Success Responses
func TermsConditionsRetrievedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_TERMS_CONDITIONS_RETRIEVED, data)
}

func TermsConditionsUpdatedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_TERMS_CONDITIONS_UPDATED, data)
}

// Privacy Policy Success Responses
func PrivacyPolicyRetrievedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_PRIVACY_POLICY_RETRIEVED, data)
}

func PrivacyPolicyUpdatedResponse(data interface{}) APIResponse {
	return NewSuccessResponse(CODE_SUCCESS, MSG_PRIVACY_POLICY_UPDATED, data)
}

// ===========================================
// ERROR RESPONSE FUNCTIONS
// ===========================================

// General/System Error Responses
func InternalServerResponse() ErrorResponse {
	return NewErrorResponse(CODE_INTERNAL_SERVER, MSG_INTERNAL_SERVER)
}

func DatabaseFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_DATABASE_FAILED, MSG_DATABASE_FAILED)
}

func ValidationFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_VALIDATION_FAILED, MSG_VALIDATION_FAILED)
}

func InvalidRequestResponse() ErrorResponse {
	return NewErrorResponse(CODE_INVALID_REQUEST, MSG_INVALID_REQUEST)
}

func NotFoundResponse() ErrorResponse {
	return NewErrorResponse(CODE_NOT_FOUND, MSG_NOT_FOUND)
}

// Authentication Error Responses
func UnauthorizedResponse() ErrorResponse {
	return NewErrorResponse(CODE_UNAUTHORIZED, MSG_UNAUTHORIZED)
}

func InvalidTokenResponse() ErrorResponse {
	return NewErrorResponse(CODE_INVALID_TOKEN, MSG_INVALID_TOKEN)
}

func TokenExpiredResponse() ErrorResponse {
	return NewErrorResponse(CODE_TOKEN_EXPIRED, MSG_TOKEN_EXPIRED)
}

func InvalidPasswordResponse() ErrorResponse {
	return NewErrorResponse(CODE_INVALID_PASSWORD, MSG_INVALID_PASSWORD)
}

func MissingTokenResponse() ErrorResponse {
	return NewErrorResponse(CODE_MISSING_TOKEN, MSG_MISSING_TOKEN)
}

func InvalidEmailResponse() ErrorResponse {
	return NewErrorResponse(CODE_INVALID_EMAIL, MSG_INVALID_EMAIL)
}

func EmailExistsResponse() ErrorResponse {
	return NewErrorResponse(CODE_EMAIL_EXISTS, MSG_EMAIL_EXISTS)
}

func RegisterFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_REGISTER_FAILED, MSG_REGISTER_FAILED)
}

func LoginFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_LOGIN_FAILED, MSG_LOGIN_FAILED)
}

func RefreshFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_REFRESH_FAILED, MSG_REFRESH_FAILED)
}

// User Management Error Responses
func UserNotFoundResponse() ErrorResponse {
	return NewErrorResponse(CODE_USER_NOT_FOUND, MSG_USER_NOT_FOUND)
}

func InvalidUserIDResponse() ErrorResponse {
	return NewErrorResponse(CODE_INVALID_USER_ID, MSG_INVALID_USER_ID)
}

func UserCreateFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_USER_CREATE_FAILED, MSG_USER_CREATE_FAILED)
}

func UserUpdateFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_USER_UPDATE_FAILED, MSG_USER_UPDATE_FAILED)
}

func UserDeleteFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_USER_DELETE_FAILED, MSG_USER_DELETE_FAILED)
}

func UserRetrieveFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_USER_RETRIEVE_FAILED, MSG_USER_RETRIEVE_FAILED)
}

func UserListFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_USER_LIST_FAILED, MSG_USER_LIST_FAILED)
}

// Profile Management Error Responses
func ProfileNotFoundResponse() ErrorResponse {
	return NewErrorResponse(CODE_PROFILE_NOT_FOUND, MSG_PROFILE_NOT_FOUND)
}

func ProfileUpdateFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_PROFILE_UPDATE_FAILED, MSG_PROFILE_UPDATE_FAILED)
}

func ProfileRetrieveFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_PROFILE_RETRIEVE_FAILED, MSG_PROFILE_RETRIEVE_FAILED)
}

func PasswordChangeFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_PASSWORD_CHANGE_FAILED, MSG_PASSWORD_CHANGE_FAILED)
}

func CurrentPasswordInvalidResponse() ErrorResponse {
	return NewErrorResponse(CODE_CURRENT_PASSWORD_INVALID, MSG_CURRENT_PASSWORD_INVALID)
}

// Config Management Error Responses
func ConfigNotFoundResponse() ErrorResponse {
	return NewErrorResponse(CODE_CONFIG_NOT_FOUND, MSG_CONFIG_NOT_FOUND)
}

func ConfigCreateFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_CONFIG_CREATE_FAILED, MSG_CONFIG_CREATE_FAILED)
}

func ConfigUpdateFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_CONFIG_UPDATE_FAILED, MSG_CONFIG_UPDATE_FAILED)
}

func ConfigDeleteFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_CONFIG_DELETE_FAILED, MSG_CONFIG_DELETE_FAILED)
}

func ConfigRetrieveFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_CONFIG_RETRIEVE_FAILED, MSG_CONFIG_RETRIEVE_FAILED)
}

func ConfigKeyInvalidResponse() ErrorResponse {
	return NewErrorResponse(CODE_CONFIG_KEY_INVALID, MSG_CONFIG_KEY_INVALID)
}

func ConfigValueInvalidResponse() ErrorResponse {
	return NewErrorResponse(CODE_CONFIG_VALUE_INVALID, MSG_CONFIG_VALUE_INVALID)
}

// Terms & Conditions Error Responses
func TermsConditionsNotFoundResponse() ErrorResponse {
	return NewErrorResponse(CODE_TERMS_CONDITIONS_NOT_FOUND, MSG_TERMS_CONDITIONS_NOT_FOUND)
}

func TermsConditionsCreateFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_TERMS_CONDITIONS_CREATE_FAILED, MSG_TERMS_CONDITIONS_CREATE_FAILED)
}

func TermsConditionsUpdateFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_TERMS_CONDITIONS_UPDATE_FAILED, MSG_TERMS_CONDITIONS_UPDATE_FAILED)
}

func TermsConditionsDeleteFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_TERMS_CONDITIONS_DELETE_FAILED, MSG_TERMS_CONDITIONS_DELETE_FAILED)
}

func TermsConditionsRetrieveFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_TERMS_CONDITIONS_RETRIEVE_FAILED, MSG_TERMS_CONDITIONS_RETRIEVE_FAILED)
}

// Privacy Policy Error Responses
func PrivacyPolicyNotFoundResponse() ErrorResponse {
	return NewErrorResponse(CODE_PRIVACY_POLICY_NOT_FOUND, MSG_PRIVACY_POLICY_NOT_FOUND)
}

func PrivacyPolicyCreateFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_PRIVACY_POLICY_CREATE_FAILED, MSG_PRIVACY_POLICY_CREATE_FAILED)
}

func PrivacyPolicyUpdateFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_PRIVACY_POLICY_UPDATE_FAILED, MSG_PRIVACY_POLICY_UPDATE_FAILED)
}

func PrivacyPolicyDeleteFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_PRIVACY_POLICY_DELETE_FAILED, MSG_PRIVACY_POLICY_DELETE_FAILED)
}

func PrivacyPolicyRetrieveFailedResponse() ErrorResponse {
	return NewErrorResponse(CODE_PRIVACY_POLICY_RETRIEVE_FAILED, MSG_PRIVACY_POLICY_RETRIEVE_FAILED)
}

// Permission/Access Error Responses
func ForbiddenResponse() ErrorResponse {
	return NewErrorResponse(CODE_FORBIDDEN, MSG_FORBIDDEN)
}

func InsufficientPermissionsResponse() ErrorResponse {
	return NewErrorResponse(CODE_INSUFFICIENT_PERMISSIONS, MSG_INSUFFICIENT_PERMISSIONS)
}

func AdminRequiredResponse() ErrorResponse {
	return NewErrorResponse(CODE_ADMIN_REQUIRED, MSG_ADMIN_REQUIRED)
}

func OwnerRequiredResponse() ErrorResponse {
	return NewErrorResponse(CODE_OWNER_REQUIRED, MSG_OWNER_REQUIRED)
}

// ===========================================
// DEPRECATED FUNCTIONS (for backward compatibility)
// ===========================================

func CreateFailedResponse() ErrorResponse {
	return UserCreateFailedResponse() // Deprecated: Use specific error responses
}

func UpdateFailedResponse() ErrorResponse {
	return UserUpdateFailedResponse() // Deprecated: Use specific error responses
}

func DeleteFailedResponse() ErrorResponse {
	return UserDeleteFailedResponse() // Deprecated: Use specific error responses
}

func RetrieveFailedResponse() ErrorResponse {
	return UserRetrieveFailedResponse() // Deprecated: Use specific error responses
}
