-- =====================================================
-- SUBSCRIPTION & PACKAGE PLANS SEEDER
-- =====================================================
-- This file contains seed data for subscription_plans table
-- Run this after database migration to populate default plans
-- =====================================================

-- Clear existing data (optional - uncomment if needed)
-- DELETE FROM subscription_plans;
-- ALTER SEQUENCE subscription_plans_id_seq RESTART WITH 1;

-- =====================================================
-- SUBSCRIPTION PLANS
-- =====================================================

-- Basic Subscription Plan
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    billing_cycle, 
    weekly_limit, 
    monthly_limit, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'subscription',
    'basic',
    'Basic Plan',
    'Perfect for getting started with event creation',
    78.00,
    'EUR',
    'monthly',
    1,
    4,
    30,
    true,
    1,
    '{"features": ["1 event per week", "4 events per month", "Basic support"], "popular": false, "color": "#3B82F6", "badge": null}',
    NOW(),
    NOW()
);

-- Plus Subscription Plan (Most Popular)
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    billing_cycle, 
    weekly_limit, 
    monthly_limit, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'subscription',
    'plus',
    'Plus Plan',
    'Great for regular event creators',
    130.00,
    'EUR',
    'monthly',
    2,
    8,
    30,
    true,
    2,
    '{"features": ["2 events per week", "8 events per month", "Priority support", "Advanced features"], "popular": true, "color": "#10B981", "badge": "Most Popular"}',
    NOW(),
    NOW()
);

-- Pro Subscription Plan
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    billing_cycle, 
    weekly_limit, 
    monthly_limit, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'subscription',
    'pro',
    'Pro Plan',
    'For professional event organizers',
    156.00,
    'EUR',
    'monthly',
    3,
    12,
    30,
    true,
    3,
    '{"features": ["3 events per week", "12 events per month", "Premium support", "Advanced analytics", "Priority listing"], "popular": false, "color": "#8B5CF6", "badge": "Professional"}',
    NOW(),
    NOW()
);

-- =====================================================
-- PACKAGE PLANS (One-time purchases)
-- =====================================================

-- Single Event Package
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    total_credits, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'package',
    'single',
    'Single Event',
    'Perfect for one-time events',
    29.00,
    'EUR',
    1,
    365,
    true,
    4,
    '{"features": ["1 event credit", "Valid for 1 year", "Basic support"], "popular": false, "color": "#6B7280", "badge": null}',
    NOW(),
    NOW()
);

-- 10 Events Package (Best Value)
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    total_credits, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'package',
    '10_events',
    '10 Events Package',
    'Better value for multiple events',
    249.99,
    'EUR',
    10,
    365,
    true,
    5,
    '{"features": ["10 event credits", "Valid for 1 year", "Better value", "Save €40"], "popular": true, "color": "#F59E0B", "badge": "Best Value", "savings": 40.01}',
    NOW(),
    NOW()
);

-- 25 Events Package (Bulk Discount)
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    total_credits, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'package',
    '25_events',
    '25 Events Package',
    'Best value for frequent event creators',
    499.99,
    'EUR',
    25,
    365,
    true,
    6,
    '{"features": ["25 event credits", "Valid for 1 year", "Bulk discount", "Save €225"], "popular": false, "color": "#EF4444", "badge": "Enterprise", "savings": 225.01}',
    NOW(),
    NOW()
);

-- =====================================================
-- ADDITIONAL PREMIUM PACKAGES (Optional)
-- =====================================================

-- 5 Events Package (Mid-tier option)
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    total_credits, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'package',
    '5_events',
    '5 Events Package',
    'Good for small event series',
    129.99,
    'EUR',
    5,
    365,
    true,
    7,
    '{"features": ["5 event credits", "Valid for 1 year", "Small discount", "Save €15"], "popular": false, "color": "#06B6D4", "badge": null, "savings": 15.01}',
    NOW(),
    NOW()
);

-- 50 Events Package (Enterprise)
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    total_credits, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'package',
    '50_events',
    '50 Events Package',
    'Enterprise solution for large organizations',
    899.99,
    'EUR',
    50,
    365,
    true,
    8,
    '{"features": ["50 event credits", "Valid for 1 year", "Enterprise discount", "Save €550", "Priority support"], "popular": false, "color": "#7C3AED", "badge": "Enterprise", "savings": 550.01}',
    NOW(),
    NOW()
);

-- =====================================================
-- WEEKLY SUBSCRIPTION PLANS (Alternative billing)
-- =====================================================

-- Weekly Basic Plan
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    billing_cycle, 
    weekly_limit, 
    monthly_limit, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'subscription',
    'weekly_basic',
    'Weekly Basic',
    'Flexible weekly subscription',
    22.00,
    'EUR',
    'weekly',
    1,
    4,
    7,
    true,
    9,
    '{"features": ["1 event per week", "Weekly billing", "Cancel anytime"], "popular": false, "color": "#3B82F6", "badge": "Flexible"}',
    NOW(),
    NOW()
);

-- Weekly Plus Plan
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    billing_cycle, 
    weekly_limit, 
    monthly_limit, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'subscription',
    'weekly_plus',
    'Weekly Plus',
    'Enhanced weekly subscription',
    38.00,
    'EUR',
    'weekly',
    2,
    8,
    7,
    true,
    10,
    '{"features": ["2 events per week", "Weekly billing", "Priority support"], "popular": false, "color": "#10B981", "badge": "Flexible"}',
    NOW(),
    NOW()
);

-- =====================================================
-- SPECIAL/PROMOTIONAL PLANS
-- =====================================================

-- Trial Package (Free trial)
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    total_credits, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'package',
    'trial',
    'Free Trial',
    'Try our platform with one free event',
    0.00,
    'EUR',
    1,
    30,
    true,
    0,
    '{"features": ["1 free event credit", "Valid for 30 days", "No payment required"], "popular": false, "color": "#22C55E", "badge": "Free Trial", "trial": true}',
    NOW(),
    NOW()
);

-- Student Discount Package
INSERT INTO subscription_plans (
    type, 
    name, 
    display_name, 
    description, 
    price, 
    currency, 
    total_credits, 
    duration_days, 
    is_active, 
    sort_order, 
    metadata, 
    created_at, 
    updated_at
) VALUES (
    'package',
    'student_5',
    'Student Package',
    'Special pricing for students and educational institutions',
    89.99,
    'EUR',
    5,
    365,
    false, -- Disabled by default, enable when needed
    11,
    '{"features": ["5 event credits", "Valid for 1 year", "Student discount", "Verification required"], "popular": false, "color": "#F97316", "badge": "Student", "discount_type": "student", "original_price": 129.99}',
    NOW(),
    NOW()
);

-- =====================================================
-- VERIFICATION QUERIES
-- =====================================================

-- Verify subscription plans
SELECT 
    id,
    type,
    name,
    display_name,
    price,
    currency,
    CASE 
        WHEN type = 'subscription' THEN CONCAT(weekly_limit, ' events/week, ', monthly_limit, ' events/month')
        WHEN type = 'package' THEN CONCAT(total_credits, ' event credits')
    END as limits,
    duration_days,
    is_active,
    sort_order
FROM subscription_plans 
WHERE type = 'subscription'
ORDER BY sort_order;

-- Verify package plans
SELECT 
    id,
    type,
    name,
    display_name,
    price,
    currency,
    total_credits,
    duration_days,
    is_active,
    sort_order
FROM subscription_plans 
WHERE type = 'package'
ORDER BY sort_order;

-- Count total plans
SELECT 
    type,
    COUNT(*) as plan_count,
    COUNT(CASE WHEN is_active = true THEN 1 END) as active_count
FROM subscription_plans 
GROUP BY type;

-- =====================================================
-- NOTES
-- =====================================================
/*
PRICING STRATEGY:
- Subscriptions: Monthly recurring with weekly/monthly limits
- Packages: One-time purchase with credit system
- Bulk discounts: Larger packages offer better per-event pricing
- Trial: Free option to attract new users
- Student: Special pricing for educational sector

METADATA FIELDS:
- features: Array of plan features for display
- popular: Boolean to highlight recommended plans
- color: UI color scheme for plan display
- badge: Special badge text (e.g., "Most Popular", "Best Value")
- savings: Amount saved compared to single event pricing
- trial: Boolean to identify trial plans
- discount_type: Type of discount applied

ACTIVATION:
- Most plans are active by default
- Special plans (student, promotional) are disabled by default
- Use is_active flag to enable/disable plans without deletion

SORT ORDER:
- 0: Trial/Free plans
- 1-3: Main subscription plans
- 4-8: Package plans
- 9-10: Alternative billing subscriptions
- 11+: Special/promotional plans
*/