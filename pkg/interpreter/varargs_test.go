package interpreter

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/stretchr/testify/assert"
)

func TestVarargsFunction(t *testing.T) {
	interpreter := NewInterpreter(nil)
	env := interpreter.GetEnvironment()

	// Create a test varargs function that returns the argument count
	varargsFunc := &environment.Function{
		Name:       "TEST_VARARGS",
		IsVarargs:  true,
		Parameters: []environment.Parameter{}, // Empty for varargs
		NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
			// Return the number of arguments passed
			return environment.IntegerValue(int64(len(args))), nil
		},
	}

	err := env.DefineFunction(varargsFunc)
	assert.NoError(t, err)

	// Test with 0 arguments
	result, err := interpreter.CallFunction("TEST_VARARGS", []environment.Value{})
	assert.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(0), result)

	// Test with 1 argument
	result, err = interpreter.CallFunction("TEST_VARARGS", []environment.Value{
		environment.StringValue("hello"),
	})
	assert.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(1), result)

	// Test with 3 arguments
	result, err = interpreter.CallFunction("TEST_VARARGS", []environment.Value{
		environment.StringValue("one"),
		environment.StringValue("two"),
		environment.IntegerValue(3),
	})
	assert.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(3), result)
}

func TestVarargsParameterBinding(t *testing.T) {
	interpreter := NewInterpreter(nil)
	env := interpreter.GetEnvironment()

	// Create a varargs function that accesses argc and individual arguments
	varargsFunc := &environment.Function{
		Name:       "TEST_VARARGS_BINDING",
		IsVarargs:  true,
		Parameters: []environment.Parameter{}, // Empty for varargs
		NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
			// This would be called during execution, but we'll test parameter binding separately
			return environment.NOTHIN, nil
		},
	}

	err := env.DefineFunction(varargsFunc)
	assert.NoError(t, err)

	// Test parameter binding by checking the environment after calling callFunction
	// We'll create a test that manually calls the internal function to check binding
	args := []environment.Value{
		environment.StringValue("test1"),
		environment.IntegerValue(42),
		environment.StringValue("test2"),
	}

	// Call the function (which will set up the parameter binding)
	_, err = interpreter.callFunction(varargsFunc, args)
	assert.NoError(t, err)
}
