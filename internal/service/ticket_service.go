package service

import (
	"context"
	"fmt"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
	"github.com/rs/zerolog"
)

type TicketService interface {
	// Basic CRUD operations
	CreateTicket(ctx context.Context, eventID int, creatorID int, req dto.CreateTicketRequest) (*dto.TicketResponse, error)
	GetTicketByID(ctx context.Context, id int) (*dto.TicketResponse, error)
	UpdateTicket(ctx context.Context, id int, req dto.UpdateTicketRequest) (*dto.TicketResponse, error)
	DeleteTicket(ctx context.Context, id int) error

	// Event-specific operations
	GetTicketsByEventID(ctx context.Context, eventID int) ([]*dto.TicketResponse, error)
	GetActiveTicketsByEventID(ctx context.Context, eventID int) ([]*dto.TicketResponse, error)
	GetAvailableTickets(ctx context.Context, eventID int) ([]*dto.TicketResponse, error)
	GetSoldOutTickets(ctx context.Context, eventID int) ([]*dto.TicketResponse, error)

	// Availability operations
	CheckTicketAvailability(ctx context.Context, ticketID int, quantity int) (bool, error)
	GetTicketAvailability(ctx context.Context, ticketID int) (int, error)

	// Sales operations
	SellTickets(ctx context.Context, ticketID int, quantity int) error
	RefundTickets(ctx context.Context, ticketID int, quantity int) error
	UpdateSoldQuantity(ctx context.Context, ticketID int, quantity int) error

	// Price operations
	GetTicketsByPriceRange(ctx context.Context, eventID int, minPrice, maxPrice float64) ([]*dto.TicketResponse, error)
	GetFreeTickets(ctx context.Context, eventID int) ([]*dto.TicketResponse, error)
	GetPaidTickets(ctx context.Context, eventID int) ([]*dto.TicketResponse, error)

	// Status operations
	ActivateTicket(ctx context.Context, ticketID int) error
	DeactivateTicket(ctx context.Context, ticketID int) error
	ActivateAllEventTickets(ctx context.Context, eventID int) error
	DeactivateAllEventTickets(ctx context.Context, eventID int) error

	// Bulk operations
	CreateMultipleTickets(ctx context.Context, eventID int, requests []dto.CreateTicketRequest) ([]*dto.TicketResponse, error)
	DeleteAllEventTickets(ctx context.Context, eventID int) error

	// Statistics operations
	GetTicketSalesStats(ctx context.Context, eventID int) (*dto.TicketSalesStatsResponse, error)
	GetTicketTypeStats(ctx context.Context, eventID int) ([]*dto.TicketTypeStatsResponse, error)
	GetTotalRevenue(ctx context.Context, eventID int) (float64, error)
	GetTotalTicketsSold(ctx context.Context, eventID int) (int, error)
	GetTotalTicketsAvailable(ctx context.Context, eventID int) (int, error)

	// Advanced filtering
	GetTicketsWithFilters(ctx context.Context, filters dto.TicketFilterRequest, pagination dto.PaginationRequest) ([]*dto.TicketResponse, *dto.PaginationResponse, error)

	// Validation operations
	ValidateTicketOwnership(ctx context.Context, ticketID, eventID int) error
	ValidateTicketData(ctx context.Context, ticket *domain.Ticket) error
}

type ticketService struct {
	ticketRepo repository.TicketRepository
	eventRepo  repository.EventRepository
	logger     zerolog.Logger
}

func NewTicketService(ticketRepo repository.TicketRepository, eventRepo repository.EventRepository, logger zerolog.Logger) TicketService {
	return &ticketService{
		ticketRepo: ticketRepo,
		eventRepo:  eventRepo,
		logger:     logger.With().Str("service", "ticket").Logger(),
	}
}

// Basic CRUD operations
func (s *ticketService) CreateTicket(ctx context.Context, eventID int, creatorID int, req dto.CreateTicketRequest) (*dto.TicketResponse, error) {
	// Validate event exists and belongs to the creator
	isOwner, err := s.eventRepo.IsEventOwner(ctx, eventID, creatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to check event ownership: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("event not found")
	}

	// Validate request
	if err := s.validateCreateTicketRequest(&req); err != nil {
		s.logger.Error().Err(err).Interface("request", req).Msg("Invalid create ticket request")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create domain entity
	ticket := &domain.Ticket{
		EventID:       eventID,
		Title:         req.Title,
		Price:         req.Price,
		TotalQuantity: req.TotalQuantity,
		SoldQuantity:  0,
		IsActive:      true,
	}

	// Validate domain entity
	if err := s.ValidateTicketData(ctx, ticket); err != nil {
		return nil, err
	}

	// Create ticket
	if err := s.ticketRepo.Create(ctx, ticket); err != nil {
		s.logger.Error().Err(err).Interface("ticket", ticket).Msg("Failed to create ticket")
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	s.logger.Info().Int("ticket_id", ticket.ID).Int("event_id", eventID).Str("title", ticket.Title).Msg("Ticket created successfully")

	return s.ticketToResponse(ticket), nil
}

func (s *ticketService) GetTicketByID(ctx context.Context, id int) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("ticket_id", id).Msg("Failed to get ticket by ID")
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return s.ticketToResponse(ticket), nil
}

func (s *ticketService) UpdateTicket(ctx context.Context, id int, req dto.UpdateTicketRequest) (*dto.TicketResponse, error) {
	// Get existing ticket
	existingTicket, err := s.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing ticket: %w", err)
	}

	// Validate request
	if err := s.validateUpdateTicketRequest(&req); err != nil {
		s.logger.Error().Err(err).Interface("request", req).Msg("Invalid update ticket request")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update fields
	if req.Title != nil {
		existingTicket.Title = *req.Title
	}
	if req.Price != nil {
		existingTicket.Price = *req.Price
	}
	if req.TotalQuantity != nil {
		// Validate that new total quantity is not less than sold quantity
		if *req.TotalQuantity < existingTicket.SoldQuantity {
			return nil, fmt.Errorf("total quantity cannot be less than sold quantity (%d)", existingTicket.SoldQuantity)
		}
		existingTicket.TotalQuantity = *req.TotalQuantity
	}

	// Validate updated ticket
	if err := s.ValidateTicketData(ctx, existingTicket); err != nil {
		return nil, err
	}

	// Update ticket
	if err := s.ticketRepo.Update(ctx, existingTicket); err != nil {
		s.logger.Error().Err(err).Int("ticket_id", id).Msg("Failed to update ticket")
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	s.logger.Info().Int("ticket_id", id).Msg("Ticket updated successfully")

	return s.ticketToResponse(existingTicket), nil
}

func (s *ticketService) DeleteTicket(ctx context.Context, id int) error {
	// Check if ticket exists
	exists, err := s.ticketRepo.ExistsByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check ticket existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("ticket not found")
	}

	// Get ticket to check if it has sales
	ticket, err := s.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	// Don't allow deletion if tickets have been sold
	if ticket.SoldQuantity > 0 {
		return fmt.Errorf("cannot delete ticket: %d tickets have been sold", ticket.SoldQuantity)
	}

	// Delete ticket
	if err := s.ticketRepo.Delete(ctx, id); err != nil {
		s.logger.Error().Err(err).Int("ticket_id", id).Msg("Failed to delete ticket")
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	s.logger.Info().Int("ticket_id", id).Msg("Ticket deleted successfully")
	return nil
}

// Event-specific operations
func (s *ticketService) GetTicketsByEventID(ctx context.Context, eventID int) ([]*dto.TicketResponse, error) {
	tickets, err := s.ticketRepo.GetByEventID(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get tickets by event ID")
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}

	var responses []*dto.TicketResponse
	for _, ticket := range tickets {
		responses = append(responses, s.ticketToResponse(ticket))
	}

	return responses, nil
}

func (s *ticketService) GetActiveTicketsByEventID(ctx context.Context, eventID int) ([]*dto.TicketResponse, error) {
	tickets, err := s.ticketRepo.GetActiveByEventID(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get active tickets by event ID")
		return nil, fmt.Errorf("failed to get active tickets: %w", err)
	}

	var responses []*dto.TicketResponse
	for _, ticket := range tickets {
		responses = append(responses, s.ticketToResponse(ticket))
	}

	return responses, nil
}

func (s *ticketService) GetAvailableTickets(ctx context.Context, eventID int) ([]*dto.TicketResponse, error) {
	tickets, err := s.ticketRepo.GetAvailableTickets(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get available tickets")
		return nil, fmt.Errorf("failed to get available tickets: %w", err)
	}

	var responses []*dto.TicketResponse
	for _, ticket := range tickets {
		responses = append(responses, s.ticketToResponse(ticket))
	}

	return responses, nil
}

func (s *ticketService) GetSoldOutTickets(ctx context.Context, eventID int) ([]*dto.TicketResponse, error) {
	tickets, err := s.ticketRepo.GetSoldOutTickets(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get sold out tickets")
		return nil, fmt.Errorf("failed to get sold out tickets: %w", err)
	}

	var responses []*dto.TicketResponse
	for _, ticket := range tickets {
		responses = append(responses, s.ticketToResponse(ticket))
	}

	return responses, nil
}

// Availability operations
func (s *ticketService) CheckTicketAvailability(ctx context.Context, ticketID int, quantity int) (bool, error) {
	available, err := s.ticketRepo.CheckAvailability(ctx, ticketID, quantity)
	if err != nil {
		s.logger.Error().Err(err).Int("ticket_id", ticketID).Int("quantity", quantity).Msg("Failed to check ticket availability")
		return false, fmt.Errorf("failed to check availability: %w", err)
	}

	return available, nil
}

func (s *ticketService) GetTicketAvailability(ctx context.Context, ticketID int) (int, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return 0, fmt.Errorf("failed to get ticket: %w", err)
	}

	available := ticket.TotalQuantity - ticket.SoldQuantity
	if available < 0 {
		available = 0
	}

	return available, nil
}

// Sales operations
func (s *ticketService) SellTickets(ctx context.Context, ticketID int, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	// Check availability
	available, err := s.CheckTicketAvailability(ctx, ticketID, quantity)
	if err != nil {
		return err
	}
	if !available {
		return fmt.Errorf("insufficient tickets available")
	}

	// Increment sold quantity
	if err := s.ticketRepo.IncrementSoldQuantity(ctx, ticketID, quantity); err != nil {
		s.logger.Error().Err(err).Int("ticket_id", ticketID).Int("quantity", quantity).Msg("Failed to sell tickets")
		return fmt.Errorf("failed to sell tickets: %w", err)
	}

	s.logger.Info().Int("ticket_id", ticketID).Int("quantity", quantity).Msg("Tickets sold successfully")
	return nil
}

func (s *ticketService) RefundTickets(ctx context.Context, ticketID int, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	// Get current ticket
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	// Check if we can refund this quantity
	if ticket.SoldQuantity < quantity {
		return fmt.Errorf("cannot refund %d tickets: only %d tickets have been sold", quantity, ticket.SoldQuantity)
	}

	// Decrement sold quantity
	if err := s.ticketRepo.DecrementSoldQuantity(ctx, ticketID, quantity); err != nil {
		s.logger.Error().Err(err).Int("ticket_id", ticketID).Int("quantity", quantity).Msg("Failed to refund tickets")
		return fmt.Errorf("failed to refund tickets: %w", err)
	}

	s.logger.Info().Int("ticket_id", ticketID).Int("quantity", quantity).Msg("Tickets refunded successfully")
	return nil
}

func (s *ticketService) UpdateSoldQuantity(ctx context.Context, ticketID int, quantity int) error {
	if quantity < 0 {
		return fmt.Errorf("sold quantity cannot be negative")
	}

	// Get ticket to validate
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	// Check if new quantity doesn't exceed total
	if quantity > ticket.TotalQuantity {
		return fmt.Errorf("sold quantity (%d) cannot exceed total quantity (%d)", quantity, ticket.TotalQuantity)
	}

	// Update sold quantity
	if err := s.ticketRepo.UpdateSoldQuantity(ctx, ticketID, quantity); err != nil {
		s.logger.Error().Err(err).Int("ticket_id", ticketID).Int("quantity", quantity).Msg("Failed to update sold quantity")
		return fmt.Errorf("failed to update sold quantity: %w", err)
	}

	s.logger.Info().Int("ticket_id", ticketID).Int("quantity", quantity).Msg("Sold quantity updated successfully")
	return nil
}

// Price operations
func (s *ticketService) GetTicketsByPriceRange(ctx context.Context, eventID int, minPrice, maxPrice float64) ([]*dto.TicketResponse, error) {
	if minPrice < 0 || maxPrice < 0 {
		return nil, fmt.Errorf("prices cannot be negative")
	}
	if minPrice > maxPrice {
		return nil, fmt.Errorf("minimum price cannot be greater than maximum price")
	}

	tickets, err := s.ticketRepo.GetTicketsByPriceRange(ctx, eventID, minPrice, maxPrice)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Float64("min_price", minPrice).Float64("max_price", maxPrice).Msg("Failed to get tickets by price range")
		return nil, fmt.Errorf("failed to get tickets by price range: %w", err)
	}

	var responses []*dto.TicketResponse
	for _, ticket := range tickets {
		responses = append(responses, s.ticketToResponse(ticket))
	}

	return responses, nil
}

func (s *ticketService) GetFreeTickets(ctx context.Context, eventID int) ([]*dto.TicketResponse, error) {
	tickets, err := s.ticketRepo.GetFreeTickets(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get free tickets")
		return nil, fmt.Errorf("failed to get free tickets: %w", err)
	}

	var responses []*dto.TicketResponse
	for _, ticket := range tickets {
		responses = append(responses, s.ticketToResponse(ticket))
	}

	return responses, nil
}

func (s *ticketService) GetPaidTickets(ctx context.Context, eventID int) ([]*dto.TicketResponse, error) {
	tickets, err := s.ticketRepo.GetPaidTickets(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get paid tickets")
		return nil, fmt.Errorf("failed to get paid tickets: %w", err)
	}

	var responses []*dto.TicketResponse
	for _, ticket := range tickets {
		responses = append(responses, s.ticketToResponse(ticket))
	}

	return responses, nil
}

// Status operations
func (s *ticketService) ActivateTicket(ctx context.Context, ticketID int) error {
	if err := s.ticketRepo.ActivateTicket(ctx, ticketID); err != nil {
		s.logger.Error().Err(err).Int("ticket_id", ticketID).Msg("Failed to activate ticket")
		return fmt.Errorf("failed to activate ticket: %w", err)
	}

	s.logger.Info().Int("ticket_id", ticketID).Msg("Ticket activated successfully")
	return nil
}

func (s *ticketService) DeactivateTicket(ctx context.Context, ticketID int) error {
	if err := s.ticketRepo.DeactivateTicket(ctx, ticketID); err != nil {
		s.logger.Error().Err(err).Int("ticket_id", ticketID).Msg("Failed to deactivate ticket")
		return fmt.Errorf("failed to deactivate ticket: %w", err)
	}

	s.logger.Info().Int("ticket_id", ticketID).Msg("Ticket deactivated successfully")
	return nil
}

func (s *ticketService) ActivateAllEventTickets(ctx context.Context, eventID int) error {
	if err := s.ticketRepo.ActivateAllEventTickets(ctx, eventID); err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to activate all event tickets")
		return fmt.Errorf("failed to activate all event tickets: %w", err)
	}

	s.logger.Info().Int("event_id", eventID).Msg("All event tickets activated successfully")
	return nil
}

func (s *ticketService) DeactivateAllEventTickets(ctx context.Context, eventID int) error {
	if err := s.ticketRepo.DeactivateAllEventTickets(ctx, eventID); err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to deactivate all event tickets")
		return fmt.Errorf("failed to deactivate all event tickets: %w", err)
	}

	s.logger.Info().Int("event_id", eventID).Msg("All event tickets deactivated successfully")
	return nil
}

// Bulk operations
func (s *ticketService) CreateMultipleTickets(ctx context.Context, eventID int, requests []dto.CreateTicketRequest) ([]*dto.TicketResponse, error) {
	if len(requests) == 0 {
		return nil, fmt.Errorf("no tickets to create")
	}

	// Validate event exists
	exists, err := s.eventRepo.ExistsByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to check event existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("event not found")
	}

	var tickets []*domain.Ticket
	for i, req := range requests {
		// Validate request
		if err := s.validateCreateTicketRequest(&req); err != nil {
			return nil, fmt.Errorf("validation failed for ticket %d: %w", i+1, err)
		}

		ticket := &domain.Ticket{
			EventID:       eventID,
			Title:         req.Title,
			Price:         req.Price,
			TotalQuantity: req.TotalQuantity,
			SoldQuantity:  0,
			IsActive:      true,
		}

		// Validate domain entity
		if err := s.ValidateTicketData(ctx, ticket); err != nil {
			return nil, fmt.Errorf("validation failed for ticket %d: %w", i+1, err)
		}

		tickets = append(tickets, ticket)
	}

	// Create all tickets
	if err := s.ticketRepo.CreateMultiple(ctx, tickets); err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Int("count", len(tickets)).Msg("Failed to create multiple tickets")
		return nil, fmt.Errorf("failed to create tickets: %w", err)
	}

	s.logger.Info().Int("event_id", eventID).Int("count", len(tickets)).Msg("Multiple tickets created successfully")

	var responses []*dto.TicketResponse
	for _, ticket := range tickets {
		responses = append(responses, s.ticketToResponse(ticket))
	}

	return responses, nil
}

func (s *ticketService) DeleteAllEventTickets(ctx context.Context, eventID int) error {
	// Check if any tickets have been sold
	totalSold, err := s.GetTotalTicketsSold(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to check sold tickets: %w", err)
	}
	if totalSold > 0 {
		return fmt.Errorf("cannot delete tickets: %d tickets have been sold", totalSold)
	}

	// Delete all event tickets
	if err := s.ticketRepo.DeleteByEventID(ctx, eventID); err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to delete all event tickets")
		return fmt.Errorf("failed to delete all event tickets: %w", err)
	}

	s.logger.Info().Int("event_id", eventID).Msg("All event tickets deleted successfully")
	return nil
}

// Statistics operations
func (s *ticketService) GetTicketSalesStats(ctx context.Context, eventID int) (*dto.TicketSalesStatsResponse, error) {
	stats, err := s.ticketRepo.GetSalesStats(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get ticket sales stats")
		return nil, fmt.Errorf("failed to get ticket sales stats: %w", err)
	}

	return stats, nil
}

func (s *ticketService) GetTicketTypeStats(ctx context.Context, eventID int) ([]*dto.TicketTypeStatsResponse, error) {
	stats, err := s.ticketRepo.GetTicketTypeStats(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get ticket type stats")
		return nil, fmt.Errorf("failed to get ticket type stats: %w", err)
	}

	return stats, nil
}

func (s *ticketService) GetTotalRevenue(ctx context.Context, eventID int) (float64, error) {
	revenue, err := s.ticketRepo.GetTotalRevenue(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get total revenue")
		return 0, fmt.Errorf("failed to get total revenue: %w", err)
	}

	return revenue, nil
}

func (s *ticketService) GetTotalTicketsSold(ctx context.Context, eventID int) (int, error) {
	sold, err := s.ticketRepo.GetTotalTicketsSold(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get total tickets sold")
		return 0, fmt.Errorf("failed to get total tickets sold: %w", err)
	}

	return sold, nil
}

func (s *ticketService) GetTotalTicketsAvailable(ctx context.Context, eventID int) (int, error) {
	available, err := s.ticketRepo.GetTotalTicketsAvailable(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("event_id", eventID).Msg("Failed to get total tickets available")
		return 0, fmt.Errorf("failed to get total tickets available: %w", err)
	}

	return available, nil
}

// Advanced filtering
func (s *ticketService) GetTicketsWithFilters(ctx context.Context, filters dto.TicketFilterRequest, pagination dto.PaginationRequest) ([]*dto.TicketResponse, *dto.PaginationResponse, error) {
	tickets, paginationResp, err := s.ticketRepo.GetTicketsWithFilters(ctx, filters, pagination)
	if err != nil {
		s.logger.Error().Err(err).Interface("filters", filters).Msg("Failed to get tickets with filters")
		return nil, nil, fmt.Errorf("failed to get tickets with filters: %w", err)
	}

	var responses []*dto.TicketResponse
	for _, ticket := range tickets {
		responses = append(responses, s.ticketToResponse(ticket))
	}

	return responses, paginationResp, nil
}

// Validation operations
func (s *ticketService) ValidateTicketOwnership(ctx context.Context, ticketID, eventID int) error {
	isOwner, err := s.ticketRepo.IsTicketOwner(ctx, ticketID, eventID)
	if err != nil {
		s.logger.Error().Err(err).Int("ticket_id", ticketID).Int("event_id", eventID).Msg("Failed to validate ticket ownership")
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("access denied: ticket does not belong to this event")
	}
	return nil
}

func (s *ticketService) ValidateTicketData(ctx context.Context, ticket *domain.Ticket) error {
	if ticket.Title == "" {
		return fmt.Errorf("ticket title is required")
	}
	if ticket.Price < 0 {
		return fmt.Errorf("ticket price cannot be negative")
	}
	if ticket.TotalQuantity <= 0 {
		return fmt.Errorf("total quantity must be positive")
	}
	if ticket.SoldQuantity < 0 {
		return fmt.Errorf("sold quantity cannot be negative")
	}
	if ticket.SoldQuantity > ticket.TotalQuantity {
		return fmt.Errorf("sold quantity cannot exceed total quantity")
	}

	return nil
}

// Helper methods
func (s *ticketService) validateCreateTicketRequest(req *dto.CreateTicketRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	if req.TotalQuantity <= 0 {
		return fmt.Errorf("total quantity must be positive")
	}

	return nil
}

func (s *ticketService) validateUpdateTicketRequest(req *dto.UpdateTicketRequest) error {
	if req.Title != nil && *req.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if req.Price != nil && *req.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	if req.TotalQuantity != nil && *req.TotalQuantity <= 0 {
		return fmt.Errorf("total quantity must be positive")
	}

	return nil
}

func (s *ticketService) ticketToResponse(ticket *domain.Ticket) *dto.TicketResponse {
	return &dto.TicketResponse{
		ID:            ticket.ID,
		EventID:       ticket.EventID,
		Title:         ticket.Title,
		Price:         ticket.Price,
		TotalQuantity: ticket.TotalQuantity,
		SoldQuantity:  ticket.SoldQuantity,
		IsActive:      ticket.IsActive,
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
	}
}
