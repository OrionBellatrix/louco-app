package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
)

type TicketRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, ticket *domain.Ticket) error
	GetByID(ctx context.Context, id int) (*domain.Ticket, error)
	Update(ctx context.Context, ticket *domain.Ticket) error
	Delete(ctx context.Context, id int) error

	// Event-specific operations
	GetByEventID(ctx context.Context, eventID int) ([]*domain.Ticket, error)
	GetActiveByEventID(ctx context.Context, eventID int) ([]*domain.Ticket, error)
	CountByEventID(ctx context.Context, eventID int) (int64, error)
	CountActiveByEventID(ctx context.Context, eventID int) (int64, error)

	// Availability operations
	GetAvailableTickets(ctx context.Context, eventID int) ([]*domain.Ticket, error)
	GetSoldOutTickets(ctx context.Context, eventID int) ([]*domain.Ticket, error)
	CheckAvailability(ctx context.Context, ticketID int, quantity int) (bool, error)

	// Sales operations
	UpdateSoldQuantity(ctx context.Context, ticketID int, quantity int) error
	IncrementSoldQuantity(ctx context.Context, ticketID int, quantity int) error
	DecrementSoldQuantity(ctx context.Context, ticketID int, quantity int) error
	GetSalesStats(ctx context.Context, eventID int) (*dto.TicketSalesStatsResponse, error)

	// Price operations
	GetTicketsByPriceRange(ctx context.Context, eventID int, minPrice, maxPrice float64) ([]*domain.Ticket, error)
	GetFreeTickets(ctx context.Context, eventID int) ([]*domain.Ticket, error)
	GetPaidTickets(ctx context.Context, eventID int) ([]*domain.Ticket, error)

	// Status operations
	ActivateTicket(ctx context.Context, ticketID int) error
	DeactivateTicket(ctx context.Context, ticketID int) error
	ActivateAllEventTickets(ctx context.Context, eventID int) error
	DeactivateAllEventTickets(ctx context.Context, eventID int) error

	// Validation operations
	ExistsByID(ctx context.Context, id int) (bool, error)
	ExistsByEventAndID(ctx context.Context, eventID, ticketID int) (bool, error)
	IsTicketOwner(ctx context.Context, ticketID, eventID int) (bool, error)

	// Bulk operations
	GetMultipleByIDs(ctx context.Context, ids []int) ([]*domain.Ticket, error)
	CreateMultiple(ctx context.Context, tickets []*domain.Ticket) error
	UpdateMultiple(ctx context.Context, tickets []*domain.Ticket) error
	DeleteMultiple(ctx context.Context, ids []int) error
	DeleteByEventID(ctx context.Context, eventID int) error

	// Statistics operations
	GetTotalRevenue(ctx context.Context, eventID int) (float64, error)
	GetTotalTicketsSold(ctx context.Context, eventID int) (int, error)
	GetTotalTicketsAvailable(ctx context.Context, eventID int) (int, error)
	GetTicketTypeStats(ctx context.Context, eventID int) ([]*dto.TicketTypeStatsResponse, error)

	// Advanced filtering
	GetTicketsWithFilters(ctx context.Context, filters dto.TicketFilterRequest, pagination dto.PaginationRequest) ([]*domain.Ticket, *dto.PaginationResponse, error)

	// Preloading operations
	PreloadEvent(ctx context.Context, tickets []*domain.Ticket) error
}
