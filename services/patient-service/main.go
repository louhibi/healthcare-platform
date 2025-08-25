package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	logging "github.com/louhibi/healthcare-logging"
)

func main() {
	// Initialize structured logging
	logging.InitLogger("patient-service")
	
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logging.LogWarn("No .env file found")
	}

	// Initialize database
	db, err := InitDB()
	if err != nil {
		logging.LogError("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	if err := RunMigrations(db); err != nil {
		logging.LogError("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize services
	patientService := NewPatientService(db)

	// Initialize handlers
	patientHandler := NewPatientHandler(patientService)

	// Setup router
	router := gin.Default()
	
	// Add request/response logging middleware
	router.Use(logging.RequestLoggingMiddleware())
	
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Email, X-User-Role, X-Healthcare-Entity-ID, sec-ch-ua, sec-ch-ua-mobile, sec-ch-ua-platform, User-Agent, Accept, Referer")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "patient-service"})
	})

	// API root group
	apiGroup := router.Group("/api")

	// Patients subgroup (requires auth)
	patients := apiGroup.Group("/patients")
	patients.Use(AuthMiddleware())
	{
		patients.POST("/", patientHandler.CreatePatient)
		patients.GET("/", patientHandler.GetPatients)
		patients.GET("/stats", patientHandler.GetPatientStats)
		patients.GET("/:id", patientHandler.GetPatient)
		patients.PUT("/:id", patientHandler.UpdatePatient)
		patients.DELETE("/:id", patientHandler.DeletePatient)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	logging.LogInfo("Patient service starting", "port", port)
	if err := router.Run(":" + port); err != nil {
		logging.LogError("Failed to start patient service", "error", err)
		os.Exit(1)
	}
}