package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type CategoryType string

const (
	CategoryTypeConcertsFestivals  CategoryType = "concerts_&_festivals"
	CategoryTypeParty              CategoryType = "party"
	CategoryTypeCulture            CategoryType = "culture"
	CategoryTypeShows              CategoryType = "shows"
	CategoryTypeSports             CategoryType = "sports"
	CategoryTypeFreetimeActivities CategoryType = "freetime_activities"
	CategoryTypeBusiness           CategoryType = "business"
	CategoryTypeEthnic             CategoryType = "ethnic"
	CategoryTypeOther              CategoryType = "other"
)

type Category struct {
	ID     int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string       `json:"name" gorm:"type:varchar(200);not null"`
	IconID *int         `json:"icon_id" gorm:"index"`
	Type   CategoryType `json:"type" gorm:"type:varchar(50);not null"`
	Slug   string       `json:"slug" gorm:"type:varchar(250);uniqueIndex;not null"`

	// Nested Set Model fields
	ParentID *int `json:"parent_id" gorm:"index"`
	Lft      int  `json:"lft" gorm:"not null;index"`
	Rgt      int  `json:"rgt" gorm:"not null;index"`
	Depth    int  `json:"depth" gorm:"not null;default:0"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Icon     *Media      `json:"icon,omitempty" gorm:"foreignKey:IconID;references:ID"`
	Parent   *Category   `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Children []*Category `json:"children,omitempty" gorm:"foreignKey:ParentID;references:ID"`
}

func NewCategory(name string, categoryType CategoryType, parentID *int) *Category {
	category := &Category{
		Name:      name,
		Type:      categoryType,
		ParentID:  parentID,
		Depth:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	category.GenerateSlug()
	return category
}

func (c *Category) GenerateSlug() {
	if c.Name == "" {
		return
	}

	// Convert to lowercase
	slug := strings.ToLower(c.Name)

	// Replace Turkish characters
	replacements := map[string]string{
		"ç": "c", "ğ": "g", "ı": "i", "ö": "o", "ş": "s", "ü": "u",
		"Ç": "c", "Ğ": "g", "İ": "i", "Ö": "o", "Ş": "s", "Ü": "u",
	}

	for turkish, english := range replacements {
		slug = strings.ReplaceAll(slug, turkish, english)
	}

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Ensure slug is not empty
	if slug == "" {
		slug = fmt.Sprintf("category-%d", time.Now().Unix())
	}

	c.Slug = slug
}

func (c *Category) SetIcon(iconID int) {
	c.IconID = &iconID
	c.UpdatedAt = time.Now()
}

func (c *Category) RemoveIcon() {
	c.IconID = nil
	c.UpdatedAt = time.Now()
}

func (c *Category) UpdateInfo(name string, categoryType CategoryType) {
	if name != "" && name != c.Name {
		c.Name = name
		c.GenerateSlug()
	}
	if categoryType != "" {
		c.Type = categoryType
	}
	c.UpdatedAt = time.Now()
}

func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}

func (c *Category) IsLeaf() bool {
	return c.Rgt-c.Lft == 1
}

func (c *Category) HasChildren() bool {
	return !c.IsLeaf()
}

func (c *Category) GetChildrenCount() int {
	return (c.Rgt - c.Lft - 1) / 2
}

func (c *Category) IsAncestorOf(other *Category) bool {
	return c.Lft < other.Lft && c.Rgt > other.Rgt
}

func (c *Category) IsDescendantOf(other *Category) bool {
	return other.IsAncestorOf(c)
}

func (c *Category) ValidateRequiredFields() error {
	if c.Name == "" {
		return ErrCategoryNameRequired
	}
	if c.Type == "" {
		return ErrCategoryTypeRequired
	}
	if c.Slug == "" {
		return ErrCategorySlugRequired
	}
	return nil
}

func (c *Category) ValidateType() error {
	validTypes := []CategoryType{
		CategoryTypeConcertsFestivals,
		CategoryTypeParty,
		CategoryTypeCulture,
		CategoryTypeShows,
		CategoryTypeSports,
		CategoryTypeFreetimeActivities,
		CategoryTypeBusiness,
		CategoryTypeEthnic,
		CategoryTypeOther,
	}

	for _, validType := range validTypes {
		if c.Type == validType {
			return nil
		}
	}

	return ErrCategoryInvalidType
}

// Category domain errors
var (
	ErrCategoryNameRequired = NewDomainError("category name is required")
	ErrCategoryTypeRequired = NewDomainError("category type is required")
	ErrCategorySlugRequired = NewDomainError("category slug is required")
	ErrCategoryInvalidType  = NewDomainError("invalid category type")
	ErrCategoryNotFound     = NewDomainError("category not found")
	ErrCategoryHasChildren  = NewDomainError("category has children and cannot be deleted")
	ErrCategoryCircularRef  = NewDomainError("circular reference detected in category hierarchy")
)
