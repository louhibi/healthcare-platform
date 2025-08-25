# Healthcare Platform - AI Agent Context

This document provides essential context for AI agents (like Claude) working on this healthcare management platform.

## Project Overview

**Healthcare Management Platform MVP** - A production-ready multi-tenant microservices platform for hospitals, clinics, and medical practices with international support.

### Core Business Logic
- **Multi-Tenant Architecture**: Entity-based data isolation for hospitals, clinics, and doctor offices
- **International Support**: Multi-country operations (Canada, USA, Morocco, France)
- **Patient Management**: CRUD operations, medical records, insurance tracking with country-specific validation
- **Appointment Scheduling**: Conflict detection, availability management, doctor schedules
- **User Authentication**: JWT-based auth with role-based access control and entity association
- **API Gateway**: Centralized routing, rate limiting, auth middleware

### User Roles & Permissions
- **Admin**: Full system access, user management
- **Doctor**: Patient access, appointment management, medical records
- **Nurse**: Patient care access, limited appointment access
- **Staff**: Basic patient info, appointment scheduling

## Technical Architecture

### Microservices Structure
```
healthcare-platform/
├── services/
│   ├── api-gateway/          # Port 8080 - Request routing & auth
│   ├── user-service/         # Port 8081 - Authentication & users
│   ├── patient-service/      # Port 8082 - Patient management
│   └── appointment-service/  # Port 8083 - Scheduling
├── frontend/patient-portal/  # Vue.js frontend
├── infrastructure/
│   ├── database/            # PostgreSQL schemas & migrations
│   └── docker/              # Container configurations
└── docs/                    # Documentation
```

### Technology Stack
- **Backend**: Go 1.21, Gin framework, PostgreSQL, Redis, JWT
- **Frontend**: Vue.js 3, Tailwind CSS, Pinia, Axios
- **Infrastructure**: Docker, Docker Compose, NGINX

### Database Design
- **Multi-Tenant Architecture**: Healthcare entity-based data isolation
- **Separate databases per service** (microservices pattern)
- **PostgreSQL** with proper indexing and relationships
- **International Data Support**: Country-specific fields, multi-language content
- **Sample data** included for testing (20+ patients across 4 countries)
- **Migrations** handled in each service

## Multi-Tenant & International Features

### Healthcare Entity Management
- **Entity Types**: Hospital, Clinic, Doctor Office
- **Supported Countries**: Canada, USA, Morocco, France
- **Data Isolation**: All patient data scoped to healthcare entity
- **User Association**: Users belong to specific healthcare entities

### International Support
- **Multi-Language**: English, French, Arabic content support
- **Country-Specific Validation**: Postal codes, phone formats, national IDs
- **Currency & Timezone**: Per-country configuration
- **Insurance Systems**: OHIP (Canada), Medicare (USA), AMO/RAMED (Morocco), etc.

### Database Schema (Multi-Tenant)
```sql
-- Healthcare Entities
healthcare_entities (
  id, name, type, country, address, phone, email,
  timezone, currency, language, created_at, updated_at
)

-- Users with Entity Association  
users (
  id, healthcare_entity_id, email, password_hash, role,
  license_number, specialization, preferred_language,
  first_name, last_name, created_at, updated_at
)

-- Patients with International Fields
patients (
  id, healthcare_entity_id, patient_id, first_name, last_name,
  date_of_birth, gender, phone, email, address, city, state,
  postal_code, country, nationality, preferred_language,
  marital_status, occupation, insurance, policy_number,
  insurance_provider, national_id, emergency_contact_*,
  medical_history, allergies, medications, blood_type,
  is_active, created_at, updated_at, created_by
)
```

### JWT Token Structure (Multi-Tenant)
```json
{
  "user_id": 1,
  "healthcare_entity_id": 1,
  "email": "user@example.com",
  "role": "doctor",
  "exp": 1234567890
}
```

## Key Files & Patterns

### Backend Service Structure (All services follow this pattern):
```
service-name/
├── main.go           # Service entry point & routing
├── models.go         # Data structures & validation
├── database.go       # DB connection & migrations  
├── service.go        # Business logic
├── handlers.go       # HTTP handlers
├── Dockerfile        # Container build
├── go.mod           # Dependencies
└── .env.example     # Environment template
```

### API Gateway Specific:
```
api-gateway/
├── main.go              # Gateway entry point
├── config.go           # Service configuration
├── proxy.go            # Request proxying logic
├── auth_middleware.go  # JWT authentication
├── rate_limiter.go     # Rate limiting
└── stats.go            # Statistics collection
```

### Frontend Structure:
```
frontend/patient-portal/
├── src/
│   ├── components/     # Vue components
│   ├── views/         # Page components
│   ├── stores/        # Pinia stores
│   ├── router/        # Vue Router config
│   └── api/           # API client
├── package.json       # Dependencies
└── vite.config.js     # Build configuration
```

## Development Guidelines

### Code Standards
- **Go**: Standard Go formatting, error handling patterns
- **Vue.js**: Composition API, TypeScript-ready
- **Database**: Parameterized queries, proper indexing
- **Security**: JWT tokens, input validation, CORS

### Common Patterns
1. **Error Handling**: Consistent JSON error responses
2. **Validation**: Struct tags for input validation
3. **Authentication**: JWT middleware in gateway, user headers in services
4. **Database**: Service-specific databases with migrations
5. **Logging**: Structured logging with request context

### Environment Configuration

The platform uses comprehensive environment configuration across all services with `.env.example` files providing templates for all required variables.

#### Environment File Structure
```
healthcare-platform/
├── .env.example                                    # Docker Compose & global config
├── services/
│   ├── api-gateway/.env.example                   # Gateway-specific config
│   ├── user-service/.env.example                  # User service config
│   ├── patient-service/.env.example               # Patient service config
│   └── appointment-service/.env.example           # Appointment service config
└── frontend/patient-portal/.env.example           # Frontend config
```

#### Critical Configuration Categories

1. **Database Configuration**
   - Each service connects to its own PostgreSQL database
   - Connection strings, credentials, and pool settings
   - **Example**: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`

2. **Authentication & Security**
   - **JWT_SECRET**: Must be identical across user-service and api-gateway
   - Token expiry settings, CORS origins, rate limiting
   - **Example**: `JWT_SECRET`, `JWT_ACCESS_TOKEN_EXPIRY`, `CORS_ALLOWED_ORIGINS`

3. **Service Communication**
   - Internal service URLs for microservice communication  
   - Timeout settings, health check intervals
   - **Example**: `USER_SERVICE_URL`, `PATIENT_SERVICE_URL`, `APPOINTMENT_SERVICE_URL`

4. **Healthcare Platform Specifics**
   - Default country, language, timezone for new entities
   - Multi-tenant and international features
   - **Example**: `DEFAULT_COUNTRY`, `DEFAULT_LANGUAGE`, `DEFAULT_TIMEZONE`

5. **Timezone Management (CRITICAL)**
   - **NO FALLBACK TIMEZONES ALLOWED** - Entity timezone must be properly configured
   - Frontend validation enforces strict timezone requirements
   - **Frontend**: `VITE_TIMEZONE_STRICT_MODE=true` (no `VITE_FALLBACK_TIMEZONE`)
   - **Backend**: All datetime storage in UTC, display in entity timezone

#### Frontend Environment Variables (Vue.js with Vite)
All frontend variables are prefixed with `VITE_` for Vite bundler:
- **API Configuration**: `VITE_API_URL`, `VITE_API_TIMEOUT`
- **Feature Flags**: `VITE_FEATURE_APPOINTMENTS`, `VITE_FEATURE_ADMIN_PANEL`
- **Validation**: `VITE_CONFIG_VALIDATION_ENABLED=true` (required)
- **Timezone**: `VITE_TIMEZONE_STRICT_MODE=true` (no fallbacks)

#### Environment Setup Process
1. **Copy template files**: `cp .env.example .env` in each directory
2. **Update credentials**: Change default passwords, JWT secrets
3. **Verify service URLs**: Ensure service-to-service communication works
4. **Test timezone config**: Validate healthcare entity timezone settings
5. **Run config validation**: Frontend validates configuration on startup

#### Production Considerations
- **Secrets Management**: Use container orchestration secrets, not .env files
- **Database URLs**: Use connection pooling and read replicas
- **JWT Secrets**: Generate cryptographically secure 256-bit keys
- **HTTPS**: Enable TLS for all services in production
- **Rate Limiting**: Adjust limits based on expected traffic

## API Conventions

### Authentication Flow (Multi-Tenant)
1. POST `/api/auth/login` → JWT tokens with healthcare_entity_id
2. Include `Authorization: Bearer <token>` in requests
3. Gateway validates & forwards user context headers to services:
   - `X-User-ID`: User identifier
   - `X-User-Email`: User email address
   - `X-User-Role`: User role (admin, doctor, nurse, staff)
   - `X-Healthcare-Entity-ID`: Healthcare entity for data isolation

### Request/Response Format
```json
// Success Response
{
  "data": {...},
  "message": "Success",
  "timestamp": "2024-01-01T00:00:00Z"
}

// Error Response  
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "message": "Detailed message",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Pagination
```json
{
  "patients": [...],
  "total_count": 100,
  "limit": 10,
  "offset": 0,
  "has_more": true
}
```

## Common Tasks & Solutions

### Adding New API Endpoint
1. Add route to service's `main.go`
2. Create handler in `handlers.go`
3. Add business logic to `service.go`
4. Update API Gateway route config if needed
5. Add validation to models if needed

### Database Changes
1. Update models in `models.go`
2. Add migration in `database.go`
3. Update service methods in `service.go`
4. Test with sample data

### Frontend Feature
1. Create/update Vue component
2. Add API calls in stores/api
3. Update routing if needed
4. Style with Tailwind classes

## Testing & Development

### Authentication Info
**Important**: All sample users have the password `admin123`

Sample users by entity:
- Entity 1 (Toronto General): doctor@healthcare.local, john.smith@test.com
- Entity 3 (Johns Hopkins): dr.sarah.johnson@jh.edu

### Local Development Setup (Complete Guide for New Developers)

#### Prerequisites
- Docker and Docker Compose installed
- Go 1.21+ installed
- Node.js 18+ and npm installed

#### 1. Initial Setup
```bash
# Clone the repository
git clone <repository-url>
cd healthcare-platform

# Start all databases (PostgreSQL instances + Redis)
docker compose up -d user-db patient-db appointment-db redis

# Wait for databases to be ready (about 30 seconds)
docker compose logs -f user-db patient-db appointment-db
```

#### 2. Database Schema & Sample Data Setup
The database schemas are automatically created when services start via built-in migrations. To add sample data:

```bash
# Apply sample healthcare entities data
docker cp infrastructure/database/sample-data.sql healthcare-user-db:/tmp/
docker exec healthcare-user-db psql -U postgres -d user_service_db -f /tmp/sample-data.sql

# Apply sample patients data  
docker cp infrastructure/database/sample-patients.sql healthcare-patient-db:/tmp/
docker exec healthcare-patient-db psql -U postgres -d patient_service_db -f /tmp/sample-patients.sql

# Verify data was loaded
docker exec healthcare-user-db psql -U postgres -d user_service_db -c "SELECT COUNT(*) FROM healthcare_entities;"
docker exec healthcare-patient-db psql -U postgres -d patient_service_db -c "SELECT COUNT(*) FROM patients;"
```

#### 3. Start Services
```bash
# Option A: Run with Docker (Recommended for full setup)
docker compose up -d

# Option B: Run services locally (for development)
# Start each service in separate terminals:
cd services/user-service && go run .        # Port 8081
cd services/patient-service && go run .     # Port 8082  
cd services/appointment-service && go run . # Port 8083
cd services/api-gateway && go run .         # Port 8080

# Start frontend
cd frontend/patient-portal && npm install && npm run dev  # Port 3000
```

#### 4. Verify Setup
```bash
# Test API Gateway health
curl http://localhost:8080/health

# Test individual services
curl http://localhost:8081/health  # User service
curl http://localhost:8082/health  # Patient service
curl http://localhost:8083/health  # Appointment service

# Test frontend
open http://localhost:3000
```

#### 5. Sample Data Overview
After setup, you'll have:
- **16 Healthcare entities** across 4 countries (Canada, USA, Morocco, France)
  - 4 hospitals, 8 clinics, 4 doctor offices
  - Each with country-specific address formats, phone numbers, currencies
- **20+ Sample patients** with international data
  - Canadian patients with OHIP insurance and SIN numbers
  - US patients with Medicare/Blue Cross and SSN
  - Moroccan patients with AMO/RAMED insurance and CIN
  - French patients with European health insurance
- **Multi-language content** (English, French, Arabic)
- **Country-specific validation** patterns
- **Realistic medical data** (blood types, allergies, medications, conditions)

#### Sample Healthcare Entities
```bash
# View sample entities by country
docker exec healthcare-user-db psql -U postgres -d user_service_db -c "
SELECT name, type, country, city FROM healthcare_entities ORDER BY country, type;
"

# View sample patients by entity
docker exec healthcare-patient-db psql -U postgres -d patient_service_db -c "
SELECT patient_id, first_name, last_name, country, healthcare_entity_id 
FROM patients ORDER BY healthcare_entity_id, id;
"
```

### Common Commands
```bash
# Build all services
docker compose build

# View logs
docker compose logs -f [service-name]

# Run tests
go test ./...

# Database access
docker exec -it healthcare-user-db psql -U postgres -d user_service_db
```

## Troubleshooting

### Common Issues

#### Database Setup Issues
1. **Empty databases after startup**: 
   ```bash
   # Check if migrations ran
   docker exec healthcare-user-db psql -U postgres -d user_service_db -c "\dt"
   # Should show healthcare_entities and users tables
   
   # If no tables, restart services to trigger migrations
   docker compose restart user-service patient-service
   ```

2. **Sample data not loading**:
   ```bash
   # Ensure files are copied correctly
   docker cp infrastructure/database/sample-data.sql healthcare-user-db:/tmp/
   docker exec healthcare-user-db psql -U postgres -d user_service_db -f /tmp/sample-data.sql
   
   # Check for SQL errors in logs
   docker compose logs user-service patient-service
   ```

3. **Database connection failures**:
   ```bash
   # Check database containers are healthy
   docker compose ps
   
   # Check database logs
   docker compose logs user-db patient-db appointment-db
   
   # Test direct connection
   docker exec -it healthcare-user-db psql -U postgres -d user_service_db
   ```

#### Service Issues
4. **Port conflicts**: Check if ports 8080-8083, 3000, 5432-5434 are available
5. **CORS errors**: Check API Gateway CORS configuration
6. **Auth failures**: Verify JWT secret consistency across services
7. **Service communication**: Check service URLs in environment config

#### Multi-Tenant Data Issues
8. **Patients not showing**: Ensure healthcare_entity_id is properly set in JWT tokens
9. **Permission errors**: Check user role and entity association in database

### Health Checks
- Gateway: `curl http://localhost:8080/health`
- Individual services: `curl http://localhost:808[1-3]/health`
- Frontend: `curl http://localhost:3000`

## Timezone Management Strategy

### CRITICAL: Universal Timezone Handling
The healthcare platform operates internationally with healthcare entities in different timezones. **ALL** date/time handling MUST follow this pattern:

#### Core Principles
1. **Storage**: Always store in UTC in PostgreSQL `TIMESTAMPTZ` columns
2. **API**: Always send/receive ISO 8601 UTC format (`2025-01-15T14:30:00Z`)
3. **Display**: Always show in healthcare entity's local timezone
4. **Input**: Convert user input from entity timezone to UTC before storage

#### Implementation Pattern
```javascript
// Frontend: Always convert display times to entity timezone
const displayTime = convertUTCToEntityTime(utcDateTime, entityTimezone)
const displayDate = convertUTCToEntityDate(utcDateTime, entityTimezone)

// Frontend: Always convert user input to UTC before API calls
const utcDateTime = convertEntityTimeToUTC(userInput, entityTimezone)
```

#### Database Schema
- All datetime columns use `TIMESTAMPTZ` (timestamp with timezone)
- All timestamps stored as UTC (`2025-01-15 14:30:00+00`)
- No separate date/time columns - use full UTC timestamps only

#### API Response Format
```json
{
  "start_datetime": "2025-01-15T14:30:00Z",
  "end_datetime": "2025-01-15T16:30:00Z",
  "entity_timezone": "America/Toronto"
}
```

#### Healthcare Entity Timezone Mapping
```sql
-- Each entity has its timezone
healthcare_entities (
  id, name, country, timezone, -- e.g., "America/Toronto", "Europe/Paris"
  created_at, updated_at
)
```

#### Frontend Timezone Utils (Required)
```javascript
// /utils/timezoneUtils.js - MUST BE USED BY ALL COMPONENTS
export function convertUTCToEntityTime(utcDateTime, entityTimezone)
export function convertUTCToEntityDate(utcDateTime, entityTimezone)  
export function convertEntityTimeToUTC(localDateTime, entityTimezone)
export function formatEntityDateTime(utcDateTime, entityTimezone, options)
```

#### Backend Timezone Service (Required)
```go
// All services MUST use this for timezone conversions
type TimezoneConverter struct {
    EntityID int
    Timezone string
}

func (tc *TimezoneConverter) ConvertUTCToEntity(utcTime time.Time) time.Time
func (tc *TimezoneConverter) ConvertEntityToUTC(entityTime time.Time) time.Time
```

#### Example Use Cases
- **Doctor Availability**: User sees "9:00 AM - 5:00 PM EST" but stored as UTC
- **Appointments**: Patient books "2:30 PM EST" but stored/queried as UTC
- **Reports**: All times shown in entity timezone, filtered by UTC ranges
- **Multi-Entity**: Admin viewing multiple entities sees each in their timezone

#### NEVER Do This
❌ Store local times without timezone info  
❌ Mix UTC and local times in same database  
❌ Show UTC times to end users  
❌ Hardcode timezones in frontend components  
❌ Convert timezones in SQL queries  

#### Always Do This
✅ Store UTC, display in entity timezone  
✅ Use TimezoneConverter service in backends  
✅ Use timezone utils in frontend components  
✅ Include entity timezone info in API responses  
✅ Test with multiple timezones during development  

## Security Considerations

### IMPORTANT Security Measures
- **JWT Secret**: Must be same across user-service and api-gateway
- **Input Validation**: All inputs validated with struct tags
- **SQL Injection**: Only parameterized queries used
- **CORS**: Configured for frontend domain
- **Rate Limiting**: Implemented in gateway
- **Password Hashing**: bcrypt with appropriate cost

### HIPAA Compliance Notes
- Patient data encrypted at rest (database level)
- Access logging implemented
- Role-based access control enforced
- Audit trails in database triggers

## Future Enhancements

### Planned Features
- Kafka for event streaming
- Advanced scheduling algorithms
- Email/SMS notifications
- Reporting and analytics
- Mobile application
- Third-party integrations

### Architecture Improvements
- Service mesh (Istio)
- Message queues
- Distributed tracing
- Advanced monitoring
- Auto-scaling

---

**For AI Agents**: This codebase follows standard microservices patterns with Go backends and Vue.js frontend. Focus on maintaining consistency with existing patterns, proper error handling, and security best practices. Always validate inputs and use parameterized database queries.

**IMPORTANT - Development Environment**: 
- **ALWAYS use Docker for all operations** - never attempt to run `go build`, `go run`, `npm install`, etc. directly
- Use `docker compose build [service-name]` to build services
- Use `docker compose up -d [service-name]` to start services  
- Use `docker compose restart [service-name]` to restart services
- Use `docker compose logs [service-name]` to view logs
- This project is fully containerized and Go/Node binaries are not available on the host system