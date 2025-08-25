package logging

import (
	"log/slog"
	"os"
	"strings"
)

// InitLogger configures slog with appropriate level and handler
func InitLogger(serviceName string) {
	// Get log level from environment variable, default to INFO
	logLevel := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	if logLevel == "" {
		logLevel = "INFO"
	}

	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Configure JSON handler for structured logging
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		AddSource: true, // Add source code location
	})

	// Create logger with service context
	logger := slog.New(handler).With("service", serviceName)
	slog.SetDefault(logger)

	slog.Info("Logger initialized", 
		"service", serviceName,
		"level", level,
		"format", "json")
}

// LogRequest logs incoming HTTP request details at debug level
func LogRequest(method, path, userAgent, clientIP string, headers map[string]string, body interface{}) {
	slog.Debug("Incoming request",
		"method", method,
		"path", path,
		"user_agent", userAgent,
		"client_ip", clientIP,
		"headers", headers,
		"body", body)
}

// LogResponse logs HTTP response details at debug level
func LogResponse(statusCode int, responseTime string, body interface{}) {
	slog.Debug("Outgoing response",
		"status_code", statusCode,
		"response_time", responseTime,
		"body", body)
}

// LogHealthcareContext logs healthcare-specific context (entity, user, etc.)
func LogHealthcareContext(operation string, userID, entityID, role string, additionalFields map[string]interface{}) {
	fields := []interface{}{
		"operation", operation,
		"user_id", userID,
		"healthcare_entity_id", entityID,
		"user_role", role,
	}
	
	// Add additional fields
	for key, value := range additionalFields {
		fields = append(fields, key, value)
	}
	
	slog.Info("Healthcare operation", fields...)
}

// LogError logs errors with context
func LogError(operation string, err error, fields map[string]interface{}) {
	logFields := []interface{}{
		"operation", operation,
		"error", err.Error(),
	}
	
	for key, value := range fields {
		logFields = append(logFields, key, value)
	}
	
	slog.Error("Operation failed", logFields...)
}

// LogDatabaseOperation logs database operations for audit trail
func LogDatabaseOperation(operation, table string, entityID string, recordID interface{}, userID string) {
	slog.Info("Database operation",
		"db_operation", operation,
		"table", table,
		"healthcare_entity_id", entityID,
		"record_id", recordID,
		"user_id", userID)
}