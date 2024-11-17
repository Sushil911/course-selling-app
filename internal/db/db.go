package db

import (
	"database/sql"
	"fmt"
	"os"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = DB.Ping(); err != nil {
		return nil, err
	}
	return DB, nil
}
