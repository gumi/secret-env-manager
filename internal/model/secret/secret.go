// Package secret provides secret management related models and utilities
package secret

import (
	"time"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/uri"
)

// Secret represents a secret value retrieved from a cloud provider
type Secret struct {
	URI       uri.SecretURI // The URI that identifies this secret
	Value     string        // The secret value
	Timestamp time.Time     // When this secret was retrieved
	Raw       interface{}   // Provider-specific raw data
}

// NewSecret creates a new secret with the specified URI and value
func NewSecret(uri uri.SecretURI, value string) Secret {
	return Secret{
		URI:       uri,
		Value:     value,
		Timestamp: time.Now(),
	}
}

// WithRaw returns a new Secret with the specified raw data
func (s Secret) WithRaw(raw interface{}) Secret {
	result := s
	result.Raw = raw
	return result
}

// WithTimestamp returns a new Secret with the specified timestamp
func (s Secret) WithTimestamp(t time.Time) Secret {
	result := s
	result.Timestamp = t
	return result
}

// AsOption converts a Secret to an Option type
func (s Secret) AsOption() functional.Option[Secret] {
	if s.Value == "" {
		return functional.None[Secret]()
	}
	return functional.Some(s)
}

// SecretMap is a type alias for a map of secret URIs to their resolved values
type SecretMap map[string]Secret

// Get retrieves a secret by its cache key
func (m SecretMap) Get(key string) functional.Option[Secret] {
	if val, ok := m[key]; ok {
		return functional.Some(val)
	}
	return functional.None[Secret]()
}

// Set adds or updates a secret in the map
func (m SecretMap) Set(key string, secret Secret) SecretMap {
	result := make(SecretMap, len(m)+1)
	for k, v := range m {
		result[k] = v
	}
	result[key] = secret
	return result
}

// IsExpired checks if a secret is older than the given duration
func (s Secret) IsExpired(maxAge time.Duration) bool {
	return time.Since(s.Timestamp) > maxAge
}

// HasJsonKey checks if this secret has a specific JSON key in its URI
func (s Secret) HasJsonKey() bool {
	return s.URI.Key != ""
}
