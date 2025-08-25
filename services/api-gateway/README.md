# API Gateway

Central API gateway for the Healthcare Management Platform, providing request routing, authentication, rate limiting, and monitoring.

## Overview

The API Gateway serves as the single entry point for all client requests, handling authentication, authorization, request routing, rate limiting, and providing comprehensive monitoring and statistics for the healthcare platform.

## Features

- **Request Routing**: Intelligent routing to appropriate microservices
- **JWT Authentication**: Token validation and user context forwarding
- **Rate Limiting**: Token bucket algorithm for request throttling
- **CORS Handling**: Cross-origin resource sharing configuration
- **Health Monitoring**: Comprehensive health checks for all services
- **Statistics Collection**: Request metrics and performance monitoring
- **Error Handling**: Standardized error responses across services
- **Load Balancing**: Request distribution with retry logic

## API Endpoints

### Gateway Management
```http
GET  /health               # Gateway and service health status
GET  /stats                # Request statistics and metrics
GET  /info                 # Gateway configuration information
```

### Routed Endpoints
```http
# Authentication (User Service)
POST /api/auth/register    # User registration
POST /api/auth/login       # User login
POST /api/auth/refresh     # Token refresh

# User Management (User Service)
GET  /api/users/profile    # Get user profile
PUT  /api/users/profile    # Update user profile
GET  /api/users/           # Get all users (admin)

# Patient Management (Patient Service)
GET  /api/patients/        # Get patients
POST /api/patients/        # Create patient
GET  /api/patients/:id     # Get patient by ID
PUT  /api/patients/:id     # Update patient
DELETE /api/patients/:id   # Delete patient

# Appointment Management (Appointment Service)
GET  /api/appointments/    # Get appointments
POST /api/appointments/    # Create appointment
GET  /api/appointments/:id # Get appointment by ID
PUT  /api/appointments/:id # Update appointment
DELETE /api/appointments/:id # Delete appointment

# Locations (Location Service)
GET /api/locations/countries
GET /api/locations/countries/:code/cities
```

## Environment Variables

```bash
# Gateway Configuration
PORT=8080

# Service URLs
USER_SERVICE_URL=http://localhost:8081
PATIENT_SERVICE_URL=http://localhost:8082
APPOINTMENT_SERVICE_URL=http://localhost:8083

# Service Timeouts (seconds)
USER_SERVICE_TIMEOUT=30
PATIENT_SERVICE_TIMEOUT=30
APPOINTMENT_SERVICE_TIMEOUT=30
GATEWAY_TIMEOUT=30

# JWT Configuration
JWT_SECRET=your-very-secure-secret-key

# Rate Limiting
RATE_LIMIT_RPM=100      # Requests per minute
RATE_LIMIT_BURST=20     # Burst capacity

# Logging
LOG_LEVEL=info
ENV=development
```

## Architecture

### Request Flow
```
Client Request
      ↓
  API Gateway
      ↓
┌─────────────┐
│ Middleware  │
│ - CORS      │
│ - Logging   │
│ - Rate Limit│
└─────────────┘
      ↓
┌─────────────┐
│ Auth Check  │
│ (if required)│
└─────────────┘
      ↓
┌─────────────┐
│ Route Match │
│ & Forward   │
└─────────────┘
      ↓
  Backend Service
      ↓
  Response to Client
```

### Service Configuration
```go
type ServiceConfig struct {
    Name    string `json:"name"`
    BaseURL string `json:"base_url"`
    Timeout int    `json:"timeout"`
}

type RouteConfig struct {
    Path         string   `json:"path"`
    Service      string   `json:"service"`
    StripPrefix  bool     `json:"strip_prefix"`
    AuthRequired bool     `json:"auth_required"`
    RolesAllowed []string `json:"roles_allowed"`
}
```

## Authentication & Authorization

### JWT Token Validation
The gateway validates JWT tokens and forwards user context to backend services:

```go
// Headers forwarded to services
X-User-ID: 123
X-User-Email: doctor@hospital.com
X-User-Role: doctor
```

### Route Protection
```go
// Route configuration example
{
    Path:        "/api/patients/*",
    Service:     "patient-service",
    AuthRequired: true,
    RolesAllowed: ["admin", "doctor", "nurse", "staff"],
}
```

### Authentication Flow
1. Client sends request with `Authorization: Bearer <token>`
2. Gateway validates JWT token signature and expiration
3. Gateway extracts user claims from token
4. Gateway forwards user context to backend service
5. Backend service processes request with user context

## Rate Limiting

### Token Bucket Algorithm
- **Capacity**: Maximum number of tokens (burst size)
- **Refill Rate**: Tokens added per minute
- **Per-IP Limiting**: Separate bucket for each client IP
- **Automatic Cleanup**: Removes unused buckets after 10 minutes

### Configuration
```bash
RATE_LIMIT_RPM=100    # 100 requests per minute
RATE_LIMIT_BURST=20   # 20 request burst capacity
```

### Rate Limit Headers
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 85
X-RateLimit-Reset: 1640995200
```

## Health Monitoring

### Gateway Health Check
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "services": {
    "user-service": {
      "service": "user-service",
      "status": "healthy",
      "timestamp": "2024-01-01T12:00:00Z"
    },
    "patient-service": {
      "service": "patient-service",
      "status": "healthy",
      "timestamp": "2024-01-01T12:00:00Z"
    },
    "appointment-service": {
      "service": "appointment-service",
      "status": "healthy",
      "timestamp": "2024-01-01T12:00:00Z"
    }
  },
  "gateway": {
    "service": "api-gateway",
    "status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z",
    "version": "1.0.0"
  }
}
```

### Service Health Monitoring
- **Individual Checks**: Each service health is checked separately
- **Timeout Handling**: 5-second timeout for health checks
- **Status Aggregation**: Overall status reflects all service health
- **Failure Detection**: Unhealthy services are clearly identified

## Statistics & Monitoring

### Request Statistics
```bash
curl http://localhost:8080/stats
```

Response:
```json
{
  "total_requests": 10000,
  "success_requests": 9500,
  "error_requests": 500,
  "avg_response_time_ms": 125.5,
  "service_stats": {
    "user-service": {
      "requests": 3000,
      "successes": 2950,
      "errors": 50,
      "avg_latency_ms": 45.2,
      "last_request": "2024-01-01T12:00:00Z"
    },
    "patient-service": {
      "requests": 4000,
      "successes": 3900,
      "errors": 100,
      "avg_latency_ms": 85.1,
      "last_request": "2024-01-01T12:00:00Z"
    },
    "appointment-service": {
      "requests": 3000,
      "successes": 2650,
      "errors": 350,
      "avg_latency_ms": 195.3,
      "last_request": "2024-01-01T12:00:00Z"
    }
  }
}
```

### Metrics Collected
- **Total Requests**: All requests processed
- **Success/Error Rates**: HTTP status code based categorization
- **Response Times**: Average latency per service
- **Service Performance**: Individual service metrics
- **Last Activity**: Timestamp of most recent request

## Development

### Local Development
```bash
# Clone repository
git clone <repo-url>
cd healthcare-platform/services/api-gateway

# Install dependencies
go mod download

# Set environment variables
cp .env.example .env

# Ensure backend services are running
docker compose up -d user-service patient-service appointment-service

# Run gateway
go run .
```

### Testing
```bash
# Run unit tests
go test ./...

# Test rate limiting
for i in {1..25}; do curl http://localhost:8080/health; done

# Test authentication
curl -H "Authorization: Bearer invalid-token" http://localhost:8080/api/users/profile
```

### Code Structure

```
api-gateway/
├── main.go              # Gateway entry point and routing
├── config.go            # Configuration management
├── proxy.go             # Request proxying logic
├── auth_middleware.go   # JWT authentication middleware
├── rate_limiter.go      # Rate limiting implementation
├── stats.go             # Statistics collection
├── models.go            # Data structures
├── go.mod              # Go dependencies
├── .env.example        # Environment template
└── Dockerfile          # Container build instructions
```

## Request Proxying

### Proxy Logic
1. **Route Matching**: Find matching route configuration
2. **Service Discovery**: Get target service URL
3. **Request Preparation**: Copy headers and body
4. **Authentication**: Add user context headers
5. **Request Execution**: Forward to backend service
6. **Response Handling**: Copy response back to client
7. **Statistics**: Record request metrics

### Header Handling
```go
// Headers excluded from forwarding (hop-by-hop)
var hopByHopHeaders = []string{
    "Connection",
    "Keep-Alive", 
    "Proxy-Authenticate",
    "Proxy-Authorization",
    "Te",
    "Trailers",
    "Transfer-Encoding",
    "Upgrade",
}
```

### Retry Logic
- **Retry Count**: 3 attempts for failed requests
- **Backoff Strategy**: Exponential backoff (1s, 2s, 4s)
- **Max Wait Time**: 5 seconds maximum wait
- **Timeout**: 30 seconds per request

## Error Handling

### Standardized Error Format
```json
{
  "error": "Service unavailable",
  "code": "SERVICE_UNAVAILABLE", 
  "message": "Failed to connect to patient-service",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Error Categories
- **4xx Client Errors**: Authentication, validation, not found
- **5xx Server Errors**: Service unavailable, internal errors
- **Gateway Errors**: Rate limiting, routing failures

### Error Codes
- `ROUTE_NOT_FOUND` - No matching route
- `SERVICE_UNAVAILABLE` - Backend service down
- `AUTH_HEADER_MISSING` - No authorization header
- `TOKEN_INVALID` - JWT validation failed
- `RATE_LIMIT_EXCEEDED` - Too many requests
- `SERVICE_TIMEOUT` - Backend service timeout

## CORS Configuration

### CORS Headers
```go
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH
Access-Control-Allow-Headers: Content-Type, Authorization, X-User-ID, X-User-Email, X-User-Role
Access-Control-Expose-Headers: Content-Length
Access-Control-Allow-Credentials: true
```

### Preflight Handling
- **OPTIONS Requests**: Automatically handled
- **Status Code**: Returns 204 No Content
- **CORS Headers**: Appropriate headers sent
- **No Processing**: Preflight requests bypass business logic

## Performance Optimizations

### Connection Management
- **HTTP Client**: Reuse connections with connection pooling
- **Timeouts**: Configurable timeouts per service
- **Keep-Alive**: Persistent connections to backend services
- **Circuit Breaker**: Future enhancement for resilience

### Caching Strategy
- **Statistics**: In-memory statistics cache
- **Rate Limiting**: In-memory token buckets
- **Route Configuration**: Cached route matching
- **Health Status**: Cached health check results (future)

## Security Features

### Request Security
- **JWT Validation**: Cryptographic signature verification
- **Input Sanitization**: Request validation and sanitization
- **Header Filtering**: Hop-by-hop header removal
- **HTTPS Enforcement**: TLS termination (production)

### Rate Limiting Security
- **DDoS Protection**: Request rate limiting per IP
- **Burst Protection**: Prevents request spikes
- **IP-based Limiting**: Individual limits per client
- **Automatic Cleanup**: Memory leak prevention

## Production Considerations

### Scalability
- **Horizontal Scaling**: Multiple gateway instances
- **Load Balancing**: External load balancer (NGINX)
- **Session Management**: Stateless operation
- **Auto-scaling**: Kubernetes HPA based on CPU/memory

### High Availability
- **Health Checks**: Kubernetes liveness/readiness probes
- **Graceful Shutdown**: Proper connection draining
- **Circuit Breaker**: Service failure isolation
- **Fallback Responses**: Cached responses for critical endpoints

### Monitoring Integration
- **Metrics Export**: Prometheus metrics endpoint
- **Distributed Tracing**: Request correlation IDs
- **Log Aggregation**: Structured logging for ELK stack
- **Alerting**: Critical error rate thresholds

## Future Enhancements

### Phase 2 Features
- **API Versioning**: Support for multiple API versions
- **Request Transformation**: Request/response transformation
- **Advanced Auth**: OAuth2, SAML integration
- **Caching Layer**: Response caching with Redis

### Phase 3 Features
- **Service Mesh**: Istio integration for advanced traffic management
- **GraphQL Gateway**: GraphQL to REST translation
- **WebSocket Support**: Real-time communication proxy
- **Advanced Analytics**: Request pattern analysis

## Troubleshooting

### Common Issues

#### Service Connection Failures
```bash
# Check service health
curl http://localhost:8080/health

# Check individual service
curl http://localhost:8081/health

# View gateway logs
docker compose logs api-gateway
```

#### Authentication Issues
- Verify `JWT_SECRET` matches user-service configuration
- Check token expiration time
- Ensure Authorization header format: `Bearer <token>`

#### Rate Limiting Issues
```bash
# Check current rate limit status
curl -v http://localhost:8080/health

# Look for rate limit headers in response
X-RateLimit-Remaining: 0
```

#### Performance Issues
- Monitor service response times in `/stats` endpoint
- Check for slow backend services
- Verify connection pool settings

### Debug Mode
```bash
export LOG_LEVEL=debug
go run .
```

### Configuration Validation
```bash
# Check configuration
curl http://localhost:8080/info
```

## Contributing

Please read the [Contributing Guidelines](../../docs/CONTRIBUTING.md) before making changes to this service.

### Key Guidelines
- Test routing changes thoroughly
- Ensure authentication security
- Monitor performance impact
- Maintain backward compatibility

---

**Maintainer**: Healthcare Platform Team  
**Last Updated**: 2024-01-01