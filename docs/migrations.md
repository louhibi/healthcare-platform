# Database Migrations

This document explains the versioned database migration system implemented for the Healthcare Platform.

## Overview

The Healthcare Platform uses a versioned migration system to manage database schema changes across all microservices. Each service maintains its own migrations and schema version tracking.

## Migration Structure

### Migration Files
- **User Service**: `services/user-service/migrations.go`
- **Patient Service**: `services/patient-service/migrations.go`
- **Appointment Service**: `services/appointment-service/migrations.go` (future)

### Migration Format
Each migration contains:
- **Version**: Unique integer identifier
- **Description**: Human-readable description
- **Up**: SQL statements to apply the migration
- **Down**: SQL statements to rollback the migration

```go
type Migration struct {
    Version     int
    Description string
    Up          string
    Down        string
}
```

## Migration Tracking

Each database contains a `schema_migrations` table that tracks applied migrations:

```sql
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    description TEXT NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Current Migrations

### User Service (Version 1)
- Creates `healthcare_entities` table
- Creates `users` table with multi-tenant support
- Adds indexes and triggers
- Supports healthcare entities across Canada, USA, Morocco, France

### Patient Service (Version 1)
- Creates `patients` table with comprehensive fields
- Multi-tenant isolation by healthcare entity
- International support (addresses, phone formats, etc.)
- Optional field constraints (blood_type, preferred_language, marital_status)
- Unique constraints for patient_id within entities

## Migration Management

### Automatic Migrations
Migrations run automatically when services start:

```go
// In main.go
if err := RunMigrations(db); err != nil {
    log.Fatal("Failed to run migrations:", err)
}
```

### Manual Migration Management
Use the migration CLI tool:

```bash
# Run all pending migrations
go run scripts/migrate.go -service=patient -action=up

# Check migration status
go run scripts/migrate.go -service=patient -action=status

# Rollback a specific migration (development only)
go run scripts/migrate.go -service=patient -action=down -version=1
```

### Migration CLI Options
- `-service`: Service name (user, patient, appointment)
- `-action`: Action to perform (up, down, status)
- `-version`: Migration version (required for down action)
- `-host`: Database host (default: localhost)
- `-port`: Database port (default: 5432)
- `-user`: Database user (default: postgres)
- `-password`: Database password (default: postgres)

## Adding New Migrations

### Step 1: Add Migration to Service
Add a new migration to the service's `migrations.go` file:

```go
{
    Version:     2,
    Description: "Add patient medical records table",
    Up: `
        CREATE TABLE patient_medical_records (
            id SERIAL PRIMARY KEY,
            patient_id INTEGER NOT NULL REFERENCES patients(id),
            -- ... other fields
        );
    `,
    Down: `
        DROP TABLE IF EXISTS patient_medical_records CASCADE;
    `,
}
```

### Step 2: Rebuild and Restart Service
The new migration will be applied automatically when the service starts:

```bash
docker compose build patient-service
docker compose restart patient-service
```

### Step 3: Verify Migration
Check that the migration was applied:

```bash
docker exec healthcare-patient-db psql -U postgres -d patient_service_db -c "SELECT * FROM schema_migrations;"
```

## Best Practices

### Migration Guidelines
1. **Always test migrations** on a copy of production data
2. **Make migrations reversible** when possible
3. **Use transactions** for complex migrations
4. **Avoid data loss** - use additive changes when possible
5. **Document breaking changes** in migration descriptions

### Schema Changes
1. **Adding columns**: Safe, use `ALTER TABLE ADD COLUMN`
2. **Removing columns**: Dangerous, consider deprecation first
3. **Changing constraints**: Use `ALTER TABLE` with appropriate checks
4. **Renaming**: Create new column, migrate data, drop old column

### Multi-Service Considerations
1. **Coordinate schema changes** across services
2. **Maintain backward compatibility** during deployments
3. **Use feature flags** for major changes
4. **Test service interactions** after schema changes

## Production Deployment

### Pre-Deployment
1. **Backup databases** before applying migrations
2. **Test migrations** on staging environment
3. **Review migration impact** on existing data
4. **Plan rollback strategy** if needed

### Deployment Process
1. **Stop services** (if required for breaking changes)
2. **Apply migrations** using the CLI tool
3. **Verify schema changes** 
4. **Start services** with new code
5. **Monitor applications** for issues

### Rollback Process
1. **Stop affected services**
2. **Rollback migrations** using CLI tool
3. **Restore database backup** (if needed)
4. **Deploy previous service version**
5. **Verify system functionality**

## Troubleshooting

### Common Issues

#### Migration Fails to Apply
```bash
# Check database connection
docker exec healthcare-patient-db psql -U postgres -d patient_service_db -c "SELECT 1;"

# Check current schema state
docker exec healthcare-patient-db psql -U postgres -d patient_service_db -c "\dt"

# Check applied migrations
docker exec healthcare-patient-db psql -U postgres -d patient_service_db -c "SELECT * FROM schema_migrations;"
```

#### Service Won't Start After Migration
```bash
# Check service logs
docker compose logs patient-service

# Check database logs
docker compose logs patient-db

# Verify database connectivity
docker exec healthcare-patient-service ping healthcare-patient-db
```

#### Data Constraint Violations
```bash
# Check constraint violations
docker exec healthcare-patient-db psql -U postgres -d patient_service_db -c "
SELECT conname, conrelid::regclass 
FROM pg_constraint 
WHERE contype = 'c' AND NOT convalidated;
"

# Fix data before reapplying migration
# ... fix data issues ...

# Retry migration
docker compose restart patient-service
```

## Examples

### Example 1: Adding a New Optional Field
```go
{
    Version:     2,
    Description: "Add emergency_contact_address to patients",
    Up: `
        ALTER TABLE patients 
        ADD COLUMN emergency_contact_address TEXT;
        
        CREATE INDEX IF NOT EXISTS idx_patients_emergency_address 
        ON patients(emergency_contact_address);
    `,
    Down: `
        DROP INDEX IF EXISTS idx_patients_emergency_address;
        ALTER TABLE patients DROP COLUMN IF EXISTS emergency_contact_address;
    `,
}
```

### Example 2: Modifying Constraints
```go
{
    Version:     3,
    Description: "Add new blood types to constraint",
    Up: `
        ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_blood_type_check;
        ALTER TABLE patients ADD CONSTRAINT patients_blood_type_check 
        CHECK (blood_type IS NULL OR blood_type = '' OR 
               blood_type IN ('A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-', 'Rh-null'));
    `,
    Down: `
        ALTER TABLE patients DROP CONSTRAINT IF EXISTS patients_blood_type_check;
        ALTER TABLE patients ADD CONSTRAINT patients_blood_type_check 
        CHECK (blood_type IS NULL OR blood_type = '' OR 
               blood_type IN ('A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-'));
    `,
}
```

### Example 3: Data Migration
```go
{
    Version:     4,
    Description: "Migrate patient phone format to international standard",
    Up: `
        -- Add new column
        ALTER TABLE patients ADD COLUMN phone_international VARCHAR(25);
        
        -- Migrate existing data
        UPDATE patients SET phone_international = 
            CASE 
                WHEN country = 'Canada' THEN '+1' || REGEXP_REPLACE(phone, '[^0-9]', '', 'g')
                WHEN country = 'USA' THEN '+1' || REGEXP_REPLACE(phone, '[^0-9]', '', 'g')
                WHEN country = 'Morocco' THEN '+212' || REGEXP_REPLACE(phone, '[^0-9]', '', 'g')
                WHEN country = 'France' THEN '+33' || REGEXP_REPLACE(phone, '[^0-9]', '', 'g')
                ELSE phone
            END;
        
        -- Drop old column and rename new one
        ALTER TABLE patients DROP COLUMN phone;
        ALTER TABLE patients RENAME COLUMN phone_international TO phone;
    `,
    Down: `
        -- This rollback assumes we still have the original phone data
        -- In practice, you might need to store the original format
        ALTER TABLE patients ADD COLUMN phone_domestic VARCHAR(20);
        
        UPDATE patients SET phone_domestic = 
            REGEXP_REPLACE(phone, '^\+[0-9]+', '', 'g');
        
        ALTER TABLE patients DROP COLUMN phone;
        ALTER TABLE patients RENAME COLUMN phone_domestic TO phone;
    `,
}
```

## Future Enhancements

### Planned Improvements
1. **Database seeding** integration with migrations
2. **Cross-service migration** coordination
3. **Migration performance** optimization
4. **Automated testing** of migrations
5. **Blue-green deployment** support

### Advanced Features
1. **Schema versioning** across services
2. **Migration dependencies** between services
3. **Automated rollback** on failure
4. **Migration notifications** and monitoring
5. **Database diff** generation