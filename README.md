# MBankingCore - Mobile Banking Core API

Go RESTful API dengan Banking Authentication, JWT, Multi-Device Session Management menggunakan Gin Framework, GORM, dan PostgreSQL.

> � **Mobile Banking Core API** dengan 2-step OTP Authentication
>
> 📋 **Untuk dokumentasi API lengkap:** [docs/MBankingCore-API.md](./docs/MBankingCore-API.md)
>
> 🔄 **Banking Authentication:** Sistem autentikasi banking dengan OTP 2-langkah menggunakan login_token

## 🏗️ Gambaran Arsitektur

### Fitur Utama

- 🏦 **Banking Authentication** (2-step OTP process dengan login_token)
- 📱 **Multi-Device Session Management** (Login dari multiple devices)
- � **Multi-Account Banking Support** (CRUD bank accounts)
- � **JWT Authentication** dengan refresh token
- 🎯 **Selective Logout** (Per device atau semua device)
- 👥 **User Management** dengan role-based access (User, Admin, Owner)
- 📝 **Content Management** (Articles, Photos, Onboarding)
- ⚙️ **Configuration Management** (Dynamic app configuration)
- 📋 **Terms & Conditions** dan **Privacy Policy** management
- ⚡ **RESTful API** dengan response format konsisten (44 endpoints)
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
│   ├── article.go               # Article CRUD handlers
│   ├── auth.go                  # Banking authentication handlers
│   ├── bank_account.go          # Bank account management
│   ├── config.go                # Configuration handlers
│   ├── onboarding.go            # Onboarding content handlers
│   ├── photo.go                 # Photo management handlers
│   ├── privacy_policy.go        # Privacy policy handlers
│   ├── terms_conditions.go      # Terms & conditions handlers
│   └── user.go                  # User management handlers
├── middleware/
│   └── auth.go                  # JWT authentication middleware
├── models/
│   ├── article.go               # Article model & structures
│   ├── bank_account.go          # Bank account model
│   ├── config.go                # Configuration model
│   ├── constants.go             # Response codes & messages
│   ├── device_session.go        # Device session model
│   ├── onboarding.go            # Onboarding model
│   ├── photo.go                 # Photo model
│   ├── responses.go             # Response helper functions
│   └── user.go                  # User model & request structures
├── utils/
│   ├── auth.go                  # JWT utilities & password hashing
│   └── session.go               # Session management utilities
├── postman/
│   ├── MBankingCore-API.postman_collection.json    # Postman collection (9 endpoints)
│   └── MBankingCore-API.postman_environment.json   # Environment variables
├── docs/
│   ├── API-Endpoint-Reference.md     # Complete endpoint reference (44 endpoints)
│   ├── MBankingCore-API.md          # Full API documentation
│   ├── LOGIN_TOKEN_IMPLEMENTATION.md # Login token security documentation
│   ├── MIGRATIONS.md                # Database migration guide
│   ├── POSTMAN_UPDATE_LOG.md        # Postman collection update history
│   ├── SIMPLIFIED_LOGIN_VERIFY.md   # Login verification guide
│   └── VALIDATION_IMPLEMENTATION.md # Validation system documentation
├── .env                              # Environment variables
├── .env.example                      # Environment template
├── .gitignore                        # Git ignore rules
├── go.mod                           # Go modules
├── go.sum                           # Go modules checksum
├── main.go                          # Application entry point
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
👉 **[docs/MBankingCore-API.md](./docs/MBankingCore-API.md)**

**🔗 Referensi endpoint lengkap (44 endpoints):**
👉 **[docs/API-Endpoint-Reference.md](./docs/API-Endpoint-Reference.md)**

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

### Postman Testing

Import koleksi Postman untuk testing yang lebih komprehensif:

1. **Import Collection**: `postman/MBankingCore-API.postman_collection.json`
2. **Import Environment**: `postman/MBankingCore-API.postman_environment.json`
3. **Update Environment Variables**: Pastikan `banking_account_number` unik
4. **Run Collection**: Test semua endpoints dengan automated token management

**Fitur Postman Collection:**

- ✅ **Banking Authentication Flow** (2-step OTP process)
- ✅ **Automated token handling** & refresh
- 📱 **Multi-device scenarios** (Android, iOS, Web, Desktop)
- 🔄 **Session management** testing
- 🏦 **Bank account management** (CRUD operations)
- 📝 **Content management** (Articles, Photos, Onboarding)
- 🧪 **9 ready-to-use endpoints** dari total 44 available
- 📊 **Test result validation**

**Environment Variables yang Diperlukan:**
- `banking_account_number`: Gunakan nomor unik 16-digit
- `banking_phone`: Nomor telepon untuk OTP
- `banking_name`: Nama lengkap (min. 8 karakter)
- `banking_mother_name`: Nama ibu (min. 8 karakter)
- `banking_pin_atm`: PIN 6-digit
- `banking_otp_code`: Kode OTP (untuk testing, gunakan 6-digit apapun)

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

## 📚 Additional Resources & Documentation

### 📖 Documentation Files

- **[docs/MBankingCore-API.md](./docs/MBankingCore-API.md)** - Complete API documentation with examples
- **[docs/API-Endpoint-Reference.md](./docs/API-Endpoint-Reference.md)** - Quick reference for all 44 endpoints
- **[docs/LOGIN_TOKEN_IMPLEMENTATION.md](./docs/LOGIN_TOKEN_IMPLEMENTATION.md)** - Banking authentication security details
- **[docs/MIGRATIONS.md](./docs/MIGRATIONS.md)** - Database migration guide
- **[docs/POSTMAN_UPDATE_LOG.md](./docs/POSTMAN_UPDATE_LOG.md)** - Postman collection update history
- **[docs/SIMPLIFIED_LOGIN_VERIFY.md](./docs/SIMPLIFIED_LOGIN_VERIFY.md)** - Login verification process
- **[docs/VALIDATION_IMPLEMENTATION.md](./docs/VALIDATION_IMPLEMENTATION.md)** - Input validation system

### 🔗 External Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [JWT Go Library](https://github.com/golang-jwt/jwt)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 👥 Authors

- **Gustaman** - Initial work

---

**📋 Complete API Documentation:** [docs/MBankingCore-API.md](./docs/MBankingCore-API.md)

**🔗 Quick Endpoint Reference:** [docs/API-Endpoint-Reference.md](./docs/API-Endpoint-Reference.md)

---

## Happy Banking! 🏦🚀
