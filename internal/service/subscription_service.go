package service

import (
	"context"
	"fmt"
	"time"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"github.com/louco-event/pkg/logger"
)

type SubscriptionService interface {
	// Plan management
	GetAvailablePlans(ctx context.Context) ([]*domain.SubscriptionPlan, error)
	GetSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error)
	GetPackagePlans(ctx context.Context) ([]*domain.SubscriptionPlan, error)
	GetPlanByID(ctx context.Context, planID uint) (*domain.SubscriptionPlan, error)
	GetPlanBySlug(ctx context.Context, slug string) (*domain.SubscriptionPlan, error)

	// User subscription management
	GetUserSubscriptions(ctx context.Context, userID uint) ([]*domain.UserSubscription, error)
	GetActiveSubscription(ctx context.Context, userID uint) (*domain.UserSubscription, error)
	GetActivePackages(ctx context.Context, userID uint) ([]*domain.UserSubscription, error)
	GetSubscriptionHistory(ctx context.Context, userID uint, limit, offset int) ([]*domain.UserSubscription, int64, error)

	// Publishing rights and validation
	GetPublishingRights(ctx context.Context, userID uint) (*domain.PublishingRights, error)
	CanPublishEvent(ctx context.Context, userID uint) (bool, error)
	ValidateEventPublishing(ctx context.Context, userID uint) error

	// Usage tracking and consumption
	ConsumeEventCredit(ctx context.Context, userID uint) error
	GetUsageStats(ctx context.Context, userID uint) (*domain.SubscriptionStats, error)

	// Subscription lifecycle
	CreateSubscription(ctx context.Context, userID uint, planID uint, stripeSubscriptionID string) (*domain.UserSubscription, error)
	CreatePackage(ctx context.Context, userID uint, planID uint, stripePaymentIntentID string) (*domain.UserSubscription, error)
	ActivateSubscription(ctx context.Context, subscriptionID uint) error
	CancelSubscription(ctx context.Context, subscriptionID uint) error
	ExpireSubscription(ctx context.Context, subscriptionID uint) error

	// Payment integration
	HandlePaymentSuccess(ctx context.Context, stripeID string) error
	HandlePaymentFailed(ctx context.Context, stripeID string) error
	HandlePaymentRefunded(ctx context.Context, stripeID string) error
	ActivateSubscriptionByStripeID(ctx context.Context, stripeID string) error
	CancelSubscriptionByStripeID(ctx context.Context, stripeID string) error

	// Administrative functions
	ResetWeeklyLimits(ctx context.Context) error
	ResetMonthlyLimits(ctx context.Context) error
	ProcessExpiredSubscriptions(ctx context.Context) error

	// Plan seeding (for initial setup)
	SeedDefaultPlans(ctx context.Context) error
}

type subscriptionService struct {
	userSubscriptionRepo repository.UserSubscriptionRepository
	subscriptionPlanRepo repository.SubscriptionPlanRepository
	logger               *logger.Logger
}

func NewSubscriptionService(
	userSubscriptionRepo repository.UserSubscriptionRepository,
	subscriptionPlanRepo repository.SubscriptionPlanRepository,
	logger *logger.Logger,
) SubscriptionService {
	return &subscriptionService{
		userSubscriptionRepo: userSubscriptionRepo,
		subscriptionPlanRepo: subscriptionPlanRepo,
		logger:               logger,
	}
}

// Plan management
func (s *subscriptionService) GetAvailablePlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	return s.subscriptionPlanRepo.GetAvailablePlans(ctx)
}

func (s *subscriptionService) GetSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	return s.subscriptionPlanRepo.GetActiveSubscriptionPlans(ctx)
}

func (s *subscriptionService) GetPackagePlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	return s.subscriptionPlanRepo.GetActivePackagePlans(ctx)
}

func (s *subscriptionService) GetPlanByID(ctx context.Context, planID uint) (*domain.SubscriptionPlan, error) {
	plan, err := s.subscriptionPlanRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, fmt.Errorf("subscription plan not found")
	}
	return plan, nil
}

func (s *subscriptionService) GetPlanBySlug(ctx context.Context, slug string) (*domain.SubscriptionPlan, error) {
	plan, err := s.subscriptionPlanRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, fmt.Errorf("subscription plan not found")
	}
	return plan, nil
}

// User subscription management
func (s *subscriptionService) GetUserSubscriptions(ctx context.Context, userID uint) ([]*domain.UserSubscription, error) {
	return s.userSubscriptionRepo.GetByUserID(ctx, userID)
}

func (s *subscriptionService) GetActiveSubscription(ctx context.Context, userID uint) (*domain.UserSubscription, error) {
	return s.userSubscriptionRepo.GetActiveSubscriptionByUserID(ctx, userID)
}

func (s *subscriptionService) GetActivePackages(ctx context.Context, userID uint) ([]*domain.UserSubscription, error) {
	return s.userSubscriptionRepo.GetActivePackagesByUserID(ctx, userID)
}

func (s *subscriptionService) GetSubscriptionHistory(ctx context.Context, userID uint, limit, offset int) ([]*domain.UserSubscription, int64, error) {
	return s.userSubscriptionRepo.GetUserSubscriptionHistory(ctx, userID, limit, offset)
}

// Publishing rights and validation
func (s *subscriptionService) GetPublishingRights(ctx context.Context, userID uint) (*domain.PublishingRights, error) {
	return s.userSubscriptionRepo.GetPublishingRights(ctx, userID)
}

func (s *subscriptionService) CanPublishEvent(ctx context.Context, userID uint) (bool, error) {
	return s.userSubscriptionRepo.CanPublishEvent(ctx, userID)
}

func (s *subscriptionService) ValidateEventPublishing(ctx context.Context, userID uint) error {
	canPublish, err := s.CanPublishEvent(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check publishing rights: %w", err)
	}

	if !canPublish {
		rights, err := s.GetPublishingRights(ctx, userID)
		if err != nil {
			return fmt.Errorf("no active subscription or package found")
		}

		if rights.RestrictionReason != "" {
			return fmt.Errorf("cannot publish event: %s", rights.RestrictionReason)
		}

		return fmt.Errorf("cannot publish event: insufficient credits or limits reached")
	}

	return nil
}

// Usage tracking and consumption
func (s *subscriptionService) ConsumeEventCredit(ctx context.Context, userID uint) error {
	// First validate that user can publish
	if err := s.ValidateEventPublishing(ctx, userID); err != nil {
		return err
	}

	// Consume the credit
	if err := s.userSubscriptionRepo.ConsumeEventCredit(ctx, userID); err != nil {
		s.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to consume event credit")
		return fmt.Errorf("failed to consume event credit: %w", err)
	}

	s.logger.Info().Uint("user_id", userID).Msg("Event credit consumed successfully")
	return nil
}

func (s *subscriptionService) GetUsageStats(ctx context.Context, userID uint) (*domain.SubscriptionStats, error) {
	return s.userSubscriptionRepo.GetSubscriptionStats(ctx, userID)
}

// Subscription lifecycle
func (s *subscriptionService) CreateSubscription(ctx context.Context, userID uint, planID uint, stripeSubscriptionID string) (*domain.UserSubscription, error) {
	// Get the plan
	plan, err := s.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	if !plan.IsSubscription() {
		return nil, fmt.Errorf("plan is not a subscription type")
	}

	// Create user subscription from plan
	subscription := plan.CreateUserSubscription(int(userID))
	subscription.StripeID = &stripeSubscriptionID

	// Validate the subscription
	if err := subscription.Validate(); err != nil {
		return nil, fmt.Errorf("subscription validation failed: %w", err)
	}

	// Save to database
	if err := s.userSubscriptionRepo.Create(ctx, subscription); err != nil {
		s.logger.Error().Err(err).Uint("user_id", userID).Uint("plan_id", planID).Msg("Failed to create subscription")
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	s.logger.Info().Uint("user_id", userID).Uint("plan_id", planID).Str("stripe_id", stripeSubscriptionID).Msg("Subscription created successfully")
	return subscription, nil
}

func (s *subscriptionService) CreatePackage(ctx context.Context, userID uint, planID uint, stripePaymentIntentID string) (*domain.UserSubscription, error) {
	// Get the plan
	plan, err := s.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	if !plan.IsPackage() {
		return nil, fmt.Errorf("plan is not a package type")
	}

	// Create user subscription from plan
	subscription := plan.CreateUserSubscription(int(userID))
	subscription.StripeID = &stripePaymentIntentID

	// Validate the subscription
	if err := subscription.Validate(); err != nil {
		return nil, fmt.Errorf("package validation failed: %w", err)
	}

	// Save to database
	if err := s.userSubscriptionRepo.Create(ctx, subscription); err != nil {
		s.logger.Error().Err(err).Uint("user_id", userID).Uint("plan_id", planID).Msg("Failed to create package")
		return nil, fmt.Errorf("failed to create package: %w", err)
	}

	s.logger.Info().Uint("user_id", userID).Uint("plan_id", planID).Str("stripe_id", stripePaymentIntentID).Msg("Package created successfully")
	return subscription, nil
}

func (s *subscriptionService) ActivateSubscription(ctx context.Context, subscriptionID uint) error {
	subscription, err := s.userSubscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return err
	}
	if subscription == nil {
		return fmt.Errorf("subscription not found")
	}

	subscription.Status = domain.SubscriptionStatusActive
	now := time.Now()
	subscription.StartedAt = &now

	if err := s.userSubscriptionRepo.Update(ctx, subscription); err != nil {
		s.logger.Error().Err(err).Uint("subscription_id", subscriptionID).Msg("Failed to activate subscription")
		return fmt.Errorf("failed to activate subscription: %w", err)
	}

	s.logger.Info().Uint("subscription_id", subscriptionID).Msg("Subscription activated successfully")
	return nil
}

func (s *subscriptionService) CancelSubscription(ctx context.Context, subscriptionID uint) error {
	subscription, err := s.userSubscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return err
	}
	if subscription == nil {
		return fmt.Errorf("subscription not found")
	}

	subscription.Status = domain.SubscriptionStatusCancelled

	if err := s.userSubscriptionRepo.Update(ctx, subscription); err != nil {
		s.logger.Error().Err(err).Uint("subscription_id", subscriptionID).Msg("Failed to cancel subscription")
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	s.logger.Info().Uint("subscription_id", subscriptionID).Msg("Subscription cancelled successfully")
	return nil
}

func (s *subscriptionService) ExpireSubscription(ctx context.Context, subscriptionID uint) error {
	return s.userSubscriptionRepo.MarkAsExpired(ctx, subscriptionID)
}

// Payment integration
func (s *subscriptionService) HandlePaymentSuccess(ctx context.Context, stripeID string) error {
	// Try to find by subscription ID first
	subscription, err := s.userSubscriptionRepo.GetByStripeSubscriptionID(ctx, stripeID)
	if err != nil {
		return err
	}

	// If not found, try by payment intent ID
	if subscription == nil {
		subscription, err = s.userSubscriptionRepo.GetByStripePaymentIntentID(ctx, stripeID)
		if err != nil {
			return err
		}
	}

	if subscription == nil {
		return fmt.Errorf("subscription not found for Stripe ID: %s", stripeID)
	}

	// Update payment status to succeeded (which will activate the subscription)
	if err := s.userSubscriptionRepo.UpdatePaymentStatus(ctx, uint(subscription.ID), domain.PaymentStatusSucceeded); err != nil {
		s.logger.Error().Err(err).Str("stripe_id", stripeID).Msg("Failed to handle payment success")
		return fmt.Errorf("failed to handle payment success: %w", err)
	}

	s.logger.Info().Str("stripe_id", stripeID).Msg("Payment success handled")
	return nil
}

func (s *subscriptionService) HandlePaymentFailed(ctx context.Context, stripeID string) error {
	// Try to find by subscription ID first
	subscription, err := s.userSubscriptionRepo.GetByStripeSubscriptionID(ctx, stripeID)
	if err != nil {
		return err
	}

	// If not found, try by payment intent ID
	if subscription == nil {
		subscription, err = s.userSubscriptionRepo.GetByStripePaymentIntentID(ctx, stripeID)
		if err != nil {
			return err
		}
	}

	if subscription == nil {
		return fmt.Errorf("subscription not found for Stripe ID: %s", stripeID)
	}

	// Update payment status to failed
	if err := s.userSubscriptionRepo.UpdatePaymentStatus(ctx, uint(subscription.ID), domain.PaymentStatusFailed); err != nil {
		s.logger.Error().Err(err).Str("stripe_id", stripeID).Msg("Failed to handle payment failure")
		return fmt.Errorf("failed to handle payment failure: %w", err)
	}

	s.logger.Info().Str("stripe_id", stripeID).Msg("Payment failure handled")
	return nil
}

func (s *subscriptionService) HandlePaymentRefunded(ctx context.Context, stripeID string) error {
	// Try to find by subscription ID first
	subscription, err := s.userSubscriptionRepo.GetByStripeSubscriptionID(ctx, stripeID)
	if err != nil {
		return err
	}

	// If not found, try by payment intent ID
	if subscription == nil {
		subscription, err = s.userSubscriptionRepo.GetByStripePaymentIntentID(ctx, stripeID)
		if err != nil {
			return err
		}
	}

	if subscription == nil {
		return fmt.Errorf("subscription not found for Stripe ID: %s", stripeID)
	}

	// Update payment status to refunded and cancel the subscription
	if err := s.userSubscriptionRepo.UpdatePaymentStatus(ctx, uint(subscription.ID), domain.PaymentStatusRefunded); err != nil {
		s.logger.Error().Err(err).Str("stripe_id", stripeID).Msg("Failed to handle payment refund")
		return fmt.Errorf("failed to handle payment refund: %w", err)
	}

	s.logger.Info().Str("stripe_id", stripeID).Msg("Payment refund handled")
	return nil
}

// Administrative functions
func (s *subscriptionService) ResetWeeklyLimits(ctx context.Context) error {
	if err := s.userSubscriptionRepo.ResetWeeklyLimits(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to reset weekly limits")
		return fmt.Errorf("failed to reset weekly limits: %w", err)
	}

	s.logger.Info().Msg("Weekly limits reset successfully")
	return nil
}

func (s *subscriptionService) ResetMonthlyLimits(ctx context.Context) error {
	if err := s.userSubscriptionRepo.ResetMonthlyLimits(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to reset monthly limits")
		return fmt.Errorf("failed to reset monthly limits: %w", err)
	}

	s.logger.Info().Msg("Monthly limits reset successfully")
	return nil
}

func (s *subscriptionService) ProcessExpiredSubscriptions(ctx context.Context) error {
	// Get expired subscriptions
	expiredSubscriptions, err := s.userSubscriptionRepo.GetExpiredSubscriptions(ctx)
	if err != nil {
		return err
	}

	// Get expired packages
	expiredPackages, err := s.userSubscriptionRepo.GetExpiredPackages(ctx)
	if err != nil {
		return err
	}

	// Mark all as expired
	totalExpired := 0
	for _, subscription := range expiredSubscriptions {
		if err := s.userSubscriptionRepo.MarkAsExpired(ctx, uint(subscription.ID)); err != nil {
			s.logger.Error().Err(err).Int("subscription_id", subscription.ID).Msg("Failed to mark subscription as expired")
			continue
		}
		totalExpired++
	}

	for _, pkg := range expiredPackages {
		if err := s.userSubscriptionRepo.MarkAsExpired(ctx, uint(pkg.ID)); err != nil {
			s.logger.Error().Err(err).Int("package_id", pkg.ID).Msg("Failed to mark package as expired")
			continue
		}
		totalExpired++
	}

	s.logger.Info().Int("total_expired", totalExpired).Msg("Processed expired subscriptions and packages")
	return nil
}

// Plan seeding (for initial setup)
func (s *subscriptionService) SeedDefaultPlans(ctx context.Context) error {
	// Check if plans already exist
	existingPlans, err := s.subscriptionPlanRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(existingPlans) > 0 {
		s.logger.Info().Int("existing_plans", len(existingPlans)).Msg("Plans already exist, skipping seeding")
		return nil
	}

	// Create default subscription plans
	subscriptionPlans := domain.GetDefaultSubscriptionPlans()
	for _, plan := range subscriptionPlans {
		if err := s.subscriptionPlanRepo.Create(ctx, &plan); err != nil {
			s.logger.Error().Err(err).Str("plan_name", string(plan.Name)).Msg("Failed to create subscription plan")
			return fmt.Errorf("failed to create subscription plan %s: %w", plan.Name, err)
		}
	}

	// Create default package plans
	packagePlans := domain.GetDefaultPackagePlans()
	for _, plan := range packagePlans {
		if err := s.subscriptionPlanRepo.Create(ctx, &plan); err != nil {
			s.logger.Error().Err(err).Str("plan_name", string(plan.Name)).Msg("Failed to create package plan")
			return fmt.Errorf("failed to create package plan %s: %w", plan.Name, err)
		}
	}

	totalPlans := len(subscriptionPlans) + len(packagePlans)
	s.logger.Info().Int("total_plans", totalPlans).Msg("Default plans seeded successfully")
	return nil
}

// Additional methods for Stripe webhook integration
func (s *subscriptionService) ActivateSubscriptionByStripeID(ctx context.Context, stripeID string) error {
	// Try to find by subscription ID first
	subscription, err := s.userSubscriptionRepo.GetByStripeSubscriptionID(ctx, stripeID)
	if err != nil {
		return err
	}

	// If not found, try by payment intent ID
	if subscription == nil {
		subscription, err = s.userSubscriptionRepo.GetByStripePaymentIntentID(ctx, stripeID)
		if err != nil {
			return err
		}
	}

	if subscription == nil {
		return fmt.Errorf("subscription not found for Stripe ID: %s", stripeID)
	}

	// Activate the subscription
	return s.ActivateSubscription(ctx, uint(subscription.ID))
}

func (s *subscriptionService) CancelSubscriptionByStripeID(ctx context.Context, stripeID string) error {
	// Try to find by subscription ID first
	subscription, err := s.userSubscriptionRepo.GetByStripeSubscriptionID(ctx, stripeID)
	if err != nil {
		return err
	}

	// If not found, try by payment intent ID
	if subscription == nil {
		subscription, err = s.userSubscriptionRepo.GetByStripePaymentIntentID(ctx, stripeID)
		if err != nil {
			return err
		}
	}

	if subscription == nil {
		return fmt.Errorf("subscription not found for Stripe ID: %s", stripeID)
	}

	// Cancel the subscription
	return s.CancelSubscription(ctx, uint(subscription.ID))
}
