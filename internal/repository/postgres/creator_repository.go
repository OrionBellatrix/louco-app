package postgres

import (
	"context"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"gorm.io/gorm"
)

type creatorRepository struct {
	db *gorm.DB
}

func NewCreatorRepository(db *gorm.DB) repository.CreatorRepository {
	return &creatorRepository{
		db: db,
	}
}

func (r *creatorRepository) Create(ctx context.Context, creator *domain.Creator) error {
	return r.db.WithContext(ctx).Create(creator).Error
}

func (r *creatorRepository) GetByID(ctx context.Context, id int) (*domain.Creator, error) {
	var creator domain.Creator

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Industries").
		Where("id = ?", id).
		First(&creator).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &creator, nil
}

func (r *creatorRepository) GetByUserID(ctx context.Context, userID int) (*domain.Creator, error) {
	var creator domain.Creator

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Industries").
		Where("user_id = ?", userID).
		First(&creator).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &creator, nil
}

func (r *creatorRepository) Update(ctx context.Context, creator *domain.Creator) error {
	return r.db.WithContext(ctx).Save(creator).Error
}

func (r *creatorRepository) Delete(ctx context.Context, id int) error {
	// First remove all industry associations
	if err := r.db.WithContext(ctx).Exec("DELETE FROM creator_industries WHERE creator_id = ?", id).Error; err != nil {
		return err
	}

	// Then delete the creator
	return r.db.WithContext(ctx).Delete(&domain.Creator{}, id).Error
}

func (r *creatorRepository) List(ctx context.Context, limit, offset int) ([]*domain.Creator, error) {
	var creators []*domain.Creator

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Industries").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&creators).Error

	return creators, err
}

func (r *creatorRepository) ListByIndustryID(ctx context.Context, industryID, limit, offset int) ([]*domain.Creator, error) {
	var creators []*domain.Creator

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Industries").
		Joins("JOIN creator_industries ON creators.id = creator_industries.creator_id").
		Where("creator_industries.industry_id = ?", industryID).
		Limit(limit).
		Offset(offset).
		Order("creators.created_at DESC").
		Find(&creators).Error

	return creators, err
}

func (r *creatorRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Creator{}).Count(&count).Error
	return int(count), err
}

func (r *creatorRepository) CountByIndustryID(ctx context.Context, industryID int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Creator{}).
		Joins("JOIN creator_industries ON creators.id = creator_industries.creator_id").
		Where("creator_industries.industry_id = ?", industryID).
		Count(&count).Error
	return int(count), err
}

func (r *creatorRepository) AddIndustries(ctx context.Context, creatorID int, industryIDs []int) error {
	for _, industryID := range industryIDs {
		// Use raw SQL to handle conflicts properly
		if err := r.db.WithContext(ctx).
			Exec("INSERT INTO creator_industries (creator_id, industry_id, created_at) VALUES (?, ?, NOW()) ON CONFLICT (creator_id, industry_id) DO NOTHING",
				creatorID, industryID).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *creatorRepository) RemoveIndustries(ctx context.Context, creatorID int, industryIDs []int) error {
	return r.db.WithContext(ctx).
		Where("creator_id = ? AND industry_id IN ?", creatorID, industryIDs).
		Delete(&domain.CreatorIndustry{}).Error
}

func (r *creatorRepository) SetIndustries(ctx context.Context, creatorID int, industryIDs []int) error {
	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Remove all existing associations
	if err := tx.Where("creator_id = ?", creatorID).Delete(&domain.CreatorIndustry{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add new associations
	for _, industryID := range industryIDs {
		creatorIndustry := &domain.CreatorIndustry{
			CreatorID:  creatorID,
			IndustryID: industryID,
		}
		if err := tx.Create(creatorIndustry).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *creatorRepository) GetIndustries(ctx context.Context, creatorID int) ([]*domain.Industry, error) {
	var industries []*domain.Industry

	err := r.db.WithContext(ctx).
		Model(&domain.Industry{}).
		Joins("JOIN creator_industries ON industries.id = creator_industries.industry_id").
		Where("creator_industries.creator_id = ?", creatorID).
		Order("industries.name ASC").
		Find(&industries).Error

	return industries, err
}
