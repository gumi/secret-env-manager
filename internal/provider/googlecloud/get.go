// Package googlecloud provides functionality for interacting with Google Cloud Secret Manager.
package googlecloud

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/uri"
	"github.com/gumi-tsd/secret-env-manager/internal/secret"
)

// SecretRequest encapsulates all parameters needed to retrieve a secret
type SecretRequest struct {
	URI uri.SecretURI
	Ctx context.Context
}

// NewSecretRequest creates a new SecretRequest with default context
func NewSecretRequest(uri uri.SecretURI) SecretRequest {
	return SecretRequest{
		URI: uri,
		Ctx: context.Background(),
	}
}

// WithContext returns a new SecretRequest with the specified context
func (r SecretRequest) WithContext(ctx context.Context) SecretRequest {
	return SecretRequest{
		URI: r.URI,
		Ctx: ctx,
	}
}

// GetCacheKey returns a unique identifier for caching
func (r SecretRequest) GetCacheKey() string {
	return uri.BuildCacheKey(r.URI.Account, r.URI.Service, r.URI.SecretName, r.URI.Version, r.URI.Region)
}

// GetSecrets retrieves secrets using a background context
func (p *GoogleCloudProvider) GetSecrets(uri uri.SecretURI) (string, error) {
	req := NewSecretRequest(uri)
	result := p.GetSecretValue(req)

	if result.IsFailure() {
		return "", result.GetError()
	}

	return result.Unwrap(), nil
}

// GetSecretValue is a functional implementation that fetches and processes secrets
func (p *GoogleCloudProvider) GetSecretValue(req SecretRequest) functional.Result[string] {
	// Check cache first
	cacheKey := req.GetCacheKey()
	cachedSecret := p.GetCachedSecret(cacheKey)

	// If found in cache, parse and return
	if cachedSecret.IsSome() {
		valueResult := secret.ParseValueResult(cachedSecret.Unwrap(), req.URI.Key)
		if valueResult.IsFailure() {
			err := valueResult.GetError()
			if strings.Contains(err.Error(), "key not found") {
				return functional.Failure[string](
					fmt.Errorf("specified key '%s' does not exist in secret '%s': %w",
						req.URI.Key, req.URI.SecretName, err))
			}
			return valueResult
		}
		return valueResult
	}

	// Fetch the secret if not in cache
	secretResult := p.RetrieveSecret(req)
	if secretResult.IsFailure() {
		return secretResult
	}

	// Store in cache
	secretVal := secretResult.Unwrap()
	p.CacheSecret(cacheKey, secretVal)

	// Parse and return the value
	valueResult := secret.ParseValueResult(secretVal, req.URI.Key)
	if valueResult.IsFailure() {
		err := valueResult.GetError()
		if strings.Contains(err.Error(), "key not found") {
			return functional.Failure[string](
				fmt.Errorf("specified key '%s' does not exist in secret '%s': %w",
					req.URI.Key, req.URI.SecretName, err))
		}
		return valueResult
	}
	return valueResult
}

// RetrieveSecret fetches a secret from Google Cloud Secret Manager
func (p *GoogleCloudProvider) RetrieveSecret(req SecretRequest) functional.Result[string] {
	// Get or create client
	client, err := p.GetClient(req.Ctx)
	if err != nil {
		return functional.Failure[string](
			fmt.Errorf("failed to get Google Cloud client - project: %s: %w",
				req.URI.Account, err))
	}

	// Fetch secret
	secretResult := FetchSecret(req.Ctx, client, req.URI)
	if secretResult.IsFailure() {
		return functional.Failure[string](
			fmt.Errorf("failed to retrieve secret [%s/%s] - project: %s: %w",
				req.URI.Service, req.URI.SecretName, req.URI.Account, secretResult.GetError()))
	}

	return secretResult
}

// FetchSecret calls Google Cloud Secret Manager API to get a secret value
func FetchSecret(ctx context.Context, client *secretmanager.Client, uri uri.SecretURI) functional.Result[string] {
	// Construct the resource name
	// Format: projects/{project}/secrets/{secret}/versions/{version}
	resourceName := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", uri.Account, uri.SecretName, uri.Version)

	// Use proper logging format instead of fmt.Println
	secret.LogInfoMsg(fmt.Sprintf("Accessing secret: %s", uri.GetUri()))

	// Build the request
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: resourceName,
	}

	// Call Google Cloud API
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return functional.Failure[string](
			fmt.Errorf("Google Cloud Secret Manager API error [%s] - version: %s: %w",
				uri.SecretName, uri.Version, err))
	}

	// Validate response
	if result.Payload == nil || result.Payload.Data == nil {
		return functional.Failure[string](
			fmt.Errorf("empty secret value [%s] - version: %s",
				uri.SecretName, uri.Version))
	}

	// Return the secret string
	return functional.Success(string(result.Payload.Data))
}
