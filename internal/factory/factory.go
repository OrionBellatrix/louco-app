package factory

import (
	"fmt"

	"github.com/louco-event/internal/config"
	"github.com/louco-event/internal/i18n"
	"github.com/louco-event/internal/repository"
	"github.com/louco-event/internal/repository/postgres"
	"github.com/louco-event/internal/service"
	"github.com/louco-event/pkg/cache"
	"github.com/louco-event/pkg/database"
	"github.com/louco-event/pkg/email"
	"github.com/louco-event/pkg/logger"
	"github.com/louco-event/pkg/stripe"
	"github.com/louco-event/pkg/twilio"
)

type Dependencies struct {
	// Database
	DB *database.Database

	// Repositories
	UserRepo             repository.UserRepository
	MediaRepo            repository.MediaRepository
	IndustryRepo         repository.IndustryRepository
	CreatorRepo          repository.CreatorRepository
	CategoryRepo         repository.CategoryRepository
	VerificationRepo     repository.VerificationRepository
	FollowRepo           repository.FollowRepository
	EventRepo            repository.EventRepository
	AddressRepo          repository.AddressRepository
	TicketRepo           repository.TicketRepository
	InvitationRepo       repository.InvitationRepository
	UserSubscriptionRepo repository.UserSubscriptionRepository
	SubscriptionPlanRepo repository.SubscriptionPlanRepository

	// Services
	UserService         service.UserService
	MediaService        service.MediaService
	JWTService          service.JWTService
	IndustryService     service.IndustryService
	CreatorService      service.CreatorService
	CategoryService     service.CategoryService
	VerificationService service.VerificationService
	FollowService       *service.FollowService
	EventService        service.EventService
	AddressService      service.AddressService
	TicketService       service.TicketService
	InvitationService   service.InvitationService
	SubscriptionService service.SubscriptionService

	// External Services
	StripeService *stripe.StripeService

	// I18n
	I18n *i18n.I18n

	// Logger
	Logger *logger.Logger
}

func NewDependencies(cfg *config.Config, logger *logger.Logger) (*Dependencies, error) {
	// Initialize database
	db, err := database.New(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Redis cache
	redisAddr := cfg.Redis.Host + ":" + cfg.Redis.Port
	redisCache := cache.NewRedisCache(redisAddr, cfg.Redis.Password, cfg.Redis.DB)

	// Run auto migration
	if err := db.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to run auto migration: %w", err)
	}

	// Initialize i18n
	i18nService, err := i18n.New("internal/i18n/locales", "en")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize i18n: %w", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db.DB)
	mediaRepo := postgres.NewMediaRepository(db.DB)
	industryRepo := postgres.NewIndustryRepository(db.DB)
	creatorRepo := postgres.NewCreatorRepository(db.DB)
	categoryRepo := postgres.NewCategoryRepository(db.DB)
	verificationRepo := postgres.NewVerificationRepository(db.DB)
	followRepo := postgres.NewFollowRepository(db.DB)
	eventRepo := postgres.NewEventRepository(db.DB)
	addressRepo := postgres.NewAddressRepository(db.DB)
	ticketRepo := postgres.NewTicketRepository(db.DB)
	invitationRepo := postgres.NewInvitationRepository(db.DB)
	userSubscriptionRepo := postgres.NewUserSubscriptionRepository(db.DB, logger)
	subscriptionPlanRepo := postgres.NewSubscriptionPlanRepository(db.DB, logger)

	// Initialize services
	jwtService := service.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)
	industryService := service.NewIndustryService(industryRepo, *logger.Logger)
	categoryService := service.NewCategoryService(categoryRepo, mediaRepo, redisCache, logger)
	creatorService := service.NewCreatorService(creatorRepo, userRepo, industryRepo, mediaRepo, logger)

	// Initialize email service
	emailService := email.NewEmailService(email.EmailConfig{
		SMTPHost:     cfg.Email.SMTPHost,
		SMTPPort:     cfg.Email.SMTPPort,
		SMTPUsername: cfg.Email.SMTPUsername,
		SMTPPassword: cfg.Email.SMTPPassword,
		FromEmail:    cfg.Email.FromEmail,
		FromName:     cfg.Email.FromName,
	})

	// Initialize SMS service
	smsService := twilio.NewSMSService(twilio.SMSConfig{
		AccountSID:    cfg.Twilio.AccountSID,
		AuthToken:     cfg.Twilio.AuthToken,
		ServiceSID:    cfg.Twilio.ServiceSID,
		ReviewerPhone: cfg.Twilio.ReviewerPhone,
		ReviewerOTP:   cfg.Twilio.ReviewerOTP,
		MaxAttempts:   cfg.Twilio.MaxAttempts,
	})

	// Initialize verification service
	verificationService := service.NewVerificationService(
		verificationRepo,
		userRepo,
		emailService,
		smsService,
		3, // max attempts
	)

	userService := service.NewUserService(userRepo, mediaRepo, creatorService, jwtService, logger)

	// AWS config for media service
	awsConfig := service.AWSConfig{
		Endpoint:             cfg.AWS.Endpoint,
		AccessKeyID:          cfg.AWS.AccessKeyID,
		SecretAccessKey:      cfg.AWS.SecretAccessKey,
		DefaultRegion:        cfg.AWS.DefaultRegion,
		Bucket:               cfg.AWS.Bucket,
		UsePathStyleEndpoint: cfg.AWS.UsePathStyleEndpoint,
	}
	mediaService := service.NewMediaService(mediaRepo, awsConfig, logger)

	// Initialize follow service
	followService := service.NewFollowService(followRepo, userRepo)

	// Initialize subscription service first (needed by event service)
	subscriptionService := service.NewSubscriptionService(userSubscriptionRepo, subscriptionPlanRepo, logger)

	// Initialize event-related services
	addressService := service.NewAddressService(addressRepo, *logger.Logger)
	ticketService := service.NewTicketService(ticketRepo, eventRepo, *logger.Logger)
	invitationService := service.NewInvitationService(invitationRepo, eventRepo, userRepo, *logger.Logger)
	eventService := service.NewEventService(eventRepo, addressRepo, ticketRepo, invitationRepo, categoryRepo, creatorRepo, mediaRepo, subscriptionService, logger)

	// Initialize Stripe service
	stripeConfig := stripe.StripeConfig{
		SecretKey:      cfg.Stripe.SecretKey,
		PublishableKey: cfg.Stripe.PublishableKey,
		WebhookSecret:  cfg.Stripe.WebhookSecret,
		Currency:       cfg.Stripe.Currency,
		Environment:    cfg.Stripe.Environment,
	}
	stripeService := stripe.NewStripeService(stripeConfig, logger)

	return &Dependencies{
		DB:                   db,
		UserRepo:             userRepo,
		MediaRepo:            mediaRepo,
		IndustryRepo:         industryRepo,
		CreatorRepo:          creatorRepo,
		CategoryRepo:         categoryRepo,
		VerificationRepo:     verificationRepo,
		FollowRepo:           followRepo,
		EventRepo:            eventRepo,
		AddressRepo:          addressRepo,
		TicketRepo:           ticketRepo,
		InvitationRepo:       invitationRepo,
		UserSubscriptionRepo: userSubscriptionRepo,
		SubscriptionPlanRepo: subscriptionPlanRepo,
		UserService:          userService,
		MediaService:         mediaService,
		JWTService:           jwtService,
		IndustryService:      industryService,
		CreatorService:       creatorService,
		CategoryService:      categoryService,
		VerificationService:  verificationService,
		FollowService:        followService,
		EventService:         eventService,
		AddressService:       addressService,
		TicketService:        ticketService,
		InvitationService:    invitationService,
		SubscriptionService:  subscriptionService,
		StripeService:        stripeService,
		I18n:                 i18nService,
		Logger:               logger,
	}, nil
}

func (d *Dependencies) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}
