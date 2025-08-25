-- Appointment Service Database Initialization
-- This script sets up the initial database schema and sample data

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create appointments table
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

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_appointments_patient_id ON appointments(patient_id);
CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id ON appointments(doctor_id);
CREATE INDEX IF NOT EXISTS idx_appointments_date_time ON appointments(date_time);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status);
CREATE INDEX IF NOT EXISTS idx_appointments_type ON appointments(type);
CREATE INDEX IF NOT EXISTS idx_appointments_active ON appointments(is_active);
CREATE INDEX IF NOT EXISTS idx_appointments_doctor_date ON appointments(doctor_id, date_time);

-- Composite index for conflict checking
CREATE INDEX IF NOT EXISTS idx_appointments_conflict_check 
    ON appointments(doctor_id, date_time, duration, status) 
    WHERE is_active = true AND status IN ('scheduled', 'confirmed', 'in-progress');

-- Update updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger
DROP TRIGGER IF EXISTS update_appointments_updated_at ON appointments;
CREATE TRIGGER update_appointments_updated_at
    BEFORE UPDATE ON appointments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Table for doctor schedules (future enhancement)
CREATE TABLE IF NOT EXISTS doctor_schedules (
    id SERIAL PRIMARY KEY,
    doctor_id INTEGER NOT NULL,
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6), -- 0 = Sunday
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    break_start TIME,
    break_end TIME,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(doctor_id, day_of_week)
);

CREATE INDEX IF NOT EXISTS idx_doctor_schedules_doctor_id ON doctor_schedules(doctor_id);
CREATE INDEX IF NOT EXISTS idx_doctor_schedules_day ON doctor_schedules(day_of_week);

DROP TRIGGER IF EXISTS update_doctor_schedules_updated_at ON doctor_schedules;
CREATE TRIGGER update_doctor_schedules_updated_at
    BEFORE UPDATE ON doctor_schedules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert sample appointments (assuming patient IDs 1-3 and doctor ID 2 from other services)
INSERT INTO appointments (
    patient_id, doctor_id, date_time, duration, type, status, reason, notes, created_by
) VALUES 
(
    1, 2, CURRENT_TIMESTAMP + INTERVAL '1 day', 30, 'consultation', 'scheduled',
    'Annual checkup', 'Patient requested morning appointment', 1
),
(
    2, 2, CURRENT_TIMESTAMP + INTERVAL '2 days', 45, 'follow-up', 'confirmed',
    'Blood pressure follow-up', 'Review medication effectiveness', 1
),
(
    3, 2, CURRENT_TIMESTAMP + INTERVAL '3 days', 30, 'consultation', 'scheduled',
    'Asthma consultation', 'Patient reports increased symptoms', 1
),
(
    1, 2, CURRENT_TIMESTAMP + INTERVAL '1 week', 30, 'follow-up', 'scheduled',
    'Lab results review', 'Discuss recent test results', 1
) ON CONFLICT DO NOTHING;

-- Insert default doctor schedule (Monday-Friday, 9 AM - 5 PM)
INSERT INTO doctor_schedules (doctor_id, day_of_week, start_time, end_time, break_start, break_end) VALUES
(2, 1, '09:00:00', '17:00:00', '12:00:00', '13:00:00'), -- Monday
(2, 2, '09:00:00', '17:00:00', '12:00:00', '13:00:00'), -- Tuesday
(2, 3, '09:00:00', '17:00:00', '12:00:00', '13:00:00'), -- Wednesday
(2, 4, '09:00:00', '17:00:00', '12:00:00', '13:00:00'), -- Thursday
(2, 5, '09:00:00', '17:00:00', '12:00:00', '13:00:00')  -- Friday
ON CONFLICT (doctor_id, day_of_week) DO NOTHING;