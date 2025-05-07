// Package secret provides utilities for handling secrets and their values.
package secret

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/logging"
	"github.com/gumi-tsd/secret-env-manager/internal/model/uri"
	"github.com/gumi-tsd/secret-env-manager/internal/text"
)

// Logger instance for this package
var logger = logging.DefaultLogger()

// LogInfoMsg logs information messages with consistent formatting
func LogInfoMsg(message string) {
	logger.Info(message)
}

// Error definitions
var (
	// ErrEmptySecretValue indicates that the secret value is empty
	ErrEmptySecretValue = errors.New("secret value is empty")

	// ErrInvalidSecretValue indicates that the secret value is invalid
	ErrInvalidSecretValue = errors.New("invalid secret value format")

	// ErrKeyNotFound indicates that the requested key was not found in the secret
	ErrKeyNotFound = errors.New("key not found in secret")
)

// ValueOptions configures how secrets are processed
type ValueOptions struct {
	CleanControlChars bool // Whether to clean control characters from values
}

// NewValueOptions creates a new ValueOptions with custom settings
func NewValueOptions(cleanControlChars bool) ValueOptions {
	return ValueOptions{
		CleanControlChars: cleanControlChars,
	}
}

// DefaultValueOptions returns the default options for secret value processing
func DefaultValueOptions() ValueOptions {
	return NewValueOptions(true)
}

// ParseValueResult extracts a field from a secret string with Result monad
func ParseValueResult(secretString string, key string) functional.Result[string] {
	return ParseValueWithOptionsResult(secretString, key, DefaultValueOptions())
}

// ParseValueWithOptionsResult extracts a field with custom options, returns Result monad
func ParseValueWithOptionsResult(secretString string, key string, options ValueOptions) functional.Result[string] {
	// Handle empty secrets
	if secretString == "" {
		return functional.Failure[string](ErrEmptySecretValue)
	}

	// If there's no specific key, just return the whole secret
	if key == "" {
		rawValue := secretString
		if options.CleanControlChars {
			rawValue = text.CleanControlChars(rawValue)
		}
		return functional.Success(rawValue)
	}

	// Try to parse as JSON first
	jsonResult := text.ParseJSONResult(secretString)

	// If parse fails and a key was specified, return an error
	if jsonResult.IsFailure() {
		// Return error when key is specified for non-JSON data
		return functional.Failure[string](
			fmt.Errorf("%w: cannot retrieve key '%s' because the specified secret is not in JSON format",
				ErrInvalidSecretValue, key))
	}

	// JSON parsed successfully, now extract the requested key
	jsonData := jsonResult.Unwrap()

	// Split the key path if it contains dots
	keyParts := strings.Split(key, ".")

	// Navigate to the specified key path in the JSON data
	valueResult := text.NavigatePathResult(jsonData, keyParts)
	if valueResult.IsFailure() {
		return functional.Failure[string](fmt.Errorf("%w: %v", ErrKeyNotFound, valueResult.GetError()))
	}

	// Format the extracted value as a string
	formattedValue, err := formatSecretValue(valueResult.Unwrap(), options)
	if err != nil {
		return functional.Failure[string](fmt.Errorf("failed to format secret value: %w", err))
	}

	return functional.Success(formattedValue)
}

// formatSecretValue converts a value to appropriate string representation based on options
func formatSecretValue(value interface{}, options ValueOptions) (string, error) {
	var result string

	switch v := value.(type) {
	case string:
		result = v
		if options.CleanControlChars {
			result = text.CleanControlChars(result)
		}

	case float64, int, bool:
		result = fmt.Sprintf("%v", v)

	default:
		// For complex types (objects, arrays), convert to JSON string
		jsonBytes, err := text.MarshalToJSONString(value)
		if err != nil {
			return "", err
		}
		result = jsonBytes
	}

	return result, nil
}

// FormatValueResult converts a value to appropriate string representation using default options and returns Result monad.
func FormatValueResult(value interface{}) functional.Result[string] {
	result, err := formatSecretValue(value, DefaultValueOptions())
	if err != nil {
		return functional.Failure[string](fmt.Errorf("failed to format value: %w", err))
	}
	return functional.Success(result)
}

// FormatValue converts a value to appropriate string representation using default options.
func FormatValue(value interface{}) string {
	result, err := formatSecretValue(value, DefaultValueOptions())
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	return result
}

// IsURI checks if a string is a secret URI (starts with sem://)
func IsURI(s string) bool {
	return strings.HasPrefix(s, "sem://")
}

// BuildCacheKey creates a consistent cache key from secret information
func BuildCacheKey(account, service, secretName, version, region string) string {
	return uri.BuildCacheKey(account, service, secretName, version, region)
}

// ParseValue is a non-monadic wrapper for ParseValueResult
func ParseValue(secretString string, key string) (string, error) {
	result := ParseValueResult(secretString, key)
	if result.IsFailure() {
		return "", result.GetError()
	}
	return result.Unwrap(), nil
}

// ParseValueWithOptions is a non-monadic wrapper for ParseValueWithOptionsResult
func ParseValueWithOptions(secretString string, key string, options ValueOptions) (string, error) {
	result := ParseValueWithOptionsResult(secretString, key, options)
	if result.IsFailure() {
		return "", result.GetError()
	}
	return result.Unwrap(), nil
}

// WithCleanControlChars returns a new ValueOptions with the CleanControlChars setting modified
func (o ValueOptions) WithCleanControlChars(clean bool) ValueOptions {
	return ValueOptions{
		CleanControlChars: clean,
	}
}

// ValidateJSON validates that a string is valid JSON and returns a Result monad
func ValidateJSON(s string) functional.Result[bool] {
	if text.IsJSON(s) {
		return functional.Success(true)
	}
	return functional.Failure[bool](ErrInvalidSecretValue)
}

// ExtractKeyFromPath extracts the last segment from a dot-separated path
func ExtractKeyFromPath(path string) string {
	parts := strings.Split(path, ".")
	return parts[len(parts)-1]
}

// CombinePaths joins path segments with dots
func CombinePaths(segments ...string) string {
	return strings.Join(functional.Filter(segments, func(s string) bool {
		return s != ""
	}), ".")
}
