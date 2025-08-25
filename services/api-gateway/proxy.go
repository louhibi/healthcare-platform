package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type ProxyService struct {
	config           *Config
	client           *resty.Client
	stats            *StatsCollector
	locationEnhancer *LocationEnhancer
}

func NewProxyService(config *Config, stats *StatsCollector, locationEnhancer *LocationEnhancer) *ProxyService {
	client := resty.New()
	client.SetTimeout(config.Timeout)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)

	return &ProxyService{
		config:           config,
		client:           client,
		stats:            stats,
		locationEnhancer: locationEnhancer,
	}
}

// ProxyHandlerForRoute handles proxying requests with route config from context
func (p *ProxyService) ProxyHandlerForRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Get route config from context (set by setupRoute middleware)
		routeInterface, exists := c.Get("route_config")
		if !exists {
			p.stats.RecordRequest("unknown", false, time.Since(startTime))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Route configuration missing",
				Code:    "ROUTE_CONFIG_MISSING", 
				Message: "Route configuration not found in context for " + c.Request.URL.Path,
				Time:    time.Now(),
			})
			return
		}

		route, ok := routeInterface.(RouteConfig)
		if !ok {
			p.stats.RecordRequest("unknown", false, time.Since(startTime))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Invalid route configuration",
				Code:    "INVALID_ROUTE_CONFIG",
				Message: "Invalid route configuration type for " + c.Request.URL.Path,
				Time:    time.Now(),
			})
			return
		}

		log.Printf("Processing route: %s -> %s (path: %s)", route.Path, route.Service, c.Request.URL.Path)
		p.proxyRequest(c, &route, startTime)
	}
}

// ProxyHandler handles proxying requests to appropriate services (legacy method - kept for compatibility)
func (p *ProxyService) ProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Find matching route
		route := p.findMatchingRoute(c.Request.URL.Path)
		if route == nil {
			p.stats.RecordRequest("unknown", false, time.Since(startTime))
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Route not found",
				Code:    "ROUTE_NOT_FOUND",
				Message: "No matching route found for " + c.Request.URL.Path,
				Time:    time.Now(),
			})
			return
		}

		p.proxyRequest(c, route, startTime)
	}
}

// proxyRequest handles the actual proxying logic
func (p *ProxyService) proxyRequest(c *gin.Context, route *RouteConfig, startTime time.Time) {
	// Get service configuration
	serviceConfig, exists := p.config.Services[route.Service]
	if !exists {
		p.stats.RecordRequest(route.Service, false, time.Since(startTime))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Service not configured",
			Code:    "SERVICE_NOT_CONFIGURED",
			Message: "Service " + route.Service + " is not properly configured",
			Time:    time.Now(),
		})
		return
	}

	// Build target URL
	targetPath := c.Request.URL.Path
	if route.StripPrefix {
		targetPath = strings.TrimPrefix(targetPath, route.Path)
	}
	targetURL := serviceConfig.BaseURL + targetPath

	// Add query parameters
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	log.Printf("Proxying %s %s to %s", c.Request.Method, c.Request.URL.Path, targetURL)

	// Read request body
	var requestBody []byte
	if c.Request.Body != nil {
		var err error
		requestBody, err = io.ReadAll(c.Request.Body)
		if err != nil {
			p.stats.RecordRequest(route.Service, false, time.Since(startTime))
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Failed to read request body",
				Code:    "BODY_READ_ERROR",
				Message: err.Error(),
				Time:    time.Now(),
			})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	}

	// Create proxy request
	req := p.client.R()

	// Copy headers (excluding hop-by-hop headers)
	for key, values := range c.Request.Header {
		if !isHopByHopHeader(key) {
			req.SetHeader(key, strings.Join(values, ","))
		}
	}

	// Set body if present
	if len(requestBody) > 0 {
		req.SetBody(requestBody)
	}

	// Add user information to headers for downstream services
	if userID, exists := c.Get("user_id"); exists {
		req.SetHeader("X-User-ID", strconv.Itoa(userID.(int)))
	}
	if userEmail, exists := c.Get("user_email"); exists {
		req.SetHeader("X-User-Email", userEmail.(string))
	}
	if userRole, exists := c.Get("user_role"); exists {
		req.SetHeader("X-User-Role", userRole.(string))
	}

	// Execute request
	var resp *resty.Response
	var err error

	switch c.Request.Method {
	case http.MethodGet:
		resp, err = req.Get(targetURL)
	case http.MethodPost:
		resp, err = req.Post(targetURL)
	case http.MethodPut:
		resp, err = req.Put(targetURL)
	case http.MethodPatch:
		resp, err = req.Patch(targetURL)
	case http.MethodDelete:
		resp, err = req.Delete(targetURL)
	case http.MethodOptions:
		resp, err = req.Options(targetURL)
	case http.MethodHead:
		resp, err = req.Head(targetURL)
	default:
		p.stats.RecordRequest(route.Service, false, time.Since(startTime))
		c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			Error:   "Method not allowed",
			Code:    "METHOD_NOT_ALLOWED",
			Message: "HTTP method " + c.Request.Method + " is not supported",
			Time:    time.Now(),
		})
		return
	}

	// Handle request error
	if err != nil {
		log.Printf("Proxy error for %s: %v", targetURL, err)
		p.stats.RecordRequest(route.Service, false, time.Since(startTime))
		c.JSON(http.StatusBadGateway, ErrorResponse{
			Error:   "Service unavailable",
			Code:    "SERVICE_UNAVAILABLE",
			Message: "Failed to connect to " + route.Service,
			Time:    time.Now(),
		})
		return
	}

	// Record stats
	success := resp.StatusCode() < 400
	p.stats.RecordRequest(route.Service, success, time.Since(startTime))

	// Copy response headers (excluding hop-by-hop headers)
	for key, values := range resp.Header() {
		if !isHopByHopHeader(key) {
			for _, value := range values {
				c.Header(key, value)
			}
		}
	}

	// Enhance response with location data if applicable
	responseBody := resp.Body()
	if success && p.locationEnhancer != nil && p.shouldEnhanceLocationData(c.Request.URL.Path, c.Request.Method, resp.Header()) {
		// Get user locale from JWT claims or default to 'en'
		userLocale := "en"
		if locale, exists := c.Get("user_locale"); exists {
			if localeStr, ok := locale.(string); ok {
				userLocale = localeStr
			}
		}

		if enhancedBody, err := p.locationEnhancer.enhanceLocationData(responseBody, userLocale); err == nil {
			responseBody = enhancedBody
			// Update Content-Length header
			c.Header("Content-Length", strconv.Itoa(len(enhancedBody)))
		} else {
			log.Printf("Failed to enhance location data for %s: %v", c.Request.URL.Path, err)
		}
	}

	// Set response status and body
	c.Status(resp.StatusCode())
	c.Writer.Write(responseBody)
}

// findMatchingRoute finds the route that matches the given path
func (p *ProxyService) findMatchingRoute(path string) *RouteConfig {
	log.Printf("Looking for route matching path: %s", path)
	log.Printf("Available routes: %d", len(p.config.Routes))
	
	for i, route := range p.config.Routes {
		log.Printf("  Route %d: %s -> %s", i+1, route.Path, route.Service)
		if p.matchesRoute(path, route.Path) {
			log.Printf("  ✓ Found matching route: %s -> %s", route.Path, route.Service)
			return &route
		}
	}
	log.Printf("  ✗ No matching route found for: %s", path)
	return nil
}

// matchesRoute checks if a path matches a route pattern
func (p *ProxyService) matchesRoute(path, pattern string) bool {
	log.Printf("    Matching '%s' against pattern '%s'", path, pattern)
	
	// Handle exact match first
	if path == pattern {
		log.Printf("    ✓ Exact match")
		return true
	}

	// Handle wildcard patterns ending with /*path
	if strings.HasSuffix(pattern, "/*path") {
		prefix := strings.TrimSuffix(pattern, "/*path")
		if strings.HasPrefix(path, prefix) {
			log.Printf("    ✓ Wildcard match with /*path")
			return true
		}
	}

	// Handle wildcard patterns ending with *
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		if strings.HasPrefix(path, prefix) {
			log.Printf("    ✓ Wildcard match")
			return true
		}
	}

	// Handle path parameters (simple implementation)
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		log.Printf("    ✗ Different number of path parts: %d vs %d", len(pathParts), len(patternParts))
		return false
	}

	for i, patternPart := range patternParts {
		if patternPart != pathParts[i] && !strings.HasPrefix(patternPart, ":") {
			log.Printf("    ✗ Part mismatch at index %d: '%s' vs '%s'", i, pathParts[i], patternPart)
			return false
		}
	}

	log.Printf("    ✓ Parameter-based match")
	return true
}

// isHopByHopHeader checks if a header should not be forwarded
func isHopByHopHeader(header string) bool {
	hopByHopHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	headerLower := strings.ToLower(header)
	for _, hopByHop := range hopByHopHeaders {
		if strings.ToLower(hopByHop) == headerLower {
			return true
		}
	}
	return false
}

// HealthCheckHandler handles health checks for all services
func (p *ProxyService) HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		healthChecks := make(map[string]HealthCheck)

		for serviceName, serviceConfig := range p.config.Services {
			healthURL := serviceConfig.BaseURL + "/health"
			
			client := p.client.SetTimeout(5 * time.Second)
			resp, err := client.R().Get(healthURL)

			if err != nil || resp.StatusCode() != 200 {
				healthChecks[serviceName] = HealthCheck{
					Service:   serviceName,
					Status:    "unhealthy",
					Timestamp: time.Now(),
				}
			} else {
				healthChecks[serviceName] = HealthCheck{
					Service:   serviceName,
					Status:    "healthy",
					Timestamp: time.Now(),
				}
			}
		}

		// Determine overall health
		overallStatus := "healthy"
		for _, health := range healthChecks {
			if health.Status != "healthy" {
				overallStatus = "unhealthy"
				break
			}
		}

		statusCode := http.StatusOK
		if overallStatus != "healthy" {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, gin.H{
			"status":    overallStatus,
			"timestamp": time.Now(),
			"services":  healthChecks,
			"gateway": HealthCheck{
				Service:   "api-gateway",
				Status:    "healthy",
				Timestamp: time.Now(),
				Version:   "1.0.0",
			},
		})
	}
}

// shouldEnhanceLocationData determines if response should be enhanced with location data
func (p *ProxyService) shouldEnhanceLocationData(path, method string, headers map[string][]string) bool {
	// Only enhance GET requests for now (read operations)
	if method != http.MethodGet {
		return false
	}

	// Only enhance JSON responses
	contentType := ""
	if ct, exists := headers["Content-Type"]; exists && len(ct) > 0 {
		contentType = ct[0]
	}
	if !strings.Contains(contentType, "application/json") {
		return false
	}

	// Enhance responses from services that return location data
	enhancePaths := []string{
		"/api/users/healthcare-entities",     // Healthcare entity endpoints
		"/api/patients",                      // Patient endpoints
		"/api/internal/healthcare-entities",  // Internal healthcare entity endpoints
	}

	for _, enhancePath := range enhancePaths {
		if strings.HasPrefix(path, enhancePath) {
			return true
		}
	}

	return false
}