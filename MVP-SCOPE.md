# Healthcare Management Platform - MVP Scope

## MVP Core Features

### 1. User Management (Basic)
- User registration and login
- Basic role-based access (Doctor, Staff, Admin)
- JWT authentication
- Password reset functionality

### 2. Patient Management
- Create, read, update, delete patient profiles
- Basic patient information (name, DOB, contact, insurance)
- Patient search functionality
- Medical history notes

### 3. Appointment Scheduling
- Create appointments with date/time
- View appointments by day/week
- Basic appointment status (scheduled, completed, cancelled)
- Simple conflict detection

### 4. Basic Dashboard
- Today's appointments overview
- Patient count statistics
- Recent activities feed

## Technical MVP Stack

### Backend
- Go with Gin framework
- PostgreSQL database
- JWT authentication
- RESTful APIs
- Basic logging

### Frontend
- Vue.js 3 with Composition API
- Vue Router for navigation
- Axios for API calls
- Tailwind CSS for styling
- Basic responsive design

### Infrastructure
- Docker containers
- Docker Compose for local development
- Environment-based configuration

## MVP Limitations (Future Features)
- No advanced scheduling features
- No inventory management
- No complex role permissions
- No file uploads
- No notifications/messaging
- No reporting/analytics
- No third-party integrations
- No advanced security features

## Success Criteria
- Users can register and log in
- Staff can manage patient profiles
- Staff can schedule and view appointments
- Basic conflict prevention works
- System runs reliably in Docker environment