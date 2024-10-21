package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate // Extracting validator
var db *sql.DB

type User struct {
	Username string `validate:"required,min=3,max=255"`
	Password string `validate:"required,min=8"`
	Email    string `validate:"required,email"`
}
type Admin struct {
	Username string `validate:"required,min=3, max=255"`
	Password string `validate:"required,min=8"`
	Email    string `validate:"required,email"`
}
type Course struct {
	Title       string `validate:"required,min=10,max=255"`
	Description string `validate:"required,min=100,max=2500"`
	Image_link  string `validate:"omitempty,url"`
	AdminID     int    `validate:"required"`
}
type Purchase struct {
	UserID   int `validate:"required"`
	CourseID int `validate:"required"`
}

func main() {
	validate = validator.New() // Initializing new validator

	// Loading environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Successfully connected to the database")

	router := mux.NewRouter()
	userRouter := router.PathPrefix("/user").Subrouter() // user sub-router

	//    /user subroutes
	userRouter.HandleFunc("/signup", UserSignup).Methods("POST")
	userRouter.HandleFunc("/login", UserLogin).Methods("POST")
	userRouter.HandleFunc("/courses", SeeAllCourse).Methods("GET")
	userRouter.HandleFunc("/purchase", PurchaseCourse).Methods("POST")
	userRouter.HandleFunc("/purchases-courses", SeeAllPurchasesCourse).Methods("GET")

	adminRouter := router.PathPrefix("/admin").Subrouter() // admin sub-router

	//   /admin subroutes
	adminRouter.HandleFunc("/signup", AdminSignup).Methods("POST")
	adminRouter.HandleFunc("/login", AdminLogin).Methods("POST")
	adminRouter.HandleFunc("/create", CreateCourse).Methods("POST")
	adminRouter.HandleFunc("/delete", DeleteCourse).Methods("DELETE")
	adminRouter.HandleFunc("/add", AddCourseContent).Methods("POST")

	fmt.Println("Starting the server at port 8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Error while starting the server")
	}
}

// User Handlers
func UserSignup(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}
	err = validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			http.Error(w, fmt.Sprintf("Error happened: %v", err), http.StatusBadRequest)
		}
		return
	}

	// Hashing the user inputted password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Input validated user into the database
	query := `INSERT INTO users(username,email,password_hash) VALUES ($1,$2,$3)`
	_, err = db.Exec(query, user.Email, user.Username, hashedPassword)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
	}
	fmt.Fprintln(w, "User signed successfully")

}
func UserLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User Login")
}
func SeeAllCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Courses")
}
func PurchaseCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Buy from these")
}
func SeeAllPurchasesCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Your Purchased Courses")
}

// Admin Handlers
func AdminSignup(w http.ResponseWriter, r *http.Request) {
	var admin Admin
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}
	err = validate.Struct(admin)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			http.Error(w, fmt.Sprintf("Invalid input: %s", err), http.StatusBadRequest)
		}
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	query := `INSERT INTO admin(username,email,password_hash) VALUES($1,$2,$3)`
	_, err = db.Exec(query, admin.Username, admin.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User signed up successfully")
}
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Admin Login")
}
func CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "Error while decoding", http.StatusBadRequest)
	}
	err = validate.Struct(course)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			http.Error(w, fmt.Sprintf("Validation error:%v", err), http.StatusInternalServerError)
			return
		}
	}
	query := `INSERT INTO courses(title,description,image_link,admin_id) VALUES($1,$2,$3,$4)`
	_, err = db.Exec(query, course.Title, course.Description, course.Image_link, course.AdminID)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}
	fmt.Sprintln("Course created successfully")
}
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Course Deleted")
}
func AddCourseContent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Course Content Added")
}
