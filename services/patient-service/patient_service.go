package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type PatientService struct {
	db *sql.DB
}

func NewPatientService(db *sql.DB) *PatientService {
	return &PatientService{db: db}
}

// Form validation structures
type FormField struct {
	FieldID      int                    `json:"field_id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	FieldType    string                 `json:"field_type"`
	IsEnabled    bool                   `json:"is_enabled"`
	IsRequired   bool                   `json:"is_required"`
	IsCore       bool                   `json:"is_core"`
	Options      []string               `json:"options"`
	ValidationRules map[string]interface{} `json:"validation_rules"`
}

type FormMetadata struct {
	FormType    string      `json:"form_type"`
	DisplayName string      `json:"display_name"`
	Description string      `json:"description"`
	Fields      []FormField `json:"fields"`
	EntityID    int         `json:"entity_id"`
}

type FormMetadataResponse struct {
	Data      FormMetadata `json:"data"`
	Message   string       `json:"message"`
	Timestamp string       `json:"timestamp"`
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidatePatientRequest validates a patient request against form configuration
func (s *PatientService) ValidatePatientRequest(req *PatientRequest, entityID int, authHeaders map[string]string) ([]ValidationError, error) {
	// Get form configuration from user-service
	formConfig, err := s.getFormConfiguration("patient", entityID, authHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to get form configuration: %w", err)
	}

	var errors []ValidationError
	
	// Create a map for easy field lookup
	fieldMap := make(map[string]FormField)
	for _, field := range formConfig.Fields {
		if field.IsEnabled {
			fieldMap[field.Name] = field
		}
	}

	// Validate each field based on configuration
	patientData := map[string]interface{}{
		"first_name":                    req.FirstName,
		"last_name":                     req.LastName,
		"date_of_birth":                 req.DateOfBirth,
		"gender":                        req.Gender,
		"phone":                         req.Phone,
		"email":                         req.Email,
		"address":                       req.Address,
		"country_id":                    req.CountryID,
		"state_id":                      req.StateID,
		"city_id":                       req.CityID,
		"postal_code":                   req.PostalCode,
		"nationality_id":                req.NationalityID,
		"preferred_language":            req.PreferredLanguage,
		"marital_status":                req.MaritalStatus,
		"occupation":                    req.Occupation,
		"insurance_type_id":             req.InsuranceTypeID,
		"policy_number":                 req.PolicyNumber,
		"insurance_provider_id":         req.InsuranceProviderID,
		"national_id":                   req.NationalID,
		"emergency_contact_name":        req.EmergencyContactName,
		"emergency_contact_phone":       req.EmergencyContactPhone,
		"emergency_contact_relationship": req.EmergencyContactRelationship,
		"medical_history":               req.MedicalHistory,
		"allergies":                     req.Allergies,
		"medications":                   req.Medications,
		"blood_type":                    req.BloodType,
	}

	// Validate each configured field
	for fieldName, fieldConfig := range fieldMap {
		value, exists := patientData[fieldName]
		if !exists {
			continue
		}

		strValue := fmt.Sprintf("%v", value)
		
		// Check if required field is empty
		if fieldConfig.IsRequired && strings.TrimSpace(strValue) == "" {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("%s is required", fieldConfig.DisplayName),
			})
			continue
		}

		// Skip validation for empty optional fields
		if strings.TrimSpace(strValue) == "" {
			continue
		}

		// Validate field type and options
		if err := s.validateFieldValue(fieldName, strValue, fieldConfig); err != nil {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: err.Error(),
			})
		}
	}

	return errors, nil
}

// validateFieldValue validates a specific field value against its configuration
func (s *PatientService) validateFieldValue(fieldName, value string, field FormField) error {
	// Email validation
	if field.FieldType == "email" && value != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(value) {
			return fmt.Errorf("%s must be a valid email address", field.DisplayName)
		}
	}

	// Phone validation
	if field.FieldType == "phone" && value != "" {
		// Basic phone validation - at least 10 digits
		phoneRegex := regexp.MustCompile(`\d`)
		digits := phoneRegex.FindAllString(value, -1)
		if len(digits) < 10 {
			return fmt.Errorf("%s must contain at least 10 digits", field.DisplayName)
		}
	}

	// Select field validation (check if value is in options)
	if field.FieldType == "select" && len(field.Options) > 0 {
		for _, option := range field.Options {
			if option == value {
				return nil // Valid option found
			}
		}
		return fmt.Errorf("%s must be one of: %s", field.DisplayName, strings.Join(field.Options, ", "))
	}

	return nil
}

// getFormConfiguration fetches form configuration from user-service
func (s *PatientService) getFormConfiguration(formType string, entityID int, authHeaders map[string]string) (*FormMetadata, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://user-service:8081"
	}

	url := fmt.Sprintf("%s/api/forms/%s/metadata", userServiceURL, formType)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add entity ID header
	req.Header.Set("X-Healthcare-Entity-ID", fmt.Sprintf("%d", entityID))
	req.Header.Set("Content-Type", "application/json")
	
	// Forward authentication headers
	for key, value := range authHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch form configuration: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	var response FormMetadataResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

// nullableString returns nil for blank strings (after trim), otherwise the trimmed string
func nullableString(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return strings.TrimSpace(s)
}

// CreatePatient creates a new patient
func (s *PatientService) CreatePatient(patient *Patient) error {
	query := `
		INSERT INTO patients (
			healthcare_entity_id, patient_id, first_name, last_name, date_of_birth, gender, phone, email,
			address, country_id, state_id, city_id, postal_code, nationality_id, preferred_language, marital_status,
			occupation, insurance_type_id, policy_number, insurance_provider_id, national_id,
			emergency_contact_name, emergency_contact_phone, emergency_contact_relationship,
			medical_history, allergies, medications, blood_type, is_active, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32
		) RETURNING id, created_at, updated_at
	`
	
	now := time.Now()
	patient.IsActive = true
	patient.CreatedAt = now
	patient.UpdatedAt = now

	err := s.db.QueryRow(
		query,
		patient.HealthcareEntityID,
		patient.PatientID,
		patient.FirstName,
		patient.LastName,
		patient.DateOfBirth,
		patient.Gender,
		patient.Phone,
		nullableString(patient.Email),
		patient.Address,
		patient.CountryID,
		patient.StateID,
		patient.CityID,
		patient.PostalCode,
		patient.NationalityID,
		patient.PreferredLanguage,
		patient.MaritalStatus,
		patient.Occupation,
		patient.InsuranceTypeID,
		patient.PolicyNumber,
		patient.InsuranceProviderID,
		patient.NationalID,
		patient.EmergencyContactName,
		patient.EmergencyContactPhone,
		patient.EmergencyContactRelationship,
		patient.MedicalHistory,
		patient.Allergies,
		patient.Medications,
		patient.BloodType,
		patient.IsActive,
		patient.CreatedAt,
		patient.UpdatedAt,
		patient.CreatedBy,
	).Scan(&patient.ID, &patient.CreatedAt, &patient.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetPatientByID gets patient by ID
func (s *PatientService) GetPatientByID(id int) (*Patient, error) {
	patient := &Patient{}
	query := `
		SELECT 
			id, healthcare_entity_id, patient_id, first_name, last_name, date_of_birth, gender, phone, COALESCE(email, '') AS email,
			address, country_id, state_id, city_id, postal_code, nationality_id, preferred_language, marital_status,
			occupation, insurance_type_id, policy_number, insurance_provider_id, national_id,
			emergency_contact_name, emergency_contact_phone, emergency_contact_relationship,
			medical_history, allergies, medications, blood_type, is_active, created_at, updated_at, created_by
		FROM patients
		WHERE id = $1 AND is_active = true
	`

	err := s.db.QueryRow(query, id).Scan(
		&patient.ID,
		&patient.HealthcareEntityID,
		&patient.PatientID,
		&patient.FirstName,
		&patient.LastName,
		&patient.DateOfBirth,
		&patient.Gender,
		&patient.Phone,
		&patient.Email,
		&patient.Address,
		&patient.CountryID,
		&patient.StateID,
		&patient.CityID,
		&patient.PostalCode,
		&patient.NationalityID,
		&patient.PreferredLanguage,
		&patient.MaritalStatus,
		&patient.Occupation,
		&patient.InsuranceTypeID,
		&patient.PolicyNumber,
		&patient.InsuranceProviderID,
		&patient.NationalID,
		&patient.EmergencyContactName,
		&patient.EmergencyContactPhone,
		&patient.EmergencyContactRelationship,
		&patient.MedicalHistory,
		&patient.Allergies,
		&patient.Medications,
		&patient.BloodType,
		&patient.IsActive,
		&patient.CreatedAt,
		&patient.UpdatedAt,
		&patient.CreatedBy,
	)

	// Debug logging
	log.Printf("DEBUG: Patient %d loaded - CountryID: %d, StateID: %v, CityID: %v", 
		patient.ID, patient.CountryID, patient.StateID, patient.CityID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("patient not found")
		}
		return nil, err
	}

	return patient, nil
}

// UpdatePatient updates patient information
func (s *PatientService) UpdatePatient(patient *Patient) error {
	query := `
		UPDATE patients SET
			first_name = $1, last_name = $2, date_of_birth = $3, gender = $4,
			phone = $5, email = $6, address = $7, country_id = $8, state_id = $9, city_id = $10, postal_code = $11,
			nationality_id = $12, preferred_language = $13, marital_status = $14,
			occupation = $15, insurance = $16, policy_number = $17, insurance_provider = $18,
			national_id = $19, emergency_contact_name = $20, emergency_contact_phone = $21,
			emergency_contact_relationship = $22, medical_history = $23, allergies = $24,
			medications = $25, blood_type = $26, updated_at = CURRENT_TIMESTAMP
		WHERE id = $27 AND is_active = true
		RETURNING updated_at
	`

	err := s.db.QueryRow(
		query,
		patient.FirstName,
		patient.LastName,
		patient.DateOfBirth,
		patient.Gender,
		patient.Phone,
		nullableString(patient.Email),
		patient.Address,
		patient.CountryID,
		patient.StateID,
		patient.CityID,
		patient.PostalCode,
		patient.NationalityID,
		patient.PreferredLanguage,
		patient.MaritalStatus,
		patient.Occupation,
		patient.InsuranceTypeID,
		patient.PolicyNumber,
		patient.InsuranceProviderID,
		patient.NationalID,
		patient.EmergencyContactName,
		patient.EmergencyContactPhone,
		patient.EmergencyContactRelationship,
		patient.MedicalHistory,
		patient.Allergies,
		patient.Medications,
		patient.BloodType,
		patient.ID,
	).Scan(&patient.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("patient not found")
		}
		return err
	}

	return nil
}

// GetPatients gets patients with pagination and filtering
func (s *PatientService) GetPatients(searchReq PatientSearchRequest) ([]Patient, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Base query with healthcare entity filtering for multi-tenant isolation
	baseQuery := `
		SELECT 
			id, healthcare_entity_id, patient_id, first_name, last_name, date_of_birth, gender, phone, COALESCE(email, '') AS email,
			address, country_id, state_id, city_id, postal_code, nationality_id, preferred_language, marital_status,
			occupation, insurance_type_id, policy_number, insurance_provider_id, national_id,
			emergency_contact_name, emergency_contact_phone, emergency_contact_relationship,
			medical_history, allergies, medications, blood_type, is_active, created_at, updated_at, created_by
		FROM patients
		WHERE is_active = true AND healthcare_entity_id = $1
	`
	
	// Add healthcare entity ID as the first parameter
	args = append(args, searchReq.HealthcareEntityID)
	argIndex++

	// Add comprehensive search conditions
	if searchReq.Query != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(searchReq.Query))
		conditions = append(conditions, fmt.Sprintf(`
			(LOWER(first_name) LIKE $%d OR 
			 LOWER(last_name) LIKE $%d OR 
			 LOWER(email) LIKE $%d OR 
			 LOWER(patient_id) LIKE $%d OR
			 phone LIKE $%d OR
			 LOWER(address) LIKE $%d OR
			 postal_code LIKE $%d OR
			 LOWER(occupation) LIKE $%d OR
			 LOWER(insurance) LIKE $%d OR
			 LOWER(insurance_provider) LIKE $%d OR
			 national_id LIKE $%d OR
			 LOWER(medical_history) LIKE $%d OR
			 LOWER(allergies) LIKE $%d OR
			 LOWER(medications) LIKE $%d)
		`, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex))
		args = append(args, searchQuery)
		argIndex++
	}

	// Build final query
	finalQuery := baseQuery
	if len(conditions) > 0 {
		finalQuery += " AND " + strings.Join(conditions, " AND ")
	}

	finalQuery += " ORDER BY created_at DESC"

	// Add pagination
	if searchReq.Limit <= 0 || searchReq.Limit > 100 {
		searchReq.Limit = 10
	}
	if searchReq.Offset < 0 {
		searchReq.Offset = 0
	}

	finalQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, searchReq.Limit, searchReq.Offset)

	// Execute query
	rows, err := s.db.Query(finalQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []Patient
	for rows.Next() {
		var patient Patient
		err := rows.Scan(
			&patient.ID,
			&patient.HealthcareEntityID,
			&patient.PatientID,
			&patient.FirstName,
			&patient.LastName,
			&patient.DateOfBirth,
			&patient.Gender,
			&patient.Phone,
			&patient.Email,
			&patient.Address,
			&patient.CountryID,
			&patient.StateID,
			&patient.CityID,
			&patient.PostalCode,
			&patient.NationalityID,
			&patient.PreferredLanguage,
			&patient.MaritalStatus,
			&patient.Occupation,
			&patient.InsuranceTypeID,
			&patient.PolicyNumber,
			&patient.InsuranceProviderID,
			&patient.NationalID,
			&patient.EmergencyContactName,
			&patient.EmergencyContactPhone,
			&patient.EmergencyContactRelationship,
			&patient.MedicalHistory,
			&patient.Allergies,
			&patient.Medications,
			&patient.BloodType,
			&patient.IsActive,
			&patient.CreatedAt,
			&patient.UpdatedAt,
			&patient.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}

	return patients, nil
}

// DeletePatient soft deletes a patient
func (s *PatientService) DeletePatient(id int) error {
	query := `
		UPDATE patients
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
		return errors.New("patient not found")
	}

	return nil
}

// EmailExists checks if email already exists for another patient
func (s *PatientService) EmailExists(email string, healthcareEntityID int, excludeID int) (bool, error) {
	e := strings.TrimSpace(email)
	if e == "" {
		return false, nil
	}
	var count int
	query := `SELECT COUNT(*) FROM patients WHERE is_active = true AND healthcare_entity_id = $3 AND id != $2 AND email IS NOT NULL AND LOWER(email) = LOWER($1)`

	err := s.db.QueryRow(query, e, excludeID, healthcareEntityID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetPatientCount gets total count of active patients for a healthcare entity
func (s *PatientService) GetPatientCount(healthcareEntityID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM patients WHERE is_active = true AND healthcare_entity_id = $1`
	
	err := s.db.QueryRow(query, healthcareEntityID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}