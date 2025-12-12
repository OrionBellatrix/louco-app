package domain

import (
	"encoding/json"
	"time"
)

// SubscriptionType represents the type of subscription or package
type SubscriptionType string

const (
	SubscriptionTypeSubscription SubscriptionType = "subscription"
	SubscriptionTypePackage      SubscriptionType = "package"
)

// SubscriptionStatus represents the status of a subscription or package
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusPending   SubscriptionStatus = "pending"
)

// SubscriptionName represents predefined subscription and package names
type SubscriptionName string

const (
	// Subscription names
	SubscriptionNameBasic SubscriptionName = "basic"
	SubscriptionNamePlus  SubscriptionName = "plus"
	SubscriptionNamePro   SubscriptionName = "pro"

	// Package names
	PackageNameSingle PackageName = "single_event"
	PackageName10     PackageName = "10_events"
	PackageName25     PackageName = "25_events"
)

type PackageName = SubscriptionName

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PlanType represents the type of plan (alias for SubscriptionType)
type PlanType = SubscriptionType

// PublishingRights represents user's current publishing rights
type PublishingRights struct {
	CanPublish         bool                `json:"can_publish"`
	WeeklyLimit        int                 `json:"weekly_limit"`
	WeeklyUsed         int                 `json:"weekly_used"`
	WeeklyRemaining    int                 `json:"weekly_remaining"`
	MonthlyLimit       int                 `json:"monthly_limit"`
	MonthlyUsed        int                 `json:"monthly_used"`
	MonthlyRemaining   int                 `json:"monthly_remaining"`
	TotalCredits       int                 `json:"total_credits"`
	UsedCredits        int                 `json:"used_credits"`
	RemainingCredits   int                 `json:"remaining_credits"`
	ActiveSubscription *UserSubscription   `json:"active_subscription,omitempty"`
	ActivePackages     []*UserSubscription `json:"active_packages,omitempty"`
	RestrictionReason  string              `json:"restriction_reason,omitempty"`
}

// SubscriptionStats represents user's subscription statistics
type SubscriptionStats struct {
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

// UserSubscription represents both subscriptions and packages in a unified table
type UserSubscription struct {
	ID           int                `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int                `gorm:"not null;index" json:"user_id"`
	Type         SubscriptionType   `gorm:"type:varchar(20);not null;index" json:"type"`
	Name         SubscriptionName   `gorm:"type:varchar(50);not null" json:"name"`
	Price        float64            `gorm:"type:decimal(10,2);not null" json:"price"`
	Currency     string             `gorm:"type:varchar(3);not null;default:'EUR'" json:"currency"`
	WeeklyLimit  *int               `gorm:"default:null" json:"weekly_limit,omitempty"`  // Only for subscriptions
	MonthlyLimit *int               `gorm:"default:null" json:"monthly_limit,omitempty"` // Only for subscriptions
	TotalCredits *int               `gorm:"default:null" json:"total_credits,omitempty"` // Only for packages
	UsedCredits  int                `gorm:"default:0" json:"used_credits"`               // Used credits for packages
	WeeklyUsed   int                `gorm:"default:0" json:"weekly_used"`                // Used this week for subscriptions
	MonthlyUsed  int                `gorm:"default:0" json:"monthly_used"`               // Used this month for subscriptions
	Status       SubscriptionStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	StartedAt    *time.Time         `gorm:"default:null" json:"started_at,omitempty"`
	ExpiredAt    *time.Time         `gorm:"default:null" json:"expired_at,omitempty"`
	StripeID     *string            `gorm:"type:varchar(255);default:null;index" json:"stripe_id,omitempty"` // Stripe subscription/payment ID
	Metadata     json.RawMessage    `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	CreatedAt    time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time          `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName returns the table name for UserSubscription
func (UserSubscription) TableName() string {
	return "user_subscriptions"
}

// IsActive checks if the subscription/package is currently active
func (us *UserSubscription) IsActive() bool {
	if us.Status != SubscriptionStatusActive {
		return false
	}

	if us.ExpiredAt != nil && us.ExpiredAt.Before(time.Now()) {
		return false
	}

	return true
}

// IsExpired checks if the subscription/package has expired
func (us *UserSubscription) IsExpired() bool {
	return us.ExpiredAt != nil && us.ExpiredAt.Before(time.Now())
}

// HasWeeklyLimit checks if user has reached weekly limit (for subscriptions)
func (us *UserSubscription) HasWeeklyLimit() bool {
	if us.Type != SubscriptionTypeSubscription || us.WeeklyLimit == nil {
		return false
	}
	return us.WeeklyUsed >= *us.WeeklyLimit
}

// HasMonthlyLimit checks if user has reached monthly limit (for subscriptions)
func (us *UserSubscription) HasMonthlyLimit() bool {
	if us.Type != SubscriptionTypeSubscription || us.MonthlyLimit == nil {
		return false
	}
	return us.MonthlyUsed >= *us.MonthlyLimit
}

// HasCredits checks if package has remaining credits (for packages)
func (us *UserSubscription) HasCredits() bool {
	if us.Type != SubscriptionTypePackage || us.TotalCredits == nil {
		return false
	}
	return us.UsedCredits < *us.TotalCredits
}

// RemainingCredits returns remaining credits for packages
func (us *UserSubscription) RemainingCredits() int {
	if us.Type != SubscriptionTypePackage || us.TotalCredits == nil {
		return 0
	}
	remaining := *us.TotalCredits - us.UsedCredits
	if remaining < 0 {
		return 0
	}
	return remaining
}

// RemainingWeeklyLimit returns remaining weekly limit for subscriptions
func (us *UserSubscription) RemainingWeeklyLimit() int {
	if us.Type != SubscriptionTypeSubscription || us.WeeklyLimit == nil {
		return 0
	}
	remaining := *us.WeeklyLimit - us.WeeklyUsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// RemainingMonthlyLimit returns remaining monthly limit for subscriptions
func (us *UserSubscription) RemainingMonthlyLimit() int {
	if us.Type != SubscriptionTypeSubscription || us.MonthlyLimit == nil {
		return 0
	}
	remaining := *us.MonthlyLimit - us.MonthlyUsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// CanPublishEvent checks if user can publish an event based on subscription/package rules
func (us *UserSubscription) CanPublishEvent() bool {
	if !us.IsActive() {
		return false
	}

	switch us.Type {
	case SubscriptionTypeSubscription:
		// Check both weekly and monthly limits
		return !us.HasWeeklyLimit() && !us.HasMonthlyLimit()
	case SubscriptionTypePackage:
		// Check remaining credits
		return us.HasCredits()
	default:
		return false
	}
}

// ConsumeUsage consumes one usage (for event publishing)
func (us *UserSubscription) ConsumeUsage() {
	switch us.Type {
	case SubscriptionTypeSubscription:
		us.WeeklyUsed++
		us.MonthlyUsed++
	case SubscriptionTypePackage:
		us.UsedCredits++
	}
}

// ResetWeeklyUsage resets weekly usage counter (called every Monday)
func (us *UserSubscription) ResetWeeklyUsage() {
	if us.Type == SubscriptionTypeSubscription {
		us.WeeklyUsed = 0
	}
}

// ResetMonthlyUsage resets monthly usage counter (called on 1st of each month)
func (us *UserSubscription) ResetMonthlyUsage() {
	if us.Type == SubscriptionTypeSubscription {
		us.MonthlyUsed = 0
	}
}

// Validate validates the UserSubscription entity
func (us *UserSubscription) Validate() error {
	if us.UserID <= 0 {
		return NewDomainError("User ID is required")
	}

	if us.Type == "" {
		return NewDomainError("Type is required")
	}

	if us.Type != SubscriptionTypeSubscription && us.Type != SubscriptionTypePackage {
		return NewDomainError("Type must be 'subscription' or 'package'")
	}

	if us.Name == "" {
		return NewDomainError("Name is required")
	}

	if us.Price < 0 {
		return NewDomainError("Price cannot be negative")
	}

	if us.Currency == "" {
		us.Currency = "EUR" // Default currency
	}

	// Validate subscription-specific fields
	if us.Type == SubscriptionTypeSubscription {
		if us.WeeklyLimit == nil || *us.WeeklyLimit <= 0 {
			return NewDomainError("Weekly limit is required for subscriptions")
		}
		if us.MonthlyLimit == nil || *us.MonthlyLimit <= 0 {
			return NewDomainError("Monthly limit is required for subscriptions")
		}
		if us.TotalCredits != nil {
			return NewDomainError("Total credits should not be set for subscriptions")
		}
	}

	// Validate package-specific fields
	if us.Type == SubscriptionTypePackage {
		if us.TotalCredits == nil || *us.TotalCredits <= 0 {
			return NewDomainError("Total credits is required for packages")
		}
		if us.WeeklyLimit != nil || us.MonthlyLimit != nil {
			return NewDomainError("Weekly/Monthly limits should not be set for packages")
		}
	}

	if us.UsedCredits < 0 {
		return NewDomainError("Used credits cannot be negative")
	}

	if us.WeeklyUsed < 0 {
		return NewDomainError("Weekly used cannot be negative")
	}

	if us.MonthlyUsed < 0 {
		return NewDomainError("Monthly used cannot be negative")
	}

	// Validate metadata JSON
	if len(us.Metadata) > 0 {
		var temp interface{}
		if err := json.Unmarshal(us.Metadata, &temp); err != nil {
			return NewDomainError("Invalid JSON format in metadata")
		}
	}

	return nil
}
