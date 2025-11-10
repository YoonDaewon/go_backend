package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go_backend/database"
	"go_backend/model"

	"github.com/redis/go-redis/v9"
)

// RedisCache handles caching operations using Redis
type RedisCache interface {
	SetUser(ctx context.Context, user *model.User, ttl time.Duration) error
	GetUser(ctx context.Context, id int) (*model.User, error)
	DeleteUser(ctx context.Context, id int) error
	SetUsers(ctx context.Context, users []*model.User, ttl time.Duration) error
	GetUsers(ctx context.Context) ([]*model.User, error)
	DeleteUsers(ctx context.Context) error
}

type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache() RedisCache {
	return &redisCache{
		client: database.RedisClient,
	}
}

// SetUser caches a user
func (r *redisCache) SetUser(ctx context.Context, user *model.User, ttl time.Duration) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	key := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

// GetUser retrieves a cached user
func (r *redisCache) GetUser(ctx context.Context, id int) (*model.User, error) {
	if r.client == nil {
		return nil, fmt.Errorf("redis client not available")
	}

	key := fmt.Sprintf("user:%d", id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser removes a user from cache
func (r *redisCache) DeleteUser(ctx context.Context, id int) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	key := fmt.Sprintf("user:%d", id)
	return r.client.Del(ctx, key).Err()
}

// SetUsers caches the list of all users
func (r *redisCache) SetUsers(ctx context.Context, users []*model.User, ttl time.Duration) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	key := "users:all"
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

// GetUsers retrieves the cached list of all users
func (r *redisCache) GetUsers(ctx context.Context) ([]*model.User, error) {
	if r.client == nil {
		return nil, fmt.Errorf("redis client not available")
	}

	key := "users:all"
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var users []*model.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// DeleteUsers removes the cached list of all users
func (r *redisCache) DeleteUsers(ctx context.Context) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	key := "users:all"
	return r.client.Del(ctx, key).Err()
}

