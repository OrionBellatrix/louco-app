package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
	"github.com/louco-event/pkg/logger"
)

type MediaService interface {
	UploadFile(ctx context.Context, userID int, file *multipart.FileHeader, fileContent io.Reader) (*dto.UploadResponse, error)
	GetMediaByID(ctx context.Context, mediaID int) (*dto.MediaResponse, error)
	GetUserMedia(ctx context.Context, userID int, req *dto.MediaListRequest) (*dto.MediaListResponse, error)
	DeleteMedia(ctx context.Context, userID int, mediaID int) error
	UpdateMedia(ctx context.Context, userID int, mediaID int, req *dto.MediaUpdateRequest) error
}

type mediaService struct {
	mediaRepo repository.MediaRepository
	s3Client  *s3.S3
	bucket    string
	logger    *logger.Logger
}

type AWSConfig struct {
	Endpoint             string
	AccessKeyID          string
	SecretAccessKey      string
	DefaultRegion        string
	Bucket               string
	UsePathStyleEndpoint bool
}

func NewMediaService(mediaRepo repository.MediaRepository, awsConfig AWSConfig, logger *logger.Logger) MediaService {
	// Create AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:         aws.String(awsConfig.Endpoint),
		Region:           aws.String(awsConfig.DefaultRegion),
		Credentials:      credentials.NewStaticCredentials(awsConfig.AccessKeyID, awsConfig.SecretAccessKey, ""),
		S3ForcePathStyle: aws.Bool(awsConfig.UsePathStyleEndpoint),
	}))

	s3Client := s3.New(sess)

	return &mediaService{
		mediaRepo: mediaRepo,
		s3Client:  s3Client,
		bucket:    awsConfig.Bucket,
		logger:    logger,
	}
}

func (s *mediaService) UploadFile(ctx context.Context, userID int, file *multipart.FileHeader, fileContent io.Reader) (*dto.UploadResponse, error) {
	// Validate file type
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = getContentTypeFromExtension(file.Filename)
	}

	mediaType := domain.GetMediaTypeFromMime(mimeType)
	if mediaType == "" {
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}

	// Validate file size
	if err := s.validateFileSize(file.Size, mediaType); err != nil {
		return nil, err
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := fmt.Sprintf("uploads/%d/%s", userID, fileName)

	// Upload to S3
	_, err := s.s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(filePath),
		Body:        aws.ReadSeekCloser(fileContent),
		ContentType: aws.String(mimeType),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to upload file to S3")
		return nil, fmt.Errorf("failed to upload file")
	}

	// Generate file URL
	endpoint := ""
	if s.s3Client.Config.Endpoint != nil {
		endpoint = strings.TrimPrefix(*s.s3Client.Config.Endpoint, "https://")
	}
	fileURL := fmt.Sprintf("https://%s.%s/%s", s.bucket, endpoint, filePath)

	// Create media record
	media := domain.NewMedia(userID, file.Filename, fileName, filePath, fileURL, mediaType, mimeType, file.Size)

	// Handle format conversion for HEIC files
	isConverted := false
	newFormat := ""
	if mimeType == "image/heic" || mimeType == "image/heif" {
		// In a real implementation, you would convert HEIC to JPEG here
		// For now, we'll just mark it as needing conversion
		newFormat = "jpeg"
		// media.MarkAsConverted() // Would be called after actual conversion
	}

	// Set dimensions for images (in a real implementation, you'd extract these from the file)
	if mediaType == domain.MediaTypeImage {
		// Placeholder dimensions - in reality, you'd use an image processing library
		width, height := 800, 600 // These would be calculated from the actual image
		media.SetDimensions(width, height)
	}

	err = s.mediaRepo.Create(ctx, media)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create media record")
		// Clean up uploaded file
		s.deleteFromS3(ctx, filePath)
		return nil, fmt.Errorf("failed to save media record")
	}

	response := &dto.UploadResponse{
		MediaID:      media.ID,
		FileType:     string(mediaType),
		OriginalName: media.OriginalName,
		MimeType:     media.MimeType,
		FileSize:     media.FileSize,
		FileURL:      media.FileURL,
		IsConverted:  isConverted,
	}

	if media.Width != nil {
		response.Width = media.Width
	}
	if media.Height != nil {
		response.Height = media.Height
	}
	if media.Duration != nil {
		response.Duration = media.Duration
	}
	if newFormat != "" {
		response.NewFormat = &newFormat
	}

	return response, nil
}

func (s *mediaService) GetMediaByID(ctx context.Context, mediaID int) (*dto.MediaResponse, error) {
	media, err := s.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		return nil, fmt.Errorf("media not found")
	}

	return s.mapMediaToResponse(media), nil
}

func (s *mediaService) GetUserMedia(ctx context.Context, userID int, req *dto.MediaListRequest) (*dto.MediaListResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	mediaList, err := s.mediaRepo.GetByUserID(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user media")
	}

	mediaResponses := make([]dto.MediaResponse, len(mediaList))
	for i, media := range mediaList {
		mediaResponses[i] = *s.mapMediaToResponse(media)
	}

	// Calculate total pages (simplified)
	totalPages := 1
	if len(mediaList) == pageSize {
		totalPages = page + 1
	}

	return &dto.MediaListResponse{
		Media:      mediaResponses,
		Total:      len(mediaList),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *mediaService) DeleteMedia(ctx context.Context, userID int, mediaID int) error {
	media, err := s.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		return fmt.Errorf("media not found")
	}

	// Check ownership
	if media.UserID != userID {
		return fmt.Errorf("unauthorized to delete this media")
	}

	// Delete from S3
	err = s.deleteFromS3(ctx, media.FilePath)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to delete file from S3")
		// Continue with database deletion even if S3 deletion fails
	}

	// Delete from database
	err = s.mediaRepo.Delete(ctx, mediaID)
	if err != nil {
		return fmt.Errorf("failed to delete media record")
	}

	return nil
}

func (s *mediaService) UpdateMedia(ctx context.Context, userID int, mediaID int, req *dto.MediaUpdateRequest) error {
	media, err := s.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		return fmt.Errorf("media not found")
	}

	// Check ownership
	if media.UserID != userID {
		return fmt.Errorf("unauthorized to update this media")
	}

	// Update fields
	if req.OriginalName != nil {
		media.OriginalName = *req.OriginalName
	}
	if req.Width != nil {
		media.SetDimensions(*req.Width, *media.Height)
	}
	if req.Height != nil {
		media.SetDimensions(*media.Width, *req.Height)
	}
	if req.Duration != nil {
		media.SetDuration(*req.Duration)
	}

	media.UpdatedAt = time.Now()

	return s.mediaRepo.Update(ctx, media)
}

func (s *mediaService) deleteFromS3(ctx context.Context, filePath string) error {
	_, err := s.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})
	return err
}

func (s *mediaService) validateFileSize(size int64, mediaType domain.MediaType) error {
	switch mediaType {
	case domain.MediaTypeImage:
		if size > dto.MaxImageSize {
			return fmt.Errorf("image file too large (max %d bytes)", dto.MaxImageSize)
		}
	case domain.MediaTypeVideo:
		if size > dto.MaxVideoSize {
			return fmt.Errorf("video file too large (max %d bytes)", dto.MaxVideoSize)
		}
	}
	return nil
}

func (s *mediaService) mapMediaToResponse(media *domain.Media) *dto.MediaResponse {
	return &dto.MediaResponse{
		ID:           media.ID,
		UserID:       media.UserID,
		OriginalName: media.OriginalName,
		FileName:     media.FileName,
		FileURL:      media.FileURL,
		MediaType:    string(media.MediaType),
		MimeType:     media.MimeType,
		FileSize:     media.FileSize,
		Width:        media.Width,
		Height:       media.Height,
		Duration:     media.Duration,
		IsConverted:  media.IsConverted,
		CreatedAt:    media.CreatedAt,
		UpdatedAt:    media.UpdatedAt,
	}
}

func getContentTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".heic":
		return "image/heic"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/mov"
	case ".webm":
		return "video/webm"
	case ".avi":
		return "video/avi"
	default:
		return "application/octet-stream"
	}
}
