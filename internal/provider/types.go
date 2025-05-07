// Package provider supplies interfaces and implementations for retrieving secrets
package provider

import (
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/uri"
)

// ProviderType represents the type of secret provider
type ProviderType string

const (
	// AwsProvider is the provider type for AWS Secrets Manager
	AwsProvider ProviderType = "aws"
	// GoogleCloudProvider is the provider type for Google Cloud Secret Manager
	GoogleCloudProvider ProviderType = "googlecloud"
)

// AwsSecret represents an AWS secret with metadata
type AwsSecret struct {
	Name      string // Secret name
	ARN       string // AWS Resource Name
	CreatedAt string // Creation timestamp in RFC3339 format
	Version   string // Version identifier (e.g., "AWSCURRENT", "AWSPREVIOUS")
}

// GoogleCloudSecret represents a Google Cloud secret with metadata
type GoogleCloudSecret struct {
	Name      string // Secret name
	ProjectID string // Google Cloud project ID
	CreatedAt string // Creation timestamp in RFC3339 format
	Version   string // Version identifier (e.g., "latest", "1", "2")
}

// ProviderOption represents a configuration option for a provider
type ProviderOption func(*ProviderConfig) *ProviderConfig

// SecretProvider defines the interface for secret retrieval operations
type SecretProvider interface {
	// GetSecrets retrieves a secret value
	GetSecrets(uri uri.SecretURI) (string, error)
	// GetSecretsResult is a monadic version of GetSecrets
	GetSecretsResult(uri uri.SecretURI) functional.Result[string]
	// GetConfig returns the provider configuration
	GetConfig() ProviderConfig
}

// FunctionalSecretProvider defines a provider interface using Result monad
type FunctionalSecretProvider interface {
	SecretProvider
	GetSecretsResult(uri uri.SecretURI) functional.Result[string]
}

// WithEndpointURL sets the endpoint URL for a provider
func WithEndpointURL(url string) ProviderOption {
	return func(config *ProviderConfig) *ProviderConfig {
		config.EndpointURL = url
		return config
	}
}
