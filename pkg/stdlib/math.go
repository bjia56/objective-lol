package stdlib

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
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
				Name: "PI",
				Documentation: []string{
					"The mathematical constant π (pi).",
				},
				Type:     "DUBBLE",
				Value:    environment.DoubleValue(math.Pi),
				IsLocked: true,
				IsPublic: true,
			},
			"E": {
				Name: "E",
				Documentation: []string{
					"Euler's number e.",
				},
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
				Name: "ABS",
				Documentation: []string{
					"Returns the absolute value of a number.",
					"Removes the sign and returns the positive magnitude.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Abs(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "ABS: invalid numeric type"}
				},
			},
			"MAX": {
				Name: "MAX",
				Documentation: []string{
					"Returns the larger of two numbers.",
					"Compares two values and returns the maximum.",
				},
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

					return environment.NOTHIN, runtime.Exception{Message: "MAX: invalid numeric environment"}
				},
			},
			"MIN": {
				Name: "MIN",
				Documentation: []string{
					"Returns the smaller of two numbers.",
					"Compares two values and returns the minimum.",
				},
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

					return environment.NOTHIN, runtime.Exception{Message: "MIN: invalid numeric environment"}
				},
			},
			"SQRT": {
				Name: "SQRT",
				Documentation: []string{
					"Returns the square root of a number.",
					"Input must be non-negative. Throws error for negative values.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						if float64(doubleVal) < 0 {
							return environment.NOTHIN, runtime.Exception{Message: "SQRT: negative argument"}
						}
						return environment.DoubleValue(math.Sqrt(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "SQRT: invalid numeric type"}
				},
			},
			"POW": {
				Name: "POW",
				Documentation: []string{
					"Returns base raised to the power of exponent (base^exponent).",
					"Performs exponentiation using floating-point arithmetic.",
				},
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

					return environment.NOTHIN, runtime.Exception{Message: "POW: invalid numeric environment"}
				},
			},
			"SIN": {
				Name: "SIN",
				Documentation: []string{
					"Returns the sine of an angle in radians.",
					"Input angle should be in radians, not degrees.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Sin(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "SIN: invalid numeric type"}
				},
			},
			"COS": {
				Name: "COS",
				Documentation: []string{
					"Returns the cosine of an angle in radians.",
					"Input angle should be in radians, not degrees.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Cos(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "COS: invalid numeric type"}
				},
			},
			"TAN": {
				Name: "TAN",
				Documentation: []string{
					"Returns the tangent of an angle in radians.",
					"Input angle should be in radians. Undefined at π/2 + nπ.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Tan(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "TAN: invalid numeric type"}
				},
			},
			"ASIN": {
				Name: "ASIN",
				Documentation: []string{
					"Returns the arc sine (inverse sine) of a value in radians.",
					"Input must be in range [-1, 1]. Result is in range [-π/2, π/2].",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						val := float64(doubleVal)
						if val < -1 || val > 1 {
							return environment.NOTHIN, runtime.Exception{Message: "ASIN: input out of range [-1, 1]"}
						}
						return environment.DoubleValue(math.Asin(val)), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "ASIN: invalid numeric type"}
				},
			},
			"ACOS": {
				Name: "ACOS",
				Documentation: []string{
					"Returns the arc cosine (inverse cosine) of a value in radians.",
					"Input must be in range [-1, 1]. Result is in range [0, π].",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						val := float64(doubleVal)
						if val < -1 || val > 1 {
							return environment.NOTHIN, runtime.Exception{Message: "ACOS: input out of range [-1, 1]"}
						}
						return environment.DoubleValue(math.Acos(val)), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "ACOS: invalid numeric type"}
				},
			},
			"ATAN": {
				Name: "ATAN",
				Documentation: []string{
					"Returns the arc tangent (inverse tangent) of a value in radians.",
					"Result is in range [-π/2, π/2].",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Atan(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "ATAN: invalid numeric type"}
				},
			},
			"ATAN2": {
				Name: "ATAN2",
				Documentation: []string{
					"Returns the arc tangent of y/x in radians, considering quadrant.",
					"Result is in range [-π, π]. More robust than ATAN for coordinate conversion.",
				},
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

					return environment.NOTHIN, runtime.Exception{Message: "ATAN2: invalid numeric environment"}
				},
			},
			"LOG": {
				Name: "LOG",
				Documentation: []string{
					"Returns the natural logarithm (base e) of a number.",
					"Input must be positive. Throws error for zero or negative values.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						val := float64(doubleVal)
						if val <= 0 {
							return environment.NOTHIN, runtime.Exception{Message: "LOG: input must be positive"}
						}
						return environment.DoubleValue(math.Log(val)), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "LOG: invalid numeric type"}
				},
			},
			"LOG10": {
				Name: "LOG10",
				Documentation: []string{
					"Returns the base-10 logarithm of a number.",
					"Input must be positive. Common logarithm for scientific calculations.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						val := float64(doubleVal)
						if val <= 0 {
							return environment.NOTHIN, runtime.Exception{Message: "LOG10: input must be positive"}
						}
						return environment.DoubleValue(math.Log10(val)), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "LOG10: invalid numeric type"}
				},
			},
			"LOG2": {
				Name: "LOG2",
				Documentation: []string{
					"Returns the base-2 logarithm of a number.",
					"Input must be positive. Useful for binary and computer science calculations.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						val := float64(doubleVal)
						if val <= 0 {
							return environment.NOTHIN, runtime.Exception{Message: "LOG2: input must be positive"}
						}
						return environment.DoubleValue(math.Log2(val)), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "LOG2: invalid numeric type"}
				},
			},
			"EXP": {
				Name: "EXP",
				Documentation: []string{
					"Returns e raised to the power of the given value (e^value).",
					"The exponential function, inverse of natural logarithm.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Exp(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "EXP: invalid numeric type"}
				},
			},
			"CEIL": {
				Name: "CEIL",
				Documentation: []string{
					"Returns the smallest integer greater than or equal to the value (ceiling).",
					"Rounds up to the next whole number.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Ceil(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "CEIL: invalid numeric type"}
				},
			},
			"FLOOR": {
				Name: "FLOOR",
				Documentation: []string{
					"Returns the largest integer less than or equal to the value (floor).",
					"Rounds down to the previous whole number.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Floor(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "FLOOR: invalid numeric type"}
				},
			},
			"ROUND": {
				Name: "ROUND",
				Documentation: []string{
					"Returns the value rounded to the nearest integer.",
					"Rounds 0.5 up to the next integer (round half up).",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Round(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "ROUND: invalid numeric type"}
				},
			},
			"TRUNC": {
				Name: "TRUNC",
				Documentation: []string{
					"Returns the integer part of a number by removing the fractional part.",
					"Truncates towards zero, different from floor for negative numbers.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "value", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					value := args[0]

					if doubleVal, ok := value.(environment.DoubleValue); ok {
						return environment.DoubleValue(math.Trunc(float64(doubleVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "TRUNC: invalid numeric type"}
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
			env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, variable.Documentation)
		}
		return nil
	}

	// Otherwise, import only specified functions and variables
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if fn, exists := mathFunctions[declUpper]; exists {
			env.DefineFunction(fn)
		} else if variable, exists := mathVariables[declUpper]; exists {
			env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, variable.Documentation)
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown MATH declaration: %s", decl)}
		}
	}

	return nil
}
