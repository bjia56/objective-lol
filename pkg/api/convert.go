package api

import (
	"fmt"
	"reflect"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/stdlib"
)

// ToGoValue converts an Objective-LOL value to a Go value
func ToGoValue(val environment.Value) (GoValue, error) {
	if val == nil || val.IsNothing() {
		return WrapAny(nil), nil
	}

	switch v := val.(type) {
	case environment.IntegerValue:
		return WrapAny(int64(v)), nil
	case environment.DoubleValue:
		return WrapAny(float64(v)), nil
	case environment.StringValue:
		return WrapAny(string(v)), nil
	case environment.BoolValue:
		return WrapAny(bool(v)), nil
	case *environment.ObjectInstance:
		// Check if it's a BUKKIT (array) object
		if v.Class.Name == "BUKKIT" {
			bukkitSlice, err := bukkitToGoSlice(v)
			if err != nil {
				return WrapAny(nil), err
			}
			return WrapAny(bukkitSlice), nil
		}
		// Check if it's a BASKIT (map) object
		if v.Class.Name == "BASKIT" {
			baskitMap, err := baskitToGoMap(v)
			if err != nil {
				return WrapAny(nil), err
			}
			return WrapAny(baskitMap), nil
		}
		return WrapObject(v), nil
	default:
		return WrapAny(nil), NewConversionError(
			fmt.Sprintf("cannot convert Objective-LOL type %s to Go value", val.Type()),
			nil,
		)
	}
}

// FromGoValue converts a Go value to an Objective-LOL value
func FromGoValue(val GoValue) (environment.Value, error) {
	if val.Get() == nil {
		return environment.NOTHIN, nil
	}

	if obj, ok := val.Get().(*environment.ObjectInstance); ok {
		return obj, nil
	}

	rv := reflect.ValueOf(val.Get())

	// Handle pointers
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return environment.NOTHIN, nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return environment.YEZ, nil
		}
		return environment.NO, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return environment.IntegerValue(rv.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return environment.IntegerValue(rv.Uint()), nil

	case reflect.Float32, reflect.Float64:
		return environment.DoubleValue(rv.Float()), nil

	case reflect.String:
		return environment.StringValue(rv.String()), nil

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
func bukkitToGoSlice(bukkit *environment.ObjectInstance) ([]GoValue, error) {
	underlying := bukkit.NativeData.(stdlib.BukkitSlice)
	result := make([]GoValue, len(underlying))
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
func baskitToGoMap(baskit *environment.ObjectInstance) (map[string]GoValue, error) {
	underlying := baskit.NativeData.(stdlib.BaskitMap)
	result := make(map[string]GoValue)
	for key, val := range underlying {
		goVal, err := ToGoValue(val)
		if err != nil {
			return nil, err
		}
		result[key] = goVal
	}
	return result, nil
}

// sliceToArray converts a Go slice/array to an Objective-LOL BUKKIT object
func sliceToArray(rv reflect.Value) (environment.Value, error) {
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
		objElem, err := FromGoValue(WrapAny(elem.Interface()))
		if err != nil {
			return nil, err
		}
		instance.NativeData = append(instance.NativeData.(stdlib.BukkitSlice), objElem)
	}

	return instance, nil
}

// mapToObject converts a Go map to an Objective-LOL BASKIT object
func mapToObject(rv reflect.Value) (environment.Value, error) {
	if rv.Kind() != reflect.Map {
		return nil, NewConversionError(
			fmt.Sprintf("cannot convert Go type %s to BASKIT", rv.Type()),
			nil,
		)
	}

	// Check for existing constructed object to avoid duplication
	if val := rv.MapIndex(reflect.ValueOf(GoValueIDKey)); val.IsValid() {
		value := val.Interface().(string)
		instance, ok := constructedObjects[value]
		if ok {
			return instance, nil
		}
		return nil, NewConversionError(
			fmt.Sprintf("found %s as %s but no corresponding constructed object", value, GoValueIDKey),
			nil,
		)
	}

	// Create a new BASKIT object
	instance := stdlib.NewBaskitInstance()

	// Populate the BASKIT object with the map elements
	for _, key := range rv.MapKeys() {
		val := rv.MapIndex(key)
		objKey, err := FromGoValue(WrapAny(key.Interface()))
		if err != nil {
			return nil, err
		}
		objVal, err := FromGoValue(WrapAny(val.Interface()))
		if err != nil {
			return nil, err
		}
		underlying := instance.NativeData.(stdlib.BaskitMap)
		underlying[objKey.String()] = objVal
	}

	return instance, nil
}

// structToObject converts a Go struct to an Objective-LOL BASKIT object
func structToObject(rv reflect.Value) (environment.Value, error) {
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

		// Convert field name and value to Objective-LOL environment
		objKey, err := FromGoValue(WrapAny(field.Name))
		if err != nil {
			return nil, err
		}
		objVal, err := FromGoValue(WrapAny(value.Interface()))
		if err != nil {
			return nil, err
		}

		// Add the key-value pair to the BASKIT object
		underlying := instance.NativeData.(stdlib.BaskitMap)
		underlying[objKey.String()] = objVal
	}
	instance.Variables["SIZ"].Value = environment.IntegerValue(rv.NumField())

	return instance, nil
}

// ConvertArguments converts a slice of Go values to Objective-LOL values
func ConvertArguments(args []GoValue) ([]environment.Value, error) {
	result := make([]environment.Value, len(args))
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
