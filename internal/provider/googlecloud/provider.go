// Package googlecloud provides functionality for interacting with Google Cloud Secret Manager.
package googlecloud

import (
	"context"
	"fmt"
	"sync"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// GoogleCloudProvider handles interactions with Google Cloud Secret Manager,
// including caching of secrets and API clients.
type GoogleCloudProvider struct {
	// Cache of retrieved secrets to avoid repeated API calls
	secretCache    map[string]string
	secretCacheMux sync.RWMutex

	// Cache of API clients
	clientCache    map[string]*secretmanager.Client
	clientCacheMux sync.RWMutex
}

// NewGoogleCloudProvider creates a new Google Cloud secrets provider with initialized caches.
// It sets up empty caches for secrets and API clients that will be populated
// as secrets are requested.
func NewGoogleCloudProvider() *GoogleCloudProvider {
	return &GoogleCloudProvider{
		secretCache: make(map[string]string),
		clientCache: make(map[string]*secretmanager.Client),
	}
}

// GetCachedSecret attempts to retrieve a secret from the cache
func (p *GoogleCloudProvider) GetCachedSecret(cacheKey string) functional.Option[string] {
	p.secretCacheMux.RLock()
	defer p.secretCacheMux.RUnlock()

	if value, exists := p.secretCache[cacheKey]; exists {
		return functional.Some(value)
	}
	return functional.None[string]()
}

// CacheSecret stores a secret in the cache for future use
func (p *GoogleCloudProvider) CacheSecret(cacheKey string, value string) {
	p.secretCacheMux.Lock()
	defer p.secretCacheMux.Unlock()

	p.secretCache[cacheKey] = value
}

// GetCachedClient attempts to retrieve a client from the cache
func (p *GoogleCloudProvider) GetCachedClient(cacheKey string) functional.Option[*secretmanager.Client] {
	p.clientCacheMux.RLock()
	defer p.clientCacheMux.RUnlock()

	if client, exists := p.clientCache[cacheKey]; exists {
		return functional.Some(client)
	}
	return functional.None[*secretmanager.Client]()
}

// CacheClient stores a client in the cache for future use
func (p *GoogleCloudProvider) CacheClient(cacheKey string, client *secretmanager.Client) {
	p.clientCacheMux.Lock()
	defer p.clientCacheMux.Unlock()

	p.clientCache[cacheKey] = client
}

// CreateGoogleCloudClient creates a new Google Cloud Secret Manager client
func CreateGoogleCloudClient(ctx context.Context) functional.Result[*secretmanager.Client] {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return functional.Failure[*secretmanager.Client](
			fmt.Errorf("failed to create secretmanager client: %w", err))
	}
	return functional.Success(client)
}

// GetClient returns a cached or new *secretmanager.Client
func (p *GoogleCloudProvider) GetClient(ctx context.Context) (*secretmanager.Client, error) {
	const cacheKey = "default" // Using default authentication credentials for GoogleCloud

	// Try to get from cache first
	cachedClientOpt := p.GetCachedClient(cacheKey)
	if cachedClientOpt.IsSome() {
		return cachedClientOpt.Unwrap(), nil
	}

	// Create a new client
	clientResult := CreateGoogleCloudClient(ctx)
	if clientResult.IsFailure() {
		return nil, clientResult.GetError()
	}

	// Cache the client for future use
	client := clientResult.Unwrap()
	p.CacheClient(cacheKey, client)

	return client, nil
}
