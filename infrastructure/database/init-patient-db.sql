-- Patient Service Database Initialization
-- This script sets up the initial database schema and sample data aligned with patient-service migrations

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create patients table
CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    healthcare_entity_id INTEGER NOT NULL,
    patient_id VARCHAR(50),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(10) NOT NULL CHECK (gender IN ('Male', 'Female', 'Other')),
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255), -- nullable, uniqueness enforced via partial index in migrations
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2) NOT NULL, -- ISO 3166-1 alpha-2
    nationality VARCHAR(100),
    preferred_language VARCHAR(10) CHECK (preferred_language IS NULL OR preferred_language = '' OR preferred_language IN ('English', 'French', 'Arabic')),
    marital_status VARCHAR(20) CHECK (marital_status IS NULL OR marital_status = '' OR marital_status IN ('Single', 'Married', 'Divorced', 'Widowed', 'Other')),
    occupation VARCHAR(200),
    insurance VARCHAR(255),
    policy_number VARCHAR(100),
    insurance_provider VARCHAR(255),
    national_id VARCHAR(50),
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

-- Indexes
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

-- Update updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger
DROP TRIGGER IF EXISTS update_patients_updated_at ON patients;
CREATE TRIGGER update_patients_updated_at
    BEFORE UPDATE ON patients
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Sample data minimal (aligning fields)
INSERT INTO patients (
    healthcare_entity_id, patient_id, first_name, last_name, date_of_birth, gender, phone, email,
    address, city, state, postal_code, country, nationality, preferred_language, marital_status,
    occupation, insurance, policy_number, insurance_provider, national_id,
    emergency_contact_name, emergency_contact_phone, emergency_contact_relationship,
    medical_history, allergies, medications, blood_type, created_by
) VALUES 
(
    1, 'P-ALICE-001', 'Alice', 'Johnson', '1985-03-15', 'Female', '555-0101', 'alice.johnson@email.com',
    '123 Main St', 'Anytown', 'CA', '12345', 'US', 'American', 'English', 'Married',
    'Engineer', 'Blue Cross Blue Shield', 'BC123456789', 'BCBS', '123-45-6789',
    'Bob Johnson', '555-0102', 'Spouse',
    'No significant medical history', 'None known', 'None', 'O+', 1
),
(
    1, 'P-BOB-002', 'Bob', 'Smith', '1978-07-22', 'Male', '555-0201', 'bob.smith@email.com',
    '456 Oak Ave', 'Somewhere', 'CA', '12346', 'US', 'American', 'English', 'Married',
    'Teacher', 'Aetna', 'AE987654321', 'Aetna', '987-65-4321',
    'Carol Smith', '555-0202', 'Spouse',
    'Hypertension, managed with medication', 'Penicillin', 'Lisinopril 10mg daily', 'A+', 1
)
ON CONFLICT DO NOTHING;