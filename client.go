package redis

import (
	"context"
	"errors"
	"fmt"

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

// DeleteValue deletes a key-value pair from Redis
func (cli *Client) DeleteValue(ctx context.Context, key string) error {
	// Call the Del method to delete the key
	result, err := cli.redisClient.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete key from Redis: %w", err)
	}

	if result == 0 {
		// If the result is 0, it means the key does not exist
		return fmt.Errorf("key not found: %s", key)
	}

	return nil
}
