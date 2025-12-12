package service

import (
	"context"
	"fmt"
	"time"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
	"github.com/rs/zerolog"
)

type InvitationService interface {
	// Basic CRUD operations
	CreateInvitation(ctx context.Context, eventID int, req dto.CreateInvitationRequest) (*dto.InvitationResponse, error)
	GetInvitationByID(ctx context.Context, id int) (*dto.InvitationResponse, error)
	GetInvitationByIDWithRelations(ctx context.Context, id int) (*dto.InvitationResponse, error)
	UpdateInvitationStatus(ctx context.Context, id int, req dto.UpdateInvitationStatusRequest) (*dto.InvitationResponse, error)
	DeleteInvitation(ctx context.Context, id int) error

	// Event-specific operations
	GetInvitationsByEventID(ctx context.Context, eventID int, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)
	GetInvitationsByEventIDWithRelations(ctx context.Context, eventID int, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)
	GetEventInvitationStats(ctx context.Context, eventID int) (*dto.InvitationStatsResponse, error)

	// User-specific operations
	GetInvitationsByUserID(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)
	GetInvitationsByEmail(ctx context.Context, email string, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)
	GetUserInvitationStats(ctx context.Context, userID int) (*dto.UserInvitationStatsResponse, error)

	// Status-specific operations
	GetInvitationsByStatus(ctx context.Context, status domain.InvitationStatus, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)
	GetPendingInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)
	GetApprovedInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)
	GetRejectedInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)

	// Status management operations
	ApproveInvitation(ctx context.Context, id int) (*dto.InvitationResponse, error)
	RejectInvitation(ctx context.Context, id int) (*dto.InvitationResponse, error)
	ResetInvitationToPending(ctx context.Context, id int) (*dto.InvitationResponse, error)

	// Bulk operations
	CreateMultipleInvitations(ctx context.Context, eventID int, req dto.BulkCreateInvitationRequest) ([]*dto.InvitationResponse, error)
	DeleteAllEventInvitations(ctx context.Context, eventID int) error

	// Email-based operations for external users
	GetPendingInvitationsByEmail(ctx context.Context, email string) ([]*dto.InvitationResponse, error)
	GetEventInvitationByEmail(ctx context.Context, eventID int, email string) (*dto.InvitationResponse, error)
	RespondToInvitationByEmail(ctx context.Context, eventID int, email string, status domain.InvitationStatus) (*dto.InvitationResponse, error)

	// Expiration operations
	GetExpiredInvitations(ctx context.Context, expirationHours int, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)
	CleanupExpiredInvitations(ctx context.Context, expirationHours int) (int64, error)
	GetInvitationsExpiringIn(ctx context.Context, hours int, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)

	// Advanced filtering
	GetInvitationsWithFilters(ctx context.Context, filters dto.InvitationFilterRequest, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error)

	// Statistics operations
	GetSystemInvitationStats(ctx context.Context) (*dto.SystemInvitationStatsResponse, error)

	// Validation operations
	ValidateInvitationAccess(ctx context.Context, invitationID int, userID *int, email *string) error
	ValidateInvitationOwnership(ctx context.Context, invitationID, eventID int) error
	ValidateInvitationData(ctx context.Context, invitation *domain.Invitation) error
	CanUserRespondToInvitation(ctx context.Context, invitationID int, userID *int, email *string) (bool, error)
}

type invitationService struct {
	invitationRepo repository.InvitationRepository
	eventRepo      repository.EventRepository
	userRepo       repository.UserRepository
	logger         zerolog.Logger
}

func NewInvitationService(
	invitationRepo repository.InvitationRepository,
	eventRepo repository.EventRepository,
	userRepo repository.UserRepository,
	logger zerolog.Logger,
) InvitationService {
	return &invitationService{
		invitationRepo: invitationRepo,
		eventRepo:      eventRepo,
		userRepo:       userRepo,
		logger:         logger.With().Str("service", "invitation").Logger(),
	}
}

// Basic CRUD operations
func (s *invitationService) CreateInvitation(ctx context.Context, eventID int, req dto.CreateInvitationRequest) (*dto.InvitationResponse, error) {
	// Validate event exists
	exists, err := s.eventRepo.ExistsByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to check event existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("event not found")
	}

	// Validate request
	if err := s.validateCreateInvitationRequest(&req); err != nil {
		s.logger.Error().Err(err).Interface("request", req).Msg("Invalid create invitation request")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check for duplicate invitation
	if req.InvitedUserID != nil {
		exists, err := s.invitationRepo.ExistsByEventAndUser(ctx, eventID, *req.InvitedUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing invitation: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("user is already invited to this event")
		}
	}

	exists, err = s.invitationRepo.ExistsByEventAndEmail(ctx, eventID, req.InvitedEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing invitation: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email is already invited to this event")
	}

	// Create domain entity
	invitation := &domain.Invitation{
		EventID:       eventID,
		InvitedUserID: req.InvitedUserID,
		InvitedEmail:  req.InvitedEmail,
		Status:        domain.InvitationStatusPending,
		InvitedAt:     time.Now(),
	}

	// Validate domain entity
	if err := s.ValidateInvitationData(ctx, invitation); err != nil {
		return nil, err
	}

	// Create invitation
	if err := s.invitationRepo.Create(ctx, invitation); err != nil {
		s.logger.Error().Err(err).Interface("invitation", invitation).Msg("Failed to create invitation")
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	s.logger.Info().Int("invitation_id", invitation.ID).Int("event_id", eventID).Str("email", invitation.InvitedEmail).Msg("Invitation created successfully")

	return s.invitationToResponse(invitation), nil
}

func (s *invitationService) GetInvitationByID(ctx context.Context, id int) (*dto.InvitationResponse, error) {
	invitation, err := s.invitationRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("invitation_id", id).Msg("Failed to get invitation by ID")
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	return s.invitationToResponse(invitation), nil
}

func (s *invitationService) GetInvitationByIDWithRelations(ctx context.Context, id int) (*dto.InvitationResponse, error) {
	invitation, err := s.invitationRepo.GetByIDWithRelations(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("invitation_id", id).Msg("Failed to get invitation by ID with relations")
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	return s.invitationToResponse(invitation), nil
}

func (s *invitationService) UpdateInvitationStatus(ctx context.Context, id int, req dto.UpdateInvitationStatusRequest) (*dto.InvitationResponse, error) {
	// Get existing invitation
	existingInvitation, err := s.invitationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing invitation: %w", err)
	}

	// Validate status transition
	if err := s.validateStatusTransition(existingInvitation.Status, req.Status); err != nil {
		return nil, err
	}

	// Update status
	if err := s.invitationRepo.UpdateStatus(ctx, id, req.Status); err != nil {
		s.logger.Error().Err(err).Int("invitation_id", id).Str("status", string(req.Status)).Msg("Failed to update invitation status")
		return nil, fmt.Errorf("failed to update invitation status: %w", err)
	}

	// Set responded_at timestamp if status is not pending
	if req.Status != domain.InvitationStatusPending && existingInvitation.RespondedAt == nil {
		now := time.Now()
		existingInvitation.RespondedAt = &now
		if err := s.invitationRepo.Update(ctx, existingInvitation); err != nil {
			s.logger.Warn().Err(err).Int("invitation_id", id).Msg("Failed to update responded_at timestamp")
		}
	}

	s.logger.Info().Int("invitation_id", id).Str("status", string(req.Status)).Msg("Invitation status updated successfully")

	// Get updated invitation
	updatedInvitation, err := s.invitationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated invitation: %w", err)
	}

	return s.invitationToResponse(updatedInvitation), nil
}

func (s *invitationService) DeleteInvitation(ctx context.Context, id int) error {
	// Check if invitation exists
	exists, err := s.invitationRepo.ExistsByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check invitation existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("invitation not found")
	}

	// Delete invitation
	if err := s.invitationRepo.Delete(ctx, id); err != nil {
		s.logger.Error().Err(err).Int("invitation_id", id).Msg("Failed to delete invitation")
		return fmt.Errorf("failed to delete invitation: %w", err)
	}

	s.logger.Info().Int("invitation_id", id).Msg("Invitation deleted successfully")
	return nil
}

// Event-specific operations
func (s *invitationService) GetInvitationsByEventID(ctx context.Context, eventID int, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetByEventID(ctx, eventID, pagination)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get invitations by event ID")
		return nil, nil, fmt.Errorf("failed to get invitations: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

func (s *invitationService) GetInvitationsByEventIDWithRelations(ctx context.Context, eventID int, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetByEventIDWithRelations(ctx, eventID, pagination)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get invitations by event ID with relations")
		return nil, nil, fmt.Errorf("failed to get invitations: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

func (s *invitationService) GetEventInvitationStats(ctx context.Context, eventID int) (*dto.InvitationStatsResponse, error) {
	stats, err := s.invitationRepo.GetInvitationStats(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get event invitation stats")
		return nil, fmt.Errorf("failed to get invitation stats: %w", err)
	}

	return stats, nil
}

// User-specific operations
func (s *invitationService) GetInvitationsByUserID(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetByUserID(ctx, userID, pagination)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get invitations by user ID")
		return nil, nil, fmt.Errorf("failed to get invitations: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

func (s *invitationService) GetInvitationsByEmail(ctx context.Context, email string, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetByEmail(ctx, email, pagination)
	if err != nil {
		s.logger.Error().Err(err).Str("email", email).Msg("Failed to get invitations by email")
		return nil, nil, fmt.Errorf("failed to get invitations: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

func (s *invitationService) GetUserInvitationStats(ctx context.Context, userID int) (*dto.UserInvitationStatsResponse, error) {
	stats, err := s.invitationRepo.GetUserInvitationStats(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get user invitation stats")
		return nil, fmt.Errorf("failed to get user invitation stats: %w", err)
	}

	return stats, nil
}

// Status-specific operations
func (s *invitationService) GetInvitationsByStatus(ctx context.Context, status domain.InvitationStatus, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetByStatus(ctx, status, pagination)
	if err != nil {
		s.logger.Error().Err(err).Str("status", string(status)).Msg("Failed to get invitations by status")
		return nil, nil, fmt.Errorf("failed to get invitations: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

func (s *invitationService) GetPendingInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetPendingInvitations(ctx, pagination)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get pending invitations")
		return nil, nil, fmt.Errorf("failed to get pending invitations: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

func (s *invitationService) GetApprovedInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetApprovedInvitations(ctx, pagination)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get approved invitations")
		return nil, nil, fmt.Errorf("failed to get approved invitations: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

func (s *invitationService) GetRejectedInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetRejectedInvitations(ctx, pagination)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get rejected invitations")
		return nil, nil, fmt.Errorf("failed to get rejected invitations: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

// Status management operations
func (s *invitationService) ApproveInvitation(ctx context.Context, id int) (*dto.InvitationResponse, error) {
	return s.UpdateInvitationStatus(ctx, id, dto.UpdateInvitationStatusRequest{
		Status: domain.InvitationStatusApproved,
	})
}

func (s *invitationService) RejectInvitation(ctx context.Context, id int) (*dto.InvitationResponse, error) {
	return s.UpdateInvitationStatus(ctx, id, dto.UpdateInvitationStatusRequest{
		Status: domain.InvitationStatusRejected,
	})
}

func (s *invitationService) ResetInvitationToPending(ctx context.Context, id int) (*dto.InvitationResponse, error) {
	return s.UpdateInvitationStatus(ctx, id, dto.UpdateInvitationStatusRequest{
		Status: domain.InvitationStatusPending,
	})
}

// Bulk operations
func (s *invitationService) CreateMultipleInvitations(ctx context.Context, eventID int, req dto.BulkCreateInvitationRequest) ([]*dto.InvitationResponse, error) {
	if len(req.Invitations) == 0 {
		return nil, fmt.Errorf("no invitations to create")
	}

	// Validate event exists
	exists, err := s.eventRepo.ExistsByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to check event existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("event not found")
	}

	var invitations []*domain.Invitation
	for i, invReq := range req.Invitations {
		// Validate request
		if err := s.validateCreateInvitationRequest(&invReq); err != nil {
			return nil, fmt.Errorf("validation failed for invitation %d: %w", i+1, err)
		}

		// Check for duplicate invitation
		if invReq.InvitedUserID != nil {
			exists, err := s.invitationRepo.ExistsByEventAndUser(ctx, eventID, *invReq.InvitedUserID)
			if err != nil {
				return nil, fmt.Errorf("failed to check existing invitation for invitation %d: %w", i+1, err)
			}
			if exists {
				return nil, fmt.Errorf("user in invitation %d is already invited to this event", i+1)
			}
		}

		exists, err := s.invitationRepo.ExistsByEventAndEmail(ctx, eventID, invReq.InvitedEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing invitation for invitation %d: %w", i+1, err)
		}
		if exists {
			return nil, fmt.Errorf("email in invitation %d is already invited to this event", i+1)
		}

		invitation := &domain.Invitation{
			EventID:       eventID,
			InvitedUserID: invReq.InvitedUserID,
			InvitedEmail:  invReq.InvitedEmail,
			Status:        domain.InvitationStatusPending,
			InvitedAt:     time.Now(),
		}

		// Validate domain entity
		if err := s.ValidateInvitationData(ctx, invitation); err != nil {
			return nil, fmt.Errorf("validation failed for invitation %d: %w", i+1, err)
		}

		invitations = append(invitations, invitation)
	}

	// Create all invitations
	if err := s.invitationRepo.CreateMultiple(ctx, invitations); err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Int("count", len(invitations)).Msg("Failed to create multiple invitations")
		return nil, fmt.Errorf("failed to create invitations: %w", err)
	}

	s.logger.Info().Int("event_id", eventID).Int("count", len(invitations)).Msg("Multiple invitations created successfully")

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, nil
}

func (s *invitationService) DeleteAllEventInvitations(ctx context.Context, eventID int) error {
	// Delete all event invitations
	if err := s.invitationRepo.DeleteByEventID(ctx, eventID); err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to delete all event invitations")
		return fmt.Errorf("failed to delete all event invitations: %w", err)
	}

	s.logger.Info().Int("event_id", eventID).Msg("All event invitations deleted successfully")
	return nil
}

// Email-based operations for external users
func (s *invitationService) GetPendingInvitationsByEmail(ctx context.Context, email string) ([]*dto.InvitationResponse, error) {
	invitations, err := s.invitationRepo.GetPendingInvitationsByEmail(ctx, email)
	if err != nil {
		s.logger.Error().Err(err).Str("email", email).Msg("Failed to get pending invitations by email")
		return nil, fmt.Errorf("failed to get pending invitations: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, nil
}

func (s *invitationService) GetEventInvitationByEmail(ctx context.Context, eventID int, email string) (*dto.InvitationResponse, error) {
	invitation, err := s.invitationRepo.GetEventInvitationByEmail(ctx, eventID, email)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Str("email", email).Msg("Failed to get event invitation by email")
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	return s.invitationToResponse(invitation), nil
}

func (s *invitationService) RespondToInvitationByEmail(ctx context.Context, eventID int, email string, status domain.InvitationStatus) (*dto.InvitationResponse, error) {
	// Get invitation by event and email
	invitation, err := s.invitationRepo.GetEventInvitationByEmail(ctx, eventID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	// Validate status transition
	if err := s.validateStatusTransition(invitation.Status, status); err != nil {
		return nil, err
	}

	// Update status
	if err := s.invitationRepo.UpdateStatus(ctx, invitation.ID, status); err != nil {
		s.logger.Error().Err(err).Int("invitation_id", invitation.ID).Str("status", string(status)).Msg("Failed to update invitation status")
		return nil, fmt.Errorf("failed to update invitation status: %w", err)
	}

	// Set responded_at timestamp
	if status != domain.InvitationStatusPending && invitation.RespondedAt == nil {
		now := time.Now()
		invitation.RespondedAt = &now
		if err := s.invitationRepo.Update(ctx, invitation); err != nil {
			s.logger.Warn().Err(err).Int("invitation_id", invitation.ID).Msg("Failed to update responded_at timestamp")
		}
	}

	s.logger.Info().Int("invitation_id", invitation.ID).Str("email", email).Str("status", string(status)).Msg("Invitation responded by email successfully")

	// Get updated invitation
	updatedInvitation, err := s.invitationRepo.GetByID(ctx, invitation.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated invitation: %w", err)
	}

	return s.invitationToResponse(updatedInvitation), nil
}

// Expiration operations
func (s *invitationService) GetExpiredInvitations(ctx context.Context, expirationHours int, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetExpiredInvitations(ctx, expirationHours, pagination)
	if err != nil {
		s.logger.Error().Err(err).Int("expiration_hours", expirationHours).Msg("Failed to get expired invitations")
		return nil, nil, fmt.Errorf("failed to get expired invitations: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

func (s *invitationService) CleanupExpiredInvitations(ctx context.Context, expirationHours int) (int64, error) {
	deletedCount, err := s.invitationRepo.DeleteExpiredInvitations(ctx, expirationHours)
	if err != nil {
		s.logger.Error().Err(err).Int("expiration_hours", expirationHours).Msg("Failed to cleanup expired invitations")
		return 0, fmt.Errorf("failed to cleanup expired invitations: %w", err)
	}

	s.logger.Info().Int64("deleted_count", deletedCount).Int("expiration_hours", expirationHours).Msg("Expired invitations cleaned up successfully")
	return deletedCount, nil
}

func (s *invitationService) GetInvitationsExpiringIn(ctx context.Context, hours int, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetInvitationsExpiringIn(ctx, hours, pagination)
	if err != nil {
		s.logger.Error().Err(err).Int("hours", hours).Msg("Failed to get invitations expiring in hours")
		return nil, nil, fmt.Errorf("failed to get invitations expiring in hours: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

// Advanced filtering
func (s *invitationService) GetInvitationsWithFilters(ctx context.Context, filters dto.InvitationFilterRequest, pagination dto.PaginationRequest) ([]*dto.InvitationResponse, *dto.PaginationResponse, error) {
	invitations, paginationResp, err := s.invitationRepo.GetInvitationsWithFilters(ctx, filters, pagination)
	if err != nil {
		s.logger.Error().Err(err).Interface("filters", filters).Msg("Failed to get invitations with filters")
		return nil, nil, fmt.Errorf("failed to get invitations with filters: %w", err)
	}

	var responses []*dto.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, s.invitationToResponse(invitation))
	}

	return responses, paginationResp, nil
}

// Statistics operations
func (s *invitationService) GetSystemInvitationStats(ctx context.Context) (*dto.SystemInvitationStatsResponse, error) {
	stats, err := s.invitationRepo.GetSystemInvitationStats(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get system invitation stats")
		return nil, fmt.Errorf("failed to get system invitation stats: %w", err)
	}

	return stats, nil
}

// Validation operations
func (s *invitationService) ValidateInvitationAccess(ctx context.Context, invitationID int, userID *int, email *string) error {
	canAccess, err := s.invitationRepo.CanUserAccessInvitation(ctx, invitationID, getUserIDOrZero(userID), getEmailOrEmpty(email))
	if err != nil {
		s.logger.Error().Err(err).Int("invitation_id", invitationID).Msg("Failed to validate invitation access")
		return fmt.Errorf("failed to validate access: %w", err)
	}
	if !canAccess {
		return fmt.Errorf("access denied: you cannot access this invitation")
	}
	return nil
}

func (s *invitationService) ValidateInvitationOwnership(ctx context.Context, invitationID, eventID int) error {
	isOwner, err := s.invitationRepo.IsInvitationOwner(ctx, invitationID, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("invitation_id", invitationID).Int("event_id", eventID).Msg("Failed to validate invitation ownership")
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("access denied: invitation does not belong to this event")
	}
	return nil
}

func (s *invitationService) ValidateInvitationData(ctx context.Context, invitation *domain.Invitation) error {
	if invitation.EventID <= 0 {
		return fmt.Errorf("event ID is required")
	}
	if invitation.InvitedEmail == "" {
		return fmt.Errorf("invited email is required")
	}

	// Validate email format (basic validation)
	if len(invitation.InvitedEmail) < 5 || !contains(invitation.InvitedEmail, "@") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func (s *invitationService) CanUserRespondToInvitation(ctx context.Context, invitationID int, userID *int, email *string) (bool, error) {
	canAccess, err := s.invitationRepo.CanUserAccessInvitation(ctx, invitationID, getUserIDOrZero(userID), getEmailOrEmpty(email))
	if err != nil {
		s.logger.Error().Err(err).Int("invitation_id", invitationID).Msg("Failed to check if user can respond to invitation")
		return false, fmt.Errorf("failed to check access: %w", err)
	}

	return canAccess, nil
}

// Helper methods
func (s *invitationService) validateCreateInvitationRequest(req *dto.CreateInvitationRequest) error {
	if req.InvitedEmail == "" {
		return fmt.Errorf("invited email is required")
	}

	// Basic email validation
	if len(req.InvitedEmail) < 5 || !contains(req.InvitedEmail, "@") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func (s *invitationService) validateStatusTransition(currentStatus, newStatus domain.InvitationStatus) error {
	validTransitions := map[domain.InvitationStatus][]domain.InvitationStatus{
		domain.InvitationStatusPending: {
			domain.InvitationStatusApproved,
			domain.InvitationStatusRejected,
		},
		domain.InvitationStatusApproved: {
			domain.InvitationStatusPending, // Allow reset to pending
		},
		domain.InvitationStatusRejected: {
			domain.InvitationStatusPending, // Allow reset to pending
		},
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

func (s *invitationService) invitationToResponse(invitation *domain.Invitation) *dto.InvitationResponse {
	response := &dto.InvitationResponse{
		ID:            invitation.ID,
		EventID:       invitation.EventID,
		InvitedUserID: invitation.InvitedUserID,
		InvitedEmail:  invitation.InvitedEmail,
		Status:        invitation.Status,
		InvitedAt:     invitation.InvitedAt,
		RespondedAt:   invitation.RespondedAt,
		CreatedAt:     invitation.CreatedAt,
		UpdatedAt:     invitation.UpdatedAt,
	}

	// Add relations if loaded
	if invitation.InvitedUser != nil {
		response.InvitedUser = &dto.UserBasicResponse{
			ID:       invitation.InvitedUser.ID,
			FullName: invitation.InvitedUser.FullName,
			Username: invitation.InvitedUser.Username,
		}
	}

	return response
}

// Utility functions
func getUserIDOrZero(userID *int) int {
	if userID == nil {
		return 0
	}
	return *userID
}

func getEmailOrEmpty(email *string) string {
	if email == nil {
		return ""
	}
	return *email
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
