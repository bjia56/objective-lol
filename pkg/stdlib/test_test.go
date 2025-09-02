package stdlib

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterTEST(t *testing.T) {
	env := environment.NewEnvironment(nil)

	err := RegisterTESTInEnv(env)
	require.NoError(t, err)

	// Test that TEST functions are registered
	testFunctions := []string{"ASSERT"}

	for _, funcName := range testFunctions {
		_, err := env.GetFunction(funcName)
		assert.NoError(t, err, "Function %s should be registered", funcName)
	}
}

func TestASSERTFunction(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterTESTInEnv(env)

	assertFunc, err := env.GetFunction("ASSERT")
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition types.Value
		shouldErr bool
	}{
		// Truthy values - should not error
		{
			"True boolean",
			types.YEZ,
			false,
		},
		{
			"Non-zero integer",
			types.IntegerValue(42),
			false,
		},
		{
			"Negative integer",
			types.IntegerValue(-1),
			false,
		},
		{
			"Non-zero double",
			types.DoubleValue(3.14),
			false,
		},
		{
			"Non-empty string",
			types.StringValue("hello"),
			false,
		},

		// Falsy values - should error
		{
			"False boolean",
			types.NO,
			true,
		},
		{
			"Zero integer",
			types.IntegerValue(0),
			true,
		},
		{
			"Zero double",
			types.DoubleValue(0.0),
			true,
		},
		{
			"Empty string",
			types.StringValue(""),
			true,
		},
		{
			"Nothing value",
			types.NOTHIN,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []types.Value{test.condition}
			result, err := assertFunc.NativeImpl(nil, nil, args)

			if test.shouldErr {
				assert.Error(t, err, "ASSERT should fail for falsy condition")

				// Check if it's an Exception error
				if exception, ok := err.(ast.Exception); ok {
					assert.Contains(t, exception.Message, "Assertion failed")
				} else {
					t.Errorf("Expected ast.Exception, got %T", err)
				}
			} else {
				assert.NoError(t, err, "ASSERT should pass for truthy condition")
				assert.Equal(t, types.NOTHIN, result)
			}
		})
	}
}

func TestIsTruthyHelper(t *testing.T) {
	tests := []struct {
		name     string
		value    types.Value
		expected bool
	}{
		// Truthy cases
		{"True boolean", types.YEZ, true},
		{"Non-zero integer", types.IntegerValue(1), true},
		{"Negative integer", types.IntegerValue(-1), true},
		{"Non-zero double", types.DoubleValue(0.1), true},
		{"Non-empty string", types.StringValue("test"), true},

		// Falsy cases
		{"False boolean", types.NO, false},
		{"Zero integer", types.IntegerValue(0), false},
		{"Zero double", types.DoubleValue(0.0), false},
		{"Empty string", types.StringValue(""), false},
		{"Nothing value", types.NOTHIN, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isTruthy(test.value)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestTESTSelectiveImport(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Import only ASSERT function
	err := RegisterTESTInEnv(env, "ASSERT")
	require.NoError(t, err)

	// ASSERT should be available
	_, err = env.GetFunction("ASSERT")
	assert.NoError(t, err)
}

func TestTESTInvalidFunction(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Try to import non-existent function
	err := RegisterTESTInEnv(env, "INVALID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown TEST function: INVALID")
}

func TestTESTCaseInsensitive(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Import using lowercase
	err := RegisterTESTInEnv(env, "assert")
	require.NoError(t, err)

	// Function should be available (stored in uppercase)
	_, err = env.GetFunction("ASSERT")
	assert.NoError(t, err)
}
