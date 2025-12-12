package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/i18n"
	"github.com/louco-event/internal/middleware"
	"github.com/louco-event/internal/service"
)

type UserHandler struct {
	userService service.UserService
	i18n        *i18n.I18n
}

func NewUserHandler(userService service.UserService, i18n *i18n.I18n) *UserHandler {
	return &UserHandler{
		userService: userService,
		i18n:        i18n,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userProfile, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		var message string
		if err.Error() == "user not found" {
			message = middleware.Translate(c, "user.not_found")
		} else {
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "common.success"),
		userProfile,
	)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		var message string
		switch err.Error() {
		case "user not found":
			message = middleware.Translate(c, "user.not_found")
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
		middleware.Translate(c, "user.profile_updated"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateContact(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.userService.UpdateContact(c.Request.Context(), userID, &req)
	if err != nil {
		var message string
		switch err.Error() {
		case "user not found":
			message = middleware.Translate(c, "user.not_found")
		case "email already exists":
			message = middleware.Translate(c, "user.email_already_exists")
		case "phone already exists":
			message = middleware.Translate(c, "user.phone_already_exists")
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

func (h *UserHandler) SetProfilePic(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.SetProfilePicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.userService.SetProfilePic(c.Request.Context(), userID, &req)
	if err != nil {
		var message string
		switch err.Error() {
		case "user not found":
			message = middleware.Translate(c, "user.not_found")
		default:
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "user.profile_pic_updated"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) SetCoverPic(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.SetCoverPicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.userService.SetCoverPic(c.Request.Context(), userID, &req)
	if err != nil {
		var message string
		switch err.Error() {
		case "user not found":
			message = middleware.Translate(c, "user.not_found")
		default:
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "user.cover_pic_updated"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) DeactivateAccount(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	err := h.userService.DeactivateAccount(c.Request.Context(), userID)
	if err != nil {
		var message string
		if err.Error() == "user not found" {
			message = middleware.Translate(c, "user.not_found")
		} else {
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "common.deleted"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	var req dto.UserListRequest

	// Parse query parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			req.PageSize = pageSize
		}
	}

	if userType := c.Query("user_type"); userType != "" {
		req.UserType = userType
	}

	result, err := h.userService.GetUserList(c.Request.Context(), &req)
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
