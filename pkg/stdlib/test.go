package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// Global TEST function definitions - created once and reused
var testFunctionsOnce = sync.Once{}
var testFunctions map[string]*environment.Function

func getTestFunctions() map[string]*environment.Function {
	testFunctionsOnce.Do(func() {
		testFunctions = map[string]*environment.Function{
			"ASSERT": {
				Name: "ASSERT",
				Documentation: []string{
					"Asserts that a condition is truthy, throwing an exception if NO.",
					"Accepts any type and evaluates truthiness according to Objective-LOL rules.",
				},
				Parameters: []environment.Parameter{{Name: "CONDITION", Type: ""}}, // Accept any type
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					condition := args[0]

					// Check if condition is truthy
					if !condition.ToBool() {
						return environment.NOTHIN, runtime.Exception{Message: "Assertion failed"}
					}

					return environment.NOTHIN, nil
				},
			},
		}
	})
	return testFunctions
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
			return runtime.Exception{Message: fmt.Sprintf("unknown TEST declaration: %s", decl)}
		}
	}

	return nil
}
