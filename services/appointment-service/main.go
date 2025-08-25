package main

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/louhibi/healthcare-logging"
)

func main() {
	// Initialize structured logging
	logging.InitLogger("appointment-service")
	
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
	appointmentService := NewAppointmentService(db)

	// Initialize handlers
	appointmentHandler := NewAppointmentHandler(appointmentService)

	// Setup router
	router := gin.Default()
	
	// Add request/response logging middleware
	router.Use(logging.RequestLoggingMiddleware())
	
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Email, X-User-Role, X-Healthcare-Entity-ID, sec-ch-ua, sec-ch-ua-mobile, sec-ch-ua-platform, User-Agent, Accept, Referer")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Authentication middleware
	authMiddleware := func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(401, gin.H{"error": "User authentication required"})
			c.Abort()
			return
		}

		// Convert to int and set in context
		id, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		c.Set("user_id", id)
		c.Next()
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "appointment-service"})
	})

	// Appointment routes (all protected via API Gateway)
	appointments := router.Group("/api/appointments", authMiddleware)
	{
		appointments.GET("/", appointmentHandler.GetAppointments)
		appointments.POST("/", appointmentHandler.CreateAppointment)
		appointments.GET("/:id", appointmentHandler.GetAppointment)
		appointments.PUT("/:id", appointmentHandler.UpdateAppointment)
		appointments.DELETE("/:id", appointmentHandler.DeleteAppointment)
		appointments.PATCH("/:id/status", appointmentHandler.UpdateAppointmentStatus)
		
		// Smart booking endpoints
		appointments.POST("/book", appointmentHandler.BookAppointment)
		appointments.GET("/slots", appointmentHandler.GetTimeSlots)
		
		// Duration options for appointment booking (moved from admin)
		appointments.GET("/duration-options", appointmentHandler.GetDurationOptions)
		
		// Rooms for appointment booking (moved from admin)
		appointments.GET("/rooms", appointmentHandler.GetRooms)
	}

	// Doctor schedule routes (for appointment availability slots)
	schedules := router.Group("/api/schedules", authMiddleware)
	{
		schedules.GET("/", appointmentHandler.GetDoctorSchedules)
	}

	// Doctor availability routes (for managing doctor working hours and status)
	availability := router.Group("/api/availability", authMiddleware)
	{
		availability.GET("/", appointmentHandler.GetDoctorAvailability)
		availability.POST("/", appointmentHandler.CreateDoctorAvailability)
		availability.GET("/:id", appointmentHandler.GetDoctorAvailabilityByID)
		availability.PUT("/:id", appointmentHandler.UpdateDoctorAvailability)
		availability.DELETE("/:id", appointmentHandler.DeleteDoctorAvailability)
		availability.GET("/calendar", appointmentHandler.GetAvailabilityCalendar)
		availability.POST("/bulk", appointmentHandler.CreateBulkAvailability)
	}

	// Doctor management routes
	doctors := router.Group("/api/doctors", authMiddleware)
	{
		doctors.GET("/", appointmentHandler.GetDoctorsByEntity)
	}

	// Admin management routes
	admin := router.Group("/api/admin", authMiddleware)
	{
		// Duration settings management
		admin.GET("/duration-settings", appointmentHandler.GetAppointmentDurationSettings)
		admin.POST("/duration-settings", appointmentHandler.CreateAppointmentDurationSetting)
		
		// Room management
		admin.GET("/rooms", appointmentHandler.GetRooms)
		admin.POST("/rooms", appointmentHandler.CreateRoom)
		admin.PUT("/rooms/:id", appointmentHandler.UpdateRoom)
		admin.DELETE("/rooms/:id", appointmentHandler.DeleteRoom)
		
		// Duration options management
		admin.GET("/duration-options", appointmentHandler.GetDurationOptions)
		admin.POST("/duration-options", appointmentHandler.CreateDurationOption)
		admin.PUT("/duration-options/:id", appointmentHandler.UpdateDurationOption)
		admin.DELETE("/duration-options/:id", appointmentHandler.DeleteDurationOption)
	}

	// Available rooms endpoint (for appointment booking)
	router.GET("/api/appointments/available-rooms", authMiddleware, appointmentHandler.GetAvailableRooms)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	logging.LogInfo("Appointment service starting", "port", port)
	if err := router.Run(":" + port); err != nil {
		logging.LogError("Failed to start appointment service", "error", err)
		os.Exit(1)
	}
}