package api

import (
	"encoding/json"
	"fmt"

	"github.com/bjia56/objective-lol/pkg/runtime"
)

// VMCompatibilityShim is a shim to provide compatibility for external
// languages that cannot interact with the standard VM interface through
// Go types. Message passing is done through JSON strings.
type VMCompatibilityShim struct {
	vm *VM
}

// DefineFunction defines a global function with maximum compatibility,
// wrapping arguments and return values as JSON strings.
// An optional id cookie is passed back to the function to identify it.
// jsonArgs is a JSON array string of the arguments.
// The function should return a JSON object string with "result" and "error" fields.
func (shim *VMCompatibilityShim) DefineFunction(id, name string, argc int, function func(id, jsonArgs string) string) error {
	fn := func(args []GoValue) (GoValue, error) {
		if len(args) != argc {
			return WrapAny(nil), fmt.Errorf("expected %d arguments, got %d", argc, len(args))
		}

		// Convert args to JSON array string
		jsonBytes, err := json.Marshal(args)
		if err != nil {
			return WrapAny(nil), fmt.Errorf("error marshaling arguments to JSON: %v", err)
		}

		// Call the provided function
		jsonResult := function(id, string(jsonBytes))

		// Parse JSON result
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(jsonResult), &result); err != nil {
			return WrapAny(nil), fmt.Errorf("error unmarshaling result from JSON: %v", err)
		}
		if errVal, ok := result["error"]; ok && errVal != nil {
			return WrapAny(nil), runtime.Exception{Message: fmt.Sprintf("%v", errVal)}
		}
		if resultVal, ok := result["result"]; ok {
			return WrapAny(resultVal), nil
		}
		return WrapAny(nil), nil
	}

	return shim.vm.DefineFunction(name, argc, fn)
}
