# MBankingCore Database Documentation

Dokumentasi lengkap struktur database untuk aplikasi MBankingCore dengan PostgreSQL sebagai database utama.

## 📋 Database Overview

**Database Engine:** PostgreSQL 12+  
**ORM:** GORM (Go ORM)  
**Auto Migration:** Enabled  
**Connection Pool:** Configured  
**SSL Mode:** Configurable (disable/require)

### Database Schema Summary

| Table | Description | Records | Primary Key |
|-------|-------------|---------|-------------|
| `users` | User accounts with banking authentication | Dynamic | `id` (uint) |
| `bank_accounts` | Multi-account banking support | Dynamic | `id` (uint) |
| `device_sessions` | Multi-device session management | Dynamic | `id` (uint) |
| `articles` | Content management articles | Dynamic | `id` (uint) |
| `photos` | Photo management system | Dynamic | `id` (uint) |
| `onboardings` | App onboarding content | Dynamic | `id` (uint) |
| `configs` | Dynamic application configuration | Dynamic | `key` (string) |

---

## 🏗️ Database Tables

### 1. users

**Purpose:** Core user accounts with banking authentication system

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255) UNIQUE NOT NULL,
    mother_name VARCHAR(255) NOT NULL,
    pin_atm VARCHAR(255) NOT NULL,  -- bcrypt hashed
    role VARCHAR(20) DEFAULT 'user',
    avatar VARCHAR(500),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Indexes
CREATE UNIQUE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_role ON users(role);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment user ID |
| `name` | VARCHAR(255) | NOT NULL | Full name (min 8 characters) |
| `phone` | VARCHAR(255) | UNIQUE, NOT NULL | Phone number (unique identifier) |
| `mother_name` | VARCHAR(255) | NOT NULL | Mother's name (min 8 characters) |
| `pin_atm` | VARCHAR(255) | NOT NULL | Hashed PIN ATM (6 digits, bcrypt) |
| `role` | VARCHAR(20) | DEFAULT 'user' | User role: 'user', 'admin', 'owner' |
| `avatar` | VARCHAR(500) | NULLABLE | Avatar image URL |
| `created_at` | TIMESTAMP | NOT NULL | Record creation time |
| `updated_at` | TIMESTAMP | NOT NULL | Last update time |

**Role Values:**

- `user` - Standard banking user (default)
- `admin` - Administrative user
- `owner` - System owner (full access)

**Relationships:**

- **One-to-Many** with `bank_accounts` (via `user_id`)
- **One-to-Many** with `device_sessions` (via `user_id`)
- **One-to-Many** with `articles` (via `user_id`)

---

### 2. bank_accounts

**Purpose:** Multi-account banking support for users

```sql
CREATE TABLE bank_accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    account_name VARCHAR(100) NOT NULL,
    bank_name VARCHAR(100),
    bank_code VARCHAR(10),
    account_type VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_bank_accounts_user_id ON bank_accounts(user_id);
CREATE UNIQUE INDEX idx_user_account ON bank_accounts(user_id, account_number);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment account ID |
| `user_id` | INTEGER | NOT NULL, FK | Reference to users table |
| `account_number` | VARCHAR(50) | NOT NULL | Bank account number (8-20 digits) |
| `account_name` | VARCHAR(100) | NOT NULL | Account holder name (3-100 chars) |
| `bank_name` | VARCHAR(100) | NULLABLE | Bank institution name |
| `bank_code` | VARCHAR(10) | NULLABLE | Bank code (e.g., "014" for BCA) |
| `account_type` | VARCHAR(20) | NULLABLE | Account type: 'saving', 'checking', 'current' |
| `is_active` | BOOLEAN | DEFAULT true | Account status |
| `is_primary` | BOOLEAN | DEFAULT false | Primary account flag |
| `created_at` | TIMESTAMP | NOT NULL | Record creation time |
| `updated_at` | TIMESTAMP | NOT NULL | Last update time |

**Unique Constraints:**

- User can have multiple accounts, but each account number must be unique per user
- Composite unique index: `(user_id, account_number)`

**Business Rules:**

- Each user can have multiple bank accounts
- Only one account can be marked as primary per user
- Account numbers must be 8-20 characters
- Account deletion is soft (set `is_active = false`)

---

### 3. device_sessions

**Purpose:** Multi-device session management with JWT tokens

```sql
CREATE TABLE device_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    device_id VARCHAR(255),
    device_name VARCHAR(255),
    provider VARCHAR(50) NOT NULL,
    provider_id VARCHAR(255),
    ip_address VARCHAR(45),
    is_active BOOLEAN DEFAULT true,
    last_activity TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_device_sessions_user_id ON device_sessions(user_id);
CREATE INDEX idx_device_sessions_device_id ON device_sessions(device_id);
CREATE UNIQUE INDEX idx_session_token ON device_sessions(session_token);
CREATE UNIQUE INDEX idx_refresh_token ON device_sessions(refresh_token);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment session ID |
| `user_id` | INTEGER | NOT NULL, FK | Reference to users table |
| `session_token` | VARCHAR(255) | UNIQUE, NOT NULL | JWT access token |
| `refresh_token` | VARCHAR(255) | UNIQUE, NOT NULL | JWT refresh token |
| `device_type` | VARCHAR(50) | NOT NULL | Device type enum |
| `device_id` | VARCHAR(255) | NULLABLE | Unique device identifier |
| `device_name` | VARCHAR(255) | NULLABLE | Human-readable device name |
| `provider` | VARCHAR(50) | NOT NULL | Authentication provider |
| `provider_id` | VARCHAR(255) | NULLABLE | Provider-specific user ID |
| `ip_address` | VARCHAR(45) | NULLABLE | Client IP address (IPv4/IPv6) |
| `is_active` | BOOLEAN | DEFAULT true | Session status |
| `last_activity` | TIMESTAMP | NOT NULL | Last activity timestamp |
| `expires_at` | TIMESTAMP | NOT NULL | Session expiration time |
| `created_at` | TIMESTAMP | NOT NULL | Record creation time |
| `updated_at` | TIMESTAMP | NOT NULL | Last update time |

**Device Types:**

- `android` - Android mobile app
- `ios` - iOS mobile app
- `web` - Web browser
- `desktop` - Desktop application
- `google_sso` - Google SSO login
- `apple_sso` - Apple SSO login
- `facebook_sso` - Facebook SSO login

**Authentication Providers:**

- `email` - Banking authentication (primary)
- `google` - Google OAuth
- `apple` - Apple Sign-In
- `facebook` - Facebook Login

**Session Management:**

- Multiple concurrent sessions per user
- Device-specific token pairs
- Automatic cleanup of expired sessions
- Selective logout (per device or all devices)

---

### 4. articles

**Purpose:** Content management system for articles

```sql
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    image VARCHAR(500),
    content TEXT,
    is_active BOOLEAN DEFAULT true,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_articles_user_id ON articles(user_id);
CREATE INDEX idx_articles_is_active ON articles(is_active);
CREATE INDEX idx_articles_created_at ON articles(created_at);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment article ID |
| `title` | VARCHAR(255) | NOT NULL | Article title (1-255 chars) |
| `image` | VARCHAR(500) | NULLABLE | Featured image URL |
| `content` | TEXT | NULLABLE | Article content (unlimited text) |
| `is_active` | BOOLEAN | DEFAULT true | Publication status |
| `user_id` | INTEGER | NOT NULL, FK | Author (reference to users) |
| `created_at` | TIMESTAMP | NOT NULL | Publication date |
| `updated_at` | TIMESTAMP | NOT NULL | Last modification date |

**Access Control:**

- **Users:** Can create, read, update, delete own articles
- **Admins/Owners:** Can create articles for any user
- **Public:** Can read active articles only

---

### 5. photos

**Purpose:** Photo management system

```sql
CREATE TABLE photos (
    id SERIAL PRIMARY KEY,
    image VARCHAR(500) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Indexes
CREATE INDEX idx_photos_created_at ON photos(created_at);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment photo ID |
| `image` | VARCHAR(500) | NOT NULL | Photo image URL |
| `created_at` | TIMESTAMP | NOT NULL | Upload date |
| `updated_at` | TIMESTAMP | NOT NULL | Last modification date |

**Access Control:**

- **Users:** Can read all photos, update/delete own photos
- **Admins/Owners:** Can create photos, full CRUD access

---

### 6. onboardings

**Purpose:** App onboarding content management

```sql
CREATE TABLE onboardings (
    id SERIAL PRIMARY KEY,
    image VARCHAR(500) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Indexes
CREATE INDEX idx_onboardings_is_active ON onboardings(is_active);
CREATE INDEX idx_onboardings_created_at ON onboardings(created_at);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment onboarding ID |
| `image` | VARCHAR(500) | NOT NULL | Onboarding image URL |
| `title` | VARCHAR(255) | NOT NULL | Onboarding title (1-255 chars) |
| `description` | TEXT | NOT NULL | Onboarding description |
| `is_active` | BOOLEAN | DEFAULT true | Visibility status |
| `created_at` | TIMESTAMP | NOT NULL | Creation date |
| `updated_at` | TIMESTAMP | NOT NULL | Last modification date |

**Access Control:**

- **Public:** Can read active onboarding content
- **Admins/Owners:** Full CRUD access

---

### 7. configs

**Purpose:** Dynamic application configuration

```sql
CREATE TABLE configs (
    key VARCHAR(128) PRIMARY KEY,
    value TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `key` | VARCHAR(128) | PRIMARY KEY | Configuration key |
| `value` | TEXT | NOT NULL | Configuration value (unlimited text) |
| `created_at` | TIMESTAMP | NOT NULL | Creation date |
| `updated_at` | TIMESTAMP | NOT NULL | Last modification date |

**Special Configuration Keys:**

- `tnc` - Terms & Conditions content
- `privacy_policy` - Privacy Policy content
- `app_version` - Current app version
- `maintenance_mode` - Maintenance status

**Access Control:**

- **Users:** Can read config values by key
- **Admins/Owners:** Full CRUD access to all configs

---

## 🔗 Database Relationships

### Entity Relationship Diagram (ERD)

```
┌─────────────┐       ┌─────────────────┐       ┌─────────────────┐
│    users    │──────▶│  bank_accounts  │       │ device_sessions │
│             │ 1:N   │                 │       │                 │
│ • id (PK)   │       │ • id (PK)       │       │ • id (PK)       │
│ • name      │       │ • user_id (FK)  │       │ • user_id (FK)  │
│ • phone     │       │ • account_number│       │ • session_token │
│ • mother_name│      │ • account_name  │       │ • refresh_token │
│ • pin_atm   │       │ • bank_name     │       │ • device_type   │
│ • role      │       │ • bank_code     │       │ • device_id     │
│ • avatar    │       │ • account_type  │       │ • provider      │
└─────────────┘       │ • is_active     │       │ • expires_at    │
       │               │ • is_primary    │       └─────────────────┘
       │               └─────────────────┘              ▲
       │                                                │
       │ 1:N                                           │ 1:N
       ▼                                                │
┌─────────────┐       ┌─────────────────┐              │
│   articles  │       │   onboardings   │              │
│             │       │                 │              │
│ • id (PK)   │       │ • id (PK)       │              │
│ • title     │       │ • image         │              │
│ • image     │       │ • title         │              │
│ • content   │       │ • description   │              │
│ • is_active │       │ • is_active     │              │
│ • user_id(FK)│      └─────────────────┘              │
└─────────────┘                                        │
                                                       │
┌─────────────┐       ┌─────────────────┐              │
│   photos    │       │     configs     │              │
│             │       │                 │              │
│ • id (PK)   │       │ • key (PK)      │              │
│ • image     │       │ • value         │              │
└─────────────┘       └─────────────────┘              │
                                                       │
└──────────────────────────────────────────────────────┘
```

### Relationship Summary

1. **users ↔ bank_accounts** (One-to-Many)
   - One user can have multiple bank accounts
   - Foreign key: `bank_accounts.user_id → users.id`
   - Cascade delete: Delete user → Delete all bank accounts

2. **users ↔ device_sessions** (One-to-Many)
   - One user can have multiple active sessions
   - Foreign key: `device_sessions.user_id → users.id`
   - Cascade delete: Delete user → Delete all sessions

3. **users ↔ articles** (One-to-Many)
   - One user can author multiple articles
   - Foreign key: `articles.user_id → users.id`
   - Cascade delete: Delete user → Delete all articles

4. **Independent Tables**
   - `photos` - No foreign key relationships
   - `onboardings` - No foreign key relationships  
   - `configs` - No foreign key relationships

---

## 🔧 Database Configuration

### Connection Settings

```go
// Database connection configuration
type DatabaseConfig struct {
    Host     string // Default: localhost
    Port     string // Default: 5432
    User     string // PostgreSQL username
    Password string // PostgreSQL password
    DBName   string // Database name
    SSLMode  string // disable, require, verify-ca, verify-full
}
```

### Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=mbcdb
DB_SSLMODE=disable
```

### GORM Configuration

```go
// Auto-migration configuration
db.AutoMigrate(
    &User{},
    &BankAccount{},
    &DeviceSession{},
    &Article{},
    &Photo{},
    &Onboarding{},
    &Config{},
)
```

---

## 🔐 Security Considerations

### Data Protection

1. **Password Security**
   - PIN ATM stored with bcrypt hashing
   - Cost factor: 10 (default, secure for production)
   - No plain text storage of sensitive data

2. **Token Security**
   - JWT tokens stored in database for validation
   - Unique constraints on session and refresh tokens
   - Token expiration enforced at database level

3. **Phone Number Privacy**
   - Phone numbers are unique identifiers
   - Used for OTP delivery in banking authentication
   - Consider data masking in production logs

### Database Security

1. **Connection Security**
   - SSL/TLS encryption configurable
   - Connection pooling for performance
   - Prepared statements prevent SQL injection

2. **Access Control**
   - Role-based access control implemented in application
   - Foreign key constraints maintain referential integrity
   - Cascade deletes prevent orphaned records

3. **Audit Trail**
   - All tables include `created_at` and `updated_at`
   - Session management tracks `last_activity`
   - IP address logging for security monitoring

---

## 📊 Performance Optimization

### Database Indexes

**Primary Indexes (Auto-created):**

- All `id` fields (Primary Keys)
- `users.phone` (Unique)
- `configs.key` (Primary Key)

**Secondary Indexes:**

```sql
-- User management
CREATE INDEX idx_users_role ON users(role);

-- Bank account management  
CREATE INDEX idx_bank_accounts_user_id ON bank_accounts(user_id);
CREATE UNIQUE INDEX idx_user_account ON bank_accounts(user_id, account_number);

-- Session management
CREATE INDEX idx_device_sessions_user_id ON device_sessions(user_id);
CREATE INDEX idx_device_sessions_device_id ON device_sessions(device_id);
CREATE UNIQUE INDEX idx_session_token ON device_sessions(session_token);
CREATE UNIQUE INDEX idx_refresh_token ON device_sessions(refresh_token);

-- Content management
CREATE INDEX idx_articles_user_id ON articles(user_id);
CREATE INDEX idx_articles_is_active ON articles(is_active);
CREATE INDEX idx_articles_created_at ON articles(created_at);

-- Onboarding management
CREATE INDEX idx_onboardings_is_active ON onboardings(is_active);
CREATE INDEX idx_onboardings_created_at ON onboardings(created_at);

-- Photo management
CREATE INDEX idx_photos_created_at ON photos(created_at);
```

### Query Optimization

1. **Pagination Support**
   - All list endpoints support OFFSET/LIMIT
   - Indexed on commonly filtered fields (`is_active`, `created_at`)

2. **Relationship Loading**
   - GORM preloading for related data
   - Selective loading to avoid N+1 queries

3. **Connection Pooling**
   - Configured max open/idle connections
   - Connection lifetime management

---

## 🧪 Testing & Migration

### Database Migration

```bash
# Automatic migration on app start
go run main.go

# Manual migration utility
go run cmd/migrate/main.go
```

### Test Data Setup

```sql
-- Create test user
INSERT INTO users (name, phone, mother_name, pin_atm, role) 
VALUES ('Test User', '081234567890', 'Test Mother', '$2a$10$...', 'user');

-- Create test bank account
INSERT INTO bank_accounts (user_id, account_number, account_name, bank_name, is_primary)
VALUES (1, '1234567890123456', 'Test User', 'Test Bank', true);

-- Create test config
INSERT INTO configs (key, value) 
VALUES ('tnc', 'Terms and conditions content');
```

### Database Cleanup

```sql
-- Clean expired sessions
DELETE FROM device_sessions WHERE expires_at < NOW();

-- Clean inactive records
UPDATE articles SET is_active = false WHERE updated_at < NOW() - INTERVAL '1 year';
UPDATE onboardings SET is_active = false WHERE updated_at < NOW() - INTERVAL '1 year';
```

---

## 📈 Monitoring & Maintenance

### Database Health Checks

1. **Connection Monitoring**
   - Health check endpoint: `GET /health`
   - Database connectivity verification
   - Response time monitoring

2. **Performance Metrics**
   - Query execution time tracking
   - Connection pool utilization
   - Index usage statistics

3. **Storage Management**
   - Regular VACUUM operations for PostgreSQL
   - Index maintenance and optimization
   - Log rotation and cleanup

### Backup Strategy

```bash
# Full database backup
pg_dump -h localhost -U username -d mbcdb > backup_$(date +%Y%m%d).sql

# Restore from backup
psql -h localhost -U username -d mbcdb < backup_20250730.sql
```

---

## 🔧 Development Guidelines

### Schema Changes

1. **Migration Best Practices**
   - Always test migrations on development data first
   - Use transactions for complex migrations
   - Backup before production migrations

2. **Backward Compatibility**
   - Add columns as nullable initially
   - Use default values for new required fields
   - Deprecate columns before removing

3. **Version Control**
   - Track schema changes in version control
   - Document breaking changes
   - Coordinate with API version changes

### Data Integrity

1. **Constraints**
   - Use foreign key constraints for relationships
   - Implement unique constraints for business rules
   - Add check constraints for data validation

2. **Validation**
   - Application-level validation for complex rules
   - Database-level constraints for critical rules
   - Regular data integrity checks

---

**Last Updated:** July 30, 2025  
**Database Version:** PostgreSQL 12+  
**ORM Version:** GORM v1.25+  
**Tables:** 7 core tables  
**Total Relationships:** 3 foreign key relationships
