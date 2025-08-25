package main

import (
	"database/sql"
	"strings"

	"github.com/lib/pq"
)

// Country represents a country with multi-locale support
type Country struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`         // ISO 3166-1 alpha-2 code
	Name        string `json:"name"`         // Localized name based on request
	NameEN      string `json:"name_en"`      // English name
	NameFR      string `json:"name_fr,omitempty"` // French name
	NameAR      string `json:"name_ar,omitempty"` // Arabic name
	ISOAlpha3   string `json:"iso_alpha3,omitempty"`
	NumericCode string `json:"numeric_code,omitempty"`
	Region      string `json:"region,omitempty"`
	Subregion   string `json:"subregion,omitempty"`
}

// State represents a state/province with multi-locale support
type State struct {
	ID        int    `json:"id"`
	CountryID int    `json:"country_id"`
	Code      string `json:"code"`
	Name      string `json:"name"`         // Localized name based on request
	NameEN    string `json:"name_en"`      // English name
	NameFR    string `json:"name_fr,omitempty"` // French name
	NameAR    string `json:"name_ar,omitempty"` // Arabic name
	Type      string `json:"type,omitempty"`     // state, province, region, etc.
}

// City represents a city with multi-locale support and geographic data
type City struct {
	ID         int     `json:"id"`
	CountryID  int     `json:"country_id"`
	StateID    *int    `json:"state_id,omitempty"`
	Code       string  `json:"code,omitempty"`
	Name       string  `json:"name"`         // Localized name based on request
	NameEN     string  `json:"name_en"`      // English name
	NameFR     string  `json:"name_fr,omitempty"` // French name
	NameAR     string  `json:"name_ar,omitempty"` // Arabic name
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
	Population *int    `json:"population,omitempty"`
	Timezone   string  `json:"timezone,omitempty"`
	IsCapital  bool    `json:"is_capital,omitempty"`
}

// Nationality represents a nationality with multi-locale support
type Nationality struct {
	ID        int    `json:"id"`
	CountryID int    `json:"country_id"`
	Name      string `json:"name"`         // Localized name based on request
	NameEN    string `json:"name_en"`      // English name
	NameFR    string `json:"name_fr,omitempty"` // French name
	NameAR    string `json:"name_ar,omitempty"` // Arabic name
	Code      string `json:"code,omitempty"`    // Nationality code
	IsPrimary bool   `json:"is_primary"`   // Is this the primary nationality for the country
}

// InsuranceType represents an insurance type with multi-locale support
type InsuranceType struct {
	ID        int    `json:"id"`
	CountryID int    `json:"country_id"`
	Code      string `json:"code"`         // Insurance type code (e.g., "public", "private", "other")
	Name      string `json:"name"`         // Localized name based on request
	NameEN    string `json:"name_en"`      // English name
	NameFR    string `json:"name_fr,omitempty"` // French name
	NameAR    string `json:"name_ar,omitempty"` // Arabic name
	IsDefault bool   `json:"is_default"`   // Is this the default type for the country
	SortOrder int    `json:"sort_order"`   // Display order
}

// InsuranceProvider represents an insurance provider with multi-locale support
type InsuranceProvider struct {
	ID                int    `json:"id"`
	InsuranceTypeID   int    `json:"insurance_type_id"`  // Foreign key to insurance_types
	Code              string `json:"code"`               // Provider code (e.g., "ohip", "medicare", "other")
	Name              string `json:"name"`               // Localized name based on request
	NameEN            string `json:"name_en"`            // English name
	NameFR            string `json:"name_fr,omitempty"`  // French name
	NameAR            string `json:"name_ar,omitempty"`  // Arabic name
	IsDefault         bool   `json:"is_default"`         // Is this the default provider for the type
	SortOrder         int    `json:"sort_order"`         // Display order
}

type LocationStore struct {
	db *sql.DB
}

func NewLocationStore() *LocationStore {
	return &LocationStore{
		db: db, // Use global db connection
	}
}

// GetCountries returns all active countries with locale support
func (s *LocationStore) GetCountries(locale ...string) ([]Country, error) {
	localeCode := "en" // Default to English
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	// Determine which name column to use for the localized name
	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(name_fr, name_en)"
	case "ar":
		nameColumn = "COALESCE(name_ar, name_en)"
	default:
		nameColumn = "name_en"
	}

	query := `
		SELECT id, code, ` + nameColumn + ` as name, 
		       name_en, name_fr, name_ar, iso_alpha3, numeric_code, region, subregion
		FROM countries 
		WHERE is_active = TRUE 
		ORDER BY ` + nameColumn + `
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []Country
	for rows.Next() {
		var country Country
		var nameFR, nameAR, isoAlpha3, numericCode, region, subregion sql.NullString
		
		if err := rows.Scan(
			&country.ID, &country.Code, &country.Name,
			&country.NameEN, &nameFR, &nameAR, 
			&isoAlpha3, &numericCode, &region, &subregion,
		); err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if nameFR.Valid {
			country.NameFR = nameFR.String
		}
		if nameAR.Valid {
			country.NameAR = nameAR.String
		}
		if isoAlpha3.Valid {
			country.ISOAlpha3 = isoAlpha3.String
		}
		if numericCode.Valid {
			country.NumericCode = numericCode.String
		}
		if region.Valid {
			country.Region = region.String
		}
		if subregion.Valid {
			country.Subregion = subregion.String
		}
		
		countries = append(countries, country)
	}
	
	return countries, rows.Err()
}

// GetStatesByCountry returns all states for a country with locale support
// Can accept either country code or country ID
func (s *LocationStore) GetStatesByCountry(countryIdentifier interface{}, locale ...string) ([]State, error) {
	localeCode := "en" // Default to English
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	// Determine which name column to use for the localized name
	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(s.name_fr, s.name_en)"
	case "ar":
		nameColumn = "COALESCE(s.name_ar, s.name_en)"
	default:
		nameColumn = "s.name_en"
	}

	// Build query based on identifier type
	var query string
	var args []interface{}
	
	query = `
		SELECT s.id, s.country_id, s.code, ` + nameColumn + ` as name,
		       s.name_en, s.name_fr, s.name_ar, s.type
		FROM states s
		JOIN countries c ON s.country_id = c.id
		WHERE s.is_active = TRUE AND c.is_active = TRUE
	`
	
	// Check if identifier is string (country code) or int (country ID)
	switch v := countryIdentifier.(type) {
	case string:
		query += ` AND c.code = $1`
		args = append(args, strings.ToUpper(v))
	case int:
		query += ` AND c.id = $1`
		args = append(args, v)
	default:
		// Try to convert to string as fallback
		query += ` AND c.code = $1`
		args = append(args, strings.ToUpper(countryIdentifier.(string)))
	}
	
	query += ` ORDER BY ` + nameColumn
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []State
	for rows.Next() {
		var state State
		var nameFR, nameAR, stateType sql.NullString
		
		if err := rows.Scan(
			&state.ID, &state.CountryID, &state.Code, &state.Name,
			&state.NameEN, &nameFR, &nameAR, &stateType,
		); err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if nameFR.Valid {
			state.NameFR = nameFR.String
		}
		if nameAR.Valid {
			state.NameAR = nameAR.String
		}
		if stateType.Valid {
			state.Type = stateType.String
		}
		
		states = append(states, state)
	}
	
	return states, rows.Err()
}

// GetCitiesByCountry returns cities for a country with optional state filter and search
func (s *LocationStore) GetCitiesByCountry(countryIdentifier interface{}, stateIdentifier interface{}, searchQuery string, limit int, locale ...string) ([]City, error) {
	localeCode := "en" // Default to English
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	// Determine which name column to use for the localized name
	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(c.name_fr, c.name_en)"
	case "ar":
		nameColumn = "COALESCE(c.name_ar, c.name_en)"
	default:
		nameColumn = "c.name_en"
	}

	var querySQL strings.Builder
	var args []interface{}
	argIndex := 1

	// Base query with full city information
	querySQL.WriteString(`
		SELECT c.id, c.country_id, c.state_id, COALESCE(c.code, ''), 
		       ` + nameColumn + ` as name, c.name_en, c.name_fr, c.name_ar,
		       c.latitude, c.longitude, c.population, c.timezone, c.is_capital
		FROM cities c
		JOIN countries co ON c.country_id = co.id
		WHERE c.is_active = TRUE AND co.is_active = TRUE
	`)

	// Add country filter
	switch v := countryIdentifier.(type) {
	case string:
		querySQL.WriteString(` AND co.code = $` + string(rune('0'+argIndex)))
		args = append(args, strings.ToUpper(v))
	case int:
		querySQL.WriteString(` AND co.id = $` + string(rune('0'+argIndex)))
		args = append(args, v)
	default:
		// Try to convert to string as fallback
		querySQL.WriteString(` AND co.code = $` + string(rune('0'+argIndex)))
		args = append(args, strings.ToUpper(countryIdentifier.(string)))
	}
	argIndex++

	// Add state filter if provided
	if stateIdentifier != nil {
		switch v := stateIdentifier.(type) {
		case string:
			if v != "" {
				querySQL.WriteString(` AND c.state_id = (SELECT id FROM states WHERE country_id = co.id AND code = $` + string(rune('0'+argIndex)) + `)`)
				args = append(args, strings.ToUpper(v))
				argIndex++
			}
		case int:
			if v > 0 {
				querySQL.WriteString(` AND c.state_id = $` + string(rune('0'+argIndex)))
				args = append(args, v)
				argIndex++
			}
		}
	}

	// Add search filter if provided
	if searchQuery != "" {
		trimmedQuery := strings.TrimSpace(searchQuery)
		if trimmedQuery != "" {
			// Use full-text search if available, otherwise use LIKE
			querySQL.WriteString(` AND (
				LOWER(c.name_en) LIKE LOWER($` + string(rune('0'+argIndex)) + `) OR
				LOWER(COALESCE(c.name_fr, '')) LIKE LOWER($` + string(rune('0'+argIndex)) + `) OR
				LOWER(COALESCE(c.name_ar, '')) LIKE LOWER($` + string(rune('0'+argIndex)) + `)
			)`)
			args = append(args, "%"+trimmedQuery+"%")
			argIndex++
		}
	}

	// Order by localized name
	querySQL.WriteString(` ORDER BY ` + nameColumn)

	// Add limit if provided
	if limit > 0 {
		querySQL.WriteString(` LIMIT $` + string(rune('0'+argIndex)))
		args = append(args, limit)
	}

	rows, err := s.db.Query(querySQL.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []City
	for rows.Next() {
		var city City
		var stateID sql.NullInt64
		var code, nameFR, nameAR, timezone sql.NullString
		var latitude, longitude sql.NullFloat64
		var population sql.NullInt64
		var isCapital bool
		
		if err := rows.Scan(
			&city.ID, &city.CountryID, &stateID, &code,
			&city.Name, &city.NameEN, &nameFR, &nameAR,
			&latitude, &longitude, &population, &timezone, &isCapital,
		); err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if stateID.Valid {
			stateIDInt := int(stateID.Int64)
			city.StateID = &stateIDInt
		}
		if code.Valid {
			city.Code = code.String
		}
		if nameFR.Valid {
			city.NameFR = nameFR.String
		}
		if nameAR.Valid {
			city.NameAR = nameAR.String
		}
		if latitude.Valid {
			city.Latitude = &latitude.Float64
		}
		if longitude.Valid {
			city.Longitude = &longitude.Float64
		}
		if population.Valid {
			populationInt := int(population.Int64)
			city.Population = &populationInt
		}
		if timezone.Valid {
			city.Timezone = timezone.String
		}
		city.IsCapital = isCapital
		
		cities = append(cities, city)
	}
	
	return cities, rows.Err()
}

// GetCitiesByState returns cities for a specific state with optional search
func (s *LocationStore) GetCitiesByState(stateID int, searchQuery string, limit int, locale ...string) ([]City, error) {
	localeCode := "en" // Default to English
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	// Determine which name column to use for the localized name
	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(c.name_fr, c.name_en)"
	case "ar":
		nameColumn = "COALESCE(c.name_ar, c.name_en)"
	default:
		nameColumn = "c.name_en"
	}

	querySQL := `
		SELECT c.id, c.country_id, c.state_id, c.code,
		       ` + nameColumn + ` as name,
		       c.name_en, c.name_fr, c.name_ar,
		       c.latitude, c.longitude, c.population, c.timezone, c.is_capital
		FROM cities c
		WHERE c.is_active = TRUE AND c.state_id = $1`
	
	args := []interface{}{stateID}
	argIndex := 2

	// Add search filter if provided
	if searchQuery != "" {
		trimmedQuery := strings.TrimSpace(searchQuery)
		if trimmedQuery != "" {
			querySQL += ` AND (
				LOWER(c.name_en) LIKE LOWER($` + string(rune('0'+argIndex)) + `) OR
				LOWER(COALESCE(c.name_fr, '')) LIKE LOWER($` + string(rune('0'+argIndex)) + `) OR
				LOWER(COALESCE(c.name_ar, '')) LIKE LOWER($` + string(rune('0'+argIndex)) + `)
			)`
			args = append(args, "%"+trimmedQuery+"%")
			argIndex++
		}
	}

	// Order by localized name
	querySQL += ` ORDER BY ` + nameColumn

	// Add limit if provided
	if limit > 0 {
		querySQL += ` LIMIT $` + string(rune('0'+argIndex))
		args = append(args, limit)
	}

	rows, err := s.db.Query(querySQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []City
	for rows.Next() {
		var city City
		var stateID sql.NullInt64
		var code, nameFR, nameAR, timezone sql.NullString
		var latitude, longitude sql.NullFloat64
		var population sql.NullInt64
		var isCapital bool
		
		if err := rows.Scan(
			&city.ID, &city.CountryID, &stateID, &code,
			&city.Name, &city.NameEN, &nameFR, &nameAR,
			&latitude, &longitude, &population, &timezone, &isCapital,
		); err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if stateID.Valid {
			stateIDInt := int(stateID.Int64)
			city.StateID = &stateIDInt
		}
		if code.Valid {
			city.Code = code.String
		}
		if nameFR.Valid {
			city.NameFR = nameFR.String
		}
		if nameAR.Valid {
			city.NameAR = nameAR.String
		}
		if latitude.Valid {
			city.Latitude = &latitude.Float64
		}
		if longitude.Valid {
			city.Longitude = &longitude.Float64
		}
		if population.Valid {
			populationInt := int(population.Int64)
			city.Population = &populationInt
		}
		if timezone.Valid {
			city.Timezone = timezone.String
		}
		city.IsCapital = isCapital
		
		cities = append(cities, city)
	}
	
	return cities, rows.Err()
}

// GetCountryByID returns a country by its ID
func (s *LocationStore) GetCountryByID(id int, locale ...string) (*Country, error) {
	countries, err := s.GetCountries(locale...)
	if err != nil {
		return nil, err
	}
	
	for _, country := range countries {
		if country.ID == id {
			return &country, nil
		}
	}
	
	return nil, sql.ErrNoRows
}

// GetStateByID returns a state by its ID
func (s *LocationStore) GetStateByID(id int, locale ...string) (*State, error) {
	localeCode := "en"
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(name_fr, name_en)"
	case "ar":
		nameColumn = "COALESCE(name_ar, name_en)"
	default:
		nameColumn = "name_en"
	}

	query := `
		SELECT id, country_id, code, ` + nameColumn + ` as name,
		       name_en, name_fr, name_ar, type
		FROM states 
		WHERE id = $1 AND is_active = TRUE
	`
	
	var state State
	var nameFR, nameAR, stateType sql.NullString
	
	err := s.db.QueryRow(query, id).Scan(
		&state.ID, &state.CountryID, &state.Code, &state.Name,
		&state.NameEN, &nameFR, &nameAR, &stateType,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Handle nullable fields
	if nameFR.Valid {
		state.NameFR = nameFR.String
	}
	if nameAR.Valid {
		state.NameAR = nameAR.String
	}
	if stateType.Valid {
		state.Type = stateType.String
	}
	
	return &state, nil
}

// GetCityByID returns a city by its ID
func (s *LocationStore) GetCityByID(id int, locale ...string) (*City, error) {
	localeCode := "en"
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(name_fr, name_en)"
	case "ar":
		nameColumn = "COALESCE(name_ar, name_en)"
	default:
		nameColumn = "name_en"
	}

	query := `
		SELECT id, country_id, state_id, COALESCE(code, ''),
		       ` + nameColumn + ` as name, name_en, name_fr, name_ar,
		       latitude, longitude, population, timezone, is_capital
		FROM cities 
		WHERE id = $1 AND is_active = TRUE
	`
	
	var city City
	var stateID sql.NullInt64
	var code, nameFR, nameAR, timezone sql.NullString
	var latitude, longitude sql.NullFloat64
	var population sql.NullInt64
	var isCapital bool
	
	err := s.db.QueryRow(query, id).Scan(
		&city.ID, &city.CountryID, &stateID, &code,
		&city.Name, &city.NameEN, &nameFR, &nameAR,
		&latitude, &longitude, &population, &timezone, &isCapital,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Handle nullable fields
	if stateID.Valid {
		stateIDInt := int(stateID.Int64)
		city.StateID = &stateIDInt
	}
	if code.Valid {
		city.Code = code.String
	}
	if nameFR.Valid {
		city.NameFR = nameFR.String
	}
	if nameAR.Valid {
		city.NameAR = nameAR.String
	}
	if latitude.Valid {
		city.Latitude = &latitude.Float64
	}
	if longitude.Valid {
		city.Longitude = &longitude.Float64
	}
	if population.Valid {
		populationInt := int(population.Int64)
		city.Population = &populationInt
	}
	if timezone.Valid {
		city.Timezone = timezone.String
	}
	city.IsCapital = isCapital
	
	return &city, nil
}

// GetCountriesByIDs retrieves multiple countries by their IDs with locale support
func (s *LocationStore) GetCountriesByIDs(ids []int, locale string) ([]Country, error) {
	if len(ids) == 0 {
		return []Country{}, nil
	}

	query := `
		SELECT 
			id, code, name_en, name_fr, name_ar, iso_alpha3, numeric_code, region, subregion
		FROM countries 
		WHERE id = ANY($1)
		ORDER BY name_en
	`

	// Convert to PostgreSQL array format
	rows, err := s.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	countries := make([]Country, 0)
	for rows.Next() {
		var country Country
		var nameFR, nameAR, isoAlpha3, numericCode, region, subregion sql.NullString

		err := rows.Scan(
			&country.ID,
			&country.Code,
			&country.NameEN,
			&nameFR,
			&nameAR,
			&isoAlpha3,
			&numericCode,
			&region,
			&subregion,
		)
		if err != nil {
			return nil, err
		}

		// Set localized name based on requested locale
		switch locale {
		case "fr":
			if nameFR.Valid && nameFR.String != "" {
				country.Name = nameFR.String
			} else {
				country.Name = country.NameEN
			}
		case "ar":
			if nameAR.Valid && nameAR.String != "" {
				country.Name = nameAR.String
			} else {
				country.Name = country.NameEN
			}
		default:
			country.Name = country.NameEN
		}

		// Set optional fields
		if nameFR.Valid {
			country.NameFR = nameFR.String
		}
		if nameAR.Valid {
			country.NameAR = nameAR.String
		}
		if isoAlpha3.Valid {
			country.ISOAlpha3 = isoAlpha3.String
		}
		if numericCode.Valid {
			country.NumericCode = numericCode.String
		}
		if region.Valid {
			country.Region = region.String
		}
		if subregion.Valid {
			country.Subregion = subregion.String
		}

		countries = append(countries, country)
	}

	return countries, nil
}

// GetStatesByIDs retrieves multiple states by their IDs with locale support
func (s *LocationStore) GetStatesByIDs(ids []int, locale string) ([]State, error) {
	if len(ids) == 0 {
		return []State{}, nil
	}

	query := `
		SELECT 
			id, country_id, code, name_en, name_fr, name_ar, type
		FROM states 
		WHERE id = ANY($1)
		ORDER BY name_en
	`

	rows, err := s.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	states := make([]State, 0)
	for rows.Next() {
		var state State
		var nameFR, nameAR, stateType sql.NullString

		err := rows.Scan(
			&state.ID,
			&state.CountryID,
			&state.Code,
			&state.NameEN,
			&nameFR,
			&nameAR,
			&stateType,
		)
		if err != nil {
			return nil, err
		}

		// Set localized name based on requested locale
		switch locale {
		case "fr":
			if nameFR.Valid && nameFR.String != "" {
				state.Name = nameFR.String
			} else {
				state.Name = state.NameEN
			}
		case "ar":
			if nameAR.Valid && nameAR.String != "" {
				state.Name = nameAR.String
			} else {
				state.Name = state.NameEN
			}
		default:
			state.Name = state.NameEN
		}

		// Set optional fields
		if nameFR.Valid {
			state.NameFR = nameFR.String
		}
		if nameAR.Valid {
			state.NameAR = nameAR.String
		}
		if stateType.Valid {
			state.Type = stateType.String
		}

		states = append(states, state)
	}

	return states, nil
}

// GetCitiesByIDs retrieves multiple cities by their IDs with locale support
func (s *LocationStore) GetCitiesByIDs(ids []int, locale string) ([]City, error) {
	if len(ids) == 0 {
		return []City{}, nil
	}

	query := `
		SELECT 
			id, country_id, state_id, code, name_en, name_fr, name_ar, 
			latitude, longitude, population, timezone, is_capital
		FROM cities 
		WHERE id = ANY($1)
		ORDER BY name_en
	`

	rows, err := s.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cities := make([]City, 0)
	for rows.Next() {
		var city City
		var stateID sql.NullInt64
		var code, nameFR, nameAR, timezone sql.NullString
		var latitude, longitude sql.NullFloat64
		var population sql.NullInt64
		var isCapital bool

		err := rows.Scan(
			&city.ID,
			&city.CountryID,
			&stateID,
			&code,
			&city.NameEN,
			&nameFR,
			&nameAR,
			&latitude,
			&longitude,
			&population,
			&timezone,
			&isCapital,
		)
		if err != nil {
			return nil, err
		}

		// Set localized name based on requested locale
		switch locale {
		case "fr":
			if nameFR.Valid && nameFR.String != "" {
				city.Name = nameFR.String
			} else {
				city.Name = city.NameEN
			}
		case "ar":
			if nameAR.Valid && nameAR.String != "" {
				city.Name = nameAR.String
			} else {
				city.Name = city.NameEN
			}
		default:
			city.Name = city.NameEN
		}

		// Set optional fields
		if stateID.Valid {
			stateIDInt := int(stateID.Int64)
			city.StateID = &stateIDInt
		}
		if code.Valid {
			city.Code = code.String
		}
		if nameFR.Valid {
			city.NameFR = nameFR.String
		}
		if nameAR.Valid {
			city.NameAR = nameAR.String
		}
		if latitude.Valid {
			city.Latitude = &latitude.Float64
		}
		if longitude.Valid {
			city.Longitude = &longitude.Float64
		}
		if population.Valid {
			populationInt := int(population.Int64)
			city.Population = &populationInt
		}
		if timezone.Valid {
			city.Timezone = timezone.String
		}
		city.IsCapital = isCapital

		cities = append(cities, city)
	}

	return cities, nil
}

// GetNationalities returns all active nationalities with locale support
func (s *LocationStore) GetNationalities(locale ...string) ([]Nationality, error) {
	localeCode := "en" // Default to English
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	// Determine which name column to use for the localized name
	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(name_fr, name_en)"
	case "ar":
		nameColumn = "COALESCE(name_ar, name_en)"
	default:
		nameColumn = "name_en"
	}

	query := `
		SELECT id, country_id, ` + nameColumn + ` as name,
		       name_en, name_fr, name_ar, code, is_primary
		FROM nationalities
		WHERE is_active = TRUE
		ORDER BY ` + nameColumn + `
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nationalities []Nationality
	for rows.Next() {
		var nationality Nationality
		var nameFR, nameAR, code sql.NullString
		
		if err := rows.Scan(
			&nationality.ID, &nationality.CountryID, &nationality.Name,
			&nationality.NameEN, &nameFR, &nameAR, &code, &nationality.IsPrimary,
		); err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if nameFR.Valid {
			nationality.NameFR = nameFR.String
		}
		if nameAR.Valid {
			nationality.NameAR = nameAR.String
		}
		if code.Valid {
			nationality.Code = code.String
		}
		
		nationalities = append(nationalities, nationality)
	}
	
	return nationalities, rows.Err()
}

// GetNationalitiesByCountry returns nationalities for a specific country
func (s *LocationStore) GetNationalitiesByCountry(countryIdentifier interface{}, locale ...string) ([]Nationality, error) {
	localeCode := "en" // Default to English
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	// Determine which name column to use for the localized name
	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(n.name_fr, n.name_en)"
	case "ar":
		nameColumn = "COALESCE(n.name_ar, n.name_en)"
	default:
		nameColumn = "n.name_en"
	}

	// Build query based on identifier type
	var query string
	var args []interface{}
	
	query = `
		SELECT n.id, n.country_id, ` + nameColumn + ` as name,
		       n.name_en, n.name_fr, n.name_ar, n.code, n.is_primary
		FROM nationalities n
		JOIN countries c ON n.country_id = c.id
		WHERE n.is_active = TRUE AND c.is_active = TRUE
	`
	
	// Check if identifier is string (country code) or int (country ID)
	switch v := countryIdentifier.(type) {
	case string:
		query += ` AND c.code = $1`
		args = append(args, strings.ToUpper(v))
	case int:
		query += ` AND c.id = $1`
		args = append(args, v)
	default:
		// Try to convert to string as fallback
		query += ` AND c.code = $1`
		args = append(args, strings.ToUpper(countryIdentifier.(string)))
	}
	
	query += ` ORDER BY n.is_primary DESC, ` + nameColumn
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nationalities []Nationality
	for rows.Next() {
		var nationality Nationality
		var nameFR, nameAR, code sql.NullString
		
		if err := rows.Scan(
			&nationality.ID, &nationality.CountryID, &nationality.Name,
			&nationality.NameEN, &nameFR, &nameAR, &code, &nationality.IsPrimary,
		); err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if nameFR.Valid {
			nationality.NameFR = nameFR.String
		}
		if nameAR.Valid {
			nationality.NameAR = nameAR.String
		}
		if code.Valid {
			nationality.Code = code.String
		}
		
		nationalities = append(nationalities, nationality)
	}
	
	return nationalities, rows.Err()
}

// GetNationalitiesByIDs retrieves multiple nationalities by their IDs with locale support
func (s *LocationStore) GetNationalitiesByIDs(ids []int, locale string) ([]Nationality, error) {
	if len(ids) == 0 {
		return []Nationality{}, nil
	}

	query := `
		SELECT 
			id, country_id, name_en, name_fr, name_ar, code, is_primary
		FROM nationalities 
		WHERE id = ANY($1)
		ORDER BY name_en
	`

	// Convert to PostgreSQL array format
	rows, err := s.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nationalities := make([]Nationality, 0)
	for rows.Next() {
		var nationality Nationality
		var nameFR, nameAR, code sql.NullString

		err := rows.Scan(
			&nationality.ID,
			&nationality.CountryID,
			&nationality.NameEN,
			&nameFR,
			&nameAR,
			&code,
			&nationality.IsPrimary,
		)
		if err != nil {
			return nil, err
		}

		// Set localized name based on requested locale
		switch locale {
		case "fr":
			if nameFR.Valid && nameFR.String != "" {
				nationality.Name = nameFR.String
			} else {
				nationality.Name = nationality.NameEN
			}
		case "ar":
			if nameAR.Valid && nameAR.String != "" {
				nationality.Name = nameAR.String
			} else {
				nationality.Name = nationality.NameEN
			}
		default:
			nationality.Name = nationality.NameEN
		}

		// Set optional fields
		if nameFR.Valid {
			nationality.NameFR = nameFR.String
		}
		if nameAR.Valid {
			nationality.NameAR = nameAR.String
		}
		if code.Valid {
			nationality.Code = code.String
		}

		nationalities = append(nationalities, nationality)
	}

	return nationalities, nil
}

// GetNationalityByID returns a nationality by its ID
func (s *LocationStore) GetNationalityByID(id int, locale ...string) (*Nationality, error) {
	localeCode := "en"
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(name_fr, name_en)"
	case "ar":
		nameColumn = "COALESCE(name_ar, name_en)"
	default:
		nameColumn = "name_en"
	}

	query := `
		SELECT id, country_id, ` + nameColumn + ` as name,
		       name_en, name_fr, name_ar, code, is_primary
		FROM nationalities 
		WHERE id = $1 AND is_active = TRUE
	`
	
	var nationality Nationality
	var nameFR, nameAR, code sql.NullString
	
	err := s.db.QueryRow(query, id).Scan(
		&nationality.ID, &nationality.CountryID, &nationality.Name,
		&nationality.NameEN, &nameFR, &nameAR, &code, &nationality.IsPrimary,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Handle nullable fields
	if nameFR.Valid {
		nationality.NameFR = nameFR.String
	}
	if nameAR.Valid {
		nationality.NameAR = nameAR.String
	}
	if code.Valid {
		nationality.Code = code.String
	}
	
	return &nationality, nil
}

// GetInsuranceTypesByCountry returns insurance types for a specific country with locale support
func (s *LocationStore) GetInsuranceTypesByCountry(countryIdentifier interface{}, locale ...string) ([]InsuranceType, error) {
	localeCode := "en" // Default to English
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	// Determine which name column to use for the localized name
	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(it.name_fr, it.name_en)"
	case "ar":
		nameColumn = "COALESCE(it.name_ar, it.name_en)"
	default:
		nameColumn = "it.name_en"
	}

	// Build query based on identifier type
	var query string
	var args []interface{}
	
	query = `
		SELECT it.id, it.country_id, it.code, ` + nameColumn + ` as name,
		       it.name_en, it.name_fr, it.name_ar, it.is_default, it.sort_order
		FROM insurance_types it
		JOIN countries c ON it.country_id = c.id
		WHERE it.is_active = TRUE AND c.is_active = TRUE
	`
	
	// Check if identifier is string (country code) or int (country ID)
	switch v := countryIdentifier.(type) {
	case string:
		query += ` AND c.code = $1`
		args = append(args, strings.ToUpper(v))
	case int:
		query += ` AND c.id = $1`
		args = append(args, v)
	default:
		// Try to convert to string as fallback
		query += ` AND c.code = $1`
		args = append(args, strings.ToUpper(countryIdentifier.(string)))
	}
	
	query += ` ORDER BY it.sort_order, ` + nameColumn
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var insuranceTypes []InsuranceType
	for rows.Next() {
		var insuranceType InsuranceType
		var nameFR, nameAR sql.NullString
		
		if err := rows.Scan(
			&insuranceType.ID, &insuranceType.CountryID, &insuranceType.Code, &insuranceType.Name,
			&insuranceType.NameEN, &nameFR, &nameAR, &insuranceType.IsDefault, &insuranceType.SortOrder,
		); err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if nameFR.Valid {
			insuranceType.NameFR = nameFR.String
		}
		if nameAR.Valid {
			insuranceType.NameAR = nameAR.String
		}
		
		insuranceTypes = append(insuranceTypes, insuranceType)
	}
	
	return insuranceTypes, rows.Err()
}

// GetInsuranceProvidersByType returns insurance providers for a specific insurance type with locale support
func (s *LocationStore) GetInsuranceProvidersByType(insuranceTypeID int, locale ...string) ([]InsuranceProvider, error) {
	localeCode := "en" // Default to English
	if len(locale) > 0 && locale[0] != "" {
		localeCode = strings.ToLower(locale[0])
	}

	// Determine which name column to use for the localized name
	var nameColumn string
	switch localeCode {
	case "fr":
		nameColumn = "COALESCE(ip.name_fr, ip.name_en)"
	case "ar":
		nameColumn = "COALESCE(ip.name_ar, ip.name_en)"
	default:
		nameColumn = "ip.name_en"
	}

	query := `
		SELECT ip.id, ip.insurance_type_id, ip.code, ` + nameColumn + ` as name,
		       ip.name_en, ip.name_fr, ip.name_ar, ip.is_default, ip.sort_order
		FROM insurance_providers ip
		WHERE ip.is_active = TRUE AND ip.insurance_type_id = $1
		ORDER BY ip.sort_order, ` + nameColumn
	
	rows, err := s.db.Query(query, insuranceTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var insuranceProviders []InsuranceProvider
	for rows.Next() {
		var insuranceProvider InsuranceProvider
		var nameFR, nameAR sql.NullString
		
		if err := rows.Scan(
			&insuranceProvider.ID, &insuranceProvider.InsuranceTypeID, &insuranceProvider.Code, &insuranceProvider.Name,
			&insuranceProvider.NameEN, &nameFR, &nameAR, &insuranceProvider.IsDefault, &insuranceProvider.SortOrder,
		); err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if nameFR.Valid {
			insuranceProvider.NameFR = nameFR.String
		}
		if nameAR.Valid {
			insuranceProvider.NameAR = nameAR.String
		}
		
		insuranceProviders = append(insuranceProviders, insuranceProvider)
	}
	
	return insuranceProviders, rows.Err()
}

// GetInsuranceTypesByIDs retrieves multiple insurance types by their IDs with locale support
func (s *LocationStore) GetInsuranceTypesByIDs(ids []int, locale string) ([]InsuranceType, error) {
	if len(ids) == 0 {
		return []InsuranceType{}, nil
	}

	query := `
		SELECT id, country_id, code, name_en, name_fr, name_ar, is_default, sort_order
		FROM insurance_types
		WHERE id = ANY($1)
		ORDER BY sort_order, name_en
	`

	rows, err := s.db.Query(query, pq.Array(ids))
	if err != nil { return nil, err }
	defer rows.Close()

	result := make([]InsuranceType, 0)
	for rows.Next() {
		var it InsuranceType
		var nameFR, nameAR sql.NullString
		if err := rows.Scan(&it.ID, &it.CountryID, &it.Code, &it.NameEN, &nameFR, &nameAR, &it.IsDefault, &it.SortOrder); err != nil {
			return nil, err
		}
		switch locale {
		case "fr":
			if nameFR.Valid && nameFR.String != "" { it.Name = nameFR.String } else { it.Name = it.NameEN }
		case "ar":
			if nameAR.Valid && nameAR.String != "" { it.Name = nameAR.String } else { it.Name = it.NameEN }
		default:
			it.Name = it.NameEN
		}
		if nameFR.Valid { it.NameFR = nameFR.String }
		if nameAR.Valid { it.NameAR = nameAR.String }
		result = append(result, it)
	}
	return result, nil
}

// GetInsuranceProvidersByIDs retrieves multiple insurance providers by their IDs with locale support
func (s *LocationStore) GetInsuranceProvidersByIDs(ids []int, locale string) ([]InsuranceProvider, error) {
	if len(ids) == 0 {
		return []InsuranceProvider{}, nil
	}

	query := `
		SELECT id, insurance_type_id, code, name_en, name_fr, name_ar, is_default, sort_order
		FROM insurance_providers
		WHERE id = ANY($1)
		ORDER BY sort_order, name_en
	`

	rows, err := s.db.Query(query, pq.Array(ids))
	if err != nil { return nil, err }
	defer rows.Close()

	result := make([]InsuranceProvider, 0)
	for rows.Next() {
		var ip InsuranceProvider
		var nameFR, nameAR sql.NullString
		if err := rows.Scan(&ip.ID, &ip.InsuranceTypeID, &ip.Code, &ip.NameEN, &nameFR, &nameAR, &ip.IsDefault, &ip.SortOrder); err != nil {
			return nil, err
		}
		switch locale {
		case "fr":
			if nameFR.Valid && nameFR.String != "" { ip.Name = nameFR.String } else { ip.Name = ip.NameEN }
		case "ar":
			if nameAR.Valid && nameAR.String != "" { ip.Name = nameAR.String } else { ip.Name = ip.NameEN }
		default:
			ip.Name = ip.NameEN
		}
		if nameFR.Valid { ip.NameFR = nameFR.String }
		if nameAR.Valid { ip.NameAR = nameAR.String }
		result = append(result, ip)
	}
	return result, nil
}
