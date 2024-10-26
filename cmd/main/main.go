package main

import (
	"course-selling-app/internal/config"
	"course-selling-app/internal/db"
	"course-selling-app/internal/routes"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Initialize the database
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// Create Echo instance
	e := echo.New()

	// Default Middleware in echo for logging and recovering from any panic
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register routes
	routes.RegisterRoutes(e)

	// Start the server
	fmt.Println("Starting the server at port 8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Error while starting the server: %v", err)
	}
}
