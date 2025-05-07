package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// Secrets is a collection of Secret objects
type Secrets struct {
	Secrets []Secret `json:"secrets"`
}

// ListSecretsResult lists all secrets in a specific AWS account and region
func (p *AwsProvider) ListSecretsResult(ctx context.Context, account string, region string) functional.Result[[]string] {
	// Get client for listing secrets
	client, err := p.GetClient(ctx, account, region, "")
	if err != nil {
		return functional.Failure[[]string](fmt.Errorf("failed to get AWS client: %w", err))
	}

	// Call AWS API
	input := &secretsmanager.ListSecretsInput{
		MaxResults: nil, // Use default limit
	}

	fmt.Println("Listing secrets for account:", account, "in region:", region)

	result, err := client.ListSecrets(ctx, input)
	if err != nil {
		return functional.Failure[[]string](fmt.Errorf("AWS ListSecrets API error: %w", err))
	}

	// Extract the names
	secrets := make([]string, 0, len(result.SecretList))
	for _, secret := range result.SecretList {
		if secret.Name != nil {
			secrets = append(secrets, *secret.Name)
		}
	}

	return functional.Success(secrets)
}

// ListSecrets lists all secrets in a specific AWS account and region
func (p *AwsProvider) ListSecrets(ctx context.Context, account string, region string) ([]string, error) {
	result := p.ListSecretsResult(ctx, account, region)
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.Unwrap(), nil
}

// ListSecretVersions lists all the versions for a specific secret
func (p *AwsProvider) ListSecretVersions(ctx context.Context, profile string, region string, secretName string) ([]string, error) {
	return p.ListSecretVersionsWithEndpoint(ctx, profile, region, secretName, "")
}

// ListSecretVersionsWithEndpoint lists all the versions for a specific secret with a custom endpoint
func (p *AwsProvider) ListSecretVersionsWithEndpoint(ctx context.Context, profile string, region string, secretName string, endpoint string) ([]string, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}

	client, err := p.GetClient(ctx, profile, region, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS client: %w", err)
	}

	// Create input for ListSecretVersionIds
	input := &secretsmanager.ListSecretVersionIdsInput{
		SecretId:   aws.String(secretName),
		MaxResults: aws.Int32(100), // Adjust as needed
	}

	result, err := client.ListSecretVersionIds(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list secret versions: %w", err)
	}

	versions := []string{}
	for _, version := range result.Versions {
		versions = append(versions, version.VersionStages...)
	}

	return versions, nil
}
