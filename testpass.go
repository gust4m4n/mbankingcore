package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	hash := "$2a$10$hLHeEMx.qZ/v/Jp4kAE5zO0/f8RdKZ2HctCy9VTncZYmtZ/U3t8/G"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Password does not match: %v\n", err)
	} else {
		fmt.Println("Password matches!")
	}
}
