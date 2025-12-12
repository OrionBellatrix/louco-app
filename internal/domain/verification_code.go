package domain

import (
	"crypto/rand"
	"fmt"
	"time"
)

type VerificationType string

const (
	VerificationTypeEmail VerificationType = "email"
	VerificationTypePhone VerificationType = "phone"
)

type VerificationCode struct {
	ID         uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	Identifier string           `json:"identifier" gorm:"not null;index"` // email or phone
	Code       string           `json:"code" gorm:"not null"`
	Type       VerificationType `json:"type" gorm:"not null;index"`
	Attempts   int              `json:"attempts" gorm:"default:0"`
	UsedAt     *time.Time       `json:"used_at" gorm:"index"`
	ExpiresAt  time.Time        `json:"expires_at" gorm:"not null;index"`
	CreatedAt  time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

func (VerificationCode) TableName() string {
	return "verification_codes"
}

// NewVerificationCode creates a new verification code instance
func NewVerificationCode(identifier, code string, codeType VerificationType) *VerificationCode {
	return &VerificationCode{
		Identifier: identifier,
		Code:       code,
		Type:       codeType,
		Attempts:   0,
		ExpiresAt:  time.Now().Add(10 * time.Minute), // 10 minutes expiry
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// NewEmailVerificationCode creates a new email verification code
func NewEmailVerificationCode(identifier string) *VerificationCode {
	code := generateSixDigitCode()
	return NewVerificationCode(identifier, code, VerificationTypeEmail)
}

// NewPhoneVerificationCode creates a new phone verification code
func NewPhoneVerificationCode(identifier string) *VerificationCode {
	code := generateSixDigitCode()
	return NewVerificationCode(identifier, code, VerificationTypePhone)
}

// IsExpired checks if the verification code has expired
func (vc *VerificationCode) IsExpired() bool {
	return time.Now().After(vc.ExpiresAt)
}

// IsUsed checks if the code has been used
func (vc *VerificationCode) IsUsed() bool {
	return vc.UsedAt != nil
}

// CanAttempt checks if more attempts are allowed (max 5 attempts)
func (vc *VerificationCode) CanAttempt() bool {
	return vc.Attempts < 5
}

// IncrementAttempts increments the attempt counter
func (vc *VerificationCode) IncrementAttempts() {
	vc.Attempts++
	vc.UpdatedAt = time.Now()
}

// MarkAsUsed marks the code as used
func (vc *VerificationCode) MarkAsUsed() {
	now := time.Now()
	vc.UsedAt = &now
	vc.UpdatedAt = now
}

// IsValid checks if the code is valid for verification
func (vc *VerificationCode) IsValid() bool {
	return !vc.IsExpired() && !vc.IsUsed() && vc.CanAttempt()
}

// VerifyCode verifies the provided code against the stored code
func (vc *VerificationCode) VerifyCode(inputCode string) bool {
	if !vc.IsValid() {
		return false
	}

	vc.IncrementAttempts()

	if vc.Code == inputCode {
		vc.MarkAsUsed()
		return true
	}

	return false
}

// generateSixDigitCode generates a random 6-digit verification code
func generateSixDigitCode() string {
	// Generate a random number between 100000 and 999999
	bytes := make([]byte, 3)
	rand.Read(bytes)

	// Convert to 6-digit number
	num := int(bytes[0])<<16 | int(bytes[1])<<8 | int(bytes[2])
	code := (num % 900000) + 100000

	return fmt.Sprintf("%06d", code)
}

// Domain errors for verification
var (
	ErrVerificationCodeExpired  = NewDomainError("verification code has expired")
	ErrVerificationCodeUsed     = NewDomainError("verification code has already been used")
	ErrVerificationCodeInvalid  = NewDomainError("invalid verification code")
	ErrVerificationMaxAttempts  = NewDomainError("maximum verification attempts exceeded")
	ErrVerificationCodeNotFound = NewDomainError("verification code not found")
)
