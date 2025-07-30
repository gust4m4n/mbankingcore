# 📮 Postman Collection Update - Admin Management APIs

## 🎉 What's New

The Postman collection has been updated to include comprehensive admin management functionality with **7 new admin endpoints**.

### ✅ **New Admin Management Section**

**Section:** `👥 Admin Management`

| Endpoint | Method | Path | Description |
|----------|--------|------|-------------|
| Admin Login | `POST` | `/api/admin/login` | Authenticate admin users |
| Admin Logout | `POST` | `/api/admin/logout` | Logout admin session |
| Get All Admins | `GET` | `/api/admin/admins` | List all admin accounts (paginated) |
| Get Admin by ID | `GET` | `/api/admin/admins/:id` | Get specific admin details |
| Create Admin | `POST` | `/api/admin/admins` | Create new admin account |
| Update Admin | `PUT` | `/api/admin/admins/:id` | Update admin information |
| Delete Admin | `DELETE` | `/api/admin/admins/:id` | Delete admin account |

### 🔧 **Updated Environment Variables**

**New Admin Variables Added:**
```json
{
  "admin_token": "",              // Admin JWT token (auto-saved)
  "admin_id": "1",               // Current admin ID
  "admin_email": "admin@mbankingcore.com",     // Default admin email
  "admin_password": "admin123",   // Default admin password
  "admin_role": "super",         // Admin role
  "new_admin_name": "Test Admin", // For creating new admins
  "new_admin_email": "test@mbankingcore.com",
  "new_admin_password": "password123",
  "new_admin_role": "admin",
  "updated_admin_name": "Updated Test Admin",
  "updated_admin_email": "test@mbankingcore.com",
  "updated_admin_role": "admin",
  "updated_admin_status": "1",
  "created_admin_id": ""         // Auto-saved from create operations
}
```

### 🎯 **Testing Features**

**Automated Testing:**
- ✅ Response status validation
- ✅ Response structure verification
- ✅ Automatic token management
- ✅ Admin ID tracking for CRUD operations
- ✅ Error handling validation

**Authentication Flow:**
1. **Admin Login** → Get admin token + profile
2. **Admin Operations** → Use token for CRUD operations
3. **Admin Logout** → Invalidate session

### 📊 **Updated Collection Stats**

- **Total Endpoints:** `51` (was 44)
- **Admin Endpoints:** `7` (new)
- **Protected Endpoints:** `25` (was 18)

### 🚀 **Usage Instructions**

1. **Import Updated Collection:**
   - Import `MBankingCore-API.postman_collection.json`
   - Import `MBankingCore-API.postman_environment.json`

2. **Default Admin Credentials:**
   - Email: `admin@mbankingcore.com`
   - Password: `admin123`
   - Role: `super`

3. **Testing Order:**
   ```
   1. Admin Login (get token)
   2. Get All Admins (verify access)
   3. Create Admin (test creation)
   4. Update Admin (modify created admin)
   5. Get Admin by ID (verify changes)
   6. Delete Admin (cleanup)
   7. Admin Logout (end session)
   ```

### 💡 **Admin Roles & Permissions**

- **Super Admin:** Full access to all admin operations
- **Admin:** Limited access (cannot manage other admins)

### 🔐 **Security Features**

- JWT-based authentication for admin sessions
- Role-based access control
- Password encryption with bcrypt
- Token expiration (24 hours)
- Authorization middleware protection

### 📋 **Response Examples**

**Admin Login Response:**
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
      "last_login": "2025-07-30T10:05:20.484234+07:00"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400
  }
}
```

**Get Admins Response:**
```json
{
  "code": 200,
  "message": "Admins retrieved successfully",
  "data": {
    "admins": [...],
    "total": 1,
    "page": 1,
    "per_page": 10
  }
}
```

## 🎉 Ready to Use!

The updated Postman collection is now ready for comprehensive admin management testing. All endpoints are pre-configured with proper authentication, validation, and automated variable management.

**Happy Testing! 🚀**
