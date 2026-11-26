package stdlib

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleJSONCategories defines the order that categories should be rendered in documentation
var moduleJSONCategories = []string{
	"json-serialization",
	"json-deserialization",
}

// Global JSON function definitions - created once and reused
var jsonFunctionsOnce = sync.Once{}
var jsonFunctions map[string]*environment.Function

func getJSONFunctions() map[string]*environment.Function {
	jsonFunctionsOnce.Do(func() {
		jsonFunctions = map[string]*environment.Function{
			"TO_JSON": {
				Name: "TO_JSON",
				Documentation: []string{
					"Serializes any Objective-LOL value to a JSON string.",
					"NOTHIN becomes null, BOOL becomes true/false, INTEGR and DUBBLE become numbers,",
					"STRIN becomes a JSON string, BUKKIT becomes a JSON array, and BASKIT becomes",
					"a JSON object. Other object types are not supported and will raise an exception.",
					"",
					"@param VALUE - The value to serialize",
					"@returns {STRIN} The JSON representation",
					"@category json-serialization",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					native, err := valueToNative(args[0])
					if err != nil {
						return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("TO_JSON: %v", err)}
					}
					data, err := json.Marshal(native)
					if err != nil {
						return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("TO_JSON: marshal error: %v", err)}
					}
					return environment.StringValue(data), nil
				},
			},
			"FROM_JSON": {
				Name: "FROM_JSON",
				Documentation: []string{
					"Parses a JSON string and returns the corresponding Objective-LOL value.",
					"JSON null becomes NOTHIN, booleans become BOOL, numbers become DUBBLE,",
					"strings become STRIN, arrays become BUKKIT, and objects become BASKIT.",
					"",
					"@param STR - The JSON string to parse",
					"@returns The parsed value (NOTHIN, BOOL, DUBBLE, STRIN, BUKKIT, or BASKIT)",
					"@category json-deserialization",
				},
				ReturnType: "",
				Parameters: []environment.Parameter{{Name: "STR", Type: "STRIN"}},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					str, ok := args[0].(environment.StringValue)
					if !ok {
						return environment.NOTHIN, runtime.Exception{Message: "FROM_JSON: argument is not a string"}
					}

					var raw interface{}
					if err := json.Unmarshal([]byte(str), &raw); err != nil {
						return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("FROM_JSON: parse error: %v", err)}
					}

					return nativeToValue(raw)
				},
			},
		}
	})
	return jsonFunctions
}

// valueToNative converts an environment.Value to a Go native type suitable for json.Marshal.
func valueToNative(v environment.Value) (interface{}, error) {
	switch val := v.(type) {
	case environment.NothingValue:
		return nil, nil
	case environment.BoolValue:
		return bool(val), nil
	case environment.IntegerValue:
		return int64(val), nil
	case environment.DoubleValue:
		return float64(val), nil
	case environment.StringValue:
		return string(val), nil
	case *environment.ObjectInstance:
		switch val.Class.Name {
		case "BUKKIT":
			slice, ok := val.NativeData.(BukkitSlice)
			if !ok {
				return nil, fmt.Errorf("BUKKIT has unexpected internal data")
			}
			result := make([]interface{}, len(slice))
			for i, elem := range slice {
				native, err := valueToNative(elem)
				if err != nil {
					return nil, fmt.Errorf("BUKKIT element %d: %v", i, err)
				}
				result[i] = native
			}
			return result, nil
		case "BASKIT":
			m, ok := val.NativeData.(BaskitMap)
			if !ok {
				return nil, fmt.Errorf("BASKIT has unexpected internal data")
			}
			result := make(map[string]interface{}, len(m))
			for k, elem := range m {
				native, err := valueToNative(elem)
				if err != nil {
					return nil, fmt.Errorf("BASKIT key %q: %v", k, err)
				}
				result[k] = native
			}
			return result, nil
		default:
			return nil, fmt.Errorf("cannot serialize object of type %s to JSON", val.Class.Name)
		}
	default:
		return nil, fmt.Errorf("cannot serialize value of type %T to JSON", v)
	}
}

// nativeToValue converts a Go value produced by json.Unmarshal into an environment.Value.
func nativeToValue(v interface{}) (environment.Value, error) {
	if v == nil {
		return environment.NOTHIN, nil
	}
	switch val := v.(type) {
	case bool:
		return environment.BoolValue(val), nil
	case float64:
		return environment.DoubleValue(val), nil
	case string:
		return environment.StringValue(val), nil
	case []interface{}:
		bukkit := NewBukkitInstance()
		slice := make(BukkitSlice, 0, len(val))
		for i, elem := range val {
			lolVal, err := nativeToValue(elem)
			if err != nil {
				return environment.NOTHIN, fmt.Errorf("array element %d: %v", i, err)
			}
			slice = append(slice, lolVal)
		}
		bukkit.NativeData = slice
		return bukkit, nil
	case map[string]interface{}:
		baskit := NewBaskitInstance()
		m := make(BaskitMap, len(val))
		for k, elem := range val {
			lolVal, err := nativeToValue(elem)
			if err != nil {
				return environment.NOTHIN, fmt.Errorf("object key %q: %v", k, err)
			}
			m[k] = lolVal
		}
		baskit.NativeData = m
		return baskit, nil
	default:
		return environment.NOTHIN, fmt.Errorf("unexpected JSON value type %T", v)
	}
}

// RegisterJSONInEnv registers JSON functions in the given environment.
// declarations: empty slice means import all, otherwise import only specified functions.
func RegisterJSONInEnv(env *environment.Environment, declarations ...string) error {
	fns := getJSONFunctions()

	if len(declarations) == 0 {
		for _, fn := range fns {
			env.DefineFunction(fn)
		}
		return nil
	}

	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if fn, exists := fns[declUpper]; exists {
			env.DefineFunction(fn)
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown JSON declaration: %s", decl)}
		}
	}

	return nil
}
