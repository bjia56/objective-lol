package stdlib

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
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
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Abs(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("ABS: invalid numeric type")
				},
			},
			"MAX": {
				Name:       "MAX",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "a", Type: "DUBBLE"},
					{Name: "b", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					a, b := args[0], args[1]

					if val1, ok := a.(environment.DoubleValue); ok {
						if val2, ok := b.(environment.DoubleValue); ok {
							return environment.DoubleValue(math.Max(float64(val1), float64(val2))), nil
						}
					}

					return environment.NOTHIN, fmt.Errorf("MAX: invalid numeric environment")
				},
			},
			"MIN": {
				Name:       "MIN",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "a", Type: "DUBBLE"},
					{Name: "b", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					a, b := args[0], args[1]

					if val1, ok := a.(environment.DoubleValue); ok {
						if val2, ok := b.(environment.DoubleValue); ok {
							return environment.DoubleValue(math.Min(float64(val1), float64(val2))), nil
						}
					}

					return environment.NOTHIN, fmt.Errorf("MIN: invalid numeric environment")
				},
			},
			"SQRT": {
				Name:       "SQRT",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						if float64(doubleVal) < 0 {
							return environment.NOTHIN, fmt.Errorf("SQRT: negative argument")
						}
						return environment.DoubleValue(math.Sqrt(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("SQRT: invalid numeric type")
				},
			},
			"POW": {
				Name:       "POW",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "base", Type: "DUBBLE"},
					{Name: "exponent", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					base, exp := args[0], args[1]

					if baseVal, ok := base.(environment.DoubleValue); ok {
						if expVal, ok := exp.(environment.DoubleValue); ok {
							result := math.Pow(float64(baseVal), float64(expVal))
							return environment.DoubleValue(result), nil
						}
					}

					return environment.NOTHIN, fmt.Errorf("POW: invalid numeric environment")
				},
			},
			"RANDOM": {
				Name:       "RANDOM",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					return environment.DoubleValue(rand.Float64()), nil
				},
			},
			"RANDINT": {
				Name:       "RANDINT",
				ReturnType: "INTEGR",
				Parameters: []environment.Parameter{
					{Name: "min", Type: "INTEGR"},
					{Name: "max", Type: "INTEGR"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					min, max := args[0], args[1]

					if minVal, ok := min.(environment.IntegerValue); ok {
						if maxVal, ok := max.(environment.IntegerValue); ok {
							if minVal >= maxVal {
								return environment.NOTHIN, fmt.Errorf("RANDINT: min must be less than max")
							}
							result := rand.Int63n(int64(maxVal-minVal)) + int64(minVal)
							return environment.IntegerValue(result), nil
						}
					}

					return environment.NOTHIN, fmt.Errorf("RANDINT: invalid integer environment")
				},
			},
			"SIN": {
				Name:       "SIN",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Sin(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("SIN: invalid numeric type")
				},
			},
			"COS": {
				Name:       "COS",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Cos(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("COS: invalid numeric type")
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
