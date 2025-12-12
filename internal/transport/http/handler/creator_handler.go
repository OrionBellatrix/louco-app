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

type CreatorHandler struct {
	creatorService service.CreatorService
	i18n           *i18n.I18n
}

func NewCreatorHandler(creatorService service.CreatorService, i18n *i18n.I18n) *CreatorHandler {
	return &CreatorHandler{
		creatorService: creatorService,
		i18n:           i18n,
	}
}

// CreateCreator godoc
// @Summary Create creator profile
// @Description Create a new creator profile for authenticated user
// @Tags creators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCreatorRequest true "Creator creation request"
// @Success 201 {object} dto.APIResponse{data=dto.CreatorResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 409 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /creators [post]
func (h *CreatorHandler) CreateCreator(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.CreateCreatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	creator, err := h.creatorService.CreateCreator(c.Request.Context(), userID, &req)
	if err != nil {
		var message string
		var statusCode int

		switch err.Error() {
		case "user is not a creator type":
			message = middleware.Translate(c, "creator.invalid_user_type")
			statusCode = http.StatusBadRequest
		case "creator profile already exists":
			message = middleware.Translate(c, "creator.already_exists")
			statusCode = http.StatusConflict
		default:
			message = middleware.Translate(c, "creator.create_failed")
			statusCode = http.StatusInternalServerError
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(statusCode, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "creator.created"),
		creator,
	)
	c.JSON(http.StatusCreated, response)
}

// GetCreator godoc
// @Summary Get creator by ID
// @Description Get creator profile by creator ID
// @Tags creators
// @Produce json
// @Param id path int true "Creator ID"
// @Success 200 {object} dto.APIResponse{data=dto.CreatorResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /creators/{id} [get]
func (h *CreatorHandler) GetCreator(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid creator ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	creator, err := h.creatorService.GetCreatorByID(c.Request.Context(), id)
	if err != nil {
		var message string
		var statusCode int

		if err.Error() == "creator not found" {
			message = middleware.Translate(c, "creator.not_found")
			statusCode = http.StatusNotFound
		} else {
			message = middleware.Translate(c, "creator.get_failed")
			statusCode = http.StatusInternalServerError
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(statusCode, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "common.success"),
		creator,
	)
	c.JSON(http.StatusOK, response)
}

// GetMyCreatorProfile godoc
// @Summary Get my creator profile
// @Description Get creator profile for authenticated user
// @Tags creators
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.CreatorProfileResponse}
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /creators/me [get]
func (h *CreatorHandler) GetMyCreatorProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	profile, err := h.creatorService.GetCreatorProfile(c.Request.Context(), userID)
	if err != nil {
		var message string
		var statusCode int

		if err.Error() == "creator profile not found" {
			message = middleware.Translate(c, "creator.profile_not_found")
			statusCode = http.StatusNotFound
		} else {
			message = middleware.Translate(c, "creator.get_failed")
			statusCode = http.StatusInternalServerError
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(statusCode, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "common.success"),
		profile,
	)
	c.JSON(http.StatusOK, response)
}

// UpdateCreator godoc
// @Summary Update creator profile
// @Description Update creator profile for authenticated user
// @Tags creators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateCreatorRequest true "Creator update request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /creators/me [put]
func (h *CreatorHandler) UpdateCreator(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.UpdateCreatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.creatorService.UpdateCreator(c.Request.Context(), userID, &req)
	if err != nil {
		var message string
		var statusCode int

		if err.Error() == "creator profile not found" {
			message = middleware.Translate(c, "creator.profile_not_found")
			statusCode = http.StatusNotFound
		} else {
			message = middleware.Translate(c, "creator.update_failed")
			statusCode = http.StatusInternalServerError
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(statusCode, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "creator.updated"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

// SetWeeztixToken godoc
// @Summary Set Weeztix token
// @Description Set or update Weeztix integration token for creator
// @Tags creators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SetWeeztixTokenRequest true "Weeztix token request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /creators/me/weeztix-token [put]
func (h *CreatorHandler) SetWeeztixToken(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.SetWeeztixTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.creatorService.SetWeeztixToken(c.Request.Context(), userID, &req)
	if err != nil {
		var message string
		var statusCode int

		if err.Error() == "creator profile not found" {
			message = middleware.Translate(c, "creator.profile_not_found")
			statusCode = http.StatusNotFound
		} else {
			message = middleware.Translate(c, "creator.token_update_failed")
			statusCode = http.StatusInternalServerError
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(statusCode, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "creator.token_updated"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

// GetCreatorList godoc
// @Summary Get creator list
// @Description Get paginated list of creators with optional industry filter
// @Tags creators
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param industry_id query int false "Industry ID filter"
// @Success 200 {object} dto.APIResponse{data=dto.CreatorListResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /creators [get]
func (h *CreatorHandler) GetCreatorList(c *gin.Context) {
	var req dto.CreatorListRequest

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
	if industryIDStr := c.Query("industry_id"); industryIDStr != "" {
		if industryID, err := strconv.Atoi(industryIDStr); err == nil {
			req.IndustryID = industryID
		}
	}

	creators, err := h.creatorService.GetCreatorList(c.Request.Context(), &req)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "creator.list_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "common.success"),
		creators,
	)
	c.JSON(http.StatusOK, response)
}

// DeleteCreator godoc
// @Summary Delete creator profile
// @Description Delete creator profile for authenticated user
// @Tags creators
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /creators/me [delete]
func (h *CreatorHandler) DeleteCreator(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	err := h.creatorService.DeleteCreator(c.Request.Context(), userID)
	if err != nil {
		var message string
		var statusCode int

		if err.Error() == "creator profile not found" {
			message = middleware.Translate(c, "creator.profile_not_found")
			statusCode = http.StatusNotFound
		} else {
			message = middleware.Translate(c, "creator.delete_failed")
			statusCode = http.StatusInternalServerError
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(statusCode, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "creator.deleted"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}
