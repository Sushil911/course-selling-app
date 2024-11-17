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
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables %v", err)
	}

	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	}
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.RegisterRoutes(e)

	fmt.Println("Starting the server at port 8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Error starting the server %v", err)
	}
}
