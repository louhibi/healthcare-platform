package main

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService() *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}
	
	return &AuthService{
		jwtSecret: []byte(secret),
	}
}

// HashPassword hashes a password using bcrypt
func (a *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword verifies password against hash
func (a *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateTokens generates access and refresh tokens
func (a *AuthService) GenerateTokens(user *User) (string, string, error) {
	// Access token (expires in 15 minutes)
	accessClaims := jwt.MapClaims{
		"user_id":             user.ID,
		"email":               user.Email,
		"role":                user.Role,
		"healthcare_entity_id": user.HealthcareEntityID,
		"is_temp_password":    user.IsTempPassword,
		"exp":                 time.Now().Add(time.Minute * 15).Unix(),
		"iat":                 time.Now().Unix(),
		"type":                "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(a.jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh token (expires in 7 days)
	refreshClaims := jwt.MapClaims{
		"user_id":             user.ID,
		"email":               user.Email,
		"healthcare_entity_id": user.HealthcareEntityID,
		"exp":                 time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":                 time.Now().Unix(),
		"type":                "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(a.jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// ValidateToken validates and parses JWT token
func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return a.jwtSecret, nil
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

	// is_temp_password is optional (default false for existing tokens)
	isTempPassword := false
	if tempPass, exists := claims["is_temp_password"]; exists {
		if tempPassBool, ok := tempPass.(bool); ok {
			isTempPassword = tempPassBool
		}
	}

	return &Claims{
		UserID:             int(userID),
		Email:              email,
		Role:               role,
		HealthcareEntityID: int(healthcareEntityID),
		IsTempPassword:     isTempPassword,
	}, nil
}

// ValidateRefreshToken validates refresh token
func (a *AuthService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return a.jwtSecret, nil
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
	if !ok || tokenType != "refresh" {
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

	healthcareEntityID, ok := claims["healthcare_entity_id"].(float64)
	if !ok {
		return nil, errors.New("invalid healthcare_entity_id in token")
	}

	return &Claims{
		UserID:             int(userID),
		Email:              email,
		HealthcareEntityID: int(healthcareEntityID),
	}, nil
}