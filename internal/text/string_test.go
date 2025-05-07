package text

import (
	"reflect"
	"strings"
	"testing"
)

func TestIdentity(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Empty string", "", ""},
		{"Regular string", "hello", "hello"},
		{"String with spaces", "hello world", "hello world"},
		{"String with special chars", "!@#$%^&*()", "!@#$%^&*()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Identity(tt.input); got != tt.want {
				t.Errorf("Identity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanControlCharsResult(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"Empty string", "", "", false},
		{"No control chars", "hello world", "hello world", false},
		{"With spaces", "hello world", "hello world", false},
		{"With tabs", "hello\tworld", "helloworld", false},
		{"With newlines", "hello\nworld", "helloworld", false},
		{"With carriage returns", "hello\rworld", "helloworld", false},
		{"With control chars", "hello\u0000world", "helloworld", false},
		{"Only control chars", "\u0000\t\n\r", "", false},
		{"Mixed control chars", "hello\u0000\t\n\rworld", "helloworld", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanControlCharsResult(tt.input)
			if (result.IsFailure()) != tt.wantErr {
				t.Errorf("CleanControlCharsResult() error = %v, wantErr %v", result.GetError(), tt.wantErr)
				return
			}
			if !tt.wantErr && result.Unwrap() != tt.want {
				t.Errorf("CleanControlCharsResult() = %v, want %v", result.Unwrap(), tt.want)
			}
		})
	}
}

func TestCleanControlChars(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Empty string", "", ""},
		{"No control chars", "hello world", "hello world"},
		{"With control chars", "hello\u0000world", "helloworld"},
		{"Mixed control chars", "hello\u0000\t\n\rworld", "helloworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanControlChars(tt.input); got != tt.want {
				t.Errorf("CleanControlChars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinWithSeparatorResult(t *testing.T) {
	tests := []struct {
		name      string
		separator string
		parts     []string
		want      string
		wantErr   bool
	}{
		{"Empty parts", ",", []string{}, "", false},
		{"Single part", ",", []string{"hello"}, "hello", false},
		{"Multiple parts", ",", []string{"hello", "world"}, "hello,world", false},
		{"With empty parts", ",", []string{"hello", "", "world"}, "hello,world", false},
		{"All empty parts", ",", []string{"", "", ""}, "", false},
		{"Different separator", " | ", []string{"hello", "world"}, "hello | world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinWithSeparatorResult(tt.separator, tt.parts...)
			if (result.IsFailure()) != tt.wantErr {
				t.Errorf("JoinWithSeparatorResult() error = %v, wantErr %v", result.GetError(), tt.wantErr)
				return
			}
			if !tt.wantErr && result.Unwrap() != tt.want {
				t.Errorf("JoinWithSeparatorResult() = %v, want %v", result.Unwrap(), tt.want)
			}
		})
	}
}

func TestJoinWithSeparator(t *testing.T) {
	tests := []struct {
		name      string
		separator string
		parts     []string
		want      string
	}{
		{"Empty parts", ",", []string{}, ""},
		{"Single part", ",", []string{"hello"}, "hello"},
		{"Multiple parts", ",", []string{"hello", "world"}, "hello,world"},
		{"With empty parts", ",", []string{"hello", "", "world"}, "hello,world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinWithSeparator(tt.separator, tt.parts...); got != tt.want {
				t.Errorf("JoinWithSeparator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitAndTrimResult(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		separator string
		want      []string
		wantErr   bool
	}{
		{"Empty string", "", ",", nil, true},
		{"Simple split", "hello,world", ",", []string{"hello", "world"}, false},
		{"With spaces", " hello , world ", ",", []string{"hello", "world"}, false},
		{"Multiple separators", "hello,world,again", ",", []string{"hello", "world", "again"}, false},
		{"Single item", "hello", ",", []string{"hello"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitAndTrimResult(tt.s, tt.separator)
			if (result.IsFailure()) != tt.wantErr {
				t.Errorf("SplitAndTrimResult() error = %v, wantErr %v", result.GetError(), tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result.Unwrap(), tt.want) {
				t.Errorf("SplitAndTrimResult() = %v, want %v", result.Unwrap(), tt.want)
			}
		})
	}
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		separator string
		want      []string
	}{
		{"Empty string", "", ",", []string{""}},
		{"Simple split", "hello,world", ",", []string{"hello", "world"}},
		{"With spaces", " hello , world ", ",", []string{"hello", "world"}},
		{"Multiple separators", "hello,world,again", ",", []string{"hello", "world", "again"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitAndTrim(tt.s, tt.separator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitAndTrim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapStrings(t *testing.T) {
	upperFunc := func(s string) string { return strings.ToUpper(s) }
	appendFunc := func(s string) string { return s + "!" }

	tests := []struct {
		name        string
		strs        []string
		transformer StringTransformer
		want        []string
	}{
		{"Empty slice", []string{}, upperFunc, []string{}},
		{"Upper transform", []string{"hello", "world"}, upperFunc, []string{"HELLO", "WORLD"}},
		{"Append transform", []string{"hello", "world"}, appendFunc, []string{"hello!", "world!"}},
		{"Mixed case", []string{"Hello", "World"}, upperFunc, []string{"HELLO", "WORLD"}},
		{"With empty string", []string{"hello", "", "world"}, upperFunc, []string{"HELLO", "", "WORLD"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapStrings(tt.strs, tt.transformer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapStringsWithIndex(t *testing.T) {
	indexedFunc := func(i int, s string) string { return s + string(rune('0'+i)) }
	emptySlice := []string{}

	tests := []struct {
		name        string
		strs        []string
		transformer func(int, string) string
		want        []string
	}{
		{"Empty slice", emptySlice, indexedFunc, emptySlice},
		{"Index transform", []string{"item", "item"}, indexedFunc, []string{"item0", "item1"}},
		{"Mixed values", []string{"a", "b", "c"}, indexedFunc, []string{"a0", "b1", "c2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapStringsWithIndex(tt.strs, tt.transformer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapStringsWithIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterStrings(t *testing.T) {
	notEmptyFunc := func(s string) bool { return s != "" }
	containsAFunc := func(s string) bool { return strings.Contains(s, "a") }

	tests := []struct {
		name      string
		strs      []string
		predicate StringPredicate
		want      []string
	}{
		{"Empty slice", []string{}, notEmptyFunc, []string{}},
		{"Filter empty", []string{"hello", "", "world", ""}, notEmptyFunc, []string{"hello", "world"}},
		{"Filter by content", []string{"apple", "banana", "cherry"}, containsAFunc, []string{"apple", "banana"}},
		{"No matches", []string{"dog", "egg"}, containsAFunc, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterStrings(tt.strs, tt.predicate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty string", "", false},
		{"Non-empty string", "hello", true},
		{"Space only", " ", true},
		{"Tab only", "\t", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotEmpty(tt.input); got != tt.want {
				t.Errorf("IsNotEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChain(t *testing.T) {
	upperFunc := func(s string) string { return strings.ToUpper(s) }
	trimFunc := func(s string) string { return strings.TrimSpace(s) }

	tests := []struct {
		name         string
		input        string
		transformers []StringTransformer
		want         string
	}{
		{"Empty string", "", []StringTransformer{upperFunc, trimFunc}, ""},
		{"Upper then trim", " hello ", []StringTransformer{upperFunc, trimFunc}, "HELLO"},
		{"Trim then upper", " hello ", []StringTransformer{trimFunc, upperFunc}, "HELLO"},
		{"Single transform", "hello", []StringTransformer{upperFunc}, "HELLO"},
		{"No transforms", "hello", []StringTransformer{}, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := Chain(tt.transformers...)
			if got := transformer(tt.input); got != tt.want {
				t.Errorf("Chain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComposePredicate(t *testing.T) {
	isLongerThan3 := func(s string) bool { return len(s) > 3 }
	hasVowel := func(s string) bool { return strings.ContainsAny(s, "aeiouAEIOU") }

	tests := []struct {
		name       string
		input      string
		predicates []StringPredicate
		want       bool
	}{
		{"Empty predicates", "word", nil, true},
		{"Single true predicate", "word", []StringPredicate{isLongerThan3}, true},
		{"Single false predicate", "hi", []StringPredicate{isLongerThan3}, false},
		{"Both true predicates", "word", []StringPredicate{isLongerThan3, hasVowel}, true},
		{"One false predicate", "xyz", []StringPredicate{isLongerThan3, hasVowel}, false},
		{"Both false predicates", "xy", []StringPredicate{isLongerThan3, hasVowel}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			predicate := ComposePredicate(tt.predicates...)
			if got := predicate(tt.input); got != tt.want {
				t.Errorf("ComposePredicate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComposePredicateOr(t *testing.T) {
	isLongerThan3 := func(s string) bool { return len(s) > 3 }
	hasVowel := func(s string) bool { return strings.ContainsAny(s, "aeiouAEIOU") }

	tests := []struct {
		name       string
		input      string
		predicates []StringPredicate
		want       bool
	}{
		{"Empty predicates", "word", nil, false},
		{"Single true predicate", "word", []StringPredicate{isLongerThan3}, true},
		{"Single false predicate", "hi", []StringPredicate{isLongerThan3}, false},
		{"Both true predicates", "word", []StringPredicate{isLongerThan3, hasVowel}, true},
		{"One true predicate", "xyzw", []StringPredicate{isLongerThan3, hasVowel}, true},
		{"Both false predicates", "xy", []StringPredicate{isLongerThan3, hasVowel}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			predicate := ComposePredicateOr(tt.predicates...)
			if got := predicate(tt.input); got != tt.want {
				t.Errorf("ComposePredicateOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitOption(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		separator string
		wantSome  bool
		wantVal   []string
	}{
		{"Empty string", "", ",", false, nil},
		{"Simple split", "hello,world", ",", true, []string{"hello", "world"}},
		{"Single value", "hello", ",", true, []string{"hello"}},
		{"Multiple separators", "a,b,c", ",", true, []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := SplitOption(tt.s, tt.separator)
			if opt.IsSome() != tt.wantSome {
				t.Errorf("SplitOption() IsSome = %v, want %v", opt.IsSome(), tt.wantSome)
				return
			}
			if tt.wantSome && !reflect.DeepEqual(opt.Unwrap(), tt.wantVal) {
				t.Errorf("SplitOption() value = %v, want %v", opt.Unwrap(), tt.wantVal)
			}
		})
	}
}

func TestSplitResult(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		separator string
		wantErr   bool
		wantVal   []string
	}{
		{"Empty string", "", ",", true, nil},
		{"Simple split", "hello,world", ",", false, []string{"hello", "world"}},
		{"Single value", "hello", ",", false, []string{"hello"}},
		{"Multiple separators", "a,b,c", ",", false, []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitResult(tt.s, tt.separator)
			if result.IsFailure() != tt.wantErr {
				t.Errorf("SplitResult() error = %v, wantErr %v", result.GetError(), tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result.Unwrap(), tt.wantVal) {
				t.Errorf("SplitResult() value = %v, want %v", result.Unwrap(), tt.wantVal)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		cutset   string
		expected string
	}{
		{
			name:     "Empty string",
			str:      "",
			cutset:   " ",
			expected: "",
		},
		{
			name:     "No characters to trim",
			str:      "Hello World",
			cutset:   "_",
			expected: "Hello World",
		},
		{
			name:     "Trim spaces",
			str:      "  Hello World  ",
			cutset:   " ",
			expected: "Hello World",
		},
		{
			name:     "Trim multiple characters",
			str:      "...Hello World...",
			cutset:   ".",
			expected: "Hello World",
		},
		{
			name:     "Trim mixed characters",
			str:      "...---Hello World...---",
			cutset:   ".-",
			expected: "Hello World",
		},
		{
			name:     "Only characters to trim",
			str:      "   ",
			cutset:   " ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Trim(tt.str, tt.cutset)
			if result != tt.expected {
				t.Errorf("Trim(%q, %q) = %q, want %q",
					tt.str, tt.cutset, result, tt.expected)
			}
		})
	}
}

func TestTrimSpace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Empty string", "", ""},
		{"No whitespace", "hello", "hello"},
		{"Leading whitespace", "  hello", "hello"},
		{"Trailing whitespace", "hello  ", "hello"},
		{"Both sides", "  hello  ", "hello"},
		{"Internal whitespace", "hello world", "hello world"},
		{"Only whitespace", "   ", ""},
		{"Tabs and newlines", "\t\nhello\r\n", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimSpace(tt.input); got != tt.want {
				t.Errorf("TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		old      string
		new      string
		n        int
		expected string
	}{
		{
			name:     "Empty string",
			str:      "",
			old:      "a",
			new:      "b",
			n:        -1,
			expected: "",
		},
		{
			name:     "No occurrences",
			str:      "Hello World",
			old:      "z",
			new:      "x",
			n:        -1,
			expected: "Hello World",
		},
		{
			name:     "Replace all occurrences",
			str:      "Hello World, Hello Universe",
			old:      "Hello",
			new:      "Hi",
			n:        -1,
			expected: "Hi World, Hi Universe",
		},
		{
			name:     "Replace limited occurrences",
			str:      "Hello World, Hello Universe, Hello Galaxy",
			old:      "Hello",
			new:      "Hi",
			n:        2,
			expected: "Hi World, Hi Universe, Hello Galaxy",
		},
		{
			name:     "Replace with empty string",
			str:      "Hello World",
			old:      "o",
			new:      "",
			n:        -1,
			expected: "Hell Wrld",
		},
		{
			name:     "Replace empty string (zero n)",
			str:      "Hello",
			old:      "",
			new:      "x",
			n:        0,
			expected: "Hello",
		},
		{
			name:     "Replace empty string (n > 0)",
			str:      "Hello",
			old:      "",
			new:      "x",
			n:        -1,
			expected: "xHxexlxlxox",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Replace(tt.str, tt.old, tt.new, tt.n)
			if result != tt.expected {
				t.Errorf("Replace(%q, %q, %q, %d) = %q, want %q",
					tt.str, tt.old, tt.new, tt.n, result, tt.expected)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Empty string", "", ""},
		{"Already lowercase", "hello", "hello"},
		{"All uppercase", "HELLO", "hello"},
		{"Mixed case", "HeLLo", "hello"},
		{"With numbers", "Hello123", "hello123"},
		{"With symbols", "HELLO!", "hello!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToLower(tt.input); got != tt.want {
				t.Errorf("ToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected string
	}{
		{
			name:     "Empty string",
			str:      "",
			expected: "",
		},
		{
			name:     "All lowercase",
			str:      "hello world",
			expected: "HELLO WORLD",
		},
		{
			name:     "All uppercase",
			str:      "HELLO WORLD",
			expected: "HELLO WORLD",
		},
		{
			name:     "Mixed case",
			str:      "Hello World",
			expected: "HELLO WORLD",
		},
		{
			name:     "With numbers and symbols",
			str:      "Hello123!@#",
			expected: "HELLO123!@#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToUpper(tt.str)
			if result != tt.expected {
				t.Errorf("ToUpper(%q) = %q, want %q",
					tt.str, result, tt.expected)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		substr   string
		expected int
	}{
		{
			name:     "Empty string",
			str:      "",
			substr:   "a",
			expected: -1,
		},
		{
			name:     "Empty substring",
			str:      "Hello",
			substr:   "",
			expected: 0,
		},
		{
			name:     "Substring at beginning",
			str:      "Hello World",
			substr:   "Hello",
			expected: 0,
		},
		{
			name:     "Substring in middle",
			str:      "Hello World",
			substr:   "o W",
			expected: 4,
		},
		{
			name:     "Substring at end",
			str:      "Hello World",
			substr:   "World",
			expected: 6,
		},
		{
			name:     "Substring not found",
			str:      "Hello World",
			substr:   "Universe",
			expected: -1,
		},
		{
			name:     "Case sensitive search",
			str:      "Hello World",
			substr:   "world",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Index(tt.str, tt.substr)
			if result != tt.expected {
				t.Errorf("Index(%q, %q) = %d, want %d",
					tt.str, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestLastIndex(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		substr   string
		expected int
	}{
		{
			name:     "Empty string",
			str:      "",
			substr:   "a",
			expected: -1,
		},
		{
			name:     "Empty substring",
			str:      "Hello",
			substr:   "",
			expected: 5,
		},
		{
			name:     "Single occurrence",
			str:      "Hello World",
			substr:   "World",
			expected: 6,
		},
		{
			name:     "Multiple occurrences",
			str:      "Hello World Hello Universe",
			substr:   "Hello",
			expected: 12,
		},
		{
			name:     "Substring not found",
			str:      "Hello World",
			substr:   "Universe",
			expected: -1,
		},
		{
			name:     "Case sensitive search",
			str:      "Hello World",
			substr:   "world",
			expected: -1,
		},
		{
			name:     "Overlapping occurrences",
			str:      "ababababa",
			substr:   "aba",
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LastIndex(tt.str, tt.substr)
			if result != tt.expected {
				t.Errorf("LastIndex(%q, %q) = %d, want %d",
					tt.str, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name     string
		elems    []string
		sep      string
		expected string
	}{
		{
			name:     "Empty slice",
			elems:    []string{},
			sep:      ",",
			expected: "",
		},
		{
			name:     "Single element",
			elems:    []string{"Hello"},
			sep:      ",",
			expected: "Hello",
		},
		{
			name:     "Multiple elements with comma",
			elems:    []string{"Hello", "World", "Test"},
			sep:      ",",
			expected: "Hello,World,Test",
		},
		{
			name:     "Multiple elements with space",
			elems:    []string{"Hello", "World", "Test"},
			sep:      " ",
			expected: "Hello World Test",
		},
		{
			name:     "Empty separator",
			elems:    []string{"Hello", "World", "Test"},
			sep:      "",
			expected: "HelloWorldTest",
		},
		{
			name:     "With empty elements",
			elems:    []string{"Hello", "", "Test", ""},
			sep:      ",",
			expected: "Hello,,Test,",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Join(tt.elems, tt.sep)
			if result != tt.expected {
				t.Errorf("Join(%v, %q) = %q, want %q",
					tt.elems, tt.sep, result, tt.expected)
			}
		})
	}
}

func TestStartsWith(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		prefix   string
		expected bool
	}{
		{
			name:     "Empty string and prefix",
			str:      "",
			prefix:   "",
			expected: true,
		},
		{
			name:     "Empty string with prefix",
			str:      "",
			prefix:   "prefix",
			expected: false,
		},
		{
			name:     "String starts with prefix",
			str:      "Hello World",
			prefix:   "Hello",
			expected: true,
		},
		{
			name:     "String does not start with prefix",
			str:      "Hello World",
			prefix:   "World",
			expected: false,
		},
		{
			name:     "Case sensitive match",
			str:      "Hello World",
			prefix:   "hello",
			expected: false,
		},
		{
			name:     "Prefix longer than string",
			str:      "Hi",
			prefix:   "Hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StartsWith(tt.str, tt.prefix)
			if result != tt.expected {
				t.Errorf("StartsWith(%q, %q) = %v, want %v",
					tt.str, tt.prefix, result, tt.expected)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		suffix   string
		expected bool
	}{
		{
			name:     "Empty string and suffix",
			str:      "",
			suffix:   "",
			expected: true,
		},
		{
			name:     "Empty string with suffix",
			str:      "",
			suffix:   "suffix",
			expected: false,
		},
		{
			name:     "String ends with suffix",
			str:      "Hello World",
			suffix:   "World",
			expected: true,
		},
		{
			name:     "String does not end with suffix",
			str:      "Hello World",
			suffix:   "Hello",
			expected: false,
		},
		{
			name:     "Case sensitive match",
			str:      "Hello World",
			suffix:   "world",
			expected: false,
		},
		{
			name:     "Suffix longer than string",
			str:      "Hi",
			suffix:   "Hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EndsWith(tt.str, tt.suffix)
			if result != tt.expected {
				t.Errorf("EndsWith(%q, %q) = %v, want %v",
					tt.str, tt.suffix, result, tt.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		substr string
		input  string
		want   bool
	}{
		{"Empty substring", "", "hello", true},
		{"Empty string", "a", "", false},
		{"Both empty", "", "", true},
		{"Contains", "lo", "hello", true},
		{"Case sensitive", "LO", "hello", false},
		{"Start of string", "he", "hello", true},
		{"End of string", "lo", "hello", true},
		{"Not contains", "xyz", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			predicate := Contains(tt.substr)
			if got := predicate(tt.input); got != tt.want {
				t.Errorf("Contains(%q)(%q) = %v, want %v", tt.substr, tt.input, got, tt.want)
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		input  string
		want   bool
	}{
		{"Empty prefix", "", "hello", true},
		{"Empty string", "a", "", false},
		{"Both empty", "", "", true},
		{"Has prefix", "he", "hello", true},
		{"Case sensitive", "HE", "hello", false},
		{"Middle of string", "ll", "hello", false},
		{"End of string", "lo", "hello", false},
		{"Exact match", "hello", "hello", true},
		{"Too long", "hello world", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			predicate := HasPrefix(tt.prefix)
			if got := predicate(tt.input); got != tt.want {
				t.Errorf("HasPrefix(%q)(%q) = %v, want %v", tt.prefix, tt.input, got, tt.want)
			}
		})
	}
}

func TestHasSuffix(t *testing.T) {
	tests := []struct {
		name   string
		suffix string
		input  string
		want   bool
	}{
		{"Empty suffix", "", "hello", true},
		{"Empty string", "a", "", false},
		{"Both empty", "", "", true},
		{"Has suffix", "lo", "hello", true},
		{"Case sensitive", "LO", "hello", false},
		{"Start of string", "he", "hello", false},
		{"Middle of string", "ll", "hello", false},
		{"Exact match", "hello", "hello", true},
		{"Too long", "world hello", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			predicate := HasSuffix(tt.suffix)
			if got := predicate(tt.input); got != tt.want {
				t.Errorf("HasSuffix(%q)(%q) = %v, want %v", tt.suffix, tt.input, got, tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		sep      string
		expected []string
	}{
		{
			name:     "Empty string",
			str:      "",
			sep:      ",",
			expected: []string{""},
		},
		{
			name:     "No separator in string",
			str:      "Hello World",
			sep:      ",",
			expected: []string{"Hello World"},
		},
		{
			name:     "Simple split by space",
			str:      "Hello World",
			sep:      " ",
			expected: []string{"Hello", "World"},
		},
		{
			name:     "Multiple separators",
			str:      "a,b,c,d",
			sep:      ",",
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "Trailing separator",
			str:      "a,b,c,",
			sep:      ",",
			expected: []string{"a", "b", "c", ""},
		},
		{
			name:     "Leading separator",
			str:      ",a,b,c",
			sep:      ",",
			expected: []string{"", "a", "b", "c"},
		},
		{
			name:     "Empty separator",
			str:      "abc",
			sep:      "",
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Split(tt.str, tt.sep)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Split(%q, %q) = %v, want %v",
					tt.str, tt.sep, result, tt.expected)
			}
		})
	}
}

func TestJoinNonEmptyResult(t *testing.T) {
	tests := []struct {
		name      string
		separator string
		strs      []string
		want      string
		wantErr   bool
	}{
		{"Empty slice", ",", []string{}, "", false},
		{"Single item", ",", []string{"hello"}, "hello", false},
		{"Multiple items", ",", []string{"hello", "world"}, "hello,world", false},
		{"With empty strings", ",", []string{"hello", "", "world"}, "hello,world", false},
		{"All empty strings", ",", []string{"", "", ""}, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinNonEmptyResult(tt.separator, tt.strs...)
			if (result.IsFailure()) != tt.wantErr {
				t.Errorf("JoinNonEmptyResult() error = %v, wantErr %v", result.GetError(), tt.wantErr)
				return
			}
			if !tt.wantErr && result.Unwrap() != tt.want {
				t.Errorf("JoinNonEmptyResult() = %v, want %v", result.Unwrap(), tt.want)
			}
		})
	}
}

func TestJoinNonEmpty(t *testing.T) {
	tests := []struct {
		name      string
		separator string
		strs      []string
		want      string
	}{
		{"Empty slice", ",", []string{}, ""},
		{"Single item", ",", []string{"hello"}, "hello"},
		{"Multiple items", ",", []string{"hello", "world"}, "hello,world"},
		{"With empty strings", ",", []string{"hello", "", "world"}, "hello,world"},
		{"All empty strings", ",", []string{"", "", ""}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinNonEmpty(tt.separator, tt.strs...); got != tt.want {
				t.Errorf("JoinNonEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
