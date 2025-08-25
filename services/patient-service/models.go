package main

import (
	"time"
)

// Patient represents a patient in the system
type Patient struct {
	ID                   int       `json:"id" db:"id"`
	HealthcareEntityID   int       `json:"healthcare_entity_id" db:"healthcare_entity_id" validate:"required"`
	PatientID            string    `json:"patient_id" db:"patient_id"` // Custom patient identifier
	FirstName            string    `json:"first_name" db:"first_name" validate:"required"`
	LastName             string    `json:"last_name" db:"last_name" validate:"required"`
	DateOfBirth          time.Time `json:"date_of_birth" db:"date_of_birth" validate:"required"`
	Gender               string    `json:"gender" db:"gender"`
	Phone                string    `json:"phone" db:"phone"`
	Email                string    `json:"email" db:"email" validate:"omitempty,email"`
	Address              string    `json:"address" db:"address"`
	CountryID            int       `json:"country_id" db:"country_id" validate:"required"`    // Foreign key to location service
	StateID              *int      `json:"state_id,omitempty" db:"state_id"`                  // Foreign key to location service (nullable)
	CityID               *int      `json:"city_id,omitempty" db:"city_id"`                    // Foreign key to location service (nullable)
	PostalCode           string    `json:"postal_code" db:"postal_code"`
	NationalityID        *int      `json:"nationality_id,omitempty" db:"nationality_id"` // Foreign key to location service
	PreferredLanguage    string    `json:"preferred_language" db:"preferred_language"`
	MaritalStatus        string    `json:"marital_status" db:"marital_status"`
	Occupation           string    `json:"occupation" db:"occupation"`
	InsuranceTypeID      *int      `json:"insurance_type_id,omitempty" db:"insurance_type_id"`     // Foreign key to location service
	InsuranceProviderID  *int      `json:"insurance_provider_id,omitempty" db:"insurance_provider_id"` // Foreign key to location service
	PolicyNumber         string    `json:"policy_number" db:"policy_number"`
	NationalID           string    `json:"national_id" db:"national_id"` // SSN, SIN, CIN, etc.
	EmergencyContactName string    `json:"emergency_contact_name" db:"emergency_contact_name"`
	EmergencyContactPhone string   `json:"emergency_contact_phone" db:"emergency_contact_phone"`
	EmergencyContactRelationship string `json:"emergency_contact_relationship" db:"emergency_contact_relationship"`
	MedicalHistory       string    `json:"medical_history" db:"medical_history"`
	Allergies            string    `json:"allergies" db:"allergies"`
	Medications          string    `json:"medications" db:"medications"`
	BloodType            string    `json:"blood_type" db:"blood_type" validate:"oneof=A+ A- B+ B- AB+ AB- O+ O-"`
	IsActive             bool      `json:"is_active" db:"is_active"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy            int       `json:"created_by" db:"created_by"`
}

// PatientRequest represents patient creation/update request
type PatientRequest struct {
	PatientID            string    `json:"patient_id"`
	FirstName            string    `json:"first_name" validate:"required"`
	LastName             string    `json:"last_name" validate:"required"`
	DateOfBirth          string    `json:"date_of_birth" validate:"required"` // Accept as string first
	Gender               string    `json:"gender"`
	Phone                string    `json:"phone"`
	Email                string    `json:"email" validate:"omitempty,email"`
	Address              string    `json:"address"`
	CountryID            int       `json:"country_id" validate:"required"`                  // Location service country ID
	StateID              *int      `json:"state_id,omitempty"`                               // Location service state ID (optional)
	CityID               *int      `json:"city_id,omitempty"`                                // Location service city ID (optional)
	PostalCode           string    `json:"postal_code"`
	NationalityID        *int      `json:"nationality_id,omitempty"`                        // Location service nationality ID (optional)
	PreferredLanguage    string    `json:"preferred_language"`
	MaritalStatus        string    `json:"marital_status"`
	Occupation           string    `json:"occupation"`
	InsuranceTypeID      *int      `json:"insurance_type_id,omitempty"`
	InsuranceProviderID  *int      `json:"insurance_provider_id,omitempty"`
	PolicyNumber         string    `json:"policy_number"`
	NationalID           string    `json:"national_id"`
	EmergencyContactName string    `json:"emergency_contact_name"`
	EmergencyContactPhone string   `json:"emergency_contact_phone"`
	EmergencyContactRelationship string `json:"emergency_contact_relationship"`
	MedicalHistory       string    `json:"medical_history"`
	Allergies            string    `json:"allergies"`
	Medications          string    `json:"medications"`
	BloodType            string    `json:"blood_type" validate:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
}

// PatientSearchRequest represents search parameters
type PatientSearchRequest struct {
	HealthcareEntityID int    `form:"healthcare_entity_id"`
	Query             string `form:"q"`
	Gender            string `form:"gender"`
	Country           string `form:"country"`
	State             string `form:"state"`
	InsuranceTypeID   *int   `form:"insurance_type_id"`
	BloodType         string `form:"blood_type"`
	PreferredLanguage string `form:"preferred_language"`
	AgeMin            int    `form:"age_min"`
	AgeMax            int    `form:"age_max"`
	Limit             int    `form:"limit"`
	Offset            int    `form:"offset"`
}

// PatientResponse represents patient data returned to client
type PatientResponse struct {
	ID                   int       `json:"id"`
	HealthcareEntityID   int       `json:"healthcare_entity_id"`
	PatientID            string    `json:"patient_id"`
	FirstName            string    `json:"first_name"`
	LastName             string    `json:"last_name"`
	DateOfBirth          time.Time `json:"date_of_birth"`
	Gender               string    `json:"gender"`
	Phone                string    `json:"phone"`
	Email                string    `json:"email"`
	Address              string    `json:"address"`
	CountryID            int       `json:"country_id"`                                       // Location service country ID
	StateID              *int      `json:"state_id,omitempty"`                              // Location service state ID
	CityID               *int      `json:"city_id,omitempty"`                               // Location service city ID
	PostalCode           string    `json:"postal_code"`
	NationalityID        *int      `json:"nationality_id,omitempty"`                       // Location service nationality ID
	PreferredLanguage    string    `json:"preferred_language"`
	MaritalStatus        string    `json:"marital_status"`
	Occupation           string    `json:"occupation"`
	InsuranceTypeID      *int      `json:"insurance_type_id,omitempty"`
	InsuranceProviderID  *int      `json:"insurance_provider_id,omitempty"`
	PolicyNumber         string    `json:"policy_number"`
	NationalID           string    `json:"national_id"`
	EmergencyContactName string    `json:"emergency_contact_name"`
	EmergencyContactPhone string   `json:"emergency_contact_phone"`
	EmergencyContactRelationship string `json:"emergency_contact_relationship"`
	MedicalHistory       string    `json:"medical_history"`
	Allergies            string    `json:"allergies"`
	Medications          string    `json:"medications"`
	BloodType            string    `json:"blood_type"`
	Age                  int       `json:"age"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ToPatientResponse converts Patient to PatientResponse
func (p *Patient) ToPatientResponse() PatientResponse {
	return PatientResponse{
		ID:                 p.ID,
		HealthcareEntityID: p.HealthcareEntityID,
		PatientID:          p.PatientID,
		FirstName:          p.FirstName,
		LastName:           p.LastName,
		DateOfBirth:        p.DateOfBirth,
		Gender:             p.Gender,
		Phone:              p.Phone,
		Email:              p.Email,
		Address:            p.Address,
		CountryID:          p.CountryID,
		StateID:            p.StateID,
		CityID:             p.CityID,
		PostalCode:         p.PostalCode,
		NationalityID:      p.NationalityID,
		PreferredLanguage:  p.PreferredLanguage,
		MaritalStatus:      p.MaritalStatus,
		Occupation:         p.Occupation,
		InsuranceTypeID:    p.InsuranceTypeID,
		InsuranceProviderID: p.InsuranceProviderID,
		PolicyNumber:       p.PolicyNumber,
		NationalID:         p.NationalID,
		EmergencyContactName: p.EmergencyContactName,
		EmergencyContactPhone: p.EmergencyContactPhone,
		EmergencyContactRelationship: p.EmergencyContactRelationship,
		MedicalHistory:     p.MedicalHistory,
		Allergies:          p.Allergies,
		Medications:        p.Medications,
		BloodType:          p.BloodType,
		Age:                calculateAge(p.DateOfBirth),
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}

// calculateAge calculates age from date of birth
func calculateAge(birthDate time.Time) int {
	today := time.Now()
	age := today.Year() - birthDate.Year()
	
	// Adjust if birthday hasn't occurred this year
	if today.YearDay() < birthDate.YearDay() {
		age--
	}
	
	return age
}