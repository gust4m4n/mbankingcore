package config

import (
	"log"
	"mbankingcore/models"

	"golang.org/x/crypto/bcrypt"
)

// SetupDatabase initializes database with migrations and seeds initial data
func SetupDatabase() error {
	log.Println("Setting up database...")

	// Run Auto Migration
	err := DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Config{},
		&models.Article{},
		&models.BankAccount{},
		&models.Transaction{},
		&models.Onboarding{},
		&models.AuditLog{},
		&models.Photo{},
		&models.DeviceSession{},
		&models.PendingTransaction{},
		&models.PendingUserStatusChange{},
		&models.ApprovalThreshold{},
	)
	if err != nil {
		log.Printf("Failed to auto-migrate models: %v", err)
		return err
	}
	log.Println("✅ Database tables created successfully")

	// Seed initial data
	if err := seedInitialData(); err != nil {
		log.Printf("Failed to seed initial data: %v", err)
		return err
	}

	log.Println("🚀 Database setup completed successfully!")
	return nil
}

// seedInitialData creates essential initial data for new project
func seedInitialData() error {
	log.Println("Seeding initial data...")

	// Seed initial admin users
	if err := seedInitialAdmins(); err != nil {
		return err
	}

	// Seed demo user for testing
	if err := seedDemoUser(); err != nil {
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

	// Seed initial approval thresholds
	if err := seedInitialApprovalThresholds(); err != nil {
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

	// Create admin user 1: admin@mbankingcore.com
	hashedPassword1, err := bcrypt.GenerateFromPassword([]byte("Admin123?"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return err
	}

	admin1 := models.Admin{
		Name:     "Admin User",
		Email:    "admin@mbankingcore.com",
		Password: string(hashedPassword1),
		Role:     "admin",
		Status:   1, // active
	}

	if err := DB.Create(&admin1).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return err
	}

	// Create super admin user 2: super@mbankingcore.com
	hashedPassword2, err := bcrypt.GenerateFromPassword([]byte("Super123?"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash super admin password: %v", err)
		return err
	}

	admin2 := models.Admin{
		Name:     "Super Admin",
		Email:    "super@mbankingcore.com",
		Password: string(hashedPassword2),
		Role:     "super",
		Status:   1, // active
	}

	if err := DB.Create(&admin2).Error; err != nil {
		log.Printf("Failed to create super admin: %v", err)
		return err
	}

	log.Println("✅ Created essential admin users:")
	log.Println("   - admin@mbankingcore.com (role: admin)")
	log.Println("   - super@mbankingcore.com (role: super)")
	return nil
}

// seedDemoUser creates a demo user for testing purposes
func seedDemoUser() error {
	log.Println("Seeding demo user...")

	// Check if demo user already exists
	var existingUser models.User
	err := DB.Where("phone = ?", "+621234567890").First(&existingUser).Error
	if err == nil {
		// Demo user exists, update PIN with hashed version if needed
		hashedPIN, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash demo user PIN: %v", err)
			return err
		}

		// Update the PIN to hashed version
		if err := DB.Model(&existingUser).Update("pin_atm", string(hashedPIN)).Error; err != nil {
			log.Printf("Failed to update demo user PIN: %v", err)
			return err
		}

		log.Println("✅ Demo user PIN updated with hashed version")
		return nil
	}

	// Hash the PIN ATM
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash demo user PIN: %v", err)
		return err
	}

	// Create demo user
	demoUser := models.User{
		Name:       "Demo User",
		Phone:      "+621234567890",
		MotherName: "Demo Mother",
		PinAtm:     string(hashedPIN),
		Balance:    1000000, // 1 juta rupiah initial balance
		Status:     1,       // active
	}

	if err := DB.Create(&demoUser).Error; err != nil {
		log.Printf("Failed to create demo user: %v", err)
		return err
	}

	// Create demo bank account for the user
	demoBankAccount := models.BankAccount{
		UserID:        demoUser.ID,
		AccountNumber: "1234567890",
		AccountName:   "Demo User",
		BankName:      "Demo Bank",
		BankCode:      "001",
		AccountType:   "saving",
		IsActive:      true,
		IsPrimary:     true,
	}

	if err := DB.Create(&demoBankAccount).Error; err != nil {
		log.Printf("Failed to create demo bank account: %v", err)
		return err
	}

	log.Println("✅ Created demo user for testing:")
	log.Println("   - Name: Demo User")
	log.Println("   - Phone: +621234567890")
	log.Println("   - Mother Name: Demo Mother")
	log.Println("   - PIN ATM: 123456")
	log.Println("   - Account Number: 1234567890")
	log.Println("   - Initial Balance: 1,000,000")
	return nil
}

// seedInitialConfigs creates default configuration values
func seedInitialConfigs() error {
	log.Println("Seeding initial configuration values...")

	// Check if configs already exist
	var count int64
	DB.Model(&models.Config{}).Count(&count)

	if count > 0 {
		log.Println("✅ Configuration values already exist")
		return nil
	}

	initialConfigs := []models.Config{
		{Key: "app_name", Value: "MBankingCore"},
		{Key: "app_version", Value: "0.9"},
		{Key: "terms_conditions", Value: getTermsConditionsContent()},
		{Key: "privacy_policy", Value: getPrivacyPolicyContent()},
		{Key: "admin_terms_conditions", Value: getAdminTermsConditionsContent()},
		{Key: "admin_privacy_policy", Value: getAdminPrivacyPolicyContent()},
		{Key: "contact_email", Value: "support@mbankingcore.com"},
		{Key: "contact_phone", Value: "+62-21-12345678"},
		{Key: "maintenance_mode", Value: "false"},
		{Key: "max_sessions_per_user", Value: "5"},
	}

	for _, config := range initialConfigs {
		if err := DB.Create(&config).Error; err != nil {
			log.Printf("Failed to create config %s: %v", config.Key, err)
			return err
		}
	}

	log.Printf("✅ Created %d initial configuration values", len(initialConfigs))
	return nil
}

// seedInitialOnboarding creates default onboarding slides
func seedInitialOnboarding() error {
	log.Println("Seeding initial onboarding content...")

	// Check if onboarding content already exists
	var count int64
	DB.Model(&models.Onboarding{}).Count(&count)

	if count > 0 {
		log.Println("✅ Onboarding content already exists")
		return nil
	}

	initialOnboarding := []models.Onboarding{
		{
			Title:       "Welcome to MBankingCore",
			Description: "Your secure and reliable mobile banking solution",
			Image:       "https://example.com/welcome.png",
			IsActive:    true,
		},
		{
			Title:       "Secure Banking",
			Description: "Bank safely with advanced security features",
			Image:       "https://example.com/security.png",
			IsActive:    true,
		},
		{
			Title:       "Easy Transactions",
			Description: "Send money and pay bills with just a few taps",
			Image:       "https://example.com/transactions.png",
			IsActive:    true,
		},
		{
			Title:       "24/7 Support",
			Description: "Get help whenever you need it",
			Image:       "https://example.com/support.png",
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

// seedInitialApprovalThresholds creates default approval thresholds for checker-maker system
func seedInitialApprovalThresholds() error {
	log.Println("Seeding initial approval thresholds...")

	// Check if approval thresholds already exist
	var count int64
	DB.Model(&models.ApprovalThreshold{}).Count(&count)

	if count > 0 {
		log.Println("✅ Approval thresholds already exist")
		return nil
	}

	initialThresholds := []models.ApprovalThreshold{
		{
			TransactionType:       "topup",
			AmountThreshold:       5000000, // 5 million IDR
			RequiresDualApproval:  true,
			DualApprovalThreshold: 50000000, // 50 million IDR
			AutoExpireHours:       24,
			IsActive:              true,
		},
		{
			TransactionType:       "withdraw",
			AmountThreshold:       2000000, // 2 million IDR
			RequiresDualApproval:  true,
			DualApprovalThreshold: 20000000, // 20 million IDR
			AutoExpireHours:       12,
			IsActive:              true,
		},
		{
			TransactionType:       "transfer",
			AmountThreshold:       10000000, // 10 million IDR
			RequiresDualApproval:  true,
			DualApprovalThreshold: 100000000, // 100 million IDR
			AutoExpireHours:       24,
			IsActive:              true,
		},
		{
			TransactionType:       "balance_adjustment",
			AmountThreshold:       1000000, // 1 million IDR
			RequiresDualApproval:  true,
			DualApprovalThreshold: 10000000, // 10 million IDR
			AutoExpireHours:       48,
			IsActive:              true,
		},
		{
			TransactionType:       "balance_set",
			AmountThreshold:       5000000, // 5 million IDR
			RequiresDualApproval:  true,
			DualApprovalThreshold: 50000000, // 50 million IDR
			AutoExpireHours:       48,
			IsActive:              true,
		},
	}

	for _, threshold := range initialThresholds {
		if err := DB.Create(&threshold).Error; err != nil {
			log.Printf("Failed to create approval threshold for %s: %v", threshold.TransactionType, err)
			return err
		}
	}

	log.Printf("✅ Created %d initial approval thresholds", len(initialThresholds))
	log.Println("   Approval thresholds configured:")
	log.Println("   - Topup: 5M IDR (dual: 50M IDR, expires: 24h)")
	log.Println("   - Withdraw: 2M IDR (dual: 20M IDR, expires: 12h)")
	log.Println("   - Transfer: 10M IDR (dual: 100M IDR, expires: 24h)")
	log.Println("   - Balance Adjustment: 1M IDR (dual: 10M IDR, expires: 48h)")
	log.Println("   - Balance Set: 5M IDR (dual: 50M IDR, expires: 48h)")
	return nil
}

// Simple content functions without emoji characters

func getTermsConditionsContent() string {
	return `<h1>SYARAT DAN KETENTUAN PENGGUNAAN MBANKINGCORE</h1>
<p><b>Efektif sejak:</b> 2 Agustus 2025</p>
<h2>1. PENERIMAAN SYARAT</h2>
<p>Dengan menggunakan aplikasi MBankingCore, Anda menyetujui syarat dan ketentuan ini.</p>
<h2>2. LAYANAN</h2>
<ul>
<li>Transfer dana</li>
<li>Pembayaran tagihan</li>
<li>Cek saldo dan mutasi</li>
</ul>
<h2>3. KEAMANAN</h2>
<ul>
<li>Jaga kerahasiaan PIN</li>
<li>Logout setelah penggunaan</li>
<li>Laporkan aktivitas mencurigakan</li>
</ul>`
}

func getPrivacyPolicyContent() string {
	return `<h1>KEBIJAKAN PRIVASI MBANKINGCORE</h1>
<p><b>Efektif sejak:</b> 2 Agustus 2025</p>
<h2>1. INFORMASI YANG DIKUMPULKAN</h2>
<ul>
<li>Informasi identitas</li>
<li>Data transaksi</li>
<li>Data perangkat</li>
</ul>
<h2>2. PENGGUNAAN INFORMASI</h2>
<ul>
<li>Penyediaan layanan</li>
<li>Keamanan akun</li>
<li>Komunikasi dengan pengguna</li>
</ul>
<h2>3. PERLINDUNGAN DATA</h2>
<ul>
<li>Enkripsi data</li>
<li>Akses terbatas</li>
<li>Monitoring keamanan</li>
</ul>`
}

func getAdminTermsConditionsContent() string {
	return `<h1>SYARAT DAN KETENTUAN ADMINISTRATOR MBANKINGCORE</h1>
<p><b>Efektif sejak:</b> 2 Agustus 2025</p>
<h2>1. KETENTUAN AKSES</h2>
<ul>
<li>Akses sesuai role dan wewenang</li>
<li>Kerahasiaan kredensial login</li>
<li>Compliance terhadap kebijakan</li>
</ul>
<h2>2. KEWENANGAN ADMINISTRATOR</h2>
<ul>
<li>Manajemen pengguna</li>
<li>Monitoring transaksi</li>
<li>Konfigurasi sistem</li>
</ul>
<h2>3. KEWAJIBAN</h2>
<ul>
<li>Menjaga keamanan sistem</li>
<li>Melaporkan insiden</li>
<li>Audit compliance</li>
</ul>`
}

func getAdminPrivacyPolicyContent() string {
	return `<h1>KEBIJAKAN PRIVASI ADMINISTRATOR MBANKINGCORE</h1>
<p><b>Efektif sejak:</b> 2 Agustus 2025</p>
<h2>1. RUANG LINGKUP</h2>
<p>Kebijakan ini mengatur perlindungan data administrator sistem MBankingCore.</p>
<h2>2. DATA YANG DIKUMPULKAN</h2>
<ul>
<li>Identitas administrator</li>
<li>Log aktivitas sistem</li>
<li>Data akses dan otorisasi</li>
</ul>
<h2>3. KEAMANAN DATA</h2>
<ul>
<li>Enkripsi end-to-end</li>
<li>Multi-factor authentication</li>
<li>Session monitoring</li>
</ul>
<h2>4. KEWAJIBAN ADMINISTRATOR</h2>
<ul>
<li>Menjaga kerahasiaan data</li>
<li>Melaporkan insiden keamanan</li>
<li>Mematuhi protokol keamanan</li>
</ul>`
}
