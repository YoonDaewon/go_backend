package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go_backend/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// ConnectRedis connects to Redis database
func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("✅ Redis connected successfully")

	RedisClient = client
	return client, nil
}

// CloseRedis closes Redis connection
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

