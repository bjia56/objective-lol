package stdlib

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// RegisterMATHInEnv registers all MATH functions directly in the given environment
func RegisterMATHInEnv(env *environment.Environment) {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// ABS function - absolute value
	abs := &environment.Function{
		Name:       "ABS",
		ReturnType: "DUBBLE",
		Parameters: []environment.Parameter{
			{Name: "value", Type: "DUBBLE"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			value := args[0]

			if doubleVal, ok := value.(types.DoubleValue); ok {
				return types.DoubleValue(math.Abs(float64(doubleVal))), nil
			}

			return types.NOTHIN, fmt.Errorf("ABS: invalid numeric type")
		},
	}
	env.DefineFunction(abs)

	// MAX function - maximum of two values
	max := &environment.Function{
		Name:       "MAX",
		ReturnType: "DUBBLE",
		Parameters: []environment.Parameter{
			{Name: "a", Type: "DUBBLE"},
			{Name: "b", Type: "DUBBLE"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			a, b := args[0], args[1]

			if val1, ok := a.(types.DoubleValue); ok {
				if val2, ok := b.(types.DoubleValue); ok {
					return types.DoubleValue(math.Max(float64(val1), float64(val2))), nil
				}
			}

			return types.NOTHIN, fmt.Errorf("MAX: invalid numeric types")
		},
	}
	env.DefineFunction(max)

	// MIN function - minimum of two values
	min := &environment.Function{
		Name:       "MIN",
		ReturnType: "DUBBLE",
		Parameters: []environment.Parameter{
			{Name: "a", Type: "DUBBLE"},
			{Name: "b", Type: "DUBBLE"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			a, b := args[0], args[1]

			if val1, ok := a.(types.DoubleValue); ok {
				if val2, ok := b.(types.DoubleValue); ok {
					return types.DoubleValue(math.Min(float64(val1), float64(val2))), nil
				}
			}

			return types.NOTHIN, fmt.Errorf("MIN: invalid numeric types")
		},
	}
	env.DefineFunction(min)

	// SQRT function - square root
	sqrt := &environment.Function{
		Name:       "SQRT",
		ReturnType: "DUBBLE",
		Parameters: []environment.Parameter{
			{Name: "value", Type: "DUBBLE"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			value := args[0]

			if doubleVal, ok := value.(types.DoubleValue); ok {
				if float64(doubleVal) < 0 {
					return types.NOTHIN, fmt.Errorf("SQRT: negative argument")
				}
				return types.DoubleValue(math.Sqrt(float64(doubleVal))), nil
			}

			return types.NOTHIN, fmt.Errorf("SQRT: invalid numeric type")
		},
	}
	env.DefineFunction(sqrt)

	// POW function - power
	pow := &environment.Function{
		Name:       "POW",
		ReturnType: "DUBBLE",
		Parameters: []environment.Parameter{
			{Name: "base", Type: "DUBBLE"},
			{Name: "exponent", Type: "DUBBLE"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			base, exp := args[0], args[1]

			if baseVal, ok := base.(types.DoubleValue); ok {
				if expVal, ok := exp.(types.DoubleValue); ok {
					result := math.Pow(float64(baseVal), float64(expVal))
					return types.DoubleValue(result), nil
				}
			}

			return types.NOTHIN, fmt.Errorf("POW: invalid numeric types")
		},
	}
	env.DefineFunction(pow)

	// RANDOM function - random number between 0 and 1
	random := &environment.Function{
		Name:       "RANDOM",
		ReturnType: "DUBBLE",
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.DoubleValue(rand.Float64()), nil
		},
	}
	env.DefineFunction(random)

	// RANDINT function - random integer in range
	randint := &environment.Function{
		Name:       "RANDINT",
		ReturnType: "INTEGR",
		Parameters: []environment.Parameter{
			{Name: "min", Type: "INTEGR"},
			{Name: "max", Type: "INTEGR"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			min, max := args[0], args[1]

			if minVal, ok := min.(types.IntegerValue); ok {
				if maxVal, ok := max.(types.IntegerValue); ok {
					if minVal >= maxVal {
						return types.NOTHIN, fmt.Errorf("RANDINT: min must be less than max")
					}
					result := rand.Int63n(int64(maxVal-minVal)) + int64(minVal)
					return types.IntegerValue(result), nil
				}
			}

			return types.NOTHIN, fmt.Errorf("RANDINT: invalid integer types")
		},
	}
	env.DefineFunction(randint)

	// SIN function - sine
	sin := &environment.Function{
		Name:       "SIN",
		ReturnType: "DUBBLE",
		Parameters: []environment.Parameter{
			{Name: "value", Type: "DUBBLE"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			value := args[0]

			if doubleVal, ok := value.(types.DoubleValue); ok {
				return types.DoubleValue(math.Sin(float64(doubleVal))), nil
			}

			return types.NOTHIN, fmt.Errorf("SIN: invalid numeric type")
		},
	}
	env.DefineFunction(sin)

	// COS function - cosine
	cos := &environment.Function{
		Name:       "COS",
		ReturnType: "DUBBLE",
		Parameters: []environment.Parameter{
			{Name: "value", Type: "DUBBLE"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			value := args[0]

			if doubleVal, ok := value.(types.DoubleValue); ok {
				return types.DoubleValue(math.Cos(float64(doubleVal))), nil
			}

			return types.NOTHIN, fmt.Errorf("COS: invalid numeric type")
		},
	}
	env.DefineFunction(cos)
}
