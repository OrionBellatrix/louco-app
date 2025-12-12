package dto

import "github.com/louco-event/internal/domain"

// SendVerificationRequest represents the request to send verification code
type SendVerificationRequest struct {
	Identifier string `json:"identifier" validate:"required" example:"user@example.com"`
}

// VerifyCodeRequest represents the request to verify a code
type VerifyCodeRequest struct {
	Identifier string `json:"identifier" validate:"required" example:"user@example.com"`
	Code       string `json:"code" validate:"required,len=6,numeric" example:"123456"`
}

// ResendVerificationRequest represents the request to resend verification code
type ResendVerificationRequest struct {
	Identifier string `json:"identifier" validate:"required" example:"user@example.com"`
}

// VerificationStatusResponse represents the verification status
type VerificationStatusResponse struct {
	EmailVerified bool `json:"email_verified" example:"true"`
	PhoneVerified bool `json:"phone_verified" example:"false"`
}

// Helper function to convert string to VerificationType
func StringToVerificationType(s string) domain.VerificationType {
	switch s {
	case "email":
		return domain.VerificationTypeEmail
	case "phone":
		return domain.VerificationTypePhone
	default:
		return domain.VerificationTypeEmail // default fallback
	}
}
