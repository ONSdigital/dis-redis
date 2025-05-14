package redis

import (
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

// ClientConfig exposes the optional configurable parameters for a client to overwrite default redis options values.
// Any value that is not provided will use the default redis options value.
type ClientConfig struct {
	// go-redis config overrides
	Address  string
	Database *int
}

// Get creates a default redis options and overwrites with any values provided in ClientConfig
func (c *ClientConfig) Get() (*redis.Options, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Get default redis config and apply overrides
	cfg := getDefaultConfig()

	if c.Database != nil {
		cfg.DB = *c.Database
	}

	if c.Address != "" {
		cfg.Addr = c.Address
	}

	return cfg, nil
}

// getDefaultConfig returns a default set of redis.Options
func getDefaultConfig() *redis.Options {
	return &redis.Options{
		DB:   0,
		Addr: "localhost:6379",
	}
}

// Validate will validate that compulsory values are provided in config
func (c *ClientConfig) Validate() error {
	// Validations on the config can be added here.
	return nil
}
