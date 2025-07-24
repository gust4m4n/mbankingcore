package utils

import (
	"errors"
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
	Email  string `json:"email"`
	Role   string `json:"role"`
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
func GenerateJWT(userID uint, email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
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
