# MBankingCore API Complete Endpoint Reference

This document provides a complete list of all 44 available API endpoints in the MBankingCore system.

## 📊 API Summary

**Total Endpoints: 44**

### 🔓 Public APIs (7 endpoints)
- `GET /health` - Health check
- `GET /api/terms-conditions` - Get terms & conditions
- `POST /api/terms-conditions` - Set terms & conditions (Admin)
- `GET /api/privacy-policy` - Get privacy policy
- `POST /api/privacy-policy` - Set privacy policy (Admin)
- `GET /api/onboardings` - Get all onboardings
- `GET /api/onboardings/:id` - Get onboarding by ID

### 🔐 Banking Authentication APIs (3 endpoints)
- `POST /api/login` - Banking login step 1 (send OTP)
- `POST /api/login/verify` - Banking login step 2 (verify OTP)
- `POST /api/refresh` - Refresh access token

### 🛡️ Protected APIs (18 endpoints)

#### User Profile Management (3 endpoints)
- `GET /api/profile` - Get user profile
- `PUT /api/profile` - Update user profile
- `PUT /api/change-pin` - Change PIN ATM

#### Session Management (3 endpoints)
- `GET /api/sessions` - Get active sessions
- `POST /api/logout` - Logout (current or all devices)
- `POST /api/logout-others` - Logout other sessions

#### Bank Account Management (5 endpoints)
- `GET /api/bank-accounts` - Get user's bank accounts
- `POST /api/bank-accounts` - Create new bank account
- `PUT /api/bank-accounts/:id` - Update bank account
- `DELETE /api/bank-accounts/:id` - Delete bank account
- `PUT /api/bank-accounts/:id/primary` - Set primary account

#### Article Management (5 endpoints)
- `GET /api/articles` - Get all articles (with pagination)
- `GET /api/articles/:id` - Get article by ID
- `PUT /api/articles/:id` - Update article (own only)
- `DELETE /api/articles/:id` - Delete article (own only)
- `GET /api/my-articles` - Get my articles

#### Photo Management (4 endpoints)
- `GET /api/photos` - Get all photos (with pagination)
- `GET /api/photos/:id` - Get photo by ID
- `PUT /api/photos/:id` - Update photo (own only)
- `DELETE /api/photos/:id` - Delete photo (own only)

#### Configuration (1 endpoint)
- `GET /api/config/:key` - Get config value by key

### 👑 Admin APIs (14 endpoints)

#### Article Management (1 endpoint)
- `POST /api/articles` - Create article (Admin/Owner only)

#### Onboarding Management (3 endpoints)
- `POST /api/onboardings` - Create onboarding (Admin/Owner only)
- `PUT /api/onboardings/:id` - Update onboarding (Admin/Owner only)
- `DELETE /api/onboardings/:id` - Delete onboarding (Admin/Owner only)

#### Photo Management (1 endpoint)
- `POST /api/photos` - Create photo (Admin/Owner only)

#### User Management (4 endpoints)
- `GET /api/users` - List all users (Admin/Owner only)
- `GET /api/admin/users` - List admin and owner users (Admin/Owner only)
- `GET /api/users/:id` - Get user by ID (Admin/Owner only)
- `DELETE /api/users/:id` - Delete user by ID (Admin/Owner only)

#### Configuration Management (3 endpoints)
- `POST /api/config` - Set config value (Admin/Owner only)
- `GET /api/admin/configs` - Get all configs (Admin/Owner only)
- `DELETE /api/config/:key` - Delete config by key (Admin/Owner only)

#### Terms & Conditions (1 endpoint)
*Note: This is actually listed in Public APIs section*

#### Privacy Policy (1 endpoint)
*Note: This is actually listed in Public APIs section*

### 👨‍💼 Owner-Only APIs (2 endpoints)

#### User Management with Roles (2 endpoints)
- `POST /api/users` - Create user with any role (Owner only)
- `PUT /api/users/:id` - Update user including role changes (Owner only)

## 🔑 Authentication Levels

1. **Public** - No authentication required
2. **Protected** - Requires Bearer token (any authenticated user)
3. **Admin** - Requires Bearer token + Admin or Owner role
4. **Owner** - Requires Bearer token + Owner role only

## 📱 Testing Setup

### Environment Variables Required:
- `base_url` - API base URL (e.g., http://localhost:8080)
- `banking_account_number` - Unique 16-digit account number
- `banking_phone` - Phone number for registration
- `banking_name` - Full name (8+ characters)
- `banking_mother_name` - Mother's name (8+ characters)
- `banking_pin_atm` - 6-digit PIN
- `banking_otp_code` - OTP code (for testing, use any 6-digit number)
- `device_id_banking` - Unique device identifier

### Auto-Generated Variables:
- `access_token` - Generated after successful login
- `refresh_token` - Generated after successful login
- `login_token` - Generated during step 1 of banking login
- `user_id` - User ID after authentication
- `session_id` - Session ID after authentication

## 🎯 Testing Flow

1. **Start with Public APIs** to verify basic connectivity
2. **Banking Authentication** to get tokens
3. **Test Protected APIs** with user token
4. **Admin APIs** (if user has admin role)
5. **Owner APIs** (if user has owner role)

## 📄 Postman Collections Available

1. **MBankingCore-API.postman_collection.json** - Basic collection with core endpoints
2. **MBankingCore-API-Complete.postman_collection.json** - Complete collection with all endpoints
3. **MBankingCore-API.postman_environment.json** - Environment variables

## 🚀 Quick Start

1. Import both Postman collection and environment files
2. Update environment variables with unique test data
3. Run "Banking Login (Step 1)" to get login_token
4. Run "Banking Login Verification (Step 2)" to get access tokens
5. Test any protected endpoint with automatic token handling

---

*Last Updated: July 30, 2025*  
*Total Endpoints: 44*  
*Authentication: JWT-based with multi-device session management*
