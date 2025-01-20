package data

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

var DB *sql.DB // Global variable to hold the database connection

// InitDB initializes the database connection
func InitDB(dsn string) {
	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to the database successfully!")
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}
