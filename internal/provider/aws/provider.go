// Package aws provides functionality for interacting with AWS Secrets Manager.
package aws

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// ClientConfig holds configuration for creating an AWS client
type ClientConfig struct {
	Profile  string
	Region   string
	Endpoint string
}

// NewClientConfig creates a new ClientConfig with the provided values
func NewClientConfig(profile string, region string, endpoint string) ClientConfig {
	return ClientConfig{
		Profile:  profile,
		Region:   region,
		Endpoint: endpoint,
	}
}

// WithEndpoint returns a new ClientConfig with the specified endpoint
func (c ClientConfig) WithEndpoint(endpoint string) ClientConfig {
	return ClientConfig{
		Profile:  c.Profile,
		Region:   c.Region,
		Endpoint: endpoint,
	}
}

// GetCacheKey returns a unique identifier for this configuration
func (c ClientConfig) GetCacheKey() string {
	return fmt.Sprintf("%s:%s:%s", c.Profile, c.Region, c.Endpoint)
}

// AwsProvider handles interactions with AWS Secrets Manager,
// including caching of secrets and API clients.
type AwsProvider struct {
	// Cache of retrieved secrets to avoid repeated API calls
	secretCache    map[string]string
	secretCacheMux sync.RWMutex

	// Cache of API clients for different profiles/regions
	clientCache    map[string]*secretsmanager.Client
	clientCacheMux sync.RWMutex
}

// NewAwsProvider creates a new AWS secrets provider with initialized caches.
// It sets up empty caches for secrets and API clients that will be populated
// as secrets are requested.
func NewAwsProvider() *AwsProvider {
	return &AwsProvider{
		secretCache: make(map[string]string),
		clientCache: make(map[string]*secretsmanager.Client),
	}
}

// GetCachedSecret attempts to retrieve a secret from the cache
func (p *AwsProvider) GetCachedSecret(cacheKey string) functional.Option[string] {
	p.secretCacheMux.RLock()
	defer p.secretCacheMux.RUnlock()

	if value, exists := p.secretCache[cacheKey]; exists {
		return functional.Some(value)
	}
	return functional.None[string]()
}

// CacheSecret stores a secret in the cache for future use
func (p *AwsProvider) CacheSecret(cacheKey string, value string) {
	p.secretCacheMux.Lock()
	defer p.secretCacheMux.Unlock()

	p.secretCache[cacheKey] = value
}

// GetCachedClient attempts to retrieve a client from the cache
func (p *AwsProvider) GetCachedClient(cacheKey string) functional.Option[*secretsmanager.Client] {
	p.clientCacheMux.RLock()
	defer p.clientCacheMux.RUnlock()

	if client, exists := p.clientCache[cacheKey]; exists {
		return functional.Some(client)
	}
	return functional.None[*secretsmanager.Client]()
}

// CacheClient stores a client in the cache for future use
func (p *AwsProvider) CacheClient(cacheKey string, client *secretsmanager.Client) {
	p.clientCacheMux.Lock()
	defer p.clientCacheMux.Unlock()

	p.clientCache[cacheKey] = client
}

// CreateAwsClient creates a new AWS Secrets Manager client based on the configuration
func CreateAwsClient(ctx context.Context, clientConfig ClientConfig) functional.Result[*secretsmanager.Client] {
	var cfg aws.Config
	var err error

	// If endpoint is specified, use LocalStack mode with dummy credentials
	if clientConfig.Endpoint != "" {
		// For LocalStack, use dummy credentials
		staticProvider := credentials.NewStaticCredentialsProvider("test", "test", "")

		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(clientConfig.Region),
			config.WithCredentialsProvider(staticProvider),
		)
	} else {
		// For real AWS, use profile-based configuration
		var options []func(*config.LoadOptions) error
		options = append(options, config.WithRegion(clientConfig.Region))
		if clientConfig.Profile != "" {
			options = append(options, config.WithSharedConfigProfile(clientConfig.Profile))
		}

		cfg, err = config.LoadDefaultConfig(ctx, options...)
	}

	if err != nil {
		return functional.Failure[*secretsmanager.Client](
			fmt.Errorf("failed to load AWS config: %w", err))
	}

	// Create a new client with options
	var clientOptions []func(*secretsmanager.Options)
	if clientConfig.Endpoint != "" {
		clientOptions = append(clientOptions, func(o *secretsmanager.Options) {
			o.BaseEndpoint = &clientConfig.Endpoint
		})
	}

	// Create a new client
	client := secretsmanager.NewFromConfig(cfg, clientOptions...)
	return functional.Success(client)
}

// GetClient returns a cached or new secretsmanager.Client
// It creates a new client for the given profile and region if not found in cache
func (p *AwsProvider) GetClient(ctx context.Context, profile string, region string, endpoint string) (*secretsmanager.Client, error) {
	// Create a configuration for the client
	clientConfig := NewClientConfig(profile, region, endpoint)
	cacheKey := clientConfig.GetCacheKey()

	// Try to get from cache first
	cachedClientOpt := p.GetCachedClient(cacheKey)
	if cachedClientOpt.IsSome() {
		return cachedClientOpt.Unwrap(), nil
	}

	// Create a new client
	clientResult := CreateAwsClient(ctx, clientConfig)
	if clientResult.IsFailure() {
		return nil, clientResult.GetError()
	}

	// Cache the client for future use
	client := clientResult.Unwrap()
	p.CacheClient(cacheKey, client)

	return client, nil
}
