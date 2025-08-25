package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type TranslationService struct {
	db *sql.DB
}

func NewTranslationService(db *sql.DB) *TranslationService {
	return &TranslationService{db: db}
}

// GetSupportedLocales returns all active locales
func (s *TranslationService) GetSupportedLocales() ([]Locale, error) {
	query := `
		SELECT code, language_name, native_name, country_code, is_active, created_at
		FROM locales 
		WHERE is_active = true
		ORDER BY code
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query locales: %w", err)
	}
	defer rows.Close()
	
	var locales []Locale
	for rows.Next() {
		var locale Locale
		err := rows.Scan(
			&locale.Code,
			&locale.LanguageName,
			&locale.NativeName,
			&locale.CountryCode,
			&locale.IsActive,
			&locale.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan locale: %w", err)
		}
		locales = append(locales, locale)
	}
	
	return locales, rows.Err()
}

// GetTranslation gets a specific translation by key and locale
func (s *TranslationService) GetTranslation(key, locale string) (*Translation, error) {
	query := `
		SELECT id, translation_key, locale, content, context, created_at, updated_at
		FROM translations 
		WHERE translation_key = $1 AND locale = $2
	`
	
	var translation Translation
	err := s.db.QueryRow(query, key, locale).Scan(
		&translation.ID,
		&translation.TranslationKey,
		&translation.Locale,
		&translation.Content,
		&translation.Context,
		&translation.CreatedAt,
		&translation.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("translation not found for key '%s' and locale '%s'", key, locale)
		}
		return nil, fmt.Errorf("failed to get translation: %w", err)
	}
	
	return &translation, nil
}

// GetTranslations gets all translations for a locale with optional context filter
func (s *TranslationService) GetTranslations(locale string, context string) (map[string]string, error) {
	query := `
		SELECT translation_key, content
		FROM translations 
		WHERE locale = $1
	`
	args := []interface{}{locale}
	
	if context != "" {
		query += ` AND context = $2`
		args = append(args, context)
	}
	
	query += ` ORDER BY translation_key`
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query translations: %w", err)
	}
	defer rows.Close()
	
	translations := make(map[string]string)
	for rows.Next() {
		var key, content string
		if err := rows.Scan(&key, &content); err != nil {
			return nil, fmt.Errorf("failed to scan translation: %w", err)
		}
		translations[key] = content
	}
	
	return translations, rows.Err()
}

// CreateOrUpdateTranslation creates or updates a translation
func (s *TranslationService) CreateOrUpdateTranslation(req TranslationRequest) (*Translation, error) {
	query := `
		INSERT INTO translations (translation_key, locale, content, context, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (translation_key, locale)
		DO UPDATE SET
			content = EXCLUDED.content,
			context = EXCLUDED.context,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, translation_key, locale, content, context, created_at, updated_at
	`
	
	var translation Translation
	err := s.db.QueryRow(query, req.TranslationKey, req.Locale, req.Content, req.Context).Scan(
		&translation.ID,
		&translation.TranslationKey,
		&translation.Locale,
		&translation.Content,
		&translation.Context,
		&translation.CreatedAt,
		&translation.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create/update translation: %w", err)
	}
	
	return &translation, nil
}

// GetFieldTranslations gets localized field metadata for a specific locale
func (s *TranslationService) GetFieldTranslations(locale string) (map[int]FieldTranslation, error) {
	query := `
		SELECT id, field_definition_id, locale, display_name, description, placeholder_text, created_at, updated_at
		FROM field_translations 
		WHERE locale = $1
		ORDER BY field_definition_id
	`
	
	rows, err := s.db.Query(query, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to query field translations: %w", err)
	}
	defer rows.Close()
	
	translations := make(map[int]FieldTranslation)
	for rows.Next() {
		var ft FieldTranslation
		err := rows.Scan(
			&ft.ID,
			&ft.FieldDefinitionID,
			&ft.Locale,
			&ft.DisplayName,
			&ft.Description,
			&ft.PlaceholderText,
			&ft.CreatedAt,
			&ft.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan field translation: %w", err)
		}
		translations[ft.FieldDefinitionID] = ft
	}
	
	return translations, rows.Err()
}

// CreateOrUpdateFieldTranslation creates or updates field translation
func (s *TranslationService) CreateOrUpdateFieldTranslation(req FieldTranslationRequest) (*FieldTranslation, error) {
	query := `
		INSERT INTO field_translations (field_definition_id, locale, display_name, description, placeholder_text, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (field_definition_id, locale)
		DO UPDATE SET
			display_name = EXCLUDED.display_name,
			description = EXCLUDED.description,
			placeholder_text = EXCLUDED.placeholder_text,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, field_definition_id, locale, display_name, description, placeholder_text, created_at, updated_at
	`
	
	var ft FieldTranslation
	err := s.db.QueryRow(query, req.FieldDefinitionID, req.Locale, req.DisplayName, req.Description, req.PlaceholderText).Scan(
		&ft.ID,
		&ft.FieldDefinitionID,
		&ft.Locale,
		&ft.DisplayName,
		&ft.Description,
		&ft.PlaceholderText,
		&ft.CreatedAt,
		&ft.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create/update field translation: %w", err)
	}
	
	return &ft, nil
}

// GetLocalizedFormMetadata returns form metadata with translations for the specified locale
func (s *TranslationService) GetLocalizedFormMetadata(formType string, entityID int, locale string) (*LocalizedFormMetadataResponse, error) {
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

	// Get field configurations with translations
	query := `
		SELECT 
			fd.id,
			fd.name,
			COALESCE(ft.display_name, fd.display_name) as display_name,
			fd.field_type,
			fd.default_required,
			fd.is_core,
			fd.validation_rules,
			fd.options,
			fd.sort_order,
			fd.category,
			COALESCE(ft.description, fd.description) as description,
			COALESCE(ft.placeholder_text, fd.placeholder_text) as placeholder_text,
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
		LEFT JOIN field_translations ft
			ON fd.id = ft.field_definition_id
			AND ft.locale = $3
		WHERE fd.form_type_id = $1 
			AND fd.is_active = true
		ORDER BY effective_sort_order, fd.id
	`

	rows, err := s.db.Query(query, formInfo.ID, entityID, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get localized field configurations: %w", err)
	}
	defer rows.Close()

	var fields []LocalizedFieldConfigurationResponse
	var lastModified time.Time

	for rows.Next() {
		var field LocalizedFieldConfigurationResponse
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
			return nil, fmt.Errorf("failed to scan localized field configuration: %w", err)
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
		if field.Name == "city" && field.FieldType != "select" {
			field.FieldType = "select"
		}
		if field.Name == "country" && field.FieldType != "select" {
			field.FieldType = "select"
		}

		field.Locale = locale

		// Track latest modification time
		if rowLastModified.After(lastModified) {
			lastModified = rowLastModified
		}

		fields = append(fields, field)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating localized field configurations: %w", err)
	}

	return &LocalizedFormMetadataResponse{
		FormType:     formType,
		DisplayName:  formInfo.DisplayName,
		Description:  formInfo.Description,
		Fields:       fields,
		EntityID:     entityID,
		Locale:       locale,
		LastModified: lastModified,
	}, nil
}

// ValidateLocale checks if the provided locale is supported
func (s *TranslationService) ValidateLocale(locale string) error {
	var exists bool
	err := s.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM locales WHERE code = $1 AND is_active = true)
	`, locale).Scan(&exists)
	
	if err != nil {
		return fmt.Errorf("failed to validate locale: %w", err)
	}
	
	if !exists {
		return fmt.Errorf("unsupported locale: %s", locale)
	}
	
	return nil
}

// GetDefaultLocaleForCountry returns the default locale for a country
func (s *TranslationService) GetDefaultLocaleForCountry(country string) (string, error) {
	// Map countries to default locales
	countryLocaleMap := map[string]string{
		"Canada":  "en-CA",
		"USA":     "en-US",
		"Morocco": "ar-MA",
		"France":  "fr-FR",
	}
	
	if locale, exists := countryLocaleMap[country]; exists {
		return locale, nil
	}
	
	// Default to en-US if country not found
	return "en-US", nil
}
