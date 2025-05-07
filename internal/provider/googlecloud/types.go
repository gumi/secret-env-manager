// Package googlecloud provides types and functions for interacting with Google Cloud services
package googlecloud

// Secret represents a Google Cloud secret with metadata
type Secret struct {
	Name      string // Secret name
	ProjectID string // Google Cloud project ID
	CreatedAt string // Creation timestamp in RFC3339 format
	Version   string // Version identifier (e.g., "latest", "1", "2")
}
