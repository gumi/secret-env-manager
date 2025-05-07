// Package functional provides functional programming utilities and patterns
package functional

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T any, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Filter returns a new slice containing only the elements that satisfy the predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
	// Initialize with empty slice instead of nil to match test expectations
	result := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce combines all elements in a slice into a single value using a reducer function
func Reduce[T any, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

// ForEach applies a function to each element in a slice (for side effects)
// This is not a purely functional operation but is often useful
func ForEach[T any](slice []T, fn func(T)) {
	// Special handling for empty slices to match test expectations
	if len(slice) == 0 {
		return
	}

	for _, v := range slice {
		fn(v)
	}
}

// All returns true if all elements in the slice satisfy the predicate
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Any returns true if at least one element in the slice satisfies the predicate
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

// Contains returns true if the slice contains the specified value
func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Find returns the first element that satisfies the predicate, wrapped in an Option
func Find[T any](slice []T, predicate func(T) bool) Option[T] {
	for _, v := range slice {
		if predicate(v) {
			return Some(v)
		}
	}
	return None[T]()
}

// Zip combines two slices into a single slice of pairs
func Zip[T any, U any](slice1 []T, slice2 []U) []Pair[T, U] {
	minLen := len(slice1)
	if len(slice2) < minLen {
		minLen = len(slice2)
	}

	result := make([]Pair[T, U], minLen)
	for i := 0; i < minLen; i++ {
		result[i] = MakePair(slice1[i], slice2[i])
	}
	return result
}

// Pair represents a pair of values
type Pair[T any, U any] struct {
	First  T
	Second U
}

// MakePair creates a new pair
func MakePair[T any, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}
