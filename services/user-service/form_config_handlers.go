package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if user has admin role
func (s *UserService) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetHeader("X-User-Role")
		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":     "Access denied",
				"message":   "Only administrators can manage form configurations",
				"timestamp": time.Now().UTC(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetFormMetadata returns the complete form configuration for a healthcare entity
func (s *UserService) GetFormMetadata(c *gin.Context) {
	formType := c.Param("formType")
	if formType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Form type is required",
			"message":   "Please specify a valid form type (patient, appointment)",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Validate form type
	if formType != "patient" && formType != "appointment" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid form type",
			"message":   "Form type must be 'patient' or 'appointment'",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get healthcare entity ID from header
	entityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if entityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Healthcare entity ID is required",
			"message":   "Missing X-Healthcare-Entity-ID header",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a valid integer",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get form metadata
	metadata, err := s.formConfigService.GetFormMetadata(formType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get form metadata",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      metadata,
		"message":   "Form metadata retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// UpdateFieldConfiguration updates a single field configuration
func (s *UserService) UpdateFieldConfiguration(c *gin.Context) {
	formType := c.Param("formType")
	fieldIDStr := c.Param("fieldId")

	// Validate form type
	if formType != "patient" && formType != "appointment" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid form type",
			"message":   "Form type must be 'patient' or 'appointment'",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Parse field ID
	fieldID, err := strconv.Atoi(fieldIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid field ID",
			"message":   "Field ID must be a valid integer",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get healthcare entity ID from header
	entityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if entityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Healthcare entity ID is required",
			"message":   "Missing X-Healthcare-Entity-ID header",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a valid integer",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Parse request body
	var req UpdateFieldConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request body",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Update field configuration
	err = s.formConfigService.UpdateFieldConfiguration(entityID, fieldID, req)
	if err != nil {
		// Check if it's a validation error
		if err.Error() == "cannot disable core field" || 
		   err.Error() == "cannot require a disabled field" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Validation error",
				"message":   err.Error(),
				"timestamp": time.Now().UTC(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to update field configuration",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Field configuration updated successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetFieldConfigurations returns all field configurations for a form type
func (s *UserService) GetFieldConfigurations(c *gin.Context) {
	formType := c.Param("formType")

	// Validate form type
	if formType != "patient" && formType != "appointment" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid form type",
			"message":   "Form type must be 'patient' or 'appointment'",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get healthcare entity ID from header
	entityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if entityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Healthcare entity ID is required",
			"message":   "Missing X-Healthcare-Entity-ID header",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a valid integer",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get form metadata
	metadata, err := s.formConfigService.GetFormMetadata(formType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get field configurations",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      metadata.Fields,
		"message":   "Field configurations retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// ResetFormConfiguration resets form configuration to defaults
func (s *UserService) ResetFormConfiguration(c *gin.Context) {
	formType := c.Param("formType")

	// Validate form type
	if formType != "patient" && formType != "appointment" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid form type",
			"message":   "Form type must be 'patient' or 'appointment'",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get healthcare entity ID from header
	entityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if entityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Healthcare entity ID is required",
			"message":   "Missing X-Healthcare-Entity-ID header",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a valid integer",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Reset form configuration
	err = s.formConfigService.ResetFormConfiguration(formType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to reset form configuration",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Form configuration reset to defaults successfully",
		"timestamp": time.Now().UTC(),
	})
}

// UpdateMultipleFieldConfigurations updates multiple field configurations atomically
func (s *UserService) UpdateMultipleFieldConfigurations(c *gin.Context) {
	formType := c.Param("formType")

	// Validate form type
	if formType != "patient" && formType != "appointment" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid form type",
			"message":   "Form type must be 'patient' or 'appointment'",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get healthcare entity ID from header
	entityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if entityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Healthcare entity ID is required",
			"message":   "Missing X-Healthcare-Entity-ID header",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a valid integer",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Parse request body
	var req FormConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request body",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Validate request
	if len(req.Fields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No fields to update",
			"message":   "Request must contain at least one field configuration",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Update multiple field configurations
	err = s.formConfigService.UpdateMultipleFieldConfigurations(entityID, req)
	if err != nil {
		// Check if it's a validation error
		if err.Error() == "cannot disable core field" || 
		   err.Error() == "cannot require a disabled field" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Validation error",
				"message":   err.Error(),
				"timestamp": time.Now().UTC(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to update field configurations",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Field configurations updated successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetFormTypes returns all available form types
func (s *UserService) GetFormTypes(c *gin.Context) {
	formTypes, err := s.formConfigService.GetFormTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get form types",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      formTypes,
		"message":   "Form types retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// UpdateFieldOrder updates field sort order for multiple fields
func (s *UserService) UpdateFieldOrder(c *gin.Context) {
	formType := c.Param("formType")

	// Validate form type
	if formType != "patient" && formType != "appointment" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid form type",
			"message":   "Form type must be 'patient' or 'appointment'",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get healthcare entity ID from header
	entityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if entityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Healthcare entity ID is required",
			"message":   "Missing X-Healthcare-Entity-ID header",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a valid integer",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Parse request body
	var req struct {
		FieldOrders []struct {
			FieldID   int `json:"field_id"`
			SortOrder int `json:"sort_order"`
		} `json:"field_orders"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request body",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Validate request
	if len(req.FieldOrders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "No field orders to update",
			"message":   "Request must contain at least one field order",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Update field orders
	err = s.formConfigService.UpdateFieldOrders(entityID, req.FieldOrders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to update field orders",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Field orders updated successfully",
		"timestamp": time.Now().UTC(),
	})
}

// Internationalization Handlers

// GetSupportedLocales returns all supported locales
func (s *UserService) GetSupportedLocales(c *gin.Context) {
	locales, err := s.translationService.GetSupportedLocales()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get supported locales",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      locales,
		"message":   "Supported locales retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// UpdateUserLocale updates user's preferred locale
func (s *UserService) UpdateUserLocale(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "User ID is required",
			"message":   "Missing X-User-ID header",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid user ID",
			"message":   "User ID must be a valid integer",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	var req UserLocaleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request body",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Validate locale
	if err := s.translationService.ValidateLocale(req.PreferredLocale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid locale",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Update user's preferred locale
	query := `UPDATE users SET preferred_locale = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err = s.db.Exec(query, req.PreferredLocale, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to update user locale",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "User locale updated successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetTranslations returns all translations for a locale
func (s *UserService) GetTranslations(c *gin.Context) {
	locale := c.Param("locale")
	if locale == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Locale is required",
			"message":   "Please specify a locale",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	context := c.Query("context") // Optional context filter

	translations, err := s.translationService.GetTranslations(locale, context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get translations",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      translations,
		"locale":    locale,
		"context":   context,
		"message":   "Translations retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}

// CreateOrUpdateTranslation creates or updates a translation
func (s *UserService) CreateOrUpdateTranslation(c *gin.Context) {
	var req TranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request body",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Validate locale
	if err := s.translationService.ValidateLocale(req.Locale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid locale",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	translation, err := s.translationService.CreateOrUpdateTranslation(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to create/update translation",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      translation,
		"message":   "Translation created/updated successfully",
		"timestamp": time.Now().UTC(),
	})
}

// CreateOrUpdateFieldTranslation creates or updates a field translation
func (s *UserService) CreateOrUpdateFieldTranslation(c *gin.Context) {
	var req FieldTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid request body",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Validate locale
	if err := s.translationService.ValidateLocale(req.Locale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid locale",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	fieldTranslation, err := s.translationService.CreateOrUpdateFieldTranslation(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to create/update field translation",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      fieldTranslation,
		"message":   "Field translation created/updated successfully",
		"timestamp": time.Now().UTC(),
	})
}

// GetLocalizedFormMetadata returns form metadata with translations
func (s *UserService) GetLocalizedFormMetadata(c *gin.Context) {
	formType := c.Param("formType")
	if formType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Form type is required",
			"message":   "Please specify a valid form type (patient, appointment)",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Validate form type
	if formType != "patient" && formType != "appointment" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid form type",
			"message":   "Form type must be 'patient' or 'appointment'",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get locale from URL param (for internal API) or query param (for regular API)
	locale := c.Param("locale")
	if locale == "" {
		locale = c.Query("locale")
	}
	
	// If no locale specified, use user's preferred locale or entity locale as fallback
	if locale == "" {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr != "" {
			if userID, err := strconv.Atoi(userIDStr); err == nil {
				// Get user's preferred locale
				var userLocale sql.NullString
				err := s.db.QueryRow(`SELECT preferred_locale FROM users WHERE id = $1`, userID).Scan(&userLocale)
				if err == nil && userLocale.Valid && userLocale.String != "" {
					locale = userLocale.String
				}
			}
		}
		
		// If still no locale, get entity locale or default to en-US
		if locale == "" {
			entityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
			if entityIDStr != "" {
				if entityID, err := strconv.Atoi(entityIDStr); err == nil {
					var entityLocale sql.NullString
					err := s.db.QueryRow(`SELECT locale FROM healthcare_entities WHERE id = $1`, entityID).Scan(&entityLocale)
					if err == nil && entityLocale.Valid && entityLocale.String != "" {
						locale = entityLocale.String
					}
				}
			}
		}
		
		// Final fallback
		if locale == "" {
			locale = "en-US"
		}
	}

	// Get healthcare entity ID from header
	entityIDStr := c.GetHeader("X-Healthcare-Entity-ID")
	if entityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Healthcare entity ID is required",
			"message":   "Missing X-Healthcare-Entity-ID header",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid healthcare entity ID",
			"message":   "Healthcare entity ID must be a valid integer",
			"timestamp": time.Now().UTC(),
		})
		return
	}

	// Get localized form metadata
	metadata, err := s.translationService.GetLocalizedFormMetadata(formType, entityID, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Failed to get localized form metadata",
			"message":   err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      metadata,
		"message":   "Localized form metadata retrieved successfully",
		"timestamp": time.Now().UTC(),
	})
}