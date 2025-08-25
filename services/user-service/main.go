package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	logging "github.com/louhibi/healthcare-logging"
)

func main() {
	// Initialize structured logging
	logging.InitLogger("user-service")
	
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
	userService := NewUserService(db)
	authService := NewAuthService()
	formConfigService := NewFormConfigService(db)
	translationService := NewTranslationService(db)

	// Add services to user service
	userService.formConfigService = formConfigService
	userService.translationService = translationService

	// Initialize handlers
	userHandler := NewUserHandler(userService, authService)

	// Setup router
	router := gin.Default()
	
	// Add request/response logging middleware
	router.Use(logging.RequestLoggingMiddleware())
	
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "user-service"})
	})

	// Root API group
	api := router.Group("/api")
	apiProtected := api.Group("")
	apiProtected.Use(AuthMiddleware(authService))

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.RefreshToken)
	}

	// Protected auth routes
	authProtected := apiProtected.Group("/auth")
	{
		authProtected.POST("/change-password", userHandler.ChangePassword)
	}

	// doctors routes
	doctors := apiProtected.Group("/doctors")
	{
		doctors.GET("", userHandler.GetDoctors) // Get doctors for appointment service

	}
	// doctors routes
	entities := apiProtected.Group("/entities")
	{
		entities.GET("/:id", userHandler.GetEntityComplete) // Get entities for appointment service
		entities.GET(":id/room-requirement", userHandler.GetEntityRoomRequirement) // Get entity room requirement setting
	}

	forms := apiProtected.Group("/forms")
	{
		forms.GET(":formType/metadata", userService.GetFormMetadata) // Form metadata
		forms.GET(":formType/metadata/:locale", userService.GetLocalizedFormMetadata) // Localized form metadata
		forms.GET("/types", userService.GetFormTypes)                         // Get all form types
	    forms.GET("/:formType/fields", userService.GetFieldConfigurations)    // Get field configurations
	}

		// Internationalization routes (authenticated)
	i18n := apiProtected.Group("/i18n")
	{
		i18n.GET("/locales", userService.GetSupportedLocales)                 // Get supported locales
		i18n.PUT("/user/locale", userService.UpdateUserLocale)               // Update user's preferred locale
		i18n.GET("/forms/:formType/metadata", userService.GetLocalizedFormMetadata) // Localized form metadata for forms
	}

	// Protected user routes
	users := apiProtected.Group("/users")
	{
		users.GET("/profile", userHandler.GetProfile)
		users.PUT("/profile", userHandler.UpdateProfile)
		users.GET("/", userHandler.GetUsers) // Admin only (enforced in handler/middleware)
	}

	// Admin base group (protected)
	admin := apiProtected.Group("/admin")
	admin.Use(userService.AdminMiddleware())
	{
		// Doctor admin routes
		admin.POST("/doctors", userHandler.AdminCreateDoctor) // Admin creates doctors with temp passwords
		admin.GET("/doctors", userHandler.AdminGetDoctors)    // Admin lists doctors

		// Admin-only form configuration management
		formsAdmin := admin.Group("/forms")
		{
			formsAdmin.PUT("/:formType/fields", userService.UpdateMultipleFieldConfigurations) // Update multiple field configurations
			formsAdmin.PUT("/:formType/fields/:fieldId", userService.UpdateFieldConfiguration) // Update single field configuration
			formsAdmin.PUT("/:formType/fields/order", userService.UpdateFieldOrder)            // Update field sort order
			formsAdmin.POST("/:formType/reset", userService.ResetFormConfiguration)            // Reset form to defaults
		}

		// Translation management (Admin only)
		i18nAdmin := admin.Group("/i18n")
		{
			i18nAdmin.GET("/translations/:locale", userService.GetTranslations)                 // Get all translations for locale
			i18nAdmin.POST("/translations", userService.CreateOrUpdateTranslation)             // Create/update translation
			i18nAdmin.POST("/field-translations", userService.CreateOrUpdateFieldTranslation) // Create/update field translation
		}

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	logging.LogInfo("User service starting", "port", port)
	if err := router.Run(":" + port); err != nil {
		logging.LogError("Failed to start user service", "error", err)
		os.Exit(1)
	}
}