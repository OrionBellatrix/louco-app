package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/louco-event/pkg/logger"
)

func Logger(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Get status code
		statusCode := c.Writer.Status()

		// Get user ID if available
		userID := ""
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(int); ok {
				userID = string(rune(id))
			}
		}

		// Build log entry
		logEntry := logger.Info().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", raw).
			Str("client_ip", clientIP).
			Int("status_code", statusCode).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent())

		if userID != "" {
			logEntry = logEntry.Str("user_id", userID)
		}

		// Log based on status code
		if statusCode >= 400 {
			if statusCode >= 500 {
				logEntry.Msg("Server error")
			} else {
				logEntry.Msg("Client error")
			}
		} else {
			logEntry.Msg("Request completed")
		}
	}
}

// GetRequestID extracts request ID from context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
