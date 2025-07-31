package config

import (
	"fmt"
	"log"
	"math/rand"
	"mbankingcore/models"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RunMigrations handles database setup for new project
func RunMigrations() error {
	log.Println("Setting up database for new project...")

	// Auto-migrate all models
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
		&models.Transaction{},
		&models.AuditLog{},
		&models.LoginAudit{},
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

	// Seed initial configuration values
	if err := seedInitialConfigs(); err != nil {
		return err
	}

	// Seed initial onboarding content
	if err := seedInitialOnboarding(); err != nil {
		return err
	}

	// Seed initial users for testing
	if err := seedInitialUsers(); err != nil {
		return err
	}

	// Seed initial transactions for testing
	if err := seedInitialTransactions(); err != nil {
		return err
	}

	// Seed 1000 dummy transactions for yearly data
	if err := seedDummyTransactionsYearly(); err != nil {
		return err
	}

	// Seed 1000 dummy users for yearly data
	if err := seedDummyUsersYearly(); err != nil {
		return err
	}

	// Seed 32 dummy admins for yearly data
	if err := seedDummyAdminsYearly(); err != nil {
		return err
	}

	// Seed 10000 dummy transactions for yearly data
	if err := seedDummyTransactionsYearlyLarge(); err != nil {
		return err
	}

	// Seed 20000 dummy transactions from 2020-2025
	if err := seedDummyTransactions2020To2025(); err != nil {
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

	// Create 16 demo admin users with common names
	demoAdmins := []struct {
		Name  string
		Email string
	}{
		{"John Smith", "john.smith@mbankingcore.com"},
		{"Sarah Johnson", "sarah.johnson@mbankingcore.com"},
		{"Michael Brown", "michael.brown@mbankingcore.com"},
		{"Emily Davis", "emily.davis@mbankingcore.com"},
		{"David Wilson", "david.wilson@mbankingcore.com"},
		{"Lisa Miller", "lisa.miller@mbankingcore.com"},
		{"Robert Garcia", "robert.garcia@mbankingcore.com"},
		{"Jessica Martinez", "jessica.martinez@mbankingcore.com"},
		{"William Anderson", "william.anderson@mbankingcore.com"},
		{"Ashley Taylor", "ashley.taylor@mbankingcore.com"},
		{"James Thomas", "james.thomas@mbankingcore.com"},
		{"Amanda Jackson", "amanda.jackson@mbankingcore.com"},
		{"Christopher White", "christopher.white@mbankingcore.com"},
		{"Jennifer Harris", "jennifer.harris@mbankingcore.com"},
		{"Matthew Clark", "matthew.clark@mbankingcore.com"},
		{"Nicole Lewis", "nicole.lewis@mbankingcore.com"},
	}

	// Hash password for demo admins (using default password: Admin123!)
	demoPassword := "Admin123!"
	hashedDemoPassword, err := bcrypt.GenerateFromPassword([]byte(demoPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash demo admin password: %v", err)
		return err
	}

	// Create demo admin users
	for _, demo := range demoAdmins {
		demoAdmin := models.Admin{
			Name:     demo.Name,
			Email:    demo.Email,
			Password: string(hashedDemoPassword),
			Role:     "admin",
			Status:   1, // active
		}

		if err := DB.Create(&demoAdmin).Error; err != nil {
			log.Printf("Failed to create demo admin %s: %v", demo.Name, err)
			return err
		}
	}

	log.Println("✅ Created admin users:")
	log.Println("   - admin@mbankingcore.com (role: admin)")
	log.Println("   - super@mbankingcore.com (role: super)")
	log.Printf("   - %d demo admin users (role: admin, password: %s)", len(demoAdmins), demoPassword)
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
		{Key: "app_version", Value: "1.0.0"},
		{Key: "terms_conditions", Value: getTermsConditionsContent()},
		{Key: "privacy_policy", Value: getPrivacyPolicyContent()},
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

// seedInitialUsers creates default users for testing
func seedInitialUsers() error {
	log.Println("Seeding initial users...")

	// Check if users already exist
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count > 0 {
		log.Println("✅ Users already exist")
		return nil
	}

	initialUsers := []models.User{
		// Original 3 users
		{
			Name:       "Demo User",
			Phone:      "081234567890",
			MotherName: "Maria Sari",
			PinAtm:     "123456", // Note: In real app, this should be hashed
			Balance:    1000000,  // 1 million IDR for demo
			Status:     1,        // active
		},
		{
			Name:       "Test User",
			Phone:      "082345678901",
			MotherName: "Siti Nurhaliza",
			PinAtm:     "654321", // Note: In real app, this should be hashed
			Balance:    500000,   // 500 thousand IDR for demo
			Status:     1,        // active
		},
		{
			Name:       "John Doe",
			Phone:      "083456789012",
			MotherName: "Jane Smith",
			PinAtm:     "111111", // Note: In real app, this should be hashed
			Balance:    2000000,  // 2 million IDR for demo
			Status:     1,        // active
		},
		// 64 additional demo users with common Indonesian and international names
		{Name: "Ahmad Wijaya", Phone: "081100000001", MotherName: "Siti Rahayu", PinAtm: "100001", Balance: 750000, Status: 1},
		{Name: "Sari Dewi", Phone: "081100000002", MotherName: "Indira Sari", PinAtm: "100002", Balance: 850000, Status: 1},
		{Name: "Budi Santoso", Phone: "081100000003", MotherName: "Kartini Budi", PinAtm: "100003", Balance: 920000, Status: 1},
		{Name: "Maya Putri", Phone: "081100000004", MotherName: "Dewi Maya", PinAtm: "100004", Balance: 680000, Status: 1},
		{Name: "Rizki Pratama", Phone: "081100000005", MotherName: "Ningsih Rizki", PinAtm: "100005", Balance: 1200000, Status: 1},
		{Name: "Dina Anggraini", Phone: "081100000006", MotherName: "Sari Dina", PinAtm: "100006", Balance: 450000, Status: 1},
		{Name: "Eko Setiawan", Phone: "081100000007", MotherName: "Wati Eko", PinAtm: "100007", Balance: 880000, Status: 1},
		{Name: "Lia Permatasari", Phone: "081100000008", MotherName: "Indah Lia", PinAtm: "100008", Balance: 950000, Status: 1},
		{Name: "Agus Hermawan", Phone: "081100000009", MotherName: "Ratna Agus", PinAtm: "100009", Balance: 720000, Status: 1},
		{Name: "Nina Salsabila", Phone: "081100000010", MotherName: "Dewi Nina", PinAtm: "100010", Balance: 1100000, Status: 1},
		{Name: "Yoga Prasetyo", Phone: "081100000011", MotherName: "Sari Yoga", PinAtm: "100011", Balance: 630000, Status: 1},
		{Name: "Rina Handayani", Phone: "081100000012", MotherName: "Indira Rina", PinAtm: "100012", Balance: 990000, Status: 1},
		{Name: "Dedi Kurniawan", Phone: "081100000013", MotherName: "Wati Dedi", PinAtm: "100013", Balance: 820000, Status: 1},
		{Name: "Sinta Maharani", Phone: "081100000014", MotherName: "Dewi Sinta", PinAtm: "100014", Balance: 1050000, Status: 1},
		{Name: "Andi Firmansyah", Phone: "081100000015", MotherName: "Sari Andi", PinAtm: "100015", Balance: 770000, Status: 1},
		{Name: "Putri Ramadhani", Phone: "081100000016", MotherName: "Indah Putri", PinAtm: "100016", Balance: 580000, Status: 1},
		{Name: "Bayu Nugroho", Phone: "081100000017", MotherName: "Ratna Bayu", PinAtm: "100017", Balance: 1350000, Status: 1},
		{Name: "Fitri Yuliana", Phone: "081100000018", MotherName: "Dewi Fitri", PinAtm: "100018", Balance: 460000, Status: 1},
		{Name: "Hendra Gunawan", Phone: "081100000019", MotherName: "Sari Hendra", PinAtm: "100019", Balance: 890000, Status: 1},
		{Name: "Nadia Safitri", Phone: "081100000020", MotherName: "Indira Nadia", PinAtm: "100020", Balance: 1180000, Status: 1},
		{Name: "Roni Setiadi", Phone: "081100000021", MotherName: "Wati Roni", PinAtm: "100021", Balance: 640000, Status: 1},
		{Name: "Yuni Astuti", Phone: "081100000022", MotherName: "Dewi Yuni", PinAtm: "100022", Balance: 750000, Status: 1},
		{Name: "Fandi Hidayat", Phone: "081100000023", MotherName: "Sari Fandi", PinAtm: "100023", Balance: 1020000, Status: 1},
		{Name: "Dewi Lestari", Phone: "081100000024", MotherName: "Indah Dewi", PinAtm: "100024", Balance: 830000, Status: 1},
		{Name: "Irwan Maulana", Phone: "081100000025", MotherName: "Ratna Irwan", PinAtm: "100025", Balance: 970000, Status: 1},
		{Name: "Tari Oktavia", Phone: "081100000026", MotherName: "Dewi Tari", PinAtm: "100026", Balance: 520000, Status: 1},
		{Name: "Gilang Ramadhan", Phone: "081100000027", MotherName: "Sari Gilang", PinAtm: "100027", Balance: 1280000, Status: 1},
		{Name: "Vera Susanti", Phone: "081100000028", MotherName: "Indira Vera", PinAtm: "100028", Balance: 680000, Status: 1},
		{Name: "Dani Pradipta", Phone: "081100000029", MotherName: "Wati Dani", PinAtm: "100029", Balance: 760000, Status: 1},
		{Name: "Mira Kusuma", Phone: "081100000030", MotherName: "Dewi Mira", PinAtm: "100030", Balance: 1150000, Status: 1},
		{Name: "Wahyu Nugraha", Phone: "081100000031", MotherName: "Sari Wahyu", PinAtm: "100031", Balance: 590000, Status: 1},
		{Name: "Lani Puspita", Phone: "081100000032", MotherName: "Indah Lani", PinAtm: "100032", Balance: 940000, Status: 1},
		{Name: "Hadi Wijono", Phone: "081100000033", MotherName: "Ratna Hadi", PinAtm: "100033", Balance: 870000, Status: 1},
		{Name: "Sari Melati", Phone: "081100000034", MotherName: "Dewi Sari", PinAtm: "100034", Balance: 1080000, Status: 1},
		{Name: "Tommy Hartono", Phone: "081100000035", MotherName: "Sari Tommy", PinAtm: "100035", Balance: 720000, Status: 1},
		{Name: "Ratih Permana", Phone: "081100000036", MotherName: "Indira Ratih", PinAtm: "100036", Balance: 650000, Status: 1},
		{Name: "Ardi Saputra", Phone: "081100000037", MotherName: "Wati Ardi", PinAtm: "100037", Balance: 1300000, Status: 1},
		{Name: "Lina Novita", Phone: "081100000038", MotherName: "Dewi Lina", PinAtm: "100038", Balance: 480000, Status: 1},
		{Name: "Ferry Kurniadi", Phone: "081100000039", MotherName: "Sari Ferry", PinAtm: "100039", Balance: 800000, Status: 1},
		{Name: "Diah Ayu", Phone: "081100000040", MotherName: "Indah Diah", PinAtm: "100040", Balance: 1120000, Status: 1},
		{Name: "Eko Wardana", Phone: "081100000041", MotherName: "Ratna Eko", PinAtm: "100041", Balance: 660000, Status: 1},
		{Name: "Mega Wulandari", Phone: "081100000042", MotherName: "Dewi Mega", PinAtm: "100042", Balance: 780000, Status: 1},
		{Name: "Ryan Adiputra", Phone: "081100000043", MotherName: "Sari Ryan", PinAtm: "100043", Balance: 1250000, Status: 1},
		{Name: "Ika Suryani", Phone: "081100000044", MotherName: "Indira Ika", PinAtm: "100044", Balance: 570000, Status: 1},
		{Name: "Donny Pranata", Phone: "081100000045", MotherName: "Wati Donny", PinAtm: "100045", Balance: 910000, Status: 1},
		{Name: "Winda Cahyani", Phone: "081100000046", MotherName: "Dewi Winda", PinAtm: "100046", Balance: 840000, Status: 1},
		{Name: "Fajar Mulyadi", Phone: "081100000047", MotherName: "Sari Fajar", PinAtm: "100047", Balance: 1190000, Status: 1},
		{Name: "Tina Marlina", Phone: "081100000048", MotherName: "Indah Tina", PinAtm: "100048", Balance: 620000, Status: 1},
		{Name: "Ardi Nurhadi", Phone: "081100000049", MotherName: "Ratna Ardi", PinAtm: "100049", Balance: 730000, Status: 1},
		{Name: "Siska Amelia", Phone: "081100000050", MotherName: "Dewi Siska", PinAtm: "100050", Balance: 1060000, Status: 1},
		{Name: "Benny Wijaya", Phone: "081100000051", MotherName: "Sari Benny", PinAtm: "100051", Balance: 790000, Status: 1},
		{Name: "Intan Permata", Phone: "081100000052", MotherName: "Indira Intan", PinAtm: "100052", Balance: 540000, Status: 1},
		{Name: "Joko Widodo", Phone: "081100000053", MotherName: "Wati Joko", PinAtm: "100053", Balance: 1380000, Status: 1},
		{Name: "Rini Susilo", Phone: "081100000054", MotherName: "Dewi Rini", PinAtm: "100054", Balance: 450000, Status: 1},
		{Name: "Andi Guntur", Phone: "081100000055", MotherName: "Sari Andi", PinAtm: "100055", Balance: 860000, Status: 1},
		{Name: "Lely Prastiwi", Phone: "081100000056", MotherName: "Indah Lely", PinAtm: "100056", Balance: 1140000, Status: 1},
		{Name: "Reza Pratama", Phone: "081100000057", MotherName: "Ratna Reza", PinAtm: "100057", Balance: 700000, Status: 1},
		{Name: "Sinta Dewi", Phone: "081100000058", MotherName: "Dewi Sinta", PinAtm: "100058", Balance: 820000, Status: 1},
		{Name: "Iwan Setiawan", Phone: "081100000059", MotherName: "Sari Iwan", PinAtm: "100059", Balance: 1200000, Status: 1},
		{Name: "Astrid Indira", Phone: "081100000060", MotherName: "Indira Astrid", PinAtm: "100060", Balance: 610000, Status: 1},
		{Name: "Candra Kusuma", Phone: "081100000061", MotherName: "Wati Candra", PinAtm: "100061", Balance: 980000, Status: 1},
		{Name: "Widya Sari", Phone: "081100000062", MotherName: "Dewi Widya", PinAtm: "100062", Balance: 750000, Status: 1},
		{Name: "Herman Sutanto", Phone: "081100000063", MotherName: "Sari Herman", PinAtm: "100063", Balance: 1320000, Status: 1},
		{Name: "Kartika Sari", Phone: "081100000064", MotherName: "Indah Kartika", PinAtm: "100064", Balance: 580000, Status: 1},
	}

	for _, user := range initialUsers {
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Name, err)
			return err
		}
	}

	log.Printf("✅ Created %d initial users for testing", len(initialUsers))
	return nil
}

// seedInitialTransactions creates sample transactions for testing
func seedInitialTransactions() error {
	log.Println("Seeding initial transactions...")

	// Check if transactions already exist
	var count int64
	DB.Model(&models.Transaction{}).Count(&count)

	if count > 0 {
		log.Println("✅ Transactions already exist")
		return nil
	}

	// Get all users to assign transactions
	var users []models.User
	if err := DB.Find(&users).Error; err != nil {
		log.Printf("Failed to get users for transaction seeding: %v", err)
		return err
	}

	if len(users) == 0 {
		log.Println("No users found for transaction seeding")
		return nil
	}

	// Create a map of user IDs for easier access
	userIDs := make([]uint, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	// Create sample transactions with realistic Indonesian mobile banking scenarios
	// Using actual user IDs from the database
	transactions := []models.Transaction{
		// First 3 users (Demo User, Test User, John Doe)
		{UserID: userIDs[0], Type: "topup", Amount: 500000, BalanceBefore: 1000000, BalanceAfter: 1500000, Status: "completed", Description: "Top up via ATM BCA"},
		{UserID: userIDs[0], Type: "transfer_out", Amount: 50000, BalanceBefore: 1500000, BalanceAfter: 1450000, Status: "completed", Description: "Transfer ke Sari - Bayar makan siang"},
		{UserID: userIDs[0], Type: "withdraw", Amount: 100000, BalanceBefore: 1450000, BalanceAfter: 1350000, Status: "completed", Description: "Tarik tunai untuk belanja"},

		{UserID: userIDs[1], Type: "transfer_in", Amount: 50000, BalanceBefore: 500000, BalanceAfter: 550000, Status: "completed", Description: "Transfer dari Demo User - Bayar makan siang"},
		{UserID: userIDs[1], Type: "topup", Amount: 200000, BalanceBefore: 550000, BalanceAfter: 750000, Status: "completed", Description: "Top up via m-banking"},
		{UserID: userIDs[1], Type: "transfer_out", Amount: 75000, BalanceBefore: 750000, BalanceAfter: 675000, Status: "completed", Description: "Transfer ke Ahmad - Bayar hutang"},

		{UserID: userIDs[2], Type: "topup", Amount: 1000000, BalanceBefore: 2000000, BalanceAfter: 3000000, Status: "completed", Description: "Top up gaji bulanan"},
		{UserID: userIDs[2], Type: "transfer_out", Amount: 500000, BalanceBefore: 3000000, BalanceAfter: 2500000, Status: "completed", Description: "Transfer ke keluarga"},
		{UserID: userIDs[2], Type: "withdraw", Amount: 200000, BalanceBefore: 2500000, BalanceAfter: 2300000, Status: "completed", Description: "Tarik tunai untuk keperluan harian"},

		// Continue with the rest of the users
		{UserID: userIDs[3], Type: "transfer_in", Amount: 75000, BalanceBefore: 750000, BalanceAfter: 825000, Status: "completed", Description: "Transfer dari Test User - Bayar hutang"},
		{UserID: userIDs[3], Type: "topup", Amount: 300000, BalanceBefore: 825000, BalanceAfter: 1125000, Status: "completed", Description: "Top up via internet banking"},
		{UserID: userIDs[3], Type: "transfer_out", Amount: 100000, BalanceBefore: 1125000, BalanceAfter: 1025000, Status: "completed", Description: "Bayar tagihan listrik"},

		{UserID: userIDs[4], Type: "topup", Amount: 250000, BalanceBefore: 850000, BalanceAfter: 1100000, Status: "completed", Description: "Top up via ATM Mandiri"},
		{UserID: userIDs[4], Type: "transfer_out", Amount: 80000, BalanceBefore: 1100000, BalanceAfter: 1020000, Status: "completed", Description: "Transfer ke Maya - Split bill restoran"},
		{UserID: userIDs[4], Type: "withdraw", Amount: 50000, BalanceBefore: 1020000, BalanceAfter: 970000, Status: "completed", Description: "Tarik tunai untuk ongkos"},

		{UserID: userIDs[5], Type: "topup", Amount: 400000, BalanceBefore: 920000, BalanceAfter: 1320000, Status: "completed", Description: "Top up bonus kerja"},
		{UserID: userIDs[5], Type: "transfer_out", Amount: 150000, BalanceBefore: 1320000, BalanceAfter: 1170000, Status: "completed", Description: "Bayar asuransi bulanan"},
		{UserID: userIDs[5], Type: "transfer_out", Amount: 200000, BalanceBefore: 1170000, BalanceAfter: 970000, Status: "completed", Description: "Transfer ke orang tua"},

		{UserID: userIDs[6], Type: "transfer_in", Amount: 80000, BalanceBefore: 680000, BalanceAfter: 760000, Status: "completed", Description: "Transfer dari Sari - Split bill restoran"},
		{UserID: userIDs[6], Type: "topup", Amount: 150000, BalanceBefore: 760000, BalanceAfter: 910000, Status: "completed", Description: "Top up untuk belanja online"},
		{UserID: userIDs[6], Type: "withdraw", Amount: 100000, BalanceBefore: 910000, BalanceAfter: 810000, Status: "completed", Description: "Tarik tunai untuk shopping"},

		{UserID: userIDs[7], Type: "topup", Amount: 600000, BalanceBefore: 1200000, BalanceAfter: 1800000, Status: "completed", Description: "Top up dari rekening utama"},
		{UserID: userIDs[7], Type: "transfer_out", Amount: 300000, BalanceBefore: 1800000, BalanceAfter: 1500000, Status: "completed", Description: "Bayar sewa kos bulanan"},
		{UserID: userIDs[7], Type: "transfer_out", Amount: 100000, BalanceBefore: 1500000, BalanceAfter: 1400000, Status: "completed", Description: "Transfer ke Dina - Bayar patungan gift"},

		{UserID: userIDs[8], Type: "transfer_in", Amount: 100000, BalanceBefore: 450000, BalanceAfter: 550000, Status: "completed", Description: "Transfer dari Rizki - Bayar patungan gift"},
		{UserID: userIDs[8], Type: "topup", Amount: 200000, BalanceBefore: 550000, BalanceAfter: 750000, Status: "completed", Description: "Top up untuk traveling"},
		{UserID: userIDs[8], Type: "withdraw", Amount: 75000, BalanceBefore: 750000, BalanceAfter: 675000, Status: "completed", Description: "Tarik tunai untuk jajan"},

		// Continue with more users if available
	}

	// Add more transactions for remaining users dynamically
	for i := 9; i < len(userIDs) && i < 30; i++ {
		baseTransactions := []models.Transaction{
			{UserID: userIDs[i], Type: "topup", Amount: int64(200000 + (i * 10000)), BalanceBefore: int64(800000 + (i * 50000)), BalanceAfter: int64(1000000 + (i * 60000)), Status: "completed", Description: "Top up via mobile banking"},
			{UserID: userIDs[i], Type: "withdraw", Amount: int64(50000 + (i * 5000)), BalanceBefore: int64(1000000 + (i * 60000)), BalanceAfter: int64(950000 + (i * 55000)), Status: "completed", Description: "Tarik tunai untuk keperluan"},
			{UserID: userIDs[i], Type: "transfer_out", Amount: int64(75000 + (i * 3000)), BalanceBefore: int64(950000 + (i * 55000)), BalanceAfter: int64(875000 + (i * 52000)), Status: "completed", Description: "Transfer untuk pembayaran"},
		}
		transactions = append(transactions, baseTransactions...)
	}

	// Add some failed transactions for realistic scenarios
	if len(userIDs) > 20 {
		failedTransactions := []models.Transaction{
			{UserID: userIDs[20], Type: "transfer_out", Amount: 2000000, BalanceBefore: 520000, BalanceAfter: 520000, Status: "failed", Description: "Transfer gagal - Saldo tidak mencukupi"},
			{UserID: userIDs[21], Type: "withdraw", Amount: 5000000, BalanceBefore: 680000, BalanceAfter: 680000, Status: "failed", Description: "Penarikan gagal - Melebihi limit harian"},
		}
		transactions = append(transactions, failedTransactions...)
	}

	// Create all transactions
	for _, transaction := range transactions {
		if err := DB.Create(&transaction).Error; err != nil {
			log.Printf("Failed to create transaction for user %d: %v", transaction.UserID, err)
			return err
		}
	}

	log.Printf("✅ Created %d initial transactions for testing", len(transactions))
	return nil
}

// seedDummyTransactionsYearly creates 1000 dummy transactions over 1 year period
func seedDummyTransactionsYearly() error {
	log.Println("Seeding 1000 dummy transactions for yearly data...")

	// Initialize random seed for reproducible results
	rand.Seed(time.Now().UnixNano())

	// Check if we already have enough transactions (threshold: 1000+)
	var count int64
	DB.Model(&models.Transaction{}).Count(&count)

	if count >= 1000 {
		log.Printf("✅ Already have %d transactions, skipping yearly dummy data", count)
		return nil
	}

	// Clear existing transactions to regenerate with yearly data
	if count > 0 {
		log.Println("Clearing existing transactions to regenerate with yearly data...")
		if err := DB.Where("1 = 1").Delete(&models.Transaction{}).Error; err != nil {
			log.Printf("Failed to clear existing transactions: %v", err)
			return err
		}
		log.Println("✅ Existing transactions cleared")
	}

	// Get all users to assign transactions
	var users []models.User
	if err := DB.Find(&users).Error; err != nil {
		log.Printf("Failed to get users for yearly transaction seeding: %v", err)
		return err
	}

	if len(users) == 0 {
		log.Println("No users found for yearly transaction seeding")
		return nil
	}

	// Create a map of user IDs for easier access
	userIDs := make([]uint, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	// Base time: 1 year ago from now
	now := time.Now()
	baseTime := now.AddDate(-1, 0, 0)

	// Transaction types and their relative frequencies
	transactionTypes := []string{"topup", "withdraw", "transfer_out", "transfer_in"}
	typeWeights := []int{25, 25, 30, 20} // Percentages

	// Indonesian banking descriptions for realism
	descriptions := map[string][]string{
		"topup": {
			"Top up via ATM BCA", "Top up via m-banking", "Top up gaji bulanan",
			"Top up via internet banking", "Top up via ATM Mandiri", "Top up bonus kerja",
			"Top up dari rekening utama", "Top up untuk traveling", "Top up via mobile banking",
			"Top up untuk belanja online", "Setor tunai via teller", "Transfer masuk dari rekening lain",
		},
		"withdraw": {
			"Tarik tunai untuk belanja", "Tarik tunai untuk keperluan harian", "Tarik tunai untuk ongkos",
			"Tarik tunai untuk shopping", "Tarik tunai untuk jajan", "Tarik tunai untuk keperluan",
			"Penarikan tunai di ATM", "Tarik tunai untuk bensin", "Penarikan untuk bayar warung",
			"Tarik tunai darurat", "Penarikan untuk belanja bulanan", "Tarik tunai weekend",
		},
		"transfer_out": {
			"Transfer ke keluarga", "Bayar tagihan listrik", "Transfer ke Maya - Split bill restoran",
			"Bayar asuransi bulanan", "Transfer ke orang tua", "Bayar sewa kos bulanan",
			"Transfer ke Dina - Bayar patungan gift", "Transfer untuk pembayaran", "Bayar tagihan internet",
			"Transfer ke teman - Bayar hutang", "Bayar SPP kuliah", "Transfer donasi",
			"Bayar cicilan motor", "Transfer ke adik - Uang saku", "Bayar tagihan kartu kredit",
		},
		"transfer_in": {
			"Transfer dari Demo User - Bayar makan siang", "Transfer dari Test User - Bayar hutang",
			"Transfer dari Sari - Split bill restoran", "Transfer dari Rizki - Bayar patungan gift",
			"Terima gaji bulanan", "Transfer dari orang tua", "Terima bonus penjualan",
			"Transfer dari klien", "Terima refund belanja", "Transfer dari teman",
			"Terima cashback", "Transfer komisi", "Terima hadiah ulang tahun",
		},
	}

	var transactions []models.Transaction
	successfulCount := 0

	// Generate 1000 transactions
	for i := 0; i < 1000; i++ {
		// Random user
		userID := userIDs[i%len(userIDs)]

		// Random time within the year
		dayOffset := rand.Intn(365)
		hourOffset := rand.Intn(24)
		minuteOffset := rand.Intn(60)
		transactionTime := baseTime.AddDate(0, 0, dayOffset).
			Add(time.Duration(hourOffset) * time.Hour).
			Add(time.Duration(minuteOffset) * time.Minute)

		// Determine transaction type based on weights
		randomNum := rand.Intn(100)
		var transactionType string
		cumulative := 0
		for j, weight := range typeWeights {
			cumulative += weight
			if randomNum < cumulative {
				transactionType = transactionTypes[j]
				break
			}
		}

		// Random amount based on transaction type
		var amount int64
		switch transactionType {
		case "topup":
			amount = int64(rand.Intn(2000000) + 100000) // 100k - 2.1M
		case "withdraw":
			amount = int64(rand.Intn(1000000) + 50000) // 50k - 1.05M
		case "transfer_out", "transfer_in":
			amount = int64(rand.Intn(1500000) + 25000) // 25k - 1.525M
		}

		// Random balance before (simulate realistic balances)
		balanceBefore := int64(rand.Intn(5000000) + 100000) // 100k - 5.1M

		// Calculate balance after based on transaction type
		var balanceAfter int64
		switch transactionType {
		case "topup", "transfer_in":
			balanceAfter = balanceBefore + amount
		case "withdraw", "transfer_out":
			balanceAfter = balanceBefore - amount
			// Ensure balance doesn't go negative for successful transactions
			if balanceAfter < 0 {
				balanceAfter = balanceBefore // Keep original balance for failed transaction
			}
		}

		// Determine transaction status (95% success rate)
		status := "completed"
		if rand.Intn(100) < 5 { // 5% failure rate
			status = "failed"
			balanceAfter = balanceBefore // Keep original balance for failed transactions
		}

		// Random description
		descOptions := descriptions[transactionType]
		description := descOptions[rand.Intn(len(descOptions))]
		if status == "failed" {
			switch transactionType {
			case "withdraw", "transfer_out":
				description += " - Saldo tidak mencukupi"
			default:
				description += " - Transaksi gagal"
			}
		}

		transaction := models.Transaction{
			UserID:        userID,
			Type:          transactionType,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Status:        status,
			Description:   description,
			CreatedAt:     transactionTime,
			UpdatedAt:     transactionTime,
		}

		transactions = append(transactions, transaction)

		if status == "completed" {
			successfulCount++
		}
	}

	// Create all transactions in batches for better performance
	batchSize := 100
	for i := 0; i < len(transactions); i += batchSize {
		end := i + batchSize
		if end > len(transactions) {
			end = len(transactions)
		}

		batch := transactions[i:end]
		if err := DB.Create(&batch).Error; err != nil {
			log.Printf("Failed to create transaction batch %d-%d: %v", i, end-1, err)
			return err
		}

		log.Printf("Created transaction batch %d-%d", i+1, end)
	}

	log.Printf("✅ Created %d dummy transactions (%d successful, %d failed) over 1 year period",
		len(transactions), successfulCount, len(transactions)-successfulCount)
	return nil
}

// seedDummyUsersYearly creates 5000 dummy users distributed from 2020-2025
func seedDummyUsersYearly() error {
	log.Println("Seeding 5000 dummy users with registration dates from 2020-2025...")

	// Check if dummy users already exist
	var count int64
	DB.Model(&models.User{}).Where("phone LIKE ?", "08%").Count(&count)
	if count >= 5000 {
		log.Printf("✅ %d dummy users already exist", count)
		return nil
	}

	// Indonesian names for variety
	firstNames := []string{
		"Andi", "Budi", "Citra", "Dewi", "Eko", "Fitri", "Gita", "Hadi",
		"Indra", "Joko", "Kiki", "Luna", "Maya", "Nita", "Omar", "Putri",
		"Qori", "Rina", "Sari", "Toni", "Udin", "Vina", "Wati", "Yani",
		"Zaki", "Agus", "Bayu", "Cinta", "Doni", "Ella", "Farid", "Gina",
		"Hendra", "Intan", "Johan", "Kartika", "Laras", "Mira", "Nando", "Ovi",
		"Pras", "Quincy", "Rudi", "Sinta", "Tari", "Ucok", "Vera", "Wina",
		"Yoga", "Zahra", "Arief", "Bella", "Cahyo", "Diana", "Edy", "Fara",
		"Galih", "Hana", "Ivan", "Jelita", "Koko", "Lina", "Mega", "Noval",
	}

	lastNames := []string{
		"Pratama", "Sari", "Utomo", "Wijaya", "Santoso", "Kurniawan", "Handoko", "Susanto",
		"Wibowo", "Setiawan", "Rahayu", "Lestari", "Mahendra", "Permana", "Saputra", "Kartini",
		"Hidayat", "Fitriani", "Nugroho", "Anggraini", "Prasetyo", "Safitri", "Gunawan", "Melati",
		"Adriansyah", "Puspita", "Ramadhan", "Cahyani", "Firmansyah", "Dewanti", "Syahputra", "Maharani",
		"Maulana", "Pertiwi", "Wardana", "Salsabila", "Putranto", "Kusuma", "Saputri", "Hakim",
		"Aditya", "Novita", "Rachman", "Wahyuni", "Subagyo", "Kirana", "Budiman", "Amelia",
		"Iskandar", "Mentari", "Rahman", "Purnama", "Surya", "Mawar", "Firdaus", "Sekarsari",
		"Ardiansyah", "Ananda", "Nurhayati", "Sutrisno", "Damayanti", "Purwanto", "Safira", "Hartono",
	}

	motherNames := []string{
		"Siti", "Sri", "Umi", "Ibu", "Nyai", "Ratu", "Dewi", "Putri",
		"Ratna", "Indira", "Kartika", "Melati", "Mawar", "Cempaka", "Anggrek", "Dahlia",
		"Kenanga", "Seruni", "Tulip", "Sakura", "Lily", "Rose", "Jasmine", "Orchid",
		"Bunga", "Cantika", "Anggun", "Jelita", "Ayu", "Endah", "Permata", "Intan",
	}

	// Time range: 2020-2025
	startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	timeRange := endTime.Unix() - startTime.Unix()

	// Track used phone numbers to ensure uniqueness
	usedPhones := make(map[string]bool)

	var users []models.User
	batchSize := 100 // Insert in batches for better performance

	for i := 0; i < 5000; i++ {
		// Generate unique phone number
		var phone string
		for {
			// Indonesian phone format: 08xxxxxxxxx (10-12 digits)
			phone = fmt.Sprintf("08%d", 1000000000+rand.Intn(9000000000))
			if !usedPhones[phone] {
				usedPhones[phone] = true
				break
			}
		}

		// Random registration time between 2020-2025
		randomTime := startTime.Add(time.Duration(rand.Int63n(timeRange)) * time.Second)

		// Generate random balance (0 to 50,000,000 rupiah)
		balance := int64(rand.Intn(50000000))

		// Random status (mostly active)
		status := models.USER_STATUS_ACTIVE
		if rand.Float32() < 0.1 { // 10% chance of inactive
			status = models.USER_STATUS_INACTIVE
		} else if rand.Float32() < 0.02 { // 2% chance of blocked
			status = models.USER_STATUS_BLOCKED
		}

		// Create user
		user := models.User{
			Name:       firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))],
			Phone:      phone,
			MotherName: motherNames[rand.Intn(len(motherNames))] + " " + lastNames[rand.Intn(len(lastNames))],
			PinAtm:     "123456", // Default PIN for dummy users
			Balance:    balance,
			Status:     status,
			Avatar:     "", // Empty avatar for dummy users
			CreatedAt:  randomTime,
			UpdatedAt:  randomTime,
		}

		users = append(users, user)

		// Insert in batches
		if len(users) >= batchSize || i == 5000-1 {
			if err := DB.CreateInBatches(&users, batchSize).Error; err != nil {
				log.Printf("Error inserting batch: %v", err)
				continue
			}

			log.Printf("Inserted batch: %d/5000 users", i+1)
			users = []models.User{} // Reset batch
		}
	}

	log.Printf("✅ Created 5000 dummy users with registration dates from 2020-2025")

	// Show summary statistics
	showUserStatistics()
	return nil
}

// showUserStatistics displays user statistics after seeding
func showUserStatistics() {
	log.Println("\n=== User Statistics ===")

	var totalUsers int64
	DB.Model(&models.User{}).Count(&totalUsers)
	log.Printf("Total users in database: %d", totalUsers)

	// Count by year
	for year := 2020; year <= 2025; year++ {
		var yearCount int64
		startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

		DB.Model(&models.User{}).Where("created_at BETWEEN ? AND ?", startOfYear, endOfYear).Count(&yearCount)
		log.Printf("Users registered in %d: %d", year, yearCount)
	}

	// Count by status
	var activeUsers, inactiveUsers, blockedUsers int64
	DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_ACTIVE).Count(&activeUsers)
	DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_INACTIVE).Count(&inactiveUsers)
	DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_BLOCKED).Count(&blockedUsers)

	log.Printf("Active users: %d", activeUsers)
	log.Printf("Inactive users: %d", inactiveUsers)
	log.Printf("Blocked users: %d", blockedUsers)
}

// seedDummyAdminsYearly creates 32 dummy admins distributed over 1 year
func seedDummyAdminsYearly() error {
	log.Println("Seeding 32 dummy admins for yearly data...")

	// Check if dummy admins already exist (excluding initial admins)
	var count int64
	DB.Model(&models.Admin{}).Where("email LIKE ?", "admin.%@mbankingcore.com").Count(&count)
	if count >= 32 {
		log.Println("✅ 32+ dummy admins already exist")
		return nil
	}

	// Clear existing dummy admins if any exist but less than 32
	if count > 0 {
		DB.Where("email LIKE ?", "admin.%@mbankingcore.com").Delete(&models.Admin{})
		log.Println("Cleared existing dummy admins")
	}

	// Generate start and end dates for 1 year period
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0) // 1 year ago

	// Admin names and roles
	adminNames := []string{
		"Ahmad Rizki", "Bella Sari", "Candra Wijaya", "Dina Kusuma", "Erik Pratama",
		"Farah Putri", "Gunawan Santoso", "Hana Maharani", "Indra Firmansyah", "Jihan Lestari",
		"Kiki Setiawan", "Linda Handayani", "Maya Nugroho", "Novi Rahayu", "Omar Kurniawan",
		"Putri Wulandari", "Qori Purwanto", "Rina Safitri", "Sari Irawan", "Tono Anggraini",
		"Umi Susanto", "Vina Permata", "Wati Hakim", "Xena Melati", "Yuni Adiputra",
		"Zahra Kartini", "Andi Baskara", "Citra Cahyani", "Dewi Dharma", "Eko Fadilla",
		"Gita Hapsari", "Heri Mahendra",
	}

	roles := []string{"admin", "admin", "admin", "super"} // Mostly admin, some super

	var admins []models.Admin
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 32; i++ {
		// Generate random creation time within the year
		randomDuration := time.Duration(rand.Int63n(int64(endDate.Sub(startDate))))
		createdAt := startDate.Add(randomDuration)

		// Generate email
		email := fmt.Sprintf("admin.%s@mbankingcore.com",
			strings.ToLower(strings.ReplaceAll(adminNames[i], " ", ".")))

		// Generate password hash
		password := fmt.Sprintf("Admin%d!", 1000+i)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for admin %s: %v", adminNames[i], err)
			return err
		}

		// Generate role (75% admin, 25% super)
		role := roles[rand.Intn(len(roles))]

		// Generate status (95% active, 5% inactive)
		status := 1
		if rand.Float64() < 0.05 {
			status = 0
		}

		admin := models.Admin{
			Name:     adminNames[i],
			Email:    email,
			Password: string(hashedPassword),
			Role:     role,
			Status:   status,
		}
		admin.CreatedAt = createdAt
		admin.UpdatedAt = createdAt

		admins = append(admins, admin)
	}

	// Create all admins in batches for better performance
	batchSize := 10
	for i := 0; i < len(admins); i += batchSize {
		end := i + batchSize
		if end > len(admins) {
			end = len(admins)
		}

		batch := admins[i:end]
		if err := DB.Create(&batch).Error; err != nil {
			log.Printf("Failed to create admin batch %d-%d: %v", i, end-1, err)
			return err
		}

		log.Printf("Created admin batch %d-%d", i+1, end)
	}

	log.Printf("✅ Created 32 dummy admins over 1 year period")
	return nil
}

// seedDummyTransactionsYearlyLarge creates 50000 dummy transactions distributed over 1 year
func seedDummyTransactionsYearlyLarge() error {
	log.Println("Seeding 50000 dummy transactions for yearly data...")

	// Check if large dummy transactions already exist
	var count int64
	DB.Model(&models.Transaction{}).Count(&count)
	if count >= 50000 {
		log.Println("✅ 50000+ transactions already exist")
		return nil
	}

	// Get all users for transaction assignment
	var users []models.User
	if err := DB.Find(&users).Error; err != nil {
		log.Printf("Failed to get users: %v", err)
		return err
	}

	if len(users) == 0 {
		log.Println("No users found, creating basic user first")
		return fmt.Errorf("no users available for transactions")
	}

	// Generate start and end dates for 1 year period
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0) // 1 year ago

	// Transaction types with weights
	transactionTypes := []string{"topup", "withdraw", "transfer_out", "transfer_in"}
	typeWeights := []float64{0.25, 0.25, 0.30, 0.20} // 25% topup, 25% withdraw, 30% transfer_out, 20% transfer_in

	// Indonesian transaction descriptions
	descriptions := map[string][]string{
		"topup": {
			"Top up via ATM BCA", "Top up via ATM Mandiri", "Top up via ATM BNI", "Top up via ATM BRI",
			"Top up via mobile banking", "Top up via internet banking", "Top up via m-banking",
			"Setor tunai via teller", "Transfer masuk dari rekening lain", "Top up untuk belanja online",
			"Top up untuk traveling", "Top up gaji bulanan", "Top up bonus kerja", "Top up dari rekening utama",
		},
		"withdraw": {
			"Tarik tunai untuk belanja", "Tarik tunai untuk bensin", "Tarik tunai untuk jajan",
			"Tarik tunai untuk ongkos", "Tarik tunai weekend", "Tarik tunai darurat",
			"Penarikan tunai di ATM", "Tarik tunai untuk keperluan", "Tarik tunai untuk shopping",
			"Penarikan untuk belanja bulanan", "Penarikan untuk bayar warung", "Tarik tunai untuk keperluan harian",
		},
		"transfer_out": {
			"Transfer ke orang tua", "Transfer donasi", "Bayar tagihan listrik", "Bayar tagihan internet",
			"Bayar SPP kuliah", "Bayar sewa kos bulanan", "Bayar cicilan motor", "Bayar tagihan kartu kredit",
			"Bayar asuransi bulanan", "Transfer untuk pembayaran", "Transfer ke teman - Bayar hutang",
			"Transfer ke Maya - Split bill restoran", "Transfer ke Dina - Bayar patungan gift",
		},
		"transfer_in": {
			"Terima gaji bulanan", "Terima bonus penjualan", "Terima cashback", "Transfer dari teman",
			"Transfer komisi", "Terima hadiah ulang tahun", "Terima refund belanja", "Transfer dari klien",
			"Transfer dari Test User - Bayar makan siang", "Transfer dari Sari - Split bill restoran",
			"Transfer dari Demo User - Bayar makan siang", "Transfer dari Rizki - Bayar patungan gift",
		},
	}

	var transactions []models.Transaction
	var successfulCount int
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 50000; i++ {
		// Generate random creation time within the year
		randomDuration := time.Duration(rand.Int63n(int64(endDate.Sub(startDate))))
		createdAt := startDate.Add(randomDuration)

		// Select random user
		user := users[rand.Intn(len(users))]

		// Select transaction type based on weights
		random := rand.Float64()
		var transactionType string
		cumulative := 0.0
		for j, weight := range typeWeights {
			cumulative += weight
			if random <= cumulative {
				transactionType = transactionTypes[j]
				break
			}
		}

		// Generate amount (25,000 to 2,500,000)
		amount := int64(25000 + rand.Float64()*2475000)

		// Generate description
		typeDescriptions := descriptions[transactionType]
		description := typeDescriptions[rand.Intn(len(typeDescriptions))]

		// Generate realistic balance before/after
		balanceBefore := int64(1000000 + rand.Float64()*4000000) // 1M to 5M
		var balanceAfter int64
		status := "completed"

		switch transactionType {
		case "topup", "transfer_in":
			balanceAfter = balanceBefore + amount
		case "withdraw", "transfer_out":
			if balanceBefore >= amount {
				balanceAfter = balanceBefore - amount
			} else {
				balanceAfter = balanceBefore
				status = "failed"
				description += " - Saldo tidak mencukupi"
			}
		}

		// 10% chance for random failure
		if rand.Float64() < 0.1 {
			status = "failed"
			balanceAfter = balanceBefore
			description += " - Transaksi gagal"
		}

		if status == "completed" {
			successfulCount++
		}

		transaction := models.Transaction{
			UserID:        user.ID,
			Type:          transactionType,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Description:   description,
			Status:        status,
		}
		transaction.CreatedAt = createdAt
		transaction.UpdatedAt = createdAt

		transactions = append(transactions, transaction)
	}

	// Create all transactions in batches for better performance
	batchSize := 100
	for i := 0; i < len(transactions); i += batchSize {
		end := i + batchSize
		if end > len(transactions) {
			end = len(transactions)
		}

		batch := transactions[i:end]
		if err := DB.Create(&batch).Error; err != nil {
			log.Printf("Failed to create transaction batch %d-%d: %v", i, end-1, err)
			return err
		}

		log.Printf("Created transaction batch %d-%d", i+1, end)
	}

	log.Printf("✅ Created 50000 dummy transactions (%d successful, %d failed) over 1 year period",
		successfulCount, 50000-successfulCount)
	return nil
}

// getTermsConditionsContent returns comprehensive Terms & Conditions content for Indonesian banking app
func getTermsConditionsContent() string {
	return `<h1>SYARAT DAN KETENTUAN PENGGUNAAN MBANKINGCORE</h1>

<p><b>Terakhir diperbarui: 31 Juli 2025</b></p>

<h2>1. PENERIMAAN SYARAT DAN KETENTUAN</h2>
<p>Dengan mengunduh, menginstal, mengakses atau menggunakan aplikasi MBankingCore ("<b>Aplikasi</b>"), Anda menyetujui untuk terikat oleh Syarat dan Ketentuan ini ("<b>S&K</b>"). Jika Anda tidak menyetujui S&K ini, harap tidak menggunakan Aplikasi.</p>

<h2>2. DEFINISI</h2>
<ul>
<li><b>"Aplikasi"</b> adalah aplikasi mobile banking MBankingCore</li>
<li><b>"Layanan"</b> adalah semua fitur dan layanan yang disediakan melalui Aplikasi</li>
<li><b>"Pengguna"</b> adalah individu yang telah terdaftar dan menggunakan Layanan</li>
<li><b>"Kami"</b> adalah penyedia Aplikasi MBankingCore</li>
<li><b>"Rekening"</b> adalah rekening bank yang terhubung dengan Aplikasi</li>
</ul>

<h2>3. PERSYARATAN PENGGUNAAN</h2>
<h3>3.1 Kelayakan</h3>
<ul>
<li>Anda harus berusia minimal 17 tahun</li>
<li>Memiliki rekening bank yang valid di Indonesia</li>
<li>Memiliki nomor telepon aktif yang terdaftar di bank</li>
<li>Menyediakan informasi yang akurat dan lengkap</li>
</ul>

<h3>3.2 Registrasi Akun</h3>
<ul>
<li>Satu nomor telepon hanya dapat digunakan untuk satu akun</li>
<li>Anda bertanggung jawab menjaga kerahasiaan PIN dan password</li>
<li>Segera laporkan jika terjadi penyalahgunaan akun</li>
</ul>

<h2>4. LAYANAN YANG TERSEDIA</h2>
<h3>4.1 Informasi Rekening</h3>
<ul>
<li>Cek saldo</li>
<li>Riwayat transaksi</li>
<li>Informasi rekening</li>
</ul>

<h3>4.2 Transfer Dana</h3>
<ul>
<li>Transfer antar bank</li>
<li>Transfer antar pengguna MBankingCore</li>
<li>Transfer terjadwal</li>
</ul>

<h3>4.3 Pembayaran</h3>
<ul>
<li>Pembayaran tagihan (listrik, air, telepon)</li>
<li>Pembayaran merchant</li>
<li>Top-up e-wallet</li>
</ul>

<h2>5. KEAMANAN DAN PERLINDUNGAN</h2>
<h3>5.1 Kewajiban Pengguna</h3>
<ul>
<li>Menjaga kerahasiaan PIN/password</li>
<li>Menggunakan koneksi internet yang aman</li>
<li>Melakukan logout setelah selesai menggunakan</li>
<li>Tidak berbagi informasi akun dengan pihak lain</li>
</ul>

<h3>5.2 Sistem Keamanan</h3>
<ul>
<li>Enkripsi data end-to-end</li>
<li>Otentikasi dua faktor (2FA)</li>
<li>Monitoring transaksi real-time</li>
<li>Notifikasi setiap transaksi</li>
</ul>

<h2>6. BATASAN DAN LARANGAN</h2>
<h3>6.1 Larangan Penggunaan</h3>
<ul>
<li>Menggunakan Aplikasi untuk kegiatan ilegal</li>
<li>Melakukan transaksi fiktif atau penipuan</li>
<li>Mengganggu sistem atau server Aplikasi</li>
<li>Menyalahgunakan fitur atau layanan</li>
</ul>

<h3>6.2 Batas Transaksi</h3>
<ul>
<li>Transfer harian: <b>Rp 25.000.000</b></li>
<li>Transfer per transaksi: <b>Rp 5.000.000</b></li>
<li>Pembayaran harian: <b>Rp 10.000.000</b></li>
<li>Batas dapat disesuaikan sesuai profil risiko</li>
</ul>

<h2>7. BIAYA DAN TARIF</h2>
<h3>7.1 Biaya Layanan</h3>
<ul>
<li>Transfer antar bank: <b>Rp 6.500</b></li>
<li>Transfer antar pengguna MBankingCore: <b>GRATIS</b></li>
<li>Cek saldo dan mutasi: <b>GRATIS</b></li>
<li>Pembayaran tagihan: <b>Rp 2.500</b></li>
</ul>

<h3>7.2 Perubahan Tarif</h3>
<ul>
<li>Kami berhak mengubah tarif dengan pemberitahuan 30 hari sebelumnya</li>
<li>Perubahan akan diberitahukan melalui Aplikasi atau email</li>
</ul>

<h2>8. PRIVASI DAN PERLINDUNGAN DATA</h2>
<h3>8.1 Pengumpulan Data</h3>
<ul>
<li>Data pribadi dikumpulkan sesuai keperluan layanan</li>
<li>Data transaksi disimpan untuk audit dan compliance</li>
<li>Lokasi perangkat untuk keamanan tambahan</li>
</ul>

<h3>8.2 Penggunaan Data</h3>
<ul>
<li>Memproses transaksi dan layanan</li>
<li>Analisis risiko dan fraud detection</li>
<li>Peningkatan layanan dan fitur</li>
<li>Compliance dengan regulasi</li>
</ul>

<h2>9. TANGGUNG JAWAB DAN GANTI RUGI</h2>
<h3>9.1 Batasan Tanggung Jawab</h3>
<ul>
<li>Tidak bertanggung jawab atas kerugian akibat kelalaian pengguna</li>
<li>Tidak bertanggung jawab atas gangguan jaringan atau sistem bank</li>
<li>Tanggung jawab terbatas pada jumlah transaksi yang bermasalah</li>
</ul>

<h3>9.2 Force Majeure</h3>
<ul>
<li>Tidak bertanggung jawab atas kejadian di luar kendali</li>
<li>Termasuk bencana alam, perang, atau gangguan pemerintah</li>
</ul>

<h2>10. PENANGGUHAN DAN PENGHENTIAN</h2>
<h3>10.1 Penangguhan Akun</h3>
<ul>
<li>Akun dapat ditangguhkan jika melanggar S&K</li>
<li>Penangguhan karena aktivitas mencurigakan</li>
<li>Pemberitahuan akan diberikan jika memungkinkan</li>
</ul>

<h3>10.2 Penghentian Layanan</h3>
<ul>
<li>Pengguna dapat menghentikan layanan kapan saja</li>
<li>Kami dapat menghentikan layanan dengan pemberitahuan 30 hari</li>
<li>Saldo akan dikembalikan sesuai prosedur bank</li>
</ul>

<h2>11. PERUBAHAN SYARAT DAN KETENTUAN</h2>
<ul>
<li>S&K dapat diubah sewaktu-waktu</li>
<li>Perubahan akan diberitahukan melalui Aplikasi</li>
<li>Penggunaan Aplikasi setelah perubahan dianggap sebagai persetujuan</li>
<li>Versi terbaru selalu tersedia di dalam Aplikasi</li>
</ul>

<h2>12. PENYELESAIAN SENGKETA</h2>
<h3>12.1 Hukum yang Berlaku</h3>
<ul>
<li>S&K ini tunduk pada hukum Republik Indonesia</li>
<li>Penyelesaian sengketa melalui pengadilan di Jakarta</li>
</ul>

<h3>12.2 Mediasi</h3>
<ul>
<li>Upaya penyelesaian secara kekeluargaan terlebih dahulu</li>
<li>Mediasi melalui Otoritas Jasa Keuangan (OJK) jika diperlukan</li>
</ul>

<h2>13. KONTAK DAN BANTUAN</h2>
<h3>Customer Service:</h3>
<ul>
<li>Email: <b>support@mbankingcore.com</b></li>
<li>Telepon: <b>1500-888 (24/7)</b></li>
<li>WhatsApp: <b>+62-812-3456-7890</b></li>
<li>Live Chat: Tersedia di dalam Aplikasi</li>
</ul>

<h3>Jam Operasional:</h3>
<ul>
<li>Senin - Jumat: 06.00 - 22.00 WIB</li>
<li>Sabtu - Minggu: 08.00 - 20.00 WIB</li>
<li>Emergency Support: 24/7</li>
</ul>

<hr>
<p><i>Dengan menggunakan Aplikasi MBankingCore, Anda menyatakan telah membaca, memahami, dan menyetujui seluruh Syarat dan Ketentuan di atas.</i></p>`
}

// getPrivacyPolicyContent returns comprehensive Privacy Policy content for Indonesian banking app
func getPrivacyPolicyContent() string {
	return `<h1>KEBIJAKAN PRIVASI MBANKINGCORE</h1>

<p><b>Terakhir diperbarui: 31 Juli 2025</b></p>

<h2>1. PENDAHULUAN</h2>
<p>MBankingCore ("<b>kami</b>", "<b>Aplikasi</b>") berkomitmen untuk melindungi privasi dan keamanan data pribadi Anda. Kebijakan Privasi ini menjelaskan bagaimana kami mengumpulkan, menggunakan, menyimpan, dan melindungi informasi pribadi Anda saat menggunakan layanan kami.</p>

<h2>2. INFORMASI YANG KAMI KUMPULKAN</h2>
<h3>2.1 Data Identitas Pribadi</h3>
<ul>
<li><b>Informasi Dasar:</b> Nama lengkap, tanggal lahir, nomor KTP/NIK</li>
<li><b>Kontak:</b> Nomor telepon, alamat email, alamat rumah</li>
<li><b>Foto:</b> Foto selfie, foto KTP untuk verifikasi identitas</li>
<li><b>Biometrik:</b> Sidik jari, face recognition (jika diaktifkan)</li>
</ul>

<h3>2.2 Data Keuangan</h3>
<ul>
<li><b>Informasi Rekening:</b> Nomor rekening, nama bank, saldo</li>
<li><b>Riwayat Transaksi:</b> Transfer, pembayaran, top-up, withdrawals</li>
<li><b>Pola Penggunaan:</b> Frekuensi transaksi, merchant favorit</li>
<li><b>Data Kredit:</b> Riwayat kredit untuk scoring (jika tersedia)</li>
</ul>

<h3>2.3 Data Teknis</h3>
<ul>
<li><b>Informasi Perangkat:</b> Model, OS, versi aplikasi, device ID</li>
<li><b>Lokasi:</b> GPS location untuk keamanan transaksi</li>
<li><b>Log Aktivitas:</b> Waktu login, IP address, aktivitas dalam aplikasi</li>
<li><b>Cookies:</b> Preferensi pengguna, session management</li>
</ul>

<h3>2.4 Data Komunikasi</h3>
<ul>
<li><b>Customer Service:</b> Rekaman chat, email, telepon</li>
<li><b>Notifikasi:</b> Preferensi push notification, SMS</li>
<li><b>Feedback:</b> Rating, review, saran perbaikan</li>
</ul>

<h2>3. CARA KAMI MENGUMPULKAN DATA</h2>
<h3>3.1 Langsung dari Anda</h3>
<ul>
<li>Saat registrasi akun baru</li>
<li>Pengisian profil dan verifikasi KYC</li>
<li>Melakukan transaksi atau menggunakan fitur</li>
<li>Menghubungi customer service</li>
</ul>

<h3>3.2 Otomatis dari Aplikasi</h3>
<ul>
<li>Log aktivitas dan penggunaan aplikasi</li>
<li>Data lokasi (dengan izin)</li>
<li>Informasi perangkat dan jaringan</li>
<li>Cookies dan teknologi tracking</li>
</ul>

<h3>3.3 Dari Pihak Ketiga</h3>
<ul>
<li><b>Bank Partner:</b> Informasi rekening dan transaksi</li>
<li><b>Credit Bureau:</b> Riwayat kredit dan scoring</li>
<li><b>Anti-Fraud Provider:</b> Verifikasi identitas dan risk assessment</li>
<li><b>Analytics Provider:</b> Data agregat untuk improvement</li>
</ul>

<h2>4. TUJUAN PENGGUNAAN DATA</h2>
<h3>4.1 Penyediaan Layanan</h3>
<ul>
<li><b>Verifikasi Identitas:</b> KYC, AML compliance, fraud prevention</li>
<li><b>Memproses Transaksi:</b> Transfer, pembayaran, top-up</li>
<li><b>Customer Support:</b> Bantuan teknis dan layanan pelanggan</li>
<li><b>Personalisasi:</b> Rekomendasi produk dan fitur yang relevan</li>
</ul>

<h3>4.2 Keamanan dan Kepatuhan</h3>
<ul>
<li><b>Fraud Detection:</b> Monitoring transaksi mencurigakan</li>
<li><b>Risk Management:</b> Penilaian risiko dan credit scoring</li>
<li><b>Regulatory Compliance:</b> Pelaporan ke Bank Indonesia, OJK</li>
<li><b>Audit Trail:</b> Jejak audit untuk investigasi</li>
</ul>

<h3>4.3 Peningkatan Layanan</h3>
<ul>
<li><b>Analytics:</b> Analisis penggunaan untuk improvement</li>
<li><b>A/B Testing:</b> Testing fitur baru untuk user experience</li>
<li><b>Machine Learning:</b> AI untuk fraud detection dan personalization</li>
<li><b>Research:</b> Riset pasar untuk pengembangan produk</li>
</ul>

<h2>5. BERBAGI DATA DENGAN PIHAK KETIGA</h2>
<h3>5.1 Bank dan Financial Institution</h3>
<ul>
<li><b>Bank Partner:</b> Untuk memproses transaksi perbankan</li>
<li><b>Payment Gateway:</b> Pembayaran merchant dan e-commerce</li>
<li><b>E-wallet Provider:</b> Top-up dan transfer antar e-wallet</li>
<li><b>Credit Bureau:</b> Credit scoring dan risk assessment</li>
</ul>

<h3>5.2 Technology Provider</h3>
<ul>
<li><b>Cloud Provider:</b> AWS, Google Cloud untuk data storage</li>
<li><b>Security Provider:</b> Anti-fraud, cybersecurity services</li>
<li><b>Analytics Provider:</b> Google Analytics, Firebase</li>
<li><b>Communication:</b> SMS gateway, email provider, push notification</li>
</ul>

<h2>6. KEAMANAN DATA</h2>
<h3>6.1 Enkripsi</h3>
<ul>
<li><b>Data at Rest:</b> AES-256 encryption untuk data storage</li>
<li><b>Data in Transit:</b> TLS 1.3 untuk komunikasi</li>
<li><b>Database:</b> Field-level encryption untuk data sensitif</li>
<li><b>Backup:</b> Encrypted backup dengan secure key management</li>
</ul>

<h3>6.2 Access Control</h3>
<ul>
<li><b>Role-based Access:</b> Akses terbatas sesuai job function</li>
<li><b>Multi-factor Authentication:</b> 2FA untuk admin access</li>
<li><b>Audit Log:</b> Complete logging untuk semua akses data</li>
<li><b>Privileged Access Management:</b> Secure admin access</li>
</ul>

<h2>7. HAK-HAK ANDA</h2>
<h3>7.1 Hak Akses</h3>
<ul>
<li><b>Data Portability:</b> Ekspor data dalam format standar</li>
<li><b>Data Transparency:</b> Informasi lengkap data yang dimiliki</li>
<li><b>Processing Activity:</b> Detail bagaimana data digunakan</li>
<li><b>Third Party Sharing:</b> List pihak ketiga yang menerima data</li>
</ul>

<h3>7.2 Hak Koreksi</h3>
<ul>
<li><b>Update Profile:</b> Self-service untuk update data pribadi</li>
<li><b>Data Correction:</b> Request koreksi data yang tidak akurat</li>
<li><b>Verification Process:</b> Proses verifikasi untuk data sensitif</li>
<li><b>Notification:</b> Pemberitahuan perubahan ke pihak ketiga</li>
</ul>

<h3>7.3 Hak Penghapusan</h3>
<ul>
<li><b>Account Deletion:</b> Penghapusan akun dan data terkait</li>
<li><b>Selective Deletion:</b> Penghapusan data spesifik</li>
<li><b>Retention Override:</b> Penghapusan sebelum periode retensi</li>
<li><b>Legal Basis:</b> Consideration terhadap kewajiban hukum</li>
</ul>

<h2>8. RETENSI DATA</h2>
<h3>8.1 Periode Penyimpanan</h3>
<ul>
<li><b>Data Transaksi:</b> 10 tahun (sesuai regulasi BI)</li>
<li><b>Data Identitas:</b> Selama akun aktif + 5 tahun</li>
<li><b>Log Komunikasi:</b> 3 tahun untuk audit purpose</li>
<li><b>Analytics Data:</b> 2 tahun dalam bentuk agregat</li>
</ul>

<h3>8.2 Penghapusan Data</h3>
<ul>
<li><b>Account Closure:</b> Data dihapus setelah periode retensi</li>
<li><b>Right to be Forgotten:</b> Penghapusan atas permintaan</li>
<li><b>Secure Deletion:</b> Cryptographic erasure dan overwriting</li>
<li><b>Certificate of Destruction:</b> Bukti penghapusan data</li>
</ul>

<h2>9. COOKIES DAN TRACKING</h2>
<h3>9.1 Jenis Cookies</h3>
<ul>
<li><b>Essential Cookies:</b> Untuk fungsi dasar aplikasi</li>
<li><b>Performance Cookies:</b> Analytics dan monitoring</li>
<li><b>Functional Cookies:</b> Preferensi dan personalization</li>
<li><b>Marketing Cookies:</b> Targeted advertising</li>
</ul>

<h3>9.2 Cookie Management</h3>
<ul>
<li><b>Cookie Settings:</b> Control di aplikasi settings</li>
<li><b>Browser Settings:</b> Disable cookies di browser</li>
<li><b>Third Party Cookies:</b> Opt-out dari advertising cookies</li>
<li><b>Cookie Policy:</b> Detail lengkap di cookie policy page</li>
</ul>

<h2>10. ANAK DI BAWAH UMUR</h2>
<ul>
<li>Layanan tidak ditujukan untuk anak di bawah 17 tahun</li>
<li>Verifikasi usia saat registrasi</li>
<li>Immediate deletion jika ditemukan data anak</li>
<li>Parental consent untuk usia 17-21 tahun</li>
</ul>

<h2>11. PERUBAHAN KEBIJAKAN PRIVASI</h2>
<h3>11.1 Notifikasi Perubahan</h3>
<ul>
<li><b>Material Changes:</b> Email notification 30 hari sebelumnya</li>
<li><b>Minor Updates:</b> In-app notification</li>
<li><b>Version History:</b> Archive versi sebelumnya</li>
<li><b>Continued Use:</b> Deemed acceptance jika tetap menggunakan</li>
</ul>

<h2>12. KONTAK DATA PROTECTION</h2>
<h3>12.1 Data Protection Officer</h3>
<ul>
<li><b>Email:</b> privacy@mbankingcore.com</li>
<li><b>Telepon:</b> +62-21-5000-1234</li>
<li><b>Alamat:</b> Menara MBankingCore, Jakarta 12345</li>
<li><b>Response Time:</b> Maksimal 30 hari untuk complex request</li>
</ul>

<h3>12.2 Complaint Process</h3>
<ul>
<li><b>Internal Complaint:</b> Melalui customer service</li>
<li><b>Regulator Complaint:</b> Ke Kementerian Kominfo</li>
<li><b>International:</b> GDPR representative untuk EU residents</li>
<li><b>Escalation:</b> Clear escalation process untuk unresolved issues</li>
</ul>

<h2>13. KETENTUAN KHUSUS</h2>
<h3>13.1 Emerging Technology</h3>
<ul>
<li><b>AI/ML Ethics:</b> Responsible AI untuk decision making</li>
<li><b>Blockchain:</b> Privacy consideration untuk DLT</li>
<li><b>IoT Integration:</b> Security untuk connected devices</li>
<li><b>Quantum Computing:</b> Quantum-safe cryptography preparation</li>
</ul>

<hr>
<p><b>Efektif sejak:</b> 31 Juli 2025<br>
<b>Versi:</b> 2.1<br>
<b>Bahasa:</b> Bahasa Indonesia (versi resmi), English (reference)</p>

<p><i>Dengan menggunakan layanan MBankingCore, Anda menyatakan telah membaca, memahami, dan menyetujui Kebijakan Privasi ini.</i></p>`
}

// seedDummyTransactions2020To2025 creates 20000 dummy transactions from 2020-2025
func seedDummyTransactions2020To2025() error {
	log.Println("Seeding 20000 dummy transactions from 2020-2025...")

	// Check if transactions already exist (check for total count >= 20000)
	var count int64
	DB.Model(&models.Transaction{}).Count(&count)
	if count >= 20000 {
		log.Printf("✅ %d transactions already exist, skipping 2020-2025 dummy data", count)
		return nil
	}

	// Get all user IDs for transaction generation
	var userIDs []uint
	if err := DB.Model(&models.User{}).Pluck("id", &userIDs).Error; err != nil {
		log.Printf("Failed to get user IDs: %v", err)
		return err
	}

	if len(userIDs) == 0 {
		log.Println("No users found, skipping transaction generation")
		return nil
	}

	// Time range: 2020-2025
	startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	timeRange := endTime.Unix() - startTime.Unix()

	// Transaction types and their weights
	transactionTypes := []struct {
		Type   string
		Weight int
	}{
		{"topup", 30},        // 30% topup
		{"withdraw", 25},     // 25% withdraw
		{"transfer_out", 20}, // 20% transfer out
		{"transfer_in", 20},  // 20% transfer in
		{"payment", 5},       // 5% payment
	}

	// Transaction descriptions by type
	descriptions := map[string][]string{
		"topup": {
			"Top up via m-banking",
			"Top up gaji bulanan",
			"Top up bonus kerja",
			"Isi saldo untuk transaksi",
			"Transfer dari rekening bank",
			"Top up via ATM",
			"Isi saldo otomatis",
		},
		"withdraw": {
			"Tarik tunai untuk keperluan harian",
			"Withdraw untuk ongkos",
			"Tarik tunai di ATM",
			"Withdraw untuk belanja",
			"Tarik tunai emergency",
			"Withdraw untuk bayar parkir",
			"Tarik tunai weekend",
		},
		"transfer_out": {
			"Transfer ke keluarga",
			"Transfer untuk split bill",
			"Transfer ke teman",
			"Bayar hutang",
			"Transfer untuk hadiah",
			"Transfer ke orang tua",
			"Transfer untuk investasi",
		},
		"transfer_in": {
			"Transfer dari keluarga",
			"Terima transfer teman",
			"Transfer gaji",
			"Terima bayaran hutang",
			"Transfer bonus",
			"Terima hadiah ulang tahun",
			"Transfer dari anak",
		},
		"payment": {
			"Bayar tagihan listrik",
			"Bayar pulsa",
			"Bayar internet",
			"Bayar PDAM",
			"Bayar asuransi",
			"Bayar cicilan",
			"Bayar belanja online",
		},
	}

	// Status distribution
	statuses := []struct {
		Status string
		Weight int
	}{
		{"completed", 85}, // 85% completed
		{"pending", 10},   // 10% pending
		{"failed", 5},     // 5% failed
	}

	var transactions []models.Transaction
	batchSize := 500 // Process in batches

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 20000; i++ {
		// Random user
		userID := userIDs[rand.Intn(len(userIDs))]

		// Random transaction time between 2020-2025
		randomTime := startTime.Add(time.Duration(rand.Int63n(timeRange)) * time.Second)

		// Random transaction type based on weights
		var txType string
		totalWeight := 0
		for _, t := range transactionTypes {
			totalWeight += t.Weight
		}
		randomWeight := rand.Intn(totalWeight)
		currentWeight := 0
		for _, t := range transactionTypes {
			currentWeight += t.Weight
			if randomWeight < currentWeight {
				txType = t.Type
				break
			}
		}

		// Random amount based on transaction type
		var amount int64
		switch txType {
		case "topup":
			// Topup: 50K - 5M IDR
			amount = int64(rand.Intn(4950000) + 50000)
		case "withdraw":
			// Withdraw: 20K - 2M IDR
			amount = int64(rand.Intn(1980000) + 20000)
		case "transfer_out", "transfer_in":
			// Transfer: 10K - 3M IDR
			amount = int64(rand.Intn(2990000) + 10000)
		case "payment":
			// Payment: 5K - 1M IDR
			amount = int64(rand.Intn(995000) + 5000)
		default:
			amount = int64(rand.Intn(500000) + 50000)
		}

		// Random balance before/after (simplified for dummy data)
		balanceBefore := int64(rand.Intn(10000000) + 100000) // 100K - 10M IDR
		var balanceAfter int64
		switch txType {
		case "topup", "transfer_in":
			balanceAfter = balanceBefore + amount
		case "withdraw", "transfer_out", "payment":
			balanceAfter = balanceBefore - amount
			if balanceAfter < 0 {
				balanceAfter = 0
			}
		default:
			balanceAfter = balanceBefore
		}

		// Random status based on weights
		var status string
		totalStatusWeight := 0
		for _, s := range statuses {
			totalStatusWeight += s.Weight
		}
		randomStatusWeight := rand.Intn(totalStatusWeight)
		currentStatusWeight := 0
		for _, s := range statuses {
			currentStatusWeight += s.Weight
			if randomStatusWeight < currentStatusWeight {
				status = s.Status
				break
			}
		}

		// Random description
		descList := descriptions[txType]
		description := descList[rand.Intn(len(descList))]

		// Create transaction
		transaction := models.Transaction{
			UserID:        userID,
			Type:          txType,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Description:   description,
			Status:        status,
			CreatedAt:     randomTime,
			UpdatedAt:     randomTime,
		}

		transactions = append(transactions, transaction)

		// Insert in batches
		if len(transactions) >= batchSize || i == 20000-1 {
			if err := DB.CreateInBatches(&transactions, batchSize).Error; err != nil {
				log.Printf("Error inserting transaction batch: %v", err)
				continue
			}

			log.Printf("Inserted transaction batch: %d/20000 transactions", i+1)
			transactions = []models.Transaction{} // Reset batch
		}
	}

	log.Printf("✅ Created 20000 dummy transactions from 2020-2025")

	// Show summary statistics
	showTransactionStatistics()
	return nil
}

// showTransactionStatistics displays transaction statistics after seeding
func showTransactionStatistics() {
	log.Println("\n=== Transaction Statistics ===")

	var totalTransactions int64
	DB.Model(&models.Transaction{}).Count(&totalTransactions)
	log.Printf("Total transactions in database: %d", totalTransactions)

	// Count by year
	for year := 2020; year <= 2025; year++ {
		var yearCount int64
		startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

		DB.Model(&models.Transaction{}).Where("created_at BETWEEN ? AND ?", startOfYear, endOfYear).Count(&yearCount)
		log.Printf("Transactions in %d: %d", year, yearCount)
	}

	// Count by type
	var topupCount, withdrawCount, transferOutCount, transferInCount, paymentCount int64
	DB.Model(&models.Transaction{}).Where("type = ?", "topup").Count(&topupCount)
	DB.Model(&models.Transaction{}).Where("type = ?", "withdraw").Count(&withdrawCount)
	DB.Model(&models.Transaction{}).Where("type = ?", "transfer_out").Count(&transferOutCount)
	DB.Model(&models.Transaction{}).Where("type = ?", "transfer_in").Count(&transferInCount)
	DB.Model(&models.Transaction{}).Where("type = ?", "payment").Count(&paymentCount)

	log.Printf("Topup transactions: %d", topupCount)
	log.Printf("Withdraw transactions: %d", withdrawCount)
	log.Printf("Transfer out transactions: %d", transferOutCount)
	log.Printf("Transfer in transactions: %d", transferInCount)
	log.Printf("Payment transactions: %d", paymentCount)

	// Count by status
	var completedCount, pendingCount, failedCount int64
	DB.Model(&models.Transaction{}).Where("status = ?", "completed").Count(&completedCount)
	DB.Model(&models.Transaction{}).Where("status = ?", "pending").Count(&pendingCount)
	DB.Model(&models.Transaction{}).Where("status = ?", "failed").Count(&failedCount)

	log.Printf("Completed transactions: %d", completedCount)
	log.Printf("Pending transactions: %d", pendingCount)
	log.Printf("Failed transactions: %d", failedCount)
}
