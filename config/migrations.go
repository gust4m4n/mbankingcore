package config

import (
	"log"
	"mbankingcore/models"
)

// RunMigrations handles all database migrations and initial setup
func RunMigrations() error {
	log.Println("Starting database migrations...")

	// Step 1: Auto-migrate all models
	err := DB.AutoMigrate(
		&models.User{},
		&models.DeviceSession{},
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

	log.Println("✅ Custom migrations completed")
	return nil
}

// migrateUserRoles ensures all users have proper role values
func migrateUserRoles() error {
	log.Println("Migrating user roles...")

	// Update existing users without roles to have default 'user' role
	result := DB.Model(&models.User{}).Where("role IS NULL OR role = ''").Update("role", models.ROLE_USER)
	if result.Error != nil {
		log.Printf("Error updating user roles: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected > 0 {
		log.Printf("✅ Updated %d users with default role", result.RowsAffected)
	} else {
		log.Println("✅ All users already have proper roles")
	}

	return nil
}

// seedInitialData creates essential initial data for the application
func seedInitialData() error {
	log.Println("Seeding initial data...")

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
