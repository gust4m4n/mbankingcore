package main

import (
	"log"

	"mbankingcore/config"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("MBankingCore - Database Migration Tool")
	log.Println("=======================================")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	// Connect to database
	config.ConnectDatabase()

	log.Println("Migration completed successfully!")
	log.Println("")
	log.Println("Your database is now ready for the MBX Backend application.")
	log.Println("You can now start the main application with: go run main.go")
}
