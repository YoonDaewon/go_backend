package database

import (
	"fmt"
	"log"

	"go_backend/config"
	"go_backend/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var PostgresDB *gorm.DB

// ConnectPostgres connects to PostgreSQL database
func ConnectPostgres(cfg *config.PostgresConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ PostgreSQL connected successfully")

	PostgresDB = db
	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.User{},
		// Add more models here
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("✅ Database migration completed")
	return nil
}

// ClosePostgres closes PostgreSQL connection
func ClosePostgres() error {
	if PostgresDB != nil {
		sqlDB, err := PostgresDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

