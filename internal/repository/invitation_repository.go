package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
)

type InvitationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, invitation *domain.Invitation) error
	GetByID(ctx context.Context, id int) (*domain.Invitation, error)
	GetByIDWithRelations(ctx context.Context, id int) (*domain.Invitation, error)
	Update(ctx context.Context, invitation *domain.Invitation) error
	Delete(ctx context.Context, id int) error

	// Event-specific operations
	GetByEventID(ctx context.Context, eventID int, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)
	GetByEventIDWithRelations(ctx context.Context, eventID int, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)
	CountByEventID(ctx context.Context, eventID int) (int64, error)
	CountByEventIDAndStatus(ctx context.Context, eventID int, status domain.InvitationStatus) (int64, error)

	// User-specific operations
	GetByUserID(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)
	GetByEmail(ctx context.Context, email string, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)
	GetByUserIDAndStatus(ctx context.Context, userID int, status domain.InvitationStatus, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)
	GetByEmailAndStatus(ctx context.Context, email string, status domain.InvitationStatus, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)

	// Status-specific operations
	GetByStatus(ctx context.Context, status domain.InvitationStatus, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)
	GetPendingInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)
	GetApprovedInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)
	GetRejectedInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)

	// Status management operations
	UpdateStatus(ctx context.Context, id int, status domain.InvitationStatus) error
	ApproveInvitation(ctx context.Context, id int) error
	RejectInvitation(ctx context.Context, id int) error
	ResetToPending(ctx context.Context, id int) error

	// Duplicate and validation operations
	ExistsByEventAndEmail(ctx context.Context, eventID int, email string) (bool, error)
	ExistsByEventAndUser(ctx context.Context, eventID int, userID int) (bool, error)
	GetByEventAndEmail(ctx context.Context, eventID int, email string) (*domain.Invitation, error)
	GetByEventAndUser(ctx context.Context, eventID int, userID int) (*domain.Invitation, error)

	// Validation operations
	ExistsByID(ctx context.Context, id int) (bool, error)
	IsInvitationOwner(ctx context.Context, invitationID, eventID int) (bool, error)
	CanUserAccessInvitation(ctx context.Context, invitationID, userID int, email string) (bool, error)

	// Bulk operations
	GetMultipleByIDs(ctx context.Context, ids []int) ([]*domain.Invitation, error)
	CreateMultiple(ctx context.Context, invitations []*domain.Invitation) error
	UpdateMultiple(ctx context.Context, invitations []*domain.Invitation) error
	DeleteMultiple(ctx context.Context, ids []int) error
	DeleteByEventID(ctx context.Context, eventID int) error

	// Statistics operations
	GetInvitationStats(ctx context.Context, eventID int) (*dto.InvitationStatsResponse, error)
	GetUserInvitationStats(ctx context.Context, userID int) (*dto.UserInvitationStatsResponse, error)
	GetSystemInvitationStats(ctx context.Context) (*dto.SystemInvitationStatsResponse, error)

	// Expiration operations
	GetExpiredInvitations(ctx context.Context, expirationHours int, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)
	DeleteExpiredInvitations(ctx context.Context, expirationHours int) (int64, error)
	GetInvitationsExpiringIn(ctx context.Context, hours int, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)

	// Advanced filtering
	GetInvitationsWithFilters(ctx context.Context, filters dto.InvitationFilterRequest, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error)

	// Preloading operations
	PreloadEvent(ctx context.Context, invitations []*domain.Invitation) error
	PreloadInvitedUser(ctx context.Context, invitations []*domain.Invitation) error
	PreloadAllRelations(ctx context.Context, invitations []*domain.Invitation) error

	// Email-based operations for external users
	GetPendingInvitationsByEmail(ctx context.Context, email string) ([]*domain.Invitation, error)
	GetEventInvitationByEmail(ctx context.Context, eventID int, email string) (*domain.Invitation, error)
	UpdateInvitedUserByEmail(ctx context.Context, email string, userID int) error
}
