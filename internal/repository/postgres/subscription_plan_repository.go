package postgres

import (
	"context"
	"fmt"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"github.com/louco-event/pkg/logger"
	"gorm.io/gorm"
)

type subscriptionPlanRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewSubscriptionPlanRepository(db *gorm.DB, logger *logger.Logger) repository.SubscriptionPlanRepository {
	return &subscriptionPlanRepository{
		db:     db,
		logger: logger,
	}
}

// Basic CRUD operations
func (r *subscriptionPlanRepository) Create(ctx context.Context, plan *domain.SubscriptionPlan) error {
	if err := r.db.WithContext(ctx).Create(plan).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to create subscription plan")
		return fmt.Errorf("failed to create subscription plan: %w", err)
	}
	return nil
}

func (r *subscriptionPlanRepository) GetByID(ctx context.Context, id uint) (*domain.SubscriptionPlan, error) {
	var plan domain.SubscriptionPlan
	if err := r.db.WithContext(ctx).First(&plan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Failed to get subscription plan by ID")
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}
	return &plan, nil
}

func (r *subscriptionPlanRepository) GetBySlug(ctx context.Context, slug string) (*domain.SubscriptionPlan, error) {
	var plan domain.SubscriptionPlan
	if err := r.db.WithContext(ctx).Where("name = ?", slug).First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("slug", slug).Msg("Failed to get subscription plan by slug")
		return nil, fmt.Errorf("failed to get subscription plan by slug: %w", err)
	}
	return &plan, nil
}

func (r *subscriptionPlanRepository) Update(ctx context.Context, plan *domain.SubscriptionPlan) error {
	if err := r.db.WithContext(ctx).Save(plan).Error; err != nil {
		r.logger.Error().Err(err).Int("id", plan.ID).Msg("Failed to update subscription plan")
		return fmt.Errorf("failed to update subscription plan: %w", err)
	}
	return nil
}

func (r *subscriptionPlanRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.SubscriptionPlan{}, id).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Failed to delete subscription plan")
		return fmt.Errorf("failed to delete subscription plan: %w", err)
	}
	return nil
}

// Plan listing and filtering
func (r *subscriptionPlanRepository) GetAll(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	var plans []*domain.SubscriptionPlan
	if err := r.db.WithContext(ctx).Order("sort_order ASC, created_at ASC").Find(&plans).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to get all subscription plans")
		return nil, fmt.Errorf("failed to get all subscription plans: %w", err)
	}
	return plans, nil
}

func (r *subscriptionPlanRepository) GetByType(ctx context.Context, planType domain.PlanType) ([]*domain.SubscriptionPlan, error) {
	var plans []*domain.SubscriptionPlan
	if err := r.db.WithContext(ctx).
		Where("type = ?", planType).
		Order("sort_order ASC, created_at ASC").
		Find(&plans).Error; err != nil {
		r.logger.Error().Err(err).Str("type", string(planType)).Msg("Failed to get subscription plans by type")
		return nil, fmt.Errorf("failed to get subscription plans by type: %w", err)
	}
	return plans, nil
}

func (r *subscriptionPlanRepository) GetActiveSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	var plans []*domain.SubscriptionPlan
	if err := r.db.WithContext(ctx).
		Where("type = ? AND is_active = ?", domain.SubscriptionTypeSubscription, true).
		Order("sort_order ASC, created_at ASC").
		Find(&plans).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to get active subscription plans")
		return nil, fmt.Errorf("failed to get active subscription plans: %w", err)
	}
	return plans, nil
}

func (r *subscriptionPlanRepository) GetActivePackagePlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	var plans []*domain.SubscriptionPlan
	if err := r.db.WithContext(ctx).
		Where("type = ? AND is_active = ?", domain.SubscriptionTypePackage, true).
		Order("sort_order ASC, created_at ASC").
		Find(&plans).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to get active package plans")
		return nil, fmt.Errorf("failed to get active package plans: %w", err)
	}
	return plans, nil
}

// Plan availability
func (r *subscriptionPlanRepository) GetAvailablePlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	var plans []*domain.SubscriptionPlan
	if err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("type ASC, sort_order ASC, created_at ASC").
		Find(&plans).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to get available plans")
		return nil, fmt.Errorf("failed to get available plans: %w", err)
	}
	return plans, nil
}

func (r *subscriptionPlanRepository) IsActive(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.SubscriptionPlan{}).
		Where("id = ? AND is_active = ?", id, true).
		Count(&count).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Failed to check if plan is active")
		return false, fmt.Errorf("failed to check if plan is active: %w", err)
	}
	return count > 0, nil
}

// Stripe integration
func (r *subscriptionPlanRepository) GetByStripeProductID(ctx context.Context, stripeProductID string) (*domain.SubscriptionPlan, error) {
	var plan domain.SubscriptionPlan
	if err := r.db.WithContext(ctx).
		Where("stripe_id = ?", stripeProductID).
		First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("stripe_product_id", stripeProductID).Msg("Failed to get plan by Stripe product ID")
		return nil, fmt.Errorf("failed to get plan by Stripe product ID: %w", err)
	}
	return &plan, nil
}

func (r *subscriptionPlanRepository) GetByStripePriceID(ctx context.Context, stripePriceID string) (*domain.SubscriptionPlan, error) {
	var plan domain.SubscriptionPlan
	if err := r.db.WithContext(ctx).
		Where("stripe_id = ?", stripePriceID).
		First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("stripe_price_id", stripePriceID).Msg("Failed to get plan by Stripe price ID")
		return nil, fmt.Errorf("failed to get plan by Stripe price ID: %w", err)
	}
	return &plan, nil
}
