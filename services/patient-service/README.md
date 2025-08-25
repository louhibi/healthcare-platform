# Patient Service

Patient data management microservice for the Healthcare Management Platform.

## Overview

The Patient Service manages all patient-related data including personal information, medical records, insurance details, and emergency contacts. It provides comprehensive CRUD operations with advanced search and filtering capabilities.

## Features

- **Patient Management**: Complete patient profile CRUD operations
- **Medical Records**: Medical history, allergies, and current medications
- **Insurance Tracking**: Insurance provider and policy information
- **Emergency Contacts**: Emergency contact management
- **Advanced Search**: Multi-field search with filtering
- **Data Validation**: Comprehensive input validation and sanitization
- **Age Calculation**: Automatic age calculation from date of birth

## API Endpoints

### Patient Operations
```http
GET    /api/patients/        # Get patients with search/filter
POST   /api/patients/        # Create new patient
GET    /api/patients/:id     # Get patient by ID
PUT    /api/patients/:id     # Update patient
DELETE /api/patients/:id     # Delete patient (soft delete)
GET    /api/patients/stats   # Patient statistics
```

### Health Check
```http
GET    /health               # Service health status
```

## Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=patient_service_db

# Server Configuration
PORT=8082
ENV=development
```

## Database Schema

### Patients Table
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

### Indexes
```sql
CREATE INDEX idx_patients_name ON patients(first_name, last_name);
CREATE INDEX idx_patients_email ON patients(email);
CREATE INDEX idx_patients_phone ON patients(phone);
CREATE INDEX idx_patients_dob ON patients(date_of_birth);
CREATE INDEX idx_patients_active ON patients(is_active);
```

## Patient Data Model

### Core Information
- **Personal Details**: Name, date of birth, gender
- **Contact Information**: Phone, email, address
- **Demographics**: City, state, zip code

### Medical Information
- **Medical History**: Free-text medical background
- **Allergies**: Known allergies and reactions
- **Current Medications**: Active medications and dosages

### Insurance & Emergency
- **Insurance**: Provider name and policy number
- **Emergency Contact**: Name, phone, relationship

## API Examples

### Create Patient
```bash
curl -X POST http://localhost:8082/api/patients/ \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "first_name": "Alice",
    "last_name": "Johnson",
    "date_of_birth": "1985-03-15T00:00:00Z",
    "gender": "female",
    "phone": "555-0101",
    "email": "alice.johnson@email.com",
    "address": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "zip_code": "12345",
    "insurance": "Blue Cross Blue Shield",
    "policy_number": "BC123456789",
    "emergency_contact": {
      "name": "Bob Johnson",
      "phone": "555-0102",
      "relationship": "spouse"
    },
    "medical_history": "No significant medical history",
    "allergies": "None known",
    "medications": "None"
  }'
```

### Search Patients
```bash
# Search by name
curl "http://localhost:8082/api/patients/?q=alice&limit=10&offset=0" \
  -H "X-User-ID: 1"

# Filter by gender
curl "http://localhost:8082/api/patients/?gender=female&limit=10" \
  -H "X-User-ID: 1"

# Filter by insurance
curl "http://localhost:8082/api/patients/?insurance=blue%20cross" \
  -H "X-User-ID: 1"
```

### Get Patient by ID
```bash
curl http://localhost:8082/api/patients/1 \
  -H "X-User-ID: 1"
```

### Update Patient
```bash
curl -X PUT http://localhost:8082/api/patients/1 \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "first_name": "Alice",
    "last_name": "Smith",
    "date_of_birth": "1985-03-15T00:00:00Z",
    "gender": "female",
    "phone": "555-0101",
    "email": "alice.smith@email.com",
    "medical_history": "Updated medical history"
  }'
```

## Development

### Local Development
```bash
# Clone repository
git clone <repo-url>
cd healthcare-platform/services/patient-service

# Install dependencies
go mod download

# Set environment variables
cp .env.example .env

# Start database (from project root)
docker compose up -d patient-db

# Run service
go run .
```

### Testing
```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Test specific functionality
go test -run TestPatientCreation
```

### Code Structure

```
patient-service/
├── main.go              # Service entry point and routing
├── models.go            # Patient data structures and validation
├── database.go          # Database connection and migrations
├── patient_service.go   # Patient business logic
├── handlers.go          # HTTP request handlers
├── go.mod              # Go dependencies
├── .env.example        # Environment template
└── Dockerfile          # Container build instructions
```

## Search and Filtering

### Search Capabilities
The service supports advanced search across multiple fields:

```go
// Search query searches across:
// - First name (case-insensitive)
// - Last name (case-insensitive)
// - Email (case-insensitive)
// - Phone number
```

### Filter Options
- **Gender**: Filter by male, female, or other
- **Insurance**: Partial match on insurance provider
- **Date Range**: Filter by date of birth range (future enhancement)
- **Active Status**: Filter active/inactive patients

### Pagination
```json
{
  "patients": [...],
  "total_count": 150,
  "limit": 10,
  "offset": 0,
  "has_more": true
}
```

## Authentication & Authorization

### User Context
The service receives user information via HTTP headers from the API Gateway:
```
X-User-ID: 123
X-User-Email: doctor@hospital.com
X-User-Role: doctor
```

### Access Control
- **All Operations**: Require authentication (via API Gateway)
- **Data Isolation**: Patients created by user are tracked (`created_by` field)
- **Role-Based Access**: All authenticated roles can access patients
- **Future Enhancement**: Fine-grained permissions per role

## Data Validation

### Required Fields
- First name, last name
- Date of birth
- Gender (male, female, other)
- Phone number
- Email address

### Validation Rules
```go
type PatientRequest struct {
    FirstName   string `validate:"required"`
    LastName    string `validate:"required"`
    DateOfBirth time.Time `validate:"required"`
    Gender      string `validate:"required,oneof=male female other"`
    Phone       string `validate:"required"`
    Email       string `validate:"required,email"`
}
```

### Business Rules
- **Email Uniqueness**: Email must be unique across all patients
- **Age Calculation**: Automatically calculated from date of birth
- **Soft Delete**: Patients are marked inactive, not physically deleted

## Error Handling

### Common Error Responses
```json
{
  "error": "Patient not found",
  "code": "PATIENT_NOT_FOUND",
  "message": "No patient found with ID 123",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Error Codes
- `PATIENT_NOT_FOUND` - Patient does not exist
- `EMAIL_EXISTS` - Email already used by another patient
- `VALIDATION_ERROR` - Input validation failed
- `UNAUTHORIZED` - Authentication required
- `PERMISSION_DENIED` - Insufficient permissions

## Performance Optimizations

### Database Optimization
- **Strategic Indexing**: Indexes on search fields (name, email, phone)
- **Query Optimization**: Efficient search queries with proper WHERE clauses
- **Pagination**: Limit/offset to handle large datasets
- **Connection Pooling**: Efficient database connection management

### Response Optimization
- **Field Selection**: Only return necessary fields
- **Computed Fields**: Age calculated at response time
- **Caching**: Future enhancement for frequently accessed data

## Monitoring

### Health Check
```bash
curl http://localhost:8082/health
```

Response:
```json
{
  "status": "healthy",
  "service": "patient-service",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Statistics Endpoint
```bash
curl http://localhost:8082/api/patients/stats \
  -H "X-User-ID: 1"
```

Response:
```json
{
  "total_patients": 1250
}
```

### Logging
- **CRUD Operations**: All create, update, delete operations logged
- **Search Queries**: Search parameters and result counts
- **Error Tracking**: Failed operations and validation errors
- **Performance Metrics**: Query execution times

## Security Features

### Data Protection
- **Input Sanitization**: All inputs validated and sanitized
- **SQL Injection Prevention**: Parameterized queries only
- **XSS Prevention**: Output encoding for web display
- **Data Encryption**: Database-level encryption (production)

### HIPAA Compliance
- **Audit Logging**: All data access logged with user context
- **Access Controls**: Authentication required for all operations
- **Data Retention**: Configurable retention policies
- **Secure Transmission**: HTTPS in production

## Integration

### API Gateway Integration
- **Request Routing**: All requests routed through API Gateway
- **Authentication**: JWT validation handled by gateway
- **User Context**: User information forwarded via headers
- **Rate Limiting**: Implemented at gateway level

### Future Integrations
- **Document Service**: Medical document attachments
- **Audit Service**: Enhanced audit logging
- **Notification Service**: Patient communication
- **Analytics Service**: Patient demographics and insights

## Production Considerations

### Scalability
- **Horizontal Scaling**: Multiple service instances
- **Database Scaling**: Read replicas for query optimization
- **Caching Layer**: Redis for frequently accessed data
- **Load Balancing**: API Gateway handles load distribution

### Backup & Recovery
- **Automated Backups**: Daily database backups
- **Point-in-Time Recovery**: Transaction log backups
- **Data Archival**: Long-term storage for historical data
- **Disaster Recovery**: Cross-region replication

## Future Enhancements

### Phase 2 Features
- **Document Management**: Medical document uploads
- **Advanced Search**: Full-text search capabilities
- **Data Export**: CSV/PDF export functionality
- **Bulk Operations**: Bulk patient imports/updates

### Phase 3 Features
- **Medical Imaging**: Integration with DICOM systems
- **HL7 Integration**: Healthcare data exchange standards
- **Analytics**: Patient demographics and trends
- **Mobile API**: Mobile-optimized endpoints

## Troubleshooting

### Common Issues

#### Database Connection Failed
```bash
# Check database status
docker compose ps patient-db

# View database logs
docker compose logs patient-db

# Test connection
docker exec -it healthcare-patient-db psql -U postgres -d patient_service_db
```

#### Search Performance Issues
- Check database indexes: `EXPLAIN ANALYZE` on slow queries
- Monitor query execution times
- Consider adding additional indexes for common search patterns

#### Validation Errors
- Verify input format matches validation rules
- Check required fields are provided
- Ensure email format is valid

### Debug Mode
```bash
# Enable debug logging
export LOG_LEVEL=debug
go run .
```

## Contributing

Please read the [Contributing Guidelines](../../docs/CONTRIBUTING.md) before making changes to this service.

### Key Guidelines
- Follow Go formatting standards (`gofmt`)
- Add tests for new search/filter functionality
- Update API documentation for changes
- Ensure HIPAA compliance for patient data

---

**Maintainer**: Healthcare Platform Team  
**Last Updated**: 2024-01-01