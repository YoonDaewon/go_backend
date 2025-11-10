package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go_backend/config"
	"go_backend/database"
	"go_backend/router"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to databases
	if err := database.ConnectAll(cfg); err != nil {
		log.Printf("Warning: Some database connections failed: %v", err)
	}

	// Setup graceful shutdown
	setupGracefulShutdown()

	// Setup router
	r := router.SetupRouter()

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("\n🛑 Shutting down gracefully...")
		database.CloseAll()
		os.Exit(0)
	}()
}
