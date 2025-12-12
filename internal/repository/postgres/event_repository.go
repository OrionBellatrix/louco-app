package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) repository.EventRepository {
	return &eventRepository{db: db}
}

// Basic CRUD operations
func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *eventRepository) GetByID(ctx context.Context, id int) (*domain.Event, error) {
	var event domain.Event
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Creator.Industries").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Preload("Invitations").
		Preload("Invitations.InvitedUser").
		First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) GetByIDWithRelations(ctx context.Context, id int) (*domain.Event, error) {
	var event domain.Event
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Creator.Industries").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Preload("Invitations").
		Preload("Invitations.InvitedUser").
		First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) Update(ctx context.Context, event *domain.Event) error {
	return r.db.WithContext(ctx).Save(event).Error
}

func (r *eventRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Event{}, id).Error
}

// Creator-specific operations
func (r *eventRepository) GetByCreatorID(ctx context.Context, creatorID int, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ?", creatorID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.Industries").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Preload("Invitations").
		Preload("Invitations.InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) GetByCreatorIDAndStatus(ctx context.Context, creatorID int, status domain.EventStatus, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("creator_id = ? AND status = ?", creatorID, status)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.Industries").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Preload("Invitations").
		Preload("Invitations.InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) CountByCreatorID(ctx context.Context, creatorID int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("creator_id = ?", creatorID).
		Count(&count).Error
	return count, err
}

func (r *eventRepository) CountByCreatorIDAndStatus(ctx context.Context, creatorID int, status domain.EventStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("creator_id = ? AND status = ?", creatorID, status).
		Count(&count).Error
	return count, err
}

// Public event operations
func (r *eventRepository) GetPublicEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("type = ? AND status = ?", domain.EventTypePublic, domain.EventStatusPublished)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) GetPublicEventsByCategory(ctx context.Context, categoryID int, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Joins("JOIN event_categories ON events.id = event_categories.event_id").
		Where("event_categories.category_id = ? AND events.type = ? AND events.status = ?",
			categoryID, domain.EventTypePublic, domain.EventStatusPublished)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Offset(offset).
		Limit(pageSize).
		Order("events.created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) GetPublicEventsByLocation(ctx context.Context, city string, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Joins("JOIN addresses ON events.address_id = addresses.id").
		Where("addresses.city ILIKE ? AND events.type = ? AND events.status = ?",
			"%"+city+"%", domain.EventTypePublic, domain.EventStatusPublished)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Offset(offset).
		Limit(pageSize).
		Order("events.created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) SearchPublicEvents(ctx context.Context, query string, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	searchQuery := r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("(name ILIKE ? OR description ILIKE ?) AND type = ? AND status = ?",
			"%"+query+"%", "%"+query+"%", domain.EventTypePublic, domain.EventStatusPublished)

	// Count total
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := searchQuery.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

// Status management operations
func (r *eventRepository) GetByStatus(ctx context.Context, status domain.EventStatus, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).Where("status = ?", status)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Preload("Invitations").
		Preload("Invitations.InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) UpdateStatus(ctx context.Context, id int, status domain.EventStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *eventRepository) GetPendingEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	return r.GetByStatus(ctx, domain.EventStatusPending, pagination)
}

// Category operations
func (r *eventRepository) AddCategories(ctx context.Context, eventID int, categoryIDs []int) error {
	var eventCategories []domain.EventCategory
	for _, categoryID := range categoryIDs {
		eventCategories = append(eventCategories, domain.EventCategory{
			EventID:    eventID,
			CategoryID: categoryID,
			CreatedAt:  time.Now(),
		})
	}
	return r.db.WithContext(ctx).Create(&eventCategories).Error
}

func (r *eventRepository) RemoveCategories(ctx context.Context, eventID int, categoryIDs []int) error {
	return r.db.WithContext(ctx).
		Where("event_id = ? AND category_id IN ?", eventID, categoryIDs).
		Delete(&domain.EventCategory{}).Error
}

func (r *eventRepository) UpdateCategories(ctx context.Context, eventID int, categoryIDs []int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remove all existing categories
		if err := tx.Where("event_id = ?", eventID).Delete(&domain.EventCategory{}).Error; err != nil {
			return err
		}

		// Add new categories
		if len(categoryIDs) > 0 {
			var eventCategories []domain.EventCategory
			for _, categoryID := range categoryIDs {
				eventCategories = append(eventCategories, domain.EventCategory{
					EventID:    eventID,
					CategoryID: categoryID,
					CreatedAt:  time.Now(),
				})
			}
			return tx.Create(&eventCategories).Error
		}

		return nil
	})
}

func (r *eventRepository) GetEventsByCategories(ctx context.Context, categoryIDs []int, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Joins("JOIN event_categories ON events.id = event_categories.event_id").
		Where("event_categories.category_id IN ? AND events.status = ?", categoryIDs, domain.EventStatusPublished)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Offset(offset).
		Limit(pageSize).
		Order("events.created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

// Private event operations
func (r *eventRepository) GetPrivateEventsByInvitedUser(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Joins("JOIN invitations ON events.id = invitations.event_id").
		Where("invitations.invited_user_id = ? AND events.type = ? AND events.status = ?",
			userID, domain.EventTypePrivate, domain.EventStatusPublished)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.Industries").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Preload("Invitations").
		Preload("Invitations.InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("events.created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) GetPrivateEventsByInvitedEmail(ctx context.Context, email string, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Joins("JOIN invitations ON events.id = invitations.event_id").
		Where("invitations.invited_email = ? AND events.type = ? AND events.status = ?",
			email, domain.EventTypePrivate, domain.EventStatusPublished)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.Industries").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Preload("Invitations").
		Preload("Invitations.InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("events.created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

// Validation operations
func (r *eventRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Event{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *eventRepository) ExistsByCreatorAndID(ctx context.Context, creatorID, eventID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("id = ? AND creator_id = ?", eventID, creatorID).
		Count(&count).Error
	return count > 0, err
}

func (r *eventRepository) IsEventOwner(ctx context.Context, eventID, creatorID int) (bool, error) {
	return r.ExistsByCreatorAndID(ctx, creatorID, eventID)
}

// Additional helper methods will be implemented in the next part...
// This is a partial implementation focusing on the core functionality

// Date-based operations
func (r *eventRepository) GetEventsByDateRange(ctx context.Context, startDate, endDate string, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("start_date >= ? AND start_date <= ? AND status = ?", startDate, endDate, domain.EventStatusPublished)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Address").
		Preload("Categories").
		Offset(offset).
		Limit(pageSize).
		Order("start_date ASC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) GetUpcomingEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	today := time.Now().Format("2006-01-02")
	return r.GetEventsByDateRange(ctx, today, "2099-12-31", pagination)
}

func (r *eventRepository) GetPastEvents(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	today := time.Now().Format("2006-01-02")
	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("start_date < ? AND status = ?", today, domain.EventStatusPublished)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Address").
		Preload("Categories").
		Offset(offset).
		Limit(pageSize).
		Order("start_date DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

// Bulk operations
func (r *eventRepository) GetMultipleByIDs(ctx context.Context, ids []int) ([]*domain.Event, error) {
	var events []*domain.Event
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Image").
		Preload("Address").
		Preload("Categories").
		Where("id IN ?", ids).
		Find(&events).Error
	return events, err
}

func (r *eventRepository) UpdateMultipleStatus(ctx context.Context, ids []int, status domain.EventStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *eventRepository) DeleteMultiple(ctx context.Context, ids []int) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.Event{}).Error
}

// Statistics operations
func (r *eventRepository) GetEventStats(ctx context.Context, creatorID int) (*dto.EventStatsResponse, error) {
	stats := &dto.EventStatsResponse{}

	// Total events
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ?", creatorID).Count(&stats.TotalEvents)

	// Events by status
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ? AND status = ?", creatorID, domain.EventStatusDraft).Count(&stats.DraftEvents)
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ? AND status = ?", creatorID, domain.EventStatusPending).Count(&stats.PendingEvents)
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ? AND status = ?", creatorID, domain.EventStatusPublished).Count(&stats.PublishedEvents)
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ? AND status = ?", creatorID, domain.EventStatusRejected).Count(&stats.RejectedEvents)
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ? AND status = ?", creatorID, domain.EventStatusCancelled).Count(&stats.CancelledEvents)

	// Events by type
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ? AND type = ?", creatorID, domain.EventTypePublic).Count(&stats.PublicEvents)
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ? AND type = ?", creatorID, domain.EventTypePrivate).Count(&stats.PrivateEvents)

	// Events by location type
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ? AND location_type = ?", creatorID, domain.EventLocationTypeLocation).Count(&stats.LocationEvents)
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ? AND location_type = ?", creatorID, domain.EventLocationTypeOnline).Count(&stats.OnlineEvents)
	r.db.WithContext(ctx).Model(&domain.Event{}).Where("creator_id = ? AND location_type = ?", creatorID, domain.EventLocationTypeAnnouncement).Count(&stats.AnnouncementEvents)

	return stats, nil
}

func (r *eventRepository) GetSystemEventStats(ctx context.Context) (*dto.SystemEventStatsResponse, error) {
	stats := &dto.SystemEventStatsResponse{}

	// Get basic event stats
	eventStats, err := r.GetEventStats(ctx, 0) // 0 means all creators
	if err != nil {
		return nil, err
	}
	stats.EventStatsResponse = *eventStats

	// Additional system stats
	r.db.WithContext(ctx).Model(&domain.Creator{}).Count(&stats.TotalCreators)
	r.db.WithContext(ctx).Model(&domain.Creator{}).Joins("JOIN users ON creators.user_id = users.id").Where("users.is_active = ?", true).Count(&stats.ActiveCreators)

	return stats, nil
}

// Advanced filtering with complex where conditions
func (r *eventRepository) GetEventsWithFilters(ctx context.Context, filters dto.EventFilterRequest, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{})

	// Apply filters
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.LocationType != nil {
		query = query.Where("location_type = ?", *filters.LocationType)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.CreatorID != nil {
		query = query.Where("creator_id = ?", *filters.CreatorID)
	}
	if filters.StartDate != nil {
		query = query.Where("start_date >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("start_date <= ?", *filters.EndDate)
	}
	if filters.Query != nil && *filters.Query != "" {
		query = query.Where("(name ILIKE ? OR description ILIKE ?)", "%"+*filters.Query+"%", "%"+*filters.Query+"%")
	}

	// Category filter
	if len(filters.CategoryIDs) > 0 {
		query = query.Joins("JOIN event_categories ON events.id = event_categories.event_id").
			Where("event_categories.category_id IN ?", filters.CategoryIDs)
	}

	// Location filters
	if filters.City != nil {
		query = query.Joins("JOIN addresses ON events.address_id = addresses.id").
			Where("addresses.city ILIKE ?", "%"+*filters.City+"%")
	}
	if filters.Country != nil {
		query = query.Joins("JOIN addresses ON events.address_id = addresses.id").
			Where("addresses.country ILIKE ?", "%"+*filters.Country+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Address").
		Preload("Categories").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

// Preloading operations
func (r *eventRepository) PreloadCategories(ctx context.Context, events []*domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	var eventIDs []int
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

	return r.db.WithContext(ctx).Preload("Categories").Where("id IN ?", eventIDs).Find(events).Error
}

func (r *eventRepository) PreloadTickets(ctx context.Context, events []*domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	var eventIDs []int
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

	return r.db.WithContext(ctx).Preload("Tickets").Where("id IN ?", eventIDs).Find(events).Error
}

func (r *eventRepository) PreloadInvitations(ctx context.Context, events []*domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	var eventIDs []int
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

	return r.db.WithContext(ctx).Preload("Invitations").Where("id IN ?", eventIDs).Find(events).Error
}

func (r *eventRepository) PreloadAddress(ctx context.Context, events []*domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	var eventIDs []int
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

	return r.db.WithContext(ctx).Preload("Address").Where("id IN ?", eventIDs).Find(events).Error
}

func (r *eventRepository) PreloadMedia(ctx context.Context, events []*domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	var eventIDs []int
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

	return r.db.WithContext(ctx).Preload("Image").Preload("Video").Where("id IN ?", eventIDs).Find(events).Error
}

func (r *eventRepository) PreloadCreator(ctx context.Context, events []*domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	var eventIDs []int
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

	return r.db.WithContext(ctx).Preload("Creator").Preload("Creator.User").Where("id IN ?", eventIDs).Find(events).Error
}

func (r *eventRepository) PreloadAllRelations(ctx context.Context, events []*domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	var eventIDs []int
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

	return r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Preload("Invitations").
		Where("id IN ?", eventIDs).
		Find(events).Error
}

// Additional methods for location and type operations
func (r *eventRepository) GetEventsByCity(ctx context.Context, city string, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	return r.GetPublicEventsByLocation(ctx, city, pagination)
}

func (r *eventRepository) GetEventsByCoordinates(ctx context.Context, latitude, longitude float64, radiusKm int, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	// Using Haversine formula for distance calculation
	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Joins("JOIN addresses ON events.address_id = addresses.id").
		Where("events.status = ? AND (6371 * acos(cos(radians(?)) * cos(radians(addresses.latitude)) * cos(radians(addresses.longitude) - radians(?)) + sin(radians(?)) * sin(radians(addresses.latitude)))) <= ?",
			domain.EventStatusPublished, latitude, longitude, latitude, radiusKm)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Address").
		Preload("Categories").
		Offset(offset).
		Limit(pageSize).
		Order("events.created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) GetEventsByType(ctx context.Context, eventType domain.EventType, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("type = ? AND status = ?", eventType, domain.EventStatusPublished)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Video").
		Preload("Address").
		Preload("Categories").
		Preload("Tickets").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) GetEventsByLocationType(ctx context.Context, locationType domain.EventLocationType, pagination dto.PaginationRequest) ([]*domain.Event, *dto.PaginationResponse, error) {
	var events []*domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("location_type = ? AND status = ?", locationType, domain.EventStatusPublished)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Address").
		Preload("Categories").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return events, paginationResponse, nil
}

// Featured and trending events
func (r *eventRepository) GetFeaturedEvents(ctx context.Context, limit int) ([]*domain.Event, error) {
	var events []*domain.Event
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Address").
		Preload("Categories").
		Where("status = ? AND type = ?", domain.EventStatusPublished, domain.EventTypePublic).
		Order("created_at DESC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *eventRepository) GetTrendingEvents(ctx context.Context, limit int) ([]*domain.Event, error) {
	var events []*domain.Event
	// For now, we'll use creation date as trending indicator
	// In a real implementation, this could be based on views, likes, ticket sales, etc.
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Creator.User").
		Preload("Image").
		Preload("Address").
		Preload("Categories").
		Where("status = ? AND type = ? AND created_at >= ?",
			domain.EventStatusPublished, domain.EventTypePublic, time.Now().AddDate(0, 0, -30)).
		Order("created_at DESC").
		Limit(limit).
		Find(&events).Error
	return events, err
}
