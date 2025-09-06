package api

import "encoding/json"

type GoValue struct {
	value interface{}
}

func (v GoValue) Get() interface{} {
	return v.value
}

func (v GoValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
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
