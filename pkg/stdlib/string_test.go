package stdlib

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterSTRING(t *testing.T) {
	env := environment.NewEnvironment(nil)

	err := RegisterSTRINGInEnv(env)
	require.NoError(t, err)

	// Test that STRING functions are registered
	stringFunctions := []string{"LEN", "CONCAT", "SUBSTR", "TRIM", "LTRIM", "RTRIM", "REPEAT", "UPPER", "LOWER", "TITLE", "CAPITALIZE", "SPLIT", "REPLACE", "REPLACE_ALL", "CONTAINS", "INDEX_OF"}

	for _, funcName := range stringFunctions {
		_, err = env.GetFunction(funcName)
		assert.NoError(t, err, "Function %s should be registered", funcName)
	}
}

func TestRegisterSTRINGSelective(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test selective import
	err := RegisterSTRINGInEnv(env, "LEN")
	require.NoError(t, err)

	// Should have LEN
	_, err = env.GetFunction("LEN")
	assert.NoError(t, err, "LEN function should be registered")

	// Should not have CONCAT
	_, err = env.GetFunction("CONCAT")
	assert.Error(t, err, "CONCAT function should not be registered")
}

func TestStringLEN(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)

	lenFunc, err := env.GetFunction("LEN")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Single character", "a", 1},
		{"Short string", "Hello", 5},
		{"Long string", "Hello, World!", 13},
		{"String with spaces", "Hello World", 11},
		{"String with special characters", "Hello, World! 123", 17},
		{"Unicode string", "café", 5},
		{"String with newlines", "line1\nline2", 11},
		{"String with tabs", "tab\there", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := lenFunc.NativeImpl(nil, nil, []environment.Value{environment.StringValue(tt.input)})
			require.NoError(t, err)

			intResult, ok := result.(environment.IntegerValue)
			require.True(t, ok, "LEN should return an integer")
			assert.Equal(t, tt.expected, int(intResult))
		})
	}
}

func TestStringLENErrorHandling(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)

	lenFunc, err := env.GetFunction("LEN")
	require.NoError(t, err)

	tests := []struct {
		name  string
		args  []environment.Value
		error string
	}{
		{
			"Wrong argument type - integer",
			[]environment.Value{environment.IntegerValue(42)},
			"LEN: argument is not a string",
		},
		{
			"Wrong argument type - boolean",
			[]environment.Value{environment.YEZ},
			"LEN: argument is not a string",
		},
		{
			"Wrong argument type - double",
			[]environment.Value{environment.DoubleValue(3.14)},
			"LEN: argument is not a string",
		},
		{
			"Wrong argument type - nothing",
			[]environment.Value{environment.NOTHIN},
			"LEN: argument is not a string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := lenFunc.NativeImpl(nil, nil, tt.args)
			if tt.error != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.error)
				assert.Equal(t, environment.NOTHIN, result)
			}
		})
	}
}

func TestStringCONCAT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)

	concatFunc, err := env.GetFunction("CONCAT")
	require.NoError(t, err)

	tests := []struct {
		name     string
		str1     string
		str2     string
		expected string
	}{
		{"Two non-empty strings", "Hello", " World", "Hello World"},
		{"First empty", "", "World", "World"},
		{"Second empty", "Hello", "", "Hello"},
		{"Both empty", "", "", ""},
		{"With numbers", "Count: ", "42", "Count: 42"},
		{"With special characters", "Hello!", " How are you?", "Hello! How are you?"},
		{"Unicode strings", "café", " ☕", "café ☕"},
		{"Long strings", "This is a long string", " and this is another long string", "This is a long string and this is another long string"},
		{"Newlines and tabs", "line1\n", "line2\t", "line1\nline2\t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := concatFunc.NativeImpl(nil, nil, []environment.Value{
				environment.StringValue(tt.str1),
				environment.StringValue(tt.str2),
			})
			require.NoError(t, err)

			strResult, ok := result.(environment.StringValue)
			require.True(t, ok, "CONCAT should return a string")
			assert.Equal(t, tt.expected, string(strResult))
		})
	}
}

func TestStringCONCATErrorHandling(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)

	concatFunc, err := env.GetFunction("CONCAT")
	require.NoError(t, err)

	tests := []struct {
		name  string
		args  []environment.Value
		error string
	}{
		{
			"First argument not string",
			[]environment.Value{environment.IntegerValue(42), environment.StringValue("world")},
			"CONCAT: first argument is not a string",
		},
		{
			"Second argument not string",
			[]environment.Value{environment.StringValue("hello"), environment.IntegerValue(42)},
			"CONCAT: second argument is not a string",
		},
		{
			"Both arguments not string",
			[]environment.Value{environment.IntegerValue(1), environment.IntegerValue(2)},
			"CONCAT: first argument is not a string",
		},
		{
			"First argument boolean",
			[]environment.Value{environment.YEZ, environment.StringValue("test")},
			"CONCAT: first argument is not a string",
		},
		{
			"Second argument boolean",
			[]environment.Value{environment.StringValue("test"), environment.NO},
			"CONCAT: second argument is not a string",
		},
		{
			"First argument nothing",
			[]environment.Value{environment.NOTHIN, environment.StringValue("test")},
			"CONCAT: first argument is not a string",
		},
		{
			"Second argument nothing",
			[]environment.Value{environment.StringValue("test"), environment.NOTHIN},
			"CONCAT: second argument is not a string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := concatFunc.NativeImpl(nil, nil, tt.args)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.error)
			assert.Equal(t, environment.NOTHIN, result)
		})
	}
}

func TestSTRINGSelectiveImport(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test importing only CONCAT
	err := RegisterSTRINGInEnv(env, "CONCAT")
	require.NoError(t, err)

	// Should have CONCAT
	_, err = env.GetFunction("CONCAT")
	assert.NoError(t, err, "CONCAT function should be registered")

	// Should not have LEN
	_, err = env.GetFunction("LEN")
	assert.Error(t, err, "LEN function should not be registered")
}

func TestSTRINGInvalidFunction(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test importing non-existent function
	err := RegisterSTRINGInEnv(env, "NONEXISTENT")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown STRING function: NONEXISTENT")
}

func TestStringUPPER(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)

	upperFunc, err := env.GetFunction("UPPER")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Lowercase string", "hello", "HELLO"},
		{"Mixed case string", "Hello World", "HELLO WORLD"},
		{"Already uppercase", "HELLO", "HELLO"},
		{"Empty string", "", ""},
		{"Numbers and symbols", "hello123!", "HELLO123!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := upperFunc.NativeImpl(nil, nil, []environment.Value{environment.StringValue(tt.input)})
			require.NoError(t, err)

			strResult, ok := result.(environment.StringValue)
			require.True(t, ok, "UPPER should return a string")
			assert.Equal(t, tt.expected, string(strResult))
		})
	}
}

func TestStringLOWER(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)

	lowerFunc, err := env.GetFunction("LOWER")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Uppercase string", "HELLO", "hello"},
		{"Mixed case string", "Hello World", "hello world"},
		{"Already lowercase", "hello", "hello"},
		{"Empty string", "", ""},
		{"Numbers and symbols", "HELLO123!", "hello123!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := lowerFunc.NativeImpl(nil, nil, []environment.Value{environment.StringValue(tt.input)})
			require.NoError(t, err)

			strResult, ok := result.(environment.StringValue)
			require.True(t, ok, "LOWER should return a string")
			assert.Equal(t, tt.expected, string(strResult))
		})
	}
}

func TestStringTRIM(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)

	trimFunc, err := env.GetFunction("TRIM")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Leading spaces", "  hello", "hello"},
		{"Trailing spaces", "hello  ", "hello"},
		{"Both sides", "  hello  ", "hello"},
		{"No spaces", "hello", "hello"},
		{"Only spaces", "   ", ""},
		{"Tabs and newlines", "\t\nhello\n\t", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := trimFunc.NativeImpl(nil, nil, []environment.Value{environment.StringValue(tt.input)})
			require.NoError(t, err)

			strResult, ok := result.(environment.StringValue)
			require.True(t, ok, "TRIM should return a string")
			assert.Equal(t, tt.expected, string(strResult))
		})
	}
}

func TestStringSPLIT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)
	RegisterArraysInEnv(env)

	splitFunc, err := env.GetFunction("SPLIT")
	require.NoError(t, err)

	tests := []struct {
		name      string
		input     string
		separator string
		expected  []string
	}{
		{"Comma separated", "a,b,c", ",", []string{"a", "b", "c"}},
		{"Space separated", "hello world test", " ", []string{"hello", "world", "test"}},
		{"Single character", "abc", "", []string{"a", "b", "c"}},
		{"No separator found", "hello", ",", []string{"hello"}},
		{"Empty string", "", ",", []string{""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := splitFunc.NativeImpl(nil, nil, []environment.Value{
				environment.StringValue(tt.input),
				environment.StringValue(tt.separator),
			})
			require.NoError(t, err)

			// Result should be a BUKKIT object
			bukkitObj, ok := result.(*environment.ObjectInstance)
			require.True(t, ok, "SPLIT should return a BUKKIT object")

			// Check the native data
			slice, ok := bukkitObj.NativeData.(BukkitSlice)
			require.True(t, ok, "BUKKIT should contain BukkitSlice")

			// Check length
			assert.Equal(t, len(tt.expected), len(slice), "Split result should have correct length")

			// Check each element
			for i, expected := range tt.expected {
				if i < len(slice) {
					strVal, ok := slice[i].(environment.StringValue)
					require.True(t, ok, "Split element should be string")
					assert.Equal(t, expected, string(strVal))
				}
			}
		})
	}
}

func TestStringREPLACE(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)

	replaceFunc, err := env.GetFunction("REPLACE")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		old      string
		new      string
		expected string
	}{
		{"Simple replace", "hello world", "world", "universe", "hello universe"},
		{"Replace first occurrence", "hello hello", "hello", "hi", "hi hello"},
		{"No match", "hello world", "foo", "bar", "hello world"},
		{"Empty replacement", "hello world", "world", "", "hello "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := replaceFunc.NativeImpl(nil, nil, []environment.Value{
				environment.StringValue(tt.input),
				environment.StringValue(tt.old),
				environment.StringValue(tt.new),
			})
			require.NoError(t, err)

			strResult, ok := result.(environment.StringValue)
			require.True(t, ok, "REPLACE should return a string")
			assert.Equal(t, tt.expected, string(strResult))
		})
	}
}

func TestSTRINGCaseInsensitive(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test case insensitive import
	err := RegisterSTRINGInEnv(env, "len", "concat")
	require.NoError(t, err)

	// Both functions should be available
	_, err = env.GetFunction("LEN")
	assert.NoError(t, err, "LEN function should be registered with lowercase import")

	_, err = env.GetFunction("CONCAT")
	assert.NoError(t, err, "CONCAT function should be registered with lowercase import")
}

func TestStringCONTAINS(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)

	containsFunc, err := env.GetFunction("CONTAINS")
	require.NoError(t, err)

	tests := []struct {
		name     string
		str      string
		substr   string
		expected bool
	}{
		{"Contains substring", "hello world", "world", true},
		{"Does not contain", "hello world", "foo", false},
		{"Empty substring", "hello", "", true},
		{"Empty string", "", "hello", false},
		{"Case sensitive", "Hello", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := containsFunc.NativeImpl(nil, nil, []environment.Value{
				environment.StringValue(tt.str),
				environment.StringValue(tt.substr),
			})
			require.NoError(t, err)

			boolResult, ok := result.(environment.BoolValue)
			require.True(t, ok, "CONTAINS should return a bool")
			assert.Equal(t, tt.expected, bool(boolResult))
		})
	}
}

func TestStringINDEX_OF(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTRINGInEnv(env)

	indexFunc, err := env.GetFunction("INDEX_OF")
	require.NoError(t, err)

	tests := []struct {
		name     string
		str      string
		substr   string
		expected int
	}{
		{"Found at beginning", "hello world", "hello", 0},
		{"Found in middle", "hello world", "lo wo", 3},
		{"Not found", "hello world", "foo", -1},
		{"Empty substring", "hello", "", 0},
		{"Empty string", "", "hello", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := indexFunc.NativeImpl(nil, nil, []environment.Value{
				environment.StringValue(tt.str),
				environment.StringValue(tt.substr),
			})
			require.NoError(t, err)

			intResult, ok := result.(environment.IntegerValue)
			require.True(t, ok, "INDEX_OF should return an integer")
			assert.Equal(t, tt.expected, int(intResult))
		})
	}
}
