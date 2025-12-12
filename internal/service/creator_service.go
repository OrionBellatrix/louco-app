package service

import (
	"context"
	"fmt"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
	"github.com/louco-event/pkg/logger"
)

type CreatorService interface {
	CreateCreator(ctx context.Context, userID int, req *dto.CreateCreatorRequest) (*dto.CreatorResponse, error)
	GetCreatorByID(ctx context.Context, id int) (*dto.CreatorResponse, error)
	GetCreatorByUserID(ctx context.Context, userID int) (*dto.CreatorResponse, error)
	UpdateCreator(ctx context.Context, userID int, req *dto.UpdateCreatorRequest) error
	SetWeeztixToken(ctx context.Context, userID int, req *dto.SetWeeztixTokenRequest) error
	GetCreatorProfile(ctx context.Context, userID int) (*dto.CreatorProfileResponse, error)
	GetCreatorList(ctx context.Context, req *dto.CreatorListRequest) (*dto.CreatorListResponse, error)
	DeleteCreator(ctx context.Context, userID int) error
}

type creatorService struct {
	creatorRepo  repository.CreatorRepository
	userRepo     repository.UserRepository
	industryRepo repository.IndustryRepository
	mediaRepo    repository.MediaRepository
	logger       *logger.Logger
}

func NewCreatorService(
	creatorRepo repository.CreatorRepository,
	userRepo repository.UserRepository,
	industryRepo repository.IndustryRepository,
	mediaRepo repository.MediaRepository,
	logger *logger.Logger,
) CreatorService {
	return &creatorService{
		creatorRepo:  creatorRepo,
		userRepo:     userRepo,
		industryRepo: industryRepo,
		mediaRepo:    mediaRepo,
		logger:       logger,
	}
}

func (s *creatorService) CreateCreator(ctx context.Context, userID int, req *dto.CreateCreatorRequest) (*dto.CreatorResponse, error) {
	// Check if user exists and is creator type
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get user")
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsCreator() {
		return nil, fmt.Errorf("user is not a creator type")
	}

	// Check if creator profile already exists
	existingCreator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to check existing creator")
		return nil, fmt.Errorf("failed to check existing creator profile")
	}
	if existingCreator != nil {
		return nil, fmt.Errorf("creator profile already exists")
	}

	// Validate industries exist
	for _, industryID := range req.IndustryIDs {
		industry, err := s.industryRepo.GetByID(ctx, industryID)
		if err != nil {
			s.logger.Error().Err(err).Int("industry_id", industryID).Msg("Failed to validate industry")
			return nil, fmt.Errorf("failed to validate industry")
		}
		if industry == nil {
			return nil, fmt.Errorf("industry with ID %d not found", industryID)
		}
	}

	// Create creator
	creator := domain.NewCreator(
		userID,
		req.CompanyName,
		req.Address,
		req.EstimatedTickets,
		req.EstimatedEvents,
	)

	// Validate industry IDs are provided
	if len(req.IndustryIDs) == 0 {
		return nil, fmt.Errorf("at least one industry selection is required")
	}

	if err := creator.ValidateRequiredFields(); err != nil {
		return nil, err
	}

	if err := s.creatorRepo.Create(ctx, creator); err != nil {
		s.logger.Error().Err(err).Msg("Failed to create creator")
		return nil, fmt.Errorf("failed to create creator profile")
	}

	// Set industries
	if err := s.creatorRepo.SetIndustries(ctx, creator.ID, req.IndustryIDs); err != nil {
		s.logger.Error().Err(err).Int("creator_id", creator.ID).Msg("Failed to set industries")
		return nil, fmt.Errorf("failed to set industries")
	}

	// Get creator with relations
	createdCreator, err := s.creatorRepo.GetByID(ctx, creator.ID)
	if err != nil {
		s.logger.Error().Err(err).Int("creator_id", creator.ID).Msg("Failed to get created creator")
		return nil, fmt.Errorf("failed to get created creator")
	}

	return s.mapCreatorToResponse(createdCreator), nil
}

func (s *creatorService) GetCreatorByID(ctx context.Context, id int) (*dto.CreatorResponse, error) {
	creator, err := s.creatorRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("creator_id", id).Msg("Failed to get creator")
		return nil, fmt.Errorf("failed to get creator")
	}
	if creator == nil {
		return nil, fmt.Errorf("creator not found")
	}

	return s.mapCreatorToResponse(creator), nil
}

func (s *creatorService) GetCreatorByUserID(ctx context.Context, userID int) (*dto.CreatorResponse, error) {
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator by user ID")
		return nil, fmt.Errorf("failed to get creator")
	}
	if creator == nil {
		return nil, fmt.Errorf("creator profile not found")
	}

	return s.mapCreatorToResponse(creator), nil
}

func (s *creatorService) UpdateCreator(ctx context.Context, userID int, req *dto.UpdateCreatorRequest) error {
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator")
		return fmt.Errorf("creator not found")
	}
	if creator == nil {
		return fmt.Errorf("creator profile not found")
	}

	// Validate industries if provided
	if len(req.IndustryIDs) > 0 {
		for _, industryID := range req.IndustryIDs {
			industry, err := s.industryRepo.GetByID(ctx, industryID)
			if err != nil {
				s.logger.Error().Err(err).Int("industry_id", industryID).Msg("Failed to validate industry")
				return fmt.Errorf("failed to validate industry")
			}
			if industry == nil {
				return fmt.Errorf("industry with ID %d not found", industryID)
			}
		}
	}

	// Update creator fields
	creator.UpdateProfile(req.CompanyName, req.Address, req.EstimatedTickets, req.EstimatedEvents)

	if err := creator.ValidateRequiredFields(); err != nil {
		return err
	}

	if err := s.creatorRepo.Update(ctx, creator); err != nil {
		s.logger.Error().Err(err).Int("creator_id", creator.ID).Msg("Failed to update creator")
		return fmt.Errorf("failed to update creator profile")
	}

	// Update industries if provided
	if len(req.IndustryIDs) > 0 {
		if err := s.creatorRepo.SetIndustries(ctx, creator.ID, req.IndustryIDs); err != nil {
			s.logger.Error().Err(err).Int("creator_id", creator.ID).Msg("Failed to update industries")
			return fmt.Errorf("failed to update industries")
		}
	}

	return nil
}

func (s *creatorService) SetWeeztixToken(ctx context.Context, userID int, req *dto.SetWeeztixTokenRequest) error {
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator")
		return fmt.Errorf("creator not found")
	}
	if creator == nil {
		return fmt.Errorf("creator profile not found")
	}

	creator.SetWeeztixToken(req.WeeztixToken)

	if err := s.creatorRepo.Update(ctx, creator); err != nil {
		s.logger.Error().Err(err).Int("creator_id", creator.ID).Msg("Failed to update weeztix token")
		return fmt.Errorf("failed to update weeztix token")
	}

	return nil
}

func (s *creatorService) GetCreatorProfile(ctx context.Context, userID int) (*dto.CreatorProfileResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get user")
		return nil, fmt.Errorf("user not found")
	}

	// Get creator
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator")
		return nil, fmt.Errorf("failed to get creator")
	}
	if creator == nil {
		return nil, fmt.Errorf("creator profile not found")
	}

	userResponse := s.mapUserToResponse(user)

	// Get profile picture if exists
	if user.ProfilePicID != nil {
		profilePic, err := s.mediaRepo.GetByID(ctx, *user.ProfilePicID)
		if err == nil {
			profilePicResponse := s.mapMediaToResponse(profilePic)
			userResponse.ProfilePic = &profilePicResponse
		}
	}

	// Get cover picture if exists
	if user.CoverPicID != nil {
		coverPic, err := s.mediaRepo.GetByID(ctx, *user.CoverPicID)
		if err == nil {
			coverPicResponse := s.mapMediaToResponse(coverPic)
			userResponse.CoverPic = &coverPicResponse
		}
	}

	return &dto.CreatorProfileResponse{
		User:    userResponse,
		Creator: *s.mapCreatorToResponse(creator),
	}, nil
}

func (s *creatorService) GetCreatorList(ctx context.Context, req *dto.CreatorListRequest) (*dto.CreatorListResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var creators []*domain.Creator
	var total int
	var err error

	if req.IndustryID > 0 {
		creators, err = s.creatorRepo.ListByIndustryID(ctx, req.IndustryID, pageSize, offset)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to get creators by industry")
			return nil, fmt.Errorf("failed to get creators")
		}
		total, err = s.creatorRepo.CountByIndustryID(ctx, req.IndustryID)
	} else {
		creators, err = s.creatorRepo.List(ctx, pageSize, offset)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to get creators")
			return nil, fmt.Errorf("failed to get creators")
		}
		total, err = s.creatorRepo.Count(ctx)
	}

	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get creator count")
		return nil, fmt.Errorf("failed to get creator count")
	}

	creatorResponses := make([]dto.CreatorResponse, len(creators))
	for i, creator := range creators {
		creatorResponses[i] = *s.mapCreatorToResponse(creator)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &dto.CreatorListResponse{
		Creators:   creatorResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *creatorService) DeleteCreator(ctx context.Context, userID int) error {
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator")
		return fmt.Errorf("creator not found")
	}
	if creator == nil {
		return fmt.Errorf("creator profile not found")
	}

	if err := s.creatorRepo.Delete(ctx, creator.ID); err != nil {
		s.logger.Error().Err(err).Int("creator_id", creator.ID).Msg("Failed to delete creator")
		return fmt.Errorf("failed to delete creator profile")
	}

	return nil
}

func (s *creatorService) mapCreatorToResponse(creator *domain.Creator) *dto.CreatorResponse {
	industries := make([]dto.IndustryResponse, len(creator.Industries))
	for i, industry := range creator.Industries {
		industries[i] = dto.IndustryResponse{
			ID:   industry.ID,
			Name: industry.Name,
			Slug: industry.Slug,
		}
	}

	return &dto.CreatorResponse{
		ID:               creator.ID,
		UserID:           creator.UserID,
		WeeztixToken:     creator.WeeztixToken,
		CompanyName:      creator.CompanyName,
		Address:          creator.Address,
		EstimatedTickets: creator.EstimatedTickets,
		EstimatedEvents:  creator.EstimatedEvents,
		Industries:       industries,
		CreatedAt:        creator.CreatedAt,
		UpdatedAt:        creator.UpdatedAt,
	}
}

func (s *creatorService) mapUserToResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:           user.ID,
		FullName:     user.FullName,
		Username:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		UserType:     string(user.UserType),
		Biography:    user.Biography,
		BirthDate:    user.BirthDate,
		ProfilePicID: user.ProfilePicID,
		CoverPicID:   user.CoverPicID,
		IsActive:     user.IsActive,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func (s *creatorService) mapMediaToResponse(media *domain.Media) dto.MediaResponse {
	return dto.MediaResponse{
		ID:           media.ID,
		UserID:       media.UserID,
		OriginalName: media.OriginalName,
		FileName:     media.FileName,
		FileURL:      media.FileURL,
		MediaType:    string(media.MediaType),
		MimeType:     media.MimeType,
		FileSize:     media.FileSize,
		Width:        media.Width,
		Height:       media.Height,
		Duration:     media.Duration,
		IsConverted:  media.IsConverted,
		CreatedAt:    media.CreatedAt,
		UpdatedAt:    media.UpdatedAt,
	}
}
