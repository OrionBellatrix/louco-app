package domain

import (
	"time"
)

type Ticket struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	EventID       int       `json:"event_id" gorm:"not null;index"`
	Title         string    `json:"title" gorm:"type:varchar(200);not null"`
	Price         float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	TotalQuantity int       `json:"total_quantity" gorm:"not null"`
	SoldQuantity  int       `json:"sold_quantity" gorm:"default:0"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Event Event `json:"event" gorm:"foreignKey:EventID;references:ID"`
}

func NewTicket(eventID int, title string, price float64, totalQuantity int) *Ticket {
	return &Ticket{
		EventID:       eventID,
		Title:         title,
		Price:         price,
		TotalQuantity: totalQuantity,
		SoldQuantity:  0,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (t *Ticket) UpdateInfo(title string, price float64, totalQuantity int) error {
	// Cannot reduce total quantity below sold quantity
	if totalQuantity < t.SoldQuantity {
		return ErrTicketQuantityBelowSold
	}

	if title != "" {
		t.Title = title
	}
	if price >= 0 {
		t.Price = price
	}
	if totalQuantity > 0 {
		t.TotalQuantity = totalQuantity
	}
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Ticket) IncrementSoldQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrTicketInvalidQuantity
	}

	if t.SoldQuantity+quantity > t.TotalQuantity {
		return ErrTicketInsufficientQuantity
	}

	t.SoldQuantity += quantity
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Ticket) DecrementSoldQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrTicketInvalidQuantity
	}

	if t.SoldQuantity-quantity < 0 {
		return ErrTicketInvalidQuantityDecrement
	}

	t.SoldQuantity -= quantity
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Ticket) Activate() {
	t.IsActive = true
	t.UpdatedAt = time.Now()
}

func (t *Ticket) Deactivate() {
	t.IsActive = false
	t.UpdatedAt = time.Now()
}

func (t *Ticket) GetAvailableQuantity() int {
	return t.TotalQuantity - t.SoldQuantity
}

func (t *Ticket) IsAvailable() bool {
	return t.IsActive && t.GetAvailableQuantity() > 0
}

func (t *Ticket) IsSoldOut() bool {
	return t.SoldQuantity >= t.TotalQuantity
}

func (t *Ticket) GetSoldPercentage() float64 {
	if t.TotalQuantity == 0 {
		return 0
	}
	return (float64(t.SoldQuantity) / float64(t.TotalQuantity)) * 100
}

func (t *Ticket) IsFree() bool {
	return t.Price == 0
}

func (t *Ticket) ValidateRequiredFields() error {
	if t.Title == "" {
		return ErrTicketTitleRequired
	}
	if t.Price < 0 {
		return ErrTicketInvalidPrice
	}
	if t.TotalQuantity <= 0 {
		return ErrTicketInvalidTotalQuantity
	}
	return nil
}

func (t *Ticket) CanBeUpdated() bool {
	// Tickets can be updated if no tickets have been sold yet
	return t.SoldQuantity == 0
}

func (t *Ticket) CanBeDeleted() bool {
	// Tickets can be deleted if no tickets have been sold yet
	return t.SoldQuantity == 0
}

// Ticket domain errors
var (
	ErrTicketTitleRequired            = NewDomainError("ticket title is required")
	ErrTicketInvalidPrice             = NewDomainError("ticket price cannot be negative")
	ErrTicketInvalidTotalQuantity     = NewDomainError("ticket total quantity must be greater than 0")
	ErrTicketInvalidQuantity          = NewDomainError("quantity must be greater than 0")
	ErrTicketInsufficientQuantity     = NewDomainError("insufficient ticket quantity available")
	ErrTicketQuantityBelowSold        = NewDomainError("total quantity cannot be less than sold quantity")
	ErrTicketInvalidQuantityDecrement = NewDomainError("cannot decrement sold quantity below 0")
	ErrTicketNotFound                 = NewDomainError("ticket not found")
	ErrTicketCannotBeUpdated          = NewDomainError("ticket cannot be updated after sales have started")
	ErrTicketCannotBeDeleted          = NewDomainError("ticket cannot be deleted after sales have started")
	ErrTicketSoldOut                  = NewDomainError("ticket is sold out")
	ErrTicketNotActive                = NewDomainError("ticket is not active")
)
