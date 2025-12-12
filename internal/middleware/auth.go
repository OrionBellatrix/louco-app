package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/service"
)

func JWTAuth(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := dto.NewErrorResponse("auth.token_required", nil)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response := dto.NewErrorResponse("auth.invalid_token_format", nil)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			response := dto.NewErrorResponse("auth.token_invalid", nil)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)
		c.Set("user_type", claims.UserType)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// OptionalJWTAuth middleware that doesn't require authentication but extracts user info if token is present
func OptionalJWTAuth(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				token := tokenParts[1]
				if claims, err := jwtService.ValidateToken(token); err == nil {
					c.Set("user_id", claims.UserID)
					c.Set("user_email", claims.Email)
					c.Set("user_username", claims.Username)
					c.Set("user_type", claims.UserType)
					c.Set("jwt_claims", claims)
				}
			}
		}
		c.Next()
	}
}

// RequireUserType middleware that checks if user has specific user type
func RequireUserType(userType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserType, exists := c.Get("user_type")
		if !exists {
			response := dto.NewErrorResponse("common.unauthorized", nil)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		if currentUserType != userType {
			response := dto.NewErrorResponse("common.forbidden", nil)
			c.JSON(http.StatusForbidden, response)
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUserID extracts current user ID from context
func GetCurrentUserID(c *gin.Context) (int, bool) {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(int); ok {
			return id, true
		}
	}
	return 0, false
}

// GetCurrentUserType extracts current user type from context
func GetCurrentUserType(c *gin.Context) (string, bool) {
	if userType, exists := c.Get("user_type"); exists {
		if uType, ok := userType.(string); ok {
			return uType, true
		}
	}
	return "", false
}
