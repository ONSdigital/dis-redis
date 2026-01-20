package awsauth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

// TokenGenerator generates AWS authentication tokens for AWS.
type TokenGenerator struct {
	creds   aws.CredentialsProvider
	host    string
	region  string
	service string
	signer  *v4.Signer
}

// NewTokenGenerator creates a new TokenGenerator for the specified region and host.
func NewTokenGenerator(ctx context.Context, region, host, service string) (*TokenGenerator, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	return &TokenGenerator{
		creds:   cfg.Credentials,
		host:    host,
		region:  region,
		service: service,
		signer:  v4.NewSigner(),
	}, nil
}

// Generate generates a fresh authentication token for AWS based on
// a dummy request
func (t *TokenGenerator) Generate(ctx context.Context) (string, error) {
	// Create a dummy HTTP request to sign.
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s", t.host), http.NoBody)
	if err != nil {
		return "", err
	}

	currentCreds, err := t.creds.Retrieve(ctx)
	if err != nil {
		return "", err
	}

	err = t.signer.SignHTTP(
		ctx,
		currentCreds,
		req,
		"",
		t.service,
		t.region,
		time.Now(),
	)
	if err != nil {
		return "", err
	}

	return req.Header.Get("Authorization"), nil
}
