package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type AppointmentService struct {
	db             *sql.DB
	timezoneCache  map[int]*TimezoneConverter // Cache for entity timezone converters
}

func NewAppointmentService(db *sql.DB) *AppointmentService {
	return &AppointmentService{
		db:            db,
		timezoneCache: make(map[int]*TimezoneConverter),
	}
}

// CreateAppointment creates a new appointment with conflict checking
func (s *AppointmentService) CreateAppointment(appointment *Appointment) error {
	// Get room ID from appointment
	var roomID int
	if appointment.RoomID.Valid {
		roomID = int(appointment.RoomID.Int32)
	}

	// Check for conflicts first
	hasConflict, err := s.CheckConflict(ConflictCheck{
		DoctorID:           appointment.DoctorID,
		DateTime:           appointment.DateTime,
		Duration:           appointment.Duration,
		RoomID:             roomID,
		HealthcareEntityID: appointment.HealthcareEntityID,
	})
	if err != nil {
		return err
	}
	if hasConflict {
		// Determine the type of conflict for a better error message
		doctorConflict, _ := s.checkDoctorConflict(ConflictCheck{
			DoctorID:           appointment.DoctorID,
			DateTime:           appointment.DateTime,
			Duration:           appointment.Duration,
			HealthcareEntityID: appointment.HealthcareEntityID,
		})
		
		if doctorConflict {
			return errors.New("doctor is not available at this time")
		} else if roomID > 0 {
			return errors.New("room is not available at this time")
		} else {
			return errors.New("appointment conflicts with existing appointment")
		}
	}

	query := `
		INSERT INTO appointments (
			healthcare_entity_id, patient_id, doctor_id, date_time, duration, type, status,
			reason, notes, priority, room_id, is_active, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) RETURNING id, created_at, updated_at
	`
	
	now := time.Now()
	appointment.IsActive = true
	if appointment.Status == "" {
		appointment.Status = "scheduled"
	}
	if appointment.Priority == "" {
		appointment.Priority = "normal"
	}
	appointment.CreatedAt = now
	appointment.UpdatedAt = now

	err = s.db.QueryRow(
		query,
		appointment.HealthcareEntityID,
		appointment.PatientID,
		appointment.DoctorID,
		appointment.DateTime,
		appointment.Duration,
		appointment.Type,
		appointment.Status,
		appointment.Reason,
		appointment.Notes,
		appointment.Priority,
		appointment.RoomID,
		appointment.IsActive,
		appointment.CreatedAt,
		appointment.UpdatedAt,
		appointment.CreatedBy,
	).Scan(&appointment.ID, &appointment.CreatedAt, &appointment.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetAppointmentByID gets appointment by ID
func (s *AppointmentService) GetAppointmentByID(id int) (*Appointment, error) {
	appointment := &Appointment{}
	query := `
		SELECT 
			id, healthcare_entity_id, patient_id, doctor_id, date_time, duration, type, status,
			reason, notes, priority, room_id, is_active, created_at, updated_at, created_by
		FROM appointments
		WHERE id = $1 AND is_active = true
	`

	err := s.db.QueryRow(query, id).Scan(
		&appointment.ID,
		&appointment.HealthcareEntityID,
		&appointment.PatientID,
		&appointment.DoctorID,
		&appointment.DateTime,
		&appointment.Duration,
		&appointment.Type,
		&appointment.Status,
		&appointment.Reason,
		&appointment.Notes,
		&appointment.Priority,
		&appointment.RoomID,
		&appointment.IsActive,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
		&appointment.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}

	return appointment, nil
}

// UpdateAppointment updates appointment information
func (s *AppointmentService) UpdateAppointment(appointment *Appointment) error {
	// Get room ID from appointment
	roomID := 0
	if appointment.RoomID.Valid {
		roomID = int(appointment.RoomID.Int32)
	}

	// Check for conflicts if date/time or duration changed
	hasConflict, err := s.CheckConflict(ConflictCheck{
		DoctorID:           appointment.DoctorID,
		DateTime:           appointment.DateTime,
		Duration:           appointment.Duration,
		ExcludeID:          appointment.ID,
		RoomID:             roomID,
		HealthcareEntityID: appointment.HealthcareEntityID,
	})
	if err != nil {
		return err
	}
	if hasConflict {
		return errors.New("appointment conflicts with existing appointment")
	}

	query := `
		UPDATE appointments SET
			patient_id = $1, doctor_id = $2, date_time = $3, duration = $4,
			type = $5, status = $6, reason = $7, notes = $8, priority = $9, room_id = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $11 AND is_active = true
		RETURNING updated_at
	`

	err = s.db.QueryRow(
		query,
		appointment.PatientID,
		appointment.DoctorID,
		appointment.DateTime,
		appointment.Duration,
		appointment.Type,
		appointment.Status,
		appointment.Reason,
		appointment.Notes,
		appointment.Priority,
		appointment.RoomID,
		appointment.ID,
	).Scan(&appointment.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("appointment not found")
		}
		return err
	}

	return nil
}

// UpdateAppointmentStatus updates only the appointment status and notes
func (s *AppointmentService) UpdateAppointmentStatus(id int, status, notes string) error {
	query := `
		UPDATE appointments SET
			status = $1, notes = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND is_active = true
	`

	result, err := s.db.Exec(query, status, notes, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("appointment not found")
	}

	return nil
}

// GetAppointments gets appointments with filtering and pagination
func (s *AppointmentService) GetAppointments(search AppointmentSearch) ([]Appointment, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Base query
	baseQuery := `
		SELECT 
			id, healthcare_entity_id, patient_id, doctor_id, date_time, duration, type, status,
			reason, notes, priority, room_id, is_active, created_at, updated_at, created_by
		FROM appointments
		WHERE is_active = true AND healthcare_entity_id = $1
	`

	// Start with healthcare entity filter
	args = append(args, search.HealthcareEntityID)
	argIndex++

	// Add additional filter conditions
	if search.PatientID > 0 {
		conditions = append(conditions, fmt.Sprintf("patient_id = $%d", argIndex))
		args = append(args, search.PatientID)
		argIndex++
	}

	if search.DoctorID > 0 {
		conditions = append(conditions, fmt.Sprintf("doctor_id = $%d", argIndex))
		args = append(args, search.DoctorID)
		argIndex++
	}

	if search.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, search.Status)
		argIndex++
	}

	if search.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, search.Type)
		argIndex++
	}

	if !search.DateFrom.IsZero() {
		conditions = append(conditions, fmt.Sprintf("date_time >= $%d", argIndex))
		args = append(args, search.DateFrom)
		argIndex++
	}

	if !search.DateTo.IsZero() {
		conditions = append(conditions, fmt.Sprintf("date_time <= $%d", argIndex))
		args = append(args, search.DateTo)
		argIndex++
	}

	// Build final query
	finalQuery := baseQuery
	if len(conditions) > 0 {
		finalQuery += " AND " + strings.Join(conditions, " AND ")
	}

	finalQuery += " ORDER BY date_time ASC"

	// Add pagination
	if search.Limit <= 0 || search.Limit > 100 {
		search.Limit = 20
	}
	if search.Offset < 0 {
		search.Offset = 0
	}

	finalQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, search.Limit, search.Offset)

	// Execute query
	rows, err := s.db.Query(finalQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var appointment Appointment
		err := rows.Scan(
			&appointment.ID,
			&appointment.HealthcareEntityID,
			&appointment.PatientID,
			&appointment.DoctorID,
			&appointment.DateTime,
			&appointment.Duration,
			&appointment.Type,
			&appointment.Status,
			&appointment.Reason,
			&appointment.Notes,
			&appointment.Priority,
			&appointment.RoomID,
			&appointment.IsActive,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
			&appointment.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

// DeleteAppointment soft deletes an appointment
func (s *AppointmentService) DeleteAppointment(id int) error {
	query := `
		UPDATE appointments
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_active = true
	`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("appointment not found")
	}

	return nil
}

// CheckConflict checks if an appointment conflicts with existing ones
func (s *AppointmentService) CheckConflict(check ConflictCheck) (bool, error) {
	// Check doctor conflict
	doctorConflict, err := s.checkDoctorConflict(check)
	if err != nil {
		return false, err
	}
	if doctorConflict {
		return true, nil
	}

	// Check room conflict if room is specified
	if check.RoomID > 0 {
		roomConflict, err := s.checkRoomConflict(check)
		if err != nil {
			return false, err
		}
		if roomConflict {
			return true, nil
		}
	}

	return false, nil
}

func (s *AppointmentService) checkDoctorConflict(check ConflictCheck) (bool, error) {
	endTime := check.DateTime.Add(time.Duration(check.Duration) * time.Minute)
	
	query := `
		SELECT COUNT(*)
		FROM appointments
		WHERE doctor_id = $1
		AND healthcare_entity_id = $2
		AND is_active = true
		AND status IN ('scheduled', 'confirmed', 'in-progress')
		AND id != $3
		AND (
			(date_time < $4 AND date_time + (duration || ' minutes')::interval > $5) OR
			(date_time < $6 AND date_time + (duration || ' minutes')::interval > $4)
		)
	`

	var count int
	err := s.db.QueryRow(
		query,
		check.DoctorID,
		check.HealthcareEntityID,
		check.ExcludeID,
		endTime,
		check.DateTime,
		check.DateTime,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *AppointmentService) checkRoomConflict(check ConflictCheck) (bool, error) {
	endTime := check.DateTime.Add(time.Duration(check.Duration) * time.Minute)
	
	query := `
		SELECT COUNT(*)
		FROM appointments
		WHERE room_id = $1
		AND healthcare_entity_id = $2
		AND is_active = true
		AND status IN ('scheduled', 'confirmed', 'in-progress')
		AND id != $3
		AND (
			(date_time < $4 AND date_time + (duration || ' minutes')::interval > $5) OR
			(date_time < $6 AND date_time + (duration || ' minutes')::interval > $4)
		)
	`

	var count int
	err := s.db.QueryRow(
		query,
		check.RoomID,
		check.HealthcareEntityID,
		check.ExcludeID,
		endTime,
		check.DateTime,
		check.DateTime,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetDoctorSchedule gets doctor's schedule for a specific date
func (s *AppointmentService) GetDoctorSchedule(doctorID int, date time.Time) (*DoctorSchedule, error) {
	// Get start and end of the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Get appointments for the day
	search := AppointmentSearch{
		DoctorID: doctorID,
		DateFrom: startOfDay,
		DateTo:   endOfDay,
		Limit:    100, // High limit for a single day
	}

	appointments, err := s.GetAppointments(search)
	if err != nil {
		return nil, err
	}

	// Filter active appointments
	var activeAppointments []Appointment
	for _, appointment := range appointments {
		if appointment.Status == "scheduled" || appointment.Status == "confirmed" || appointment.Status == "in-progress" {
			activeAppointments = append(activeAppointments, appointment)
		}
	}

	// Get working hours (using default for now)
	workingHours := GetDefaultWorkingHours(date)

	// Generate available slots
	availableSlots := s.generateAvailableSlots(workingHours, activeAppointments)

	schedule := &DoctorSchedule{
		DoctorID:       doctorID,
		Date:           date,
		WorkingHours:   workingHours,
		Appointments:   activeAppointments,
		AvailableSlots: availableSlots,
	}

	return schedule, nil
}

// generateAvailableSlots generates available time slots based on working hours and existing appointments
func (s *AppointmentService) generateAvailableSlots(workingHours WorkingHours, appointments []Appointment) []AvailabilitySlot {
	var slots []AvailabilitySlot
	slotDuration := 30 // 30-minute slots

	// Create time slots from start to end of working hours
	current := workingHours.StartTime
	
	for current.Before(workingHours.EndTime) {
		slotEnd := current.Add(time.Duration(slotDuration) * time.Minute)
		
		// Skip if slot overlaps with break time
		if !(current.Before(workingHours.BreakStart) || current.After(workingHours.BreakEnd)) {
			current = slotEnd
			continue
		}
		
		// Check if slot conflicts with any appointment
		isAvailable := true
		for _, appointment := range appointments {
			appointmentEnd := appointment.DateTime.Add(time.Duration(appointment.Duration) * time.Minute)
			
			// Check for overlap
			if current.Before(appointmentEnd) && appointment.DateTime.Before(slotEnd) {
				isAvailable = false
				break
			}
		}
		
		if isAvailable && slotEnd.Before(workingHours.EndTime) {
			slots = append(slots, AvailabilitySlot{
				DateTime: current,
				Duration: slotDuration,
			})
		}
		
		current = slotEnd
	}

	return slots
}

// GetAppointmentStats gets appointment statistics
func (s *AppointmentService) GetAppointmentStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total appointments
	var totalCount int
	err := s.db.QueryRow("SELECT COUNT(*) FROM appointments WHERE is_active = true").Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	stats["total_appointments"] = totalCount

	// Appointments by status
	statusQuery := `
		SELECT status, COUNT(*) 
		FROM appointments 
		WHERE is_active = true 
		GROUP BY status
	`
	rows, err := s.db.Query(statusQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statusCounts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		statusCounts[status] = count
	}
	stats["by_status"] = statusCounts

	// Today's appointments
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var todayCount int
	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM appointments WHERE is_active = true AND date_time >= $1 AND date_time < $2",
		startOfDay, endOfDay,
	).Scan(&todayCount)
	if err != nil {
		return nil, err
	}
	stats["today_appointments"] = todayCount

	return stats, nil
}

// Doctor Availability Service Methods

// CreateDoctorAvailability creates or updates doctor availability using UTC timestamps
func (s *AppointmentService) CreateDoctorAvailability(availability *DoctorAvailability) error {
	query := `
		INSERT INTO doctor_availability (
			healthcare_entity_id, doctor_id, status,
			start_datetime, end_datetime, break_start_datetime, break_end_datetime,
			notes, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id, created_at, updated_at
	`

	now := time.Now()
	availability.CreatedAt = now
	availability.UpdatedAt = now

	err := s.db.QueryRow(
		query,
		availability.HealthcareEntityID,
		availability.DoctorID,
		availability.Status,
		availability.StartDateTime,
		availability.EndDateTime,
		availability.BreakStartDateTime,
		availability.BreakEndDateTime,
		availability.Notes,
		now,
		now,
		availability.CreatedBy,
	).Scan(&availability.ID, &availability.CreatedAt, &availability.UpdatedAt)

	return err
}

// GetDoctorAvailability gets doctor availability for a date range
func (s *AppointmentService) GetDoctorAvailability(search AvailabilitySearch) ([]DoctorAvailabilityResponse, error) {

	query := `
		SELECT 
			da.id, da.doctor_id, da.status,
			da.start_datetime, da.end_datetime, da.break_start_datetime, da.break_end_datetime,
			da.notes, da.created_at, da.updated_at
		FROM doctor_availability da
		WHERE da.healthcare_entity_id = $1
	`
	args := []interface{}{search.HealthcareEntityID}
	argCount := 1

	if search.DoctorID > 0 {
		argCount++
		query += fmt.Sprintf(" AND da.doctor_id = $%d", argCount)
		args = append(args, search.DoctorID)
	}

	if search.DateFrom != "" {
		argCount++
		query += fmt.Sprintf(" AND DATE(da.start_datetime) >= $%d", argCount)
		args = append(args, search.DateFrom)
	}

	if search.DateTo != "" {
		argCount++
		query += fmt.Sprintf(" AND DATE(da.start_datetime) <= $%d", argCount)
		args = append(args, search.DateTo)
	}

	if search.Status != "" {
		argCount++
		query += fmt.Sprintf(" AND da.status = $%d", argCount)
		args = append(args, search.Status)
	}

	query += " ORDER BY da.start_datetime DESC, da.doctor_id"

	if search.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, search.Limit)
	}

	if search.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, search.Offset)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var availabilities []DoctorAvailabilityResponse
	doctorCache := make(map[int]string) // Cache doctor names to avoid multiple API calls
	
	for rows.Next() {
		var availability DoctorAvailabilityResponse
		var startDateTime, endDateTime, breakStartDateTime, breakEndDateTime *time.Time
		
		err := rows.Scan(
			&availability.ID,
			&availability.DoctorID,
			&availability.Status,
			&startDateTime,
			&endDateTime,
			&breakStartDateTime,
			&breakEndDateTime,
			&availability.Notes,
			&availability.CreatedAt,
			&availability.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set the UTC datetime fields
		availability.StartDateTime = startDateTime
		availability.EndDateTime = endDateTime
		availability.BreakStartDateTime = breakStartDateTime
		availability.BreakEndDateTime = breakEndDateTime
		
		// Get doctor name from cache or fetch from user service
		doctorName, exists := doctorCache[availability.DoctorID]
		if !exists {
			doctor, err := s.GetDoctorByID(availability.DoctorID, search.HealthcareEntityID)
			if err != nil {
				doctorName = fmt.Sprintf("Doctor %d", availability.DoctorID)
			} else {
				doctorName = fmt.Sprintf("Dr. %s %s", doctor.FirstName, doctor.LastName)
			}
			doctorCache[availability.DoctorID] = doctorName
		}
		availability.DoctorName = doctorName
		
		availabilities = append(availabilities, availability)
	}

	return availabilities, nil
}

// GetDoctorAvailabilityByID gets a specific doctor availability record
func (s *AppointmentService) GetDoctorAvailabilityByID(id int) (*DoctorAvailability, error) {
	query := `
		SELECT id, healthcare_entity_id, doctor_id, status,
		       start_datetime, end_datetime, break_start_datetime, break_end_datetime,
		       notes, created_at, updated_at, created_by
		FROM doctor_availability
		WHERE id = $1
	`

	var availability DoctorAvailability
	err := s.db.QueryRow(query, id).Scan(
		&availability.ID,
		&availability.HealthcareEntityID,
		&availability.DoctorID,
		&availability.Status,
		&availability.StartDateTime,
		&availability.EndDateTime,
		&availability.BreakStartDateTime,
		&availability.BreakEndDateTime,
		&availability.Notes,
		&availability.CreatedAt,
		&availability.UpdatedAt,
		&availability.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("availability record not found")
		}
		return nil, err
	}

	return &availability, nil
}

// UpdateDoctorAvailability updates doctor availability using UTC timestamps
func (s *AppointmentService) UpdateDoctorAvailability(availability *DoctorAvailability) error {
	query := `
		UPDATE doctor_availability
		SET status = $1, start_datetime = $2, end_datetime = $3,
		    break_start_datetime = $4, break_end_datetime = $5, notes = $6,
		    updated_at = $7
		WHERE id = $8
	`

	result, err := s.db.Exec(
		query,
		availability.Status,
		availability.StartDateTime,
		availability.EndDateTime,
		availability.BreakStartDateTime,
		availability.BreakEndDateTime,
		availability.Notes,
		time.Now(),
		availability.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("availability record not found")
	}

	return nil
}

// DeleteDoctorAvailability soft deletes doctor availability
func (s *AppointmentService) DeleteDoctorAvailability(id int) error {
	query := `DELETE FROM doctor_availability WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("availability record not found")
	}

	return nil
}

// GetDoctorsByEntity gets all doctors for a healthcare entity by calling the user service
func (s *AppointmentService) GetDoctorsByEntity(healthcareEntityID int) ([]DoctorInfo, error) {
	return s.fetchDoctorsFromUserService(healthcareEntityID)
}

// GetDoctorByID gets a specific doctor by ID from the user service
func (s *AppointmentService) GetDoctorByID(doctorID int, healthcareEntityID int) (*DoctorInfo, error) {
	doctors, err := s.fetchDoctorsFromUserService(healthcareEntityID) // Use the entity ID from the authenticated context
	if err != nil {
		return nil, err
	}
	
	for _, doctor := range doctors {
		if doctor.ID == doctorID {
			return &doctor, nil
		}
	}
	
	return nil, errors.New("doctor not found")
}

// fetchDoctorsFromUserService makes HTTP call to user service to get doctor information
func (s *AppointmentService) fetchDoctorsFromUserService(healthcareEntityID int) ([]DoctorInfo, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://user-service:8081"
	}
	
	// Build URL - if healthcareEntityID is 0, get all doctors
	url := fmt.Sprintf("%s/api/internal/doctors", userServiceURL)
	if healthcareEntityID > 0 {
		url = fmt.Sprintf("%s?healthcare_entity_id=%d", url, healthcareEntityID)
	}
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call user service: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned status %d", resp.StatusCode)
	}
	
	var response struct {
		Data []struct {
			ID             int    `json:"id"`
			FirstName      string `json:"first_name"`
			LastName       string `json:"last_name"`
			Email          string `json:"email"`
			Specialization string `json:"specialization"`
			Role           string `json:"role"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode user service response: %v", err)
	}
	
	var doctors []DoctorInfo
	for _, user := range response.Data {
		if user.Role == "doctor" {
			doctors = append(doctors, DoctorInfo{
				ID:             user.ID,
				FirstName:      user.FirstName,
				LastName:       user.LastName,
				Email:          user.Email,
				Specialization: user.Specialization,
				Role:           user.Role,
			})
		}
	}
	
	return doctors, nil
}

// GetTimezoneConverter gets or creates a timezone converter for a healthcare entity
func (s *AppointmentService) GetTimezoneConverter(healthcareEntityID int) (*TimezoneConverter, error) {
	// Check cache first
	if converter, exists := s.timezoneCache[healthcareEntityID]; exists {
		return converter, nil
	}
	
	// Fetch timezone from user service
	timezoneInfo, err := s.fetchEntityTimezone(healthcareEntityID)
	if err != nil {
		// Fallback to UTC if we can't fetch timezone
		converter := NewTimezoneConverter("UTC")
		s.timezoneCache[healthcareEntityID] = converter
		return converter, nil
	}
	
	// Create and cache the converter
	converter := NewTimezoneConverter(timezoneInfo.Timezone)
	s.timezoneCache[healthcareEntityID] = converter
	return converter, nil
}

// fetchEntityTimezone fetches timezone information from user service
func (s *AppointmentService) fetchEntityTimezone(healthcareEntityID int) (*TimezoneInfo, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://user-service:8081"
	}
	
	url := fmt.Sprintf("%s/api/internal/entity/%d", userServiceURL, healthcareEntityID)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call user service for timezone: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned status %d for timezone", resp.StatusCode)
	}
	
	var response struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Timezone string `json:"timezone"`
		Country  string `json:"country"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode timezone response: %v", err)
	}
	
	return &TimezoneInfo{
		ID:       response.ID,
		Name:     response.Name,
		Timezone: response.Timezone,
		Country:  response.Country,
	}, nil
}

// GetAvailabilityCalendar gets calendar view of doctor availability
func (s *AppointmentService) GetAvailabilityCalendar(healthcareEntityID, doctorID int, yearMonth string) (map[string]interface{}, error) {
	// Parse year-month (YYYY-MM)
	dateFrom := yearMonth + "-01"
	
	// Calculate end date (last day of month)
	parsedDate, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return nil, err
	}
	
	// Get first day of next month and subtract a day
	nextMonth := parsedDate.AddDate(0, 1, 0)
	lastDay := nextMonth.AddDate(0, 0, -1)
	dateTo := lastDay.Format("2006-01-02")

	search := AvailabilitySearch{
		HealthcareEntityID: healthcareEntityID,
		DoctorID:           doctorID,
		DateFrom:           dateFrom,
		DateTo:             dateTo,
		Limit:              100,
	}

	availabilities, err := s.GetDoctorAvailability(search)
	if err != nil {
		return nil, err
	}

	// Organize by date
	calendarData := make(map[string]interface{})
	days := make(map[string]DoctorAvailabilityResponse)

	for _, avail := range availabilities {
		// Extract date from start_datetime for calendar organization
		var dateKey string
		if avail.StartDateTime != nil {
			dateKey = avail.StartDateTime.Format("2006-01-02")
		} else {
			// Fallback - skip entries without datetime
			continue
		}
		days[dateKey] = avail
	}

	calendarData["month"] = yearMonth
	calendarData["days"] = days
	calendarData["doctor_id"] = doctorID

	return calendarData, nil
}

// CreateBulkAvailability creates availability for multiple dates using a template with UTC timestamps
func (s *AppointmentService) CreateBulkAvailability(doctorID, healthcareEntityID, createdBy int, dateFrom, dateTo string, template DoctorAvailabilityRequest) error {
	startDate, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return err
	}

	endDate, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		return err
	}

	// Parse template datetime strings to UTC timestamps
	startDateTime, err := time.Parse(time.RFC3339, template.StartDateTime)
	if err != nil {
		return fmt.Errorf("invalid start datetime: %v", err)
	}

	endDateTime, err := time.Parse(time.RFC3339, template.EndDateTime)
	if err != nil {
		return fmt.Errorf("invalid end datetime: %v", err)
	}

	var breakStartDateTime, breakEndDateTime *time.Time
	if template.BreakStartDateTime != "" {
		breakStart, err := time.Parse(time.RFC3339, template.BreakStartDateTime)
		if err != nil {
			return fmt.Errorf("invalid break start datetime: %v", err)
		}
		breakStartDateTime = &breakStart
	}
	if template.BreakEndDateTime != "" {
		breakEnd, err := time.Parse(time.RFC3339, template.BreakEndDateTime)
		if err != nil {
			return fmt.Errorf("invalid break end datetime: %v", err)
		}
		breakEndDateTime = &breakEnd
	}

	// Create availability for each date
	current := startDate
	for !current.After(endDate) {
		// Calculate date offset from the original template date
		dayOffset := int(current.Sub(startDate).Hours() / 24)
		
		// Adjust the template timestamps for this date
		dateStartDateTime := startDateTime.AddDate(0, 0, dayOffset)
		dateEndDateTime := endDateTime.AddDate(0, 0, dayOffset)
		
		var dateBreakStartDateTime, dateBreakEndDateTime *time.Time
		if breakStartDateTime != nil {
			tempBreakStart := breakStartDateTime.AddDate(0, 0, dayOffset)
			dateBreakStartDateTime = &tempBreakStart
		}
		if breakEndDateTime != nil {
			tempBreakEnd := breakEndDateTime.AddDate(0, 0, dayOffset)
			dateBreakEndDateTime = &tempBreakEnd
		}
		
		availability := &DoctorAvailability{
			HealthcareEntityID:    healthcareEntityID,
			DoctorID:              doctorID,
			Status:                template.Status,
			StartDateTime:         &dateStartDateTime,
			EndDateTime:           &dateEndDateTime,
			BreakStartDateTime:    dateBreakStartDateTime,
			BreakEndDateTime:      dateBreakEndDateTime,
			Notes:                 template.Notes,
			CreatedBy:             createdBy,
		}

		err := s.CreateDoctorAvailability(availability)
		if err != nil {
			return err
		}

		current = current.AddDate(0, 0, 1)
	}

	return nil
}

// BookAppointmentWithConflictCheck performs smart appointment booking with conflict detection
func (s *AppointmentService) BookAppointmentWithConflictCheck(request BookingRequest, healthcareEntityID, userID int) (*BookingResponse, error) {
	// Parse UTC datetime string (same as availability)
	dateTime, err := time.Parse(time.RFC3339, request.DateTime)
	if err != nil {
		return &BookingResponse{
			Success: false,
			Message: "Invalid datetime format. DateTime must be in ISO 8601 UTC format (2006-01-02T15:04:05Z)",
		}, nil
	}

	// Check doctor availability first - extract date from UTC datetime
	dateStr := dateTime.Format("2006-01-02")
	availability, err := s.GetDoctorAvailabilityForDate(request.DoctorID, healthcareEntityID, dateStr)
	if err != nil {
		return &BookingResponse{
			Success: false,
			Message: "Failed to check doctor availability",
		}, err
	}

	var conflicts []ConflictInfo
	var alternatives []AvailabilitySlot

	// Check if doctor is available that day
	if availability == nil || !availability.IsAvailable() {
		conflict := ConflictInfo{
			ConflictType: "doctor_unavailable",
			ConflictTime: dateTime,
			ConflictEnd:  dateTime.Add(time.Duration(request.Duration) * time.Minute),
			Description:  "Doctor is not available on this date",
		}
		conflicts = append(conflicts, conflict)
		
		// Generate alternative slots for nearby dates
		alternatives = s.generateAlternativeSlots(request.DoctorID, healthcareEntityID, dateTime, request.Duration)
		
		return &BookingResponse{
			Success:          false,
			Conflicts:        conflicts,
			AlternativeSlots: alternatives,
			Message:          "Doctor is not available on the selected date. Please choose an alternative time slot.",
		}, nil
	}

	// Check for appointment conflicts
	if request.CheckConflicts {
		hasConflict, err := s.CheckConflict(ConflictCheck{
			DoctorID:           request.DoctorID,
			DateTime:           dateTime,
			Duration:           request.Duration,
			RoomID:             request.RoomID,
			HealthcareEntityID: healthcareEntityID,
		})
		if err != nil {
			return &BookingResponse{
				Success: false,
				Message: "Failed to check for conflicts",
			}, err
		}

		if hasConflict {
			// Get conflicting appointments
			existingAppointments, err := s.getConflictingAppointments(request.DoctorID, dateTime, request.Duration)
			if err == nil && len(existingAppointments) > 0 {
				for _, existing := range existingAppointments {
					conflict := ConflictInfo{
						ConflictType:        "doctor_busy",
						ExistingAppointment: &existing,
						ConflictTime:        existing.DateTime,
						ConflictEnd:         existing.DateTime.Add(time.Duration(existing.Duration) * time.Minute),
						Description:         "Doctor has another appointment at this time",
					}
					conflicts = append(conflicts, conflict)
				}
			}

			// Generate alternative slots
			alternatives = s.generateAlternativeSlots(request.DoctorID, healthcareEntityID, dateTime, request.Duration)
			
			return &BookingResponse{
				Success:          false,
				Conflicts:        conflicts,
				AlternativeSlots: alternatives,
				Message:          "Time slot is not available. Please choose an alternative time.",
			}, nil
		}
	}

	// Create the appointment
	appointment := &Appointment{
		HealthcareEntityID: healthcareEntityID,
		PatientID:          request.PatientID,
		DoctorID:           request.DoctorID,
		DateTime:           dateTime,
		Duration:           request.Duration,
		Type:               request.Type,
		Reason:             request.Reason,
		Notes:              request.Notes,
		Priority:           request.Priority,
		RoomID:             sql.NullInt32{Int32: int32(request.RoomID), Valid: request.RoomID > 0},
		CreatedBy:          userID,
	}

	err = s.CreateAppointment(appointment)
	if err != nil {
		return &BookingResponse{
			Success: false,
			Message: "Failed to create appointment: " + err.Error(),
		}, err
	}

	appointmentResponse := appointment.ToAppointmentResponse()
	
	return &BookingResponse{
		Success:       true,
		AppointmentID: appointment.ID,
		Appointment:   &appointmentResponse,
		Message:       "Appointment booked successfully",
	}, nil
}


// GetAvailableTimeSlots generates available time slots based on doctor availability
func (s *AppointmentService) GetAvailableTimeSlots(doctorID, healthcareEntityID int, date string, duration int) ([]AvailabilitySlot, error) {
	log.Printf("GetAvailableTimeSlots called for doctor %d, entity %d, date %s", doctorID, healthcareEntityID, date)
	
	// Get doctor availability for the date
	availability, err := s.GetDoctorAvailabilityForDate(doctorID, healthcareEntityID, date)
	if err != nil {
		log.Printf("Error getting availability: %v", err)
		return nil, err
	}

	// If doctor is not available, return empty slots
	if availability == nil || !availability.IsAvailable() {
		log.Printf("Doctor not available or availability is nil")
		return []AvailabilitySlot{}, nil
	}

	// Debug: Check what we got from the database
	log.Printf("DEBUG: StartDateTime: %v, EndDateTime: %v", availability.StartDateTime, availability.EndDateTime)

	// Ensure we have UTC datetime fields
	if availability.StartDateTime == nil || availability.EndDateTime == nil {
		log.Printf("DEBUG: No UTC datetime fields available, returning empty slots")
		return []AvailabilitySlot{}, nil
	}

	// Parse the target date for comparison
	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	// Get start and end of day in UTC for appointment search
	startOfDayUTC := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDayUTC := startOfDayUTC.Add(24 * time.Hour)

	// Get existing appointments for the day
	search := AppointmentSearch{
		HealthcareEntityID: healthcareEntityID,
		DoctorID:           doctorID,
		DateFrom:           startOfDayUTC,
		DateTo:             endOfDayUTC,
		Limit:              100,
	}

	appointments, err := s.GetAppointments(search)
	if err != nil {
		return nil, err
	}

	// Use the UTC working hours directly
	workingStart := *availability.StartDateTime
	workingEnd := *availability.EndDateTime
	
	var breakStart, breakEnd time.Time
	if availability.BreakStartDateTime != nil && availability.BreakEndDateTime != nil {
		breakStart = *availability.BreakStartDateTime
		breakEnd = *availability.BreakEndDateTime
	}

	// Generate time slots
	var slots []AvailabilitySlot
	slotDuration := time.Duration(duration) * time.Minute
	current := workingStart

	for current.Add(slotDuration).Before(workingEnd) || current.Add(slotDuration).Equal(workingEnd) {
		slotEnd := current.Add(slotDuration)
		
		// Skip break time
		if !breakStart.IsZero() && !breakEnd.IsZero() {
			if (current.Before(breakEnd) && slotEnd.After(breakStart)) {
				current = breakEnd
				continue
			}
		}
		
		// Check for conflicts with existing appointments
		isAvailable := true
		for _, appointment := range appointments {
			if appointment.Status == "cancelled" || appointment.Status == "no-show" {
				continue
			}
			
			appointmentEnd := appointment.DateTime.Add(time.Duration(appointment.Duration) * time.Minute)
			
			// Check for overlap (all times are in UTC)
			if current.Before(appointmentEnd) && appointment.DateTime.Before(slotEnd) {
				isAvailable = false
				break
			}
		}
		
		// Determine slot type based on UTC hour (will be converted to local time by frontend)
		slotType := "morning"
		hour := current.Hour()
		if hour >= 12 && hour < 17 {
			slotType = "afternoon" 
		} else if hour >= 17 {
			slotType = "evening"
		}
		
		// Store UTC time directly (no conversion needed)
		slots = append(slots, AvailabilitySlot{
			DateTime:    current,
			Duration:    duration,
			IsAvailable: isAvailable,
			SlotType:    slotType,
		})
		
		current = current.Add(time.Duration(30) * time.Minute) // 30-minute increments
	}

	return slots, nil
}

// GetDoctorAvailabilityForDate gets doctor availability for a specific date
func (s *AppointmentService) GetDoctorAvailabilityForDate(doctorID, healthcareEntityID int, date string) (*DoctorAvailability, error) {
	query := `
		SELECT id, healthcare_entity_id, doctor_id, status,
		       start_datetime, end_datetime, break_start_datetime, break_end_datetime,
		       notes, created_at, updated_at, created_by
		FROM doctor_availability
		WHERE doctor_id = $1 AND healthcare_entity_id = $2 AND DATE(start_datetime) = $3
	`

	var availability DoctorAvailability
	err := s.db.QueryRow(query, doctorID, healthcareEntityID, date).Scan(
		&availability.ID,
		&availability.HealthcareEntityID,
		&availability.DoctorID,
		&availability.Status,
		&availability.StartDateTime,
		&availability.EndDateTime,
		&availability.BreakStartDateTime,
		&availability.BreakEndDateTime,
		&availability.Notes,
		&availability.CreatedAt,
		&availability.UpdatedAt,
		&availability.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No availability record found
		}
		return nil, err
	}

	return &availability, nil
}

// generateAlternativeSlots generates alternative time slots when conflicts occur
func (s *AppointmentService) generateAlternativeSlots(doctorID, healthcareEntityID int, requestedTime time.Time, duration int) []AvailabilitySlot {
	var alternatives []AvailabilitySlot
	
	// Check next 7 days for alternatives
	for i := 0; i < 7; i++ {
		checkDate := requestedTime.AddDate(0, 0, i)
		dateStr := checkDate.Format("2006-01-02")
		
		slots, err := s.GetAvailableTimeSlots(doctorID, healthcareEntityID, dateStr, duration)
		if err != nil {
			continue
		}
		
		// Take first 3 available slots per day
		availableCount := 0
		for _, slot := range slots {
			if slot.IsAvailable && availableCount < 3 {
				alternatives = append(alternatives, slot)
				availableCount++
			}
		}
		
		// Limit total alternatives to 10
		if len(alternatives) >= 10 {
			break
		}
	}
	
	return alternatives
}

// getConflictingAppointments gets appointments that conflict with the given time slot
func (s *AppointmentService) getConflictingAppointments(doctorID int, dateTime time.Time, duration int) ([]AppointmentResponse, error) {
	endTime := dateTime.Add(time.Duration(duration) * time.Minute)
	
	query := `
		SELECT 
			id, healthcare_entity_id, patient_id, doctor_id, date_time, duration, type, status,
			reason, notes, priority, room_id, created_at, updated_at, created_by, is_active
		FROM appointments
		WHERE doctor_id = $1
		AND is_active = true
		AND status IN ('scheduled', 'confirmed', 'in-progress')
		AND (
			(date_time < $2 AND date_time + (duration || ' minutes')::interval > $3) OR
			(date_time < $4 AND date_time + (duration || ' minutes')::interval > $2)
		)
	`

	rows, err := s.db.Query(query, doctorID, endTime, dateTime, dateTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []AppointmentResponse
	for rows.Next() {
		var appointment Appointment
		err := rows.Scan(
			&appointment.ID,
			&appointment.HealthcareEntityID,
			&appointment.PatientID,
			&appointment.DoctorID,
			&appointment.DateTime,
			&appointment.Duration,
			&appointment.Type,
			&appointment.Status,
			&appointment.Reason,
			&appointment.Notes,
			&appointment.Priority,
			&appointment.RoomID,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
			&appointment.CreatedBy,
			&appointment.IsActive,
		)
		if err != nil {
			return nil, err
		}
		
		appointmentResponse := appointment.ToAppointmentResponse()
		appointments = append(appointments, appointmentResponse)
	}

	return appointments, nil
}

// Admin Management Service Methods

// Appointment Duration Settings Management

// GetAppointmentDurationSettings gets duration settings for a healthcare entity
func (s *AppointmentService) GetAppointmentDurationSettings(healthcareEntityID int) ([]AppointmentDurationResponse, error) {
	query := `
		SELECT id, healthcare_entity_id, appointment_type, duration_minutes, is_default, is_active, created_at, updated_at
		FROM appointment_duration_settings
		WHERE healthcare_entity_id = $1 AND is_active = true
		ORDER BY appointment_type
	`

	rows, err := s.db.Query(query, healthcareEntityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []AppointmentDurationResponse
	for rows.Next() {
		var setting AppointmentDurationSetting
		err := rows.Scan(
			&setting.ID,
			&setting.HealthcareEntityID,
			&setting.AppointmentType,
			&setting.DurationMinutes,
			&setting.IsDefault,
			&setting.IsActive,
			&setting.CreatedAt,
			&setting.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting.ToAppointmentDurationResponse())
	}

	return settings, nil
}

// CreateAppointmentDurationSetting creates or updates duration setting
func (s *AppointmentService) CreateAppointmentDurationSetting(healthcareEntityID, userID int, req AppointmentDurationRequest) (*AppointmentDurationResponse, error) {
	setting := &AppointmentDurationSetting{
		HealthcareEntityID: healthcareEntityID,
		AppointmentType:    req.AppointmentType,
		DurationMinutes:    req.DurationMinutes,
		IsDefault:          req.IsDefault,
		IsActive:           true,
		CreatedBy:          userID,
	}

	query := `
		INSERT INTO appointment_duration_settings (
			healthcare_entity_id, appointment_type, duration_minutes, is_default, is_active, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) ON CONFLICT (healthcare_entity_id, appointment_type)
		DO UPDATE SET
			duration_minutes = EXCLUDED.duration_minutes,
			is_default = EXCLUDED.is_default,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := s.db.QueryRow(
		query,
		setting.HealthcareEntityID,
		setting.AppointmentType,
		setting.DurationMinutes,
		setting.IsDefault,
		setting.IsActive,
		now,
		now,
		setting.CreatedBy,
	).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)

	if err != nil {
		return nil, err
	}

	response := setting.ToAppointmentDurationResponse()
	return &response, nil
}

// GetDefaultDurationForType gets the default duration for an appointment type
func (s *AppointmentService) GetDefaultDurationForType(healthcareEntityID int, appointmentType string) (int, error) {
	query := `
		SELECT duration_minutes
		FROM appointment_duration_settings
		WHERE healthcare_entity_id = $1 AND appointment_type = $2 AND is_active = true
	`

	var duration int
	err := s.db.QueryRow(query, healthcareEntityID, appointmentType).Scan(&duration)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return default durations if no custom setting exists
			defaults := map[string]int{
				"consultation": 30,
				"follow-up":    15,
				"procedure":    60,
				"emergency":    20,
			}
			if defaultDuration, exists := defaults[appointmentType]; exists {
				return defaultDuration, nil
			}
			return 30, nil // fallback default
		}
		return 0, err
	}

	return duration, nil
}

// Room Management Service Methods

// GetRooms gets all rooms for a healthcare entity
func (s *AppointmentService) GetRooms(healthcareEntityID int, roomType string, floor int) ([]RoomResponse, error) {
	query := `
		SELECT id, healthcare_entity_id, room_number, room_name, room_type, floor, department, capacity, equipment, is_active, notes, created_at, updated_at
		FROM rooms
		WHERE healthcare_entity_id = $1 AND is_active = true
	`
	var args []interface{}
	args = append(args, healthcareEntityID)
	argCount := 1

	if roomType != "" {
		argCount++
		query += fmt.Sprintf(" AND room_type = $%d", argCount)
		args = append(args, roomType)
	}

	if floor > 0 {
		argCount++
		query += fmt.Sprintf(" AND floor = $%d", argCount)
		args = append(args, floor)
	}

	query += " ORDER BY floor, room_number"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []RoomResponse
	for rows.Next() {
		var room Room
		err := rows.Scan(
			&room.ID,
			&room.HealthcareEntityID,
			&room.RoomNumber,
			&room.RoomName,
			&room.RoomType,
			&room.Floor,
			&room.Department,
			&room.Capacity,
			&room.Equipment,
			&room.IsActive,
			&room.Notes,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		roomResponse := room.ToRoomResponse()
		// Check availability for the room (simplified - in real implementation, check against specific time)
		roomResponse.IsAvailable = true
		rooms = append(rooms, roomResponse)
	}

	return rooms, nil
}

// GetRoomByID gets room by ID and healthcare entity ID
func (s *AppointmentService) GetRoomByID(roomID, healthcareEntityID int) (*Room, error) {
	query := `
		SELECT id, healthcare_entity_id, room_number, room_name, room_type, floor, department, capacity, equipment, is_active, notes, created_at, updated_at
		FROM rooms
		WHERE id = $1 AND healthcare_entity_id = $2
	`
	
	var room Room
	err := s.db.QueryRow(query, roomID, healthcareEntityID).Scan(
		&room.ID,
		&room.HealthcareEntityID,
		&room.RoomNumber,
		&room.RoomName,
		&room.RoomType,
		&room.Floor,
		&room.Department,
		&room.Capacity,
		&room.Equipment,
		&room.IsActive,
		&room.Notes,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &room, nil
}

// CreateRoom creates a new room
func (s *AppointmentService) CreateRoom(healthcareEntityID, userID int, req RoomRequest) (*RoomResponse, error) {
	// Convert equipment array to comma-separated string
	equipmentStr := strings.Join(req.Equipment, ",")

	room := &Room{
		HealthcareEntityID: healthcareEntityID,
		RoomNumber:         req.RoomNumber,
		RoomName:           sql.NullString{String: req.RoomName, Valid: req.RoomName != ""},
		RoomType:           req.RoomType,
		Floor:              req.Floor,
		Department:         sql.NullString{String: req.Department, Valid: req.Department != ""},
		Capacity:           req.Capacity,
		Equipment:          sql.NullString{String: equipmentStr, Valid: equipmentStr != ""},
		IsActive:           true,
		Notes:              sql.NullString{String: req.Notes, Valid: req.Notes != ""},
		CreatedBy:          userID,
	}

	query := `
		INSERT INTO rooms (
			healthcare_entity_id, room_number, room_name, room_type, floor, department, capacity, equipment, is_active, notes, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := s.db.QueryRow(
		query,
		room.HealthcareEntityID,
		room.RoomNumber,
		room.RoomName,
		room.RoomType,
		room.Floor,
		room.Department,
		room.Capacity,
		room.Equipment,
		room.IsActive,
		room.Notes,
		now,
		now,
		room.CreatedBy,
	).Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		return nil, err
	}

	response := room.ToRoomResponse()
	return &response, nil
}

// UpdateRoom updates an existing room
func (s *AppointmentService) UpdateRoom(roomID int, req RoomRequest) (*RoomResponse, error) {
	equipmentStr := strings.Join(req.Equipment, ",")

	query := `
		UPDATE rooms SET
			room_number = $1, room_name = $2, room_type = $3, floor = $4, department = $5,
			capacity = $6, equipment = $7, notes = $8, updated_at = $9
		WHERE id = $10 AND is_active = true
		RETURNING id, healthcare_entity_id, room_number, room_name, room_type, floor, department, capacity, equipment, is_active, notes, created_at, updated_at
	`

	var room Room
	err := s.db.QueryRow(
		query,
		req.RoomNumber,
		req.RoomName,
		req.RoomType,
		req.Floor,
		req.Department,
		req.Capacity,
		equipmentStr,
		req.Notes,
		time.Now(),
		roomID,
	).Scan(
		&room.ID,
		&room.HealthcareEntityID,
		&room.RoomNumber,
		&room.RoomName,
		&room.RoomType,
		&room.Floor,
		&room.Department,
		&room.Capacity,
		&room.Equipment,
		&room.IsActive,
		&room.Notes,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("room not found")
		}
		return nil, err
	}

	response := room.ToRoomResponse()
	return &response, nil
}

// DeleteRoom soft deletes a room
func (s *AppointmentService) DeleteRoom(roomID int) error {
	query := `
		UPDATE rooms
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_active = true
	`

	result, err := s.db.Exec(query, roomID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("room not found")
	}

	return nil
}

// GetAvailableRooms gets all rooms with their availability status for a specific time slot
func (s *AppointmentService) GetAvailableRooms(healthcareEntityID int, dateTime time.Time, duration int, roomType string) ([]RoomResponse, error) {
	endTime := dateTime.Add(time.Duration(duration) * time.Minute)

	query := `
		SELECT 
			r.id, r.healthcare_entity_id, r.room_number, r.room_name, r.room_type, r.floor, r.department, r.capacity, r.equipment, r.is_active, r.notes, r.created_at, r.updated_at,
			CASE 
				WHEN EXISTS (
					SELECT 1 FROM appointments a
					WHERE a.room_id = r.id
					AND a.healthcare_entity_id = r.healthcare_entity_id
					AND a.is_active = true
					AND a.status IN ('scheduled', 'confirmed', 'in-progress')
					AND (
						(a.date_time < $4 AND a.date_time + (a.duration || ' minutes')::interval > $3) OR
						(a.date_time < $3 AND a.date_time + (a.duration || ' minutes')::interval > $4)
					)
				) THEN false
				ELSE true
			END as is_available
		FROM rooms r
		WHERE r.healthcare_entity_id = $1 
		AND r.is_active = true
		AND ($2 = '' OR r.room_type = $2)
		ORDER BY r.floor, r.room_number
	`

	rows, err := s.db.Query(query, healthcareEntityID, roomType, dateTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []RoomResponse
	for rows.Next() {
		var room Room
		var isAvailable bool
		err := rows.Scan(
			&room.ID,
			&room.HealthcareEntityID,
			&room.RoomNumber,
			&room.RoomName,
			&room.RoomType,
			&room.Floor,
			&room.Department,
			&room.Capacity,
			&room.Equipment,
			&room.IsActive,
			&room.Notes,
			&room.CreatedAt,
			&room.UpdatedAt,
			&isAvailable,
		)
		if err != nil {
			return nil, err
		}

		roomResponse := room.ToRoomResponse()
		roomResponse.IsAvailable = isAvailable
		rooms = append(rooms, roomResponse)
	}

	return rooms, nil
}
// GetEntityRoomRequirement fetches room requirement setting from user service
func (s *AppointmentService) GetEntityRoomRequirement(healthcareEntityID int) (bool, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://user-service:8081"
	}
	
	url := fmt.Sprintf("%s/api/internal/entity/%d/room-requirement", userServiceURL, healthcareEntityID)
	
	resp, err := http.Get(url)
	if err != nil {
		// Fallback to false if service is unavailable
		return false, nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		// Fallback to false if error
		return false, nil
	}
	
	var response struct {
		RequireRoomAssignment bool `json:"require_room_assignment"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		// Fallback to false if parse error
		return false, nil
	}
	
	return response.RequireRoomAssignment, nil
}

// GetPatientNameByID gets patient name by ID from the patient service
func (s *AppointmentService) GetPatientNameByID(patientID int) (string, error) {
	patientServiceURL := os.Getenv("PATIENT_SERVICE_URL")
	if patientServiceURL == "" {
		patientServiceURL = "http://patient-service:8082"
	}
	
	url := fmt.Sprintf("%s/api/internal/patients/%d", patientServiceURL, patientID)
	
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to call patient service: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("patient service returned status %d", resp.StatusCode)
	}
	
	var response struct {
		Data struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode patient service response: %v", err)
	}
	
	return fmt.Sprintf("%s %s", response.Data.FirstName, response.Data.LastName), nil
}

// GetDurationOptions returns available duration options for an appointment type
func (s *AppointmentService) GetDurationOptions(healthcareEntityID int, appointmentType string) ([]AppointmentDurationOption, error) {
	query := `
		SELECT id, healthcare_entity_id, appointment_type, duration_minutes, is_default, is_active, display_order, created_at, updated_at, created_by
		FROM appointment_duration_options 
		WHERE healthcare_entity_id = $1 AND is_active = true
	`
	args := []interface{}{healthcareEntityID}

	if appointmentType != "" {
		query += ` AND appointment_type = $2`
		args = append(args, appointmentType)
	}

	query += ` ORDER BY appointment_type, display_order, duration_minutes`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []AppointmentDurationOption
	for rows.Next() {
		var option AppointmentDurationOption
		err := rows.Scan(
			&option.ID,
			&option.HealthcareEntityID,
			&option.AppointmentType,
			&option.DurationMinutes,
			&option.IsDefault,
			&option.IsActive,
			&option.DisplayOrder,
			&option.CreatedAt,
			&option.UpdatedAt,
			&option.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		options = append(options, option)
	}

	return options, nil
}

// CreateDurationOption creates a new duration option
func (s *AppointmentService) CreateDurationOption(healthcareEntityID, userID int, req DurationOptionRequest) (*AppointmentDurationOption, error) {
	option := &AppointmentDurationOption{
		HealthcareEntityID: healthcareEntityID,
		AppointmentType:    req.AppointmentType,
		DurationMinutes:    req.DurationMinutes,
		IsDefault:          req.IsDefault,
		IsActive:           req.IsActive,
		DisplayOrder:       req.DisplayOrder,
		CreatedBy:          userID,
	}

	query := `
		INSERT INTO appointment_duration_options (
			healthcare_entity_id, appointment_type, duration_minutes, is_default, is_active, display_order, created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $7)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(
		query,
		option.HealthcareEntityID,
		option.AppointmentType,
		option.DurationMinutes,
		option.IsDefault,
		option.IsActive,
		option.DisplayOrder,
		option.CreatedBy,
	).Scan(&option.ID, &option.CreatedAt, &option.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return option, nil
}

// UpdateDurationOption updates a duration option
func (s *AppointmentService) UpdateDurationOption(healthcareEntityID, optionID int, req DurationOptionRequest) (*AppointmentDurationOption, error) {
	query := `
		UPDATE appointment_duration_options 
		SET appointment_type = $3, duration_minutes = $4, is_default = $5, is_active = $6, display_order = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND healthcare_entity_id = $2
		RETURNING id, healthcare_entity_id, appointment_type, duration_minutes, is_default, is_active, display_order, created_at, updated_at, created_by
	`

	var option AppointmentDurationOption
	err := s.db.QueryRow(
		query,
		optionID,
		healthcareEntityID,
		req.AppointmentType,
		req.DurationMinutes,
		req.IsDefault,
		req.IsActive,
		req.DisplayOrder,
	).Scan(
		&option.ID,
		&option.HealthcareEntityID,
		&option.AppointmentType,
		&option.DurationMinutes,
		&option.IsDefault,
		&option.IsActive,
		&option.DisplayOrder,
		&option.CreatedAt,
		&option.UpdatedAt,
		&option.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &option, nil
}

// DeleteDurationOption deletes a duration option
func (s *AppointmentService) DeleteDurationOption(healthcareEntityID, optionID int) error {
	query := `DELETE FROM appointment_duration_options WHERE id = $1 AND healthcare_entity_id = $2`
	
	result, err := s.db.Exec(query, optionID, healthcareEntityID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
