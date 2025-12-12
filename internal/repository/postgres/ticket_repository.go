package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
)

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) repository.TicketRepository {
	return &ticketRepository{db: db}
}

// Basic CRUD operations
func (r *ticketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ticketRepository) GetByID(ctx context.Context, id int) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).Preload("Event").First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

func (r *ticketRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Ticket{}, id).Error
}

// Event-specific operations
func (r *ticketRepository) GetByEventID(ctx context.Context, eventID int) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetByEventIDWithPagination(ctx context.Context, eventID int, pagination dto.PaginationRequest) ([]*domain.Ticket, *dto.PaginationResponse, error) {
	var tickets []*domain.Ticket
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Ticket{}).Where("event_id = ?", eventID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("created_at ASC").
		Find(&tickets).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return tickets, paginationResponse, nil
}

func (r *ticketRepository) CountByEventID(ctx context.Context, eventID int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}

func (r *ticketRepository) DeleteByEventID(ctx context.Context, eventID int) error {
	return r.db.WithContext(ctx).Where("event_id = ?", eventID).Delete(&domain.Ticket{}).Error
}

// Active ticket operations (missing methods from interface)
func (r *ticketRepository) GetActiveByEventID(ctx context.Context, eventID int) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	err := r.db.WithContext(ctx).Where("event_id = ? AND is_active = ?", eventID, true).Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) CountActiveByEventID(ctx context.Context, eventID int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).Where("event_id = ? AND is_active = ?", eventID, true).Count(&count).Error
	return count, err
}

// Availability operations
func (r *ticketRepository) GetAvailableTickets(ctx context.Context, eventID int) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND sold_quantity < total_quantity", eventID).
		Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetSoldOutTickets(ctx context.Context, eventID int) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND sold_quantity >= total_quantity", eventID).
		Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) CheckAvailability(ctx context.Context, ticketID int, quantity int) (bool, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).First(&ticket, ticketID).Error
	if err != nil {
		return false, err
	}

	availableQuantity := ticket.TotalQuantity - ticket.SoldQuantity
	return availableQuantity >= quantity, nil
}

func (r *ticketRepository) GetAvailableQuantity(ctx context.Context, ticketID int) (int, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).First(&ticket, ticketID).Error
	if err != nil {
		return 0, err
	}

	return ticket.TotalQuantity - ticket.SoldQuantity, nil
}

// Sales operations
func (r *ticketRepository) UpdateSoldQuantity(ctx context.Context, ticketID int, quantity int) error {
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("id = ?", ticketID).
		Update("sold_quantity", gorm.Expr("sold_quantity + ?", quantity)).Error
}

func (r *ticketRepository) IncrementSoldQuantity(ctx context.Context, ticketID int, quantity int) error {
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("id = ? AND sold_quantity + ? <= total_quantity", ticketID, quantity).
		Update("sold_quantity", gorm.Expr("sold_quantity + ?", quantity)).Error
}

func (r *ticketRepository) DecrementSoldQuantity(ctx context.Context, ticketID int, quantity int) error {
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("id = ? AND sold_quantity >= ?", ticketID, quantity).
		Update("sold_quantity", gorm.Expr("sold_quantity - ?", quantity)).Error
}

func (r *ticketRepository) GetTotalSoldQuantity(ctx context.Context, eventID int) (int, error) {
	var totalSold int
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ?", eventID).
		Select("COALESCE(SUM(sold_quantity), 0)").
		Scan(&totalSold).Error
	return totalSold, err
}

func (r *ticketRepository) GetTotalRevenue(ctx context.Context, eventID int) (float64, error) {
	var totalRevenue float64
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ?", eventID).
		Select("COALESCE(SUM(price * sold_quantity), 0)").
		Scan(&totalRevenue).Error
	return totalRevenue, err
}

// Price operations
func (r *ticketRepository) GetTicketsByPriceRange(ctx context.Context, eventID int, minPrice, maxPrice float64) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	query := r.db.WithContext(ctx).Where("event_id = ?", eventID)

	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	err := query.Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetFreeTickets(ctx context.Context, eventID int) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	err := r.db.WithContext(ctx).Where("event_id = ? AND price = 0", eventID).Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetPaidTickets(ctx context.Context, eventID int) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	err := r.db.WithContext(ctx).Where("event_id = ? AND price > 0", eventID).Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetCheapestTicket(ctx context.Context, eventID int) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND sold_quantity < total_quantity", eventID).
		Order("price ASC").
		First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetMostExpensiveTicket(ctx context.Context, eventID int) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND sold_quantity < total_quantity", eventID).
		Order("price DESC").
		First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetAveragePrice(ctx context.Context, eventID int) (float64, error) {
	var avgPrice float64
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ?", eventID).
		Select("COALESCE(AVG(price), 0)").
		Scan(&avgPrice).Error
	return avgPrice, err
}

// Status operations (missing methods from interface)
func (r *ticketRepository) ActivateTicket(ctx context.Context, ticketID int) error {
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).Where("id = ?", ticketID).Update("is_active", true).Error
}

func (r *ticketRepository) DeactivateTicket(ctx context.Context, ticketID int) error {
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).Where("id = ?", ticketID).Update("is_active", false).Error
}

func (r *ticketRepository) ActivateAllEventTickets(ctx context.Context, eventID int) error {
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).Where("event_id = ?", eventID).Update("is_active", true).Error
}

func (r *ticketRepository) DeactivateAllEventTickets(ctx context.Context, eventID int) error {
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).Where("event_id = ?", eventID).Update("is_active", false).Error
}

// Search operations
func (r *ticketRepository) SearchByTitle(ctx context.Context, eventID int, title string) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND title ILIKE ?", eventID, "%"+title+"%").
		Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetTicketsWithFilters(ctx context.Context, filters dto.TicketFilterRequest, pagination dto.PaginationRequest) ([]*domain.Ticket, *dto.PaginationResponse, error) {
	var tickets []*domain.Ticket
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Ticket{})

	// Apply filters
	if filters.EventID != nil {
		query = query.Where("event_id = ?", *filters.EventID)
	}
	if filters.Query != nil && *filters.Query != "" {
		query = query.Where("title ILIKE ?", "%"+*filters.Query+"%")
	}
	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}
	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.IsSoldOut != nil && *filters.IsSoldOut {
		query = query.Where("sold_quantity >= total_quantity")
	}
	if filters.IsFree != nil && *filters.IsFree {
		query = query.Where("price = 0")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("price ASC").
		Find(&tickets).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return tickets, paginationResponse, nil
}

// Statistics operations (interface methods)
func (r *ticketRepository) GetSalesStats(ctx context.Context, eventID int) (*dto.TicketSalesStatsResponse, error) {
	var stats dto.TicketSalesStatsResponse
	stats.EventID = eventID

	// Total tickets count
	var totalTickets int
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ?", eventID).
		Select("COALESCE(SUM(total_quantity), 0)").
		Scan(&totalTickets).Error
	if err != nil {
		return nil, err
	}
	stats.TotalTickets = totalTickets

	// Sold tickets count
	var soldTickets int
	err = r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ?", eventID).
		Select("COALESCE(SUM(sold_quantity), 0)").
		Scan(&soldTickets).Error
	if err != nil {
		return nil, err
	}
	stats.SoldTickets = soldTickets

	// Available tickets
	stats.AvailableTickets = totalTickets - soldTickets

	// Total revenue
	var totalRevenue float64
	err = r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ?", eventID).
		Select("COALESCE(SUM(price * sold_quantity), 0)").
		Scan(&totalRevenue).Error
	if err != nil {
		return nil, err
	}
	stats.TotalRevenue = totalRevenue

	// Average price
	var avgPrice float64
	err = r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ?", eventID).
		Select("COALESCE(AVG(price), 0)").
		Scan(&avgPrice).Error
	if err != nil {
		return nil, err
	}
	stats.AveragePrice = avgPrice

	// Sold percentage
	if totalTickets > 0 {
		stats.SoldPercentage = (float64(soldTickets) / float64(totalTickets)) * 100
	}

	return &stats, nil
}

func (r *ticketRepository) GetTotalTicketsSold(ctx context.Context, eventID int) (int, error) {
	var totalSold int
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ?", eventID).
		Select("COALESCE(SUM(sold_quantity), 0)").
		Scan(&totalSold).Error
	return totalSold, err
}

func (r *ticketRepository) GetTotalTicketsAvailable(ctx context.Context, eventID int) (int, error) {
	var totalQuantity, soldQuantity int

	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ?", eventID).
		Select("COALESCE(SUM(total_quantity), 0), COALESCE(SUM(sold_quantity), 0)").
		Row().Scan(&totalQuantity, &soldQuantity)
	if err != nil {
		return 0, err
	}

	return totalQuantity - soldQuantity, nil
}

func (r *ticketRepository) GetTicketTypeStats(ctx context.Context, eventID int) ([]*dto.TicketTypeStatsResponse, error) {
	var tickets []*domain.Ticket
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&tickets).Error
	if err != nil {
		return nil, err
	}

	var stats []*dto.TicketTypeStatsResponse
	for _, ticket := range tickets {
		stat := &dto.TicketTypeStatsResponse{
			TicketID:          ticket.ID,
			Title:             ticket.Title,
			Price:             ticket.Price,
			TotalQuantity:     ticket.TotalQuantity,
			SoldQuantity:      ticket.SoldQuantity,
			AvailableQuantity: ticket.TotalQuantity - ticket.SoldQuantity,
			Revenue:           ticket.Price * float64(ticket.SoldQuantity),
		}

		if ticket.TotalQuantity > 0 {
			stat.SoldPercentage = (float64(ticket.SoldQuantity) / float64(ticket.TotalQuantity)) * 100
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

// Validation operations (interface methods)
func (r *ticketRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *ticketRepository) ExistsByEventAndID(ctx context.Context, eventID, ticketID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("event_id = ? AND id = ?", eventID, ticketID).
		Count(&count).Error
	return count > 0, err
}

func (r *ticketRepository) IsTicketOwner(ctx context.Context, ticketID, eventID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("id = ? AND event_id = ?", ticketID, eventID).
		Count(&count).Error
	return count > 0, err
}

// Bulk operations
func (r *ticketRepository) GetMultipleByIDs(ctx context.Context, ids []int) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) CreateMultiple(ctx context.Context, tickets []*domain.Ticket) error {
	return r.db.WithContext(ctx).Create(&tickets).Error
}

func (r *ticketRepository) DeleteMultiple(ctx context.Context, ids []int) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.Ticket{}).Error
}

func (r *ticketRepository) UpdateMultiple(ctx context.Context, tickets []*domain.Ticket) error {
	return r.db.WithContext(ctx).Save(&tickets).Error
}

func (r *ticketRepository) UpdateMultiplePrices(ctx context.Context, ticketIDs []int, newPrice float64) error {
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("id IN ?", ticketIDs).
		Update("price", newPrice).Error
}

// Preloading operations (interface method)
func (r *ticketRepository) PreloadEvent(ctx context.Context, tickets []*domain.Ticket) error {
	for i := range tickets {
		err := r.db.WithContext(ctx).Preload("Event").First(tickets[i], tickets[i].ID).Error
		if err != nil {
			return err
		}
	}
	return nil
}
