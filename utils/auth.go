package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// getJWTSecret returns JWT secret from environment variable or default
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-jwt-secret-change-this-in-production"
	}
	return []byte(secret)
}

// JWT Claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

// HashPassword applies bcrypt to the SHA256 hash sent by client
func HashPassword(sha256Hash string) (string, error) {
	// Client sends SHA256 hash, we apply bcrypt on top for additional security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(sha256Hash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares bcrypt hashed password with client-provided SHA256 hash
func CheckPassword(storedBcryptHash, clientSHA256Hash string) error {
	// Compare stored bcrypt hash with client's SHA256 hash
	return bcrypt.CompareHashAndPassword([]byte(storedBcryptHash), []byte(clientSHA256Hash))
}

// GenerateJWT generates a JWT token for the user
func GenerateJWT(userID uint, phone string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		UserID: userID,
		Phone:  phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates and parses a JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GenerateOTP generates a 6-digit OTP code
func GenerateOTP() string {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback to timestamp-based OTP if crypto/rand fails
		return fmt.Sprintf("%06d", time.Now().Unix()%1000000)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

// GenerateLoginToken generates a unique login token for OTP sessions
func GenerateLoginToken() string {
	// Generate 32 bytes of random data
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to timestamp-based token if crypto/rand fails
		return fmt.Sprintf("login_token_%d_%d", time.Now().UnixNano(), time.Now().Unix())
	}

	// Convert to hex string
	return fmt.Sprintf("%x", bytes)
}
