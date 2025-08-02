package main

import (
	"log"
	"mbankingcore/config"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	envPath := filepath.Join("..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("No .env file found at %s: %v", envPath, err)
	}

	// Initialize database connection
	config.ConnectDatabase()

	log.Println("🔄 Running soft delete migration...")

	// Add DeletedAt columns to users and admins tables if they don't exist
	err := config.DB.Exec(`
		DO $$
		BEGIN
			-- Add deleted_at column to users table if it doesn't exist
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name = 'users' AND column_name = 'deleted_at'
			) THEN
				ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
				CREATE INDEX idx_users_deleted_at ON users(deleted_at);
				RAISE NOTICE 'Added deleted_at column to users table';
			ELSE
				RAISE NOTICE 'deleted_at column already exists in users table';
			END IF;

			-- Add deleted_at column to admins table if it doesn't exist
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name = 'admins' AND column_name = 'deleted_at'
			) THEN
				ALTER TABLE admins ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
				CREATE INDEX idx_admins_deleted_at ON admins(deleted_at);
				RAISE NOTICE 'Added deleted_at column to admins table';
			ELSE
				RAISE NOTICE 'deleted_at column already exists in admins table';
			END IF;
		END
		$$;
	`).Error

	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Soft delete migration completed successfully!")
	log.Println("📋 Summary:")
	log.Println("   - Added deleted_at column to users table (if not exists)")
	log.Println("   - Added deleted_at column to admins table (if not exists)")
	log.Println("   - Created indexes for soft delete queries")
	log.Println("   - Existing data remains unchanged")
}
