package service

import (
	"context"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"github.com/rs/zerolog"
)

type IndustryService interface {
	GetAllIndustries(ctx context.Context) ([]*domain.Industry, error)
	GetIndustryByID(ctx context.Context, id int) (*domain.Industry, error)
	GetIndustryBySlug(ctx context.Context, slug string) (*domain.Industry, error)
}

type industryService struct {
	industryRepo repository.IndustryRepository
	logger       zerolog.Logger
}

func NewIndustryService(industryRepo repository.IndustryRepository, logger zerolog.Logger) IndustryService {
	return &industryService{
		industryRepo: industryRepo,
		logger:       logger,
	}
}

func (s *industryService) GetAllIndustries(ctx context.Context) ([]*domain.Industry, error) {
	s.logger.Info().Msg("Getting all industries")

	industries, err := s.industryRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get all industries")
		return nil, err
	}

	s.logger.Info().Int("count", len(industries)).Msg("Successfully retrieved industries")
	return industries, nil
}

func (s *industryService) GetIndustryByID(ctx context.Context, id int) (*domain.Industry, error) {
	s.logger.Info().Int("industry_id", id).Msg("Getting industry by ID")

	industry, err := s.industryRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("industry_id", id).Msg("Failed to get industry by ID")
		return nil, err
	}

	if industry == nil {
		s.logger.Warn().Int("industry_id", id).Msg("Industry not found")
		return nil, nil
	}

	s.logger.Info().Int("industry_id", id).Str("name", industry.Name).Msg("Successfully retrieved industry")
	return industry, nil
}

func (s *industryService) GetIndustryBySlug(ctx context.Context, slug string) (*domain.Industry, error) {
	s.logger.Info().Str("slug", slug).Msg("Getting industry by slug")

	industry, err := s.industryRepo.GetBySlug(ctx, slug)
	if err != nil {
		s.logger.Error().Err(err).Str("slug", slug).Msg("Failed to get industry by slug")
		return nil, err
	}

	if industry == nil {
		s.logger.Warn().Str("slug", slug).Msg("Industry not found")
		return nil, nil
	}

	s.logger.Info().Str("slug", slug).Str("name", industry.Name).Msg("Successfully retrieved industry")
	return industry, nil
}
