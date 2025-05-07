// Package cmd implements command-line commands for the secret-env-manager
package cmd

import (
	"fmt"

	"github.com/gumi-tsd/secret-env-manager/internal/env"
	"github.com/gumi-tsd/secret-env-manager/internal/fileio"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	modelenv "github.com/gumi-tsd/secret-env-manager/internal/model/env"
	"github.com/gumi-tsd/secret-env-manager/internal/provider"
	"github.com/urfave/cli/v2"
)

// UpdateParams contains parameters for the Update command
type UpdateParams struct {
	InputFileName string
	EndpointURL   string
	NoQuotes      bool
	NoExpandJson  bool
}

// WithUpdateParams creates a new UpdateParams with provided values
// Pure function: Always returns the same output for the same input
func WithUpdateParams(inputFileName, endpointURL string, noQuotes, noExpandJson bool) UpdateParams {
	return UpdateParams{
		InputFileName: inputFileName,
		EndpointURL:   endpointURL,
		NoQuotes:      noQuotes,
		NoExpandJson:  noExpandJson,
	}
}

// UpdateResult represents the result of an update operation
type UpdateResult struct {
	OutputFileName string
	EntryCount     int
}

// WithUpdateResult creates a new UpdateResult
// Pure function: Always returns the same output for the same input
func WithUpdateResult(outputFileName string, entryCount int) UpdateResult {
	return UpdateResult{
		OutputFileName: outputFileName,
		EntryCount:     entryCount,
	}
}

// AcquiredSecrets holds the results of secret acquisition
type AcquiredSecrets struct {
	Values map[string]string
	Keys   []string
}

// WithAcquiredSecrets creates a new AcquiredSecrets
// Pure function: Always returns the same output for the same input
func WithAcquiredSecrets(values map[string]string, keys []string) AcquiredSecrets {
	return AcquiredSecrets{
		Values: values,
		Keys:   keys,
	}
}

// Update retrieves secrets from cloud providers and updates the environment variables
// Entry point function for CLI that gets necessary information from cli.Context and executes processing
func Update(c *cli.Context) error {
	// Validate input parameters
	paramsResult := validateUpdateParams(c)
	if paramsResult.IsFailure() {
		return paramsResult.GetError()
	}

	// Execute update process
	result := performUpdate(paramsResult.Unwrap())
	if result.IsFailure() {
		return result.GetError()
	}

	// Display success message
	updateResult := result.Unwrap()
	logSuccessInfo(fmt.Sprintf("Successfully updated %d environment variables in %s",
		updateResult.EntryCount, updateResult.OutputFileName))

	return nil
}

// validateUpdateParams validates CLI parameters and returns a Result monad
// Pure function: Validates command-line arguments and returns a parameter object
func validateUpdateParams(c *cli.Context) functional.Result[UpdateParams] {
	inputFileName := c.String("input")
	if inputFileName == "" {
		return withFailure[UpdateParams]("input file path required (-i or --input)")
	}

	endpointURL := c.String("endpoint-url")
	noQuotes := c.Bool("no-quotes")
	noExpandJson := c.Bool("no-expand-json")

	return withSuccess(WithUpdateParams(
		inputFileName,
		endpointURL,
		noQuotes,
		noExpandJson,
	))
}

// performUpdate executes the update process using function composition and Result monad
// Composite function: Chains each processing step and returns the result
func performUpdate(params UpdateParams) functional.Result[UpdateResult] {
	outputFileName := fileio.GenerateCacheFileName(params.InputFileName)

	// Log the input file name being processed
	logInfoMsg(fmt.Sprintf("Reading input file: %s", params.InputFileName))

	// Step 1: Ensure the file is git-ignored
	gitIgnoreResult := fileio.IsFileIgnored(outputFileName)
	if gitIgnoreResult.IsFailure() {
		return withFailure[UpdateResult](gitIgnoreResult.GetError().Error())
	}

	// If the file is not git-ignored, fail with error
	if !gitIgnoreResult.Unwrap() {
		fileio.DisplaySecurityWarning(outputFileName)
		return withFailure[UpdateResult](
			fmt.Sprintf("output file '%s' is not ignored by git, which poses a security risk",
				outputFileName))
	}

	// Step 2: Read and parse input file
	entriesResult := readInputFile(params.InputFileName)
	if entriesResult.IsFailure() {
		return withFailure[UpdateResult](entriesResult.GetError().Error())
	}
	entries := entriesResult.Unwrap()

	// Log entries for debugging
	logDebugInfo(fmt.Sprintf("Found %d entries in input file", len(entries)))

	// Step 3: Acquire secrets - NoExpandJsonパラメータを正しく渡す
	secretsResult := acquireSecrets(entries, params.EndpointURL, params.NoExpandJson)
	if secretsResult.IsFailure() {
		return withFailure[UpdateResult](secretsResult.GetError().Error())
	}
	secrets := secretsResult.Unwrap()

	// Step 4: Write to output file
	writeResult := writeOutputFile(outputFileName, secrets.Values, secrets.Keys, params.NoQuotes, params.NoExpandJson)
	if writeResult.IsFailure() {
		return withFailure[UpdateResult](writeResult.GetError().Error())
	}

	// Step 5: Secure output file (only show warning on failure)
	secureResult := secureOutputFileWithWarning(outputFileName)
	if secureResult.IsFailure() {
		// This should never happen due to how secureOutputFileWithWarning is implemented
		return withFailure[UpdateResult](secureResult.GetError().Error())
	}

	// Step 6: Create and return result
	return withSuccess(WithUpdateResult(
		outputFileName,
		len(entries),
	))
}

// readInputFile reads and parses the input file, returning env entries
// Composite function: Combines file reading and parsing
func readInputFile(fileName string) functional.Result[[]modelenv.Entry] {
	// Use function chaining pattern
	return functional.Chain(
		fileio.ReadFile(fileName),
		fileio.ParseFileContent,
	)
}

// acquireSecrets fetches secrets from providers and organizes them by key
// Composite function: Retrieves secrets and organizes key order
func acquireSecrets(entries []modelenv.Entry, endpointURL string, noExpandJson bool) functional.Result[AcquiredSecrets] {
	// Create provider configuration with endpoint URL and JSON expansion setting
	config := provider.NewProviderConfig(endpointURL)
	config.NoExpandJson = noExpandJson

	providers := provider.CreateProviderMap(config)

	// Use provider's ProcessEntriesResult for secret processing
	processResult := provider.ProcessEntriesResult(entries, providers)
	if !processResult.IsSuccess() {
		return withFailure[AcquiredSecrets](processResult.Error.Error())
	}

	// Extract values and organize keys in order using the utility function
	values := processResult.Values
	orderedKeys := env.OrganizeKeyOrder(entries, values)

	return withSuccess(WithAcquiredSecrets(
		values,
		orderedKeys,
	))
}

// writeOutputFile writes environment variables to a file
// Composite function: Writes environment variables to a file
func writeOutputFile(fileName string, values map[string]string, orderedKeys []string, noQuotes bool, noExpandJson bool) functional.Result[bool] {
	// Create output file with custom options
	output := fileio.NewEnvFileOutputWithOptions(
		fileName,
		values,
		orderedKeys,
		!noQuotes, // If noQuotes is true, UseQuotes is false
	)

	// Set JSON expansion flag
	output.NoExpandJson = noExpandJson

	return fileio.WriteOutputFile(output)
}

// secureOutputFileWithWarning sets appropriate permissions, converting errors to warnings
// Composite function: Continues with only a warning if file permission setting fails
func secureOutputFileWithWarning(fileName string) functional.Result[bool] {
	result := fileio.SecureOutputFile(fileName)
	if result.IsFailure() {
		logWarning(fmt.Sprintf("Unable to secure output file: %v", result.GetError()))
		// Even if securing fails, consider it a non-critical warning
		return withSuccess(true)
	}
	return result
}
