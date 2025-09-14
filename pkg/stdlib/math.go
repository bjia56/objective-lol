package stdlib

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleMATHCategories defines the order that categories should be rendered in documentation
var moduleMATHCategories = []string{
	"mathematical-constants",
	"basic-math",
	"advanced-math",
	"trigonometry",
	"logarithmic",
	"rounding",
}

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
					"The mathematical constant π (pi) ≈ 3.14159.",
					"Represents the ratio of a circle's circumference to its diameter.",
					"",
					"@type {DUBBLE}",
					"@value 3.14159265359",
					"@example Circle area calculation",
					"I HAS A VARIABLE RADIUS TEH DUBBLE ITZ 2.0",
					"I HAS A VARIABLE AREA TEH DUBBLE ITZ PI TIEMZ RADIUS TIEMZ RADIUS",
					"BTW Result: 12.566370614",
					"@example Convert degrees to radians",
					"I HAS A VARIABLE DEGREES TEH DUBBLE ITZ 180.0",
					"I HAS A VARIABLE RADIANS TEH DUBBLE ITZ DEGREES TIEMZ PI DIVIDEZ 180.0",
					"BTW Result: 3.14159265359 (π radians)",
					"@see E, SIN, COS, TAN",
					"@category mathematical-constants",
				},
				Type:     "DUBBLE",
				Value:    environment.DoubleValue(math.Pi),
				IsLocked: true,
				IsPublic: true,
			},
			"E": {
				Name: "E",
				Documentation: []string{
					"Euler's number e ≈ 2.71828.",
					"The base of natural logarithms, fundamental mathematical constant.",
					"",
					"@type {DUBBLE}",
					"@value 2.71828182846",
					"@example Natural logarithm base verification",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG WIT E",
					"BTW Result: 1.0 (ln(e) = 1)",
					"@example Exponential function",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ EXP WIT 2.0",
					"BTW Result: 7.389056099 (e^2)",
					"@see PI, LOG, EXP",
					"@category mathematical-constants",
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
					"",
					"@syntax ABS WIT <number>",
					"@param {DUBBLE} value - The number to get absolute value of",
					"@returns {DUBBLE} The absolute value (always positive or zero)",
					"@example Basic absolute value",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ABS WIT -5.5",
					"BTW Result: 5.5",
					"@example Positive number unchanged",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ABS WIT 42.0",
					"BTW Result: 42.0",
					"@example Zero unchanged",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ABS WIT 0.0",
					"BTW Result: 0.0",
					"@note Works with both positive and negative numbers",
					"@see MAX, MIN",
					"@category basic-math",
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
					"",
					"@syntax MAX WIT <value1> AN WIT <value2>",
					"@param {DUBBLE} a - First number to compare",
					"@param {DUBBLE} b - Second number to compare",
					"@returns {DUBBLE} The larger of the two values",
					"@example Compare positive numbers",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ MAX WIT 10.5 AN WIT 7.2",
					"BTW Result: 10.5",
					"@example Compare negative numbers",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ MAX WIT -3.0 AN WIT -8.0",
					"BTW Result: -3.0",
					"@example Equal values",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ MAX WIT 5.0 AN WIT 5.0",
					"BTW Result: 5.0",
					"@note Returns the first value if both are equal",
					"@see MIN, ABS",
					"@category basic-math",
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
					"",
					"@syntax MIN WIT <value1> AN WIT <value2>",
					"@param {DUBBLE} a - First number to compare",
					"@param {DUBBLE} b - Second number to compare",
					"@returns {DUBBLE} The smaller of the two values",
					"@example Compare positive numbers",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ MIN WIT 10.5 AN WIT 7.2",
					"BTW Result: 7.2",
					"@example Compare negative numbers",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ MIN WIT -3.0 AN WIT -8.0",
					"BTW Result: -8.0",
					"@example Equal values",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ MIN WIT 5.0 AN WIT 5.0",
					"BTW Result: 5.0",
					"@note Returns the first value if both are equal",
					"@see MAX, ABS",
					"@category basic-math",
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
					"",
					"@syntax SQRT WIT <number>",
					"@param {DUBBLE} value - The number to get square root of (must be ≥ 0)",
					"@returns {DUBBLE} The square root of the input value",
					"@example Perfect square",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ SQRT WIT 16.0",
					"BTW Result: 4.0",
					"@example Decimal result",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ SQRT WIT 2.0",
					"BTW Result: 1.4142135623",
					"@example Zero input",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ SQRT WIT 0.0",
					"BTW Result: 0.0",
					"@throws Negative argument error if value < 0",
					"@note Input must be non-negative",
					"@see POW, ABS",
					"@category advanced-math",
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
					"",
					"@syntax POW WIT <base> AN WIT <exponent>",
					"@param {DUBBLE} base - The base number",
					"@param {DUBBLE} exponent - The power to raise the base to",
					"@returns {DUBBLE} The result of base^exponent",
					"@example Integer exponent",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ POW WIT 2.0 AN WIT 3.0",
					"BTW Result: 8.0",
					"@example Fractional exponent (square root)",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ POW WIT 4.0 AN WIT 0.5",
					"BTW Result: 2.0",
					"@example Power of ten",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ POW WIT 10.0 AN WIT 2.0",
					"BTW Result: 100.0",
					"@example Zero exponent",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ POW WIT 5.0 AN WIT 0.0",
					"BTW Result: 1.0",
					"@note Any number to the power of 0 equals 1",
					"@see SQRT, EXP, LOG",
					"@category advanced-math",
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
					"",
					"@syntax SIN WIT <angle_radians>",
					"@param {DUBBLE} value - The angle in radians",
					"@returns {DUBBLE} The sine of the angle (-1 to 1)",
					"@example Sine of 0",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ SIN WIT 0.0",
					"BTW Result: 0.0",
					"@example Sine of π/2 (90 degrees)",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ SIN WIT PI DIVIDEZ 2",
					"BTW Result: 1.0",
					"@example Sine of π (180 degrees)",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ SIN WIT PI",
					"BTW Result: ≈0.0",
					"@note Input must be in radians, not degrees",
					"@note Result is always between -1 and 1",
					"@see COS, TAN, ASIN, PI",
					"@category trigonometry",
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
					"",
					"@syntax COS WIT <angle_radians>",
					"@param {DUBBLE} value - The angle in radians",
					"@returns {DUBBLE} The cosine of the angle (-1 to 1)",
					"@example Cosine of 0",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ COS WIT 0.0",
					"BTW Result: 1.0",
					"@example Cosine of π/2 (90 degrees)",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ COS WIT PI DIVIDEZ 2",
					"BTW Result: ≈0.0",
					"@example Cosine of π (180 degrees)",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ COS WIT PI",
					"BTW Result: -1.0",
					"@note Input must be in radians, not degrees",
					"@note Result is always between -1 and 1",
					"@see SIN, TAN, ACOS, PI",
					"@category trigonometry",
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
					"",
					"@syntax TAN WIT <angle_radians>",
					"@param {DUBBLE} value - The angle in radians",
					"@returns {DUBBLE} The tangent of the angle",
					"@example Tangent of 0",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ TAN WIT 0.0",
					"BTW Result: 0.0",
					"@example Tangent of π/4 (45 degrees)",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ TAN WIT PI DIVIDEZ 4",
					"BTW Result: 1.0",
					"@example Tangent of π (180 degrees)",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ TAN WIT PI",
					"BTW Result: ≈0.0",
					"@note Input must be in radians, not degrees",
					"@note Undefined at π/2 + nπ (90°, 270°, etc.)",
					"@see SIN, COS, ATAN, PI",
					"@category trigonometry",
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
					"",
					"@syntax ASIN WIT <value>",
					"@param {DUBBLE} value - Input value (must be between -1 and 1)",
					"@returns {DUBBLE} The arc sine in radians (-π/2 to π/2)",
					"@example Arc sine of 0",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ASIN WIT 0.0",
					"BTW Result: 0.0",
					"@example Arc sine of 1",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ASIN WIT 1.0",
					"BTW Result: π/2 (≈1.5708)",
					"@example Arc sine of -1",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ASIN WIT -1.0",
					"BTW Result: -π/2 (≈-1.5708)",
					"@throws Input out of range error if value < -1 or value > 1",
					"@note Input must be in range [-1, 1]",
					"@see SIN, ACOS, ATAN, PI",
					"@category trigonometry",
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
					"",
					"@syntax ACOS WIT <value>",
					"@param {DUBBLE} value - Input value (must be between -1 and 1)",
					"@returns {DUBBLE} The arc cosine in radians (0 to π)",
					"@example Arc cosine of 1",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ACOS WIT 1.0",
					"BTW Result: 0.0",
					"@example Arc cosine of 0",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ACOS WIT 0.0",
					"BTW Result: π/2 (≈1.5708)",
					"@example Arc cosine of -1",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ACOS WIT -1.0",
					"BTW Result: π (≈3.1416)",
					"@throws Input out of range error if value < -1 or value > 1",
					"@note Input must be in range [-1, 1]",
					"@see COS, ASIN, ATAN, PI",
					"@category trigonometry",
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
					"",
					"@syntax ATAN WIT <value>",
					"@param {DUBBLE} value - Input value (any real number)",
					"@returns {DUBBLE} The arc tangent in radians (-π/2 to π/2)",
					"@example Arc tangent of 0",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN WIT 0.0",
					"BTW Result: 0.0",
					"@example Arc tangent of 1",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN WIT 1.0",
					"BTW Result: π/4 (≈0.7854)",
					"@example Arc tangent of -1",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN WIT -1.0",
					"BTW Result: -π/4 (≈-0.7854)",
					"@note Input can be any real number",
					"@see TAN, ATAN2, ASIN, ACOS, PI",
					"@category trigonometry",
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
					"",
					"@syntax ATAN2 WIT <y> AN WIT <x>",
					"@param {DUBBLE} y - Y coordinate",
					"@param {DUBBLE} x - X coordinate",
					"@returns {DUBBLE} The angle from x-axis to point (x,y) in radians (-π to π)",
					"@example Point in first quadrant",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN2 WIT 1.0 AN WIT 1.0",
					"BTW Result: π/4 (≈0.7854)",
					"@example Point on positive y-axis",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN2 WIT 1.0 AN WIT 0.0",
					"BTW Result: π/2 (≈1.5708)",
					"@example Point on positive x-axis",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN2 WIT 0.0 AN WIT 1.0",
					"BTW Result: 0.0",
					"@note More robust than ATAN for coordinate conversion",
					"@note Handles all quadrants correctly",
					"@see ATAN, TAN, PI",
					"@category trigonometry",
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
					"",
					"@syntax LOG WIT <number>",
					"@param {DUBBLE} value - The number to get natural logarithm of (must be > 0)",
					"@returns {DUBBLE} The natural logarithm (ln) of the input",
					"@example Natural log of e",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG WIT E",
					"BTW Result: 1.0",
					"@example Natural log of 1",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG WIT 1.0",
					"BTW Result: 0.0",
					"@example Natural log of 10",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG WIT 10.0",
					"BTW Result: 2.302585093",
					"@throws Input must be positive error if value ≤ 0",
					"@note Input must be positive (greater than zero)",
					"@see LOG10, LOG2, EXP, E",
					"@category logarithmic",
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
					"",
					"@syntax LOG10 WIT <number>",
					"@param {DUBBLE} value - The number to get base-10 logarithm of (must be > 0)",
					"@returns {DUBBLE} The base-10 logarithm of the input",
					"@example Log base 10 of 10",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG10 WIT 10.0",
					"BTW Result: 1.0",
					"@example Log base 10 of 100",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG10 WIT 100.0",
					"BTW Result: 2.0",
					"@example Log base 10 of 1",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG10 WIT 1.0",
					"BTW Result: 0.0",
					"@throws Input must be positive error if value ≤ 0",
					"@note Input must be positive (greater than zero)",
					"@note Common logarithm for scientific calculations",
					"@see LOG, LOG2, EXP",
					"@category logarithmic",
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
					"",
					"@syntax LOG2 WIT <number>",
					"@param {DUBBLE} value - The number to get base-2 logarithm of (must be > 0)",
					"@returns {DUBBLE} The base-2 logarithm of the input",
					"@example Log base 2 of 2",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG2 WIT 2.0",
					"BTW Result: 1.0",
					"@example Log base 2 of 8",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG2 WIT 8.0",
					"BTW Result: 3.0",
					"@example Log base 2 of 1",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG2 WIT 1.0",
					"BTW Result: 0.0",
					"@throws Input must be positive error if value ≤ 0",
					"@note Input must be positive (greater than zero)",
					"@note Useful for binary and computer science calculations",
					"@see LOG, LOG10, EXP",
					"@category logarithmic",
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
					"",
					"@syntax EXP WIT <number>",
					"@param {DUBBLE} value - The exponent to raise e to",
					"@returns {DUBBLE} e raised to the power of the input (e^value)",
					"@example e to the power of 1",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ EXP WIT 1.0",
					"BTW Result: 2.718281828 (e)",
					"@example e to the power of 0",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ EXP WIT 0.0",
					"BTW Result: 1.0",
					"@example e to the power of 2",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ EXP WIT 2.0",
					"BTW Result: 7.389056099 (e²)",
					"@note The exponential function, inverse of natural logarithm",
					"@note EXP(LOG(x)) = x for positive x",
					"@see LOG, LOG10, LOG2, E, POW",
					"@category logarithmic",
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
					"",
					"@syntax CEIL WIT <number>",
					"@param {DUBBLE} value - The number to round up",
					"@returns {DUBBLE} The smallest integer ≥ the input value",
					"@example Ceiling of positive decimal",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ CEIL WIT 3.2",
					"BTW Result: 4.0",
					"@example Ceiling of whole number",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ CEIL WIT 3.0",
					"BTW Result: 3.0",
					"@example Ceiling of negative number",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ CEIL WIT -3.7",
					"BTW Result: -3.0",
					"@note Always rounds up to the next whole number",
					"@note For negative numbers, rounds towards zero",
					"@see FLOOR, ROUND, TRUNC",
					"@category rounding",
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
					"",
					"@syntax FLOOR WIT <number>",
					"@param {DUBBLE} value - The number to round down",
					"@returns {DUBBLE} The largest integer ≤ the input value",
					"@example Floor of positive decimal",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ FLOOR WIT 3.7",
					"BTW Result: 3.0",
					"@example Floor of whole number",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ FLOOR WIT 3.0",
					"BTW Result: 3.0",
					"@example Floor of negative number",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ FLOOR WIT -3.2",
					"BTW Result: -4.0",
					"@note Always rounds down to the previous whole number",
					"@note For negative numbers, rounds away from zero",
					"@see CEIL, ROUND, TRUNC",
					"@category rounding",
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
					"",
					"@syntax ROUND WIT <number>",
					"@param {DUBBLE} value - The number to round",
					"@returns {DUBBLE} The rounded integer value",
					"@example Round down",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ROUND WIT 3.4",
					"BTW Result: 3.0",
					"@example Round up",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ROUND WIT 3.6",
					"BTW Result: 4.0",
					"@example Round half up",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ROUND WIT 3.5",
					"BTW Result: 4.0",
					"@example Negative numbers",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ ROUND WIT -3.5",
					"BTW Result: -3.0",
					"@note Uses \"round half up\" strategy for 0.5 values",
					"@see CEIL, FLOOR, TRUNC",
					"@category rounding",
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
					"",
					"@syntax TRUNC WIT <number>",
					"@param {DUBBLE} value - The number to truncate",
					"@returns {DUBBLE} The integer part with fractional part removed",
					"@example Truncate positive decimal",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ TRUNC WIT 3.7",
					"BTW Result: 3.0",
					"@example Truncate negative decimal",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ TRUNC WIT -3.7",
					"BTW Result: -3.0",
					"@example Truncate whole number",
					"I HAS A VARIABLE RESULT TEH DUBBLE ITZ TRUNC WIT 5.0",
					"BTW Result: 5.0",
					"@note Always truncates towards zero",
					"@note Different from FLOOR for negative numbers",
					"@see CEIL, FLOOR, ROUND",
					"@category rounding",
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
