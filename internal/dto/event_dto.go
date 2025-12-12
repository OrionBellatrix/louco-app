package dto

import (
	"time"

	"github.com/louco-event/internal/domain"
)

// Event creation and update requests
type CreateEventRequest struct {
	Name             string                   `json:"name" validate:"required,min=3,max=200"`
	Description      *string                  `json:"description" validate:"omitempty,max=2000"`
	ImageID          *int                     `json:"image_id" validate:"omitempty,gt=0"`
	VideoID          *int                     `json:"video_id" validate:"omitempty,gt=0"`
	Type             domain.EventType         `json:"type" validate:"required,oneof=public private"`
	LocationType     domain.EventLocationType `json:"location_type" validate:"required,oneof=location online announcement"`
	StartDate        *string                  `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	StartTime        *string                  `json:"start_time" validate:"omitempty,datetime=15:04"`
	EndDate          *string                  `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	EndTime          *string                  `json:"end_time" validate:"omitempty,datetime=15:04"`
	AddressID        *int                     `json:"address_id" validate:"omitempty,gt=0"`
	OnlineEventURL   *string                  `json:"online_event_url" validate:"omitempty,url,max=500"`
	OnlineEventType  *string                  `json:"online_event_type" validate:"omitempty,max=50"`
	TicketURL        *string                  `json:"ticket_url" validate:"omitempty,url,max=500"`
	HasSystemTickets bool                     `json:"has_system_tickets"`
	AdditionalInfo   *string                  `json:"additional_info" validate:"omitempty,max=2000"`
	CategoryIDs      []int                    `json:"category_ids" validate:"omitempty,dive,gt=0"`
}

type UpdateEventRequest struct {
	Name             *string                   `json:"name" validate:"omitempty,min=3,max=200"`
	Description      *string                   `json:"description" validate:"omitempty,max=2000"`
	ImageID          *int                      `json:"image_id" validate:"omitempty,gt=0"`
	VideoID          *int                      `json:"video_id" validate:"omitempty,gt=0"`
	Type             *domain.EventType         `json:"type" validate:"omitempty,oneof=public private"`
	LocationType     *domain.EventLocationType `json:"location_type" validate:"omitempty,oneof=location online announcement"`
	StartDate        *string                   `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	StartTime        *string                   `json:"start_time" validate:"omitempty,datetime=15:04"`
	EndDate          *string                   `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	EndTime          *string                   `json:"end_time" validate:"omitempty,datetime=15:04"`
	AddressID        *int                      `json:"address_id" validate:"omitempty,gt=0"`
	OnlineEventURL   *string                   `json:"online_event_url" validate:"omitempty,url,max=500"`
	OnlineEventType  *string                   `json:"online_event_type" validate:"omitempty,max=50"`
	TicketURL        *string                   `json:"ticket_url" validate:"omitempty,url,max=500"`
	HasSystemTickets *bool                     `json:"has_system_tickets"`
	AdditionalInfo   *string                   `json:"additional_info" validate:"omitempty,max=2000"`
	CategoryIDs      []int                     `json:"category_ids" validate:"omitempty,dive,gt=0"`
}

type UpdateEventStatusRequest struct {
	Status domain.EventStatus `json:"status" validate:"required,oneof=draft pending rejected stopped cancelled published"`
}

// Event response DTOs
type EventResponse struct {
	ID               int                      `json:"id"`
	CreatorID        int                      `json:"creator_id"`
	Name             string                   `json:"name"`
	Description      *string                  `json:"description"`
	ImageID          *int                     `json:"image_id"`
	VideoID          *int                     `json:"video_id"`
	Type             domain.EventType         `json:"type"`
	LocationType     domain.EventLocationType `json:"location_type"`
	Status           domain.EventStatus       `json:"status"`
	StartDate        *string                  `json:"start_date"`
	StartTime        *string                  `json:"start_time"`
	EndDate          *string                  `json:"end_date"`
	EndTime          *string                  `json:"end_time"`
	AddressID        *int                     `json:"address_id"`
	OnlineEventURL   *string                  `json:"online_event_url"`
	OnlineEventType  *string                  `json:"online_event_type"`
	TicketURL        *string                  `json:"ticket_url"`
	HasSystemTickets bool                     `json:"has_system_tickets"`
	AdditionalInfo   *string                  `json:"additional_info"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`

	// Relations
	Creator     *CreatorResponse     `json:"creator,omitempty"`
	Image       *MediaResponse       `json:"image,omitempty"`
	Video       *MediaResponse       `json:"video,omitempty"`
	Address     *AddressResponse     `json:"address,omitempty"`
	Categories  []CategoryResponse   `json:"categories,omitempty"`
	Tickets     []TicketResponse     `json:"tickets,omitempty"`
	Invitations []InvitationResponse `json:"invitations,omitempty"`
}

type EventListResponse struct {
	ID               int                      `json:"id"`
	CreatorID        int                      `json:"creator_id"`
	Name             string                   `json:"name"`
	Description      *string                  `json:"description"`
	Type             domain.EventType         `json:"type"`
	LocationType     domain.EventLocationType `json:"location_type"`
	Status           domain.EventStatus       `json:"status"`
	StartDate        *string                  `json:"start_date"`
	StartTime        *string                  `json:"start_time"`
	HasSystemTickets bool                     `json:"has_system_tickets"`
	CreatedAt        time.Time                `json:"created_at"`

	// Basic relations for list view
	Creator     *CreatorBasicResponse `json:"creator,omitempty"`
	Image       *MediaResponse        `json:"image,omitempty"`
	Address     *AddressBasicResponse `json:"address,omitempty"`
	Categories  []CategoryResponse    `json:"categories,omitempty"`
	TicketCount int                   `json:"ticket_count"`
}

// Address DTOs
type AddressResponse struct {
	ID          int       `json:"id"`
	PlaceID     string    `json:"place_id"`
	FullAddress string    `json:"full_address"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	District    *string   `json:"district"`
	Street      *string   `json:"street"`
	PostalCode  *string   `json:"postal_code"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	DoorNumber  *string   `json:"door_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AddressBasicResponse struct {
	ID          int     `json:"id"`
	FullAddress string  `json:"full_address"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	DoorNumber  *string `json:"door_number"`
}

type CreateAddressRequest struct {
	PlaceID     string  `json:"place_id" validate:"required,max=255"`
	FullAddress string  `json:"full_address" validate:"required,max=1000"`
	Country     string  `json:"country" validate:"required,max=100"`
	City        string  `json:"city" validate:"required,max=100"`
	District    *string `json:"district" validate:"omitempty,max=100"`
	Street      *string `json:"street" validate:"omitempty,max=200"`
	PostalCode  *string `json:"postal_code" validate:"omitempty,max=20"`
	Latitude    float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude   float64 `json:"longitude" validate:"required,min=-180,max=180"`
	DoorNumber  *string `json:"door_number" validate:"omitempty,max=50"`
}

// Ticket DTOs
type TicketResponse struct {
	ID            int       `json:"id"`
	EventID       int       `json:"event_id"`
	Title         string    `json:"title"`
	Price         float64   `json:"price"`
	TotalQuantity int       `json:"total_quantity"`
	SoldQuantity  int       `json:"sold_quantity"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateTicketRequest struct {
	Title         string  `json:"title" validate:"required,min=3,max=200"`
	Price         float64 `json:"price" validate:"required,min=0"`
	TotalQuantity int     `json:"total_quantity" validate:"required,min=1"`
}

type UpdateTicketRequest struct {
	Title         *string  `json:"title" validate:"omitempty,min=3,max=200"`
	Price         *float64 `json:"price" validate:"omitempty,min=0"`
	TotalQuantity *int     `json:"total_quantity" validate:"omitempty,min=1"`
}

// Invitation DTOs
type InvitationResponse struct {
	ID            int                     `json:"id"`
	EventID       int                     `json:"event_id"`
	InvitedUserID *int                    `json:"invited_user_id"`
	InvitedEmail  string                  `json:"invited_email"`
	Status        domain.InvitationStatus `json:"status"`
	InvitedAt     time.Time               `json:"invited_at"`
	RespondedAt   *time.Time              `json:"responded_at"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`

	// Relations
	InvitedUser *UserBasicResponse `json:"invited_user,omitempty"`
}

type CreateInvitationRequest struct {
	InvitedEmail  string `json:"invited_email" validate:"required,email,max=255"`
	InvitedUserID *int   `json:"invited_user_id" validate:"omitempty,gt=0"`
}

type UpdateInvitationStatusRequest struct {
	Status domain.InvitationStatus `json:"status" validate:"required,oneof=pending approved rejected"`
}

type BulkCreateInvitationRequest struct {
	Invitations []CreateInvitationRequest `json:"invitations" validate:"required,min=1,max=100,dive"`
}

// Filter and search DTOs
type EventFilterRequest struct {
	Type         *domain.EventType         `json:"type" validate:"omitempty,oneof=public private"`
	LocationType *domain.EventLocationType `json:"location_type" validate:"omitempty,oneof=location online announcement"`
	Status       *domain.EventStatus       `json:"status" validate:"omitempty,oneof=draft pending rejected stopped cancelled published"`
	CategoryIDs  []int                     `json:"category_ids" validate:"omitempty,dive,gt=0"`
	City         *string                   `json:"city" validate:"omitempty,max=100"`
	Country      *string                   `json:"country" validate:"omitempty,max=100"`
	StartDate    *string                   `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate      *string                   `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	MinPrice     *float64                  `json:"min_price" validate:"omitempty,min=0"`
	MaxPrice     *float64                  `json:"max_price" validate:"omitempty,min=0"`
	HasTickets   *bool                     `json:"has_tickets"`
	CreatorID    *int                      `json:"creator_id" validate:"omitempty,gt=0"`
	Query        *string                   `json:"query" validate:"omitempty,max=200"`
}

type EventSearchRequest struct {
	Query        string                    `json:"query" validate:"required,min=2,max=200"`
	Type         *domain.EventType         `json:"type" validate:"omitempty,oneof=public private"`
	LocationType *domain.EventLocationType `json:"location_type" validate:"omitempty,oneof=location online announcement"`
	CategoryIDs  []int                     `json:"category_ids" validate:"omitempty,dive,gt=0"`
	City         *string                   `json:"city" validate:"omitempty,max=100"`
}

// Statistics DTOs
type EventStatsResponse struct {
	TotalEvents        int64 `json:"total_events"`
	DraftEvents        int64 `json:"draft_events"`
	PendingEvents      int64 `json:"pending_events"`
	PublishedEvents    int64 `json:"published_events"`
	RejectedEvents     int64 `json:"rejected_events"`
	CancelledEvents    int64 `json:"cancelled_events"`
	PublicEvents       int64 `json:"public_events"`
	PrivateEvents      int64 `json:"private_events"`
	LocationEvents     int64 `json:"location_events"`
	OnlineEvents       int64 `json:"online_events"`
	AnnouncementEvents int64 `json:"announcement_events"`
}

type SystemEventStatsResponse struct {
	EventStatsResponse
	TotalCreators    int64   `json:"total_creators"`
	ActiveCreators   int64   `json:"active_creators"`
	TotalTicketsSold int64   `json:"total_tickets_sold"`
	TotalRevenue     float64 `json:"total_revenue"`
}

// Helper response DTOs
type CreatorBasicResponse struct {
	ID          int                `json:"id"`
	UserID      int                `json:"user_id"`
	CompanyName string             `json:"company_name"`
	User        *UserBasicResponse `json:"user,omitempty"`
}

type UserBasicResponse struct {
	ID       int     `json:"id"`
	FullName string  `json:"full_name"`
	Username *string `json:"username"`
}

// Utility functions for DTO conversion
func (r *CreateEventRequest) ToStartDateTime() (*time.Time, error) {
	if r.StartDate == nil || r.StartTime == nil {
		return nil, nil
	}

	dateStr := *r.StartDate + " " + *r.StartTime
	t, err := time.Parse("2006-01-02 15:04", dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *CreateEventRequest) ToEndDateTime() (*time.Time, error) {
	if r.EndDate == nil || r.EndTime == nil {
		return nil, nil
	}

	dateStr := *r.EndDate + " " + *r.EndTime
	t, err := time.Parse("2006-01-02 15:04", dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func EventToResponse(event *domain.Event) *EventResponse {
	response := &EventResponse{
		ID:               event.ID,
		CreatorID:        event.CreatorID,
		Name:             event.Name,
		Description:      event.Description,
		ImageID:          event.ImageID,
		VideoID:          event.VideoID,
		Type:             event.Type,
		LocationType:     event.LocationType,
		Status:           event.Status,
		AddressID:        event.AddressID,
		OnlineEventURL:   event.OnlineEventURL,
		OnlineEventType:  event.OnlineEventType,
		TicketURL:        event.TicketURL,
		HasSystemTickets: event.HasSystemTickets,
		AdditionalInfo:   event.AdditionalInfo,
		CreatedAt:        event.CreatedAt,
		UpdatedAt:        event.UpdatedAt,
	}

	// Format dates
	if event.StartDate != nil {
		startDate := event.StartDate.Format("2006-01-02")
		response.StartDate = &startDate
	}
	if event.StartTime != nil {
		startTime := event.StartTime.Format("15:04")
		response.StartTime = &startTime
	}
	if event.EndDate != nil {
		endDate := event.EndDate.Format("2006-01-02")
		response.EndDate = &endDate
	}
	if event.EndTime != nil {
		endTime := event.EndTime.Format("15:04")
		response.EndTime = &endTime
	}

	// Add relations
	if event.Creator.ID != 0 {
		response.Creator = CreatorToResponse(&event.Creator)
	}

	if event.Image != nil {
		response.Image = MediaToResponse(event.Image)
	}

	if event.Video != nil {
		response.Video = MediaToResponse(event.Video)
	}

	if event.Address != nil {
		response.Address = AddressToResponse(event.Address)
	}

	// Add categories
	if len(event.Categories) > 0 {
		response.Categories = make([]CategoryResponse, len(event.Categories))
		for i, category := range event.Categories {
			response.Categories[i] = CategoryResponse{
				ID:       category.ID,
				Name:     category.Name,
				Slug:     category.Slug,
				ParentID: category.ParentID,
				Depth:    category.Depth,
			}
		}
	}

	// Add tickets
	if len(event.Tickets) > 0 {
		response.Tickets = make([]TicketResponse, len(event.Tickets))
		for i, ticket := range event.Tickets {
			response.Tickets[i] = TicketToResponse(&ticket)
		}
	}

	// Add invitations
	if len(event.Invitations) > 0 {
		response.Invitations = make([]InvitationResponse, len(event.Invitations))
		for i, invitation := range event.Invitations {
			response.Invitations[i] = InvitationToResponse(&invitation)
		}
	}

	return response
}

// Additional DTOs for filtering and statistics
type AddressFilterRequest struct {
	City         *string  `json:"city" validate:"omitempty,max=100"`
	Country      *string  `json:"country" validate:"omitempty,max=100"`
	District     *string  `json:"district" validate:"omitempty,max=100"`
	MinLatitude  *float64 `json:"min_latitude" validate:"omitempty,min=-90,max=90"`
	MaxLatitude  *float64 `json:"max_latitude" validate:"omitempty,min=-90,max=90"`
	MinLongitude *float64 `json:"min_longitude" validate:"omitempty,min=-180,max=180"`
	MaxLongitude *float64 `json:"max_longitude" validate:"omitempty,min=-180,max=180"`
	Query        *string  `json:"query" validate:"omitempty,max=200"`
}

type TicketFilterRequest struct {
	EventID   *int     `json:"event_id" validate:"omitempty,gt=0"`
	MinPrice  *float64 `json:"min_price" validate:"omitempty,min=0"`
	MaxPrice  *float64 `json:"max_price" validate:"omitempty,min=0"`
	IsActive  *bool    `json:"is_active"`
	IsSoldOut *bool    `json:"is_sold_out"`
	IsFree    *bool    `json:"is_free"`
	Query     *string  `json:"query" validate:"omitempty,max=200"`
}

type InvitationFilterRequest struct {
	EventID       *int                     `json:"event_id" validate:"omitempty,gt=0"`
	InvitedUserID *int                     `json:"invited_user_id" validate:"omitempty,gt=0"`
	Status        *domain.InvitationStatus `json:"status" validate:"omitempty,oneof=pending approved rejected"`
	Email         *string                  `json:"email" validate:"omitempty,email,max=255"`
	IsExpired     *bool                    `json:"is_expired"`
	HasResponded  *bool                    `json:"has_responded"`
	Query         *string                  `json:"query" validate:"omitempty,max=200"`
}

// Statistics DTOs
type TicketSalesStatsResponse struct {
	EventID          int     `json:"event_id"`
	TotalTickets     int     `json:"total_tickets"`
	SoldTickets      int     `json:"sold_tickets"`
	AvailableTickets int     `json:"available_tickets"`
	TotalRevenue     float64 `json:"total_revenue"`
	AveragePrice     float64 `json:"average_price"`
	SoldPercentage   float64 `json:"sold_percentage"`
}

type TicketTypeStatsResponse struct {
	TicketID          int     `json:"ticket_id"`
	Title             string  `json:"title"`
	Price             float64 `json:"price"`
	TotalQuantity     int     `json:"total_quantity"`
	SoldQuantity      int     `json:"sold_quantity"`
	AvailableQuantity int     `json:"available_quantity"`
	Revenue           float64 `json:"revenue"`
	SoldPercentage    float64 `json:"sold_percentage"`
}

type InvitationStatsResponse struct {
	EventID                 int     `json:"event_id"`
	TotalInvitations        int     `json:"total_invitations"`
	PendingInvitations      int     `json:"pending_invitations"`
	ApprovedInvitations     int     `json:"approved_invitations"`
	RejectedInvitations     int     `json:"rejected_invitations"`
	SystemUserInvitations   int     `json:"system_user_invitations"`
	ExternalUserInvitations int     `json:"external_user_invitations"`
	ResponseRate            float64 `json:"response_rate"`
	ApprovalRate            float64 `json:"approval_rate"`
}

type UserInvitationStatsResponse struct {
	UserID              int     `json:"user_id"`
	TotalInvitations    int     `json:"total_invitations"`
	PendingInvitations  int     `json:"pending_invitations"`
	ApprovedInvitations int     `json:"approved_invitations"`
	RejectedInvitations int     `json:"rejected_invitations"`
	ResponseRate        float64 `json:"response_rate"`
	ApprovalRate        float64 `json:"approval_rate"`
}

type SystemInvitationStatsResponse struct {
	TotalInvitations        int     `json:"total_invitations"`
	PendingInvitations      int     `json:"pending_invitations"`
	ApprovedInvitations     int     `json:"approved_invitations"`
	RejectedInvitations     int     `json:"rejected_invitations"`
	SystemUserInvitations   int     `json:"system_user_invitations"`
	ExternalUserInvitations int     `json:"external_user_invitations"`
	ExpiredInvitations      int     `json:"expired_invitations"`
	ResponseRate            float64 `json:"response_rate"`
	ApprovalRate            float64 `json:"approval_rate"`
}

func EventToListResponse(event *domain.Event) *EventListResponse {
	response := &EventListResponse{
		ID:               event.ID,
		CreatorID:        event.CreatorID,
		Name:             event.Name,
		Description:      event.Description,
		Type:             event.Type,
		LocationType:     event.LocationType,
		Status:           event.Status,
		HasSystemTickets: event.HasSystemTickets,
		CreatedAt:        event.CreatedAt,
		TicketCount:      len(event.Tickets),
	}

	// Format dates
	if event.StartDate != nil {
		startDate := event.StartDate.Format("2006-01-02")
		response.StartDate = &startDate
	}
	if event.StartTime != nil {
		startTime := event.StartTime.Format("15:04")
		response.StartTime = &startTime
	}

	return response
}

// Location and Address related DTOs
type LocationRequest struct {
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
	RadiusKm  int     `json:"radius_km" validate:"omitempty,min=1,max=100"`
}

// Address Statistics DTOs
type AddressStatsResponse struct {
	TotalAddresses   int64   `json:"total_addresses"`
	TotalCities      int64   `json:"total_cities"`
	TotalCountries   int64   `json:"total_countries"`
	MostUsedCity     string  `json:"most_used_city"`
	MostUsedCountry  string  `json:"most_used_country"`
	AverageLatitude  float64 `json:"average_latitude"`
	AverageLongitude float64 `json:"average_longitude"`
}

type CityStatsResponse struct {
	City         string `json:"city"`
	Country      string `json:"country"`
	AddressCount int64  `json:"address_count"`
	EventCount   int64  `json:"event_count"`
}

type CountryStatsResponse struct {
	Country      string `json:"country"`
	AddressCount int64  `json:"address_count"`
	EventCount   int64  `json:"event_count"`
	CityCount    int64  `json:"city_count"`
}

// Helper conversion functions
func CreatorToResponse(creator *domain.Creator) *CreatorResponse {
	if creator == nil || creator.ID == 0 {
		return nil
	}

	response := &CreatorResponse{
		ID:               creator.ID,
		UserID:           creator.UserID,
		WeeztixToken:     creator.WeeztixToken,
		CompanyName:      creator.CompanyName,
		Address:          creator.Address,
		EstimatedTickets: creator.EstimatedTickets,
		EstimatedEvents:  creator.EstimatedEvents,
		CreatedAt:        creator.CreatedAt,
		UpdatedAt:        creator.UpdatedAt,
	}

	// Add industries
	if len(creator.Industries) > 0 {
		response.Industries = make([]IndustryResponse, len(creator.Industries))
		for i, industry := range creator.Industries {
			response.Industries[i] = IndustryResponse{
				ID:   industry.ID,
				Name: industry.Name,
				Slug: industry.Slug,
			}
		}
	}

	return response
}

func MediaToResponse(media *domain.Media) *MediaResponse {
	if media == nil {
		return nil
	}

	return &MediaResponse{
		ID:           media.ID,
		UserID:       media.UserID,
		OriginalName: media.OriginalName,
		FileName:     media.FileName,
		FileURL:      media.FileURL,
		MediaType:    string(media.MediaType),
		MimeType:     media.MimeType,
		FileSize:     media.FileSize,
		Width:        media.Width,
		Height:       media.Height,
		Duration:     media.Duration,
		IsConverted:  media.IsConverted,
		CreatedAt:    media.CreatedAt,
		UpdatedAt:    media.UpdatedAt,
	}
}

func AddressToResponse(address *domain.Address) *AddressResponse {
	if address == nil {
		return nil
	}

	return &AddressResponse{
		ID:          address.ID,
		PlaceID:     address.PlaceID,
		FullAddress: address.FullAddress,
		Country:     address.Country,
		City:        address.City,
		District:    address.District,
		Street:      address.Street,
		PostalCode:  address.PostalCode,
		Latitude:    address.Latitude,
		Longitude:   address.Longitude,
		DoorNumber:  address.DoorNumber,
		CreatedAt:   address.CreatedAt,
		UpdatedAt:   address.UpdatedAt,
	}
}

func TicketToResponse(ticket *domain.Ticket) TicketResponse {
	return TicketResponse{
		ID:            ticket.ID,
		EventID:       ticket.EventID,
		Title:         ticket.Title,
		Price:         ticket.Price,
		TotalQuantity: ticket.TotalQuantity,
		SoldQuantity:  ticket.SoldQuantity,
		IsActive:      ticket.IsActive,
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
	}
}

func InvitationToResponse(invitation *domain.Invitation) InvitationResponse {
	response := InvitationResponse{
		ID:            invitation.ID,
		EventID:       invitation.EventID,
		InvitedUserID: invitation.InvitedUserID,
		InvitedEmail:  invitation.InvitedEmail,
		Status:        invitation.Status,
		InvitedAt:     invitation.InvitedAt,
		RespondedAt:   invitation.RespondedAt,
		CreatedAt:     invitation.CreatedAt,
		UpdatedAt:     invitation.UpdatedAt,
	}

	// Add invited user if exists
	if invitation.InvitedUser != nil {
		response.InvitedUser = &UserBasicResponse{
			ID:       invitation.InvitedUser.ID,
			FullName: invitation.InvitedUser.FullName,
			Username: invitation.InvitedUser.Username,
		}
	}

	return response
}
