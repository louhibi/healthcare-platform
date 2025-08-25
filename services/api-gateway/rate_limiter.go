package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	mu          sync.RWMutex
	buckets     map[string]*TokenBucket
	config      RateLimitConfig
	cleanupTick *time.Ticker
}

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	tokens       int
	lastRefill   time.Time
	capacity     int
	refillRate   int // tokens per minute
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*TokenBucket),
		config:  config,
	}

	// Start cleanup goroutine to remove old buckets
	rl.cleanupTick = time.NewTicker(5 * time.Minute)
	go rl.cleanup()

	return rl
}

// RateLimitMiddleware returns a rate limiting middleware
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as the key for rate limiting
		clientIP := c.ClientIP()
		
		if !rl.allowRequest(clientIP) {
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error:   "Rate limit exceeded",
				Code:    "RATE_LIMIT_EXCEEDED",
				Message: "Too many requests. Please try again later.",
				Time:    time.Now(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// allowRequest checks if a request should be allowed based on rate limiting
func (rl *RateLimiter) allowRequest(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[key]
	if !exists {
		bucket = &TokenBucket{
			tokens:     rl.config.BurstSize,
			lastRefill: time.Now(),
			capacity:   rl.config.BurstSize,
			refillRate: rl.config.RequestsPerMinute,
		}
		rl.buckets[key] = bucket
	}

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)
	tokensToAdd := int(elapsed.Minutes() * float64(bucket.refillRate))
	
	if tokensToAdd > 0 {
		bucket.tokens += tokensToAdd
		if bucket.tokens > bucket.capacity {
			bucket.tokens = bucket.capacity
		}
		bucket.lastRefill = now
	}

	// Check if request can be allowed
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// cleanup removes old unused buckets
func (rl *RateLimiter) cleanup() {
	for range rl.cleanupTick.C {
		rl.mu.Lock()
		now := time.Now()
		
		for key, bucket := range rl.buckets {
			// Remove buckets that haven't been used for 10 minutes
			if now.Sub(bucket.lastRefill) > 10*time.Minute {
				delete(rl.buckets, key)
			}
		}
		
		rl.mu.Unlock()
	}
}

// Stop stops the rate limiter cleanup goroutine
func (rl *RateLimiter) Stop() {
	if rl.cleanupTick != nil {
		rl.cleanupTick.Stop()
	}
}

// GetBucketInfo returns information about a specific bucket (for debugging)
func (rl *RateLimiter) GetBucketInfo(key string) (TokenBucket, bool) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	bucket, exists := rl.buckets[key]
	if !exists {
		return TokenBucket{}, false
	}

	return TokenBucket{
		tokens:     bucket.tokens,
		lastRefill: bucket.lastRefill,
		capacity:   bucket.capacity,
		refillRate: bucket.refillRate,
	}, true
}