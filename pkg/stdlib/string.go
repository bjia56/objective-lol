package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// Global STRING function definitions - created once and reused
var stringFunctionsOnce = sync.Once{}
var stringFunctions map[string]*environment.Function

func getStringFunctions() map[string]*environment.Function {
	stringFunctionsOnce.Do(func() {
		stringFunctions = map[string]*environment.Function{
			"LEN": {
				Name:       "LEN",
				ReturnType: "INTEGR",
				Parameters: []environment.Parameter{{Name: "STR", Type: "STRIN"}},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					str := args[0]

					if strVal, ok := str.(types.StringValue); ok {
						return types.IntegerValue(len(string(strVal))), nil
					}

					return types.NOTHIN, fmt.Errorf("LEN: argument is not a string")
				},
			},
			"CONCAT": {
				Name:       "CONCAT",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR1", Type: "STRIN"},
					{Name: "STR2", Type: "STRIN"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					str1, ok1 := args[0].(types.StringValue)
					str2, ok2 := args[1].(types.StringValue)

					if !ok1 {
						return types.NOTHIN, fmt.Errorf("CONCAT: first argument is not a string")
					}
					if !ok2 {
						return types.NOTHIN, fmt.Errorf("CONCAT: second argument is not a string")
					}

					return types.StringValue(string(str1) + string(str2)), nil
				},
			},
		}
	})
	return stringFunctions
}

// RegisterSTRINGInEnv registers STRING functions in the given environment
// declarations: empty slice means import all, otherwise import only specified functions
func RegisterSTRINGInEnv(env *environment.Environment, declarations ...string) error {
	stringFunctions := getStringFunctions()

	// If declarations is empty, import all functions
	if len(declarations) == 0 {
		for _, fn := range stringFunctions {
			env.DefineFunction(fn)
		}
		return nil
	}

	// Otherwise, import only specified functions
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if fn, exists := stringFunctions[declUpper]; exists {
			env.DefineFunction(fn)
		} else {
			return fmt.Errorf("unknown STRING function: %s", decl)
		}
	}

	return nil
}
