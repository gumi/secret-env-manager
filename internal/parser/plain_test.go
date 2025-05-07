package parser

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/env"
)

func TestNewLine(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		number   int
		expected Line
	}{
		{
			name:    "Empty line",
			content: "",
			number:  1,
			expected: Line{
				Content: "",
				Number:  1,
				Type:    EmptyLine,
				Trimmed: "",
			},
		},
		{
			name:    "Comment line",
			content: "# This is a comment",
			number:  2,
			expected: Line{
				Content: "# This is a comment",
				Number:  2,
				Type:    CommentLine,
				Trimmed: "# This is a comment",
			},
		},
		{
			name:    "Secret URI line",
			content: "sem://aws/secretsmanager/my-secret",
			number:  3,
			expected: Line{
				Content: "sem://aws/secretsmanager/my-secret",
				Number:  3,
				Type:    SecretURILine,
				Trimmed: "sem://aws/secretsmanager/my-secret",
			},
		},
		{
			name:    "Key-value line",
			content: "API_KEY=12345",
			number:  4,
			expected: Line{
				Content: "API_KEY=12345",
				Number:  4,
				Type:    KeyValueLine,
				Trimmed: "API_KEY=12345",
			},
		},
		{
			name:    "Key only line",
			content: "DATABASE_URL",
			number:  5,
			expected: Line{
				Content: "DATABASE_URL",
				Number:  5,
				Type:    KeyOnlyLine,
				Trimmed: "DATABASE_URL",
			},
		},
		{
			name:    "Whitespace trimming",
			content: "  API_KEY=12345  ",
			number:  6,
			expected: Line{
				Content: "  API_KEY=12345  ",
				Number:  6,
				Type:    KeyValueLine,
				Trimmed: "API_KEY=12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewLine(tt.content, tt.number)

			if result.Content != tt.expected.Content {
				t.Errorf("Content = %q, want %q", result.Content, tt.expected.Content)
			}
			if result.Number != tt.expected.Number {
				t.Errorf("Number = %d, want %d", result.Number, tt.expected.Number)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Type = %d, want %d", result.Type, tt.expected.Type)
			}
			if result.Trimmed != tt.expected.Trimmed {
				t.Errorf("Trimmed = %q, want %q", result.Trimmed, tt.expected.Trimmed)
			}
		})
	}
}

func TestLine_IsEmpty(t *testing.T) {
	testCases := []struct {
		name     string
		lineType LineType
		expected bool
	}{
		{"Empty line", EmptyLine, true},
		{"Comment line", CommentLine, false},
		{"Secret URI line", SecretURILine, false},
		{"Key-value line", KeyValueLine, false},
		{"Key only line", KeyOnlyLine, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			line := Line{Type: tc.lineType}
			if result := line.IsEmpty(); result != tc.expected {
				t.Errorf("IsEmpty() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestLine_IsComment(t *testing.T) {
	testCases := []struct {
		name     string
		lineType LineType
		expected bool
	}{
		{"Empty line", EmptyLine, false},
		{"Comment line", CommentLine, true},
		{"Secret URI line", SecretURILine, false},
		{"Key-value line", KeyValueLine, false},
		{"Key only line", KeyOnlyLine, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			line := Line{Type: tc.lineType}
			if result := line.IsComment(); result != tc.expected {
				t.Errorf("IsComment() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestLine_IsSecret(t *testing.T) {
	testCases := []struct {
		name     string
		lineType LineType
		expected bool
	}{
		{"Empty line", EmptyLine, false},
		{"Comment line", CommentLine, false},
		{"Secret URI line", SecretURILine, true},
		{"Key-value line", KeyValueLine, false},
		{"Key only line", KeyOnlyLine, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			line := Line{Type: tc.lineType}
			if result := line.IsSecret(); result != tc.expected {
				t.Errorf("IsSecret() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestLine_IsKeyValue(t *testing.T) {
	testCases := []struct {
		name     string
		lineType LineType
		expected bool
	}{
		{"Empty line", EmptyLine, false},
		{"Comment line", CommentLine, false},
		{"Secret URI line", SecretURILine, false},
		{"Key-value line", KeyValueLine, true},
		{"Key only line", KeyOnlyLine, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			line := Line{Type: tc.lineType}
			if result := line.IsKeyValue(); result != tc.expected {
				t.Errorf("IsKeyValue() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestLine_IsKeyOnly(t *testing.T) {
	testCases := []struct {
		name     string
		lineType LineType
		expected bool
	}{
		{"Empty line", EmptyLine, false},
		{"Comment line", CommentLine, false},
		{"Secret URI line", SecretURILine, false},
		{"Key-value line", KeyValueLine, false},
		{"Key only line", KeyOnlyLine, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			line := Line{Type: tc.lineType}
			if result := line.IsKeyOnly(); result != tc.expected {
				t.Errorf("IsKeyOnly() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestLine_IsValid(t *testing.T) {
	testCases := []struct {
		name     string
		lineType LineType
		expected bool
	}{
		{"Empty line", EmptyLine, false},
		{"Comment line", CommentLine, false},
		{"Secret URI line", SecretURILine, true},
		{"Key-value line", KeyValueLine, true},
		{"Key only line", KeyOnlyLine, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			line := Line{Type: tc.lineType}
			if result := line.IsValid(); result != tc.expected {
				t.Errorf("IsValid() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestLine_ToEnvEntry(t *testing.T) {
	tests := []struct {
		name     string
		line     Line
		expected functional.Result[env.Entry]
		wantErr  bool
	}{
		{
			name: "Empty line",
			line: Line{
				Content: "",
				Number:  1,
				Type:    EmptyLine,
				Trimmed: "",
			},
			expected: functional.Success(env.Entry{}),
			wantErr:  false,
		},
		{
			name: "Comment line",
			line: Line{
				Content: "# This is a comment",
				Number:  2,
				Type:    CommentLine,
				Trimmed: "# This is a comment",
			},
			expected: functional.Success(env.Entry{}),
			wantErr:  false,
		},
		{
			name: "Secret URI line",
			line: Line{
				Content: "sem://aws/secretsmanager/my-secret",
				Number:  3,
				Type:    SecretURILine,
				Trimmed: "sem://aws/secretsmanager/my-secret",
			},
			expected: functional.Success(env.NewEntry(3, "sem://aws/secretsmanager/my-secret", "")),
			wantErr:  false,
		},
		{
			name: "Key-value line",
			line: Line{
				Content: "API_KEY=12345",
				Number:  4,
				Type:    KeyValueLine,
				Trimmed: "API_KEY=12345",
			},
			expected: functional.Success(env.NewEntry(4, "API_KEY", "12345")),
			wantErr:  false,
		},
		{
			name: "Key only line",
			line: Line{
				Content: "DATABASE_URL",
				Number:  5,
				Type:    KeyOnlyLine,
				Trimmed: "DATABASE_URL",
			},
			expected: functional.Success(env.NewEntry(5, "DATABASE_URL", "")),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.line.ToEnvEntry()

			if tt.wantErr {
				if !result.IsFailure() {
					t.Errorf("Expected error, but got success")
				}
				return
			}

			if result.IsFailure() {
				t.Errorf("Unexpected error: %v", result.GetError())
				return
			}

			expectedEntry := tt.expected.Unwrap()
			actualEntry := result.Unwrap()

			if !reflect.DeepEqual(actualEntry, expectedEntry) {
				t.Errorf("Entry = %+v, want %+v", actualEntry, expectedEntry)
			}
		})
	}
}

func TestParseEnvLine(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		lineNum  int
		expected env.Entry
		wantErr  bool
	}{
		{
			name:     "Empty line",
			content:  "",
			lineNum:  1,
			expected: env.Entry{},
			wantErr:  false,
		},
		{
			name:     "Comment line",
			content:  "# This is a comment",
			lineNum:  2,
			expected: env.Entry{},
			wantErr:  false,
		},
		{
			name:     "Secret URI line",
			content:  "sem://aws/secretsmanager/my-secret",
			lineNum:  3,
			expected: env.NewEntry(3, "sem://aws/secretsmanager/my-secret", ""),
			wantErr:  false,
		},
		{
			name:     "Key-value line",
			content:  "API_KEY=12345",
			lineNum:  4,
			expected: env.NewEntry(4, "API_KEY", "12345"),
			wantErr:  false,
		},
		{
			name:     "Key-value with spaces",
			content:  "API_KEY = 12345",
			lineNum:  5,
			expected: env.NewEntry(5, "API_KEY", " 12345"),
			wantErr:  false,
		},
		{
			name:     "Key only line",
			content:  "DATABASE_URL",
			lineNum:  6,
			expected: env.NewEntry(6, "DATABASE_URL", ""),
			wantErr:  false,
		},
		{
			name:     "Key with empty value",
			content:  "EMPTY_KEY=",
			lineNum:  7,
			expected: env.NewEntry(7, "EMPTY_KEY", ""),
			wantErr:  false,
		},
		{
			name:     "Key-value with equals in value",
			content:  "CONNECTION=host=localhost port=5432",
			lineNum:  8,
			expected: env.NewEntry(8, "CONNECTION", "host=localhost port=5432"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseEnvLine(tt.content, tt.lineNum)

			if tt.wantErr {
				if !result.IsFailure() {
					t.Errorf("Expected error, but got success")
				}
				return
			}

			if result.IsFailure() {
				t.Errorf("Unexpected error: %v", result.GetError())
				return
			}

			entry := result.Unwrap()
			if !reflect.DeepEqual(entry, tt.expected) {
				t.Errorf("Entry = %+v, want %+v", entry, tt.expected)
			}
		})
	}
}

func TestPreprocessContent(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected []byte
	}{
		{
			name:     "No BOM, Unix line endings",
			content:  []byte("line1\nline2\nline3"),
			expected: []byte("line1\nline2\nline3"),
		},
		{
			name:     "With BOM, Unix line endings",
			content:  []byte{0xEF, 0xBB, 0xBF, 'l', 'i', 'n', 'e', '1', '\n', 'l', 'i', 'n', 'e', '2'},
			expected: []byte("line1\nline2"),
		},
		{
			name:     "No BOM, Windows line endings",
			content:  []byte("line1\r\nline2\r\nline3"),
			expected: []byte("line1\nline2\nline3"),
		},
		{
			name:     "No BOM, Mac line endings",
			content:  []byte("line1\rline2\rline3"),
			expected: []byte("line1\nline2\nline3"),
		},
		{
			name:     "With BOM, Mixed line endings",
			content:  []byte{0xEF, 0xBB, 0xBF, 'l', 'i', 'n', 'e', '1', '\r', '\n', 'l', 'i', 'n', 'e', '2', '\r', 'l', 'i', 'n', 'e', '3'},
			expected: []byte("line1\nline2\nline3"),
		},
		{
			name:     "Empty content",
			content:  []byte{},
			expected: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PreprocessContent(tt.content)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("PreprocessContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParsePlainFileContentResult(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected []env.Entry
		wantErr  bool
	}{
		{
			name:     "Empty file",
			content:  []byte(""),
			expected: []env.Entry{},
			wantErr:  false,
		},
		{
			name:     "File with comments and empty lines",
			content:  []byte("# Comment 1\n\n# Comment 2"),
			expected: []env.Entry{},
			wantErr:  false,
		},
		{
			name:    "File with key-value pairs",
			content: []byte("KEY1=value1\nKEY2=value2"),
			expected: []env.Entry{
				env.NewEntry(1, "KEY1", "value1"),
				env.NewEntry(2, "KEY2", "value2"),
			},
			wantErr: false,
		},
		{
			name:    "Mixed content file",
			content: []byte("# Header\nKEY1=value1\n\nsem://aws/secret\nKEY2="),
			expected: []env.Entry{
				env.NewEntry(2, "KEY1", "value1"),
				env.NewEntry(4, "sem://aws/secret", ""),
				env.NewEntry(5, "KEY2", ""),
			},
			wantErr: false,
		},
		{
			name:    "File with BOM and different line endings",
			content: []byte{0xEF, 0xBB, 0xBF, 'K', 'E', 'Y', '1', '=', 'v', 'a', 'l', '1', '\r', '\n', 'K', 'E', 'Y', '2', '=', 'v', 'a', 'l', '2'},
			expected: []env.Entry{
				env.NewEntry(1, "KEY1", "val1"),
				env.NewEntry(2, "KEY2", "val2"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParsePlainFileContentResult(tt.content)

			if tt.wantErr {
				if !result.IsFailure() {
					t.Errorf("Expected error, but got success")
				}
				return
			}

			if result.IsFailure() {
				t.Errorf("Unexpected error: %v", result.GetError())
				return
			}

			entries := result.Unwrap()
			if !reflect.DeepEqual(entries, tt.expected) {
				t.Errorf("Entries = %+v, want %+v", entries, tt.expected)
			}
		})
	}
}

func TestClassifyLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected LineType
	}{
		{"Empty line", "", EmptyLine},
		{"Comment line", "# This is a comment", CommentLine},
		{"Secret URI line", "sem://aws/secretsmanager/my-secret", SecretURILine},
		{"Key-value line", "API_KEY=12345", KeyValueLine},
		{"Key only line", "DATABASE_URL", KeyOnlyLine},
		// In the actual code flow, whitespace is trimmed before classifyLine is called
		{"Whitespace only", "", EmptyLine}, // using "" to represent trimmed whitespace
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyLine(tt.line)
			if result != tt.expected {
				t.Errorf("classifyLine(%q) = %d, want %d", tt.line, result, tt.expected)
			}
		})
	}
}

func TestParseEnvLineWithLargeJSONValue(t *testing.T) {
	// Test with a large JSON value similar to .cache.env
	content := `large_secret='{"api_keys":{"aws":"aws-api-key-67890","github":"github-api-key-abcdef","google":"google-api-key-12345"},"app_name":"test-application","cache":{"enabled":true,"max_size_mb":512,"ttl_seconds":3600},"contacts":[{"email":"admin@example.com","name":"Admin Team","phone":"+81-3-1234-5678"},{"email":"support@example.com","name":"Support Team","phone":"+81-3-8765-4321"}],"database":{"host":"db.example.com","max_connections":100,"password":"very_secure_password_123","port":5432,"ssl":true,"timeout_seconds":30,"username":"prod_db_user"},"environment":"production","feature_flags":{"beta_features":false,"maintenance_mode":false,"new_ui":true},"log_levels":{"development":"DEBUG","production":"ERROR","staging":"INFO"}}}'`
	lineNum := 1

	// Parse the line
	result := ParseEnvLine(content, lineNum)

	// Check for successful parsing
	if result.IsFailure() {
		t.Errorf("Failed to parse large JSON value: %v", result.GetError())
		return
	}

	entry := result.Unwrap()

	// Verify key and value were properly extracted
	if entry.Key != "large_secret" {
		t.Errorf("Expected key 'large_secret', got '%s'", entry.Key)
	}

	// Verify the value contains expected JSON structure (check a few key elements)
	expectedSubstrings := []string{
		"api_keys",
		"app_name",
		"very_secure_password_123",
		"maintenance_mode",
	}

	for _, substr := range expectedSubstrings {
		if !strings.Contains(entry.Value, substr) {
			t.Errorf("Expected value to contain '%s', but it doesn't", substr)
		}
	}

	// Check that single quotes were properly handled (they should be part of the value)
	if !strings.HasPrefix(entry.Value, "'") || !strings.HasSuffix(entry.Value, "'") {
		t.Errorf("Single quotes not properly preserved in value: %s", entry.Value)
	}
}

// TestParsePlainFileContentWithLargeJSON tests parsing a file containing a large JSON value
func TestParsePlainFileContentWithLargeJSON(t *testing.T) {
	// Create file content similar to .cache.env
	fileContent := []byte(`large_secret='{"api_keys":{"aws":"aws-api-key-67890","github":"github-api-key-abcdef","google":"google-api-key-12345"},"app_name":"test-application","cache":{"enabled":true,"max_size_mb":512,"ttl_seconds":3600},"contacts":[{"email":"admin@example.com","name":"Admin Team","phone":"+81-3-1234-5678"},{"email":"support@example.com","name":"Support Team","phone":"+81-3-8765-4321"}],"database":{"host":"db.example.com","max_connections":100,"password":"very_secure_password_123","port":5432,"ssl":true,"timeout_seconds":30,"username":"prod_db_user"},"environment":"production","feature_flags":{"beta_features":false,"maintenance_mode":false,"new_ui":true},"log_levels":{"development":"DEBUG","production":"ERROR","staging":"INFO"}}'`)

	// Parse the content
	result := ParsePlainFileContentResult(fileContent)

	// Check for successful parsing
	if result.IsFailure() {
		t.Errorf("Failed to parse file with large JSON: %v", result.GetError())
		return
	}

	entries := result.Unwrap()

	// Verify we got exactly one entry
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
		return
	}

	// Verify key and line number
	entry := entries[0]
	if entry.Key != "large_secret" {
		t.Errorf("Expected key 'large_secret', got '%s'", entry.Key)
	}

	if entry.Index != 1 {
		t.Errorf("Expected line number 1, got %d", entry.Index)
	}

	// Verify the JSON value was preserved including quotes
	if !strings.HasPrefix(entry.Value, "'") || !strings.HasSuffix(entry.Value, "'") {
		t.Errorf("Single quotes not properly preserved in value")
	}

	// Verify the JSON content is intact
	if !strings.Contains(entry.Value, "very_secure_password_123") {
		t.Errorf("Expected value to contain sensitive data, but it's missing")
	}
}
