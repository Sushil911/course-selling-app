package handlers

import (
	"course-selling-app/internal/config"
	"course-selling-app/internal/db"
	"course-selling-app/internal/models"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

// user handlers
func HandleUserSignup(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid Input"})
	}
	if err := validate.Struct(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation failed"})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error hashing password"})
	}
	query := `INSERT INTO users(username,password_hash,email) VALUES($1,$2,$3)`
	_, err = db.DB.Exec(query, user.Username, hashedPassword, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error connecting to database"})
	}
	token, err := config.GenerateJWT(user.Username, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "error generating JWT"})
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})

}

func HandleUserLogin(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "UserLogin"})
}
func HandleSeeAllCourses(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "Courses"})
}
func HandlePurchaseCourses(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "PurchaseCourses"})
}
func HandleSeeAllPurchasedCourses(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "SeeAllPurchasedCourses"})
}

//admin handlers

func HandleAdminSignup(c echo.Context) error {
	var admin models.Admin
	if err := c.Bind(&admin); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid Input"})
	}
	if err := validate.Struct(&admin); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation failed"})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error hashing password"})
	}
	query := `INSERT INTO admin(username,password_hash,email) VALUES($1,$2,$3)`
	_, err = db.DB.Exec(query, admin.Username, hashedPassword, admin.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error creating admin"})
	}
	token, err := config.GenerateJWT(admin.Username, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error generating JWT"})
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func HandleAdminLogin(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "admin loggedin successfully"})
}
func HandleCreateCourse(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "Course created successfully"})
}
func HandleDeleteCourse(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "Course deleted successfully"})
}
func HandleAddCourseContent(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "Course content added successfully"})
}
