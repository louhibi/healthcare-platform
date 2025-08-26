package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// InitDB initializes database connection for config service
func InitDB() (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "postgres")
	name := getEnv("DB_NAME", "config_service_db")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// RunMigrations executes versioned migrations
func RunMigrations(db *sql.DB) error { return RunMigrationsWithVersioning(db) }
