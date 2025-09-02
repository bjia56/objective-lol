package stdlib

import (
	"math"
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
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
		input    types.Value
		expected types.Value
	}{
		{
			"Positive number",
			types.DoubleValue(5.0),
			types.DoubleValue(5.0),
		},
		{
			"Negative number",
			types.DoubleValue(-5.0),
			types.DoubleValue(5.0),
		},
		{
			"Positive double",
			types.DoubleValue(3.14),
			types.DoubleValue(3.14),
		},
		{
			"Negative double",
			types.DoubleValue(-3.14),
			types.DoubleValue(3.14),
		},
		{
			"Zero",
			types.DoubleValue(0.0),
			types.DoubleValue(0.0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []types.Value{test.input}
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
		a        types.Value
		b        types.Value
		expected types.Value
	}{
		{
			"First larger",
			types.DoubleValue(10),
			types.DoubleValue(5),
			types.DoubleValue(10.0),
		},
		{
			"Second larger",
			types.DoubleValue(3),
			types.DoubleValue(7),
			types.DoubleValue(7.0),
		},
		{
			"Equal values",
			types.DoubleValue(5.5),
			types.DoubleValue(5.5),
			types.DoubleValue(5.5),
		},
		{
			"Mixed types",
			types.DoubleValue(4),
			types.DoubleValue(4.1),
			types.DoubleValue(4.1),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []types.Value{test.a, test.b}
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
		a        types.Value
		b        types.Value
		expected types.Value
	}{
		{
			"First smaller",
			types.DoubleValue(3),
			types.DoubleValue(8),
			types.DoubleValue(3.0),
		},
		{
			"Second smaller",
			types.DoubleValue(10),
			types.DoubleValue(2),
			types.DoubleValue(2.0),
		},
		{
			"Negative values",
			types.DoubleValue(-5),
			types.DoubleValue(-2),
			types.DoubleValue(-5.0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []types.Value{test.a, test.b}
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
		input    types.Value
		expected float64
	}{
		{
			"Perfect square",
			types.DoubleValue(9),
			3.0,
		},
		{
			"Non-perfect square",
			types.DoubleValue(2),
			math.Sqrt(2),
		},
		{
			"Zero",
			types.DoubleValue(0.0),
			0.0,
		},
		{
			"Decimal",
			types.DoubleValue(6.25),
			2.5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []types.Value{test.input}
			result, err := sqrtFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(types.DoubleValue)
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
		base     types.Value
		exponent types.Value
		expected float64
	}{
		{
			"Square",
			types.DoubleValue(3),
			types.DoubleValue(2),
			9.0,
		},
		{
			"Cube",
			types.DoubleValue(2),
			types.DoubleValue(3),
			8.0,
		},
		{
			"Power of zero",
			types.DoubleValue(5),
			types.DoubleValue(0),
			1.0,
		},
		{
			"Fractional exponent",
			types.DoubleValue(4),
			types.DoubleValue(0.5),
			2.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []types.Value{test.base, test.exponent}
			result, err := powFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(types.DoubleValue)
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
		input    types.Value
		expected float64
	}{
		{
			"Sin of 0",
			types.DoubleValue(0.0),
			0.0,
		},
		{
			"Sin of π/2",
			types.DoubleValue(math.Pi / 2),
			1.0,
		},
		{
			"Sin of π",
			types.DoubleValue(math.Pi),
			0.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []types.Value{test.input}
			result, err := sinFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(types.DoubleValue)
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
		input    types.Value
		expected float64
	}{
		{
			"Cos of 0",
			types.DoubleValue(0.0),
			1.0,
		},
		{
			"Cos of π/2",
			types.DoubleValue(math.Pi / 2),
			0.0,
		},
		{
			"Cos of π",
			types.DoubleValue(math.Pi),
			-1.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []types.Value{test.input}
			result, err := cosFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			doubleResult, ok := result.(types.DoubleValue)
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
		args := []types.Value{}
		result, err := randomFunc.NativeImpl(nil, nil, args)
		require.NoError(t, err)

		doubleResult, ok := result.(types.DoubleValue)
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
		args := []types.Value{types.IntegerValue(1), types.IntegerValue(10)}
		result, err := randintFunc.NativeImpl(nil, nil, args)
		require.NoError(t, err)

		intResult, ok := result.(types.IntegerValue)
		require.True(t, ok, "RANDINT should return IntegerValue")

		val := int64(intResult)
		assert.GreaterOrEqual(t, val, int64(1), "RANDINT should be >= min")
		assert.Less(t, val, int64(10), "RANDINT should be < max")
	}

	// Test error case: min >= max
	args := []types.Value{types.IntegerValue(10), types.IntegerValue(5)}
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

	args := []types.Value{types.DoubleValue(-4.0)}
	_, err = sqrtFunc.NativeImpl(nil, nil, args)

	// Should return error for negative numbers in this implementation
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SQRT: negative argument")
}
