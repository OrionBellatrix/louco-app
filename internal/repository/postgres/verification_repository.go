package postgres

import (
	"context"
	"time"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"gorm.io/gorm"
)

type verificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) repository.VerificationRepository {
	return &verificationRepository{
		db: db,
	}
}

func (r *verificationRepository) Create(ctx context.Context, code *domain.VerificationCode) error {
	return r.db.WithContext(ctx).Create(code).Error
}

func (r *verificationRepository) GetByIdentifierAndType(ctx context.Context, identifier string, codeType domain.VerificationType) (*domain.VerificationCode, error) {
	var code domain.VerificationCode
	err := r.db.WithContext(ctx).
		Where("identifier = ? AND type = ?", identifier, codeType).
		Order("created_at DESC").
		First(&code).Error

	if err != nil {
		return nil, err
	}

	return &code, nil
}

func (r *verificationRepository) GetActiveByIdentifier(ctx context.Context, identifier string, codeType domain.VerificationType) (*domain.VerificationCode, error) {
	var code domain.VerificationCode
	err := r.db.WithContext(ctx).
		Where("identifier = ? AND type = ? AND used_at IS NULL AND expires_at > ?",
			identifier, codeType, time.Now()).
		Order("created_at DESC").
		First(&code).Error

	if err != nil {
		return nil, err
	}

	return &code, nil
}

func (r *verificationRepository) UpdateAttempts(ctx context.Context, id uint, attempts int) error {
	return r.db.WithContext(ctx).
		Model(&domain.VerificationCode{}).
		Where("id = ?", id).
		Update("attempts", attempts).Error
}

func (r *verificationRepository) MarkAsUsed(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.VerificationCode{}).
		Where("id = ?", id).
		Update("used_at", &now).Error
}

func (r *verificationRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&domain.VerificationCode{}).Error
}
