package api

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/stdlib"
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
		// Check if it's a BASKIT (map) object
		if v.ClassName == "BASKIT" {
			return baskitToGoMap(v)
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
	underlying := bukkit.Instance.(*environment.ObjectInstance).NativeData.(stdlib.BukkitSlice)
	result := make([]interface{}, len(underlying))
	for i, val := range underlying {
		goVal, err := ToGoValue(val)
		if err != nil {
			return nil, err
		}
		result[i] = goVal
	}
	return result, nil
}

// baskitToGoMap converts an Objective-LOL BASKIT object to a Go map
func baskitToGoMap(baskit types.ObjectValue) (map[string]interface{}, error) {
	underlying := baskit.Instance.(*environment.ObjectInstance).NativeData.(stdlib.BaskitMap)
	result := make(map[string]interface{})
	for key, val := range underlying {
		goVal, err := ToGoValue(val)
		if err != nil {
			return nil, err
		}
		result[key] = goVal
	}
	return result, nil
}

// objectToGoMap converts an Objective-LOL object to a Go map
func objectToGoMap(obj types.ObjectValue) (map[string]interface{}, error) {
	instance := obj.Instance.(*environment.ObjectInstance)
	result := make(map[string]interface{})
	for key, val := range instance.Variables {
		goVal, err := ToGoValue(val.Value)
		if err != nil {
			return nil, err
		}
		result[key] = goVal
	}
	for key, val := range instance.SharedVariables {
		goVal, err := ToGoValue(val.Value)
		if err != nil {
			return nil, err
		}
		result[key] = goVal
	}
	return result, nil
}

// sliceToArray converts a Go slice/array to an Objective-LOL BUKKIT object
func sliceToArray(rv reflect.Value) (types.Value, error) {
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, NewConversionError(
			fmt.Sprintf("cannot convert Go type %s to BUKKIT", rv.Type()),
			nil,
		)
	}

	// Create a new BUKKIT object
	instance := stdlib.NewBukkitInstance()

	// Populate the BUKKIT object with the slice elements
	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i)
		objElem, err := FromGoValue(elem.Interface())
		if err != nil {
			return nil, err
		}
		instance.NativeData = append(instance.NativeData.(stdlib.BukkitSlice), objElem)
	}
	instance.Variables["SIZ"].Value = types.IntegerValue(rv.Len())

	return types.NewObjectValue(instance, "BUKKIT"), nil
}

// mapToObject converts a Go map to an Objective-LOL BASKIT object
func mapToObject(rv reflect.Value) (types.Value, error) {
	if rv.Kind() != reflect.Map {
		return nil, NewConversionError(
			fmt.Sprintf("cannot convert Go type %s to BUKKIT", rv.Type()),
			nil,
		)
	}

	// Create a new BASKIT object
	instance := stdlib.NewBaskitInstance()

	// Populate the BASKIT object with the map elements
	for _, key := range rv.MapKeys() {
		val := rv.MapIndex(key)
		objKey, err := FromGoValue(key.Interface())
		if err != nil {
			return nil, err
		}
		objVal, err := FromGoValue(val.Interface())
		if err != nil {
			return nil, err
		}
		underlying := instance.NativeData.(stdlib.BaskitMap)
		underlying[objKey.String()] = objVal
	}
	instance.Variables["SIZ"].Value = types.IntegerValue(rv.Len())

	return types.NewObjectValue(instance, "BASKIT"), nil
}

// structToObject converts a Go struct to an Objective-LOL BASKIT object
func structToObject(rv reflect.Value) (types.Value, error) {
	if rv.Kind() != reflect.Struct {
		return nil, NewConversionError(
			fmt.Sprintf("cannot convert Go type %s to BASKIT", rv.Type()),
			nil,
		)
	}

	// Create a new BASKIT object
	instance := stdlib.NewBaskitInstance()

	// Populate the BASKIT object with the struct fields
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		value := rv.Field(i)

		// Convert field name and value to Objective-LOL types
		objKey, err := FromGoValue(field.Name)
		if err != nil {
			return nil, err
		}
		objVal, err := FromGoValue(value.Interface())
		if err != nil {
			return nil, err
		}

		// Add the key-value pair to the BASKIT object
		underlying := instance.NativeData.(stdlib.BaskitMap)
		underlying[objKey.String()] = objVal
	}
	instance.Variables["SIZ"].Value = types.IntegerValue(rv.NumField())

	return types.NewObjectValue(instance, "BASKIT"), nil
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
