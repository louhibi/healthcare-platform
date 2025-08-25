# Appointment Service

Appointment scheduling and management microservice for the Healthcare Management Platform.

## Overview

The Appointment Service handles all appointment-related operations including scheduling, conflict detection, availability management, and appointment status tracking. It provides intelligent scheduling capabilities with automated conflict prevention.

## Features

- **Smart Scheduling**: Appointment creation with conflict detection
- **Availability Management**: Doctor schedule and availability tracking
- **Status Tracking**: Complete appointment lifecycle management
- **Conflict Prevention**: Automatic detection of scheduling conflicts
- **Calendar Integration**: Weekly/monthly calendar views
- **Advanced Filtering**: Multi-criteria appointment filtering
- **Statistics**: Appointment analytics and reporting

## API Endpoints

### Appointment Operations
```http
GET    /api/appointments/           # Get appointments with filters
POST   /api/appointments/           # Create new appointment
GET    /api/appointments/:id        # Get appointment by ID
PUT    /api/appointments/:id        # Update appointment
DELETE /api/appointments/:id        # Delete appointment (soft delete)
PATCH  /api/appointments/:id/status # Update appointment status
```

### Health Check
```http
GET    /health                      # Service health status
```

## Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=appointment_service_db

# Server Configuration
PORT=8083
ENV=development
```

## Database Schema

### Appointments Table
```sql
CREATE TABLE appointments (
    id SERIAL PRIMARY KEY,
    patient_id INTEGER NOT NULL,
    doctor_id INTEGER NOT NULL,
    date_time TIMESTAMP NOT NULL,
    duration INTEGER NOT NULL CHECK (duration >= 15 AND duration <= 480),
    type VARCHAR(50) NOT NULL CHECK (type IN ('consultation', 'follow-up', 'procedure', 'emergency')),
    status VARCHAR(50) NOT NULL CHECK (status IN ('scheduled', 'confirmed', 'in-progress', 'completed', 'cancelled', 'no-show')),
    reason TEXT NOT NULL,
    notes TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL
);
```

### Doctor Schedules Table
```sql
CREATE TABLE doctor_schedules (
    id SERIAL PRIMARY KEY,
    doctor_id INTEGER NOT NULL,
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    break_start TIME,
    break_end TIME,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(doctor_id, day_of_week)
);
```

### Indexes
```sql
CREATE INDEX idx_appointments_patient_id ON appointments(patient_id);
CREATE INDEX idx_appointments_doctor_id ON appointments(doctor_id);
CREATE INDEX idx_appointments_date_time ON appointments(date_time);
CREATE INDEX idx_appointments_status ON appointments(status);
CREATE INDEX idx_appointments_conflict_check 
    ON appointments(doctor_id, date_time, duration, status) 
    WHERE is_active = true AND status IN ('scheduled', 'confirmed', 'in-progress');
```

## Appointment Types

| Type | Description | Typical Duration |
|------|-------------|------------------|
| `consultation` | Initial patient consultation | 30-60 minutes |
| `follow-up` | Follow-up appointment | 15-30 minutes |
| `procedure` | Medical procedure | 60-240 minutes |
| `emergency` | Emergency appointment | 15-30 minutes |

## Appointment Status Flow

```
scheduled → confirmed → in-progress → completed
    ↓           ↓            ↓
cancelled   cancelled   cancelled
    ↓           ↓            ↓
no-show     no-show     no-show
```

### Status Descriptions
- **scheduled**: Appointment created and pending confirmation
- **confirmed**: Patient confirmed attendance
- **in-progress**: Appointment currently happening
- **completed**: Appointment finished successfully
- **cancelled**: Appointment cancelled by patient or staff
- **no-show**: Patient did not show up for appointment

## API Examples

### Create Appointment
```bash
curl -X POST http://localhost:8083/api/appointments/ \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "patient_id": 1,
    "doctor_id": 2,
    "date_time": "2024-01-15T10:00:00Z",
    "duration": 30,
    "type": "consultation",
    "reason": "Annual checkup",
    "notes": "Patient requested morning appointment"
  }'
```

### Get Appointments with Filters
```bash
# Get doctor's appointments for specific date
curl "http://localhost:8083/api/appointments/?doctor_id=2&date_from=2024-01-15&date_to=2024-01-15" \
  -H "X-User-ID: 1"

# Get patient's appointments
curl "http://localhost:8083/api/appointments/?patient_id=1&status=scheduled" \
  -H "X-User-ID: 1"

# Get appointments by type
curl "http://localhost:8083/api/appointments/?type=consultation&limit=20" \
  -H "X-User-ID: 1"
```

### Update Appointment Status
```bash
curl -X PATCH http://localhost:8083/api/appointments/1/status \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "status": "confirmed",
    "notes": "Patient confirmed via phone"
  }'
```

### Get Doctor Schedule
```bash
curl "http://localhost:8083/api/doctors/2/schedule?date=2024-01-15" \
  -H "X-User-ID: 1"
```

## Development

### Local Development
```bash
# Clone repository
git clone <repo-url>
cd healthcare-platform/services/appointment-service

# Install dependencies
go mod download

# Set environment variables
cp .env.example .env

# Start database (from project root)
docker compose up -d appointment-db

# Run service
go run .
```

### Testing
```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Test conflict detection
go test -run TestConflictDetection
```

### Code Structure

```
appointment-service/
├── main.go                  # Service entry point and routing
├── models.go               # Appointment data structures
├── database.go             # Database connection and migrations
├── appointment_service.go  # Appointment business logic
├── handlers.go             # HTTP request handlers
├── go.mod                  # Go dependencies
├── .env.example           # Environment template
└── Dockerfile             # Container build instructions
```

## Conflict Detection

### Algorithm
The service implements intelligent conflict detection:

1. **Time Overlap Check**: Detects overlapping appointment times
2. **Doctor Availability**: Checks doctor's working hours
3. **Break Time Avoidance**: Respects doctor's break periods
4. **Status Consideration**: Only active appointments cause conflicts

### Conflict Detection Logic
```go
// Two appointments conflict if:
// 1. Same doctor
// 2. Both are active (scheduled, confirmed, in-progress)
// 3. Time periods overlap

func (a *Appointment) IsConflicting(other *Appointment) bool {
    if a.DoctorID != other.DoctorID {
        return false
    }
    
    aStart := a.DateTime
    aEnd := a.DateTime.Add(time.Duration(a.Duration) * time.Minute)
    bStart := other.DateTime
    bEnd := other.DateTime.Add(time.Duration(other.Duration) * time.Minute)
    
    return aStart.Before(bEnd) && bStart.Before(aEnd)
}
```

## Doctor Availability

### Working Hours
Default working schedule:
- **Monday-Friday**: 9:00 AM - 5:00 PM
- **Lunch Break**: 12:00 PM - 1:00 PM
- **Weekend**: Not available (configurable)

### Availability Slots
The service generates available time slots:
```json
{
  "doctor_id": 2,
  "date": "2024-01-15",
  "working_hours": {
    "start_time": "09:00:00",
    "end_time": "17:00:00",
    "break_start": "12:00:00",
    "break_end": "13:00:00"
  },
  "available_slots": [
    {
      "date_time": "2024-01-15T09:00:00Z",
      "duration": 30
    },
    {
      "date_time": "2024-01-15T09:30:00Z",
      "duration": 30
    }
  ]
}
```

## Search and Filtering

### Filter Options
- **Patient ID**: Filter by specific patient
- **Doctor ID**: Filter by specific doctor
- **Status**: Filter by appointment status
- **Type**: Filter by appointment type
- **Date Range**: Filter by date range
- **Duration**: Filter by appointment duration

### Advanced Queries
```bash
# Complex filter example
curl "http://localhost:8083/api/appointments/?doctor_id=2&status=scheduled&type=consultation&date_from=2024-01-01&date_to=2024-01-31&limit=50" \
  -H "X-User-ID: 1"
```

## Authentication & Authorization

### User Context
Receives user information via HTTP headers:
```
X-User-ID: 123
X-User-Email: doctor@hospital.com
X-User-Role: doctor
```

### Access Control
- **All Operations**: Require authentication
- **Data Tracking**: Appointments tracked by creator (`created_by`)
- **Role-Based Access**: All authenticated roles can manage appointments
- **Future Enhancement**: Doctor-specific appointment access

## Validation

### Appointment Validation
```go
type AppointmentRequest struct {
    PatientID int       `validate:"required"`
    DoctorID  int       `validate:"required"`
    DateTime  time.Time `validate:"required"`
    Duration  int       `validate:"required,min=15,max=480"`
    Type      string    `validate:"required,oneof=consultation follow-up procedure emergency"`
    Reason    string    `validate:"required"`
}
```

### Business Rules
- **Future Appointments**: Cannot schedule appointments in the past
- **Duration Limits**: 15 minutes minimum, 8 hours maximum
- **Working Hours**: Must be within doctor's working hours
- **Conflict Prevention**: No overlapping appointments for same doctor

## Error Handling

### Common Error Responses
```json
{
  "error": "Appointment conflicts with existing appointment",
  "code": "APPOINTMENT_CONFLICT",
  "message": "The selected time slot conflicts with an existing appointment",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Error Codes
- `APPOINTMENT_NOT_FOUND` - Appointment does not exist
- `APPOINTMENT_CONFLICT` - Scheduling conflict detected
- `INVALID_TIME_SLOT` - Outside working hours
- `VALIDATION_ERROR` - Input validation failed
- `STATUS_TRANSITION_INVALID` - Invalid status change

## Performance Optimizations

### Database Performance
- **Conflict Index**: Optimized index for conflict detection queries
- **Date Range Queries**: Efficient date-based filtering
- **Status Filtering**: Index on status for quick filtering
- **Pagination**: Efficient offset-based pagination

### Caching Strategy
- **Doctor Schedules**: Cache working hours per doctor
- **Availability Slots**: Cache generated time slots
- **Frequent Queries**: Cache common filter combinations

## Monitoring

### Health Check
```bash
curl http://localhost:8083/health
```

### Statistics
```bash
curl http://localhost:8083/api/appointments/stats \
  -H "X-User-ID: 1"
```

Response:
```json
{
  "total_appointments": 1250,
  "by_status": {
    "scheduled": 300,
    "confirmed": 150,
    "completed": 700,
    "cancelled": 80,
    "no-show": 20
  },
  "today_appointments": 25
}
```

### Logging
- **Appointment Operations**: Create, update, delete operations
- **Conflict Detection**: Conflicts found and resolved
- **Status Changes**: Appointment status transitions
- **Performance Metrics**: Query execution times

## Security Features

### Data Protection
- **Input Validation**: Comprehensive validation of all inputs
- **SQL Injection Prevention**: Parameterized queries only
- **Access Control**: Authentication required for all operations
- **Audit Trail**: All operations logged with user context

### HIPAA Compliance
- **Patient Privacy**: Patient information properly protected
- **Access Logging**: All data access logged
- **Data Retention**: Configurable retention policies
- **Secure Communication**: HTTPS in production

## Integration

### Service Dependencies
- **Patient Service**: Validates patient IDs exist
- **User Service**: Validates doctor IDs exist
- **API Gateway**: Handles authentication and routing

### Future Integrations
- **Notification Service**: Appointment reminders and confirmations
- **Calendar Service**: External calendar synchronization
- **Billing Service**: Appointment-based billing
- **Analytics Service**: Scheduling analytics and optimization

## Production Considerations

### Scalability
- **Read Replicas**: Separate read/write database operations
- **Horizontal Scaling**: Multiple service instances
- **Caching Layer**: Redis for performance optimization
- **Connection Pooling**: Efficient database connections

### High Availability
- **Health Checks**: Continuous health monitoring
- **Graceful Shutdown**: Proper service shutdown handling
- **Circuit Breaker**: Resilience patterns (future)
- **Backup Services**: Automated backup procedures

## Future Enhancements

### Phase 2 Features
- **Recurring Appointments**: Weekly/monthly recurring schedules
- **Waitlist Management**: Automatic notification for cancellations
- **Online Booking**: Patient self-scheduling portal
- **Calendar Sync**: Google Calendar, Outlook integration

### Phase 3 Features
- **AI Scheduling**: Machine learning for optimal scheduling
- **Resource Management**: Room and equipment scheduling
- **Telemedicine**: Virtual appointment support
- **Advanced Analytics**: Scheduling optimization insights

## Troubleshooting

### Common Issues

#### Conflict Detection Not Working
- Check database indexes are properly created
- Verify appointment status values
- Ensure timezone handling is correct

#### Performance Issues
```bash
# Check slow queries
EXPLAIN ANALYZE SELECT * FROM appointments WHERE doctor_id = 2 AND date_time >= '2024-01-01';

# Monitor database connections
SELECT * FROM pg_stat_activity WHERE datname = 'appointment_service_db';
```

#### Schedule Generation Issues
- Verify doctor_schedules table has data
- Check working hours configuration
- Ensure timezone consistency

### Debug Mode
```bash
export LOG_LEVEL=debug
go run .
```

## Contributing

Please read the [Contributing Guidelines](../../docs/CONTRIBUTING.md) before making changes to this service.

### Key Guidelines
- Test conflict detection thoroughly
- Ensure timezone handling is correct
- Add tests for new scheduling features
- Maintain performance for large datasets

---

**Maintainer**: Healthcare Platform Team  
**Last Updated**: 2024-01-01