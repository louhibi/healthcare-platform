# Environment Setup Guide

This guide explains how to set up your development environment for the Healthcare Platform.

## Environment Files

The project uses environment variables for configuration. Each service and the frontend have their own `.env.example` files that should be copied to `.env` and customized for your environment.

### Quick Setup

```bash
# Copy all environment files
cp .env.example .env
cp services/user-service/.env.example services/user-service/.env
cp services/patient-service/.env.example services/patient-service/.env
cp frontend/patient-portal/.env.example frontend/patient-portal/.env
```

## Environment Files Overview

### Root Environment (`.env`)
The main environment file contains global configuration for all services:

- **Database Configuration**: PostgreSQL connection settings for all services
- **Service URLs**: Inter-service communication endpoints
- **Authentication**: JWT secrets and token expiry settings
- **Security**: CORS, rate limiting, encryption keys
- **Compliance**: HIPAA audit settings, data retention policies
- **Feature Flags**: Enable/disable platform features
- **Monitoring**: Metrics, tracing, health check configuration

### Service-Specific Environments

#### User Service (`.env`)
- Database connection for user data
- JWT configuration for authentication
- Email service settings for notifications
- Healthcare entity defaults
- CORS and rate limiting settings

#### Patient Service (`.env`)
- Database connection for patient data
- Multi-tenant configuration
- International support settings
- File upload configuration
- HIPAA compliance settings
- Patient data validation rules

#### Frontend (`.env`)
- API endpoint URLs
- Application branding
- Feature flags for UI components
- Default user preferences
- Authentication storage settings

## Git Configuration

### `.gitignore` File
The project includes a comprehensive `.gitignore` file that excludes:

- **Sensitive Files**: `.env` files, secrets, certificates
- **Build Artifacts**: Compiled binaries, node_modules, dist folders
- **Database Files**: Local database data, dumps, backups
- **IDE Files**: Editor configurations, workspace files
- **OS Files**: System-generated files (`.DS_Store`, `Thumbs.db`)
- **Logs**: Application logs, debug files
- **Temporary Files**: Cache, temporary uploads, test artifacts

### Files to Commit vs. Ignore

#### ✅ Commit These Files
- `.env.example` files (templates)
- Source code (`.go`, `.js`, `.vue`, etc.)
- Configuration templates
- Documentation (`.md` files)
- Docker configuration (`Dockerfile`, `docker-compose.yml`)
- Package definitions (`go.mod`, `package.json`)

#### ❌ Never Commit These Files
- `.env` files (contain secrets)
- Database files (`*.db`, `pgdata/`)
- Build artifacts (`dist/`, `bin/`)
- Node modules (`node_modules/`)
- IDE configurations (`.vscode/`, `.idea/`)
- Logs (`*.log`, `logs/`)
- Temporary files (`tmp/`, `temp/`)

## Environment Variables Reference

### Database Configuration
```bash
# PostgreSQL settings
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=service_name_db
```

### Authentication
```bash
# JWT configuration (use strong secrets in production)
JWT_SECRET=your-256-bit-secret-key-change-in-production
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=7d
```

### Service Communication
```bash
# Service URLs for internal communication
USER_SERVICE_URL=http://user-service:8081
PATIENT_SERVICE_URL=http://patient-service:8082
APPOINTMENT_SERVICE_URL=http://appointment-service:8083
```

### Security
```bash
# CORS configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080

# Rate limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=100

# Data encryption
DATA_ENCRYPTION_KEY=your-32-character-encryption-key
PHI_ENCRYPTION_ENABLED=true
```

### Healthcare Configuration
```bash
# Multi-country support
SUPPORTED_COUNTRIES=Canada,USA,Morocco,France
DEFAULT_COUNTRY=Canada
DEFAULT_LANGUAGE=en
DEFAULT_TIMEZONE=America/Toronto

# HIPAA compliance
HIPAA_AUDIT_ENABLED=true
DATA_RETENTION_YEARS=7
```

### Development Settings
```bash
# Development environment
ENV=development
LOG_LEVEL=debug
GIN_MODE=debug

# Debug features
ENABLE_DEBUG_ROUTES=true
ENABLE_SQL_LOGGING=false
ENABLE_REQUEST_LOGGING=true
```

## Security Best Practices

### Environment Security
1. **Never commit `.env` files** to version control
2. **Use strong secrets** in production (minimum 256 bits for JWT)
3. **Rotate secrets regularly** in production environments
4. **Use different secrets** for each environment (dev, staging, prod)
5. **Restrict database access** with proper user permissions

### Development Security
1. **Use HTTPS** in production environments
2. **Enable rate limiting** to prevent abuse
3. **Configure CORS** properly for your domain
4. **Enable audit logging** for compliance
5. **Encrypt sensitive data** (PHI) at rest

### Production Considerations
1. **Use environment variables** instead of config files
2. **Store secrets** in secure secret management systems
3. **Enable monitoring** and alerting
4. **Configure backup** and disaster recovery
5. **Regular security audits** and penetration testing

## Environment-Specific Configuration

### Development Environment
```bash
# Relaxed security for development
ENV=development
LOG_LEVEL=debug
GIN_MODE=debug
CORS_ALLOWED_ORIGINS=*
RATE_LIMIT_ENABLED=false
ENABLE_DEBUG_ROUTES=true
```

### Staging Environment
```bash
# Production-like settings with debugging
ENV=staging
LOG_LEVEL=info
GIN_MODE=release
CORS_ALLOWED_ORIGINS=https://staging.healthcare-platform.com
RATE_LIMIT_ENABLED=true
ENABLE_DEBUG_ROUTES=false
```

### Production Environment
```bash
# Strict security settings
ENV=production
LOG_LEVEL=warn
GIN_MODE=release
CORS_ALLOWED_ORIGINS=https://healthcare-platform.com
RATE_LIMIT_ENABLED=true
ENABLE_DEBUG_ROUTES=false
HIPAA_AUDIT_ENABLED=true
PHI_ENCRYPTION_ENABLED=true
```

## Docker Configuration

### Docker Compose Environment
When using Docker Compose, environment variables are automatically loaded from `.env` files:

```yaml
# docker-compose.yml uses .env automatically
services:
  user-service:
    environment:
      - DB_HOST=${USER_DB_HOST}
      - DB_PORT=${USER_DB_PORT}
      - JWT_SECRET=${JWT_SECRET}
```

### Container Environment
Each service container loads its own `.env` file:

```dockerfile
# Dockerfile can use build-time variables
ARG ENV=development
ENV APP_ENV=${ENV}
```

## Troubleshooting

### Common Issues

#### Environment Variables Not Loading
```bash
# Check if .env file exists
ls -la .env

# Verify environment variables are set
printenv | grep DB_HOST

# Check Docker container environment
docker exec container_name printenv
```

#### Database Connection Issues
```bash
# Test database connectivity
docker exec postgres_container psql -U postgres -l

# Check service logs
docker compose logs service_name

# Verify environment variables
docker compose config
```

#### Service Communication Problems
```bash
# Test service endpoints
curl http://localhost:8081/health
curl http://localhost:8082/health

# Check network connectivity
docker network ls
docker network inspect network_name
```

### Getting Help

If you encounter issues with environment setup:

1. **Check the logs**: `docker compose logs service_name`
2. **Verify configuration**: `docker compose config`
3. **Test connectivity**: Use curl to test endpoints
4. **Check documentation**: Review service-specific README files
5. **Ask for help**: Create an issue with error logs and configuration details

## Development Workflow

### Initial Setup
1. Clone the repository
2. Copy all `.env.example` files to `.env`
3. Update database passwords and secrets
4. Start services with `docker compose up`
5. Verify all services are healthy

### Making Changes
1. Update `.env.example` files when adding new variables
2. Document changes in this file
3. Test configuration changes locally
4. Coordinate environment changes with team

### Deployment
1. Update production environment variables
2. Test configuration in staging first
3. Apply database migrations if needed
4. Deploy services in correct order
5. Verify all services are healthy

This setup ensures consistent, secure, and maintainable environment configuration across all development and production environments.