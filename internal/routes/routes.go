package routes

import (
	"course-selling-app/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(c *echo.Echo) {
	// user handlers
	userGroup := c.Group("/user")
	userGroup.POST("/signup", handlers.HandleUserSignup)
	userGroup.POST("/login", handlers.HandleUserLogin)
	userGroup.GET("/courses", handlers.HandleSeeAllCourses)
	userGroup.POST("/purchase-courses", handlers.HandlePurchaseCourses)
	userGroup.POST("/purchased-courses", handlers.HandleSeeAllPurchasedCourses)

	adminGroup := c.Group("/admin")
	adminGroup.POST("/signup", handlers.HandleAdminSignup)
	adminGroup.POST("/login", handlers.HandleAdminLogin)
	adminGroup.POST("/create-course", handlers.HandleCreateCourse)
	adminGroup.DELETE("/delete-course", handlers.HandleDeleteCourse)
	adminGroup.POST("/add-course-content", handlers.HandleAddCourseContent)
}
