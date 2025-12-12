package dto

import "time"

// Register Step 1 - Account Creation
type RegisterStep1Request struct {
	Identifier string `json:"identifier" validate:"required"` // email or phone
	Password   string `json:"password" validate:"required,min=8"`
	UserType   string `json:"user_type" validate:"required,oneof=user creator"`
}

type RegisterStep1Response struct {
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}

// Register Step 3 - Username
type SetUsernameRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum_underscore_dot"`
}

type CheckUsernameRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum_underscore_dot"`
}

type CheckUsernameResponse struct {
	Exists bool `json:"exists"`
}

// Register Step 4 - Profile Details
type RegisterStep4Request struct {
	FullName  string     `json:"full_name" validate:"required,min=2,max=100"`
	Email     *string    `json:"email" validate:"omitempty,email"`
	Phone     *string    `json:"phone" validate:"omitempty,e164"`
	Biography *string    `json:"biography" validate:"omitempty,max=500"`
	BirthDate *time.Time `json:"birth_date" validate:"omitempty"`
	// Creator-specific fields
	Address          *string `json:"address" validate:"omitempty,max=500"`
	CompanyName      *string `json:"company_name" validate:"omitempty,max=200"`
	EstimatedTickets *int    `json:"estimated_tickets" validate:"omitempty,min=1"`
	EstimatedEvents  *int    `json:"estimated_events" validate:"omitempty,min=1"`
	IndustryIDs      []int   `json:"industry_ids" validate:"omitempty,dive,min=1"`
}

// Login
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // email, phone, or username
	Password   string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// Social Login
type SocialLoginRequest struct {
	Provider string  `json:"provider" validate:"required,oneof=apple google"`
	SocialID string  `json:"social_id" validate:"required"`
	Email    *string `json:"email" validate:"omitempty,email"`
	FullName *string `json:"full_name" validate:"omitempty"`
	UserType string  `json:"user_type" validate:"required,oneof=user creator"`
}

// JWT Claims
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	UserType string `json:"user_type"`
}

// Password Reset
type ForgotPasswordRequest struct {
	Identifier string `json:"identifier" validate:"required"` // email or phone
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// Change Password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}
