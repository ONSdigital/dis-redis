package redis

import (
	"context"
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
