package postgres

import (
	"context"
	"fmt"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"gorm.io/gorm"
)

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) repository.MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(ctx context.Context, media *domain.Media) error {
	if err := r.db.WithContext(ctx).Create(media).Error; err != nil {
		return fmt.Errorf("failed to create media: %w", err)
	}
	return nil
}

func (r *mediaRepository) GetByID(ctx context.Context, id int) (*domain.Media, error) {
	var media domain.Media
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&media).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media not found")
		}
		return nil, fmt.Errorf("failed to get media: %w", err)
	}
	return &media, nil
}

func (r *mediaRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Media{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check media existence: %w", err)
	}
	return count > 0, nil
}

func (r *mediaRepository) GetByUserID(ctx context.Context, userID int, limit, offset int) ([]*domain.Media, error) {
	var mediaList []*domain.Media
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&mediaList).Error; err != nil {
		return nil, fmt.Errorf("failed to get media by user ID: %w", err)
	}
	return mediaList, nil
}

func (r *mediaRepository) CountByUserID(ctx context.Context, userID int) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Media{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count media by user ID: %w", err)
	}
	return int(count), nil
}

func (r *mediaRepository) Update(ctx context.Context, media *domain.Media) error {
	result := r.db.WithContext(ctx).Save(media)
	if result.Error != nil {
		return fmt.Errorf("failed to update media: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("media not found")
	}
	return nil
}

func (r *mediaRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&domain.Media{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete media: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("media not found")
	}
	return nil
}

func (r *mediaRepository) GetByFileName(ctx context.Context, fileName string) (*domain.Media, error) {
	var media domain.Media
	if err := r.db.WithContext(ctx).Where("file_name = ?", fileName).First(&media).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media not found")
		}
		return nil, fmt.Errorf("failed to get media: %w", err)
	}
	return &media, nil
}

func (r *mediaRepository) GetByMediaType(ctx context.Context, mediaType domain.MediaType, limit, offset int) ([]*domain.Media, error) {
	var mediaList []*domain.Media
	if err := r.db.WithContext(ctx).Where("media_type = ?", mediaType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&mediaList).Error; err != nil {
		return nil, fmt.Errorf("failed to get media by type: %w", err)
	}
	return mediaList, nil
}

func (r *mediaRepository) GetUnconvertedMedia(ctx context.Context, limit int) ([]*domain.Media, error) {
	var mediaList []*domain.Media
	if err := r.db.WithContext(ctx).Where("is_converted = ?", false).
		Order("created_at ASC").
		Limit(limit).
		Find(&mediaList).Error; err != nil {
		return nil, fmt.Errorf("failed to get unconverted media: %w", err)
	}
	return mediaList, nil
}
