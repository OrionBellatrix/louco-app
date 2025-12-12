package dto

import (
	"time"

	"github.com/louco-event/internal/domain"
)

// CategoryResponse represents a single category in API responses
type CategoryResponse struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Icon      *MediaResponse      `json:"icon,omitempty"`
	Type      string              `json:"type"`
	Slug      string              `json:"slug"`
	ParentID  *int                `json:"parent_id,omitempty"`
	Depth     int                 `json:"depth"`
	Children  []*CategoryResponse `json:"children,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// CategoryTreeResponse represents the complete category tree
type CategoryTreeResponse struct {
	Categories []*CategoryResponse `json:"categories"`
	Total      int                 `json:"total"`
}

// CategoryListRequest represents request parameters for category listing
type CategoryListRequest struct {
	Type   string `form:"type" json:"type"`
	Search string `form:"search" json:"search"`
	Depth  *int   `form:"depth" json:"depth"`
}

// CategorySearchRequest represents search parameters
type CategorySearchRequest struct {
	Query string `form:"q" json:"query" binding:"required,min=2"`
	Type  string `form:"type" json:"type"`
}

// CategoryByTypeRequest represents request for categories by type
type CategoryByTypeRequest struct {
	Type string `uri:"type" json:"type" binding:"required"`
}

// CategoryBySlugRequest represents request for category by slug
type CategoryBySlugRequest struct {
	Slug string `uri:"slug" json:"slug" binding:"required"`
}

// CategoryByIDRequest represents request for category by ID
type CategoryByIDRequest struct {
	ID int `uri:"id" json:"id" binding:"required,min=1"`
}

// Helper function to convert domain.Category to CategoryResponse
func CategoryToResponse(category *domain.Category) *CategoryResponse {
	if category == nil {
		return nil
	}

	response := &CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Type:      string(category.Type),
		Slug:      category.Slug,
		ParentID:  category.ParentID,
		Depth:     category.Depth,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	// Add icon if exists
	if category.Icon != nil {
		response.Icon = &MediaResponse{
			ID:           category.Icon.ID,
			UserID:       category.Icon.UserID,
			OriginalName: category.Icon.OriginalName,
			FileName:     category.Icon.FileName,
			FileURL:      category.Icon.FileURL,
			MediaType:    string(category.Icon.MediaType),
			MimeType:     category.Icon.MimeType,
			FileSize:     category.Icon.FileSize,
			Width:        category.Icon.Width,
			Height:       category.Icon.Height,
			Duration:     category.Icon.Duration,
			IsConverted:  category.Icon.IsConverted,
			CreatedAt:    category.Icon.CreatedAt,
			UpdatedAt:    category.Icon.UpdatedAt,
		}
	}

	return response
}

// Helper function to convert slice of domain.Category to slice of CategoryResponse
func CategoriesToResponse(categories []*domain.Category) []*CategoryResponse {
	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = CategoryToResponse(category)
	}
	return responses
}

// Helper function to build nested tree structure from flat list
func BuildCategoryTree(categories []*domain.Category) []*CategoryResponse {
	if len(categories) == 0 {
		return []*CategoryResponse{}
	}

	// Create a map for quick lookup
	categoryMap := make(map[int]*CategoryResponse)
	var rootCategories []*CategoryResponse

	// First pass: create all category responses
	for _, category := range categories {
		response := CategoryToResponse(category)
		categoryMap[category.ID] = response

		if category.ParentID == nil {
			rootCategories = append(rootCategories, response)
		}
	}

	// Second pass: build parent-child relationships
	for _, category := range categories {
		if category.ParentID != nil {
			parent := categoryMap[*category.ParentID]
			child := categoryMap[category.ID]
			if parent != nil && child != nil {
				if parent.Children == nil {
					parent.Children = []*CategoryResponse{}
				}
				parent.Children = append(parent.Children, child)
			}
		}
	}

	return rootCategories
}

// Helper function to flatten tree structure (for search results)
func FlattenCategoryTree(treeCategories []*CategoryResponse) []*CategoryResponse {
	var flattened []*CategoryResponse

	var flatten func([]*CategoryResponse)
	flatten = func(categories []*CategoryResponse) {
		for _, category := range categories {
			// Create a copy without children for flat representation
			flatCategory := &CategoryResponse{
				ID:        category.ID,
				Name:      category.Name,
				Icon:      category.Icon,
				Type:      category.Type,
				Slug:      category.Slug,
				ParentID:  category.ParentID,
				Depth:     category.Depth,
				CreatedAt: category.CreatedAt,
				UpdatedAt: category.UpdatedAt,
			}
			flattened = append(flattened, flatCategory)

			if category.Children != nil {
				flatten(category.Children)
			}
		}
	}

	flatten(treeCategories)
	return flattened
}
