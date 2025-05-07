package formatting

import (
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

func TestFormatKeyValues(t *testing.T) {
	tests := []struct {
		name      string
		keys      []string
		values    map[string]string
		useQuotes bool
		expected  []string
	}{
		{
			name:      "Empty map",
			keys:      []string{},
			values:    map[string]string{},
			useQuotes: false,
			expected:  []string{},
		},
		{
			name: "Single entry",
			keys: []string{"USER"},
			values: map[string]string{
				"USER": "admin",
			},
			useQuotes: false,
			expected:  []string{"USER=admin"},
		},
		{
			name: "Multiple entries",
			keys: []string{"DEBUG", "PASSWORD", "USER"},
			values: map[string]string{
				"USER":     "admin",
				"PASSWORD": "secret",
				"DEBUG":    "true",
			},
			useQuotes: false,
			expected: []string{
				"DEBUG=true",
				"PASSWORD=secret",
				"USER=admin",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatKeyValues(tt.keys, tt.values, tt.useQuotes)

			// Sort both slices to ensure consistent comparison
			if !compareStringSlices(result, tt.expected) {
				t.Errorf("FormatKeyValues(%v, %v, %v) = %v, want %v",
					tt.keys, tt.values, tt.useQuotes, result, tt.expected)
			}
		})
	}
}

func TestIndent(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		prefix   string
		expected string
	}{
		{
			name:     "Empty string",
			text:     "",
			prefix:   "  ",
			expected: "",
		},
		{
			name:     "No indent",
			text:     "Hello",
			prefix:   "",
			expected: "Hello",
		},
		{
			name:     "Simple indent",
			text:     "Hello",
			prefix:   "  ",
			expected: "  Hello",
		},
		{
			name:     "Multiline text",
			text:     "Line 1\nLine 2\nLine 3",
			prefix:   "    ",
			expected: "    Line 1\n    Line 2\n    Line 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Indent(tt.text, tt.prefix)
			if result != tt.expected {
				t.Errorf("Indent(%q, %q) = %q, want %q",
					tt.text, tt.prefix, result, tt.expected)
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

func TestEscapeSingleQuotes(t *testing.T) {
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
			name:     "No single quotes",
			text:     "Hello world",
			expected: "Hello world",
		},
		{
			name:     "With single quote",
			text:     "It's a test",
			expected: "It\\'s a test",
		},
		{
			name:     "Multiple single quotes",
			text:     "Don't say 'hello' like that",
			expected: "Don\\'t say \\'hello\\' like that",
		},
		{
			name:     "With double quotes (unchanged)",
			text:     "\"Hello world\"",
			expected: "\"Hello world\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeSingleQuotes(tt.text)
			if result != tt.expected {
				t.Errorf("EscapeSingleQuotes(%q) = %q, want %q",
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

func TestFormatPlainLine(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected string
	}{
		{
			name:     "Basic line",
			key:      "USER",
			value:    "admin",
			expected: "USER=admin",
		},
		{
			name:     "Empty value",
			key:      "DEBUG",
			value:    "",
			expected: "DEBUG=",
		},
		{
			name:     "With spaces in value",
			key:      "DESCRIPTION",
			value:    "This is a test",
			expected: "DESCRIPTION=This is a test",
		},
		{
			name:     "With special chars in value",
			key:      "PATH",
			value:    "/usr/bin:/bin",
			expected: "PATH=/usr/bin:/bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPlainLine(tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("FormatPlainLine(%q, %q) = %q, want %q",
					tt.key, tt.value, result, tt.expected)
			}
		})
	}
}

// Helper function to compare two string slices ignoring order
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Convert slices to maps for easier comparison
	mapA := make(map[string]bool)
	mapB := make(map[string]bool)

	for _, val := range a {
		mapA[val] = true
	}

	for _, val := range b {
		mapB[val] = true
	}

	// Compare maps
	for key := range mapA {
		if !mapB[key] {
			return false
		}
	}

	return true
}
