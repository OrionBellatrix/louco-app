package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
)

// FollowRepository defines the interface for follow operations
type FollowRepository interface {
	// Follow operations
	Create(ctx context.Context, follow *domain.Follow) error
	Delete(ctx context.Context, followerID, followingID int) error

	// Check operations
	Exists(ctx context.Context, followerID, followingID int) (bool, error)
	GetByFollowerAndFollowing(ctx context.Context, followerID, followingID int) (*domain.Follow, error)

	// List operations with pagination
	GetFollowers(ctx context.Context, userID int, limit, offset int) ([]*domain.Follow, error)
	GetFollowing(ctx context.Context, userID int, limit, offset int) ([]*domain.Follow, error)

	// Count operations
	CountFollowers(ctx context.Context, userID int) (int, error)
	CountFollowing(ctx context.Context, userID int) (int, error)

	// Batch operations for user profile
	GetFollowersWithUserInfo(ctx context.Context, userID int, limit, offset int) ([]*domain.Follow, error)
	GetFollowingWithUserInfo(ctx context.Context, userID int, limit, offset int) ([]*domain.Follow, error)

	// Mutual follows
	GetMutualFollows(ctx context.Context, userID1, userID2 int, limit, offset int) ([]*domain.Follow, error)
	CountMutualFollows(ctx context.Context, userID1, userID2 int) (int, error)
}
