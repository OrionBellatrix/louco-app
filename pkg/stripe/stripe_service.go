package stripe

import (
	"context"
	"fmt"
	"time"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/pkg/logger"
	"github.com/stripe/stripe-go/v76"
	checkoutsession "github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
)

type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
	Currency       string
	Environment    string
	SuccessURL     string
	CancelURL      string
	WebhookURL     string
}

type StripeService struct {
	config StripeConfig
	logger *logger.Logger
}

type CreateCustomerRequest struct {
	Email string
	Name  string
}

type CreateSubscriptionRequest struct {
	CustomerID string
	PriceID    string
	PlanID     uint
}

type CreatePaymentIntentRequest struct {
	Amount     int64 // in cents
	Currency   string
	CustomerID string
	PlanID     uint
}

type StripeCustomer struct {
	ID    string
	Email string
	Name  string
}

type StripeSubscription struct {
	ID                 string
	CustomerID         string
	Status             string
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	PriceID            string
	ClientSecret       string // For payment confirmation
}

type StripePaymentIntent struct {
	ID           string
	Amount       int64
	Currency     string
	Status       string
	ClientSecret string
}

type StripeCheckoutSession struct {
	ID  string
	URL string
}

func NewStripeService(config StripeConfig, logger *logger.Logger) *StripeService {
	stripe.Key = config.SecretKey
	return &StripeService{
		config: config,
		logger: logger,
	}
}

// Customer Management
func (s *StripeService) CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*StripeCustomer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(req.Email),
		Name:  stripe.String(req.Name),
	}

	c, err := customer.New(params)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create Stripe customer")
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	s.logger.Info().Str("customer_id", c.ID).Msg("Stripe customer created")

	return &StripeCustomer{
		ID:    c.ID,
		Email: c.Email,
		Name:  c.Name,
	}, nil
}

func (s *StripeService) GetCustomer(ctx context.Context, customerID string) (*StripeCustomer, error) {
	c, err := customer.Get(customerID, nil)
	if err != nil {
		s.logger.Error().Err(err).Str("customer_id", customerID).Msg("Failed to get Stripe customer")
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &StripeCustomer{
		ID:    c.ID,
		Email: c.Email,
		Name:  c.Name,
	}, nil
}

// Product and Price Management
func (s *StripeService) CreateProductAndPrice(ctx context.Context, plan *domain.SubscriptionPlan) (string, string, error) {
	// Create product
	productParams := &stripe.ProductParams{
		Name:        stripe.String(string(plan.Name)),
		Description: stripe.String(plan.Description),
		Metadata: map[string]string{
			"plan_id":   fmt.Sprintf("%d", plan.ID),
			"plan_type": string(plan.Type),
		},
	}

	prod, err := product.New(productParams)
	if err != nil {
		s.logger.Error().Err(err).Int("plan_id", plan.ID).Msg("Failed to create Stripe product")
		return "", "", fmt.Errorf("failed to create product: %w", err)
	}

	// Create price
	var priceParams *stripe.PriceParams
	if plan.Type == domain.SubscriptionTypeSubscription {
		// Recurring price for subscriptions
		var interval string
		if plan.BillingCycle != nil {
			switch *plan.BillingCycle {
			case domain.BillingCycleWeekly:
				interval = "week"
			case domain.BillingCycleMonthly:
				interval = "month"
			default:
				interval = "month"
			}
		} else {
			interval = "month"
		}

		priceParams = &stripe.PriceParams{
			Product:    stripe.String(prod.ID),
			UnitAmount: stripe.Int64(int64(plan.Price * 100)), // Convert to cents
			Currency:   stripe.String(s.config.Currency),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String(interval),
			},
			Metadata: map[string]string{
				"plan_id":   fmt.Sprintf("%d", plan.ID),
				"plan_type": string(plan.Type),
			},
		}
	} else {
		// One-time price for packages
		priceParams = &stripe.PriceParams{
			Product:    stripe.String(prod.ID),
			UnitAmount: stripe.Int64(int64(plan.Price * 100)), // Convert to cents
			Currency:   stripe.String(s.config.Currency),
			Metadata: map[string]string{
				"plan_id":   fmt.Sprintf("%d", plan.ID),
				"plan_type": string(plan.Type),
			},
		}
	}

	p, err := price.New(priceParams)
	if err != nil {
		s.logger.Error().Err(err).Int("plan_id", plan.ID).Msg("Failed to create Stripe price")
		return "", "", fmt.Errorf("failed to create price: %w", err)
	}

	s.logger.Info().
		Str("product_id", prod.ID).
		Str("price_id", p.ID).
		Int("plan_id", plan.ID).
		Msg("Stripe product and price created")

	return prod.ID, p.ID, nil
}

// Subscription Management
func (s *StripeService) CreateSubscription(ctx context.Context, req CreateSubscriptionRequest) (*StripeSubscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(req.CustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(req.PriceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String("on_subscription"),
		},
		Expand: []*string{
			stripe.String("latest_invoice.payment_intent"),
		},
		Metadata: map[string]string{
			"plan_id": fmt.Sprintf("%d", req.PlanID),
		},
	}

	sub, err := subscription.New(params)
	if err != nil {
		s.logger.Error().Err(err).
			Str("customer_id", req.CustomerID).
			Str("price_id", req.PriceID).
			Msg("Failed to create Stripe subscription")
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	s.logger.Info().
		Str("subscription_id", sub.ID).
		Str("customer_id", req.CustomerID).
		Msg("Stripe subscription created")

	// Extract client secret from latest invoice payment intent
	var clientSecret string
	if sub.LatestInvoice != nil && sub.LatestInvoice.PaymentIntent != nil {
		clientSecret = sub.LatestInvoice.PaymentIntent.ClientSecret
	}

	return &StripeSubscription{
		ID:                 sub.ID,
		CustomerID:         sub.Customer.ID,
		Status:             string(sub.Status),
		CurrentPeriodStart: time.Unix(sub.CurrentPeriodStart, 0),
		CurrentPeriodEnd:   time.Unix(sub.CurrentPeriodEnd, 0),
		PriceID:            req.PriceID,
		ClientSecret:       clientSecret,
	}, nil
}

func (s *StripeService) CancelSubscription(ctx context.Context, subscriptionID string) error {
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err := subscription.Update(subscriptionID, params)
	if err != nil {
		s.logger.Error().Err(err).Str("subscription_id", subscriptionID).Msg("Failed to cancel Stripe subscription")
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	s.logger.Info().Str("subscription_id", subscriptionID).Msg("Stripe subscription cancelled")
	return nil
}

func (s *StripeService) GetSubscription(ctx context.Context, subscriptionID string) (*StripeSubscription, error) {
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		s.logger.Error().Err(err).Str("subscription_id", subscriptionID).Msg("Failed to get Stripe subscription")
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &StripeSubscription{
		ID:                 sub.ID,
		CustomerID:         sub.Customer.ID,
		Status:             string(sub.Status),
		CurrentPeriodStart: time.Unix(sub.CurrentPeriodStart, 0),
		CurrentPeriodEnd:   time.Unix(sub.CurrentPeriodEnd, 0),
	}, nil
}

// Payment Intent Management (for one-time payments/packages)
func (s *StripeService) CreatePaymentIntent(ctx context.Context, req CreatePaymentIntentRequest) (*StripePaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(req.Currency),
		Customer: stripe.String(req.CustomerID),
		Metadata: map[string]string{
			"plan_id": fmt.Sprintf("%d", req.PlanID),
			"type":    "package_purchase",
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		s.logger.Error().Err(err).
			Str("customer_id", req.CustomerID).
			Int64("amount", req.Amount).
			Msg("Failed to create Stripe payment intent")
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	s.logger.Info().
		Str("payment_intent_id", pi.ID).
		Str("customer_id", req.CustomerID).
		Int64("amount", req.Amount).
		Msg("Stripe payment intent created")

	return &StripePaymentIntent{
		ID:           pi.ID,
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
		Status:       string(pi.Status),
		ClientSecret: pi.ClientSecret,
	}, nil
}

func (s *StripeService) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*StripePaymentIntent, error) {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		s.logger.Error().Err(err).Str("payment_intent_id", paymentIntentID).Msg("Failed to get Stripe payment intent")
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	return &StripePaymentIntent{
		ID:           pi.ID,
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
		Status:       string(pi.Status),
		ClientSecret: pi.ClientSecret,
	}, nil
}

// Webhook signature verification
func (s *StripeService) VerifyWebhookSignature(payload []byte, signature string) error {
	// Skip signature verification in development mode for testing
	if s.config.Environment == "development" {
		s.logger.Warn().Msg("Skipping webhook signature verification in development mode")
		return nil
	}

	_, err := webhook.ConstructEvent(payload, signature, s.config.WebhookSecret)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to verify Stripe webhook signature")
		return fmt.Errorf("failed to verify webhook signature: %w", err)
	}
	return nil
}

// Checkout Session Management
func (s *StripeService) CreateCheckoutSessionForSubscription(ctx context.Context, plan *domain.SubscriptionPlan, customerEmail, customerName string, userID uint) (*StripeCheckoutSession, error) {
	// Create product and price first
	_, priceID, err := s.CreateProductAndPrice(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("failed to create product and price: %w", err)
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:          stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL:    stripe.String(s.config.SuccessURL),
		CancelURL:     stripe.String(s.config.CancelURL),
		CustomerEmail: stripe.String(customerEmail),
		Metadata: map[string]string{
			"plan_id": fmt.Sprintf("%d", plan.ID),
			"user_id": fmt.Sprintf("%d", userID),
			"type":    "subscription",
		},
	}

	sess, err := checkoutsession.New(params)
	if err != nil {
		s.logger.Error().Err(err).Int("plan_id", plan.ID).Msg("Failed to create Stripe checkout session for subscription")
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	s.logger.Info().
		Str("session_id", sess.ID).
		Str("url", sess.URL).
		Int("plan_id", plan.ID).
		Msg("Stripe checkout session created for subscription")

	return &StripeCheckoutSession{
		ID:  sess.ID,
		URL: sess.URL,
	}, nil
}

func (s *StripeService) CreateCheckoutSessionForPackage(ctx context.Context, plan *domain.SubscriptionPlan, customerEmail, customerName string, userID uint) (*StripeCheckoutSession, error) {
	// Create product and price first
	_, priceID, err := s.CreateProductAndPrice(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("failed to create product and price: %w", err)
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:    stripe.String(s.config.SuccessURL),
		CancelURL:     stripe.String(s.config.CancelURL),
		CustomerEmail: stripe.String(customerEmail),
		Metadata: map[string]string{
			"plan_id": fmt.Sprintf("%d", plan.ID),
			"user_id": fmt.Sprintf("%d", userID),
			"type":    "package",
		},
	}

	sess, err := checkoutsession.New(params)
	if err != nil {
		s.logger.Error().Err(err).Int("plan_id", plan.ID).Msg("Failed to create Stripe checkout session for package")
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	s.logger.Info().
		Str("session_id", sess.ID).
		Str("url", sess.URL).
		Int("plan_id", plan.ID).
		Msg("Stripe checkout session created for package")

	return &StripeCheckoutSession{
		ID:  sess.ID,
		URL: sess.URL,
	}, nil
}

// Helper methods
func (s *StripeService) ConvertDollarsToCents(dollars float64) int64 {
	return int64(dollars * 100)
}

func (s *StripeService) ConvertCentsToDollars(cents int64) float64 {
	return float64(cents) / 100
}
