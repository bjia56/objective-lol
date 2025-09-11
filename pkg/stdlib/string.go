package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
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
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str := args[0]

					if strVal, ok := str.(environment.StringValue); ok {
						return environment.IntegerValue(len(string(strVal))), nil
					}

					return environment.NOTHIN, fmt.Errorf("LEN: argument is not a string")
				},
			},
			"CONCAT": {
				Name:       "CONCAT",
				ReturnType: "STRIN",
				IsVarargs:  true,
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					var strBuilder strings.Builder
					for _, arg := range args {
						arg, err := arg.Cast("STRIN")
						if err != nil {
							return environment.NOTHIN, fmt.Errorf("CONCAT: all arguments must be strings")
						}
						strBuilder.WriteString(string(arg.(environment.StringValue)))
					}
					return environment.StringValue(strBuilder.String()), nil
				},
			},
			"SUBSTR": {
				Name:       "SUBSTR",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "START", Type: "INTEGR"},
					{Name: "LENGTH", Type: "INTEGR"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("SUBSTR: first argument is not a string")
					}

					start, ok := args[1].(environment.IntegerValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("SUBSTR: second argument is not an integer")
					}

					length, ok := args[2].(environment.IntegerValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("SUBSTR: third argument is not an integer")
					}

					s := string(str)
					startIdx := int(start)
					lengthVal := int(length)

					if startIdx < 0 || startIdx >= len(s) {
						return environment.StringValue(""), nil
					}

					endIdx := startIdx + lengthVal
					if endIdx > len(s) {
						endIdx = len(s)
					}

					return environment.StringValue(s[startIdx:endIdx]), nil
				},
			},
			"TRIM": {
				Name:       "TRIM",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("TRIM: argument is not a string")
					}

					return environment.StringValue(strings.TrimSpace(string(str))), nil
				},
			},
			"LTRIM": {
				Name:       "LTRIM",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("LTRIM: argument is not a string")
					}

					return environment.StringValue(strings.TrimLeftFunc(string(str), func(r rune) bool {
						return r == ' ' || r == '\t' || r == '\n' || r == '\r'
					})), nil
				},
			},
			"RTRIM": {
				Name:       "RTRIM",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("RTRIM: argument is not a string")
					}

					return environment.StringValue(strings.TrimRightFunc(string(str), func(r rune) bool {
						return r == ' ' || r == '\t' || r == '\n' || r == '\r'
					})), nil
				},
			},
			"REPEAT": {
				Name:       "REPEAT",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "COUNT", Type: "INTEGR"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("REPEAT: first argument is not a string")
					}

					count, ok := args[1].(environment.IntegerValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("REPEAT: second argument is not an integer")
					}

					if count < 0 {
						return environment.NOTHIN, fmt.Errorf("REPEAT: count must be non-negative")
					}

					return environment.StringValue(strings.Repeat(string(str), int(count))), nil
				},
			},
			"UPPER": {
				Name:       "UPPER",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("UPPER: argument is not a string")
					}

					return environment.StringValue(strings.ToUpper(string(str))), nil
				},
			},
			"LOWER": {
				Name:       "LOWER",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("LOWER: argument is not a string")
					}

					return environment.StringValue(strings.ToLower(string(str))), nil
				},
			},
			"TITLE": {
				Name:       "TITLE",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("TITLE: argument is not a string")
					}

					return environment.StringValue(strings.Title(string(str))), nil
				},
			},
			"CAPITALIZE": {
				Name:       "CAPITALIZE",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("CAPITALIZE: argument is not a string")
					}

					s := string(str)
					if len(s) == 0 {
						return environment.StringValue(""), nil
					}

					return environment.StringValue(strings.ToUpper(s[:1]) + strings.ToLower(s[1:])), nil
				},
			},
			"SPLIT": {
				Name:       "SPLIT",
				ReturnType: "BUKKIT",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "SEPARATOR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("SPLIT: first argument is not a string")
					}

					separator, ok := args[1].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("SPLIT: second argument is not a string")
					}

					parts := strings.Split(string(str), string(separator))

					// Create a new BUKKIT array with the split parts
					bukkitObj := NewBukkitInstance()
					bukkitSlice := make(BukkitSlice, 0, len(parts))
					for _, part := range parts {
						bukkitSlice = append(bukkitSlice, environment.StringValue(part))
					}
					bukkitObj.NativeData = bukkitSlice

					return bukkitObj, nil
				},
			},
			"REPLACE": {
				Name:       "REPLACE",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "OLD", Type: "STRIN"},
					{Name: "NEW", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("REPLACE: first argument is not a string")
					}

					old, ok := args[1].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("REPLACE: second argument is not a string")
					}

					new, ok := args[2].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("REPLACE: third argument is not a string")
					}

					result := strings.Replace(string(str), string(old), string(new), 1)
					return environment.StringValue(result), nil
				},
			},
			"REPLACE_ALL": {
				Name:       "REPLACE_ALL",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "OLD", Type: "STRIN"},
					{Name: "NEW", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("REPLACE_ALL: first argument is not a string")
					}

					old, ok := args[1].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("REPLACE_ALL: second argument is not a string")
					}

					new, ok := args[2].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("REPLACE_ALL: third argument is not a string")
					}

					result := strings.ReplaceAll(string(str), string(old), string(new))
					return environment.StringValue(result), nil
				},
			},
			"CONTAINS": {
				Name:       "CONTAINS",
				ReturnType: "BOOL",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "SUBSTR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("CONTAINS: first argument is not a string")
					}

					substr, ok := args[1].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("CONTAINS: second argument is not a string")
					}

					result := strings.Contains(string(str), string(substr))
					return environment.BoolValue(result), nil
				},
			},
			"INDEX_OF": {
				Name:       "INDEX_OF",
				ReturnType: "INTEGR",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "SUBSTR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("INDEX_OF: first argument is not a string")
					}

					substr, ok := args[1].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, fmt.Errorf("INDEX_OF: second argument is not a string")
					}

					index := strings.Index(string(str), string(substr))
					return environment.IntegerValue(index), nil
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
