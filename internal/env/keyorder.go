// Package env provides utilities for working with environment variables.
package env

import (
	"sort"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	modelenv "github.com/gumi-tsd/secret-env-manager/internal/model/env"
)

// FilterAndSortKeys ensures that the environment variables are processed
// in a consistent order with priority for keys in the sortedKeys list.
// This is a pure function.
func FilterAndSortKeys(values EnvVarMap, sortedKeys []string) []string {
	// 必ず最終的なキーをアルファベット順にソートする
	// 順序を一貫させるため、入力のsortedKeysはプレフィックスとして尊重するが、最終的に全キーをソートする
	allKeys := SortKeys(values)

	// If no specific order is specified, just return all keys lexicographically sorted
	if len(sortedKeys) == 0 {
		return allKeys
	}

	// Ensure we only include keys that actually exist in the values map
	// and maintain their order, but also sort them for consistency
	orderedKeys := MergeSortedKeys(sortedKeys, values)
	sort.Strings(orderedKeys)
	return orderedKeys
}

// SortKeys returns a lexicographically sorted list of all keys in the map
// This is a pure function.
func SortKeys(values EnvVarMap) []string {
	keys := make([]string, 0, len(values))

	for k := range values {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

// MergeSortedKeys merges user-specified key order with actual keys from the map.
// It preserves the order of sortedKeys for keys that exist in the values map,
// and appends any remaining keys in alphabetical order.
// This is a pure function.
func MergeSortedKeys(sortedKeys []string, values EnvVarMap) []string {
	result := make([]string, 0, len(values))
	processed := make(map[string]bool, len(sortedKeys))

	// First pass: add keys in requested order if they exist in values
	for _, key := range sortedKeys {
		if _, exists := values[key]; exists {
			result = append(result, key)
			processed[key] = true
		}
	}

	// Second pass: collect remaining keys
	remaining := make([]string, 0, len(values)-len(processed))
	for key := range values {
		if !processed[key] {
			remaining = append(remaining, key)
		}
	}

	// Sort remaining keys for consistency
	sort.Strings(remaining)

	// Combine results (preserving the order)
	return append(result, remaining...)
}

// SortEnvMapKeysResult sorts environment map keys using Result monad pattern
// This is a pure function.
func SortEnvMapKeysResult(values EnvVarMap) functional.Result[[]string] {
	if values == nil {
		return functional.Success([]string{})
	}

	keys := SortKeys(values)
	return functional.Success(keys)
}

// GroupKeysByPrefix groups keys by their common prefixes for better organization
// This is a pure function.
func GroupKeysByPrefix(keys []string, separator string) map[string][]string {
	groups := make(map[string][]string)

	for _, key := range keys {
		prefix := GetKeyPrefix(key, separator)
		groups[prefix] = append(groups[prefix], key)
	}

	// Sort keys within each group
	for prefix, groupKeys := range groups {
		sort.Strings(groupKeys)
		groups[prefix] = groupKeys
	}

	return groups
}

// GetKeyPrefix extracts the prefix part of a key based on a separator
// This is a pure function.
func GetKeyPrefix(key, separator string) string {
	if separator == "" {
		return key
	}

	parts := strings.SplitN(key, separator, 2)
	if len(parts) > 1 {
		return parts[0]
	}

	return key
}

// SortPrefixGroups returns prefixes in sorted order
// This is a pure function.
func SortPrefixGroups(groups map[string][]string) []string {
	prefixes := make([]string, 0, len(groups))

	for prefix := range groups {
		prefixes = append(prefixes, prefix)
	}

	sort.Strings(prefixes)
	return prefixes
}

// GetSortedKeysByGroups returns all keys organized by their prefix groups in order
// This is a pure function.
func GetSortedKeysByGroups(values EnvVarMap, separator string) []string {
	allKeys := SortKeys(values)

	// If no separator, just return the sorted keys
	if separator == "" {
		return allKeys
	}

	groups := GroupKeysByPrefix(allKeys, separator)
	prefixes := SortPrefixGroups(groups)

	// Build result with keys sorted within each prefix group
	result := make([]string, 0, len(allKeys))
	for _, prefix := range prefixes {
		result = append(result, groups[prefix]...)
	}

	return result
}

// OrganizeKeyOrder creates an ordered list of keys based on:
// 1. Original order in the entries file
// 2. Additional keys from retrieved secrets
// This is a pure function.
func OrganizeKeyOrder(entries []modelenv.Entry, values EnvVarMap) []string {

	// 最終的なキーリストを常にアルファベット順にソートして一貫性を保つ
	// これによって実行ごとの順序のばらつきを防ぐ
	allKeys := make([]string, 0, len(values))
	for k := range values {
		allKeys = append(allKeys, k)
	}
	sort.Strings(allKeys)

	return allKeys
}

// ExtractKeysFromEntries extracts keys from a list of env entries
// This is a pure function.
func ExtractKeysFromEntries(entries []modelenv.Entry) []string {
	// Collect keys from entries, preserving order
	keys := make([]string, 0, len(entries))
	seen := make(map[string]bool, len(entries))

	for _, entry := range entries {
		// Skip entries without keys or already processed keys
		if entry.Key == "" || seen[entry.Key] {
			continue
		}

		keys = append(keys, entry.Key)
		seen[entry.Key] = true
	}

	return keys
}
