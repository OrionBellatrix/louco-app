package service

import (
	"context"
	"fmt"
	"time"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
	"github.com/louco-event/pkg/cache"
	"github.com/louco-event/pkg/logger"
)

const (
	CategoryTreeCacheKey    = "categories_nested_tree"
	CategoryTypeCacheKey    = "categories_type_%s"
	CategoryCacheExpiration = 10 * time.Minute
)

type CategoryService interface {
	// Tree operations
	GetCategoryTree(ctx context.Context) (*dto.CategoryTreeResponse, error)
	GetCategoryTreeByType(ctx context.Context, categoryType domain.CategoryType) (*dto.CategoryTreeResponse, error)

	// Single category operations
	GetCategoryByID(ctx context.Context, id int) (*dto.CategoryResponse, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*dto.CategoryResponse, error)

	// Hierarchy operations
	GetChildren(ctx context.Context, parentID int) ([]*dto.CategoryResponse, error)
	GetParents(ctx context.Context, categoryID int) ([]*dto.CategoryResponse, error)
	GetSiblings(ctx context.Context, categoryID int) ([]*dto.CategoryResponse, error)
	GetDescendants(ctx context.Context, categoryID int) ([]*dto.CategoryResponse, error)
	GetAncestors(ctx context.Context, categoryID int) ([]*dto.CategoryResponse, error)

	// Utility operations
	GetRootCategories(ctx context.Context) ([]*dto.CategoryResponse, error)
	GetLeafCategories(ctx context.Context) ([]*dto.CategoryResponse, error)
	GetCategoriesByType(ctx context.Context, categoryType domain.CategoryType) ([]*dto.CategoryResponse, error)

	// Search operations
	SearchCategories(ctx context.Context, query string) ([]*dto.CategoryResponse, error)
	SearchCategoriesByType(ctx context.Context, query string, categoryType domain.CategoryType) ([]*dto.CategoryResponse, error)

	// Count operations
	GetCategoryCount(ctx context.Context) (int, error)
	GetCategoryCountByType(ctx context.Context, categoryType domain.CategoryType) (int, error)
	GetChildrenCount(ctx context.Context, parentID int) (int, error)

	// Cache operations
	RefreshCache(ctx context.Context) error
	ClearCache(ctx context.Context) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	mediaRepo    repository.MediaRepository
	cache        *cache.RedisCache
	logger       *logger.Logger
}

func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	mediaRepo repository.MediaRepository,
	cache *cache.RedisCache,
	logger *logger.Logger,
) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		mediaRepo:    mediaRepo,
		cache:        cache,
		logger:       logger,
	}
}

// Tree operations
func (s *categoryService) GetCategoryTree(ctx context.Context) (*dto.CategoryTreeResponse, error) {
	// Try to get from cache first
	var cachedTree dto.CategoryTreeResponse
	if err := s.cache.Get(ctx, CategoryTreeCacheKey, &cachedTree); err == nil {
		s.logger.Debug().Msg("Category tree retrieved from cache")
		return &cachedTree, nil
	}

	// Get from database
	categories, err := s.categoryRepo.GetTree(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get category tree from database")
		return nil, fmt.Errorf("failed to get category tree")
	}

	// Build tree structure
	treeCategories := dto.BuildCategoryTree(categories)

	// Load media for categories with icons
	for _, category := range treeCategories {
		s.loadCategoryMedia(ctx, category)
	}

	response := &dto.CategoryTreeResponse{
		Categories: treeCategories,
		Total:      len(categories),
	}

	// Cache the result
	if err := s.cache.Set(ctx, CategoryTreeCacheKey, response, CategoryCacheExpiration); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to cache category tree")
	}

	return response, nil
}

func (s *categoryService) GetCategoryTreeByType(ctx context.Context, categoryType domain.CategoryType) (*dto.CategoryTreeResponse, error) {
	cacheKey := fmt.Sprintf(CategoryTypeCacheKey, string(categoryType))

	// Try to get from cache first
	var cachedTree dto.CategoryTreeResponse
	if err := s.cache.Get(ctx, cacheKey, &cachedTree); err == nil {
		s.logger.Debug().Str("type", string(categoryType)).Msg("Category tree by type retrieved from cache")
		return &cachedTree, nil
	}

	// Get from database
	categories, err := s.categoryRepo.GetTreeByType(ctx, categoryType)
	if err != nil {
		s.logger.Error().Err(err).Str("type", string(categoryType)).Msg("Failed to get category tree by type from database")
		return nil, fmt.Errorf("failed to get category tree by type")
	}

	// Build tree structure
	treeCategories := dto.BuildCategoryTree(categories)

	// Load media for categories with icons
	for _, category := range treeCategories {
		s.loadCategoryMedia(ctx, category)
	}

	response := &dto.CategoryTreeResponse{
		Categories: treeCategories,
		Total:      len(categories),
	}

	// Cache the result
	if err := s.cache.Set(ctx, cacheKey, response, CategoryCacheExpiration); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to cache category tree by type")
	}

	return response, nil
}

// Single category operations
func (s *categoryService) GetCategoryByID(ctx context.Context, id int) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("id", id).Msg("Failed to get category by ID")
		return nil, fmt.Errorf("failed to get category")
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	response := dto.CategoryToResponse(category)
	s.loadCategoryMedia(ctx, response)

	return response, nil
}

func (s *categoryService) GetCategoryBySlug(ctx context.Context, slug string) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		s.logger.Error().Err(err).Str("slug", slug).Msg("Failed to get category by slug")
		return nil, fmt.Errorf("failed to get category")
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	response := dto.CategoryToResponse(category)
	s.loadCategoryMedia(ctx, response)

	return response, nil
}

// Hierarchy operations
func (s *categoryService) GetChildren(ctx context.Context, parentID int) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetChildren(ctx, parentID)
	if err != nil {
		s.logger.Error().Err(err).Int("parent_id", parentID).Msg("Failed to get children")
		return nil, fmt.Errorf("failed to get children")
	}

	responses := dto.CategoriesToResponse(categories)
	for _, response := range responses {
		s.loadCategoryMedia(ctx, response)
	}

	return responses, nil
}

func (s *categoryService) GetParents(ctx context.Context, categoryID int) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetParents(ctx, categoryID)
	if err != nil {
		s.logger.Error().Err(err).Int("category_id", categoryID).Msg("Failed to get parents")
		return nil, fmt.Errorf("failed to get parents")
	}

	responses := dto.CategoriesToResponse(categories)
	for _, response := range responses {
		s.loadCategoryMedia(ctx, response)
	}

	return responses, nil
}

func (s *categoryService) GetSiblings(ctx context.Context, categoryID int) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetSiblings(ctx, categoryID)
	if err != nil {
		s.logger.Error().Err(err).Int("category_id", categoryID).Msg("Failed to get siblings")
		return nil, fmt.Errorf("failed to get siblings")
	}

	responses := dto.CategoriesToResponse(categories)
	for _, response := range responses {
		s.loadCategoryMedia(ctx, response)
	}

	return responses, nil
}

func (s *categoryService) GetDescendants(ctx context.Context, categoryID int) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetDescendants(ctx, categoryID)
	if err != nil {
		s.logger.Error().Err(err).Int("category_id", categoryID).Msg("Failed to get descendants")
		return nil, fmt.Errorf("failed to get descendants")
	}

	responses := dto.CategoriesToResponse(categories)
	for _, response := range responses {
		s.loadCategoryMedia(ctx, response)
	}

	return responses, nil
}

func (s *categoryService) GetAncestors(ctx context.Context, categoryID int) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetAncestors(ctx, categoryID)
	if err != nil {
		s.logger.Error().Err(err).Int("category_id", categoryID).Msg("Failed to get ancestors")
		return nil, fmt.Errorf("failed to get ancestors")
	}

	responses := dto.CategoriesToResponse(categories)
	for _, response := range responses {
		s.loadCategoryMedia(ctx, response)
	}

	return responses, nil
}

// Utility operations
func (s *categoryService) GetRootCategories(ctx context.Context) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get root categories")
		return nil, fmt.Errorf("failed to get root categories")
	}

	responses := dto.CategoriesToResponse(categories)
	for _, response := range responses {
		s.loadCategoryMedia(ctx, response)
	}

	return responses, nil
}

func (s *categoryService) GetLeafCategories(ctx context.Context) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetLeafCategories(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get leaf categories")
		return nil, fmt.Errorf("failed to get leaf categories")
	}

	responses := dto.CategoriesToResponse(categories)
	for _, response := range responses {
		s.loadCategoryMedia(ctx, response)
	}

	return responses, nil
}

func (s *categoryService) GetCategoriesByType(ctx context.Context, categoryType domain.CategoryType) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetByType(ctx, categoryType)
	if err != nil {
		s.logger.Error().Err(err).Str("type", string(categoryType)).Msg("Failed to get categories by type")
		return nil, fmt.Errorf("failed to get categories by type")
	}

	responses := dto.CategoriesToResponse(categories)
	for _, response := range responses {
		s.loadCategoryMedia(ctx, response)
	}

	return responses, nil
}

// Search operations
func (s *categoryService) SearchCategories(ctx context.Context, query string) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.Search(ctx, query)
	if err != nil {
		s.logger.Error().Err(err).Str("query", query).Msg("Failed to search categories")
		return nil, fmt.Errorf("failed to search categories")
	}

	responses := dto.CategoriesToResponse(categories)
	for _, response := range responses {
		s.loadCategoryMedia(ctx, response)
	}

	return responses, nil
}

func (s *categoryService) SearchCategoriesByType(ctx context.Context, query string, categoryType domain.CategoryType) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.SearchByType(ctx, query, categoryType)
	if err != nil {
		s.logger.Error().Err(err).Str("query", query).Str("type", string(categoryType)).Msg("Failed to search categories by type")
		return nil, fmt.Errorf("failed to search categories by type")
	}

	responses := dto.CategoriesToResponse(categories)
	for _, response := range responses {
		s.loadCategoryMedia(ctx, response)
	}

	return responses, nil
}

// Count operations
func (s *categoryService) GetCategoryCount(ctx context.Context) (int, error) {
	count, err := s.categoryRepo.Count(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get category count")
		return 0, fmt.Errorf("failed to get category count")
	}
	return count, nil
}

func (s *categoryService) GetCategoryCountByType(ctx context.Context, categoryType domain.CategoryType) (int, error) {
	count, err := s.categoryRepo.CountByType(ctx, categoryType)
	if err != nil {
		s.logger.Error().Err(err).Str("type", string(categoryType)).Msg("Failed to get category count by type")
		return 0, fmt.Errorf("failed to get category count by type")
	}
	return count, nil
}

func (s *categoryService) GetChildrenCount(ctx context.Context, parentID int) (int, error) {
	count, err := s.categoryRepo.CountChildren(ctx, parentID)
	if err != nil {
		s.logger.Error().Err(err).Int("parent_id", parentID).Msg("Failed to get children count")
		return 0, fmt.Errorf("failed to get children count")
	}
	return count, nil
}

// Cache operations
func (s *categoryService) RefreshCache(ctx context.Context) error {
	// Clear existing cache
	if err := s.ClearCache(ctx); err != nil {
		return err
	}

	// Refresh main tree cache
	if _, err := s.GetCategoryTree(ctx); err != nil {
		return err
	}

	// Refresh type-specific caches
	categoryTypes := []domain.CategoryType{
		domain.CategoryTypeConcertsFestivals,
		domain.CategoryTypeParty,
		domain.CategoryTypeCulture,
		domain.CategoryTypeShows,
		domain.CategoryTypeSports,
		domain.CategoryTypeFreetimeActivities,
		domain.CategoryTypeBusiness,
		domain.CategoryTypeEthnic,
		domain.CategoryTypeOther,
	}

	for _, categoryType := range categoryTypes {
		if _, err := s.GetCategoryTreeByType(ctx, categoryType); err != nil {
			s.logger.Warn().Err(err).Str("type", string(categoryType)).Msg("Failed to refresh cache for category type")
		}
	}

	s.logger.Info().Msg("Category cache refreshed successfully")
	return nil
}

func (s *categoryService) ClearCache(ctx context.Context) error {
	// Clear main tree cache
	if err := s.cache.Delete(ctx, CategoryTreeCacheKey); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to clear main category tree cache")
	}

	// Clear type-specific caches
	categoryTypes := []domain.CategoryType{
		domain.CategoryTypeConcertsFestivals,
		domain.CategoryTypeParty,
		domain.CategoryTypeCulture,
		domain.CategoryTypeShows,
		domain.CategoryTypeSports,
		domain.CategoryTypeFreetimeActivities,
		domain.CategoryTypeBusiness,
		domain.CategoryTypeEthnic,
		domain.CategoryTypeOther,
	}

	for _, categoryType := range categoryTypes {
		cacheKey := fmt.Sprintf(CategoryTypeCacheKey, string(categoryType))
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn().Err(err).Str("type", string(categoryType)).Msg("Failed to clear category type cache")
		}
	}

	s.logger.Info().Msg("Category cache cleared successfully")
	return nil
}

// Helper function to load media for category icons
func (s *categoryService) loadCategoryMedia(ctx context.Context, category *dto.CategoryResponse) {
	if category == nil {
		return
	}

	// Load media for current category
	if category.Icon != nil && category.Icon.ID > 0 {
		media, err := s.mediaRepo.GetByID(ctx, category.Icon.ID)
		if err == nil && media != nil {
			category.Icon = &dto.MediaResponse{
				ID:           media.ID,
				UserID:       media.UserID,
				OriginalName: media.OriginalName,
				FileName:     media.FileName,
				FileURL:      media.FileURL,
				MediaType:    string(media.MediaType),
				MimeType:     media.MimeType,
				FileSize:     media.FileSize,
				Width:        media.Width,
				Height:       media.Height,
				Duration:     media.Duration,
				IsConverted:  media.IsConverted,
				CreatedAt:    media.CreatedAt,
				UpdatedAt:    media.UpdatedAt,
			}
		}
	}

	// Recursively load media for children
	if category.Children != nil {
		for _, child := range category.Children {
			s.loadCategoryMedia(ctx, child)
		}
	}
}
