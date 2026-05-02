package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
}

type bucket struct {
	count    int
	resetAt  time.Time
}

var limiter = &rateLimiter{
	buckets: make(map[string]*bucket),
}

// RateLimit allows at most `limit` requests per `window` per IP.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mu.Lock()
		b, ok := limiter.buckets[ip]
		if !ok || time.Now().After(b.resetAt) {
			b = &bucket{count: 0, resetAt: time.Now().Add(window)}
			limiter.buckets[ip] = b
		}
		b.count++
		count := b.count
		limiter.mu.Unlock()

		if count > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}
