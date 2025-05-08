package formatting

import (
	"fmt"
	"testing"
)

func TestMakeColorizer(t *testing.T) {
	// Test with color enabled
	EnableColors()
	colorizer := makeColorizer(FgRed)
	result := colorizer("test")
	expected := FgRed + "test" + Reset
	if result != expected {
		t.Errorf("makeColorizer with colors enabled = %q, want %q", result, expected)
	}

	// Test with color disabled
	DisableColors()
	result = colorizer("test")
	expected = "test"
	if result != expected {
		t.Errorf("makeColorizer with colors disabled = %q, want %q", result, expected)
	}

	// Reset to default state
	EnableColors()
}

func TestColorizeKey(t *testing.T) {
	EnableColors()
	result := ColorizeKey("key")
	expected := FgCyan + "key" + Reset
	if result != expected {
		t.Errorf("ColorizeKey = %q, want %q", result, expected)
	}

	DisableColors()
	result = ColorizeKey("key")
	expected = "key"
	if result != expected {
		t.Errorf("ColorizeKey with colors disabled = %q, want %q", result, expected)
	}

	EnableColors()
}

func TestColorizeValue(t *testing.T) {
	EnableColors()
	result := ColorizeValue("value")
	expected := FgGreen + "value" + Reset
	if result != expected {
		t.Errorf("ColorizeValue = %q, want %q", result, expected)
	}

	DisableColors()
	result = ColorizeValue("value")
	expected = "value"
	if result != expected {
		t.Errorf("ColorizeValue with colors disabled = %q, want %q", result, expected)
	}

	EnableColors()
}

func TestColorizeKeyValue(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		useQuotes bool
		colored   bool
		expected  string
	}{
		{
			name:      "Basic without quotes",
			key:       "KEY",
			value:     "value",
			useQuotes: false,
			colored:   true,
			expected:  fmt.Sprintf("%sKEY%s=%svalue%s", FgCyan, Reset, FgGreen, Reset),
		},
		{
			name:      "Basic with quotes",
			key:       "KEY",
			value:     "value",
			useQuotes: true,
			colored:   true,
			expected:  fmt.Sprintf("%sKEY%s=\"%svalue%s\"", FgCyan, Reset, FgGreen, Reset),
		},
		{
			name:      "Without color, no quotes",
			key:       "KEY",
			value:     "value",
			useQuotes: false,
			colored:   false,
			expected:  "KEY=value",
		},
		{
			name:      "Without color, with quotes",
			key:       "KEY",
			value:     "value",
			useQuotes: true,
			colored:   false,
			expected:  "KEY=\"value\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.colored {
				EnableColors()
			} else {
				DisableColors()
			}

			result := ColorizeKeyValue(tt.key, tt.value, tt.useQuotes)
			if result != tt.expected {
				t.Errorf("ColorizeKeyValue(%q, %q, %v) = %q, want %q",
					tt.key, tt.value, tt.useQuotes, result, tt.expected)
			}
		})
	}

	// Reset to default state
	EnableColors()
}

func TestColorFormats(t *testing.T) {
	formats := []struct {
		name     string
		function func(string, ...interface{}) string
		color    string
	}{
		{"Success", Success, FgGreen + Bold},
		{"Error", Error, FgRed + Bold},
		{"Warning", Warning, FgYellow},
		{"Hint", Hint, FgHiBlack},
		{"Info", Info, FgBlue + Bold},
		{"FormatHeader", FormatHeader, FgHiWhite + Bold},
	}

	for _, format := range formats {
		t.Run(format.name+"_colored", func(t *testing.T) {
			EnableColors()
			result := format.function("Test %s", "message")
			expected := format.color + "Test message" + Reset
			if result != expected {
				t.Errorf("%s() = %q, want %q", format.name, result, expected)
			}
		})

		t.Run(format.name+"_no_color", func(t *testing.T) {
			DisableColors()
			result := format.function("Test %s", "message")
			expected := "Test message"
			if result != expected {
				t.Errorf("%s() with colors disabled = %q, want %q", format.name, result, expected)
			}
		})
	}

	// Reset to default state
	EnableColors()
}

func TestIsColorEnabled(t *testing.T) {
	// Test default state (should be enabled)
	EnableColors()
	if !IsColorEnabled() {
		t.Error("IsColorEnabled() = false, want true after EnableColors()")
	}

	// Test after disabling
	DisableColors()
	if IsColorEnabled() {
		t.Error("IsColorEnabled() = true, want false after DisableColors()")
	}

	// Reset to default state
	EnableColors()
}
