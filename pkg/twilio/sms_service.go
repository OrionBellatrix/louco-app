package twilio

import (
	"context"
	"fmt"
	"strings"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type SMSService interface {
	SendVerificationCode(ctx context.Context, phoneNumber string) error
	VerifyCode(ctx context.Context, phoneNumber, code string) (bool, error)
}

type smsService struct {
	client        *twilio.RestClient
	serviceSID    string
	reviewerPhone string
	reviewerOTP   string
	maxAttempts   int
}

type SMSConfig struct {
	AccountSID    string
	AuthToken     string
	ServiceSID    string
	ReviewerPhone string
	ReviewerOTP   string
	MaxAttempts   int
}

func NewSMSService(config SMSConfig) SMSService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	return &smsService{
		client:        client,
		serviceSID:    config.ServiceSID,
		reviewerPhone: config.ReviewerPhone,
		reviewerOTP:   config.ReviewerOTP,
		maxAttempts:   config.MaxAttempts,
	}
}

func (s *smsService) SendVerificationCode(ctx context.Context, phoneNumber string) error {
	// Normalize phone number (remove spaces, dashes, etc.)
	normalizedPhone := s.normalizePhoneNumber(phoneNumber)

	// Check if this is the reviewer phone number
	if normalizedPhone == s.normalizePhoneNumber(s.reviewerPhone) {
		// For reviewer phone, we don't actually send SMS, just return success
		return nil
	}

	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(normalizedPhone)
	params.SetChannel("sms")

	_, err := s.client.VerifyV2.CreateVerification(s.serviceSID, params)
	if err != nil {
		return fmt.Errorf("failed to send verification SMS: %w", err)
	}

	return nil
}

func (s *smsService) VerifyCode(ctx context.Context, phoneNumber, code string) (bool, error) {
	// Normalize phone number
	normalizedPhone := s.normalizePhoneNumber(phoneNumber)

	// Check if this is the reviewer phone number with reviewer OTP
	if normalizedPhone == s.normalizePhoneNumber(s.reviewerPhone) && code == s.reviewerOTP {
		return true, nil
	}

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(normalizedPhone)
	params.SetCode(code)

	resp, err := s.client.VerifyV2.CreateVerificationCheck(s.serviceSID, params)
	if err != nil {
		return false, fmt.Errorf("failed to verify SMS code: %w", err)
	}

	// Check if verification was successful
	if resp.Status != nil && *resp.Status == "approved" {
		return true, nil
	}

	return false, nil
}

// normalizePhoneNumber removes spaces, dashes and ensures proper format
func (s *smsService) normalizePhoneNumber(phone string) string {
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
