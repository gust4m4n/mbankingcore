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
| `admins` | Admin accounts with administrative privileges | Dynamic | `id` (uint) |
| `bank_accounts` | Multi-account banking support | Dynamic | `id` (uint) |
| `device_sessions` | Multi-device session management | Dynamic | `id` (uint) |
| `otp_sessions` | Temporary OTP session data for banking login | Dynamic | `id` (uint) |
| `transactions` | Transaction history (topup, withdraw, transfer) | Dynamic | `id` (uint) |
| `articles` | Content management articles | Dynamic | `id` (uint) |
| `photos` | Photo management system | Dynamic | `id` (uint) |
| `onboardings` | App onboarding content | Dynamic | `id` (uint) |
| `configs` | Dynamic application configuration | Dynamic | `key` (string) |
| `audit_logs` | Comprehensive system activity audit trail | Dynamic | `id` (uint) |
| `login_audits` | Authentication and login activity tracking | Dynamic | `id` (uint) |

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
    balance BIGINT DEFAULT 0,  -- User account balance (integer)
    status INTEGER DEFAULT 1,  -- User status: 0=inactive, 1=active, 2=terblokir
    avatar VARCHAR(500),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Indexes
CREATE UNIQUE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_balance ON users(balance);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment user ID |
| `name` | VARCHAR(255) | NOT NULL | Full name (min 8 characters) |
| `phone` | VARCHAR(255) | UNIQUE, NOT NULL | Phone number (unique identifier) |
| `mother_name` | VARCHAR(255) | NOT NULL | Mother's name (min 8 characters) |
| `pin_atm` | VARCHAR(255) | NOT NULL | Hashed PIN ATM (6 digits, bcrypt) |
| `balance` | BIGINT | DEFAULT 0 | User account balance (integer amount) |
| `status` | INTEGER | DEFAULT 1 | User status: 0=inactive, 1=active, 2=terblokir |
| `avatar` | VARCHAR(500) | NULLABLE | Avatar image URL |
| `created_at` | TIMESTAMP | NOT NULL | Record creation time |
| `updated_at` | TIMESTAMP | NOT NULL | Last update time |

**Status Values:**

- `0` - Inactive user (cannot login)
- `1` - Active user (default, normal operation)
- `2` - Terblokir (blocked, cannot perform transactions)

**Relationships:**

- **One-to-Many** with `bank_accounts` (via `user_id`)
- **One-to-Many** with `device_sessions` (via `user_id`)
- **One-to-Many** with `articles` (via `user_id`)

---

### 2. admins

**Purpose:** Administrative accounts with elevated privileges for system management

```sql
CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,  -- bcrypt hashed
    role VARCHAR(20) DEFAULT 'admin',  -- admin, super
    status INTEGER DEFAULT 1,  -- 0=inactive, 1=active, 2=blocked
    avatar VARCHAR(500),
    last_login TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Indexes
CREATE UNIQUE INDEX idx_admins_email ON admins(email);
CREATE INDEX idx_admins_role ON admins(role);
CREATE INDEX idx_admins_status ON admins(status);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment admin ID |
| `name` | VARCHAR(255) | NOT NULL | Full name of administrator |
| `email` | VARCHAR(255) | UNIQUE, NOT NULL | Email address (unique identifier) |
| `password` | VARCHAR(255) | NOT NULL | Hashed password (bcrypt) |
| `role` | VARCHAR(20) | DEFAULT 'admin' | Admin role: 'admin', 'super' |
| `status` | INTEGER | DEFAULT 1 | Admin status: 0=inactive, 1=active, 2=blocked |
| `avatar` | VARCHAR(500) | NULLABLE | Avatar image URL |
| `last_login` | TIMESTAMP | NULLABLE | Last login timestamp |
| `created_at` | TIMESTAMP | NOT NULL | Record creation time |
| `updated_at` | TIMESTAMP | NOT NULL | Last update time |

**Role Values:**

- `admin` - Standard administrator (default)
- `super` - Super administrator (full system access)

**Status Values:**

- `0` - Inactive admin (cannot login)
- `1` - Active admin (default, normal operation)
- `2` - Blocked admin (cannot perform any actions)

**Relationships:**

- **Independent table** - No direct foreign key relationships
- **Functional relationships** with all other tables through admin operations

---

### 3. bank_accounts

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

### 4. device_sessions

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

### 5. otp_sessions

**Purpose:** Temporary storage for OTP session data during banking login process

```sql
CREATE TABLE otp_sessions (
    id SERIAL PRIMARY KEY,
    login_token VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(255) NOT NULL,
    otp_code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    account_number VARCHAR(255) NOT NULL,
    mother_name VARCHAR(255) NOT NULL,
    pin_atm VARCHAR(255) NOT NULL,
    device_type VARCHAR(255) NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    device_name VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Indexes
CREATE UNIQUE INDEX idx_otp_sessions_login_token ON otp_sessions(login_token);
CREATE INDEX idx_otp_sessions_phone ON otp_sessions(phone);
CREATE INDEX idx_otp_sessions_expires_at ON otp_sessions(expires_at);
CREATE INDEX idx_otp_sessions_is_used ON otp_sessions(is_used);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment session ID |
| `login_token` | VARCHAR(255) | UNIQUE, NOT NULL | Unique login token for verification |
| `phone` | VARCHAR(255) | NOT NULL | Phone number for login |
| `otp_code` | VARCHAR(255) | NOT NULL | OTP code (hidden from JSON) |
| `name` | VARCHAR(255) | NOT NULL | User name |
| `account_number` | VARCHAR(255) | NOT NULL | Account number |
| `mother_name` | VARCHAR(255) | NOT NULL | Mother's name |
| `pin_atm` | VARCHAR(255) | NOT NULL | PIN ATM (hidden from JSON) |
| `device_type` | VARCHAR(255) | NOT NULL | Device type |
| `device_id` | VARCHAR(255) | NOT NULL | Device identifier |
| `device_name` | VARCHAR(255) | NOT NULL | Device name |
| `expires_at` | TIMESTAMP | NOT NULL | Session expiration time |
| `is_used` | BOOLEAN | DEFAULT FALSE | Whether session has been used |
| `created_at` | TIMESTAMP | NOT NULL | Record creation time |
| `updated_at` | TIMESTAMP | NOT NULL | Last update time |

**Relationships:**

- **Independent table** - Temporary storage, no permanent relationships
- **Functional relationship** with `users` during login verification process

**Key Features:**

- Temporary storage for 2-step banking authentication
- Automatic expiration and cleanup
- Security through unique login tokens
- Device tracking for multi-device support

---

### 6. transactions

**Purpose:** Complete transaction history for all balance operations

```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL,  -- 'topup', 'withdraw', 'transfer_out', 'transfer_in'
    amount BIGINT NOT NULL,  -- Transaction amount (integer)
    balance_before BIGINT NOT NULL,  -- Balance before transaction
    balance_after BIGINT NOT NULL,  -- Balance after transaction
    status VARCHAR(50) DEFAULT 'completed',  -- 'completed', 'failed', 'pending'
    description TEXT,  -- Transaction description
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_transactions_user_type ON transactions(user_id, type);
```

**Field Details:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Unique transaction identifier |
| `user_id` | INTEGER | NOT NULL, FK | Reference to users table |
| `type` | VARCHAR(50) | NOT NULL | Transaction type |
| `amount` | BIGINT | NOT NULL | Transaction amount (positive integer) |
| `balance_before` | BIGINT | NOT NULL | User balance before transaction |
| `balance_after` | BIGINT | NOT NULL | User balance after transaction |
| `status` | VARCHAR(50) | DEFAULT 'completed' | Transaction status |
| `description` | TEXT | NULLABLE | Optional transaction description |
| `created_at` | TIMESTAMP | NOT NULL | Transaction timestamp |
| `updated_at` | TIMESTAMP | NOT NULL | Last update time |

**Relationships:**

- **Foreign Key** → `users(id)` - Each transaction belongs to a user
- **Cascade Delete** - Transactions deleted when user is deleted

**Transaction Types:**

- `topup` - Balance addition operation
- `withdraw` - Balance deduction operation
- `transfer_out` - Outgoing transfer (sender side)
- `transfer_in` - Incoming transfer (receiver side)

**Transaction Status:**

- `completed` - Successfully processed transaction
- `failed` - Failed transaction (future implementation)
- `pending` - Processing transaction (future implementation)

**Key Features:**

- Complete audit trail for all balance changes
- Atomic transaction processing with balance tracking
- Support for multiple transaction types
- Optimized indexes for common query patterns
- Foreign key constraints ensure data integrity

**Business Rules:**

- All amounts stored as positive integers (smallest currency unit)
- Balance before/after provides complete audit trail
- Transfer operations create dual records (sender + receiver)
- Immutable records - transactions cannot be modified after creation

---

### 7. articles

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

### 7. photos

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

### 8. onboardings

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

### 9. configs

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

### 10. audit_logs

**Purpose:** Comprehensive audit trail for all system activities and changes

```sql
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    admin_id INTEGER,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    action VARCHAR(20) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent VARCHAR(255),
    api_endpoint VARCHAR(255),
    request_method VARCHAR(10),
    status_code INTEGER,
    created_at TIMESTAMP NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE SET NULL
);

-- Indexes
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_admin_id ON audit_logs(admin_id);
CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_ip_address ON audit_logs(ip_address);
CREATE INDEX idx_audit_logs_composite ON audit_logs(entity_type, action, created_at DESC);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment audit log ID |
| `user_id` | INTEGER | NULLABLE, FK | User who performed the action |
| `admin_id` | INTEGER | NULLABLE, FK | Admin who performed the action |
| `entity_type` | VARCHAR(50) | NOT NULL | Type of entity (user, transaction, etc.) |
| `entity_id` | INTEGER | NOT NULL | ID of the affected entity |
| `action` | VARCHAR(20) | NOT NULL | Action performed (CREATE, READ, UPDATE, DELETE) |
| `old_values` | JSONB | NULLABLE | Data before change (JSON format) |
| `new_values` | JSONB | NULLABLE | Data after change (JSON format) |
| `ip_address` | INET | NULLABLE | Client IP address |
| `user_agent` | VARCHAR(255) | NULLABLE | Client user agent |
| `api_endpoint` | VARCHAR(255) | NULLABLE | API endpoint called |
| `request_method` | VARCHAR(10) | NULLABLE | HTTP method (GET, POST, PUT, DELETE) |
| `status_code` | INTEGER | NULLABLE | HTTP response status code |
| `created_at` | TIMESTAMP | NOT NULL | When the action occurred |

**Entity Types:**

- `user` - User account operations
- `transaction` - Transaction operations
- `bank_account` - Bank account operations
- `admin` - Admin operations
- `article` - Article operations
- `photo` - Photo operations
- `config` - Configuration changes
- `auth` - Authentication operations
- `session` - Session management

**Action Types:**

- `CREATE` - New record creation
- `READ` - Data retrieval/viewing
- `UPDATE` - Record modification
- `DELETE` - Record deletion
- `LOGIN` - Authentication success
- `LOGOUT` - Session termination

---

### 11. login_audits

**Purpose:** Specialized audit trail for authentication activities

```sql
CREATE TABLE login_audits (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    admin_id INTEGER,
    login_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    ip_address INET,
    user_agent VARCHAR(255),
    device_info JSONB,
    failure_reason VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE SET NULL
);

-- Indexes
CREATE INDEX idx_login_audits_user_id ON login_audits(user_id);
CREATE INDEX idx_login_audits_admin_id ON login_audits(admin_id);
CREATE INDEX idx_login_audits_login_type ON login_audits(login_type);
CREATE INDEX idx_login_audits_status ON login_audits(status);
CREATE INDEX idx_login_audits_created_at ON login_audits(created_at DESC);
CREATE INDEX idx_login_audits_ip_address ON login_audits(ip_address);
CREATE INDEX idx_login_audits_composite ON login_audits(login_type, status, created_at DESC);
```

**Fields:**

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment login audit ID |
| `user_id` | INTEGER | NULLABLE, FK | User account involved |
| `admin_id` | INTEGER | NULLABLE, FK | Admin account involved |
| `login_type` | VARCHAR(20) | NOT NULL | Type of login attempt |
| `status` | VARCHAR(20) | NOT NULL | Result of login attempt |
| `ip_address` | INET | NULLABLE | Client IP address |
| `user_agent` | VARCHAR(255) | NULLABLE | Client user agent |
| `device_info` | JSONB | NULLABLE | Device information (JSON format) |
| `failure_reason` | VARCHAR(255) | NULLABLE | Reason for failed attempt |
| `created_at` | TIMESTAMP | NOT NULL | When the attempt occurred |

**Login Types:**

- `user_login` - User authentication attempt
- `admin_login` - Admin authentication attempt
- `user_logout` - User session termination
- `admin_logout` - Admin session termination

**Status Values:**

- `success` - Authentication successful
- `failed` - Authentication failed
- `blocked` - Account blocked/suspended

**Key Features:**

- **Complete Authentication Audit** - All login/logout attempts tracked
- **Security Monitoring** - Failed attempts and IP tracking
- **Device Information** - Device fingerprinting for security
- **Forensic Analysis** - Complete audit trail for security incidents
- **Automatic Logging** - Middleware-based logging for all auth events

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
│ • balance   │       │ • bank_code     │       │ • device_id     │
│ • status    │       │ • account_type  │       │ • provider      │
└─────────────┘       │ • is_active     │       │ • expires_at    │
       │               │ • is_primary    │       └─────────────────┘
       │               └─────────────────┘              ▲
       │                                                │
       │ 1:N                                           │ 1:N
       ▼                                                │
┌─────────────┐       ┌─────────────────┐              │
│transactions │       │     admins      │              │
│             │       │                 │              │
│ • id (PK)   │       │ • id (PK)       │              │
│ • user_id(FK)│      │ • name          │              │
│ • type      │       │ • email         │              │
│ • amount    │       │ • password      │              │
│ • balance_before│   │ • role          │              │
│ • balance_after │   │ • status        │              │
│ • description│      └─────────────────┘              │
└─────────────┘                 │                     │
                               │ 1:N                  │
                               ▼                       │
┌─────────────┐       ┌─────────────────┐              │
│   articles  │       │   audit_logs    │              │
│             │       │                 │              │
│ • id (PK)   │       │ • id (PK)       │              │
│ • title     │       │ • user_id (FK)  │◄─────────────┘
│ • image     │       │ • admin_id (FK) │◄─────────────┐
│ • content   │       │ • entity_type   │              │
│ • is_active │       │ • entity_id     │              │
│ • user_id(FK)│      │ • action        │              │
└─────────────┘       │ • old_values    │              │
                      │ • new_values    │              │
┌─────────────┐       │ • ip_address    │              │
│   photos    │       │ • user_agent    │              │
│             │       │ • api_endpoint  │              │
│ • id (PK)   │       │ • request_method│              │
│ • image     │       │ • status_code   │              │
│ • user_id(FK)│      └─────────────────┘              │
└─────────────┘                 │                     │
                               │                       │
┌─────────────┐       ┌─────────────────┐              │
│ onboardings │       │  login_audits   │              │
│             │       │                 │              │
│ • id (PK)   │       │ • id (PK)       │              │
│ • image     │       │ • user_id (FK)  │◄─────────────┘
│ • title     │       │ • admin_id (FK) │◄─────────────┐
│ • description│      │ • login_type    │              │
│ • is_active │       │ • status        │              │
└─────────────┘       │ • ip_address    │              │
                      │ • user_agent    │              │
┌─────────────┐       │ • device_info   │              │
│   configs   │       │ • failure_reason│              │
│             │       └─────────────────┘              │
│ • key (PK)  │                                        │
│ • value     │                                        │
└─────────────┘                                        │
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

3. **users ↔ transactions** (One-to-Many)
   - One user can have multiple transactions
   - Foreign key: `transactions.user_id → users.id`
   - Cascade delete: Delete user → Delete all transactions

4. **users ↔ articles** (One-to-Many)
   - One user can author multiple articles
   - Foreign key: `articles.user_id → users.id`
   - Cascade delete: Delete user → Delete all articles

5. **users ↔ photos** (One-to-Many)
   - One user can upload multiple photos
   - Foreign key: `photos.user_id → users.id`
   - Cascade delete: Delete user → Delete all photos

6. **users ↔ audit_logs** (One-to-Many)
   - One user can have multiple audit log entries
   - Foreign key: `audit_logs.user_id → users.id`
   - Set NULL on delete: Delete user → Set user_id to NULL

7. **admins ↔ audit_logs** (One-to-Many)
   - One admin can have multiple audit log entries
   - Foreign key: `audit_logs.admin_id → admins.id`
   - Set NULL on delete: Delete admin → Set admin_id to NULL

8. **users ↔ login_audits** (One-to-Many)
   - One user can have multiple login audit entries
   - Foreign key: `login_audits.user_id → users.id`
   - Set NULL on delete: Delete user → Set user_id to NULL

9. **admins ↔ login_audits** (One-to-Many)
   - One admin can have multiple login audit entries
   - Foreign key: `login_audits.admin_id → admins.id`
   - Set NULL on delete: Delete admin → Set admin_id to NULL

10. **Independent Tables**
    - `otp_sessions` - Temporary data, no permanent relationships
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
    &Admin{},
    &BankAccount{},
    &DeviceSession{},
    &OTPSession{},
    &Transaction{},
    &Article{},
    &Photo{},
    &Onboarding{},
    &Config{},
    &AuditLog{},
    &LoginAudit{},
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
CREATE INDEX idx_users_balance ON users(balance);
CREATE INDEX idx_users_status ON users(status);

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

## 🔍 Audit Trails & Compliance

### Comprehensive Audit System

The MBankingCore database includes a sophisticated audit trail system designed for compliance, security monitoring, and forensic analysis.

#### 1. Activity Audit Logs (`audit_logs`)

**Automatic Logging:**
- All CRUD operations across all entities
- API endpoint calls and HTTP methods
- Request/response status codes
- IP addresses and user agents
- Complete before/after data snapshots (JSON)

**Filtering Capabilities:**
- Entity type (user, transaction, admin, etc.)
- Action type (CREATE, READ, UPDATE, DELETE)
- User or Admin ID
- Date range filtering
- IP address tracking
- Combined filters with pagination

**Use Cases:**
- Compliance reporting and audits
- Security incident investigation
- Change tracking and rollback support
- Performance monitoring and analytics
- Regulatory compliance (PCI DSS, SOX, etc.)

#### 2. Authentication Audit (`login_audits`)

**Login/Logout Tracking:**
- All authentication attempts (successful/failed)
- Device information and fingerprinting
- IP address and geolocation data
- Failure reasons for security analysis
- Session duration tracking

**Security Features:**
- Brute force attack detection
- Suspicious IP monitoring
- Device authentication patterns
- Failed login attempt analysis
- Administrative access monitoring

#### 3. Compliance Features

**Data Retention:**
```sql
-- Audit log retention (example: 7 years)
DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL '7 years';
DELETE FROM login_audits WHERE created_at < NOW() - INTERVAL '7 years';
```

**Performance Optimization:**
```sql
-- Partitioning for large audit tables (PostgreSQL 10+)
CREATE TABLE audit_logs_2025 PARTITION OF audit_logs
FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

-- Archive old audit data
CREATE TABLE audit_logs_archive AS
SELECT * FROM audit_logs WHERE created_at < '2024-01-01';
```

**GDPR Compliance:**
- User data anonymization support
- Right to be forgotten implementation
- Data export capabilities
- Consent tracking through audit logs

#### 4. Monitoring & Alerting

**Critical Events:**
- Failed admin login attempts
- Unusual transaction patterns
- Bulk data modifications
- System configuration changes
- Security policy violations

**Automated Alerts:**
```sql
-- Example: Monitor failed login attempts
SELECT ip_address, COUNT(*) as failed_attempts
FROM login_audits 
WHERE status = 'failed' 
  AND created_at > NOW() - INTERVAL '1 hour'
GROUP BY ip_address
HAVING COUNT(*) >= 5;
```

### Best Practices

1. **Regular Audit Reviews**
   - Weekly security audit log analysis
   - Monthly compliance report generation
   - Quarterly data retention cleanup
   - Annual audit trail verification

2. **Access Control**
   - Admin-only access to audit endpoints
   - Role-based audit data filtering
   - Secure audit data transmission
   - Encrypted audit data storage

3. **Performance Considerations**
   - Asynchronous audit logging
   - Batch processing for bulk operations
   - Indexed columns for fast queries
   - Partitioned tables for scalability

---

**Last Updated:** July 30, 2025  
**Database Version:** PostgreSQL 12+  
**ORM Version:** GORM v1.25+  
**Tables:** 12 core tables (including audit trails)  
**Total Relationships:** 9 foreign key relationships  
**Audit Features:** Comprehensive activity & authentication tracking
