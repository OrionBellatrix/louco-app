package service

import (
	"context"
	"fmt"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
)

// FollowService handles follow-related business logic
type FollowService struct {
	followRepo repository.FollowRepository
	userRepo   repository.UserRepository
}

// NewFollowService creates a new follow service instance
func NewFollowService(followRepo repository.FollowRepository, userRepo repository.UserRepository) *FollowService {
	return &FollowService{
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

// Follow creates a follow relationship between two users
func (s *FollowService) Follow(ctx context.Context, followerID, followingID int) error {
	// Check if users exist
	follower, err := s.userRepo.GetByID(ctx, followerID)
	if err != nil {
		return fmt.Errorf("failed to get follower: %w", err)
	}
	if follower == nil {
		return domain.NewDomainError("user.not.found")
	}

	following, err := s.userRepo.GetByID(ctx, followingID)
	if err != nil {
		return fmt.Errorf("failed to get following user: %w", err)
	}
	if following == nil {
		return domain.NewDomainError("user.not.found")
	}

	// Create follow relationship
	follow, err := domain.NewFollow(followerID, followingID)
	if err != nil {
		return err
	}

	// Check if already following
	exists, err := s.followRepo.Exists(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to check follow existence: %w", err)
	}
	if exists {
		return domain.NewDomainError("follow.already.exists")
	}

	// Create the follow relationship
	if err := s.followRepo.Create(ctx, follow); err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}

	// Update follow counts
	if err := s.updateFollowCounts(ctx, followerID, followingID, true); err != nil {
		return fmt.Errorf("failed to update follow counts: %w", err)
	}

	return nil
}

// Unfollow removes a follow relationship between two users
func (s *FollowService) Unfollow(ctx context.Context, followerID, followingID int) error {
	// Validate that users cannot unfollow themselves
	if followerID == followingID {
		return domain.NewDomainError("follow.cannot.unfollow.self")
	}

	// Check if follow relationship exists
	exists, err := s.followRepo.Exists(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to check follow existence: %w", err)
	}
	if !exists {
		return domain.NewDomainError("follow.not.found")
	}

	// Delete the follow relationship
	if err := s.followRepo.Delete(ctx, followerID, followingID); err != nil {
		return fmt.Errorf("failed to delete follow: %w", err)
	}

	// Update follow counts
	if err := s.updateFollowCounts(ctx, followerID, followingID, false); err != nil {
		return fmt.Errorf("failed to update follow counts: %w", err)
	}

	return nil
}

// IsFollowing checks if one user follows another
func (s *FollowService) IsFollowing(ctx context.Context, followerID, followingID int) (bool, error) {
	return s.followRepo.Exists(ctx, followerID, followingID)
}

// GetFollowers gets users who follow the specified user with pagination
func (s *FollowService) GetFollowers(ctx context.Context, userID int, limit, offset int) ([]*domain.Follow, int, error) {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, 0, domain.NewDomainError("user.not.found")
	}

	// Get followers with user info
	followers, err := s.followRepo.GetFollowersWithUserInfo(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get followers: %w", err)
	}

	// Get total count
	total, err := s.followRepo.CountFollowers(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count followers: %w", err)
	}

	return followers, total, nil
}

// GetFollowing gets users that the specified user follows with pagination
func (s *FollowService) GetFollowing(ctx context.Context, userID int, limit, offset int) ([]*domain.Follow, int, error) {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, 0, domain.NewDomainError("user.not.found")
	}

	// Get following with user info
	following, err := s.followRepo.GetFollowingWithUserInfo(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get following: %w", err)
	}

	// Get total count
	total, err := s.followRepo.CountFollowing(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count following: %w", err)
	}

	return following, total, nil
}

// GetMutualFollows gets mutual follows between two users
func (s *FollowService) GetMutualFollows(ctx context.Context, userID1, userID2 int, limit, offset int) ([]*domain.Follow, int, error) {
	// Check if users exist
	user1, err := s.userRepo.GetByID(ctx, userID1)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user1: %w", err)
	}
	if user1 == nil {
		return nil, 0, domain.NewDomainError("user.not.found")
	}

	user2, err := s.userRepo.GetByID(ctx, userID2)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user2: %w", err)
	}
	if user2 == nil {
		return nil, 0, domain.NewDomainError("user.not.found")
	}

	// Get mutual follows
	mutualFollows, err := s.followRepo.GetMutualFollows(ctx, userID1, userID2, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get mutual follows: %w", err)
	}

	// Get total count
	total, err := s.followRepo.CountMutualFollows(ctx, userID1, userID2)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count mutual follows: %w", err)
	}

	return mutualFollows, total, nil
}

// GetFollowCounts gets follow counts for a user
func (s *FollowService) GetFollowCounts(ctx context.Context, userID int) (int, int, error) {
	followersCount, err := s.followRepo.CountFollowers(ctx, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count followers: %w", err)
	}

	followingCount, err := s.followRepo.CountFollowing(ctx, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count following: %w", err)
	}

	return followersCount, followingCount, nil
}

// updateFollowCounts updates the cached follow counts in the users table
func (s *FollowService) updateFollowCounts(ctx context.Context, followerID, followingID int, isFollow bool) error {
	// Get current users
	follower, err := s.userRepo.GetByID(ctx, followerID)
	if err != nil {
		return err
	}

	following, err := s.userRepo.GetByID(ctx, followingID)
	if err != nil {
		return err
	}

	// Update counts
	if isFollow {
		follower.IncrementFollowingCount()
		following.IncrementFollowersCount()
	} else {
		follower.DecrementFollowingCount()
		following.DecrementFollowersCount()
	}

	// Save updated counts
	if err := s.userRepo.Update(ctx, follower); err != nil {
		return err
	}

	if err := s.userRepo.Update(ctx, following); err != nil {
		return err
	}

	return nil
}
