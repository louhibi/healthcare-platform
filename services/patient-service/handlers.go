package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates user authentication (simplified for patient service)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User authentication required"})
			c.Abort()
			return
		}
		id, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}
		c.Set("user_id", id)
		c.Next()
	}
}