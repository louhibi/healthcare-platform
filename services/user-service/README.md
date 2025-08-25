# User Service

Authentication and user management microservice for the Healthcare Management Platform.

## Overview

The User Service handles all user-related operations including authentication, authorization, user profile management, and role-based access control. It serves as the foundation for security across the entire platform.

## Features

- **User Authentication**: JWT-based login and token management
- **User Registration**: New user account creation with validation
- **Role Management**: Admin, Doctor, Nurse, and Staff roles
- **Profile Management**: User profile CRUD operations
- **Token Refresh**: Secure token renewal mechanism
- **Password Security**: bcrypt hashing with appropriate cost

## API Endpoints

### Authentication
```http
POST /api/auth/register    # Register new user
POST /api/auth/login       # User login
POST /api/auth/refresh     # Refresh JWT token
```

### User Management
```http
GET  /api/users/profile    # Get current user profile
PUT  /api/users/profile    # Update user profile
GET  /api/users/           # Get all users (admin only)
```

### Health Check
```http
GET  /health               # Service health status
```

## User Roles

| Role | Description | Permissions |
|------|-------------|-------------|
| `admin` | System administrator | Full system access, user management |
| `doctor` | Medical doctor | Patient access, appointment management |
| `nurse` | Nursing staff | Patient care access, limited appointments |
| `staff` | Administrative staff | Basic patient info, scheduling |

## Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=user_service_db

# JWT Configuration
JWT_SECRET=your-very-secure-secret-key

# Server Configuration
PORT=8081
ENV=development
```

## Database Schema

### Users Table
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

## JWT Token Structure

### Access Token
```json
{
  "user_id": 1,
  "email": "doctor@hospital.com",
  "role": "doctor",
  "exp": 1640991600,
  "iat": 1640990700,
  "type": "access"
}
```

### Refresh Token
```json
{
  "user_id": 1,
  "email": "doctor@hospital.com",
  "exp": 1641595200,
  "iat": 1640990700,
  "type": "refresh"
}
```

## API Examples

### User Registration
```bash
curl -X POST http://localhost:8081/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "doctor@hospital.com",
    "password": "securepassword123",
    "first_name": "Dr. John",
    "last_name": "Smith",
    "role": "doctor"
  }'
```

### User Login
```bash
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "doctor@hospital.com",
    "password": "securepassword123"
  }'
```

### Get User Profile
```bash
curl -X GET http://localhost:8081/api/users/profile \
  -H "Authorization: Bearer <jwt_token>"
```

## Development

### Local Development
```bash
# Clone repository
git clone <repo-url>
cd healthcare-platform/services/user-service

# Install dependencies
go mod download

# Set environment variables
cp .env.example .env

# Start database (from project root)
docker compose up -d user-db

# Run service
go run .
```

### Testing
```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestUserRegistration
```

### Code Structure

```
user-service/
├── main.go              # Service entry point and routing
├── models.go            # User data structures and validation
├── database.go          # Database connection and migrations
├── auth_service.go      # JWT and password handling
├── user_service.go      # User business logic
├── handlers.go          # HTTP request handlers
├── go.mod              # Go dependencies
├── .env.example        # Environment template
└── Dockerfile          # Container build instructions
```

## Security Features

### Password Security
- **bcrypt Hashing**: Industry-standard password hashing
- **Salt Generation**: Automatic salt generation per password
- **Cost Factor**: Configurable computational cost (default: 10)

### JWT Security
- **Short Expiration**: Access tokens expire in 15 minutes
- **Refresh Tokens**: Long-lived tokens for renewal (7 days)
- **Secure Signing**: HMAC SHA-256 signing algorithm
- **Type Validation**: Token type verification (access vs refresh)

### Input Validation
- **Email Format**: RFC-compliant email validation
- **Password Strength**: Minimum 8 characters required
- **Role Validation**: Enum-based role validation
- **XSS Prevention**: Input sanitization

## Error Handling

### Common Error Responses
```json
{
  "error": "Invalid credentials",
  "code": "AUTH_FAILED",
  "message": "Email or password is incorrect",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Error Codes
- `EMAIL_EXISTS` - Email already registered
- `AUTH_FAILED` - Invalid login credentials
- `TOKEN_INVALID` - JWT token validation failed
- `VALIDATION_ERROR` - Input validation failed
- `USER_NOT_FOUND` - User does not exist
- `PERMISSION_DENIED` - Insufficient permissions

## Monitoring

### Health Check
```bash
curl http://localhost:8081/health
```

Response:
```json
{
  "status": "healthy",
  "service": "user-service",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Logging
- **Structured Logging**: JSON format for easy parsing
- **Security Events**: Login attempts, token generation
- **Error Tracking**: Failed authentication attempts
- **Audit Trail**: User creation and profile changes

## Production Considerations

### Database
- **Connection Pooling**: Configure appropriate pool size
- **Index Optimization**: Email and role indexes for performance
- **Backup Strategy**: Regular database backups
- **Migration Management**: Version-controlled schema changes

### Security
- **JWT Secret**: Use strong, unique secret in production
- **HTTPS Only**: Secure transport layer
- **Rate Limiting**: Implement at API Gateway level
- **Password Policy**: Enforce strong password requirements

### Performance
- **Response Caching**: Cache user profiles and roles
- **Database Optimization**: Query optimization and indexing
- **Connection Management**: Efficient database connections
- **Metrics Collection**: Response times and error rates

## Integration

### API Gateway Integration
The User Service is accessed through the API Gateway which handles:
- Request routing to `/api/auth/*` and `/api/users/*`
- JWT token validation for protected endpoints
- Rate limiting and throttling
- Request logging and metrics

### Service Communication
Other services receive user context via HTTP headers:
```
X-User-ID: 123
X-User-Email: doctor@hospital.com
X-User-Role: doctor
```

## Future Enhancements

### Phase 2 Features
- **Multi-factor Authentication (MFA)**: SMS/Email verification
- **Social Login**: OAuth integration (Google, Microsoft)
- **Password Reset**: Email-based password recovery
- **Account Lockout**: Brute force protection
- **Audit Logging**: Enhanced security audit trails

### Phase 3 Features
- **SAML Integration**: Enterprise SSO support
- **Advanced Roles**: Hierarchical role system
- **Session Management**: Advanced session controls
- **Compliance Features**: HIPAA audit requirements

## Troubleshooting

### Common Issues

#### Database Connection Failed
```bash
# Check database status
docker compose ps user-db

# View database logs
docker compose logs user-db

# Test connection
docker exec -it healthcare-user-db psql -U postgres -d user_service_db
```

#### JWT Token Issues
- Verify `JWT_SECRET` is consistent across services
- Check token expiration time
- Validate Authorization header format: `Bearer <token>`

#### Permission Errors
- Verify user role in database
- Check API Gateway routing configuration
- Confirm user is active (`is_active = true`)

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
- Add tests for new functionality
- Update API documentation for changes
- Ensure security best practices

---

**Maintainer**: Healthcare Platform Team  
**Last Updated**: 2024-01-01