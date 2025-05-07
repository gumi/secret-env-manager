// Package formatting provides text formatting and colorization utilities.
package formatting

import (
	"fmt"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/text"
)

// ANSI Color escape codes
const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	FgBlack     = "\033[30m"
	FgRed       = "\033[31m"
	FgGreen     = "\033[32m"
	FgYellow    = "\033[33m"
	FgBlue      = "\033[34m"
	FgMagenta   = "\033[35m"
	FgCyan      = "\033[36m"
	FgWhite     = "\033[37m"
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
	// Color formatter functions for various elements
	colorizeKey     = makeColorizer(FgCyan)
	colorizeValue   = makeColorizer(FgGreen)
	colorizeHeader  = makeColorizer(FgHiWhite + Bold)
	colorizeSuccess = makeColorizer(FgGreen + Bold)
	colorizeError   = makeColorizer(FgRed + Bold)
	colorizeWarning = makeColorizer(FgYellow)
	colorizeHint    = makeColorizer(FgHiBlack)
	colorizeSection = makeColorizer(FgMagenta + Bold)
	colorizeInfo    = makeColorizer(FgBlue + Bold)

	// Global flag to control color output
	noColor = false
)

// ColorFunc is a function type that applies color formatting to a string
type ColorFunc func(string) string

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

// ColorizeKeyValue returns a colorized key-value pair
func ColorizeKeyValue(key, value string, useQuotes bool) string {
	colorizedKey := ColorizeKey(key)
	colorizedValue := ColorizeValue(value)
	if useQuotes {
		return fmt.Sprintf("%s=\"%s\"", colorizedKey, colorizedValue)
	}
	return fmt.Sprintf("%s=%s", colorizedKey, colorizedValue)
}

// ColorizeKeyValues returns a slice of colorized key-value pairs
func ColorizeKeyValues(keys []string, values map[string]string, useQuotes bool) []string {
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		if value, ok := values[key]; ok {
			lines = append(lines, ColorizeKeyValue(key, value, useQuotes))
		}
	}
	return lines
}

// FormatList formats and colorizes a list of strings with optional prefixes
func FormatList(items []string, prefix string, itemTransform text.StringTransformer) string {
	if itemTransform == nil {
		itemTransform = text.Identity
	}

	transformedItems := text.MapStrings(items, itemTransform)
	lines := text.MapStringsWithIndex(transformedItems, func(i int, s string) string {
		return fmt.Sprintf("%s%d. %s", prefix, i+1, s)
	})

	return strings.Join(lines, "\n")
}

// FormatTitle formats a string as a title with optional padding
func FormatTitle(title string, padding int) string {
	padStr := strings.Repeat(" ", padding)
	return padStr + colorizeSection(title)
}

// Success formats text as success message
func Success(format string, a ...interface{}) string {
	return colorizeSuccess(fmt.Sprintf(format, a...))
}

// Error formats text as error message
func Error(format string, a ...interface{}) string {
	return colorizeError(fmt.Sprintf(format, a...))
}

// Warning formats text as warning message
func Warning(format string, a ...interface{}) string {
	return colorizeWarning(fmt.Sprintf(format, a...))
}

// Hint formats text as hint/help message
func Hint(format string, a ...interface{}) string {
	return colorizeHint(fmt.Sprintf(format, a...))
}

// Info formats text as information message
func Info(format string, a ...interface{}) string {
	return colorizeInfo(fmt.Sprintf(format, a...))
}

// FormatHeader formats text as header
func FormatHeader(format string, a ...interface{}) string {
	return colorizeHeader(fmt.Sprintf(format, a...))
}

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
