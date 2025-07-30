# MBankingCore API Documentation

Dokumentasi lengkap untuk RESTful API MBankingCore dengan JWT Authentication, Multi-Device Session Management, dan Multi-Account Banking.

## 📖 Response Format

Semua API menggunakan format response yang konsisten:

**Success Response:**

```json
{
    "code": 200,
    "message": "Success message",
    "data": {
        // Response data
    }
}
```

**Error Response:**

```json
{
    "code": 400,
    "message": "Error message"
}
```

## 📚 Gambaran API

API diorganisir ke dalam bagian-bagian berikut:

### 🔓 Public APIs (7 endpoints)

- **[Health Check](#1-health-check)** - Status kesehatan server (1 endpoint)
- **[Terms & Conditions APIs](#2-terms--conditions-apis)** - Manajemen syarat dan ketentuan (2 endpoints)
- **[Privacy Policy APIs](#3-privacy-policy-apis)** - Manajemen kebijakan privasi (2 endpoints)
- **[Onboarding APIs](#4-onboarding-apis)** - Konten onboarding aplikasi (2 endpoints)

### 🔐 Banking Authentication APIs (3 endpoints)

- **[Banking Authentication](#5-banking-authentication-apis)** - 2-step banking authentication dengan OTP

### 🛡️ Protected APIs (18 endpoints)

- **[User Profile Management](#6-user-profile-apis)** - Manajemen profil user (3 endpoints)
- **[Session Management](#7-session-management-apis)** - Manajemen sesi device (3 endpoints)
- **[Bank Account Management](#8-bank-account-management-apis)** - Multi-account banking CRUD (5 endpoints)
- **[Article Management](#9-article-management-apis)** - Operasi CRUD artikel (5 endpoints)
- **[Photo Management](#10-photo-management-apis)** - Sistem manajemen foto (4 endpoints)
- **[Configuration APIs](#11-configuration-apis)** - Read config (1 endpoint)

### 👑 Admin APIs (21 endpoints)

- **[Admin Management](#12-admin-management-apis)** - Admin authentication & CRUD (7 endpoints)
- **[Admin Article Management](#13-admin-article-management)** - Create artikel (1 endpoint)
- **[Admin Onboarding Management](#14-admin-onboarding-management)** - CRUD onboarding (3 endpoints)
- **[Admin Photo Management](#15-admin-photo-management)** - Create photo (1 endpoint)
- **[Admin User Management](#16-admin-user-management)** - Manajemen user (4 endpoints)
- **[Admin Configuration](#17-admin-configuration-apis)** - Full config management (3 endpoints)
- **[Admin Terms & Conditions](#18-admin-terms-conditions)** - Set T&C (1 endpoint)
- **[Admin Privacy Policy](#19-admin-privacy-policy)** - Set Privacy Policy (1 endpoint)

### 👨‍💼 Owner-Only APIs (2 endpoints)

- **[Owner User Management](#20-owner-user-management)** - Create & update users dengan roles (2 endpoints)

**Total: 51 Active Endpoints**

---

# � API Endpoint Quick Reference

This section provides a complete list of all 51 available API endpoints organized by access level.

## 🔓 Public APIs (7 endpoints)

- `GET /health` - Health check
- `GET /api/terms-conditions` - Get terms & conditions
- `POST /api/terms-conditions` - Set terms & conditions (Admin)
- `GET /api/privacy-policy` - Get privacy policy
- `POST /api/privacy-policy` - Set privacy policy (Admin)
- `GET /api/onboardings` - Get all onboardings
- `GET /api/onboardings/:id` - Get onboarding by ID

## 🔐 Banking Authentication APIs (3 endpoints)

- `POST /api/login` - Banking login step 1 (send OTP)
- `POST /api/login/verify` - Banking login step 2 (verify OTP)
- `POST /api/refresh` - Refresh access token

## 🛡️ Protected APIs (18 endpoints)

### User Profile Management (3 endpoints)

- `GET /api/profile` - Get user profile
- `PUT /api/profile` - Update user profile
- `PUT /api/change-pin` - Change PIN ATM

### Session Management (3 endpoints)

- `GET /api/sessions` - Get active sessions
- `POST /api/logout` - Logout (current or all devices)
- `POST /api/logout-others` - Logout other sessions

### Bank Account Management (5 endpoints)

- `GET /api/bank-accounts` - Get user's bank accounts
- `POST /api/bank-accounts` - Create new bank account
- `PUT /api/bank-accounts/:id` - Update bank account
- `DELETE /api/bank-accounts/:id` - Delete bank account
- `PUT /api/bank-accounts/:id/primary` - Set primary account

### Article Management (5 endpoints)

- `GET /api/articles` - Get all articles (with pagination)
- `GET /api/articles/:id` - Get article by ID
- `PUT /api/articles/:id` - Update article (own only)
- `DELETE /api/articles/:id` - Delete article (own only)
- `GET /api/my-articles` - Get my articles

### Photo Management (4 endpoints)

- `GET /api/photos` - Get all photos (with pagination)
- `GET /api/photos/:id` - Get photo by ID
- `PUT /api/photos/:id` - Update photo (own only)
- `DELETE /api/photos/:id` - Delete photo (own only)

### Configuration (1 endpoint)

- `GET /api/config/:key` - Get config value by key

## 👑 Admin APIs (21 endpoints)

### Admin Management (7 endpoints)

- `POST /api/admin/login` - Admin login
- `POST /api/admin/logout` - Admin logout
- `GET /api/admin/admins` - Get all admins (Admin only)
- `GET /api/admin/admins/:id` - Get admin by ID (Admin only)
- `POST /api/admin/admins` - Create admin (Super Admin only)
- `PUT /api/admin/admins/:id` - Update admin (Super Admin only)
- `DELETE /api/admin/admins/:id` - Delete admin (Super Admin only)

### Article Management (1 endpoint)

- `POST /api/articles` - Create article (Admin/Owner only)

### Onboarding Management (3 endpoints)

- `POST /api/onboardings` - Create onboarding (Admin/Owner only)
- `PUT /api/onboardings/:id` - Update onboarding (Admin/Owner only)
- `DELETE /api/onboardings/:id` - Delete onboarding (Admin/Owner only)

### Photo Management (1 endpoint)

- `POST /api/photos` - Create photo (Admin/Owner only)

### User Management (4 endpoints)

- `GET /api/users` - List all users (Admin/Owner only)
- `GET /api/users/:id` - Get user by ID (Admin/Owner only)
- `DELETE /api/users/:id` - Delete user by ID (Admin/Owner only)

### Configuration Management (3 endpoints)

- `POST /api/config` - Set config value (Admin/Owner only)
- `GET /api/configs` - Get all configs (Admin/Owner only)
- `DELETE /api/config/:key` - Delete config by key (Admin/Owner only)

### Terms & Conditions (1 endpoint)

*Note: This is actually listed in Public APIs section*

### Privacy Policy (1 endpoint)

*Note: This is actually listed in Public APIs section*

## 👨‍💼 Owner-Only APIs (2 endpoints)

### User Management with Roles (2 endpoints)

- `POST /api/users` - Create user with any role (Owner only)
- `PUT /api/users/:id` - Update user including role changes (Owner only)

## 🔑 Authentication Levels

1. **Public** - No authentication required
2. **Protected** - Requires Bearer token (any authenticated user)
3. **Admin** - Requires Bearer token + Admin or Owner role
4. **Owner** - Requires Bearer token + Owner role only

## 📱 Testing Setup

### Environment Variables Required

- `base_url` - API base URL (e.g., <http://localhost:8080>)
- `banking_account_number` - Unique 16-digit account number
- `banking_phone` - Phone number for registration
- `banking_name` - Full name (8+ characters)
- `banking_mother_name` - Mother's name (8+ characters)
- `banking_pin_atm` - 6-digit PIN
- `banking_otp_code` - OTP code (for testing, use any 6-digit number)
- `device_id_banking` - Unique device identifier

### Auto-Generated Variables

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
2. **MBankingCore-API.postman_environment.json** - Environment variables

## 🚀 Quick Start

1. Import both Postman collection and environment files
2. Update environment variables with unique test data
3. Run "Banking Login (Step 1)" to get login_token
4. Run "Banking Login Verification (Step 2)" to get access tokens
5. Test any protected endpoint with automatic token handling

---

# �🔌 API Endpoints

## 1. Health Check

**Endpoint:** `GET /health`  
**Description:** Check if the API server is running and healthy  
**Authentication:** None required

**Response (200):**

```json
{
  "code": 200,
  "message": "MBankingCore API is running",
  "data": {
    "status": "ok"
  }
}
```

---

## 2. Terms & Conditions APIs (Config-Based)

**🔄 Updated System**: Terms & Conditions now integrated with Config system using key `'tnc'` for simplified management.

### 2.1 Get Terms & Conditions (Public)

**Endpoint:** `GET /api/terms-conditions`  
**Access:** Public  
**Description:** Retrieve the current terms and conditions content from config

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Terms and conditions retrieved successfully",
    "data": {
        "content": "# Terms and Conditions\n\n## 1. Introduction\nWelcome to MBankingCore API...",
        "updated_at": "2024-01-20T10:30:00Z"
    }
}
```

**Response Errors:**

- `360` - Terms and conditions not found
- `361` - Failed to create terms and conditions
- `362` - Failed to update terms and conditions
- `363` - Failed to delete terms and conditions
- `364` - Failed to retrieve terms and conditions

---

### 2.2 Set Terms & Conditions (Admin/Owner Only)

**Endpoint:** `POST /api/terms-conditions`  
**Access:** Admin/Owner only  
**Description:** Set or update the terms and conditions content

**Request Headers:**

```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "content": "# Terms and Conditions\n\n## 1. Introduction\nWelcome to our service..."
}
```

**Request Fields:**

- `content` (required): Terms and conditions content (supports Markdown)

**Response Success (201 - Created / 200 - Updated):**

```json
{
    "code": 201,
    "message": "Terms and conditions created successfully",
    "data": {
        "content": "# Terms and Conditions\n\n## 1. Introduction\nWelcome to our service...",
        "updated_at": "2024-01-20T10:30:00Z"
    }
}
```

**Response Errors:**

- `253` - Invalid request data or missing content
- `300` - Authentication required
- `752` - Admin privileges required
- `361` - Failed to create terms and conditions
- `362` - Failed to update terms and conditions

---

## 3. Privacy Policy APIs (Config-Based)

**🔄 Updated System**: Privacy Policy now integrated with Config system using key `'privacy-policy'` for simplified management.

### 3.1 Get Privacy Policy (Public)

**Endpoint:** `GET /api/privacy-policy`  
**Access:** Public  
**Description:** Retrieve the current privacy policy content from config

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Privacy policy retrieved successfully",
    "data": {
        "content": "# Privacy Policy\n\n## Data Collection\nWe collect the following data:\n- User information\n- Usage analytics\n\n## Data Usage\nWe use your data to:\n- Provide our services\n- Improve user experience\n\n## Contact\nFor privacy questions, contact us at privacy@example.com",
        "updated_at": "2024-01-20T10:30:00Z"
    }
}
```

**Response Errors:**

- `370` - Privacy policy not found
- `371` - Failed to create privacy policy
- `372` - Failed to update privacy policy
- `373` - Failed to delete privacy policy
- `374` - Failed to retrieve privacy policy

---

### 3.2 Set Privacy Policy (Admin/Owner Only)

**Endpoint:** `POST /api/privacy-policy`  
**Access:** Admin/Owner only  
**Description:** Set or update the privacy policy content

**Request Headers:**

```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "content": "# Privacy Policy\n\n## Data Collection\nWe collect the following data:\n- User information\n- Usage analytics\n\n## Data Usage\nWe use your data to:\n- Provide our services\n- Improve user experience\n\n## Contact\nFor privacy questions, contact us at privacy@example.com"
}
```

**Request Fields:**

- `content` (required): Privacy policy content (supports Markdown)

**Response Success (201 - Created / 200 - Updated):**

```json
{
    "code": 201,
    "message": "Privacy policy created successfully",
    "data": {
        "content": "# Privacy Policy\n\n## Data Collection\nWe collect the following data:\n- User information\n- Usage analytics\n\n## Data Usage\nWe use your data to:\n- Provide our services\n- Improve user experience\n\n## Contact\nFor privacy questions, contact us at privacy@example.com",
        "updated_at": "2024-01-20T10:30:00Z"
    }
}
```

**Response Errors:**

- `253` - Invalid request data or missing content
- `300` - Authentication required
- `752` - Admin privileges required
- `371` - Failed to create privacy policy
- `372` - Failed to update privacy policy

---

## 4. Onboarding Management APIs

### 4.1 Get All Onboardings (Public)

**Endpoint:** `GET /api/onboardings`  
**Access:** Public  
**Description:** Retrieve all onboarding content for the app

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `per_page` (optional): Items per page (default: 10, max: 100)

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Onboardings retrieved successfully",
    "data": {
        "onboardings": [
            {
                "id": 1,
                "image": "https://example.com/onboarding1.jpg",
                "title": "Welcome to Our App",
                "description": "Discover amazing features",
                "is_active": true,
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            }
        ],
        "page": 1,
        "per_page": 10,
        "total": 1,
        "pages": 1
    }
}
```

**Response Errors:**

- `352` - Content not available

---

### 4.2 Get Onboarding by ID (Public)

**Endpoint:** `GET /api/onboardings/{id}`  
**Access:** Public  
**Description:** Retrieve specific onboarding content by ID

**Path Parameters:**

- `id`: Onboarding ID

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Onboarding retrieved successfully",
    "data": {
        "id": 1,
        "image": "https://example.com/onboarding1.jpg",
        "title": "Welcome to Our App",
        "description": "Discover amazing features",
        "is_active": true,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `352` - Content not available
- `352` - Content not available

---

### 4.3 Create Onboarding (Admin Only)

**Endpoint:** `POST /api/onboardings`  
**Access:** Admin/Owner only  
**Description:** Create new onboarding content

**Request Headers:**

```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "image": "https://example.com/onboarding1.jpg",
    "title": "Welcome to Our App",
    "description": "Discover amazing features",
    "is_active": true
}
```

**Response Success (201):**

```json
{
    "code": 201,
    "message": "Onboarding created successfully",
    "data": {
        "id": 1,
        "image": "https://example.com/onboarding1.jpg",
        "title": "Welcome to Our App",
        "description": "Discover amazing features",
        "is_active": true,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `250` - Invalid request data
- `300` - Authentication required
- `651` - Insufficient admin privileges
- `603` - Config creation failed

---

### 4.4 Update Onboarding (Admin Only)

**Endpoint:** `PUT /api/onboardings/{id}`  
**Access:** Admin/Owner only  
**Description:** Update existing onboarding content

**Path Parameters:**

- `id`: Onboarding ID

**Request Headers:**

```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "image": "https://example.com/onboarding1.jpg",
    "title": "Updated Welcome",
    "description": "Updated description",
    "is_active": true
}
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Onboarding updated successfully",
    "data": {
        "id": 1,
        "image": "https://example.com/onboarding1.jpg",
        "title": "Updated Welcome",
        "description": "Updated description",
        "is_active": true,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
    }
}
```

**Response Errors:**

- `250` - Invalid request data
- `300` - Authentication required
- `651` - Insufficient admin privileges
- `352` - Content not available
- `601` - Config update failed

---

### 4.5 Delete Onboarding (Admin Only)

**Endpoint:** `DELETE /api/onboardings/{id}`  
**Access:** Admin/Owner only  
**Description:** Delete onboarding content

**Path Parameters:**

- `id`: Onboarding ID

**Request Headers:**

```
Authorization: Bearer <jwt_token>
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Onboarding deleted successfully"
}
```

**Response Errors:**

- `300` - Authentication required
- `651` - Insufficient admin privileges
- `352` - Content not available
- `604` - Config deletion failed

---

## Banking Authentication Flow

**2-Step Banking Authentication Process:**

1. **Banking Login** - Submit banking credentials and get login_token + OTP
2. **Verify OTP** - Submit login_token + OTP code and receive session tokens
3. **Access Protected Endpoints** - Use JWT tokens in Authorization header

**Important Notes:**

- Each account number must be unique across the system
- If testing with existing data, ensure you use a unique account number
- The system will return a 500 error if trying to create a user with a duplicate account number
- For new user registration, use account numbers that don't exist in the database

---

## 5. Banking Authentication APIs

### 5.1 Banking Login (Step 1)

**Endpoint:** `POST /api/login`  
**Description:** First step of banking authentication - validates credentials and sends OTP  
**Authentication:** None required

**Request Body:**

```json
{
    "name": "John Doe Smith",
    "account_number": "1234567890123456",
    "mother_name": "Jane Doe Smith",
    "phone": "081234567890",
    "pin_atm": "123456",
    "device_info": {
        "device_type": "android",
        "device_id": "android_device_123",
        "device_name": "Samsung Galaxy S21"
    }
}
```

**Request Fields:**

- `name` (required): User's full name as per KTP (min: 8 characters)
- `account_number` (required): Bank account number (min: 8 characters)
- `mother_name` (required): Mother's maiden name (min: 8 characters)
- `phone` (required): Phone number (min: 8 characters)
- `pin_atm` (required): 6-digit ATM PIN (exactly 6 numeric digits)
- `device_info` (required): Device information object

**Validation Rules:**

- Name: minimum 8 characters
- Account number: minimum 8 characters
- Phone: minimum 8 characters
- Mother name: minimum 8 characters
- PIN ATM: exactly 6 numeric digits
- If phone is registered: validates user data and account ownership
- If phone is new: prepares for auto-registration after OTP verification

**Response Success (200):**

```json
{
    "code": 200,
    "message": "OTP sent successfully",
    "data": {
        "login_token": "a1b2c3d4e5f6789012345678901234567890abcdef123456789012345678901234",
        "message": "OTP sent to your phone number",
        "expires_in": 300
    }
}
```

**Response Errors:**

- `400` - Invalid request data / validation errors
- `401` - User information does not match records / Invalid PIN
- `401` - Account number not found for this user
- `250` - Internal server error

---

### 5.2 Banking Login Verify (Step 2)

**Endpoint:** `POST /api/login/verify`  
**Description:** Second step - verify OTP and receive session tokens  
**Authentication:** None required

**Request Body:**

```json
{
    "login_token": "a1b2c3d4e5f6789012345678901234567890abcdef123456789012345678901234",
    "otp_code": "123456"
}
```

**Request Fields:**

- `login_token` (required): Login token from first step
- `otp_code` (required): 6-digit OTP code sent to phone

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Login successful",
    "data": {
        "user": {
            "id": 1,
            "name": "John Doe Smith",
            "phone": "081234567890",
            "mother_name": "Jane Doe Smith",
            "role": "user",
            "avatar": null,
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z"
        },
        "tokens": {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "expires_in": 86400
        },
        "session": {
            "session_id": 1,
            "device_info": {
                "device_type": "android",
                "device_id": "android_device_123",
                "device_name": "Samsung Galaxy S21"
            }
        }
    }
}
```

**Note:** For development purposes, this endpoint currently accepts any valid login_token and returns success. Production implementation will include proper OTP validation.

**Response Errors:**

- `400` - Invalid request data
- `401` - Invalid or expired login token
- `401` - Invalid OTP code
- `250` - Internal server error

---

```json
{
    "code": 200,
    "message": "Login successful",
    "data": {
        "user": {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com",
            "phone": "+1234567890",
            "role": "user",
            "email_verified": false,
            "avatar": null,
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z"
        },
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_in": 86400,
        "session_id": 1,
        "device_info": {
            "device_type": "ios",
            "device_id": "ios_device_456",
            "device_name": "iPhone 13 Pro"
        }
    }
}
```

**Device Conflict (409):**

```json
{
    "code": 409,
    "message": "Device is already logged in. Please logout from this device first or use a different device.",
    "data": {
        "existing_session": {
            "device_type": "ios",
            "device_id": "ios_device_456",
            "device_name": "iPhone 13 Pro",
            "last_activity": "2023-01-01T00:00:00Z"
        }
    }
}
```

**Response Errors:**

- `303` - Invalid email or password
- `253` - Invalid request data
- `409` - Device already logged in
- `250` - Internal server error

---

### 5.3 Refresh Token

**Endpoint:** `POST /api/refresh`  
**Description:** Refresh access token using refresh token  
**Authentication:** None required

**Request Body:**

```json
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Token refreshed successfully",
    "data": {
        "tokens": {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "expires_in": 86400
        },
        "session": {
            "session_id": 1
        }
    }
}
```

**Response Errors:**

- `401` - Invalid or expired refresh token
- `253` - Invalid request data

---

## 6. User Profile APIs

### 6.1 Get User Profile

**Endpoint:** `GET /api/profile`  
**Description:** Get current user's profile information  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "User retrieved successfully",
    "data": {
        "id": 1,
        "name": "John Doe Smith",
        "phone": "081234567890",
        "mother_name": "Jane Doe Smith",
        "role": "user",
        "avatar": null,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `300` - Unauthorized access
- `400` - User not found

---

### 6.2 Update User Profile

**Endpoint:** `PUT /api/profile`  
**Description:** Update current user's profile information  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "name": "John Updated Name",
    "avatar": "base64_encoded_image"
}
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Profile updated successfully",
    "data": {
        "id": 1,
        "name": "John Updated Name",
        "phone": "081234567890",
        "mother_name": "Jane Doe Smith",
        "role": "user",
        "avatar": "https://example.com/avatar.jpg",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
    }
}
```

**Response Errors:**

- `300` - Unauthorized access
- `400` - Invalid request data
- `250` - Internal server error

---

### 6.3 Change PIN ATM

**Endpoint:** `PUT /api/change-pin`  
**Description:** Change user's PIN ATM and invalidate all sessions  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "old_pin": "123456",
    "new_pin": "654321"
}
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "PIN changed successfully. Please login again with your new PIN",
    "data": null
}
```

**Response Errors:**

- `300` - Unauthorized access
- `400` - Invalid current PIN
- `253` - Invalid request data
- `250` - Internal server error

---

## 7. Session Management APIs

### 7.1 Get Active Sessions

**Endpoint:** `GET /api/sessions`  
**Description:** Get all active sessions for current user across all devices  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Active sessions retrieved successfully",
    "data": {
        "sessions": [
            {
                "id": 1,
                "device_type": "android",
                "device_id": "device123",
                "device_name": "Samsung Galaxy S23",
                "ip_address": "192.168.1.1",
                "user_agent": "MBankingApp/1.0.0",
                "created_at": "2023-01-01T00:00:00Z",
                "last_active": "2023-01-01T12:00:00Z",
                "is_current": true
            },
            {
                "id": 2,
                "device_type": "ios",
                "device_id": "device456",
                "device_name": "iPhone 14",
                "ip_address": "192.168.1.2",
                "user_agent": "MBankingApp/1.0.0",
                "created_at": "2023-01-01T10:00:00Z",
                "last_active": "2023-01-01T11:00:00Z",
                "is_current": false
            }
        ],
        "total": 2
    }
}
```

**Response Errors:**

- `300` - Unauthorized access
- `250` - Internal server error

---

### 7.2 Logout

**Endpoint:** `POST /api/logout`  
**Description:** Logout from current device or all devices  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body (Optional):**

```json
{
    "logout_all": false
}
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Logged out successfully",
    "data": null
}
```

**Response Errors:**

- `300` - Unauthorized access
- `250` - Internal server error

---

### 7.3 Logout Other Sessions

**Endpoint:** `POST /api/logout-others`  
**Description:** Logout from all other sessions except current device  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Logged out from other devices successfully",
    "data": null
}
```

**Response Errors:**

- `300` - Unauthorized access
- `250` - Internal server error

---

## 8. Bank Account Management APIs

### 7.1 Get Bank Accounts

**Endpoint:** `GET /api/bank-accounts`  
**Description:** Get all bank accounts for the authenticated user  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Bank accounts retrieved successfully",
    "data": {
        "accounts": [
            {
                "id": 1,
                "account_number": "1234567890123456",
                "account_name": "John Doe Smith",
                "bank_name": "Bank Central Asia",
                "bank_code": "014",
                "account_type": "saving",
                "is_active": true,
                "is_primary": true,
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            },
            {
                "id": 2,
                "account_number": "9876543210987654",
                "account_name": "John Doe Smith",
                "bank_name": "Bank Mandiri",
                "bank_code": "008",
                "account_type": "checking",
                "is_active": true,
                "is_primary": false,
                "created_at": "2023-01-02T00:00:00Z",
                "updated_at": "2023-01-02T00:00:00Z"
            }
        ],
        "total": 2
    }
}
```

**Response Errors:**

- `300` - Unauthorized access
- `250` - Internal server error

---

### 7.2 Create Bank Account

**Endpoint:** `POST /api/bank-accounts`  
**Description:** Create a new bank account for the authenticated user  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "account_number": "1111222233334444",
    "account_name": "John Doe Smith",
    "bank_name": "Bank Negara Indonesia",
    "bank_code": "009",
    "account_type": "saving",
    "is_primary": false
}
```

**Request Fields:**

- `account_number` (required): Bank account number (min: 8, max: 20 characters)
- `account_name` (required): Account name (min: 3, max: 100 characters)
- `bank_name` (optional): Bank institution name (max: 100 characters)
- `bank_code` (optional): Bank code (max: 10 characters)
- `account_type` (optional): Account type (max: 20 characters)
- `is_primary` (optional): Set as primary account (boolean)

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Bank account created successfully",
    "data": {
        "id": 3,
        "account_number": "1111222233334444",
        "account_name": "John Doe Smith",
        "bank_name": "Bank Negara Indonesia",
        "bank_code": "009",
        "account_type": "saving",
        "is_active": true,
        "is_primary": false,
        "created_at": "2023-01-03T00:00:00Z",
        "updated_at": "2023-01-03T00:00:00Z"
    }
}
```

**Response Errors:**

- `300` - Unauthorized access
- `400` - Invalid request data
- `409` - Account number already exists for this user
- `250` - Internal server error

---

### 7.3 Update Bank Account

**Endpoint:** `PUT /api/bank-accounts/:id`  
**Description:** Update an existing bank account  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "account_name": "Updated Account Name",
    "bank_name": "Updated Bank Name",
    "bank_code": "010",
    "account_type": "checking"
}
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Bank account updated successfully",
    "data": {
        "id": 3,
        "account_number": "1111222233334444",
        "account_name": "Updated Account Name",
        "bank_name": "Updated Bank Name",
        "bank_code": "010",
        "account_type": "checking",
        "is_active": true,
        "is_primary": false,
        "created_at": "2023-01-03T00:00:00Z",
        "updated_at": "2023-01-03T12:00:00Z"
    }
}
```

**Response Errors:**

- `300` - Unauthorized access
- `400` - Invalid request data
- `404` - Bank account not found
- `250` - Internal server error

---

### 7.4 Delete Bank Account

**Endpoint:** `DELETE /api/bank-accounts/:id`  
**Description:** Delete (soft delete) a bank account  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Bank account deleted successfully"
}
```

**Response Errors:**

- `300` - Unauthorized access
- `404` - Bank account not found
- `400` - Cannot delete primary account (set another as primary first)
- `250` - Internal server error

---

### 7.5 Set Primary Account

**Endpoint:** `PUT /api/bank-accounts/:id/primary`  
**Description:** Set a bank account as the primary account  
**Authentication:** Bearer token required

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Primary account updated successfully",
    "data": {
        "id": 3,
        "account_number": "1111222233334444",
        "account_name": "John Doe Smith",
        "bank_name": "Bank Negara Indonesia",
        "bank_code": "009",
        "account_type": "saving",
        "is_active": true,
        "is_primary": true,
        "created_at": "2023-01-03T00:00:00Z",
        "updated_at": "2023-01-03T12:00:00Z"
    }
}
```

**Response Errors:**

- `300` - Unauthorized access
- `404` - Bank account not found
- `250` - Internal server error

---

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "name": "John Smith",
    "phone": "+1234567891"
}
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "User updated successfully",
    "data": {
        "id": 1,
        "name": "John Smith",
        "email": "john@example.com",
        "phone": "+1234567891",
        "role": "user",
        "email_verified": false,
        "avatar": null,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
    }
}
```

**Response Errors:**

- `253` - Invalid request data
- `300` - Unauthorized access
- `403` - Failed to update user

---

## 6. User Management APIs

### 6.1 List All Users (Admin Only)

**Endpoint:** `GET /api/users`  
**Access:** Admin/Owner only  
**Description:** Retrieve all users with pagination

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `per_page` (optional): Items per page (default: 10, max: 100)

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Users retrieved successfully",
    "data": {
        "users": [
            {
                "id": 1,
                "name": "John Doe",
                "email": "john@example.com",
                "phone": "+1234567890",
                "role": "user",
                "email_verified": true,
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            }
        ],
        "page": 1,
        "per_page": 10,
        "total": 1,
        "pages": 1
    }
}
```

**Response Errors:**

- `401` - Unauthorized
- `403` - Forbidden
- `500` - Failed to retrieve users

---

### 6.2 Get User by ID (Admin Only)

**Endpoint:** `GET /api/users/{id}`  
**Access:** Admin/Owner only  
**Description:** Retrieve specific user by ID

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Path Parameters:**

- `id`: User ID

**Response Success (200):**

```json
{
    "code": 200,
    "message": "User retrieved successfully",
    "data": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "phone": "+1234567890",
        "role": "user",
        "email_verified": true,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `401` - Unauthorized
- `403` - Forbidden
- `404` - User not found
- `500` - Failed to retrieve user

---

### 6.3 Delete User (Admin Only)

**Endpoint:** `DELETE /api/users/{id}`  
**Access:** Admin/Owner only  
**Description:** Delete user by ID

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Path Parameters:**

- `id`: User ID

**Response Success (200):**

```json
{
    "code": 200,
    "message": "User deleted successfully"
}
```

**Response Errors:**

- `401` - Unauthorized
- `403` - Forbidden
- `404` - User not found
- `500` - Failed to delete user

---

### 6.5 Create User (Owner Only)

**Endpoint:** `POST /api/users`  
**Access:** Owner only  
**Description:** Create new user with any role (owner can set any role)

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "name": "New User",
    "email": "newuser@example.com",
    "password": "hashed_password_from_client",
    "phone": "+1234567890",
    "role": "admin",
    "email_verified": true
}
```

**Request Fields:**

- `name` (required): User's full name
- `email` (required): User's email address
- `password` (required): SHA256 hashed password from client
- `phone` (optional): User's phone number
- `role` (required): User role ("user", "admin", "owner")
- `email_verified` (optional): Email verification status (default: false)

**Response Success (201):**

```json
{
    "code": 201,
    "message": "User created successfully",
    "data": {
        "id": 1,
        "name": "New User",
        "email": "newuser@example.com",
        "phone": "+1234567890",
        "role": "admin",
        "email_verified": true,
        "avatar": null,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `250` - Invalid request data
- `300` - Authentication required
- `753` - Owner privileges required
- `401` - User already exists
- `402` - User creation failed

---

### 6.6 Update User (Owner Only)

**Endpoint:** `PUT /api/users/{id}`  
**Access:** Owner only  
**Description:** Update user information including role changes (owner can change any user's role)

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Path Parameters:**

- `id`: User ID

**Request Body:**

```json
{
    "name": "Updated User Name",
    "email": "updated@example.com",
    "phone": "+1234567891",
    "role": "admin",
    "email_verified": true
}
```

**Request Fields:**

- `name` (optional): User's full name
- `email` (optional): User's email address
- `phone` (optional): User's phone number
- `role` (optional): User role ("user", "admin", "owner")
- `email_verified` (optional): Email verification status

**Response Success (200):**

```json
{
    "code": 200,
    "message": "User updated successfully",
    "data": {
        "id": 1,
        "name": "Updated User Name",
        "email": "updated@example.com",
        "phone": "+1234567891",
        "role": "admin",
        "email_verified": true,
        "avatar": null,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
    }
}
```

**Response Errors:**

- `250` - Invalid request data
- `300` - Authentication required
- `753` - Owner privileges required
- `400` - User not found
- `403` - User update failed

---

## 9. Article Management APIs

### 7.1 Get All Articles

**Endpoint:** `GET /api/articles`  
**Access:** Authenticated users  
**Description:** Retrieve all articles with pagination

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `per_page` (optional): Items per page (default: 10, max: 100)

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Articles retrieved successfully",
    "data": {
        "articles": [
            {
                "id": 1,
                "title": "Sample Article",
                "image": "https://example.com/article.jpg",
                "content": "Article content here...",
                "is_active": true,
                "user_id": 1,
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            }
        ],
        "page": 1,
        "per_page": 10,
        "total": 1,
        "pages": 1
    }
}
```

**Response Errors:**

- `401` - Unauthorized
- `500` - Failed to retrieve articles

---

### 7.2 Get Article by ID

**Endpoint:** `GET /api/articles/{id}`  
**Access:** Authenticated users  
**Description:** Retrieve specific article by ID

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Path Parameters:**

- `id`: Article ID

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Article retrieved successfully",
    "data": {
        "id": 1,
        "title": "Sample Article",
        "image": "https://example.com/article.jpg",
        "content": "Article content here...",
        "is_active": true,
        "user_id": 1,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `401` - Unauthorized
- `404` - Article not found
- `500` - Failed to retrieve article

---

### 7.3 Create Article (Admin Only)

**Endpoint:** `POST /api/articles`  
**Access:** Admin/Owner only  
**Description:** Create new article

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "title": "New Article Title",
    "image": "https://example.com/image.jpg",
    "content": "Article content here...",
    "is_active": true
}
```

**Response Success (201):**

```json
{
    "code": 201,
    "message": "Article created successfully",
    "data": {
        "id": 1,
        "title": "New Article Title",
        "image": "https://example.com/image.jpg",
        "content": "Article content here...",
        "is_active": true,
        "user_id": 1,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `504` - Invalid article data
- `300` - Authentication required
- `651` - Insufficient admin privileges
- `501` - Article creation failed

---

### 7.4 Update Article (Admin Only)

**Endpoint:** `PUT /api/articles/{id}`  
**Access:** Admin/Owner only  
**Description:** Update existing article

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Path Parameters:**

- `id`: Article ID

**Request Body:**

```json
{
    "title": "Updated Article Title",
    "image": "https://example.com/updated-image.jpg",
    "content": "Updated content...",
    "is_active": false
}
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Article updated successfully",
    "data": {
        "id": 1,
        "title": "Updated Article Title",
        "image": "https://example.com/updated-image.jpg",
        "content": "Updated content...",
        "is_active": false,
        "user_id": 1,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
    }
}
```

**Response Errors:**

- `504` - Invalid article data
- `300` - Authentication required
- `651` - Insufficient admin privileges
- `500` - Article not found
- `502` - Article update failed

---

### 7.5 Delete Article (Admin Only)

**Endpoint:** `DELETE /api/articles/{id}`  
**Access:** Admin/Owner only  
**Description:** Delete article by ID

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Path Parameters:**

- `id`: Article ID

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Article deleted successfully"
}
```

**Response Errors:**

- `401` - Unauthorized
- `403` - Forbidden
- `404` - Article not found
- `500` - Failed to delete article

---

### 7.6 Get My Articles

**Endpoint:** `GET /api/my-articles`  
**Access:** Authenticated users  
**Description:** Retrieve articles created by the current authenticated user

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `per_page` (optional): Items per page (default: 10, max: 100)

**Response Success (200):**

```json
{
    "code": 200,
    "message": "My articles retrieved successfully",
    "data": {
        "articles": [
            {
                "id": 1,
                "title": "My Article",
                "image": "https://example.com/my-article.jpg",
                "content": "My article content here...",
                "is_active": true,
                "user_id": 1,
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            }
        ],
        "page": 1,
        "per_page": 10,
        "total": 1,
        "pages": 1
    }
}
```

**Response Errors:**

- `300` - Authentication required
- `500` - Failed to retrieve articles

---

## 10. Photo Management APIs

### 8.1 Get All Photos

**Endpoint:** `GET /api/photos`  
**Access:** Authenticated users  
**Description:** Retrieve all photos with pagination

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `per_page` (optional): Items per page (default: 10, max: 100)

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Photos retrieved successfully",
    "data": {
        "photos": [
            {
                "id": 1,
                "image": "https://example.com/photo1.jpg",
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            }
        ],
        "page": 1,
        "per_page": 10,
        "total": 1,
        "pages": 1
    }
}
```

**Response Errors:**

- `401` - Unauthorized
- `500` - Failed to retrieve photos

---

### 8.2 Get Photo by ID

**Endpoint:** `GET /api/photos/{id}`  
**Access:** Authenticated users  
**Description:** Retrieve specific photo by ID

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Path Parameters:**

- `id`: Photo ID

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Photo retrieved successfully",
    "data": {
        "id": 1,
        "image": "https://example.com/photo1.jpg",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `401` - Unauthorized
- `404` - Photo not found
- `500` - Failed to retrieve photo

---

### 8.3 Create Photo (Admin Only)

**Endpoint:** `POST /api/photos`  
**Access:** Admin/Owner only  
**Description:** Upload new photo

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "image": "https://example.com/new-photo.jpg"
}
```

**Response Success (201):**

```json
{
    "code": 201,
    "message": "Photo created successfully",
    "data": {
        "id": 1,
        "image": "https://example.com/new-photo.jpg",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `553` - Invalid photo format
- `300` - Authentication required
- `651` - Insufficient admin privileges
- `551` - Photo upload failed

---

### 8.4 Update Photo (Admin Only)

**Endpoint:** `PUT /api/photos/{id}`  
**Access:** Admin/Owner only  
**Description:** Update existing photo

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Path Parameters:**

- `id`: Photo ID

**Request Body:**

```json
{
    "image": "https://example.com/updated-photo.jpg"
}
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Photo updated successfully",
    "data": {
        "id": 1,
        "image": "https://example.com/updated-photo.jpg",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
    }
}
```

**Response Errors:**

- `553` - Invalid photo format
- `300` - Authentication required
- `651` - Insufficient admin privileges
- `550` - Photo not found
- `555` - Photo processing failed

---

### 8.5 Delete Photo (Admin Only)

**Endpoint:** `DELETE /api/photos/{id}`  
**Access:** Admin/Owner only  
**Description:** Delete photo by ID

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Path Parameters:**

- `id`: Photo ID

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Photo deleted successfully"
}
```

**Response Errors:**

- `401` - Unauthorized
- `403` - Forbidden
- `404` - Photo not found
- `500` - Failed to delete photo

---

## 11. Configuration APIs

### 9.1 Get Config by Key

**Endpoint:** `GET /api/config/{key}`  
**Access:** Authenticated users  
**Description:** Retrieve configuration value by key

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Path Parameters:**

- `key`: Configuration key

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Configuration retrieved successfully",
    "data": {
        "key": "app_version",
        "value": "1.0.0",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `401` - Unauthorized
- `404` - Configuration not found
- `500` - Failed to retrieve configuration

---

### 9.2 Set Config (Admin Only)

**Endpoint:** `POST /api/config`  
**Access:** Admin/Owner only  
**Description:** Set or update configuration value

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "key": "app_version",
    "value": "1.0.1"
}
```

**Response Success (200/201):**

```json
{
    "code": 201,
    "message": "Configuration set successfully",
    "data": {
        "key": "app_version",
        "value": "1.0.1",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
    }
}
```

**Response Errors:**

- `602` - Invalid config data
- `300` - Authentication required
- `651` - Insufficient admin privileges
- `603` - Config creation failed

---

### 9.3 Get All Configs (Admin Only)

**Endpoint:** `GET /api/configs`  
**Access:** Admin/Owner only  
**Description:** Retrieve all configuration keys and values

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Configurations retrieved successfully",
    "data": {
        "configs": [
            {
                "key": "app_version",
                "value": "1.0.1",
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            },
            {
                "key": "tnc",
                "value": "# Terms and Conditions\n\n...",
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            }
        ]
    }
}
```

**Response Errors:**

- `401` - Unauthorized
- `403` - Forbidden
- `500` - Failed to retrieve configurations

---

### 9.4 Delete Config (Admin Only)

**Endpoint:** `DELETE /api/config/{key}`  
**Access:** Admin/Owner only  
**Description:** Delete configuration by key

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Path Parameters:**

- `key`: Configuration key

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Configuration deleted successfully"
}
```

**Response Errors:**

- `300` - Authentication required
- `651` - Insufficient admin privileges
- `600` - Config not found
- `604` - Config deletion failed

---

## 12. Admin Management APIs

**🔧 Admin System**: Complete admin authentication and CRUD management system with role-based access control.

### 12.1 Admin Login

**Endpoint:** `POST /api/admin/login`  
**Access:** Public (Admin Credentials Required)  
**Description:** Authenticate admin users and get access token

**Request Body:**

```json
{
    "email": "admin@mbankingcore.com",
    "password": "admin123"
}
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Admin login successful",
    "data": {
        "admin": {
            "id": 1,
            "name": "Super Admin",
            "email": "admin@mbankingcore.com",
            "role": "super",
            "status": 1,
            "avatar": "",
            "last_login": "2025-07-30T10:05:20.484234+07:00",
            "created_at": "2025-07-30T09:58:58.74129+07:00",
            "updated_at": "2025-07-30T10:05:20.48451+07:00"
        },
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_in": 86400
    }
}
```

**Response Errors:**

- `303` - Invalid email or password
- `750` - Admin account is inactive
- `751` - Admin account is blocked

---

### 12.2 Admin Logout

**Endpoint:** `POST /api/admin/logout`  
**Access:** Admin Authentication Required  
**Description:** Logout admin session and invalidate token

**Request Headers:**

```
Authorization: Bearer <admin_access_token>
```

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Admin logout successful",
    "data": null
}
```

**Response Errors:**

- `301` - Invalid or expired admin token
- `304` - Authorization token required

---

### 12.3 Get All Admins

**Endpoint:** `GET /api/admin/admins`  
**Access:** Admin Authentication Required  
**Description:** Retrieve paginated list of all admin accounts

**Request Headers:**

```
Authorization: Bearer <admin_access_token>
```

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `per_page` (optional): Items per page (default: 10)

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Admins retrieved successfully",
    "data": {
        "admins": [
            {
                "id": 1,
                "name": "Super Admin",
                "email": "admin@mbankingcore.com",
                "role": "super",
                "status": 1,
                "avatar": "",
                "last_login": "2025-07-30T10:05:20.484234+07:00",
                "created_at": "2025-07-30T09:58:58.74129+07:00",
                "updated_at": "2025-07-30T10:05:20.48451+07:00"
            }
        ],
        "total": 1,
        "page": 1,
        "per_page": 10
    }
}
```

**Response Errors:**

- `301` - Invalid or expired admin token
- `750` - Access forbidden - insufficient permissions

---

### 12.4 Get Admin by ID

**Endpoint:** `GET /api/admin/admins/:id`  
**Access:** Admin Authentication Required  
**Description:** Retrieve specific admin account details

**Request Headers:**

```
Authorization: Bearer <admin_access_token>
```

**Path Parameters:**

- `id`: Admin ID

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Admin retrieved successfully",
    "data": {
        "id": 1,
        "name": "Super Admin",
        "email": "admin@mbankingcore.com",
        "role": "super",
        "status": 1,
        "avatar": "",
        "last_login": "2025-07-30T10:05:20.484234+07:00",
        "created_at": "2025-07-30T09:58:58.74129+07:00",
        "updated_at": "2025-07-30T10:05:20.48451+07:00"
    }
}
```

**Response Errors:**

- `301` - Invalid or expired admin token
- `400` - Admin not found
- `750` - Access forbidden - insufficient permissions

---

### 12.5 Create Admin

**Endpoint:** `POST /api/admin/admins`  
**Access:** Super Admin Authentication Required  
**Description:** Create new admin account (Super Admin only)

**Request Headers:**

```
Authorization: Bearer <admin_access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "name": "Test Admin",
    "email": "test@mbankingcore.com",
    "password": "password123",
    "role": "admin"
}
```

**Field Validations:**

- `name`: Required, minimum 3 characters
- `email`: Required, valid email format, unique
- `password`: Required, minimum 6 characters
- `role`: Required, must be "admin" or "super"

**Response Success (201):**

```json
{
    "code": 201,
    "message": "Admin created successfully",
    "data": {
        "id": 2,
        "name": "Test Admin",
        "email": "test@mbankingcore.com",
        "role": "admin",
        "status": 1,
        "avatar": "",
        "created_at": "2025-07-30T10:01:50.720778+07:00",
        "updated_at": "2025-07-30T10:01:50.720778+07:00"
    }
}
```

**Response Errors:**

- `301` - Invalid or expired admin token
- `306` - Email already exists
- `252` - Validation failed
- `753` - Super Admin privileges required

---

### 12.6 Update Admin

**Endpoint:** `PUT /api/admin/admins/:id`  
**Access:** Super Admin Authentication Required  
**Description:** Update admin account details (Super Admin only)

**Request Headers:**

```
Authorization: Bearer <admin_access_token>
Content-Type: application/json
```

**Path Parameters:**

- `id`: Admin ID to update

**Request Body:**

```json
{
    "name": "Updated Test Admin",
    "email": "test@mbankingcore.com",
    "role": "admin",
    "status": 1,
    "password": "newpassword123"
}
```

**Field Notes:**

- `password`: Optional, only include if changing password
- `status`: 1 = Active, 0 = Inactive, 2 = Blocked

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Admin updated successfully",
    "data": {
        "id": 2,
        "name": "Updated Test Admin",
        "email": "test@mbankingcore.com",
        "role": "admin",
        "status": 1,
        "avatar": "",
        "last_login": null,
        "created_at": "2025-07-30T10:01:50.720778+07:00",
        "updated_at": "2025-07-30T10:15:30.123456+07:00"
    }
}
```

**Response Errors:**

- `301` - Invalid or expired admin token
- `400` - Admin not found
- `306` - Email already exists (if changing email)
- `252` - Validation failed
- `753` - Super Admin privileges required

---

### 12.7 Delete Admin

**Endpoint:** `DELETE /api/admin/admins/:id`  
**Access:** Super Admin Authentication Required  
**Description:** Delete admin account (Super Admin only, cannot delete self)

**Request Headers:**

```
Authorization: Bearer <admin_access_token>
```

**Path Parameters:**

- `id`: Admin ID to delete

**Response Success (200):**

```json
{
    "code": 200,
    "message": "Admin deleted successfully",
    "data": null
}
```

**Response Errors:**

- `301` - Invalid or expired admin token
- `400` - Admin not found
- `750` - Cannot delete yourself
- `753` - Super Admin privileges required

---

## Admin System Features

### Role-Based Access Control

**Super Admin:**
- Full access to all admin operations
- Can create, read, update, delete other admins
- Cannot delete themselves

**Admin:**
- Can view admin list and details
- Cannot manage other admin accounts
- Standard admin operations only

### Security Features

- JWT-based authentication with 24-hour expiration
- Password encryption using bcrypt
- Role-based middleware protection
- Self-deletion prevention
- Email uniqueness validation

### Admin Status Management

- **Active (1):** Admin can login and perform operations
- **Inactive (0):** Admin cannot login
- **Blocked (2):** Admin account is blocked

---

## Common Response Codes

The API uses a hierarchical error code system where each feature has its own unique range of 50 numbers. This provides better error categorization and debugging capabilities.

### Success Codes

- `200` - General success and default success response

### Error Code Ranges

#### General/System Errors (250-299)

- `250` - Internal server error
- `251` - Database operation failed
- `252` - Validation failed
- `253` - Invalid request data
- `254` - Resource not found

#### Authentication & Authorization (300-349) - Public Endpoints

- `300` - Unauthorized access
- `301` - Invalid or malformed token
- `302` - Token has expired
- `303` - Invalid email or password
- `304` - Authorization token required
- `305` - Invalid email format
- `306` - Email already exists
- `307` - Registration failed
- `308` - Login failed

- `309` - Token refresh failed

#### Public Content (350-399) - Terms, Privacy Policy

**Terms & Conditions (360-369):**

- `360` - Terms and conditions not found

- `361` - Failed to create terms and conditions
- `362` - Failed to update terms and conditions
- `363` - Failed to delete terms and conditions
- `364` - Failed to retrieve terms and conditions

**Privacy Policy (370-379):**

- `370` - Privacy policy not found
- `371` - Failed to create privacy policy
- `372` - Failed to update privacy policy
- `373` - Failed to delete privacy policy
- `374` - Failed to retrieve privacy policy

#### User Management (400-449) - Protected Endpoints

- `400` - User not found
- `401` - Invalid user ID
- `402` - User registration failed
- `403` - User update failed
- `404` - User deletion failed
- `405` - Failed to retrieve user
- `406` - Failed to retrieve users list

#### User Profile (450-499) - Protected Endpoints

- `450` - Profile not found
- `451` - Profile update failed
- `452` - Failed to retrieve profile
- `453` - Failed to change password
- `454` - Current password is invalid

#### Configuration (600-649) - Admin Endpoints

- `600` - Configuration not found
- `601` - Failed to create configuration
- `602` - Failed to update configuration
- `603` - Failed to delete configuration
- `604` - Failed to retrieve configuration
- `605` - Invalid configuration key
- `606` - Invalid configuration value

#### Admin Management (650-699) - Admin Endpoints

- `650` - Admin not found
- `651` - Failed to create admin
- `652` - Failed to update admin
- `653` - Failed to delete admin
- `654` - Failed to retrieve admin
- `655` - Invalid admin role
- `656` - Admin email already exists
- `657` - Cannot delete yourself
- `658` - Admin account inactive
- `659` - Admin account blocked

#### Permissions (750-799)

- `750` - Access forbidden - insufficient permissions
- `751` - Insufficient permissions to perform this action
- `752` - Admin privileges required
- `753` - Owner privileges required

## Error Response Format

All error responses follow this standard format:

```json
{
    "code": 250,
    "message": "Detailed error message explaining what went wrong",
    "errors": [
        {
            "field": "email",
            "message": "Email is required"
        }

    ]
}
```

## Rate Limiting

The API implements rate limiting to ensure fair usage:

- **Authentication endpoints**: 5 requests per minute per IP
- **General endpoints**: 100 requests per minute per authenticated user
- **Admin endpoints**: 1000 requests per minute per admin user

When rate limits are exceeded, the API returns a `429 Too Many Requests` status with retry information in headers.

## Versioning

The current API version is `v1`. All endpoints are prefixed with `/api` to maintain consistency and enable future versioning.

## Support

For technical support or questions about the API, please contact the development team or refer to the project documentation.

---

**Last Updated:** July 30, 2025  
**API Version:** 1.0  
**Documentation Version:** 4.0
