package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vashafaris/nucleus/pkg/config"
)

type Client struct {
	*redis.Client
}

// NewConnection creates a new Redis connection
func NewConnection(cfg *config.RedisConfig) (*Client, error) {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,

		// Connection pool settings
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{client}, nil
}

// Health checks the Redis connection
func (c *Client) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result := c.Ping(ctx)
	if err := result.Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

// GetContext is a wrapper around Get with context
func (c *Client) GetContext(ctx context.Context, key string) (string, error) {
	return c.Get(ctx, key).Result()
}

// SetContext is a wrapper around Set with context
func (c *Client) SetContext(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Set(ctx, key, value, expiration).Err()
}

// DeleteContext is a wrapper around Del with context
func (c *Client) DeleteContext(ctx context.Context, keys ...string) error {
	return c.Del(ctx, keys...).Err()
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.Client.Close()
}
