package stdlib

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// Global MATH function definitions - created once and reused
var mathFuncOnce = sync.Once{}
var mathFunctions map[string]*environment.Function

func getMathFunctions() map[string]*environment.Function {
	mathFuncOnce.Do(func() {
		mathFunctions = map[string]*environment.Function{
			"ABS": {
				Name:       "ABS",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(types.DoubleValue); ok {
						return types.DoubleValue(math.Abs(float64(doubleVal))), nil
					}

					return types.NOTHIN, fmt.Errorf("ABS: invalid numeric type")
				},
			},
			"MAX": {
				Name:       "MAX",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "a", Type: "DUBBLE"},
					{Name: "b", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					a, b := args[0], args[1]

					if val1, ok := a.(types.DoubleValue); ok {
						if val2, ok := b.(types.DoubleValue); ok {
							return types.DoubleValue(math.Max(float64(val1), float64(val2))), nil
						}
					}

					return types.NOTHIN, fmt.Errorf("MAX: invalid numeric types")
				},
			},
			"MIN": {
				Name:       "MIN",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "a", Type: "DUBBLE"},
					{Name: "b", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					a, b := args[0], args[1]

					if val1, ok := a.(types.DoubleValue); ok {
						if val2, ok := b.(types.DoubleValue); ok {
							return types.DoubleValue(math.Min(float64(val1), float64(val2))), nil
						}
					}

					return types.NOTHIN, fmt.Errorf("MIN: invalid numeric types")
				},
			},
			"SQRT": {
				Name:       "SQRT",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(types.DoubleValue); ok {
						if float64(doubleVal) < 0 {
							return types.NOTHIN, fmt.Errorf("SQRT: negative argument")
						}
						return types.DoubleValue(math.Sqrt(float64(doubleVal))), nil
					}

					return types.NOTHIN, fmt.Errorf("SQRT: invalid numeric type")
				},
			},
			"POW": {
				Name:       "POW",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "base", Type: "DUBBLE"},
					{Name: "exponent", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					base, exp := args[0], args[1]

					if baseVal, ok := base.(types.DoubleValue); ok {
						if expVal, ok := exp.(types.DoubleValue); ok {
							result := math.Pow(float64(baseVal), float64(expVal))
							return types.DoubleValue(result), nil
						}
					}

					return types.NOTHIN, fmt.Errorf("POW: invalid numeric types")
				},
			},
			"RANDOM": {
				Name:       "RANDOM",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					return types.DoubleValue(rand.Float64()), nil
				},
			},
			"RANDINT": {
				Name:       "RANDINT",
				ReturnType: "INTEGR",
				Parameters: []environment.Parameter{
					{Name: "min", Type: "INTEGR"},
					{Name: "max", Type: "INTEGR"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
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
			},
			"SIN": {
				Name:       "SIN",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(types.DoubleValue); ok {
						return types.DoubleValue(math.Sin(float64(doubleVal))), nil
					}

					return types.NOTHIN, fmt.Errorf("SIN: invalid numeric type")
				},
			},
			"COS": {
				Name:       "COS",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(types.DoubleValue); ok {
						return types.DoubleValue(math.Cos(float64(doubleVal))), nil
					}

					return types.NOTHIN, fmt.Errorf("COS: invalid numeric type")
				},
			},
		}
	})
	return mathFunctions
}

// RegisterMATHInEnv registers MATH functions in the given environment
// declarations: empty slice means import all, otherwise import only specified functions
func RegisterMATHInEnv(env *environment.Environment, declarations ...string) error {
	mathFunctions := getMathFunctions()

	// If declarations is empty, import all functions
	if len(declarations) == 0 {
		for _, fn := range mathFunctions {
			env.DefineFunction(fn)
		}
		return nil
	}

	// Otherwise, import only specified functions
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if fn, exists := mathFunctions[declUpper]; exists {
			env.DefineFunction(fn)
		} else {
			return fmt.Errorf("unknown MATH function: %s", decl)
		}
	}

	return nil
}
