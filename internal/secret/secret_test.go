package secret

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestIsURI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid URI",
			input:    "sem://aws/secretsmanager/my-secret",
			expected: true,
		},
		{
			name:     "Invalid URI - http",
			input:    "http://example.com",
			expected: false,
		},
		{
			name:     "Invalid URI - https",
			input:    "https://example.com",
			expected: false,
		},
		{
			name:     "Invalid URI - empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid URI - missing prefix",
			input:    "aws/secretsmanager/my-secret",
			expected: false,
		},
		{
			name:     "Invalid URI - different prefix",
			input:    "secret://aws/secretsmanager/my-secret",
			expected: false,
		},
		{
			name:     "Invalid URI - partial prefix",
			input:    "sem:/aws/secretsmanager/my-secret",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsURI(tt.input)
			if result != tt.expected {
				t.Errorf("IsURI(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildCacheKey(t *testing.T) {
	tests := []struct {
		name       string
		account    string
		service    string
		secretName string
		version    string
		region     string
		expected   string
	}{
		{
			name:       "All fields",
			account:    "123456789012",
			service:    "aws/secretsmanager",
			secretName: "my-secret",
			version:    "1",
			region:     "us-west-2",
			expected:   "123456789012|aws/secretsmanager|my-secret|1|us-west-2",
		},
		{
			name:       "Missing account",
			account:    "",
			service:    "aws/secretsmanager",
			secretName: "my-secret",
			version:    "1",
			region:     "us-west-2",
			expected:   "|aws/secretsmanager|my-secret|1|us-west-2",
		},
		{
			name:       "Missing version",
			account:    "123456789012",
			service:    "aws/secretsmanager",
			secretName: "my-secret",
			version:    "",
			region:     "us-west-2",
			expected:   "123456789012|aws/secretsmanager|my-secret||us-west-2",
		},
		{
			name:       "Missing region",
			account:    "123456789012",
			service:    "aws/secretsmanager",
			secretName: "my-secret",
			version:    "1",
			region:     "",
			expected:   "123456789012|aws/secretsmanager|my-secret|1|",
		},
		{
			name:       "All fields empty",
			account:    "",
			service:    "",
			secretName: "",
			version:    "",
			region:     "",
			expected:   "||||",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildCacheKey(tt.account, tt.service, tt.secretName, tt.version, tt.region)
			if result != tt.expected {
				t.Errorf("BuildCacheKey(%q, %q, %q, %q, %q) = %q, want %q",
					tt.account, tt.service, tt.secretName, tt.version, tt.region, result, tt.expected)
			}
		})
	}
}

func TestFormatSecretValue(t *testing.T) {
	options := DefaultValueOptions()

	tests := []struct {
		name     string
		value    interface{}
		options  ValueOptions
		expected string
		wantErr  bool
	}{
		{
			name:     "String value",
			value:    "test-value",
			options:  options,
			expected: "test-value",
			wantErr:  false,
		},
		{
			name:     "String value with control chars",
			value:    "test\nvalue\r\t",
			options:  options,
			expected: "testvalue", // Control chars are removed by default
			wantErr:  false,
		},
		{
			name:     "String value with control chars - no cleaning",
			value:    "test\nvalue\r\t",
			options:  options.WithCleanControlChars(false),
			expected: "test\nvalue\r\t", // Control chars are preserved
			wantErr:  false,
		},
		{
			name:     "Integer value",
			value:    123,
			options:  options,
			expected: "123",
			wantErr:  false,
		},
		{
			name:     "Float value",
			value:    123.45,
			options:  options,
			expected: "123.45",
			wantErr:  false,
		},
		{
			name:     "Boolean value - true",
			value:    true,
			options:  options,
			expected: "true",
			wantErr:  false,
		},
		{
			name:     "Boolean value - false",
			value:    false,
			options:  options,
			expected: "false",
			wantErr:  false,
		},
		{
			name:     "Map value",
			value:    map[string]interface{}{"key": "value"},
			options:  options,
			expected: `{"key":"value"}`,
			wantErr:  false,
		},
		{
			name:     "Slice value",
			value:    []string{"a", "b", "c"},
			options:  options,
			expected: `["a","b","c"]`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatSecretValue(tt.value, tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatSecretValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("formatSecretValue() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractKeyFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Single segment",
			path:     "key",
			expected: "key",
		},
		{
			name:     "Two segments",
			path:     "parent.child",
			expected: "child",
		},
		{
			name:     "Multiple segments",
			path:     "root.parent.child.grandchild",
			expected: "grandchild",
		},
		{
			name:     "Empty string",
			path:     "",
			expected: "",
		},
		{
			name:     "Path with trailing dot",
			path:     "parent.child.",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractKeyFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("ExtractKeyFromPath(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCombinePaths(t *testing.T) {
	tests := []struct {
		name     string
		segments []string
		expected string
	}{
		{
			name:     "Single segment",
			segments: []string{"key"},
			expected: "key",
		},
		{
			name:     "Two segments",
			segments: []string{"parent", "child"},
			expected: "parent.child",
		},
		{
			name:     "Multiple segments",
			segments: []string{"root", "parent", "child", "grandchild"},
			expected: "root.parent.child.grandchild",
		},
		{
			name:     "Empty segments are filtered",
			segments: []string{"root", "", "child", ""},
			expected: "root.child",
		},
		{
			name:     "All empty segments",
			segments: []string{"", "", ""},
			expected: "",
		},
		{
			name:     "No segments",
			segments: []string{},
			expected: "",
		},
		{
			name:     "Mixed segments",
			segments: []string{"a", "b.c", "d"},
			expected: "a.b.c.d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CombinePaths(tt.segments...)
			if result != tt.expected {
				t.Errorf("CombinePaths(%q) = %q, want %q", tt.segments, result, tt.expected)
			}
		})
	}
}

func TestValueOptionsWithCleanControlChars(t *testing.T) {
	tests := []struct {
		name       string
		options    ValueOptions
		cleanValue bool
		expected   ValueOptions
	}{
		{
			name:       "Change to true",
			options:    ValueOptions{CleanControlChars: false},
			cleanValue: true,
			expected:   ValueOptions{CleanControlChars: true},
		},
		{
			name:       "Change to false",
			options:    ValueOptions{CleanControlChars: true},
			cleanValue: false,
			expected:   ValueOptions{CleanControlChars: false},
		},
		{
			name:       "No change - true",
			options:    ValueOptions{CleanControlChars: true},
			cleanValue: true,
			expected:   ValueOptions{CleanControlChars: true},
		},
		{
			name:       "No change - false",
			options:    ValueOptions{CleanControlChars: false},
			cleanValue: false,
			expected:   ValueOptions{CleanControlChars: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.WithCleanControlChars(tt.cleanValue)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ValueOptions.WithCleanControlChars(%v) = %+v, want %+v", tt.cleanValue, result, tt.expected)
			}
		})
	}
}

func TestParseValueWithOptionsResult(t *testing.T) {
	defaultOptions := DefaultValueOptions()
	noCleanOptions := defaultOptions.WithCleanControlChars(false)

	tests := []struct {
		name         string
		secretString string
		key          string
		options      ValueOptions
		wantValue    string
		wantErr      bool
		expectedErr  error
	}{
		{
			name:         "Empty secret string",
			secretString: "",
			key:          "mykey",
			options:      defaultOptions,
			wantValue:    "",
			wantErr:      true,
			expectedErr:  ErrEmptySecretValue,
		},
		{
			name:         "No key specified - return whole secret",
			secretString: "plain text value",
			key:          "",
			options:      defaultOptions,
			wantValue:    "plain text value",
			wantErr:      false,
		},
		{
			name:         "No key specified with control chars - cleaned",
			secretString: "value\nwith\tcontrol\rchars",
			key:          "",
			options:      defaultOptions,
			wantValue:    "valuewithcontrolchars",
			wantErr:      false,
		},
		{
			name:         "No key specified with control chars - not cleaned",
			secretString: "value\nwith\tcontrol\rchars",
			key:          "",
			options:      noCleanOptions,
			wantValue:    "value\nwith\tcontrol\rchars",
			wantErr:      false,
		},
		{
			name:         "Simple JSON with existing key",
			secretString: `{"mykey": "myvalue"}`,
			key:          "mykey",
			options:      defaultOptions,
			wantValue:    "myvalue",
			wantErr:      false,
		},
		{
			name:         "Simple JSON with nonexistent key",
			secretString: `{"mykey": "myvalue"}`,
			key:          "wrongkey",
			options:      defaultOptions,
			wantValue:    "",
			wantErr:      true,
			expectedErr:  ErrKeyNotFound,
		},
		{
			name:         "Nested JSON with dot notation",
			secretString: `{"parent": {"child": "nestedvalue"}}`,
			key:          "parent.child",
			options:      defaultOptions,
			wantValue:    "nestedvalue",
			wantErr:      false,
		},
		{
			name:         "Deep nested JSON with dot notation",
			secretString: `{"level1": {"level2": {"level3": {"data": "deepvalue"}}}}`,
			key:          "level1.level2.level3.data",
			options:      defaultOptions,
			wantValue:    "deepvalue",
			wantErr:      false,
		},
		{
			name:         "Nested JSON with invalid path",
			secretString: `{"parent": {"child": "nestedvalue"}}`,
			key:          "parent.wrongchild",
			options:      defaultOptions,
			wantValue:    "",
			wantErr:      true,
			expectedErr:  ErrKeyNotFound,
		},
		{
			name:         "Non-JSON string with key specified",
			secretString: "not a json string",
			key:          "mykey",
			options:      defaultOptions,
			wantValue:    "",
			wantErr:      true,
			expectedErr:  ErrInvalidSecretValue,
		},
		{
			name:         "JSON with integer value",
			secretString: `{"mykey": 123}`,
			key:          "mykey",
			options:      defaultOptions,
			wantValue:    "123",
			wantErr:      false,
		},
		{
			name:         "JSON with boolean value",
			secretString: `{"mykey": true}`,
			key:          "mykey",
			options:      defaultOptions,
			wantValue:    "true",
			wantErr:      false,
		},
		{
			name:         "JSON with null value",
			secretString: `{"mykey": null}`,
			key:          "mykey",
			options:      defaultOptions,
			wantValue:    "null",
			wantErr:      false,
		},
		{
			name:         "JSON with array value",
			secretString: `{"mykey": [1, 2, 3]}`,
			key:          "mykey",
			options:      defaultOptions,
			wantValue:    "[1,2,3]",
			wantErr:      false,
		},
		{
			name:         "JSON with nested object value",
			secretString: `{"mykey": {"nested": "value"}}`,
			key:          "mykey",
			options:      defaultOptions,
			wantValue:    `{"nested":"value"}`,
			wantErr:      false,
		},
		{
			name:         "JSON with string containing control chars",
			secretString: `{"mykey": "value\nwith\tcontrol\rchars"}`,
			key:          "mykey",
			options:      defaultOptions,
			wantValue:    "valuewithcontrolchars",
			wantErr:      false,
		},
		{
			name:         "JSON with string containing control chars - not cleaned",
			secretString: `{"mykey": "value\nwith\tcontrol\rchars"}`,
			key:          "mykey",
			options:      noCleanOptions,
			wantValue:    "value\nwith\tcontrol\rchars",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseValueWithOptionsResult(tt.secretString, tt.key, tt.options)

			// Check if we expect an error
			if tt.wantErr {
				if !result.IsFailure() {
					t.Errorf("Expected error but got success with value: %s", result.Unwrap())
					return
				}

				err := result.GetError()
				if tt.expectedErr != nil {
					if !errors.Is(err, tt.expectedErr) && !strings.Contains(err.Error(), tt.expectedErr.Error()) {
						t.Errorf("Expected error containing %q, got %q", tt.expectedErr, err)
					}
				}
				return
			}

			// If we don't expect an error, check the result
			if result.IsFailure() {
				t.Errorf("Unexpected error: %v", result.GetError())
				return
			}

			if got := result.Unwrap(); got != tt.wantValue {
				t.Errorf("ParseValueWithOptionsResult() = %q, want %q", got, tt.wantValue)
			}
		})
	}
}
