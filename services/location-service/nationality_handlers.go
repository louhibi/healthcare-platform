package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetNationalities gets all nationalities with locale support
func (h *Handler) GetNationalities(c *gin.Context) {
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
    nationalities, err := h.store.GetNationalities(locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch nationalities"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": nationalities, "message": "Nationalities fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetNationalitiesByIDs gets multiple nationalities by their IDs
func (h *Handler) GetNationalitiesByIDs(c *gin.Context) {
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
    nationalities, err := h.store.GetNationalitiesByIDs(ids, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch nationalities"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": nationalities, "message": "Nationalities fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetNationalityByID gets a nationality by its ID
func (h *Handler) GetNationalityByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nationality ID"})
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
    nationality, err := h.store.GetNationalityByID(id, locale)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Nationality not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": nationality, "message": "Nationality fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}
