package database

import (
	"fmt"
	"time"

	"github.com/louco-event/internal/config"
	"github.com/louco-event/internal/domain"
	pkgLogger "github.com/louco-event/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

func New(cfg *config.Config, logger *pkgLogger.Logger) (*Database, error) {
	dsn := cfg.GetDSN()

	// GORM config
	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent), // Use GORM's default silent logger
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().Msg("Database connection established")

	return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate() error {
	// Run auto migration for all models
	err := d.DB.AutoMigrate(
		&domain.User{},
		&domain.Media{},
		&domain.Industry{},
		&domain.Creator{},
		&domain.CreatorIndustry{},
		&domain.Category{},
		&domain.VerificationCode{},
		&domain.Follow{},
		&domain.Address{},
		&domain.Event{},
		&domain.EventCategory{},
		&domain.Ticket{},
		&domain.Invitation{},
		&domain.SubscriptionPlan{},
		&domain.UserSubscription{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health check
func (d *Database) Health() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
