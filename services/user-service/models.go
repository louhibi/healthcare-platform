package main

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                   int       `json:"id" db:"id"`
	Email                string    `json:"email" db:"email" validate:"required,email"`
	Password             string    `json:"-" db:"password_hash"`
	FirstName            string    `json:"first_name" db:"first_name" validate:"required"`
	LastName             string    `json:"last_name" db:"last_name" validate:"required"`
	Role                 string    `json:"role" db:"role" validate:"required,oneof=admin doctor nurse staff"`
	HealthcareEntityID   int       `json:"healthcare_entity_id" db:"healthcare_entity_id" validate:"required"`
	LicenseNumber        string    `json:"license_number" db:"license_number"` // Professional license for doctors/nurses
	Specialization       string    `json:"specialization" db:"specialization"` // Medical specialization
	PreferredLocale      string    `json:"preferred_locale" db:"preferred_locale" validate:"oneof=en-US en-CA fr-CA fr-FR ar-MA"`
	IsActive             bool      `json:"is_active" db:"is_active"`
	IsTempPassword       bool      `json:"is_temp_password" db:"is_temp_password"`
	TempPasswordExpires  *time.Time `json:"temp_password_expires" db:"temp_password_expires"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// UserRegistration represents registration request
type UserRegistration struct {
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8"`
	FirstName            string `json:"first_name" validate:"required"`
	LastName             string `json:"last_name" validate:"required"`
	Role                 string `json:"role" validate:"required,oneof=admin doctor nurse staff"`
	HealthcareEntityID   int    `json:"healthcare_entity_id" validate:"required"`
	LicenseNumber        string `json:"license_number"`
	Specialization       string `json:"specialization"`
	PreferredLocale      string `json:"preferred_locale" validate:"oneof=en-US en-CA fr-CA fr-FR ar-MA"`
}

// UserLogin represents login request
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents user data returned to client (no sensitive info)
type UserResponse struct {
	ID                   int       `json:"id"`
	Email                string    `json:"email"`
	FirstName            string    `json:"first_name"`
	LastName             string    `json:"last_name"`
	Role                 string    `json:"role"`
	HealthcareEntityID   int       `json:"healthcare_entity_id"`
	LicenseNumber        string    `json:"license_number"`
	Specialization       string    `json:"specialization"`
	PreferredLocale      string    `json:"preferred_locale"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Claims represents JWT claims
type Claims struct {
	UserID             int    `json:"user_id"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	HealthcareEntityID int    `json:"healthcare_entity_id"`
	IsTempPassword     bool   `json:"is_temp_password"`
}

// AdminCreateDoctorRequest represents admin request to create a doctor
type AdminCreateDoctorRequest struct {
	Email             string `json:"email" validate:"required,email"`
	FirstName         string `json:"first_name" validate:"required"`
	LastName          string `json:"last_name" validate:"required"`
	LicenseNumber     string `json:"license_number" validate:"required"`
	Specialization    string `json:"specialization" validate:"required"`
	PreferredLocale   string `json:"preferred_locale" validate:"oneof=en-US en-CA fr-CA fr-FR ar-MA"`
}

// AdminCreateDoctorResponse represents response after creating a doctor
type AdminCreateDoctorResponse struct {
	Doctor          UserResponse `json:"doctor"`
	TempPassword    string       `json:"temp_password"`
	PasswordExpires time.Time    `json:"password_expires"`
	Message         string       `json:"message"`
}

// PasswordChangeRequest represents password change request
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// ToUserResponse converts User to UserResponse
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:                 u.ID,
		Email:              u.Email,
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		Role:               u.Role,
		HealthcareEntityID: u.HealthcareEntityID,
		LicenseNumber:      u.LicenseNumber,
		Specialization:     u.Specialization,
		PreferredLocale:    u.PreferredLocale,
		IsActive:           u.IsActive,
		CreatedAt:          u.CreatedAt,
	}
}

// Form Configuration Models

// FormType represents a form type (patient, appointment)
type FormType struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,oneof=patient appointment"`
	DisplayName string    `json:"display_name" db:"display_name" validate:"required"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// FieldDefinition represents a field definition with metadata
type FieldDefinition struct {
	ID               int                    `json:"id" db:"id"`
	FormTypeID       int                    `json:"form_type_id" db:"form_type_id" validate:"required"`
	Name             string                 `json:"name" db:"name" validate:"required"`
	DisplayName      string                 `json:"display_name" db:"display_name" validate:"required"`
	FieldType        string                 `json:"field_type" db:"field_type" validate:"required,oneof=text email phone number date datetime select textarea checkbox"`
	DefaultRequired  bool                   `json:"default_required" db:"default_required"`
	IsCore           bool                   `json:"is_core" db:"is_core"` // Core fields cannot be disabled
	ValidationRules  map[string]interface{} `json:"validation_rules" db:"validation_rules"`
	Options          []string               `json:"options" db:"options"` // For select fields
	SortOrder        int                    `json:"sort_order" db:"sort_order"`
	Category         string                 `json:"category" db:"category"`
	Description      string                 `json:"description" db:"description"`
	PlaceholderText  string                 `json:"placeholder_text" db:"placeholder_text"`
	IsActive         bool                   `json:"is_active" db:"is_active"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
}

// EntityFieldConfiguration represents per-entity field customization
type EntityFieldConfiguration struct {
	ID                 int                    `json:"id" db:"id"`
	HealthcareEntityID int                    `json:"healthcare_entity_id" db:"healthcare_entity_id" validate:"required"`
	FieldDefinitionID  int                    `json:"field_definition_id" db:"field_definition_id" validate:"required"`
	IsEnabled          bool                   `json:"is_enabled" db:"is_enabled"`
	IsRequired         bool                   `json:"is_required" db:"is_required"`
	CustomLabel        string                 `json:"custom_label" db:"custom_label"`
	CustomValidation   map[string]interface{} `json:"custom_validation" db:"custom_validation"`
	SortOrder          int                    `json:"sort_order" db:"sort_order"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"updated_at"`
}

// FormMetadataResponse represents the complete form configuration for a healthcare entity
type FormMetadataResponse struct {
	FormType      string                       `json:"form_type"`
	DisplayName   string                       `json:"display_name"`
	Description   string                       `json:"description"`
	Fields        []FieldConfigurationResponse `json:"fields"`
	EntityID      int                          `json:"entity_id"`
	LastModified  time.Time                    `json:"last_modified"`
}

// FieldConfigurationResponse represents a field configuration for API response
type FieldConfigurationResponse struct {
	FieldID         int                    `json:"field_id"`
	Name            string                 `json:"name"`
	DisplayName     string                 `json:"display_name"`
	CustomLabel     string                 `json:"custom_label,omitempty"`
	FieldType       string                 `json:"field_type"`
	IsEnabled       bool                   `json:"is_enabled"`
	IsRequired      bool                   `json:"is_required"`
	IsCore          bool                   `json:"is_core"`
	ValidationRules map[string]interface{} `json:"validation_rules"`
	CustomValidation map[string]interface{} `json:"custom_validation,omitempty"`
	Options         []string               `json:"options,omitempty"`
	SortOrder       int                    `json:"sort_order"`
	Category        string                 `json:"category"`
	Description     string                 `json:"description"`
	PlaceholderText string                 `json:"placeholder_text"`
}

// UpdateFieldConfigRequest represents a request to update field configuration
type UpdateFieldConfigRequest struct {
	IsEnabled        *bool                  `json:"is_enabled,omitempty"`
	IsRequired       *bool                  `json:"is_required,omitempty"`
	CustomLabel      string                 `json:"custom_label"`
	CustomValidation map[string]interface{} `json:"custom_validation"`
	SortOrder        int                    `json:"sort_order"`
}

// FormConfigurationRequest represents a request to update multiple field configurations
type FormConfigurationRequest struct {
	Fields []UpdateFieldWithIDRequest `json:"fields" validate:"required"`
}

// UpdateFieldWithIDRequest represents a field update request with field ID
type UpdateFieldWithIDRequest struct {
	FieldID          int                    `json:"field_id" validate:"required"`
	IsEnabled        bool                   `json:"is_enabled"`
	IsRequired       bool                   `json:"is_required"`
	CustomLabel      string                 `json:"custom_label"`
	CustomValidation map[string]interface{} `json:"custom_validation"`
	SortOrder        int                    `json:"sort_order"`
}

// Internationalization Models

// Locale represents a supported locale
type Locale struct {
	Code         string    `json:"code" db:"code"`
	LanguageName string    `json:"language_name" db:"language_name"`
	NativeName   string    `json:"native_name" db:"native_name"`
	CountryCode  string    `json:"country_code" db:"country_code"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Translation represents a translation entry
type Translation struct {
	ID             int       `json:"id" db:"id"`
	TranslationKey string    `json:"translation_key" db:"translation_key"`
	Locale         string    `json:"locale" db:"locale"`
	Content        string    `json:"content" db:"content"`
	Context        string    `json:"context" db:"context"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// FieldTranslation represents localized field metadata
type FieldTranslation struct {
	ID                int       `json:"id" db:"id"`
	FieldDefinitionID int       `json:"field_definition_id" db:"field_definition_id"`
	Locale            string    `json:"locale" db:"locale"`
	DisplayName       string    `json:"display_name" db:"display_name"`
	Description       string    `json:"description" db:"description"`
	PlaceholderText   string    `json:"placeholder_text" db:"placeholder_text"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// LocalizedFormMetadataResponse represents form metadata with translations
type LocalizedFormMetadataResponse struct {
	FormType      string                                 `json:"form_type"`
	DisplayName   string                                 `json:"display_name"`
	Description   string                                 `json:"description"`
	Fields        []LocalizedFieldConfigurationResponse `json:"fields"`
	EntityID      int                                    `json:"entity_id"`
	Locale        string                                 `json:"locale"`
	LastModified  time.Time                              `json:"last_modified"`
}

// LocalizedFieldConfigurationResponse represents a field configuration with translations
type LocalizedFieldConfigurationResponse struct {
	FieldID         int                    `json:"field_id"`
	Name            string                 `json:"name"`
	DisplayName     string                 `json:"display_name"`
	CustomLabel     string                 `json:"custom_label,omitempty"`
	FieldType       string                 `json:"field_type"`
	IsEnabled       bool                   `json:"is_enabled"`
	IsRequired      bool                   `json:"is_required"`
	IsCore          bool                   `json:"is_core"`
	ValidationRules map[string]interface{} `json:"validation_rules"`
	CustomValidation map[string]interface{} `json:"custom_validation,omitempty"`
	Options         []string               `json:"options,omitempty"`
	SortOrder       int                    `json:"sort_order"`
	Category        string                 `json:"category"`
	Description     string                 `json:"description"`
	PlaceholderText string                 `json:"placeholder_text"`
	Locale          string                 `json:"locale"`
}

// TranslationRequest represents a request to create/update translations
type TranslationRequest struct {
	TranslationKey string `json:"translation_key" validate:"required"`
	Locale         string `json:"locale" validate:"required"`
	Content        string `json:"content" validate:"required"`
	Context        string `json:"context"`
}

// FieldTranslationRequest represents a request to create/update field translations
type FieldTranslationRequest struct {
	FieldDefinitionID int    `json:"field_definition_id" validate:"required"`
	Locale            string `json:"locale" validate:"required,oneof=en-US en-CA fr-CA fr-FR ar-MA"`
	DisplayName       string `json:"display_name" validate:"required"`
	Description       string `json:"description"`
	PlaceholderText   string `json:"placeholder_text"`
}

// UserLocaleUpdateRequest represents a request to update user's locale preference
type UserLocaleUpdateRequest struct {
	PreferredLocale string `json:"preferred_locale" validate:"required,oneof=en-US en-CA fr-CA fr-FR ar-MA"`
}