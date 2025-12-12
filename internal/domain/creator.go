package domain

import (
	"time"
)

type Creator struct {
	ID               int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           int       `json:"user_id" gorm:"not null;uniqueIndex"`
	WeeztixToken     *string   `json:"weeztix_token" gorm:"type:json"`
	CompanyName      string    `json:"company_name" gorm:"type:varchar(200);not null"`
	Address          string    `json:"address" gorm:"type:varchar(500);not null"`
	EstimatedTickets int       `json:"estimated_tickets" gorm:"not null"`
	EstimatedEvents  int       `json:"estimated_events" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	User       User       `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Industries []Industry `json:"industries" gorm:"many2many:creator_industries;"`
}

// CreatorIndustry represents the many-to-many relationship between creators and industries
type CreatorIndustry struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatorID  int       `json:"creator_id" gorm:"not null"`
	IndustryID int       `json:"industry_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relations
	Creator  Creator  `json:"creator" gorm:"foreignKey:CreatorID;references:ID"`
	Industry Industry `json:"industry" gorm:"foreignKey:IndustryID;references:ID"`
}

func NewCreator(userID int, companyName, address string, estimatedTickets, estimatedEvents int) *Creator {
	return &Creator{
		UserID:           userID,
		CompanyName:      companyName,
		Address:          address,
		EstimatedTickets: estimatedTickets,
		EstimatedEvents:  estimatedEvents,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func (c *Creator) SetWeeztixToken(token string) {
	c.WeeztixToken = &token
	c.UpdatedAt = time.Now()
}

func (c *Creator) UpdateProfile(companyName, address string, estimatedTickets, estimatedEvents int) {
	if companyName != "" {
		c.CompanyName = companyName
	}
	if address != "" {
		c.Address = address
	}
	if estimatedTickets > 0 {
		c.EstimatedTickets = estimatedTickets
	}
	if estimatedEvents > 0 {
		c.EstimatedEvents = estimatedEvents
	}
	c.UpdatedAt = time.Now()
}

func (c *Creator) AddIndustry(industry Industry) {
	c.Industries = append(c.Industries, industry)
	c.UpdatedAt = time.Now()
}

func (c *Creator) SetIndustries(industries []Industry) {
	c.Industries = industries
	c.UpdatedAt = time.Now()
}

func (c *Creator) HasIndustry(industryID int) bool {
	for _, industry := range c.Industries {
		if industry.ID == industryID {
			return true
		}
	}
	return false
}

func (c *Creator) ValidateRequiredFields() error {
	if c.CompanyName == "" {
		return ErrCreatorCompanyNameRequired
	}
	if c.Address == "" {
		return ErrCreatorAddressRequired
	}
	if c.EstimatedTickets <= 0 {
		return ErrCreatorEstimatedTicketsRequired
	}
	if c.EstimatedEvents <= 0 {
		return ErrCreatorEstimatedEventsRequired
	}
	return nil
}

// Creator domain errors
var (
	ErrCreatorCompanyNameRequired      = NewDomainError("company name is required for creator")
	ErrCreatorAddressRequired          = NewDomainError("address is required for creator")
	ErrCreatorEstimatedTicketsRequired = NewDomainError("estimated tickets is required for creator")
	ErrCreatorEstimatedEventsRequired  = NewDomainError("estimated events is required for creator")
	ErrCreatorIndustryRequired         = NewDomainError("at least one industry selection is required for creator")
)
