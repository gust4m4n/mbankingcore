package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"mbankingcore/config"
	"mbankingcore/models"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("MBankingCore - Dummy Users Generator")
	log.Println("===================================")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	// Connect to database
	config.ConnectDatabase()

	// Generate 5000 dummy users
	generateDummyUsers(5000)

	log.Println("Dummy users generation completed successfully!")
}

func generateDummyUsers(count int) {
	log.Printf("Generating %d dummy users...", count)

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

	for i := 0; i < count; i++ {
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
		if len(users) >= batchSize || i == count-1 {
			if err := config.DB.CreateInBatches(&users, batchSize).Error; err != nil {
				log.Printf("Error inserting batch: %v", err)
				continue
			}

			log.Printf("Inserted batch: %d/%d users", i+1, count)
			users = []models.User{} // Reset batch
		}
	}

	log.Printf("Successfully generated %d dummy users with registration dates from 2020-2025", count)

	// Show summary statistics
	showUserStatistics()
}

func showUserStatistics() {
	log.Println("\n=== User Statistics ===")

	var totalUsers int64
	config.DB.Model(&models.User{}).Count(&totalUsers)
	log.Printf("Total users in database: %d", totalUsers)

	// Count by year
	for year := 2020; year <= 2025; year++ {
		var yearCount int64
		startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

		config.DB.Model(&models.User{}).Where("created_at BETWEEN ? AND ?", startOfYear, endOfYear).Count(&yearCount)
		log.Printf("Users registered in %d: %d", year, yearCount)
	}

	// Count by status
	var activeUsers, inactiveUsers, blockedUsers int64
	config.DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_ACTIVE).Count(&activeUsers)
	config.DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_INACTIVE).Count(&inactiveUsers)
	config.DB.Model(&models.User{}).Where("status = ?", models.USER_STATUS_BLOCKED).Count(&blockedUsers)

	log.Printf("Active users: %d", activeUsers)
	log.Printf("Inactive users: %d", inactiveUsers)
	log.Printf("Blocked users: %d", blockedUsers)
}
