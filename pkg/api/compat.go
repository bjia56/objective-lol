package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bjia56/objective-lol/pkg/runtime"
)

// VMCompatibilityShim is a shim to provide compatibility for external
// languages that cannot interact with the standard VM interface through
// Go types. Message passing is done through JSON strings.
type VMCompatibilityShim struct {
	vm *VM
}

type CompatibilityCallback func(id, jsonArgs string) string

// DefineFunction defines a global function with maximum compatibility,
// wrapping arguments and return values as JSON strings.
// An optional id cookie is passed back to the function to identify it.
// jsonArgs is a JSON array string of the arguments.
// The function should return a JSON object string with "result" and "error" fields.
func (shim *VMCompatibilityShim) DefineFunction(id, name string, argc int, function CompatibilityCallback) error {
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

func (shim *VMCompatibilityShim) BuildNewClassVariableWithGetter(variable *ClassVariable, getterID string, getter CompatibilityCallback) *ClassVariable {
	// Wrap the getter to match the expected signature
	var wrappedGetter func(this GoValue) (GoValue, error)
	if getter != nil {
		wrappedGetter = func(this GoValue) (GoValue, error) {
			args := []interface{}{
				this.ID(),
			}
			jsonBytes, err := json.Marshal(args)
			if err != nil {
				return WrapAny(nil), fmt.Errorf("error marshaling getter argument to JSON: %v", err)
			}
			jsonResult := getter(getterID, string(jsonBytes))
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(jsonResult), &result); err != nil {
				return WrapAny(nil), fmt.Errorf("error unmarshaling getter result from JSON: %v", err)
			}
			if errVal, ok := result["error"]; ok && errVal != nil {
				return WrapAny(nil), runtime.Exception{Message: fmt.Sprintf("%v", errVal)}
			}
			if resultVal, ok := result["result"]; ok {
				return WrapAny(resultVal), nil
			}
			return WrapAny(nil), nil
		}
	}

	variable.Getter = wrappedGetter
	return variable
}

func (shim *VMCompatibilityShim) BuildNewClassVariableWithSetter(variable *ClassVariable, setterID string, setter CompatibilityCallback) *ClassVariable {
	// Wrap the setter to match the expected signature
	var wrappedSetter func(this GoValue, value GoValue) error
	if setter != nil {
		wrappedSetter = func(this GoValue, value GoValue) error {
			args := []interface{}{
				this.ID(),
				value,
			}
			jsonBytes, err := json.Marshal(args)
			if err != nil {
				return fmt.Errorf("error marshaling setter argument to JSON: %v", err)
			}
			jsonResult := setter(setterID, string(jsonBytes))
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(jsonResult), &result); err != nil {
				return fmt.Errorf("error unmarshaling setter result from JSON: %v", err)
			}
			if errVal, ok := result["error"]; ok && errVal != nil {
				return runtime.Exception{Message: fmt.Sprintf("%v", errVal)}
			}
			return nil
		}
	}

	variable.Setter = wrappedSetter
	return variable
}

func (shim *VMCompatibilityShim) BuildNewClassMethod(method *ClassMethod, id string, function CompatibilityCallback) *ClassMethod {
	// Wrap the function to match the expected signature
	wrappedFunction := func(this GoValue, args []GoValue) (GoValue, error) {
		// Convert args to JSON array string
		argsList := []interface{}{this.ID()}
		for _, arg := range args {
			argsList = append(argsList, arg)
		}
		jsonBytes, err := json.Marshal(argsList)
		if err != nil {
			return WrapAny(nil), fmt.Errorf("error marshaling method arguments to JSON: %v", err)
		}

		// Call the provided function
		jsonResult := function(id, string(jsonBytes))

		// Parse JSON result
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(jsonResult), &result); err != nil {
			return WrapAny(nil), fmt.Errorf("error unmarshaling method result from JSON: %v", err)
		}
		if errVal, ok := result["error"]; ok && errVal != nil {
			return WrapAny(nil), runtime.Exception{Message: fmt.Sprintf("%v", errVal)}
		}
		if resultVal, ok := result["result"]; ok {
			return WrapAny(resultVal), nil
		}
		return WrapAny(nil), nil
	}

	method.Function = wrappedFunction
	return method
}

func (shim *VMCompatibilityShim) BuildNewUnknownFunctionHandler(id string, function CompatibilityCallback) UnknownFunctionHandler {
	return UnknownFunctionHandler{
		Handler: func(this GoValue, functionName string, fromContext string, args []GoValue) (GoValue, error) {
			// Convert args to JSON array string
			argsList := []interface{}{this.ID(), functionName, fromContext}
			for _, arg := range args {
				argsList = append(argsList, arg)
			}
			jsonBytes, err := json.Marshal(argsList)
			if err != nil {
				return WrapAny(nil), fmt.Errorf("error marshaling unknown function arguments to JSON: %v", err)
			}

			// Call the provided function
			jsonResult := function(id, string(jsonBytes))

			// Parse JSON result
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(jsonResult), &result); err != nil {
				return WrapAny(nil), fmt.Errorf("error unmarshaling unknown function result from JSON: %v", err)
			}
			if errVal, ok := result["error"]; ok && errVal != nil {
				return WrapAny(nil), runtime.Exception{Message: fmt.Sprintf("%v", errVal)}
			}
			if resultVal, ok := result["result"]; ok {
				return WrapAny(resultVal), nil
			}
			return WrapAny(nil), nil
		},
	}
}

func (shim *VMCompatibilityShim) IsClassDefined(name string) bool {
	cls, err := shim.vm.interpreter.GetEnvironment().GetClass(strings.ToUpper(name))
	return err == nil && cls != nil
}

func (shim *VMCompatibilityShim) LookupObject(id string) (GoValue, error) {
	if obj, err := LookupObject(id); err == nil {
		return WrapObject(obj), nil
	}
	return WrapAny(nil), fmt.Errorf("error in LookupObject: no object found with id %s", id)
}

func (shim *VMCompatibilityShim) GetObjectMRO(id string) ([]string, error) {
	if obj, err := LookupObject(id); err == nil {
		return obj.Class.MRO, nil
	}
	return nil, fmt.Errorf("error in GetObjectMRO: no object found with id %s", id)
}

func (shim *VMCompatibilityShim) GetObjectImmediateFunctions(id string) ([]string, error) {
	if obj, err := LookupObject(id); err == nil {
		funcs := []string{}
		for fname := range obj.Class.PublicFunctions {
			funcs = append(funcs, fname)
		}
		for fname := range obj.Class.PrivateFunctions {
			funcs = append(funcs, fname)
		}
		for fname := range obj.Class.SharedFunctions {
			funcs = append(funcs, fname)
		}
		return funcs, nil
	}
	return nil, fmt.Errorf("error in GetObjectImmediateFunctions: no object found with id %s", id)
}

func (shim *VMCompatibilityShim) GetObjectImmediateVariables(id string) ([]string, error) {
	if obj, err := LookupObject(id); err == nil {
		vars := []string{}
		for vname := range obj.Variables {
			vars = append(vars, vname)
		}
		for vname := range obj.SharedVariables {
			vars = append(vars, vname)
		}
		return vars, nil
	}
	return nil, fmt.Errorf("error in GetObjectImmediateVariables: no object found with id %s", id)
}

func (shim *VMCompatibilityShim) AddVariableToObject(id string, variable *ClassVariable) error {
	if obj, err := LookupObject(id); err == nil {
		memberVariable, err := convertClassVariable(variable.Name, variable, true)
		if err != nil {
			return fmt.Errorf("error converting ClassVariable: %v", err)
		}
		obj.Variables[variable.Name] = memberVariable
		return nil
	}
	return fmt.Errorf("error in AddVariableToObject: no object found with id %s", id)
}
