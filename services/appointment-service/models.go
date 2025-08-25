package main

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

// Appointment represents an appointment in the system
type Appointment struct {
	ID                int       `json:"id" db:"id"`
	HealthcareEntityID int      `json:"healthcare_entity_id" db:"healthcare_entity_id" validate:"required"`
	PatientID         int       `json:"patient_id" db:"patient_id" validate:"required"`
	DoctorID          int       `json:"doctor_id" db:"doctor_id" validate:"required"`
	DateTime          time.Time `json:"date_time" db:"date_time" validate:"required"`
	Duration          int       `json:"duration" db:"duration" validate:"required,min=15,max=480"` // in minutes
	Type              string    `json:"type" db:"type" validate:"required,oneof=consultation follow-up procedure emergency"`
	Status            string    `json:"status" db:"status" validate:"required,oneof=scheduled confirmed in-progress completed cancelled no-show"`
	Reason            string    `json:"reason" db:"reason" validate:"required"`
	Notes             string    `json:"notes" db:"notes"`
	Priority          string         `json:"priority" db:"priority" validate:"oneof=low normal high urgent"` // appointment priority
	RoomID            sql.NullInt32  `json:"room_id" db:"room_id"` // Use sql.NullString to handle NULL
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy         int       `json:"created_by" db:"created_by"`
	IsActive          bool      `json:"is_active" db:"is_active"`
}

// AppointmentRequest represents appointment creation/update request
type AppointmentRequest struct {
	PatientID    int    `json:"patient_id" validate:"required"`
	DoctorID     int    `json:"doctor_id" validate:"required"`
	DateTime     string `json:"date_time" validate:"required"` // ISO format: 2024-01-15T10:00:00Z
	Duration     int    `json:"duration" validate:"required,min=15,max=480"`
	Type         string `json:"type" validate:"required,oneof=consultation follow-up procedure emergency"`
	Reason       string `json:"reason" validate:"required"`
	Notes        string `json:"notes"`
	Priority     string `json:"priority" validate:"oneof=low normal high urgent"`
	RoomID       string `json:"room_id"` // Accept room_id from frontend
}

// AppointmentUpdate represents appointment status update
type AppointmentUpdate struct {
	Status string `json:"status" validate:"required,oneof=scheduled confirmed in-progress completed cancelled no-show"`
	Notes  string `json:"notes"`
}

// AppointmentSearch represents search parameters
type AppointmentSearch struct {
	HealthcareEntityID int       `form:"healthcare_entity_id"`
	PatientID          int       `form:"patient_id"`
	DoctorID           int       `form:"doctor_id"`
	Status             string    `form:"status"`
	Type               string    `form:"type"`
	DateFrom           time.Time `form:"date_from"`
	DateTo             time.Time `form:"date_to"`
	IncludePast        bool      `form:"include_past"` // If false (default), only show future/current appointments
	Limit              int       `form:"limit"`
	Offset             int       `form:"offset"`
}

// AppointmentResponse represents appointment data returned to client
type AppointmentResponse struct {
	ID                 int       `json:"id"`
	HealthcareEntityID int       `json:"healthcare_entity_id"`
	PatientID          int       `json:"patient_id"`
	PatientName        string    `json:"patient_name"`
	DoctorID           int       `json:"doctor_id"`
	DoctorName         string    `json:"doctor_name"`
	DateTime           time.Time `json:"date_time"`
	Duration           int       `json:"duration"`
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	Reason             string    `json:"reason"`
	Notes              string    `json:"notes"`
	Priority           string    `json:"priority"`
	RoomID             *int      `json:"room_id"` // Will be populated from sql.NullInt32
	RoomNumber         string    `json:"room_number"` // Will be populated by service layer
	RoomName           string    `json:"room_name"` // Will be populated by service layer
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	EndTime            time.Time `json:"end_time"`
}

// AvailabilitySlot represents an available time slot
type AvailabilitySlot struct {
	DateTime    time.Time `json:"date_time"`
	Duration    int       `json:"duration"`
	IsAvailable bool      `json:"is_available"`
	SlotType    string    `json:"slot_type"` // morning, afternoon, evening
}

// DoctorSchedule represents doctor's schedule for a day
type DoctorSchedule struct {
	DoctorID      int                `json:"doctor_id"`
	Date          time.Time          `json:"date"`
	WorkingHours  WorkingHours       `json:"working_hours"`
	Appointments  []Appointment      `json:"appointments"`
	AvailableSlots []AvailabilitySlot `json:"available_slots"`
}

// WorkingHours represents working hours for a day
type WorkingHours struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	BreakStart time.Time `json:"break_start"`
	BreakEnd   time.Time `json:"break_end"`
}

// ConflictCheck represents appointment conflict checking
type ConflictCheck struct {
	DoctorID           int       `json:"doctor_id"`
	DateTime           time.Time `json:"date_time"`
	Duration           int       `json:"duration"`
	ExcludeID          int       `json:"exclude_id"` // For updates
	RoomID             int       `json:"room_id"`
	HealthcareEntityID int       `json:"healthcare_entity_id"`
}

// ToAppointmentResponse converts Appointment to AppointmentResponse
func (a *Appointment) ToAppointmentResponse() AppointmentResponse {
	var roomID *int
	if a.RoomID.Valid {
		roomID = new(int)
		*roomID = int(a.RoomID.Int32)
	}

	return AppointmentResponse{
		ID:                 a.ID,
		HealthcareEntityID: a.HealthcareEntityID,
		PatientID:          a.PatientID,
		PatientName:        "Patient", // Will be filled by service layer
		DoctorID:           a.DoctorID,
		DoctorName:         "",  // Doctor info comes from user service
		DateTime:           a.DateTime,
		Duration:           a.Duration,
		Type:               a.Type,
		Status:             a.Status,
		Reason:             a.Reason,
		Notes:              a.Notes,
		Priority:           a.Priority,
		RoomID:             roomID,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
		EndTime:            a.DateTime.Add(time.Duration(a.Duration) * time.Minute),
	}
}

// BookingRequest represents appointment booking request with conflict checking
type BookingRequest struct {
	PatientID       int    `json:"patient_id" validate:"required"`
	DoctorID        int    `json:"doctor_id" validate:"required"`
	DateTime        string `json:"date_time" validate:"required"`      // ISO 8601 UTC format: 2006-01-02T15:04:05Z
	Duration        int    `json:"duration" validate:"required,min=15,max=480"`
	Type            string `json:"type" validate:"required,oneof=consultation follow-up procedure emergency"`
	Reason          string `json:"reason" validate:"required"`
	Notes           string `json:"notes"`
	Priority        string `json:"priority" validate:"oneof=low normal high urgent"`
	RoomID          int    `json:"room_id"`           // Room ID for appointment
	CheckConflicts  bool   `json:"check_conflicts"`  // Whether to check for conflicts
}

// BookingResponse represents appointment booking response with conflict information
type BookingResponse struct {
	Success           bool                    `json:"success"`
	AppointmentID     int                     `json:"appointment_id,omitempty"`
	Appointment       *AppointmentResponse    `json:"appointment,omitempty"`
	Conflicts         []ConflictInfo          `json:"conflicts,omitempty"`
	AlternativeSlots  []AvailabilitySlot      `json:"alternative_slots,omitempty"`
	Message           string                  `json:"message"`
}

// ConflictInfo represents information about appointment conflicts
type ConflictInfo struct {
	ConflictType    string    `json:"conflict_type"`    // doctor_busy, room_occupied, outside_hours
	ExistingAppointment *AppointmentResponse `json:"existing_appointment,omitempty"`
	ConflictTime    time.Time `json:"conflict_time"`
	ConflictEnd     time.Time `json:"conflict_end"`
	Description     string    `json:"description"`
}

// PatientInfo represents basic patient information for appointments
type PatientInfo struct {
	ID        int    `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phone" db:"phone"`
}

// AppointmentCalendarDay represents a day in the appointment calendar
type AppointmentCalendarDay struct {
	Date         string                `json:"date"`         // YYYY-MM-DD
	DayOfWeek    string                `json:"day_of_week"`
	IsToday      bool                  `json:"is_today"`
	IsAvailable  bool                  `json:"is_available"`
	Appointments []AppointmentResponse `json:"appointments"`
	AvailableSlots []AvailabilitySlot  `json:"available_slots"`
	Status       string                `json:"status"`       // available, busy, unavailable
}

// AppointmentStats represents appointment statistics
type AppointmentStats struct {
	TotalAppointments    int            `json:"total_appointments"`
	TodayAppointments    int            `json:"today_appointments"`
	UpcomingAppointments int            `json:"upcoming_appointments"`
	StatusCounts         map[string]int `json:"status_counts"`
	TypeCounts           map[string]int `json:"type_counts"`
	MonthlyStats         []MonthlyStatistic `json:"monthly_stats"`
}

// MonthlyStatistic represents monthly appointment statistics
type MonthlyStatistic struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

// IsConflicting checks if two appointments conflict
func (a *Appointment) IsConflicting(other *Appointment) bool {
	if a.DoctorID != other.DoctorID {
		return false
	}
	
	// Both appointments must be in active status to conflict
	activeStatuses := map[string]bool{
		"scheduled":   true,
		"confirmed":   true,
		"in-progress": true,
	}
	
	if !activeStatuses[a.Status] || !activeStatuses[other.Status] {
		return false
	}
	
	aStart := a.DateTime
	aEnd := a.DateTime.Add(time.Duration(a.Duration) * time.Minute)
	bStart := other.DateTime
	bEnd := other.DateTime.Add(time.Duration(other.Duration) * time.Minute)
	
	// Check for overlap
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}

// DoctorAvailability represents doctor availability and status
type DoctorAvailability struct {
	ID                    int       `json:"id" db:"id"`
	HealthcareEntityID    int       `json:"healthcare_entity_id" db:"healthcare_entity_id" validate:"required"`
	DoctorID              int       `json:"doctor_id" db:"doctor_id" validate:"required"`
	Status                string    `json:"status" db:"status" validate:"required,oneof=available unavailable vacation training sick_leave meeting"`
	// UTC datetime fields - these are now the primary fields
	StartDateTime         *time.Time `json:"start_datetime" db:"start_datetime" validate:"required"` // UTC timestamps
	EndDateTime           *time.Time `json:"end_datetime" db:"end_datetime" validate:"required"`
	BreakStartDateTime    *time.Time `json:"break_start_datetime" db:"break_start_datetime"`
	BreakEndDateTime      *time.Time `json:"break_end_datetime" db:"break_end_datetime"`
	Notes                 string    `json:"notes" db:"notes"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy             int       `json:"created_by" db:"created_by"`
}

// DoctorAvailabilityRequest represents availability creation/update request
type DoctorAvailabilityRequest struct {
	DoctorID              int    `json:"doctor_id" validate:"required"`
	StartDateTime         string `json:"start_datetime" validate:"required"` // ISO 8601 UTC format: 2006-01-02T15:04:05Z
	EndDateTime           string `json:"end_datetime" validate:"required"`   // ISO 8601 UTC format: 2006-01-02T15:04:05Z
	Status                string `json:"status" validate:"required,oneof=available unavailable vacation training sick_leave meeting"`
	BreakStartDateTime    string `json:"break_start_datetime,omitempty"`     // ISO 8601 UTC format: 2006-01-02T15:04:05Z  
	BreakEndDateTime      string `json:"break_end_datetime,omitempty"`       // ISO 8601 UTC format: 2006-01-02T15:04:05Z
	Notes                 string `json:"notes"`
}

// DoctorAvailabilityResponse represents availability data returned to client
type DoctorAvailabilityResponse struct {
	ID                    int       `json:"id"`
	DoctorID              int       `json:"doctor_id"`
	DoctorName            string    `json:"doctor_name"`
	Status                string    `json:"status"`
	StartDateTime         *time.Time `json:"start_datetime,omitempty"` // UTC timestamp (nullable)
	EndDateTime           *time.Time `json:"end_datetime,omitempty"`   // UTC timestamp (nullable)
	BreakStartDateTime    *time.Time `json:"break_start_datetime,omitempty"` // UTC timestamp (nullable)
	BreakEndDateTime      *time.Time `json:"break_end_datetime,omitempty"`   // UTC timestamp (nullable)
	EntityTimezone        string    `json:"entity_timezone"`           // Healthcare entity's IANA timezone
	Notes                 string    `json:"notes"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// DoctorInfo represents basic doctor information for schedules
type DoctorInfo struct {
	ID             int    `json:"id" db:"id"`
	FirstName      string `json:"first_name" db:"first_name"`
	LastName       string `json:"last_name" db:"last_name"`
	Email          string `json:"email" db:"email"`
	Specialization string `json:"specialization" db:"specialization"`
	Role           string `json:"role" db:"role"`
}

// AvailabilitySearch represents search parameters for availability
type AvailabilitySearch struct {
	HealthcareEntityID int    `form:"healthcare_entity_id"`
	DoctorID           int    `form:"doctor_id"`
	DateFrom           string `form:"date_from"`
	DateTo             string `form:"date_to"`
	Status             string `form:"status"`
	Limit              int    `form:"limit"`
	Offset             int    `form:"offset"`
}

// ScheduleTemplate represents recurring schedule patterns
type ScheduleTemplate struct {
	ID                int       `json:"id" db:"id"`
	HealthcareEntityID int      `json:"healthcare_entity_id" db:"healthcare_entity_id" validate:"required"`
	DoctorID          int       `json:"doctor_id" db:"doctor_id" validate:"required"`
	DayOfWeek         int       `json:"day_of_week" db:"day_of_week" validate:"min=0,max=6"` // 0=Sunday, 6=Saturday
	StartTime         string    `json:"start_time" db:"start_time" validate:"required"`
	EndTime           string    `json:"end_time" db:"end_time" validate:"required"`
	BreakStart        string    `json:"break_start" db:"break_start"`
	BreakEnd          string    `json:"break_end" db:"break_end"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// ToAvailabilityResponse converts DoctorAvailability to response
func (da *DoctorAvailability) ToAvailabilityResponse(doctorName string) DoctorAvailabilityResponse {
	return DoctorAvailabilityResponse{
		ID:                 da.ID,
		DoctorID:           da.DoctorID,
		DoctorName:         doctorName,
		Status:             da.Status,
		StartDateTime:      da.StartDateTime,
		EndDateTime:        da.EndDateTime,
		BreakStartDateTime: da.BreakStartDateTime,
		BreakEndDateTime:   da.BreakEndDateTime,
		Notes:              da.Notes,
		CreatedAt:          da.CreatedAt,
		UpdatedAt:          da.UpdatedAt,
	}
}

// GetStatusColor returns color class for UI display
func (da *DoctorAvailability) GetStatusColor() string {
	statusColors := map[string]string{
		"available":  "bg-green-100 text-green-800",
		"unavailable": "bg-gray-100 text-gray-800",
		"vacation":   "bg-blue-100 text-blue-800",
		"training":   "bg-yellow-100 text-yellow-800",
		"sick_leave": "bg-red-100 text-red-800",
		"meeting":    "bg-purple-100 text-purple-800",
	}
	if color, exists := statusColors[da.Status]; exists {
		return color
	}
	return "bg-gray-100 text-gray-800"
}

// GetStatusIcon returns icon for UI display
func (da *DoctorAvailability) GetStatusIcon() string {
	statusIcons := map[string]string{
		"available":  "✓",
		"unavailable": "✗",
		"vacation":   "🏖️",
		"training":   "📚",
		"sick_leave": "🤒",
		"meeting":    "👥",
	}
	if icon, exists := statusIcons[da.Status]; exists {
		return icon
	}
	return "?"
}

// IsAvailable checks if doctor is available for appointments
func (da *DoctorAvailability) IsAvailable() bool {
	return da.Status == "available"
}

// GetDefaultWorkingHours returns default working hours (9 AM to 5 PM with 1-hour lunch)
func GetDefaultWorkingHours(date time.Time) WorkingHours {
	year, month, day := date.Date()
	location := date.Location()
	
	return WorkingHours{
		StartTime:  time.Date(year, month, day, 9, 0, 0, 0, location),
		EndTime:    time.Date(year, month, day, 17, 0, 0, 0, location),
		BreakStart: time.Date(year, month, day, 12, 0, 0, 0, location),
		BreakEnd:   time.Date(year, month, day, 13, 0, 0, 0, location),
	}
}

// AppointmentDurationSetting represents configurable appointment durations per healthcare entity
type AppointmentDurationSetting struct {
	ID                 int       `json:"id" db:"id"`
	HealthcareEntityID int       `json:"healthcare_entity_id" db:"healthcare_entity_id" validate:"required"`
	AppointmentType    string    `json:"appointment_type" db:"appointment_type" validate:"required,oneof=consultation follow-up procedure emergency"`
	DurationMinutes    int       `json:"duration_minutes" db:"duration_minutes" validate:"required,min=15,max=480"`
	IsDefault          bool      `json:"is_default" db:"is_default"`
	IsActive           bool      `json:"is_active" db:"is_active"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy          int       `json:"created_by" db:"created_by"`
}

// AppointmentDurationRequest represents duration setting creation/update request
type AppointmentDurationRequest struct {
	AppointmentType string `json:"appointment_type" validate:"required,oneof=consultation follow-up procedure emergency"`
	DurationMinutes int    `json:"duration_minutes" validate:"required,min=15,max=480"`
	IsDefault       bool   `json:"is_default"`
}

// AppointmentDurationResponse represents duration setting data returned to client
type AppointmentDurationResponse struct {
	ID                 int    `json:"id"`
	HealthcareEntityID int    `json:"healthcare_entity_id"`
	AppointmentType    string `json:"appointment_type"`
	DurationMinutes    int    `json:"duration_minutes"`
	IsDefault          bool   `json:"is_default"`
	IsActive           bool   `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// AppointmentDurationOption represents available duration choices for appointment types
type AppointmentDurationOption struct {
	ID                 int       `json:"id" db:"id"`
	HealthcareEntityID int       `json:"healthcare_entity_id" db:"healthcare_entity_id"`
	AppointmentType    string    `json:"appointment_type" db:"appointment_type"`
	DurationMinutes    int       `json:"duration_minutes" db:"duration_minutes"`
	IsDefault          bool      `json:"is_default" db:"is_default"`
	IsActive           bool      `json:"is_active" db:"is_active"`
	DisplayOrder       int       `json:"display_order" db:"display_order"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy          int       `json:"created_by" db:"created_by"`
}

// DurationOptionRequest represents request for duration option creation/update
type DurationOptionRequest struct {
	AppointmentType string `json:"appointment_type" validate:"required,oneof=consultation follow-up procedure emergency"`
	DurationMinutes int    `json:"duration_minutes" validate:"required,min=15,max=480"`
	IsDefault       bool   `json:"is_default"`
	IsActive        bool   `json:"is_active"`
	DisplayOrder    int    `json:"display_order"`
}

// Room represents a hospital room
type Room struct {
	ID                 int            `json:"id" db:"id"`
	HealthcareEntityID int            `json:"healthcare_entity_id" db:"healthcare_entity_id" validate:"required"`
	RoomNumber         string         `json:"room_number" db:"room_number" validate:"required"`
	RoomName           sql.NullString `json:"room_name" db:"room_name"`
	RoomType           string         `json:"room_type" db:"room_type" validate:"required,oneof=consultation examination procedure operating emergency"`
	Floor              int            `json:"floor" db:"floor"`
	Department         sql.NullString `json:"department" db:"department"`
	Capacity           int            `json:"capacity" db:"capacity" validate:"min=1"`
	Equipment          sql.NullString `json:"equipment" db:"equipment"` // JSON string of equipment list
	IsActive           bool           `json:"is_active" db:"is_active"`
	Notes              sql.NullString `json:"notes" db:"notes"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
	CreatedBy          int            `json:"created_by" db:"created_by"`
}

// MarshalJSON provides custom JSON marshaling for Room
func (r Room) MarshalJSON() ([]byte, error) {
	type Alias Room
	return json.Marshal(&struct {
		*Alias
		RoomName   string `json:"room_name"`
		Department string `json:"department"`
		Equipment  string `json:"equipment"`
		Notes      string `json:"notes"`
	}{
		Alias:      (*Alias)(&r),
		RoomName:   r.RoomName.String,
		Department: r.Department.String,
		Equipment:  r.Equipment.String,
		Notes:      r.Notes.String,
	})
}

// RoomRequest represents room creation/update request
type RoomRequest struct {
	RoomNumber  string   `json:"room_number" validate:"required"`
	RoomName    string   `json:"room_name"`
	RoomType    string   `json:"room_type" validate:"required,oneof=consultation examination procedure operating emergency"`
	Floor       int      `json:"floor"`
	Department  string   `json:"department"`
	Capacity    int      `json:"capacity" validate:"min=1"`
	Equipment   []string `json:"equipment"` // Array of equipment
	Notes       string   `json:"notes"`
}

// RoomResponse represents room data returned to client
type RoomResponse struct {
	ID                 int      `json:"id"`
	HealthcareEntityID int      `json:"healthcare_entity_id"`
	RoomNumber         string   `json:"room_number"`
	RoomName           string   `json:"room_name"`
	RoomType           string   `json:"room_type"`
	Floor              int      `json:"floor"`
	Department         string   `json:"department"`
	Capacity           int      `json:"capacity"`
	Equipment          []string `json:"equipment"`
	IsActive           bool     `json:"is_active"`
	Notes              string   `json:"notes"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	IsAvailable        bool     `json:"is_available"` // Calculated field for availability checking
}

// RoomAvailability represents room availability for a specific time slot
type RoomAvailability struct {
	RoomID    int       `json:"room_id"`
	DateTime  time.Time `json:"date_time"`
	Duration  int       `json:"duration"`
	IsBooked  bool      `json:"is_booked"`
	BookedBy  string    `json:"booked_by,omitempty"` // Appointment or event details
}

// ToAppointmentDurationResponse converts AppointmentDurationSetting to response
func (ads *AppointmentDurationSetting) ToAppointmentDurationResponse() AppointmentDurationResponse {
	return AppointmentDurationResponse{
		ID:                 ads.ID,
		HealthcareEntityID: ads.HealthcareEntityID,
		AppointmentType:    ads.AppointmentType,
		DurationMinutes:    ads.DurationMinutes,
		IsDefault:          ads.IsDefault,
		IsActive:           ads.IsActive,
		CreatedAt:          ads.CreatedAt,
		UpdatedAt:          ads.UpdatedAt,
	}
}

// ToRoomResponse converts Room to response
func (r *Room) ToRoomResponse() RoomResponse {
	var equipment []string
	if r.Equipment.Valid && r.Equipment.String != "" {
		// Parse JSON equipment string to array
		// For simplicity, we'll split by comma for now
		// In production, you'd want proper JSON parsing
		equipment = strings.Split(r.Equipment.String, ",")
		for i, item := range equipment {
			equipment[i] = strings.TrimSpace(item)
		}
	}

	return RoomResponse{
		ID:                 r.ID,
		HealthcareEntityID: r.HealthcareEntityID,
		RoomNumber:         r.RoomNumber,
		RoomName:           r.RoomName.String,
		RoomType:           r.RoomType,
		Floor:              r.Floor,
		Department:         r.Department.String,
		Capacity:           r.Capacity,
		Equipment:          equipment,
		IsActive:           r.IsActive,
		Notes:              r.Notes.String,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
		IsAvailable:        true, // Will be calculated by service
	}
}

// GetRoomTypeIcon returns icon for room type
func (r *Room) GetRoomTypeIcon() string {
	roomTypeIcons := map[string]string{
		"consultation": "💬",
		"examination":  "🔍",
		"procedure":    "⚕️",
		"operating":    "🏥",
		"emergency":    "🚨",
	}
	if icon, exists := roomTypeIcons[r.RoomType]; exists {
		return icon
	}
	return "🏠"
}

// GetRoomTypeColor returns color class for room type
func (r *Room) GetRoomTypeColor() string {
	roomTypeColors := map[string]string{
		"consultation": "bg-blue-100 text-blue-800",
		"examination":  "bg-green-100 text-green-800",
		"procedure":    "bg-yellow-100 text-yellow-800",
		"operating":    "bg-red-100 text-red-800",
		"emergency":    "bg-purple-100 text-purple-800",
	}
	if color, exists := roomTypeColors[r.RoomType]; exists {
		return color
	}
	return "bg-gray-100 text-gray-800"
}

// TimezoneInfo represents healthcare entity timezone information
type TimezoneInfo struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Timezone string `json:"timezone" db:"timezone"`
	Country  string `json:"country" db:"country"`
}

// TimezoneConverter handles timezone conversions for healthcare entities
type TimezoneConverter struct {
	EntityTimezone string
}

// NewTimezoneConverter creates a new timezone converter for a healthcare entity
func NewTimezoneConverter(entityTimezone string) *TimezoneConverter {
	if entityTimezone == "" {
		entityTimezone = "UTC" // Default fallback
	}
	return &TimezoneConverter{
		EntityTimezone: entityTimezone,
	}
}

// ConvertUTCToEntity converts UTC time to entity's timezone
func (tc *TimezoneConverter) ConvertUTCToEntity(utcTime time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(tc.EntityTimezone)
	if err != nil {
		return utcTime, err
	}
	return utcTime.In(loc), nil
}

// ConvertEntityToUTC converts entity timezone to UTC
func (tc *TimezoneConverter) ConvertEntityToUTC(entityTime time.Time) time.Time {
	return entityTime.UTC()
}

// FormatDateInEntityTimezone formats a UTC date for display in entity timezone
func (tc *TimezoneConverter) FormatDateInEntityTimezone(utcTime time.Time, format string) (string, error) {
	entityTime, err := tc.ConvertUTCToEntity(utcTime)
	if err != nil {
		return utcTime.Format(format), err
	}
	return entityTime.Format(format), nil
}

// GetEntityDate returns the date in entity timezone as YYYY-MM-DD string
func (tc *TimezoneConverter) GetEntityDate(utcTime time.Time) (string, error) {
	entityTime, err := tc.ConvertUTCToEntity(utcTime)
	if err != nil {
		return utcTime.Format("2006-01-02"), err
	}
	return entityTime.Format("2006-01-02"), nil
}

// ParseEntityDateTime parses a datetime string in entity timezone and converts to UTC
func (tc *TimezoneConverter) ParseEntityDateTime(dateTimeStr string) (time.Time, error) {
	loc, err := time.LoadLocation(tc.EntityTimezone)
	if err != nil {
		return time.Time{}, err
	}
	
	// Try parsing different formats
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}
	
	for _, format := range formats {
		if t, err := time.ParseInLocation(format, dateTimeStr, loc); err == nil {
			return t.UTC(), nil
		}
	}
	
	return time.Time{}, err
}