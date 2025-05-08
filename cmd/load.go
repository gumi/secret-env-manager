// Package cmd implements command-line commands for the secret-env-manager
package cmd

import (
	"fmt"
	"os"

	"github.com/gumi-tsd/secret-env-manager/internal/env"
	"github.com/gumi-tsd/secret-env-manager/internal/fileio"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/urfave/cli/v2"
)

// LoadParams contains parameters for the Load command
type LoadParams struct {
	InputFileName   string
	OutputFileName  string
	ExportOnlyUnset bool
}

// WithLoadParams creates a new LoadParams with provided values
func WithLoadParams(inputFileName, outputFileName string, exportOnlyUnset bool) LoadParams {
	return LoadParams{
		InputFileName:   inputFileName,
		OutputFileName:  outputFileName,
		ExportOnlyUnset: exportOnlyUnset,
	}
}

// LoadResult represents the result of a load operation
type LoadResult struct {
	Lines           []string
	ExportStatement string
	OutputFileName  string
	EnvVarsCount    int
}

// WithLoadResult creates a new LoadResult with provided values
func WithLoadResult(lines []string, exportStatement, outputFileName string, envVarsCount int) LoadResult {
	return LoadResult{
		Lines:           lines,
		ExportStatement: exportStatement,
		OutputFileName:  outputFileName,
		EnvVarsCount:    envVarsCount,
	}
}

// Load loads environment variables from a file and generates a shell export statement.
func Load(c *cli.Context) error {
	// Validate parameters
	paramsResult := validateLoadParams(c)
	if paramsResult.IsFailure() {
		return paramsResult.GetError()
	}

	// Load environment variables
	loadResult := loadEnvVars(paramsResult.Unwrap())
	if loadResult.IsFailure() {
		return loadResult.GetError()
	}

	result := loadResult.Unwrap()

	// Handle result
	if c.Bool("with-export") {
		// If with-export is specified, print the export statement
		fmt.Println(result.ExportStatement)
	} else {
		// Otherwise print the variables without export statements
		for _, line := range result.Lines {
			fmt.Println(line)
		}
	}

	return nil
}

// validateLoadParams validates CLI parameters with Result monad
func validateLoadParams(c *cli.Context) functional.Result[LoadParams] {
	inputFileName := c.String("input")
	if inputFileName == "" {
		return withFailure[LoadParams]("input file path required (-i or --input)")
	}

	// Handle output file parameter
	outputFileName := c.String("output")
	if outputFileName == "" {
		// If no output file is specified, use the cache file
		outputFileName = fileio.GenerateCacheFileName(inputFileName)
	}

	// Whether to export only variables that are not already set in the environment
	exportOnlyUnset := c.Bool("only-unset")

	return withSuccess(WithLoadParams(
		inputFileName,
		outputFileName,
		exportOnlyUnset,
	))
}

// loadEnvVars loads environment variables from a file using Result monad
func loadEnvVars(params LoadParams) functional.Result[LoadResult] {
	// Read variables from file
	varsResult := readEnvVarsFromFile(params.OutputFileName)
	if varsResult.IsFailure() {
		// Convert error to LoadResult error
		return functional.Failure[LoadResult](varsResult.GetError())
	}

	variables := varsResult.Unwrap()

	// Filter variables if only-unset is specified
	if params.ExportOnlyUnset {
		variables = filterUnsetVariables(variables)
	}

	// Generate export statement and individual lines
	// Generate with export prefix for export statement, but not for individual lines
	exportWithPrefix, _ := env.FormatEnvVars(variables, true, false)
	_, linesWithoutPrefix := env.FormatEnvVars(variables, false, false)

	return withSuccess(WithLoadResult(
		linesWithoutPrefix,
		exportWithPrefix,
		params.OutputFileName,
		len(variables),
	))
}

// readEnvVarsFromFile reads environment variables from a file
func readEnvVarsFromFile(fileName string) functional.Result[map[string]string] {
	variables, err := fileio.ReadEnvVarsFromFile(fileName)
	if err != nil {
		return withFailure[map[string]string](fmt.Sprintf(
			"failed to read environment variables from %q: %v",
			fileName, err))
	}
	return withSuccess(variables)
}

// filterUnsetVariables filters out variables that are already set in the environment
func filterUnsetVariables(variables map[string]string) map[string]string {
	filteredVars := make(map[string]string)
	for key, value := range variables {
		_, exists := os.LookupEnv(key)
		if !exists {
			filteredVars[key] = value
		}
	}
	return filteredVars
}
