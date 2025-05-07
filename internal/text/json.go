// Package json provides utilities for handling JSON data.
package text

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// Error definitions for JSON processing
var (
	ErrInvalidJSON       = fmt.Errorf("invalid JSON format")
	ErrJSONPathNotFound  = fmt.Errorf("JSON path not found")
	ErrInvalidJSONPath   = fmt.Errorf("invalid JSON path syntax")
	ErrJSONValueNotArray = fmt.Errorf("JSON value is not an array")
	ErrJSONValueNotMap   = fmt.Errorf("JSON value is not a map")
	ErrJSONMarshal       = fmt.Errorf("failed to convert value to JSON string")
)

// IsJSON performs a strict check if a string is valid JSON by attempting to parse it.
// Unlike IsJSONObject, this function validates any JSON type (objects, arrays, primitives).
func IsJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// IsJSONObject checks if a string appears to be a JSON object (starts with { and ends with })
// Note: This is a fast check without full JSON parsing
func IsJSONObject(s string) bool {
	// First do a quick check with strings.TrimSpace for efficiency
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 2 || trimmed[0] != '{' || trimmed[len(trimmed)-1] != '}' {
		return false
	}

	// For better validation, we can use functional composition with predicates
	hasObjectBrackets := ComposePredicate(
		HasPrefix("{"),
		HasSuffix("}"),
	)

	return hasObjectBrackets(trimmed)
}

// IsJSONObjectResult is a monadic version of IsJSONObject that returns a Result
func IsJSONObjectResult(s string) functional.Result[bool] {
	return functional.Success(IsJSONObject(s))
}

// IsJSONArray checks if a string appears to be a JSON array (starts with [ and ends with ])
// Note: This is a fast check without full JSON parsing
func IsJSONArray(s string) bool {
	// Use a similar approach to IsJSONObject
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 2 || trimmed[0] != '[' || trimmed[len(trimmed)-1] != ']' {
		return false
	}

	hasArrayBrackets := ComposePredicate(
		HasPrefix("["),
		HasSuffix("]"),
	)

	return hasArrayBrackets(trimmed)
}

// IsJSONData checks if a string appears to be either a JSON object or array
// Note: This is a fast check without full JSON parsing
func IsJSONData(s string) bool {
	return IsJSONObject(s) || IsJSONArray(s)
}

// SafeJSONEncode safely encodes a value to JSON.
// Returns empty string if encoding fails.
func SafeJSONEncode(value interface{}) string {
	result, err := MarshalToJSONString(value)
	if err != nil {
		return ""
	}
	return result
}

// ParseJSON parses a JSON string into an interface{} value
func ParseJSON(jsonString string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(jsonString), &result)
	return result, err
}

// ParseJSONOption parses a JSON string and returns an Option
func ParseJSONOption(jsonString string) functional.Option[interface{}] {
	var result interface{}

	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		return functional.None[interface{}]()
	}

	return functional.Some(result)
}

// ParseJSONResult parses a JSON string and returns a Result
func ParseJSONResult(jsonString string) functional.Result[interface{}] {
	var result interface{}

	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		return functional.Failure[interface{}](
			fmt.Errorf("%w: %v", ErrInvalidJSON, err))
	}

	return functional.Success(result)
}

// ParseJSONMapResult parses a JSON string specifically as a map of string key-value pairs
func ParseJSONMapResult(jsonString string) functional.Result[map[string]interface{}] {
	jsonResult := ParseJSONResult(jsonString)

	if jsonResult.IsFailure() {
		return functional.Failure[map[string]interface{}](jsonResult.GetError())
	}

	// Try to convert to map
	jsonData := jsonResult.Unwrap()

	if mapData, ok := jsonData.(map[string]interface{}); ok {
		return functional.Success(mapData)
	}

	return functional.Failure[map[string]interface{}](
		fmt.Errorf("JSON is not an object: %v", jsonData))
}

// CreateStringMap transforms a map of interface{} values to string values
func CreateStringMap(data map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = v
		default:
			// For non-string values, convert to JSON
			jsonBytes, err := json.Marshal(v)
			if err == nil {
				result[key] = string(jsonBytes)
			}
		}
	}

	return result
}

// NavigatePath follows a path of keys/indices through a nested data structure
func NavigatePath(data interface{}, path []string) (interface{}, error) {
	currentData := data

	for i, segment := range path {
		switch current := currentData.(type) {
		case map[string]interface{}:
			// Access object property
			value, exists := current[segment]
			if !exists {
				return nil, fmt.Errorf("key '%s' not found at path position %d", segment, i)
			}
			currentData = value

		case []interface{}:
			// Parse array index
			index, err := parseIndex(segment, len(current))
			if err != nil {
				return nil, fmt.Errorf("at path position %d: %w", i, err)
			}

			if index < 0 || index >= len(current) {
				return nil, fmt.Errorf("array index %d out of bounds at path position %d", index, i)
			}

			currentData = current[index]

		default:
			return nil, fmt.Errorf("cannot navigate further at path position %d: not an object or array", i)
		}
	}

	return currentData, nil
}

// NavigatePathResult follows a path of keys/indices through a nested data structure (monadic version)
func NavigatePathResult(data interface{}, path []string) functional.Result[interface{}] {
	result, err := NavigatePath(data, path)
	if err != nil {
		return functional.Failure[interface{}](err)
	}
	return functional.Success(result)
}

// parseIndex converts a string to an array index
func parseIndex(index string, arrayLength int) (int, error) {
	// Handle negative indices (count from end of array)
	if strings.HasPrefix(index, "-") {
		if index == "-" {
			// Special case: "-" refers to the last element
			return arrayLength - 1, nil
		}

		// Parse numeric part
		val, err := strconv.Atoi(index)
		if err != nil {
			return 0, fmt.Errorf("invalid array index: %s", index)
		}

		// Convert negative index to positive
		if -val <= arrayLength {
			return arrayLength + val, nil
		}

		return 0, fmt.Errorf("negative index %d out of bounds", val)
	}

	// Parse regular index
	val, err := strconv.Atoi(index)
	if err != nil {
		return 0, fmt.Errorf("invalid array index: %s", index)
	}

	// Check if index is out of bounds
	if val >= arrayLength {
		return 0, fmt.Errorf("index %d out of bounds", val)
	}

	return val, nil
}

// NavigateArrayIndex accesses a specific index in an array
func NavigateArrayIndex(data interface{}, index int) functional.Result[interface{}] {
	array, ok := data.([]interface{})
	if !ok {
		return functional.Failure[interface{}](ErrJSONValueNotArray)
	}

	if index < 0 || index >= len(array) {
		return functional.Failure[interface{}](
			fmt.Errorf("index %d out of bounds for array of length %d", index, len(array)))
	}

	return functional.Success(array[index])
}

// NavigateObjectKey accesses a specific key in an object
func NavigateObjectKey(data interface{}, key string) functional.Result[interface{}] {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return functional.Failure[interface{}](ErrJSONValueNotMap)
	}

	value, exists := obj[key]
	if !exists {
		return functional.Failure[interface{}](
			fmt.Errorf("key '%s' not found in object", key))
	}

	return functional.Success(value)
}

// FormatAvailableKeys generates a string listing available keys in a map
func FormatAvailableKeys(data interface{}) string {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return "Not a JSON object - no keys available"
	}

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}

	return "Available keys: " + strings.Join(keys, ", ")
}

// MarshalToJSONString converts any value to a JSON string representation
// It follows the pure function principle and explicitly returns errors
func MarshalToJSONString(value interface{}) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(bytes), nil
}

// MarshalToJSONResult converts any value to a JSON string representation with Result monad
// This function follows functional programming principles using the Result monad pattern
func MarshalToJSONResult(value interface{}) functional.Result[string] {
	jsonStr, err := MarshalToJSONString(value)
	if err != nil {
		return functional.Failure[string](err)
	}
	return functional.Success(jsonStr)
}

// CompactJSON parses and re-serializes JSON to remove unnecessary whitespace
// while preserving the data structure. If the input is not valid JSON,
// it returns the original string unchanged.
func CompactJSON(jsonString string) string {
	// Parse the JSON
	var data interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		// Not valid JSON, return original
		return jsonString
	}

	// Re-serialize to compact form
	compacted, err := json.Marshal(data)
	if err != nil {
		// Unlikely to happen, but return original in case of error
		return jsonString
	}

	return string(compacted)
}

// CompactJSONResult converts the function to use Result monad pattern
func CompactJSONResult(jsonString string) functional.Result[string] {
	return functional.Success(CompactJSON(jsonString))
}
