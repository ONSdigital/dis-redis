package redis

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/ONSdigital/dis-redis/awsauth"
	redis "github.com/redis/go-redis/v9"
)

// ClientConfig exposes the optional configurable parameters for a client to overwrite default redis options values.
// Any value that is not provided will use the default redis options value.
type ClientConfig struct {
	ClusterName string
	Region      string
	Service     string
	Username    string
	// go-redis config overrides
	Address   string
	Database  *int
	TLSConfig *tls.Config
}

// Get creates a default redis options and overwrites with any values provided in ClientConfig
func (c *ClientConfig) Get(ctx context.Context) (*redis.Options, error) {
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

	if c.TLSConfig != nil {
		cfg.TLSConfig = c.TLSConfig
	}

	if c.Username != "" {
		cfg.Username = c.Username
	}

	if c.Region != "" && c.Username != "" && c.Service != "" {
		credsProvider, err := getAWSCredsProvider(ctx, c.ClusterName, c.Address, c.Region, c.Service, c.Username)
		if err != nil {
			return nil, fmt.Errorf("error getting AWS credentials provider: %w", err)
		}

		cfg.CredentialsProviderContext = credsProvider
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

func getAWSCredsProvider(ctx context.Context, clusterName, endpoint, region, service, username string) (func(context.Context) (string, string, error), error) {
	tokenGenerator, err := awsauth.NewTokenGenerator(ctx, clusterName, endpoint, region, service, username)
	if err != nil {
		return nil, fmt.Errorf("error creating token generator: %w", err)
	}

	credsProvider := func(credsCtx context.Context) (string, string, error) {
		token, err := tokenGenerator.Generate(credsCtx)
		return username, token, err
	}

	return credsProvider, nil
}

// Validate will validate that compulsory values are provided in config
func (c *ClientConfig) Validate() error {
	// Validations on the config can be added here.
	if c.Region != "" && c.Username == "" {
		return fmt.Errorf("username must be provided when region is set")
	}

	if c.Username != "" && c.Region == "" {
		return fmt.Errorf("region must be provided when username is set")
	}
	return nil
}
