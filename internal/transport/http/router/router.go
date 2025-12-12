package router

import (
	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/factory"
	"github.com/louco-event/internal/middleware"
	"github.com/louco-event/internal/transport/http/handler"
)

func SetupRoutes(r *gin.Engine, deps *factory.Dependencies) {
	// Initialize handlers
	authHandler := handler.NewAuthHandler(deps.UserService, deps.VerificationService, deps.I18n)
	userHandler := handler.NewUserHandler(deps.UserService, deps.I18n)
	mediaHandler := handler.NewMediaHandler(deps.MediaService, deps.I18n)
	industryHandler := handler.NewIndustryHandler(deps.IndustryService, *deps.Logger.Logger)
	creatorHandler := handler.NewCreatorHandler(deps.CreatorService, deps.I18n)
	categoryHandler := handler.NewCategoryHandler(deps.CategoryService)
	verificationHandler := handler.NewVerificationHandler(deps.VerificationService, deps.I18n)
	followHandler := handler.NewFollowHandler(deps.FollowService, deps.I18n)
	eventHandler := handler.NewEventHandler(deps.EventService, deps.AddressService, deps.TicketService, deps.InvitationService, deps.CreatorService, deps.I18n)
	addressHandler := handler.NewAddressHandler(deps.AddressService, deps.I18n)
	subscriptionHandler := handler.NewSubscriptionHandler(deps.SubscriptionService, deps.UserService, deps.StripeService, deps.I18n, deps.Logger)

	// Health check endpoint
	r.GET("/health", handler.HealthCheck(deps.DB))

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register/step1", authHandler.RegisterStep1)
			auth.POST("/login", authHandler.Login)
			auth.POST("/social-login", authHandler.SocialLogin)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// Public industry routes (no authentication required)
		industries := v1.Group("/industries")
		{
			industries.GET("", industryHandler.GetAllIndustries)
			industries.GET("/:id", industryHandler.GetIndustryByID)
			industries.GET("/slug/:slug", industryHandler.GetIndustryBySlug)
		}

		// Public category routes (no authentication required)
		categories := v1.Group("/categories")
		{
			categories.GET("/tree", categoryHandler.GetCategoryTree)
			categories.GET("/tree/type/:type", categoryHandler.GetCategoryTreeByType)
			categories.GET("/roots", categoryHandler.GetRootCategories)
			categories.GET("/leaves", categoryHandler.GetLeafCategories)
			categories.GET("/type/:type", categoryHandler.GetCategoriesByType)
			categories.GET("/search", categoryHandler.SearchCategories)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
			categories.GET("/:id/children", categoryHandler.GetCategoryChildren)
			categories.GET("/:id/parents", categoryHandler.GetCategoryParents)
		}

		// Public subscription plans (no authentication required)
		subscriptions := v1.Group("/subscriptions")
		{
			subscriptions.GET("/plans", subscriptionHandler.GetPlans)
			subscriptions.GET("/plans/subscriptions", subscriptionHandler.GetSubscriptionPlans)
			subscriptions.GET("/plans/packages", subscriptionHandler.GetPackagePlans)
		}

		// Public creator routes (no authentication required)
		creators := v1.Group("/creators")
		{
			creators.GET("", creatorHandler.GetCreatorList)
			creators.GET("/:id", creatorHandler.GetCreator)
		}

		// Username routes (require authentication)
		username := v1.Group("/username")
		username.Use(middleware.JWTAuth(deps.JWTService))
		{
			username.POST("/check", authHandler.CheckUsername)
			username.POST("/set", authHandler.SetUsername)
		}

		// Verification routes (require authentication) - only verify and resend, send is automatic in registration
		verification := v1.Group("/verification")
		verification.Use(middleware.JWTAuth(deps.JWTService))
		{
			verification.POST("/verify", verificationHandler.VerifyCode)
			verification.POST("/resend", verificationHandler.ResendVerification)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(deps.JWTService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.PUT("/contact", userHandler.UpdateContact)
				users.PUT("/profile-pic", userHandler.SetProfilePic)
				users.PUT("/cover-pic", userHandler.SetCoverPic)
				users.POST("/register/step4", authHandler.RegisterStep4)
				users.POST("/change-password", authHandler.ChangePassword)
				users.DELETE("/deactivate", userHandler.DeactivateAccount)
			}

			// Media routes
			media := protected.Group("/media")
			{
				media.POST("/upload", mediaHandler.UploadFile)
				media.GET("/:id", mediaHandler.GetMedia)
				media.GET("/user/:user_id", mediaHandler.GetUserMedia)
				media.PUT("/:id", mediaHandler.UpdateMedia)
				media.DELETE("/:id", mediaHandler.DeleteMedia)
			}

			// Creator routes (require authentication)
			creatorProtected := protected.Group("/creators")
			{
				creatorProtected.POST("", creatorHandler.CreateCreator)
				creatorProtected.GET("/me", creatorHandler.GetMyCreatorProfile)
				creatorProtected.PUT("/me", creatorHandler.UpdateCreator)
				creatorProtected.PUT("/me/weeztix-token", creatorHandler.SetWeeztixToken)
				creatorProtected.DELETE("/me", creatorHandler.DeleteCreator)
			}

			// Follow routes (require authentication)
			follows := protected.Group("/follows")
			{
				follows.POST("", followHandler.Follow)                          // Follow a user
				follows.DELETE("/:user_id", followHandler.Unfollow)             // Unfollow a user
				follows.GET("/status/:user_id", followHandler.GetFollowStatus)  // Check follow status
				follows.GET("/followers", followHandler.GetFollowers)           // Get my followers
				follows.GET("/following", followHandler.GetFollowing)           // Get who I follow
				follows.GET("/mutual/:user_id", followHandler.GetMutualFollows) // Get mutual follows
			}

			// Event management routes (require authentication and creator profile)
			eventManage := protected.Group("/events/manage")
			eventManage.Use(middleware.RequireUserType("creator"))
			{
				// Basic CRUD operations
				eventManage.POST("", eventHandler.CreateEvent)
				eventManage.GET("/:id", eventHandler.GetEvent)
				eventManage.PUT("/:id", eventHandler.UpdateEvent)
				eventManage.DELETE("/:id", eventHandler.DeleteEvent)

				// Creator-specific operations
				eventManage.GET("/my", eventHandler.GetMyEvents)
				eventManage.GET("/my/drafts", eventHandler.GetMyDraftEvents)
				eventManage.GET("/my/published", eventHandler.GetMyPublishedEvents)

				// Status management
				eventManage.PUT("/:id/status", eventHandler.UpdateEventStatus)
				eventManage.POST("/:id/submit", eventHandler.SubmitForReview)
				eventManage.POST("/:id/publish", eventHandler.PublishEvent)
				eventManage.POST("/:id/cancel", eventHandler.CancelEvent)

				// Statistics
				eventManage.GET("/stats", eventHandler.GetEventStats)
			}

			// Ticket management routes (separate to avoid route conflicts)
			tickets := protected.Group("/tickets")
			tickets.Use(middleware.RequireUserType("creator"))
			{
				tickets.POST("/event/:event_id", eventHandler.CreateTicket)
				tickets.GET("/event/:event_id", eventHandler.GetEventTickets)
			}

			// Invitation management routes (separate to avoid route conflicts)
			invitations := protected.Group("/invitations")
			invitations.Use(middleware.RequireUserType("creator"))
			{
				invitations.POST("/event/:event_id", eventHandler.CreateInvitation)
				invitations.GET("/event/:event_id", eventHandler.GetEventInvitations)
				invitations.PUT("/:invitation_id/respond", eventHandler.RespondToInvitation)
			}

			// Address management routes (require authentication and creator profile)
			addresses := protected.Group("/addresses")
			addresses.Use(middleware.RequireUserType("creator"))
			{
				addresses.POST("", addressHandler.CreateAddress)
				addresses.GET("/:id", addressHandler.GetAddress)
				addresses.GET("/search", addressHandler.SearchAddresses)
				addresses.GET("/city/:city", addressHandler.GetAddressesByCity)
			}

			// Subscription management routes (require authentication)
			subscriptionProtected := protected.Group("/subscriptions")
			{
				// User subscription management
				subscriptionProtected.GET("/my", subscriptionHandler.GetMySubscriptions)
				subscriptionProtected.GET("/publishing-rights", subscriptionHandler.GetPublishingRights)
				subscriptionProtected.GET("/stats", subscriptionHandler.GetUsageStats)
				subscriptionProtected.GET("/history", subscriptionHandler.GetSubscriptionHistory)

				// Purchase endpoints
				subscriptionProtected.POST("/purchase", subscriptionHandler.PurchaseSubscription)
				subscriptionProtected.POST("/packages/purchase", subscriptionHandler.PurchasePackage)

				// Subscription management
				subscriptionProtected.POST("/:id/cancel", subscriptionHandler.CancelSubscription)

				// Stripe webhook (no authentication required for webhooks)
				subscriptionProtected.POST("/webhook/stripe", subscriptionHandler.StripeWebhook)
			}
		}

		// Public event routes (no authentication required)
		publicEvents := v1.Group("/events")
		{
			publicEvents.GET("", eventHandler.GetPublicEvents)
			publicEvents.GET("/search", eventHandler.SearchEvents)
			publicEvents.GET("/location/:city", eventHandler.GetEventsByLocation)
			publicEvents.GET("/upcoming", eventHandler.GetUpcomingEvents)
			publicEvents.GET("/category/:category_id", eventHandler.GetEventsByCategory)
		}

		// Admin routes (require creator user type)
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(deps.JWTService))
		admin.Use(middleware.RequireUserType("creator"))
		{
			admin.GET("/users", userHandler.GetUserList)
			admin.GET("/media", mediaHandler.GetAllMedia)

			// Category cache management routes (admin only)
			adminCategories := admin.Group("/categories")
			{
				adminCategories.POST("/cache/refresh", categoryHandler.RefreshCache)
				adminCategories.DELETE("/cache/clear", categoryHandler.ClearCache)
			}
		}
	}
}
