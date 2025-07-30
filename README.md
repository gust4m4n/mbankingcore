# MBankingCore - Mobile Banking Core API

Go RESTful API denga├── models/
│   ├── admin.go                 # Admin model & structures
│   ├── article.go               # Article model & structures
│   ├── audit.go                 # Audit trails models (NEW)
│   ├── bank_account.go          # Bank account model
│   ├── config.go                # Configuration model
│   ├── constants.go             # Response codes & messages
│   ├── device_session.go        # Device session model
│   ├── onboarding.go            # Onboarding model
│   ├── photo.go                 # Photo model
│   ├── responses.go             # Response helper functions
│   ├── transaction.go           # Transaction model & structures (NEW)
│   └── user.go                  # User model & request structuresthentication, JWT, Multi-Device Session Management menggunakan Gin Framework, GORM, dan PostgreSQL.

> � **Mobile Banking Core API** dengan 2-step OTP Authentication
>
> 📋 **Untuk dokumentasi API lengkap:** [MBANKINGCORE-API.md](./MBANKINGCORE-API.md)
>
> 🔄 **Banking Authentication:** Sistem autentikasi banking dengan OTP 2-langkah menggunakan login_token

## 🏗️ Gambaran Arsitektur

### Fitur Utama

- 🏦 **Banking Authentication** (2-step OTP process dengan login_token)
- 📱 **Multi-Device Session Management** (Login dari multiple devices)
- 💼 **Multi-Account Banking Support** (CRUD bank accounts)
- 💳 **Transaction Management** (Topup, withdraw, transfer, reversal)
- 🔄 **Transaction Reversal System** (Admin-only dengan audit trail lengkap)
- 🔑 **JWT Authentication** dengan refresh token
- 🎯 **Selective Logout** (Per device atau semua device)
- 👥 **User Management** dengan role-based access (User, Admin, Owner)
- 🔧 **Admin Management System** (Admin authentication & CRUD)
- 📝 **Content Management** (Articles, Photos, Onboarding)
- ⚙️ **Configuration Management** (Dynamic app configuration)
- 📋 **Terms & Conditions** dan **Privacy Policy** management
- 🔍 **Comprehensive Audit Trails** (Activity & Login monitoring)
- 💰 **Transaction Management** dengan reversal system
- ⚡ **RESTful API** dengan response format konsisten (60 endpoints)
- 🗄️ **PostgreSQL Database** dengan GORM ORM
- 🔄 **Auto Database Migration**
- 🌐 **CORS Support**
- ⚙️ **Environment Configuration**
- 📊 **Health Check Endpoint**

## 🏗️ Struktur Proyek

```
mbankingcore/
├── cmd/
│   └── migrate/
│       └── main.go              # Database migration utility
├── config/
│   ├── database.go              # Database configuration & connection
│   └── migrations.go            # Migration management
├── handlers/
│   ├── admin.go                 # Admin management handlers (NEW)
│   ├── article.go               # Article CRUD handlers
│   ├── audit.go                 # Audit trails handlers (NEW)
│   ├── auth.go                  # Banking authentication handlers
│   ├── bank_account.go          # Bank account management
│   ├── config.go                # Configuration handlers
│   ├── onboarding.go            # Onboarding content handlers
│   ├── photo.go                 # Photo management handlers
│   ├── privacy_policy.go        # Privacy policy handlers
│   ├── terms_conditions.go      # Terms & conditions handlers
│   ├── transaction.go           # Transaction management handlers (NEW)
│   └── user.go                  # User management handlers
├── middleware/
│   ├── admin_auth.go            # Admin authentication middleware (NEW)
│   ├── audit.go                 # Audit logging middleware (NEW)
│   └── auth.go                  # JWT authentication middleware
├── models/
│   ├── admin.go                 # Admin model & structures (NEW)
│   ├── article.go               # Article model & structures
│   ├── bank_account.go          # Bank account model
│   ├── config.go                # Configuration model
│   ├── constants.go             # Response codes & messages
│   ├── device_session.go        # Device session model
│   ├── onboarding.go            # Onboarding model
│   ├── photo.go                 # Photo model
│   ├── responses.go             # Response helper functions
│   ├── transaction.go           # Transaction model & structures (NEW)
│   └── user.go                  # User model & request structures
├── utils/
│   ├── admin_auth.go            # Admin JWT utilities (NEW)
│   ├── auth.go                  # JWT utilities & password hashing
│   └── session.go               # Session management utilities
├── postman/
│   ├── MBankingCore-API.postman_collection.json    # Postman collection (60 endpoints)
│   └── MBankingCore-API.postman_environment.json   # Environment variables
├── .env                              # Environment variables
├── .env.example                      # Environment template
├── .gitignore                        # Git ignore rules
├── go.mod                           # Go modules
├── go.sum                           # Go modules checksum
├── main.go                          # Application entry point
├── MBANKINGCORE-API.md              # Complete API documentation (58 endpoints)
└── README.md                        # This documentation
```

## 📋 Prerequisites (macOS)

- **Go** 1.19+ (install via Homebrew: `brew install go`)
- **PostgreSQL** 12+ (install via Homebrew: `brew install postgresql`)
- **Homebrew** package manager
- **Git**

## 🚀 Quick Start (macOS)

### 1. Clone Repository

```bash
git clone <repository-url>
cd mbankingcore
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Database Setup (macOS)

```bash
# Using Homebrew
brew install postgresql
brew services start postgresql

# Create database
createdb mbcdb
```

### 4. Environment Configuration

Copy dan edit file `.env`:

```bash
cp .env.example .env
```

Edit `.env` file:

```properties
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=mbcdb
DB_SSLMODE=disable

# Server Configuration
PORT=8080

# JWT Configuration (Production: use secure random string)
JWT_SECRET=your-secret-key-change-this-in-production
```

### 5. Run Application

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## 📖 API Documentation

**📋 Untuk dokumentasi API lengkap dengan contoh request/response:**
👉 **[MBANKINGCORE-API.md](./MBANKINGCORE-API.md)**

## 🏦 Banking Authentication System

MBankingCore menggunakan sistem autentikasi banking dengan 2-step OTP process yang aman:

### 🔐 Authentication Flow

1. **Banking Login (Step 1)** - `POST /api/login`
   - Submit: name, account_number, mother_name, phone, pin_atm, device_info
   - Receive: login_token (expires in 5 minutes)
   - OTP dikirim ke nomor telepon

2. **Banking Login Verification (Step 2)** - `POST /api/login/verify`
   - Submit: login_token + otp_code
   - Receive: access_token, refresh_token, user info

3. **Access Protected APIs** dengan Bearer token
   - Header: `Authorization: Bearer <access_token>`

4. **Token Refresh** - `POST /api/refresh`
   - Submit: refresh_token
   - Receive: new access_token

### 🔑 Key Security Features

- **login_token**: Temporary token (5 menit) untuk verifikasi OTP
- **Unique Account Numbers**: Setiap account number harus unik
- **Multi-Device Support**: Login dari berbagai device secara bersamaan
- **Selective Logout**: Logout per device atau semua device
- **Auto-Registration**: Nomor baru otomatis terdaftar setelah verifikasi OTP

## 🔧 Admin Management System

MBankingCore dilengkapi dengan sistem manajemen admin yang komprehensif untuk mengelola administrator aplikasi.

### 👑 Admin Authentication Flow

1. **Admin Login** - `POST /api/admin/login`
   - Submit: email, password
   - Receive: admin_token (expires in 24 hours), admin profile

2. **Access Admin APIs** dengan Bearer token
   - Header: `Authorization: Bearer <admin_token>`

3. **Admin Logout** - `POST /api/admin/logout`
   - Invalidate admin session

### 🔒 Admin Security Features

- **JWT-based Authentication**: Separate token system untuk admin
- **Role-based Access Control**: Super Admin vs Admin permissions
- **Password Encryption**: bcrypt hashing untuk password security
- **Status Management**: Active, Inactive, Blocked admin states
- **Self-Protection**: Admin tidak bisa menghapus akun sendiri
- **Email Uniqueness**: Validasi email unik untuk setiap admin

### 👥 Admin Roles & Permissions

**Super Admin:**
- Full access to all admin operations
- Can create, update, delete other admins
- Can manage system configurations

**Admin:**
- Limited access to admin operations
- Cannot manage other admin accounts
- Can access admin-protected content endpoints

### 📋 Admin Management Endpoints (7 endpoints)

| Endpoint | Method | Path | Access Level |
|----------|--------|------|--------------|
| Admin Login | `POST` | `/api/admin/login` | Public (Credentials Required) |
| Admin Logout | `POST` | `/api/admin/logout` | Admin Authentication |
| Get All Admins | `GET` | `/api/admin/admins` | Admin Authentication |
| Get Admin by ID | `GET` | `/api/admin/admins/:id` | Admin Authentication |
| Create Admin | `POST` | `/api/admin/admins` | Super Admin Only |
| Update Admin | `PUT` | `/api/admin/admins/:id` | Super Admin Only |
| Delete Admin | `DELETE` | `/api/admin/admins/:id` | Super Admin Only |

### 🔑 Default Admin Credentials

**Super Admin Account:**
- Email: `admin@mbankingcore.com`
- Password: `admin123`
- Role: `super`
- Status: `active`

⚠️ **Production Warning**: Change default credentials immediately in production!

## 💰 Transaction Management

### Transaction Features

- 💵 **Topup Balance** - Add balance to user account
- 💸 **Withdraw Balance** - Deduct balance from user account  
- 🔄 **Transfer Balance** - Transfer balance between users using account numbers
- ↩️ **Transaction Reversal** - Admin-only reversal with comprehensive business logic
- 📊 **Transaction History** - Complete audit trail with pagination
- 🔒 **Atomic Operations** - Database transactions with row-level locking
- ⚡ **Real-time Balance Updates** - Immediate balance reflection
- 📋 **Admin Monitoring** - Admin dashboard for all transactions
- 🛡️ **Reversal Audit Trail** - Complete transaction relationship tracking

### 📋 Transaction Endpoints (6 endpoints)

| Endpoint | Method | Path | Access Level |
|----------|--------|------|--------------|
| Topup Balance | `POST` | `/api/transactions/topup` | User Authentication |
| Withdraw Balance | `POST` | `/api/transactions/withdraw` | User Authentication |
| Transfer Balance | `POST` | `/api/transactions/transfer` | User Authentication |
| Transaction History | `GET` | `/api/transactions/history` | User Authentication |
| Admin Transaction Monitor | `GET` | `/api/admin/transactions` | Admin Authentication |
| Transaction Reversal | `POST` | `/api/admin/transactions/reversal` | Admin Authentication |

### 🔄 Transaction Types

- **topup** - Balance addition operation
- **withdraw** - Balance deduction operation
- **transfer_out** - Outgoing transfer (sender side)
- **transfer_in** - Incoming transfer (receiver side)

### ↩️ Transaction Reversal System

**Admin-Only Feature** dengan business logic komprehensif:

#### Reversal Business Logic

1. **Topup Reversal**: Deducts the topup amount from user balance
2. **Withdraw Reversal**: Adds the withdraw amount back to user balance  
3. **Transfer Reversal**: Creates two reversal transactions:
   - Adds amount back to sender's balance (`transfer_out` → `transfer_in`)
   - Deducts amount from receiver's balance (`transfer_in` → `transfer_out`)

#### Reversal Features

- 🔐 **Admin-Only Access** - Requires admin authentication
- 🔗 **Transaction Relationships** - Links original and reversal transactions
- 📝 **Reversal Reason** - Mandatory reason for audit purposes
- ⏰ **Timestamp Tracking** - Records when reversal occurred
- 🚫 **Duplicate Prevention** - Prevents reversing already reversed transactions
- 💰 **Balance Validation** - Ensures sufficient balance for deduction reversals
- 🔒 **Atomic Operations** - Ensures data consistency during complex reversals

#### Reversal Request Example

```json
{
    "transaction_id": 123,
    "reason": "Administrative reversal - Error correction"
}
```

## 🔍 Audit Trails System

### Audit Features

- **Comprehensive Activity Logging** - Records all user and admin actions
- **Login/Logout Monitoring** - Tracks authentication activities  
- **Real-time Tracking** - Automatic logging via middleware
- **Advanced Filtering** - Filter by user, admin, entity type, date range, IP
- **Pagination Support** - Efficient handling of large audit datasets
- **Admin-Only Access** - Secure access control for audit data

### 📋 Audit Endpoints (2 endpoints)

- `GET /api/admin/audit-logs` - Get activity audit logs with filtering (Admin only)
- `GET /api/admin/login-audits` - Get login/logout audit logs (Admin only)

### 🔐 Audit Data Types

#### 1. Activity Audit Logs

Tracks all system activities including:

- **User Actions**: Profile updates, transactions, content creation
- **Admin Actions**: User management, transaction reversals, system configuration
- **Entity Operations**: CREATE, READ, UPDATE, DELETE operations
- **API Calls**: Request details, response codes, execution time

#### 2. Login Audit Logs  

Monitors authentication activities:

- **Login Attempts**: Successful and failed login attempts
- **Logout Activities**: User and admin logout tracking
- **Device Information**: IP addresses, user agents, device details
- **Security Events**: Failed attempts, blocked access, unusual activity

### 🔍 Filtering Capabilities

```bash
# Filter by entity type
GET /api/admin/audit-logs?entity_type=transaction

# Filter by user ID  
GET /api/admin/audit-logs?user_id=123

# Filter by admin ID
GET /api/admin/audit-logs?admin_id=456

# Filter by date range
GET /api/admin/audit-logs?start_date=2024-01-01&end_date=2024-01-31

# Filter by IP address
GET /api/admin/audit-logs?ip_address=192.168.1.100

# Combine multiple filters with pagination
GET /api/admin/audit-logs?entity_type=user&start_date=2024-01-01&page=1&limit=50
```

### ⚡ Automatic Logging

The audit system automatically logs:

- **All API Requests** - Method, endpoint, parameters, response codes
- **User Activities** - Profile changes, transactions, content management  
- **Admin Activities** - User management, system configuration, transaction oversight
- **Authentication Events** - Login/logout attempts with device information
- **Security Events** - Failed attempts, blocked access, suspicious activities

### 🔒 Security Features

- **Atomic Database Transactions** - Ensures data consistency
- **Row-level Locking** - Prevents concurrent balance conflicts  
- **Balance Validation** - Prevents negative balances and invalid amounts
- **Account Number Verification** - Validates recipient accounts for transfers
- **Self-transfer Prevention** - Blocks transfers to same account
- **Complete Audit Trail** - Balance before/after tracking

## 🧪 Testing Banking Authentication

### Quick Test dengan cURL

#### 1. Banking Login (Step 1) - Send OTP

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe Smith",
    "account_number": "1234567890123456",
    "mother_name": "Jane Doe Smith", 
    "phone": "081234567890",
    "pin_atm": "123456",
    "device_info": {
      "device_type": "android",
      "device_id": "android_test_123",
      "device_name": "Samsung Galaxy S23"
    }
  }'
```

#### 2. Banking Login Verification (Step 2) - Verify OTP

```bash
curl -X POST http://localhost:8080/api/login/verify \
  -H "Content-Type: application/json" \
  -d '{
    "login_token": "your_login_token_from_step1",
    "otp_code": "123456"
  }'
```

#### 3. Access Protected Endpoint

```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🧪 Testing Admin Authentication

### Quick Test dengan cURL

#### 1. Admin Login

```bash
curl -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@mbankingcore.com",
    "password": "admin123"
  }'
```

#### 2. Get All Admins

```bash
curl -X GET http://localhost:8080/api/admin/admins \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

#### 3. Create New Admin (Super Admin only)

```bash
curl -X POST http://localhost:8080/api/admin/admins \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Admin",
    "email": "test@mbankingcore.com",
    "password": "password123",
    "role": "admin"
  }'
```

#### 4. Admin Logout

```bash
curl -X POST http://localhost:8080/api/admin/logout \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### Postman Testing

Import koleksi Postman untuk testing yang lebih komprehensif:

1. **Import Collection**: `postman/MBankingCore-API.postman_collection.json`
2. **Import Environment**: `postman/MBankingCore-API.postman_environment.json`
3. **Update Environment Variables**: Pastikan `banking_account_number` unik
4. **Run Collection**: Test semua endpoints dengan automated token management

**Fitur Postman Collection:**

- ✅ **Banking Authentication Flow** (2-step OTP process)
- ✅ **Admin Authentication Flow** (Admin login/logout)
- ✅ **Automated token handling** & refresh
- 📱 **Multi-device scenarios** (Android, iOS, Web, Desktop)
- 🔄 **Session management** testing
- 🏦 **Bank account management** (CRUD operations)
- � **Admin management** (Admin CRUD operations)
- �📝 **Content management** (Articles, Photos, Onboarding)
- 🧪 **58 ready-to-use endpoints** (Complete API coverage)
- 📊 **Test result validation**

**Environment Variables yang Diperlukan:**

**Banking Variables:**
- `banking_account_number`: Gunakan nomor unik 16-digit
- `banking_phone`: Nomor telepon untuk OTP
- `banking_name`: Nama lengkap (min. 8 karakter)
- `banking_mother_name`: Nama ibu (min. 8 karakter)
- `banking_pin_atm`: PIN 6-digit
- `banking_otp_code`: Kode OTP (untuk testing, gunakan 6-digit apapun)

**Admin Variables:**
- `admin_email`: Email admin (default: admin@mbankingcore.com)
- `admin_password`: Password admin (default: admin123)
- `new_admin_name`: Nama admin baru untuk testing
- `new_admin_email`: Email admin baru untuk testing

## 🔧 Development Guide

### Hot Reload dengan Air

Untuk development yang lebih cepat dengan auto-reload:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Database Migration

Database migration dilakukan otomatis saat aplikasi start. Untuk operasi manual:

```bash
# Connect to PostgreSQL
psql -h localhost -U your_username -d mbcdb

# Check tables
\dt

# View users table structure
\d users
```

## 🔧 Development & Deployment Guide

### Build Production

```bash
# Build executable
go build -o mbankingcore

# Run production build
./mbankingcore
```

### Environment Variables untuk Production

```bash
export DB_HOST=your_production_db_host
export DB_PASSWORD=your_production_db_password
export JWT_SECRET=your_very_secure_jwt_secret_key
export GIN_MODE=release
export PORT=8080
```

### Database Setup untuk Production (macOS)

```bash
# Create production database
createdb mbankingcore_prod

# Set production environment
export DB_NAME=mbankingcore_prod
export DB_HOST=your_db_host
export DB_USER=your_db_user
export DB_PASSWORD=your_secure_password
```

## 🔐 Security Implementation

### Banking Authentication Security

#### 2-Step OTP Authentication Process

```
Step 1: Banking Login
Client → Server: Credentials + Device Info
Server → Client: login_token (5 min expiry)
Server → SMS: OTP Code

Step 2: OTP Verification  
Client → Server: login_token + OTP
Server → Client: access_token + refresh_token
```

#### Multi-Device Session Management

- **Device-Specific Sessions**: Each device gets unique session tracking
- **Session Isolation**: Sessions per device, tidak saling mempengaruhi
- **Selective Logout**: Bisa logout per device atau semua device
- **Auto Session Cleanup**: Expired sessions otomatis dibersihkan

#### JWT Token Strategy

- **Access Token**: Short-lived (24 jam)
- **Refresh Token**: Long-lived (7 hari)  
- **Device-Specific**: Token terikat dengan device_id
- **Auto-Invalidation**: PIN change invalidates semua sessions

#### Banking Security Features

- **Unique Account Numbers**: Database constraint untuk mencegah duplikasi
- **PIN ATM Protection**: bcrypt hashing untuk PIN storage
- **OTP Security**: Random 6-digit OTP dengan expiry time
- **login_token**: Temporary secure token dengan crypto/rand generation

### Security Best Practices

#### JWT Security

- Use strong, random JWT_SECRET (minimum 32 characters)
- Consider shorter token expiration for production
- Implement token refresh mechanism

#### Password Security

- bcrypt cost is set to default (sufficient for most use cases)
- Client-side SHA256 prevents plain text transmission

#### Database Security

- Use strong database passwords
- Enable SSL for production database connections
- Implement database connection pooling
- Account number uniqueness constraints

#### API Security

- Banking authentication with OTP verification
- JWT token-based authorization
- Multi-device session management
- Request validation middleware
- Use HTTPS in production
- Implement proper CORS configuration

## 🧪 Testing & Validation

### Manual Testing dengan cURL

```bash
# Start server first
go run main.go

# Test health check
curl http://localhost:8080/health

# Banking Login Step 1 (Send OTP)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","account_number":"1234567890123456","mother_name":"Test Mother","phone":"081234567890","pin_atm":"123456","device_info":{"device_type":"android","device_id":"test_device_123","device_name":"Test Device"}}'

# Banking Login Step 2 (Verify OTP) - use login_token from Step 1
curl -X POST http://localhost:8080/api/login/verify \
  -H "Content-Type: application/json" \
  -d '{"login_token":"your_login_token_here","otp_code":"123456"}'

# Admin Login
curl -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@mbankingcore.com","password":"admin123"}'

# Get All Admins (use admin_token from above)
curl -X GET http://localhost:8080/api/admin/admins \
  -H "Authorization: Bearer your_admin_token_here"
```

### Postman Collection Testing

1. Import collection: `postman/MBankingCore-API.postman_collection.json`
2. Import environment: `postman/MBankingCore-API.postman_environment.json`
3. Update `banking_account_number` dengan nomor unik
4. Run collection dengan Newman:

```bash
# Install Newman
npm install -g newman

# Run Postman tests
newman run postman/MBankingCore-API.postman_collection.json \
  -e postman/MBankingCore-API.postman_environment.json
```

## 📊 Monitoring & Logging

### Logging Best Practices

```go
// Add structured logging
import "github.com/sirupsen/logrus"

log := logrus.New()
log.SetFormatter(&logrus.JSONFormatter{})
log.SetLevel(logrus.InfoLevel)
```

### Health Checks

- `GET /health` endpoint untuk monitoring
- Database connection check
- Response time monitoring

## 🚀 Deployment Options (macOS)

### 1. Local Development

```bash
# Run directly
go run main.go

# Or build and run
go build -o mbankingcore
./mbankingcore
```

### 2. Process Manager (PM2)

```bash
# Install PM2
npm install -g pm2

# Start application with PM2
pm2 start mbankingcore --name "mbankingcore-api"
pm2 save
pm2 startup
```

### 3. Background Service with launchctl (macOS)

Create a plist file for macOS service:

```xml
<!-- ~/Library/LaunchAgents/com.mbankingcore.api.plist -->
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.mbankingcore.api</string>
    <key>ProgramArguments</key>
    <array>
        <string>/path/to/your/mbankingcore</string>
    </array>
    <key>WorkingDirectory</key>
    <string>/path/to/your/mbankingcore/directory</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

Load the service:

```bash
launchctl load ~/Library/LaunchAgents/com.mbankingcore.api.plist
launchctl start com.mbankingcore.api
```

## 🔍 Troubleshooting (macOS)

### Common Issues

**Database Connection Error:**

```bash
# Check if PostgreSQL is running
brew services list | grep postgresql

# Start PostgreSQL if not running
brew services start postgresql

# Check database status
pg_isready -h localhost -U $(whoami)

# Check environment variables
echo $DB_HOST $DB_PORT $DB_USER
```

**Port Already in Use:**

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

**Permission Issues with PostgreSQL:**

```bash
# If you get permission denied, create user and database
createuser -s $(whoami)
createdb mbcdb
```

**JWT Token Issues:**

- Verify JWT_SECRET is set correctly
- Check token expiration time
- Validate token format

## 🆕 Recent Updates

### Version 2.0 - Transaction Reversal System

**🎉 New Feature: Transaction Reversal API**

- ↩️ **Admin Transaction Reversal** - Comprehensive reversal system for all transaction types
- 🔐 **Admin-Only Access** - Secure reversal operations with admin authentication
- 🔗 **Complete Audit Trail** - Full transaction relationship tracking
- 💰 **Smart Business Logic** - Handles topup, withdraw, and transfer reversals
- 🛡️ **Data Integrity** - Atomic operations with balance validation
- 📝 **Reversal Reasons** - Mandatory documentation for all reversals
- ⏰ **Timestamp Tracking** - Complete reversal history

**Technical Enhancements:**

- Enhanced Transaction model with reversal tracking fields
- New admin endpoint: `POST /api/admin/transactions/reversal`
- Comprehensive reversal business logic for all transaction types
- Updated Postman collection with reversal testing
- Complete API documentation updates
- Database migration for reversal functionality

**Total Endpoints: 58** (Previous: 57)

## 📚 Additional Resources & Documentation

### 📖 Documentation Files

- **[MBANKINGCORE-API.md](./MBANKINGCORE-API.md)** - Complete API documentation with examples and endpoint reference (58 endpoints)

### 🔗 External Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [JWT Go Library](https://github.com/golang-jwt/jwt)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 👥 Authors

- **Gustaman** - Initial work

---

**📋 Complete API Documentation:** [MBANKINGCORE-API.md](./MBANKINGCORE-API.md)

---

## Happy Banking! 🏦🚀
