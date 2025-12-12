package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/repository"
	"github.com/louco-event/pkg/email"
	"github.com/louco-event/pkg/twilio"
	"github.com/louco-event/pkg/utils"
	"gorm.io/gorm"
)

type VerificationService interface {
	SendEmailVerification(ctx context.Context, userID int, emailAddr, language string) error
	SendPhoneVerification(ctx context.Context, userID int, phoneNumber, language string) error
	VerifyEmailCode(ctx context.Context, userID int, emailAddr, code string) error
	VerifyPhoneCode(ctx context.Context, userID int, phoneNumber, code string) error
	ResendVerification(ctx context.Context, userID int, verificationType domain.VerificationType, language string) error
}

type verificationService struct {
	verificationRepo repository.VerificationRepository
	userRepo         repository.UserRepository
	emailService     email.EmailService
	smsService       twilio.SMSService
	maxAttempts      int
}

func NewVerificationService(
	verificationRepo repository.VerificationRepository,
	userRepo repository.UserRepository,
	emailService email.EmailService,
	smsService twilio.SMSService,
	maxAttempts int,
) VerificationService {
	return &verificationService{
		verificationRepo: verificationRepo,
		userRepo:         userRepo,
		emailService:     emailService,
		smsService:       smsService,
		maxAttempts:      maxAttempts,
	}
}

func (s *verificationService) SendEmailVerification(ctx context.Context, userID int, emailAddr, language string) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if email is already verified
	if user.IsEmailVerified() {
		return fmt.Errorf("email is already verified")
	}

	// Check if there's an active verification code
	activeCode, err := s.verificationRepo.GetActiveByIdentifier(ctx, emailAddr, domain.VerificationTypeEmail)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check active verification: %w", err)
	}

	// If there's an active code, don't send a new one
	if activeCode != nil {
		return fmt.Errorf("verification code already sent, please wait before requesting a new one")
	}

	// Generate new verification code
	code := utils.GenerateRandomCode(6)
	verificationCode := domain.NewVerificationCode(emailAddr, code, domain.VerificationTypeEmail)

	// Save verification code to database
	if err := s.verificationRepo.Create(ctx, verificationCode); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	// Send email
	if err := s.emailService.SendVerificationCode(ctx, emailAddr, code, language); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (s *verificationService) SendPhoneVerification(ctx context.Context, userID int, phoneNumber, language string) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if phone is already verified
	if user.IsPhoneVerified() {
		return fmt.Errorf("phone is already verified")
	}

	// Normalize phone number
	normalizedPhone := s.normalizePhoneNumber(phoneNumber)

	// Check if there's an active verification code
	activeCode, err := s.verificationRepo.GetActiveByIdentifier(ctx, normalizedPhone, domain.VerificationTypePhone)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check active verification: %w", err)
	}

	// If there's an active code, don't send a new one
	if activeCode != nil {
		return fmt.Errorf("verification code already sent, please wait before requesting a new one")
	}

	// Generate new verification code (for database tracking)
	code := utils.GenerateRandomCode(6)
	verificationCode := domain.NewVerificationCode(normalizedPhone, code, domain.VerificationTypePhone)

	// Save verification code to database
	if err := s.verificationRepo.Create(ctx, verificationCode); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	// Send SMS via Twilio (Twilio generates its own code)
	if err := s.smsService.SendVerificationCode(ctx, normalizedPhone); err != nil {
		return fmt.Errorf("failed to send verification SMS: %w", err)
	}

	return nil
}

func (s *verificationService) VerifyEmailCode(ctx context.Context, userID int, emailAddr, code string) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if email is already verified
	if user.IsEmailVerified() {
		return fmt.Errorf("email is already verified")
	}

	// Get verification code from database
	verificationCode, err := s.verificationRepo.GetActiveByIdentifier(ctx, emailAddr, domain.VerificationTypeEmail)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("verification code not found or expired")
		}
		return fmt.Errorf("failed to get verification code: %w", err)
	}

	// Check if code is expired
	if verificationCode.IsExpired() {
		return fmt.Errorf("verification code has expired")
	}

	// Check if code is already used
	if verificationCode.IsUsed() {
		return fmt.Errorf("verification code has already been used")
	}

	// Check attempts
	if verificationCode.Attempts >= s.maxAttempts {
		return fmt.Errorf("maximum verification attempts exceeded")
	}

	// Verify code
	if !verificationCode.VerifyCode(code) {
		// Increment attempts
		verificationCode.Attempts++
		if err := s.verificationRepo.UpdateAttempts(ctx, verificationCode.ID, verificationCode.Attempts); err != nil {
			return fmt.Errorf("failed to update attempts: %w", err)
		}
		return fmt.Errorf("invalid verification code")
	}

	// Mark code as used
	if err := s.verificationRepo.MarkAsUsed(ctx, verificationCode.ID); err != nil {
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	// Mark user email as verified
	if err := s.userRepo.SetEmailVerified(ctx, userID, time.Now()); err != nil {
		return fmt.Errorf("failed to mark email as verified: %w", err)
	}

	return nil
}

func (s *verificationService) VerifyPhoneCode(ctx context.Context, userID int, phoneNumber, code string) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if phone is already verified
	if user.IsPhoneVerified() {
		return fmt.Errorf("phone is already verified")
	}

	// Normalize phone number
	normalizedPhone := s.normalizePhoneNumber(phoneNumber)

	// Get verification code from database
	verificationCode, err := s.verificationRepo.GetActiveByIdentifier(ctx, normalizedPhone, domain.VerificationTypePhone)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("verification code not found or expired")
		}
		return fmt.Errorf("failed to get verification code: %w", err)
	}

	// Check if code is expired
	if verificationCode.IsExpired() {
		return fmt.Errorf("verification code has expired")
	}

	// Check if code is already used
	if verificationCode.IsUsed() {
		return fmt.Errorf("verification code has already been used")
	}

	// Check attempts
	if verificationCode.Attempts >= s.maxAttempts {
		return fmt.Errorf("maximum verification attempts exceeded")
	}

	// Verify code with Twilio
	isValid, err := s.smsService.VerifyCode(ctx, normalizedPhone, code)
	if err != nil {
		return fmt.Errorf("failed to verify SMS code: %w", err)
	}

	if !isValid {
		// Increment attempts
		verificationCode.Attempts++
		if err := s.verificationRepo.UpdateAttempts(ctx, verificationCode.ID, verificationCode.Attempts); err != nil {
			return fmt.Errorf("failed to update attempts: %w", err)
		}
		return fmt.Errorf("invalid verification code")
	}

	// Mark code as used
	if err := s.verificationRepo.MarkAsUsed(ctx, verificationCode.ID); err != nil {
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	// Mark user phone as verified
	if err := s.userRepo.SetPhoneVerified(ctx, userID, time.Now()); err != nil {
		return fmt.Errorf("failed to mark phone as verified: %w", err)
	}

	return nil
}

func (s *verificationService) ResendVerification(ctx context.Context, userID int, verificationType domain.VerificationType, language string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	switch verificationType {
	case domain.VerificationTypeEmail:
		if user.Email == nil {
			return fmt.Errorf("user has no email address")
		}
		return s.SendEmailVerification(ctx, userID, *user.Email, language)
	case domain.VerificationTypePhone:
		if user.Phone == nil {
			return fmt.Errorf("user has no phone number")
		}
		return s.SendPhoneVerification(ctx, userID, *user.Phone, language)
	default:
		return fmt.Errorf("invalid verification type")
	}
}

// normalizePhoneNumber removes spaces, dashes and ensures proper format
func (s *verificationService) normalizePhoneNumber(phone string) string {
	// Remove spaces, dashes, parentheses
	normalized := strings.ReplaceAll(phone, " ", "")
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, "(", "")
	normalized = strings.ReplaceAll(normalized, ")", "")

	// Ensure it starts with +
	if !strings.HasPrefix(normalized, "+") {
		normalized = "+" + normalized
	}

	return normalized
}
