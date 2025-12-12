package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
)

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) repository.InvitationRepository {
	return &invitationRepository{db: db}
}

// Basic CRUD operations
func (r *invitationRepository) Create(ctx context.Context, invitation *domain.Invitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *invitationRepository) GetByID(ctx context.Context, id int) (*domain.Invitation, error) {
	var invitation domain.Invitation
	err := r.db.WithContext(ctx).First(&invitation, id).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *invitationRepository) GetByIDWithRelations(ctx context.Context, id int) (*domain.Invitation, error) {
	var invitation domain.Invitation
	err := r.db.WithContext(ctx).
		Preload("Event").
		Preload("Event.Creator").
		Preload("InvitedUser").
		First(&invitation, id).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *invitationRepository) Update(ctx context.Context, invitation *domain.Invitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}

func (r *invitationRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Invitation{}, id).Error
}

// Event-specific operations
func (r *invitationRepository) GetByEventID(ctx context.Context, eventID int, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	var invitations []*domain.Invitation
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Invitation{}).Where("event_id = ?", eventID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return invitations, paginationResponse, nil
}

func (r *invitationRepository) GetByEventIDWithRelations(ctx context.Context, eventID int, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	var invitations []*domain.Invitation
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Invitation{}).Where("event_id = ?", eventID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Event").
		Preload("Event.Creator").
		Preload("InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return invitations, paginationResponse, nil
}

func (r *invitationRepository) CountByEventID(ctx context.Context, eventID int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}

func (r *invitationRepository) CountByEventIDAndStatus(ctx context.Context, eventID int, status domain.InvitationStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("event_id = ? AND status = ?", eventID, status).
		Count(&count).Error
	return count, err
}

// User-specific operations
func (r *invitationRepository) GetByUserID(ctx context.Context, userID int, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	var invitations []*domain.Invitation
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Invitation{}).Where("invited_user_id = ?", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Event").
		Preload("Event.Creator").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return invitations, paginationResponse, nil
}

func (r *invitationRepository) GetByEmail(ctx context.Context, email string, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	var invitations []*domain.Invitation
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Invitation{}).Where("invited_email = ?", email)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Event").
		Preload("Event.Creator").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return invitations, paginationResponse, nil
}

func (r *invitationRepository) GetByUserIDAndStatus(ctx context.Context, userID int, status domain.InvitationStatus, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	var invitations []*domain.Invitation
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Invitation{}).Where("invited_user_id = ? AND status = ?", userID, status)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Event").
		Preload("Event.Creator").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return invitations, paginationResponse, nil
}

func (r *invitationRepository) GetByEmailAndStatus(ctx context.Context, email string, status domain.InvitationStatus, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	var invitations []*domain.Invitation
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Invitation{}).Where("invited_email = ? AND status = ?", email, status)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Event").
		Preload("Event.Creator").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return invitations, paginationResponse, nil
}

// Status-specific operations
func (r *invitationRepository) GetByStatus(ctx context.Context, status domain.InvitationStatus, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	var invitations []*domain.Invitation
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Invitation{}).Where("status = ?", status)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Event").
		Preload("InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return invitations, paginationResponse, nil
}

func (r *invitationRepository) GetPendingInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	return r.GetByStatus(ctx, domain.InvitationStatusPending, pagination)
}

func (r *invitationRepository) GetApprovedInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	return r.GetByStatus(ctx, domain.InvitationStatusApproved, pagination)
}

func (r *invitationRepository) GetRejectedInvitations(ctx context.Context, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	return r.GetByStatus(ctx, domain.InvitationStatusRejected, pagination)
}

// Status management operations
func (r *invitationRepository) UpdateStatus(ctx context.Context, id int, status domain.InvitationStatus) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"responded_at": &now,
			"updated_at":   now,
		}).Error
}

func (r *invitationRepository) ApproveInvitation(ctx context.Context, id int) error {
	return r.UpdateStatus(ctx, id, domain.InvitationStatusApproved)
}

func (r *invitationRepository) RejectInvitation(ctx context.Context, id int) error {
	return r.UpdateStatus(ctx, id, domain.InvitationStatusRejected)
}

func (r *invitationRepository) ResetToPending(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       domain.InvitationStatusPending,
			"responded_at": nil,
			"updated_at":   time.Now(),
		}).Error
}

// Duplicate and validation operations
func (r *invitationRepository) ExistsByEventAndEmail(ctx context.Context, eventID int, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("event_id = ? AND invited_email = ?", eventID, email).
		Count(&count).Error
	return count > 0, err
}

func (r *invitationRepository) ExistsByEventAndUser(ctx context.Context, eventID int, userID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("event_id = ? AND invited_user_id = ?", eventID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *invitationRepository) GetByEventAndEmail(ctx context.Context, eventID int, email string) (*domain.Invitation, error) {
	var invitation domain.Invitation
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND invited_email = ?", eventID, email).
		First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *invitationRepository) GetByEventAndUser(ctx context.Context, eventID int, userID int) (*domain.Invitation, error) {
	var invitation domain.Invitation
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND invited_user_id = ?", eventID, userID).
		First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// Validation operations
func (r *invitationRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *invitationRepository) IsInvitationOwner(ctx context.Context, invitationID, eventID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("id = ? AND event_id = ?", invitationID, eventID).
		Count(&count).Error
	return count > 0, err
}

func (r *invitationRepository) CanUserAccessInvitation(ctx context.Context, invitationID, userID int, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("id = ? AND (invited_user_id = ? OR invited_email = ?)", invitationID, userID, email).
		Count(&count).Error
	return count > 0, err
}

// Bulk operations
func (r *invitationRepository) GetMultipleByIDs(ctx context.Context, ids []int) ([]*domain.Invitation, error) {
	var invitations []*domain.Invitation
	err := r.db.WithContext(ctx).
		Preload("Event").
		Preload("InvitedUser").
		Where("id IN ?", ids).
		Find(&invitations).Error
	return invitations, err
}

func (r *invitationRepository) CreateMultiple(ctx context.Context, invitations []*domain.Invitation) error {
	return r.db.WithContext(ctx).Create(&invitations).Error
}

func (r *invitationRepository) UpdateMultiple(ctx context.Context, invitations []*domain.Invitation) error {
	return r.db.WithContext(ctx).Save(&invitations).Error
}

func (r *invitationRepository) DeleteMultiple(ctx context.Context, ids []int) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.Invitation{}).Error
}

func (r *invitationRepository) DeleteByEventID(ctx context.Context, eventID int) error {
	return r.db.WithContext(ctx).Where("event_id = ?", eventID).Delete(&domain.Invitation{}).Error
}

// Statistics operations
func (r *invitationRepository) GetInvitationStats(ctx context.Context, eventID int) (*dto.InvitationStatsResponse, error) {
	var stats dto.InvitationStatsResponse
	stats.EventID = eventID

	// Total invitations
	var totalInvitations int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("event_id = ?", eventID).
		Count(&totalInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.TotalInvitations = int(totalInvitations)

	// Pending invitations
	var pendingInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("event_id = ? AND status = ?", eventID, domain.InvitationStatusPending).
		Count(&pendingInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.PendingInvitations = int(pendingInvitations)

	// Approved invitations
	var approvedInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("event_id = ? AND status = ?", eventID, domain.InvitationStatusApproved).
		Count(&approvedInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.ApprovedInvitations = int(approvedInvitations)

	// Rejected invitations
	var rejectedInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("event_id = ? AND status = ?", eventID, domain.InvitationStatusRejected).
		Count(&rejectedInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.RejectedInvitations = int(rejectedInvitations)

	// System user invitations
	var systemUserInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("event_id = ? AND invited_user_id IS NOT NULL", eventID).
		Count(&systemUserInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.SystemUserInvitations = int(systemUserInvitations)

	// External user invitations
	var externalUserInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("event_id = ? AND invited_user_id IS NULL", eventID).
		Count(&externalUserInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.ExternalUserInvitations = int(externalUserInvitations)

	// Calculate rates
	if stats.TotalInvitations > 0 {
		respondedInvitations := stats.ApprovedInvitations + stats.RejectedInvitations
		stats.ResponseRate = (float64(respondedInvitations) / float64(stats.TotalInvitations)) * 100

		if respondedInvitations > 0 {
			stats.ApprovalRate = (float64(stats.ApprovedInvitations) / float64(respondedInvitations)) * 100
		}
	}

	return &stats, nil
}

func (r *invitationRepository) GetUserInvitationStats(ctx context.Context, userID int) (*dto.UserInvitationStatsResponse, error) {
	var stats dto.UserInvitationStatsResponse
	stats.UserID = userID

	// Total invitations
	var totalInvitations int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("invited_user_id = ?", userID).
		Count(&totalInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.TotalInvitations = int(totalInvitations)

	// Pending invitations
	var pendingInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("invited_user_id = ? AND status = ?", userID, domain.InvitationStatusPending).
		Count(&pendingInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.PendingInvitations = int(pendingInvitations)

	// Approved invitations
	var approvedInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("invited_user_id = ? AND status = ?", userID, domain.InvitationStatusApproved).
		Count(&approvedInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.ApprovedInvitations = int(approvedInvitations)

	// Rejected invitations
	var rejectedInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("invited_user_id = ? AND status = ?", userID, domain.InvitationStatusRejected).
		Count(&rejectedInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.RejectedInvitations = int(rejectedInvitations)

	// Calculate rates
	if stats.TotalInvitations > 0 {
		respondedInvitations := stats.ApprovedInvitations + stats.RejectedInvitations
		stats.ResponseRate = (float64(respondedInvitations) / float64(stats.TotalInvitations)) * 100

		if respondedInvitations > 0 {
			stats.ApprovalRate = (float64(stats.ApprovedInvitations) / float64(respondedInvitations)) * 100
		}
	}

	return &stats, nil
}

func (r *invitationRepository) GetSystemInvitationStats(ctx context.Context) (*dto.SystemInvitationStatsResponse, error) {
	var stats dto.SystemInvitationStatsResponse

	// Total invitations
	var totalInvitations int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Count(&totalInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.TotalInvitations = int(totalInvitations)

	// Pending invitations
	var pendingInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("status = ?", domain.InvitationStatusPending).
		Count(&pendingInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.PendingInvitations = int(pendingInvitations)

	// Approved invitations
	var approvedInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("status = ?", domain.InvitationStatusApproved).
		Count(&approvedInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.ApprovedInvitations = int(approvedInvitations)

	// Rejected invitations
	var rejectedInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("status = ?", domain.InvitationStatusRejected).
		Count(&rejectedInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.RejectedInvitations = int(rejectedInvitations)

	// System user invitations
	var systemUserInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("invited_user_id IS NOT NULL").
		Count(&systemUserInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.SystemUserInvitations = int(systemUserInvitations)

	// External user invitations
	var externalUserInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("invited_user_id IS NULL").
		Count(&externalUserInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.ExternalUserInvitations = int(externalUserInvitations)

	// Expired invitations
	expirationTime := time.Now().Add(-7 * 24 * time.Hour)
	var expiredInvitations int64
	err = r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("status = ? AND created_at < ?", domain.InvitationStatusPending, expirationTime).
		Count(&expiredInvitations).Error
	if err != nil {
		return nil, err
	}
	stats.ExpiredInvitations = int(expiredInvitations)

	// Calculate rates
	if stats.TotalInvitations > 0 {
		respondedInvitations := stats.ApprovedInvitations + stats.RejectedInvitations
		stats.ResponseRate = (float64(respondedInvitations) / float64(stats.TotalInvitations)) * 100

		if respondedInvitations > 0 {
			stats.ApprovalRate = (float64(stats.ApprovedInvitations) / float64(respondedInvitations)) * 100
		}
	}

	return &stats, nil
}

// Expiration operations
func (r *invitationRepository) GetExpiredInvitations(ctx context.Context, expirationHours int, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	var invitations []*domain.Invitation
	var total int64

	expirationTime := time.Now().Add(-time.Duration(expirationHours) * time.Hour)
	query := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("status = ? AND created_at < ?", domain.InvitationStatusPending, expirationTime)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Event").
		Preload("InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("created_at ASC").
		Find(&invitations).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return invitations, paginationResponse, nil
}

func (r *invitationRepository) DeleteExpiredInvitations(ctx context.Context, expirationHours int) (int64, error) {
	expirationTime := time.Now().Add(-time.Duration(expirationHours) * time.Hour)
	result := r.db.WithContext(ctx).
		Where("status = ? AND created_at < ?", domain.InvitationStatusPending, expirationTime).
		Delete(&domain.Invitation{})
	return result.RowsAffected, result.Error
}

func (r *invitationRepository) GetInvitationsExpiringIn(ctx context.Context, hours int, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	var invitations []*domain.Invitation
	var total int64

	expirationTime := time.Now().Add(time.Duration(hours) * time.Hour)
	query := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("status = ? AND created_at < ?", domain.InvitationStatusPending, expirationTime)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Event").
		Preload("InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("created_at ASC").
		Find(&invitations).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return invitations, paginationResponse, nil
}

// Advanced filtering
func (r *invitationRepository) GetInvitationsWithFilters(ctx context.Context, filters dto.InvitationFilterRequest, pagination dto.PaginationRequest) ([]*domain.Invitation, *dto.PaginationResponse, error) {
	var invitations []*domain.Invitation
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Invitation{})

	// Apply filters
	if filters.EventID != nil {
		query = query.Where("event_id = ?", *filters.EventID)
	}
	if filters.InvitedUserID != nil {
		query = query.Where("invited_user_id = ?", *filters.InvitedUserID)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.Email != nil && *filters.Email != "" {
		query = query.Where("invited_email ILIKE ?", "%"+*filters.Email+"%")
	}
	if filters.HasResponded != nil {
		if *filters.HasResponded {
			query = query.Where("responded_at IS NOT NULL")
		} else {
			query = query.Where("responded_at IS NULL")
		}
	}
	if filters.IsExpired != nil && *filters.IsExpired {
		expirationTime := time.Now().Add(-7 * 24 * time.Hour)
		query = query.Where("status = ? AND created_at < ?", domain.InvitationStatusPending, expirationTime)
	}
	if filters.Query != nil && *filters.Query != "" {
		query = query.Where("invited_email ILIKE ?", "%"+*filters.Query+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Preload("Event").
		Preload("InvitedUser").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return invitations, paginationResponse, nil
}

// Preloading operations
func (r *invitationRepository) PreloadEvent(ctx context.Context, invitations []*domain.Invitation) error {
	for i := range invitations {
		err := r.db.WithContext(ctx).Preload("Event").First(invitations[i], invitations[i].ID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *invitationRepository) PreloadInvitedUser(ctx context.Context, invitations []*domain.Invitation) error {
	for i := range invitations {
		err := r.db.WithContext(ctx).Preload("InvitedUser").First(invitations[i], invitations[i].ID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *invitationRepository) PreloadAllRelations(ctx context.Context, invitations []*domain.Invitation) error {
	for i := range invitations {
		err := r.db.WithContext(ctx).
			Preload("Event").
			Preload("Event.Creator").
			Preload("InvitedUser").
			First(invitations[i], invitations[i].ID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// Email-based operations for external users
func (r *invitationRepository) GetPendingInvitationsByEmail(ctx context.Context, email string) ([]*domain.Invitation, error) {
	var invitations []*domain.Invitation
	err := r.db.WithContext(ctx).
		Preload("Event").
		Preload("Event.Creator").
		Where("invited_email = ? AND status = ?", email, domain.InvitationStatusPending).
		Find(&invitations).Error
	return invitations, err
}

func (r *invitationRepository) GetEventInvitationByEmail(ctx context.Context, eventID int, email string) (*domain.Invitation, error) {
	return r.GetByEventAndEmail(ctx, eventID, email)
}

func (r *invitationRepository) UpdateInvitedUserByEmail(ctx context.Context, email string, userID int) error {
	return r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("invited_email = ? AND invited_user_id IS NULL", email).
		Update("invited_user_id", userID).Error
}
