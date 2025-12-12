package repository

import (
	"context"
	"time"

	"github.com/louco-event/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
	IsEmailExists(ctx context.Context, email string, excludeUserID *int) (bool, error)
	IsPhoneExists(ctx context.Context, phone string, excludeUserID *int) (bool, error)
	IsUsernameExists(ctx context.Context, username string, excludeUserID *int) (bool, error)
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
	GetByAppleID(ctx context.Context, appleID string) (*domain.User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error)

	// Verification methods
	SetEmailVerified(ctx context.Context, userID int, verifiedAt time.Time) error
	SetPhoneVerified(ctx context.Context, userID int, verifiedAt time.Time) error

	// Follow count methods
	IncrementFollowersCount(ctx context.Context, userID int) error
	DecrementFollowersCount(ctx context.Context, userID int) error
	IncrementFollowingCount(ctx context.Context, userID int) error
	DecrementFollowingCount(ctx context.Context, userID int) error
}
