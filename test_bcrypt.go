package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Password yang kita coba
	password := "Super123?"

	// Hash yang tersimpan di database
	storedHash := "$2a$10$YJ7g.KvNjM9sC4QKiXhF3.sF4cVDT8t8O7T5wJ5QKXXr1eLM9WbkS"

	fmt.Printf("Testing password: %s\n", password)
	fmt.Printf("Stored hash: %s\n", storedHash)

	// Test bcrypt comparison
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		fmt.Printf("❌ Password verification failed: %v\n", err)
	} else {
		fmt.Printf("✅ Password verification successful!\n")
	}

	// Generate new hash for testing
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("❌ Hash generation failed: %v\n", err)
	} else {
		fmt.Printf("New hash for comparison: %s\n", string(newHash))

		// Test with new hash
		err2 := bcrypt.CompareHashAndPassword(newHash, []byte(password))
		if err2 != nil {
			fmt.Printf("❌ New hash verification failed: %v\n", err2)
		} else {
			fmt.Printf("✅ New hash verification successful!\n")
		}
	}
}
