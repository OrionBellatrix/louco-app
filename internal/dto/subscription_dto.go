package dto

import (
	"encoding/json"
	"time"

	"github.com/louco-event/internal/domain"
)

// Plan DTOs
type SubscriptionPlanResponse struct {
	ID           int                     `json:"id"`
	Type         domain.SubscriptionType `json:"type"`
	Name         domain.SubscriptionName `json:"name"`
	DisplayName  string                  `json:"display_name"`
	Description  string                  `json:"description"`
	Price        float64                 `json:"price"`
	Currency     string                  `json:"currency"`
	WeeklyLimit  *int                    `json:"weekly_limit,omitempty"`
	MonthlyLimit *int                    `json:"monthly_limit,omitempty"`
	TotalCredits *int                    `json:"total_credits,omitempty"`
	DurationDays *int                    `json:"duration_days,omitempty"`
	IsActive     bool                    `json:"is_active"`
	SortOrder    int                     `json:"sort_order"`
	Features     []string                `json:"features,omitempty"`
	Popular      bool                    `json:"popular"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

type PlansResponse struct {
	Subscriptions []*SubscriptionPlanResponse `json:"subscriptions"`
	Packages      []*SubscriptionPlanResponse `json:"packages"`
}

// User Subscription DTOs
type UserSubscriptionResponse struct {
	ID           int                       `json:"id"`
	UserID       int                       `json:"user_id"`
	Type         domain.SubscriptionType   `json:"type"`
	Name         domain.SubscriptionName   `json:"name"`
	Price        float64                   `json:"price"`
	Currency     string                    `json:"currency"`
	WeeklyLimit  *int                      `json:"weekly_limit,omitempty"`
	MonthlyLimit *int                      `json:"monthly_limit,omitempty"`
	TotalCredits *int                      `json:"total_credits,omitempty"`
	UsedCredits  int                       `json:"used_credits"`
	WeeklyUsed   int                       `json:"weekly_used"`
	MonthlyUsed  int                       `json:"monthly_used"`
	Status       domain.SubscriptionStatus `json:"status"`
	StartedAt    *time.Time                `json:"started_at,omitempty"`
	ExpiredAt    *time.Time                `json:"expired_at,omitempty"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

type PublishingRightsResponse struct {
	CanPublish         bool                        `json:"can_publish"`
	WeeklyLimit        int                         `json:"weekly_limit"`
	WeeklyUsed         int                         `json:"weekly_used"`
	WeeklyRemaining    int                         `json:"weekly_remaining"`
	MonthlyLimit       int                         `json:"monthly_limit"`
	MonthlyUsed        int                         `json:"monthly_used"`
	MonthlyRemaining   int                         `json:"monthly_remaining"`
	TotalCredits       int                         `json:"total_credits"`
	UsedCredits        int                         `json:"used_credits"`
	RemainingCredits   int                         `json:"remaining_credits"`
	ActiveSubscription *UserSubscriptionResponse   `json:"active_subscription,omitempty"`
	ActivePackages     []*UserSubscriptionResponse `json:"active_packages,omitempty"`
	RestrictionReason  string                      `json:"restriction_reason,omitempty"`
}

type SubscriptionStatsResponse struct {
	TotalSubscriptions  int        `json:"total_subscriptions"`
	ActiveSubscriptions int        `json:"active_subscriptions"`
	TotalPackages       int        `json:"total_packages"`
	ActivePackages      int        `json:"active_packages"`
	TotalSpent          float64    `json:"total_spent"`
	EventsPublished     int        `json:"events_published"`
	EventsRemaining     int        `json:"events_remaining"`
	CurrentPeriodUsage  int        `json:"current_period_usage"`
	LastPaymentDate     *time.Time `json:"last_payment_date,omitempty"`
	NextBillingDate     *time.Time `json:"next_billing_date,omitempty"`
}

// Purchase DTOs
type PurchaseSubscriptionRequest struct {
	PlanID          uint   `json:"plan_id" validate:"required"`
	PaymentMethodID string `json:"payment_method_id" validate:"required"`
}

type PurchasePackageRequest struct {
	PlanID          uint   `json:"plan_id" validate:"required"`
	PaymentMethodID string `json:"payment_method_id" validate:"required"`
}

type PurchaseResponse struct {
	SubscriptionID  int    `json:"subscription_id"`
	ClientSecret    string `json:"client_secret,omitempty"`
	Status          string `json:"status"`
	RequiresAction  bool   `json:"requires_action"`
	PaymentIntentID string `json:"payment_intent_id,omitempty"`
	SubscriptionID2 string `json:"stripe_subscription_id,omitempty"`
}

// Subscription Management DTOs
type CancelSubscriptionRequest struct {
	Reason string `json:"reason,omitempty"`
}

type UpdateSubscriptionRequest struct {
	PlanID uint `json:"plan_id" validate:"required"`
}

// Webhook DTOs
type StripeWebhookRequest struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type StripeEventData struct {
	Object json.RawMessage `json:"object"`
}

// History DTOs
type SubscriptionHistoryRequest struct {
	PaginationRequest
	Type   string `json:"type,omitempty" form:"type"`
	Status string `json:"status,omitempty" form:"status"`
}

type SubscriptionHistoryResponse struct {
	Subscriptions []*UserSubscriptionResponse `json:"subscriptions"`
	Pagination    PaginationResponse          `json:"pagination"`
}

// Conversion functions
func ToSubscriptionPlanResponse(plan *domain.SubscriptionPlan) *SubscriptionPlanResponse {
	response := &SubscriptionPlanResponse{
		ID:           plan.ID,
		Type:         plan.Type,
		Name:         plan.Name,
		DisplayName:  plan.DisplayName,
		Description:  plan.Description,
		Price:        plan.Price,
		Currency:     plan.Currency,
		WeeklyLimit:  plan.WeeklyLimit,
		MonthlyLimit: plan.MonthlyLimit,
		TotalCredits: plan.TotalCredits,
		DurationDays: plan.DurationDays,
		IsActive:     plan.IsActive,
		SortOrder:    plan.SortOrder,
		CreatedAt:    plan.CreatedAt,
		UpdatedAt:    plan.UpdatedAt,
	}

	// Parse metadata for features and popular flag
	if len(plan.Metadata) > 0 {
		var metadata map[string]interface{}
		if err := json.Unmarshal(plan.Metadata, &metadata); err == nil {
			if features, ok := metadata["features"].([]interface{}); ok {
				for _, feature := range features {
					if featureStr, ok := feature.(string); ok {
						response.Features = append(response.Features, featureStr)
					}
				}
			}
			if popular, ok := metadata["popular"].(bool); ok {
				response.Popular = popular
			}
		}
	}

	return response
}

func ToUserSubscriptionResponse(subscription *domain.UserSubscription) *UserSubscriptionResponse {
	return &UserSubscriptionResponse{
		ID:           subscription.ID,
		UserID:       subscription.UserID,
		Type:         subscription.Type,
		Name:         subscription.Name,
		Price:        subscription.Price,
		Currency:     subscription.Currency,
		WeeklyLimit:  subscription.WeeklyLimit,
		MonthlyLimit: subscription.MonthlyLimit,
		TotalCredits: subscription.TotalCredits,
		UsedCredits:  subscription.UsedCredits,
		WeeklyUsed:   subscription.WeeklyUsed,
		MonthlyUsed:  subscription.MonthlyUsed,
		Status:       subscription.Status,
		StartedAt:    subscription.StartedAt,
		ExpiredAt:    subscription.ExpiredAt,
		CreatedAt:    subscription.CreatedAt,
		UpdatedAt:    subscription.UpdatedAt,
	}
}

func ToPublishingRightsResponse(rights *domain.PublishingRights) *PublishingRightsResponse {
	response := &PublishingRightsResponse{
		CanPublish:        rights.CanPublish,
		WeeklyLimit:       rights.WeeklyLimit,
		WeeklyUsed:        rights.WeeklyUsed,
		WeeklyRemaining:   rights.WeeklyRemaining,
		MonthlyLimit:      rights.MonthlyLimit,
		MonthlyUsed:       rights.MonthlyUsed,
		MonthlyRemaining:  rights.MonthlyRemaining,
		TotalCredits:      rights.TotalCredits,
		UsedCredits:       rights.UsedCredits,
		RemainingCredits:  rights.RemainingCredits,
		RestrictionReason: rights.RestrictionReason,
	}

	if rights.ActiveSubscription != nil {
		response.ActiveSubscription = ToUserSubscriptionResponse(rights.ActiveSubscription)
	}

	if len(rights.ActivePackages) > 0 {
		response.ActivePackages = make([]*UserSubscriptionResponse, len(rights.ActivePackages))
		for i, pkg := range rights.ActivePackages {
			response.ActivePackages[i] = ToUserSubscriptionResponse(pkg)
		}
	}

	return response
}

func ToSubscriptionStatsResponse(stats *domain.SubscriptionStats) *SubscriptionStatsResponse {
	return &SubscriptionStatsResponse{
		TotalSubscriptions:  stats.TotalSubscriptions,
		ActiveSubscriptions: stats.ActiveSubscriptions,
		TotalPackages:       stats.TotalPackages,
		ActivePackages:      stats.ActivePackages,
		TotalSpent:          stats.TotalSpent,
		EventsPublished:     stats.EventsPublished,
		EventsRemaining:     stats.EventsRemaining,
		CurrentPeriodUsage:  stats.CurrentPeriodUsage,
		LastPaymentDate:     stats.LastPaymentDate,
		NextBillingDate:     stats.NextBillingDate,
	}
}

func ToSubscriptionPlanResponses(plans []*domain.SubscriptionPlan) []*SubscriptionPlanResponse {
	responses := make([]*SubscriptionPlanResponse, len(plans))
	for i, plan := range plans {
		responses[i] = ToSubscriptionPlanResponse(plan)
	}
	return responses
}

func ToUserSubscriptionResponses(subscriptions []*domain.UserSubscription) []*UserSubscriptionResponse {
	responses := make([]*UserSubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responses[i] = ToUserSubscriptionResponse(subscription)
	}
	return responses
}

func ToPlansResponse(subscriptions, packages []*domain.SubscriptionPlan) *PlansResponse {
	return &PlansResponse{
		Subscriptions: ToSubscriptionPlanResponses(subscriptions),
		Packages:      ToSubscriptionPlanResponses(packages),
	}
}
