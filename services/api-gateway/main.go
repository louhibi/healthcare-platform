package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	logging "github.com/louhibi/healthcare-logging"
)

func main() {
	// Initialize structured logging
	logging.InitLogger("api-gateway")
	
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logging.LogWarn("No .env file found")
	}

	// Load configuration
	config := LoadConfig()

	// Initialize components
	stats := NewStatsCollector()
	rateLimiter := NewRateLimiter(config.RateLimit)
	defer rateLimiter.Stop()

	// Initialize location enhancer
	locationServiceURL := config.Services["location-service"].BaseURL
	locationEnhancer := NewLocationEnhancer(locationServiceURL)
	
	proxyService := NewProxyService(config, stats, locationEnhancer)

	// Setup router
	router := gin.Default()

	// Trust proxy headers for correct client IP detection
	router.SetTrustedProxies([]string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})

	// Global middleware
	router.Use(corsMiddleware(config.CORS))
	router.Use(logging.RequestLoggingMiddleware())
	router.Use(rateLimiter.RateLimitMiddleware())

	// Gateway health check
	router.GET("/health", proxyService.HealthCheckHandler())

	// Gateway statistics
	router.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, stats.GetStats())
	})

	// Gateway info
	router.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":   "api-gateway",
			"version":   "1.0.0",
			"timestamp": time.Now(),
			"config": gin.H{
				"services":   len(config.Services),
				"routes":     len(config.Routes),
				"rate_limit": config.RateLimit,
			},
		})
	})

	// Debug: feature flag style endpoint to indicate debug features are enabled.
	// For now this always returns true. Extend later to read from config/env.
	router.GET("/debug/enabled", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"enabled": true})
	})

	// Setup routes with authentication
	for _, route := range config.Routes {
		setupRoute(router, route, config.JWTSecret, proxyService)
	}

	// Catch-all route for undefined paths
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Route not found",
			Code:    "ROUTE_NOT_FOUND",
			Message: "The requested route " + c.Request.URL.Path + " was not found",
			Time:    time.Now(),
		})
	})

	logging.LogInfo("API Gateway starting",
		"port", config.Port,
		"services", getServiceNames(config.Services),
		"rate_limit_per_minute", config.RateLimit.RequestsPerMinute,
		"rate_limit_burst", config.RateLimit.BurstSize)

	if err := router.Run(":" + config.Port); err != nil {
		logging.LogError("Failed to start API Gateway", "error", err)
		os.Exit(1)
	}
}

// setupRoute sets up a route with appropriate middleware
func setupRoute(router *gin.Engine, route RouteConfig, jwtSecret string, proxyService *ProxyService) {
	var handlers []gin.HandlerFunc

	// Add authentication middleware if required
	if route.AuthRequired {
		handlers = append(handlers, AuthMiddleware(jwtSecret))
		
		// Add role-based authorization if roles are specified
		if len(route.RolesAllowed) > 0 {
			handlers = append(handlers, RoleMiddleware(route.RolesAllowed))
		}
	}

	// Add middleware to pass route config to proxy handler
	handlers = append(handlers, func(c *gin.Context) {
		c.Set("route_config", route)
		c.Next()
	})

	// Add proxy handler
	handlers = append(handlers, proxyService.ProxyHandlerForRoute())

	// Register route for all HTTP methods
	router.Any(route.Path, handlers...)
}

// corsMiddleware handles CORS
func corsMiddleware(corsConfig CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		env := os.Getenv("APP_ENV")
		
		logging.LogDebug("CORS check",
			"origin", origin,
			"env", env,
			"allowed_origins", corsConfig.AllowedOrigins)
		
		// Check if origin is allowed
		allowOrigin := "*"
		if len(corsConfig.AllowedOrigins) > 0 {
			allowOrigin = ""
			for _, allowedOrigin := range corsConfig.AllowedOrigins {
				logging.LogDebug("Checking CORS origin", "origin", origin, "allowed", allowedOrigin)
				if allowedOrigin == origin || allowedOrigin == "*" {
					allowOrigin = origin
					logging.LogDebug("CORS origin matched", "origin", allowedOrigin)
					break
				}
			}
			
			// In non-production, also allow local network (192.168.0.*)
			if env != "production" && origin != "" {
				logging.LogDebug("Checking local network CORS", "origin", origin)
				if strings.Contains(origin, "192.168.0.") || strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
					allowOrigin = origin
					logging.LogDebug("CORS allowed as local network", "origin", origin)
				}
			}
			
			if allowOrigin == "" {
				allowOrigin = corsConfig.AllowedOrigins[0] // Fallback to first allowed origin
				logging.LogDebug("Using fallback CORS origin", "origin", allowOrigin)
			}
		}
		
		logging.LogDebug("Final CORS origin", "allow_origin", allowOrigin)
		
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}


// getServiceNames extracts service names from config
func getServiceNames(services map[string]ServiceConfig) []string {
	names := make([]string, 0, len(services))
	for name := range services {
		names = append(names, name)
	}
	return names
}