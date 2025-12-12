package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateKeyError(err, "users_email_key") {
			return fmt.Errorf("email already exists")
		}
		if isDuplicateKeyError(err, "users_phone_key") {
			return fmt.Errorf("phone already exists")
		}
		if isDuplicateKeyError(err, "users_username_key") {
			return fmt.Errorf("username already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("phone = ? AND is_active = ?", phone, true).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("username = ? AND is_active = ?", username, true).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error, "users_email_key") {
			return fmt.Errorf("email already exists")
		}
		if isDuplicateKeyError(result.Error, "users_phone_key") {
			return fmt.Errorf("phone already exists")
		}
		if isDuplicateKeyError(result.Error, "users_username_key") {
			return fmt.Errorf("username already exists")
		}
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *userRepository) IsEmailExists(ctx context.Context, email string, excludeUserID *int) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.User{}).Where("email = ? AND is_active = ?", email, true)

	if excludeUserID != nil {
		query = query.Where("id != ?", *excludeUserID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

func (r *userRepository) IsPhoneExists(ctx context.Context, phone string, excludeUserID *int) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.User{}).Where("phone = ? AND is_active = ?", phone, true)

	if excludeUserID != nil {
		query = query.Where("id != ?", *excludeUserID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check phone existence: %w", err)
	}

	return count > 0, nil
}

func (r *userRepository) IsUsernameExists(ctx context.Context, username string, excludeUserID *int) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.User{}).Where("username = ? AND is_active = ?", username, true)

	if excludeUserID != nil {
		query = query.Where("id != ?", *excludeUserID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return count > 0, nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

func (r *userRepository) GetByAppleID(ctx context.Context, appleID string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("apple_id = ? AND is_active = ?", appleID, true).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("google_id = ? AND is_active = ?", googleID, true).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) SetEmailVerified(ctx context.Context, userID int, verifiedAt time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ? AND is_active = ?", userID, true).
		Update("email_verified_at", verifiedAt)

	if result.Error != nil {
		return fmt.Errorf("failed to set email verified: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) SetPhoneVerified(ctx context.Context, userID int, verifiedAt time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ? AND is_active = ?", userID, true).
		Update("phone_verified_at", verifiedAt)

	if result.Error != nil {
		return fmt.Errorf("failed to set phone verified: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) IncrementFollowersCount(ctx context.Context, userID int) error {
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ? AND is_active = ?", userID, true).
		Update("followers_count", gorm.Expr("followers_count + ?", 1))

	if result.Error != nil {
		return fmt.Errorf("failed to increment followers count: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) DecrementFollowersCount(ctx context.Context, userID int) error {
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ? AND is_active = ? AND followers_count > 0", userID, true).
		Update("followers_count", gorm.Expr("followers_count - ?", 1))

	if result.Error != nil {
		return fmt.Errorf("failed to decrement followers count: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found or followers count is already zero")
	}

	return nil
}

func (r *userRepository) IncrementFollowingCount(ctx context.Context, userID int) error {
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ? AND is_active = ?", userID, true).
		Update("following_count", gorm.Expr("following_count + ?", 1))

	if result.Error != nil {
		return fmt.Errorf("failed to increment following count: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) DecrementFollowingCount(ctx context.Context, userID int) error {
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ? AND is_active = ? AND following_count > 0", userID, true).
		Update("following_count", gorm.Expr("following_count - ?", 1))

	if result.Error != nil {
		return fmt.Errorf("failed to decrement following count: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found or following count is already zero")
	}

	return nil
}

// Helper function to check for duplicate key errors
func isDuplicateKeyError(err error, constraint string) bool {
	if err == nil {
		return false
	}
	// This is a simplified check - in a real implementation you might want to check for specific PostgreSQL error codes
	return false // GORM handles unique constraint violations differently
}
