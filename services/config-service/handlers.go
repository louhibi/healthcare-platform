package main

import (
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
)

type ConfigHandler struct { svc *ConfigService }

func NewConfigHandler(svc *ConfigService) *ConfigHandler { return &ConfigHandler{svc: svc} }

// Public bootstrap config
func (h *ConfigHandler) GetBootstrap(c *gin.Context) {
    env := os.Getenv("APP_ENV")
    if env == "" { env = "development" }
    c.JSON(http.StatusOK, BootstrapConfig{Environment: env})
}

// Public list of settings (public only)
func (h *ConfigHandler) GetPublicSettings(c *gin.Context) {
    list, err := h.svc.ListSettings(false); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to list settings"}); return }
    c.JSON(http.StatusOK, list)
}

// Admin: list all settings
func (h *ConfigHandler) GetAllSettings(c *gin.Context) {
    list, err := h.svc.ListSettings(true); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to list settings"}); return }
    c.JSON(http.StatusOK, list)
}

type upsertSettingRequest struct { Key string `json:"key" binding:"required"`; Value string `json:"value" binding:"required"`; IsPublic *bool `json:"is_public"`; Description *string `json:"description"` }

func (h *ConfigHandler) UpsertSetting(c *gin.Context) {
    var req upsertSettingRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid payload"}); return }
    isPublic := false; if req.IsPublic != nil { isPublic = *req.IsPublic }
    cs, err := h.svc.UpsertSetting(req.Key, req.Value, isPublic, req.Description); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to upsert"}); return }
    c.JSON(http.StatusOK, cs)
}

// Feature flags
func (h *ConfigHandler) GetPublicFlags(c *gin.Context) {
    flags, err := h.svc.ListFeatureFlags(false); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to list flags"}); return }
    c.JSON(http.StatusOK, flags)
}

func (h *ConfigHandler) GetAllFlags(c *gin.Context) {
    flags, err := h.svc.ListFeatureFlags(true); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to list flags"}); return }
    c.JSON(http.StatusOK, flags)
}

type upsertFlagRequest struct { Name string `json:"name" binding:"required"`; Enabled bool `json:"enabled"`; IsPublic *bool `json:"is_public"`; Description *string `json:"description"` }

func (h *ConfigHandler) UpsertFlag(c *gin.Context) {
    var req upsertFlagRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid payload"}); return }
    isPublic := true; if req.IsPublic != nil { isPublic = *req.IsPublic }
    ff, err := h.svc.UpsertFeatureFlag(req.Name, req.Enabled, isPublic, req.Description); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to upsert"}); return }
    c.JSON(http.StatusOK, ff)
}
