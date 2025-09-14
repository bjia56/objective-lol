package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleTestCategories defines the order that categories should be rendered in documentation
var moduleTestCategories = []string{
	"assertions",
}

// Global TEST function definitions - created once and reused
var testFunctionsOnce = sync.Once{}
var testFunctions map[string]*environment.Function

func getTestFunctions() map[string]*environment.Function {
	testFunctionsOnce.Do(func() {
		testFunctions = map[string]*environment.Function{
			"ASSERT": {
				Name: "ASSERT",
				Documentation: []string{
					"Asserts that a condition is truthy, throwing an exception if the condition evaluates to NO.",
					"Accepts any type and evaluates truthiness according to Objective-LOL truthiness rules.",
					"",
					"@syntax ASSERT WIT <condition>",
					"@param {} condition - Any value to test for truthiness",
					"@returns {NOTHIN}",
					"@throws Exception when condition is falsy",
					"@example Basic assertion",
					"ASSERT WIT YEZ",
					"SAYZ WIT \"Test passed!\"",
					"@example Assert with variables",
					"I HAS A VARIABLE COUNT TEH NUMBR ITZ 5",
					"ASSERT WIT COUNT",
					"@example Assert comparison result",
					"ASSERT WIT 2 SAEM AS 2",
					"@note Truthiness: NO, 0, 0.0, \"\", empty arrays, NOTHIN are falsy",
					"@note All other values are truthy",
					"@note Throws \"Assertion failed\" when condition is falsy",
					"@category assertions",
				},
				Parameters: []environment.Parameter{{Name: "condition", Type: ""}}, // Accept any type
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
