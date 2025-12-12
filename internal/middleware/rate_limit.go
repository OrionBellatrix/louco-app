package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/config"
	"github.com/louco-event/internal/dto"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func newRateLimiter(requestsPerMinute, burstSize int) *rateLimiter {
	return &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate.Limit(requestsPerMinute) / 60, // Convert per minute to per second
		burst:    burstSize,
	}
}

func (rl *rateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *rateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func RateLimit(cfg config.RateLimitConfig) gin.HandlerFunc {
	rl := newRateLimiter(cfg.RequestsPerMinute, cfg.BurstSize)

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return func(c *gin.Context) {
		limiter := rl.getVisitor(c.ClientIP())

		if !limiter.Allow() {
			response := dto.NewErrorResponse("common.too_many_requests", nil)
			c.JSON(http.StatusTooManyRequests, response)
			c.Abort()
			return
		}

		c.Next()
	}
}
