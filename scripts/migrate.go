package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	var (
		service = flag.String("service", "", "Service name (user, patient, appointment)")
		action  = flag.String("action", "up", "Migration action (up, down, status)")
		version = flag.Int("version", 0, "Migration version (for down action)")
		dbHost  = flag.String("host", "localhost", "Database host")
		dbPort  = flag.String("port", "5432", "Database port") 
		dbUser  = flag.String("user", "postgres", "Database user")
		dbPass  = flag.String("password", "postgres", "Database password")
	)
	flag.Parse()

	if *service == "" {
		log.Fatal("Service name is required. Use: user, patient, or appointment")
	}

	// Determine database name based on service
	var dbName string
	switch *service {
	case "user":
		dbName = "user_service_db"
	case "patient":
		dbName = "patient_service_db"
	case "appointment":
		dbName = "appointment_service_db"
	default:
		log.Fatal("Invalid service name. Use: user, patient, or appointment")
	}

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		*dbHost, *dbPort, *dbUser, *dbPass, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Execute action
	switch *action {
	case "up":
		if err := runMigrations(db, *service); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")

	case "down":
		if *version == 0 {
			log.Fatal("Version is required for down action")
		}
		if err := rollbackMigration(db, *service, *version); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		log.Printf("Migration %d rolled back successfully", *version)

	case "status":
		if err := showMigrationStatus(db); err != nil {
			log.Fatalf("Failed to show migration status: %v", err)
		}

	default:
		log.Fatal("Invalid action. Use: up, down, or status")
	}
}

func runMigrations(db *sql.DB, service string) error {
	// This is a simplified version - in a real implementation,
	// you would import the actual migration functions from each service
	switch service {
	case "user":
		return runUserMigrations(db)
	case "patient":
		return runPatientMigrations(db)
	case "appointment":
		return runAppointmentMigrations(db)
	default:
		return fmt.Errorf("unknown service: %s", service)
	}
}

func rollbackMigration(db *sql.DB, service string, version int) error {
	// This is a simplified version - in a real implementation,
	// you would import the actual rollback functions from each service
	switch service {
	case "user":
		return rollbackUserMigration(db, version)
	case "patient":
		return rollbackPatientMigration(db, version)
	case "appointment":
		return rollbackAppointmentMigration(db, version)
	default:
		return fmt.Errorf("unknown service: %s", service)
	}
}

func showMigrationStatus(db *sql.DB) error {
	rows, err := db.Query("SELECT version, description, applied_at FROM schema_migrations ORDER BY version")
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("Applied Migrations:")
	fmt.Println("Version | Description | Applied At")
	fmt.Println("--------|-------------|------------")

	for rows.Next() {
		var version int
		var description, appliedAt string
		if err := rows.Scan(&version, &description, &appliedAt); err != nil {
			return err
		}
		fmt.Printf("%-7d | %-50s | %s\n", version, description, appliedAt)
	}

	return rows.Err()
}

// Placeholder functions - in a real implementation, these would be imported
func runUserMigrations(db *sql.DB) error {
	log.Println("Running user service migrations...")
	// This would call the actual migration function from user service
	return nil
}

func runPatientMigrations(db *sql.DB) error {
	log.Println("Running patient service migrations...")
	// This would call the actual migration function from patient service
	return nil
}

func runAppointmentMigrations(db *sql.DB) error {
	log.Println("Running appointment service migrations...")
	// This would call the actual migration function from appointment service
	return nil
}

func rollbackUserMigration(db *sql.DB, version int) error {
	log.Printf("Rolling back user service migration %d...", version)
	// This would call the actual rollback function from user service
	return nil
}

func rollbackPatientMigration(db *sql.DB, version int) error {
	log.Printf("Rolling back patient service migration %d...", version)
	// This would call the actual rollback function from patient service
	return nil
}

func rollbackAppointmentMigration(db *sql.DB, version int) error {
	log.Printf("Rolling back appointment service migration %d...", version)
	// This would call the actual rollback function from appointment service
	return nil
}