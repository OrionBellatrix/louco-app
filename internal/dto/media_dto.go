package dto

import "time"

type MediaResponse struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	OriginalName string    `json:"original_name"`
	FileName     string    `json:"file_name"`
	FileURL      string    `json:"file_url"`
	MediaType    string    `json:"media_type"`
	MimeType     string    `json:"mime_type"`
	FileSize     int64     `json:"file_size"`
	Width        *int      `json:"width"`
	Height       *int      `json:"height"`
	Duration     *int      `json:"duration"`
	IsConverted  bool      `json:"is_converted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UploadResponse struct {
	MediaID      int     `json:"media_id"`
	FileType     string  `json:"file_type"`
	OriginalName string  `json:"original_name"`
	MimeType     string  `json:"mime_type"`
	FileSize     int64   `json:"file_size"`
	FileURL      string  `json:"file_url"`
	Width        *int    `json:"width,omitempty"`
	Height       *int    `json:"height,omitempty"`
	Duration     *int    `json:"duration,omitempty"`
	IsConverted  bool    `json:"is_converted"`
	NewFormat    *string `json:"new_format,omitempty"`
}

type MediaListRequest struct {
	Page      int    `json:"page" validate:"omitempty,min=1"`
	PageSize  int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	MediaType string `json:"media_type" validate:"omitempty,oneof=image video"`
	UserID    *int   `json:"user_id" validate:"omitempty,min=1"`
}

type MediaListResponse struct {
	Media      []MediaResponse `json:"media"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

type MediaUpdateRequest struct {
	OriginalName *string `json:"original_name" validate:"omitempty,max=255"`
	Width        *int    `json:"width" validate:"omitempty,min=1"`
	Height       *int    `json:"height" validate:"omitempty,min=1"`
	Duration     *int    `json:"duration" validate:"omitempty,min=1"`
}

// AWS Upload Configuration
type AWSConfig struct {
	Endpoint             string `json:"endpoint"`
	AccessKeyID          string `json:"access_key_id"`
	SecretAccessKey      string `json:"secret_access_key"`
	DefaultRegion        string `json:"default_region"`
	Bucket               string `json:"bucket"`
	UsePathStyleEndpoint bool   `json:"use_path_style_endpoint"`
}

// File validation constants
const (
	MaxImageSize = 10 * 1024 * 1024  // 10MB
	MaxVideoSize = 100 * 1024 * 1024 // 100MB
	ImageWidth   = 800               // Target width for image resize
)

var (
	AllowedImageTypes = []string{"image/jpeg", "image/jpg", "image/png", "image/webp", "image/heic", "image/heif"}
	AllowedVideoTypes = []string{"video/mp4", "video/mov", "video/webm", "video/avi", "video/quicktime"}
)
