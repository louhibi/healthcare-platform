package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// InitDB initializes database connection
func InitDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	if dbName == "" {
		dbName = "appointment_service_db"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// RunMigrations creates necessary tables and runs migrations
func RunMigrations(db *sql.DB) error {
	// Create migrations table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Run migration 1: Create basic appointments table
	if err := runMigration(db, 1, `
		CREATE TABLE IF NOT EXISTS appointments (
			id SERIAL PRIMARY KEY,
			patient_id INTEGER NOT NULL,
			doctor_id INTEGER NOT NULL,
			date_time TIMESTAMP NOT NULL,
			duration INTEGER NOT NULL CHECK (duration >= 15 AND duration <= 480),
			type VARCHAR(50) NOT NULL CHECK (type IN ('consultation', 'follow-up', 'procedure', 'emergency')),
			status VARCHAR(50) NOT NULL CHECK (status IN ('scheduled', 'confirmed', 'in-progress', 'completed', 'cancelled', 'no-show')),
			reason TEXT NOT NULL,
			notes TEXT,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER NOT NULL
		);
	`); err != nil {
		return err
	}

	// Run migration 2: Add multi-tenant support and new fields
	if err := runMigration(db, 2, `
		-- Add healthcare_entity_id column if it doesn't exist
		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'appointments' AND column_name = 'healthcare_entity_id') THEN
				ALTER TABLE appointments ADD COLUMN healthcare_entity_id INTEGER;
				UPDATE appointments SET healthcare_entity_id = 1 WHERE healthcare_entity_id IS NULL;
				ALTER TABLE appointments ALTER COLUMN healthcare_entity_id SET NOT NULL;
			END IF;
		END $$;

		-- Add priority column if it doesn't exist
		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'appointments' AND column_name = 'priority') THEN
				ALTER TABLE appointments ADD COLUMN priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent'));
			END IF;
		END $$;

		-- Add room_number column if it doesn't exist
		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'appointments' AND column_name = 'room_number') THEN
				ALTER TABLE appointments ADD COLUMN room_number VARCHAR(50);
			END IF;
		END $$;
	`); err != nil {
		return err
	}

	// Run migration 3: Create basic indexes (before adding new columns)
	if err := runMigration(db, 3, `
		-- Basic indexes for better performance
		CREATE INDEX IF NOT EXISTS idx_appointments_patient_id ON appointments(patient_id);
		CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id ON appointments(doctor_id);
		CREATE INDEX IF NOT EXISTS idx_appointments_date_time ON appointments(date_time);
		CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status);
		CREATE INDEX IF NOT EXISTS idx_appointments_type ON appointments(type);
		CREATE INDEX IF NOT EXISTS idx_appointments_active ON appointments(is_active);
		CREATE INDEX IF NOT EXISTS idx_appointments_doctor_date ON appointments(doctor_id, date_time);
	`); err != nil {
		return err
	}

	// Run migration 3.5: Create indexes for new columns (after adding healthcare_entity_id, priority)
	if err := runMigration(db, 7, `
		-- Indexes for new multi-tenant columns
		CREATE INDEX IF NOT EXISTS idx_appointments_entity ON appointments(healthcare_entity_id);
		CREATE INDEX IF NOT EXISTS idx_appointments_priority ON appointments(priority);
		CREATE INDEX IF NOT EXISTS idx_appointments_entity_date ON appointments(healthcare_entity_id, date_time);

		-- Composite index for conflict checking
		CREATE INDEX IF NOT EXISTS idx_appointments_conflict_check 
			ON appointments(healthcare_entity_id, doctor_id, date_time, duration, status) 
			WHERE is_active = true AND status IN ('scheduled', 'confirmed', 'in-progress');
	`); err != nil {
		return err
	}

	// Run migration 4: Create triggers and functions
	if err := runMigration(db, 4, `
		-- Update updated_at trigger
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql';

		DROP TRIGGER IF EXISTS update_appointments_updated_at ON appointments;
		CREATE TRIGGER update_appointments_updated_at
			BEFORE UPDATE ON appointments
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`); err != nil {
		return err
	}

	// Run migration 5: Create doctor availability tables
	if err := runMigration(db, 5, `
		-- Doctor availability table for daily status and working hours
		CREATE TABLE IF NOT EXISTS doctor_availability (
			id SERIAL PRIMARY KEY,
			healthcare_entity_id INTEGER NOT NULL,
			doctor_id INTEGER NOT NULL,
			date DATE NOT NULL,
			status VARCHAR(50) NOT NULL CHECK (status IN ('available', 'unavailable', 'vacation', 'training', 'sick_leave', 'meeting')),
			start_time VARCHAR(5), -- HH:MM format
			end_time VARCHAR(5),   -- HH:MM format
			break_start VARCHAR(5), -- HH:MM format
			break_end VARCHAR(5),   -- HH:MM format
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER NOT NULL,
			UNIQUE(healthcare_entity_id, doctor_id, date)
		);

		-- Indexes for doctor availability
		CREATE INDEX IF NOT EXISTS idx_doctor_availability_doctor_id ON doctor_availability(doctor_id);
		CREATE INDEX IF NOT EXISTS idx_doctor_availability_date ON doctor_availability(date);
		CREATE INDEX IF NOT EXISTS idx_doctor_availability_status ON doctor_availability(status);
		CREATE INDEX IF NOT EXISTS idx_doctor_availability_entity ON doctor_availability(healthcare_entity_id);
		CREATE INDEX IF NOT EXISTS idx_doctor_availability_doctor_date ON doctor_availability(doctor_id, date);
		CREATE INDEX IF NOT EXISTS idx_doctor_availability_entity_date ON doctor_availability(healthcare_entity_id, date);

		-- Trigger for doctor availability
		DROP TRIGGER IF EXISTS update_doctor_availability_updated_at ON doctor_availability;
		CREATE TRIGGER update_doctor_availability_updated_at
			BEFORE UPDATE ON doctor_availability
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`); err != nil {
		return err
	}

	// Run migration 6: Create schedule templates table
	if err := runMigration(db, 6, `
		-- Schedule templates for recurring weekly patterns
		CREATE TABLE IF NOT EXISTS schedule_templates (
			id SERIAL PRIMARY KEY,
			healthcare_entity_id INTEGER NOT NULL,
			doctor_id INTEGER NOT NULL,
			day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6), -- 0 = Sunday
			start_time VARCHAR(5) NOT NULL, -- HH:MM format
			end_time VARCHAR(5) NOT NULL,   -- HH:MM format
			break_start VARCHAR(5), -- HH:MM format
			break_end VARCHAR(5),   -- HH:MM format
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(healthcare_entity_id, doctor_id, day_of_week)
		);

		-- Indexes for schedule templates
		CREATE INDEX IF NOT EXISTS idx_schedule_templates_doctor_id ON schedule_templates(doctor_id);
		CREATE INDEX IF NOT EXISTS idx_schedule_templates_day ON schedule_templates(day_of_week);
		CREATE INDEX IF NOT EXISTS idx_schedule_templates_entity ON schedule_templates(healthcare_entity_id);
		CREATE INDEX IF NOT EXISTS idx_schedule_templates_active ON schedule_templates(is_active);

		-- Trigger for schedule templates
		DROP TRIGGER IF EXISTS update_schedule_templates_updated_at ON schedule_templates;
		CREATE TRIGGER update_schedule_templates_updated_at
			BEFORE UPDATE ON schedule_templates
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`); err != nil {
		return err
	}

	// Run migration 9: Create appointment duration settings table
	if err := runMigration(db, 9, `
		-- Appointment duration settings per healthcare entity
		CREATE TABLE IF NOT EXISTS appointment_duration_settings (
			id SERIAL PRIMARY KEY,
			healthcare_entity_id INTEGER NOT NULL,
			appointment_type VARCHAR(20) NOT NULL CHECK (appointment_type IN ('consultation', 'follow-up', 'procedure', 'emergency')),
			duration_minutes INTEGER NOT NULL CHECK (duration_minutes >= 15 AND duration_minutes <= 480),
			is_default BOOLEAN DEFAULT false,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER,
			UNIQUE(healthcare_entity_id, appointment_type)
		);

		-- Insert default duration settings for sample healthcare entities (1-5)
		INSERT INTO appointment_duration_settings (healthcare_entity_id, appointment_type, duration_minutes, is_default, created_by)
		VALUES
			-- Entity 1 duration settings
			(1, 'consultation', 30, true, 1),
			(1, 'follow-up', 15, true, 1),
			(1, 'procedure', 60, true, 1),
			(1, 'emergency', 20, true, 1),
			-- Entity 2 duration settings
			(2, 'consultation', 30, true, 1),
			(2, 'follow-up', 15, true, 1),
			(2, 'procedure', 60, true, 1),
			(2, 'emergency', 20, true, 1),
			-- Entity 3 duration settings
			(3, 'consultation', 45, true, 1),
			(3, 'follow-up', 20, true, 1),
			(3, 'procedure', 90, true, 1),
			(3, 'emergency', 30, true, 1),
			-- Entity 4 duration settings
			(4, 'consultation', 30, true, 1),
			(4, 'follow-up', 15, true, 1),
			(4, 'procedure', 60, true, 1),
			(4, 'emergency', 15, true, 1),
			-- Entity 5 duration settings
			(5, 'consultation', 30, true, 1),
			(5, 'follow-up', 15, true, 1),
			(5, 'procedure', 60, true, 1),
			(5, 'emergency', 20, true, 1)
		ON CONFLICT (healthcare_entity_id, appointment_type) DO NOTHING;

		-- Create index for performance
		CREATE INDEX IF NOT EXISTS idx_appointment_duration_settings_entity_type ON appointment_duration_settings(healthcare_entity_id, appointment_type);

		-- Trigger for appointment duration settings
		DROP TRIGGER IF EXISTS update_appointment_duration_settings_updated_at ON appointment_duration_settings;
		CREATE TRIGGER update_appointment_duration_settings_updated_at
			BEFORE UPDATE ON appointment_duration_settings
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`); err != nil {
		return err
	}

	// Run migration 8: Create rooms table
	if err := runMigration(db, 8, `
		-- Rooms table for hospital room management
		CREATE TABLE IF NOT EXISTS rooms (
			id SERIAL PRIMARY KEY,
			healthcare_entity_id INTEGER NOT NULL,
			room_number VARCHAR(20) NOT NULL,
			room_name VARCHAR(100),
			room_type VARCHAR(20) NOT NULL CHECK (room_type IN ('consultation', 'examination', 'procedure', 'operating', 'emergency')),
			floor INTEGER DEFAULT 1,
			department VARCHAR(100),
			capacity INTEGER DEFAULT 1 CHECK (capacity >= 1),
			equipment TEXT, -- Comma-separated list or JSON
			is_active BOOLEAN DEFAULT true,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER,
			UNIQUE(healthcare_entity_id, room_number)
		);

		-- Insert sample rooms for common healthcare entity IDs (1-5)
		-- These correspond to entities that should exist in the user service
		INSERT INTO rooms (healthcare_entity_id, room_number, room_name, room_type, floor, department, capacity, equipment, created_by)
		VALUES
			-- Entity 1 rooms
			(1, '101', 'Consultation Room 1', 'consultation', 1, 'General', 2, 'Chair, Desk, Computer', 1),
			(1, '102', 'Consultation Room 2', 'consultation', 1, 'General', 2, 'Chair, Desk, Computer', 1),
			(1, '201', 'Examination Room 1', 'examination', 2, 'General', 1, 'Examination Table, Medical Equipment', 1),
			(1, '301', 'Procedure Room 1', 'procedure', 3, 'Surgery', 4, 'Surgical Table, Anesthesia Machine', 1),
			-- Entity 2 rooms
			(2, '101', 'Consultation Room 1', 'consultation', 1, 'General', 2, 'Chair, Desk, Computer', 1),
			(2, '102', 'Examination Room 1', 'examination', 1, 'General', 1, 'Examination Table, Medical Equipment', 1),
			(2, '201', 'Procedure Room 1', 'procedure', 2, 'Surgery', 4, 'Surgical Table, Monitors', 1),
			-- Entity 3 rooms
			(3, '101', 'Consultation Room 1', 'consultation', 1, 'General', 2, 'Chair, Desk, Computer', 1),
			(3, '201', 'Operating Room 1', 'operating', 2, 'Surgery', 8, 'Full Surgical Suite, Anesthesia, Monitors', 1),
			-- Entity 4 rooms
			(4, '101', 'Emergency Room 1', 'emergency', 1, 'Emergency', 4, 'Emergency Equipment, Monitors', 1),
			(4, '102', 'Consultation Room 1', 'consultation', 1, 'General', 2, 'Chair, Desk, Computer', 1),
			-- Entity 5 rooms
			(5, '101', 'Consultation Room 1', 'consultation', 1, 'General', 2, 'Chair, Desk, Computer', 1),
			(5, '102', 'Examination Room 1', 'examination', 1, 'General', 1, 'Examination Table, Medical Equipment', 1)
		ON CONFLICT (healthcare_entity_id, room_number) DO NOTHING;

		-- Create indexes for performance
		CREATE INDEX IF NOT EXISTS idx_rooms_entity_active ON rooms(healthcare_entity_id, is_active);
		CREATE INDEX IF NOT EXISTS idx_rooms_type ON rooms(room_type);
		CREATE INDEX IF NOT EXISTS idx_rooms_floor ON rooms(healthcare_entity_id, floor);

		-- Trigger for rooms
		DROP TRIGGER IF EXISTS update_rooms_updated_at ON rooms;
		CREATE TRIGGER update_rooms_updated_at
			BEFORE UPDATE ON rooms
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`); err != nil {
		return err
	}

	// Run migration 11: Fix timezone handling - convert all timestamps to UTC
	if err := runMigration(db, 11, `
		-- Ensure all timestamp columns use TIMESTAMPTZ for proper UTC storage
		-- Update existing tables to use TIMESTAMPTZ
		
		-- Fix appointments table timestamps
		DO $$ 
		BEGIN 
			-- Check if column is already TIMESTAMPTZ
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'appointments' 
				AND column_name = 'created_at' 
				AND data_type = 'timestamp with time zone'
			) THEN
				-- Convert to TIMESTAMPTZ (assumes current data is already in UTC)
				ALTER TABLE appointments 
				ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';
				
				ALTER TABLE appointments 
				ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC';
				
				ALTER TABLE appointments 
				ALTER COLUMN date_time TYPE TIMESTAMPTZ USING date_time AT TIME ZONE 'UTC';
			END IF;
		END $$;

		-- Fix doctor_availability table timestamps  
		DO $$ 
		BEGIN 
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'doctor_availability' 
				AND column_name = 'created_at' 
				AND data_type = 'timestamp with time zone'
			) THEN
				ALTER TABLE doctor_availability 
				ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';
				
				ALTER TABLE doctor_availability 
				ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC';
			END IF;
		END $$;

		-- Fix other tables
		DO $$ 
		BEGIN 
			-- Schedule templates
			IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'schedule_templates') THEN
				IF NOT EXISTS (
					SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'schedule_templates' 
					AND column_name = 'created_at' 
					AND data_type = 'timestamp with time zone'
				) THEN
					ALTER TABLE schedule_templates 
					ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';
					
					ALTER TABLE schedule_templates 
					ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC';
				END IF;
			END IF;

			-- Rooms table
			IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'rooms') THEN
				IF NOT EXISTS (
					SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'rooms' 
					AND column_name = 'created_at' 
					AND data_type = 'timestamp with time zone'
				) THEN
					ALTER TABLE rooms 
					ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';
					
					ALTER TABLE rooms 
					ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC';
				END IF;
			END IF;

			-- Duration settings tables
			IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'appointment_duration_settings') THEN
				IF NOT EXISTS (
					SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'appointment_duration_settings' 
					AND column_name = 'created_at' 
					AND data_type = 'timestamp with time zone'
				) THEN
					ALTER TABLE appointment_duration_settings 
					ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';
					
					ALTER TABLE appointment_duration_settings 
					ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC';
				END IF;
			END IF;

			-- Duration options table
			IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'appointment_duration_options') THEN
				IF NOT EXISTS (
					SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'appointment_duration_options' 
					AND column_name = 'created_at' 
					AND data_type = 'timestamp with time zone'
				) THEN
					ALTER TABLE appointment_duration_options 
					ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';
					
					ALTER TABLE appointment_duration_options 
					ALTER COLUMN updated_at TYPE TIMESTAMPTZ USING updated_at AT TIME ZONE 'UTC';
				END IF;
			END IF;
		END $$;

		-- Create timezone conversion helper functions
		CREATE OR REPLACE FUNCTION convert_utc_to_entity_timezone(
			utc_timestamp TIMESTAMPTZ, 
			entity_timezone TEXT
		) RETURNS TIMESTAMPTZ AS $$
		BEGIN
			RETURN utc_timestamp AT TIME ZONE entity_timezone;
		END;
		$$ LANGUAGE plpgsql IMMUTABLE;

		CREATE OR REPLACE FUNCTION convert_entity_to_utc_timezone(
			local_timestamp TIMESTAMPTZ, 
			entity_timezone TEXT
		) RETURNS TIMESTAMPTZ AS $$
		BEGIN
			RETURN local_timestamp AT TIME ZONE entity_timezone AT TIME ZONE 'UTC';
		END;
		$$ LANGUAGE plpgsql IMMUTABLE;
	`); err != nil {
		return err
	}

	// Run migration 10: Create appointment duration options table for multiple durations per type
	if err := runMigration(db, 10, `
		-- Appointment duration options - allows multiple duration choices per appointment type
		CREATE TABLE IF NOT EXISTS appointment_duration_options (
			id SERIAL PRIMARY KEY,
			healthcare_entity_id INTEGER NOT NULL,
			appointment_type VARCHAR(20) NOT NULL CHECK (appointment_type IN ('consultation', 'follow-up', 'procedure', 'emergency')),
			duration_minutes INTEGER NOT NULL CHECK (duration_minutes >= 15 AND duration_minutes <= 480),
			is_default BOOLEAN DEFAULT false,
			is_active BOOLEAN DEFAULT true,
			display_order INTEGER DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER,
			UNIQUE(healthcare_entity_id, appointment_type, duration_minutes)
		);

		-- Migrate existing duration settings to new options table
		INSERT INTO appointment_duration_options (healthcare_entity_id, appointment_type, duration_minutes, is_default, is_active, display_order, created_by)
		SELECT 
			healthcare_entity_id,
			appointment_type,
			duration_minutes,
			is_default,
			is_active,
			1,
			created_by
		FROM appointment_duration_settings
		ON CONFLICT (healthcare_entity_id, appointment_type, duration_minutes) DO NOTHING;

		-- Add additional common duration options for each type
		INSERT INTO appointment_duration_options (healthcare_entity_id, appointment_type, duration_minutes, is_default, is_active, display_order, created_by)
		VALUES
			-- Consultation additional options
			(1, 'consultation', 15, false, true, 2, 1),
			(1, 'consultation', 45, false, true, 3, 1),
			(1, 'consultation', 60, false, true, 4, 1),
			-- Follow-up additional options  
			(1, 'follow-up', 20, false, true, 2, 1),
			(1, 'follow-up', 25, false, true, 3, 1),
			(1, 'follow-up', 30, false, true, 4, 1),
			-- Procedure additional options
			(1, 'procedure', 30, false, true, 2, 1),
			(1, 'procedure', 90, false, true, 3, 1),
			(1, 'procedure', 120, false, true, 4, 1),
			-- Emergency additional options
			(1, 'emergency', 15, false, true, 2, 1),
			(1, 'emergency', 30, false, true, 3, 1),
			(1, 'emergency', 45, false, true, 4, 1),
			
			-- Entity 2 options (similar pattern)
			(2, 'consultation', 15, false, true, 2, 1),
			(2, 'consultation', 45, false, true, 3, 1),
			(2, 'follow-up', 20, false, true, 2, 1),
			(2, 'follow-up', 25, false, true, 3, 1),
			(2, 'procedure', 30, false, true, 2, 1),
			(2, 'procedure', 90, false, true, 3, 1),
			(2, 'emergency', 15, false, true, 2, 1),
			(2, 'emergency', 30, false, true, 3, 1)
		ON CONFLICT (healthcare_entity_id, appointment_type, duration_minutes) DO NOTHING;

		-- Create indexes for performance
		CREATE INDEX IF NOT EXISTS idx_appointment_duration_options_entity_type ON appointment_duration_options(healthcare_entity_id, appointment_type);
		CREATE INDEX IF NOT EXISTS idx_appointment_duration_options_active ON appointment_duration_options(healthcare_entity_id, appointment_type, is_active);
		CREATE INDEX IF NOT EXISTS idx_appointment_duration_options_order ON appointment_duration_options(healthcare_entity_id, appointment_type, display_order);

		-- Trigger for appointment duration options
		DROP TRIGGER IF EXISTS update_appointment_duration_options_updated_at ON appointment_duration_options;
		CREATE TRIGGER update_appointment_duration_options_updated_at
			BEFORE UPDATE ON appointment_duration_options
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`); err != nil {
		return err
	}

	// Run migration 12: Convert doctor availability to use UTC timestamps
	if err := runMigration(db, 12, `
		-- Add new TIMESTAMPTZ columns for working hours
		ALTER TABLE doctor_availability 
		ADD COLUMN IF NOT EXISTS start_datetime TIMESTAMPTZ,
		ADD COLUMN IF NOT EXISTS end_datetime TIMESTAMPTZ,
		ADD COLUMN IF NOT EXISTS break_start_datetime TIMESTAMPTZ,
		ADD COLUMN IF NOT EXISTS break_end_datetime TIMESTAMPTZ;

		-- We'll populate these in the application code during transition
		-- For now, keep both old and new columns to ensure compatibility
	`); err != nil {
		return err
	}

	// Run migration 13: Remove legacy date/time columns and enforce UTC-only storage
	if err := runMigration(db, 13, `
		-- First, populate UTC datetime values for records that only have legacy values
		UPDATE doctor_availability 
		SET start_datetime = (date + start_time::time) AT TIME ZONE 'UTC'
		WHERE start_datetime IS NULL AND start_time IS NOT NULL;

		UPDATE doctor_availability 
		SET end_datetime = (date + end_time::time) AT TIME ZONE 'UTC'
		WHERE end_datetime IS NULL AND end_time IS NOT NULL;

		UPDATE doctor_availability 
		SET break_start_datetime = (date + break_start::time) AT TIME ZONE 'UTC'
		WHERE break_start_datetime IS NULL AND break_start IS NOT NULL;

		UPDATE doctor_availability 
		SET break_end_datetime = (date + break_end::time) AT TIME ZONE 'UTC'
		WHERE break_end_datetime IS NULL AND break_end IS NOT NULL;

		-- Remove the unique constraint that uses the legacy date column
		ALTER TABLE doctor_availability 
		DROP CONSTRAINT IF EXISTS doctor_availability_healthcare_entity_id_doctor_id_date_key;

		-- Drop old indexes that referenced the legacy date column
		DROP INDEX IF EXISTS idx_doctor_availability_date;
		DROP INDEX IF EXISTS idx_doctor_availability_doctor_date;
		DROP INDEX IF EXISTS idx_doctor_availability_entity_date;

		-- Drop the legacy columns
		ALTER TABLE doctor_availability 
		DROP COLUMN IF EXISTS date,
		DROP COLUMN IF EXISTS start_time,
		DROP COLUMN IF EXISTS end_time,
		DROP COLUMN IF EXISTS break_start,
		DROP COLUMN IF EXISTS break_end;

		-- Create simple indexes on the UTC datetime columns
		CREATE INDEX IF NOT EXISTS idx_doctor_availability_start_datetime 
		ON doctor_availability(start_datetime);
		
		CREATE INDEX IF NOT EXISTS idx_doctor_availability_doctor_start 
		ON doctor_availability(doctor_id, start_datetime);
		
		CREATE INDEX IF NOT EXISTS idx_doctor_availability_entity_start 
		ON doctor_availability(healthcare_entity_id, start_datetime);
	`); err != nil {
		return err
	}

	// Run migration 8: Update appointments to use room_id instead of room_number
	if err := runMigration(db, 8, `
		-- Add room_id column
		ALTER TABLE appointments ADD COLUMN IF NOT EXISTS room_id INTEGER;
		
		-- Create foreign key constraint
		ALTER TABLE appointments ADD CONSTRAINT IF NOT EXISTS fk_appointments_room_id 
			FOREIGN KEY (room_id) REFERENCES rooms(id);
		
		-- Create index for room_id
		CREATE INDEX IF NOT EXISTS idx_appointments_room_id ON appointments(room_id);
		
		-- Create composite index for room conflict checking
		CREATE INDEX IF NOT EXISTS idx_appointments_room_conflict 
			ON appointments(healthcare_entity_id, room_id, date_time, duration, status) 
			WHERE is_active = true AND status IN ('scheduled', 'confirmed', 'in-progress') AND room_id IS NOT NULL;
			
		-- Remove room_number column (if it exists)
		ALTER TABLE appointments DROP COLUMN IF EXISTS room_number;
	`); err != nil {
		return err
	}

	return nil
}

// runMigration executes a migration if it hasn't been applied yet
func runMigration(db *sql.DB, version int, query string) error {
	// Check if migration has already been applied
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Migration already applied
	}

	// Execute the migration
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	// Record that migration was applied
	_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
	return err
}