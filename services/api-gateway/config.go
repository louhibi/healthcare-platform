package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents the gateway configuration
type Config struct {
	Port         string                    `json:"port"`
	Services     map[string]ServiceConfig  `json:"services"`
	Routes       []RouteConfig            `json:"routes"`
	RateLimit    RateLimitConfig          `json:"rate_limit"`
	JWTSecret    string                   `json:"jwt_secret"`
	Timeout      time.Duration            `json:"timeout"`
	LogLevel     string                   `json:"log_level"`
	CORS         CORSConfig               `json:"cors"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Port: getEnv("PORT", "8080"),
		Services: map[string]ServiceConfig{
			"user-service": {
				Name:    "user-service",
				BaseURL: getEnv("USER_SERVICE_URL", "http://user-service:8081"),
				Timeout: getEnvInt("USER_SERVICE_TIMEOUT", 30),
			},
			"patient-service": {
				Name:    "patient-service", 
				BaseURL: getEnv("PATIENT_SERVICE_URL", "http://patient-service:8082"),
				Timeout: getEnvInt("PATIENT_SERVICE_TIMEOUT", 30),
			},
			"appointment-service": {
				Name:    "appointment-service",
				BaseURL: getEnv("APPOINTMENT_SERVICE_URL", "http://appointment-service:8083"),
				Timeout: getEnvInt("APPOINTMENT_SERVICE_TIMEOUT", 30),
			},
			"location-service": {
				Name:    "location-service",
				BaseURL: getEnv("LOCATION_SERVICE_URL", "http://location-service:8084"),
				Timeout: getEnvInt("LOCATION_SERVICE_TIMEOUT", 15),
			},
		},
		Routes: []RouteConfig{
			// Auth routes (no auth required) - specific endpoints
			{
				Path:        "/api/auth/login",
				Service:     "user-service",
				StripPrefix: false,
				AuthRequired: false,
				RolesAllowed: []string{},
			},
			{
				Path:        "/api/auth/register",
				Service:     "user-service",
				StripPrefix: false,
				AuthRequired: false,
				RolesAllowed: []string{},
			},
			{
				Path:        "/api/auth/refresh",
				Service:     "user-service",
				StripPrefix: false,
				AuthRequired: false,
				RolesAllowed: []string{},
			},
			// entities routes (auth required)
			{
				Path:        "/api/entities/*path",
				Service:     "user-service",
				StripPrefix: false,
				AuthRequired: true,
				RolesAllowed: []string{"admin", "doctor", "nurse", "staff"},
			},
			// User routes (auth required)
			{
				Path:        "/api/users/profile",
				Service:     "user-service",
				StripPrefix: false,
				AuthRequired: true,
				RolesAllowed: []string{"admin", "doctor", "nurse", "staff"},
			},
			{
				Path:        "/api/users/",
				Service:     "user-service",
				StripPrefix: false,
				AuthRequired: true,
				RolesAllowed: []string{"admin"},
			},
			// Patient routes (auth required)
			{
				Path:        "/api/patients/*any",
				Service:     "patient-service",
				StripPrefix: false,
				AuthRequired: true,
				RolesAllowed: []string{"admin", "doctor", "nurse", "staff"},
			},
			// Appointment routes (auth required)
			{
				Path:        "/api/appointments/*path",
				Service:     "appointment-service",
				StripPrefix: false,
				AuthRequired: true,
				RolesAllowed: []string{"admin", "doctor", "nurse", "staff"},
			},
			// Admin routes for appointment service (admin only)
			{
				Path:        "/api/admin/*path",
				Service:     "appointment-service",
				StripPrefix: false,
				AuthRequired: true,
				RolesAllowed: []string{"admin"},
			},
			// Doctor routes for appointment service (auth required)
			{
				Path:        "/api/doctors/*path",
				Service:     "appointment-service",
				StripPrefix: false,
				AuthRequired: true,
				RolesAllowed: []string{"admin", "doctor"},
			},
			// Availability routes for appointment service (auth required)
			{
				Path:        "/api/availability/*path",
				Service:     "appointment-service",
				StripPrefix: false,
				AuthRequired: true,
				RolesAllowed: []string{"admin", "doctor", "nurse", "staff"},
			},
			// Form configuration routes (admin only)
			{
				Path:        "/api/forms/*path",
				Service:     "user-service",
				StripPrefix: false,
				AuthRequired: true,
				RolesAllowed: []string{"admin", "doctor", "nurse", "staff"},

			},
			// Internationalization routes (auth required)
			{
				Path:        "/api/i18n/*path",
				Service:     "user-service",
				StripPrefix: false,
				AuthRequired: true,
				RolesAllowed: []string{"admin", "doctor", "nurse", "staff"},
			},
			// Location routes (no auth required)
			{
				Path:        "/api/locations/*path",
				Service:     "location-service",
				StripPrefix: false,
				AuthRequired: false,
				RolesAllowed: []string{},
			},
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getEnvInt("RATE_LIMIT_RPM", 100),
			BurstSize:         getEnvInt("RATE_LIMIT_BURST", 20),
		},
		JWTSecret: getEnv("JWT_SECRET", "your-very-secure-secret-key-change-in-production"),
		Timeout:   time.Duration(getEnvInt("GATEWAY_TIMEOUT", 30)) * time.Second,
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ","),
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders: []string{"Content-Type", "Authorization", "X-User-ID", "X-User-Email", "X-User-Role", "X-Healthcare-Entity-ID", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform", "User-Agent", "Accept", "Referer"},
		},
	}

	return config
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets environment variable as integer with default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvSlice gets environment variable as comma-separated slice with default value
func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}