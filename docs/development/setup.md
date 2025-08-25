# Development Setup Guide

Complete guide for setting up the Healthcare Platform development environment.

## Prerequisites

### Required Software
- **Docker** 20.10+ with **Compose V2** (integrated `docker compose` command)
- **Git** for version control
- **Go** 1.21+ for backend development
- **Node.js** 18+ and **npm** for frontend development
- **PostgreSQL** 15+ (optional, for direct database access)

### System Requirements
- **OS**: Linux, macOS, or Windows with WSL2
- **RAM**: 8GB minimum, 16GB recommended
- **Storage**: 10GB free space
- **Ports**: 3000, 5050, 5432-5434, 6379, 8080-8083 must be available

## Quick Setup (Recommended)

### 1. Clone Repository
```bash
git clone <repository-url>
cd healthcare-platform
```

### 2. Environment Configuration
```bash
# Copy environment templates
cp services/user-service/.env.example services/user-service/.env
cp services/patient-service/.env.example services/patient-service/.env
cp services/appointment-service/.env.example services/appointment-service/.env
cp services/api-gateway/.env.example services/api-gateway/.env
```

### 3. Start All Services
```bash
# Build and start all services
docker compose up -d

# Check service status
docker compose ps

# View logs
docker compose logs -f
```

### 4. Verify Installation
```bash
# Test API Gateway
curl http://localhost:8080/health

# Test individual services
curl http://localhost:8081/health  # User service
curl http://localhost:8082/health  # Patient service  
curl http://localhost:8083/health  # Appointment service

# Test frontend
curl http://localhost:3000
```

### 5. Access Applications
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **PgAdmin**: http://localhost:5050

## Manual Development Setup

For active development, you may want to run services individually.

### 1. Database Setup
```bash
# Start only databases and Redis
docker compose up -d user-db patient-db appointment-db redis

# Verify databases are running
docker compose ps
```

### 2. Backend Development

#### User Service
```bash
cd services/user-service

# Install dependencies
go mod download

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=user_service_db
export DB_USER=postgres
export DB_PASSWORD=postgres
export JWT_SECRET=your-secret-key
export PORT=8081

# Run service
go run .
```

#### Patient Service
```bash
cd services/patient-service

# Install dependencies
go mod download

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5433
export DB_NAME=patient_service_db
export DB_USER=postgres
export DB_PASSWORD=postgres
export PORT=8082

# Run service
go run .
```

#### Appointment Service
```bash
cd services/appointment-service

# Install dependencies
go mod download

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5434
export DB_NAME=appointment_service_db
export DB_USER=postgres
export DB_PASSWORD=postgres
export PORT=8083

# Run service
go run .
```

#### API Gateway
```bash
cd services/api-gateway

# Install dependencies
go mod download

# Set environment variables
export PORT=8080
export USER_SERVICE_URL=http://localhost:8081
export PATIENT_SERVICE_URL=http://localhost:8082
export APPOINTMENT_SERVICE_URL=http://localhost:8083
export JWT_SECRET=your-secret-key

# Run gateway
go run .
```

### 3. Frontend Development
```bash
cd frontend/patient-portal

# Install dependencies
npm install

# Set environment variables
echo "VITE_API_URL=http://localhost:8080" > .env.local

# Start development server
npm run dev
```

## Development Tools

### Database Management
```bash
# Access databases directly
docker exec -it healthcare-user-db psql -U postgres -d user_service_db
docker exec -it healthcare-patient-db psql -U postgres -d patient_service_db
docker exec -it healthcare-appointment-db psql -U postgres -d appointment_service_db

# Or use PgAdmin
# URL: http://localhost:5050
# Email: admin@healthcare.local
# Password: admin123
```

### Code Quality Tools

#### Go Tools
```bash
# Format code
gofmt -w .

# Run linter (install golangci-lint first)
golangci-lint run

# Run tests
go test ./...

# Generate test coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Frontend Tools
```bash
cd frontend/patient-portal

# Lint JavaScript/Vue
npm run lint

# Run tests
npm test

# Build for production
npm run build
```

### Useful Development Commands

#### Docker Commands
```bash
# Build specific service
docker compose build user-service

# Restart service
docker compose restart user-service

# View service logs
docker compose logs -f user-service

# Execute commands in container
docker compose exec user-service sh

# Remove all containers and volumes
docker compose down -v
```

#### Database Commands
```bash
# Reset specific database
docker compose stop user-db
docker volume rm healthcare-platform_user_db_data
docker compose up -d user-db

# Backup database
docker exec healthcare-user-db pg_dump -U postgres user_service_db > backup.sql

# Restore database
cat backup.sql | docker exec -i healthcare-user-db psql -U postgres -d user_service_db
```

## Environment Variables Reference

### User Service
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=user_service_db
JWT_SECRET=your-very-secure-secret-key
PORT=8081
ENV=development
```

### Patient Service
```bash
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=patient_service_db
PORT=8082
ENV=development
```

### Appointment Service
```bash
DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=appointment_service_db
PORT=8083
ENV=development
```

### API Gateway
```bash
PORT=8080
USER_SERVICE_URL=http://localhost:8081
PATIENT_SERVICE_URL=http://localhost:8082
APPOINTMENT_SERVICE_URL=http://localhost:8083
JWT_SECRET=your-very-secure-secret-key
RATE_LIMIT_RPM=100
RATE_LIMIT_BURST=20
LOG_LEVEL=info
ENV=development
```

### Frontend
```bash
VITE_API_URL=http://localhost:8080
```

## Testing

### Backend Testing
```bash
# Run all tests
cd services/user-service && go test ./...
cd services/patient-service && go test ./...
cd services/appointment-service && go test ./...
cd services/api-gateway && go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestUserRegistration
```

### Integration Testing
```bash
# Start services
docker compose up -d

# Run integration tests
./scripts/run-integration-tests.sh
```

### Frontend Testing
```bash
cd frontend/patient-portal

# Unit tests
npm test

# E2E tests (install Cypress first)
npm run test:e2e
```

## Common Development Issues

### Port Conflicts
```bash
# Check what's using a port
lsof -i :8080

# Kill process using port
kill -9 $(lsof -t -i:8080)
```

### Database Connection Issues
```bash
# Check database status
docker compose ps

# Check database logs
docker compose logs user-db

# Reset database
docker compose stop user-db
docker volume rm healthcare-platform_user_db_data
docker compose up -d user-db
```

### JWT Token Issues
- Ensure `JWT_SECRET` is same in user-service and api-gateway
- Check token expiration time (15 minutes by default)
- Verify Authorization header format: `Bearer <token>`

### CORS Issues
- Check API Gateway CORS configuration
- Verify frontend is making requests to correct URL
- Ensure preflight requests are handled properly

## IDE Configuration

### VS Code
Recommended extensions:
- Go extension
- Vue.js Extension Pack
- Docker extension
- PostgreSQL extension

### GoLand/IntelliJ
- Enable Go modules support
- Configure code style for Go
- Set up run configurations for each service

## Performance Optimization

### Development Performance
```bash
# Use air for hot reloading Go services
go install github.com/cosmtrek/air@latest
cd services/user-service && air

# Use Vite HMR for frontend
cd frontend/patient-portal && npm run dev
```

### Database Performance
- Use connection pooling
- Monitor query performance
- Use database indexes appropriately
- Regular VACUUM and ANALYZE operations

## Next Steps

1. Read [Architecture Documentation](../architecture/README.md)
2. Review [API Documentation](../api/README.md)
3. Check [Contributing Guidelines](../CONTRIBUTING.md)
4. Set up your preferred IDE/editor
5. Start with small feature additions or bug fixes

## Getting Help

- Check [Troubleshooting Guide](../troubleshooting.md)
- Review [FAQ](../FAQ.md)
- Create an issue in the repository
- Join the development team chat