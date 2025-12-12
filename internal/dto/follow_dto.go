package dto

import "time"

// FollowRequest represents a follow request (target user to follow)
type FollowRequest struct {
	UserID int `json:"user_id" validate:"required,min=1" example:"2"`
}

// FollowStatusRequest represents a request to check follow status
type FollowStatusRequest struct {
	UserID int `uri:"user_id" validate:"required,min=1" example:"2"`
}

// FollowStatusResponse represents the follow status between two users
type FollowStatusResponse struct {
	IsFollowing bool `json:"is_following" example:"true"`
	IsFollower  bool `json:"is_follower" example:"false"`
}

// UserFollowInfo represents basic user info for follow lists
type UserFollowInfo struct {
	ID             int            `json:"id" example:"1"`
	FirstName      string         `json:"first_name" example:"John"`
	LastName       string         `json:"last_name" example:"Doe"`
	Username       string         `json:"username" example:"johndoe"`
	ProfilePicture *MediaResponse `json:"profile_picture,omitempty"`
	FollowersCount int            `json:"followers_count" example:"150"`
	FollowingCount int            `json:"following_count" example:"75"`
}

// FollowResponse represents a follow relationship with user info
type FollowResponse struct {
	ID          int             `json:"id" example:"1"`
	User        *UserFollowInfo `json:"user"`
	FollowedAt  time.Time       `json:"followed_at" example:"2024-01-15T10:30:00Z"`
	IsFollowing bool            `json:"is_following,omitempty" example:"true"`
}

// FollowersResponse represents the response for followers list
type FollowersResponse struct {
	Followers  []*FollowResponse `json:"followers"`
	Total      int               `json:"total" example:"150"`
	Page       int               `json:"page" example:"1"`
	Limit      int               `json:"limit" example:"20"`
	TotalPages int               `json:"total_pages" example:"8"`
}

// FollowingResponse represents the response for following list
type FollowingResponse struct {
	Following  []*FollowResponse `json:"following"`
	Total      int               `json:"total" example:"75"`
	Page       int               `json:"page" example:"1"`
	Limit      int               `json:"limit" example:"20"`
	TotalPages int               `json:"total_pages" example:"4"`
}

// MutualFollowsResponse represents the response for mutual follows
type MutualFollowsResponse struct {
	MutualFollows []*FollowResponse `json:"mutual_follows"`
	Total         int               `json:"total" example:"25"`
	Page          int               `json:"page" example:"1"`
	Limit         int               `json:"limit" example:"20"`
	TotalPages    int               `json:"total_pages" example:"2"`
}

// FollowCountsResponse represents follow counts for a user
type FollowCountsResponse struct {
	FollowersCount int `json:"followers_count" example:"150"`
	FollowingCount int `json:"following_count" example:"75"`
}

// PaginationQuery represents pagination parameters for follow lists
type PaginationQuery struct {
	Page  int `form:"page" validate:"min=1" example:"1"`
	Limit int `form:"limit" validate:"min=1,max=100" example:"20"`
}

// GetDefaultPagination returns default pagination values
func (p *PaginationQuery) GetDefaultPagination() (int, int) {
	page := p.Page
	if page < 1 {
		page = 1
	}

	limit := p.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return page, limit
}

// GetOffset calculates the offset for database queries
func (p *PaginationQuery) GetOffset() int {
	page, limit := p.GetDefaultPagination()
	return (page - 1) * limit
}

// CalculateTotalPages calculates total pages based on total count and limit
func CalculateTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
