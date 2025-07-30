# MBankingCore Migration Guide

## Overview

This is a **NEW PROJECT** - no complex migrations are needed. The system uses GORM AutoMigrate for clean database setup.

## 🚀 Quick Setup

For a new project, simply run:

```bash
# Option 1: Run migrations only
go run cmd/migrate/main.go

# Option 2: Start main application (includes auto-migration)
go run main.go
```

## 📋 What Gets Created

### Database Tables
The following tables are automatically created:

1. **users** - User accounts with banking authentication
2. **admins** - Admin accounts with administrative privileges  
3. **bank_accounts** - Multi-account banking support
4. **device_sessions** - Multi-device session management
5. **otp_sessions** - Temporary OTP session data
6. **articles** - Content management articles
7. **photos** - Photo management system
8. **onboardings** - App onboarding content
9. **configs** - Dynamic application configuration

### Initial Data
The system automatically seeds:

#### Default Admin Account
- **Email:** admin@mbankingcore.com
- **Password:** admin123
- **Role:** super (Super Administrator)

#### Default Configuration
- App name and version
- Terms & conditions
- Privacy policy
- Contact information
- System settings

#### Default Onboarding Slides
- Welcome screen
- Security features
- Transaction guide
- Support information

## 🔧 Configuration

### Environment Variables
Ensure your `.env` file contains:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mbankingcore
DB_USERNAME=your_username
DB_PASSWORD=your_password
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_REFRESH_SECRET=your-super-secret-refresh-key

# OTP Configuration
OTP_LENGTH=6
OTP_EXPIRY_MINUTES=5

# Server Configuration
PORT=8080
```

### Database Requirements
- **PostgreSQL 12+** (recommended)
- Database and user with full privileges
- UTF-8 encoding

## 🧪 Testing Setup

For development and testing:

```bash
# Create test database
createdb mbankingcore_test

# Run with test environment
DB_NAME=mbankingcore_test go run main.go
```

## 📝 Notes

### For New Projects
- ✅ No complex migrations needed
- ✅ Clean database setup with GORM AutoMigrate
- ✅ Automatic initial data seeding
- ✅ Ready-to-use admin account

### Database Changes
If you need to modify the database structure:

1. Update the model in `models/` directory
2. GORM will automatically handle the schema changes
3. For data migrations, add functions to `config/migrations.go`

### Security Notes
- 🔐 Change default admin password after first login
- 🔐 Use strong JWT secrets in production
- 🔐 Enable SSL for database connections in production

## 🎯 Production Deployment

For production:

1. **Database Setup:**
   ```sql
   CREATE DATABASE mbankingcore;
   CREATE USER mbankingcore_user WITH PASSWORD 'strong_password';
   GRANT ALL PRIVILEGES ON DATABASE mbankingcore TO mbankingcore_user;
   ```

2. **Environment Configuration:**
   - Set strong JWT secrets
   - Enable database SSL
   - Configure proper database credentials

3. **First Run:**
   ```bash
   go run cmd/migrate/main.go
   ```

4. **Security:**
   - Change default admin password immediately
   - Review and update configuration values
   - Set up proper backup procedures

## ✅ Verification

After setup, verify the installation:

1. **Check Database Tables:**
   ```sql
   \dt  -- List all tables
   ```

2. **Check Admin Account:**
   ```sql
   SELECT * FROM admins;
   ```

3. **Check Configuration:**
   ```sql
   SELECT * FROM configs;
   ```

4. **Test API:**
   ```bash
   curl http://localhost:8080/health
   ```

## 🆘 Troubleshooting

### Common Issues

**Connection Failed:**
- Check database credentials
- Ensure PostgreSQL is running
- Verify network connectivity

**Permission Denied:**
- Grant proper database privileges
- Check user permissions

**Migration Errors:**
- Clear database and retry
- Check for conflicting data
- Review error logs

---

**Status:** ✅ **READY FOR NEW PROJECT**

No complex migrations required - just run and go! 🚀
