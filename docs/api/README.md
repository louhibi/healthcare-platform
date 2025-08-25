# Healthcare Platform API Documentation

Complete API reference for the Healthcare Management Platform.

## Base URL

```
Development: http://localhost:8080
Production: https://your-domain.com
```

## Authentication

All API endpoints (except auth endpoints) require JWT authentication.

### Authentication Header
```
Authorization: Bearer <jwt_token>
```

### Token Structure
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "role": "doctor",
  "exp": 1640995200,
  "iat": 1640991600,
  "type": "access"
}
```

## Rate Limiting

- **Default**: 100 requests per minute per IP
- **Burst**: 20 requests
- **Headers**: Rate limit info in response headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1640995200
```

## Error Responses

All errors follow consistent format:

```json
{
  "error": "Error description",
  "code": "ERROR_CODE",
  "message": "Detailed error message",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Common Error Codes

| Code | Status | Description |
|------|--------|-------------|
| `ROUTE_NOT_FOUND` | 404 | Endpoint not found |
| `METHOD_NOT_ALLOWED` | 405 | HTTP method not allowed |
| `AUTH_HEADER_MISSING` | 401 | Authorization header required |
| `TOKEN_INVALID` | 401 | Invalid JWT token |
| `ROLE_FORBIDDEN` | 403 | Insufficient permissions |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `VALIDATION_ERROR` | 400 | Input validation failed |
| `SERVICE_UNAVAILABLE` | 502 | Backend service error |

## Gateway Endpoints

### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "services": {
    "user-service": {
      "service": "user-service",
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

### Statistics
```http
GET /stats
```

**Response:**
```json
{
  "total_requests": 1000,
  "success_requests": 950,
  "error_requests": 50,
  "avg_response_time_ms": 120.5,
  "service_stats": {
    "user-service": {
      "requests": 300,
      "successes": 295,
      "errors": 5,
      "avg_latency_ms": 45.2,
      "last_request": "2024-01-01T12:00:00Z"
    }
  }
}
```

## Authentication API

### Register User
```http
POST /api/auth/register
```

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "first_name": "John",
  "last_name": "Doe",
  "role": "doctor"
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "doctor",
    "is_active": true,
    "created_at": "2024-01-01T12:00:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

### Login
```http
POST /api/auth/login
```

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:** Same as register

### Refresh Token
```http
POST /api/auth/refresh
```

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:** Same as login

## User Management API

### Get Current User Profile
```http
GET /api/users/profile
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "doctor",
  "is_active": true,
  "created_at": "2024-01-01T12:00:00Z"
}
```

### Update User Profile
```http
PUT /api/users/profile
```

**Request:**
```json
{
  "first_name": "John",
  "last_name": "Smith"
}
```

### Get All Users (Admin Only)
```http
GET /api/users/?limit=10&offset=0
```

**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "doctor",
      "is_active": true,
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

## Patient Management API

### Get Patients
```http
GET /api/patients/?q=search&gender=male&insurance=blue&limit=10&offset=0
```

**Query Parameters:**
- `q` - Search query (name, email, phone)
- `gender` - Filter by gender
- `insurance` - Filter by insurance
- `limit` - Results per page (max 100)
- `offset` - Pagination offset

**Response:**
```json
{
  "patients": [
    {
      "id": 1,
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
      "medications": "None",
      "age": 39,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total_count": 1,
  "limit": 10,
  "offset": 0,
  "has_more": false
}
```

### Create Patient
```http
POST /api/patients/
```

**Request:**
```json
{
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
}
```

### Get Patient by ID
```http
GET /api/patients/:id
```

### Update Patient
```http
PUT /api/patients/:id
```

### Delete Patient
```http
DELETE /api/patients/:id
```

### Patient Statistics
```http
GET /api/patients/stats
```

**Response:**
```json
{
  "total_patients": 150
}
```

## Appointment Management API

### Get Appointments
```http
GET /api/appointments/?patient_id=1&doctor_id=2&status=scheduled&type=consultation&date_from=2024-01-01&date_to=2024-01-31&limit=20&offset=0
```

**Query Parameters:**
- `patient_id` - Filter by patient
- `doctor_id` - Filter by doctor
- `status` - Filter by status
- `type` - Filter by appointment type
- `date_from` - Start date filter
- `date_to` - End date filter
- `limit` - Results per page
- `offset` - Pagination offset

**Response:**
```json
{
  "appointments": [
    {
      "id": 1,
      "patient_id": 1,
      "doctor_id": 2,
      "date_time": "2024-01-15T10:00:00Z",
      "duration": 30,
      "type": "consultation",
      "status": "scheduled",
      "reason": "Annual checkup",
      "notes": "Patient requested morning appointment",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "end_time": "2024-01-15T10:30:00Z"
    }
  ],
  "total_count": 1,
  "limit": 20,
  "offset": 0,
  "has_more": false
}
```

### Create Appointment
```http
POST /api/appointments/
```

**Request:**
```json
{
  "patient_id": 1,
  "doctor_id": 2,
  "date_time": "2024-01-15T10:00:00Z",
  "duration": 30,
  "type": "consultation",
  "reason": "Annual checkup",
  "notes": "Patient requested morning appointment"
}
```

### Get Appointment by ID
```http
GET /api/appointments/:id
```

### Update Appointment
```http
PUT /api/appointments/:id
```

### Update Appointment Status
```http
PATCH /api/appointments/:id/status
```

**Request:**
```json
{
  "status": "completed",
  "notes": "Patient arrived on time, examination completed"
}
```

### Delete Appointment
```http
DELETE /api/appointments/:id
```

## Validation Rules

### User Registration
- `email`: Valid email format, unique
- `password`: Minimum 8 characters
- `first_name`: Required, max 100 chars
- `last_name`: Required, max 100 chars
- `role`: Must be one of: admin, doctor, nurse, staff

### Patient
- `first_name`: Required, max 100 chars
- `last_name`: Required, max 100 chars
- `date_of_birth`: Required, valid date
- `gender`: Must be one of: male, female, other
- `phone`: Required, valid phone format
- `email`: Required, valid email format

### Appointment
- `patient_id`: Required, valid patient ID
- `doctor_id`: Required, valid doctor ID
- `date_time`: Required, future date/time
- `duration`: Between 15 and 480 minutes
- `type`: Must be one of: consultation, follow-up, procedure, emergency
- `reason`: Required

## Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 429 | Too Many Requests |
| 500 | Internal Server Error |
| 502 | Bad Gateway |
| 503 | Service Unavailable |