package config

import (
	"log"
	"mbankingcore/models"

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
		{Key: "terms_conditions", Value: "Default terms and conditions content"},
		{Key: "privacy_policy", Value: "Default privacy policy content"},
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
