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

type MediaHandler struct {
	mediaService service.MediaService
	i18n         *i18n.I18n
}

func NewMediaHandler(mediaService service.MediaService, i18n *i18n.I18n) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
		i18n:         i18n,
	}
}

func (h *MediaHandler) UploadFile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"No file provided",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "media.upload_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer fileContent.Close()

	result, err := h.mediaService.UploadFile(c.Request.Context(), userID, file, fileContent)
	if err != nil {
		var message string
		switch err.Error() {
		case "unsupported file type":
			message = middleware.Translate(c, "media.invalid_file_type")
		default:
			if err.Error() == "image file too large" || err.Error() == "video file too large" {
				message = middleware.Translate(c, "media.file_too_large")
			} else {
				message = middleware.Translate(c, "media.upload_failed")
			}
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "media.file_uploaded"),
		result,
	)
	c.JSON(http.StatusCreated, response)
}

func (h *MediaHandler) GetMedia(c *gin.Context) {
	mediaIDStr := c.Param("id")
	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid media ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := h.mediaService.GetMediaByID(c.Request.Context(), mediaID)
	if err != nil {
		var message string
		if err.Error() == "media not found" {
			message = middleware.Translate(c, "media.not_found")
		} else {
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "common.success"),
		result,
	)
	c.JSON(http.StatusOK, response)
}

func (h *MediaHandler) GetUserMedia(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid user ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var req dto.MediaListRequest

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

	if mediaType := c.Query("media_type"); mediaType != "" {
		req.MediaType = mediaType
	}

	result, err := h.mediaService.GetUserMedia(c.Request.Context(), userID, &req)
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

func (h *MediaHandler) UpdateMedia(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	mediaIDStr := c.Param("id")
	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid media ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var req dto.MediaUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.mediaService.UpdateMedia(c.Request.Context(), userID, mediaID, &req)
	if err != nil {
		var message string
		switch err.Error() {
		case "media not found":
			message = middleware.Translate(c, "media.not_found")
		case "unauthorized to update this media":
			message = middleware.Translate(c, "media.unauthorized_access")
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

func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	mediaIDStr := c.Param("id")
	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid media ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.mediaService.DeleteMedia(c.Request.Context(), userID, mediaID)
	if err != nil {
		var message string
		switch err.Error() {
		case "media not found":
			message = middleware.Translate(c, "media.not_found")
		case "unauthorized to delete this media":
			message = middleware.Translate(c, "media.unauthorized_access")
		default:
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "media.file_deleted"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

func (h *MediaHandler) GetAllMedia(c *gin.Context) {
	var req dto.MediaListRequest

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

	if mediaType := c.Query("media_type"); mediaType != "" {
		req.MediaType = mediaType
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			req.UserID = &userID
		}
	}

	// For admin endpoint, we'll use GetUserMedia with the specified user ID
	// or implement a separate admin service method
	var result *dto.MediaListResponse
	var err error

	if req.UserID != nil {
		result, err = h.mediaService.GetUserMedia(c.Request.Context(), *req.UserID, &req)
	} else {
		// TODO: Implement GetAllMedia method in service for admin
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.not_found"),
			"Admin media listing not implemented",
		)
		c.JSON(http.StatusNotImplemented, response)
		return
	}

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
