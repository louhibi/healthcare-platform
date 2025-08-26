package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	logging "github.com/louhibi/healthcare-logging"
)

func main() {
	// Define CLI flags for migration management
	migrationCommand := flag.String("migrate", "", "Migration command: 'schema', 'data', 'status', 'flush-data'")
	environmentFlag := flag.String("env", "all", "Environment for data migrations: 'dev', 'prod', 'test', 'all'")
	targetVersion := flag.Int("target", 0, "Target migration version (0 = latest)")
	flag.Parse()

	// Initialize structured logging
	logging.InitLogger("location-service")

	// (Optional) Load .env file for local development
	// Load env
	_ = godotenv.Load()

	// Initialize database
	if err := InitDatabase(); err != nil {
		logging.LogError("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer CloseDatabase()

	// Handle migration commands
	if *migrationCommand != "" {
		if err := handleMigrationCommand(*migrationCommand, *environmentFlag, *targetVersion); err != nil {
			logging.LogError("Migration failed", "error", err)
			os.Exit(1)
		}
		return
	}

	// Run migrations automatically on startup (like other services)
	engine := NewMigrationEngine(db)
	
	// Initialize migration tables
	if err := engine.InitializeMigrationTables(); err != nil {
		logging.LogError("Failed to initialize migration tables", "error", err)
		os.Exit(1)
	}

	logging.LogInfo("Running schema migrations automatically...")
	schemaMigrations := GetSchemaMigrations()
	if err := engine.RunSchemaMigrations(schemaMigrations, 0); err != nil {
		logging.LogError("Failed to run schema migrations", "error", err)
		os.Exit(1)
	}

	// Verify schema migrations reached latest version before proceeding
	currentVersion, err := engine.GetCurrentSchemaVersion()
	if err != nil {
		logging.LogError("Failed to get current schema version", "error", err)
		os.Exit(1)
	}
	expectedVersion := schemaMigrations[len(schemaMigrations)-1].Version
	if currentVersion < expectedVersion {
		logging.LogError("Schema migrations incomplete; skipping data migrations", "current_version", currentVersion, "expected_version", expectedVersion)
		os.Exit(1)
	}

	logging.LogInfo("Running data migrations automatically...")
	dataMigrations := GetDataMigrations()
	if err := engine.RunDataMigrations(dataMigrations, "all", 0); err != nil {
		logging.LogError("Failed to run data migrations", "error", err)
		os.Exit(1)
	}

	logging.LogInfo("All migrations completed successfully")

	// Init store with database connection
	store := NewLocationStore()

	// Router
	r := gin.Default()
	
	// Add request/response logging middleware
	r.Use(logging.RequestLoggingMiddleware())

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Email, X-User-Role, X-Healthcare-Entity-ID, sec-ch-ua, sec-ch-ua-mobile, sec-ch-ua-platform, User-Agent, Accept, Referer")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "location-service"})
	})

	// API
	h := NewHandler(store)
	api := r.Group("/api/locations")
	{
		// Countries grouped routes
		countries := api.Group("/countries")
		{
			// Specific static paths first to avoid ":code" capturing them
			countries.GET("/by-ids", h.GetCountriesByIDs)
			// Base list
			countries.GET("", h.GetCountries)
			// Country-level related resources (code or ID supported in handler)
			countries.GET("/:code", h.GetCountry)
			countries.GET("/:code/states", h.GetStatesByCountry)
			countries.GET("/:code/cities", h.GetCitiesByCountry)
			countries.GET("/:code/nationalities", h.GetNationalitiesByCountry)
			countries.GET("/:code/insurance-types", h.GetInsuranceTypesByCountry)
		}

		// States grouped routes
		states := api.Group("/states")
		{
			states.GET("/by-ids", h.GetStatesByIDs)
			states.GET("/:id", h.GetStateByID)
			states.GET("/:id/cities", h.GetCitiesByState)
		}

		// Nationalities grouped routes
		nationalities := api.Group("/nationalities")
		{
			nationalities.GET("/by-ids", h.GetNationalitiesByIDs)
			nationalities.GET("", h.GetNationalities)
			nationalities.GET("/:id", h.GetNationalityByID)
		}

		// Cities grouped routes
		cities := api.Group("/cities")
		{
			cities.GET("/by-ids", h.GetCitiesByIDs)
			cities.GET("/:id", h.GetCityByID)
		}

		// Insurance types grouped routes
		insurance := api.Group("/insurance-types")
		{
			insurance.GET("/by-ids", h.GetInsuranceTypesByIDs)
			insurance.GET(":type_id/providers", h.GetInsuranceProvidersByType)
		}

		// Insurance providers grouped routes (direct by-ids lookup)
		providers := api.Group("/insurance-providers")
		{
			providers.GET("/by-ids", h.GetInsuranceProvidersByIDs)
		}
	}

	port := os.Getenv("PORT")
	if port == "" { port = "8084" }
	logging.LogInfo("Location service starting", "port", port)
	if err := r.Run(":" + port); err != nil {
		logging.LogError("Failed to start location service", "error", err)
		os.Exit(1)
	}
}

// handleMigrationCommand processes migration CLI commands
func handleMigrationCommand(command, environment string, targetVersion int) error {
	engine := NewMigrationEngine(db)
	
	// Initialize migration tables
	if err := engine.InitializeMigrationTables(); err != nil {
		return fmt.Errorf("failed to initialize migration tables: %w", err)
	}
	
	switch command {
	case "schema":
		logging.LogInfo("Running schema migrations", "target_version", targetVersion)
		schemaMigrations := GetSchemaMigrations()
		if err := engine.RunSchemaMigrations(schemaMigrations, targetVersion); err != nil {
			return fmt.Errorf("schema migrations failed: %w", err)
		}
			logging.LogInfo("Schema migrations completed successfully")
		
	case "data":
		logging.LogInfo("Running data migrations", "environment", environment, "target_version", targetVersion)
		dataMigrations := GetDataMigrations()
		if environment == "dev" || environment == "test" {
			// Include test data migrations in dev/test environments
			testMigrations := GetTestDataMigrations()
			dataMigrations = append(dataMigrations, testMigrations...)
		}
		if err := engine.RunDataMigrations(dataMigrations, environment, targetVersion); err != nil {
			return fmt.Errorf("data migrations failed: %w", err)
		}
			logging.LogInfo("Data migrations completed successfully")
		
	case "status":
		logging.LogInfo("Getting migration status")
		status, err := engine.GetMigrationStatus()
		if err != nil {
			return fmt.Errorf("failed to get migration status: %w", err)
		}
		fmt.Printf("Migration Status:\n")
		fmt.Printf("  Schema Version: %v\n", status["schema_version"])
		fmt.Printf("  Data Migrations Count: %v\n", status["data_migrations_count"])
		fmt.Printf("  Timestamp: %v\n", status["timestamp"])
		
	case "flush-data":
		logging.LogInfo("Flushing all data migrations")
		dataMigrations := GetDataMigrations()
		testMigrations := GetTestDataMigrations()
		allDataMigrations := append(dataMigrations, testMigrations...)
		
		// Run down migrations before flushing
		if err := engine.FlushDataMigrations(allDataMigrations, true); err != nil {
			return fmt.Errorf("failed to flush data migrations: %w", err)
		}
			logging.LogInfo("Data migrations flushed successfully")
		
	default:
		return fmt.Errorf("unknown migration command: %s. Available commands: 'schema', 'data', 'status', 'flush-data'", command)
	}
	
	return nil
}
