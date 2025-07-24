# Postman Collection Update Log

## Update: July 24, 2025

### ✅ Completed Updates

#### Collection File: `MBankingCore-API.postman_collection.json`
- **Collection ID**: Updated from `mbx-backend-api-collection` → `mbankingcore-api-collection`
- **Collection UID**: Updated from `12345-mbx-backend` → `12345-mbankingcore`
- **Terms & Conditions Content**: Updated references from `MBX Backend` → `MBankingCore`
- **Date Updated**: Last updated date changed from `January 2024` → `July 2025`

#### Environment File: `MBankingCore-API.postman_environment.json`
- **Environment ID**: Updated from `mbx-backend-env` → `mbankingcore-api-env`
- **Environment Name**: Remains `MBankingCore API Environment - Complete`

### 📋 Current API Coverage (37 endpoints)

#### Public APIs (7 endpoints)
- ✅ `GET /health` - Health check
- ✅ `GET /api/terms-conditions` - Get terms and conditions
- ✅ `POST /api/terms-conditions` - Set terms and conditions (admin)
- ✅ `GET /api/privacy-policy` - Get privacy policy
- ✅ `POST /api/privacy-policy` - Set privacy policy (admin)
- ✅ `GET /api/onboardings` - Get all onboardings
- ✅ `GET /api/onboardings/:id` - Get onboarding by ID

#### Authentication APIs (3 endpoints currently in collection)
- ✅ `POST /api/register` - User registration
- ✅ `POST /api/login` - User login
- ✅ `POST /api/refresh` - Refresh token

#### Protected User APIs (8 endpoints)
- ✅ `GET /api/profile` - Get user profile
- ✅ `PUT /api/profile` - Update user profile
- ✅ `GET /api/articles` - Get all articles
- ✅ `GET /api/articles/:id` - Get article by ID
- ✅ `PUT /api/articles/:id` - Update article
- ✅ `DELETE /api/articles/:id` - Delete article
- ✅ `GET /api/my-articles` - Get my articles
- ✅ `GET /api/photos` - Get all photos
- ✅ `GET /api/photos/:id` - Get photo by ID
- ✅ `PUT /api/photos/:id` - Update photo
- ✅ `DELETE /api/photos/:id` - Delete photo
- ✅ `GET /api/config/:key` - Get config value

#### Admin APIs (13 endpoints)
- ✅ `POST /api/articles` - Create article
- ✅ `POST /api/onboardings` - Create onboarding
- ✅ `PUT /api/onboardings/:id` - Update onboarding
- ✅ `DELETE /api/onboardings/:id` - Delete onboarding
- ✅ `POST /api/photos` - Create photo
- ✅ `GET /api/users` - List all users
- ✅ `GET /api/admin/users` - List admin users
- ✅ `GET /api/users/:id` - Get user by ID
- ✅ `DELETE /api/users/:id` - Delete user
- ✅ `POST /api/config` - Set config value
- ✅ `GET /api/admin/configs` - Get all configs
- ✅ `DELETE /api/config/:key` - Delete config

#### Owner APIs (2 endpoints)
- ✅ `POST /api/users` - Create user with role
- ✅ `PUT /api/users/:id` - Update user with role

### ⚠️ Missing Session Management Endpoints

Found session management handlers in `auth.go` but not registered in `main.go`:

#### Missing Endpoints (should be added):
- `GET /api/sessions` - Get active sessions
- `POST /api/logout` - Logout current session
- `POST /api/logout-others` - Logout other sessions

**Recommendation**: Add these endpoints to main.go and update Postman collection accordingly.

### 🔧 Environment Variables

Current environment includes tokens for multiple devices:
- `access_token` / `refresh_token` - Web session
- `android_access_token` / `android_refresh_token` - Android session
- `ios_access_token` / `ios_refresh_token` - iOS session  
- `desktop_access_token` / `desktop_refresh_token` - Desktop session

### 🎯 Testing Flow

The collection supports comprehensive testing:
1. **Public APIs** - Test without authentication
2. **Registration** - Create test users
3. **Multi-device Login** - Test device session management
4. **Protected APIs** - Test with authentication
5. **Admin Operations** - Test admin-only features
6. **Owner Operations** - Test owner-only features

### 📚 Next Steps

1. Consider adding missing session management endpoints to main.go
2. Update Postman collection to include session endpoints when added
3. Test collection with latest API changes
4. Ensure all 37 endpoints are properly tested

---

**Collection Status**: ✅ Updated and Ready for Use  
**Last Updated**: July 24, 2025  
**Updated By**: Gustaman
