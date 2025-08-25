package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// resolveCountryIdentifier converts country identifier (ID or code) to appropriate format
func resolveCountryIdentifier(identifier string) interface{} {
    if countryID, err := strconv.Atoi(identifier); err == nil {
        return countryID
    }
    return strings.ToUpper(identifier)
}

// GetCountries returns all countries (optionally localized)
func (h *Handler) GetCountries(c *gin.Context) {
    locale := c.Query("locale")
    if locale == "" {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang != "" {
            parts := strings.Split(acceptLang, ",")
            if len(parts) > 0 {
                langParts := strings.Split(strings.TrimSpace(parts[0]), "-")
                if len(langParts) > 0 {
                    locale = strings.ToLower(langParts[0])
                }
            }
        }
    }

    countries, err := h.store.GetCountries(locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch countries"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": countries, "message": "Countries fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetStatesByCountry returns states for a country (ID or code)
func (h *Handler) GetStatesByCountry(c *gin.Context) {
    identifier := c.Param("code")
    locale := c.Query("locale")
    if locale == "" {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang != "" {
            parts := strings.Split(acceptLang, ",")
            if len(parts) > 0 {
                langParts := strings.Split(strings.TrimSpace(parts[0]), "-")
                if len(langParts) > 0 {
                    locale = strings.ToLower(langParts[0])
                }
            }
        }
    }
    countryIdentifier := resolveCountryIdentifier(identifier)
    states, err := h.store.GetStatesByCountry(countryIdentifier, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch states"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": states, "message": "States fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetCitiesByCountry returns cities for a country (filterable by state or query)
func (h *Handler) GetCitiesByCountry(c *gin.Context) {
    identifier := c.Param("code")
    q := strings.TrimSpace(c.Query("q"))
    stateParam := strings.TrimSpace(c.Query("state"))
    locale := c.Query("locale")
    if locale == "" {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang != "" {
            parts := strings.Split(acceptLang, ",")
            if len(parts) > 0 {
                langParts := strings.Split(strings.TrimSpace(parts[0]), "-")
                if len(langParts) > 0 {
                    locale = strings.ToLower(langParts[0])
                }
            }
        }
    }
    limit := 50
    if limitStr := c.Query("limit"); limitStr != "" {
        if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
            limit = parsedLimit
        }
    }
    countryIdentifier := resolveCountryIdentifier(identifier)
    var stateIdentifier interface{}
    if stateParam != "" {
        if stateID, err := strconv.Atoi(stateParam); err == nil {
            stateIdentifier = stateID
        } else {
            stateIdentifier = strings.ToUpper(stateParam)
        }
    }
    cities, err := h.store.GetCitiesByCountry(countryIdentifier, stateIdentifier, q, limit, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cities"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": cities, "message": "Cities fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetCountryByID (kept for direct ID endpoint /country/:id)
func (h *Handler) GetCountryByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country ID"})
        return
    }
    locale := c.Query("locale")
    if locale == "" {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang != "" {
            parts := strings.Split(acceptLang, ",")
            if len(parts) > 0 {
                langParts := strings.Split(strings.TrimSpace(parts[0]), "-")
                if len(langParts) > 0 {
                    locale = strings.ToLower(langParts[0])
                }
            }
        }
    }
    country, err := h.store.GetCountryByID(id, locale)
    if err != nil {
        if err.Error() == "sql: no rows in result set" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch country"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": country, "message": "Country fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetCountriesByIDs bulk fetch
func (h *Handler) GetCountriesByIDs(c *gin.Context) {
    idsStr := c.Query("ids")
    if idsStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "IDs parameter is required"})
        return
    }
    idStrings := strings.Split(idsStr, ",")
    ids := make([]int, 0, len(idStrings))
    for _, idStr := range idStrings {
        id, err := strconv.Atoi(strings.TrimSpace(idStr))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format: " + idStr})
            return
        }
        ids = append(ids, id)
    }
    locale := c.Query("locale")
    if locale == "" { locale = "en" }
    countries, err := h.store.GetCountriesByIDs(ids, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch countries"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": countries, "message": "Countries fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetNationalitiesByCountry nationalities scoped by country identifier
func (h *Handler) GetNationalitiesByCountry(c *gin.Context) {
    identifier := c.Param("code")
    locale := c.Query("locale")
    if locale == "" {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang != "" {
            parts := strings.Split(acceptLang, ",")
            if len(parts) > 0 {
                langParts := strings.Split(strings.TrimSpace(parts[0]), "-")
                if len(langParts) > 0 {
                    locale = strings.ToLower(langParts[0])
                }
            }
        }
    }
    countryIdentifier := resolveCountryIdentifier(identifier)
    nationalities, err := h.store.GetNationalitiesByCountry(countryIdentifier, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch nationalities"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": nationalities, "message": "Nationalities fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetInsuranceTypesByCountry insurance types scoped by country identifier
func (h *Handler) GetInsuranceTypesByCountry(c *gin.Context) {
    identifier := c.Param("code")
    if identifier == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Country ID or code is required"})
        return
    }
    locale := c.Query("locale")
    if locale == "" {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang != "" {
            parts := strings.Split(acceptLang, ",")
            if len(parts) > 0 {
                langParts := strings.Split(strings.TrimSpace(parts[0]), "-")
                if len(langParts) > 0 {
                    locale = strings.ToLower(langParts[0])
                }
            }
        }
    }
    countryIdentifier := resolveCountryIdentifier(identifier)
    insuranceTypes, err := h.store.GetInsuranceTypesByCountry(countryIdentifier, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insurance types"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": insuranceTypes, "message": "Insurance types fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetCountry generic endpoint accepting either numeric ID or ISO code
func (h *Handler) GetCountry(c *gin.Context) {
    codeParam := c.Param("code")
    if codeParam == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Country code or ID is required"})
        return
    }

    locale := c.Query("locale")
    if locale == "" {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang != "" {
            parts := strings.Split(acceptLang, ",")
            if len(parts) > 0 {
                langParts := strings.Split(strings.TrimSpace(parts[0]), "-")
                if len(langParts) > 0 {
                    locale = strings.ToLower(langParts[0])
                }
            }
        }
    }

    // If numeric, reuse existing GetCountryByID logic
    if _, err := strconv.Atoi(codeParam); err == nil {
        // Temporarily wrap existing method (could refactor later)
        c.Params = append(c.Params, gin.Param{Key: "id", Value: codeParam})
        h.GetCountryByID(c)
        return
    }

    // Otherwise treat as ISO code: fetch all (could optimize with direct query later)
    countries, err := h.store.GetCountries(locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch countries"})
        return
    }
    upperCode := strings.ToUpper(codeParam)
    for _, country := range countries {
        if strings.ToUpper(country.Code) == upperCode {
            c.JSON(http.StatusOK, gin.H{"data": country, "message": "Country fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
}
