package dto

import "time"

// Creator Response
type CreatorResponse struct {
	ID               int                `json:"id"`
	UserID           int                `json:"user_id"`
	WeeztixToken     *string            `json:"weeztix_token"`
	CompanyName      string             `json:"company_name"`
	Address          string             `json:"address"`
	EstimatedTickets int                `json:"estimated_tickets"`
	EstimatedEvents  int                `json:"estimated_events"`
	Industries       []IndustryResponse `json:"industries"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// Create Creator Request
type CreateCreatorRequest struct {
	WeeztixToken     *string `json:"weeztix_token,omitempty"`
	CompanyName      string  `json:"company_name" validate:"required,min=2,max=200"`
	Address          string  `json:"address" validate:"required,min=5,max=500"`
	EstimatedTickets int     `json:"estimated_tickets" validate:"required,min=1"`
	EstimatedEvents  int     `json:"estimated_events" validate:"required,min=1"`
	IndustryIDs      []int   `json:"industry_ids" validate:"required,min=1,dive,min=1"`
}

// Update Creator Request
type UpdateCreatorRequest struct {
	CompanyName      string `json:"company_name" validate:"omitempty,min=2,max=200"`
	Address          string `json:"address" validate:"omitempty,min=5,max=500"`
	EstimatedTickets int    `json:"estimated_tickets" validate:"omitempty,min=1"`
	EstimatedEvents  int    `json:"estimated_events" validate:"omitempty,min=1"`
	IndustryIDs      []int  `json:"industry_ids" validate:"omitempty,min=1,dive,min=1"`
}

// Set Weeztix Token Request
type SetWeeztixTokenRequest struct {
	WeeztixToken string `json:"weeztix_token" validate:"required,min=10,max=255"`
}

// Creator Profile Response (includes user data)
type CreatorProfileResponse struct {
	User    UserResponse    `json:"user"`
	Creator CreatorResponse `json:"creator"`
}

// Creator List Request
type CreatorListRequest struct {
	Page       int `json:"page" validate:"omitempty,min=1"`
	PageSize   int `json:"page_size" validate:"omitempty,min=1,max=100"`
	IndustryID int `json:"industry_id" validate:"omitempty,min=1"`
}

// Creator List Response
type CreatorListResponse struct {
	Creators   []CreatorResponse `json:"creators"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
