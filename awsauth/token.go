package awsauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

const (
	connectAction string = "connect"

	hexEncodedSHA256EmptyString = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	tokenValiditySeconds = 900
)

// TokenGenerator generates AWS authentication tokens for AWS.
type TokenGenerator struct {
	creds       aws.CredentialsProvider
	clusterName string
	host        string
	region      string
	service     string
	signer      *v4.Signer
	username    string
}

// NewTokenGenerator creates a new TokenGenerator for the specified region and host.
func NewTokenGenerator(ctx context.Context, clusterName, host, region, service, username string) (*TokenGenerator, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	return &TokenGenerator{
		clusterName: clusterName,
		creds:       cfg.Credentials,
		host:        host,
		region:      region,
		service:     service,
		signer:      v4.NewSigner(),
		username:    username,
	}, nil
}

// Generate generates a fresh authentication token for AWS based on
// a dummy request
func (t *TokenGenerator) Generate(ctx context.Context) (string, error) {
	// Create a dummy request to sign
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/", t.clusterName), http.NoBody)
	if err != nil {
		return "", err
	}

	currentCreds, err := t.creds.Retrieve(ctx)
	if err != nil {
		return "", err
	}

	if currentCreds.AccessKeyID == "" || currentCreds.SecretAccessKey == "" {
		return "", errors.New("empty AWS credentials")
	}

	q := req.URL.Query()
	q.Set("Action", connectAction)
	q.Set("User", t.username)
	q.Set("X-Amz-Expires", strconv.FormatInt(int64(tokenValiditySeconds), 10))

	req.URL.RawQuery = q.Encode()

	signedURI, _, err := t.signer.PresignHTTP(
		ctx,
		currentCreds,
		req,
		hexEncodedSHA256EmptyString,
		t.service,
		t.region,
		time.Now().UTC(),
	)
	if err != nil {
		return "", err
	}

	token := strings.TrimPrefix(signedURI, "https://")

	return token, nil
}
