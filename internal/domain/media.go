package domain

import (
	"time"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

type Media struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	OriginalName string    `json:"original_name" db:"original_name"`
	FileName     string    `json:"file_name" db:"file_name"`
	FilePath     string    `json:"file_path" db:"file_path"`
	FileURL      string    `json:"file_url" db:"file_url"`
	MediaType    MediaType `json:"media_type" db:"media_type"`
	MimeType     string    `json:"mime_type" db:"mime_type"`
	FileSize     int64     `json:"file_size" db:"file_size"`
	Width        *int      `json:"width" db:"width"`
	Height       *int      `json:"height" db:"height"`
	Duration     *int      `json:"duration" db:"duration"` // for videos in seconds
	IsConverted  bool      `json:"is_converted" db:"is_converted"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewMedia(userID int, originalName, fileName, filePath, fileURL string, mediaType MediaType, mimeType string, fileSize int64) *Media {
	return &Media{
		UserID:       userID,
		OriginalName: originalName,
		FileName:     fileName,
		FilePath:     filePath,
		FileURL:      fileURL,
		MediaType:    mediaType,
		MimeType:     mimeType,
		FileSize:     fileSize,
		IsConverted:  false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (m *Media) SetDimensions(width, height int) {
	m.Width = &width
	m.Height = &height
	m.UpdatedAt = time.Now()
}

func (m *Media) SetDuration(duration int) {
	m.Duration = &duration
	m.UpdatedAt = time.Now()
}

func (m *Media) MarkAsConverted() {
	m.IsConverted = true
	m.UpdatedAt = time.Now()
}

func (m *Media) IsImage() bool {
	return m.MediaType == MediaTypeImage
}

func (m *Media) IsVideo() bool {
	return m.MediaType == MediaTypeVideo
}

func (m *Media) GetSizeInMB() float64 {
	return float64(m.FileSize) / (1024 * 1024)
}

// Supported formats
var (
	SupportedImageFormats = []string{"jpg", "jpeg", "png", "webp", "heic"}
	SupportedVideoFormats = []string{"mp4", "mov", "webm", "avi"}

	ImageMimeTypes = map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
		"image/heic": true,
		"image/heif": true,
	}

	VideoMimeTypes = map[string]bool{
		"video/mp4":       true,
		"video/mov":       true,
		"video/webm":      true,
		"video/avi":       true,
		"video/quicktime": true,
	}
)

func IsImageMimeType(mimeType string) bool {
	return ImageMimeTypes[mimeType]
}

func IsVideoMimeType(mimeType string) bool {
	return VideoMimeTypes[mimeType]
}

func GetMediaTypeFromMime(mimeType string) MediaType {
	if IsImageMimeType(mimeType) {
		return MediaTypeImage
	}
	if IsVideoMimeType(mimeType) {
		return MediaTypeVideo
	}
	return ""
}
