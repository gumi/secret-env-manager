// Package googlecloud provides functionality for interacting with Google Cloud Secret Manager.
package googlecloud

import (
	"context"
	"fmt"
	"strings"

	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"google.golang.org/api/iterator"
)

// ListSecretsResult lists all secrets in a specific GCP project
func (p *GoogleCloudProvider) ListSecretsResult(ctx context.Context, projectID string, _ string) functional.Result[[]string] {
	// Get client for listing secrets
	client, err := p.GetClient(ctx)
	if err != nil {
		return functional.Failure[[]string](fmt.Errorf("failed to get Google Cloud client: %w", err))
	}

	// Construct the parent resource name
	parent := fmt.Sprintf("projects/%s", projectID)

	// Build the request
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
	}

	fmt.Println("Listing secrets for project:", projectID)

	// Call Google Cloud API
	it := client.ListSecrets(ctx, req)
	secrets := []string{}

	// Iterate through all secrets
	for {
		secret, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return functional.Failure[[]string](fmt.Errorf("error iterating secrets: %w", err))
		}

		// Extract the secret name from the full resource name
		name := secret.Name
		parts := extractSecretNameFromResourceName(name)
		if parts != "" {
			secrets = append(secrets, parts)
		}
	}

	return functional.Success(secrets)
}

// ListSecrets lists all secrets in a specific GCP project
func (p *GoogleCloudProvider) ListSecrets(ctx context.Context, projectID string, region string) ([]string, error) {
	result := p.ListSecretsResult(ctx, projectID, region)
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.Unwrap(), nil
}

// extractSecretNameFromResourceName extracts the secret name from a resource name
// Format: projects/{project}/secrets/{secret}
func extractSecretNameFromResourceName(resourceName string) string {
	// Find the last segment after "secrets/"
	const prefix = "secrets/"
	idx := strings.LastIndex(resourceName, prefix)
	if idx < 0 {
		return ""
	}
	return resourceName[idx+len(prefix):]
}
