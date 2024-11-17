package config

import (
	"course-selling-app/internal/db"
	"course-selling-app/internal/models"
	"database/sql"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var JWTSecret []byte

// LoadEnv loads environment variables
func LoadEnv() error {
	if err := godotenv.Load("/mnt/c/Users/Acer/Desktop/course-selling-app/.env"); err != nil {
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

// function to check if the email already exists or not during login
func EmailExists(email string) (bool, *models.DatabaseInfo, error) {
	query := `SELECT email,password_hash,username FROM users WHERE email=$1`
	var user models.DatabaseInfo
	err := db.DB.QueryRow(query, email).Scan(&user.Email, &user.Password, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, nil
		}
		return false, nil, nil
	}
	return true, &user, nil
}
