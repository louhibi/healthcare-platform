package main

import (
	"time"
)

// ServiceConfig represents configuration for a microservice
type ServiceConfig struct {
	Name    string `json:"name"`
	BaseURL string `json:"base_url"`
	Timeout int    `json:"timeout"` // in seconds
}

// RouteConfig represents route configuration
type RouteConfig struct {
	Path        string `json:"path"`
	Service     string `json:"service"`
	StripPrefix bool   `json:"strip_prefix"`
	AuthRequired bool  `json:"auth_required"`
	RolesAllowed []string `json:"roles_allowed"`
}

// UserClaims represents JWT claims
type UserClaims struct {
	UserID              int    `json:"user_id"`
	Email               string `json:"email"`
	Role                string `json:"role"`
	HealthcareEntityID  int    `json:"healthcare_entity_id"`
}

// ProxyRequest represents a request to be proxied
type ProxyRequest struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

// ProxyResponse represents response from proxied service
type ProxyResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
}

// HealthCheck represents health check response
type HealthCheck struct {
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// GatewayStats represents gateway statistics
type GatewayStats struct {
	TotalRequests   int64              `json:"total_requests"`
	SuccessRequests int64              `json:"success_requests"`
	ErrorRequests   int64              `json:"error_requests"`
	AvgResponseTime float64            `json:"avg_response_time_ms"`
	ServiceStats    map[string]ServiceStats `json:"service_stats"`
}

// ServiceStats represents statistics for a specific service
type ServiceStats struct {
	Requests    int64   `json:"requests"`
	Successes   int64   `json:"successes"`
	Errors      int64   `json:"errors"`
	AvgLatency  float64 `json:"avg_latency_ms"`
	LastRequest time.Time `json:"last_request"`
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	BurstSize         int `json:"burst_size"`
}

// ErrorResponse represents error response format
type ErrorResponse struct {
	Error   string    `json:"error"`
	Code    string    `json:"code,omitempty"`
	Message string    `json:"message,omitempty"`
	Time    time.Time `json:"timestamp"`
}