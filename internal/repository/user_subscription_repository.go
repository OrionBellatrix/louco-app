package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
)

type UserSubscriptionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, subscription *domain.UserSubscription) error
	GetByID(ctx context.Context, id uint) (*domain.UserSubscription, error)
	GetByUserID(ctx context.Context, userID uint) ([]*domain.UserSubscription, error)
	Update(ctx context.Context, subscription *domain.UserSubscription) error
	Delete(ctx context.Context, id uint) error

	// Active subscription/package queries
	GetActiveSubscriptionByUserID(ctx context.Context, userID uint) (*domain.UserSubscription, error)
	GetActivePackagesByUserID(ctx context.Context, userID uint) ([]*domain.UserSubscription, error)

	// Usage tracking queries
	GetCurrentWeeklyUsage(ctx context.Context, userID uint) (int, error)
	GetCurrentMonthlyUsage(ctx context.Context, userID uint) (int, error)
	GetTotalCreditsRemaining(ctx context.Context, userID uint) (int, error)

	// Limit checking queries
	CanPublishEvent(ctx context.Context, userID uint) (bool, error)
	GetPublishingRights(ctx context.Context, userID uint) (*domain.PublishingRights, error)

	// Usage consumption
	ConsumeEventCredit(ctx context.Context, userID uint) error
	IncrementWeeklyUsage(ctx context.Context, userID uint) error
	IncrementMonthlyUsage(ctx context.Context, userID uint) error

	// Limit reset operations
	ResetWeeklyLimits(ctx context.Context) error
	ResetMonthlyLimits(ctx context.Context) error

	// Expiration management
	GetExpiredSubscriptions(ctx context.Context) ([]*domain.UserSubscription, error)
	GetExpiredPackages(ctx context.Context) ([]*domain.UserSubscription, error)
	MarkAsExpired(ctx context.Context, id uint) error

	// Payment integration
	GetByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*domain.UserSubscription, error)
	GetByStripePaymentIntentID(ctx context.Context, stripePaymentIntentID string) (*domain.UserSubscription, error)
	UpdatePaymentStatus(ctx context.Context, id uint, status domain.PaymentStatus) error

	// Statistics and reporting
	GetSubscriptionStats(ctx context.Context, userID uint) (*domain.SubscriptionStats, error)
	GetUserSubscriptionHistory(ctx context.Context, userID uint, limit, offset int) ([]*domain.UserSubscription, int64, error)
}

type SubscriptionPlanRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, plan *domain.SubscriptionPlan) error
	GetByID(ctx context.Context, id uint) (*domain.SubscriptionPlan, error)
	GetBySlug(ctx context.Context, slug string) (*domain.SubscriptionPlan, error)
	Update(ctx context.Context, plan *domain.SubscriptionPlan) error
	Delete(ctx context.Context, id uint) error

	// Plan listing and filtering
	GetAll(ctx context.Context) ([]*domain.SubscriptionPlan, error)
	GetByType(ctx context.Context, planType domain.PlanType) ([]*domain.SubscriptionPlan, error)
	GetActiveSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error)
	GetActivePackagePlans(ctx context.Context) ([]*domain.SubscriptionPlan, error)

	// Plan availability
	GetAvailablePlans(ctx context.Context) ([]*domain.SubscriptionPlan, error)
	IsActive(ctx context.Context, id uint) (bool, error)

	// Stripe integration
	GetByStripeProductID(ctx context.Context, stripeProductID string) (*domain.SubscriptionPlan, error)
	GetByStripePriceID(ctx context.Context, stripePriceID string) (*domain.SubscriptionPlan, error)
}
