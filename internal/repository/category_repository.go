package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
)

type CategoryRepository interface {
	// Basic read operations
	GetByID(ctx context.Context, id int) (*domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Category, error)

	// Tree operations
	GetTree(ctx context.Context) ([]*domain.Category, error)
	GetTreeByType(ctx context.Context, categoryType domain.CategoryType) ([]*domain.Category, error)
	GetChildren(ctx context.Context, parentID int) ([]*domain.Category, error)
	GetParents(ctx context.Context, categoryID int) ([]*domain.Category, error)
	GetSiblings(ctx context.Context, categoryID int) ([]*domain.Category, error)
	GetDescendants(ctx context.Context, categoryID int) ([]*domain.Category, error)
	GetAncestors(ctx context.Context, categoryID int) ([]*domain.Category, error)

	// Utility operations
	GetByType(ctx context.Context, categoryType domain.CategoryType) ([]*domain.Category, error)
	GetRootCategories(ctx context.Context) ([]*domain.Category, error)
	GetLeafCategories(ctx context.Context) ([]*domain.Category, error)

	// Search operations
	Search(ctx context.Context, query string) ([]*domain.Category, error)
	SearchByType(ctx context.Context, query string, categoryType domain.CategoryType) ([]*domain.Category, error)

	// Count operations
	Count(ctx context.Context) (int, error)
	CountByType(ctx context.Context, categoryType domain.CategoryType) (int, error)
	CountChildren(ctx context.Context, parentID int) (int, error)
}
