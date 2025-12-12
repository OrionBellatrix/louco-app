package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
)

type VerificationRepository interface {
	Create(ctx context.Context, code *domain.VerificationCode) error
	GetByIdentifierAndType(ctx context.Context, identifier string, codeType domain.VerificationType) (*domain.VerificationCode, error)
	UpdateAttempts(ctx context.Context, id uint, attempts int) error
	MarkAsUsed(ctx context.Context, id uint) error
	DeleteExpired(ctx context.Context) error
	GetActiveByIdentifier(ctx context.Context, identifier string, codeType domain.VerificationType) (*domain.VerificationCode, error)
}
