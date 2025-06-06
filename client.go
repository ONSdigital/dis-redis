package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
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

// GetKeyValuePairs retrieves a set of key-value pairs from Redis based on a match pattern and a given cursor.
func (cli *Client) GetKeyValuePairs(ctx context.Context, matchPattern string, count int64, cursor uint64) (keyValuePairs map[string]string, newCursor uint64, err error) {
	keyValuePairs = make(map[string]string)

	// Get the list of keys matching the pattern, starting from the provided cursor
	keys, newCursor, err := cli.redisClient.Scan(ctx, cursor, matchPattern, count).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("error scanning keys: %w", err)
	}

	// If we have keys, get the values for those keys
	if len(keys) > 0 {
		values, err := cli.redisClient.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, 0, fmt.Errorf("error fetching values for keys: %w", err)
		}

		// Add the key-value pairs to the map
		for i, key := range keys {
			if i < len(values) {
				if val, ok := values[i].(string); ok {
					keyValuePairs[key] = val
				}
			}
		}
	}

	// Return the results along with the new cursor
	return keyValuePairs, newCursor, nil
}

// GetTotalKeys returns the total number of keys in the Redis database.
func (cli *Client) GetTotalKeys(ctx context.Context) (int64, error) {
	count, err := cli.redisClient.DBSize(ctx).Result()
	if err != nil {
		return 0, fmt.Errorf("error getting total number of keys: %w", err)
	}
	return count, nil
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
