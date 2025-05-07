package text

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestIsJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty string", "", false},
		{"Valid JSON object", `{"name":"John","age":30}`, true},
		{"Valid JSON array", `[1,2,3]`, true},
		{"Valid JSON string", `"hello"`, true},
		{"Valid JSON number", `123`, true},
		{"Valid JSON boolean", `true`, true},
		{"Valid JSON null", `null`, true},
		{"Invalid JSON", `{"name":"John"`, false},
		{"Plain text", "hello world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSON(tt.input); got != tt.want {
				t.Errorf("IsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsJSONObject(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty string", "", false},
		{"Valid JSON object", `{"name":"John","age":30}`, true},
		{"Valid JSON object with whitespace", ` { "name" : "John" } `, true},
		{"JSON array", `[1,2,3]`, false},
		{"JSON string", `"hello"`, false},
		{"Invalid JSON object", `{"name":"John"`, false},
		{"Non-JSON", "hello world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSONObject(tt.input); got != tt.want {
				t.Errorf("IsJSONObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsJSONObjectResult(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Valid JSON object", `{"name":"John"}`, true},
		{"Not JSON object", `[1,2,3]`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsJSONObjectResult(tt.input)
			if result.IsFailure() {
				t.Errorf("IsJSONObjectResult() failed: %v", result.GetError())
			}
			if got := result.Unwrap(); got != tt.want {
				t.Errorf("IsJSONObjectResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsJSONArray(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty string", "", false},
		{"Valid JSON array", `[1,2,3]`, true},
		{"Valid JSON array with whitespace", ` [ 1, 2, 3 ] `, true},
		{"JSON object", `{"name":"John"}`, false},
		{"JSON string", `"hello"`, false},
		{"Invalid JSON array", `[1,2,`, false},
		{"Non-JSON", "hello world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSONArray(tt.input); got != tt.want {
				t.Errorf("IsJSONArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsJSONData(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty string", "", false},
		{"Valid JSON object", `{"name":"John"}`, true},
		{"Valid JSON array", `[1,2,3]`, true},
		{"JSON string", `"hello"`, false},
		{"JSON number", `123`, false},
		{"Invalid JSON", `{"name":"John"`, false},
		{"Non-JSON", "hello world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSONData(tt.input); got != tt.want {
				t.Errorf("IsJSONData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeJSONEncode(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{"String", "hello", `"hello"`},
		{"Number", 123, `123`},
		{"Boolean", true, `true`},
		{"Null", nil, `null`},
		{"Array", []int{1, 2, 3}, `[1,2,3]`},
		{"Object", map[string]string{"name": "John"}, `{"name":"John"}`},
		{"Nested object", map[string]interface{}{"person": map[string]string{"name": "John"}}, `{"person":{"name":"John"}}`},
		// Cannot test invalid values that would fail JSON marshaling here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafeJSONEncode(tt.input)
			// Normalize the strings to handle variations in formatting
			var gotObj, wantObj interface{}
			_ = json.Unmarshal([]byte(got), &gotObj)
			_ = json.Unmarshal([]byte(tt.want), &wantObj)
			if !reflect.DeepEqual(gotObj, wantObj) {
				t.Errorf("SafeJSONEncode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{"Empty string", "", nil, true},
		{"Valid JSON object", `{"name":"John"}`, map[string]interface{}{"name": "John"}, false},
		{"Valid JSON array", `[1,2,3]`, []interface{}{1.0, 2.0, 3.0}, false},
		{"Valid JSON string", `"hello"`, "hello", false},
		{"Valid JSON number", `123`, 123.0, false},
		{"Valid JSON boolean", `true`, true, false},
		{"Valid JSON null", `null`, nil, false},
		{"Invalid JSON", `{"name":"John"`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJSONOption(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantSome bool
		want     interface{}
	}{
		{"Empty string", "", false, nil},
		{"Valid JSON object", `{"name":"John"}`, true, map[string]interface{}{"name": "John"}},
		{"Invalid JSON", `{"name":"John"`, false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := ParseJSONOption(tt.input)
			if option.IsSome() != tt.wantSome {
				t.Errorf("ParseJSONOption().IsSome() = %v, want %v", option.IsSome(), tt.wantSome)
				return
			}
			if tt.wantSome && !reflect.DeepEqual(option.Unwrap(), tt.want) {
				t.Errorf("ParseJSONOption().Unwrap() = %v, want %v", option.Unwrap(), tt.want)
			}
		})
	}
}

func TestParseJSONResult(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    interface{}
	}{
		{"Empty string", "", true, nil},
		{"Valid JSON object", `{"name":"John"}`, false, map[string]interface{}{"name": "John"}},
		{"Invalid JSON", `{"name":"John"`, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseJSONResult(tt.input)
			if result.IsFailure() != tt.wantErr {
				t.Errorf("ParseJSONResult().IsFailure() = %v, want %v", result.IsFailure(), tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result.Unwrap(), tt.want) {
				t.Errorf("ParseJSONResult().Unwrap() = %v, want %v", result.Unwrap(), tt.want)
			}
		})
	}
}

func TestParseJSONMapResult(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    map[string]interface{}
	}{
		{"Empty string", "", true, nil},
		{"Valid JSON object", `{"name":"John"}`, false, map[string]interface{}{"name": "John"}},
		{"Valid JSON array (not object)", `[1,2,3]`, true, nil},
		{"Invalid JSON", `{"name":"John"`, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseJSONMapResult(tt.input)
			if result.IsFailure() != tt.wantErr {
				t.Errorf("ParseJSONMapResult().IsFailure() = %v, want %v", result.IsFailure(), tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result.Unwrap(), tt.want) {
				t.Errorf("ParseJSONMapResult().Unwrap() = %v, want %v", result.Unwrap(), tt.want)
			}
		})
	}
}

func TestCreateStringMap(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		want  map[string]string
	}{
		{"Empty map", map[string]interface{}{}, map[string]string{}},
		{
			"Map with string values",
			map[string]interface{}{"name": "John", "city": "Tokyo"},
			map[string]string{"name": "John", "city": "Tokyo"},
		},
		{
			"Map with mixed types",
			map[string]interface{}{
				"name":  "John",
				"age":   30,
				"admin": true,
				"data":  []int{1, 2, 3},
			},
			map[string]string{
				"name":  "John",
				"age":   "30",
				"admin": "true",
				"data":  "[1,2,3]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateStringMap(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateStringMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigatePath(t *testing.T) {
	testData := map[string]interface{}{
		"person": map[string]interface{}{
			"name": "John",
			"age":  30,
			"address": map[string]interface{}{
				"city":  "Tokyo",
				"phone": "123-456-7890",
			},
		},
		"tags": []interface{}{"tag1", "tag2", "tag3"},
		"items": []interface{}{
			map[string]interface{}{"id": 1, "name": "item1"},
			map[string]interface{}{"id": 2, "name": "item2"},
		},
	}

	tests := []struct {
		name    string
		data    interface{}
		path    []string
		want    interface{}
		wantErr bool
	}{
		{"Empty path", testData, []string{}, testData, false},
		{"Single level", testData, []string{"person"}, testData["person"], false},
		{"Two levels", testData, []string{"person", "name"}, "John", false},
		{"Three levels", testData, []string{"person", "address", "city"}, "Tokyo", false},
		{"Array index", testData, []string{"tags", "1"}, "tag2", false},
		{"Array negative index", testData, []string{"tags", "-1"}, "tag3", false},
		{"Array last element", testData, []string{"tags", "-"}, "tag3", false},
		{"Nested array and object", testData, []string{"items", "0", "name"}, "item1", false},
		{"Key not found", testData, []string{"person", "email"}, nil, true},
		{"Index out of bounds", testData, []string{"tags", "5"}, nil, true},
		{"Negative index out of bounds", testData, []string{"tags", "-5"}, nil, true},
		{"Not an object", testData, []string{"person", "name", "first"}, nil, true},
		{"Not an array", testData, []string{"person", "0"}, nil, true},
		{"Invalid array index", testData, []string{"tags", "abc"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NavigatePath(tt.data, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NavigatePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NavigatePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigatePathResult(t *testing.T) {
	testData := map[string]interface{}{
		"person": map[string]interface{}{"name": "John"},
		"tags":   []interface{}{"tag1", "tag2"},
	}

	tests := []struct {
		name    string
		data    interface{}
		path    []string
		wantErr bool
		want    interface{}
	}{
		{"Valid path", testData, []string{"person", "name"}, false, "John"},
		{"Invalid path", testData, []string{"person", "email"}, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NavigatePathResult(tt.data, tt.path)
			if result.IsFailure() != tt.wantErr {
				t.Errorf("NavigatePathResult().IsFailure() = %v, want %v", result.IsFailure(), tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result.Unwrap(), tt.want) {
				t.Errorf("NavigatePathResult().Unwrap() = %v, want %v", result.Unwrap(), tt.want)
			}
		})
	}
}

func TestParseIndex(t *testing.T) {
	tests := []struct {
		name       string
		index      string
		arrayLen   int
		want       int
		wantErr    bool
		errorCheck func(error) bool
	}{
		{"Simple index", "2", 5, 2, false, nil},
		{"Zero index", "0", 5, 0, false, nil},
		{"Index at boundary", "4", 5, 4, false, nil},
		{"Simple negative", "-1", 5, 4, false, nil},
		{"Negative notation", "-", 5, 4, false, nil},
		{"Negative at boundary", "-5", 5, 0, false, nil},
		{"Index out of bounds", "5", 5, 0, true, nil},
		{"Negative out of bounds", "-6", 5, 0, true, nil},
		{"Invalid index", "abc", 5, 0, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIndex(tt.index, tt.arrayLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseIndex() = %v, want %v", got, tt.want)
			}
			if tt.wantErr && tt.errorCheck != nil && !tt.errorCheck(err) {
				t.Errorf("parseIndex() error = %v does not match expected error condition", err)
			}
		})
	}
}

func TestNavigateArrayIndex(t *testing.T) {
	array := []interface{}{"one", "two", "three"}
	notArray := map[string]string{"key": "value"}

	tests := []struct {
		name    string
		data    interface{}
		index   int
		wantErr bool
		want    interface{}
	}{
		{"Valid index", array, 1, false, "two"},
		{"Index zero", array, 0, false, "one"},
		{"Index at boundary", array, 2, false, "three"},
		{"Index out of bounds", array, 3, true, nil},
		{"Negative index", array, -1, true, nil},
		{"Not an array", notArray, 0, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NavigateArrayIndex(tt.data, tt.index)
			if result.IsFailure() != tt.wantErr {
				t.Errorf("NavigateArrayIndex().IsFailure() = %v, want %v", result.IsFailure(), tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result.Unwrap(), tt.want) {
				t.Errorf("NavigateArrayIndex().Unwrap() = %v, want %v", result.Unwrap(), tt.want)
			}
		})
	}
}

func TestNavigateObjectKey(t *testing.T) {
	obj := map[string]interface{}{"name": "John", "age": 30}
	notObj := []string{"one", "two"}

	tests := []struct {
		name    string
		data    interface{}
		key     string
		wantErr bool
		want    interface{}
	}{
		{"Valid key", obj, "name", false, "John"},
		{"Another valid key", obj, "age", false, 30},
		{"Key not found", obj, "email", true, nil},
		{"Not an object", notObj, "key", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NavigateObjectKey(tt.data, tt.key)
			if result.IsFailure() != tt.wantErr {
				t.Errorf("NavigateObjectKey().IsFailure() = %v, want %v", result.IsFailure(), tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result.Unwrap(), tt.want) {
				t.Errorf("NavigateObjectKey().Unwrap() = %v, want %v", result.Unwrap(), tt.want)
			}
		})
	}
}

func TestFormatAvailableKeys(t *testing.T) {
	tests := []struct {
		name  string
		data  interface{}
		wants []string
	}{
		{
			"Valid object",
			map[string]interface{}{"name": "John", "age": 30},
			[]string{"Available keys: name, age", "Available keys: age, name"}, // Order might vary
		},
		{
			"Not an object",
			[]string{"one", "two"},
			[]string{"Not a JSON object - no keys available"},
		},
		{
			"Empty object",
			map[string]interface{}{},
			[]string{"Available keys: "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatAvailableKeys(tt.data)
			found := false
			for _, want := range tt.wants {
				if got == want {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("FormatAvailableKeys() = %v, want one of %v", got, tt.wants)
			}
		})
	}
}

func TestMarshalToJSONString(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    string
		wantErr bool
	}{
		{"String", "hello", `"hello"`, false},
		{"Number", 123, `123`, false},
		{"Boolean", true, `true`, false},
		{"Null", nil, `null`, false},
		{"Array", []int{1, 2, 3}, `[1,2,3]`, false},
		{"Object", map[string]string{"name": "John"}, `{"name":"John"}`, false},
		{"Nested object", map[string]interface{}{"person": map[string]string{"name": "John"}}, `{"person":{"name":"John"}}`, false},
		// Cannot test invalid values that would fail JSON marshaling (e.g., channels) easily
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalToJSONString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalToJSONString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Normalize the strings to handle variations in formatting
				var gotObj, wantObj interface{}
				_ = json.Unmarshal([]byte(got), &gotObj)
				_ = json.Unmarshal([]byte(tt.want), &wantObj)
				if !reflect.DeepEqual(gotObj, wantObj) {
					t.Errorf("MarshalToJSONString() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMarshalToJSONResult(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"Valid value", map[string]string{"name": "John"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MarshalToJSONResult(tt.input)
			if result.IsFailure() != tt.wantErr {
				t.Errorf("MarshalToJSONResult().IsFailure() = %v, want %v", result.IsFailure(), tt.wantErr)
			}
			if !tt.wantErr {
				// Verify that the result is valid JSON
				var parsed interface{}
				err := json.Unmarshal([]byte(result.Unwrap()), &parsed)
				if err != nil {
					t.Errorf("MarshalToJSONResult() did not produce valid JSON: %v", err)
				}
			}
		})
	}
}

func TestCompactJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Empty string", "", ""},
		{"Already compact", `{"name":"John"}`, `{"name":"John"}`},
		{"With whitespace", `{ "name" : "John" }`, `{"name":"John"}`},
		{"Invalid JSON", `{ "name": "John" `, `{ "name": "John" `}, // Returns original if invalid
		{"Complex JSON", `{
			"name": "John",
			"age": 30,
			"address": {
				"city": "Tokyo"
			}
		}`, `{"name":"John","age":30,"address":{"city":"Tokyo"}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompactJSON(tt.input)

			// For valid JSON, compare by parsing first to handle formatting variations
			if IsJSON(tt.input) {
				var gotObj, wantObj interface{}
				_ = json.Unmarshal([]byte(got), &gotObj)
				_ = json.Unmarshal([]byte(tt.want), &wantObj)
				if !reflect.DeepEqual(gotObj, wantObj) {
					t.Errorf("CompactJSON() = %v, want %v", got, tt.want)
				}
			} else {
				// For invalid JSON, compare strings directly
				if got != tt.want {
					t.Errorf("CompactJSON() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestCompactJSONResult(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid JSON", `{ "name": "John" }`, false},
		{"Invalid JSON", `{ "name": "John" `, false}, // Doesn't error, returns original
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompactJSONResult(tt.input)
			if result.IsFailure() != tt.wantErr {
				t.Errorf("CompactJSONResult().IsFailure() = %v, want %v", result.IsFailure(), tt.wantErr)
			}
			if !tt.wantErr {
				// Check that result is the same as CompactJSON()
				expected := CompactJSON(tt.input)
				if result.Unwrap() != expected {
					t.Errorf("CompactJSONResult() = %v, want %v", result.Unwrap(), expected)
				}
			}
		})
	}
}
