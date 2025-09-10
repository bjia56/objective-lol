package stdlib

import (
	"math"
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterMath(t *testing.T) {
	env := environment.NewEnvironment(nil)

	err := RegisterMATHInEnv(env) // Empty slice imports all
	require.NoError(t, err)

	// Test that math functions and variables are registered
	mathFunctions := []string{"ABS", "MAX", "MIN", "SQRT", "POW", "SIN", "COS", "TAN", "ASIN", "ACOS", "ATAN", "ATAN2", "LOG", "LOG10", "LOG2", "EXP", "CEIL", "FLOOR", "ROUND", "TRUNC"}
	mathVariables := []string{"PI", "E"}

	for _, funcName := range mathFunctions {
		_, err := env.GetFunction(funcName)
		assert.NoError(t, err, "Function %s should be registered", funcName)
	}

	// Test that math variables are registered
	for _, varName := range mathVariables {
		_, err := env.GetVariable(varName)
		assert.NoError(t, err, "Variable %s should be registered", varName)
	}
}

func TestMathABS(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	absFunc, err := env.GetFunction("ABS")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    environment.Value
		expected environment.Value
	}{
		{
			"Positive number",
			environment.DoubleValue(5.0),
			environment.DoubleValue(5.0),
		},
		{
			"Negative number",
			environment.DoubleValue(-5.0),
			environment.DoubleValue(5.0),
		},
		{
			"Positive double",
			environment.DoubleValue(3.14),
			environment.DoubleValue(3.14),
		},
		{
			"Negative double",
			environment.DoubleValue(-3.14),
			environment.DoubleValue(3.14),
		},
		{
			"Zero",
			environment.DoubleValue(0.0),
			environment.DoubleValue(0.0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.input}
			result, err := absFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestMathMAX(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	maxFunc, err := env.GetFunction("MAX")
	require.NoError(t, err)

	tests := []struct {
		name     string
		a        environment.Value
		b        environment.Value
		expected environment.Value
	}{
		{
			"First larger",
			environment.DoubleValue(10),
			environment.DoubleValue(5),
			environment.DoubleValue(10.0),
		},
		{
			"Second larger",
			environment.DoubleValue(3),
			environment.DoubleValue(7),
			environment.DoubleValue(7.0),
		},
		{
			"Equal values",
			environment.DoubleValue(5.5),
			environment.DoubleValue(5.5),
			environment.DoubleValue(5.5),
		},
		{
			"Mixed environment",
			environment.DoubleValue(4),
			environment.DoubleValue(4.1),
			environment.DoubleValue(4.1),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.a, test.b}
			result, err := maxFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestMathMIN(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	minFunc, err := env.GetFunction("MIN")
	require.NoError(t, err)

	tests := []struct {
		name     string
		a        environment.Value
		b        environment.Value
		expected environment.Value
	}{
		{
			"First smaller",
			environment.DoubleValue(3),
			environment.DoubleValue(8),
			environment.DoubleValue(3.0),
		},
		{
			"Second smaller",
			environment.DoubleValue(10),
			environment.DoubleValue(2),
			environment.DoubleValue(2.0),
		},
		{
			"Negative values",
			environment.DoubleValue(-5),
			environment.DoubleValue(-2),
			environment.DoubleValue(-5.0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.a, test.b}
			result, err := minFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestMathSQRT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	sqrtFunc, err := env.GetFunction("SQRT")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    environment.Value
		expected float64
	}{
		{
			"Perfect square",
			environment.DoubleValue(9),
			3.0,
		},
		{
			"Non-perfect square",
			environment.DoubleValue(2),
			math.Sqrt(2),
		},
		{
			"Zero",
			environment.DoubleValue(0.0),
			0.0,
		},
		{
			"Decimal",
			environment.DoubleValue(6.25),
			2.5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.input}
			result, err := sqrtFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(environment.DoubleValue)
			require.True(t, ok, "SQRT should return DoubleValue")
			assert.InDelta(t, test.expected, float64(doubleResult), 0.0001)
		})
	}
}

func TestMathPOW(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	powFunc, err := env.GetFunction("POW")
	require.NoError(t, err)

	tests := []struct {
		name     string
		base     environment.Value
		exponent environment.Value
		expected float64
	}{
		{
			"Square",
			environment.DoubleValue(3),
			environment.DoubleValue(2),
			9.0,
		},
		{
			"Cube",
			environment.DoubleValue(2),
			environment.DoubleValue(3),
			8.0,
		},
		{
			"Power of zero",
			environment.DoubleValue(5),
			environment.DoubleValue(0),
			1.0,
		},
		{
			"Fractional exponent",
			environment.DoubleValue(4),
			environment.DoubleValue(0.5),
			2.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.base, test.exponent}
			result, err := powFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(environment.DoubleValue)
			require.True(t, ok, "POW should return DoubleValue")
			assert.InDelta(t, test.expected, float64(doubleResult), 0.0001)
		})
	}
}

func TestMathSIN(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	sinFunc, err := env.GetFunction("SIN")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    environment.Value
		expected float64
	}{
		{
			"Sin of 0",
			environment.DoubleValue(0.0),
			0.0,
		},
		{
			"Sin of π/2",
			environment.DoubleValue(math.Pi / 2),
			1.0,
		},
		{
			"Sin of π",
			environment.DoubleValue(math.Pi),
			0.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.input}
			result, err := sinFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(environment.DoubleValue)
			require.True(t, ok, "SIN should return DoubleValue")
			assert.InDelta(t, test.expected, float64(doubleResult), 0.0001)
		})
	}
}

func TestMathCOS(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	cosFunc, err := env.GetFunction("COS")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    environment.Value
		expected float64
	}{
		{
			"Cos of 0",
			environment.DoubleValue(0.0),
			1.0,
		},
		{
			"Cos of π/2",
			environment.DoubleValue(math.Pi / 2),
			0.0,
		},
		{
			"Cos of π",
			environment.DoubleValue(math.Pi),
			-1.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.input}
			result, err := cosFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(environment.DoubleValue)
			require.True(t, ok, "COS should return DoubleValue")
			assert.InDelta(t, test.expected, float64(doubleResult), 0.0001)
		})
	}
}

func TestMathErrorHandling(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	// Test SQRT with negative number
	sqrtFunc, err := env.GetFunction("SQRT")
	require.NoError(t, err)

	args := []environment.Value{environment.DoubleValue(-4.0)}
	_, err = sqrtFunc.NativeImpl(nil, nil, args)

	// Should return error for negative numbers in this implementation
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SQRT: negative argument")
}

func TestMathTAN(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	tanFunc, err := env.GetFunction("TAN")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    environment.Value
		expected float64
	}{
		{
			"Tan of 0",
			environment.DoubleValue(0.0),
			0.0,
		},
		{
			"Tan of π/4",
			environment.DoubleValue(math.Pi / 4),
			1.0,
		},
		{
			"Tan of π",
			environment.DoubleValue(math.Pi),
			0.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.input}
			result, err := tanFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(environment.DoubleValue)
			require.True(t, ok, "TAN should return DoubleValue")
			assert.InDelta(t, test.expected, float64(doubleResult), 0.0001)
		})
	}
}

func TestMathLOG(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	logFunc, err := env.GetFunction("LOG")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    environment.Value
		expected float64
	}{
		{
			"Log of E",
			environment.DoubleValue(math.E),
			1.0,
		},
		{
			"Log of 1",
			environment.DoubleValue(1.0),
			0.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.input}
			result, err := logFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(environment.DoubleValue)
			require.True(t, ok, "LOG should return DoubleValue")
			assert.InDelta(t, test.expected, float64(doubleResult), 0.0001)
		})
	}

	// Test error case: negative input
	args := []environment.Value{environment.DoubleValue(-1.0)}
	_, err = logFunc.NativeImpl(nil, nil, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LOG: input must be positive")
}

func TestMathCEIL(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	ceilFunc, err := env.GetFunction("CEIL")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    environment.Value
		expected float64
	}{
		{
			"Ceil of 3.2",
			environment.DoubleValue(3.2),
			4.0,
		},
		{
			"Ceil of 3.0",
			environment.DoubleValue(3.0),
			3.0,
		},
		{
			"Ceil of -3.7",
			environment.DoubleValue(-3.7),
			-3.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.input}
			result, err := ceilFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(environment.DoubleValue)
			require.True(t, ok, "CEIL should return DoubleValue")
			assert.Equal(t, test.expected, float64(doubleResult))
		})
	}
}

func TestMathFLOOR(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	floorFunc, err := env.GetFunction("FLOOR")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    environment.Value
		expected float64
	}{
		{
			"Floor of 3.7",
			environment.DoubleValue(3.7),
			3.0,
		},
		{
			"Floor of 3.0",
			environment.DoubleValue(3.0),
			3.0,
		},
		{
			"Floor of -3.2",
			environment.DoubleValue(-3.2),
			-4.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.input}
			result, err := floorFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(environment.DoubleValue)
			require.True(t, ok, "FLOOR should return DoubleValue")
			assert.Equal(t, test.expected, float64(doubleResult))
		})
	}
}

func TestMathROUND(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	roundFunc, err := env.GetFunction("ROUND")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    environment.Value
		expected float64
	}{
		{
			"Round of 3.4",
			environment.DoubleValue(3.4),
			3.0,
		},
		{
			"Round of 3.5",
			environment.DoubleValue(3.5),
			4.0,
		},
		{
			"Round of 3.6",
			environment.DoubleValue(3.6),
			4.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.input}
			result, err := roundFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(environment.DoubleValue)
			require.True(t, ok, "ROUND should return DoubleValue")
			assert.Equal(t, test.expected, float64(doubleResult))
		})
	}
}

func TestMathConstants(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	// Test PI constant
	piVar, err := env.GetVariable("PI")
	require.NoError(t, err)
	piValue, ok := piVar.Value.(environment.DoubleValue)
	require.True(t, ok, "PI should be a DoubleValue")
	assert.InDelta(t, math.Pi, float64(piValue), 0.0001)

	// Test E constant
	eVar, err := env.GetVariable("E")
	require.NoError(t, err)
	eValue, ok := eVar.Value.(environment.DoubleValue)
	require.True(t, ok, "E should be a DoubleValue")
	assert.InDelta(t, math.E, float64(eValue), 0.0001)
}
