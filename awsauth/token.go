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
	creds  aws.CredentialsProvider
	region string
	host   string
}

// NewTokenGenerator creates a new TokenGenerator for the specified region and host.
func NewTokenGenerator(ctx context.Context, region, host string) (*TokenGenerator, error) {
	creds, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	return &TokenGenerator{
		creds:  creds.Credentials,
		region: region,
		host:   host,
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

	signer := v4.NewSigner()

	err = signer.SignHTTP(
		ctx,
		currentCreds,
		req,
		"",
		"memorydb",
		t.region,
		time.Now(),
	)
	if err != nil {
		return "", err
	}

	return req.Header.Get("Authorization"), nil
}
