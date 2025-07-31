package config

import (
	"fmt"
	"log"
	"math/rand"
	"mbankingcore/models"
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

	// DISABLED: Seed initial users for testing
	// if err := seedInitialUsers(); err != nil {
	// 	return err
	// }

	// DISABLED: Seed initial transactions for testing
	// if err := seedInitialTransactions(); err != nil {
	// 	return err
	// }

	// DISABLED: Seed 1000 dummy transactions for yearly data
	// if err := seedDummyTransactionsYearly(); err != nil {
	// 	return err
	// }

	// DISABLED: Clear and seed 2000 dummy users for 2020-2025
	// if err := clearAndSeed2000Users(); err != nil {
	// 	return err
	// }

	// DISABLED: Seed 10000 dummy transactions for yearly data
	// if err := seedDummyTransactionsYearlyLarge(); err != nil {
	// 	return err
	// }

	// DISABLED: Clear existing transactions and seed new 20000 transactions from 2020-2025 with 300% growth
	// if err := clearAndSeedTransactionsWithGrowth(); err != nil {
	// 	return err
	// }

	// Generate 10,000 dummy transactions from 2010 to 2025
	if err := seed10KTransactions2010to2025(); err != nil {
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

// clearAndSeed2000Users clears existing users and creates 2000 dummy users distributed from 2020-2025
func clearAndSeed2000Users() error {
	log.Println("Clearing existing users and seeding 2000 dummy users from 2020-2025...")

	// Clear existing transactions first (due to foreign key constraint)
	if err := DB.Exec("DELETE FROM transactions").Error; err != nil {
		log.Printf("Failed to clear existing transactions: %v", err)
		return err
	}
	log.Println("✅ Cleared existing transactions")

	// Clear existing users
	if err := DB.Exec("DELETE FROM users").Error; err != nil {
		log.Printf("Failed to clear existing users: %v", err)
		return err
	}
	log.Println("✅ Cleared existing users")

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
		"Adrian", "Bianca", "Carlos", "Dinda", "Elena", "Felix", "Grace", "Henry",
		"Irene", "Javier", "Karina", "Lucas", "Monica", "Nathan", "Olivia", "Patrick",
		"Queen", "Roberto", "Sofia", "Thomas", "Ursula", "Victor", "Wendy", "Xavier",
		"Yasmin", "Zidane", "Amanda", "Bryan", "Charlotte", "David", "Emma", "Fabio",
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
		"Baskara", "Cantika", "Darmawan", "Ester", "Farhan", "Ghina", "Hasibuan", "Indira",
		"Jefri", "Kencana", "Langit", "Maheswari", "Nirwana", "Oktavia", "Pradana", "Qadira",
	}

	motherNames := []string{
		"Siti", "Sri", "Umi", "Ibu", "Nyai", "Ratu", "Dewi", "Putri",
		"Ratna", "Indira", "Kartika", "Melati", "Mawar", "Cempaka", "Anggrek", "Dahlia",
		"Kenanga", "Seruni", "Tulip", "Sakura", "Lily", "Rose", "Jasmine", "Orchid",
		"Bunga", "Cantika", "Anggun", "Jelita", "Ayu", "Endah", "Permata", "Intan",
		"Sari", "Wulan", "Lestari", "Cahaya", "Bintang", "Fajar", "Mega", "Citra",
		"Rina", "Nina", "Dina", "Tina", "Mira", "Dira", "Kira", "Lira",
	}

	// Time range: 2020-2025
	startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	timeRange := endTime.Unix() - startTime.Unix()

	// Track used phone numbers to ensure uniqueness
	usedPhones := make(map[string]bool)

	var users []models.User
	batchSize := 100 // Insert in batches for better performance

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 2000; i++ {
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
		if len(users) >= batchSize || i == 2000-1 {
			if err := DB.CreateInBatches(&users, batchSize).Error; err != nil {
				log.Printf("Error inserting batch: %v", err)
				continue
			}

			log.Printf("Inserted batch: %d/2000 users", i+1)
			users = []models.User{} // Reset batch
		}
	}

	log.Printf("✅ Created 2000 dummy users with registration dates from 2020-2025")

	// Show summary statistics
	show2000UserStatistics()
	return nil
}

// show2000UserStatistics displays user statistics after seeding 2000 users
func show2000UserStatistics() {
	log.Println("\n=== User Statistics (2000 Users) ===")

	var totalUsers int64
	DB.Model(&models.User{}).Count(&totalUsers)
	log.Printf("Total users in database: %d", totalUsers)

	// Count by year
	var totalGrowth float64
	var baselineYear int64
	for year := 2020; year <= 2025; year++ {
		var yearCount int64
		startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

		DB.Model(&models.User{}).Where("created_at BETWEEN ? AND ?", startOfYear, endOfYear).Count(&yearCount)
		log.Printf("Users registered in %d: %d", year, yearCount)

		if year == 2020 {
			baselineYear = yearCount
		} else if year == 2025 && baselineYear > 0 {
			totalGrowth = (float64(yearCount)/float64(baselineYear) - 1) * 100
		}
	}

	if baselineYear > 0 {
		log.Printf("Total Growth (2020-2025): %.1f%%", totalGrowth)
	}

	// Count by status
	var activeUsers, inactiveUsers, blockedUsers int64
	DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_ACTIVE).Count(&activeUsers)
	DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_INACTIVE).Count(&inactiveUsers)
	DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_BLOCKED).Count(&blockedUsers)

	log.Printf("Active users: %d (%.1f%%)", activeUsers, float64(activeUsers)/float64(totalUsers)*100)
	log.Printf("Inactive users: %d (%.1f%%)", inactiveUsers, float64(inactiveUsers)/float64(totalUsers)*100)
	log.Printf("Blocked users: %d (%.1f%%)", blockedUsers, float64(blockedUsers)/float64(totalUsers)*100)
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

// seed10KTransactions2010to2025 creates 10,000 dummy transactions from 2010 to 2025
func seed10KTransactions2010to2025() error {
	log.Println("Creating 10,000 dummy transactions from 2010 to 2025...")

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

	log.Printf("Found %d users for transaction generation", len(userIDs))

	// Clear existing transactions first
	log.Println("Clearing existing transactions...")
	if err := DB.Exec("DELETE FROM transactions").Error; err != nil {
		log.Printf("Failed to clear transactions: %v", err)
		return err
	}
	log.Println("✅ All existing transactions cleared")

	// Calculate transaction distribution by year (2010-2025)
	totalTransactions := 10000
	yearCount := 16 // 2010 to 2025 inclusive
	basePerYear := totalTransactions / yearCount

	// Distribute transactions with some variation
	yearDistribution := make(map[int]int)
	remaining := totalTransactions

	for year := 2010; year <= 2025; year++ {
		if year == 2025 { // Last year gets remaining
			yearDistribution[year] = remaining
		} else {
			// Add some variation ±20%
			variation := rand.Float64()*0.4 - 0.2 // -20% to +20%
			count := int(float64(basePerYear) * (1 + variation))
			if count < 200 {
				count = 200 // Minimum 200 per year
			}
			if count > remaining-15+year-2010 {
				count = remaining - 15 + year - 2010
			}
			yearDistribution[year] = count
			remaining -= count
		}
	}

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

	// Enhanced transaction descriptions by type
	descriptions := map[string][]string{
		"topup": {
			"Top up via m-banking", "Top up gaji bulanan", "Top up bonus kerja", "Top up via ATM BCA",
			"Top up via ATM BRI", "Top up via ATM BNI", "Top up via ATM Mandiri", "Top up untuk traveling",
			"Top up untuk belanja online", "Top up bonus kerja", "Top up gaji bulanan", "Top up dari rekening utama",
			"Top up via mobile banking", "Top up via internet banking", "Transfer masuk dari rekening lain",
			"Setor tunai via teller", "Top up dari kartu kredit", "Top up untuk investasi",
		},
		"withdraw": {
			"Tarik tunai untuk keperluan harian", "Tarik tunai untuk ongkos", "Tarik tunai di ATM",
			"Tarik tunai untuk belanja", "Tarik tunai darurat", "Tarik tunai weekend", "Tarik tunai untuk bensin",
			"Tarik tunai untuk jajan", "Tarik tunai untuk shopping", "Penarikan tunai di ATM",
			"Penarikan untuk belanja bulanan", "Penarikan untuk bayar warung", "Tarik tunai untuk keperluan",
			"Penarikan dana darurat", "Tarik tunai untuk liburan",
		},
		"transfer_out": {
			"Transfer ke keluarga", "Transfer ke teman - Bayar hutang", "Transfer ke orang tua",
			"Transfer ke Maya - Split bill restoran", "Transfer ke Dina - Bayar patungan gift",
			"Transfer donasi", "Transfer untuk pembayaran", "Bayar tagihan listrik", "Bayar tagihan kartu kredit",
			"Bayar tagihan internet", "Bayar asuransi bulanan", "Bayar sewa kos bulanan", "Bayar SPP kuliah",
			"Bayar cicilan motor", "Transfer ke vendor", "Bayar tagihan PDAM", "Transfer ke supplier",
		},
		"transfer_in": {
			"Transfer dari keluarga", "Transfer dari teman", "Transfer dari klien", "Terima gaji bulanan",
			"Transfer dari Demo User - Bayar makan siang", "Transfer dari Test User - Bayar makan siang",
			"Transfer dari Sari - Split bill restoran", "Transfer dari Rizki - Bayar patungan gift",
			"Terima hadiah ulang tahun", "Terima bonus penjualan", "Terima cashback", "Terima refund belanja",
			"Transfer komisi", "Terima pembayaran freelance", "Terima dividen", "Terima royalti",
		},
		"payment": {
			"Bayar tagihan listrik", "Bayar pulsa", "Bayar tagihan internet", "Bayar PDAM",
			"Bayar asuransi", "Bayar cicilan", "Bayar belanja online", "Bayar tagihan kartu kredit",
			"Bayar Netflix", "Bayar Spotify", "Bayar Google Play", "Bayar Apple Store",
		},
	}

	// Status distribution
	statuses := []struct {
		Status string
		Weight int
	}{
		{"completed", 90}, // 90% completed
		{"failed", 10},    // 10% failed
	}

	rand.Seed(time.Now().UnixNano())
	batchSize := 500
	transactionID := 1

	log.Printf("Starting transaction generation across %d years...", yearCount)

	// Generate transactions year by year
	for year := 2010; year <= 2025; year++ {
		yearTransactions := yearDistribution[year]
		log.Printf("Generating %d transactions for year %d...", yearTransactions, year)

		var transactions []models.Transaction

		// Distribute transactions across months
		monthlyDistribution := make([]int, 12)
		basePerMonth := yearTransactions / 12
		remainingForYear := yearTransactions

		for month := 0; month < 12; month++ {
			if month == 11 { // Last month gets remaining
				monthlyDistribution[month] = remainingForYear
			} else {
				// Random fluctuation ±15%
				fluctuation := rand.Float64()*0.3 - 0.15 // -15% to +15%
				monthCount := int(float64(basePerMonth) * (1 + fluctuation))
				if monthCount < 1 {
					monthCount = 1
				}
				if monthCount > remainingForYear-11+month {
					monthCount = remainingForYear - 11 + month
				}
				monthlyDistribution[month] = monthCount
				remainingForYear -= monthCount
			}
		}

		// Generate transactions for each month
		for month := 0; month < 12; month++ {
			monthCount := monthlyDistribution[month]

			for i := 0; i < monthCount; i++ {
				// Random user
				userID := userIDs[rand.Intn(len(userIDs))]

				// Random day and time within the month
				startOfMonth := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
				var endOfMonth time.Time
				if month == 11 { // December
					endOfMonth = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
				} else {
					endOfMonth = time.Date(year, time.Month(month+2), 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
				}

				timeRange := endOfMonth.Unix() - startOfMonth.Unix()
				randomTime := startOfMonth.Add(time.Duration(rand.Int63n(timeRange)) * time.Second)

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

				// Amount calculation with year-based inflation (2% per year from 2010)
				inflationFactor := 1.0 + float64(year-2010)*0.02 // 2% inflation per year
				var baseAmount int64

				switch txType {
				case "topup":
					// Topup: 50K - 5M IDR
					baseAmount = int64(rand.Intn(4950000) + 50000)
				case "withdraw":
					// Withdraw: 20K - 2M IDR
					baseAmount = int64(rand.Intn(1980000) + 20000)
				case "transfer_out", "transfer_in":
					// Transfer: 10K - 3M IDR
					baseAmount = int64(rand.Intn(2990000) + 10000)
				case "payment":
					// Payment: 5K - 1M IDR
					baseAmount = int64(rand.Intn(995000) + 5000)
				default:
					baseAmount = int64(rand.Intn(500000) + 50000)
				}

				amount := int64(float64(baseAmount) * inflationFactor)

				// Random balance calculation
				balanceBefore := int64(rand.Intn(20000000) + 100000) // 100K - 20M IDR
				var balanceAfter int64

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

				// Calculate balance after based on transaction type and status
				if status == "completed" {
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
				} else {
					balanceAfter = balanceBefore // Failed transactions don't change balance
				}

				// Random description with failure reason for failed transactions
				descList := descriptions[txType]
				description := descList[rand.Intn(len(descList))]

				if status == "failed" {
					failureReasons := []string{
						" - Saldo tidak mencukupi",
						" - Transaksi gagal",
						" - Timeout",
						" - Gangguan jaringan",
						" - Sistem maintenance",
					}
					description += failureReasons[rand.Intn(len(failureReasons))]
				}

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
				if len(transactions) >= batchSize {
					if err := DB.CreateInBatches(&transactions, batchSize).Error; err != nil {
						log.Printf("Error inserting transaction batch: %v", err)
						continue
					}
					log.Printf("Created transaction batch %d-%d for year %d", transactionID, transactionID+len(transactions)-1, year)
					transactionID += len(transactions)
					transactions = []models.Transaction{} // Reset batch
				}
			}
		}

		// Insert remaining transactions for this year
		if len(transactions) > 0 {
			if err := DB.CreateInBatches(&transactions, len(transactions)).Error; err != nil {
				log.Printf("Error inserting final transaction batch for year %d: %v", year, err)
			} else {
				log.Printf("Created final transaction batch %d-%d for year %d", transactionID, transactionID+len(transactions)-1, year)
				transactionID += len(transactions)
			}
		}

		log.Printf("✅ Completed %d transactions for year %d", yearTransactions, year)
	}

	// Show statistics
	log.Println("\n=== Transaction Generation Summary ===")
	var totalCount int64
	DB.Model(&models.Transaction{}).Count(&totalCount)
	log.Printf("Total transactions created: %d", totalCount)

	// Count by type
	for _, txType := range transactionTypes {
		var count int64
		DB.Model(&models.Transaction{}).Where("type = ?", txType.Type).Count(&count)
		percentage := float64(count) / float64(totalCount) * 100
		log.Printf("%s: %d transactions (%.1f%%)", txType.Type, count, percentage)
	}

	// Count by status
	var completedCount, failedCount int64
	DB.Model(&models.Transaction{}).Where("status = ?", "completed").Count(&completedCount)
	DB.Model(&models.Transaction{}).Where("status = ?", "failed").Count(&failedCount)
	log.Printf("Completed: %d (%.1f%%)", completedCount, float64(completedCount)/float64(totalCount)*100)
	log.Printf("Failed: %d (%.1f%%)", failedCount, float64(failedCount)/float64(totalCount)*100)

	log.Printf("✅ Successfully created 10,000 dummy transactions from 2010-2025")
	return nil
}
