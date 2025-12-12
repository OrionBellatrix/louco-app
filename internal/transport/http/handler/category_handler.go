package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/middleware"
	"github.com/louco-event/internal/service"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetCategoryTree godoc
// @Summary Get complete category tree
// @Description Get the complete nested category tree structure
// @Tags categories
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.CategoryTreeResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /categories/tree [get]
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	tree, err := h.categoryService.GetCategoryTree(c.Request.Context())
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.tree_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.tree_success"),
		tree,
	)
	c.JSON(http.StatusOK, response)
}

// GetCategoryTreeByType godoc
// @Summary Get category tree by type
// @Description Get nested category tree structure filtered by category type
// @Tags categories
// @Produce json
// @Param type path string true "Category Type" Enums(concerts_&_festivals,party,culture,shows,sports,freetime_activities,business,ethnic,other)
// @Success 200 {object} dto.APIResponse{data=dto.CategoryTreeResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /categories/tree/type/{type} [get]
func (h *CategoryHandler) GetCategoryTreeByType(c *gin.Context) {
	var req dto.CategoryByTypeRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	categoryType := domain.CategoryType(req.Type)
	if err := (&domain.Category{Type: categoryType}).ValidateType(); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.invalid_type"),
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	tree, err := h.categoryService.GetCategoryTreeByType(c.Request.Context(), categoryType)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.tree_by_type_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.tree_by_type_success"),
		tree,
	)
	c.JSON(http.StatusOK, response)
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a specific category by its ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.APIResponse{data=dto.CategoryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	var req dto.CategoryByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	category, err := h.categoryService.GetCategoryByID(c.Request.Context(), req.ID)
	if err != nil {
		var message string
		var statusCode int

		if err.Error() == "category not found" {
			message = middleware.Translate(c, "category.not_found")
			statusCode = http.StatusNotFound
		} else {
			message = middleware.Translate(c, "category.get_failed")
			statusCode = http.StatusInternalServerError
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(statusCode, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.get_success"),
		category,
	)
	c.JSON(http.StatusOK, response)
}

// GetCategoryBySlug godoc
// @Summary Get category by slug
// @Description Get a specific category by its slug
// @Tags categories
// @Produce json
// @Param slug path string true "Category Slug"
// @Success 200 {object} dto.APIResponse{data=dto.CategoryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	var req dto.CategoryBySlugRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	category, err := h.categoryService.GetCategoryBySlug(c.Request.Context(), req.Slug)
	if err != nil {
		var message string
		var statusCode int

		if err.Error() == "category not found" {
			message = middleware.Translate(c, "category.not_found")
			statusCode = http.StatusNotFound
		} else {
			message = middleware.Translate(c, "category.get_failed")
			statusCode = http.StatusInternalServerError
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(statusCode, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.get_success"),
		category,
	)
	c.JSON(http.StatusOK, response)
}

// GetCategoryChildren godoc
// @Summary Get category children
// @Description Get direct children of a specific category
// @Tags categories
// @Produce json
// @Param id path int true "Parent Category ID"
// @Success 200 {object} dto.APIResponse{data=[]dto.CategoryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetCategoryChildren(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid category ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	children, err := h.categoryService.GetChildren(c.Request.Context(), id)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.children_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.children_success"),
		children,
	)
	c.JSON(http.StatusOK, response)
}

// GetCategoryParents godoc
// @Summary Get category parents
// @Description Get all parent categories (ancestors) of a specific category
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.APIResponse{data=[]dto.CategoryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /categories/{id}/parents [get]
func (h *CategoryHandler) GetCategoryParents(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid category ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	parents, err := h.categoryService.GetParents(c.Request.Context(), id)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.parents_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.parents_success"),
		parents,
	)
	c.JSON(http.StatusOK, response)
}

// GetRootCategories godoc
// @Summary Get root categories
// @Description Get all root level categories (categories without parent)
// @Tags categories
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]dto.CategoryResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /categories/roots [get]
func (h *CategoryHandler) GetRootCategories(c *gin.Context) {
	categories, err := h.categoryService.GetRootCategories(c.Request.Context())
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.roots_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.roots_success"),
		categories,
	)
	c.JSON(http.StatusOK, response)
}

// GetLeafCategories godoc
// @Summary Get leaf categories
// @Description Get all leaf categories (categories without children)
// @Tags categories
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]dto.CategoryResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /categories/leaves [get]
func (h *CategoryHandler) GetLeafCategories(c *gin.Context) {
	categories, err := h.categoryService.GetLeafCategories(c.Request.Context())
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.leaves_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.leaves_success"),
		categories,
	)
	c.JSON(http.StatusOK, response)
}

// GetCategoriesByType godoc
// @Summary Get categories by type
// @Description Get all categories of a specific type (flat list)
// @Tags categories
// @Produce json
// @Param type path string true "Category Type" Enums(concerts_&_festivals,party,culture,shows,sports,freetime_activities,business,ethnic,other)
// @Success 200 {object} dto.APIResponse{data=[]dto.CategoryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /categories/type/{type} [get]
func (h *CategoryHandler) GetCategoriesByType(c *gin.Context) {
	var req dto.CategoryByTypeRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	categoryType := domain.CategoryType(req.Type)
	if err := (&domain.Category{Type: categoryType}).ValidateType(); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.invalid_type"),
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	categories, err := h.categoryService.GetCategoriesByType(c.Request.Context(), categoryType)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.by_type_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.by_type_success"),
		categories,
	)
	c.JSON(http.StatusOK, response)
}

// SearchCategories godoc
// @Summary Search categories
// @Description Search categories by name or slug
// @Tags categories
// @Produce json
// @Param q query string true "Search query" minlength(2)
// @Param type query string false "Category type filter" Enums(concerts_&_festivals,party,culture,shows,sports,freetime_activities,business,ethnic,other)
// @Success 200 {object} dto.APIResponse{data=[]dto.CategoryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /categories/search [get]
func (h *CategoryHandler) SearchCategories(c *gin.Context) {
	var req dto.CategorySearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var categories []*dto.CategoryResponse
	var err error

	if req.Type != "" {
		categoryType := domain.CategoryType(req.Type)
		if err := (&domain.Category{Type: categoryType}).ValidateType(); err != nil {
			response := dto.NewErrorResponse(
				middleware.Translate(c, "category.invalid_type"),
				nil,
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		categories, err = h.categoryService.SearchCategoriesByType(c.Request.Context(), req.Query, categoryType)
	} else {
		categories, err = h.categoryService.SearchCategories(c.Request.Context(), req.Query)
	}

	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.search_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.search_success"),
		categories,
	)
	c.JSON(http.StatusOK, response)
}

// RefreshCache godoc
// @Summary Refresh category cache
// @Description Refresh the Redis cache for category tree
// @Tags categories
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /categories/cache/refresh [post]
func (h *CategoryHandler) RefreshCache(c *gin.Context) {
	err := h.categoryService.RefreshCache(c.Request.Context())
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.cache_refresh_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.cache_refresh_success"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

// ClearCache godoc
// @Summary Clear category cache
// @Description Clear the Redis cache for category tree
// @Tags categories
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /categories/cache/clear [delete]
func (h *CategoryHandler) ClearCache(c *gin.Context) {
	err := h.categoryService.ClearCache(c.Request.Context())
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "category.cache_clear_failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "category.cache_clear_success"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}
