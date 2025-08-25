package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	service *AppointmentService
}

func NewAppointmentHandler(service *AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{service: service}
}

// GetAppointments handles GET /api/appointments
func (h *AppointmentHandler) GetAppointments(c *gin.Context) {
	// Get healthcare entity ID from header
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID header is required"})
		return
	}
	
	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")
	patientIDStr := c.Query("patient_id")
	doctorIDStr := c.Query("doctor_id")
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	var patientID, doctorID int
	var dateFrom, dateTo time.Time

	if patientIDStr != "" {
		patientID, _ = strconv.Atoi(patientIDStr)
	}
	if doctorIDStr != "" {
		doctorID, _ = strconv.Atoi(doctorIDStr)
	}
	if dateFromStr != "" {
		dateFrom, _ = time.Parse("2006-01-02", dateFromStr)
	}
	if dateToStr != "" {
		dateTo, _ = time.Parse("2006-01-02", dateToStr)
	}

	search := AppointmentSearch{
		HealthcareEntityID: healthcareEntityID,
		Limit:              limit,
		Offset:             offset,
		Status:             status,
		PatientID:          patientID,
		DoctorID:           doctorID,
		DateFrom:           dateFrom,
		DateTo:             dateTo,
	}

	appointments, err := h.service.GetAppointments(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get appointments",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Convert to response format with enhanced data
	var responses []AppointmentResponse
	for _, appointment := range appointments {
		response := appointment.ToAppointmentResponse()
		
		// Fetch doctor name
		if doctor, err := h.service.GetDoctorByID(appointment.DoctorID, healthcareEntityID); err == nil {
			response.DoctorName = fmt.Sprintf("Dr. %s %s", doctor.FirstName, doctor.LastName)
		}
		
		// Fetch patient name (call patient service)
		if patientName, err := h.service.GetPatientNameByID(appointment.PatientID); err == nil {
			response.PatientName = patientName
		}

		// Fetch room information if room is assigned
		if response.RoomID != nil {
			if room, err := h.service.GetRoomByID(*response.RoomID, healthcareEntityID); err == nil {
				response.RoomNumber = room.RoomNumber
				response.RoomName = room.RoomName.String
			}
		}
		
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"appointments": responses,
			"total_count":  len(responses),
			"limit":        limit,
			"offset":       offset,
		},
		"message":   "Appointments retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// CreateAppointment handles POST /api/appointments
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	userIDStr := c.GetHeader("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)

	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	healthcareEntityID, _ := strconv.Atoi(healthcareEntityIDStr)

	// Parse date time
	dateTime, err := time.Parse("2006-01-02T15:04:05Z", req.DateTime)
	if err != nil {
		// Try alternative format
		dateTime, err = time.Parse("2006-01-02 15:04", req.DateTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Invalid date format",
				"message":   "Expected format: YYYY-MM-DDTHH:MM:SSZ or YYYY-MM-DD HH:MM",
				"timestamp": time.Now().UTC(),
			})
			return
		}
	}

	// Parse room_id if provided
	var roomID sql.NullInt32
	if req.RoomID != "" && req.RoomID != "0" {
		if parsedRoomID, err := strconv.Atoi(req.RoomID); err == nil {
			roomID = sql.NullInt32{Int32: int32(parsedRoomID), Valid: true}
		}
	}

	appointment := &Appointment{
		HealthcareEntityID: healthcareEntityID,
		PatientID:          req.PatientID,
		DoctorID:           req.DoctorID,
		DateTime:           dateTime,
		Duration:           req.Duration,
		Type:               req.Type,
		Reason:             req.Reason,
		Notes:              req.Notes,
		Priority:           req.Priority,
		RoomID:             roomID,
		CreatedBy:          userID,
	}

	err = h.service.CreateAppointment(appointment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Failed to create appointment",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":      appointment.ToAppointmentResponse(),
		"message":   "Appointment created successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetAppointment handles GET /api/appointments/:id
func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	// Get healthcare entity ID from header
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID header is required"})
		return
	}
	
	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid appointment ID",
			"message":   "Appointment ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	appointment, err := h.service.GetAppointmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Appointment not found",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Validate that the appointment belongs to the user's healthcare entity
	if appointment.HealthcareEntityID != healthcareEntityID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":     "Access denied",
			"message":   "Appointment does not belong to your healthcare entity",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Convert to response format and populate doctor name
	response := appointment.ToAppointmentResponse()
	if doctor, err := h.service.GetDoctorByID(appointment.DoctorID, healthcareEntityID); err == nil {
		response.DoctorName = fmt.Sprintf("Dr. %s %s", doctor.FirstName, doctor.LastName)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      response,
		"message":   "Appointment retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// UpdateAppointment handles PUT /api/appointments/:id
func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid appointment ID",
			"message":   "Appointment ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	var req AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get existing appointment first
	appointment, err := h.service.GetAppointmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Appointment not found",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Parse date time
	dateTime, err := time.Parse("2006-01-02T15:04:05Z", req.DateTime)
	if err != nil {
		// Try alternative format
		dateTime, err = time.Parse("2006-01-02 15:04", req.DateTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Invalid date format",
				"message":   "Expected format: YYYY-MM-DDTHH:MM:SSZ or YYYY-MM-DD HH:MM",
				"timestamp": time.Now().UTC(),
			})
			return
		}
	}

	// Update fields
	appointment.PatientID = req.PatientID
	appointment.DoctorID = req.DoctorID
	appointment.DateTime = dateTime
	appointment.Duration = req.Duration
	appointment.Type = req.Type
	appointment.Reason = req.Reason
	appointment.Notes = req.Notes
	appointment.Priority = req.Priority
	// Parse room_id if provided
	var roomID sql.NullInt32
	if req.RoomID != "" && req.RoomID != "0" {
		if parsedRoomID, err := strconv.Atoi(req.RoomID); err == nil {
			roomID = sql.NullInt32{Int32: int32(parsedRoomID), Valid: true}
		}
	}
	appointment.RoomID = roomID

	err = h.service.UpdateAppointment(appointment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Failed to update appointment",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      appointment.ToAppointmentResponse(),
		"message":   "Appointment updated successfully",
		"timestamp": time.Now().UTC(),
	})
}

// DeleteAppointment handles DELETE /api/appointments/:id
func (h *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid appointment ID",
			"message":   "Appointment ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	err = h.service.DeleteAppointment(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Failed to delete appointment",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Appointment deleted successfully",
		"timestamp": time.Now().UTC(),
	})
}

// UpdateAppointmentStatus handles PATCH /api/appointments/:id/status
func (h *AppointmentHandler) UpdateAppointmentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid appointment ID",
			"message":   "Appointment ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	var req AppointmentUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	err = h.service.UpdateAppointmentStatus(id, req.Status, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Failed to update appointment status",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get updated appointment to return
	appointment, err := h.service.GetAppointmentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get updated appointment",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      appointment.ToAppointmentResponse(),
		"message":   "Appointment status updated successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetDoctorSchedules handles GET /api/schedules
func (h *AppointmentHandler) GetDoctorSchedules(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid doctor ID",
			"message":   "Doctor ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid date format",
			"message":   "Date must be in YYYY-MM-DD format",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	schedule, err := h.service.GetDoctorSchedule(doctorID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get doctor schedule",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      schedule,
		"message":   "Doctor schedule retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// Doctor Availability Handlers

// CreateDoctorAvailability handles POST /api/availability
func (h *AppointmentHandler) CreateDoctorAvailability(c *gin.Context) {
	var req DoctorAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get user ID and healthcare entity ID from headers
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	// Parse UTC datetime strings
	startDateTime, err := time.Parse(time.RFC3339, req.StartDateTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid start datetime format",
			"message":   "StartDateTime must be in ISO 8601 UTC format (2006-01-02T15:04:05Z)",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	endDateTime, err := time.Parse(time.RFC3339, req.EndDateTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid end datetime format",
			"message":   "EndDateTime must be in ISO 8601 UTC format (2006-01-02T15:04:05Z)",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	var breakStartDateTime, breakEndDateTime *time.Time
	if req.BreakStartDateTime != "" {
		breakStart, err := time.Parse(time.RFC3339, req.BreakStartDateTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Invalid break start datetime format",
				"message":   "BreakStartDateTime must be in ISO 8601 UTC format (2006-01-02T15:04:05Z)",
				"timestamp": time.Now().UTC(),
			})
			return
		}
		breakStartDateTime = &breakStart
	}

	if req.BreakEndDateTime != "" {
		breakEnd, err := time.Parse(time.RFC3339, req.BreakEndDateTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Invalid break end datetime format",
				"message":   "BreakEndDateTime must be in ISO 8601 UTC format (2006-01-02T15:04:05Z)",
				"timestamp": time.Now().UTC(),
			})
			return
		}
		breakEndDateTime = &breakEnd
	}

	availability := &DoctorAvailability{
		HealthcareEntityID:    healthcareEntityID,
		DoctorID:              req.DoctorID,
		Status:                req.Status,
		StartDateTime:         &startDateTime,
		EndDateTime:           &endDateTime,
		BreakStartDateTime:    breakStartDateTime,
		BreakEndDateTime:      breakEndDateTime,
		Notes:                 req.Notes,
		CreatedBy:             userID.(int),
	}

	err = h.service.CreateDoctorAvailability(availability)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to create availability",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":      availability.ToAvailabilityResponse(""),
		"message":   "Doctor availability created successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetDoctorAvailability handles GET /api/availability
func (h *AppointmentHandler) GetDoctorAvailability(c *gin.Context) {
	// Get healthcare entity ID from header
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID header is required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	// Parse query parameters
	search := AvailabilitySearch{
		HealthcareEntityID: healthcareEntityID,
		DoctorID:           0,
		DateFrom:           c.Query("date_from"),
		DateTo:             c.Query("date_to"),
		Status:             c.Query("status"),
		Limit:              10,
		Offset:             0,
	}

	if doctorIDStr := c.Query("doctor_id"); doctorIDStr != "" {
		if doctorID, err := strconv.Atoi(doctorIDStr); err == nil {
			search.DoctorID = doctorID
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			search.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			search.Offset = offset
		}
	}

	availabilities, err := h.service.GetDoctorAvailability(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get availability",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      availabilities,
		"message":   "Doctor availability retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetDoctorAvailabilityByID handles GET /api/availability/:id
func (h *AppointmentHandler) GetDoctorAvailabilityByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid availability ID",
			"message":   "Availability ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	availability, err := h.service.GetDoctorAvailabilityByID(id)
	if err != nil {
		if err.Error() == "availability record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "Availability not found",
				"message":   "No availability record found with the given ID",
				"timestamp": time.Now().UTC(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get availability",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      availability.ToAvailabilityResponse(""),
		"message":   "Doctor availability retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// UpdateDoctorAvailability handles PUT /api/availability/:id
func (h *AppointmentHandler) UpdateDoctorAvailability(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid availability ID",
			"message":   "Availability ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	var req DoctorAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get existing availability
	availability, err := h.service.GetDoctorAvailabilityByID(id)
	if err != nil {
		if err.Error() == "availability record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "Availability not found",
				"message":   "No availability record found with the given ID",
				"timestamp": time.Now().UTC(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get availability",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Parse UTC datetime strings and update fields
	startDateTime, err := time.Parse(time.RFC3339, req.StartDateTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid start datetime format",
			"message":   "StartDateTime must be in ISO 8601 UTC format (2006-01-02T15:04:05Z)",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	endDateTime, err := time.Parse(time.RFC3339, req.EndDateTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid end datetime format",
			"message":   "EndDateTime must be in ISO 8601 UTC format (2006-01-02T15:04:05Z)",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	var breakStartDateTime, breakEndDateTime *time.Time
	if req.BreakStartDateTime != "" {
		breakStart, err := time.Parse(time.RFC3339, req.BreakStartDateTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Invalid break start datetime format",
				"message":   "BreakStartDateTime must be in ISO 8601 UTC format (2006-01-02T15:04:05Z)",
				"timestamp": time.Now().UTC(),
			})
			return
		}
		breakStartDateTime = &breakStart
	}

	if req.BreakEndDateTime != "" {
		breakEnd, err := time.Parse(time.RFC3339, req.BreakEndDateTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Invalid break end datetime format",
				"message":   "BreakEndDateTime must be in ISO 8601 UTC format (2006-01-02T15:04:05Z)",
				"timestamp": time.Now().UTC(),
			})
			return
		}
		breakEndDateTime = &breakEnd
	}

	// Update fields
	availability.Status = req.Status
	availability.StartDateTime = &startDateTime
	availability.EndDateTime = &endDateTime
	availability.BreakStartDateTime = breakStartDateTime
	availability.BreakEndDateTime = breakEndDateTime
	availability.Notes = req.Notes

	err = h.service.UpdateDoctorAvailability(availability)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to update availability",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      availability.ToAvailabilityResponse(""),
		"message":   "Doctor availability updated successfully",
		"timestamp": time.Now().UTC(),
	})
}

// DeleteDoctorAvailability handles DELETE /api/availability/:id
func (h *AppointmentHandler) DeleteDoctorAvailability(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid availability ID",
			"message":   "Availability ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	err = h.service.DeleteDoctorAvailability(id)
	if err != nil {
		if err.Error() == "availability record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "Availability not found",
				"message":   "No availability record found with the given ID",
				"timestamp": time.Now().UTC(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to delete availability",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Doctor availability deleted successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetAvailabilityCalendar handles GET /api/availability/calendar
func (h *AppointmentHandler) GetAvailabilityCalendar(c *gin.Context) {
	// Get healthcare entity ID from header
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID header is required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	doctorIDStr := c.Query("doctor_id")
	yearMonth := c.DefaultQuery("month", time.Now().Format("2006-01"))

	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid doctor ID",
			"message":   "Doctor ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	calendar, err := h.service.GetAvailabilityCalendar(healthcareEntityID, doctorID, yearMonth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get availability calendar",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      calendar,
		"message":   "Availability calendar retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// CreateBulkAvailability handles POST /api/availability/bulk
func (h *AppointmentHandler) CreateBulkAvailability(c *gin.Context) {
	var req struct {
		DoctorID int                        `json:"doctor_id" validate:"required"`
		DateFrom string                     `json:"date_from" validate:"required"`
		DateTo   string                     `json:"date_to" validate:"required"`
		Template DoctorAvailabilityRequest `json:"template" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get user ID and healthcare entity ID from headers
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	err = h.service.CreateBulkAvailability(
		req.DoctorID,
		healthcareEntityID,
		userID.(int),
		req.DateFrom,
		req.DateTo,
		req.Template,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to create bulk availability",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Bulk availability created successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetDoctorsByEntity handles GET /api/doctors
func (h *AppointmentHandler) GetDoctorsByEntity(c *gin.Context) {
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	doctors, err := h.service.GetDoctorsByEntity(healthcareEntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get doctors",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      doctors,
		"message":   "Doctors retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// BookAppointment handles POST /api/appointments/book - Smart booking with conflict detection
func (h *AppointmentHandler) BookAppointment(c *gin.Context) {
	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get user ID and healthcare entity ID from headers
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	// Check if room assignment is required for this healthcare entity
	requireRoom, err := h.service.GetEntityRoomRequirement(healthcareEntityID)
	if err != nil {
		// Log error but don't fail the request - room requirement check is optional
		log.Printf("Failed to check room requirement for entity %d: %v", healthcareEntityID, err)
	}

	// Validate room selection if required
	if requireRoom && req.RoomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Room assignment required",
			"message":   "This healthcare entity requires room assignment for all appointments. Please select a room.",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Book the appointment with conflict checking
	response, err := h.service.BookAppointmentWithConflictCheck(req, healthcareEntityID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to process booking request",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Return appropriate status based on success
	statusCode := http.StatusOK
	if response.Success {
		statusCode = http.StatusCreated
	}

	c.JSON(statusCode, gin.H{
		"data":      response,
		"message":   response.Message,
		"timestamp": time.Now().UTC(),
	})
}

// GetTimeSlots handles GET /api/appointments/slots - Get available time slots for booking
func (h *AppointmentHandler) GetTimeSlots(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	date := c.Query("date")
	durationStr := c.DefaultQuery("duration", "30")

	// Validate parameters
	if doctorIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Missing required parameter",
			"message":   "doctor_id is required",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Missing required parameter",
			"message":   "date is required (YYYY-MM-DD format)",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid doctor ID",
			"message":   "Doctor ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration < 15 || duration > 480 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid duration",
			"message":   "Duration must be between 15 and 480 minutes",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Validate date format
	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid date format",
			"message":   "Date must be in YYYY-MM-DD format",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get healthcare entity ID from headers
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	// Get available time slots
	slots, err := h.service.GetAvailableTimeSlots(doctorID, healthcareEntityID, date, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get available time slots",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Organize slots by type for better UI presentation
	slotsByType := make(map[string][]AvailabilitySlot)
	availableCount := 0
	
	for _, slot := range slots {
		slotsByType[slot.SlotType] = append(slotsByType[slot.SlotType], slot)
		if slot.IsAvailable {
			availableCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"date":            date,
			"doctor_id":       doctorID,
			"duration":        duration,
			"total_slots":     len(slots),
			"available_slots": availableCount,
			"slots":           slots,
			"slots_by_type":   slotsByType,
		},
		"message":   "Available time slots retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// Admin Management Handlers

// Appointment Duration Settings Handlers

// GetAppointmentDurationSettings handles GET /api/admin/duration-settings
func (h *AppointmentHandler) GetAppointmentDurationSettings(c *gin.Context) {
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	settings, err := h.service.GetAppointmentDurationSettings(healthcareEntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get duration settings",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      settings,
		"message":   "Duration settings retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// CreateAppointmentDurationSetting handles POST /api/admin/duration-settings
func (h *AppointmentHandler) CreateAppointmentDurationSetting(c *gin.Context) {
	var req AppointmentDurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	setting, err := h.service.CreateAppointmentDurationSetting(healthcareEntityID, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to create duration setting",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":      setting,
		"message":   "Duration setting created successfully",
		"timestamp": time.Now().UTC(),
	})
}

// Room Management Handlers

// GetRooms handles GET /api/admin/rooms
func (h *AppointmentHandler) GetRooms(c *gin.Context) {
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	roomType := c.Query("room_type")
	floorStr := c.Query("floor")
	floor := 0
	if floorStr != "" {
		floor, _ = strconv.Atoi(floorStr)
	}

	rooms, err := h.service.GetRooms(healthcareEntityID, roomType, floor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get rooms",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      rooms,
		"message":   "Rooms retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// CreateRoom handles POST /api/admin/rooms
func (h *AppointmentHandler) CreateRoom(c *gin.Context) {
	var req RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	room, err := h.service.CreateRoom(healthcareEntityID, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to create room",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":      room,
		"message":   "Room created successfully",
		"timestamp": time.Now().UTC(),
	})
}

// UpdateRoom handles PUT /api/admin/rooms/:id
func (h *AppointmentHandler) UpdateRoom(c *gin.Context) {
	idStr := c.Param("id")
	roomID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid room ID",
			"message":   "Room ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	var req RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request format",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	room, err := h.service.UpdateRoom(roomID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to update room",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      room,
		"message":   "Room updated successfully",
		"timestamp": time.Now().UTC(),
	})
}

// DeleteRoom handles DELETE /api/admin/rooms/:id
func (h *AppointmentHandler) DeleteRoom(c *gin.Context) {
	idStr := c.Param("id")
	roomID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid room ID",
			"message":   "Room ID must be a number",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	err = h.service.DeleteRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to delete room",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Room deleted successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetAvailableRooms handles GET /api/appointments/available-rooms
func (h *AppointmentHandler) GetAvailableRooms(c *gin.Context) {
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required"})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
		return
	}

	dateTimeStr := c.Query("date_time")
	durationStr := c.DefaultQuery("duration", "30")
	roomType := c.Query("room_type")

	if dateTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Missing required parameter",
			"message":   "date_time is required (YYYY-MM-DDTHH:MM format)",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration < 15 || duration > 480 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid duration",
			"message":   "Duration must be between 15 and 480 minutes",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Parse UTC datetime (same as appointment booking)
	dateTime, err := time.Parse("2006-01-02T15:04:05Z", dateTimeStr)
	if err != nil {
		// Try alternative format
		dateTime, err = time.Parse("2006-01-02T15:04", dateTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Invalid date format",
				"message":   "Date must be in YYYY-MM-DDTHH:MM:SSZ or YYYY-MM-DDTHH:MM format",
				"timestamp": time.Now().UTC(),
			})
			return
		}
	}

	rooms, err := h.service.GetAvailableRooms(healthcareEntityID, dateTime, duration, roomType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get available rooms",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"date_time":       dateTimeStr,
			"duration":        duration,
			"room_type":       roomType,
			"total_rooms":     len(rooms),
			"available_rooms": rooms,
		},
		"message":   "Available rooms retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}
// GetDurationOptions handles GET /api/admin/duration-options
func (h *AppointmentHandler) GetDurationOptions(c *gin.Context) {
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Missing healthcare entity ID",
			"message":   "Healthcare entity ID is required",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a number",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	appointmentType := c.Query("appointment_type")

	options, err := h.service.GetDurationOptions(healthcareEntityID, appointmentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get duration options",
			"message":   err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": options})
}

// CreateDurationOption handles POST /api/admin/duration-options
func (h *AppointmentHandler) CreateDurationOption(c *gin.Context) {
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Missing healthcare entity ID",
			"message":   "Healthcare entity ID is required",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a number",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Missing user ID",
			"message":   "User ID is required",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid user ID",
			"message":   "User ID must be a number",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	var req DurationOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request data",
			"message":   err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	option, err := h.service.CreateDurationOption(healthcareEntityID, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to create duration option",
			"message":   err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": option})
}

// UpdateDurationOption handles PUT /api/admin/duration-options/:id
func (h *AppointmentHandler) UpdateDurationOption(c *gin.Context) {
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Missing healthcare entity ID",
			"message":   "Healthcare entity ID is required",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a number",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}
	
	optionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid duration option ID",
			"message":   "Duration option ID must be a number",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	var req DurationOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request data",
			"message":   err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	option, err := h.service.UpdateDurationOption(healthcareEntityID, optionID, req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "Duration option not found",
				"message":   "Duration option not found",
				"timestamp": time.Now().Format(time.RFC3339),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to update duration option",
			"message":   err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": option})
}

// DeleteDurationOption handles DELETE /api/admin/duration-options/:id
func (h *AppointmentHandler) DeleteDurationOption(c *gin.Context) {
	healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if healthcareEntityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Missing healthcare entity ID",
			"message":   "Healthcare entity ID is required",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a number",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}
	
	optionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid duration option ID",
			"message":   "Duration option ID must be a number",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	err = h.service.DeleteDurationOption(healthcareEntityID, optionID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error":     "Duration option not found",
				"message":   "Duration option not found",
				"timestamp": time.Now().Format(time.RFC3339),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to delete duration option",
			"message":   err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Duration option deleted successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
