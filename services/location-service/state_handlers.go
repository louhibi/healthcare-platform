package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetStateByID returns a specific state by its ID
func (h *Handler) GetStateByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state ID"})
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
    state, err := h.store.GetStateByID(id, locale)
    if err != nil {
        if err.Error() == "sql: no rows in result set" {
            c.JSON(http.StatusNotFound, gin.H{"error": "State not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch state"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": state, "message": "State fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetStatesByIDs gets multiple states by their IDs
func (h *Handler) GetStatesByIDs(c *gin.Context) {
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
    states, err := h.store.GetStatesByIDs(ids, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch states"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": states, "message": "States fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

// GetCitiesByState returns cities for a specific state ID
func (h *Handler) GetCitiesByState(c *gin.Context) {
    stateIDStr := c.Param("id")
    stateID, err := strconv.Atoi(stateIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state ID"})
        return
    }
    q := strings.TrimSpace(c.Query("q"))
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
        if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 { limit = parsedLimit }
    }
    cities, err := h.store.GetCitiesByState(stateID, q, limit, locale)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cities"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": cities, "message": "Cities fetched successfully", "timestamp": time.Now().Format(time.RFC3339)})
}

