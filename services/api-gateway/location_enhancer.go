package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// LocationEnhancer handles location data enhancement
type LocationEnhancer struct {
	client            *resty.Client
	locationServiceURL string
}

// Location represents a location from the location service
type Location struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	LocalizedNames map[string]string  `json:"localized_names,omitempty"`
}

// LocationResponse represents location service API response
type LocationResponse struct {
	Data []Location `json:"data"`
}

// NewLocationEnhancer creates a new location enhancer
func NewLocationEnhancer(locationServiceURL string) *LocationEnhancer {
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)
	client.SetRetryWaitTime(500 * time.Millisecond)

	return &LocationEnhancer{
		client:            client,
		locationServiceURL: locationServiceURL,
	}
}

// LocationEnhancementMiddleware creates middleware for location data enhancement
func (le *LocationEnhancer) LocationEnhancementMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if this is a route that needs location enhancement
		if !le.shouldEnhanceResponse(c.Request.URL.Path, c.Request.Method) {
			c.Next()
			return
		}

		// Create a custom ResponseWriter to capture response
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:          &bytes.Buffer{},
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Only enhance successful responses with JSON content
		if writer.status < 400 && le.isJSONResponse(writer.Header()) {
			enhanced, err := le.enhanceLocationData(writer.body.Bytes(), c.GetString("user_locale"))
			if err != nil {
				log.Printf("Failed to enhance location data: %v", err)
				// Write original response if enhancement fails
				writer.ResponseWriter.Write(writer.body.Bytes())
				return
			}

			// Write enhanced response
			writer.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(enhanced)))
			writer.ResponseWriter.Write(enhanced)
		} else {
			// Write original response
			writer.ResponseWriter.Write(writer.body.Bytes())
		}
	}
}

// responseWriter captures response data for processing
type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return len(data), nil
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// shouldEnhanceResponse determines if a response should be enhanced with location data
func (le *LocationEnhancer) shouldEnhanceResponse(path, method string) bool {
	// Only enhance GET requests for now (read operations)
	if method != http.MethodGet {
		return false
	}

	// Enhance responses from services that return location data
	enhancePaths := []string{
		"/api/users/healthcare-entities",  // Healthcare entity endpoints
		"/api/patients",                   // Patient endpoints
		"/api/internal/healthcare-entities", // Internal healthcare entity endpoints
	}

	for _, enhancePath := range enhancePaths {
		if strings.HasPrefix(path, enhancePath) {
			return true
		}
	}

	return false
}

// isJSONResponse checks if response is JSON
func (le *LocationEnhancer) isJSONResponse(headers http.Header) bool {
	contentType := headers.Get("Content-Type")
	return strings.Contains(contentType, "application/json")
}

// enhanceLocationData enhances JSON response with location names
func (le *LocationEnhancer) enhanceLocationData(responseData []byte, locale string) ([]byte, error) {
	// Parse JSON response
	var response map[string]interface{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return responseData, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Extract location IDs from response
	locationIDs := le.extractLocationIDs(response)
	if len(locationIDs) == 0 {
		// No location IDs found, return original response
		return responseData, nil
	}

	// Fetch location data from location service
	locationMap, err := le.fetchLocationData(locationIDs, locale)
	if err != nil {
		return responseData, fmt.Errorf("failed to fetch location data: %w", err)
	}

	// Enhance response with location names
	le.enhanceResponseWithLocations(response, locationMap)

	// Marshal enhanced response
	enhancedData, err := json.Marshal(response)
	if err != nil {
		return responseData, fmt.Errorf("failed to marshal enhanced response: %w", err)
	}

	return enhancedData, nil
}

// extractLocationIDs extracts all location IDs from a JSON response
func (le *LocationEnhancer) extractLocationIDs(data interface{}) map[string][]int {
	locationIDs := map[string][]int{
		"countries":     {},
		"states":        {},
		"cities":        {},
		"nationalities": {},
		"insurance_types": {},
		"insurance_providers": {},
	}

	le.extractLocationIDsRecursive(data, locationIDs)

	// Remove duplicates
	for key, ids := range locationIDs {
		locationIDs[key] = le.removeDuplicateIDs(ids)
	}

	return locationIDs
}

// extractLocationIDsRecursive recursively extracts location IDs from nested data
func (le *LocationEnhancer) extractLocationIDsRecursive(data interface{}, locationIDs map[string][]int) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			switch key {
			case "country_id":
				if id, ok := value.(float64); ok && id > 0 {
					locationIDs["countries"] = append(locationIDs["countries"], int(id))
				}
			case "state_id":
				if id, ok := value.(float64); ok && id > 0 {
					locationIDs["states"] = append(locationIDs["states"], int(id))
				}
			case "city_id":
				if id, ok := value.(float64); ok && id > 0 {
					locationIDs["cities"] = append(locationIDs["cities"], int(id))
				}
			case "nationality_id":
				if id, ok := value.(float64); ok && id > 0 {
					locationIDs["nationalities"] = append(locationIDs["nationalities"], int(id))
				}
			case "insurance_type_id":
				if id, ok := value.(float64); ok && id > 0 {
					locationIDs["insurance_types"] = append(locationIDs["insurance_types"], int(id))
				}
			case "insurance_provider_id":
				if id, ok := value.(float64); ok && id > 0 {
					locationIDs["insurance_providers"] = append(locationIDs["insurance_providers"], int(id))
				}
			default:
				le.extractLocationIDsRecursive(value, locationIDs)
			}
		}
	case []interface{}:
		for _, item := range v {
			le.extractLocationIDsRecursive(item, locationIDs)
		}
	}
}

// removeDuplicateIDs removes duplicate IDs from slice
func (le *LocationEnhancer) removeDuplicateIDs(ids []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}

	return result
}

// fetchLocationData fetches location data from location service
func (le *LocationEnhancer) fetchLocationData(locationIDs map[string][]int, locale string) (map[string]map[int]Location, error) {
	locationMap := map[string]map[int]Location{
		"countries":     make(map[int]Location),
		"states":        make(map[int]Location),
		"cities":        make(map[int]Location),
		"nationalities": make(map[int]Location),
		"insurance_types": make(map[int]Location),
		"insurance_providers": make(map[int]Location),
	}

	// Set default locale if not provided
	if locale == "" {
		locale = "en"
	}

	// Fetch countries
	if len(locationIDs["countries"]) > 0 {
		countries, err := le.fetchLocationsByType("countries", locationIDs["countries"], locale)
		if err != nil {
			log.Printf("Failed to fetch countries: %v", err)
		} else {
			locationMap["countries"] = countries
		}
	}

	// Fetch states
	if len(locationIDs["states"]) > 0 {
		states, err := le.fetchLocationsByType("states", locationIDs["states"], locale)
		if err != nil {
			log.Printf("Failed to fetch states: %v", err)
		} else {
			locationMap["states"] = states
		}
	}

	// Fetch cities
	if len(locationIDs["cities"]) > 0 {
		cities, err := le.fetchLocationsByType("cities", locationIDs["cities"], locale)
		if err != nil {
			log.Printf("Failed to fetch cities: %v", err)
		} else {
			locationMap["cities"] = cities
		}
	}

	// Fetch nationalities
	if len(locationIDs["nationalities"]) > 0 {
		nationalities, err := le.fetchLocationsByType("nationalities", locationIDs["nationalities"], locale)
		if err != nil {
			log.Printf("Failed to fetch nationalities: %v", err)
		} else {
			locationMap["nationalities"] = nationalities
		}
	}

	// Fetch insurance types
	if len(locationIDs["insurance_types"]) > 0 {
		insuranceTypes, err := le.fetchLocationsByType("insurance-types", locationIDs["insurance_types"], locale)
		if err != nil {
			log.Printf("Failed to fetch insurance types: %v", err)
		} else {
			locationMap["insurance_types"] = insuranceTypes
		}
	}

	// Fetch insurance providers
	if len(locationIDs["insurance_providers"]) > 0 {
		insuranceProviders, err := le.fetchLocationsByType("insurance-providers", locationIDs["insurance_providers"], locale)
		if err != nil {
			log.Printf("Failed to fetch insurance providers: %v", err)
		} else {
			locationMap["insurance_providers"] = insuranceProviders
		}
	}

	return locationMap, nil
}

// fetchLocationsByType fetches locations of a specific type by IDs
func (le *LocationEnhancer) fetchLocationsByType(locationType string, ids []int, locale string) (map[int]Location, error) {
	if len(ids) == 0 {
		return make(map[int]Location), nil
	}

	// Convert IDs to string array for query parameter
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = strconv.Itoa(id)
	}

	url := fmt.Sprintf("%s/api/locations/%s/by-ids", le.locationServiceURL, locationType)
	
	resp, err := le.client.R().
		SetQueryParam("ids", strings.Join(idStrings, ",")).
		SetQueryParam("locale", locale).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", locationType, err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("location service returned status %d for %s", resp.StatusCode(), locationType)
	}

	var locationResponse LocationResponse
	if err := json.Unmarshal(resp.Body(), &locationResponse); err != nil {
		return nil, fmt.Errorf("failed to parse location response: %w", err)
	}

	// Convert to map for easy lookup
	locationMap := make(map[int]Location)
	for _, location := range locationResponse.Data {
		locationMap[location.ID] = location
	}

	return locationMap, nil
}

// enhanceResponseWithLocations enhances response data with location names
func (le *LocationEnhancer) enhanceResponseWithLocations(data interface{}, locationMap map[string]map[int]Location) {
	le.enhanceLocationNamesRecursive(data, locationMap)
}

// enhanceLocationNamesRecursive recursively enhances data with location names
func (le *LocationEnhancer) enhanceLocationNamesRecursive(data interface{}, locationMap map[string]map[int]Location) {
	switch v := data.(type) {
	case map[string]interface{}:
		// Enhance country name
		if countryID, ok := v["country_id"].(float64); ok {
			if country, exists := locationMap["countries"][int(countryID)]; exists {
				v["country"] = country.Name
			}
		}

		// Enhance state name
		if stateID, ok := v["state_id"].(float64); ok {
			if state, exists := locationMap["states"][int(stateID)]; exists {
				v["state"] = state.Name
			}
		}

		// Enhance city name
		if cityID, ok := v["city_id"].(float64); ok {
			if city, exists := locationMap["cities"][int(cityID)]; exists {
				v["city"] = city.Name
			}
		}

		// Enhance nationality name
		if nationalityID, ok := v["nationality_id"].(float64); ok {
			if nationality, exists := locationMap["nationalities"][int(nationalityID)]; exists {
				v["nationality"] = nationality.Name
			}
		}

		// Enhance insurance type name
		if insuranceTypeID, ok := v["insurance_type_id"].(float64); ok {
			if insuranceType, exists := locationMap["insurance_types"][int(insuranceTypeID)]; exists {
				v["insurance_type"] = insuranceType.Name
			}
		}

		// Enhance insurance provider name
		if insuranceProviderID, ok := v["insurance_provider_id"].(float64); ok {
			if insuranceProvider, exists := locationMap["insurance_providers"][int(insuranceProviderID)]; exists {
				v["insurance_provider"] = insuranceProvider.Name
			}
		}

		// Recursively process nested objects
		for _, value := range v {
			le.enhanceLocationNamesRecursive(value, locationMap)
		}

	case []interface{}:
		// Recursively process array items
		for _, item := range v {
			le.enhanceLocationNamesRecursive(item, locationMap)
		}
	}
}