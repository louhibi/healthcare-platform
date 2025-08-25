package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authorization header required",
				Code:    "AUTH_HEADER_MISSING",
				Message: "Please provide Authorization header with Bearer token",
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid authorization header format",
				Code:    "AUTH_HEADER_INVALID",
				Message: "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Validate token
		claims, err := validateJWT(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Invalid token",
				Code:    "TOKEN_INVALID",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		// Set claims in context
		c.Set("user_claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("healthcare_entity_id", claims.HealthcareEntityID)
		
		// Add user headers for downstream services  
		c.Request.Header.Set("X-User-ID", fmt.Sprintf("%d", claims.UserID))
		c.Request.Header.Set("X-User-Email", claims.Email)
		c.Request.Header.Set("X-User-Role", claims.Role)
		c.Request.Header.Set("X-Healthcare-Entity-ID", fmt.Sprintf("%d", claims.HealthcareEntityID))
		
		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(allowedRoles) == 0 {
			c.Next()
			return
		}

		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "Role information not found",
				Code:    "ROLE_MISSING",
				Message: "User role information is missing from request context",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Insufficient permissions",
			Code:    "ROLE_FORBIDDEN",
			Message: "User role '" + role + "' is not allowed to access this resource",
		})
		c.Abort()
	}
}

// validateJWT validates JWT token and returns claims
func validateJWT(tokenString, jwtSecret string) (*UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role in token")
	}

	healthcareEntityID, ok := claims["healthcare_entity_id"].(float64)
	if !ok {
		return nil, errors.New("invalid healthcare_entity_id in token")
	}

	return &UserClaims{
		UserID:             int(userID),
		Email:              email,
		Role:               role,
		HealthcareEntityID: int(healthcareEntityID),
	}, nil
}