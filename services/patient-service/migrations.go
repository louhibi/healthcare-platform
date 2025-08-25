package main

import (
	"database/sql"
	"fmt"
	
	"github.com/louhibi/healthcare-logging"
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	Up          string
	Down        string
}

// GetMigrations returns all migrations in order
func GetMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Description: "Create patients table with multi-tenant support",
			Up: `
				CREATE TABLE IF NOT EXISTS patients (
					id SERIAL PRIMARY KEY,
					healthcare_entity_id INTEGER NOT NULL,
					patient_id VARCHAR(50), -- Custom patient identifier
					first_name VARCHAR(100) NOT NULL,
					last_name VARCHAR(100) NOT NULL,
					date_of_birth DATE NOT NULL,
					gender VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
					phone VARCHAR(20) NOT NULL,
					email VARCHAR(255) NOT NULL,
					address TEXT,
					city VARCHAR(100),
					state VARCHAR(100) -- Province/State/Region
					postal_code VARCHAR(20),
					country VARCHAR(2) NOT NULL,
					nationality VARCHAR(100),
					preferred_language VARCHAR(10) CHECK (preferred_language IS NULL OR preferred_language = '' OR preferred_language IN ('en', 'fr', 'ar')),
					marital_status VARCHAR(20) CHECK (marital_status IS NULL OR marital_status = '' OR marital_status IN ('single', 'married', 'divorced', 'widowed')),
					occupation VARCHAR(200),
					insurance VARCHAR(255),
					policy_number VARCHAR(100),
					insurance_provider VARCHAR(255),
					national_id VARCHAR(50), -- SSN, SIN, CIN, etc.
					emergency_contact_name VARCHAR(200),
					emergency_contact_phone VARCHAR(20),
					emergency_contact_relationship VARCHAR(50),
					medical_history TEXT,
					allergies TEXT,
					medications TEXT,
					blood_type VARCHAR(10) CHECK (blood_type IS NULL OR blood_type = '' OR blood_type IN ('A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-')),
					is_active BOOLEAN DEFAULT true,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					created_by INTEGER NOT NULL
				);

				-- Indexes for multi-tenant queries and performance
				CREATE INDEX IF NOT EXISTS idx_patients_entity ON patients(healthcare_entity_id);
				CREATE INDEX IF NOT EXISTS idx_patients_entity_active ON patients(healthcare_entity_id, is_active);
				CREATE INDEX IF NOT EXISTS idx_patients_name ON patients(first_name, last_name);
				CREATE INDEX IF NOT EXISTS idx_patients_email ON patients(email);
				CREATE INDEX IF NOT EXISTS idx_patients_phone ON patients(phone);
				CREATE INDEX IF NOT EXISTS idx_patients_dob ON patients(date_of_birth);
				CREATE INDEX IF NOT EXISTS idx_patients_country ON patients(country);
				CREATE INDEX IF NOT EXISTS idx_patients_patient_id ON patients(patient_id);
				CREATE INDEX IF NOT EXISTS idx_patients_national_id ON patients(national_id);
				CREATE INDEX IF NOT EXISTS idx_patients_created_by ON patients(created_by);
				CREATE INDEX IF NOT EXISTS idx_patients_blood_type ON patients(blood_type);
				
				-- Composite index for multi-tenant searches
				CREATE INDEX IF NOT EXISTS idx_patients_entity_search ON patients(healthcare_entity_id, first_name, last_name, email);
				
				-- Unique constraint for patient_id within entity (if patient_id is provided)
				CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_entity_patient_id ON patients(healthcare_entity_id, patient_id) WHERE patient_id IS NOT NULL AND patient_id <> '';

				-- Update updated_at trigger
				CREATE OR REPLACE FUNCTION update_updated_at_column()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.updated_at = CURRENT_TIMESTAMP;
					RETURN NEW;
				END;
				$$ language 'plpgsql';

				DROP TRIGGER IF EXISTS update_patients_updated_at ON patients;
				CREATE TRIGGER update_patients_updated_at
					BEFORE UPDATE ON patients
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();
			`,
			Down: `
				DROP TRIGGER IF EXISTS update_patients_updated_at ON patients;
				DROP FUNCTION IF EXISTS update_updated_at_column();
				DROP TABLE IF EXISTS patients CASCADE;
			`,
		},
		{
			Version:     2,
			Description: "Update database constraints to match form configuration values",
			Up: `
				-- First, update any existing data to match new format (before changing constraints)
				UPDATE patients SET gender = 'Male' WHERE gender = 'male';
				UPDATE patients SET gender = 'Female' WHERE gender = 'female';
				UPDATE patients SET gender = 'Other' WHERE gender = 'other';
				
				UPDATE patients SET marital_status = 'Single' WHERE marital_status = 'single';
				UPDATE patients SET marital_status = 'Married' WHERE marital_status = 'married';
				UPDATE patients SET marital_status = 'Divorced' WHERE marital_status = 'divorced';
				UPDATE patients SET marital_status = 'Widowed' WHERE marital_status = 'widowed';
				
				UPDATE patients SET preferred_language = 'English' WHERE preferred_language = 'en';
				UPDATE patients SET preferred_language = 'French' WHERE preferred_language = 'fr';
				UPDATE patients SET preferred_language = 'Arabic' WHERE preferred_language = 'ar';
				
				-- Now update constraints to match form configuration (Male, Female, Other)
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_gender_check;
				ALTER TABLE patients ADD CONSTRAINT patients_gender_check CHECK (gender IN ('Male', 'Female', 'Other'));
				
				-- Update marital_status constraint to match form configuration (Single, Married, Divorced, Widowed, Other)
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_marital_status_check;
				ALTER TABLE patients ADD CONSTRAINT patients_marital_status_check CHECK (marital_status IS NULL OR marital_status = '' OR marital_status IN ('Single', 'Married', 'Divorced', 'Widowed', 'Other'));
				
				-- Update preferred_language constraint to match form configuration (English, French, Arabic)
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_preferred_language_check;
				ALTER TABLE patients ADD CONSTRAINT patients_preferred_language_check CHECK (preferred_language IS NULL OR preferred_language = '' OR preferred_language IN ('English', 'French', 'Arabic'));
			`,
			Down: `
				-- Revert data to lowercase format
				UPDATE patients SET gender = 'male' WHERE gender = 'Male';
				UPDATE patients SET gender = 'female' WHERE gender = 'Female';
				UPDATE patients SET gender = 'other' WHERE gender = 'Other';
				
				UPDATE patients SET marital_status = 'single' WHERE marital_status = 'Single';
				UPDATE patients SET marital_status = 'married' WHERE marital_status = 'Married';
				UPDATE patients SET marital_status = 'divorced' WHERE marital_status = 'Divorced';
				UPDATE patients SET marital_status = 'widowed' WHERE marital_status = 'Widowed';
				
				UPDATE patients SET preferred_language = 'en' WHERE preferred_language = 'English';
				UPDATE patients SET preferred_language = 'fr' WHERE preferred_language = 'French';
				UPDATE patients SET preferred_language = 'ar' WHERE preferred_language = 'Arabic';
				
				-- Revert constraints to lowercase format
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_gender_check;
				ALTER TABLE patients ADD CONSTRAINT patients_gender_check CHECK (gender IN ('male', 'female', 'other'));
				
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_marital_status_check;
				ALTER TABLE patients ADD CONSTRAINT patients_marital_status_check CHECK (marital_status IS NULL OR marital_status = '' OR marital_status IN ('single', 'married', 'divorced', 'widowed'));
				
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_preferred_language_check;
				ALTER TABLE patients ADD CONSTRAINT patients_preferred_language_check CHECK (preferred_language IS NULL OR preferred_language = '' OR preferred_language IN ('en', 'fr', 'ar'));
			`,
		},
		{
			Version:     3,
			Description: "Fix existing data and update database constraints properly",
			Up: `
				-- First, update any existing data to match new format
				UPDATE patients SET gender = 'Male' WHERE LOWER(gender) = 'male';
				UPDATE patients SET gender = 'Female' WHERE LOWER(gender) = 'female';
				UPDATE patients SET gender = 'Other' WHERE LOWER(gender) = 'other';
				
				UPDATE patients SET marital_status = 'Single' WHERE LOWER(marital_status) = 'single';
				UPDATE patients SET marital_status = 'Married' WHERE LOWER(marital_status) = 'married';
				UPDATE patients SET marital_status = 'Divorced' WHERE LOWER(marital_status) = 'divorced';
				UPDATE patients SET marital_status = 'Widowed' WHERE LOWER(marital_status) = 'widowed';
				
				UPDATE patients SET preferred_language = 'English' WHERE preferred_language = 'en';
				UPDATE patients SET preferred_language = 'French' WHERE preferred_language = 'fr';
				UPDATE patients SET preferred_language = 'Arabic' WHERE preferred_language = 'ar';
				
				-- Drop old constraints
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_gender_check;
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_marital_status_check;
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_preferred_language_check;
				
				-- Add new constraints with proper values
				ALTER TABLE patients ADD CONSTRAINT patients_gender_check CHECK (gender IN ('Male', 'Female', 'Other'));
				ALTER TABLE patients ADD CONSTRAINT patients_marital_status_check CHECK (marital_status IS NULL OR marital_status = '' OR marital_status IN ('Single', 'Married', 'Divorced', 'Widowed', 'Other'));
				ALTER TABLE patients ADD CONSTRAINT patients_preferred_language_check CHECK (preferred_language IS NULL OR preferred_language = '' OR preferred_language IN ('English', 'French', 'Arabic'));
			`,
			Down: `
				-- Revert constraints to original format
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_gender_check;
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_marital_status_check;
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_preferred_language_check;
				
				-- Revert data to original format
				UPDATE patients SET gender = 'male' WHERE gender = 'Male';
				UPDATE patients SET gender = 'female' WHERE gender = 'Female';
				UPDATE patients SET gender = 'other' WHERE gender = 'Other';
				
				UPDATE patients SET marital_status = 'single' WHERE marital_status = 'Single';
				UPDATE patients SET marital_status = 'married' WHERE marital_status = 'Married';
				UPDATE patients SET marital_status = 'divorced' WHERE marital_status = 'Divorced';
				UPDATE patients SET marital_status = 'widowed' WHERE marital_status = 'Widowed';
				
				UPDATE patients SET preferred_language = 'en' WHERE preferred_language = 'English';
				UPDATE patients SET preferred_language = 'fr' WHERE preferred_language = 'French';
				UPDATE patients SET preferred_language = 'ar' WHERE preferred_language = 'Arabic';
				
				-- Restore original constraints
				ALTER TABLE patients ADD CONSTRAINT patients_gender_check CHECK (gender IN ('male', 'female', 'other'));
				ALTER TABLE patients ADD CONSTRAINT patients_marital_status_check CHECK (marital_status IS NULL OR marital_status = '' OR marital_status IN ('single', 'married', 'divorced', 'widowed'));
				ALTER TABLE patients ADD CONSTRAINT patients_preferred_language_check CHECK (preferred_language IS NULL OR preferred_language = '' OR preferred_language IN ('en', 'fr', 'ar'));
			`,
		},
		{
			Version:     4,
			Description: "Make email nullable and enforce uniqueness only for non-empty emails (case-insensitive)",
			Up: `
				-- Allow NULL emails (drop NOT NULL before setting any NULLs)
				ALTER TABLE patients ALTER COLUMN email DROP NOT NULL;

				-- Normalize existing blank emails to NULL
				UPDATE patients SET email = NULL WHERE email = '';

				-- Drop any previous unique constraints/indexes on email if they exist
				DO $$
				DECLARE r RECORD;
				BEGIN
					FOR r IN (
						SELECT indexname FROM pg_indexes 
						WHERE schemaname = 'public' AND tablename = 'patients' AND indexname LIKE '%email%unique%'
					) LOOP
						EXECUTE format('DROP INDEX IF EXISTS %I', r.indexname);
					END LOOP;
				END$$;

				-- Create partial unique index for provided emails only (case-insensitive)
				CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_email_unique
				ON patients (LOWER(email))
				WHERE email IS NOT NULL AND email <> '';
			`,
			Down: `
				-- Best-effort rollback: drop partial unique index and set NOT NULL back
				DROP INDEX IF EXISTS idx_patients_email_unique;
				UPDATE patients SET email = '' WHERE email IS NULL;
				ALTER TABLE patients ALTER COLUMN email SET NOT NULL;
			`,
		},
		{
			Version:     5,
			Description: "Enforce per-entity unique emails only when provided (case-insensitive)",
			Up: `
				-- Drop previous global email unique index if present
				DROP INDEX IF EXISTS idx_patients_email_unique;

				-- Create per-entity partial unique index for provided emails only
				CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_entity_email_unique
				ON patients (healthcare_entity_id, LOWER(email))
				WHERE email IS NOT NULL AND email <> '';
			`,
			Down: `
				-- Rollback to global unique index behavior
				DROP INDEX IF EXISTS idx_patients_entity_email_unique;
				CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_email_unique
				ON patients (LOWER(email))
				WHERE email IS NOT NULL AND email <> '';
			`,
		},
		{
			Version:     6,
			Description: "Ensure country is ISO 3166-1 alpha-2 uppercase",
			Up: `
				-- Map existing country names to proper ISO 3166-1 alpha-2 codes
				UPDATE patients SET country = 'MA' WHERE TRIM(UPPER(country)) = 'MOROCCO';
				UPDATE patients SET country = 'CA' WHERE TRIM(UPPER(country)) = 'CANADA';
				UPDATE patients SET country = 'FR' WHERE TRIM(UPPER(country)) = 'FRANCE';
				UPDATE patients SET country = 'US' WHERE TRIM(UPPER(country)) IN ('USA', 'US');
				
				-- Normalize any remaining values: trim and uppercase
				UPDATE patients SET country = UPPER(TRIM(country)) WHERE country IS NOT NULL;
				
				-- Add length check constraint (drop old if exists)
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_country_check;
				ALTER TABLE patients ADD CONSTRAINT patients_country_check CHECK (char_length(country) = 2);
			`,
			Down: `
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_country_check;
				-- Revert to original country names
				UPDATE patients SET country = 'Morocco' WHERE country = 'MA';
				UPDATE patients SET country = 'Canada' WHERE country = 'CA';
				UPDATE patients SET country = 'France' WHERE country = 'FR';
				UPDATE patients SET country = 'USA' WHERE country = 'US';
			`,
		},
		{
			Version:     7,
			Description: "Update patients table to use location service IDs",
			Up: `
				-- Add new location ID columns for integration with location service
				ALTER TABLE patients 
				ADD COLUMN country_id INTEGER,
				ADD COLUMN state_id INTEGER,
				ADD COLUMN city_id INTEGER;

				-- Add indexes for the new foreign key columns
				CREATE INDEX IF NOT EXISTS idx_patients_country_id ON patients(country_id);
				CREATE INDEX IF NOT EXISTS idx_patients_state_id ON patients(state_id);
				CREATE INDEX IF NOT EXISTS idx_patients_city_id ON patients(city_id);

				-- Note: String columns (country, state, city) kept for backward compatibility
				-- They will be populated by API Gateway from location service
				-- In a future migration, we can make country_id NOT NULL and drop string columns
			`,
			Down: `
				-- Remove the location ID columns
				DROP INDEX IF EXISTS idx_patients_country_id;
				DROP INDEX IF EXISTS idx_patients_state_id;
				DROP INDEX IF EXISTS idx_patients_city_id;
				
				ALTER TABLE patients 
				DROP COLUMN IF EXISTS country_id,
				DROP COLUMN IF EXISTS state_id,
				DROP COLUMN IF EXISTS city_id;
			`,
		},
		{
			Version:     8,
			Description: "Remove string location columns and use only location service IDs",
			Up: `
				-- Make country_id NOT NULL since it's required
				ALTER TABLE patients ALTER COLUMN country_id SET NOT NULL;

				-- Drop the old string-based location columns as we only use IDs now
				ALTER TABLE patients 
				DROP COLUMN IF EXISTS country,
				DROP COLUMN IF EXISTS state,
				DROP COLUMN IF EXISTS city;

				-- Drop related constraints and indexes that reference the string columns
				ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_country_check;
				DROP INDEX IF EXISTS idx_patients_country;
			`,
			Down: `
				-- Restore the string-based location columns for rollback
				ALTER TABLE patients 
				ADD COLUMN country VARCHAR(2),
				ADD COLUMN state VARCHAR(100),
				ADD COLUMN city VARCHAR(100);

				-- Restore constraints
				ALTER TABLE patients ADD CONSTRAINT patients_country_check CHECK (char_length(country) = 2);
				CREATE INDEX IF NOT EXISTS idx_patients_country ON patients(country);

				-- Make country_id nullable again
				ALTER TABLE patients ALTER COLUMN country_id DROP NOT NULL;
			`,
		},
		{
			Version:     9,
			Description: "Add nationality_id column to link to location service nationalities",
			Up: `
				-- Add nationality_id column to reference location service
				ALTER TABLE patients 
				ADD COLUMN nationality_id INTEGER;

				-- Add index for the new foreign key column
				CREATE INDEX IF NOT EXISTS idx_patients_nationality_id ON patients(nationality_id);

				-- Note: nationality string column kept temporarily for data migration
				-- Will be removed in a future migration after data conversion
			`,
			Down: `
				-- Remove the nationality_id column and index
				DROP INDEX IF EXISTS idx_patients_nationality_id;
				ALTER TABLE patients DROP COLUMN IF EXISTS nationality_id;
			`,
		},
		{
			Version:     10,
			Description: "Convert text nationality values to nationality_id references",
			Up: `
				-- Convert existing nationality text values to nationality_id references
				-- This migration will be handled by the application layer since we can't 
				-- cross-reference other microservice databases directly in SQL
				
				-- For now, just add a comment that conversion will happen via API calls
				-- The actual conversion will be handled by a data migration script or 
				-- during the first API call after the location service is populated
				
				-- Placeholder: nationality_id values will be populated via application logic
				-- that calls the location service API to resolve nationality text to IDs
				SELECT 1; -- No-op migration, handled at application level
			`,
			Down: `
				-- Revert nationality_id values back to NULL (text values remain in nationality column)
				UPDATE patients SET nationality_id = NULL;
			`,
		},
		{
			Version:     11,
			Description: "Drop old nationality text column after conversion to nationality_id",
			Up: `
				-- Remove the old text nationality column
				ALTER TABLE patients DROP COLUMN IF EXISTS nationality;
			`,
			Down: `
				-- Restore the nationality text column for rollback
				ALTER TABLE patients ADD COLUMN nationality VARCHAR(100);
				
				-- Restore text values from nationality_id (best effort)
				UPDATE patients 
				SET nationality = (
					SELECT n.name_en FROM location_service.nationalities n 
					WHERE n.id = patients.nationality_id
				)
				WHERE nationality_id IS NOT NULL;
			`,
		},
		{
			Version:     12,
			Description: "Replace insurance string fields with insurance ID references",
			Up: `
				-- Add new insurance ID columns to reference location service
				ALTER TABLE patients 
				ADD COLUMN insurance_type_id INTEGER,
				ADD COLUMN insurance_provider_id INTEGER;

				-- Add indexes for the new foreign key columns
				CREATE INDEX IF NOT EXISTS idx_patients_insurance_type_id ON patients(insurance_type_id);
				CREATE INDEX IF NOT EXISTS idx_patients_insurance_provider_id ON patients(insurance_provider_id);

				-- Remove old string-based insurance columns
				ALTER TABLE patients 
				DROP COLUMN IF EXISTS insurance,
				DROP COLUMN IF EXISTS insurance_provider;

				-- Note: policy_number column is kept as it remains a string field
			`,
			Down: `
				-- Restore the old string-based insurance columns
				ALTER TABLE patients 
				ADD COLUMN insurance VARCHAR(255),
				ADD COLUMN insurance_provider VARCHAR(255);

				-- Remove the insurance ID columns and indexes
				DROP INDEX IF EXISTS idx_patients_insurance_type_id;
				DROP INDEX IF EXISTS idx_patients_insurance_provider_id;
				
				ALTER TABLE patients 
				DROP COLUMN IF EXISTS insurance_type_id,
				DROP COLUMN IF EXISTS insurance_provider_id;
			`,
		},
	}
}

// CreateMigrationsTable creates the schema_migrations table to track applied migrations
func CreateMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

// GetAppliedMigrations returns a map of applied migration versions
func GetAppliedMigrations(db *sql.DB) (map[int]bool, error) {
	applied := make(map[int]bool)
	
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return applied, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return applied, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// RunMigrations applies all pending migrations
func RunMigrationsWithVersioning(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	if err := CreateMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	migrations := GetMigrations()
	
	for _, migration := range migrations {
		if applied[migration.Version] {
			logging.LogInfo("Migration already applied", "version", migration.Version, "description", migration.Description)
			continue
		}

		logging.LogInfo("Applying migration", "version", migration.Version, "description", migration.Description)
		
		// Begin transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", migration.Version, err)
		}

		// Execute migration
		if _, err := tx.Exec(migration.Up); err != nil {
			tx.Rollback()
			logging.LogError("Failed to execute migration", "version", migration.Version, "error", err)
			return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
		}

		// Record migration
		if _, err := tx.Exec("INSERT INTO schema_migrations (version, description) VALUES ($1, $2)", 
			migration.Version, migration.Description); err != nil {
			tx.Rollback()
			logging.LogError("Failed to record migration", "version", migration.Version, "error", err)
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			logging.LogError("Failed to commit migration", "version", migration.Version, "error", err)
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		logging.LogInfo("Successfully applied migration", "version", migration.Version)
	}

	return nil
}

// RollbackMigration rolls back a specific migration (for development use)
func RollbackMigration(db *sql.DB, version int) error {
	migrations := GetMigrations()
	
	var targetMigration *Migration
	for _, migration := range migrations {
		if migration.Version == version {
			targetMigration = &migration
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration %d not found", version)
	}

	logging.LogInfo("Rolling back migration", "version", version, "description", targetMigration.Description)
	
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for rollback %d: %w", version, err)
	}

	// Execute rollback
	if _, err := tx.Exec(targetMigration.Down); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute rollback %d: %w", version, err)
	}

	// Remove migration record
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %d: %w", version, err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback %d: %w", version, err)
	}

	logging.LogInfo("Successfully rolled back migration", "version", version)
	return nil
}