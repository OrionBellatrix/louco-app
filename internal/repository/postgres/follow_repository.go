package postgres

import (
	"context"
	"errors"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"gorm.io/gorm"
)

type followRepository struct {
	db *gorm.DB
}

// NewFollowRepository creates a new follow repository instance
func NewFollowRepository(db *gorm.DB) repository.FollowRepository {
	return &followRepository{
		db: db,
	}
}

// Create creates a new follow relationship
func (r *followRepository) Create(ctx context.Context, follow *domain.Follow) error {
	if err := r.db.WithContext(ctx).Create(follow).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a follow relationship
func (r *followRepository) Delete(ctx context.Context, followerID, followingID int) error {
	result := r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&domain.Follow{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("follow relationship not found")
	}

	return nil
}

// Exists checks if a follow relationship exists
func (r *followRepository) Exists(ctx context.Context, followerID, followingID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetByFollowerAndFollowing gets a specific follow relationship
func (r *followRepository) GetByFollowerAndFollowing(ctx context.Context, followerID, followingID int) (*domain.Follow, error) {
	var follow domain.Follow
	err := r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		First(&follow).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &follow, nil
}

// GetFollowers gets users who follow the specified user
func (r *followRepository) GetFollowers(ctx context.Context, userID int, limit, offset int) ([]*domain.Follow, error) {
	var follows []*domain.Follow
	err := r.db.WithContext(ctx).
		Where("following_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&follows).Error

	if err != nil {
		return nil, err
	}

	return follows, nil
}

// GetFollowing gets users that the specified user follows
func (r *followRepository) GetFollowing(ctx context.Context, userID int, limit, offset int) ([]*domain.Follow, error) {
	var follows []*domain.Follow
	err := r.db.WithContext(ctx).
		Where("follower_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&follows).Error

	if err != nil {
		return nil, err
	}

	return follows, nil
}

// CountFollowers counts how many users follow the specified user
func (r *followRepository) CountFollowers(ctx context.Context, userID int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Follow{}).
		Where("following_id = ?", userID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// CountFollowing counts how many users the specified user follows
func (r *followRepository) CountFollowing(ctx context.Context, userID int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Follow{}).
		Where("follower_id = ?", userID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetFollowersWithUserInfo gets followers with user information preloaded
func (r *followRepository) GetFollowersWithUserInfo(ctx context.Context, userID int, limit, offset int) ([]*domain.Follow, error) {
	var follows []*domain.Follow
	err := r.db.WithContext(ctx).
		Preload("Follower").
		Preload("Follower.ProfilePicture").
		Where("following_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&follows).Error

	if err != nil {
		return nil, err
	}

	return follows, nil
}

// GetFollowingWithUserInfo gets following with user information preloaded
func (r *followRepository) GetFollowingWithUserInfo(ctx context.Context, userID int, limit, offset int) ([]*domain.Follow, error) {
	var follows []*domain.Follow
	err := r.db.WithContext(ctx).
		Preload("Following").
		Preload("Following.ProfilePicture").
		Where("follower_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&follows).Error

	if err != nil {
		return nil, err
	}

	return follows, nil
}

// GetMutualFollows gets mutual follows between two users
func (r *followRepository) GetMutualFollows(ctx context.Context, userID1, userID2 int, limit, offset int) ([]*domain.Follow, error) {
	var follows []*domain.Follow

	// Find users that both userID1 and userID2 follow
	err := r.db.WithContext(ctx).
		Preload("Following").
		Preload("Following.ProfilePicture").
		Where("follower_id = ? AND following_id IN (?)", userID1,
			r.db.Select("following_id").
				Where("follower_id = ?", userID2).
				Table("follows")).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&follows).Error

	if err != nil {
		return nil, err
	}

	return follows, nil
}

// CountMutualFollows counts mutual follows between two users
func (r *followRepository) CountMutualFollows(ctx context.Context, userID1, userID2 int) (int, error) {
	var count int64

	// Count users that both userID1 and userID2 follow
	err := r.db.WithContext(ctx).
		Model(&domain.Follow{}).
		Where("follower_id = ? AND following_id IN (?)", userID1,
			r.db.Select("following_id").
				Where("follower_id = ?", userID2).
				Table("follows")).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}
