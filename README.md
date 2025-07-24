# MBankingCore - Panduan Setup macOS

Go RESTful API dengan JWT Authentication menggunakan Gin Framework, GORM, dan PostgreSQL.

> 🍎 **Panduan setup ini khusus untuk macOS menggunakan Homebrew**
>
> 📋 **Untuk dokumentasi API lengkap:** [MBANKINGCORE-APIS.md](./MBANKINGCORE-APIS.md)
>
> ⚠️ **Pembaruan API Endpoints:** API endpoints telah disederhanakan - `/api/register`, `/api/login`, `/api/refresh` (sebelumnya `/api/v1/auth/multi-*`)

## 🏗️ Gambaran Arsitektur

### Fitur Utama

- 🔐 **Multi-Platform JWT Authentication** (Android, iOS, Web, Desktop)
- 📱 **Multi-Device Session Management** (Login dari multiple devices)
- 🔐 **SSO Provider Support** (Google, Apple, Facebook - Siap)
- 🔒 **Double-Layer Security** (SHA256 + bcrypt password hashing)
- 🎯 **Selective Logout** (Per device atau semua device)
- 👥 **User Management** (Operasi CRUD)
- ⚡ **RESTful API** dengan response format konsisten
- 🗄️ **PostgreSQL Database** dengan GORM ORM
- 🔄 **Auto Database Migration**
- 🌐 **CORS Support**
- ⚙️ **Environment Configuration**
- 📊 **Health Check Endpoint**

## 🏗️ Struktur Proyek

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
├── MBankingCore-API.md          # Dokumentasi API
├── MIGRATIONS.md                # Panduan migrasi database
└── README.md                    # Panduan setup & dokumentasi
```

## 📋 Prasyarat (macOS)

- **Go** 1.19+ (install via Homebrew: `brew install go`)
- **PostgreSQL** 12+ (install via Homebrew: `brew install postgresql`)
- **Homebrew** package manager
- **Git**

## 🚀 Mulai Cepat (macOS)

### 1. Clone Repository

```bash
git clone <repository-url>
cd mbankingcore
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Setup Database (macOS)

```bash
# Menggunakan Homebrew
brew install postgresql
brew services start postgresql

# Buat database
createdb mbcdb
```

### 4. Konfigurasi Environment

Copy dan edit file `.env`:

```bash
cp .env.example .env
```

Edit file `.env`:

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

### 5. Jalankan Aplikasi

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## 📖 Dokumentasi API

**📋 Untuk dokumentasi API lengkap dengan contoh request/response:**
👉 **[MBANKINGCORE-APIS.md](./MBANKINGCORE-APIS.md)**

## 🧪 Testing Multi-Platform Authentication

### Test Cepat dengan cURL

#### 1. Daftar User

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

#### 3. Dapatkan Active Sessions

```bash
curl -X GET http://localhost:8080/api/sessions \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Testing Postman

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

## 🔧 Panduan Development

### Hot Reload dengan Air

Untuk development yang lebih cepat dengan auto-reload:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Jalankan dengan hot reload
air
```

### Migrasi Database

Migrasi database dilakukan otomatis saat aplikasi start. Untuk operasi manual:

```bash
# Koneksi ke PostgreSQL
psql -h localhost -U your_username -d mbcdb

# Cek tabel
\dt

# Lihat struktur tabel users
\d users
```

## 🔧 Panduan Development & Deployment

### Build Production

```bash
# Build executable
go build -o mbankingcore

# Jalankan production build
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

### Setup Database untuk Production (macOS)

```bash
# Buat production database
createdb mbankingcore_prod

# Set production environment
export DB_NAME=mbankingcore_prod
export DB_HOST=your_db_host
export DB_USER=your_db_user
export DB_PASSWORD=your_secure_password
```

## 🔐 Implementasi Keamanan

### Keamanan Multi-Platform Authentication

#### Perlindungan Password Double-Layer

```
Client-Side: SHA256 Hash
 ↓
Server-Side: bcrypt Hash + Salt
 ↓ 
Database: bcrypt(SHA256(password))
```

#### Manajemen Device Session

- **Unique Device IDs**: Setiap device dilacak secara individual
- **Session Isolation**: Sessions per device, tidak saling mempengaruhi
- **Selective Logout**: Bisa logout per device atau semua device
- **Token Refresh**: Access token + Refresh token per device

#### Strategi JWT Token

- **Access Token**: Short-lived (15 menit)
- **Refresh Token**: Long-lived (7 hari)  
- **Device-Specific**: Token terikat dengan device_id
- **Auto-Invalidation**: Password change membatalkan semua sessions

### Best Practices Keamanan

#### Keamanan JWT

- Gunakan JWT_SECRET yang kuat dan random (minimum 32 karakter)
- Pertimbangkan token expiration yang lebih pendek untuk production
- Implementasikan token refresh mechanism

#### Keamanan Password

- bcrypt cost di-set ke default (cukup untuk kebanyakan use cases)
- Client-side SHA256 mencegah transmisi plain text

#### Keamanan Database

- Gunakan password database yang kuat
- Aktifkan SSL untuk koneksi database production
- Implementasikan database connection pooling

#### Keamanan API

- Implementasikan rate limiting
- Tambahkan request validation middleware
- Gunakan HTTPS di production
- Implementasikan konfigurasi CORS yang tepat

## 🧪 Testing & Validasi

### Unit Testing

```bash
# Jalankan tests
go test ./...

# Jalankan tests dengan coverage
go test -cover ./...

# Jalankan specific test
go test ./handlers -v
```

### Load Testing

```bash
# Install artillery untuk load testing
npm install -g artillery

# Buat artillery test config
# Kemudian jalankan load test
artillery run load-test.yml
```

### Testing API dengan Postman

1. Import collection: `postman/MBankingCore-API.postman_collection.json`
2. Import environment: `postman/MBankingCore-API.postman_environment.json`
3. Jalankan collection dengan Newman:

```bash
# Install Newman
npm install -g newman

# Jalankan Postman tests
newman run postman/MBankingCore-API.postman_collection.json \
  -e postman/MBankingCore-API.postman_environment.json
```

### Testing Cepat dengan cURL

```bash
# Jalankan server terlebih dahulu
go run main.go

# Test health check
curl http://localhost:8080/health

# Daftarkan test user (multi-platform)
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f","phone":"08123456789","provider":"email","device_info":{"device_type":"web","device_id":"web_browser_123","device_name":"Chrome","user_agent":"Mozilla/5.0"}}'

# Login dan dapatkan JWT token (multi-platform)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f","provider":"email","device_info":{"device_type":"web","device_id":"web_browser_123","device_name":"Chrome","user_agent":"Mozilla/5.0"}}'
```

## 📊 Monitoring & Logging

### Best Practices Logging

```go
// Tambahkan structured logging
import "github.com/sirupsen/logrus"

log := logrus.New()
log.SetFormatter(&logrus.JSONFormatter{})
log.SetLevel(logrus.InfoLevel)
```

### Health Checks

- `GET /health` endpoint untuk monitoring
- Pengecekan koneksi database
- Monitoring waktu response

## 🚀 Opsi Deployment (macOS)

### 1. Development Lokal

```bash
# Jalankan secara langsung
go run main.go

# Atau build dan jalankan
go build -o mbankingcore
./mbankingcore
```

### 2. Process Manager (PM2)

```bash
# Install PM2
npm install -g pm2

# Jalankan aplikasi dengan PM2
pm2 start mbankingcore --name "mbankingcore-api"
pm2 save
pm2 startup
```

### 3. Background Service dengan launchctl (macOS)

Buat file plist untuk service macOS:

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

Load service:

```bash
launchctl load ~/Library/LaunchAgents/com.mbankingcore.api.plist
launchctl start com.mbankingcore.api
```

## 🔍 Troubleshooting (macOS)

### Masalah Umum

**Error Koneksi Database:**

```bash
# Cek apakah PostgreSQL sedang berjalan
brew services list | grep postgresql

# Jalankan PostgreSQL jika belum berjalan
brew services start postgresql

# Cek status database
pg_isready -h localhost -U $(whoami)

# Cek environment variables
echo $DB_HOST $DB_PORT $DB_USER
```

**Port Sudah Digunakan:**

```bash
# Cari process yang menggunakan port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

**Masalah Permission PostgreSQL:**

```bash
# Jika mendapat permission denied, buat user dan database
createuser -s $(whoami)
createdb mbcdb
```

**Masalah JWT Token:**

- Verifikasi JWT_SECRET sudah di-set dengan benar
- Cek waktu expiration token
- Validasi format token

## 📚 Sumber Daya Tambahan

- [Dokumentasi Go](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [Dokumentasi GORM](https://gorm.io/)
- [JWT Go Library](https://github.com/golang-jwt/jwt)
- [Dokumentasi PostgreSQL](https://www.postgresql.org/docs/)

## 👥 Penulis

- **Gustaman** - Initial work

## 🔗 Link Referensi

- [Dokumentasi Go](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [JWT Go](https://github.com/golang-jwt/jwt)

---

**📋 Dokumentasi API:** [MBankingCore-API.md](./MBankingCore-API.md)

---

## Selamat Coding! 🚀
