package functional

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) string
		expected []string
	}{
		{
			name:     "Empty slice",
			input:    []int{},
			fn:       func(i int) string { return strconv.Itoa(i) },
			expected: []string{},
		},
		{
			name:     "Convert int to string",
			input:    []int{1, 2, 3, 4, 5},
			fn:       func(i int) string { return strconv.Itoa(i) },
			expected: []string{"1", "2", "3", "4", "5"},
		},
		{
			name:     "Double values",
			input:    []int{1, 2, 3},
			fn:       func(i int) string { return strconv.Itoa(i * 2) },
			expected: []string{"2", "4", "6"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.fn)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Map() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "Empty slice",
			input:     []int{},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  []int{},
		},
		{
			name:      "Filter even numbers",
			input:     []int{1, 2, 3, 4, 5, 6},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  []int{2, 4, 6},
		},
		{
			name:      "Filter odd numbers",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(i int) bool { return i%2 != 0 },
			expected:  []int{1, 3, 5},
		},
		{
			name:      "Filter none",
			input:     []int{1, 3, 5, 7},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  []int{},
		},
		{
			name:      "Filter all",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(i int) bool { return i > 0 },
			expected:  []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.input, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Filter() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		initial  int
		reducer  func(int, int) int
		expected int
	}{
		{
			name:     "Empty slice",
			input:    []int{},
			initial:  0,
			reducer:  func(acc, val int) int { return acc + val },
			expected: 0,
		},
		{
			name:     "Sum numbers",
			input:    []int{1, 2, 3, 4, 5},
			initial:  0,
			reducer:  func(acc, val int) int { return acc + val },
			expected: 15,
		},
		{
			name:     "Multiply numbers",
			input:    []int{1, 2, 3, 4},
			initial:  1,
			reducer:  func(acc, val int) int { return acc * val },
			expected: 24,
		},
		{
			name:    "Max value",
			input:   []int{3, 7, 2, 9, 5},
			initial: 0,
			reducer: func(acc, val int) int {
				if val > acc {
					return val
				}
				return acc
			},
			expected: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reduce(tt.input, tt.initial, tt.reducer)
			if result != tt.expected {
				t.Errorf("Reduce() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestForEach(t *testing.T) {
	tests := []struct {
		name  string
		input []string
	}{
		{
			name:  "Empty slice",
			input: []string{},
		},
		{
			name:  "String slice",
			input: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use a slice to capture the values seen by ForEach
			// Initialize with empty slice instead of nil
			var seen = make([]string, 0)

			// Call ForEach with a function that appends values to our collector
			ForEach(tt.input, func(s string) {
				seen = append(seen, s)
			})

			// Verify that each element was processed exactly once
			if !reflect.DeepEqual(seen, tt.input) {
				t.Errorf("ForEach() processed %v, want %v", seen, tt.input)
			}
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  bool
	}{
		{
			name:      "Empty slice",
			input:     []int{},
			predicate: func(i int) bool { return i > 0 },
			expected:  true, // All returns true for empty slices
		},
		{
			name:      "All positive",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(i int) bool { return i > 0 },
			expected:  true,
		},
		{
			name:      "Not all positive",
			input:     []int{1, 2, 0, 4, 5},
			predicate: func(i int) bool { return i > 0 },
			expected:  false,
		},
		{
			name:      "All even",
			input:     []int{2, 4, 6, 8},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  true,
		},
		{
			name:      "Not all even",
			input:     []int{2, 4, 5, 8},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := All(tt.input, tt.predicate)
			if result != tt.expected {
				t.Errorf("All() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  bool
	}{
		{
			name:      "Empty slice",
			input:     []int{},
			predicate: func(i int) bool { return i > 0 },
			expected:  false, // Any returns false for empty slices
		},
		{
			name:      "Some positive",
			input:     []int{-1, -2, 3, -4, -5},
			predicate: func(i int) bool { return i > 0 },
			expected:  true,
		},
		{
			name:      "None positive",
			input:     []int{-1, -2, -3, -4, -5},
			predicate: func(i int) bool { return i > 0 },
			expected:  false,
		},
		{
			name:      "Some even",
			input:     []int{1, 3, 4, 7},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  true,
		},
		{
			name:      "None even",
			input:     []int{1, 3, 5, 7},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Any(tt.input, tt.predicate)
			if result != tt.expected {
				t.Errorf("Any() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		value    string
		expected bool
	}{
		{
			name:     "Empty slice",
			input:    []string{},
			value:    "a",
			expected: false,
		},
		{
			name:     "Contains value",
			input:    []string{"a", "b", "c", "d"},
			value:    "c",
			expected: true,
		},
		{
			name:     "Does not contain value",
			input:    []string{"a", "b", "c", "d"},
			value:    "z",
			expected: false,
		},
		{
			name:     "Contains empty string",
			input:    []string{"a", "b", "", "d"},
			value:    "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.input, tt.value)
			if result != tt.expected {
				t.Errorf("Contains() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Test with int slice
	intTests := []struct {
		name     string
		input    []int
		value    int
		expected bool
	}{
		{
			name:     "Contains int",
			input:    []int{1, 2, 3, 4, 5},
			value:    3,
			expected: true,
		},
		{
			name:     "Does not contain int",
			input:    []int{1, 2, 3, 4, 5},
			value:    6,
			expected: false,
		},
	}

	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.input, tt.value)
			if result != tt.expected {
				t.Errorf("Contains() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		predicate func(string) bool
		expected  Option[string]
		isSome    bool
	}{
		{
			name:      "Empty slice",
			input:     []string{},
			predicate: func(s string) bool { return strings.Contains(s, "a") },
			expected:  None[string](),
			isSome:    false,
		},
		{
			name:      "Find element",
			input:     []string{"dog", "cat", "rabbit", "hamster"},
			predicate: func(s string) bool { return strings.Contains(s, "a") },
			expected:  Some("cat"), // First element with 'a'
			isSome:    true,
		},
		{
			name:      "No matching element",
			input:     []string{"dog", "fox", "wolf"},
			predicate: func(s string) bool { return strings.Contains(s, "z") },
			expected:  None[string](),
			isSome:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Find(tt.input, tt.predicate)

			// Check IsSome/IsNone
			if result.IsSome() != tt.isSome {
				t.Errorf("Find().IsSome() = %v, want %v", result.IsSome(), tt.isSome)
			}

			// If expected to find something, check the value
			if tt.isSome {
				if result.Unwrap() != tt.expected.Unwrap() {
					t.Errorf("Find() = %v, want %v", result.Unwrap(), tt.expected.Unwrap())
				}
			}
		})
	}
}

func TestZip(t *testing.T) {
	tests := []struct {
		name     string
		input1   []string
		input2   []int
		expected []Pair[string, int]
	}{
		{
			name:     "Empty slices",
			input1:   []string{},
			input2:   []int{},
			expected: []Pair[string, int]{},
		},
		{
			name:   "Equal length slices",
			input1: []string{"a", "b", "c"},
			input2: []int{1, 2, 3},
			expected: []Pair[string, int]{
				{First: "a", Second: 1},
				{First: "b", Second: 2},
				{First: "c", Second: 3},
			},
		},
		{
			name:   "First slice longer",
			input1: []string{"a", "b", "c", "d"},
			input2: []int{1, 2, 3},
			expected: []Pair[string, int]{
				{First: "a", Second: 1},
				{First: "b", Second: 2},
				{First: "c", Second: 3},
			},
		},
		{
			name:   "Second slice longer",
			input1: []string{"a", "b"},
			input2: []int{1, 2, 3, 4},
			expected: []Pair[string, int]{
				{First: "a", Second: 1},
				{First: "b", Second: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Zip(tt.input1, tt.input2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Zip() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMakePair(t *testing.T) {
	tests := []struct {
		name     string
		first    string
		second   int
		expected Pair[string, int]
	}{
		{
			name:     "String-Int pair",
			first:    "hello",
			second:   42,
			expected: Pair[string, int]{First: "hello", Second: 42},
		},
		{
			name:     "Empty-Zero pair",
			first:    "",
			second:   0,
			expected: Pair[string, int]{First: "", Second: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MakePair(tt.first, tt.second)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MakePair() = %v, want %v", result, tt.expected)
			}
		})
	}
}
