// Package uri provides functionality for parsing and handling secret URIs.
package uri

import (
	"net/url"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// Default version constants
const (
	AwsDefaultVersion         = "AWSCURRENT" // Default AWS version label
	GoogleCloudDefaultVersion = "latest"     // Default Google Cloud version label
)

// SecretURI represents a parsed secret URI
type SecretURI struct {
	Platform   string // Cloud platform (aws, googlecloud)
	Service    string // Service type (secretsmanager, secretmanager)
	Account    string // Account identifier (e.g., AWS profile, Google Cloud project ID)
	SecretName string // Name of the secret
	Key        string // Optional key name for JSON secrets
	Version    string // Version of the secret
	Region     string // Region (mainly for AWS)
}

// Methods for SecretURI type

// WithRegion returns a copy of the SecretURI with the specified region
func (s SecretURI) WithRegion(region string) SecretURI {
	result := s
	result.Region = region
	return result
}

// WithVersion returns a copy of the SecretURI with the specified version
func (s SecretURI) WithVersion(version string) SecretURI {
	result := s
	result.Version = version
	return result
}

// WithKey returns a copy of the SecretURI with the specified key
func (s SecretURI) WithKey(key string) SecretURI {
	result := s
	result.Key = key
	return result
}

// IsComplete checks if the URI has all required fields
func (s SecretURI) IsComplete() bool {
	return s.Platform != "" && s.Service != "" &&
		s.Account != "" && s.SecretName != ""
}

// GetCacheKey returns a unique identifier for caching purposes
func (s SecretURI) GetCacheKey() string {
	return BuildCacheKey(s.Account, s.Service, s.SecretName, s.Version, s.Region)
}

// AsOption converts a SecretURI to an Option type
func (s SecretURI) AsOption() functional.Option[SecretURI] {
	if !s.IsComplete() {
		return functional.None[SecretURI]()
	}
	return functional.Some(s)
}

// GetUri returns the formatted secret URI as a string
func (s SecretURI) GetUri() string {
	// Using function composition pattern for a more functional approach
	return functional.ApplyAll(
		URIPrefix+s.Platform+":"+s.Service+"/"+s.Account+"/"+s.SecretName,
		func(base string) string {
			return appendQueryParamsIfNeeded(base, s)
		},
	)
}

// Exported Functions

// NewSecretURI creates a new SecretURI with the required fields
func NewSecretURI(platform, service, account, secretName string) SecretURI {
	return SecretURI{
		Platform:   platform,
		Service:    service,
		Account:    account,
		SecretName: secretName,
	}
}

// BuildCacheKey creates a consistent cache key from secret information
func BuildCacheKey(account, service, secretName, version, region string) string {
	parts := []string{account, service, secretName, version, region}
	return strings.Join(parts, "|")
}

// Unexported Helper Functions

// appendQueryParamsIfNeeded adds query parameters to the URI if any exist
func appendQueryParamsIfNeeded(baseUri string, s SecretURI) string {
	params := collectQueryParams(s)
	if len(params) == 0 {
		return baseUri
	}
	return baseUri + "?" + strings.Join(params, "&")
}

// collectQueryParams gathers all non-empty query parameters
// Pure function that transforms URI components into query parameters
func collectQueryParams(s SecretURI) []string {
	// Use a more functional approach with Map and Filter
	possibleParams := []struct {
		name  string
		value string
	}{
		{"version", s.Version},
		{"key", s.Key},
		{"region", s.Region},
	}

	// Filter out empty values and map to parameter strings
	return functional.Filter(
		functional.Map(
			possibleParams,
			func(p struct {
				name  string
				value string
			}) string {
				return p.name + "=" + url.QueryEscape(p.value)
			},
		),
		func(param string) bool {
			return !strings.HasSuffix(param, "=")
		},
	)
}
