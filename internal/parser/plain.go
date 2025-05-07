// Package parser provides utilities for parsing environment files
package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/env"
)

// LineType represents the classification of a line in an environment file
type LineType int

const (
	EmptyLine LineType = iota
	CommentLine
	SecretURILine
	KeyValueLine
	KeyOnlyLine
)

// Line represents a single line from a file with its properties
type Line struct {
	Content string
	Number  int
	Type    LineType
	Trimmed string
}

// IsEmpty returns whether the line is empty
func (l Line) IsEmpty() bool {
	return l.Type == EmptyLine
}

// IsComment returns whether the line is a comment
func (l Line) IsComment() bool {
	return l.Type == CommentLine
}

// IsSecret returns whether the line is a secret URI
func (l Line) IsSecret() bool {
	return l.Type == SecretURILine
}

// IsKeyValue returns whether the line is a key-value pair
func (l Line) IsKeyValue() bool {
	return l.Type == KeyValueLine
}

// IsKeyOnly returns whether the line contains only a key
func (l Line) IsKeyOnly() bool {
	return l.Type == KeyOnlyLine
}

// IsValid returns whether the line represents a valid entry
func (l Line) IsValid() bool {
	return l.Type == SecretURILine || l.Type == KeyValueLine || l.Type == KeyOnlyLine
}

// ToEnvEntry converts a Line to an EnvEntry
func (l Line) ToEnvEntry() functional.Result[env.Entry] {
	switch l.Type {
	case EmptyLine, CommentLine:
		return functional.Success(env.Entry{})

	case SecretURILine:
		return functional.Success(env.NewEntry(l.Number, l.Trimmed, ""))

	case KeyValueLine:
		return parseKeyValueLine(l)

	case KeyOnlyLine:
		return functional.Success(env.NewEntry(l.Number, l.Trimmed, ""))
	}

	// This should never happen due to exhaustive handling
	return functional.Success(env.Entry{})
}

// ContentLines represents the lines of content with their properties
type ContentLines struct {
	Lines []Line
	Error error
}

// Exported functions

// NewLine creates a new Line instance
func NewLine(content string, number int) Line {
	trimmed := strings.TrimSpace(content)
	return Line{
		Content: content,
		Number:  number,
		Type:    classifyLine(trimmed),
		Trimmed: trimmed,
	}
}

// ParseEnvLine parses a single line from an environment file.
// It handles empty lines, comments, secret URIs (sem://), key-value pairs, and keys without values.
// Returns an empty EnvEntry for empty or comment lines.
func ParseEnvLine(content string, idx int) functional.Result[env.Entry] {
	line := NewLine(content, idx)
	return line.ToEnvEntry()
}

// ParsePlainFileContentResult parses the entire content of an environment file using the Result monad.
func ParsePlainFileContentResult(content []byte) functional.Result[[]env.Entry] {
	processedContent := PreprocessContent(content)
	return parseLinesResult(processedContent)
}

// ParsePlainFileContent parses the entire content of an environment file.
// This is a compatibility wrapper for the monadic version.
func ParsePlainFileContent(content []byte) ([]env.Entry, error) {
	result := ParsePlainFileContentResult(content)
	if result.IsFailure() {
		return nil, result.GetError()
	}
	return result.Unwrap(), nil
}

// FilterValidEntries filters out empty and comment lines
func FilterValidEntries(entries []env.Entry) []env.Entry {
	return functional.Filter(entries, func(entry env.Entry) bool {
		return entry != (env.Entry{})
	})
}

// Unexported helper functions

// classifyLine determines the type of line in an environment file
func classifyLine(line string) LineType {
	// Note that line has already been processed with strings.TrimSpace() at this point
	if line == "" {
		return EmptyLine
	}

	if strings.HasPrefix(line, "#") {
		return CommentLine
	}

	if strings.HasPrefix(line, "sem://") {
		return SecretURILine
	}

	if strings.Contains(line, "=") {
		return KeyValueLine
	}

	return KeyOnlyLine
}

// parseKeyValueLine parses a line containing key=value format
func parseKeyValueLine(line Line) functional.Result[env.Entry] {
	eqIndex := strings.Index(line.Trimmed, "=")
	if eqIndex == -1 {
		return functional.Failure[env.Entry](
			fmt.Errorf("invalid key-value line: %s", line.Content))
	}

	key := strings.TrimSpace(line.Trimmed[:eqIndex])
	value := line.Trimmed[eqIndex+1:] // Take everything after the first equals sign

	return functional.Success(env.NewEntry(line.Number, key, value))
}

// PreprocessContent prepares file content for parsing by removing BOM and normalizing line endings
func PreprocessContent(content []byte) []byte {
	return functional.ApplyAll(
		content,
		removeBOM,
		normalizeLineEndings,
	)
}

// removeBOM strips the UTF-8 Byte Order Mark (BOM) if present.
func removeBOM(b []byte) []byte {
	if bytes.HasPrefix(b, []byte{0xEF, 0xBB, 0xBF}) {
		return b[3:]
	}
	return b
}

// normalizeLineEndings converts CRLF (Windows) and CR (Classic Mac) line endings to LF (Unix).
func normalizeLineEndings(b []byte) []byte {
	// Replace CRLF with LF
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
	// Replace remaining CR with LF
	return bytes.ReplaceAll(b, []byte("\r"), []byte("\n"))
}

// NewContentLinesResult creates ContentLines from raw content using Result monad
func NewContentLinesResult(content []byte) functional.Result[[]Line] {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	lineNumber := 0
	var lines []Line

	for scanner.Scan() {
		lineNumber++
		lineContent := scanner.Text()
		lines = append(lines, NewLine(lineContent, lineNumber))
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return functional.Failure[[]Line](
			fmt.Errorf("error scanning file content: %w", scanErr))
	}

	return functional.Success(lines)
}

// IsValidEntry is a predicate that checks if an entry is valid (not empty)
func IsValidEntry(entry env.Entry) bool {
	return entry != (env.Entry{})
}

// parseLinesResult processes each line in the content and builds an array of environment entries
// using the Result monad for consistent error handling
func parseLinesResult(content []byte) functional.Result[[]env.Entry] {
	// Parse all lines
	linesResult := NewContentLinesResult(content)
	if linesResult.IsFailure() {
		return functional.Failure[[]env.Entry](linesResult.GetError())
	}

	lines := linesResult.Unwrap()

	// Convert each line to an entry and collect results
	entryResults := functional.Map(lines, func(line Line) functional.Result[env.Entry] {
		return line.ToEnvEntry()
	})

	// Combine all results or return first error
	return combineEntryResults(entryResults)
}

// combineEntryResults combines multiple Result[env.Entry] into a single Result[[]env.Entry]
// Returns the first error encountered, if any
func combineEntryResults(results []functional.Result[env.Entry]) functional.Result[[]env.Entry] {
	entries := make([]env.Entry, 0, len(results))

	for i, result := range results {
		if result.IsFailure() {
			return functional.Failure[[]env.Entry](
				fmt.Errorf("error parsing line %d: %w", i+1, result.GetError()))
		}
		entries = append(entries, result.Unwrap())
	}

	// Filter out invalid entries
	return functional.Success(FilterValidEntries(entries))
}
