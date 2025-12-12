package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
	"github.com/louco-event/pkg/logger"
)

type EventService interface {
	// Basic CRUD operations
	CreateEvent(ctx context.Context, userID int, req dto.CreateEventRequest) (*dto.EventResponse, error)
	GetEventByID(ctx context.Context, id int, userID *int) (*dto.EventResponse, error)
	GetEventByIDWithRelations(ctx context.Context, id int, userID *int) (*dto.EventResponse, error)
	UpdateEvent(ctx context.Context, id, userID int, req dto.UpdateEventRequest) (*dto.EventResponse, error)
	DeleteEvent(ctx context.Context, id, userID int) error

	// Creator-specific operations
	GetCreatorEvents(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)
	GetCreatorEventsByStatus(ctx context.Context, userID int, status domain.EventStatus, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)
	GetCreatorDraftEvents(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)
	GetCreatorPublishedEvents(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)

	// Public event operations
	GetPublicEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)
	GetPublicEventsByCategory(ctx context.Context, categoryID int, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)
	GetPublicEventsByLocation(ctx context.Context, city, country string, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)
	SearchPublicEvents(ctx context.Context, req dto.EventSearchRequest, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)

	// Status management
	UpdateEventStatus(ctx context.Context, id, userID int, req dto.UpdateEventStatusRequest) (*dto.EventResponse, error)
	SubmitEventForReview(ctx context.Context, id, userID int) (*dto.EventResponse, error)
	PublishEvent(ctx context.Context, id, userID int) (*dto.EventResponse, error)
	CancelEvent(ctx context.Context, id, userID int) (*dto.EventResponse, error)

	// Advanced filtering and search
	GetEventsWithFilters(ctx context.Context, filters dto.EventFilterRequest, pagination dto.PaginationRequest, userID *int) ([]*dto.EventListResponse, *dto.PaginationResponse, error)
	GetEventsByDateRange(ctx context.Context, startDate, endDate time.Time, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)
	GetUpcomingEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)
	GetPastEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error)

	// Statistics operations
	GetEventStats(ctx context.Context, userID int) (*dto.EventStatsResponse, error)
	GetSystemEventStats(ctx context.Context) (*dto.SystemEventStatsResponse, error)

	// Validation operations
	ValidateEventOwnership(ctx context.Context, eventID, userID int) error
	ValidateEventAccess(ctx context.Context, eventID int, userID *int) error
	CanUserAccessEvent(ctx context.Context, eventID int, userID *int) (bool, error)
}

type eventService struct {
	eventRepo           repository.EventRepository
	addressRepo         repository.AddressRepository
	ticketRepo          repository.TicketRepository
	invitationRepo      repository.InvitationRepository
	categoryRepo        repository.CategoryRepository
	creatorRepo         repository.CreatorRepository
	mediaRepo           repository.MediaRepository
	subscriptionService SubscriptionService
	logger              *logger.Logger
}

func NewEventService(
	eventRepo repository.EventRepository,
	addressRepo repository.AddressRepository,
	ticketRepo repository.TicketRepository,
	invitationRepo repository.InvitationRepository,
	categoryRepo repository.CategoryRepository,
	creatorRepo repository.CreatorRepository,
	mediaRepo repository.MediaRepository,
	subscriptionService SubscriptionService,
	logger *logger.Logger,
) EventService {
	return &eventService{
		eventRepo:           eventRepo,
		addressRepo:         addressRepo,
		ticketRepo:          ticketRepo,
		invitationRepo:      invitationRepo,
		categoryRepo:        categoryRepo,
		creatorRepo:         creatorRepo,
		mediaRepo:           mediaRepo,
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

// Basic CRUD operations
func (s *eventService) CreateEvent(ctx context.Context, userID int, req dto.CreateEventRequest) (*dto.EventResponse, error) {
	s.logger.Info().Int("user_id", userID).Str("event_name", req.Name).Msg("Creating new event")

	// Get creator by user ID
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator by user ID")
		return nil, fmt.Errorf("creator profile not found")
	}

	creatorID := creator.ID
	s.logger.Info().Int("user_id", userID).Int("creator_id", creatorID).Msg("Found creator for user")

	// Validate address if provided
	if req.AddressID != nil {
		addressExists, err := s.addressRepo.ExistsByID(ctx, *req.AddressID)
		if err != nil {
			s.logger.Error().Err(err).Int("address_id", *req.AddressID).Msg("Failed to check address existence")
			return nil, fmt.Errorf("failed to validate address: %w", err)
		}
		if !addressExists {
			return nil, errors.New("event.address_not_found")
		}
	}

	// Validate categories if provided
	if len(req.CategoryIDs) > 0 {
		for _, categoryID := range req.CategoryIDs {
			_, err := s.categoryRepo.GetByID(ctx, categoryID)
			if err != nil {
				s.logger.Error().Err(err).Int("category_id", categoryID).Msg("Failed to check category existence")
				return nil, fmt.Errorf("failed to validate category: %w", err)
			}
		}
	}

	// Validate image if provided
	if req.ImageID != nil {
		mediaExists, err := s.mediaRepo.ExistsByID(ctx, *req.ImageID)
		if err != nil {
			s.logger.Error().Err(err).Int("image_id", *req.ImageID).Msg("Failed to check image existence")
			return nil, fmt.Errorf("failed to validate image: %w", err)
		}
		if !mediaExists {
			return nil, errors.New("event.image_not_found")
		}
	}

	// Validate video if provided
	if req.VideoID != nil {
		mediaExists, err := s.mediaRepo.ExistsByID(ctx, *req.VideoID)
		if err != nil {
			s.logger.Error().Err(err).Int("video_id", *req.VideoID).Msg("Failed to check video existence")
			return nil, fmt.Errorf("failed to validate video: %w", err)
		}
		if !mediaExists {
			return nil, errors.New("event.video_not_found")
		}
	}

	// Parse dates
	var startDate, endDate *time.Time
	var startTime, endTime *time.Time

	if req.StartDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}
		startDate = &parsed
	}

	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
		endDate = &parsed
	}

	if req.StartTime != nil {
		parsed, err := time.Parse("15:04", *req.StartTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start time format: %w", err)
		}
		startTime = &parsed
	}

	if req.EndTime != nil {
		parsed, err := time.Parse("15:04", *req.EndTime)
		if err != nil {
			return nil, fmt.Errorf("invalid end time format: %w", err)
		}
		endTime = &parsed
	}

	// Create event domain entity
	event := &domain.Event{
		CreatorID:        creatorID,
		Name:             req.Name,
		Description:      req.Description,
		ImageID:          req.ImageID,
		VideoID:          req.VideoID,
		Type:             req.Type,
		LocationType:     req.LocationType,
		Status:           domain.EventStatusDraft, // Always start as draft
		StartDate:        startDate,
		StartTime:        startTime,
		EndDate:          endDate,
		EndTime:          endTime,
		AddressID:        req.AddressID,
		OnlineEventURL:   req.OnlineEventURL,
		OnlineEventType:  req.OnlineEventType,
		TicketURL:        req.TicketURL,
		HasSystemTickets: req.HasSystemTickets,
		AdditionalInfo:   req.AdditionalInfo,
	}

	// Validate business rules
	if err := s.validateEventBusinessRules(event); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create event
	if err := s.eventRepo.Create(ctx, event); err != nil {
		s.logger.Error().Err(err).Int("creator_id", creatorID).Msg("Failed to create event")
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Add categories if provided
	if len(req.CategoryIDs) > 0 {
		if err := s.eventRepo.AddCategories(ctx, event.ID, req.CategoryIDs); err != nil {
			s.logger.Error().Err(err).Int("event_id", event.ID).Msg("Failed to add categories to event")
			return nil, fmt.Errorf("failed to add categories: %w", err)
		}
	}

	s.logger.Info().Int("event_id", event.ID).Int("creator_id", creatorID).Msg("Event created successfully")

	// Get event with relations for response
	createdEvent, err := s.eventRepo.GetByIDWithRelations(ctx, event.ID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", event.ID).Msg("Failed to get created event")
		return nil, fmt.Errorf("failed to get created event: %w", err)
	}

	return dto.EventToResponse(createdEvent), nil
}

func (s *eventService) GetEventByID(ctx context.Context, id int, userID *int) (*dto.EventResponse, error) {
	s.logger.Info().Int("event_id", id).Interface("user_id", userID).Msg("Getting event by ID")

	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to get event")
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Check access permissions
	if err := s.ValidateEventAccess(ctx, id, userID); err != nil {
		return nil, err
	}

	return dto.EventToResponse(event), nil
}

func (s *eventService) GetEventByIDWithRelations(ctx context.Context, id int, userID *int) (*dto.EventResponse, error) {
	s.logger.Info().Int("event_id", id).Interface("user_id", userID).Msg("Getting event with relations")

	event, err := s.eventRepo.GetByIDWithRelations(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to get event with relations")
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Check access permissions
	if err := s.ValidateEventAccess(ctx, id, userID); err != nil {
		return nil, err
	}

	return dto.EventToResponse(event), nil
}

func (s *eventService) UpdateEvent(ctx context.Context, id, userID int, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	s.logger.Info().Int("event_id", id).Int("user_id", userID).Msg("Updating event")

	// Get creator by user ID
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator by user ID")
		return nil, fmt.Errorf("creator profile not found")
	}

	creatorID := creator.ID

	// Validate ownership
	if err := s.ValidateEventOwnership(ctx, id, userID); err != nil {
		return nil, err
	}

	// Get existing event
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to get event for update")
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Check if event can be updated (only draft and pending events can be updated)
	if event.Status != domain.EventStatusDraft && event.Status != domain.EventStatusPending {
		return nil, errors.New("only draft and pending events can be updated")
	}

	// Update fields if provided
	if req.Name != nil {
		event.Name = *req.Name
	}
	if req.Description != nil {
		event.Description = req.Description
	}
	if req.ImageID != nil {
		event.ImageID = req.ImageID
	}
	if req.VideoID != nil {
		event.VideoID = req.VideoID
	}
	if req.Type != nil {
		event.Type = *req.Type
	}
	if req.LocationType != nil {
		event.LocationType = *req.LocationType
	}
	if req.AddressID != nil {
		// Validate address exists
		addressExists, err := s.addressRepo.ExistsByID(ctx, *req.AddressID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate address: %w", err)
		}
		if !addressExists {
			return nil, errors.New("address not found")
		}
		event.AddressID = req.AddressID
	}
	if req.OnlineEventURL != nil {
		event.OnlineEventURL = req.OnlineEventURL
	}
	if req.OnlineEventType != nil {
		event.OnlineEventType = req.OnlineEventType
	}
	if req.TicketURL != nil {
		event.TicketURL = req.TicketURL
	}
	if req.HasSystemTickets != nil {
		event.HasSystemTickets = *req.HasSystemTickets
	}
	if req.AdditionalInfo != nil {
		event.AdditionalInfo = req.AdditionalInfo
	}

	// Parse and update dates if provided
	if req.StartDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}
		event.StartDate = &parsed
	}
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
		event.EndDate = &parsed
	}
	if req.StartTime != nil {
		parsed, err := time.Parse("15:04", *req.StartTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start time format: %w", err)
		}
		event.StartTime = &parsed
	}
	if req.EndTime != nil {
		parsed, err := time.Parse("15:04", *req.EndTime)
		if err != nil {
			return nil, fmt.Errorf("invalid end time format: %w", err)
		}
		event.EndTime = &parsed
	}

	// Validate business rules
	if err := s.validateEventBusinessRules(event); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update event
	if err := s.eventRepo.Update(ctx, event); err != nil {
		s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to update event")
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	// Update categories if provided
	if req.CategoryIDs != nil {
		// Remove existing categories
		if err := s.eventRepo.RemoveCategories(ctx, id, []int{}); err != nil {
			s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to remove existing categories")
			return nil, fmt.Errorf("failed to update categories: %w", err)
		}

		// Add new categories
		if len(req.CategoryIDs) > 0 {
			// Validate categories
			for _, categoryID := range req.CategoryIDs {
				_, err := s.categoryRepo.GetByID(ctx, categoryID)
				if err != nil {
					return nil, fmt.Errorf("failed to validate category: %w", err)
				}
			}

			if err := s.eventRepo.AddCategories(ctx, id, req.CategoryIDs); err != nil {
				s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to add new categories")
				return nil, fmt.Errorf("failed to add categories: %w", err)
			}
		}
	}

	s.logger.Info().Int("event_id", id).Int("creator_id", creatorID).Msg("Event updated successfully")

	// Get updated event with relations
	updatedEvent, err := s.eventRepo.GetByIDWithRelations(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to get updated event")
		return nil, fmt.Errorf("failed to get updated event: %w", err)
	}

	return dto.EventToResponse(updatedEvent), nil
}

func (s *eventService) DeleteEvent(ctx context.Context, id, userID int) error {
	s.logger.Info().Int("event_id", id).Int("user_id", userID).Msg("Deleting event")

	// Get creator by user ID
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator by user ID")
		return fmt.Errorf("creator profile not found")
	}

	creatorID := creator.ID

	// Validate ownership
	if err := s.ValidateEventOwnership(ctx, id, userID); err != nil {
		return err
	}

	// Get event to check status
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to get event for deletion")
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Check if event can be deleted (only draft events can be deleted)
	if event.Status != domain.EventStatusDraft {
		return errors.New("only draft events can be deleted")
	}

	// Delete related data first
	// Delete invitations
	if err := s.invitationRepo.DeleteByEventID(ctx, id); err != nil {
		s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to delete event invitations")
		return fmt.Errorf("failed to delete invitations: %w", err)
	}

	// Delete tickets
	if err := s.ticketRepo.DeleteByEventID(ctx, id); err != nil {
		s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to delete event tickets")
		return fmt.Errorf("failed to delete tickets: %w", err)
	}

	// Delete event
	if err := s.eventRepo.Delete(ctx, id); err != nil {
		s.logger.Error().Err(err).Int("event_id", id).Msg("Failed to delete event")
		return fmt.Errorf("failed to delete event: %w", err)
	}

	s.logger.Info().Int("event_id", id).Int("creator_id", creatorID).Msg("Event deleted successfully")
	return nil
}

// Creator-specific operations
func (s *eventService) GetCreatorEvents(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	// Get creator by user ID
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator by user ID")
		return nil, nil, fmt.Errorf("creator profile not found")
	}

	creatorID := creator.ID
	events, paginationResp, err := s.eventRepo.GetByCreatorID(ctx, creatorID, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get creator events: %w", err)
	}

	var responses []*dto.EventListResponse
	for _, event := range events {
		responses = append(responses, dto.EventToListResponse(event))
	}

	return responses, paginationResp, nil
}

func (s *eventService) GetCreatorEventsByStatus(ctx context.Context, userID int, status domain.EventStatus, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	// Get creator by user ID
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator by user ID")
		return nil, nil, fmt.Errorf("creator profile not found")
	}

	creatorID := creator.ID
	events, paginationResp, err := s.eventRepo.GetByCreatorIDAndStatus(ctx, creatorID, status, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get creator events by status: %w", err)
	}

	var responses []*dto.EventListResponse
	for _, event := range events {
		responses = append(responses, dto.EventToListResponse(event))
	}

	return responses, paginationResp, nil
}

func (s *eventService) GetCreatorDraftEvents(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	return s.GetCreatorEventsByStatus(ctx, userID, domain.EventStatusDraft, pagination)
}

func (s *eventService) GetCreatorPublishedEvents(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	return s.GetCreatorEventsByStatus(ctx, userID, domain.EventStatusPublished, pagination)
}

// Public event operations
func (s *eventService) GetPublicEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	events, paginationResp, err := s.eventRepo.GetPublicEvents(ctx, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get public events: %w", err)
	}

	var responses []*dto.EventListResponse
	for _, event := range events {
		responses = append(responses, dto.EventToListResponse(event))
	}

	return responses, paginationResp, nil
}

func (s *eventService) GetPublicEventsByCategory(ctx context.Context, categoryID int, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	events, paginationResp, err := s.eventRepo.GetPublicEventsByCategory(ctx, categoryID, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get public events by category: %w", err)
	}

	var responses []*dto.EventListResponse
	for _, event := range events {
		responses = append(responses, dto.EventToListResponse(event))
	}

	return responses, paginationResp, nil
}

func (s *eventService) GetPublicEventsByLocation(ctx context.Context, city, country string, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	events, paginationResp, err := s.eventRepo.GetPublicEventsByLocation(ctx, city, pagination)
	if err != nil {
		s.logger.Error().Err(err).Str("city", city).Msg("Failed to get public events by location")
		return nil, nil, fmt.Errorf("failed to get public events by location: %w", err)
	}

	var responses []*dto.EventListResponse
	for _, event := range events {
		responses = append(responses, dto.EventToListResponse(event))
	}

	return responses, paginationResp, nil
}

func (s *eventService) SearchPublicEvents(ctx context.Context, req dto.EventSearchRequest, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	// Search public events - use advanced filtering instead
	publishedStatus := domain.EventStatusPublished
	publicType := domain.EventTypePublic

	filters := dto.EventFilterRequest{
		Query:        &req.Query,
		Type:         &publicType,
		LocationType: req.LocationType,
		CategoryIDs:  req.CategoryIDs,
		City:         req.City,
		Status:       &publishedStatus,
	}

	events, paginationResp, err := s.eventRepo.GetEventsWithFilters(ctx, filters, pagination)
	if err != nil {
		s.logger.Error().Err(err).Interface("request", req).Msg("Failed to search public events")
		return nil, nil, fmt.Errorf("failed to search events: %w", err)
	}

	var responses []*dto.EventListResponse
	for _, event := range events {
		responses = append(responses, dto.EventToListResponse(event))
	}

	return responses, paginationResp, nil
}

// Status management
func (s *eventService) UpdateEventStatus(ctx context.Context, id, userID int, req dto.UpdateEventStatusRequest) (*dto.EventResponse, error) {
	// Validate ownership
	if err := s.ValidateEventOwnership(ctx, id, userID); err != nil {
		return nil, err
	}

	// Get current event
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Validate status transition
	if err := s.validateStatusTransition(event.Status, req.Status); err != nil {
		return nil, err
	}

	// Update status
	if err := s.eventRepo.UpdateStatus(ctx, id, req.Status); err != nil {
		return nil, fmt.Errorf("failed to update event status: %w", err)
	}

	// Get updated event
	updatedEvent, err := s.eventRepo.GetByIDWithRelations(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated event: %w", err)
	}

	return dto.EventToResponse(updatedEvent), nil
}

func (s *eventService) SubmitEventForReview(ctx context.Context, id, userID int) (*dto.EventResponse, error) {
	s.logger.Info().Int("event_id", id).Int("user_id", userID).Msg("Submitting event for review with subscription validation")

	// Validate ownership first
	if err := s.ValidateEventOwnership(ctx, id, userID); err != nil {
		return nil, err
	}

	// Check publishing rights before allowing draft->pending transition
	canPublish, err := s.subscriptionService.CanPublishEvent(ctx, uint(userID))
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to check publishing rights")
		return nil, fmt.Errorf("failed to check publishing rights: %w", err)
	}

	if !canPublish {
		s.logger.Warn().Int("user_id", userID).Msg("User does not have publishing rights")
		return nil, errors.New("subscription.insufficient_publishing_rights")
	}

	// If user has rights, proceed with status update and consume usage
	response, err := s.UpdateEventStatus(ctx, id, userID, dto.UpdateEventStatusRequest{
		Status: domain.EventStatusPending,
	})
	if err != nil {
		return nil, err
	}

	// Track usage after successful status update
	if err := s.subscriptionService.ConsumeEventCredit(ctx, uint(userID)); err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to consume event publishing usage")
		// Don't fail the request, just log the error
		// The event status has already been updated successfully
	} else {
		s.logger.Info().Int("user_id", userID).Msg("Event publishing usage consumed successfully")
	}

	return response, nil
}

func (s *eventService) PublishEvent(ctx context.Context, id, userID int) (*dto.EventResponse, error) {
	return s.UpdateEventStatus(ctx, id, userID, dto.UpdateEventStatusRequest{
		Status: domain.EventStatusPublished,
	})
}

func (s *eventService) CancelEvent(ctx context.Context, id, userID int) (*dto.EventResponse, error) {
	return s.UpdateEventStatus(ctx, id, userID, dto.UpdateEventStatusRequest{
		Status: domain.EventStatusCancelled,
	})
}

// Advanced filtering and search
func (s *eventService) GetEventsWithFilters(ctx context.Context, filters dto.EventFilterRequest, pagination dto.PaginationRequest, userID *int) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	events, paginationResp, err := s.eventRepo.GetEventsWithFilters(ctx, filters, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get events with filters: %w", err)
	}

	var responses []*dto.EventListResponse
	for _, event := range events {
		// Check access for each event
		if canAccess, _ := s.CanUserAccessEvent(ctx, event.ID, userID); canAccess {
			responses = append(responses, dto.EventToListResponse(event))
		}
	}

	return responses, paginationResp, nil
}

func (s *eventService) GetEventsByDateRange(ctx context.Context, startDate, endDate time.Time, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	// Convert dates to string format for repository
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	events, paginationResp, err := s.eventRepo.GetEventsByDateRange(ctx, startDateStr, endDateStr, pagination)
	if err != nil {
		s.logger.Error().Err(err).Time("start_date", startDate).Time("end_date", endDate).Msg("Failed to get events by date range")
		return nil, nil, fmt.Errorf("failed to get events by date range: %w", err)
	}

	var responses []*dto.EventListResponse
	for _, event := range events {
		responses = append(responses, dto.EventToListResponse(event))
	}

	return responses, paginationResp, nil
}

func (s *eventService) GetUpcomingEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	events, paginationResp, err := s.eventRepo.GetUpcomingEvents(ctx, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}

	var responses []*dto.EventListResponse
	for _, event := range events {
		responses = append(responses, dto.EventToListResponse(event))
	}

	return responses, paginationResp, nil
}

func (s *eventService) GetPastEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.EventListResponse, *dto.PaginationResponse, error) {
	events, paginationResp, err := s.eventRepo.GetPastEvents(ctx, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get past events: %w", err)
	}

	var responses []*dto.EventListResponse
	for _, event := range events {
		responses = append(responses, dto.EventToListResponse(event))
	}

	return responses, paginationResp, nil
}

// Statistics operations
func (s *eventService) GetEventStats(ctx context.Context, userID int) (*dto.EventStatsResponse, error) {
	// Get creator by user ID
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator by user ID")
		return nil, fmt.Errorf("creator profile not found")
	}

	creatorID := creator.ID
	stats, err := s.eventRepo.GetEventStats(ctx, creatorID)
	if err != nil {
		s.logger.Error().Err(err).Int("creator_id", creatorID).Msg("Failed to get creator event stats")
		return nil, fmt.Errorf("failed to get event stats: %w", err)
	}
	return stats, nil
}

func (s *eventService) GetSystemEventStats(ctx context.Context) (*dto.SystemEventStatsResponse, error) {
	stats, err := s.eventRepo.GetSystemEventStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get system event stats: %w", err)
	}
	return stats, nil
}

// Business rule validation
func (s *eventService) validateEventBusinessRules(event *domain.Event) error {
	// Location type specific validations
	switch event.LocationType {
	case domain.EventLocationTypeLocation:
		if event.AddressID == nil {
			return errors.New("address is required for location-based events")
		}
		if event.OnlineEventURL != nil {
			return errors.New("online event URL should not be provided for location-based events")
		}
	case domain.EventLocationTypeOnline:
		if event.OnlineEventURL == nil || *event.OnlineEventURL == "" {
			return errors.New("online event URL is required for online events")
		}
		if event.AddressID != nil {
			return errors.New("address should not be provided for online events")
		}
	case domain.EventLocationTypeAnnouncement:
		// Announcements don't require address or URL
		break
	}

	// Date validations
	if event.StartDate != nil && event.EndDate != nil {
		if event.EndDate.Before(*event.StartDate) {
			return errors.New("end date cannot be before start date")
		}
	}

	// Time validations (if same date)
	if event.StartDate != nil && event.EndDate != nil &&
		event.StartTime != nil && event.EndTime != nil {
		if event.StartDate.Equal(*event.EndDate) && event.EndTime.Before(*event.StartTime) {
			return errors.New("end time cannot be before start time on the same date")
		}
	}

	// System tickets validation
	if event.HasSystemTickets && event.TicketURL != nil {
		return errors.New("cannot have both system tickets and external ticket URL")
	}

	return nil
}

// Status transition validation
func (s *eventService) validateStatusTransition(currentStatus, newStatus domain.EventStatus) error {
	validTransitions := map[domain.EventStatus][]domain.EventStatus{
		domain.EventStatusDraft: {
			domain.EventStatusPending,
			domain.EventStatusPublished, // Direct publish for creators
		},
		domain.EventStatusPending: {
			domain.EventStatusDraft,     // Back to draft
			domain.EventStatusPublished, // Approved
			domain.EventStatusRejected,  // Rejected
		},
		domain.EventStatusPublished: {
			domain.EventStatusStopped,   // Temporarily stop
			domain.EventStatusCancelled, // Cancel permanently
		},
		domain.EventStatusRejected: {
			domain.EventStatusDraft,   // Back to draft for fixes
			domain.EventStatusPending, // Resubmit
		},
		domain.EventStatusStopped: {
			domain.EventStatusPublished, // Resume
			domain.EventStatusCancelled, // Cancel permanently
		},
		// Cancelled is final - no transitions allowed
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("no transitions allowed from status %s", currentStatus)
	}

	for _, allowedStatus := range allowedStatuses {
		if newStatus == allowedStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
}

// Validation operations
func (s *eventService) ValidateEventOwnership(ctx context.Context, eventID, userID int) error {
	// Get creator by user ID
	creator, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get creator by user ID")
		return fmt.Errorf("creator profile not found")
	}

	creatorID := creator.ID
	isOwner, err := s.eventRepo.IsEventOwner(ctx, eventID, creatorID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Int("creator_id", creatorID).Msg("Failed to validate event ownership")
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	if !isOwner {
		return errors.New("access denied: you don't own this event")
	}
	return nil
}

func (s *eventService) ValidateEventAccess(ctx context.Context, eventID int, userID *int) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Public events are accessible to everyone if published
	if event.Type == domain.EventTypePublic && event.Status == domain.EventStatusPublished {
		return nil
	}

	// Private events require user to be logged in
	if event.Type == domain.EventTypePrivate {
		if userID == nil {
			return errors.New("authentication required for private events")
		}

		// Check if user is the creator
		creator, err := s.creatorRepo.GetByUserID(ctx, *userID)
		if err == nil && creator.ID == event.CreatorID {
			return nil
		}

		// Check if user is invited
		hasInvitation, err := s.invitationRepo.ExistsByEventAndUser(ctx, eventID, *userID)
		if err != nil {
			return fmt.Errorf("failed to check invitation: %w", err)
		}
		if !hasInvitation {
			return errors.New("access denied: you are not invited to this private event")
		}
	}

	// For non-published events, only creator can access
	if event.Status != domain.EventStatusPublished {
		if userID == nil {
			return errors.New("authentication required")
		}

		creator, err := s.creatorRepo.GetByUserID(ctx, *userID)
		if err != nil {
			return errors.New("creator profile required")
		}
		if creator.ID != event.CreatorID {
			return errors.New("access denied: event not published")
		}
	}

	return nil
}

func (s *eventService) CanUserAccessEvent(ctx context.Context, eventID int, userID *int) (bool, error) {
	err := s.ValidateEventAccess(ctx, eventID, userID)
	return err == nil, nil
}
