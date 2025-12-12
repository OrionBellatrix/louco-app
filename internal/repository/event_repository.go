package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
)

type EventRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, event *domain.Event) error
	GetByID(ctx context.Context, id int) (*domain.Event, error)
	GetByIDWithRelations(ctx context.Context, id int) (*domain.Event, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id int) error

	// Creator-specific operations
	GetByCreatorID(ctx context.Context, creatorID int, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	GetByCreatorIDAndStatus(ctx context.Context, creatorID int, status domain.EventStatus, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	CountByCreatorID(ctx context.Context, creatorID int) (int64, error)
	CountByCreatorIDAndStatus(ctx context.Context, creatorID int, status domain.EventStatus) (int64, error)

	// Public event operations
	GetPublicEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	GetPublicEventsByCategory(ctx context.Context, categoryID int, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	GetPublicEventsByLocation(ctx context.Context, city string, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	SearchPublicEvents(ctx context.Context, query string, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)

	// Status management operations
	GetByStatus(ctx context.Context, status domain.EventStatus, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	UpdateStatus(ctx context.Context, id int, status domain.EventStatus) error
	GetPendingEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)

	// Category operations
	AddCategories(ctx context.Context, eventID int, categoryIDs []int) error
	RemoveCategories(ctx context.Context, eventID int, categoryIDs []int) error
	UpdateCategories(ctx context.Context, eventID int, categoryIDs []int) error
	GetEventsByCategories(ctx context.Context, categoryIDs []int, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)

	// Private event operations
	GetPrivateEventsByInvitedUser(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	GetPrivateEventsByInvitedEmail(ctx context.Context, email string, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)

	// Date-based operations
	GetEventsByDateRange(ctx context.Context, startDate, endDate string, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	GetUpcomingEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	GetPastEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)

	// Location-based operations
	GetEventsByCity(ctx context.Context, city string, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	GetEventsByCoordinates(ctx context.Context, latitude, longitude float64, radiusKm int, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)

	// Type-based operations
	GetEventsByType(ctx context.Context, eventType domain.EventType, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	GetEventsByLocationType(ctx context.Context, locationType domain.EventLocationType, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)

	// Statistics operations
	GetEventStats(ctx context.Context, creatorID int) (*dto.EventStatsResponse, error)
	GetSystemEventStats(ctx context.Context) (*dto.SystemEventStatsResponse, error)

	// Validation operations
	ExistsByID(ctx context.Context, id int) (bool, error)
	ExistsByCreatorAndID(ctx context.Context, creatorID, eventID int) (bool, error)
	IsEventOwner(ctx context.Context, eventID, creatorID int) (bool, error)

	// Bulk operations
	GetMultipleByIDs(ctx context.Context, ids []int) ([]*domain.Event, error)
	UpdateMultipleStatus(ctx context.Context, ids []int, status domain.EventStatus) error
	DeleteMultiple(ctx context.Context, ids []int) error

	// Advanced filtering
	GetEventsWithFilters(ctx context.Context, filters dto.EventFilterRequest, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error)
	GetFeaturedEvents(ctx context.Context, limit int) ([]*domain.Event, error)
	GetTrendingEvents(ctx context.Context, limit int) ([]*domain.Event, error)

	// Preloading operations
	PreloadCategories(ctx context.Context, events []*domain.Event) error
	PreloadTickets(ctx context.Context, events []*domain.Event) error
	PreloadInvitations(ctx context.Context, events []*domain.Event) error
	PreloadAddress(ctx context.Context, events []*domain.Event) error
	PreloadMedia(ctx context.Context, events []*domain.Event) error
	PreloadCreator(ctx context.Context, events []*domain.Event) error
	PreloadAllRelations(ctx context.Context, events []*domain.Event) error
}
