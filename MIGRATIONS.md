# MBX Backend - Migration System

This project uses an integrated migration system that automatically sets up the database and initial data when the application starts.

## How It Works

### Automatic Migration on Startup
When you run the main application (`go run main.go`), the system automatically:

1. **Connects to the database**
2. **Runs all migrations** (table creation, schema updates)
3. **Seeds initial data** (default configurations, sample content)
4. **Starts the API server**

### Manual Migration Tool
For database setup without starting the server, use the migration tool:

```bash
# Build the migration tool
go build -o migrate ./cmd/migrate

# Run migrations only
./migrate
```

## Migration Components

### 1. Auto-Migration
- Creates all database tables based on model structures
- Updates existing tables when models change
- Handles: Users, DeviceSessions, Articles, Onboarding, Photos, Config

### 2. Custom Migrations
- **User Roles**: Ensures all users have proper role assignments
- **Data Consistency**: Fixes any data integrity issues
- **Schema Updates**: Handles complex schema changes

### 3. Initial Data Seeding

#### Configuration Data
- `app_version`: Application version info
- `tnc`: Default Terms and Conditions content
- `privacy-policy`: Default Privacy Policy content  
- `maintenance_mode`: Maintenance flag
- `max_upload_size`: File upload size limits

#### Onboarding Content
- Welcome slide
- Authentication features
- API management info
- Getting started guide

## Database Configuration

The system uses the following environment variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=mbankingcore
DB_SSLMODE=disable
```

## Migration Safety

### First Time Setup
- ✅ **Safe**: Creates all tables and initial data
- ✅ **Idempotent**: Can be run multiple times safely
- ✅ **Non-destructive**: Never deletes existing data

### Updates
- ✅ **Automatic**: Runs on every application start
- ✅ **Smart**: Only creates missing data, skips existing
- ✅ **Logged**: All actions are logged for debugging

## File Structure

```
config/
├── database.go    # Database connection and migration trigger
└── migrations.go  # Migration logic and data seeding

cmd/
└── migrate/
    └── main.go    # Standalone migration tool

models/            # Database models (auto-migrated)
├── user.go
├── onboarding.go
├── config.go
└── ...
```

## Usage Examples

### Development Setup
```bash
# 1. Setup environment
cp .env.example .env
# Edit .env with your database credentials

# 2. Start application (migrations run automatically)
go run main.go
```

### Production Deployment
```bash
# 1. Run migrations first
go build -o migrate ./cmd/migrate
./migrate

# 2. Start application
go build -o mbankingcore .
./mbankingcore
```

### Database Reset (Development)
```bash
# Drop database, recreate, and run migrations
dropdb mbankingcore
createdb mbankingcore
go run main.go
```

## Troubleshooting

### Migration Fails
- Check database credentials in `.env`
- Ensure PostgreSQL is running
- Verify database exists and is accessible

### Duplicate Data
- The system is idempotent - safe to re-run
- Existing data is preserved
- Only missing data is created

### Schema Issues
- GORM AutoMigrate handles most schema changes
- For complex changes, add custom migration functions
- Backup database before major updates

## Adding New Migrations

To add new migration logic:

1. **Edit `config/migrations.go`**
2. **Add function to `runCustomMigrations()`**
3. **Test with migration tool**: `./migrate`
4. **Deploy**: Migrations run automatically on startup

Example:
```go
func runCustomMigrations() error {
    // Existing migrations...
    
    // Add new migration
    if err := migrateNewFeature(); err != nil {
        return err
    }
    
    return nil
}

func migrateNewFeature() error {
    // Your migration logic here
    return nil
}
```

## Benefits

✅ **Zero-Config**: Works out of the box  
✅ **Developer Friendly**: No manual SQL scripts  
✅ **Production Ready**: Safe for production deployments  
✅ **Maintainable**: All migration logic in one place  
✅ **Flexible**: Easy to add new migrations  
✅ **Logged**: Complete visibility into migration process
