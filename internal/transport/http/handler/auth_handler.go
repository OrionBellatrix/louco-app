package handler

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/i18n"
	"github.com/louco-event/internal/middleware"
	"github.com/louco-event/internal/service"
)

type AuthHandler struct {
	userService         service.UserService
	verificationService service.VerificationService
	i18n                *i18n.I18n
}

func NewAuthHandler(userService service.UserService, verificationService service.VerificationService, i18n *i18n.I18n) *AuthHandler {
	return &AuthHandler{
		userService:         userService,
		verificationService: verificationService,
		i18n:                i18n,
	}
}

// Helper functions for validation
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(phone)
}

func getLanguageFromContext(c *gin.Context) string {
	acceptLang := c.GetHeader("Accept-Language")
	return i18n.ExtractLanguageFromHeader(acceptLang)
}

func (h *AuthHandler) RegisterStep1(c *gin.Context) {
	var req dto.RegisterStep1Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := h.userService.RegisterStep1(c.Request.Context(), &req)
	if err != nil {
		var message string
		switch err.Error() {
		case "email already exists":
			message = middleware.Translate(c, "user.email_already_exists")
		case "phone already exists":
			message = middleware.Translate(c, "user.phone_already_exists")
		case "either email or phone is required":
			message = middleware.Translate(c, "auth.email_or_phone_required")
		default:
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// After successful registration, automatically send verification code
	language := getLanguageFromContext(c)
	userID := result.UserID

	// Determine if identifier is email or phone and send verification code
	var verificationErr error
	var verificationMessage string

	if isValidEmail(req.Identifier) {
		// Send email verification
		verificationErr = h.verificationService.SendEmailVerification(c.Request.Context(), userID, req.Identifier, language)
		if verificationErr == nil {
			verificationMessage = middleware.Translate(c, "verification.email_sent")
		}
	} else if isValidPhone(req.Identifier) {
		// Send phone verification
		verificationErr = h.verificationService.SendPhoneVerification(c.Request.Context(), userID, req.Identifier, language)
		if verificationErr == nil {
			verificationMessage = middleware.Translate(c, "verification.sms_sent")
		}
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"user_id": result.UserID,
		"token":   result.Token,
	}

	// Add verification info to response if verification was sent
	if verificationErr == nil && verificationMessage != "" {
		responseData["verification"] = map[string]interface{}{
			"sent":               true,
			"message":            verificationMessage,
			"expires_in_minutes": 10,
		}
	} else if verificationErr != nil {
		// Log verification error but don't fail the registration
		responseData["verification"] = map[string]interface{}{
			"sent":    false,
			"message": middleware.Translate(c, "verification.send_failed_but_registration_success"),
		}
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "auth.registration_success"),
		responseData,
	)
	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) RegisterStep4(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.RegisterStep4Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.userService.RegisterStep4(c.Request.Context(), userID, &req)
	if err != nil {
		var message string
		switch err.Error() {
		case "email already exists":
			message = middleware.Translate(c, "user.email_already_exists")
		case "phone already exists":
			message = middleware.Translate(c, "user.phone_already_exists")
		case "address is required for creator users":
			message = middleware.Translate(c, "user.address_required")
		case "company name is required for creator users":
			message = middleware.Translate(c, "user.company_name_required")
		default:
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "common.updated"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) SetUsername(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.SetUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.userService.SetUsername(c.Request.Context(), userID, &req)
	if err != nil {
		var message string
		if err.Error() == "username already exists" {
			message = middleware.Translate(c, "user.username_already_exists")
		} else {
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "auth.username_set_success"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) CheckUsername(c *gin.Context) {
	var req dto.CheckUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := h.userService.CheckUsername(c.Request.Context(), &req)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.internal_server_error"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "common.success"),
		result,
	)
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		var message string
		if err.Error() == "invalid credentials" {
			message = middleware.Translate(c, "auth.invalid_credentials")
		} else {
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "auth.login_success"),
		result,
	)
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) SocialLogin(c *gin.Context) {
	var req dto.SocialLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := h.userService.SocialLogin(c.Request.Context(), &req)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.internal_server_error"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "auth.social_login_success"),
		result,
	)
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.userService.ChangePassword(c.Request.Context(), userID, &req)
	if err != nil {
		var message string
		if err.Error() == "current password is incorrect" {
			message = middleware.Translate(c, "auth.current_password_incorrect")
		} else {
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "auth.password_changed"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// TODO: Implement forgot password logic
	response := dto.NewSuccessResponse(
		middleware.Translate(c, "common.success"),
		map[string]string{"message": "Password reset email sent"},
	)
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// TODO: Implement reset password logic
	response := dto.NewSuccessResponse(
		middleware.Translate(c, "auth.password_changed"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}
