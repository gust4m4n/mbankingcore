package main

import (
	"fmt"
	"log"
	"math/rand"
	"mbankingcore/config"
	"mbankingcore/models"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("🚀 Starting user generation script...")

	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("⚠️  Warning: Error loading .env file: %v", err)
		log.Println("Attempting to use system environment variables...")
	}

	// Connect to database
	config.ConnectDatabase()
	if config.DB == nil {
		log.Fatal("❌ Failed to connect to database")
	}
	log.Println("✅ Database connected successfully")

	// Generate users
	if err := generateRandomUsers(); err != nil {
		log.Fatalf("❌ Failed to generate users: %v", err)
	}

	log.Println("🎉 User generation completed successfully!")
}

// generateRandomUsers creates 2000 dummy users with random creation dates from 2010-2025
func generateRandomUsers() error {
	log.Println("🔄 Generating 2000 dummy users with random creation dates (2010-2025)...")

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Indonesian names for realistic data
	firstNames := []string{
		"Adi", "Ahmad", "Andi", "Agus", "Budi", "Bambang", "Chandra", "Dani", "Dedi", "Eko",
		"Fadli", "Firman", "Gani", "Hadi", "Iwan", "Joko", "Krisna", "Lukas", "Made", "Nanang",
		"Oki", "Panji", "Rizki", "Sari", "Tono", "Udin", "Vicky", "Wawan", "Yanto", "Zaki",
		"Ayu", "Bella", "Citra", "Dewi", "Eka", "Fitri", "Gita", "Hana", "Indri", "Jihan",
		"Kartika", "Lina", "Maya", "Nina", "Okta", "Putri", "Rika", "Sinta", "Tari", "Utami",
		"Sinta", "Ratna", "Wulan", "Dewi", "Sari", "Indah", "Lestari", "Rahayu", "Permatasari",
		"Kusuma", "Handayani", "Maharani", "Anggraini", "Pratiwi", "Wulandari", "Oktaviani",
	}

	lastNames := []string{
		"Pratama", "Santoso", "Wijaya", "Kusuma", "Purnama", "Sari", "Dewi", "Lestari", "Wati",
		"Putri", "Rahayu", "Indah", "Permana", "Saputra", "Kurniawan", "Nugroho", "Setiawan",
		"Handoko", "Susanto", "Mulyadi", "Hakim", "Rahman", "Hidayat", "Maulana", "Syahputra",
		"Gunawan", "Utomo", "Wardana", "Suryadi", "Firdaus", "Munawar", "Sutrisno", "Hartono",
		"Sulistyo", "Pranoto", "Suharto", "Suroso", "Budiono", "Supriadi", "Darmawan", "Syamsudin",
		"Prasetyo", "Iskandar", "Priyanto", "Hermawan", "Widodo", "Suryanto", "Wahyudi", "Subagyo",
		"Purwanto", "Riyanto", "Haryanto", "Wahyu", "Agung", "Baskoro", "Cahyadi", "Firmansyah",
	}

	motherNames := []string{
		"Siti Aminah", "Fatimah", "Khadijah", "Aisyah", "Maryam", "Zainab", "Ruqayyah", "Ummu Kulsum",
		"Sri Rejeki", "Tuti Handayani", "Endang Susilowati", "Niken Sari", "Ratna Dewi", "Umi Kalsum",
		"Wulan Dari", "Sumiati", "Murniati", "Sulastri", "Nurhayati", "Suryani", "Indrawati",
		"Purwanti", "Susanti", "Hartini", "Lestari", "Ratnawati", "Widowati", "Mulyani", "Suhartini",
		"Nur Hayati", "Tri Wahyuni", "Dwi Astuti", "Eni Suryani", "Yuni Rahayu", "Lies Susanti",
		"Rini Handayani", "Rina Marlina", "Wiwik Purwanti", "Nanik Sulistyowati", "Tutik Rahayu",
	}

	// Check current max ID to avoid conflicts
	var maxID uint
	config.DB.Model(&models.User{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID)
	log.Printf("📊 Current max user ID: %d", maxID)

	users := make([]models.User, 2000)

	for i := 0; i < 2000; i++ {
		// Generate random creation date between 2010-2025
		startYear := 2010
		endYear := 2025
		randomYear := startYear + rand.Intn(endYear-startYear+1)
		randomMonth := 1 + rand.Intn(12)
		var randomDay int

		// Handle different days per month
		switch randomMonth {
		case 2: // February
			if randomYear%4 == 0 && (randomYear%100 != 0 || randomYear%400 == 0) {
				randomDay = 1 + rand.Intn(29) // Leap year
			} else {
				randomDay = 1 + rand.Intn(28)
			}
		case 4, 6, 9, 11: // April, June, September, November
			randomDay = 1 + rand.Intn(30)
		default: // All other months
			randomDay = 1 + rand.Intn(31)
		}

		randomHour := rand.Intn(24)
		randomMinute := rand.Intn(60)
		randomSecond := rand.Intn(60)

		createdAt := time.Date(randomYear, time.Month(randomMonth), randomDay, randomHour, randomMinute, randomSecond, 0, time.UTC)

		// Generate user data
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		fullName := firstName + " " + lastName

		// Generate unique phone number with retry logic
		var phoneNumber string
		for attempts := 0; attempts < 10; attempts++ {
			phoneNumber = fmt.Sprintf("08%d%d%d%d%d%d%d%d",
				rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
				rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))

			// Check if phone number already exists
			var count int64
			config.DB.Model(&models.User{}).Where("phone = ?", phoneNumber).Count(&count)
			if count == 0 {
				break // Phone number is unique
			}
		}

		// Generate mother name
		motherName := motherNames[rand.Intn(len(motherNames))]

		// Generate PIN (6 digits)
		pin := fmt.Sprintf("%06d", rand.Intn(1000000))
		hashedPIN, _ := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)

		// Generate balance (100k to 50M rupiah)
		balance := int64(100000 + rand.Intn(49900000))

		// Generate status (95% active, 3% inactive, 2% blocked)
		var status int
		randomStatus := rand.Intn(100)
		if randomStatus < 95 {
			status = models.USER_STATUS_ACTIVE
		} else if randomStatus < 98 {
			status = models.USER_STATUS_INACTIVE
		} else {
			status = models.USER_STATUS_BLOCKED
		}

		users[i] = models.User{
			Name:       fullName,
			Phone:      phoneNumber,
			MotherName: motherName,
			PinAtm:     string(hashedPIN),
			Balance:    balance,
			Status:     status,
			Avatar:     "",
			CreatedAt:  createdAt,
			UpdatedAt:  createdAt,
		}
	}

	// Batch insert users in chunks of 100
	batchSize := 100
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		batch := users[i:end]
		if err := config.DB.Create(&batch).Error; err != nil {
			log.Printf("❌ Failed to create user batch %d-%d: %v", i+1, end, err)
			return err
		}
		log.Printf("✅ Created users %d-%d", i+1, end)
	}

	// Display statistics
	log.Println("\n📊 User generation statistics:")

	// Count users by year
	for year := 2010; year <= 2025; year++ {
		var count int64
		startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)

		config.DB.Model(&models.User{}).Where("created_at BETWEEN ? AND ?", startOfYear, endOfYear).Count(&count)
		percentage := float64(count) / 20.0
		log.Printf("Year %d: %d users (%.1f%%)", year, count, percentage)
	}

	// Count by status
	var activeCount, inactiveCount, blockedCount int64
	config.DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_ACTIVE).Count(&activeCount)
	config.DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_INACTIVE).Count(&inactiveCount)
	config.DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_BLOCKED).Count(&blockedCount)

	log.Printf("\nStatus distribution:")
	log.Printf("- Active: %d (%.1f%%)", activeCount, float64(activeCount)/20.0)
	log.Printf("- Inactive: %d (%.1f%%)", inactiveCount, float64(inactiveCount)/20.0)
	log.Printf("- Blocked: %d (%.1f%%)", blockedCount, float64(blockedCount)/20.0)

	// Total count verification
	var totalCount int64
	config.DB.Model(&models.User{}).Count(&totalCount)
	log.Printf("\n✅ Total users in database: %d", totalCount)
	log.Printf("✅ New users created: 2000")

	return nil
}
