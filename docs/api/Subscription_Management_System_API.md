# Subscription & Package Management System API Documentation

## Overview

The Subscription & Package Management System provides comprehensive event publishing rights management through subscription plans and one-time packages. It integrates with Stripe for payment processing and includes automatic usage tracking and validation.

## Table of Contents

1. [System Architecture](#system-architecture)
2. [Authentication](#authentication)
3. [Subscription Plans](#subscription-plans)
4. [API Endpoints](#api-endpoints)
5. [Event Publishing Integration](#event-publishing-integration)
6. [Stripe Integration](#stripe-integration)
7. [Error Handling](#error-handling)
8. [Usage Examples](#usage-examples)
9. [Testing Scenarios](#testing-scenarios)

## System Architecture

### Subscription Types

#### 1. Recurring Subscriptions
- **Basic Plan**: 5 events/week or 20 events/month
- **Plus Plan**: 15 events/week or 60 events/month  
- **Pro Plan**: Unlimited events
- **Billing Cycles**: Weekly (Monday reset) or Monthly (1st day reset)

#### 2. One-time Packages
- **1 Event Package**: Single event credit
- **10 Events Package**: 10 event credits
- **25 Events Package**: 25 event credits
- **No Expiration**: Credits never expire until used

### Publishing Rights Validation

The system validates publishing rights before allowing event submissions:

1. **Check Active Subscription/Package**: User must have at least one active subscription or package
2. **Validate Limits/Credits**: Ensure sufficient remaining limits or credits
3. **Consume Usage**: Automatically deduct from limits/credits after successful submission
4. **Error Handling**: Return user-friendly error messages for insufficient rights

## Authentication

All endpoints (except webhooks) require Bearer token authentication:

```
Authorization: Bearer {access_token}
```

Language preference can be set via header:

```
Accept-Language: en|tr
```

## Subscription Plans

### Plan Structure

```json
{
  "id": 1,
  "name": "Basic",
  "description": "Perfect for small events",
  "plan_type": "subscription",
  "billing_cycle": "monthly",
  "price": 29.99,
  "currency": "USD",
  "weekly_limit": 5,
  "monthly_limit": 20,
  "credits": null,
  "is_active": true,
  "stripe_price_id": "price_basic_monthly"
}
```

### Available Plans

| ID | Name | Type | Weekly Limit | Monthly Limit | Credits | Price |
|----|------|------|--------------|---------------|---------|-------|
| 1 | Basic | subscription | 5 | 20 | - | $29.99/month |
| 2 | Plus | subscription | 15 | 60 | - | $79.99/month |
| 3 | Pro | subscription | unlimited | unlimited | - | $199.99/month |
| 4 | 1 Event | package | - | - | 1 | $9.99 |
| 5 | 10 Events | package | - | - | 10 | $79.99 |
| 6 | 25 Events | package | - | - | 25 | $179.99 |

## API Endpoints

### Subscription Plans

#### GET /api/v1/subscriptions/plans
Get all available subscription plans.

**Response:**
```json
{
  "success": true,
  "message": "subscription.plan.list.success",
  "data": [
    {
      "id": 1,
      "name": "Basic",
      "description": "Perfect for small events",
      "plan_type": "subscription",
      "billing_cycle": "monthly",
      "price": 29.99,
      "currency": "USD",
      "weekly_limit": 5,
      "monthly_limit": 20,
      "credits": null,
      "is_active": true
    }
  ]
}
```

#### GET /api/v1/subscriptions/plans/{id}
Get specific subscription plan details.

**Parameters:**
- `id` (path): Plan ID

**Response:**
```json
{
  "success": true,
  "message": "subscription.plan.get.success",
  "data": {
    "id": 1,
    "name": "Basic",
    "description": "Perfect for small events",
    "plan_type": "subscription",
    "billing_cycle": "monthly",
    "price": 29.99,
    "currency": "USD",
    "weekly_limit": 5,
    "monthly_limit": 20,
    "credits": null,
    "is_active": true
  }
}
```

### User Subscriptions

#### GET /api/v1/subscriptions/my
Get current user's active subscriptions and packages.

**Response:**
```json
{
  "success": true,
  "message": "subscription.user.list.success",
  "data": [
    {
      "id": 1,
      "user_id": 123,
      "plan_id": 1,
      "plan_name": "Basic",
      "plan_type": "subscription",
      "billing_cycle": "monthly",
      "status": "active",
      "current_period_start": "2024-01-01T00:00:00Z",
      "current_period_end": "2024-02-01T00:00:00Z",
      "weekly_limit": 5,
      "monthly_limit": 20,
      "weekly_usage": 2,
      "monthly_usage": 8,
      "credits": null,
      "credits_used": null,
      "stripe_subscription_id": "sub_stripe_id",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### GET /api/v1/subscriptions/publishing-rights
Check current user's event publishing rights.

**Response:**
```json
{
  "success": true,
  "message": "subscription.rights.success",
  "data": {
    "can_publish": true,
    "active_subscriptions": 1,
    "active_packages": 0,
    "weekly_remaining": 3,
    "monthly_remaining": 12,
    "total_credits": 0,
    "used_credits": 0,
    "remaining_credits": 0,
    "next_reset": "2024-02-01T00:00:00Z"
  }
}
```

#### GET /api/v1/subscriptions/usage-stats
Get detailed usage statistics.

**Response:**
```json
{
  "success": true,
  "message": "subscription.stats.success",
  "data": {
    "current_period": {
      "start": "2024-01-01T00:00:00Z",
      "end": "2024-02-01T00:00:00Z",
      "weekly_usage": 2,
      "monthly_usage": 8,
      "events_published": 8
    },
    "subscriptions": [
      {
        "id": 1,
        "plan_name": "Basic",
        "status": "active",
        "weekly_limit": 5,
        "monthly_limit": 20,
        "weekly_usage": 2,
        "monthly_usage": 8
      }
    ],
    "packages": [],
    "total_events_published": 8,
    "total_credits_purchased": 0,
    "total_credits_used": 0
  }
}
```

### Subscription Purchase

#### POST /api/v1/subscriptions/purchase
Purchase a subscription plan.

**Request Body:**
```json
{
  "plan_id": 1,
  "billing_cycle": "monthly",
  "payment_method_id": "pm_card_visa"
}
```

**Parameters:**
- `plan_id`: Subscription plan ID (1-3)
- `billing_cycle`: "weekly" or "monthly"
- `payment_method_id`: Stripe payment method ID

**Response:**
```json
{
  "success": true,
  "message": "subscription.purchase.success",
  "data": {
    "id": 1,
    "user_id": 123,
    "plan_id": 1,
    "status": "active",
    "stripe_subscription_id": "sub_stripe_id",
    "stripe_customer_id": "cus_stripe_id",
    "current_period_start": "2024-01-01T00:00:00Z",
    "current_period_end": "2024-02-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### POST /api/v1/subscriptions/purchase-package
Purchase a one-time event package.

**Request Body:**
```json
{
  "plan_id": 5,
  "payment_method_id": "pm_card_visa"
}
```

**Parameters:**
- `plan_id`: Package plan ID (4-6)
- `payment_method_id`: Stripe payment method ID

**Response:**
```json
{
  "success": true,
  "message": "subscription.purchase.success",
  "data": {
    "id": 2,
    "user_id": 123,
    "plan_id": 5,
    "status": "active",
    "credits": 10,
    "credits_used": 0,
    "stripe_payment_intent_id": "pi_stripe_id",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Subscription Management

#### POST /api/v1/subscriptions/cancel
Cancel an active subscription.

**Request Body:**
```json
{
  "subscription_id": 1,
  "reason": "No longer needed"
}
```

**Response:**
```json
{
  "success": true,
  "message": "subscription.user.cancel.success",
  "data": {
    "id": 1,
    "status": "cancelled",
    "cancelled_at": "2024-01-15T00:00:00Z",
    "cancellation_reason": "No longer needed"
  }
}
```

#### PUT /api/v1/subscriptions/update
Update an existing subscription (upgrade/downgrade).

**Request Body:**
```json
{
  "subscription_id": 1,
  "new_plan_id": 2,
  "billing_cycle": "monthly"
}
```

**Response:**
```json
{
  "success": true,
  "message": "subscription.user.update.success",
  "data": {
    "id": 1,
    "plan_id": 2,
    "plan_name": "Plus",
    "status": "active",
    "updated_at": "2024-01-15T00:00:00Z"
  }
}
```

## Event Publishing Integration

### Publishing Rights Validation

The system automatically validates publishing rights when events are submitted for review:

#### POST /api/v1/events/{id}/submit
Submit event for review (with rights validation).

**Enhanced Behavior:**
1. **Pre-validation**: Check if user has active subscription/package
2. **Limit Check**: Validate remaining weekly/monthly limits or credits
3. **Status Update**: Change event status from "draft" to "pending"
4. **Usage Consumption**: Automatically deduct from limits/credits
5. **Error Handling**: Return detailed error messages for insufficient rights

**Success Response:**
```json
{
  "success": true,
  "message": "event.submit.success",
  "data": {
    "id": 1,
    "status": "pending",
    "submitted_at": "2024-01-15T00:00:00Z",
    "usage_consumed": {
      "type": "monthly_limit",
      "amount": 1,
      "remaining": 19
    }
  }
}
```

**Error Response (Insufficient Rights):**
```json
{
  "success": false,
  "message": "subscription.insufficient_publishing_rights",
  "errors": [
    {
      "field": "subscription",
      "message": "You don't have sufficient publishing rights. Please purchase a subscription or package to publish events."
    }
  ]
}
```

#### GET /api/v1/subscriptions/can-publish
Test publishing rights without submitting an event.

**Response:**
```json
{
  "success": true,
  "message": "subscription.rights.success",
  "data": {
    "can_publish": true,
    "reason": "Active subscription with remaining limits",
    "details": {
      "active_subscription": true,
      "weekly_remaining": 3,
      "monthly_remaining": 12,
      "credits_remaining": 0
    }
  }
}
```

## Stripe Integration

### Webhook Handling

#### POST /api/v1/webhooks/stripe
Handle Stripe webhook events.

**Supported Events:**
- `customer.subscription.created`
- `customer.subscription.updated`
- `customer.subscription.deleted`
- `invoice.payment_succeeded`
- `invoice.payment_failed`
- `payment_intent.succeeded`
- `payment_intent.payment_failed`

**Headers:**
```
Content-Type: application/json
Stripe-Signature: {webhook_signature}
```

**Response:**
```json
{
  "success": true,
  "message": "subscription.webhook.success",
  "data": {
    "event_id": "evt_stripe_id",
    "event_type": "customer.subscription.created",
    "processed_at": "2024-01-15T00:00:00Z"
  }
}
```

### Payment Processing Flow

1. **Frontend**: Collect payment method using Stripe.js
2. **Backend**: Create subscription/payment intent via Stripe API
3. **Webhook**: Process payment confirmation
4. **Database**: Update subscription status
5. **User**: Receive confirmation and access to publishing rights

## Error Handling

### Common Error Codes

| Code | Message Key | Description |
|------|-------------|-------------|
| 400 | `subscription.purchase.invalid_plan` | Invalid subscription plan ID |
| 400 | `subscription.validation.invalid_billing_cycle` | Invalid billing cycle |
| 402 | `subscription.insufficient_publishing_rights` | Insufficient publishing rights |
| 404 | `subscription.plan.not_found` | Subscription plan not found |
| 409 | `subscription.purchase.already_active` | User already has active subscription |
| 409 | `subscription.rights.limit_exceeded` | Publishing limit exceeded |
| 409 | `subscription.rights.credits_exhausted` | Package credits exhausted |
| 500 | `subscription.purchase.payment_failed` | Payment processing failed |

### Error Response Format

```json
{
  "success": false,
  "message": "subscription.insufficient_publishing_rights",
  "errors": [
    {
      "field": "subscription",
      "message": "You don't have sufficient publishing rights. Please purchase a subscription or package to publish events."
    }
  ],
  "data": null
}
```

## Usage Examples

### Example 1: Purchase Basic Subscription

```bash
# 1. Get available plans
curl -X GET "http://localhost:8080/api/v1/subscriptions/plans" \
  -H "Accept-Language: en"

# 2. Purchase Basic monthly subscription
curl -X POST "http://localhost:8080/api/v1/subscriptions/purchase" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en" \
  -d '{
    "plan_id": 1,
    "billing_cycle": "monthly",
    "payment_method_id": "pm_card_visa"
  }'

# 3. Check publishing rights
curl -X GET "http://localhost:8080/api/v1/subscriptions/publishing-rights" \
  -H "Authorization: Bearer {token}" \
  -H "Accept-Language: en"
```

### Example 2: Purchase Event Package

```bash
# 1. Purchase 10 events package
curl -X POST "http://localhost:8080/api/v1/subscriptions/purchase-package" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en" \
  -d '{
    "plan_id": 5,
    "payment_method_id": "pm_card_visa"
  }'

# 2. Check usage statistics
curl -X GET "http://localhost:8080/api/v1/subscriptions/usage-stats" \
  -H "Authorization: Bearer {token}" \
  -H "Accept-Language: en"
```

### Example 3: Submit Event with Rights Validation

```bash
# 1. Try submitting without subscription (should fail)
curl -X POST "http://localhost:8080/api/v1/events/1/submit" \
  -H "Authorization: Bearer {token}" \
  -H "Accept-Language: en"

# Expected response: 402 Insufficient publishing rights

# 2. Purchase package
curl -X POST "http://localhost:8080/api/v1/subscriptions/purchase-package" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "plan_id": 4,
    "payment_method_id": "pm_card_visa"
  }'

# 3. Submit event (should succeed)
curl -X POST "http://localhost:8080/api/v1/events/1/submit" \
  -H "Authorization: Bearer {token}" \
  -H "Accept-Language: en"

# 4. Check remaining credits
curl -X GET "http://localhost:8080/api/v1/subscriptions/usage-stats" \
  -H "Authorization: Bearer {token}" \
  -H "Accept-Language: en"
```

## Testing Scenarios

### Scenario 1: Subscription Workflow
1. **Register** as creator user
2. **Create** draft event
3. **Try submitting** without subscription (expect error)
4. **Purchase** Basic monthly subscription
5. **Submit** event successfully
6. **Check** usage statistics
7. **Submit** more events until limit reached
8. **Try submitting** beyond limit (expect error)

### Scenario 2: Package Workflow
1. **Purchase** 1 Event package
2. **Submit** event successfully (consumes 1 credit)
3. **Try submitting** another event (expect error - no credits)
4. **Purchase** 10 Events package
5. **Submit** multiple events
6. **Check** remaining credits after each submission

### Scenario 3: Mixed Usage
1. **Purchase** Basic weekly subscription
2. **Use** weekly limit
3. **Purchase** 10 Events package for additional events
4. **Submit** events using package credits
5. **Wait** for weekly reset
6. **Submit** events using renewed weekly limit

### Scenario 4: Subscription Management
1. **Purchase** Basic subscription
2. **Upgrade** to Plus subscription
3. **Check** updated limits
4. **Cancel** subscription
5. **Try submitting** after cancellation (expect error)

## Rate Limiting

- **Subscription endpoints**: 100 requests per minute per user
- **Webhook endpoints**: 1000 requests per minute (no user limit)
- **Publishing rights check**: 200 requests per minute per user

## Security Considerations

1. **Webhook Signature Validation**: All Stripe webhooks are validated using signature
2. **Payment Method Security**: Payment methods are handled entirely by Stripe
3. **User Authorization**: All subscription operations require valid user authentication
4. **Usage Tracking**: All usage consumption is logged for audit purposes
5. **Idempotency**: Payment operations include idempotency keys to prevent duplicates

## Monitoring and Logging

The system provides comprehensive logging for:
- **Subscription purchases** and cancellations
- **Usage consumption** and limit resets
- **Publishing rights validation** attempts
- **Payment processing** events
- **Webhook processing** status
- **Error conditions** and failures

All logs include structured data for easy monitoring and alerting.