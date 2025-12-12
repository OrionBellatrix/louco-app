package domain

import (
	"encoding/json"
	"time"
)

// BillingCycle represents the billing cycle for subscriptions
type BillingCycle string

const (
	BillingCycleWeekly  BillingCycle = "weekly"
	BillingCycleMonthly BillingCycle = "monthly"
)

// SubscriptionPlan represents available subscription and package plans
type SubscriptionPlan struct {
	ID           int              `gorm:"primaryKey;autoIncrement" json:"id"`
	Type         SubscriptionType `gorm:"type:varchar(20);not null;index" json:"type"`
	Name         SubscriptionName `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	DisplayName  string           `gorm:"type:varchar(100);not null" json:"display_name"`
	Description  string           `gorm:"type:text" json:"description"`
	Price        float64          `gorm:"type:decimal(10,2);not null" json:"price"`
	Currency     string           `gorm:"type:varchar(3);not null;default:'EUR'" json:"currency"`
	BillingCycle *BillingCycle    `gorm:"type:varchar(20);default:null" json:"billing_cycle,omitempty"` // Only for subscriptions
	WeeklyLimit  *int             `gorm:"default:null" json:"weekly_limit,omitempty"`                   // Only for subscriptions
	MonthlyLimit *int             `gorm:"default:null" json:"monthly_limit,omitempty"`                  // Only for subscriptions
	TotalCredits *int             `gorm:"default:null" json:"total_credits,omitempty"`                  // Only for packages
	DurationDays *int             `gorm:"default:null" json:"duration_days,omitempty"`                  // Duration in days (365 for packages, 30 for subscriptions)
	IsActive     bool             `gorm:"default:true" json:"is_active"`
	SortOrder    int              `gorm:"default:0" json:"sort_order"`                                     // For display ordering
	StripeID     *string          `gorm:"type:varchar(255);default:null;index" json:"stripe_id,omitempty"` // Stripe price/product ID
	Metadata     json.RawMessage  `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	CreatedAt    time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for SubscriptionPlan
func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

// Validate validates the SubscriptionPlan entity
func (sp *SubscriptionPlan) Validate() error {
	if sp.Type == "" {
		return NewDomainError("Type is required")
	}

	if sp.Type != SubscriptionTypeSubscription && sp.Type != SubscriptionTypePackage {
		return NewDomainError("Type must be 'subscription' or 'package'")
	}

	if sp.Name == "" {
		return NewDomainError("Name is required")
	}

	if sp.DisplayName == "" {
		return NewDomainError("Display name is required")
	}

	if sp.Price < 0 {
		return NewDomainError("Price cannot be negative")
	}

	if sp.Currency == "" {
		sp.Currency = "EUR" // Default currency
	}

	// Validate subscription-specific fields
	if sp.Type == SubscriptionTypeSubscription {
		if sp.WeeklyLimit == nil || *sp.WeeklyLimit <= 0 {
			return NewDomainError("Weekly limit is required for subscriptions")
		}
		if sp.MonthlyLimit == nil || *sp.MonthlyLimit <= 0 {
			return NewDomainError("Monthly limit is required for subscriptions")
		}
		if sp.TotalCredits != nil {
			return NewDomainError("Total credits should not be set for subscriptions")
		}
		if sp.DurationDays == nil {
			defaultDuration := 30 // 30 days for monthly subscriptions
			sp.DurationDays = &defaultDuration
		}
	}

	// Validate package-specific fields
	if sp.Type == SubscriptionTypePackage {
		if sp.TotalCredits == nil || *sp.TotalCredits <= 0 {
			return NewDomainError("Total credits is required for packages")
		}
		if sp.WeeklyLimit != nil || sp.MonthlyLimit != nil {
			return NewDomainError("Weekly/Monthly limits should not be set for packages")
		}
		if sp.DurationDays == nil {
			defaultDuration := 365 // 1 year for packages
			sp.DurationDays = &defaultDuration
		}
	}

	// Validate metadata JSON
	if len(sp.Metadata) > 0 {
		var temp interface{}
		if err := json.Unmarshal(sp.Metadata, &temp); err != nil {
			return NewDomainError("Invalid JSON format in metadata")
		}
	}

	return nil
}

// IsSubscription checks if the plan is a subscription
func (sp *SubscriptionPlan) IsSubscription() bool {
	return sp.Type == SubscriptionTypeSubscription
}

// IsPackage checks if the plan is a package
func (sp *SubscriptionPlan) IsPackage() bool {
	return sp.Type == SubscriptionTypePackage
}

// GetDurationDays returns the duration in days
func (sp *SubscriptionPlan) GetDurationDays() int {
	if sp.DurationDays == nil {
		if sp.IsSubscription() {
			return 30 // Default 30 days for subscriptions
		}
		return 365 // Default 1 year for packages
	}
	return *sp.DurationDays
}

// CreateUserSubscription creates a UserSubscription from this plan for a user
func (sp *SubscriptionPlan) CreateUserSubscription(userID int) *UserSubscription {
	now := time.Now()
	expiredAt := now.AddDate(0, 0, sp.GetDurationDays())

	us := &UserSubscription{
		UserID:    userID,
		Type:      sp.Type,
		Name:      sp.Name,
		Price:     sp.Price,
		Currency:  sp.Currency,
		Status:    SubscriptionStatusPending,
		StartedAt: &now,
		ExpiredAt: &expiredAt,
		Metadata:  sp.Metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Copy plan-specific fields
	if sp.IsSubscription() {
		us.WeeklyLimit = sp.WeeklyLimit
		us.MonthlyLimit = sp.MonthlyLimit
	} else if sp.IsPackage() {
		us.TotalCredits = sp.TotalCredits
	}

	return us
}

// GetDefaultSubscriptionPlans returns the default subscription plans
func GetDefaultSubscriptionPlans() []SubscriptionPlan {
	return []SubscriptionPlan{
		{
			Type:         SubscriptionTypeSubscription,
			Name:         SubscriptionNameBasic,
			DisplayName:  "Basic Plan",
			Description:  "Perfect for getting started with event creation",
			Price:        78.00,
			Currency:     "EUR",
			BillingCycle: billingCyclePtr(BillingCycleMonthly),
			WeeklyLimit:  intPtr(1),
			MonthlyLimit: intPtr(4),
			DurationDays: intPtr(30),
			IsActive:     true,
			SortOrder:    1,
			Metadata:     json.RawMessage(`{"features": ["1 event per week", "4 events per month", "Basic support"], "popular": false}`),
		},
		{
			Type:         SubscriptionTypeSubscription,
			Name:         SubscriptionNamePlus,
			DisplayName:  "Plus Plan",
			Description:  "Great for regular event creators",
			Price:        130.00,
			Currency:     "EUR",
			BillingCycle: billingCyclePtr(BillingCycleMonthly),
			WeeklyLimit:  intPtr(2),
			MonthlyLimit: intPtr(8),
			DurationDays: intPtr(30),
			IsActive:     true,
			SortOrder:    2,
			Metadata:     json.RawMessage(`{"features": ["2 events per week", "8 events per month", "Priority support"], "popular": true}`),
		},
		{
			Type:         SubscriptionTypeSubscription,
			Name:         SubscriptionNamePro,
			DisplayName:  "Pro Plan",
			Description:  "For professional event organizers",
			Price:        156.00,
			Currency:     "EUR",
			BillingCycle: billingCyclePtr(BillingCycleMonthly),
			WeeklyLimit:  intPtr(3),
			MonthlyLimit: intPtr(12),
			DurationDays: intPtr(30),
			IsActive:     true,
			SortOrder:    3,
			Metadata:     json.RawMessage(`{"features": ["3 events per week", "12 events per month", "Premium support", "Advanced analytics"], "popular": false}`),
		},
	}
}

// GetDefaultPackagePlans returns the default package plans
func GetDefaultPackagePlans() []SubscriptionPlan {
	return []SubscriptionPlan{
		{
			Type:         SubscriptionTypePackage,
			Name:         PackageNameSingle,
			DisplayName:  "Single Event",
			Description:  "Perfect for one-time events",
			Price:        29.00,
			Currency:     "EUR",
			TotalCredits: intPtr(1),
			DurationDays: intPtr(365),
			IsActive:     true,
			SortOrder:    1,
			Metadata:     json.RawMessage(`{"features": ["1 event credit", "Valid for 1 year"], "popular": false}`),
		},
		{
			Type:         SubscriptionTypePackage,
			Name:         PackageName10,
			DisplayName:  "10 Events Package",
			Description:  "Better value for multiple events",
			Price:        249.99,
			Currency:     "EUR",
			TotalCredits: intPtr(10),
			DurationDays: intPtr(365),
			IsActive:     true,
			SortOrder:    2,
			Metadata:     json.RawMessage(`{"features": ["10 event credits", "Valid for 1 year", "Better value"], "popular": true}`),
		},
		{
			Type:         SubscriptionTypePackage,
			Name:         PackageName25,
			DisplayName:  "25 Events Package",
			Description:  "Best value for frequent event creators",
			Price:        499.99,
			Currency:     "EUR",
			TotalCredits: intPtr(25),
			DurationDays: intPtr(365),
			IsActive:     true,
			SortOrder:    3,
			Metadata:     json.RawMessage(`{"features": ["25 event credits", "Valid for 1 year", "Best value", "Bulk discount"], "popular": false}`),
		},
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

// Helper function to create BillingCycle pointer
func billingCyclePtr(bc BillingCycle) *BillingCycle {
	return &bc
}
