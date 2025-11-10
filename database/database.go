package database

import (
	"log"

	"go_backend/config"
)

// ConnectAll connects to all databases
func ConnectAll(cfg *config.Config) error {
	// Connect PostgreSQL
	if _, err := ConnectPostgres(&cfg.Postgres); err != nil {
		log.Printf("⚠️  PostgreSQL connection failed: %v", err)
		// Continue even if PostgreSQL fails (optional)
	} else {
		// Run migrations if PostgreSQL is connected
		if err := AutoMigrate(PostgresDB); err != nil {
			log.Printf("⚠️  PostgreSQL migration failed: %v", err)
		}
	}

	// Connect MongoDB
	if _, _, err := ConnectMongoDB(&cfg.MongoDB); err != nil {
		log.Printf("⚠️  MongoDB connection failed: %v", err)
		// Continue even if MongoDB fails (optional)
	}

	// Connect Redis
	if _, err := ConnectRedis(&cfg.Redis); err != nil {
		log.Printf("⚠️  Redis connection failed: %v", err)
		// Continue even if Redis fails (optional)
	}

	return nil
}

// CloseAll closes all database connections
func CloseAll() {
	log.Println("Closing database connections...")

	if err := ClosePostgres(); err != nil {
		log.Printf("Error closing PostgreSQL: %v", err)
	}

	if err := CloseMongoDB(); err != nil {
		log.Printf("Error closing MongoDB: %v", err)
	}

	if err := CloseRedis(); err != nil {
		log.Printf("Error closing Redis: %v", err)
	}

	log.Println("✅ All database connections closed")
}

