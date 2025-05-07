// Package aws provides types and functions for interacting with AWS services
package aws

// Secret represents an AWS secret with metadata
type Secret struct {
	Name      string // Secret name
	ARN       string // AWS Resource Name
	CreatedAt string // Creation timestamp in RFC3339 format
	Version   string // Version identifier (e.g., "AWSCURRENT", "AWSPREVIOUS")
}
