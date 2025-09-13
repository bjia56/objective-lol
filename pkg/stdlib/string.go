package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// Global STRING function definitions - created once and reused
var stringFunctionsOnce = sync.Once{}
var stringFunctions map[string]*environment.Function

func getStringFunctions() map[string]*environment.Function {
	stringFunctionsOnce.Do(func() {
		stringFunctions = map[string]*environment.Function{
			"LEN": {
				Name: "LEN",
				Documentation: []string{
					"Returns the length of a STRIN in characters.",
					"Counts the number of UTF-8 characters in the STRIN.",
				},
				ReturnType: "INTEGR",
				Parameters: []environment.Parameter{{Name: "STR", Type: "STRIN"}},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str := args[0]

					if strVal, ok := str.(environment.StringValue); ok {
						return environment.IntegerValue(len(string(strVal))), nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "LEN: argument is not a string"}
				},
			},
			"CONCAT": {
				Name: "CONCAT",
				Documentation: []string{
					"Concatenates multiple values into a single STRIN.",
				},
				ReturnType: "STRIN",
				IsVarargs:  true,
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					var strBuilder strings.Builder
					for _, arg := range args {
						arg, err := arg.Cast("STRIN")
						if err != nil {
							return environment.NOTHIN, runtime.Exception{Message: "CONCAT: all arguments must be strings"}
						}
						strBuilder.WriteString(string(arg.(environment.StringValue)))
					}
					return environment.StringValue(strBuilder.String()), nil
				},
			},
			"SUBSTR": {
				Name: "SUBSTR",
				Documentation: []string{
					"Extracts a substring from a STRIN starting at the given position.",
					"Returns substring from START index for LENGTH characters. Bounds are checked.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "START", Type: "INTEGR"},
					{Name: "LENGTH", Type: "INTEGR"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "SUBSTR: first argument is not a string"}
					}

					start, ok := args[1].(environment.IntegerValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "SUBSTR: second argument is not an integer"}
					}

					length, ok := args[2].(environment.IntegerValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "SUBSTR: third argument is not an integer"}
					}

					s := string(str)
					startIdx := int(start)
					lengthVal := int(length)

					if startIdx < 0 || startIdx >= len(s) {
						return environment.StringValue(""), nil
					}

					endIdx := min(startIdx+lengthVal, len(s))

					return environment.StringValue(s[startIdx:endIdx]), nil
				},
			},
			"TRIM": {
				Name: "TRIM",
				Documentation: []string{
					"Removes whitespace from both ends of a STRIN.",
					"Trims spaces, tabs, newlines, and carriage returns.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "TRIM: argument is not a string"}
					}

					return environment.StringValue(strings.TrimSpace(string(str))), nil
				},
			},
			"LTRIM": {
				Name: "LTRIM",
				Documentation: []string{
					"Removes whitespace from the left end of a STRIN.",
					"Trims spaces, tabs, newlines, and carriage returns.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "LTRIM: argument is not a string"}
					}

					return environment.StringValue(strings.TrimLeftFunc(string(str), func(r rune) bool {
						return r == ' ' || r == '\t' || r == '\n' || r == '\r'
					})), nil
				},
			},
			"RTRIM": {
				Name: "RTRIM",
				Documentation: []string{
					"Removes whitespace from the right end of a STRIN.",
					"Trims spaces, tabs, newlines, and carriage returns.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "RTRIM: argument is not a string"}
					}

					return environment.StringValue(strings.TrimRightFunc(string(str), func(r rune) bool {
						return r == ' ' || r == '\t' || r == '\n' || r == '\r'
					})), nil
				},
			},
			"REPEAT": {
				Name: "REPEAT",
				Documentation: []string{
					"Repeats a STRIN a specified number of times.",
					"Returns a new STRIN consisting of the original STRIN repeated COUNT times.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "COUNT", Type: "INTEGR"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "REPEAT: first argument is not a string"}
					}

					count, ok := args[1].(environment.IntegerValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "REPEAT: second argument is not an integer"}
					}

					if count < 0 {
						return environment.NOTHIN, runtime.Exception{Message: "REPEAT: count must be non-negative"}
					}

					return environment.StringValue(strings.Repeat(string(str), int(count))), nil
				},
			},
			"UPPER": {
				Name: "UPPER",
				Documentation: []string{
					"Converts all characters in a STRIN to uppercase.",
					"Returns a new STRIN with all letters converted to upper case.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "UPPER: argument is not a string"}
					}

					return environment.StringValue(strings.ToUpper(string(str))), nil
				},
			},
			"LOWER": {
				Name: "LOWER",
				Documentation: []string{
					"Converts all characters in a STRIN to lowercase.",
					"Returns a new STRIN with all letters converted to lower case.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "LOWER: argument is not a string"}
					}

					return environment.StringValue(strings.ToLower(string(str))), nil
				},
			},
			"TITLE": {
				Name: "TITLE",
				Documentation: []string{
					"Converts the first character of each word to uppercase.",
					"Returns a new STRIN with the first letter of each word capitalized.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "TITLE: argument is not a string"}
					}

					return environment.StringValue(strings.Title(string(str))), nil
				},
			},
			"CAPITALIZE": {
				Name: "CAPITALIZE",
				Documentation: []string{
					"Capitalizes the first character of a STRIN and makes the rest lowercase.",
					"Returns a new STRIN with the first letter capitalized and the rest in lower case.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "CAPITALIZE: argument is not a string"}
					}

					s := string(str)
					if len(s) == 0 {
						return environment.StringValue(""), nil
					}

					return environment.StringValue(strings.ToUpper(s[:1]) + strings.ToLower(s[1:])), nil
				},
			},
			"SPLIT": {
				Name: "SPLIT",
				Documentation: []string{
					"Splits a STRIN into a BUKKIT array using the specified separator.",
					"Returns array of substrings divided by the separator string.",
				},
				ReturnType: "BUKKIT",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "SEPARATOR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "SPLIT: first argument is not a string"}
					}

					separator, ok := args[1].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "SPLIT: second argument is not a string"}
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
				Name: "REPLACE",
				Documentation: []string{
					"Replaces the first occurrence of OLD substring with NEW substring in STR.",
					"Returns a new STRIN with the first occurrence of OLD replaced by NEW.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "OLD", Type: "STRIN"},
					{Name: "NEW", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "REPLACE: first argument is not a string"}
					}

					old, ok := args[1].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "REPLACE: second argument is not a string"}
					}

					new, ok := args[2].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "REPLACE: third argument is not a string"}
					}

					result := strings.Replace(string(str), string(old), string(new), 1)
					return environment.StringValue(result), nil
				},
			},
			"REPLACE_ALL": {
				Name: "REPLACE_ALL",
				Documentation: []string{
					"Replaces all occurrences of OLD substring with NEW substring in STR.",
					"Returns a new STRIN with all occurrences of OLD replaced by NEW.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "OLD", Type: "STRIN"},
					{Name: "NEW", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "REPLACE_ALL: first argument is not a string"}
					}

					old, ok := args[1].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "REPLACE_ALL: second argument is not a string"}
					}

					new, ok := args[2].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "REPLACE_ALL: third argument is not a string"}
					}

					result := strings.ReplaceAll(string(str), string(old), string(new))
					return environment.StringValue(result), nil
				},
			},
			"CONTAINS": {
				Name: "CONTAINS",
				Documentation: []string{
					"Checks if STR contains the substring SUBSTR.",
					"Returns TRUE if SUBSTR is found within STR, otherwise FALSE.",
				},
				ReturnType: "BOOL",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "SUBSTR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "CONTAINS: first argument is not a string"}
					}

					substr, ok := args[1].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "CONTAINS: second argument is not a string"}
					}

					result := strings.Contains(string(str), string(substr))
					return environment.BoolValue(result), nil
				},
			},
			"INDEX_OF": {
				Name: "INDEX_OF",
				Documentation: []string{
					"Finds the index of the first occurrence of SUBSTR in STR.",
					"Returns the zero-based index of SUBSTR in STR, or -1 if not found.",
				},
				ReturnType: "INTEGR",
				Parameters: []environment.Parameter{
					{Name: "STR", Type: "STRIN"},
					{Name: "SUBSTR", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "INDEX_OF: first argument is not a string"}
					}

					substr, ok := args[1].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "INDEX_OF: second argument is not a string"}
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
			return runtime.Exception{Message: fmt.Sprintf("unknown STRING declaration: %s", decl)}
		}
	}

	return nil
}
