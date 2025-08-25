# Healthcare Platform - MVP Services Architecture

## Core Services

### 1. user-service
**Purpose**: User authentication, authorization, and profile management
- User registration/login
- JWT token management
- Role-based access control (RBAC)
- Password management
- User profile CRUD

### 2. patient-service
**Purpose**: Patient data management and medical records
- Patient profile CRUD operations
- Medical history management
- Patient search and filtering
- Basic patient data validation

### 3. appointment-service
**Purpose**: Appointment scheduling and calendar management
- Appointment CRUD operations
- Schedule conflict detection
- Availability management
- Appointment status tracking

### 4. notification-service
**Purpose**: Communication and alerts (future enhancement for MVP)
- Email notifications
- SMS alerts
- In-app notifications
- Appointment reminders

### 5. api-gateway
**Purpose**: Single entry point for all client requests
- Request routing to appropriate services
- Authentication middleware
- Rate limiting
- Request/response logging

## Support Services

### 6. config-service
**Purpose**: Centralized configuration management
- Environment-specific configs
- Feature flags
- Service discovery endpoints

### 7. audit-service
**Purpose**: Compliance and activity tracking
- User activity logging
- Data access audit trails
- HIPAA compliance logging

## MVP Project Structure
```
healthcare-platform/
├── services/
│   ├── user-service/
│   ├── patient-service/
│   ├── appointment-service/
│   ├── api-gateway/
│   └── audit-service/
├── frontend/
│   └── patient-portal/
├── infrastructure/
│   ├── docker/
│   ├── database/
│   └── monitoring/
├── shared/
│   ├── models/
│   ├── utils/
│   └── middleware/
└── docs/
```

## Service Communication
- **Synchronous**: REST APIs for real-time operations
- **Asynchronous**: Kafka for event-driven updates
- **Database**: Each service has its own database (microservices pattern)

## MVP Phase Services
For the initial MVP, we'll implement:
1. **user-service** (authentication)
2. **patient-service** (patient management)
3. **appointment-service** (scheduling)
4. **api-gateway** (routing)

The notification-service and audit-service will be added in later phases.