package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/middleware"
	"github.com/louco-event/internal/service"
	"github.com/rs/zerolog"
)

type IndustryHandler struct {
	industryService service.IndustryService
	logger          zerolog.Logger
}

func NewIndustryHandler(industryService service.IndustryService, logger zerolog.Logger) *IndustryHandler {
	return &IndustryHandler{
		industryService: industryService,
		logger:          logger,
	}
}

// GetAllIndustries godoc
// @Summary Get all industries
// @Description Get list of all available industries
// @Tags Industries
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.IndustriesResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /industries [get]
func (h *IndustryHandler) GetAllIndustries(c *gin.Context) {
	ctx := c.Request.Context()

	industries, err := h.industryService.GetAllIndustries(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get all industries")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "industry.get_all.failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Convert domain entities to DTOs
	industryDTOs := make([]*dto.IndustryResponse, len(industries))
	for i, industry := range industries {
		industryDTOs[i] = h.domainToDTO(industry)
	}

	response := &dto.IndustriesResponse{
		Industries: industryDTOs,
		Total:      len(industryDTOs),
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: middleware.Translate(c, "industry.get_all.success"),
		Data:    response,
		Errors:  nil,
	})
}

// GetIndustryByID godoc
// @Summary Get industry by ID
// @Description Get a specific industry by its ID
// @Tags Industries
// @Accept json
// @Produce json
// @Param id path int true "Industry ID"
// @Success 200 {object} dto.APIResponse{data=dto.IndustryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /industries/{id} [get]
func (h *IndustryHandler) GetIndustryByID(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("Invalid industry ID")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "industry.invalid_id"),
			Data:    nil,
			Errors:  []string{"Invalid industry ID"},
		})
		return
	}

	industry, err := h.industryService.GetIndustryByID(ctx, id)
	if err != nil {
		h.logger.Error().Err(err).Int("id", id).Msg("Failed to get industry by ID")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "industry.get_by_id.failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	if industry == nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "industry.not_found"),
			Data:    nil,
			Errors:  []string{"Industry not found"},
		})
		return
	}

	response := h.domainToDTO(industry)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: middleware.Translate(c, "industry.get_by_id.success"),
		Data:    response,
		Errors:  nil,
	})
}

// GetIndustryBySlug godoc
// @Summary Get industry by slug
// @Description Get a specific industry by its slug
// @Tags Industries
// @Accept json
// @Produce json
// @Param slug path string true "Industry Slug"
// @Success 200 {object} dto.APIResponse{data=dto.IndustryResponse}
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /industries/slug/{slug} [get]
func (h *IndustryHandler) GetIndustryBySlug(c *gin.Context) {
	ctx := c.Request.Context()

	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "industry.invalid_slug"),
			Data:    nil,
			Errors:  []string{"Slug is required"},
		})
		return
	}

	industry, err := h.industryService.GetIndustryBySlug(ctx, slug)
	if err != nil {
		h.logger.Error().Err(err).Str("slug", slug).Msg("Failed to get industry by slug")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "industry.get_by_slug.failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	if industry == nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: middleware.Translate(c, "industry.not_found"),
			Data:    nil,
			Errors:  []string{"Industry not found"},
		})
		return
	}

	response := h.domainToDTO(industry)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: middleware.Translate(c, "industry.get_by_slug.success"),
		Data:    response,
		Errors:  nil,
	})
}

// domainToDTO converts domain.Industry to dto.IndustryResponse
func (h *IndustryHandler) domainToDTO(industry *domain.Industry) *dto.IndustryResponse {
	return &dto.IndustryResponse{
		ID:   industry.ID,
		Name: industry.Name,
		Slug: industry.Slug,
	}
}
