// Package formatting provides text formatting and colorization utilities.
//
// color.go handles ANSI terminal color formatting for console output.
package formatting

import (
	"fmt"
)

// ANSI Color escape codes
const (
	Reset = "\033[0m" // Reset terminal to default state
	Bold  = "\033[1m" // Bold text style
	// Foreground colors
	FgBlack   = "\033[30m"
	FgRed     = "\033[31m"
	FgGreen   = "\033[32m"
	FgYellow  = "\033[33m"
	FgBlue    = "\033[34m"
	FgMagenta = "\033[35m"
	FgCyan    = "\033[36m"
	FgWhite   = "\033[37m"
	// High-intensity colors
	FgHiBlack   = "\033[90m"
	FgHiRed     = "\033[91m"
	FgHiGreen   = "\033[92m"
	FgHiYellow  = "\033[93m"
	FgHiBlue    = "\033[94m"
	FgHiMagenta = "\033[95m"
	FgHiCyan    = "\033[96m"
	FgHiWhite   = "\033[97m"
)

var (
	noColor = false // Controls whether colors are enabled
)

// ColorFunc is a function type that applies color formatting to a string
type ColorFunc func(string) string

// Pre-defined colorizers for different purposes
var (
	// Content type colorizers
	colorizeKey   = makeColorizer(FgCyan)
	colorizeValue = makeColorizer(FgGreen)

	// Message type colorizers
	colorizeHeader  = makeColorizer(FgHiWhite + Bold)
	colorizeSuccess = makeColorizer(FgGreen + Bold)
	colorizeError   = makeColorizer(FgRed + Bold)
	colorizeWarning = makeColorizer(FgYellow)
	colorizeHint    = makeColorizer(FgHiBlack)
	colorizeSection = makeColorizer(FgMagenta + Bold)
	colorizeInfo    = makeColorizer(FgBlue + Bold)
)

// makeColorizer creates a color formatting function for the given ANSI color code
func makeColorizer(colorCode string) ColorFunc {
	return func(s string) string {
		if noColor {
			return s
		}
		return colorCode + s + Reset
	}
}

// ColorizeKey returns a colorized key string
func ColorizeKey(key string) string {
	return colorizeKey(key)
}

// ColorizeValue returns a colorized value string
func ColorizeValue(value string) string {
	return colorizeValue(value)
}

// ColorizeKeyValue returns a colorized key-value pair with optional quotes
func ColorizeKeyValue(key, value string, useQuotes bool) string {
	// If colors are disabled, use plain formatting
	if noColor {
		if useQuotes {
			// Use double quotes for compatibility with the colored version
			return fmt.Sprintf("%s=\"%s\"", key, value)
		}
		return fmt.Sprintf("%s=%s", key, value)
	}

	// Apply colors
	colorizedKey := ColorizeKey(key)
	colorizedValue := ColorizeValue(value)

	if useQuotes {
		return fmt.Sprintf("%s=\"%s\"", colorizedKey, colorizedValue)
	}
	return fmt.Sprintf("%s=%s", colorizedKey, colorizedValue)
}

// Message formatting functions

// Success formats text as success message (green bold)
func Success(format string, a ...interface{}) string {
	return colorizeSuccess(fmt.Sprintf(format, a...))
}

// Error formats text as error message (red bold)
func Error(format string, a ...interface{}) string {
	return colorizeError(fmt.Sprintf(format, a...))
}

// Warning formats text as warning message (yellow)
func Warning(format string, a ...interface{}) string {
	return colorizeWarning(fmt.Sprintf(format, a...))
}

// Hint formats text as hint/help message (gray)
func Hint(format string, a ...interface{}) string {
	return colorizeHint(fmt.Sprintf(format, a...))
}

// Info formats text as information message (blue bold)
func Info(format string, a ...interface{}) string {
	return colorizeInfo(fmt.Sprintf(format, a...))
}

// FormatHeader formats text as header (white bold)
func FormatHeader(format string, a ...interface{}) string {
	return colorizeHeader(fmt.Sprintf(format, a...))
}

// Color control functions

// DisableColors turns off all color output
func DisableColors() {
	noColor = true
}

// EnableColors turns on color output
func EnableColors() {
	noColor = false
}

// IsColorEnabled returns whether colors are currently enabled
func IsColorEnabled() bool {
	return !noColor
}
