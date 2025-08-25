package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	logging "github.com/louhibi/healthcare-logging"
)

// PatientHandler groups all patient related handlers
type PatientHandler struct {
    patientService *PatientService
    validator      *validator.Validate
}

// NewPatientHandler constructs a new PatientHandler
func NewPatientHandler(patientService *PatientService) *PatientHandler {
    return &PatientHandler{patientService: patientService, validator: validator.New()}
}

// parseAndValidateRequest is a generic helper to parse and validate JSON requests with logging
func (h *PatientHandler) parseAndValidateRequest(c *gin.Context, req interface{}, operation string) error {
    if err := c.ShouldBindJSON(req); err != nil {
        logging.LogWarn("Invalid request format",
            "operation", operation,
            "error", err.Error(),
            "method", c.Request.Method,
            "endpoint", c.Request.URL.Path,
            "user_id", c.GetHeader("X-User-ID"),
            "healthcare_entity_id", c.GetHeader("X-Healthcare-Entity-ID"))

        c.JSON(http.StatusBadRequest, gin.H{
            "error":        "Invalid request format",
            "details":      err.Error(),
            "received_data": req,
        })
        return err
    }

    if err := h.validator.Struct(req); err != nil {
        logging.LogWarn("Validation failed",
            "operation", operation,
            "error", err.Error(),
            "method", c.Request.Method,
            "endpoint", c.Request.URL.Path,
            "user_id", c.GetHeader("X-User-ID"),
            "healthcare_entity_id", c.GetHeader("X-Healthcare-Entity-ID"),
            "request_data", req)

        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Validation failed",
            "details": err.Error(),
        })
        return err
    }

    return nil
}

// CreatePatient handles patient creation
func (h *PatientHandler) CreatePatient(c *gin.Context) {
    var req PatientRequest
    if err := h.parseAndValidateRequest(c, &req, "CreatePatient"); err != nil {
        return
    }

    healthcareEntityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
    if healthcareEntityIDStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Healthcare entity ID required", "details": "X-Healthcare-Entity-ID header is missing"})
        return
    }

    healthcareEntityID, err := strconv.Atoi(healthcareEntityIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID", "details": "X-Healthcare-Entity-ID must be a valid integer"})
        return
    }

    authHeaders := map[string]string{
        "Authorization":  c.GetHeader("Authorization"),
        "X-User-ID":      c.GetHeader("X-User-ID"),
        "X-User-Email":   c.GetHeader("X-User-Email"),
        "X-User-Role":    c.GetHeader("X-User-Role"),
    }

    validationErrors, err := h.patientService.ValidatePatientRequest(&req, healthcareEntityID, authHeaders)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation system error", "details": err.Error()})
        return
    }
    if len(validationErrors) > 0 {
        var msgs []string
        for _, ve := range validationErrors { msgs = append(msgs, ve.Message) }
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "validation_errors": msgs, "received_data": req})
        return
    }

    userID, exists := c.Get("user_id")
    if !exists { c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"}); return }

    req.Email = strings.TrimSpace(req.Email)
    if req.CountryID <= 0 { c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country", "details": "Country ID must be a positive integer"}); return }

    dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format", "details": fmt.Sprintf("Date of birth must be in YYYY-MM-DD format, received: '%s'", req.DateOfBirth), "parse_error": err.Error()}); return }

    if strings.TrimSpace(req.Email) != "" {
        emailExists, err := h.patientService.EmailExists(req.Email, healthcareEntityID, 0)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": "Failed to check email existence", "internal_error": err.Error()}); return }
        if emailExists { c.JSON(http.StatusConflict, gin.H{"error": "Email already exists", "details": fmt.Sprintf("A patient with email '%s' already exists", req.Email)}); return }
    }

    patient := &Patient{
        HealthcareEntityID: healthcareEntityID,
        PatientID:          req.PatientID,
        FirstName:          req.FirstName,
        LastName:           req.LastName,
        DateOfBirth:        dateOfBirth,
        Gender:             req.Gender,
        Phone:              req.Phone,
        Email:              req.Email,
        Address:            req.Address,
        CountryID:          req.CountryID,
        StateID:            req.StateID,
        CityID:             req.CityID,
        PostalCode:         req.PostalCode,
        NationalityID:      req.NationalityID,
        PreferredLanguage:  req.PreferredLanguage,
        MaritalStatus:      req.MaritalStatus,
        Occupation:         req.Occupation,
        InsuranceTypeID:    req.InsuranceTypeID,
        PolicyNumber:       req.PolicyNumber,
        InsuranceProviderID: req.InsuranceProviderID,
        NationalID:         req.NationalID,
        EmergencyContactName:         req.EmergencyContactName,
        EmergencyContactPhone:        req.EmergencyContactPhone,
        EmergencyContactRelationship: req.EmergencyContactRelationship,
        MedicalHistory: req.MedicalHistory,
        Allergies:      req.Allergies,
        Medications:    req.Medications,
        BloodType:      req.BloodType,
        CreatedBy:      userID.(int),
    }

    if err := h.patientService.CreatePatient(patient); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient", "details": "Database operation failed", "internal_error": err.Error(), "patient_data": patient})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Patient created successfully", "patient": patient.ToPatientResponse()})
}

// GetPatient returns a patient by ID
func (h *PatientHandler) GetPatient(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"}); return }
    patient, err := h.patientService.GetPatientByID(id)
    if err != nil { if err.Error()=="patient not found" { c.JSON(http.StatusNotFound, gin.H{"error":"Patient not found"}); return }; c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to get patient"}); return }
    c.JSON(http.StatusOK, patient.ToPatientResponse())
}

// UpdatePatient updates a patient record
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"}); return }
    var req PatientRequest
    if err := h.parseAndValidateRequest(c, &req, "UpdatePatient"); err != nil { return }
    existing, err := h.patientService.GetPatientByID(id)
    if err != nil { if err.Error()=="patient not found" { c.JSON(http.StatusNotFound, gin.H{"error":"Patient not found"}); return }; c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to get patient"}); return }
    dob, err := time.Parse("2006-01-02", req.DateOfBirth)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid date format", "details": fmt.Sprintf("Date of birth must be in YYYY-MM-DD format, received: '%s'", req.DateOfBirth), "parse_error": err.Error()}); return }
    req.Email = strings.TrimSpace(req.Email)
    if req.Email != "" { emailExists, err := h.patientService.EmailExists(req.Email, existing.HealthcareEntityID, id); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"Internal server error"}); return }; if emailExists { c.JSON(http.StatusConflict, gin.H{"error":"Email already exists"}); return } }
    if req.CountryID <= 0 { c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid country", "details":"Country ID must be a positive integer"}); return }
    existing.FirstName = req.FirstName; existing.LastName=req.LastName; existing.DateOfBirth=dob; existing.Gender=req.Gender; existing.Phone=req.Phone; existing.Email=req.Email; existing.Address=req.Address; existing.CountryID=req.CountryID; existing.StateID=req.StateID; existing.CityID=req.CityID; existing.PostalCode=req.PostalCode; existing.NationalityID=req.NationalityID; existing.PreferredLanguage=req.PreferredLanguage; existing.MaritalStatus=req.MaritalStatus; existing.Occupation=req.Occupation; existing.InsuranceTypeID=req.InsuranceTypeID; existing.PolicyNumber=req.PolicyNumber; existing.InsuranceProviderID=req.InsuranceProviderID; existing.NationalID=req.NationalID; existing.EmergencyContactName=req.EmergencyContactName; existing.EmergencyContactPhone=req.EmergencyContactPhone; existing.EmergencyContactRelationship=req.EmergencyContactRelationship; existing.MedicalHistory=req.MedicalHistory; existing.Allergies=req.Allergies; existing.Medications=req.Medications; existing.BloodType=req.BloodType
    if err := h.patientService.UpdatePatient(existing); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to update patient"}); return }
    c.JSON(http.StatusOK, existing.ToPatientResponse())
}

// DeletePatient deletes a patient
func (h *PatientHandler) DeletePatient(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid patient ID"}); return }
    if err := h.patientService.DeletePatient(id); err != nil { if err.Error()=="patient not found" { c.JSON(http.StatusNotFound, gin.H{"error":"Patient not found"}); return }; c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to delete patient"}); return }
    c.JSON(http.StatusOK, gin.H{"message":"Patient deleted successfully"})
}

// GetPatients returns paginated patients
func (h *PatientHandler) GetPatients(c *gin.Context) {
    var search PatientSearchRequest
    if err := c.ShouldBindQuery(&search); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid query parameters"}); return }
    heIDStr := c.GetHeader("X-Healthcare-Entity-ID")
    if heIDStr == "" { c.JSON(http.StatusBadRequest, gin.H{"error":"Healthcare entity ID header is required"}); return }
    heID, err := strconv.Atoi(heIDStr); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid healthcare entity ID"}); return }
    search.HealthcareEntityID = heID
    if search.Limit <=0 || search.Limit>100 { search.Limit = 10 }
    if search.Offset <0 { search.Offset = 0 }
    patients, err := h.patientService.GetPatients(search); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to get patients"}); return }
    var resp []PatientResponse
    for _, p := range patients { resp = append(resp, p.ToPatientResponse()) }
    total, err := h.patientService.GetPatientCount(heID); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to get patient count"}); return }
    c.JSON(http.StatusOK, gin.H{"patients":resp, "total_count": total, "limit": search.Limit, "offset": search.Offset, "has_more": search.Offset+search.Limit < total})
}

// GetPatientStats returns patient statistics
func (h *PatientHandler) GetPatientStats(c *gin.Context) {
    heIDStr := c.GetHeader("X-Healthcare-Entity-ID")
    if heIDStr == "" { c.JSON(http.StatusBadRequest, gin.H{"error":"Healthcare entity ID header is required"}); return }
    heID, err := strconv.Atoi(heIDStr); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid healthcare entity ID"}); return }
    total, err := h.patientService.GetPatientCount(heID); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to get patient stats"}); return }
    c.JSON(http.StatusOK, gin.H{"total_patients": total})
}
