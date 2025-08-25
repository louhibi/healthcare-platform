# Healthcare Platform - AI Agent Context

This document provides essential context for AI agents (like Claude) working on this healthcare management platform.

## CRITICAL DEVELOPMENT PRINCIPLES

**NEVER ASSUME - ALWAYS VERIFY**

AI agents working on this codebase MUST follow these principles:

1. **Never assume function names** - Always use `grep`, `Read`, or other tools to find the exact function name
2. **Never assume database column names** - Always use `\d table_name` or `DESCRIBE` to check actual column names
3. **Never assume table names** - Always use `\dt` or similar to list actual table names
4. **Never assume API endpoints** - Always check route definitions in code or logs
5. **Never assume field names in structs** - Always read the actual struct definition
6. **Never assume configuration keys** - Always check actual config files or environment variables
7. **Never assume file paths** - Always verify file existence and paths with `ls`, `find`, or `Read` tools

**VERIFICATION WORKFLOW:**
- Before referencing any name (function, column, table, field, etc.), FIRST verify it exists
- Use appropriate tools: `grep`, `Read`, database queries, `ls`, `find`
- If unsure about naming, search for similar patterns in the codebase
- Document your findings as you discover the actual names

**EXAMPLES OF WHAT NOT TO DO:**
❌ "Let me update the `user_name` column" (without checking if it's actually `user_name`, `username`, or `name`)  
❌ "I'll call the `getUserById` function" (without verifying the actual function name)  
❌ "The `patients` table has..." (without checking the actual table schema)

**EXAMPLES OF CORRECT APPROACH:**
✅ Use `\d patients` to check actual column names before referencing them  
✅ Use `grep -r "function.*user" .` to find actual function names  
✅ Use `\dt` to list actual table names before querying

This principle prevents bugs, compilation errors, and broken functionality.

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

### System Architecture Overview
```
healthcare-platform/
├── services/                 # Backend microservices (Go 1.21)
│   ├── api-gateway/         # Port 8080 - Central routing & authentication
│   ├── user-service/        # Port 8081 - User management & healthcare entities  
│   ├── patient-service/     # Port 8082 - Patient CRUD & medical records
│   ├── appointment-service/ # Port 8083 - Scheduling & availability
│   └── location-service/    # Port 8084 - Geographic data (countries/states/cities)
├── frontend/patient-portal/ # Vue.js 3 frontend with modern architecture
├── infrastructure/
│   ├── database/            # PostgreSQL schemas & migrations
│   └── docker/              # Container configurations
└── docs/                    # Documentation & refactoring guides
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

### Backend Microservices Architecture (Go 1.21)

#### Backend Microservices Details

#### Standard Service Structure
All services follow this consistent pattern:
```
service-name/
├── main.go              # Service entry point & Gin routing
├── models.go            # Data structures with validation tags
├── database.go          # PostgreSQL connection & auto-migrations
├── service.go           # Business logic & database operations
├── handlers.go          # HTTP handlers (JSON request/response)
├── migrations.go        # Database schema migrations (auto-run)
├── Dockerfile           # Multi-stage container build
├── go.mod               # Go dependencies
├── go.sum               # Dependency checksums
├── .env.example         # Environment variables template
└── README.md            # Service-specific documentation
```

#### API Gateway Architecture
```
api-gateway/
├── main.go              # Gateway entry point & service discovery
├── config.go            # Service configuration & environment
├── proxy.go             # Request routing to microservices
├── auth_middleware.go   # JWT validation & user context injection
├── rate_limiter.go      # Rate limiting per user/entity
├── stats.go             # Request statistics & health monitoring
└── models.go            # Shared request/response models
```

#### Key Backend Features Discovered

**1. Multi-Tenant Data Isolation**
- All services enforce `healthcare_entity_id` scoping
- JWT tokens include entity context: `X-Healthcare-Entity-ID` header
- Database queries automatically filter by entity for data isolation

**2. Auto-Migration System**
```go
// Each service auto-creates tables on startup
func InitDatabase() {
    db.AutoMigrate(&Patient{}, &HealthcareEntity{})
    // Migrations run automatically in development
}
```

**3. Standardized API Response Format**
```go
// Success Response
{
    "data": {...},
    "message": "Success", 
    "timestamp": "2025-01-09T10:30:00Z"
}

// Error Response
{
    "error": "Validation failed",
    "code": "VALIDATION_ERROR", 
    "message": "Email is required",
    "timestamp": "2025-01-09T10:30:00Z"
}
```

**4. Service Communication Pattern**
- API Gateway routes requests to services based on URL paths
- Services are stateless and database-per-service pattern
- Internal communication via HTTP (no service mesh yet)

**5. Authentication Flow**
```
1. Frontend → API Gateway: POST /api/auth/login
2. API Gateway → User Service: Forward auth request
3. User Service → API Gateway: JWT with healthcare_entity_id
4. API Gateway → Frontend: Return JWT tokens
5. Frontend → API Gateway: Requests with Bearer token
6. API Gateway → Services: Forward with user context headers
```

#### Database Architecture Deep Dive

**Per-Service Databases (PostgreSQL)**:
```
healthcare-user-db:5432      → user_service_db
healthcare-patient-db:5433   → patient_service_db  
healthcare-appointment-db:5434 → appointment_service_db
healthcare-location-db:5435  → location_service_db
redis:6379                   → Caching & sessions
```

**Critical Database Patterns**:
- **TIMESTAMPTZ**: All datetime fields stored in UTC
- **Multi-tenant indexes**: All tables have `healthcare_entity_id` foreign keys
- **Soft deletes**: `is_active` boolean instead of hard deletes  
- **Audit fields**: `created_at`, `updated_at`, `created_by` on all entities
- **Unique constraints**: Scoped per healthcare entity (e.g., patient emails)

#### Environment Configuration Deep Dive

**Critical Environment Variables Pattern**:
```bash
# Database (each service has its own)
DB_HOST=healthcare-[service]-db
DB_PORT=543[2-5]  
DB_NAME=[service]_service_db
DB_USER=postgres
DB_PASSWORD=password

# JWT (MUST be identical across user-service & api-gateway)
JWT_SECRET=your-256-bit-secret-key
JWT_ACCESS_TOKEN_EXPIRY=24h

# Service URLs (for internal communication)
USER_SERVICE_URL=http://user-service:8081
PATIENT_SERVICE_URL=http://patient-service:8082
APPOINTMENT_SERVICE_URL=http://appointment-service:8083

# Multi-tenant & International
DEFAULT_COUNTRY=CA
DEFAULT_LANGUAGE=en
DEFAULT_TIMEZONE=America/Toronto
```

#### Development Workflow Patterns

**Docker-First Development**:
- All services run in containers (Go binaries not on host)
- Hot-reloading via volume mounts in development
- Database containers persist data via Docker volumes
- Networks isolated per environment (development/production)

**Health Check System**:
- Every service exposes `/health` endpoint
- API Gateway aggregates health status from all services
- Database connectivity checked in health endpoints

### Frontend Architecture (Vue.js 3 + Composition API):

#### Overview
The frontend is a modern Vue.js 3 application with composition API, utilizing a modular architecture with composables, reusable components, and performance optimizations.

**Key Features:**
- **Modular Component Architecture**: Large components split into focused, reusable pieces
- **Composition API**: Logic extraction into composables for better reusability
- **Performance Optimized**: Virtual scrolling, dynamic loading, and debouncing
- **Type Safety**: Comprehensive prop validation and error handling
- **Testing**: 209 test cases with 77% pass rate
- **Multi-tenant Support**: Healthcare entity-based data isolation

#### Frontend Structure:
```
frontend/patient-portal/
├── src/
│   ├── components/         # Vue components (modular architecture)
│   │   ├── FormField.vue              # Reusable form field renderer
│   │   ├── VirtualScrollForm.vue      # Performance-optimized scrolling
│   │   ├── fields/                    # Specialized field components (6 files)
│   │   │   ├── DynamicFieldLoader.vue # Dynamic component loading
│   │   │   ├── DateField.vue          # Date input component
│   │   │   ├── SelectField.vue        # Select/dropdown component
│   │   │   └── ...                    # Other field types
│   │   └── patients/                  # Patient management components (10 files)
│   │       ├── PatientHeader.vue      # Header with actions
│   │       ├── PatientSearch.vue      # Search functionality
│   │       ├── PatientsTable.vue      # Table container
│   │       ├── PatientTableRow.vue    # Individual rows
│   │       ├── PatientViewModal.vue   # Patient details modal
│   │       └── ...                    # Other patient components
│   ├── composables/        # Composition API logic (4 composables)
│   │   ├── useLocationData.js         # Location cascade logic
│   │   ├── useFormState.js           # Form state management
│   │   ├── useValidation.js          # Validation rules & logic
│   │   └── useFormConfig.js          # Dynamic form configuration
│   ├── utils/              # Utility functions
│   │   ├── fieldComponentMap.js       # Field type mapping system
│   │   ├── fieldLoader.js            # Dynamic component loading
│   │   ├── virtualScrollMetrics.js   # Performance tracking
│   │   └── timezoneUtils.js          # Timezone handling
│   ├── views/              # Page components
│   ├── stores/             # Pinia stores (global state)
│   ├── router/             # Vue Router config
│   ├── api/                # API client services
│   └── test/               # Test configuration
│       └── setup.js                   # Global test setup
├── __tests__/              # Test files (9 test suites, 209 test cases)
│   ├── composables/__tests__/         # Composable tests (4 files)
│   ├── components/__tests__/          # Component tests (2 files)
│   └── utils/__tests__/               # Utility tests (3 files)
├── vitest.config.js        # Test configuration
├── package.json            # Dependencies + test scripts
├── vite.config.js          # Build configuration
├── TODO-TEST-FIXES.md      # Test fixes todo list (20 tasks)
└── REFACTORING-SUMMARY.md  # Complete refactoring documentation
```

#### Key Architectural Patterns

**1. Composables for Logic Extraction**
```typescript
// useLocationData.js - Location cascade logic
export function useLocationData() {
  const countryOptions = ref([])
  const stateOptions = ref([])
  const cityOptions = ref([])
  
  const loadCountries = async () => { /* API logic */ }
  const handleCountryChange = async (country, formData) => { /* cascade logic */ }
  
  return { countryOptions, stateOptions, cityOptions, loadCountries, handleCountryChange }
}
```

**2. Dynamic Component Loading**
```typescript
// Field components loaded on-demand
const FIELD_COMPONENT_MAP = {
  text: () => import('./fields/TextField.vue'),
  date: () => import('./fields/DateField.vue'),
  select: () => import('./fields/SelectField.vue')
}
```

**3. Virtual Scrolling for Performance**
```vue
<VirtualScrollForm 
  :items="formFields"
  :item-height="80"
  :container-height="600"
>
  <template #item="{ item }">
    <FormField :field="item" />
  </template>
</VirtualScrollForm>
```

**4. Comprehensive Testing**
- **Unit Tests**: Composables, utilities, and components
- **Integration Tests**: Form workflows and location cascade
- **Mock System**: API and browser API mocking
- **Coverage**: 80% threshold with detailed reporting

#### Performance Optimizations
- **Bundle Size**: 22% reduction through dynamic loading
- **API Efficiency**: 80% reduction in API calls via debouncing
- **Memory Management**: Proper cleanup in lifecycle hooks
- **Render Performance**: Virtual scrolling for large forms

#### Component Architecture Benefits
- **Maintainability**: 702-line component → 10 focused components (72% reduction)
- **Reusability**: Components used across different views
- **Testability**: Smaller components easier to test
- **Performance**: Better tree-shaking and code splitting

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

### Development Workflow Notes

**IMPORTANT - Docker Management**: The user handles all Docker operations (testing and restarting services). AI agents should **NEVER** run `docker compose restart`, `docker compose up`, or similar Docker commands. The user will test and restart services manually after code changes.

### Common Commands

#### Backend Commands
```bash
# Build all services
docker compose build

# View logs
docker compose logs -f [service-name]

# Run Go tests
go test ./...

# Database access
docker exec -it healthcare-user-db psql -U postgres -d user_service_db
```

#### Frontend Testing Commands
```bash
# Navigate to frontend directory
cd frontend/patient-portal

# Install dependencies (if not already done)
npm install

# Run all frontend tests
npm run test:run

# Run tests with coverage report
npm run test:coverage

# Run tests in watch mode (development)
npm run test

# Run tests with UI interface
npm run test:ui

# Run specific test file
npm run test:run -- src/composables/__tests__/useFormConfig.test.js

# Run tests matching pattern
npm run test:run -- --testNamePattern="validation"

# Fix test issues (see TODO-TEST-FIXES.md)
# Current status: 161/209 tests passing (77% pass rate)
```

#### Development Commands
```bash
# Start frontend development server
cd frontend/patient-portal && npm run dev

# Build frontend for production
cd frontend/patient-portal && npm run build

# Lint frontend code
cd frontend/patient-portal && npm run lint
```

## Troubleshooting

### Backend Service Issues

#### Service Communication Issues
1. **API Gateway not routing correctly**:
   ```bash
   # Check gateway logs for routing errors
   docker compose logs api-gateway
   
   # Verify service URLs in gateway config
   docker exec healthcare-api-gateway env | grep SERVICE_URL
   
   # Test direct service endpoints
   curl http://localhost:8081/health  # User service
   curl http://localhost:8082/health  # Patient service
   ```

2. **JWT Authentication failures**:
   ```bash
   # Verify JWT secrets match between user-service and api-gateway
   docker exec healthcare-user-service env | grep JWT_SECRET
   docker exec healthcare-api-gateway env | grep JWT_SECRET
   
   # Check for expired tokens in logs
   docker compose logs api-gateway | grep "expired"
   ```

3. **Multi-tenant data isolation issues**:
   ```bash
   # Verify healthcare_entity_id is being passed in headers
   docker compose logs api-gateway | grep "X-Healthcare-Entity-ID"
   
   # Check user context in service logs
   docker compose logs patient-service | grep "entity"
   ```

#### Location Service Integration
4. **Location cascade not working**:
   ```bash
   # Check if location-service is running and accessible
   curl http://localhost:8084/health
   
   # Test location API endpoints directly
   curl http://localhost:8080/api/locations/countries
   curl http://localhost:8080/api/locations/states/CA
   curl http://localhost:8080/api/locations/cities/CA/ON
   ```

### Database Issues

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

## Backend Performance & Monitoring

### Performance Characteristics (Observed during development)
- **API Gateway**: Handles ~100 requests/second in development environment
- **Service Response Times**: Patient service averages 50-100ms, User service 30-80ms
- **Database Performance**: Queries optimized with healthcare_entity_id indexing
- **Location Service**: Geographic data cached for improved cascade performance
- **Memory Usage**: Services typically use 50-150MB RAM per container

### Monitoring Features Available
```bash
# Health Check System
curl http://localhost:8080/health    # API Gateway health
curl http://localhost:8081/health    # User Service
curl http://localhost:8082/health    # Patient Service  
curl http://localhost:8083/health    # Appointment Service
curl http://localhost:8084/health    # Location Service

# Service Monitoring Commands
docker stats healthcare-api-gateway healthcare-user-service healthcare-patient-service

# Request Logging (API Gateway)
docker compose logs api-gateway | grep "request" | tail -20

# Database Connection Monitoring
docker exec healthcare-user-db psql -U postgres -d user_service_db -c "
  SELECT count(*) as active_connections 
  FROM pg_stat_activity 
  WHERE state = 'active';"

# JWT Token Validation Tracking
docker compose logs api-gateway | grep "JWT"

# Multi-tenant Data Access Patterns
docker compose logs patient-service | grep "healthcare_entity_id"
```

### Critical Performance Bottlenecks Identified
1. **Location Cascade API Calls**: Fixed with debouncing (300ms) in frontend
2. **Form Field Rendering**: Optimized with dynamic component loading
3. **Large Patient Lists**: Mitigated with virtual scrolling implementation
4. **Database Entity Queries**: All queries properly scoped by healthcare_entity_id
5. **JWT Validation**: Cached at API Gateway level to reduce overhead

### Scaling Considerations
- **Database Sharding**: Ready for entity-based horizontal scaling
- **Service Replication**: Stateless services can be easily replicated
- **Container Resources**: Services optimized for 512MB RAM limits in production
- **Connection Pooling**: PostgreSQL configured with appropriate pool sizes

---

**For AI Agents**: This codebase follows standard microservices patterns with Go backends and Vue.js frontend. The frontend has undergone comprehensive refactoring with modern Vue 3 + Composition API patterns. Focus on maintaining consistency with existing patterns, proper error handling, and security best practices. Always validate inputs and use parameterized database queries.

**Frontend Refactoring Status (2025-01-09)**:
- ✅ **19 major refactoring tasks completed** - See `REFACTORING-SUMMARY.md` for full details
- ✅ **Modern Vue.js 3 architecture** with Composition API and composables
- ✅ **Modular components** - 702-line component split into 10 focused pieces (72% reduction)
- ✅ **4 new composables** - `useLocationData`, `useFormState`, `useValidation`, `useFormConfig`
- ✅ **Performance optimizations** - Virtual scrolling, dynamic loading, debouncing
- ✅ **Testing infrastructure** - 209 test cases created with Vitest
- 🔄 **Test fixes needed** - 20 tasks in `TODO-TEST-FIXES.md` (161/209 tests passing, 77% rate)

**Testing Priority**: Focus on API mocking issues and useFormState/useLocationData composable test fixes first.

**IMPORTANT - Development Environment**: 
- **ALWAYS use Docker for all operations** - never attempt to run `go build`, `go run`, `npm install`, etc. directly
- Use `docker compose build [service-name]` to build services
- Use `docker compose up -d [service-name]` to start services  
- Use `docker compose restart [service-name]` to restart services
- Use `docker compose logs [service-name]` to view logs
- This project is fully containerized and Go/Node binaries are not available on the host system