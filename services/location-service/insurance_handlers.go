package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetInsuranceProvidersByType gets insurance providers for a specific insurance type
func (h *Handler) GetInsuranceProvidersByType(c *gin.Context) {
    typeIDStr := c.Param("type_id")
    typeID, err := strconv.Atoi(typeIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insurance type ID"})
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
    providers, err := h.store.GetInsuranceProvidersByType(typeID, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insurance providers"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": providers, "message": "Insurance providers fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetInsuranceTypesByIDs bulk fetch insurance types
func (h *Handler) GetInsuranceTypesByIDs(c *gin.Context) {
    idsParam := c.Query("ids")
    if idsParam == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ids query parameter is required"})
        return
    }
    parts := strings.Split(idsParam, ",")
    ids := make([]int, 0, len(parts))
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p == "" { continue }
        id, err := strconv.Atoi(p)
        if err != nil { continue }
        ids = append(ids, id)
    }
    if len(ids) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no valid IDs provided"})
        return
    }
    locale := c.Query("locale")
    if locale == "" {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang != "" {
            parts := strings.Split(acceptLang, ",")
            if len(parts) > 0 {
                langParts := strings.Split(strings.TrimSpace(parts[0]), "-")
                if len(langParts) > 0 { locale = strings.ToLower(langParts[0]) }
            }
        }
    }
    types, err := h.store.GetInsuranceTypesByIDs(ids, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insurance types"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": types, "message": "Insurance types fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetInsuranceProvidersByIDs bulk fetch insurance providers
func (h *Handler) GetInsuranceProvidersByIDs(c *gin.Context) {
    idsParam := c.Query("ids")
    if idsParam == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ids query parameter is required"})
        return
    }
    parts := strings.Split(idsParam, ",")
    ids := make([]int, 0, len(parts))
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p == "" { continue }
        id, err := strconv.Atoi(p)
        if err != nil { continue }
        ids = append(ids, id)
    }
    if len(ids) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no valid IDs provided"})
        return
    }
    locale := c.Query("locale")
    if locale == "" {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang != "" {
            parts := strings.Split(acceptLang, ",")
            if len(parts) > 0 {
                langParts := strings.Split(strings.TrimSpace(parts[0]), "-")
                if len(langParts) > 0 { locale = strings.ToLower(langParts[0]) }
            }
        }
    }
    providers, err := h.store.GetInsuranceProvidersByIDs(ids, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insurance providers"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": providers, "message": "Insurance providers fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}
