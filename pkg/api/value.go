package api

import (
	"encoding/json"
	"fmt"

	"github.com/bjia56/objective-lol/pkg/environment"
)

const (
	GoValueIDKey = "__GoValue_id"
)

type GoValue struct {
	value interface{}
}

func (v GoValue) Get() interface{} {
	return v.value
}

func (v GoValue) ID() string {
	if obj, ok := v.value.(*environment.ObjectInstance); ok {
		return fmt.Sprintf("%p", obj)
	}
	return ""
}

func (v GoValue) MarshalJSON() ([]byte, error) {
	if id := v.ID(); id != "" {
		return json.Marshal(map[string]string{
			GoValueIDKey: id,
		})
	}
	return json.Marshal(v.value)
}

func (v GoValue) Type() string {
	switch v.value.(type) {
	case nil:
		return "NOTHIN"
	case string:
		return "STRIN"
	case int64:
		return "INTEGR"
	case float64:
		return "DUBBLE"
	case bool:
		return "BOOL"
	case []GoValue:
		return "BUKKIT"
	case map[string]GoValue:
		return "BASKIT"
	case *environment.ObjectInstance:
		return v.value.(*environment.ObjectInstance).Class.Name
	default:
		return "UNKNOWN"
	}
}

func (v GoValue) Int() (int64, error) {
	if val, ok := v.value.(int64); ok {
		return val, nil
	}
	return 0, fmt.Errorf("value is not an int64")
}

func (v GoValue) Float() (float64, error) {
	if val, ok := v.value.(float64); ok {
		return val, nil
	}
	return 0, fmt.Errorf("value is not a float64")
}

func (v GoValue) String() (string, error) {
	if val, ok := v.value.(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("value is not a string")
}

func (v GoValue) Bool() (bool, error) {
	if val, ok := v.value.(bool); ok {
		return val, nil
	}
	return false, fmt.Errorf("value is not a bool")
}

func (v GoValue) Slice() ([]GoValue, error) {
	if val, ok := v.value.([]GoValue); ok {
		return val, nil
	}
	return nil, fmt.Errorf("value is not a slice")
}

func (v GoValue) Map() (map[string]GoValue, error) {
	if val, ok := v.value.(map[string]GoValue); ok {
		return val, nil
	}
	return nil, fmt.Errorf("value is not a map")
}

func (v GoValue) Object() (*environment.ObjectInstance, error) {
	if val, ok := v.value.(*environment.ObjectInstance); ok {
		return val, nil
	}
	return nil, fmt.Errorf("value is not an object instance")
}

func WrapAny(value interface{}) GoValue {
	return GoValue{value: value}
}

func WrapString(value string) GoValue {
	return GoValue{value: value}
}

func WrapInt(value int64) GoValue {
	return GoValue{value: value}
}

func WrapFloat(value float64) GoValue {
	return GoValue{value: value}
}

func WrapBool(value bool) GoValue {
	return GoValue{value: value}
}

func WrapObject(value *environment.ObjectInstance) GoValue {
	return GoValue{value: value}
}
