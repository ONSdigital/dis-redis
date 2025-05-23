package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Client struct {
	redisClient redis.UniversalClient
}

// NewClient returns a new Client with the provided config
func NewClient(ctx context.Context, clientConfig *ClientConfig) (*Client, error) {
	client, err := generateClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("error generating client: %w", err)
	}

	return NewClientWithCustomClient(ctx, clientConfig, client), nil
}

// NewClientWithCustomClient returns a new Client with the provided Redis Client
func NewClientWithCustomClient(ctx context.Context, clientConfig *ClientConfig, client redis.UniversalClient) *Client {
	return &Client{
		redisClient: client,
	}
}

// generateClient creates a Redis Client using the provided configuration
func generateClient(clientConfig *ClientConfig) (redis.UniversalClient, error) {
	options, err := clientConfig.Get()
	if err != nil {
		return nil, fmt.Errorf("error getting client config: %w", err)
	}
	return redis.NewClient(options), nil
}

// GetValue retrieves the value for a given key from Redis and returns it as a string.
func (cli *Client) GetValue(ctx context.Context, key string) (string, error) {
	val, err := cli.redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("key %s not found", key)
	} else if err != nil {
		return "", fmt.Errorf("error getting value for key %s: %w", key, err)
	}

	return val, nil
}

// SetValue sets a key-value pair in Redis with an optional expiration time.
func (cli *Client) SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := cli.redisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		// Wrap and return the error from Redis
		return fmt.Errorf("failed to set value in Redis: %w", err)
	}
	return nil
}
