# MBankingCore - macOS Setup Guide

Go RESTful API dengan JWT Authentication menggunakan Gin Framework, GORM, dan PostgreSQL.

> 🍎 **Setup guide ini khusus untuk macOS menggunakan Homebrew**
>
> 📋 **Untuk dokumentasi API lengkap:** [MBANKINGCORE-APIS.md](./MBANKINGCORE-APIS.md)
>
> ⚠️ **Updated API Endpoints:** API endpoints have been simplified - `/api/register`, `/api/login`, `/api/refresh` (previously `/api/v1/auth/multi-*`)

## 🏗️ Architecture Overview

### Core Features

- 🔐 **Multi-Platform JWT Authentication** (Android, iOS, Web, Desktop)
- � **Multi-Device Session Management** (Login dari multiple devices)
- � **SSO Provider Support** (Google, Apple, Facebook - Ready)
- 🔒 **Double-Layer Security** (SHA256 + bcrypt password hashing)
- 🎯 **Selective Logout** (Per device atau semua device)
- 👥 **User Management** (CRUD Operations)
- ⚡ **RESTful API** dengan response format konsisten
- 🗄️ **PostgreSQL Database** dengan GORM ORM
- 🔄 **Auto Database Migration**
- 🌐 **CORS Support**
- ⚙️ **Environment Configuration**
- 📊 **Health Check Endpoint**

## 🏗️ Project Structure

```
mbankingcore/
├── config/
│   └── database.go              # Database configuration & connection
├── handlers/
│   ├── auth.go                  # Multi-platform authentication handlers (consolidated)
│   └── user.go                  # User CRUD handlers
├── middleware/
│   └── auth.go                  # JWT authentication middleware
├── models/
│   ├── constants.go             # Response codes & messages constants
│   ├── responses.go             # Response helper functions
│   └── user.go                  # User model & request structures
├── utils/
│   ├── auth.go                  # JWT utilities & password hashing
│   └── session.go               # Session management utilities
├── postman/
│   ├── MBankingCore-API.postman_collection.json
│   └── MBankingCore-API.postman_environment.json
├── .env                         # Environment variables
├── .env.example                 # Environment template
├── .gitignore                   # Git ignore rules
├── go.mod                       # Go modules
├── go.sum                       # Go modules checksum
├── main.go                      # Application entry point
├── MBankingCore-API.md          # API Documentation
├── MIGRATIONS.md                # Database migration guide
└── README.md                    # Setup guide & documentation
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
👉 **[MBANKINGCORE-APIS.md](./MBANKINGCORE-APIS.md)**

## 🧪 Testing Multi-Platform Authentication

### Quick Test dengan cURL

#### 1. Register User

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
    "device_info": {
      "device_type": "web",
      "device_id": "browser_123",
      "device_name": "Chrome Browser",
      "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)"
    }
  }'
```

#### 2. Login User

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
    "device_info": {
      "device_type": "android",
      "device_id": "android_456",
      "device_name": "Samsung Galaxy",
      "user_agent": "MBankingCore-Android-App/1.0.0"
    }
  }'
```

#### 3. Get Active Sessions

```bash
curl -X GET http://localhost:8080/api/sessions \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Postman Testing

Import koleksi Postman untuk testing yang lebih komprehensif:

1. **Import Collection**: `postman/MBankingCore-API.postman_collection.json`
2. **Import Environment**: `postman/MBankingCore-API.postman_environment.json`
3. **Run Collection**: Test semua endpoints dengan automated token management

**Fitur Postman Collection:**

- ✅ Automated token handling & refresh
- 📱 Multi-device scenarios (Android, iOS, Web)
- 🔄 Session management testing
- 🧪 Comprehensive API coverage
- 📊 Test result validation

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

### Multi-Platform Authentication Security

#### Double-Layer Password Protection

```
Client-Side: SHA256 Hash
 ↓
Server-Side: bcrypt Hash + Salt
 ↓ 
Database: bcrypt(SHA256(password))
```

#### Device Session Management

- **Unique Device IDs**: Each device gets tracked individually
- **Session Isolation**: Sessions per device, tidak saling mempengaruhi
- **Selective Logout**: Bisa logout per device atau semua device
- **Token Refresh**: Access token + Refresh token per device

#### JWT Token Strategy

- **Access Token**: Short-lived (15 menit)
- **Refresh Token**: Long-lived (7 hari)  
- **Device-Specific**: Token terikat dengan device_id
- **Auto-Invalidation**: Password change invalidates semua sessions

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

#### API Security

- Implement rate limiting
- Add request validation middleware
- Use HTTPS in production
- Implement proper CORS configuration

## 🧪 Testing & Validation

### Unit Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./handlers -v
```

### Load Testing

```bash
# Install artillery for load testing
npm install -g artillery

# Create artillery test config
# Then run load test
artillery run load-test.yml
```

### API Testing dengan Postman

1. Import collection: `postman/MBankingCore-API.postman_collection.json`
2. Import environment: `postman/MBankingCore-API.postman_environment.json`
3. Run collection dengan Newman:

```bash
# Install Newman
npm install -g newman

# Run Postman tests
newman run postman/MBankingCore-API.postman_collection.json \
  -e postman/MBankingCore-API.postman_environment.json
```

### Quick Testing dengan cURL

```bash
# Start server first
go run main.go

# Test health check
curl http://localhost:8080/health

# Register a test user (multi-platform)
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f","phone":"08123456789","provider":"email","device_info":{"device_type":"web","device_id":"web_browser_123","device_name":"Chrome","user_agent":"Mozilla/5.0"}}'

# Login and get JWT token (multi-platform)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f","provider":"email","device_info":{"device_type":"web","device_id":"web_browser_123","device_name":"Chrome","user_agent":"Mozilla/5.0"}}'
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

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [JWT Go Library](https://github.com/golang-jwt/jwt)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 👥 Authors

- **Gustaman** - Initial work

## 🔗 Links

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [JWT Go](https://github.com/golang-jwt/jwt)

---

**📋 API Documentation:** [MBANKINGCORE-APIS.md](./MBANKINGCORE-APIS.md)

---

## Happy Coding! 🚀
