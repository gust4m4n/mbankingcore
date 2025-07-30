package config

import (
	"log"
	"mbankingcore/models"

	"golang.org/x/crypto/bcrypt"
)

// RunMigrations handles database migrations and initial setup
func RunMigrations() error {
	log.Println("Starting database migrations...")

	// Step 0: Run pre-migration cleanup
	if err := cleanupOTPSessions(); err != nil {
		log.Printf("Failed to run pre-migration cleanup: %v", err)
		return err
	}

	// Step 1: Auto-migrate all models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.BankAccount{},
		&models.DeviceSession{},
		&models.OTPSession{},
		&models.Article{},
		&models.Onboarding{},
		&models.Photo{},
		&models.Config{},
	)
	if err != nil {
		log.Printf("Failed to auto-migrate models: %v", err)
		return err
	}
	log.Println("✅ Auto-migration completed successfully")

	// Step 2: Run custom migrations
	if err := runCustomMigrations(); err != nil {
		log.Printf("Failed to run custom migrations: %v", err)
		return err
	}

	// Step 3: Seed initial data
	if err := seedInitialData(); err != nil {
		log.Printf("Failed to seed initial data: %v", err)
		return err
	}

	log.Println("🚀 All migrations and initial setup completed successfully!")
	return nil
}

// runCustomMigrations handles specific migration tasks
func runCustomMigrations() error {
	log.Println("Running custom migrations...")

	// Migration: Ensure user role column exists and has default values
	if err := migrateUserRoles(); err != nil {
		return err
	}

	// Migration: Remove user_agent column from device_sessions table
	if err := removeUserAgentColumn(); err != nil {
		return err
	}

	// Step 3: Remove email column
	if err := removeEmailColumn(); err != nil {
		return err
	}

	// Step 4: Migrate account numbers to separate table
	if err := migrateAccountNumbers(); err != nil {
		return err
	}

	log.Println("✅ Custom migrations completed")
	return nil
}

// migrateUserRoles - deprecated, roles removed from system
func migrateUserRoles() error {
	log.Println("Migrating user roles...")
	log.Println("✅ Role migration skipped - roles removed from system")
	return nil
}

// removeUserAgentColumn removes the user_agent column from device_sessions table
func removeUserAgentColumn() error {
	log.Println("Removing user_agent column from device_sessions table...")

	// Check if the column exists first
	var columnExists bool
	err := DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'device_sessions' AND column_name = 'user_agent')").Scan(&columnExists).Error
	if err != nil {
		log.Printf("Error checking if user_agent column exists: %v", err)
		return err
	}

	if columnExists {
		// Drop the user_agent column
		err = DB.Exec("ALTER TABLE device_sessions DROP COLUMN user_agent").Error
		if err != nil {
			log.Printf("Error dropping user_agent column: %v", err)
			return err
		}
		log.Println("✅ Successfully removed user_agent column from device_sessions table")
	} else {
		log.Println("✅ user_agent column does not exist, no action needed")
	}

	return nil
}

// removeEmailColumn removes the email and email_verified columns from users table and ensures phone is unique
func removeEmailColumn() error {
	log.Println("Removing email column from users table...")

	// Check if the email column exists first
	var emailColumnExists bool
	err := DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'email')").Scan(&emailColumnExists).Error
	if err != nil {
		log.Printf("Error checking if email column exists: %v", err)
		return err
	}

	if emailColumnExists {
		// First, ensure phone has unique constraint
		log.Println("Adding unique constraint to phone column...")
		err = DB.Exec("ALTER TABLE users ADD CONSTRAINT uni_users_phone UNIQUE (phone)").Error
		if err != nil {
			// If constraint already exists, that's fine
			log.Printf("Note: Unique constraint on phone may already exist: %v", err)
		}

		// Drop the email column
		err = DB.Exec("ALTER TABLE users DROP COLUMN email").Error
		if err != nil {
			log.Printf("Error dropping email column: %v", err)
			return err
		}
		log.Println("✅ Successfully removed email column from users table")

		// Check if email_verified column exists and remove it too
		var emailVerifiedExists bool
		err = DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'email_verified')").Scan(&emailVerifiedExists).Error
		if err == nil && emailVerifiedExists {
			err = DB.Exec("ALTER TABLE users DROP COLUMN email_verified").Error
			if err != nil {
				log.Printf("Error dropping email_verified column: %v", err)
				return err
			}
			log.Println("✅ Successfully removed email_verified column from users table")
		}
	} else {
		log.Println("✅ email column does not exist, no action needed")
	}

	return nil
}

// migrateAccountNumbers migrates account numbers from users table to bank_accounts table
func migrateAccountNumbers() error {
	log.Println("Migrating account numbers to separate bank_accounts table...")

	// Check if the account_number column exists in users table
	var accountColumnExists bool
	err := DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'account_number')").Scan(&accountColumnExists).Error
	if err != nil {
		log.Printf("Error checking if account_number column exists: %v", err)
		return err
	}

	if accountColumnExists {
		// Get all users with account numbers
		var users []struct {
			ID            uint
			Name          string
			AccountNumber string
		}

		err = DB.Table("users").Select("id, name, account_number").Where("account_number IS NOT NULL AND account_number != ''").Find(&users).Error
		if err != nil {
			log.Printf("Error fetching users with account numbers: %v", err)
			return err
		}

		// Create bank accounts for each user
		for _, user := range users {
			var existingAccount models.BankAccount
			err = DB.Where("user_id = ? AND account_number = ?", user.ID, user.AccountNumber).First(&existingAccount).Error
			if err != nil {
				// Create new bank account
				bankAccount := models.BankAccount{
					UserID:        user.ID,
					AccountNumber: user.AccountNumber,
					AccountName:   user.Name, // Use user's name as account name
					BankName:      "Unknown Bank",
					AccountType:   "saving",
					IsActive:      true,
					IsPrimary:     true, // First account is primary
				}

				err = DB.Create(&bankAccount).Error
				if err != nil {
					log.Printf("Error creating bank account for user %d: %v", user.ID, err)
					continue
				}
				log.Printf("✅ Created bank account for user %d with account number %s", user.ID, user.AccountNumber)
			}
		}

		// Drop the account_number column from users table
		err = DB.Exec("ALTER TABLE users DROP COLUMN account_number").Error
		if err != nil {
			log.Printf("Error dropping account_number column: %v", err)
			return err
		}
		log.Println("✅ Successfully removed account_number column from users table")
	} else {
		log.Println("✅ account_number column does not exist, no action needed")
	}

	return nil
}

// seedInitialData creates essential initial data for the application
func seedInitialData() error {
	log.Println("Seeding initial data...")

	// Seed initial admin users
	if err := seedInitialAdmins(); err != nil {
		return err
	}

	// Seed initial configuration values
	if err := seedInitialConfigs(); err != nil {
		return err
	}

	// Seed initial onboarding content
	if err := seedInitialOnboarding(); err != nil {
		return err
	}

	log.Println("✅ Initial data seeding completed")
	return nil
}

// seedInitialAdmins creates default admin users
func seedInitialAdmins() error {
	log.Println("Seeding initial admin users...")

	// Check if admin users already exist
	var count int64
	DB.Model(&models.Admin{}).Count(&count)

	if count > 0 {
		log.Println("✅ Admin users already exist")
		return nil
	}

	// Hash default password
	hashedPassword, err := hashPassword("admin123")
	if err != nil {
		log.Printf("Failed to hash default admin password: %v", err)
		return err
	}

	// Create default super admin
	superAdmin := models.Admin{
		Name:     "Super Admin",
		Email:    "admin@mbankingcore.com",
		Password: hashedPassword,
		Role:     models.ADMIN_ROLE_SUPER,
		Status:   models.ADMIN_STATUS_ACTIVE,
	}

	if err := DB.Create(&superAdmin).Error; err != nil {
		log.Printf("Failed to create super admin: %v", err)
		return err
	}

	log.Println("✅ Created default super admin (email: admin@mbankingcore.com, password: admin123)")
	log.Println("⚠️  IMPORTANT: Change the default admin password in production!")

	return nil
}

// hashPassword is a helper function for password hashing
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// hashPasswordBytes performs the actual bcrypt hashing
func hashPasswordBytes(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// hashPasswordWithCost performs bcrypt hashing with specified cost
func hashPasswordWithCost(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

// seedInitialConfigs creates essential configuration entries
func seedInitialConfigs() error {
	log.Println("Seeding initial configurations...")

	initialConfigs := []struct {
		Key   string
		Value string
		Desc  string
	}{
		{
			Key:   "app_version",
			Value: "1.0.0",
			Desc:  "Application version",
		},
		{
			Key:   "tnc",
			Value: "# Terms and Conditions\n\n## 1. Introduction\nWelcome to MBX Backend API. By using our service, you agree to comply with and be bound by the following terms and conditions.\n\n## 2. Use License\nPermission is granted to temporarily access our service for personal, non-commercial transitory viewing only.\n\n## 3. Disclaimer\nThe materials on our service are provided on an 'as is' basis. We make no warranties, expressed or implied.\n\n## 4. Limitations\nIn no event shall MBX Backend or its suppliers be liable for any damages arising out of the use or inability to use our service.\n\n## 5. Accuracy of Materials\nThe materials appearing on our service could include technical, typographical, or photographic errors.\n\n## 6. Contact Information\nFor questions about these Terms and Conditions, please contact our support team.\n\n---\n*Last updated: January 2024*",
			Desc:  "Terms and Conditions content",
		},
		{
			Key:   "privacy-policy",
			Value: "# Privacy Policy\n\n## 1. Information We Collect\nWe collect information you provide directly to us, such as when you create an account, make a purchase, or contact us.\n\n## 2. How We Use Your Information\nWe use the information we collect to provide, maintain, and improve our services.\n\n## 3. Information Sharing\nWe do not sell, trade, or otherwise transfer your personal information to third parties without your consent.\n\n## 4. Data Security\nWe implement appropriate security measures to protect your personal information.\n\n## 5. Your Rights\nYou have the right to access, update, or delete your personal information.\n\n## 6. Changes to This Policy\nWe may update this privacy policy from time to time. We will notify you of any changes.\n\n## 7. Contact Us\nIf you have any questions about this Privacy Policy, please contact us.\n\n---\n*Last updated: January 2024*",
			Desc:  "Privacy Policy content",
		},
		{
			Key:   "maintenance_mode",
			Value: "false",
			Desc:  "Application maintenance mode flag",
		},
		{
			Key:   "max_upload_size",
			Value: "10485760", // 10MB in bytes
			Desc:  "Maximum file upload size in bytes",
		},
	}

	for _, configData := range initialConfigs {
		var existingConfig models.Config
		err := DB.Where("key = ?", configData.Key).First(&existingConfig).Error

		if err != nil {
			if err.Error() == "record not found" {
				// Create new config
				newConfig := models.Config{
					Key:   configData.Key,
					Value: configData.Value,
				}

				if err := DB.Create(&newConfig).Error; err != nil {
					log.Printf("Failed to create config %s: %v", configData.Key, err)
					return err
				}
				log.Printf("✅ Created initial config: %s", configData.Key)
			} else {
				log.Printf("Error checking config %s: %v", configData.Key, err)
				return err
			}
		} else {
			log.Printf("✅ Config already exists: %s", configData.Key)
		}
	}

	return nil
}

// seedInitialOnboarding creates default onboarding content
func seedInitialOnboarding() error {
	log.Println("Seeding initial onboarding content...")

	// Check if onboarding content already exists
	var count int64
	DB.Model(&models.Onboarding{}).Count(&count)

	if count > 0 {
		log.Println("✅ Onboarding content already exists")
		return nil
	}

	// Create initial onboarding slides
	initialOnboarding := []models.Onboarding{
		{
			Title:       "Welcome to MBX Backend",
			Description: "Your comprehensive backend solution for modern applications",
			Image:       "https://via.placeholder.com/400x300/4F46E5/FFFFFF?text=Welcome",
			IsActive:    true,
		},
		{
			Title:       "Secure Authentication",
			Description: "Built-in JWT authentication system with device session management",
			Image:       "https://via.placeholder.com/400x300/7C3AED/FFFFFF?text=Security",
			IsActive:    true,
		},
		{
			Title:       "Easy API Management",
			Description: "RESTful APIs with comprehensive documentation and testing tools",
			Image:       "https://via.placeholder.com/400x300/059669/FFFFFF?text=API",
			IsActive:    true,
		},
		{
			Title:       "Get Started",
			Description: "Register your account and start building amazing applications",
			Image:       "https://via.placeholder.com/400x300/DC2626/FFFFFF?text=Start",
			IsActive:    true,
		},
	}

	for _, onboarding := range initialOnboarding {
		if err := DB.Create(&onboarding).Error; err != nil {
			log.Printf("Failed to create onboarding slide: %v", err)
			return err
		}
	}

	log.Printf("✅ Created %d initial onboarding slides", len(initialOnboarding))
	return nil
}

// cleanupOTPSessions removes all existing OTP sessions to allow adding the new name column
func cleanupOTPSessions() error {
	log.Println("Cleaning up OTP sessions table for new name column...")

	// Check if otp_sessions table exists
	var tableExists bool
	err := DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'otp_sessions')").Scan(&tableExists).Error
	if err != nil {
		log.Printf("Error checking if otp_sessions table exists: %v", err)
		return err
	}

	if tableExists {
		// Check if name column already exists
		var nameColumnExists bool
		err = DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'otp_sessions' AND column_name = 'name')").Scan(&nameColumnExists).Error
		if err != nil {
			log.Printf("Error checking if name column exists: %v", err)
			return err
		}

		if !nameColumnExists {
			// Clear all existing OTP sessions since they don't have the name field
			err = DB.Exec("DELETE FROM otp_sessions").Error
			if err != nil {
				log.Printf("Error clearing otp_sessions table: %v", err)
				return err
			}
			log.Println("✅ Cleared existing OTP sessions to allow adding name column")
		} else {
			log.Println("✅ name column already exists in otp_sessions table")
		}
	}

	return nil
}
