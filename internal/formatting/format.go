// Package formatting provides text formatting and colorization utilities.
//
// format.go handles pure text formatting operations without color.
package formatting

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// ---- Basic formatting functions ----

// FormatKeyValuePair formats a key-value pair with optional quotes
func FormatKeyValuePair(key, value string, useQuotes bool) string {
	if useQuotes {
		return fmt.Sprintf("%s='%s'", key, value)
	}
	return fmt.Sprintf("%s=%s", key, value)
}

// FormatExportLine formats an environment variable with export prefix
// Handles various quoting scenarios for shell compatibility
func FormatExportLine(key, value string) string {
	// Input validation
	if key == "" {
		// Return default value for backward compatibility
		return "export=''"
	}

	// Handle value already wrapped in single quotes
	if HasSingleQuotes(value) {
		return fmt.Sprintf("export %s=%s", key, value)
	}

	// Handle embedded single quotes with shell-safe escaping
	if strings.Contains(value, "'") {
		// Use shell escape pattern: 'It'\''s a test'
		parts := strings.Split(value, "'")
		escaped := strings.Join(parts, "'\\''")
		return fmt.Sprintf("export %s='%s'", key, escaped)
	}

	// Normal case
	return fmt.Sprintf("export %s='%s'", key, value)
}

// ---- Quote and string manipulation functions ----

// HasSingleQuotes checks if a string is surrounded by single quotes
func HasSingleQuotes(s string) bool {
	return strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")
}

// UnwrapQuotes removes surrounding single or double quotes from a value
func UnwrapQuotes(value string) string {
	if len(value) < 2 {
		return value
	}

	// Check for matching quotes at beginning and end
	if (value[0] == '"' && value[len(value)-1] == '"') ||
		(value[0] == '\'' && value[len(value)-1] == '\'') {
		return value[1 : len(value)-1]
	}

	return value
}

// ---- Map and collection functions ----

// SortMapKeys returns a sorted slice of all keys in a map
func SortMapKeys(m map[string]string) []string {
	result := SortMapKeysResult(m)
	return result.Unwrap()
}

// ---- Result monad wrapper functions ----

// FormatLineResult formats a key-value pair using Result monad for error handling
func FormatLineResult(key, value string, useExport bool, useQuotes bool) functional.Result[string] {
	if key == "" {
		return functional.Failure[string](fmt.Errorf("invalid environment variable format: key cannot be empty"))
	}

	if useExport {
		return functional.Success(FormatExportLine(key, value))
	} else {
		return functional.Success(FormatKeyValuePair(key, value, useQuotes))
	}
}

// SortMapKeysResult returns a sorted slice of all keys in a map with Result monad
func SortMapKeysResult(m map[string]string) functional.Result[[]string] {
	// Handle nil map
	if m == nil {
		return functional.Success([]string{})
	}

	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return functional.Success(keys)
}

// ---- JSON and complex value formatting ----

// FormatJSONKeyValueResult formats a nested JSON key-value pair with Result monad
func FormatJSONKeyValueResult(parentKey, jsonKey string, jsonValue interface{}, useQuotes bool) functional.Result[string] {
	// Validate keys
	if parentKey == "" || jsonKey == "" {
		return functional.Failure[string](fmt.Errorf("parent key and JSON key must not be empty"))
	}

	valueResult := FormatValueForEnvResult(jsonValue)
	if valueResult.IsFailure() {
		return functional.Failure[string](
			fmt.Errorf("failed to format JSON value for key '%s': %w", jsonKey, valueResult.GetError()))
	}

	// Remove any quotes that might have been added
	valueStr := UnwrapQuotes(valueResult.Unwrap())

	// Create composite key (parent_jsonKey)
	compositeKey := fmt.Sprintf("%s_%s", parentKey, jsonKey)

	// Format with or without quotes
	return functional.Success(FormatKeyValuePair(compositeKey, valueStr, useQuotes))
}

// FormatValueForEnvResult formats a value with Result monad
// Handles nil, basic types, and complex types via JSON marshaling
func FormatValueForEnvResult(value interface{}) functional.Result[string] {
	if value == nil {
		return functional.Success("null")
	}

	switch v := value.(type) {
	case string:
		return functional.Success(v)
	case float64, int, bool:
		return functional.Success(fmt.Sprintf("%v", v))
	default:
		// For complex types, marshal to JSON
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return functional.Failure[string](fmt.Errorf("failed to marshal JSON: %w", err))
		}
		return functional.Success(string(jsonBytes))
	}
}

// FormatValueForEnv formats a value for environment variables
// Convenience wrapper around FormatValueForEnvResult
func FormatValueForEnv(value interface{}) (string, error) {
	result := FormatValueForEnvResult(value)
	if result.IsFailure() {
		return "", result.GetError()
	}
	return result.Unwrap(), nil
}
