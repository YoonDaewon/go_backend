package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go_backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClient *mongo.Client
var MongoDB *mongo.Database

// ConnectMongoDB connects to MongoDB database
func ConnectMongoDB(cfg *config.MongoDBConfig) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := cfg.GetURI()
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("✅ MongoDB connected successfully")

	MongoDBClient = client
	MongoDB = client.Database(cfg.DBName)

	return client, MongoDB, nil
}

// CloseMongoDB closes MongoDB connection
func CloseMongoDB() error {
	if MongoDBClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return MongoDBClient.Disconnect(ctx)
	}
	return nil
}

