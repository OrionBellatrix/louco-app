package domain

import (
	"time"
)

type EventType string
type EventLocationType string
type EventStatus string

const (
	// Event Types
	EventTypePublic  EventType = "public"
	EventTypePrivate EventType = "private"

	// Event Location Types
	EventLocationTypeLocation     EventLocationType = "location"
	EventLocationTypeOnline       EventLocationType = "online"
	EventLocationTypeAnnouncement EventLocationType = "announcement"

	// Event Status
	EventStatusDraft     EventStatus = "draft"
	EventStatusPending   EventStatus = "pending"
	EventStatusRejected  EventStatus = "rejected"
	EventStatusStopped   EventStatus = "stopped"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusPublished EventStatus = "published"
)

type Event struct {
	ID           int               `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatorID    int               `json:"creator_id" gorm:"not null;index"`
	Name         string            `json:"name" gorm:"type:varchar(200);not null"`
	Description  *string           `json:"description" gorm:"type:text"`
	ImageID      *int              `json:"image_id" gorm:"index"`
	VideoID      *int              `json:"video_id" gorm:"index"`
	Type         EventType         `json:"type" gorm:"type:varchar(20);not null"`
	LocationType EventLocationType `json:"location_type" gorm:"type:varchar(20);not null"`
	Status       EventStatus       `json:"status" gorm:"type:varchar(20);not null;default:'draft'"`

	// Date and Time fields
	StartDate *time.Time `json:"start_date" gorm:"type:date"`
	StartTime *time.Time `json:"start_time" gorm:"type:time"`
	EndDate   *time.Time `json:"end_date" gorm:"type:date"`
	EndTime   *time.Time `json:"end_time" gorm:"type:time"`

	// Location specific fields
	AddressID *int `json:"address_id" gorm:"index"`

	// Online specific fields
	OnlineEventURL  *string `json:"online_event_url" gorm:"type:varchar(500)"`
	OnlineEventType *string `json:"online_event_type" gorm:"type:varchar(50)"` // zoom, google_meet, teams, custom, etc.

	// Ticket fields
	TicketURL        *string `json:"ticket_url" gorm:"type:varchar(500)"`
	HasSystemTickets bool    `json:"has_system_tickets" gorm:"default:false"`

	// Additional info
	AdditionalInfo *string `json:"additional_info" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Creator     Creator      `json:"creator" gorm:"foreignKey:CreatorID;references:ID"`
	Image       *Media       `json:"image,omitempty" gorm:"foreignKey:ImageID;references:ID"`
	Video       *Media       `json:"video,omitempty" gorm:"foreignKey:VideoID;references:ID"`
	Address     *Address     `json:"address,omitempty" gorm:"foreignKey:AddressID;references:ID"`
	Categories  []Category   `json:"categories" gorm:"many2many:event_categories;"`
	Tickets     []Ticket     `json:"tickets" gorm:"foreignKey:EventID;references:ID"`
	Invitations []Invitation `json:"invitations" gorm:"foreignKey:EventID;references:ID"`
}

// EventCategory represents the many-to-many relationship between events and categories
type EventCategory struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	EventID    int       `json:"event_id" gorm:"not null"`
	CategoryID int       `json:"category_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relations
	Event    Event    `json:"event" gorm:"foreignKey:EventID;references:ID"`
	Category Category `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
}

func NewEvent(creatorID int, name string, eventType EventType, locationType EventLocationType) *Event {
	return &Event{
		CreatorID:    creatorID,
		Name:         name,
		Type:         eventType,
		LocationType: locationType,
		Status:       EventStatusDraft,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (e *Event) SetDescription(description string) {
	e.Description = &description
	e.UpdatedAt = time.Now()
}

func (e *Event) SetAdditionalInfo(info string) {
	e.AdditionalInfo = &info
	e.UpdatedAt = time.Now()
}

func (e *Event) SetImage(imageID int) {
	e.ImageID = &imageID
	e.UpdatedAt = time.Now()
}

func (e *Event) SetVideo(videoID int) {
	e.VideoID = &videoID
	e.UpdatedAt = time.Now()
}

func (e *Event) SetAddress(addressID int) {
	e.AddressID = &addressID
	e.UpdatedAt = time.Now()
}

func (e *Event) SetOnlineEventURL(url, eventType string) {
	e.OnlineEventURL = &url
	e.OnlineEventType = &eventType
	e.UpdatedAt = time.Now()
}

func (e *Event) RemoveOnlineEventURL() {
	e.OnlineEventURL = nil
	e.OnlineEventType = nil
	e.UpdatedAt = time.Now()
}

func (e *Event) SetTicketURL(url string) {
	e.TicketURL = &url
	e.UpdatedAt = time.Now()
}

func (e *Event) EnableSystemTickets() {
	e.HasSystemTickets = true
	e.UpdatedAt = time.Now()
}

func (e *Event) DisableSystemTickets() {
	e.HasSystemTickets = false
	e.UpdatedAt = time.Now()
}

func (e *Event) SetDateTime(startDate, endDate *time.Time, startTime, endTime *time.Time) {
	e.StartDate = startDate
	e.EndDate = endDate
	e.StartTime = startTime
	e.EndTime = endTime
	e.UpdatedAt = time.Now()
}

func (e *Event) SetCategories(categories []Category) {
	e.Categories = categories
	e.UpdatedAt = time.Now()
}

func (e *Event) AddCategory(category Category) {
	e.Categories = append(e.Categories, category)
	e.UpdatedAt = time.Now()
}

func (e *Event) HasCategory(categoryID int) bool {
	for _, category := range e.Categories {
		if category.ID == categoryID {
			return true
		}
	}
	return false
}

// Status management methods
func (e *Event) SubmitForReview() error {
	if e.Status != EventStatusDraft {
		return ErrEventInvalidStatusTransition
	}

	if err := e.ValidateForSubmission(); err != nil {
		return err
	}

	e.Status = EventStatusPending
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) Approve() error {
	if e.Status != EventStatusPending {
		return ErrEventInvalidStatusTransition
	}

	e.Status = EventStatusPublished
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) Reject() error {
	if e.Status != EventStatusPending {
		return ErrEventInvalidStatusTransition
	}

	e.Status = EventStatusRejected
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) Stop() error {
	if e.Status != EventStatusPublished {
		return ErrEventInvalidStatusTransition
	}

	e.Status = EventStatusStopped
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) Cancel() error {
	if e.Status == EventStatusCancelled {
		return ErrEventAlreadyCancelled
	}

	e.Status = EventStatusCancelled
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) BackToDraft() error {
	if e.Status != EventStatusRejected {
		return ErrEventInvalidStatusTransition
	}

	e.Status = EventStatusDraft
	e.UpdatedAt = time.Now()
	return nil
}

// Validation methods
func (e *Event) ValidateForDraft() error {
	if e.Name == "" {
		return ErrEventNameRequired
	}
	if e.ImageID == nil {
		return ErrEventImageRequired
	}
	if e.VideoID == nil {
		return ErrEventVideoRequired
	}
	return nil
}

func (e *Event) ValidateForSubmission() error {
	// First validate draft requirements
	if err := e.ValidateForDraft(); err != nil {
		return err
	}

	// Validate based on location type
	switch e.LocationType {
	case EventLocationTypeLocation:
		if err := e.ValidateLocationEvent(); err != nil {
			return err
		}
	case EventLocationTypeOnline:
		if err := e.ValidateOnlineEvent(); err != nil {
			return err
		}
	case EventLocationTypeAnnouncement:
		if err := e.ValidateAnnouncementEvent(); err != nil {
			return err
		}
	default:
		return ErrEventInvalidLocationType
	}

	// Validate categories
	if len(e.Categories) == 0 {
		return ErrEventCategoriesRequired
	}

	return nil
}

func (e *Event) ValidateLocationEvent() error {
	// Location events require address and tickets
	if e.AddressID == nil {
		return ErrEventAddressRequired
	}

	if err := e.ValidateTicketRequirement(); err != nil {
		return err
	}

	// Date and time are required for location events
	if e.StartDate == nil || e.StartTime == nil {
		return ErrEventDateTimeRequired
	}

	return nil
}

func (e *Event) ValidateOnlineEvent() error {
	// Online events require tickets but not address
	if err := e.ValidateTicketRequirement(); err != nil {
		return err
	}

	// Date and time are required for online events
	if e.StartDate == nil || e.StartTime == nil {
		return ErrEventDateTimeRequired
	}

	return nil
}

func (e *Event) ValidateAnnouncementEvent() error {
	// Announcement events don't require date, tickets, or address
	// Only basic fields are required (already validated in ValidateForDraft)
	return nil
}

func (e *Event) ValidateTicketRequirement() error {
	// Either system tickets or external ticket URL is required
	if !e.HasSystemTickets && (e.TicketURL == nil || *e.TicketURL == "") {
		return ErrEventTicketRequired
	}
	return nil
}

// Helper methods
func (e *Event) IsPublic() bool {
	return e.Type == EventTypePublic
}

func (e *Event) IsPrivate() bool {
	return e.Type == EventTypePrivate
}

func (e *Event) IsLocationEvent() bool {
	return e.LocationType == EventLocationTypeLocation
}

func (e *Event) IsOnlineEvent() bool {
	return e.LocationType == EventLocationTypeOnline
}

func (e *Event) IsAnnouncementEvent() bool {
	return e.LocationType == EventLocationTypeAnnouncement
}

func (e *Event) IsDraft() bool {
	return e.Status == EventStatusDraft
}

func (e *Event) IsPending() bool {
	return e.Status == EventStatusPending
}

func (e *Event) IsPublished() bool {
	return e.Status == EventStatusPublished
}

func (e *Event) IsCancelled() bool {
	return e.Status == EventStatusCancelled
}

func (e *Event) CanBeEdited() bool {
	return e.Status == EventStatusDraft || e.Status == EventStatusRejected
}

func (e *Event) RequiresInvitations() bool {
	return e.IsPrivate()
}

func (e *Event) HasTickets() bool {
	return e.HasSystemTickets || (e.TicketURL != nil && *e.TicketURL != "")
}

func (e *Event) GetFullStartDateTime() *time.Time {
	if e.StartDate == nil || e.StartTime == nil {
		return nil
	}

	// Combine date and time
	year, month, day := e.StartDate.Date()
	hour, min, sec := e.StartTime.Clock()
	combined := time.Date(year, month, day, hour, min, sec, 0, e.StartDate.Location())
	return &combined
}

func (e *Event) GetFullEndDateTime() *time.Time {
	if e.EndDate == nil || e.EndTime == nil {
		return nil
	}

	// Combine date and time
	year, month, day := e.EndDate.Date()
	hour, min, sec := e.EndTime.Clock()
	combined := time.Date(year, month, day, hour, min, sec, 0, e.EndDate.Location())
	return &combined
}

// Event domain errors
var (
	ErrEventNameRequired            = NewDomainError("event name is required")
	ErrEventImageRequired           = NewDomainError("event image is required")
	ErrEventVideoRequired           = NewDomainError("event video is required")
	ErrEventAddressRequired         = NewDomainError("address is required for location events")
	ErrEventTicketRequired          = NewDomainError("ticket information or ticket URL is required")
	ErrEventDateTimeRequired        = NewDomainError("start date and time are required")
	ErrEventCategoriesRequired      = NewDomainError("at least one category is required")
	ErrEventInvalidLocationType     = NewDomainError("invalid location type")
	ErrEventInvalidStatusTransition = NewDomainError("invalid status transition")
	ErrEventAlreadyCancelled        = NewDomainError("event is already cancelled")
	ErrEventNotFound                = NewDomainError("event not found")
	ErrEventUnauthorized            = NewDomainError("unauthorized to access this event")
	ErrEventCannotBeEdited          = NewDomainError("event cannot be edited in current status")
)
