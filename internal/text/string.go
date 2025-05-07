// Package text provides string manipulation and text processing utilities.
package text

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// StringTransformer represents a function that transforms one string into another
type StringTransformer func(string) string

// StringPredicate represents a function that tests a string
type StringPredicate func(string) bool

// Identity returns the input string unchanged
func Identity(s string) string {
	return s
}

// CleanControlCharsResult removes control characters from a string efficiently, using Result monad
func CleanControlCharsResult(s string) functional.Result[string] {
	if s == "" {
		return functional.Success("")
	}

	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		if r == ' ' {
			b.WriteRune(r)
		} else if !unicode.IsControl(r) && !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return functional.Success(b.String())
}

// CleanControlChars removes control characters from a string efficiently.
func CleanControlChars(s string) string {
	result := CleanControlCharsResult(s)
	return result.Unwrap()
}

// JoinWithSeparatorResult joins strings with a separator, skipping empty strings with Result monad
func JoinWithSeparatorResult(separator string, parts ...string) functional.Result[string] {
	var nonEmpty []string
	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}
	return functional.Success(strings.Join(nonEmpty, separator))
}

// JoinWithSeparator joins strings with a separator, skipping empty strings
func JoinWithSeparator(separator string, parts ...string) string {
	var nonEmpty []string
	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}
	return strings.Join(nonEmpty, separator)
}

// SplitAndTrimResult splits a string by separator and trims each part with Result monad
func SplitAndTrimResult(s, separator string) functional.Result[[]string] {
	return functional.MapResultTo(
		SplitResult(s, separator),
		func(parts []string) []string {
			return MapStrings(parts, strings.TrimSpace)
		},
	)
}

// SplitAndTrim splits a string by separator and trims each part
func SplitAndTrim(s, separator string) []string {
	parts := strings.Split(s, separator)
	return MapStrings(parts, strings.TrimSpace)
}

// MapStrings applies a transformer function to each string in a slice
func MapStrings(strs []string, transformer StringTransformer) []string {
	result := make([]string, len(strs))
	for i, s := range strs {
		result[i] = transformer(s)
	}
	return result
}

// MapStringsWithIndex applies a transformer function with index to each string in a slice
func MapStringsWithIndex(strs []string, transformer func(int, string) string) []string {
	result := make([]string, len(strs))
	for i, s := range strs {
		result[i] = transformer(i, s)
	}
	return result
}

// FilterStrings returns a new slice containing only strings that satisfy the predicate
func FilterStrings(strs []string, predicate StringPredicate) []string {
	// Initialize with empty slice instead of nil
	result := []string{}
	for _, s := range strs {
		if predicate(s) {
			result = append(result, s)
		}
	}
	return result
}

// IsNotEmpty returns true if the string is not empty
func IsNotEmpty(s string) bool {
	return len(s) > 0
}

// Chain combines multiple string transformers into one
func Chain(transformers ...StringTransformer) StringTransformer {
	return func(s string) string {
		result := s
		for _, transformer := range transformers {
			result = transformer(result)
		}
		return result
	}
}

// ComposePredicate combines predicates with logical AND
func ComposePredicate(predicates ...StringPredicate) StringPredicate {
	return func(s string) bool {
		for _, predicate := range predicates {
			if !predicate(s) {
				return false
			}
		}
		return true
	}
}

// ComposePredicateOr combines predicates with logical OR
func ComposePredicateOr(predicates ...StringPredicate) StringPredicate {
	return func(s string) bool {
		// If no predicates provided, return false (nothing matches)
		if len(predicates) == 0 {
			return false
		}

		for _, predicate := range predicates {
			if predicate(s) {
				return true
			}
		}
		return false
	}
}

// SplitOption splits a string and returns None if the result is empty
func SplitOption(s, separator string) functional.Option[[]string] {
	if s == "" {
		return functional.None[[]string]()
	}
	return functional.Some(strings.Split(s, separator))
}

// SplitResult splits a string and returns a Result monad
func SplitResult(s, separator string) functional.Result[[]string] {
	if s == "" {
		return functional.Failure[[]string](
			fmt.Errorf("cannot split an empty string"))
	}
	return functional.Success(strings.Split(s, separator))
}

// TrimSpace is a functional version of strings.TrimSpace
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// ToLower is a functional version of strings.ToLower
func ToLower(s string) string {
	return strings.ToLower(s)
}

// Contains returns a predicate that tests if a string contains the substring
func Contains(substr string) StringPredicate {
	return func(s string) bool {
		return strings.Contains(s, substr)
	}
}

// HasPrefix returns a predicate that tests if a string has the given prefix
func HasPrefix(prefix string) StringPredicate {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}

// HasSuffix returns a predicate that tests if a string has the given suffix
func HasSuffix(suffix string) StringPredicate {
	return func(s string) bool {
		return strings.HasSuffix(s, suffix)
	}
}

// JoinNonEmptyResult joins only non-empty strings with a separator, returning a Result monad
func JoinNonEmptyResult(separator string, strs ...string) functional.Result[string] {
	return JoinWithSeparatorResult(separator, strs...)
}

// JoinNonEmpty joins only non-empty strings with a separator
func JoinNonEmpty(separator string, strs ...string) string {
	return JoinWithSeparator(separator, strs...)
}

// Trim removes cutset from both ends of the string
func Trim(s, cutset string) string {
	return strings.Trim(s, cutset)
}

// Replace replaces occurrences of old with new in s.
// If n < 0, there is no limit on the number of replacements.
func Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

// ToUpper converts all characters in s to uppercase
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// Index returns the index of the first instance of substr in s, or -1 if not found
func Index(s, substr string) int {
	return strings.Index(s, substr)
}

// LastIndex returns the index of the last instance of substr in s, or -1 if not found
func LastIndex(s, substr string) int {
	return strings.LastIndex(s, substr)
}

// Join concatenates the elements of a to create a single string with sep between elements
func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

// StartsWith checks if the string starts with the given prefix
func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// EndsWith checks if the string ends with the given suffix
func EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Split splits the string s around sep
func Split(s, sep string) []string {
	return strings.Split(s, sep)
}
