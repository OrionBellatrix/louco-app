package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"github.com/louco-event/pkg/logger"
	"gorm.io/gorm"
)

type userSubscriptionRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewUserSubscriptionRepository(db *gorm.DB, logger *logger.Logger) repository.UserSubscriptionRepository {
	return &userSubscriptionRepository{
		db:     db,
		logger: logger,
	}
}

// Basic CRUD operations
func (r *userSubscriptionRepository) Create(ctx context.Context, subscription *domain.UserSubscription) error {
	if err := r.db.WithContext(ctx).Create(subscription).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to create user subscription")
		return fmt.Errorf("failed to create user subscription: %w", err)
	}
	return nil
}

func (r *userSubscriptionRepository) GetByID(ctx context.Context, id uint) (*domain.UserSubscription, error) {
	var subscription domain.UserSubscription
	if err := r.db.WithContext(ctx).Preload("User").First(&subscription, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Failed to get user subscription by ID")
		return nil, fmt.Errorf("failed to get user subscription: %w", err)
	}
	return &subscription, nil
}

func (r *userSubscriptionRepository) GetByUserID(ctx context.Context, userID uint) ([]*domain.UserSubscription, error) {
	var subscriptions []*domain.UserSubscription
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&subscriptions).Error; err != nil {
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get user subscriptions by user ID")
		return nil, fmt.Errorf("failed to get user subscriptions: %w", err)
	}
	return subscriptions, nil
}

func (r *userSubscriptionRepository) Update(ctx context.Context, subscription *domain.UserSubscription) error {
	if err := r.db.WithContext(ctx).Save(subscription).Error; err != nil {
		r.logger.Error().Err(err).Int("id", subscription.ID).Msg("Failed to update user subscription")
		return fmt.Errorf("failed to update user subscription: %w", err)
	}
	return nil
}

func (r *userSubscriptionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.UserSubscription{}, id).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Failed to delete user subscription")
		return fmt.Errorf("failed to delete user subscription: %w", err)
	}
	return nil
}

// Active subscription/package queries
func (r *userSubscriptionRepository) GetActiveSubscriptionByUserID(ctx context.Context, userID uint) (*domain.UserSubscription, error) {
	var subscription domain.UserSubscription
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND status = ? AND (expired_at IS NULL OR expired_at > ?)",
			userID, domain.SubscriptionTypeSubscription, domain.SubscriptionStatusActive, time.Now()).
		Order("created_at DESC").
		First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get active subscription")
		return nil, fmt.Errorf("failed to get active subscription: %w", err)
	}
	return &subscription, nil
}

func (r *userSubscriptionRepository) GetActivePackagesByUserID(ctx context.Context, userID uint) ([]*domain.UserSubscription, error) {
	var packages []*domain.UserSubscription
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND status = ? AND (expired_at IS NULL OR expired_at > ?) AND used_credits < total_credits",
						userID, domain.SubscriptionTypePackage, domain.SubscriptionStatusActive, time.Now()).
		Order("created_at ASC"). // Use oldest packages first
		Find(&packages).Error; err != nil {
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get active packages")
		return nil, fmt.Errorf("failed to get active packages: %w", err)
	}
	return packages, nil
}

// Usage tracking queries
func (r *userSubscriptionRepository) GetCurrentWeeklyUsage(ctx context.Context, userID uint) (int, error) {
	var totalUsage int64
	if err := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, domain.SubscriptionTypeSubscription, domain.SubscriptionStatusActive).
		Select("COALESCE(SUM(weekly_used), 0)").
		Scan(&totalUsage).Error; err != nil {
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get current weekly usage")
		return 0, fmt.Errorf("failed to get weekly usage: %w", err)
	}
	return int(totalUsage), nil
}

func (r *userSubscriptionRepository) GetCurrentMonthlyUsage(ctx context.Context, userID uint) (int, error) {
	var totalUsage int64
	if err := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, domain.SubscriptionTypeSubscription, domain.SubscriptionStatusActive).
		Select("COALESCE(SUM(monthly_used), 0)").
		Scan(&totalUsage).Error; err != nil {
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get current monthly usage")
		return 0, fmt.Errorf("failed to get monthly usage: %w", err)
	}
	return int(totalUsage), nil
}

func (r *userSubscriptionRepository) GetTotalCreditsRemaining(ctx context.Context, userID uint) (int, error) {
	var totalRemaining int64
	if err := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND type = ? AND status = ? AND (expired_at IS NULL OR expired_at > ?)",
			userID, domain.SubscriptionTypePackage, domain.SubscriptionStatusActive, time.Now()).
		Select("COALESCE(SUM(total_credits - used_credits), 0)").
		Scan(&totalRemaining).Error; err != nil {
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get total credits remaining")
		return 0, fmt.Errorf("failed to get credits remaining: %w", err)
	}
	return int(totalRemaining), nil
}

// Limit checking queries
func (r *userSubscriptionRepository) CanPublishEvent(ctx context.Context, userID uint) (bool, error) {
	// Check active subscription limits
	subscription, err := r.GetActiveSubscriptionByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	if subscription != nil && subscription.CanPublishEvent() {
		return true, nil
	}

	// Check active packages
	packages, err := r.GetActivePackagesByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, pkg := range packages {
		if pkg.CanPublishEvent() {
			return true, nil
		}
	}

	return false, nil
}

func (r *userSubscriptionRepository) GetPublishingRights(ctx context.Context, userID uint) (*domain.PublishingRights, error) {
	rights := &domain.PublishingRights{
		CanPublish: false,
	}

	// Get active subscription
	subscription, err := r.GetActiveSubscriptionByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if subscription != nil {
		rights.ActiveSubscription = subscription
		rights.WeeklyLimit = subscription.RemainingWeeklyLimit()
		rights.WeeklyUsed = subscription.WeeklyUsed
		rights.WeeklyRemaining = subscription.RemainingWeeklyLimit()
		rights.MonthlyLimit = subscription.RemainingMonthlyLimit()
		rights.MonthlyUsed = subscription.MonthlyUsed
		rights.MonthlyRemaining = subscription.RemainingMonthlyLimit()

		if subscription.CanPublishEvent() {
			rights.CanPublish = true
		}
	}

	// Get active packages
	packages, err := r.GetActivePackagesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(packages) > 0 {
		rights.ActivePackages = packages
		totalCredits := 0
		usedCredits := 0

		for _, pkg := range packages {
			if pkg.TotalCredits != nil {
				totalCredits += *pkg.TotalCredits
			}
			usedCredits += pkg.UsedCredits

			if pkg.CanPublishEvent() {
				rights.CanPublish = true
			}
		}

		rights.TotalCredits = totalCredits
		rights.UsedCredits = usedCredits
		rights.RemainingCredits = totalCredits - usedCredits
	}

	// Set restriction reason if cannot publish
	if !rights.CanPublish {
		if subscription == nil && len(packages) == 0 {
			rights.RestrictionReason = "No active subscription or package"
		} else if subscription != nil && (subscription.HasWeeklyLimit() || subscription.HasMonthlyLimit()) {
			rights.RestrictionReason = "Subscription limit reached"
		} else if len(packages) > 0 && rights.RemainingCredits <= 0 {
			rights.RestrictionReason = "No credits remaining"
		} else {
			rights.RestrictionReason = "Subscription expired"
		}
	}

	return rights, nil
}

// Usage consumption
func (r *userSubscriptionRepository) ConsumeEventCredit(ctx context.Context, userID uint) error {
	// Try to consume from active subscription first
	subscription, err := r.GetActiveSubscriptionByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if subscription != nil && subscription.CanPublishEvent() {
		subscription.ConsumeUsage()
		return r.Update(ctx, subscription)
	}

	// Try to consume from active packages
	packages, err := r.GetActivePackagesByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		if pkg.CanPublishEvent() {
			pkg.ConsumeUsage()
			return r.Update(ctx, pkg)
		}
	}

	return fmt.Errorf("no available credits to consume")
}

func (r *userSubscriptionRepository) IncrementWeeklyUsage(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, domain.SubscriptionTypeSubscription, domain.SubscriptionStatusActive).
		Update("weekly_used", gorm.Expr("weekly_used + 1")).Error
}

func (r *userSubscriptionRepository) IncrementMonthlyUsage(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, domain.SubscriptionTypeSubscription, domain.SubscriptionStatusActive).
		Update("monthly_used", gorm.Expr("monthly_used + 1")).Error
}

// Limit reset operations
func (r *userSubscriptionRepository) ResetWeeklyLimits(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("type = ? AND status = ?", domain.SubscriptionTypeSubscription, domain.SubscriptionStatusActive).
		Update("weekly_used", 0).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to reset weekly limits")
		return fmt.Errorf("failed to reset weekly limits: %w", err)
	}
	return nil
}

func (r *userSubscriptionRepository) ResetMonthlyLimits(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("type = ? AND status = ?", domain.SubscriptionTypeSubscription, domain.SubscriptionStatusActive).
		Update("monthly_used", 0).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to reset monthly limits")
		return fmt.Errorf("failed to reset monthly limits: %w", err)
	}
	return nil
}

// Expiration management
func (r *userSubscriptionRepository) GetExpiredSubscriptions(ctx context.Context) ([]*domain.UserSubscription, error) {
	var subscriptions []*domain.UserSubscription
	if err := r.db.WithContext(ctx).
		Where("status = ? AND expired_at IS NOT NULL AND expired_at <= ?",
			domain.SubscriptionStatusActive, time.Now()).
		Find(&subscriptions).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to get expired subscriptions")
		return nil, fmt.Errorf("failed to get expired subscriptions: %w", err)
	}
	return subscriptions, nil
}

func (r *userSubscriptionRepository) GetExpiredPackages(ctx context.Context) ([]*domain.UserSubscription, error) {
	var packages []*domain.UserSubscription
	if err := r.db.WithContext(ctx).
		Where("type = ? AND status = ? AND expired_at IS NOT NULL AND expired_at <= ?",
			domain.SubscriptionTypePackage, domain.SubscriptionStatusActive, time.Now()).
		Find(&packages).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to get expired packages")
		return nil, fmt.Errorf("failed to get expired packages: %w", err)
	}
	return packages, nil
}

func (r *userSubscriptionRepository) MarkAsExpired(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("id = ?", id).
		Update("status", domain.SubscriptionStatusExpired).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Failed to mark subscription as expired")
		return fmt.Errorf("failed to mark as expired: %w", err)
	}
	return nil
}

// Payment integration
func (r *userSubscriptionRepository) GetByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*domain.UserSubscription, error) {
	var subscription domain.UserSubscription
	if err := r.db.WithContext(ctx).
		Where("stripe_id = ?", stripeSubscriptionID).
		First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("stripe_id", stripeSubscriptionID).Msg("Failed to get subscription by Stripe ID")
		return nil, fmt.Errorf("failed to get subscription by Stripe ID: %w", err)
	}
	return &subscription, nil
}

func (r *userSubscriptionRepository) GetByStripePaymentIntentID(ctx context.Context, stripePaymentIntentID string) (*domain.UserSubscription, error) {
	var subscription domain.UserSubscription
	if err := r.db.WithContext(ctx).
		Where("stripe_id = ?", stripePaymentIntentID).
		First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("stripe_payment_intent_id", stripePaymentIntentID).Msg("Failed to get subscription by Stripe Payment Intent ID")
		return nil, fmt.Errorf("failed to get subscription by Stripe Payment Intent ID: %w", err)
	}
	return &subscription, nil
}

func (r *userSubscriptionRepository) UpdatePaymentStatus(ctx context.Context, id uint, status domain.PaymentStatus) error {
	var subscriptionStatus domain.SubscriptionStatus

	switch status {
	case domain.PaymentStatusSucceeded:
		subscriptionStatus = domain.SubscriptionStatusActive
	case domain.PaymentStatusFailed, domain.PaymentStatusCancelled:
		subscriptionStatus = domain.SubscriptionStatusCancelled
	default:
		subscriptionStatus = domain.SubscriptionStatusPending
	}

	if err := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("id = ?", id).
		Update("status", subscriptionStatus).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Str("status", string(status)).Msg("Failed to update payment status")
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	return nil
}

// Statistics and reporting
func (r *userSubscriptionRepository) GetSubscriptionStats(ctx context.Context, userID uint) (*domain.SubscriptionStats, error) {
	stats := &domain.SubscriptionStats{}

	// Get subscription counts
	var subscriptionCount, activeSubscriptionCount int64
	r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND type = ?", userID, domain.SubscriptionTypeSubscription).
		Count(&subscriptionCount)
	r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, domain.SubscriptionTypeSubscription, domain.SubscriptionStatusActive).
		Count(&activeSubscriptionCount)

	// Get package counts
	var packageCount, activePackageCount int64
	r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND type = ?", userID, domain.SubscriptionTypePackage).
		Count(&packageCount)
	r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, domain.SubscriptionTypePackage, domain.SubscriptionStatusActive).
		Count(&activePackageCount)

	// Get total spent
	var totalSpent float64
	r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND status IN (?)", userID, []domain.SubscriptionStatus{domain.SubscriptionStatusActive, domain.SubscriptionStatusExpired}).
		Select("COALESCE(SUM(price), 0)").
		Scan(&totalSpent)

	// Get events remaining
	creditsRemaining, _ := r.GetTotalCreditsRemaining(ctx, userID)

	// Get current period usage
	weeklyUsage, _ := r.GetCurrentWeeklyUsage(ctx, userID)
	monthlyUsage, _ := r.GetCurrentMonthlyUsage(ctx, userID)
	currentUsage := weeklyUsage
	if monthlyUsage > weeklyUsage {
		currentUsage = monthlyUsage
	}

	// Get last payment and next billing dates
	var lastPayment, nextBilling *time.Time
	var lastPaymentTime, nextBillingTime time.Time

	if err := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND status = ?", userID, domain.SubscriptionStatusActive).
		Select("MAX(created_at)").
		Scan(&lastPaymentTime).Error; err == nil && !lastPaymentTime.IsZero() {
		lastPayment = &lastPaymentTime
	}

	if err := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, domain.SubscriptionTypeSubscription, domain.SubscriptionStatusActive).
		Select("MIN(expired_at)").
		Scan(&nextBillingTime).Error; err == nil && !nextBillingTime.IsZero() {
		nextBilling = &nextBillingTime
	}

	stats.TotalSubscriptions = int(subscriptionCount)
	stats.ActiveSubscriptions = int(activeSubscriptionCount)
	stats.TotalPackages = int(packageCount)
	stats.ActivePackages = int(activePackageCount)
	stats.TotalSpent = totalSpent
	stats.EventsRemaining = creditsRemaining
	stats.CurrentPeriodUsage = currentUsage
	stats.LastPaymentDate = lastPayment
	stats.NextBillingDate = nextBilling

	return stats, nil
}

func (r *userSubscriptionRepository) GetUserSubscriptionHistory(ctx context.Context, userID uint, limit, offset int) ([]*domain.UserSubscription, int64, error) {
	var subscriptions []*domain.UserSubscription
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to count user subscription history")
		return nil, 0, fmt.Errorf("failed to count subscription history: %w", err)
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&subscriptions).Error; err != nil {
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get user subscription history")
		return nil, 0, fmt.Errorf("failed to get subscription history: %w", err)
	}

	return subscriptions, total, nil
}
