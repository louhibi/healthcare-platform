package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    logging "github.com/louhibi/healthcare-logging"
)

func getEnv(key, def string) string { if v := os.Getenv(key); v != "" { return v }; return def }

func main() {
    logging.InitLogger("config-service")
    _ = godotenv.Load()

    db, err := InitDB(); if err != nil { logging.LogError("Failed to connect DB", "error", err); os.Exit(1) }
    defer db.Close()

    if err := RunMigrations(db); err != nil { logging.LogError("Migrations failed", "error", err); os.Exit(1) }

    svc := NewConfigService(db)
    h := NewConfigHandler(svc)


    r := gin.Default()
    r.Use(logging.RequestLoggingMiddleware())
    r.Use(func(c *gin.Context){
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Email, X-User-Role")
        if c.Request.Method == "OPTIONS" { c.AbortWithStatus(204); return }
        c.Next()
    })

    r.GET("/health", func(c *gin.Context){ c.JSON(200, gin.H{"status":"healthy", "service":"config-service"}) })

    // Public endpoints (API gateway marks auth not required)
    pub := r.Group("/api/config")
    {
        pub.GET("/bootstrap", h.GetBootstrap)
        pub.GET("/public/settings", h.GetPublicSettings)
        pub.GET("/public/flags", h.GetPublicFlags)
    }

    // Protected endpoints (gateway enforces auth & role) - treat same paths here
    admin := r.Group("/api/config/admin")
    {
        admin.GET("/settings", h.GetAllSettings)
        admin.POST("/settings", h.UpsertSetting)
        admin.GET("/flags", h.GetAllFlags)
        admin.POST("/flags", h.UpsertFlag)
    }

    port := getEnv("PORT", "8085")
    logging.LogInfo("Config service starting", "port", port)
    if err := r.Run(":"+port); err != nil { logging.LogError("Failed to start", "error", err); os.Exit(1) }
}
