package dto

// Standard API Response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Success response helper
func NewSuccessResponse(message string, data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// Error response helper
func NewErrorResponse(message string, errors interface{}) *APIResponse {
	return &APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}

// Validation Error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// Pagination Request
type PaginationRequest struct {
	Page     int `json:"page" validate:"omitempty,min=1" form:"page" query:"page"`
	PageSize int `json:"page_size" validate:"omitempty,min=1,max=100" form:"page_size" query:"page_size"`
}

// Pagination Response (alias for PaginationMeta for consistency)
type PaginationResponse = PaginationMeta

// Pagination
type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Helper functions for pagination
func (p *PaginationRequest) GetPageWithDefault() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *PaginationRequest) GetPageSizeWithDefault() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

func (p *PaginationRequest) GetOffset() int {
	page := p.GetPageWithDefault()
	pageSize := p.GetPageSizeWithDefault()
	return (page - 1) * pageSize
}

// Health Check
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	Database  string `json:"database"`
}

// Generic List Response
type ListResponse struct {
	Items      interface{}    `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}

// Error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeInternalServer = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest     = "BAD_REQUEST"
)

// Common error messages (i18n keys)
const (
	MsgSuccess               = "common.success"
	MsgCreated               = "common.created"
	MsgUpdated               = "common.updated"
	MsgDeleted               = "common.deleted"
	MsgValidationFailed      = "common.validation_failed"
	MsgUnauthorized          = "common.unauthorized"
	MsgForbidden             = "common.forbidden"
	MsgNotFound              = "common.not_found"
	MsgInternalServerError   = "common.internal_server_error"
	MsgBadRequest            = "common.bad_request"
	MsgEmailAlreadyExists    = "user.email_already_exists"
	MsgPhoneAlreadyExists    = "user.phone_already_exists"
	MsgUsernameAlreadyExists = "user.username_already_exists"
	MsgInvalidCredentials    = "auth.invalid_credentials"
	MsgLoginSuccess          = "auth.login_success"
	MsgLogoutSuccess         = "auth.logout_success"
	MsgRegistrationSuccess   = "auth.registration_success"
	MsgProfileUpdated        = "user.profile_updated"
	MsgPasswordChanged       = "auth.password_changed"
	MsgFileUploaded          = "media.file_uploaded"
	MsgFileDeleted           = "media.file_deleted"
	MsgInvalidFileType       = "media.invalid_file_type"
	MsgFileTooLarge          = "media.file_too_large"
)
