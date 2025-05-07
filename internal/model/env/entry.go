// Package env provides environment variable related models and utilities
package env

import "github.com/gumi-tsd/secret-env-manager/internal/functional"

// Entry represents a key-value pair from an environment file line.
// This is used when parsing input configuration files.
type Entry struct {
	Index int    // Position in the original file
	Key   string // Environment variable name
	Value string // Environment variable value or secret URI
}

// NewEntry creates a new Entry with the specified values
func NewEntry(index int, key, value string) Entry {
	return Entry{
		Index: index,
		Key:   key,
		Value: value,
	}
}

// WithIndex returns a new Entry with the specified index
func (e Entry) WithIndex(index int) Entry {
	return Entry{
		Index: index,
		Key:   e.Key,
		Value: e.Value,
	}
}

// WithKey returns a new Entry with the specified key
func (e Entry) WithKey(key string) Entry {
	return Entry{
		Index: e.Index,
		Key:   key,
		Value: e.Value,
	}
}

// WithValue returns a new Entry with the specified value
func (e Entry) WithValue(value string) Entry {
	return Entry{
		Index: e.Index,
		Key:   e.Key,
		Value: value,
	}
}

// IsEmpty checks if both key and value are empty
func (e Entry) IsEmpty() bool {
	return e.Key == "" && e.Value == ""
}

// HasValue checks if the entry has a non-empty value
func (e Entry) HasValue() bool {
	return e.Value != ""
}

// HasKey checks if the entry has a non-empty key
func (e Entry) HasKey() bool {
	return e.Key != ""
}

// Map applies a function to transform an Entry
func (e Entry) Map(f func(Entry) Entry) Entry {
	return f(e)
}

// AsOption converts an Entry to an Option
func (e Entry) AsOption() functional.Option[Entry] {
	if e.IsEmpty() {
		return functional.None[Entry]()
	}
	return functional.Some(e)
}

// MapEntries applies a function to each entry in a slice and returns a new slice
func MapEntries(entries []Entry, f func(Entry) Entry) []Entry {
	result := make([]Entry, len(entries))
	for i, entry := range entries {
		result[i] = f(entry)
	}
	return result
}

// FilterEntries returns a new slice containing only entries that satisfy the predicate
func FilterEntries(entries []Entry, predicate func(Entry) bool) []Entry {
	var result []Entry
	for _, entry := range entries {
		if predicate(entry) {
			result = append(result, entry)
		}
	}
	return result
}

// ContainsKey checks if the slice of entries contains an entry with the given key
func ContainsKey(entries []Entry, key string) bool {
	for _, entry := range entries {
		if entry.Key == key {
			return true
		}
	}
	return false
}
