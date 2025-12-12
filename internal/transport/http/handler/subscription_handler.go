package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/i18n"
	"github.com/louco-event/internal/service"
	"github.com/louco-event/pkg/logger"
	"github.com/louco-event/pkg/stripe"
)

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
	userService         service.UserService
	stripeService       *stripe.StripeService
	i18n                *i18n.I18n
	logger              *logger.Logger
}

func NewSubscriptionHandler(
	subscriptionService service.SubscriptionService,
	userService service.UserService,
	stripeService *stripe.StripeService,
	i18n *i18n.I18n,
	logger *logger.Logger,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		userService:         userService,
		stripeService:       stripeService,
		i18n:                i18n,
		logger:              logger,
	}
}

// GetPlans godoc
// @Summary Get available subscription and package plans
// @Description Get all available subscription and package plans
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.PlansResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/plans [get]
func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
	lang := c.GetString("lang")

	// Get subscription plans
	subscriptionPlans, err := h.subscriptionService.GetSubscriptionPlans(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get subscription plans")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.plans.fetch_failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Get package plans
	packagePlans, err := h.subscriptionService.GetPackagePlans(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get package plans")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.plans.fetch_failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	response := dto.ToPlansResponse(subscriptionPlans, packagePlans)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: h.i18n.Translate(lang, "subscription.plans.fetch_success"),
		Data:    response,
		Errors:  nil,
	})
}

// GetSubscriptionPlans godoc
// @Summary Get subscription plans only
// @Description Get all available subscription plans
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]dto.SubscriptionPlanResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/plans/subscriptions [get]
func (h *SubscriptionHandler) GetSubscriptionPlans(c *gin.Context) {
	lang := c.GetString("lang")

	plans, err := h.subscriptionService.GetSubscriptionPlans(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get subscription plans")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.plans.fetch_failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	response := dto.ToSubscriptionPlanResponses(plans)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: h.i18n.Translate(lang, "subscription.plans.fetch_success"),
		Data:    response,
		Errors:  nil,
	})
}

// GetPackagePlans godoc
// @Summary Get package plans only
// @Description Get all available package plans
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]dto.SubscriptionPlanResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/plans/packages [get]
func (h *SubscriptionHandler) GetPackagePlans(c *gin.Context) {
	lang := c.GetString("lang")

	plans, err := h.subscriptionService.GetPackagePlans(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get package plans")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.plans.fetch_failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	response := dto.ToSubscriptionPlanResponses(plans)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: h.i18n.Translate(lang, "subscription.plans.fetch_success"),
		Data:    response,
		Errors:  nil,
	})
}

// GetMySubscriptions godoc
// @Summary Get user's subscriptions
// @Description Get all subscriptions and packages for the authenticated user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=[]dto.UserSubscriptionResponse}
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/my [get]
func (h *SubscriptionHandler) GetMySubscriptions(c *gin.Context) {
	lang := c.GetString("lang")
	userIDInt, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "auth.token_required"),
			Data:    nil,
			Errors:  []string{"User ID not found in token"},
		})
		return
	}

	userID := uint(userIDInt.(int))

	subscriptions, err := h.subscriptionService.GetUserSubscriptions(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get user subscriptions")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.fetch_failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	response := dto.ToUserSubscriptionResponses(subscriptions)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: h.i18n.Translate(lang, "subscription.fetch_success"),
		Data:    response,
		Errors:  nil,
	})
}

// GetPublishingRights godoc
// @Summary Get user's publishing rights
// @Description Get current publishing rights and limits for the authenticated user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.PublishingRightsResponse}
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/publishing-rights [get]
func (h *SubscriptionHandler) GetPublishingRights(c *gin.Context) {
	lang := c.GetString("lang")
	userIDInt, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "auth.token_required"),
			Data:    nil,
			Errors:  []string{"User ID not found in token"},
		})
		return
	}

	userID := uint(userIDInt.(int))

	rights, err := h.subscriptionService.GetPublishingRights(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get publishing rights")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.rights.fetch_failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	response := dto.ToPublishingRightsResponse(rights)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: h.i18n.Translate(lang, "subscription.rights.fetch_success"),
		Data:    response,
		Errors:  nil,
	})
}

// GetUsageStats godoc
// @Summary Get user's usage statistics
// @Description Get detailed usage statistics for the authenticated user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.SubscriptionStatsResponse}
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/stats [get]
func (h *SubscriptionHandler) GetUsageStats(c *gin.Context) {
	lang := c.GetString("lang")
	userIDInt, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "auth.token_required"),
			Data:    nil,
			Errors:  []string{"User ID not found in token"},
		})
		return
	}

	userID := uint(userIDInt.(int))

	stats, err := h.subscriptionService.GetUsageStats(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get usage stats")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.stats.fetch_failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	response := dto.ToSubscriptionStatsResponse(stats)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: h.i18n.Translate(lang, "subscription.stats.fetch_success"),
		Data:    response,
		Errors:  nil,
	})
}

// GetSubscriptionHistory godoc
// @Summary Get user's subscription history
// @Description Get paginated subscription history for the authenticated user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param type query string false "Filter by type (subscription/package)"
// @Param status query string false "Filter by status (active/cancelled/expired/pending)"
// @Success 200 {object} dto.APIResponse{data=dto.SubscriptionHistoryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/history [get]
func (h *SubscriptionHandler) GetSubscriptionHistory(c *gin.Context) {
	lang := c.GetString("lang")
	userIDInt, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "auth.token_required"),
			Data:    nil,
			Errors:  []string{"User ID not found in token"},
		})
		return
	}

	userID := uint(userIDInt.(int))

	// Parse pagination parameters
	var req dto.SubscriptionHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "validation.invalid_request"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	limit := 10
	if req.PaginationRequest.PageSize > 0 {
		limit = req.PaginationRequest.PageSize
	}

	offset := (req.Page - 1) * limit

	subscriptions, total, err := h.subscriptionService.GetSubscriptionHistory(c.Request.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get subscription history")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.history.fetch_failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	response := &dto.SubscriptionHistoryResponse{
		Subscriptions: dto.ToUserSubscriptionResponses(subscriptions),
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			PageSize:   limit,
			Total:      int(total),
			TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		},
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: h.i18n.Translate(lang, "subscription.history.fetch_success"),
		Data:    response,
		Errors:  nil,
	})
}

// PurchaseSubscription godoc
// @Summary Purchase a subscription
// @Description Purchase a subscription plan with Stripe payment
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.PurchaseSubscriptionRequest true "Purchase subscription request"
// @Success 200 {object} dto.APIResponse{data=dto.PurchaseResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/purchase [post]
func (h *SubscriptionHandler) PurchaseSubscription(c *gin.Context) {
	lang := c.GetString("lang")
	userIDInt, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "auth.token_required"),
			Data:    nil,
			Errors:  []string{"User ID not found in token"},
		})
		return
	}

	userID := uint(userIDInt.(int))

	var req dto.PurchaseSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "validation.invalid_request"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Validate request
	if req.PlanID == 0 {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "validation.failed"),
			Data:    nil,
			Errors:  []string{"Plan ID is required"},
		})
		return
	}
	if req.PaymentMethodID == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "validation.failed"),
			Data:    nil,
			Errors:  []string{"Payment method ID is required"},
		})
		return
	}

	// Get subscription plan
	plan, err := h.subscriptionService.GetPlanByID(c.Request.Context(), req.PlanID)
	if err != nil {
		h.logger.Error().Err(err).Uint("plan_id", req.PlanID).Msg("Failed to get subscription plan")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.plan.not_found"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Get user details for customer creation
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get user details")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.purchase.failed"),
			Data:    nil,
			Errors:  []string{"Failed to get user details"},
		})
		return
	}

	// Get user email for checkout session
	var email string
	if user.Email != nil {
		email = *user.Email
	}

	// Create Stripe Checkout Session for subscription
	checkoutSession, err := h.stripeService.CreateCheckoutSessionForSubscription(c.Request.Context(), plan, email, user.FullName, userID)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", userID).Uint("plan_id", req.PlanID).Msg("Failed to create Stripe checkout session")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.purchase.failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	response := &dto.PurchaseResponse{
		SubscriptionID: 0, // Will be created after successful payment
		Status:         "pending_payment",
		RequiresAction: true,
		ClientSecret:   checkoutSession.URL, // Return checkout URL instead of client secret
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: h.i18n.Translate(lang, "subscription.purchase.success"),
		Data:    response,
		Errors:  nil,
	})
}

// PurchasePackage godoc
// @Summary Purchase a package
// @Description Purchase a package plan with Stripe payment
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.PurchasePackageRequest true "Purchase package request"
// @Success 200 {object} dto.APIResponse{data=dto.PurchaseResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/packages/purchase [post]
func (h *SubscriptionHandler) PurchasePackage(c *gin.Context) {
	lang := c.GetString("lang")
	userIDInt, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "auth.token_required"),
			Data:    nil,
			Errors:  []string{"User ID not found in token"},
		})
		return
	}

	userID := uint(userIDInt.(int))

	var req dto.PurchasePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "validation.invalid_request"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Validate request
	if req.PlanID == 0 {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "validation.failed"),
			Data:    nil,
			Errors:  []string{"Plan ID is required"},
		})
		return
	}
	if req.PaymentMethodID == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "validation.failed"),
			Data:    nil,
			Errors:  []string{"Payment method ID is required"},
		})
		return
	}

	// Get package plan
	plan, err := h.subscriptionService.GetPlanByID(c.Request.Context(), req.PlanID)
	if err != nil {
		h.logger.Error().Err(err).Uint("plan_id", req.PlanID).Msg("Failed to get package plan")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.plan.not_found"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Get user details for customer creation
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", userID).Msg("Failed to get user details")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.package.purchase.failed"),
			Data:    nil,
			Errors:  []string{"Failed to get user details"},
		})
		return
	}

	// Get user email for checkout session
	var email string
	if user.Email != nil {
		email = *user.Email
	}

	// Create Stripe Checkout Session for package
	checkoutSession, err := h.stripeService.CreateCheckoutSessionForPackage(c.Request.Context(), plan, email, user.FullName, userID)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", userID).Uint("plan_id", req.PlanID).Msg("Failed to create Stripe checkout session")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.package.purchase.failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	response := &dto.PurchaseResponse{
		SubscriptionID: 0, // Will be created after successful payment
		Status:         "pending_payment",
		RequiresAction: true,
		ClientSecret:   checkoutSession.URL, // Return checkout URL instead of client secret
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: h.i18n.Translate(lang, "subscription.package.purchase.success"),
		Data:    response,
		Errors:  nil,
	})
}

// CancelSubscription godoc
// @Summary Cancel a subscription
// @Description Cancel an active subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subscription ID"
// @Param request body dto.CancelSubscriptionRequest false "Cancel subscription request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/{id}/cancel [post]
func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	lang := c.GetString("lang")

	subscriptionIDStr := c.Param("id")
	subscriptionID, err := strconv.ParseUint(subscriptionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "validation.invalid_id"),
			Data:    nil,
			Errors:  []string{"Invalid subscription ID"},
		})
		return
	}

	var req dto.CancelSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Reason is optional, so we can ignore binding errors
	}

	if err := h.subscriptionService.CancelSubscription(c.Request.Context(), uint(subscriptionID)); err != nil {
		h.logger.Error().Err(err).Uint("subscription_id", uint(subscriptionID)).Msg("Failed to cancel subscription")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: h.i18n.Translate(lang, "subscription.cancel.failed"),
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: h.i18n.Translate(lang, "subscription.cancel.success"),
		Data:    nil,
		Errors:  nil,
	})
}

// StripeWebhook godoc
// @Summary Handle Stripe webhooks
// @Description Handle Stripe webhook events for payment processing
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.StripeWebhookRequest true "Stripe webhook request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /subscriptions/webhook/stripe [post]
func (h *SubscriptionHandler) StripeWebhook(c *gin.Context) {
	// Get raw body for signature verification
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read webhook body")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to read request body",
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Verify webhook signature
	signature := c.GetHeader("Stripe-Signature")
	if err := h.stripeService.VerifyWebhookSignature(body, signature); err != nil {
		h.logger.Error().Err(err).Msg("Failed to verify webhook signature")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid webhook signature",
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Parse webhook request
	var req dto.StripeWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to parse webhook payload")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid webhook payload",
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Handle different event types
	switch req.Type {
	case "checkout.session.completed":
		// Handle successful checkout session completion (both subscriptions and packages)
		h.logger.Info().Str("event_type", req.Type).Msg("Handling Stripe webhook: checkout.session.completed")
		if err := h.handleCheckoutSessionCompleted(c.Request.Context(), req.Data); err != nil {
			h.logger.Error().Err(err).Msg("Failed to handle checkout.session.completed")
		}

	case "payment_intent.succeeded":
		// Handle successful payment for packages (legacy support)
		h.logger.Info().Str("event_type", req.Type).Msg("Handling Stripe webhook: payment_intent.succeeded")
		if err := h.handlePaymentIntentSucceeded(c.Request.Context(), req.Data); err != nil {
			h.logger.Error().Err(err).Msg("Failed to handle payment_intent.succeeded")
		}

	case "invoice.payment_succeeded":
		// Handle successful subscription payment (legacy support)
		h.logger.Info().Str("event_type", req.Type).Msg("Handling Stripe webhook: invoice.payment_succeeded")
		if err := h.handleInvoicePaymentSucceeded(c.Request.Context(), req.Data); err != nil {
			h.logger.Error().Err(err).Msg("Failed to handle invoice.payment_succeeded")
		}

	case "invoice.payment_failed":
		// Handle failed subscription payment
		h.logger.Info().Str("event_type", req.Type).Msg("Handling Stripe webhook: invoice.payment_failed")
		if err := h.handleInvoicePaymentFailed(c.Request.Context(), req.Data); err != nil {
			h.logger.Error().Err(err).Msg("Failed to handle invoice.payment_failed")
		}

	case "customer.subscription.deleted":
		// Handle subscription cancellation
		h.logger.Info().Str("event_type", req.Type).Msg("Handling Stripe webhook: customer.subscription.deleted")
		if err := h.handleSubscriptionDeleted(c.Request.Context(), req.Data); err != nil {
			h.logger.Error().Err(err).Msg("Failed to handle customer.subscription.deleted")
		}

	default:
		h.logger.Info().Str("event_type", req.Type).Msg("Unhandled Stripe webhook event")
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Webhook processed",
		Data:    nil,
		Errors:  nil,
	})
}

// Helper methods for webhook event handling
func (h *SubscriptionHandler) handleCheckoutSessionCompleted(ctx context.Context, data interface{}) error {
	h.logger.Info().Msg("Processing checkout.session.completed webhook")

	// Parse the checkout session data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to marshal checkout session data")
		return err
	}

	var checkoutSessionData struct {
		Object struct {
			ID           string            `json:"id"`
			Mode         string            `json:"mode"`
			Subscription string            `json:"subscription"`
			Metadata     map[string]string `json:"metadata"`
		} `json:"object"`
	}

	if err := json.Unmarshal(dataBytes, &checkoutSessionData); err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal checkout session data")
		return err
	}

	sessionID := checkoutSessionData.Object.ID
	mode := checkoutSessionData.Object.Mode
	userIDStr, exists := checkoutSessionData.Object.Metadata["user_id"]
	if !exists {
		h.logger.Error().Str("session_id", sessionID).Msg("User ID not found in checkout session metadata")
		return nil
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userIDStr).Msg("Invalid user ID in checkout session metadata")
		return err
	}

	planIDStr, exists := checkoutSessionData.Object.Metadata["plan_id"]
	if !exists {
		h.logger.Error().Str("session_id", sessionID).Msg("Plan ID not found in checkout session metadata")
		return nil
	}

	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("plan_id", planIDStr).Msg("Invalid plan ID in checkout session metadata")
		return err
	}

	// Handle based on mode (subscription or payment)
	switch mode {
	case "subscription":
		// Handle subscription activation
		subscriptionID := checkoutSessionData.Object.Subscription
		if subscriptionID == "" {
			h.logger.Error().Str("session_id", sessionID).Msg("Subscription ID not found in checkout session")
			return nil
		}

		// Create subscription in our database
		if _, err := h.subscriptionService.CreateSubscription(ctx, uint(userID), uint(planID), subscriptionID); err != nil {
			h.logger.Error().Err(err).Uint("user_id", uint(userID)).Uint("plan_id", uint(planID)).Str("stripe_subscription_id", subscriptionID).Msg("Failed to create subscription")
			return err
		}

		h.logger.Info().Str("session_id", sessionID).Uint("user_id", uint(userID)).Uint("plan_id", uint(planID)).Str("stripe_subscription_id", subscriptionID).Msg("Subscription created successfully")

	case "payment":
		// Handle package activation
		// Create package in our database
		if _, err := h.subscriptionService.CreatePackage(ctx, uint(userID), uint(planID), sessionID); err != nil {
			h.logger.Error().Err(err).Uint("user_id", uint(userID)).Uint("plan_id", uint(planID)).Str("stripe_session_id", sessionID).Msg("Failed to create package")
			return err
		}

		h.logger.Info().Str("session_id", sessionID).Uint("user_id", uint(userID)).Uint("plan_id", uint(planID)).Msg("Package created successfully")

	default:
		h.logger.Warn().Str("session_id", sessionID).Str("mode", mode).Msg("Unknown checkout session mode")
	}

	return nil
}

func (h *SubscriptionHandler) handlePaymentIntentSucceeded(ctx context.Context, data interface{}) error {
	h.logger.Info().Msg("Processing payment_intent.succeeded webhook")

	// Parse the payment intent data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to marshal payment intent data")
		return err
	}

	var paymentIntentData struct {
		Object struct {
			ID       string            `json:"id"`
			Metadata map[string]string `json:"metadata"`
		} `json:"object"`
	}

	if err := json.Unmarshal(dataBytes, &paymentIntentData); err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal payment intent data")
		return err
	}

	paymentIntentID := paymentIntentData.Object.ID
	planIDStr, exists := paymentIntentData.Object.Metadata["plan_id"]
	if !exists {
		h.logger.Error().Str("payment_intent_id", paymentIntentID).Msg("Plan ID not found in payment intent metadata")
		return nil
	}

	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("plan_id", planIDStr).Msg("Invalid plan ID in payment intent metadata")
		return err
	}

	// Find the user subscription by Stripe payment intent ID and activate it
	if err := h.subscriptionService.ActivateSubscriptionByStripeID(ctx, paymentIntentID); err != nil {
		h.logger.Error().Err(err).Str("payment_intent_id", paymentIntentID).Uint("plan_id", uint(planID)).Msg("Failed to activate package")
		return err
	}

	h.logger.Info().Str("payment_intent_id", paymentIntentID).Uint("plan_id", uint(planID)).Msg("Package activated successfully")
	return nil
}

func (h *SubscriptionHandler) handleInvoicePaymentSucceeded(ctx context.Context, data interface{}) error {
	h.logger.Info().Msg("Processing invoice.payment_succeeded webhook")

	// Parse the invoice data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to marshal invoice data")
		return err
	}

	var invoiceData struct {
		Object struct {
			ID           string `json:"id"`
			Subscription string `json:"subscription"`
		} `json:"object"`
	}

	if err := json.Unmarshal(dataBytes, &invoiceData); err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal invoice data")
		return err
	}

	subscriptionID := invoiceData.Object.Subscription
	if subscriptionID == "" {
		h.logger.Error().Str("invoice_id", invoiceData.Object.ID).Msg("Subscription ID not found in invoice")
		return nil
	}

	// Find the user subscription by Stripe subscription ID and activate it
	if err := h.subscriptionService.ActivateSubscriptionByStripeID(ctx, subscriptionID); err != nil {
		h.logger.Error().Err(err).Str("subscription_id", subscriptionID).Msg("Failed to activate subscription")
		return err
	}

	h.logger.Info().Str("subscription_id", subscriptionID).Msg("Subscription activated successfully")
	return nil
}

func (h *SubscriptionHandler) handleInvoicePaymentFailed(ctx context.Context, data interface{}) error {
	h.logger.Info().Msg("Processing invoice.payment_failed webhook")

	// Parse the invoice data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to marshal invoice data")
		return err
	}

	var invoiceData struct {
		Object struct {
			ID           string `json:"id"`
			Subscription string `json:"subscription"`
		} `json:"object"`
	}

	if err := json.Unmarshal(dataBytes, &invoiceData); err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal invoice data")
		return err
	}

	subscriptionID := invoiceData.Object.Subscription
	if subscriptionID == "" {
		h.logger.Error().Str("invoice_id", invoiceData.Object.ID).Msg("Subscription ID not found in invoice")
		return nil
	}

	// Handle payment failure - could suspend subscription or retry payment
	h.logger.Warn().Str("subscription_id", subscriptionID).Msg("Subscription payment failed")

	// For now, just log the failure. In a real implementation, you might:
	// - Suspend the subscription
	// - Send notification to user
	// - Retry payment after some time
	// - Cancel subscription after multiple failures

	return nil
}

func (h *SubscriptionHandler) handleSubscriptionDeleted(ctx context.Context, data interface{}) error {
	h.logger.Info().Msg("Processing customer.subscription.deleted webhook")

	// Parse the subscription data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to marshal subscription data")
		return err
	}

	var subscriptionData struct {
		Object struct {
			ID string `json:"id"`
		} `json:"object"`
	}

	if err := json.Unmarshal(dataBytes, &subscriptionData); err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal subscription data")
		return err
	}

	subscriptionID := subscriptionData.Object.ID
	if subscriptionID == "" {
		h.logger.Error().Msg("Subscription ID not found in subscription data")
		return nil
	}

	// Find the user subscription by Stripe subscription ID and cancel it
	if err := h.subscriptionService.CancelSubscriptionByStripeID(ctx, subscriptionID); err != nil {
		h.logger.Error().Err(err).Str("subscription_id", subscriptionID).Msg("Failed to cancel subscription")
		return err
	}

	h.logger.Info().Str("subscription_id", subscriptionID).Msg("Subscription cancelled successfully")
	return nil
}
