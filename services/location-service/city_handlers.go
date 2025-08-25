package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetCityByID returns a specific city by its ID
func (h *Handler) GetCityByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
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
    city, err := h.store.GetCityByID(id, locale)
    if err != nil {
        if err.Error() == "sql: no rows in result set" {
            c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch city"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": city, "message": "City fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetCitiesByIDs gets multiple cities by their IDs
func (h *Handler) GetCitiesByIDs(c *gin.Context) {
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
    cities, err := h.store.GetCitiesByIDs(ids, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cities"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": cities, "message": "Cities fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}
