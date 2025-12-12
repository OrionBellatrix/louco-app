package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/louco-event/internal/dto"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validators
	validate.RegisterValidation("alphanum_underscore_dot", validateAlphanumUnderscoreDot)
	validate.RegisterValidation("e164", validateE164Phone)
}

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(s interface{}) []dto.ValidationError {
	var errors []dto.ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, dto.ValidationError{
				Field:   strings.ToLower(err.Field()),
				Message: getErrorMessage(err),
				Value:   fmt.Sprintf("%v", err.Value()),
			})
		}
	}

	return errors
}

// Custom validator for username format (alphanumeric, underscore, dot)
func validateAlphanumUnderscoreDot(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._]+$`, username)
	return matched
}

// Custom validator for E.164 phone format
func validateE164Phone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(`^\+[1-9]\d{1,14}$`, phone)
	return matched
}

// Get user-friendly error messages
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Please enter a valid email address"
	case "min":
		return fmt.Sprintf("Minimum length is %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s characters", fe.Param())
	case "alphanum_underscore_dot":
		return "Username can only contain letters, numbers, dots and underscores"
	case "e164":
		return "Please enter a valid phone number with country code (e.g., +1234567890)"
	case "oneof":
		return fmt.Sprintf("Value must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Field validation failed on '%s' tag", fe.Tag())
	}
}

// Validate specific structs with custom logic
func ValidateRegisterStep1(req *dto.RegisterStep1Request) []dto.ValidationError {
	errors := ValidateStruct(req)

	// Custom validation: identifier must be valid email or phone
	if req.Identifier != "" {
		isEmail := isValidEmail(req.Identifier)
		isPhone := isValidE164Phone(req.Identifier)

		if !isEmail && !isPhone {
			errors = append(errors, dto.ValidationError{
				Field:   "identifier",
				Message: "Identifier must be a valid email address or phone number with country code",
			})
		}
	}

	return errors
}

func ValidateRegisterStep4(req *dto.RegisterStep4Request, userType string) []dto.ValidationError {
	errors := ValidateStruct(req)

	// Custom validation for creator users
	if userType == "creator" {
		if req.Address == nil || *req.Address == "" {
			errors = append(errors, dto.ValidationError{
				Field:   "address",
				Message: "Address is required for creator users",
			})
		}

		if req.CompanyName == nil || *req.CompanyName == "" {
			errors = append(errors, dto.ValidationError{
				Field:   "company_name",
				Message: "Company name is required for creator users",
			})
		}
	}

	return errors
}

// Example usage in handlers
func ValidateAndRespond(data interface{}) []dto.ValidationError {
	switch v := data.(type) {
	case *dto.RegisterStep1Request:
		return ValidateRegisterStep1(v)
	default:
		return ValidateStruct(data)
	}
}

// Helper functions for validation
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidE164Phone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(strings.TrimSpace(phone))
}
