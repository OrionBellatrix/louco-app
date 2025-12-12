package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/pkg/logger"
)

func Recovery(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)

				logger.Error().
					Str("request_id", requestID).
					Interface("error", err).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Str("client_ip", c.ClientIP()).
					Msg("Panic recovered")

				// Return error response
				response := dto.NewErrorResponse("common.internal_server_error", nil)
				c.JSON(http.StatusInternalServerError, response)
				c.Abort()
			}
		}()

		c.Next()
	}
}
