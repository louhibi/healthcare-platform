package main

import (
	"database/sql"
	"fmt"
	"log"
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
			Description: "Create healthcare_entities and users tables with multi-tenant support",
			Up: `
				-- Healthcare entities table (hospitals, clinics, doctor offices)
				CREATE TABLE IF NOT EXISTS healthcare_entities (
					id SERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					type VARCHAR(50) NOT NULL CHECK (type IN ('hospital', 'clinic', 'doctor_office')),
					country VARCHAR(50) NOT NULL CHECK (country IN ('Canada', 'USA', 'Morocco', 'France')),
					address VARCHAR(500) NOT NULL,
					city VARCHAR(100) NOT NULL,
					state VARCHAR(100),
					postal_code VARCHAR(20) NOT NULL,
					phone VARCHAR(50) NOT NULL,
					email VARCHAR(255) NOT NULL,
					website VARCHAR(255),
					license VARCHAR(100),
					tax_id VARCHAR(100),
					timezone VARCHAR(50) NOT NULL,
					language VARCHAR(10) NOT NULL CHECK (language IN ('en', 'fr', 'ar')),
					currency VARCHAR(10) NOT NULL CHECK (currency IN ('CAD', 'USD', 'MAD', 'EUR')),
					is_active BOOLEAN DEFAULT true,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Users table with entity association
				CREATE TABLE IF NOT EXISTS users (
					id SERIAL PRIMARY KEY,
					email VARCHAR(255) UNIQUE NOT NULL,
					password_hash VARCHAR(255) NOT NULL,
					first_name VARCHAR(100) NOT NULL,
					last_name VARCHAR(100) NOT NULL,
					role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'doctor', 'nurse', 'staff')),
					healthcare_entity_id INTEGER NOT NULL REFERENCES healthcare_entities(id),
					license_number VARCHAR(100),
					specialization VARCHAR(200),
					preferred_language VARCHAR(10) CHECK (preferred_language IN ('en', 'fr', 'ar')),
					is_active BOOLEAN DEFAULT true,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Indexes
				CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
				CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
				CREATE INDEX IF NOT EXISTS idx_users_entity ON users(healthcare_entity_id);
				CREATE INDEX IF NOT EXISTS idx_healthcare_entities_country ON healthcare_entities(country);
			`,
			Down: `
				DROP INDEX IF EXISTS idx_healthcare_entities_country;
				DROP INDEX IF EXISTS idx_users_entity;
				DROP INDEX IF EXISTS idx_users_role;
				DROP INDEX IF EXISTS idx_users_email;
				DROP TABLE IF EXISTS users;
				DROP TABLE IF EXISTS healthcare_entities;
			`,
		},
		{
			Version:     2,
			Description: "Create form configuration tables for dynamic forms",
			Up: `
				-- Form types (patient, appointment, etc.)
				CREATE TABLE IF NOT EXISTS form_types (
					id SERIAL PRIMARY KEY,
					name VARCHAR(50) UNIQUE NOT NULL,
					display_name VARCHAR(100) NOT NULL,
					description TEXT,
					is_active BOOLEAN DEFAULT true,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Field definitions (shared across all entities)
				CREATE TABLE IF NOT EXISTS field_definitions (
					id SERIAL PRIMARY KEY,
					form_type_id INTEGER NOT NULL REFERENCES form_types(id),
					name VARCHAR(100) NOT NULL, -- field name like 'first_name'
					display_name VARCHAR(200) NOT NULL,
					field_type VARCHAR(50) NOT NULL CHECK (field_type IN ('text', 'email', 'phone', 'number', 'date', 'datetime', 'select', 'textarea', 'checkbox')),
					default_required BOOLEAN DEFAULT false,
					is_core BOOLEAN DEFAULT false, -- Core fields cannot be disabled
					validation_rules JSONB DEFAULT '{}',
					options TEXT[] DEFAULT '{}', -- For select fields
					sort_order INTEGER DEFAULT 0,
					category VARCHAR(100), -- Personal Information, Contact Information, etc.
					description TEXT,
					placeholder_text VARCHAR(200),
					is_active BOOLEAN DEFAULT true,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(form_type_id, name)
				);

				-- Per-entity field configuration (customizations)
				CREATE TABLE IF NOT EXISTS entity_field_configurations (
					id SERIAL PRIMARY KEY,
					healthcare_entity_id INTEGER NOT NULL REFERENCES healthcare_entities(id),
					field_definition_id INTEGER NOT NULL REFERENCES field_definitions(id),
					is_enabled BOOLEAN DEFAULT true,
					is_required BOOLEAN DEFAULT false,
					custom_label VARCHAR(200), -- Override display name
					custom_validation JSONB DEFAULT '{}', -- Additional validation rules
					sort_order INTEGER DEFAULT 0, -- Override sort order
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(healthcare_entity_id, field_definition_id)
				);

				-- Add updated_at triggers
				CREATE OR REPLACE FUNCTION update_updated_at_column()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.updated_at = CURRENT_TIMESTAMP;
					RETURN NEW;
				END;
				$$ language 'plpgsql';

				DROP TRIGGER IF EXISTS update_healthcare_entities_updated_at ON healthcare_entities;
				CREATE TRIGGER update_healthcare_entities_updated_at
					BEFORE UPDATE ON healthcare_entities
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();

				DROP TRIGGER IF EXISTS update_users_updated_at ON users;
				CREATE TRIGGER update_users_updated_at
					BEFORE UPDATE ON users
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();

				DROP TRIGGER IF EXISTS update_form_types_updated_at ON form_types;
				CREATE TRIGGER update_form_types_updated_at
					BEFORE UPDATE ON form_types
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();

				DROP TRIGGER IF EXISTS update_field_definitions_updated_at ON field_definitions;
				CREATE TRIGGER update_field_definitions_updated_at
					BEFORE UPDATE ON field_definitions
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();

				DROP TRIGGER IF EXISTS update_entity_field_configurations_updated_at ON entity_field_configurations;
				CREATE TRIGGER update_entity_field_configurations_updated_at
					BEFORE UPDATE ON entity_field_configurations
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();

				-- Insert default form types
				INSERT INTO form_types (name, display_name, description) VALUES
				('patient', 'Patient Registration', 'Patient registration and profile management form'),
				('appointment', 'Appointment Booking', 'Appointment scheduling and management form')
				ON CONFLICT (name) DO NOTHING;

				-- Insert patient field definitions
				INSERT INTO field_definitions (form_type_id, name, display_name, field_type, default_required, is_core, category, description, placeholder_text, sort_order) VALUES
				((SELECT id FROM form_types WHERE name = 'patient'), 'first_name', 'First Name', 'text', true, true, 'Personal Information', 'Patient first name', 'Enter first name', 1),
				((SELECT id FROM form_types WHERE name = 'patient'), 'last_name', 'Last Name', 'text', true, true, 'Personal Information', 'Patient last name', 'Enter last name', 2),
				((SELECT id FROM form_types WHERE name = 'patient'), 'date_of_birth', 'Date of Birth', 'date', true, true, 'Personal Information', 'Patient date of birth', 'YYYY-MM-DD', 3),
				((SELECT id FROM form_types WHERE name = 'patient'), 'gender', 'Gender', 'select', true, false, 'Personal Information', 'Patient gender', '', 4),
				((SELECT id FROM form_types WHERE name = 'patient'), 'email', 'Email Address', 'email', false, false, 'Contact Information', 'Patient email address', 'patient@example.com', 5),
				((SELECT id FROM form_types WHERE name = 'patient'), 'phone', 'Phone Number', 'phone', true, false, 'Contact Information', 'Patient phone number', '+1 (xxx) xxx-xxxx', 6),
				((SELECT id FROM form_types WHERE name = 'patient'), 'address', 'Address', 'textarea', true, false, 'Contact Information', 'Patient address', 'Street address', 7),
				((SELECT id FROM form_types WHERE name = 'patient'), 'country', 'Country', 'select', true, false, 'Contact Information', 'Patient country', '', 8),
				((SELECT id FROM form_types WHERE name = 'patient'), 'state', 'State/Province', 'select', false, false, 'Contact Information', 'Patient state or province', 'State/Province', 9),
				((SELECT id FROM form_types WHERE name = 'patient'), 'city', 'City', 'select', true, false, 'Contact Information', 'Patient city', 'Enter city', 10),
				((SELECT id FROM form_types WHERE name = 'patient'), 'postal_code', 'Postal Code', 'text', true, false, 'Contact Information', 'Patient postal code', 'Postal/ZIP code', 11),
				((SELECT id FROM form_types WHERE name = 'patient'), 'nationality', 'Nationality', 'text', false, false, 'Personal Information', 'Patient nationality', 'Enter nationality', 12),
				((SELECT id FROM form_types WHERE name = 'patient'), 'preferred_language', 'Preferred Language', 'select', false, false, 'Personal Information', 'Patient preferred language', '', 13),
				((SELECT id FROM form_types WHERE name = 'patient'), 'marital_status', 'Marital Status', 'select', false, false, 'Personal Information', 'Patient marital status', '', 14),
				((SELECT id FROM form_types WHERE name = 'patient'), 'occupation', 'Occupation', 'text', false, false, 'Personal Information', 'Patient occupation', 'Enter occupation', 15),
				((SELECT id FROM form_types WHERE name = 'patient'), 'insurance', 'Insurance Type', 'select', false, false, 'Insurance Information', 'Patient insurance type', '', 16),
				((SELECT id FROM form_types WHERE name = 'patient'), 'policy_number', 'Policy Number', 'text', false, false, 'Insurance Information', 'Insurance policy number', 'Enter policy number', 17),
				((SELECT id FROM form_types WHERE name = 'patient'), 'insurance_provider', 'Insurance Provider', 'text', false, false, 'Insurance Information', 'Insurance provider name', 'Enter provider', 18),
				((SELECT id FROM form_types WHERE name = 'patient'), 'national_id', 'National ID', 'text', false, false, 'Personal Information', 'National identification number', 'Enter national ID', 19),
				((SELECT id FROM form_types WHERE name = 'patient'), 'emergency_contact_name', 'Emergency Contact Name', 'text', false, false, 'Emergency Contact', 'Emergency contact full name', 'Enter contact name', 20),
				((SELECT id FROM form_types WHERE name = 'patient'), 'emergency_contact_relationship', 'Emergency Contact Relationship', 'text', false, false, 'Emergency Contact', 'Relationship to patient', 'Enter relationship', 21),
				((SELECT id FROM form_types WHERE name = 'patient'), 'emergency_contact_phone', 'Emergency Contact Phone', 'phone', false, false, 'Emergency Contact', 'Emergency contact phone number', '+1 (xxx) xxx-xxxx', 22),
				((SELECT id FROM form_types WHERE name = 'patient'), 'medical_history', 'Medical History', 'textarea', false, false, 'Medical Information', 'Patient medical history', 'Enter medical history', 23),
				((SELECT id FROM form_types WHERE name = 'patient'), 'allergies', 'Allergies', 'textarea', false, false, 'Medical Information', 'Patient allergies', 'Enter allergies', 24),
				((SELECT id FROM form_types WHERE name = 'patient'), 'medications', 'Current Medications', 'textarea', false, false, 'Medical Information', 'Current medications', 'Enter current medications', 25),
				((SELECT id FROM form_types WHERE name = 'patient'), 'blood_type', 'Blood Type', 'select', false, false, 'Medical Information', 'Patient blood type', '', 26)
				ON CONFLICT (form_type_id, name) DO NOTHING;

				-- Insert appointment field definitions
				INSERT INTO field_definitions (form_type_id, name, display_name, field_type, default_required, is_core, category, description, placeholder_text, sort_order) VALUES
				((SELECT id FROM form_types WHERE name = 'appointment'), 'appointment_date', 'Appointment Date', 'date', true, true, 'Appointment Details', 'Date of appointment', 'YYYY-MM-DD', 1),
				((SELECT id FROM form_types WHERE name = 'appointment'), 'appointment_time', 'Appointment Time', 'datetime', true, true, 'Appointment Details', 'Time of appointment', 'HH:MM', 2),
				((SELECT id FROM form_types WHERE name = 'appointment'), 'doctor', 'Doctor', 'select', true, true, 'Appointment Details', 'Attending physician', '', 3),
				((SELECT id FROM form_types WHERE name = 'appointment'), 'appointment_type', 'Appointment Type', 'select', true, false, 'Appointment Details', 'Type of appointment', '', 4),
				((SELECT id FROM form_types WHERE name = 'appointment'), 'reason', 'Reason for Visit', 'textarea', false, false, 'Appointment Details', 'Reason for the appointment', 'Enter reason for visit', 5),
				((SELECT id FROM form_types WHERE name = 'appointment'), 'notes', 'Additional Notes', 'textarea', false, false, 'Appointment Details', 'Any additional notes', 'Enter additional notes', 6)
				ON CONFLICT (form_type_id, name) DO NOTHING;

				-- Create indexes
				CREATE INDEX IF NOT EXISTS idx_field_definitions_form_type ON field_definitions(form_type_id);
				CREATE INDEX IF NOT EXISTS idx_entity_field_configurations_entity ON entity_field_configurations(healthcare_entity_id);
				CREATE INDEX IF NOT EXISTS idx_entity_field_configurations_field ON entity_field_configurations(field_definition_id);
			`,
			Down: `
				DROP INDEX IF EXISTS idx_entity_field_configurations_field;
				DROP INDEX IF EXISTS idx_entity_field_configurations_entity;
				DROP INDEX IF EXISTS idx_field_definitions_form_type;
				
				DROP TABLE IF EXISTS entity_field_configurations;
				DROP TABLE IF EXISTS field_definitions;
				DROP TABLE IF EXISTS form_types;
				
				DROP FUNCTION IF EXISTS update_updated_at_column();
			`,
		},
		{
			Version:     3,
			Description: "Add internationalization support with locales and translations",
			Up: `
				-- Add healthcare entity room assignment requirement
				ALTER TABLE healthcare_entities 
				ADD COLUMN require_room_assignment BOOLEAN DEFAULT false;
				
				-- Locales table for supported languages/countries
				CREATE TABLE IF NOT EXISTS locales (
					code VARCHAR(10) PRIMARY KEY, -- en-US, en-CA, fr-CA, fr-FR, ar-MA
					language_name VARCHAR(100) NOT NULL, -- English, French, Arabic
					native_name VARCHAR(100) NOT NULL, -- English, Français, العربية
					country_code VARCHAR(2), -- US, CA, FR, MA
					is_active BOOLEAN DEFAULT true,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Insert supported locales for healthcare platform countries
				INSERT INTO locales (code, language_name, native_name, country_code, is_active) VALUES
				('en-US', 'English (United States)', 'English (United States)', 'US', true),
				('en-CA', 'English (Canada)', 'English (Canada)', 'CA', true),
				('fr-CA', 'French (Canada)', 'Français (Canada)', 'CA', true),
				('fr-FR', 'French (France)', 'Français (France)', 'FR', true),
				('ar-MA', 'Arabic (Morocco)', 'العربية (المغرب)', 'MA', true)
				ON CONFLICT (code) DO NOTHING;
				
				-- Create translations table for dynamic content
				CREATE TABLE IF NOT EXISTS translations (
					id SERIAL PRIMARY KEY,
					translation_key VARCHAR(255) NOT NULL,
					locale VARCHAR(10) NOT NULL REFERENCES locales(code),
					content TEXT NOT NULL,
					context VARCHAR(100), -- form_field, validation_message, ui_text, etc.
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(translation_key, locale)
				);
				
				-- Add translations for form field display names
				CREATE TABLE IF NOT EXISTS field_translations (
					id SERIAL PRIMARY KEY,
					field_definition_id INTEGER NOT NULL REFERENCES field_definitions(id) ON DELETE CASCADE,
					locale VARCHAR(10) NOT NULL REFERENCES locales(code),
					display_name VARCHAR(200) NOT NULL,
					description TEXT,
					placeholder_text VARCHAR(200),
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(field_definition_id, locale)
				);

				-- Update triggers for new tables
				DROP TRIGGER IF EXISTS update_translations_updated_at ON translations;
				CREATE TRIGGER update_translations_updated_at
					BEFORE UPDATE ON translations
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();

				DROP TRIGGER IF EXISTS update_field_translations_updated_at ON field_translations;
				CREATE TRIGGER update_field_translations_updated_at
					BEFORE UPDATE ON field_translations
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();

				-- Create indexes for performance
				CREATE INDEX IF NOT EXISTS idx_translations_key ON translations(translation_key);
				CREATE INDEX IF NOT EXISTS idx_translations_locale ON translations(locale);
				CREATE INDEX IF NOT EXISTS idx_field_translations_field ON field_translations(field_definition_id);
				CREATE INDEX IF NOT EXISTS idx_field_translations_locale ON field_translations(locale);
			`,
			Down: `
				DROP INDEX IF EXISTS idx_field_translations_locale;
				DROP INDEX IF EXISTS idx_field_translations_field;
				DROP INDEX IF EXISTS idx_translations_locale;
				DROP INDEX IF EXISTS idx_translations_key;
				
				DROP TABLE IF EXISTS field_translations;
				DROP TABLE IF EXISTS translations;
				DROP TABLE IF EXISTS locales;
				
				ALTER TABLE healthcare_entities DROP COLUMN IF EXISTS require_room_assignment;
			`,
		},
		{
			Version:     4,
			Description: "Rename locale to language for consistency",
			Up: `
				-- Rename columns for consistency with healthcare platform standards
				ALTER TABLE healthcare_entities DROP CONSTRAINT IF EXISTS healthcare_entities_language_check;
				ALTER TABLE healthcare_entities RENAME COLUMN language TO locale;
				ALTER TABLE healthcare_entities ADD CONSTRAINT healthcare_entities_locale_check 
					CHECK (locale IN ('en', 'fr', 'ar'));
				
				ALTER TABLE users DROP CONSTRAINT IF EXISTS users_preferred_language_check;
				ALTER TABLE users RENAME COLUMN preferred_language TO preferred_locale;
				ALTER TABLE users ADD CONSTRAINT users_preferred_locale_check 
					CHECK (preferred_locale IS NULL OR preferred_locale = '' OR preferred_locale IN ('en', 'fr', 'ar'));
			`,
			Down: `
				-- Revert column names and constraints
				ALTER TABLE healthcare_entities DROP CONSTRAINT IF EXISTS healthcare_entities_locale_check;
				ALTER TABLE healthcare_entities RENAME COLUMN locale TO language;
				ALTER TABLE healthcare_entities ADD CONSTRAINT healthcare_entities_language_check 
					CHECK (language IN ('en', 'fr', 'ar'));
				
				ALTER TABLE users DROP CONSTRAINT IF EXISTS users_preferred_locale_check;
				ALTER TABLE users RENAME COLUMN preferred_locale TO preferred_language;
				ALTER TABLE users ADD CONSTRAINT users_preferred_language_check 
					CHECK (preferred_language IS NULL OR preferred_language = '' OR preferred_language IN ('en', 'fr', 'ar'));
			`,
		},
		{
			Version:     6,
			Description: "Update healthcare_entities to use location service IDs",
			Up: `
				-- Add new location ID columns for integration with location service
				ALTER TABLE healthcare_entities 
				ADD COLUMN country_id INTEGER,
				ADD COLUMN state_id INTEGER,
				ADD COLUMN city_id INTEGER;

				-- Add indexes for the new foreign key columns
				CREATE INDEX IF NOT EXISTS idx_healthcare_entities_country_id ON healthcare_entities(country_id);
				CREATE INDEX IF NOT EXISTS idx_healthcare_entities_state_id ON healthcare_entities(state_id);
				CREATE INDEX IF NOT EXISTS idx_healthcare_entities_city_id ON healthcare_entities(city_id);

				-- Note: String columns (country, state, city) kept for backward compatibility
				-- They will be populated by API Gateway from location service
				-- In a future migration, we can make country_id NOT NULL and drop string columns
			`,
			Down: `
				-- Remove the location ID columns
				DROP INDEX IF EXISTS idx_healthcare_entities_country_id;
				DROP INDEX IF EXISTS idx_healthcare_entities_state_id;
				DROP INDEX IF EXISTS idx_healthcare_entities_city_id;
				
				ALTER TABLE healthcare_entities 
				DROP COLUMN IF EXISTS country_id,
				DROP COLUMN IF EXISTS state_id,
				DROP COLUMN IF EXISTS city_id;
			`,
		},
		{
			Version:     5,
			Description: "Update nationality field to use dropdown instead of text input and rename to nationality_id",
			Up: `
				-- Change nationality field from text to select type and update field name
				UPDATE form_fields 
				SET field_type = 'select',
					name = 'nationality_id',
					placeholder_text = ''
				WHERE form_type_id = (SELECT id FROM form_types WHERE name = 'patient') 
					AND name = 'nationality';
			`,
			Down: `
				-- Revert nationality field back to text input with original name
				UPDATE form_fields 
				SET field_type = 'text',
					name = 'nationality',
					placeholder_text = 'Enter nationality'
				WHERE form_type_id = (SELECT id FROM form_types WHERE name = 'patient') 
					AND name = 'nationality_id';
			`,
		},
		{
			Version:     6,
			Description: "Update healthcare entities to use country_id instead of country string",
			Up: `
				-- Add country_id column to healthcare_entities
				ALTER TABLE healthcare_entities ADD COLUMN country_id INTEGER;

				-- Convert existing country strings to country_id references
				-- This will be handled by application logic since we can't cross-reference
				-- other microservice databases directly in SQL
				
				-- Manually set country_id based on known country codes for existing data
				UPDATE healthcare_entities SET country_id = 1 WHERE UPPER(TRIM(country)) = 'CA';
				UPDATE healthcare_entities SET country_id = 2 WHERE UPPER(TRIM(country)) = 'US';
				UPDATE healthcare_entities SET country_id = 3 WHERE UPPER(TRIM(country)) = 'MA';
				UPDATE healthcare_entities SET country_id = 4 WHERE UPPER(TRIM(country)) = 'FR';

				-- Add index for the new foreign key column
				CREATE INDEX IF NOT EXISTS idx_healthcare_entities_country_id ON healthcare_entities(country_id);

				-- Note: Keep country string column temporarily for backward compatibility
				-- Will be removed in a future migration after verification
			`,
			Down: `
				-- Remove the country_id column and index
				DROP INDEX IF EXISTS idx_healthcare_entities_country_id;
				ALTER TABLE healthcare_entities DROP COLUMN IF EXISTS country_id;
			`,
		},
	}
}

// CreateMigrationsTable creates the schema_migrations table to track applied migrations
func CreateMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
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

	// Get all migrations
	migrations := GetMigrations()

	// Apply pending migrations
	for _, migration := range migrations {
		if applied[migration.Version] {
			log.Printf("Migration %d already applied: %s", migration.Version, migration.Description)
			continue
		}

		// Apply migration in a transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", migration.Version, err)
		}

		// Execute migration
		if _, err := tx.Exec(migration.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
		}

		// Record migration
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		log.Printf("Successfully applied migration %d", migration.Version)
	}

	return nil
}

// RollbackMigration rolls back a specific migration (for development use)
func RollbackMigration(db *sql.DB, version int) error {
	migrations := GetMigrations()
	
	// Find the migration to rollback
	var migration *Migration
	for _, m := range migrations {
		if m.Version == version {
			migration = &m
			break
		}
	}
	
	if migration == nil {
		return fmt.Errorf("migration version %d not found", version)
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for rollback %d: %w", version, err)
	}
	defer tx.Rollback()

	// Execute rollback
	if _, err := tx.Exec(migration.Down); err != nil {
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

	log.Printf("Successfully rolled back migration %d", version)
	return nil
}