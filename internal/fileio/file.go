// Package fileio provides functions for file input/output operations
package fileio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/env"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	modelenv "github.com/gumi-tsd/secret-env-manager/internal/model/env"
	"github.com/gumi-tsd/secret-env-manager/internal/parser"
)

const secureFilePerm = 0600

// FileType represents the type of file (enumeration)
type FileType int

const (
	EnvFile FileType = iota
	UnknownFile
)

// FileContent represents the content of a file
type FileContent struct {
	Data     []byte
	FilePath string
	Type     FileType
}

// EnvFileOutput represents output file parameters
type EnvFileOutput struct {
	FilePath     string
	Values       map[string]string
	OrderedKeys  []string
	UseQuotes    bool
	NoExpandJson bool
}

// NewEnvFileOutput creates a new EnvFileOutput instance with default settings
func NewEnvFileOutput(filePath string, values map[string]string, orderedKeys []string) EnvFileOutput {
	return NewEnvFileOutputWithOptions(filePath, values, orderedKeys, true)
}

// NewEnvFileOutputWithOptions creates a new EnvFileOutput instance with options
func NewEnvFileOutputWithOptions(filePath string, values map[string]string, orderedKeys []string, useQuotes bool) EnvFileOutput {
	// Create copies of maps and slices to ensure immutability
	valuesCopy := make(map[string]string, len(values))
	for k, v := range values {
		valuesCopy[k] = v
	}

	orderedKeysCopy := make([]string, len(orderedKeys))
	copy(orderedKeysCopy, orderedKeys)

	return EnvFileOutput{
		FilePath:     filePath,
		Values:       valuesCopy,
		OrderedKeys:  orderedKeysCopy,
		UseQuotes:    useQuotes,
		NoExpandJson: false, // Default is to expand JSON
	}
}

// WithValues returns a new EnvFileOutput with the specified values
func (e EnvFileOutput) WithValues(values map[string]string) EnvFileOutput {
	return NewEnvFileOutputWithOptions(e.FilePath, values, e.OrderedKeys, e.UseQuotes)
}

// WithOrderedKeys returns a new EnvFileOutput with the specified ordered keys
func (e EnvFileOutput) WithOrderedKeys(orderedKeys []string) EnvFileOutput {
	return NewEnvFileOutputWithOptions(e.FilePath, e.Values, orderedKeys, e.UseQuotes)
}

// WithUseQuotes returns a new EnvFileOutput with the specified useQuotes setting
func (e EnvFileOutput) WithUseQuotes(useQuotes bool) EnvFileOutput {
	return NewEnvFileOutputWithOptions(e.FilePath, e.Values, e.OrderedKeys, useQuotes)
}

// GenerateCacheFileName creates a cache filename from input filename
// Pure function: Always returns the same output for the same input
func GenerateCacheFileName(inputFileName string) string {
	normalizedPath := filepath.ToSlash(inputFileName)
	safeFileName := strings.ReplaceAll(normalizedPath, "/", "_")

	// Handle filenames that start with a dot (like .env)
	if strings.HasPrefix(safeFileName, ".") {
		// Remove the leading dot and add it back in the correct format
		return fmt.Sprintf(".cache%s", safeFileName)
	} else {
		return fmt.Sprintf(".cache.%s", safeFileName)
	}
}

// ReadFile reads a file and returns its content
// Returns a Result monad containing FileContent
func ReadFile(fileName string) functional.Result[FileContent] {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return functional.Failure[FileContent](
			fmt.Errorf("failed to read file '%s': %w", fileName, err))
	}

	fileType := DetermineFileType(fileName)
	return functional.Success(FileContent{
		Data:     data,
		FilePath: fileName,
		Type:     fileType,
	})
}

// DetermineFileType identifies the type of file based on extension
// Pure function: Always returns the same output for the same input
func DetermineFileType(fileName string) FileType {
	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".env", "":
		return EnvFile
	default:
		return UnknownFile
	}
}

// ParseFileContent parses the content based on its type
// Returns a Result monad containing parsed entries
func ParseFileContent(content FileContent) functional.Result[[]modelenv.Entry] {
	switch content.Type {
	case EnvFile:
		entries, err := parser.ParsePlainFileContent(content.Data)
		if err != nil {
			return functional.Failure[[]modelenv.Entry](
				fmt.Errorf("failed to parse env file '%s': %w", content.FilePath, err))
		}
		return functional.Success(entries)
	default:
		return functional.Failure[[]modelenv.Entry](
			fmt.Errorf("unsupported file format: %s", filepath.Ext(content.FilePath)))
	}
}

// WriteStringToFile writes a string to a file
// Returns a Result monad indicating success or failure
func WriteStringToFile(filePath string, content string) functional.Result[bool] {
	err := os.WriteFile(filePath, []byte(content), secureFilePerm)
	if err != nil {
		return functional.Failure[bool](
			fmt.Errorf("failed to write to output file '%s': %w", filePath, err))
	}
	return functional.Success(true)
}

// SecureOutputFile sets owner-only permissions on output file
// Returns a Result monad indicating success or failure
func SecureOutputFile(outputFileName string) functional.Result[bool] {
	if err := os.Chmod(outputFileName, secureFilePerm); err != nil {
		return functional.Failure[bool](
			fmt.Errorf("failed to set secure permissions on output file '%s': %w", outputFileName, err))
	}
	return functional.Success(true)
}

// WriteOutputFile writes env variables to a file
// Returns a Result monad indicating success or failure
func WriteOutputFile(output EnvFileOutput) functional.Result[bool] {
	// Format options using env.EnvVarOptions
	options := env.EnvVarOptions{
		UseQuotes:    output.UseQuotes,
		SortedKeys:   output.OrderedKeys,
		IncludeURIs:  false,
		NoExpandJson: output.NoExpandJson, // Pass the JSON expansion setting
	}

	// Delegate formatting to env package function
	formatResult := env.FormatEnvVarContent(output.Values, options)
	if formatResult.IsFailure() {
		return functional.Failure[bool](formatResult.GetError())
	}

	// Handle warnings if any
	result := formatResult.Unwrap()
	for _, warning := range result.Warnings {
		fmt.Printf("Warning: %s\n", warning)
	}

	// Write formatted content to file
	return WriteStringToFile(output.FilePath, result.Content)
}

// ReadAndParseInputFile reads and parses the input file
// Uses function composition for a more functional approach
func ReadAndParseInputFile(fileName string) functional.Result[[]modelenv.Entry] {
	return functional.Chain(
		ReadFile(fileName),
		ParseFileContent,
	)
}

// ReadEnvVarsAsMap reads environment variables from a file and converts to a map
// Uses function composition for a more functional approach
func ReadEnvVarsAsMap(fileName string) functional.Result[map[string]string] {
	return functional.MapResultTo(
		ReadAndParseInputFile(fileName),
		env.EnvsToMap,
	)
}

// ReadEnvVarsFromFile reads environment variables from a file
// Compatibility version that returns unwrapped result and error
func ReadEnvVarsFromFile(fileName string) (map[string]string, error) {
	result := ReadEnvVarsAsMap(fileName)
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.Unwrap(), nil
}
