package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDatabase initializes the database connection
func InitDatabase() error {
	var err error
	
	// Database configuration from environment
	host := getEnvWithDefault("DB_HOST", "localhost")
	port := getEnvWithDefault("DB_PORT", "5435")
	user := getEnvWithDefault("DB_USER", "postgres")
	password := getEnvWithDefault("DB_PASSWORD", "postgres")
	dbname := getEnvWithDefault("DB_NAME", "location_service_db")
	sslmode := getEnvWithDefault("DB_SSLMODE", "disable")

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Connected to location database: %s:%s/%s", host, port, dbname)

	// Note: Old migrations system disabled in favor of new dual migration system
	// Use CLI commands: -migrate=schema and -migrate=data
	// if err = RunMigrations(db); err != nil {
	// 	return fmt.Errorf("failed to run migrations: %v", err)
	// }

	return nil
}

// getEnvWithDefault gets environment variable with default fallback
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if db != nil {
		return db.Close()
	}
	return nil
}