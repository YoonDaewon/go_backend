package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	MongoDB  MongoDBConfig
	Redis    RedisConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Env  string
}

// PostgresConfig holds PostgreSQL configuration
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	URI      string
	DBName   string
	Username string
	Password string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_DB", "go_backend"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			DBName:   getEnv("MONGODB_DB", "go_backend"),
			Username: getEnv("MONGODB_USERNAME", ""),
			Password: getEnv("MONGODB_PASSWORD", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}

	AppConfig = config
	return config, nil
}

// GetPostgresDSN returns PostgreSQL connection string
func (c *PostgresConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetMongoURI returns MongoDB connection URI
func (c *MongoDBConfig) GetURI() string {
	if c.URI != "" {
		return c.URI
	}
	if c.Username != "" && c.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s/%s",
			c.Username, c.Password, c.URI, c.DBName)
	}
	return fmt.Sprintf("%s/%s", c.URI, c.DBName)
}

// GetRedisAddr returns Redis address
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	var value int
	fmt.Sscanf(valueStr, "%d", &value)
	return value
}
