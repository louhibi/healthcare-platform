package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type FormConfigService struct {
	db *sql.DB
}

func NewFormConfigService(db *sql.DB) *FormConfigService {
	return &FormConfigService{db: db}
}

// GetFormMetadata returns the complete form configuration for a healthcare entity
func (s *FormConfigService) GetFormMetadata(formType string, entityID int) (*FormMetadataResponse, error) {
	// Get form type info
	var formInfo struct {
		ID          int
		DisplayName string
		Description string
	}
	
	err := s.db.QueryRow(`
		SELECT id, display_name, description 
		FROM form_types 
		WHERE name = $1 AND is_active = true
	`, formType).Scan(&formInfo.ID, &formInfo.DisplayName, &formInfo.Description)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("form type '%s' not found", formType)
		}
		return nil, fmt.Errorf("failed to get form type: %w", err)
	}

	// Get field configurations for this entity
	query := `
		SELECT 
			fd.id,
			fd.name,
			fd.display_name,
			fd.field_type,
			fd.default_required,
			fd.is_core,
			fd.validation_rules,
			fd.options,
			fd.sort_order,
			fd.category,
			fd.description,
			fd.placeholder_text,
			COALESCE(efc.is_enabled, true) as is_enabled,
			COALESCE(efc.is_required, fd.default_required) as is_required,
			COALESCE(efc.custom_label, '') as custom_label,
			COALESCE(efc.custom_validation, '{}') as custom_validation,
			COALESCE(efc.sort_order, fd.sort_order) as effective_sort_order,
			COALESCE(efc.updated_at, fd.updated_at) as last_modified
		FROM field_definitions fd
		LEFT JOIN entity_field_configurations efc 
			ON fd.id = efc.field_definition_id 
			AND efc.healthcare_entity_id = $2
		WHERE fd.form_type_id = $1 
			AND fd.is_active = true
		ORDER BY effective_sort_order, fd.id
	`

	rows, err := s.db.Query(query, formInfo.ID, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get field configurations: %w", err)
	}
	defer rows.Close()

	var fields []FieldConfigurationResponse
	var lastModified time.Time

	for rows.Next() {
		var field FieldConfigurationResponse
		var validationRulesJSON, customValidationJSON, optionsJSON []byte
		var rowLastModified time.Time

		err := rows.Scan(
			&field.FieldID,
			&field.Name,
			&field.DisplayName,
			&field.FieldType,
			&field.IsRequired, // This gets overridden below
			&field.IsCore,
			&validationRulesJSON,
			&optionsJSON,
			&field.SortOrder, // This gets overridden below
			&field.Category,
			&field.Description,
			&field.PlaceholderText,
			&field.IsEnabled,
			&field.IsRequired, // Effective is_required from COALESCE
			&field.CustomLabel,
			&customValidationJSON,
			&field.SortOrder, // Effective sort_order from COALESCE
			&rowLastModified,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan field configuration: %w", err)
		}

		// Parse JSON fields
		if len(validationRulesJSON) > 0 {
			if err := json.Unmarshal(validationRulesJSON, &field.ValidationRules); err != nil {
				log.Printf("Warning: failed to parse validation rules for field %s: %v", field.Name, err)
				field.ValidationRules = make(map[string]interface{})
			}
		} else {
			field.ValidationRules = make(map[string]interface{})
		}

		if len(customValidationJSON) > 0 && string(customValidationJSON) != "{}" {
			if err := json.Unmarshal(customValidationJSON, &field.CustomValidation); err != nil {
				log.Printf("Warning: failed to parse custom validation for field %s: %v", field.Name, err)
			}
		}

		if len(optionsJSON) > 0 {
			if err := json.Unmarshal(optionsJSON, &field.Options); err != nil {
				log.Printf("Warning: failed to parse options for field %s: %v", field.Name, err)
				field.Options = []string{}
			}
		} else {
			field.Options = []string{}
		}

		// Backend override to correct field types for dependent selects
		// Ensure city is exposed as a select so frontend renders it properly
		if field.Name == "city" && field.FieldType != "select" {
			field.FieldType = "select"
		}
		// Country should remain/select
		if field.Name == "country" && field.FieldType != "select" {
			field.FieldType = "select"
		}

		// Track latest modification time
		if rowLastModified.After(lastModified) {
			lastModified = rowLastModified
		}

		fields = append(fields, field)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating field configurations: %w", err)
	}

	return &FormMetadataResponse{
		FormType:     formType,
		DisplayName:  formInfo.DisplayName,
		Description:  formInfo.Description,
		Fields:       fields,
		EntityID:     entityID,
		LastModified: lastModified,
	}, nil
}

// UpdateFieldConfiguration updates a single field configuration for an entity
func (s *FormConfigService) UpdateFieldConfiguration(entityID, fieldID int, req UpdateFieldConfigRequest) error {
	// Validate that the field exists and get its metadata
	var fieldInfo struct {
		Name   string
		IsCore bool
	}
	
	err := s.db.QueryRow(`
		SELECT name, is_core 
		FROM field_definitions 
		WHERE id = $1 AND is_active = true
	`, fieldID).Scan(&fieldInfo.Name, &fieldInfo.IsCore)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("field with ID %d not found", fieldID)
		}
		return fmt.Errorf("failed to get field information: %w", err)
	}

	// Get current field configuration to determine current state
	var currentEnabled, currentRequired bool
	err = s.db.QueryRow(`
		SELECT COALESCE(efc.is_enabled, true), COALESCE(efc.is_required, fd.default_required)
		FROM field_definitions fd
		LEFT JOIN entity_field_configurations efc 
			ON fd.id = efc.field_definition_id AND efc.healthcare_entity_id = $2
		WHERE fd.id = $1 AND fd.is_active = true
	`, fieldID, entityID).Scan(&currentEnabled, &currentRequired)
	
	if err != nil {
		return fmt.Errorf("failed to get current field configuration: %w", err)
	}

	// Determine the effective values (use current values if not provided in request)
	effectiveEnabled := currentEnabled
	effectiveRequired := currentRequired
	
	if req.IsEnabled != nil {
		effectiveEnabled = *req.IsEnabled
	}
	if req.IsRequired != nil {
		effectiveRequired = *req.IsRequired
	}

	// Validate core field protection
	if fieldInfo.IsCore && !effectiveEnabled {
		return fmt.Errorf("cannot disable core field: %s", fieldInfo.Name)
	}

	// Validate that required fields must be enabled
	if effectiveRequired && !effectiveEnabled {
		return fmt.Errorf("cannot require a disabled field: %s", fieldInfo.Name)
	}

	// Convert custom validation to JSON
	customValidationJSON, err := json.Marshal(req.CustomValidation)
	if err != nil {
		return fmt.Errorf("failed to marshal custom validation: %w", err)
	}

	// Upsert entity field configuration using effective values
	query := `
		INSERT INTO entity_field_configurations 
		(healthcare_entity_id, field_definition_id, is_enabled, is_required, custom_label, custom_validation, sort_order, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		ON CONFLICT (healthcare_entity_id, field_definition_id)
		DO UPDATE SET
			is_enabled = EXCLUDED.is_enabled,
			is_required = EXCLUDED.is_required,
			custom_label = EXCLUDED.custom_label,
			custom_validation = EXCLUDED.custom_validation,
			sort_order = EXCLUDED.sort_order,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = s.db.Exec(query, entityID, fieldID, effectiveEnabled, effectiveRequired, 
		req.CustomLabel, customValidationJSON, req.SortOrder)
	
	if err != nil {
		return fmt.Errorf("failed to update field configuration: %w", err)
	}

	return nil
}

// CreateDefaultFieldConfigurations creates default field configurations for a new healthcare entity
func (s *FormConfigService) CreateDefaultFieldConfigurations(entityID int) error {
	// Get all field definitions
	rows, err := s.db.Query(`
		SELECT id, default_required, sort_order
		FROM field_definitions 
		WHERE is_active = true
		ORDER BY form_type_id, sort_order
	`)
	if err != nil {
		return fmt.Errorf("failed to get field definitions: %w", err)
	}
	defer rows.Close()

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert default configurations
	insertQuery := `
		INSERT INTO entity_field_configurations 
		(healthcare_entity_id, field_definition_id, is_enabled, is_required, sort_order, created_at, updated_at)
		VALUES ($1, $2, true, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (healthcare_entity_id, field_definition_id) DO NOTHING
	`

	for rows.Next() {
		var fieldID, sortOrder int
		var defaultRequired bool

		if err := rows.Scan(&fieldID, &defaultRequired, &sortOrder); err != nil {
			return fmt.Errorf("failed to scan field definition: %w", err)
		}

		_, err = tx.Exec(insertQuery, entityID, fieldID, defaultRequired, sortOrder)
		if err != nil {
			return fmt.Errorf("failed to create default configuration for field %d: %w", fieldID, err)
		}
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating field definitions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit default configurations: %w", err)
	}

	return nil
}

// UpdateMultipleFieldConfigurations updates multiple field configurations atomically
func (s *FormConfigService) UpdateMultipleFieldConfigurations(entityID int, req FormConfigurationRequest) error {
	if len(req.Fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Validate all fields first
	for _, fieldReq := range req.Fields {
		var fieldInfo struct {
			Name   string
			IsCore bool
		}
		
		err := tx.QueryRow(`
			SELECT name, is_core 
			FROM field_definitions 
			WHERE id = $1 AND is_active = true
		`, fieldReq.FieldID).Scan(&fieldInfo.Name, &fieldInfo.IsCore)
		
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("field with ID %d not found", fieldReq.FieldID)
			}
			return fmt.Errorf("failed to get field information for ID %d: %w", fieldReq.FieldID, err)
		}

		// Validate core field protection
		if fieldInfo.IsCore && !fieldReq.IsEnabled {
			return fmt.Errorf("cannot disable core field: %s", fieldInfo.Name)
		}

		// Validate that required fields must be enabled
		if fieldReq.IsRequired && !fieldReq.IsEnabled {
			return fmt.Errorf("cannot require a disabled field: %s", fieldInfo.Name)
		}
	}

	// Update all fields
	query := `
		INSERT INTO entity_field_configurations 
		(healthcare_entity_id, field_definition_id, is_enabled, is_required, custom_label, custom_validation, sort_order, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		ON CONFLICT (healthcare_entity_id, field_definition_id)
		DO UPDATE SET
			is_enabled = EXCLUDED.is_enabled,
			is_required = EXCLUDED.is_required,
			custom_label = EXCLUDED.custom_label,
			custom_validation = EXCLUDED.custom_validation,
			sort_order = EXCLUDED.sort_order,
			updated_at = CURRENT_TIMESTAMP
	`

	for _, fieldReq := range req.Fields {
		customValidationJSON, err := json.Marshal(fieldReq.CustomValidation)
		if err != nil {
			return fmt.Errorf("failed to marshal custom validation for field %d: %w", fieldReq.FieldID, err)
		}

		_, err = tx.Exec(query, entityID, fieldReq.FieldID, fieldReq.IsEnabled, 
			fieldReq.IsRequired, fieldReq.CustomLabel, customValidationJSON, fieldReq.SortOrder)
		
		if err != nil {
			return fmt.Errorf("failed to update field configuration for ID %d: %w", fieldReq.FieldID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit field configuration updates: %w", err)
	}

	return nil
}

// ResetFormConfiguration resets form configuration to defaults for an entity
func (s *FormConfigService) ResetFormConfiguration(formType string, entityID int) error {
	// Get form type ID
	var formTypeID int
	err := s.db.QueryRow(`
		SELECT id FROM form_types WHERE name = $1 AND is_active = true
	`, formType).Scan(&formTypeID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("form type '%s' not found", formType)
		}
		return fmt.Errorf("failed to get form type: %w", err)
	}

	// Delete all custom configurations for this form type and entity
	_, err = s.db.Exec(`
		DELETE FROM entity_field_configurations 
		WHERE healthcare_entity_id = $1 
		AND field_definition_id IN (
			SELECT id FROM field_definitions WHERE form_type_id = $2
		)
	`, entityID, formTypeID)

	if err != nil {
		return fmt.Errorf("failed to reset form configuration: %w", err)
	}

	return nil
}

// GetFormTypes returns all available form types
func (s *FormConfigService) GetFormTypes() ([]FormType, error) {
	rows, err := s.db.Query(`
		SELECT id, name, display_name, description, is_active, created_at, updated_at
		FROM form_types 
		WHERE is_active = true
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get form types: %w", err)
	}
	defer rows.Close()

	var formTypes []FormType
	for rows.Next() {
		var ft FormType
		if err := rows.Scan(&ft.ID, &ft.Name, &ft.DisplayName, &ft.Description, 
			&ft.IsActive, &ft.CreatedAt, &ft.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan form type: %w", err)
		}
		formTypes = append(formTypes, ft)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating form types: %w", err)
	}

	return formTypes, nil
}

// UpdateFieldOrders updates the sort order for multiple fields
func (s *FormConfigService) UpdateFieldOrders(entityID int, fieldOrders []struct {
	FieldID   int `json:"field_id"`
	SortOrder int `json:"sort_order"`
}) error {
	if len(fieldOrders) == 0 {
		return fmt.Errorf("no field orders to update")
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update each field's sort order
	updateQuery := `
		INSERT INTO entity_field_configurations 
		(healthcare_entity_id, field_definition_id, is_enabled, is_required, sort_order, updated_at)
		VALUES ($1, $2, COALESCE((SELECT is_enabled FROM entity_field_configurations WHERE healthcare_entity_id = $1 AND field_definition_id = $2), true), 
		        COALESCE((SELECT is_required FROM entity_field_configurations WHERE healthcare_entity_id = $1 AND field_definition_id = $2), 
		                 (SELECT default_required FROM field_definitions WHERE id = $2)), $3, CURRENT_TIMESTAMP)
		ON CONFLICT (healthcare_entity_id, field_definition_id)
		DO UPDATE SET
			sort_order = EXCLUDED.sort_order,
			updated_at = CURRENT_TIMESTAMP
	`

	for _, fieldOrder := range fieldOrders {
		// Verify field exists and is active
		var exists bool
		err := tx.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM field_definitions WHERE id = $1 AND is_active = true)
		`, fieldOrder.FieldID).Scan(&exists)
		
		if err != nil {
			return fmt.Errorf("failed to verify field %d: %w", fieldOrder.FieldID, err)
		}
		
		if !exists {
			return fmt.Errorf("field with ID %d not found or not active", fieldOrder.FieldID)
		}

		// Update the field's sort order
		_, err = tx.Exec(updateQuery, entityID, fieldOrder.FieldID, fieldOrder.SortOrder)
		if err != nil {
			return fmt.Errorf("failed to update sort order for field %d: %w", fieldOrder.FieldID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit field order updates: %w", err)
	}

	return nil
}
