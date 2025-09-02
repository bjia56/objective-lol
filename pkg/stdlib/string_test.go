package stdlib

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterSTRING(t *testing.T) {
	env := environment.NewEnvironment(nil)

	err := RegisterSTRINGInEnv(env)
	require.NoError(t, err)

	// Test that STRING functions are registered
	_, err = env.GetFunction("LEN")
	assert.NoError(t, err, "LEN function should be registered")

	_, err = env.GetFunction("CONCAT")
	assert.NoError(t, err, "CONCAT function should be registered")
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
			result, err := lenFunc.NativeImpl(nil, nil, []types.Value{types.StringValue(tt.input)})
			require.NoError(t, err)

			intResult, ok := result.(types.IntegerValue)
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
		args  []types.Value
		error string
	}{
		{
			"Wrong argument type - integer",
			[]types.Value{types.IntegerValue(42)},
			"LEN: argument is not a string",
		},
		{
			"Wrong argument type - boolean",
			[]types.Value{types.YEZ},
			"LEN: argument is not a string",
		},
		{
			"Wrong argument type - double",
			[]types.Value{types.DoubleValue(3.14)},
			"LEN: argument is not a string",
		},
		{
			"Wrong argument type - nothing",
			[]types.Value{types.NOTHIN},
			"LEN: argument is not a string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := lenFunc.NativeImpl(nil, nil, tt.args)
			if tt.error != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.error)
				assert.Equal(t, types.NOTHIN, result)
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
			result, err := concatFunc.NativeImpl(nil, nil, []types.Value{
				types.StringValue(tt.str1),
				types.StringValue(tt.str2),
			})
			require.NoError(t, err)

			strResult, ok := result.(types.StringValue)
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
		args  []types.Value
		error string
	}{
		{
			"First argument not string",
			[]types.Value{types.IntegerValue(42), types.StringValue("world")},
			"CONCAT: first argument is not a string",
		},
		{
			"Second argument not string",
			[]types.Value{types.StringValue("hello"), types.IntegerValue(42)},
			"CONCAT: second argument is not a string",
		},
		{
			"Both arguments not string",
			[]types.Value{types.IntegerValue(1), types.IntegerValue(2)},
			"CONCAT: first argument is not a string",
		},
		{
			"First argument boolean",
			[]types.Value{types.YEZ, types.StringValue("test")},
			"CONCAT: first argument is not a string",
		},
		{
			"Second argument boolean",
			[]types.Value{types.StringValue("test"), types.NO},
			"CONCAT: second argument is not a string",
		},
		{
			"First argument nothing",
			[]types.Value{types.NOTHIN, types.StringValue("test")},
			"CONCAT: first argument is not a string",
		},
		{
			"Second argument nothing",
			[]types.Value{types.StringValue("test"), types.NOTHIN},
			"CONCAT: second argument is not a string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := concatFunc.NativeImpl(nil, nil, tt.args)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.error)
			assert.Equal(t, types.NOTHIN, result)
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
