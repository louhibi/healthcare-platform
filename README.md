# Healthcare Management Platform

A comprehensive, production-ready healthcare management platform built with microservices architecture for hospitals, clinics, and medical practices.

## 🏥 Overview

This platform provides essential healthcare management capabilities including patient management, appointment scheduling, user authentication, and role-based access control. Built with modern technologies and designed for scalability, security, and ease of use.

### Key Features

- **Patient Management**: Complete patient profiles with medical history, insurance, and emergency contacts
- **Appointment Scheduling**: Smart scheduling with conflict detection and availability management
- **User Management**: Role-based authentication (Admin, Doctor, Nurse, Staff)
- **API Gateway**: Centralized routing, authentication, and rate limiting
- **Microservices Architecture**: Scalable and maintainable service separation
- **Production Ready**: Docker containerization, health checks, and monitoring

## 🏗️ Architecture

### Microservices

- **API Gateway** (Port 8080): Request routing, authentication, rate limiting
- **User Service** (Port 8081): Authentication, user management, JWT tokens
- **Patient Service** (Port 8082): Patient CRUD operations, search, medical records
- **Appointment Service** (Port 8083): Scheduling, conflict detection, availability

### Technology Stack

**Backend**
- Go 1.21 with Gin framework
- PostgreSQL (separate DB per service)
- JWT authentication
- Redis caching
- RESTful APIs

**Frontend**
- Vue.js 3 with Composition API
- Tailwind CSS for styling
- Pinia for state management
- Axios for API communication

**Infrastructure**
- Docker & Docker Compose
- NGINX load balancing
- Health checks and monitoring
- PgAdmin for database management

## 🚀 Quick Start

### Prerequisites

- Docker with Compose V2 (integrated `docker compose` command)
- Git
- Node.js 18+ (for frontend development)
- Go 1.21+ (for backend development)

### 1. Clone Repository

```bash
git clone <repository-url>
cd healthcare-platform
```

### 2. Environment Setup

Copy environment files:
```bash
# Copy environment templates
cp services/user-service/.env.example services/user-service/.env
cp services/patient-service/.env.example services/patient-service/.env
cp services/appointment-service/.env.example services/appointment-service/.env
cp services/api-gateway/.env.example services/api-gateway/.env
```

### 3. Start All Services

```bash
# Start all services with Docker Compose
docker compose up -d

# View logs
docker compose logs -f

# Check service health
curl http://localhost:8080/health
```

### 4. Access Applications

- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **PgAdmin**: http://localhost:5050 (admin@healthcare.local / admin123)

### 5. Default Credentials

```
Admin User:
- Email: admin@healthcare.local
- Password: admin123

Doctor User:
- Email: doctor@healthcare.local  
- Password: admin123

Nurse User:
- Email: nurse@healthcare.local
- Password: admin123
```

## 📚 Documentation

### API Documentation
- [API Reference](docs/api/README.md) - Complete API documentation
- [Authentication Guide](docs/api/authentication.md) - JWT auth flow
- [Rate Limiting](docs/api/rate-limiting.md) - Rate limit configuration

### Development
- [Development Setup](docs/development/setup.md) - Local development environment
- [Architecture Guide](docs/architecture/README.md) - System architecture overview
- [Database Schema](docs/database/schema.md) - Database design and relationships
- [Contributing](docs/CONTRIBUTING.md) - Contribution guidelines

### Deployment
- [Production Deployment](docs/deployment/production.md) - Production setup guide
- [Docker Guide](docs/deployment/docker.md) - Container deployment
- [Monitoring](docs/deployment/monitoring.md) - Health checks and monitoring

## 🔗 Service Endpoints

### API Gateway (Port 8080)
```
GET  /health              - Gateway health check
GET  /stats               - Gateway statistics
GET  /info                - Gateway information
```

### Authentication
```
POST /api/auth/register   - User registration
POST /api/auth/login      - User login
POST /api/auth/refresh    - Refresh JWT token
```

### Users
```
GET  /api/users/profile   - Get current user profile
PUT  /api/users/profile   - Update user profile
GET  /api/users/          - Get all users (admin only)
```

### Patients
```
GET    /api/patients/        - Get patients (with search/filter)
POST   /api/patients/        - Create new patient
GET    /api/patients/:id     - Get patient by ID
PUT    /api/patients/:id     - Update patient
DELETE /api/patients/:id     - Delete patient
GET    /api/patients/stats   - Patient statistics
```

### Appointments
```
GET    /api/appointments/           - Get appointments (with filters)
POST   /api/appointments/           - Create appointment
GET    /api/appointments/:id        - Get appointment by ID
PUT    /api/appointments/:id        - Update appointment
DELETE /api/appointments/:id        - Delete appointment
PATCH  /api/appointments/:id/status - Update appointment status
```

## 🛠️ Development

### Local Development

1. **Backend Services**:
```bash
# Start databases
docker compose up -d user-db patient-db appointment-db redis

# Run services locally
cd services/user-service && go run .
cd services/patient-service && go run .
cd services/appointment-service && go run .
cd services/api-gateway && go run .
```

2. **Frontend Development**:
```bash
cd frontend/patient-portal
npm install
npm run dev
```

### Testing

```bash
# Run backend tests
cd services/user-service && go test ./...
cd services/patient-service && go test ./...
cd services/appointment-service && go test ./...

# Run frontend tests
cd frontend/patient-portal && npm test
```

### Code Quality

```bash
# Go formatting and linting
gofmt -w .
golangci-lint run

# Frontend linting
cd frontend/patient-portal && npm run lint
```

## 📊 Database Schema

### Users Service
- `users` - User accounts with roles and authentication

### Patients Service  
- `patients` - Patient profiles and medical information

### Appointments Service
- `appointments` - Appointment scheduling and management
- `doctor_schedules` - Doctor availability schedules

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Role-Based Access Control**: Admin, Doctor, Nurse, Staff roles
- **Rate Limiting**: Configurable request rate limiting
- **Input Validation**: Comprehensive input validation
- **SQL Injection Protection**: Parameterized queries
- **CORS Configuration**: Proper cross-origin resource sharing
- **Health Checks**: Service health monitoring

## 🚀 Production Deployment

### Docker Production

```bash
# Build production images
docker compose -f docker-compose.prod.yml build

# Deploy to production
docker compose -f docker-compose.prod.yml up -d
```

### Environment Variables

Key environment variables for production:

```bash
# Security
JWT_SECRET=your-production-secret-key
DB_PASSWORD=secure-database-password

# Services
USER_SERVICE_URL=http://user-service:8081
PATIENT_SERVICE_URL=http://patient-service:8082
APPOINTMENT_SERVICE_URL=http://appointment-service:8083

# Rate Limiting
RATE_LIMIT_RPM=1000
RATE_LIMIT_BURST=100
```

## 📈 Monitoring & Observability

- **Health Checks**: All services expose `/health` endpoints
- **Service Statistics**: Gateway provides request metrics
- **Database Monitoring**: PgAdmin for database management
- **Logging**: Structured logging with request tracing

## 🤝 Contributing

Please read [CONTRIBUTING.md](docs/CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

For support and questions:
- Create an issue in the repository
- Check the [troubleshooting guide](docs/troubleshooting.md)
- Review the [FAQ](docs/FAQ.md)

## 🗺️ Roadmap

### Phase 1 (Current - MVP)
- ✅ User authentication and management
- ✅ Patient management
- ✅ Basic appointment scheduling
- ✅ API Gateway with rate limiting

### Phase 2 (Next)
- 📅 Advanced scheduling features
- 📊 Reporting and analytics
- 📧 Email/SMS notifications
- 🔍 Advanced search capabilities

### Phase 3 (Future)
- 📱 Mobile application
- 🤖 AI-powered features
- 📋 Electronic health records
- 🔗 Third-party integrations

---

**Built with ❤️ for healthcare professionals**