package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/pkg/database"
)

func HealthCheck(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database health
		dbStatus := "ok"
		if err := db.Health(); err != nil {
			dbStatus = "error"
		}

		health := dto.HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().Format(time.RFC3339),
			Version:   "1.0.0",
			Database:  dbStatus,
		}

		// If database is down, return 503
		if dbStatus == "error" {
			health.Status = "error"
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		c.JSON(http.StatusOK, health)
	}
}
