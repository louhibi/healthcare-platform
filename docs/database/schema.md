# Database Schema Documentation

Comprehensive database schema documentation for the Healthcare Management Platform.

## Overview

The Healthcare Platform uses a **database-per-service** architecture with PostgreSQL databases. Each microservice maintains its own isolated database to ensure data independence and service autonomy.

## Database Architecture

### Service Databases
- **User Service**: `user_service_db` (Port 5432)
- **Patient Service**: `patient_service_db` (Port 5433)  
- **Appointment Service**: `appointment_service_db` (Port 5434)

### Design Principles
- **Data Isolation**: Each service owns its data completely
- **No Cross-Service Queries**: Services communicate via APIs only
- **Referential Integrity**: Foreign keys reference IDs from other services via API calls
- **Eventual Consistency**: Accepted for non-critical cross-service data

## User Service Database

### Users Table
**Purpose**: Store user accounts, authentication, and role information

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'doctor', 'nurse', 'staff')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Columns
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Unique user identifier |
| `email` | VARCHAR(255) | UNIQUE, NOT NULL | User email (login) |
| `password_hash` | VARCHAR(255) | NOT NULL | bcrypt hashed password |
| `first_name` | VARCHAR(100) | NOT NULL | User's first name |
| `last_name` | VARCHAR(100) | NOT NULL | User's last name |
| `role` | VARCHAR(50) | CHECK constraint | User role (admin, doctor, nurse, staff) |
| `is_active` | BOOLEAN | DEFAULT true | Account status |
| `created_at` | TIMESTAMP | DEFAULT NOW() | Account creation time |
| `updated_at` | TIMESTAMP | AUTO-UPDATE | Last modification time |

#### Indexes
```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(is_active);
```

#### Triggers
```sql
-- Auto-update updated_at column
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

## Patient Service Database

### Patients Table
**Purpose**: Store patient personal information, medical data, and emergency contacts

```sql
CREATE TABLE patients (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL,
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(50),
    zip_code VARCHAR(20),
    insurance VARCHAR(255),
    policy_number VARCHAR(100),
    emergency_contact_name VARCHAR(200),
    emergency_contact_phone VARCHAR(20),
    emergency_contact_relationship VARCHAR(50),
    medical_history TEXT,
    allergies TEXT,
    medications TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL
);
```

#### Columns
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Unique patient identifier |
| `first_name` | VARCHAR(100) | NOT NULL | Patient's first name |
| `last_name` | VARCHAR(100) | NOT NULL | Patient's last name |
| `date_of_birth` | DATE | NOT NULL | Birth date for age calculation |
| `gender` | VARCHAR(10) | CHECK constraint | Gender (male, female, other) |
| `phone` | VARCHAR(20) | NOT NULL | Primary phone number |
| `email` | VARCHAR(255) | NOT NULL | Email address |
| `address` | TEXT | | Street address |
| `city` | VARCHAR(100) | | City |
| `state` | VARCHAR(50) | | State/Province |
| `zip_code` | VARCHAR(20) | | Postal code |
| `insurance` | VARCHAR(255) | | Insurance provider name |
| `policy_number` | VARCHAR(100) | | Insurance policy number |
| `emergency_contact_name` | VARCHAR(200) | | Emergency contact name |
| `emergency_contact_phone` | VARCHAR(20) | | Emergency contact phone |
| `emergency_contact_relationship` | VARCHAR(50) | | Relationship to patient |
| `medical_history` | TEXT | | Patient's medical history |
| `allergies` | TEXT | | Known allergies |
| `medications` | TEXT | | Current medications |
| `is_active` | BOOLEAN | DEFAULT true | Patient record status |
| `created_at` | TIMESTAMP | DEFAULT NOW() | Record creation time |
| `updated_at` | TIMESTAMP | AUTO-UPDATE | Last modification time |
| `created_by` | INTEGER | NOT NULL | User ID who created record |

#### Indexes
```sql
CREATE INDEX idx_patients_name ON patients(first_name, last_name);
CREATE INDEX idx_patients_email ON patients(email);
CREATE INDEX idx_patients_phone ON patients(phone);
CREATE INDEX idx_patients_dob ON patients(date_of_birth);
CREATE INDEX idx_patients_active ON patients(is_active);
CREATE INDEX idx_patients_created_by ON patients(created_by);
CREATE INDEX idx_patients_search ON patients USING gin(to_tsvector('english', first_name || ' ' || last_name || ' ' || email));
```

#### Triggers
```sql
CREATE TRIGGER update_patients_updated_at
    BEFORE UPDATE ON patients
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

## Appointment Service Database

### Appointments Table
**Purpose**: Store appointment scheduling information and status

```sql
CREATE TABLE appointments (
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
```

#### Columns
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Unique appointment identifier |
| `patient_id` | INTEGER | NOT NULL | Patient ID (references patient service) |
| `doctor_id` | INTEGER | NOT NULL | Doctor ID (references user service) |
| `date_time` | TIMESTAMP | NOT NULL | Appointment date and time |
| `duration` | INTEGER | CHECK 15-480 | Appointment duration in minutes |
| `type` | VARCHAR(50) | CHECK constraint | Appointment type |
| `status` | VARCHAR(50) | CHECK constraint | Current appointment status |
| `reason` | TEXT | NOT NULL | Reason for appointment |
| `notes` | TEXT | | Additional notes |
| `is_active` | BOOLEAN | DEFAULT true | Appointment record status |
| `created_at` | TIMESTAMP | DEFAULT NOW() | Record creation time |
| `updated_at` | TIMESTAMP | AUTO-UPDATE | Last modification time |
| `created_by` | INTEGER | NOT NULL | User ID who created appointment |

#### Appointment Types
- **consultation**: Initial patient consultation
- **follow-up**: Follow-up appointment
- **procedure**: Medical procedure
- **emergency**: Emergency appointment

#### Appointment Status Flow
```
scheduled → confirmed → in-progress → completed
    ↓           ↓            ↓
cancelled   cancelled   cancelled
    ↓           ↓            ↓
no-show     no-show     no-show
```

#### Indexes
```sql
CREATE INDEX idx_appointments_patient_id ON appointments(patient_id);
CREATE INDEX idx_appointments_doctor_id ON appointments(doctor_id);
CREATE INDEX idx_appointments_date_time ON appointments(date_time);
CREATE INDEX idx_appointments_status ON appointments(status);
CREATE INDEX idx_appointments_type ON appointments(type);
CREATE INDEX idx_appointments_active ON appointments(is_active);
CREATE INDEX idx_appointments_doctor_date ON appointments(doctor_id, date_time);

-- Specialized index for conflict detection
CREATE INDEX idx_appointments_conflict_check 
    ON appointments(doctor_id, date_time, duration, status) 
    WHERE is_active = true AND status IN ('scheduled', 'confirmed', 'in-progress');
```

### Doctor Schedules Table
**Purpose**: Store doctor working hours and availability

```sql
CREATE TABLE doctor_schedules (
    id SERIAL PRIMARY KEY,
    doctor_id INTEGER NOT NULL,
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    break_start TIME,
    break_end TIME,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(doctor_id, day_of_week)
);
```

#### Columns
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Unique schedule identifier |
| `doctor_id` | INTEGER | NOT NULL | Doctor ID (references user service) |
| `day_of_week` | INTEGER | CHECK 0-6 | Day of week (0=Sunday, 6=Saturday) |
| `start_time` | TIME | NOT NULL | Work start time |
| `end_time` | TIME | NOT NULL | Work end time |
| `break_start` | TIME | | Break start time |
| `break_end` | TIME | | Break end time |
| `is_available` | BOOLEAN | DEFAULT true | Doctor availability |
| `created_at` | TIMESTAMP | DEFAULT NOW() | Record creation time |
| `updated_at` | TIMESTAMP | AUTO-UPDATE | Last modification time |

#### Indexes
```sql
CREATE INDEX idx_doctor_schedules_doctor_id ON doctor_schedules(doctor_id);
CREATE INDEX idx_doctor_schedules_day ON doctor_schedules(day_of_week);
CREATE INDEX idx_doctor_schedules_available ON doctor_schedules(is_available);
```

## Cross-Service Relationships

### Logical Relationships
Although databases are isolated, logical relationships exist:

```
Users (User Service)
  ↓ (doctor_id)
Appointments (Appointment Service)
  ↓ (patient_id)  
Patients (Patient Service)
```

### Data Consistency Strategy
- **API Validation**: Services validate foreign key references via API calls
- **Eventual Consistency**: Cross-service data may be temporarily inconsistent
- **Compensating Actions**: Failed operations trigger compensating transactions
- **Audit Logging**: All cross-service operations logged for troubleshooting

## Common Database Functions

### Timestamp Update Function
Used across all services for automatic `updated_at` column updates:

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
```

### Sample Data Queries

#### Get User with Role
```sql
SELECT id, email, first_name, last_name, role, is_active 
FROM users 
WHERE email = 'doctor@hospital.com' AND is_active = true;
```

#### Search Patients
```sql
SELECT id, first_name, last_name, email, phone, 
       EXTRACT(YEAR FROM AGE(date_of_birth)) as age
FROM patients 
WHERE is_active = true 
  AND (LOWER(first_name) LIKE '%alice%' OR LOWER(last_name) LIKE '%alice%')
ORDER BY last_name, first_name;
```

#### Get Doctor's Appointments
```sql
SELECT a.id, a.patient_id, a.date_time, a.duration, a.type, a.status, a.reason
FROM appointments a
WHERE a.doctor_id = 2 
  AND a.date_time >= CURRENT_DATE 
  AND a.date_time < CURRENT_DATE + INTERVAL '1 day'
  AND a.is_active = true
ORDER BY a.date_time;
```

#### Check Appointment Conflicts
```sql
SELECT COUNT(*) as conflicts
FROM appointments 
WHERE doctor_id = 2
  AND is_active = true
  AND status IN ('scheduled', 'confirmed', 'in-progress')
  AND (
    (date_time < '2024-01-15 10:30:00' AND date_time + (duration || ' minutes')::interval > '2024-01-15 10:00:00') OR
    (date_time < '2024-01-15 10:00:00' AND date_time + (duration || ' minutes')::interval > '2024-01-15 10:30:00')
  );
```

## Performance Considerations

### Query Optimization
- **Selective Indexes**: Indexes on frequently queried columns
- **Composite Indexes**: Multi-column indexes for complex queries
- **Partial Indexes**: Indexes with WHERE clauses for specific conditions
- **Full-Text Search**: GIN indexes for text search capabilities

### Connection Management
- **Connection Pooling**: Configured per service (default: 10 connections)
- **Idle Timeout**: Automatic connection cleanup
- **Max Connections**: Limited to prevent resource exhaustion

### Maintenance Tasks
```sql
-- Regular maintenance (run weekly)
VACUUM ANALYZE users;
VACUUM ANALYZE patients;
VACUUM ANALYZE appointments;

-- Index maintenance (run monthly)
REINDEX INDEX idx_patients_name;
REINDEX INDEX idx_appointments_conflict_check;

-- Statistics update (run daily)
ANALYZE users;
ANALYZE patients;
ANALYZE appointments;
```

## Security Features

### Access Control
- **Database Users**: Separate database users per service
- **Minimal Privileges**: Only necessary permissions granted
- **Connection Security**: SSL/TLS connections in production
- **Network Isolation**: Database access restricted to service containers

### Data Protection
- **Encryption at Rest**: Database-level encryption (production)
- **Backup Encryption**: Encrypted database backups
- **Audit Logging**: All data modifications logged
- **Data Masking**: PII masking for development environments

### HIPAA Compliance
- **Patient Data Protection**: All patient data encrypted
- **Access Auditing**: Complete audit trail of data access
- **Data Retention**: Configurable retention policies
- **Breach Detection**: Monitoring for unauthorized access

## Backup and Recovery

### Backup Strategy
```bash
# Daily full backup
pg_dump -U postgres -h localhost -p 5432 user_service_db > backup_user_$(date +%Y%m%d).sql

# Transaction log backups (continuous)
pg_receivewal -U postgres -h localhost -p 5432 -D /backup/wal/
```

### Recovery Procedures
```bash
# Point-in-time recovery
pg_restore -U postgres -h localhost -p 5432 -d user_service_db backup_user_20240101.sql

# Verify data integrity
SELECT COUNT(*) FROM users WHERE created_at >= '2024-01-01';
```

## Migration Management

### Schema Versioning
Each service manages its own migrations:

```sql
-- Migration version tracking
CREATE TABLE schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Example migration
-- V001__create_users_table.sql
-- V002__add_user_roles.sql
-- V003__add_user_indexes.sql
```

### Migration Best Practices
- **Backward Compatibility**: Ensure migrations don't break existing code
- **Rollback Scripts**: Provide rollback for each migration
- **Testing**: Test migrations on copy of production data
- **Incremental Changes**: Small, focused migrations

## Monitoring and Observability

### Database Metrics
```sql
-- Connection monitoring
SELECT * FROM pg_stat_activity WHERE datname = 'user_service_db';

-- Query performance
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC LIMIT 10;

-- Index usage
SELECT schemaname, tablename, indexname, idx_scan 
FROM pg_stat_user_indexes 
WHERE idx_scan = 0;
```

### Health Checks
```sql
-- Database health check
SELECT 1 as healthy;

-- Table row counts
SELECT 
  'users' as table_name, COUNT(*) as row_count FROM users
UNION ALL
SELECT 
  'patients' as table_name, COUNT(*) as row_count FROM patients
UNION ALL
SELECT 
  'appointments' as table_name, COUNT(*) as row_count FROM appointments;
```

## Development Guidelines

### Schema Changes
1. **Create Migration**: Write SQL migration script
2. **Test Locally**: Test on local development database
3. **Peer Review**: Code review for schema changes
4. **Staging Deploy**: Deploy to staging environment
5. **Production Deploy**: Careful production deployment

### Best Practices
- **Naming Conventions**: Consistent table and column naming
- **Data Types**: Appropriate data types for storage efficiency
- **Constraints**: Use database constraints for data integrity
- **Documentation**: Document all schema changes

---

**Last Updated**: 2024-01-01  
**Maintainer**: Healthcare Platform Team