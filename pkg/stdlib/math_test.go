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

	// Test that math functions are registered
	mathFunctions := []string{"ABS", "MAX", "MIN", "SQRT", "POW", "SIN", "COS", "RANDOM", "RANDINT"}

	for _, funcName := range mathFunctions {
		_, err := env.GetFunction(funcName)
		assert.NoError(t, err, "Function %s should be registered", funcName)
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

func TestMathRANDOM(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	randomFunc, err := env.GetFunction("RANDOM")
	require.NoError(t, err)

	// Test that RANDOM returns values in expected range
	for i := 0; i < 100; i++ {
		args := []environment.Value{}
		result, err := randomFunc.NativeImpl(nil, nil, args)
		require.NoError(t, err)

		doubleResult, ok := result.(environment.DoubleValue)
		require.True(t, ok, "RANDOM should return DoubleValue")

		val := float64(doubleResult)
		assert.GreaterOrEqual(t, val, 0.0, "RANDOM should be >= 0")
		assert.Less(t, val, 1.0, "RANDOM should be < 1")
	}
}

func TestMathRANDINT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterMATHInEnv(env)

	randintFunc, err := env.GetFunction("RANDINT")
	require.NoError(t, err)

	// Test that RANDINT returns values in expected range
	for i := 0; i < 50; i++ {
		args := []environment.Value{environment.IntegerValue(1), environment.IntegerValue(10)}
		result, err := randintFunc.NativeImpl(nil, nil, args)
		require.NoError(t, err)

		intResult, ok := result.(environment.IntegerValue)
		require.True(t, ok, "RANDINT should return IntegerValue")

		val := int64(intResult)
		assert.GreaterOrEqual(t, val, int64(1), "RANDINT should be >= min")
		assert.Less(t, val, int64(10), "RANDINT should be < max")
	}

	// Test error case: min >= max
	args := []environment.Value{environment.IntegerValue(10), environment.IntegerValue(5)}
	_, err = randintFunc.NativeImpl(nil, nil, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "min must be less than max")
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
