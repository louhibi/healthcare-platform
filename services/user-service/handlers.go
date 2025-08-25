package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService *UserService
	authService *AuthService
	validator   *validator.Validate
}

func NewUserHandler(userService *UserService, authService *AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
		validator:   validator.New(),
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var req UserRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	exists, err := h.userService.EmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user
	user := &User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
	}

	if err := h.userService.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := h.authService.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	response := AuthResponse{
		User:         user.ToUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var req UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !h.authService.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := h.authService.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	response := AuthResponse{
		User:         user.ToUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	claims, err := h.authService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Get user
	user, err := h.userService.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate new tokens
	accessToken, refreshToken, err := h.authService.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	response := AuthResponse{
		User:         user.ToUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile gets current user profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims := claims.(*Claims)
	user, err := h.userService.GetUserByID(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToUserResponse())
}


// UpdateProfile updates current user profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims := claims.(*Claims)
	
	var req struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user
	user, err := h.userService.GetUserByID(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update user
	user.FirstName = req.FirstName
	user.LastName = req.LastName

	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user.ToUserResponse())
}

// GetUsers gets all users (admin only)
func (h *UserHandler) GetUsers(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims := claims.(*Claims)
	if userClaims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	users, err := h.userService.GetUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	// Convert to response format
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToUserResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  userResponses,
		"limit":  limit,
		"offset": offset,
	})
}

// GetDoctors gets all doctors for the appointment service
func (h *UserHandler) GetDoctors(c *gin.Context) {
	// Get healthcare entity ID from query params (optional)
	healthcareEntityIDStr := c.Query("healthcare_entity_id")
	var healthcareEntityID int
	var err error
	
	if healthcareEntityIDStr != "" {
		healthcareEntityID, err = strconv.Atoi(healthcareEntityIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid healthcare entity ID"})
			return
		}
	}

	doctors, err := h.userService.GetDoctorsByEntity(healthcareEntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get doctors", "details": err.Error()})
		return
	}

	// Convert to simplified response format for appointment service
	var doctorResponses []gin.H
	for _, doctor := range doctors {
		doctorResponses = append(doctorResponses, gin.H{
			"id":             doctor.ID,
			"first_name":     doctor.FirstName,
			"last_name":      doctor.LastName,
			"email":          doctor.Email,
			"role":           doctor.Role,
			"specialization": doctor.Specialization,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": doctorResponses,
	})
}

// GetEntityComplete gets complete information for a healthcare entity
func (h *UserHandler) GetEntityComplete(c *gin.Context) {
	entityIDStr := c.Param("id")
	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	entity, err := h.userService.GetHealthcareEntityByID(entityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Healthcare entity not found"})
		return
	}

	// Use the conversion method
	response := entity.ToHealthcareEntityResponse()

	c.JSON(http.StatusOK, response)
}

// GetEntityRoomRequirement gets room requirement setting for a healthcare entity
func (h *UserHandler) GetEntityRoomRequirement(c *gin.Context) {
	entityIDStr := c.Param("id")
	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	entity, err := h.userService.GetHealthcareEntityByID(entityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Healthcare entity not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"require_room_assignment": entity.RequireRoomAssignment,
	})
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(authService *AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set claims in context
		c.Set("claims", claims)
		c.Next()
	}
}

// AdminCreateDoctor handles POST /api/admin/doctors - Admin creates a doctor with temp password
func (h *UserHandler) AdminCreateDoctor(c *gin.Context) {
	var req AdminCreateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check authorization - only admins can create doctors
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	
	userClaims := claims.(*Claims)
	if userClaims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Check if email already exists
	exists, err := h.userService.EmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Create doctor with temporary password
	response, err := h.userService.CreateDoctorWithTempPassword(req, userClaims.HealthcareEntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// ChangePassword handles POST /api/auth/change-password - Change user password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req PasswordChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	
	userClaims := claims.(*Claims)

	// Change password
	err := h.userService.ChangePassword(userClaims.UserID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "invalid current password" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
		"timestamp": "now",
	})
}

// AdminGetDoctors handles GET /api/admin/doctors - Admin lists doctors
func (h *UserHandler) AdminGetDoctors(c *gin.Context) {
	// Check authorization - only admins can list doctors
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userClaims := claims.(*Claims)
	if userClaims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Get all doctors for the healthcare entity
	doctors, err := h.userService.GetDoctorsByEntity(userClaims.HealthcareEntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to full response format for admin
	var doctorResponses []UserResponse
	for _, doctor := range doctors {
		doctorResponses = append(doctorResponses, doctor.ToUserResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"doctors":     doctorResponses,
		"total_count": len(doctorResponses),
		"message":     "Doctors retrieved successfully",
	})
}