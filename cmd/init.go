// Package cmd implements command-line commands for the secret-env-manager
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/gumi-tsd/secret-env-manager/internal/formatting"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/uri"
	"github.com/gumi-tsd/secret-env-manager/internal/provider/aws"
	"github.com/gumi-tsd/secret-env-manager/internal/provider/googlecloud"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

// Provider represents a secret provider type
type Provider string

const (
	// AWSProvider represents AWS Secrets Manager
	AWSProvider Provider = "aws"
	// GoogleCloudProvider represents Google Cloud Secret Manager
	GoogleCloudProvider Provider = "googlecloud"
)

// EnvParams holds the validated environment parameters
type EnvParams struct {
	// AWS parameters
	AwsProfile  string
	AwsRegion   string
	EndpointURL string

	// Google Cloud parameters
	GoogleCloudProjectID string

	// Selected provider
	Provider Provider
}

// WithEnvParams creates a new EnvParams with provided values
// Pure function: Always returns the same output for the same input
func WithEnvParams(provider Provider, awsProfile, awsRegion, endpointURL, googleCloudProjectID string) *EnvParams {
	return &EnvParams{
		Provider:             provider,
		AwsProfile:           awsProfile,
		AwsRegion:            awsRegion,
		EndpointURL:          endpointURL,
		GoogleCloudProjectID: googleCloudProjectID,
	}
}

// InitResult represents the result of an init operation
type InitResult struct {
	SecretNames          []string
	AwsProfile           string
	AwsRegion            string
	GoogleCloudProjectID string
	Provider             Provider
}

// Init initializes the environment by listing secrets from cloud providers and allows interactive selection
// It supports both AWS Secrets Manager and Google Cloud Secret Manager.
// For AWS, it requires AWS_PROFILE and AWS_REGION environment variables.
// For Google Cloud, it requires GOOGLE_CLOUD_PROJECT environment variable.
// Custom AWS endpoints can be specified using the --endpoint-url flag.
func Init(c *cli.Context) error {
	// Get endpoint URL from flag
	endpointURL := c.String("endpoint-url")

	// Select provider
	providerResult := selectProvider()
	if providerResult.IsFailure() {
		return providerResult.GetError()
	}
	provider := providerResult.Unwrap()

	// Validate environment based on selected provider
	envResult := validateEnvironment(provider, endpointURL)
	if envResult.IsFailure() {
		return envResult.GetError()
	}

	params := envResult.Unwrap()

	// Log current context based on provider
	var logCtxMsg string
	if params.Provider == AWSProvider {
		logCtxMsg = fmt.Sprintf("Listing secrets for AWS profile '%s' in region '%s'",
			params.AwsProfile, params.AwsRegion)

		// Add endpoint URL information if specified
		if params.EndpointURL != "" {
			logCtxMsg += fmt.Sprintf(" using endpoint URL '%s'", params.EndpointURL)
		}
	} else {
		logCtxMsg = fmt.Sprintf("Listing secrets for Google Cloud project '%s'", params.GoogleCloudProjectID)
	}
	logInfoMsg(logCtxMsg + "...")

	// List secrets from the selected provider
	var secretsResult functional.Result[[]string]
	if params.Provider == AWSProvider {
		secretsResult = listAwsSecrets(params.AwsProfile, params.AwsRegion, params.EndpointURL)
	} else {
		secretsResult = listGoogleCloudSecrets(params.GoogleCloudProjectID)
	}

	if secretsResult.IsFailure() {
		return secretsResult.GetError()
	}

	secretNames := secretsResult.Unwrap()
	if len(secretNames) == 0 {
		logInfoMsg("No secrets found.")
		return nil
	}

	// Get selected secrets through interactive prompt
	selectResult := selectSecretsResult(secretNames)
	if selectResult.IsFailure() {
		return selectResult.GetError()
	}

	selectedSecrets := selectResult.Unwrap()

	// Display selected secrets with nice formatting
	outputSelectedSecrets(selectedSecrets, params)

	// Display help for next steps
	displayNextSteps()

	return nil
}

// selectProvider prompts the user to select a secret provider (AWS or Google Cloud)
// Returns a Result monad with the selected provider
func selectProvider() functional.Result[Provider] {
	prompt := promptui.Select{
		Label: "Select a secret provider",
		Items: []string{"AWS Secrets Manager", "Google Cloud Secret Manager"},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return functional.Failure[Provider](fmt.Errorf("provider selection failed: %w", err))
	}

	if idx == 0 {
		return functional.Success(AWSProvider)
	}
	return functional.Success(GoogleCloudProvider)
}

// validateEnvironment validates required environment variables based on selected provider
// For AWS: checks AWS_PROFILE and AWS_REGION
// For Google Cloud: checks GOOGLE_CLOUD_PROJECT
// Pure function: Returns a Result monad with environment parameters
func validateEnvironment(provider Provider, endpointURL string) functional.Result[*EnvParams] {
	if provider == AWSProvider {
		// Check AWS_PROFILE with Option monad
		profileOption := getEnvOption("AWS_PROFILE")
		if profileOption.IsNone() {
			return withFailure[*EnvParams]("AWS_PROFILE environment variable is not set")
		}

		// Check AWS_REGION with Option monad
		regionOption := getEnvOption("AWS_REGION")
		if regionOption.IsNone() {
			return withFailure[*EnvParams]("AWS_REGION environment variable is not set")
		}

		return withSuccess(WithEnvParams(
			provider,
			profileOption.Unwrap(),
			regionOption.Unwrap(),
			endpointURL,
			"", // No Google Cloud project ID needed
		))
	} else {
		// Check GOOGLE_CLOUD_PROJECT with Option monad
		projectOption := getEnvOption("GOOGLE_CLOUD_PROJECT")
		if projectOption.IsNone() {
			return withFailure[*EnvParams]("GOOGLE_CLOUD_PROJECT environment variable is not set")
		}

		return withSuccess(WithEnvParams(
			provider,
			"", // No AWS profile needed
			"", // No AWS region needed
			"", // No endpoint URL needed
			projectOption.Unwrap(),
		))
	}
}

// getEnvOption gets an environment variable as an Option
// Pure function: Transforms environment lookups to Option monad
func getEnvOption(name string) functional.Option[string] {
	value, exists := os.LookupEnv(name)
	if !exists || value == "" {
		return functional.None[string]()
	}
	return functional.Some(value)
}

// listAwsSecrets retrieves secrets from AWS Secrets Manager
// Pure function: Returns a Result monad with secret names
func listAwsSecrets(awsProfile, awsRegion, endpointURL string) functional.Result[[]string] {
	ctx := context.Background()

	// Create a custom AWS Provider
	provider := aws.NewAwsProvider()

	// Handle custom endpoint if specified
	if endpointURL != "" {
		return listAwsSecretsWithEndpoint(ctx, provider, awsProfile, awsRegion, endpointURL)
	} else {
		// Standard case (no endpoint URL)
		return provider.ListSecretsResult(ctx, awsProfile, awsRegion)
	}
}

// listGoogleCloudSecrets retrieves secrets from Google Cloud Secret Manager
// Pure function: Returns a Result monad with secret names
func listGoogleCloudSecrets(projectID string) functional.Result[[]string] {
	ctx := context.Background()

	// Create a Google Cloud Provider
	provider := googlecloud.NewGoogleCloudProvider()

	// Call Google Cloud API
	return provider.ListSecretsResult(ctx, projectID, "")
}

// listAwsSecretsWithEndpoint retrieves secrets using a custom endpoint
// Pure function: Returns a Result monad with secret names
func listAwsSecretsWithEndpoint(ctx context.Context, provider *aws.AwsProvider,
	awsProfile, awsRegion, endpointURL string) functional.Result[[]string] {

	// Get client with custom endpoint
	clientResult := functional.TryCatch(func() (*secretsmanager.Client, error) {
		return provider.GetClient(ctx, awsProfile, awsRegion, endpointURL)
	})

	if clientResult.IsFailure() {
		return withFailure[[]string](
			fmt.Sprintf("failed to get AWS client for profile '%s' in region '%s' with endpoint '%s': %v",
				awsProfile, awsRegion, endpointURL, clientResult.GetError()))
	}

	client := clientResult.Unwrap()

	// Create API input
	input := &secretsmanager.ListSecretsInput{
		MaxResults: nil, // Use default limit
	}

	// Call API and transform response
	secretsResult := functional.TryCatch(func() (*secretsmanager.ListSecretsOutput, error) {
		return client.ListSecrets(ctx, input)
	})

	if secretsResult.IsFailure() {
		return withFailure[[]string](
			fmt.Sprintf("AWS ListSecrets API error for profile '%s' in region '%s' with endpoint '%s': %v",
				awsProfile, awsRegion, endpointURL, secretsResult.GetError()))
	}

	// Extract secret names
	result := secretsResult.Unwrap()
	secrets := extractSecretNames(result.SecretList)

	return withSuccess(secrets)
}

// extractSecretNames extracts names from secret list
// Pure function: Returns a list of secret names
func extractSecretNames(secretList []types.SecretListEntry) []string {
	secrets := make([]string, 0, len(secretList))
	for _, secret := range secretList {
		if secret.Name != nil {
			secrets = append(secrets, *secret.Name)
		}
	}
	return secrets
}

// selectSecretsResult shows an interactive select prompt for secrets
// Returns a Result monad for consistent error handling
func selectSecretsResult(secretNames []string) functional.Result[[]string] {
	if len(secretNames) == 0 {
		return withSuccess([]string{})
	}

	// Log information about selection
	logInfoMsg(formatting.Info("Found %d secrets. Use arrow keys to navigate, Enter to select a secret.", len(secretNames)))

	// Wrap the selection process in a Result monad
	return functional.TryCatch(func() ([]string, error) {
		return selectSecrets(secretNames)
	})
}

// selectSecrets shows an interactive select prompt for secrets
// Side effect function: Interacts with the user through terminal
func selectSecrets(secretNames []string) ([]string, error) {
	if len(secretNames) == 0 {
		return []string{}, nil
	}

	prompt := promptui.Select{
		Label: "Select a secret (use arrow keys, press Enter to select)",
		Items: secretNames,
		Size:  20, // Show up to 20 items at once
	}

	// For multi-select we'd need to handle it manually since promptui doesn't support it directly
	// For now, we'll just do single select
	selectedIdx, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("secret selection failed: %w", err)
	}

	// Return only the selected secret
	return []string{secretNames[selectedIdx]}, nil
}

// outputSelectedSecrets prints the selected secret URIs to standard output
// Side effect function: Prints to standard output
func outputSelectedSecrets(secretNames []string, params *EnvParams) {
	if len(secretNames) == 0 {
		return
	}

	fmt.Println(formatting.FormatHeader("\nSelected Secret URI"))
	fmt.Println(formatting.Hint("Copy the following URI to your environment file:"))
	fmt.Println()

	for _, secretName := range secretNames {
		var secretUri uri.SecretURI

		if params.Provider == AWSProvider {
			secretUri = buildAwsSecretURI(secretName, params.AwsProfile, params.AwsRegion)
		} else {
			secretUri = buildGoogleCloudSecretURI(secretName, params.GoogleCloudProjectID)
		}

		fmt.Println(formatting.ColorizeValue(secretUri.GetUri()))
	}
	fmt.Println()
}

// displayNextSteps shows the user what to do next
// Side effect function: Prints to standard output
func displayNextSteps() {
	fmt.Println(formatting.FormatHeader("\nNext Steps"))
	fmt.Println(formatting.Hint("1. Create or edit your environment file with any name (e.g. .env, dev.env, prod.env) and add the secret URI"))
	fmt.Println(formatting.Hint("   You can use any filename that makes sense for your project environment"))
	fmt.Println()
	fmt.Println(formatting.Hint("2. Run the following command to fetch and cache your secrets:"))
	fmt.Println()
	fmt.Println(formatting.ColorizeValue("   sem update --input your-env-file"))
	fmt.Println()
	fmt.Println(formatting.Hint("3. Then load the environment variables with one of these methods:"))
	fmt.Println()
	fmt.Println(formatting.ColorizeValue("   # Using env -S"))
	fmt.Println(formatting.ColorizeValue("   env -S \"sem load -i your-env-file\" your-command"))
	fmt.Println()

	fmt.Println(formatting.FormatHeader("\nUsing with direnv"))
	fmt.Println(formatting.Hint("With direnv, environment variables will be automatically set when you enter the project directory."))
	fmt.Println(formatting.Hint("Create a .envrc file in your root directory and add the following:"))
	fmt.Println()
	fmt.Println(formatting.ColorizeValue("   # Update environment variable cache"))
	fmt.Println(formatting.ColorizeValue("   sem update -i your-env-file"))
	fmt.Println()
	fmt.Println(formatting.ColorizeValue("   # Load environment variables from the cache file"))
	fmt.Println(formatting.ColorizeValue("   dotenv .cache.your-env-file"))
	fmt.Println()
	fmt.Println(formatting.Hint("Then run 'direnv allow' to apply the changes."))
	fmt.Println(formatting.Hint("This will automatically load the environment variables whenever you enter the directory."))
	fmt.Println()

	fmt.Println(formatting.FormatHeader("\nFile Management Best Practices"))
	fmt.Println(formatting.Hint("- The environment file (containing only URIs) is safe to commit to version control"))
	fmt.Println(formatting.Hint("- We recommend tracking these files in Git for easier environment configuration"))
	fmt.Println(formatting.Hint("- Access to actual secrets is controlled by cloud provider permissions"))
	fmt.Println(formatting.Hint("- The cache file (.cache.your-env-file) contains actual secrets and should be in .gitignore"))
	fmt.Println()
}

// buildAwsSecretURI creates a SecretURI for the given AWS secret name
// Pure function: Always returns the same output for the same input
func buildAwsSecretURI(secretName, awsProfile, awsRegion string) uri.SecretURI {
	return uri.SecretURI{
		Platform:   string(AWSProvider),
		Account:    awsProfile,
		Service:    "secretsmanager",
		SecretName: secretName,
		Version:    "AWSCURRENT", // Default AWS version
		Region:     awsRegion,
	}
}

// buildGcpSecretURI creates a SecretURI for the given Google Cloud secret name
// Pure function: Always returns the same output for the same input
func buildGoogleCloudSecretURI(secretName, projectID string) uri.SecretURI {
	return uri.SecretURI{
		Platform:   string(GoogleCloudProvider),
		Account:    projectID,
		Service:    "secretmanager",
		SecretName: secretName,
		Version:    "latest", // Default Google Cloud version
	}
}
