// Package env provides utilities for working with environment variables.
package env

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/formatting"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/env"
	"github.com/gumi-tsd/secret-env-manager/internal/text"
)

// EnvVarMap represents a map of environment variables
type EnvVarMap = map[string]string

// EnvVarOptions represents options for formatting environment variables
type EnvVarOptions struct {
	UseQuotes    bool
	SortedKeys   []string
	IncludeURIs  bool
	NoExpandJson bool
}

// NewEnvVarOptions creates default environment variable options
func NewEnvVarOptions() EnvVarOptions {
	return EnvVarOptions{
		UseQuotes:    true,
		SortedKeys:   []string{},
		IncludeURIs:  false,
		NoExpandJson: false,
	}
}

// FormatResult represents the result of formatting operations
type FormatResult struct {
	Content  string
	Warnings []string
}

// FormatEnvVarContent formats environment variables as key=value lines
// Returns formatted content and warnings as a pure function without side effects
func FormatEnvVarContent(values EnvVarMap, options EnvVarOptions) functional.Result[FormatResult] {
	var builder strings.Builder
	filteredKeys := FilterAndSortKeys(values, options.SortedKeys)
	warnings := []string{}

	for i, key := range filteredKeys {
		value, exists := values[key]
		if !exists {
			warnings = append(warnings, fmt.Sprintf("Key '%s' not found in secret values", key))
			continue
		}

		if i > 0 {
			builder.WriteString("\n")
		}

		// Skip writing secret URIs directly unless explicitly requested
		if IsSecretURI(key) && !options.IncludeURIs {
			continue
		}

		// Remove any surrounding quotes
		value = formatting.UnwrapQuotes(value)

		// Check if the value is a JSON object or array and whether to expand it
		if text.IsJSONData(value) {
			if !options.NoExpandJson {
				// When expansion is enabled, parse and process the JSON
				jsonResult := ProcessJSONValue(key, value, options.UseQuotes)
				if jsonResult.IsFailure() {
					return functional.Failure[FormatResult](jsonResult.GetError())
				}

				builder.WriteString(jsonResult.Unwrap())
				continue // Skip the original key-value pair
			} else {
				// When expansion is disabled, compact the JSON
				value = text.CompactJSON(value)
			}
		}

		// Add the formatted key-value pair
		formattedLine := formatting.FormatKeyValuePair(key, value, options.UseQuotes)
		builder.WriteString(formattedLine)
	}

	result := FormatResult{
		Content:  builder.String(),
		Warnings: warnings,
	}

	return functional.Success(result)
}

// ProcessJSONValue processes a JSON value and returns formatted key-value pairs
// Pure function that uses the Result monad for composition
func ProcessJSONValue(parentKey, jsonString string, useQuotes bool) functional.Result[string] {
	// Parse the JSON string into a map
	parseResult := text.ParseJSONResult(jsonString)
	if parseResult.IsFailure() {
		// Not a valid JSON object, treat as regular string
		return functional.Success(formatting.FormatKeyValuePair(parentKey, jsonString, useQuotes))
	}

	// Process the parsed JSON data appropriately
	return processAnyJSONValue(parentKey, parseResult.Unwrap(), useQuotes)
}

// processAnyJSONValue processes JSON data appropriately based on its type
func processAnyJSONValue(parentKey string, jsonData interface{}, useQuotes bool) functional.Result[string] {
	switch data := jsonData.(type) {
	case map[string]interface{}:
		// Process each key-value pair in the object
		return processJSONObject(parentKey, data, useQuotes)

	case []interface{}:
		// Expand array elements with index-based keys
		return processJSONArray(parentKey, data, useQuotes)

	default:
		// Process as a single value if neither an object nor array
		valueStr, err := formatting.FormatValueForEnv(jsonData)
		if err != nil {
			return functional.Failure[string](err)
		}
		return functional.Success(formatting.FormatKeyValuePair(parentKey, valueStr, useQuotes))
	}
}

// processJSONObject expands JSON object keys and values
func processJSONObject(parentKey string, jsonData map[string]interface{}, useQuotes bool) functional.Result[string] {
	// Use functional.Result to accumulate formatted key-value pairs
	formattedPairs := make([]string, 0, len(jsonData))

	// Sort keys alphabetically to maintain consistent processing order
	keys := make([]string, 0, len(jsonData))
	for jsonKey := range jsonData {
		keys = append(keys, jsonKey)
	}
	sort.Strings(keys)

	// Process in sorted key order
	for _, jsonKey := range keys {
		jsonVal := jsonData[jsonKey]
		// Process child elements recursively if they are objects or arrays
		switch childVal := jsonVal.(type) {
		case map[string]interface{}, []interface{}:
			// Generate key for child element
			childKey := fmt.Sprintf("%s_%s", parentKey, jsonKey)
			// Process child element recursively
			childResult := processAnyJSONValue(childKey, childVal, useQuotes)
			if childResult.IsFailure() {
				return childResult
			}
			formattedPairs = append(formattedPairs, childResult.Unwrap())

		default:
			// Format simple values directly
			kvResult := formatting.FormatJSONKeyValueResult(parentKey, jsonKey, jsonVal, useQuotes)
			if kvResult.IsFailure() {
				return kvResult // Propagate the failure
			}
			formattedPairs = append(formattedPairs, kvResult.Unwrap())
		}
	}

	// Join the formatted pairs with newlines
	return functional.Success(strings.Join(formattedPairs, "\n"))
}

// processJSONArray expands JSON array elements with indexed keys
func processJSONArray(parentKey string, jsonArray []interface{}, useQuotes bool) functional.Result[string] {
	formattedPairs := make([]string, 0, len(jsonArray))

	// Process each array element with its index
	for i, item := range jsonArray {
		// Generate key with index
		indexKey := fmt.Sprintf("%s_%d", parentKey, i)

		// Process based on element type
		switch val := item.(type) {
		case map[string]interface{}, []interface{}:
			// Process objects or arrays recursively
			childResult := processAnyJSONValue(indexKey, val, useQuotes)
			if childResult.IsFailure() {
				return childResult
			}
			formattedPairs = append(formattedPairs, childResult.Unwrap())

		default:
			// Expand simple values directly
			valueResult := formatting.FormatValueForEnvResult(item)
			if valueResult.IsFailure() {
				return functional.Failure[string](
					fmt.Errorf("failed to format array item at index %d: %w", i, valueResult.GetError()))
			}
			valueStr := formatting.UnwrapQuotes(valueResult.Unwrap())
			formattedPairs = append(formattedPairs, formatting.FormatKeyValuePair(indexKey, valueStr, useQuotes))
		}
	}

	// Join the formatted pairs with newlines
	return functional.Success(strings.Join(formattedPairs, "\n"))
}

// EnvsToMap converts environment entries to a map
func EnvsToMap(entries []env.Entry) EnvVarMap {
	result := make(EnvVarMap, len(entries))
	for _, entry := range entries {
		result[entry.Key] = entry.Value
	}
	return result
}

// IsSecretURI checks if a key is a secret URI
func IsSecretURI(key string) bool {
	return strings.HasPrefix(key, "sem://")
}

// FormatEnvVarsResult is a monadic version of FormatEnvVars
func FormatEnvVarsResult(variables map[string]string, useExport bool, useQuotes bool) functional.Result[FormatResult] {
	// Sort keys for consistent output
	keys := formatting.SortMapKeys(variables)
	lines := make([]string, 0, len(variables))

	for _, key := range keys {
		value := variables[key]
		lineResult := formatting.FormatLineResult(key, value, useExport, useQuotes)

		if lineResult.IsFailure() {
			return functional.Failure[FormatResult](lineResult.GetError())
		}

		lines = append(lines, lineResult.Unwrap())
	}

	// Return newline-separated string regardless of export option
	// This improves readability for both export and non-export formats
	content := strings.Join(lines, "\n")

	return functional.Success(FormatResult{
		Content:  content,
		Warnings: []string{},
	})
}

// FormatEnvVars is a convenience wrapper that unwraps the Result monad
// This maintains backward compatibility with existing code
func FormatEnvVars(variables map[string]string, useExport bool, useQuotes bool) (string, []string) {
	// Use the monadic version and unwrap the result
	result := FormatEnvVarsResult(variables, useExport, useQuotes)

	// Handle potential errors (should rarely happen since we validate inputs)
	if result.IsFailure() {
		// Return empty results in case of error
		return "", []string{}
	}

	formatResult := result.Unwrap()

	// Extract the lines from the content
	// Always split by newline since we're using newlines for both formats
	lines := strings.Split(formatResult.Content, "\n")

	return formatResult.Content, lines
}

// AddPrefix adds a prefix to all keys in a map
func AddPrefix(m map[string]string, prefix string) map[string]string {
	result := make(map[string]string, len(m))
	for key, value := range m {
		result[prefix+key] = value
	}
	return result
}

// TransformMapValues applies a transformer function to all values in a map
func TransformMapValues(m map[string]string, transformer text.StringTransformer) map[string]string {
	result := make(map[string]string, len(m))
	for key, value := range m {
		result[key] = transformer(value)
	}
	return result
}

// FilterMap filters a map by key predicates
func FilterMap(m map[string]string, keyPredicate text.StringPredicate) map[string]string {
	result := make(map[string]string)
	for key, value := range m {
		if keyPredicate(key) {
			result[key] = value
		}
	}
	return result
}
