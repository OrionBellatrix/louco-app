package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
)

type CreatorRepository interface {
	Create(ctx context.Context, creator *domain.Creator) error
	GetByID(ctx context.Context, id int) (*domain.Creator, error)
	GetByUserID(ctx context.Context, userID int) (*domain.Creator, error)
	Update(ctx context.Context, creator *domain.Creator) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*domain.Creator, error)
	ListByIndustryID(ctx context.Context, industryID, limit, offset int) ([]*domain.Creator, error)
	Count(ctx context.Context) (int, error)
	CountByIndustryID(ctx context.Context, industryID int) (int, error)
	AddIndustries(ctx context.Context, creatorID int, industryIDs []int) error
	RemoveIndustries(ctx context.Context, creatorID int, industryIDs []int) error
	SetIndustries(ctx context.Context, creatorID int, industryIDs []int) error
	GetIndustries(ctx context.Context, creatorID int) ([]*domain.Industry, error)
}
