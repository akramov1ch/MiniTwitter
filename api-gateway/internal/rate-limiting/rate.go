package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	RateLimit int
	Window    time.Duration
}

func NewRateLimiter(rateLimit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		RateLimit: rateLimit,
		Window:    window,
	}
}

func (rl *RateLimiter) Limit() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(rl.Window), rl.RateLimit)

	return func(c *gin.Context) {
		// Rate limiting logic
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
