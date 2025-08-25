package main

import (
	"time"
)

// HealthcareEntity represents a hospital, clinic, or doctor's office
type HealthcareEntity struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required"`
	Type        string    `json:"type" db:"type" validate:"required,oneof=hospital clinic doctor_office"`
	CountryID   *int      `json:"country_id" db:"country_id"`                       // Foreign key to location service (nullable)
	Address     string    `json:"address" db:"address" validate:"required"`
	StateID     *int      `json:"state_id,omitempty" db:"state_id"`                  // Foreign key to location service (nullable)
	CityID      *int      `json:"city_id,omitempty" db:"city_id"`                    // Foreign key to location service (nullable)
	
	// Location names for backward compatibility and display (will be auto-populated by API Gateway)
	Country     string    `json:"country,omitempty" db:"-"`                         // Read-only, populated from location service
	State       string    `json:"state,omitempty" db:"-"`                           // Read-only, populated from location service
	City        string    `json:"city,omitempty" db:"-"`                            // Read-only, populated from location service
	PostalCode  string    `json:"postal_code" db:"postal_code" validate:"required"`
	Phone       string    `json:"phone" db:"phone" validate:"required"`
	Email       string    `json:"email" db:"email" validate:"required,email"`
	Website     string    `json:"website" db:"website"`
	License     string    `json:"license" db:"license"` // Medical license number
	TaxID       string    `json:"tax_id" db:"tax_id"`   // Tax identification number
	Timezone             string    `json:"timezone" db:"timezone" validate:"required"`
	Language             string    `json:"language" db:"language" validate:"required,oneof=en fr ar"`
	Currency             string    `json:"currency" db:"currency" validate:"required,oneof=CAD USD MAD EUR"`
	RequireRoomAssignment bool     `json:"require_room_assignment" db:"require_room_assignment"`
	IsActive             bool      `json:"is_active" db:"is_active"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// HealthcareEntityRequest represents entity creation/update request
type HealthcareEntityRequest struct {
	Name       string `json:"name" validate:"required"`
	Type       string `json:"type" validate:"required,oneof=hospital clinic doctor_office"`
	CountryID  *int   `json:"country_id,omitempty"`                              // Location service country ID (nullable)
	StateID    *int   `json:"state_id,omitempty"`                               // Location service state ID (optional)
	CityID     *int   `json:"city_id,omitempty"`                                // Location service city ID (optional)
	Address    string `json:"address" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Website    string `json:"website"`
	License    string `json:"license"`
	TaxID      string `json:"tax_id"`
	Timezone             string `json:"timezone" validate:"required"`
	Language             string `json:"language" validate:"required,oneof=en fr ar"`
	Currency             string `json:"currency" validate:"required,oneof=CAD USD MAD EUR"`
	RequireRoomAssignment bool   `json:"require_room_assignment"`
}

// HealthcareEntityResponse represents entity data returned to client
type HealthcareEntityResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	CountryID  int       `json:"country_id"`                                     // Location service country ID
	StateID    *int      `json:"state_id,omitempty"`                            // Location service state ID
	CityID     *int      `json:"city_id,omitempty"`                             // Location service city ID
	Country    string    `json:"country"`                                       // Enhanced by API Gateway
	State      string    `json:"state,omitempty"`                               // Enhanced by API Gateway
	City       string    `json:"city,omitempty"`                                // Enhanced by API Gateway
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Website    string    `json:"website"`
	License    string    `json:"license"`
	TaxID      string    `json:"tax_id"`
	Timezone             string    `json:"timezone"`
	Language             string    `json:"language"`
	Currency             string    `json:"currency"`
	RequireRoomAssignment bool     `json:"require_room_assignment"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ToHealthcareEntityResponse converts HealthcareEntity to HealthcareEntityResponse
func (he *HealthcareEntity) ToHealthcareEntityResponse() HealthcareEntityResponse {
	return HealthcareEntityResponse{
		ID:         he.ID,
		Name:       he.Name,
		Type:       he.Type,
		CountryID:  func() int { if he.CountryID != nil { return *he.CountryID } else { return 0 } }(),
		StateID:    he.StateID,
		CityID:     he.CityID,
		Country:    he.Country,    // Will be populated by API Gateway
		State:      he.State,     // Will be populated by API Gateway
		City:       he.City,      // Will be populated by API Gateway
		Address:    he.Address,
		PostalCode: he.PostalCode,
		Phone:      he.Phone,
		Email:      he.Email,
		Website:    he.Website,
		License:    he.License,
		TaxID:      he.TaxID,
		Timezone:             he.Timezone,
		Language:             he.Language,
		Currency:             he.Currency,
		RequireRoomAssignment: he.RequireRoomAssignment,
		IsActive:             he.IsActive,
		CreatedAt:            he.CreatedAt,
		UpdatedAt:            he.UpdatedAt,
	}
}

// CountryConfig holds country-specific configurations
type CountryConfig struct {
	Country         string   `json:"country"`
	States          []string `json:"states,omitempty"`          // For USA/Canada
	PostalCodeRegex string   `json:"postal_code_regex"`         // Validation pattern
	PhoneFormat     string   `json:"phone_format"`              // Format example
	Currency        string   `json:"currency"`                  // Default currency
	Language        string   `json:"language"`                  // Default language
	Timezone        []string `json:"timezones"`                 // Available timezones
	AddressFormat   string   `json:"address_format"`            // Address display format
}

// GetCountryConfigs returns configuration for all supported countries
func GetCountryConfigs() map[string]CountryConfig {
	return map[string]CountryConfig{
		"Canada": {
			Country:         "Canada",
			States:          []string{"AB", "BC", "MB", "NB", "NL", "NS", "NT", "NU", "ON", "PE", "QC", "SK", "YT"},
			PostalCodeRegex: `^[A-Za-z]\d[A-Za-z][ -]?\d[A-Za-z]\d$`,
			PhoneFormat:     "(xxx) xxx-xxxx",
			Currency:        "CAD",
			Language:        "en",
			Timezone:        []string{"America/Toronto", "America/Vancouver", "America/Halifax", "America/Winnipeg"},
			AddressFormat:   "{address}, {city}, {state} {postal_code}, Canada",
		},
		"USA": {
			Country:         "USA",
			States:          []string{"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA", "HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD", "MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ", "NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC", "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY", "DC"},
			PostalCodeRegex: `^\d{5}(-\d{4})?$`,
			PhoneFormat:     "(xxx) xxx-xxxx",
			Currency:        "USD",
			Language:        "en",
			Timezone:       []string{"America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles", "America/Anchorage", "Pacific/Honolulu"},
			AddressFormat:   "{address}, {city}, {state} {postal_code}, USA",
		},
		"Morocco": {
			Country:         "Morocco",
			States:          []string{"Casablanca-Settat", "Rabat-Salé-Kénitra", "Marrakech-Safi", "Fès-Meknès", "Tanger-Tétouan-Al Hoceïma", "Oriental", "Souss-Massa", "Drâa-Tafilalet", "Béni Mellal-Khénifra", "Laâyoune-Sakia El Hamra", "Dakhla-Oued Ed-Dahab", "Guelmim-Oued Noun"},
			PostalCodeRegex: `^\d{5}$`,
			PhoneFormat:     "+212 xxx-xxx-xxx",
			Currency:        "MAD",
			Language:        "ar",
			Timezone:       []string{"Africa/Casablanca"},
			AddressFormat:   "{address}, {city} {postal_code}, {state}, Morocco",
		},
		"France": {
			Country:         "France",
			States:          []string{"Auvergne-Rhône-Alpes", "Bourgogne-Franche-Comté", "Bretagne", "Centre-Val de Loire", "Corse", "Grand Est", "Hauts-de-France", "Île-de-France", "Normandie", "Nouvelle-Aquitaine", "Occitanie", "Pays de la Loire", "Provence-Alpes-Côte d'Azur"},
			PostalCodeRegex: `^\d{5}$`,
			PhoneFormat:     "+33 x xx xx xx xx",
			Currency:        "EUR",
			Language:        "fr",
			Timezone:       []string{"Europe/Paris"},
			AddressFormat:   "{address}, {postal_code} {city}, {state}, France",
		},
	}
}