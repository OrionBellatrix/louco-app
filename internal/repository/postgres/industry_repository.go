package postgres

import (
	"context"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"gorm.io/gorm"
)

type industryRepository struct {
	db *gorm.DB
}

func NewIndustryRepository(db *gorm.DB) repository.IndustryRepository {
	return &industryRepository{
		db: db,
	}
}

func (r *industryRepository) GetAll(ctx context.Context) ([]*domain.Industry, error) {
	var industries []*domain.Industry

	if err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&industries).Error; err != nil {
		return nil, err
	}

	return industries, nil
}

func (r *industryRepository) GetByID(ctx context.Context, id int) (*domain.Industry, error) {
	var industry domain.Industry

	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&industry).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &industry, nil
}

func (r *industryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Industry, error) {
	var industry domain.Industry

	if err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&industry).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &industry, nil
}
