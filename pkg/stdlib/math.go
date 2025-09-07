package stdlib

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
)

// Global MATH function definitions - created once and reused
var mathFuncOnce = sync.Once{}
var mathFunctions map[string]*environment.Function

// Global MATH variable definitions - created once and reused
var mathVarsOnce = sync.Once{}
var mathVariables map[string]*environment.Variable

func getMathVariables() map[string]*environment.Variable {
	mathVarsOnce.Do(func() {
		mathVariables = map[string]*environment.Variable{
			"PI": {
				Name:     "PI",
				Type:     "DUBBLE",
				Value:    environment.DoubleValue(math.Pi),
				IsLocked: true,
				IsPublic: true,
			},
			"E": {
				Name:     "E",
				Type:     "DUBBLE",
				Value:    environment.DoubleValue(math.E),
				IsLocked: true,
				IsPublic: true,
			},
		}
	})
	return mathVariables
}

func getMathFunctions() map[string]*environment.Function {
	mathFuncOnce.Do(func() {
		mathFunctions = map[string]*environment.Function{
			"ABS": {
				Name:       "ABS",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
			"SIN": {
				Name:       "SIN",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Cos(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("COS: invalid numeric type")
				},
			},
			"TAN": {
				Name:       "TAN",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Tan(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("TAN: invalid numeric type")
				},
			},
			"ASIN": {
				Name:       "ASIN",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						val := float64(doubleVal)
						if val < -1 || val > 1 {
							return environment.NOTHIN, fmt.Errorf("ASIN: input out of range [-1, 1]")
						}
						return environment.DoubleValue(math.Asin(val)), nil
					}

					return environment.NOTHIN, fmt.Errorf("ASIN: invalid numeric type")
				},
			},
			"ACOS": {
				Name:       "ACOS",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						val := float64(doubleVal)
						if val < -1 || val > 1 {
							return environment.NOTHIN, fmt.Errorf("ACOS: input out of range [-1, 1]")
						}
						return environment.DoubleValue(math.Acos(val)), nil
					}

					return environment.NOTHIN, fmt.Errorf("ACOS: invalid numeric type")
				},
			},
			"ATAN": {
				Name:       "ATAN",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Atan(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("ATAN: invalid numeric type")
				},
			},
			"ATAN2": {
				Name:       "ATAN2",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "y", Type: "DUBBLE"},
					{Name: "x", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					y, x := args[0], args[1]

					if yVal, ok := y.(environment.DoubleValue); ok {
						if xVal, ok := x.(environment.DoubleValue); ok {
							return environment.DoubleValue(math.Atan2(float64(yVal), float64(xVal))), nil
						}
					}

					return environment.NOTHIN, fmt.Errorf("ATAN2: invalid numeric environment")
				},
			},
			"LOG": {
				Name:       "LOG",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						val := float64(doubleVal)
						if val <= 0 {
							return environment.NOTHIN, fmt.Errorf("LOG: input must be positive")
						}
						return environment.DoubleValue(math.Log(val)), nil
					}

					return environment.NOTHIN, fmt.Errorf("LOG: invalid numeric type")
				},
			},
			"LOG10": {
				Name:       "LOG10",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						val := float64(doubleVal)
						if val <= 0 {
							return environment.NOTHIN, fmt.Errorf("LOG10: input must be positive")
						}
						return environment.DoubleValue(math.Log10(val)), nil
					}

					return environment.NOTHIN, fmt.Errorf("LOG10: invalid numeric type")
				},
			},
			"LOG2": {
				Name:       "LOG2",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						val := float64(doubleVal)
						if val <= 0 {
							return environment.NOTHIN, fmt.Errorf("LOG2: input must be positive")
						}
						return environment.DoubleValue(math.Log2(val)), nil
					}

					return environment.NOTHIN, fmt.Errorf("LOG2: invalid numeric type")
				},
			},
			"EXP": {
				Name:       "EXP",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Exp(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("EXP: invalid numeric type")
				},
			},
			"CEIL": {
				Name:       "CEIL",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Ceil(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("CEIL: invalid numeric type")
				},
			},
			"FLOOR": {
				Name:       "FLOOR",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Floor(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("FLOOR: invalid numeric type")
				},
			},
			"ROUND": {
				Name:       "ROUND",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Round(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("ROUND: invalid numeric type")
				},
			},
			"TRUNC": {
				Name:       "TRUNC",
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Trunc(float64(doubleVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("TRUNC: invalid numeric type")
				},
			},
		}
	})
	return mathFunctions
}

// RegisterMATHInEnv registers MATH functions and variables in the given environment
// declarations: empty slice means import all, otherwise import only specified functions/variables
func RegisterMATHInEnv(env *environment.Environment, declarations ...string) error {
	mathFunctions := getMathFunctions()
	mathVariables := getMathVariables()

	// If declarations is empty, import all functions and variables
	if len(declarations) == 0 {
		for _, fn := range mathFunctions {
			env.DefineFunction(fn)
		}
		for _, variable := range mathVariables {
			err := env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, nil)
			if err != nil {
				return fmt.Errorf("failed to define MATH variable %s: %v", variable.Name, err)
			}
		}
		return nil
	}

	// Otherwise, import only specified functions and variables
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if fn, exists := mathFunctions[declUpper]; exists {
			env.DefineFunction(fn)
		} else if variable, exists := mathVariables[declUpper]; exists {
			err := env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, nil)
			if err != nil {
				return fmt.Errorf("failed to define MATH variable %s: %v", variable.Name, err)
			}
		} else {
			return fmt.Errorf("unknown MATH declaration: %s", decl)
		}
	}

	return nil
}
