package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/i18n"
	"github.com/louco-event/internal/middleware"
	"github.com/louco-event/internal/service"
)

type VerificationHandler struct {
	verificationService service.VerificationService
	i18n                *i18n.I18n
}

func NewVerificationHandler(verificationService service.VerificationService, i18n *i18n.I18n) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
		i18n:                i18n,
	}
}

// SendVerification godoc
// @Summary Send verification code
// @Description Send verification code to email or phone based on identifier
// @Tags verification
// @Accept json
// @Produce json
// @Param request body dto.SendVerificationRequest true "Send verification request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/verification/send [post]
// @Security BearerAuth
func (h *VerificationHandler) SendVerification(c *gin.Context) {
	var req dto.SendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "validation.invalid_request"),
			Errors:  []string{err.Error()},
		})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "auth.unauthorized"),
		})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "common.internal_error"),
		})
		return
	}

	// Determine if identifier is email or phone
	var err error
	if isValidEmail(req.Identifier) {
		// Send email verification
		err = h.verificationService.SendEmailVerification(c.Request.Context(), userIDInt, req.Identifier, getLanguageFromContext(c))
	} else if isValidPhone(req.Identifier) {
		// Send phone verification
		err = h.verificationService.SendPhoneVerification(c.Request.Context(), userIDInt, req.Identifier, getLanguageFromContext(c))
	} else {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "validation.invalid_identifier"),
			Errors:  []string{"Identifier must be a valid email or phone number"},
		})
		return
	}

	if err != nil {
		if strings.Contains(err.Error(), "already verified") {
			c.JSON(http.StatusBadRequest, dto.APIResponse{
				Success: false,
				Message: middleware.Translate(c, "verification.already_verified"),
				Errors:  []string{err.Error()},
			})
			return
		}

		if strings.Contains(err.Error(), "already sent") {
			c.JSON(http.StatusBadRequest, dto.APIResponse{
				Success: false,
				Message: middleware.Translate(c, "verification.code_already_sent"),
				Errors:  []string{err.Error()},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "verification.send_failed"),
			Errors:  []string{err.Error()},
		})
		return
	}

	// Determine response message based on identifier type
	var message string
	if isValidEmail(req.Identifier) {
		message = middleware.Translate(c, "verification.email_sent")
	} else {
		message = middleware.Translate(c, "verification.sms_sent")
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: message,
		Data: map[string]interface{}{
			"expires_in_minutes": 10,
		},
	})
}

// VerifyCode godoc
// @Summary Verify code
// @Description Verify email or phone with the received code
// @Tags verification
// @Accept json
// @Produce json
// @Param request body dto.VerifyCodeRequest true "Verify code request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/verification/verify [post]
// @Security BearerAuth
func (h *VerificationHandler) VerifyCode(c *gin.Context) {
	var req dto.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "validation.invalid_request"),
			Errors:  []string{err.Error()},
		})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "auth.unauthorized"),
		})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "common.internal_error"),
		})
		return
	}

	// Determine if identifier is email or phone and verify accordingly
	var err error
	if isValidEmail(req.Identifier) {
		// Verify email code
		err = h.verificationService.VerifyEmailCode(c.Request.Context(), userIDInt, req.Identifier, req.Code)
	} else if isValidPhone(req.Identifier) {
		// Verify phone code
		err = h.verificationService.VerifyPhoneCode(c.Request.Context(), userIDInt, req.Identifier, req.Code)
	} else {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "validation.invalid_identifier"),
			Errors:  []string{"Identifier must be a valid email or phone number"},
		})
		return
	}

	if err != nil {
		if strings.Contains(err.Error(), "already verified") {
			c.JSON(http.StatusBadRequest, dto.APIResponse{
				Success: false,
				Message: middleware.Translate(c, "verification.already_verified"),
				Errors:  []string{err.Error()},
			})
			return
		}

		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "expired") {
			c.JSON(http.StatusBadRequest, dto.APIResponse{
				Success: false,
				Message: middleware.Translate(c, "verification.invalid_code"),
				Errors:  []string{err.Error()},
			})
			return
		}

		if strings.Contains(err.Error(), "maximum") {
			c.JSON(http.StatusBadRequest, dto.APIResponse{
				Success: false,
				Message: middleware.Translate(c, "verification.max_attempts"),
				Errors:  []string{err.Error()},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "verification.verify_failed"),
			Errors:  []string{err.Error()},
		})
		return
	}

	// Determine response message based on identifier type
	var message string
	if isValidEmail(req.Identifier) {
		message = middleware.Translate(c, "verification.email_verified")
	} else {
		message = middleware.Translate(c, "verification.phone_verified")
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: message,
	})
}

// ResendVerification godoc
// @Summary Resend verification code
// @Description Resend verification code to email or phone based on identifier
// @Tags verification
// @Accept json
// @Produce json
// @Param request body dto.ResendVerificationRequest true "Resend verification request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/verification/resend [post]
// @Security BearerAuth
func (h *VerificationHandler) ResendVerification(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "validation.invalid_request"),
			Errors:  []string{err.Error()},
		})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "auth.unauthorized"),
		})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "common.internal_error"),
		})
		return
	}

	// Determine if identifier is email or phone
	var err error
	if isValidEmail(req.Identifier) {
		// Resend email verification
		err = h.verificationService.SendEmailVerification(c.Request.Context(), userIDInt, req.Identifier, getLanguageFromContext(c))
	} else if isValidPhone(req.Identifier) {
		// Resend phone verification
		err = h.verificationService.SendPhoneVerification(c.Request.Context(), userIDInt, req.Identifier, getLanguageFromContext(c))
	} else {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "validation.invalid_identifier"),
			Errors:  []string{"Identifier must be a valid email or phone number"},
		})
		return
	}

	if err != nil {
		if strings.Contains(err.Error(), "already verified") {
			c.JSON(http.StatusBadRequest, dto.APIResponse{
				Success: false,
				Message: middleware.Translate(c, "verification.already_verified"),
				Errors:  []string{err.Error()},
			})
			return
		}

		if strings.Contains(err.Error(), "already sent") {
			c.JSON(http.StatusBadRequest, dto.APIResponse{
				Success: false,
				Message: middleware.Translate(c, "verification.code_already_sent"),
				Errors:  []string{err.Error()},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "verification.resend_failed"),
			Errors:  []string{err.Error()},
		})
		return
	}

	// Determine response message based on identifier type
	var message string
	if isValidEmail(req.Identifier) {
		message = middleware.Translate(c, "verification.email_resent")
	} else {
		message = middleware.Translate(c, "verification.sms_resent")
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: message,
		Data: map[string]interface{}{
			"expires_in_minutes": 10,
		},
	})
}
