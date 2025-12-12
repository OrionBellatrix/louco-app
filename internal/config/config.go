package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Logger    LoggerConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	AWS       AWSConfig
	Redis     RedisConfig
	Twilio    TwilioConfig
	Email     EmailConfig
	Stripe    StripeConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type LoggerConfig struct {
	Level  string
	Format string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

type AWSConfig struct {
	Endpoint             string
	AccessKeyID          string
	SecretAccessKey      string
	DefaultRegion        string
	Bucket               string
	UsePathStyleEndpoint bool
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type TwilioConfig struct {
	AccountSID    string
	AuthToken     string
	ServiceSID    string
	ReviewerPhone string
	ReviewerOTP   string
	MaxAttempts   int
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
	Currency       string
	Environment    string // "test" or "live"
	SuccessURL     string
	CancelURL      string
	WebhookURL     string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we don't return error if it doesn't exist
		fmt.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 8080),
			Mode: getEnv("SERVER_MODE", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "louco_event_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			Expiration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getEnvAsInt("RATE_LIMIT_RPM", 60),
			BurstSize:         getEnvAsInt("RATE_LIMIT_BURST", 10),
		},
		AWS: AWSConfig{
			Endpoint:             getEnv("AWS_ENDPOINT", "https://nbg1.your-objectstorage.com"),
			AccessKeyID:          getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey:      getEnv("AWS_SECRET_ACCESS_KEY", ""),
			DefaultRegion:        getEnv("AWS_DEFAULT_REGION", "nbg1"),
			Bucket:               getEnv("AWS_BUCKET", "louco-staging"),
			UsePathStyleEndpoint: getEnvAsBool("AWS_USE_PATH_STYLE_ENDPOINT", false),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "127.0.0.1"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Twilio: TwilioConfig{
			AccountSID:    getEnv("TWILIO_ACCOUNT_SID", ""),
			AuthToken:     getEnv("TWILIO_AUTH_TOKEN", ""),
			ServiceSID:    getEnv("TWILIO_SERVICE_SID", ""),
			ReviewerPhone: getEnv("TWILIO_REVIEWER_PHONE_NUMBER", ""),
			ReviewerOTP:   getEnv("TWILIO_REVIEWER_OTP", ""),
			MaxAttempts:   getEnvAsInt("TWILIO_MAX_ATTEMPTS", 5),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("FROM_EMAIL", "noreply@louco-event.com"),
			FromName:     getEnv("FROM_NAME", "Louco Event"),
		},
		Stripe: StripeConfig{
			SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
			Currency:       getEnv("STRIPE_CURRENCY", "usd"),
			Environment:    getEnv("STRIPE_ENVIRONMENT", "test"),
			SuccessURL:     getEnv("STRIPE_SUCCESS_URL", "https://your-domain.com/payment/success"),
			CancelURL:      getEnv("STRIPE_CANCEL_URL", "https://your-domain.com/payment/cancel"),
			WebhookURL:     getEnv("STRIPE_WEBHOOK_URL", "https://your-domain.com/api/v1/webhooks/stripe"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	if c.AWS.AccessKeyID == "" {
		return fmt.Errorf("AWS access key ID is required")
	}

	if c.AWS.SecretAccessKey == "" {
		return fmt.Errorf("AWS secret access key is required")
	}

	if c.Stripe.SecretKey == "" {
		return fmt.Errorf("Stripe secret key is required")
	}

	if c.Stripe.WebhookSecret == "" {
		return fmt.Errorf("Stripe webhook secret is required")
	}

	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
