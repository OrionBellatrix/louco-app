package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
	"github.com/louco-event/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterStep1(ctx context.Context, req *dto.RegisterStep1Request) (*dto.RegisterStep1Response, error)
	RegisterStep4(ctx context.Context, userID int, req *dto.RegisterStep4Request) error
	SetUsername(ctx context.Context, userID int, req *dto.SetUsernameRequest) error
	CheckUsername(ctx context.Context, req *dto.CheckUsernameRequest) (*dto.CheckUsernameResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	SocialLogin(ctx context.Context, req *dto.SocialLoginRequest) (*dto.LoginResponse, error)
	GetProfile(ctx context.Context, userID int) (*dto.UserProfileResponse, error)
	GetByID(ctx context.Context, userID uint) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID int, req *dto.UpdateProfileRequest) error
	UpdateContact(ctx context.Context, userID int, req *dto.UpdateContactRequest) error
	ChangePassword(ctx context.Context, userID int, req *dto.ChangePasswordRequest) error
	SetProfilePic(ctx context.Context, userID int, req *dto.SetProfilePicRequest) error
	SetCoverPic(ctx context.Context, userID int, req *dto.SetCoverPicRequest) error
	DeactivateAccount(ctx context.Context, userID int) error
	GetUserList(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error)
}

type userService struct {
	userRepo   repository.UserRepository
	mediaRepo  repository.MediaRepository
	creatorSvc CreatorService
	jwtSvc     JWTService
	logger     *logger.Logger
}

func NewUserService(userRepo repository.UserRepository, mediaRepo repository.MediaRepository, creatorSvc CreatorService, jwtSvc JWTService, logger *logger.Logger) UserService {
	return &userService{
		userRepo:   userRepo,
		mediaRepo:  mediaRepo,
		creatorSvc: creatorSvc,
		jwtSvc:     jwtSvc,
		logger:     logger,
	}
}

func (s *userService) GetByID(ctx context.Context, userID uint) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, int(userID))
	if err != nil {
		s.logger.Error().
			Err(err).
			Uint("user_id", userID).
			Msg("Failed to get user by ID")
		return nil, err
	}

	if user == nil {
		s.logger.Warn().
			Uint("user_id", userID).
			Msg("User not found")
		return nil, nil
	}

	return user, nil
}

func (s *userService) RegisterStep1(ctx context.Context, req *dto.RegisterStep1Request) (*dto.RegisterStep1Response, error) {
	// Determine if identifier is email or phone
	isEmail := s.isValidEmail(req.Identifier)
	isPhone := s.isValidPhone(req.Identifier)

	if !isEmail && !isPhone {
		return nil, fmt.Errorf("identifier must be a valid email or phone number")
	}

	// Check if identifier already exists
	if isEmail {
		exists, err := s.userRepo.IsEmailExists(ctx, req.Identifier, nil)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to check email existence")
			return nil, fmt.Errorf("failed to validate email")
		}
		if exists {
			return nil, fmt.Errorf("email already exists")
		}
	} else if isPhone {
		exists, err := s.userRepo.IsPhoneExists(ctx, req.Identifier, nil)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to check phone existence")
			return nil, fmt.Errorf("failed to validate phone")
		}
		if exists {
			return nil, fmt.Errorf("phone already exists")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to hash password")
		return nil, fmt.Errorf("failed to process password")
	}

	// Create user
	userType := domain.UserType(req.UserType)
	user := domain.NewUser("", string(hashedPassword), userType)

	if isEmail {
		user.SetEmail(req.Identifier)
	} else if isPhone {
		user.SetPhone(req.Identifier)
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create user")
		return nil, fmt.Errorf("failed to create user")
	}

	// Generate JWT token
	claims := &dto.JWTClaims{
		UserID:   user.ID,
		UserType: string(user.UserType),
	}
	if user.Email != nil {
		claims.Email = *user.Email
	}

	token, err := s.jwtSvc.GenerateToken(claims)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate token")
		return nil, fmt.Errorf("failed to generate token")
	}

	return &dto.RegisterStep1Response{
		UserID: user.ID,
		Token:  token,
	}, nil
}

func (s *userService) RegisterStep4(ctx context.Context, userID int, req *dto.RegisterStep4Request) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Update profile
	user.UpdateProfile(req.FullName, req.Biography, req.BirthDate)

	// Set email and phone if provided
	if req.Email != nil {
		// Check uniqueness
		exists, err := s.userRepo.IsEmailExists(ctx, *req.Email, &userID)
		if err != nil {
			return fmt.Errorf("failed to validate email")
		}
		if exists {
			return fmt.Errorf("email already exists")
		}
		user.SetEmail(*req.Email)
	}

	if req.Phone != nil {
		// Check uniqueness
		exists, err := s.userRepo.IsPhoneExists(ctx, *req.Phone, &userID)
		if err != nil {
			return fmt.Errorf("failed to validate phone")
		}
		if exists {
			return fmt.Errorf("phone already exists")
		}
		user.SetPhone(*req.Phone)
	}

	// Update user first
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Handle creator-specific data if user is creator
	if user.UserType == domain.UserTypeCreator {
		// Validate creator-specific fields
		if req.CompanyName == nil || *req.CompanyName == "" {
			return fmt.Errorf("company name is required for creator users")
		}
		if req.Address == nil || *req.Address == "" {
			return fmt.Errorf("address is required for creator users")
		}
		if req.EstimatedTickets == nil || *req.EstimatedTickets <= 0 {
			return fmt.Errorf("estimated tickets must be greater than 0 for creator users")
		}
		if req.EstimatedEvents == nil || *req.EstimatedEvents <= 0 {
			return fmt.Errorf("estimated events must be greater than 0 for creator users")
		}
		if req.IndustryIDs == nil || len(req.IndustryIDs) == 0 {
			return fmt.Errorf("at least one industry must be selected for creator users")
		}

		// Create creator profile
		createReq := &dto.CreateCreatorRequest{
			CompanyName:      *req.CompanyName,
			Address:          *req.Address,
			EstimatedTickets: *req.EstimatedTickets,
			EstimatedEvents:  *req.EstimatedEvents,
			IndustryIDs:      req.IndustryIDs,
		}

		_, err := s.creatorSvc.CreateCreator(ctx, userID, createReq)
		if err != nil {
			return fmt.Errorf("failed to create creator profile: %w", err)
		}
	}

	return nil
}

func (s *userService) SetUsername(ctx context.Context, userID int, req *dto.SetUsernameRequest) error {
	// Check if username already exists
	exists, err := s.userRepo.IsUsernameExists(ctx, req.Username, &userID)
	if err != nil {
		return fmt.Errorf("failed to validate username")
	}
	if exists {
		return fmt.Errorf("username already exists")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	user.SetUsername(req.Username)
	return s.userRepo.Update(ctx, user)
}

func (s *userService) CheckUsername(ctx context.Context, req *dto.CheckUsernameRequest) (*dto.CheckUsernameResponse, error) {
	exists, err := s.userRepo.IsUsernameExists(ctx, req.Username, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check username")
	}

	return &dto.CheckUsernameResponse{
		Exists: exists,
	}, nil
}

func (s *userService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user *domain.User
	var err error

	// Try to find user by email, phone, or username
	if user, err = s.userRepo.GetByEmail(ctx, req.Identifier); err != nil {
		if user, err = s.userRepo.GetByPhone(ctx, req.Identifier); err != nil {
			if user, err = s.userRepo.GetByUsername(ctx, req.Identifier); err != nil {
				return nil, fmt.Errorf("invalid credentials")
			}
		}
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	claims := &dto.JWTClaims{
		UserID:   user.ID,
		UserType: string(user.UserType),
	}
	if user.Email != nil {
		claims.Email = *user.Email
	}
	if user.Username != nil {
		claims.Username = *user.Username
	}

	token, err := s.jwtSvc.GenerateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	return &dto.LoginResponse{
		User:  s.mapUserToResponse(user),
		Token: token,
	}, nil
}

func (s *userService) SocialLogin(ctx context.Context, req *dto.SocialLoginRequest) (*dto.LoginResponse, error) {
	var user *domain.User
	var err error

	// Try to find existing user by social ID
	switch req.Provider {
	case "apple":
		user, err = s.userRepo.GetByAppleID(ctx, req.SocialID)
	case "google":
		user, err = s.userRepo.GetByGoogleID(ctx, req.SocialID)
	default:
		return nil, fmt.Errorf("unsupported provider")
	}

	// If user doesn't exist, create new one
	if err != nil {
		userType := domain.UserType(req.UserType)
		user = domain.NewUser("", "", userType) // No password for social login

		if req.FullName != nil {
			user.FullName = *req.FullName
		}
		if req.Email != nil {
			user.SetEmail(*req.Email)
		}

		user.SetSocialID(req.Provider, req.SocialID)

		err = s.userRepo.Create(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user")
		}
	}

	// Generate JWT token
	claims := &dto.JWTClaims{
		UserID:   user.ID,
		UserType: string(user.UserType),
	}
	if user.Email != nil {
		claims.Email = *user.Email
	}
	if user.Username != nil {
		claims.Username = *user.Username
	}

	token, err := s.jwtSvc.GenerateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	return &dto.LoginResponse{
		User:  s.mapUserToResponse(user),
		Token: token,
	}, nil
}

func (s *userService) GetProfile(ctx context.Context, userID int) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	userResponse := s.mapUserToResponse(user)

	// Get profile picture if exists
	if user.ProfilePicID != nil {
		profilePic, err := s.mediaRepo.GetByID(ctx, *user.ProfilePicID)
		if err == nil {
			profilePicResponse := s.mapMediaToResponse(profilePic)
			userResponse.ProfilePic = &profilePicResponse
		}
	}

	// Get cover picture if exists
	if user.CoverPicID != nil {
		coverPic, err := s.mediaRepo.GetByID(ctx, *user.CoverPicID)
		if err == nil {
			coverPicResponse := s.mapMediaToResponse(coverPic)
			userResponse.CoverPic = &coverPicResponse
		}
	}

	response := &dto.UserProfileResponse{
		User: userResponse,
	}

	return response, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID int, req *dto.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	user.UpdateProfile(req.FullName, req.Biography, req.BirthDate)

	// Validate required fields for creator users
	if err := user.ValidateRequiredFields(); err != nil {
		return err
	}

	return s.userRepo.Update(ctx, user)
}

func (s *userService) UpdateContact(ctx context.Context, userID int, req *dto.UpdateContactRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if req.Email != nil {
		exists, err := s.userRepo.IsEmailExists(ctx, *req.Email, &userID)
		if err != nil {
			return fmt.Errorf("failed to validate email")
		}
		if exists {
			return fmt.Errorf("email already exists")
		}
		user.SetEmail(*req.Email)
	}

	if req.Phone != nil {
		exists, err := s.userRepo.IsPhoneExists(ctx, *req.Phone, &userID)
		if err != nil {
			return fmt.Errorf("failed to validate phone")
		}
		if exists {
			return fmt.Errorf("phone already exists")
		}
		user.SetPhone(*req.Phone)
	}

	return s.userRepo.Update(ctx, user)
}

func (s *userService) ChangePassword(ctx context.Context, userID int, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to process new password")
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

func (s *userService) SetProfilePic(ctx context.Context, userID int, req *dto.SetProfilePicRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	user.SetProfilePicID(req.MediaID)
	return s.userRepo.Update(ctx, user)
}

func (s *userService) SetCoverPic(ctx context.Context, userID int, req *dto.SetCoverPicRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	user.SetCoverPicID(req.MediaID)
	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeactivateAccount(ctx context.Context, userID int) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *userService) GetUserList(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	users, err := s.userRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user list")
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = s.mapUserToResponse(user)
	}

	// For simplicity, we're not implementing total count here
	// In a real application, you'd want to add a Count method to the repository
	totalPages := 1
	if len(users) == pageSize {
		totalPages = page + 1 // Estimate
	}

	return &dto.UserListResponse{
		Users:      userResponses,
		Total:      len(users),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) mapUserToResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:              user.ID,
		FullName:        user.FullName,
		Username:        user.Username,
		Email:           user.Email,
		Phone:           user.Phone,
		UserType:        string(user.UserType),
		Biography:       user.Biography,
		BirthDate:       user.BirthDate,
		ProfilePicID:    user.ProfilePicID,
		CoverPicID:      user.CoverPicID,
		EmailVerifiedAt: user.EmailVerifiedAt,
		PhoneVerifiedAt: user.PhoneVerifiedAt,
		IsActive:        user.IsActive,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func (s *userService) mapMediaToResponse(media *domain.Media) dto.MediaResponse {
	return dto.MediaResponse{
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

func (s *userService) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (s *userService) isValidPhone(phone string) bool {
	// E.164 format: +[country code][number]
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(strings.TrimSpace(phone))
}
