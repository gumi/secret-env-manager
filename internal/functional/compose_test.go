package functional

import (
	"strconv"
	"testing"
)

func TestPipe(t *testing.T) {
	addOne := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }

	tests := []struct {
		name     string
		f1       func(int) int
		f2       func(int) int
		input    int
		expected int
	}{
		{
			name:     "Add one then double",
			f1:       addOne,
			f2:       double,
			input:    3,
			expected: 8, // (3+1)*2 = 8
		},
		{
			name:     "Double then add one",
			f1:       double,
			f2:       addOne,
			input:    3,
			expected: 7, // 3*2+1 = 7
		},
		{
			name:     "Identity composition",
			f1:       func(x int) int { return x },
			f2:       func(x int) int { return x },
			input:    5,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composed := Pipe(tt.f1, tt.f2)
			result := composed(tt.input)
			if result != tt.expected {
				t.Errorf("Pipe() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPipeMany(t *testing.T) {
	addOne := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }
	square := func(x int) int { return x * x }

	tests := []struct {
		name     string
		funcs    []func(int) int
		input    int
		expected int
	}{
		{
			name:     "Empty pipe",
			funcs:    []func(int) int{},
			input:    5,
			expected: 5, // No functions means identity
		},
		{
			name:     "Single function",
			funcs:    []func(int) int{double},
			input:    3,
			expected: 6,
		},
		{
			name:     "Multiple functions (add, double, square)",
			funcs:    []func(int) int{addOne, double, square},
			input:    3,
			expected: 64, // ((3+1)*2)^2 = 64
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composed := PipeMany(tt.funcs...)
			result := composed(tt.input)
			if result != tt.expected {
				t.Errorf("PipeMany() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCompose(t *testing.T) {
	addOne := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }

	tests := []struct {
		name     string
		f1       func(int) int
		f2       func(int) int
		input    int
		expected int
	}{
		{
			name:     "Add one after double",
			f1:       addOne,
			f2:       double,
			input:    3,
			expected: 7, // 3*2+1 = 7
		},
		{
			name:     "Double after add one",
			f1:       double,
			f2:       addOne,
			input:    3,
			expected: 8, // (3+1)*2 = 8
		},
		{
			name:     "Identity composition",
			f1:       func(x int) int { return x },
			f2:       func(x int) int { return x },
			input:    5,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composed := Compose(tt.f1, tt.f2)
			result := composed(tt.input)
			if result != tt.expected {
				t.Errorf("Compose() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestComposeMany(t *testing.T) {
	addOne := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }
	square := func(x int) int { return x * x }

	tests := []struct {
		name     string
		funcs    []func(int) int
		input    int
		expected int
	}{
		{
			name:     "Empty compose",
			funcs:    []func(int) int{},
			input:    5,
			expected: 5, // No functions means identity
		},
		{
			name:     "Single function",
			funcs:    []func(int) int{double},
			input:    3,
			expected: 6,
		},
		{
			name:     "Right to left: square, then double, then add one",
			funcs:    []func(int) int{addOne, double, square},
			input:    3,
			expected: 19, // 3^2*2+1 = 19
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composed := ComposeMany(tt.funcs...)
			result := composed(tt.input)
			if result != tt.expected {
				t.Errorf("ComposeMany() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		function func(int) string
		expected string
	}{
		{
			name:     "Int to string",
			value:    42,
			function: strconv.Itoa,
			expected: "42",
		},
		{
			name:     "Custom function",
			value:    5,
			function: func(x int) string { return "Number: " + strconv.Itoa(x*x) },
			expected: "Number: 25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Apply(tt.value, tt.function)
			if result != tt.expected {
				t.Errorf("Apply() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestApplyAll(t *testing.T) {
	addOne := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }
	square := func(x int) int { return x * x }

	tests := []struct {
		name     string
		value    int
		funcs    []func(int) int
		expected int
	}{
		{
			name:     "Empty apply",
			value:    5,
			funcs:    []func(int) int{},
			expected: 5, // No functions means identity
		},
		{
			name:     "Single function",
			value:    3,
			funcs:    []func(int) int{double},
			expected: 6,
		},
		{
			name:     "Multiple functions (add, double, square)",
			value:    3,
			funcs:    []func(int) int{addOne, double, square},
			expected: 64, // ((3+1)*2)^2 = 64
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyAll(tt.value, tt.funcs...)
			if result != tt.expected {
				t.Errorf("ApplyAll() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIdentity(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected interface{}
	}{
		{
			name:     "Int identity",
			value:    42,
			expected: 42,
		},
		{
			name:     "String identity",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "Bool identity",
			value:    true,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Identity(tt.value)
			if result != tt.expected {
				t.Errorf("Identity() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConstant(t *testing.T) {
	tests := []struct {
		name       string
		constValue string
		input      int
		expected   string
	}{
		{
			name:       "Constant string",
			constValue: "hello",
			input:      42,
			expected:   "hello",
		},
		{
			name:       "Empty string",
			constValue: "",
			input:      99,
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constFunc := Constant[string, int](tt.constValue)
			result := constFunc(tt.input)
			if result != tt.expected {
				t.Errorf("Constant() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMemoize(t *testing.T) {
	// Counter to track function calls
	callCount := 0

	// Expensive function that we want to memoize
	expensiveFunc := func(key string) int {
		callCount++
		return len(key)
	}

	// Create memoized version
	memoized := Memoize(expensiveFunc)

	// Test cases
	tests := []struct {
		name          string
		input         string
		expected      int
		expectedCalls int
	}{
		{
			name:          "First call",
			input:         "hello",
			expected:      5,
			expectedCalls: 1,
		},
		{
			name:          "Repeated call with same value",
			input:         "hello",
			expected:      5,
			expectedCalls: 1, // No additional call
		},
		{
			name:          "Different value",
			input:         "world",
			expected:      5,
			expectedCalls: 2, // One more call
		},
		{
			name:          "Repeat second value",
			input:         "world",
			expected:      5,
			expectedCalls: 2, // No additional call
		},
		{
			name:          "Empty string",
			input:         "",
			expected:      0,
			expectedCalls: 3, // One more call
		},
	}

	callCount = 0 // Reset call counter

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := memoized(tt.input)
			if result != tt.expected {
				t.Errorf("Memoized function result = %v, want %v", result, tt.expected)
			}

			if callCount != tt.expectedCalls {
				t.Errorf("Call count = %v, want %v", callCount, tt.expectedCalls)
			}
		})
	}
}
