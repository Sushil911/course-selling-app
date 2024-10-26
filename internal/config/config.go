package config

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var JWTSecret []byte

// LoadEnv loads environment variables
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	return nil
}

// GenerateJWT generates a JWT token
func GenerateJWT(username, role string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}
