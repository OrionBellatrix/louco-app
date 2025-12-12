package dto

import "time"

type UserResponse struct {
	ID              int            `json:"id"`
	FullName        string         `json:"full_name"`
	Username        *string        `json:"username"`
	Email           *string        `json:"email"`
	Phone           *string        `json:"phone"`
	UserType        string         `json:"user_type"`
	Biography       *string        `json:"biography"`
	BirthDate       *time.Time     `json:"birth_date"`
	ProfilePicID    *int           `json:"profile_pic_id"`
	ProfilePic      *MediaResponse `json:"profile_pic,omitempty"`
	CoverPicID      *int           `json:"cover_pic_id"`
	CoverPic        *MediaResponse `json:"cover_pic,omitempty"`
	FollowersCount  int            `json:"followers_count"`
	FollowingCount  int            `json:"following_count"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time     `json:"phone_verified_at"`
	IsActive        bool           `json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type UpdateProfileRequest struct {
	FullName  string     `json:"full_name" validate:"omitempty,min=2,max=100"`
	Biography *string    `json:"biography" validate:"omitempty,max=500"`
	BirthDate *time.Time `json:"birth_date" validate:"omitempty"`
}

type UpdateContactRequest struct {
	Email *string `json:"email" validate:"omitempty,email"`
	Phone *string `json:"phone" validate:"omitempty,e164"`
}

type SetProfilePicRequest struct {
	MediaID int `json:"media_id" validate:"required,min=1"`
}

type SetCoverPicRequest struct {
	MediaID int `json:"media_id" validate:"required,min=1"`
}

type UserListRequest struct {
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	UserType string `json:"user_type" validate:"omitempty,oneof=user creator"`
}

type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

type UserProfileResponse struct {
	User UserResponse `json:"user"`
}
