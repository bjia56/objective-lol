package api

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/bjia56/objective-lol/pkg/types"
)

// ToGoValue converts an Objective-LOL value to a Go value
func ToGoValue(val types.Value) (interface{}, error) {
	if val == nil || val.IsNothing() {
		return nil, nil
	}

	switch v := val.(type) {
	case types.IntegerValue:
		return int64(v), nil
	case types.DoubleValue:
		return float64(v), nil
	case types.StringValue:
		return string(v), nil
	case types.BoolValue:
		return bool(v), nil
	case types.ObjectValue:
		// Check if it's a BUKKIT (array) object
		if v.ClassName == "BUKKIT" {
			return bukkitToGoSlice(v)
		}
		return objectToGoMap(v)
	default:
		return nil, NewConversionError(
			fmt.Sprintf("cannot convert Objective-LOL type %s to Go value", val.Type()),
			nil,
		)
	}
}

// FromGoValue converts a Go value to an Objective-LOL value
func FromGoValue(val interface{}) (types.Value, error) {
	if val == nil {
		return types.NOTHIN, nil
	}

	rv := reflect.ValueOf(val)
	
	// Handle pointers
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return types.NOTHIN, nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return types.YEZ, nil
		}
		return types.NO, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return types.IntegerValue(rv.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return types.IntegerValue(rv.Uint()), nil

	case reflect.Float32, reflect.Float64:
		return types.DoubleValue(rv.Float()), nil

	case reflect.String:
		return types.StringValue(rv.String()), nil

	case reflect.Slice, reflect.Array:
		return sliceToArray(rv)

	case reflect.Map:
		return mapToObject(rv)

	case reflect.Struct:
		return structToObject(rv)

	default:
		return nil, NewConversionError(
			fmt.Sprintf("cannot convert Go type %s to Objective-LOL value", rv.Type()),
			nil,
		)
	}
}

// bukkitToGoSlice converts an Objective-LOL BUKKIT object to a Go slice
func bukkitToGoSlice(bukkit types.ObjectValue) ([]interface{}, error) {
	// For BUKKIT objects, we'd need to access the underlying slice data
	// This is a simplified implementation - in practice, you'd need to
	// interact with the BUKKIT's native data through method calls
	return []interface{}{}, nil
}

// objectToGoMap converts an Objective-LOL object to a Go map
func objectToGoMap(obj types.ObjectValue) (map[string]interface{}, error) {
	// This is a simplified implementation - in practice, you'd need to
	// access the object's properties through the environment system
	return map[string]interface{}{
		"__type__": obj.ClassName,
	}, nil
}

// sliceToArray converts a Go slice/array to an Objective-LOL BUKKIT object
func sliceToArray(rv reflect.Value) (types.Value, error) {
	// This is a simplified implementation - in practice, you'd need to
	// create a BUKKIT object through the environment system
	// For now, we'll just return NOTHIN to indicate we can't convert yet
	return types.NOTHIN, NewConversionError(
		"slice to BUKKIT conversion not yet implemented",
		nil,
	)
}

// mapToObject converts a Go map to an Objective-LOL object
func mapToObject(rv reflect.Value) (types.Value, error) {
	// This would need to create a dynamic object - simplified for now
	return types.NOTHIN, NewConversionError(
		"map to object conversion not yet implemented",
		nil,
	)
}

// structToObject converts a Go struct to an Objective-LOL object
func structToObject(rv reflect.Value) (types.Value, error) {
	// This would need to create a dynamic object - simplified for now
	return types.NOTHIN, NewConversionError(
		"struct to object conversion not yet implemented",
		nil,
	)
}

// ConvertArguments converts a slice of Go values to Objective-LOL values
func ConvertArguments(args []interface{}) ([]types.Value, error) {
	result := make([]types.Value, len(args))
	for i, arg := range args {
		val, err := FromGoValue(arg)
		if err != nil {
			return nil, NewConversionError(
				fmt.Sprintf("error converting argument at index %d", i),
				err,
			)
		}
		result[i] = val
	}
	return result, nil
}

// ParseGoValue attempts to parse a string into a Go value of the appropriate type
func ParseGoValue(str string) interface{} {
	// Try integer
	if val, err := strconv.ParseInt(str, 10, 64); err == nil {
		return val
	}
	
	// Try float
	if val, err := strconv.ParseFloat(str, 64); err == nil {
		return val
	}
	
	// Try boolean
	if val, err := strconv.ParseBool(str); err == nil {
		return val
	}
	
	// Default to string
	return str
}