package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// Load JWT secret from .env file
var jwtSecret []byte

// Structs for validation and binding
type User struct {
	Username string `json:"username" binding:"required,min=3,max=255"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

type Admin struct {
	Username string `json:"username" binding:"required,min=3,max=255"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

type Course struct {
	Title       string `json:"title" binding:"required,min=10,max=255"`
	Description string `json:"description" binding:"required,min=100,max=2500"`
	ImageLink   string `json:"image_link" binding:"omitempty,url"`
	AdminID     int    `json:"admin_id" binding:"required"`
}

type Purchase struct {
	UserID   int `json:"user_id" binding:"required"`
	CourseID int `json:"course_id" binding:"required"`
}

// Main function
func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	// User routes
	userGroup := r.Group("/user")
	{
		userGroup.POST("/signup", UserSignup)
		userGroup.POST("/login", UserLogin)
		userGroup.GET("/courses", SeeAllCourse)
		userGroup.POST("/purchase", PurchaseCourse)
		userGroup.GET("/purchases-courses", SeeAllPurchasesCourse)
	}

	// Admin routes
	adminGroup := r.Group("/admin")
	{
		adminGroup.POST("/signup", AdminSignup)
		adminGroup.POST("/login", AdminLogin)
		adminGroup.POST("/create", CreateCourse)
		adminGroup.DELETE("/delete", DeleteCourse)
		adminGroup.POST("/add", AddCourseContent)
	}

	// Protecting admin and user routes
	adminGroup.Use(verifyJWT("admin"))
	userGroup.Use(verifyJWT("user"))

	fmt.Println("Starting the server at port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error while starting the server: %v", err)
	}
}

// Load environment variables
func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	return nil
}

// Initialize database connection
func initDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the database")
	return db, nil
}

// JWT Verification Middleware
func verifyJWT(expectedRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			role := claims["role"].(string)
			if role != expectedRole {
				c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// User Handlers
func UserSignup(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	query := `INSERT INTO users(username, email, password_hash) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, user.Username, user.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user"})
		return
	}

	token, err := GenerateJWT(user.Username, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User signed up successfully", "token": token})
}

func UserLogin(c *gin.Context) {
	// Your login logic here
	c.JSON(http.StatusOK, gin.H{"message": "User Login"})
}

func SeeAllCourse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Courses"})
}

func PurchaseCourse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Buy from these"})
}

func SeeAllPurchasesCourse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Your Purchased Courses"})
}

// Admin Handlers
func AdminSignup(c *gin.Context) {
	var admin Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	query := `INSERT INTO admin(username, email, password_hash) VALUES($1, $2, $3)`
	_, err = db.Exec(query, admin.Username, admin.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting admin"})
		return
	}

	token, err := GenerateJWT(admin.Username, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin signed up successfully", "token": token})
}

func AdminLogin(c *gin.Context) {
	// Your login logic here
	c.JSON(http.StatusOK, gin.H{"message": "Admin Login"})
}

func CreateCourse(c *gin.Context) {
	var course Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := `INSERT INTO courses(title, description, image_link, admin_id) VALUES($1, $2, $3, $4)`
	_, err := db.Exec(query, course.Title, course.Description, course.ImageLink, course.AdminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course created successfully"})
}

func DeleteCourse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete Course"})
}

func AddCourseContent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Add Content to Course"})
}

// JWT Generation function
func GenerateJWT(username, role string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
