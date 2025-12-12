package postgres

import (
	"context"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

// Basic read operations
func (r *categoryRepository) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).
		Preload("Icon").
		Preload("Parent").
		Where("id = ?", id).
		First(&category).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).
		Preload("Icon").
		Preload("Parent").
		Where("slug = ?", slug).
		First(&category).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

// Tree operations
func (r *categoryRepository) GetTree(ctx context.Context) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.WithContext(ctx).
		Preload("Icon").
		Order("lft ASC").
		Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) GetTreeByType(ctx context.Context, categoryType domain.CategoryType) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.WithContext(ctx).
		Preload("Icon").
		Where("type = ?", categoryType).
		Order("lft ASC").
		Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) GetChildren(ctx context.Context, parentID int) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.WithContext(ctx).
		Preload("Icon").
		Where("parent_id = ?", parentID).
		Order("lft ASC").
		Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) GetParents(ctx context.Context, categoryID int) ([]*domain.Category, error) {
	// Get the category first to get its lft and rgt values
	category, err := r.GetByID(ctx, categoryID)
	if err != nil || category == nil {
		return nil, err
	}

	var parents []*domain.Category
	err = r.db.WithContext(ctx).
		Preload("Icon").
		Where("lft < ? AND rgt > ?", category.Lft, category.Rgt).
		Order("lft ASC").
		Find(&parents).Error

	return parents, err
}

func (r *categoryRepository) GetSiblings(ctx context.Context, categoryID int) ([]*domain.Category, error) {
	// Get the category first to get its parent_id
	category, err := r.GetByID(ctx, categoryID)
	if err != nil || category == nil {
		return nil, err
	}

	var siblings []*domain.Category
	query := r.db.WithContext(ctx).
		Preload("Icon").
		Where("id != ?", categoryID)

	if category.ParentID != nil {
		query = query.Where("parent_id = ?", *category.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	err = query.Order("lft ASC").Find(&siblings).Error
	return siblings, err
}

func (r *categoryRepository) GetDescendants(ctx context.Context, categoryID int) ([]*domain.Category, error) {
	// Get the category first to get its lft and rgt values
	category, err := r.GetByID(ctx, categoryID)
	if err != nil || category == nil {
		return nil, err
	}

	var descendants []*domain.Category
	err = r.db.WithContext(ctx).
		Preload("Icon").
		Where("lft > ? AND rgt < ?", category.Lft, category.Rgt).
		Order("lft ASC").
		Find(&descendants).Error

	return descendants, err
}

func (r *categoryRepository) GetAncestors(ctx context.Context, categoryID int) ([]*domain.Category, error) {
	return r.GetParents(ctx, categoryID)
}

// Utility operations
func (r *categoryRepository) GetByType(ctx context.Context, categoryType domain.CategoryType) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.WithContext(ctx).
		Preload("Icon").
		Where("type = ?", categoryType).
		Order("lft ASC").
		Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) GetRootCategories(ctx context.Context) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.WithContext(ctx).
		Preload("Icon").
		Where("parent_id IS NULL").
		Order("lft ASC").
		Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) GetLeafCategories(ctx context.Context) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.WithContext(ctx).
		Preload("Icon").
		Where("rgt - lft = 1").
		Order("lft ASC").
		Find(&categories).Error

	return categories, err
}

// Search operations
func (r *categoryRepository) Search(ctx context.Context, query string) ([]*domain.Category, error) {
	var categories []*domain.Category
	searchPattern := "%" + query + "%"

	err := r.db.WithContext(ctx).
		Preload("Icon").
		Where("name ILIKE ? OR slug ILIKE ?", searchPattern, searchPattern).
		Order("lft ASC").
		Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) SearchByType(ctx context.Context, query string, categoryType domain.CategoryType) ([]*domain.Category, error) {
	var categories []*domain.Category
	searchPattern := "%" + query + "%"

	err := r.db.WithContext(ctx).
		Preload("Icon").
		Where("type = ? AND (name ILIKE ? OR slug ILIKE ?)", categoryType, searchPattern, searchPattern).
		Order("lft ASC").
		Find(&categories).Error

	return categories, err
}

// Count operations
func (r *categoryRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Category{}).Count(&count).Error
	return int(count), err
}

func (r *categoryRepository) CountByType(ctx context.Context, categoryType domain.CategoryType) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Category{}).Where("type = ?", categoryType).Count(&count).Error
	return int(count), err
}

func (r *categoryRepository) CountChildren(ctx context.Context, parentID int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Category{}).Where("parent_id = ?", parentID).Count(&count).Error
	return int(count), err
}
