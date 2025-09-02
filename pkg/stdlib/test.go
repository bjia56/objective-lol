package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// Global TEST function definitions - created once and reused
var testFunctionsOnce = sync.Once{}
var testFunctions map[string]*environment.Function

func getTestFunctions() map[string]*environment.Function {
	testFunctionsOnce.Do(func() {
		testFunctions = map[string]*environment.Function{
			"ASSERT": {
				Name:       "ASSERT",
				Parameters: []environment.Parameter{{Name: "CONDITION", Type: ""}}, // Accept any type
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					condition := args[0]

					// Check if condition is truthy
					if !isTruthy(condition) {
						return types.NOTHIN, ast.Exception{Message: "Assertion failed"}
					}

					return types.NOTHIN, nil
				},
			},
		}
	})
	return testFunctions
}

// isTruthy determines if a value is considered truthy in Objective-LOL
func isTruthy(value types.Value) bool {
	switch v := value.(type) {
	case types.BoolValue:
		return bool(v)
	case types.IntegerValue:
		return int64(v) != 0
	case types.DoubleValue:
		return float64(v) != 0.0
	case types.StringValue:
		return string(v) != ""
	case types.NothingValue:
		return false
	default:
		return true // Objects and other types are considered truthy
	}
}

// RegisterTESTInEnv registers TEST functions in the given environment
// declarations: empty slice means import all, otherwise import only specified functions
func RegisterTESTInEnv(env *environment.Environment, declarations ...string) error {
	testFunctions := getTestFunctions()

	// If declarations is empty, import all functions
	if len(declarations) == 0 {
		for _, fn := range testFunctions {
			env.DefineFunction(fn)
		}
		return nil
	}

	// Otherwise, import only specified functions
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if fn, exists := testFunctions[declUpper]; exists {
			env.DefineFunction(fn)
		} else {
			return fmt.Errorf("unknown TEST function: %s", decl)
		}
	}

	return nil
}
