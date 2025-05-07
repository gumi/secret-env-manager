// Package secret provides secret management related models and utilities
package secret

import (
	"context"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/uri"
)

// Provider defines the interface for secret management service providers
type Provider interface {
	// GetSecret retrieves a secret value from the provider
	GetSecret(ctx context.Context, uri uri.SecretURI) functional.Result[Secret]

	// ListSecrets lists available secrets from the provider
	ListSecrets(ctx context.Context, account string, region string) functional.Result[[]string]

	// GetDefaultOptions returns default options for the provider
	GetDefaultOptions() ProviderOptions

	// SupportsListing returns true if this provider supports listing secrets
	SupportsListing() bool

	// Capabilities returns the set of capabilities this provider supports
	Capabilities() ProviderCapabilities

	// Name returns the name of this provider
	Name() string
}

// ProviderCapabilities represents the set of features a provider supports
type ProviderCapabilities struct {
	SupportsVersioning bool // Whether the provider supports versioned secrets
	SupportsRegions    bool // Whether the provider uses region-specific endpoints
	SupportsJsonKeys   bool // Whether the provider supports JSON key extraction
	SupportsListing    bool // Whether the provider supports listing secrets
}

// ProviderOptions defines configuration options for providers
type ProviderOptions struct {
	DefaultRegion  string            // Default region to use if not specified
	DefaultVersion string            // Default version to use if not specified
	ExtraConfig    map[string]string // Provider-specific extra configuration
}

// NewProviderOptions creates default provider options
func NewProviderOptions() ProviderOptions {
	return ProviderOptions{
		ExtraConfig: make(map[string]string),
	}
}

// WithDefaultRegion returns a copy with the specified default region
func (o ProviderOptions) WithDefaultRegion(region string) ProviderOptions {
	result := o
	result.DefaultRegion = region
	return result
}

// WithDefaultVersion returns a copy with the specified default version
func (o ProviderOptions) WithDefaultVersion(version string) ProviderOptions {
	result := o
	result.DefaultVersion = version
	return result
}

// WithExtraConfig returns a copy with an added extra configuration entry
func (o ProviderOptions) WithExtraConfig(key, value string) ProviderOptions {
	result := o
	if result.ExtraConfig == nil {
		result.ExtraConfig = make(map[string]string)
	} else {
		// Create a new map to avoid modifying the original
		newConfig := make(map[string]string)
		for k, v := range o.ExtraConfig {
			newConfig[k] = v
		}
		result.ExtraConfig = newConfig
	}
	result.ExtraConfig[key] = value
	return result
}
