// Package formatting provides text formatting and colorization utilities.
package formatting

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// FormatKeyValuePair formats a key-value pair with or without quotes
// Pure function that handles environment variable formatting
func FormatKeyValuePair(key, value string, useQuotes bool) string {
	if useQuotes {
		return fmt.Sprintf("%s='%s'", key, value)
	}
	return fmt.Sprintf("%s=%s", key, value)
}

// FormatKeyValues formats multiple key-value pairs into lines
func FormatKeyValues(keys []string, values map[string]string, useQuotes bool) []string {
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		if value, ok := values[key]; ok {
			lines = append(lines, FormatKeyValuePair(key, value, useQuotes))
		}
	}
	return lines
}

// Indent adds a prefix to each line in a multi-line string
func Indent(text, prefix string) string {
	if text == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

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

// EscapeSingleQuotes escapes single quotes in a string for shell usage
func EscapeSingleQuotes(value string) string {
	return strings.ReplaceAll(value, "'", "\\'")
}

// FormatExportLine formats an environment variable with export prefix
func FormatExportLine(key, value string) string {
	// Handle case of value already wrapped in single quotes
	if HasSingleQuotes(value) {
		return fmt.Sprintf("export %s=%s", key, value)
	}

	// Special handling for embedded single quotes
	if strings.Contains(value, "'") {
		// Use the shell escape pattern: 'It'\''s a test' instead of escaping with backslash
		parts := strings.Split(value, "'")
		escaped := strings.Join(parts, "'\\''")
		return fmt.Sprintf("export %s='%s'", key, escaped)
	}

	// Normal case
	return fmt.Sprintf("export %s='%s'", key, value)
}

// FormatPlainLine formats a key-value pair without export
func FormatPlainLine(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

// FormatLineResult formats a key-value pair with Result monad
func FormatLineResult(key, value string, useExport bool, useQuotes bool) functional.Result[string] {
	if key == "" {
		return functional.Failure[string](fmt.Errorf("invalid environment variable format"))
	}

	if useExport {
		return functional.Success(FormatExportLine(key, value))
	} else {
		return functional.Success(FormatKeyValuePair(key, value, useQuotes))
	}
}

// SortMapKeysResult returns a sorted slice of all keys in a map with Result monad
func SortMapKeysResult(m map[string]string) functional.Result[[]string] {
	// Guard against nil map
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

// SortMapKeys returns a sorted slice of all keys in a map
func SortMapKeys(m map[string]string) []string {
	result := SortMapKeysResult(m)
	// This should never fail as the function is pure and we guard against nil values
	if result.IsFailure() {
		return []string{}
	}
	return result.Unwrap()
}

// FormatJSONKeyValueResult formats a nested JSON key-value pair with Result monad
func FormatJSONKeyValueResult(parentKey, jsonKey string, jsonValue interface{}, useQuotes bool) functional.Result[string] {
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

// FormatJSONKeyValue formats a nested JSON key-value pair
func FormatJSONKeyValue(parentKey, jsonKey string, jsonValue interface{}, useQuotes bool) (string, error) {
	// Use the Result monad version and handle errors explicitly
	result := FormatJSONKeyValueResult(parentKey, jsonKey, jsonValue, useQuotes)
	if result.IsFailure() {
		return "", result.GetError()
	}
	return result.Unwrap(), nil
}

// FormatValueForEnvResult formats a value with Result monad
func FormatValueForEnvResult(value interface{}) functional.Result[string] {
	switch v := value.(type) {
	case string:
		return functional.Success(v)
	case float64, int, bool:
		return functional.Success(fmt.Sprintf("%v", v))
	default:
		// If not a simple type, try to marshal to JSON
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return functional.Failure[string](fmt.Errorf("failed to marshal JSON: %w", err))
		}
		return functional.Success(string(jsonBytes))
	}
}

// FormatValueForEnv formats a value based on its type for environment variables
func FormatValueForEnv(value interface{}) (string, error) {
	result := FormatValueForEnvResult(value)
	if result.IsFailure() {
		return "", result.GetError()
	}
	return result.Unwrap(), nil
}
