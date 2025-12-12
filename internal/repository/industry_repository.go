package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
)

type IndustryRepository interface {
	GetAll(ctx context.Context) ([]*domain.Industry, error)
	GetByID(ctx context.Context, id int) (*domain.Industry, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Industry, error)
}
