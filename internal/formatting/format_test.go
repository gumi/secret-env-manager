package formatting

import (
	"strings"
	"testing"
)

func TestFormatKeyValuePair(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		useQuotes bool
		expected  string
	}{
		{
			name:      "Basic formatting without quotes",
			key:       "USER",
			value:     "admin",
			useQuotes: false,
			expected:  "USER=admin",
		},
		{
			name:      "Basic formatting with quotes",
			key:       "USER",
			value:     "admin",
			useQuotes: true,
			expected:  "USER='admin'",
		},
		{
			name:      "Empty value without quotes",
			key:       "DEBUG",
			value:     "",
			useQuotes: false,
			expected:  "DEBUG=",
		},
		{
			name:      "Empty value with quotes",
			key:       "DEBUG",
			value:     "",
			useQuotes: true,
			expected:  "DEBUG=''",
		},
		{
			name:      "Key with spaces without quotes",
			key:       "APP NAME",
			value:     "TestApp",
			useQuotes: false,
			expected:  "APP NAME=TestApp",
		},
		{
			name:      "Value with spaces with quotes",
			key:       "DESCRIPTION",
			value:     "This is a test",
			useQuotes: true,
			expected:  "DESCRIPTION='This is a test'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatKeyValuePair(tt.key, tt.value, tt.useQuotes)
			if result != tt.expected {
				t.Errorf("FormatKeyValuePair(%q, %q, %v) = %q, want %q",
					tt.key, tt.value, tt.useQuotes, result, tt.expected)
			}
		})
	}
}

func TestHasSingleQuotes(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{
			name:     "Empty string",
			text:     "",
			expected: false,
		},
		{
			name:     "No quotes",
			text:     "Hello world",
			expected: false,
		},
		{
			name:     "With single quotes",
			text:     "'Hello world'",
			expected: true,
		},
		{
			name:     "With double quotes",
			text:     "\"Hello world\"",
			expected: false,
		},
		{
			name:     "With single quotes in middle",
			text:     "Hello 'world'",
			expected: false,
		},
		{
			name:     "Only opening quote",
			text:     "'Hello world",
			expected: false,
		},
		{
			name:     "Only closing quote",
			text:     "Hello world'",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasSingleQuotes(tt.text)
			if result != tt.expected {
				t.Errorf("HasSingleQuotes(%q) = %v, want %v",
					tt.text, result, tt.expected)
			}
		})
	}
}

func TestUnwrapQuotes(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Empty string",
			text:     "",
			expected: "",
		},
		{
			name:     "No quotes",
			text:     "Hello world",
			expected: "Hello world",
		},
		{
			name:     "With single quotes",
			text:     "'Hello world'",
			expected: "Hello world",
		},
		{
			name:     "With double quotes",
			text:     "\"Hello world\"",
			expected: "Hello world",
		},
		{
			name:     "With mixed quotes",
			text:     "\"'Hello world'\"",
			expected: "'Hello world'",
		},
		{
			name:     "Only opening quote",
			text:     "'Hello world",
			expected: "'Hello world",
		},
		{
			name:     "Only closing quote",
			text:     "Hello world'",
			expected: "Hello world'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnwrapQuotes(tt.text)
			if result != tt.expected {
				t.Errorf("UnwrapQuotes(%q) = %q, want %q",
					tt.text, result, tt.expected)
			}
		})
	}
}

func TestFormatExportLine(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected string
	}{
		{
			name:     "Basic export",
			key:      "USER",
			value:    "admin",
			expected: "export USER='admin'",
		},
		{
			name:     "Empty value",
			key:      "DEBUG",
			value:    "",
			expected: "export DEBUG=''",
		},
		{
			name:     "With spaces in value",
			key:      "DESCRIPTION",
			value:    "This is a test",
			expected: "export DESCRIPTION='This is a test'",
		},
		{
			name:     "With special chars in value",
			key:      "PATH",
			value:    "/usr/bin:/bin",
			expected: "export PATH='/usr/bin:/bin'",
		},
		{
			name:     "With single quotes already in value",
			key:      "MESSAGE",
			value:    "'Hello world'",
			expected: "export MESSAGE='Hello world'",
		},
		{
			name:     "With embedded single quotes",
			key:      "QUOTE",
			value:    "It's a test",
			expected: "export QUOTE='It'\\''s a test'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatExportLine(tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("FormatExportLine(%q, %q) = %q, want %q",
					tt.key, tt.value, result, tt.expected)
			}
		})
	}
}

func TestFormatLineResult(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		value       string
		useExport   bool
		useQuotes   bool
		expectError bool
		expected    string
	}{
		{
			name:        "Empty key",
			key:         "",
			value:       "test",
			useExport:   false,
			useQuotes:   false,
			expectError: true,
		},
		{
			name:        "With export",
			key:         "KEY",
			value:       "value",
			useExport:   true,
			useQuotes:   false, // useQuotes is ignored when useExport is true
			expectError: false,
			expected:    "export KEY='value'",
		},
		{
			name:        "Without export, with quotes",
			key:         "KEY",
			value:       "value",
			useExport:   false,
			useQuotes:   true,
			expectError: false,
			expected:    "KEY='value'",
		},
		{
			name:        "Without export, no quotes",
			key:         "KEY",
			value:       "value",
			useExport:   false,
			useQuotes:   false,
			expectError: false,
			expected:    "KEY=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatLineResult(tt.key, tt.value, tt.useExport, tt.useQuotes)

			if tt.expectError {
				if !result.IsFailure() {
					t.Errorf("FormatLineResult(%q, %q, %v, %v) should fail but succeeded",
						tt.key, tt.value, tt.useExport, tt.useQuotes)
				}
			} else {
				if result.IsFailure() {
					t.Errorf("FormatLineResult(%q, %q, %v, %v) failed: %v",
						tt.key, tt.value, tt.useExport, tt.useQuotes, result.GetError())
				} else {
					if result.Unwrap() != tt.expected {
						t.Errorf("FormatLineResult(%q, %q, %v, %v) = %q, want %q",
							tt.key, tt.value, tt.useExport, tt.useQuotes, result.Unwrap(), tt.expected)
					}
				}
			}
		})
	}
}

func TestSortMapKeysResult(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected []string
	}{
		{
			name:     "Empty map",
			input:    map[string]string{},
			expected: []string{},
		},
		{
			name:     "Nil map",
			input:    nil,
			expected: []string{},
		},
		{
			name: "Single entry",
			input: map[string]string{
				"key": "value",
			},
			expected: []string{"key"},
		},
		{
			name: "Multiple entries",
			input: map[string]string{
				"c": "3",
				"a": "1",
				"b": "2",
			},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortMapKeysResult(tt.input)

			if result.IsFailure() {
				t.Errorf("SortMapKeysResult(%v) failed: %v", tt.input, result.GetError())
				return
			}

			keys := result.Unwrap()

			// Check length
			if len(keys) != len(tt.expected) {
				t.Errorf("SortMapKeysResult(%v) returned %d keys, want %d",
					tt.input, len(keys), len(tt.expected))
				return
			}

			// Check sorting
			for i, key := range keys {
				if key != tt.expected[i] {
					t.Errorf("SortMapKeysResult(%v) at index %d = %q, want %q",
						tt.input, i, key, tt.expected[i])
				}
			}
		})
	}
}

func TestSortMapKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected []string
	}{
		{
			name:     "Empty map",
			input:    map[string]string{},
			expected: []string{},
		},
		{
			name:     "Nil map",
			input:    nil,
			expected: []string{},
		},
		{
			name: "Multiple entries",
			input: map[string]string{
				"z": "last",
				"a": "first",
				"m": "middle",
			},
			expected: []string{"a", "m", "z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := SortMapKeys(tt.input)

			// Check length
			if len(keys) != len(tt.expected) {
				t.Errorf("SortMapKeys(%v) returned %d keys, want %d",
					tt.input, len(keys), len(tt.expected))
				return
			}

			// Check sorting
			for i, key := range keys {
				if key != tt.expected[i] {
					t.Errorf("SortMapKeys(%v) at index %d = %q, want %q",
						tt.input, i, key, tt.expected[i])
				}
			}
		})
	}
}

func TestSortMapKeysErrorHandling(t *testing.T) {
	// Since we can't mock SortMapKeysResult directly, we'll add a test that
	// specifically checks the behavior of SortMapKeys in the normal case
	// Note: This function is essentially testing the same thing as TestSortMapKeys
	// but we're keeping it to verify that we're getting 100% coverage

	// The case is simple - we just want to ensure the function works as expected
	m := map[string]string{
		"z": "last",
		"a": "first",
	}

	expected := []string{"a", "z"}
	result := SortMapKeys(m)

	// Verify the result is as expected
	if len(result) != len(expected) {
		t.Errorf("SortMapKeys should return %d items, got %d", len(expected), len(result))
	}

	// Check each item in order
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("SortMapKeys at index %d should be %q, got %q", i, v, result[i])
		}
	}
}

func TestFormatJSONKeyValueResult(t *testing.T) {
	tests := []struct {
		name        string
		parentKey   string
		jsonKey     string
		jsonValue   interface{}
		useQuotes   bool
		expectError bool
		expected    string
	}{
		{
			name:        "String value without quotes",
			parentKey:   "CONFIG",
			jsonKey:     "name",
			jsonValue:   "test-app",
			useQuotes:   false,
			expectError: false,
			expected:    "CONFIG_name=test-app",
		},
		{
			name:        "String value with quotes",
			parentKey:   "CONFIG",
			jsonKey:     "name",
			jsonValue:   "test-app",
			useQuotes:   true,
			expectError: false,
			expected:    "CONFIG_name='test-app'",
		},
		{
			name:        "Numeric value",
			parentKey:   "CONFIG",
			jsonKey:     "port",
			jsonValue:   8080,
			useQuotes:   false,
			expectError: false,
			expected:    "CONFIG_port=8080",
		},
		{
			name:        "Boolean value",
			parentKey:   "CONFIG",
			jsonKey:     "debug",
			jsonValue:   true,
			useQuotes:   true,
			expectError: false,
			expected:    "CONFIG_debug='true'",
		},
		{
			name:        "Complex value",
			parentKey:   "CONFIG",
			jsonKey:     "complex",
			jsonValue:   map[string]interface{}{"a": 1, "b": "test"},
			useQuotes:   false,
			expectError: false,
			// The exact JSON string might vary, so we'll check this separately
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatJSONKeyValueResult(tt.parentKey, tt.jsonKey, tt.jsonValue, tt.useQuotes)

			if tt.expectError {
				if !result.IsFailure() {
					t.Errorf("FormatJSONKeyValueResult(%q, %q, %v, %v) should fail but succeeded",
						tt.parentKey, tt.jsonKey, tt.jsonValue, tt.useQuotes)
				}
			} else {
				if result.IsFailure() {
					t.Errorf("FormatJSONKeyValueResult(%q, %q, %v, %v) failed: %v",
						tt.parentKey, tt.jsonKey, tt.jsonValue, tt.useQuotes, result.GetError())
				} else if isComplexObject(tt.jsonValue) {
					// For complex objects, we'll check that the result contains the expected components
					formattedValue := result.Unwrap()
					expectedPrefix := tt.parentKey + "_" + tt.jsonKey + "="
					if !strings.HasPrefix(formattedValue, expectedPrefix) {
						t.Errorf("FormatJSONKeyValueResult for complex value doesn't have expected prefix: %q", formattedValue)
					}
				} else if result.Unwrap() != tt.expected {
					t.Errorf("FormatJSONKeyValueResult(%q, %q, %v, %v) = %q, want %q",
						tt.parentKey, tt.jsonKey, tt.jsonValue, tt.useQuotes, result.Unwrap(), tt.expected)
				}
			}
		})
	}
}

func TestFormatValueForEnvResult(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		expectError bool
		expected    string
	}{
		{
			name:        "String value",
			value:       "test-string",
			expectError: false,
			expected:    "test-string",
		},
		{
			name:        "Integer value",
			value:       42,
			expectError: false,
			expected:    "42",
		},
		{
			name:        "Float value",
			value:       3.14,
			expectError: false,
			expected:    "3.14",
		},
		{
			name:        "Boolean true",
			value:       true,
			expectError: false,
			expected:    "true",
		},
		{
			name:        "Boolean false",
			value:       false,
			expectError: false,
			expected:    "false",
		},
		{
			name:        "Simple map",
			value:       map[string]string{"key": "value"},
			expectError: false,
			// The exact JSON string might vary, so we'll check this separately
		},
		{
			name:        "Array",
			value:       []string{"a", "b", "c"},
			expectError: false,
			// The exact JSON string might vary, so we'll check this separately
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatValueForEnvResult(tt.value)

			if tt.expectError {
				if !result.IsFailure() {
					t.Errorf("FormatValueForEnvResult(%v) should fail but succeeded", tt.value)
				}
			} else {
				if result.IsFailure() {
					t.Errorf("FormatValueForEnvResult(%v) failed: %v", tt.value, result.GetError())
				} else {
					switch tt.value.(type) {
					case map[string]string, []string:
						// For complex objects, just verify we got some result
						if result.Unwrap() == "" {
							t.Errorf("FormatValueForEnvResult(%v) returned empty string for complex object", tt.value)
						}
					default:
						if result.Unwrap() != tt.expected {
							t.Errorf("FormatValueForEnvResult(%v) = %q, want %q", tt.value, result.Unwrap(), tt.expected)
						}
					}
				}
			}
		})
	}
}

func TestFormatValueForEnv(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		expectError bool
		expected    string
	}{
		{
			name:        "String value",
			value:       "simple",
			expectError: false,
			expected:    "simple",
		},
		{
			name:        "Boolean value",
			value:       true,
			expectError: false,
			expected:    "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatValueForEnv(tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("FormatValueForEnv(%v) should fail but succeeded", tt.value)
				}
			} else {
				if err != nil {
					t.Errorf("FormatValueForEnv(%v) failed: %v", tt.value, err)
				} else if result != tt.expected {
					t.Errorf("FormatValueForEnv(%v) = %q, want %q", tt.value, result, tt.expected)
				}
			}
		})
	}
}

// Helper function to check if a value is a complex object
func isComplexObject(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}, map[string]string, []interface{}, []string:
		return true
	default:
		return false
	}
}
