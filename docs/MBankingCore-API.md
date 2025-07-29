# MBankingCore API Documentation

Dokumentasi lengkap untuk RESTful API MBankingCore dengan JWT Authentication dan Multi-Device Session Management.

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

### 🔐 Authentication APIs (3 endpoints)
- **[Authentication](#5-authentication-apis)** - Registrasi, login, refresh token

### 🛡️ Protected APIs (8 endpoints) 
- **[User Profile Management](#6-user-profile-apis)** - Manajemen profil user (2 endpoints)
- **[Article Management](#7-article-management-apis)** - Operasi CRUD artikel (5 endpoints)
- **[Photo Management](#8-photo-management-apis)** - Sistem manajemen foto (4 endpoints)
- **[Configuration APIs](#9-configuration-apis)** - Read config (1 endpoint)

### 👑 Admin APIs (13 endpoints)
- **[Admin Article Management](#10-admin-article-management)** - Create artikel (1 endpoint)
- **[Admin Onboarding Management](#11-admin-onboarding-management)** - CRUD onboarding (3 endpoints)
- **[Admin Photo Management](#12-admin-photo-management)** - Create photo (1 endpoint)
- **[Admin User Management](#13-admin-user-management)** - Manajemen user (4 endpoints)
- **[Admin Configuration](#14-admin-configuration-apis)** - Full config management (3 endpoints)
- **[Admin Terms & Conditions](#15-admin-terms-conditions)** - Set T&C (1 endpoint)
- **[Admin Privacy Policy](#16-admin-privacy-policy)** - Set Privacy Policy (1 endpoint)

### 👨‍💼 Owner-Only APIs (2 endpoints)
- **[Owner User Management](#17-owner-user-management)** - Create & update users dengan roles (2 endpoints)

**Total: 36 Active Endpoints**

---

# 🔌 API Endpoints

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
        "content": "# Terms and Conditions\n\n## 1. Introduction\nWelcome to MBX Backend API...",
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

## Authentication Flow

1. **Register User** - Submit user registration data with device information
2. **Login User** - Authenticate and receive JWT tokens  
3. **Access Protected Endpoint** - Use JWT token in Authorization header

---

## 5. Authentication APIs

### 5.1 Register User

**Endpoint:** `POST /api/register`  
**Description:** Register new user with device information  
**Authentication:** None required

**Request Body:**

```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "hashed_password_from_client",
    "phone": "+1234567890",
    "role": "user",
    "provider": "email",
    "provider_id": "",
    "device_info": {
        "device_type": "android",
        "device_id": "android_device_123",
        "device_name": "Samsung Galaxy S21",
        "user_agent": "MBankingCore-Android-App/1.0.0"
    }
}
```

**Request Fields:**

- `name` (required): User's full name
- `email` (required): User's email address
- `password` (required): SHA256 hashed password from client
- `phone` (optional): User's phone number
- `provider` (required): Authentication provider ("email", "google", "apple", "facebook")
- `device_info` (required): Device information object

**Response Success (200):**

```json
{
    "code": 200,
    "message": "User registered successfully",
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
            "device_type": "android",
            "device_id": "android_device_123",
            "device_name": "Samsung Galaxy S21",
            "user_agent": "MBankingCore-Android-App/1.0.0"
        }
    }
}
```

**Response Errors:**

- `306` - Email already exists
- `253` - Invalid request data
- `250` - Internal server error

---

### 5.2 Login User

**Endpoint:** `POST /api/login`  
**Description:** Authenticate user and create session  
**Authentication:** None required

**Request Body:**

```json
{
    "email": "john@example.com",
    "password": "hashed_password_from_client",
    "provider": "email",
    "provider_id": "",
    "device_info": {
        "device_type": "ios",
        "device_id": "ios_device_456",
        "device_name": "iPhone 13 Pro",
        "user_agent": "MBankingCore-iOS-App/1.0.0"
    }
}
```

**Response Success (200):**

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
            "device_name": "iPhone 13 Pro",
            "user_agent": "MBankingCore-iOS-App/1.0.0"
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
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_in": 86400,
        "session_id": 1
    }
}
```

**Response Errors:**

- `401` - Invalid or expired refresh token
- `253` - Invalid request data

---

### 5.4 Get User Profile

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
        "name": "John Doe",
        "email": "john@example.com",
        "phone": "+1234567890",
        "role": "user",
        "email_verified": false,
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

### 5.5 Update User Profile

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

### 6.2 List Admin Users (Admin Only)

**Endpoint:** `GET /api/admin/users`  
**Access:** Admin/Owner only  
**Description:** Retrieve all admin and owner users with pagination

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
                "name": "Admin User",
                "email": "admin@example.com",
                "phone": "+1234567890",
                "role": "admin",
                "email_verified": true,
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            },
            {
                "id": 2,
                "name": "Owner User", 
                "email": "owner@example.com",
                "phone": "+1234567891",
                "role": "owner",
                "email_verified": true,
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            }
        ],
        "page": 1,
        "per_page": 10,
        "total": 2,
        "pages": 1
    }
}
```

**Response Errors:**

- `300` - Authentication required
- `752` - Admin privileges required
- `405` - Failed to retrieve users

---

### 6.3 Get User by ID (Admin Only)

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

### 6.4 Delete User (Admin Only)

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

## 7. Article Management APIs

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

## 8. Photo Gallery APIs

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

## 9. Configuration APIs

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

**Endpoint:** `GET /api/admin/configs`  
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

**Last Updated:** July 24, 2025  
**API Version:** 1.0  
**Documentation Version:** 3.1
