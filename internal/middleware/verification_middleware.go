package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/i18n"
	"github.com/louco-event/internal/service"
	"github.com/louco-event/pkg/logger"
)

type VerificationMiddleware struct {
	userService service.UserService
	i18n        *i18n.I18n
	logger      *logger.Logger
}

func NewVerificationMiddleware(userService service.UserService, i18n *i18n.I18n, logger *logger.Logger) *VerificationMiddleware {
	return &VerificationMiddleware{
		userService: userService,
		i18n:        i18n,
		logger:      logger,
	}
}

// RequireVerification middleware checks if user has completed verification
func (m *VerificationMiddleware) RequireVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get language from Accept-Language header
		acceptLang := c.GetHeader("Accept-Language")
		lang := i18n.ExtractLanguageFromHeader(acceptLang)

		// Get user ID from context (set by auth middleware)
		userIDStr, exists := c.Get("user_id")
		if !exists {
			m.logger.Error().Msg("User ID not found in context")
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "auth.unauthorized"),
			})
			c.Abort()
			return
		}

		userID, err := strconv.Atoi(userIDStr.(string))
		if err != nil {
			m.logger.Error().Str("user_id", userIDStr.(string)).Msg("Invalid user ID")
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "auth.unauthorized"),
			})
			c.Abort()
			return
		}

		// Get user profile to check verification status
		userProfile, err := m.userService.GetProfile(c.Request.Context(), userID)
		if err != nil {
			m.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get user profile")
			c.JSON(http.StatusInternalServerError, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "user.not_found"),
			})
			c.Abort()
			return
		}

		// Check if user requires verification
		requiresVerification := false

		// If user has email but it's not verified
		if userProfile.User.Email != nil && userProfile.User.EmailVerifiedAt == nil {
			requiresVerification = true
		}

		// If user has phone but it's not verified
		if userProfile.User.Phone != nil && userProfile.User.PhoneVerifiedAt == nil {
			requiresVerification = true
		}

		if requiresVerification {
			c.JSON(http.StatusForbidden, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "verification.required"),
				Data: map[string]interface{}{
					"email_verified": userProfile.User.EmailVerifiedAt != nil,
					"phone_verified": userProfile.User.PhoneVerifiedAt != nil,
					"has_email":      userProfile.User.Email != nil,
					"has_phone":      userProfile.User.Phone != nil,
				},
			})
			c.Abort()
			return
		}

		// User is verified, continue
		c.Next()
	}
}

// RequireEmailVerification middleware checks if user has verified their email
func (m *VerificationMiddleware) RequireEmailVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get language from Accept-Language header
		acceptLang := c.GetHeader("Accept-Language")
		lang := i18n.ExtractLanguageFromHeader(acceptLang)

		// Get user ID from context (set by auth middleware)
		userIDStr, exists := c.Get("user_id")
		if !exists {
			m.logger.Error().Msg("User ID not found in context")
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "auth.unauthorized"),
			})
			c.Abort()
			return
		}

		userID, err := strconv.Atoi(userIDStr.(string))
		if err != nil {
			m.logger.Error().Str("user_id", userIDStr.(string)).Msg("Invalid user ID")
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "auth.unauthorized"),
			})
			c.Abort()
			return
		}

		// Get user profile to check verification status
		userProfile, err := m.userService.GetProfile(c.Request.Context(), userID)
		if err != nil {
			m.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get user profile")
			c.JSON(http.StatusInternalServerError, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "user.not_found"),
			})
			c.Abort()
			return
		}

		// Check if user has email and it's verified
		if userProfile.User.Email == nil || userProfile.User.EmailVerifiedAt == nil {
			c.JSON(http.StatusForbidden, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "verification.email_required"),
			})
			c.Abort()
			return
		}

		// Email is verified, continue
		c.Next()
	}
}

// RequirePhoneVerification middleware checks if user has verified their phone
func (m *VerificationMiddleware) RequirePhoneVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get language from Accept-Language header
		acceptLang := c.GetHeader("Accept-Language")
		lang := i18n.ExtractLanguageFromHeader(acceptLang)

		// Get user ID from context (set by auth middleware)
		userIDStr, exists := c.Get("user_id")
		if !exists {
			m.logger.Error().Msg("User ID not found in context")
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "auth.unauthorized"),
			})
			c.Abort()
			return
		}

		userID, err := strconv.Atoi(userIDStr.(string))
		if err != nil {
			m.logger.Error().Str("user_id", userIDStr.(string)).Msg("Invalid user ID")
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "auth.unauthorized"),
			})
			c.Abort()
			return
		}

		// Get user profile to check verification status
		userProfile, err := m.userService.GetProfile(c.Request.Context(), userID)
		if err != nil {
			m.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get user profile")
			c.JSON(http.StatusInternalServerError, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "user.not_found"),
			})
			c.Abort()
			return
		}

		// Check if user has phone and it's verified
		if userProfile.User.Phone == nil || userProfile.User.PhoneVerifiedAt == nil {
			c.JSON(http.StatusForbidden, dto.APIResponse{
				Success: false,
				Message: m.i18n.Translate(lang, "verification.phone_required"),
			})
			c.Abort()
			return
		}

		// Phone is verified, continue
		c.Next()
	}
}
