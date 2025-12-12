package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
)

type MediaRepository interface {
	Create(ctx context.Context, media *domain.Media) error
	GetByID(ctx context.Context, id int) (*domain.Media, error)
	ExistsByID(ctx context.Context, id int) (bool, error)
	GetByUserID(ctx context.Context, userID int, limit, offset int) ([]*domain.Media, error)
	CountByUserID(ctx context.Context, userID int) (int, error)
	Update(ctx context.Context, media *domain.Media) error
	Delete(ctx context.Context, id int) error
	GetByFileName(ctx context.Context, fileName string) (*domain.Media, error)
	GetByMediaType(ctx context.Context, mediaType domain.MediaType, limit, offset int) ([]*domain.Media, error)
	GetUnconvertedMedia(ctx context.Context, limit int) ([]*domain.Media, error)
}
