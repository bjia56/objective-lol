package types

import (
	"fmt"
	"slices"
)

type ObjectInstance interface {
	GetHierarchy() []string
}

// ObjectValue wraps an object instance to implement the Value interface
type ObjectValue struct {
	Instance  ObjectInstance // Will hold *environment.ObjectInstance
	ClassName string
}

func (o ObjectValue) Type() string      { return o.ClassName }
func (o ObjectValue) String() string    { return fmt.Sprintf("<%s object>", o.ClassName) }
func (o ObjectValue) Copy() Value       { return o }   // Objects are reference types
func (o ObjectValue) ToBool() BoolValue { return YEZ } // Objects are always truthy
func (o ObjectValue) IsNothing() bool   { return o.Instance == nil }

func (o ObjectValue) Cast(targetType string) (Value, error) {
	if targetType == "" {
		return o, nil
	}
	if slices.Contains(o.Instance.GetHierarchy(), targetType) {
		return o, nil
	}
	if targetType == "STRIN" {
		return StringValue(o.String()), nil
	}
	if targetType == "BOOL" {
		return o.ToBool(), nil
	}
	return nil, fmt.Errorf("cannot cast %s to %s", o.ClassName, targetType)
}

func (o ObjectValue) EqualTo(other Value) (BoolValue, error) {
	if other.IsNothing() {
		return BoolValue(o.IsNothing()), nil
	}

	if otherObj, ok := other.(ObjectValue); ok {
		// Simple reference equality for now
		return BoolValue(o.Instance == otherObj.Instance), nil
	}

	return NO, nil
}

// NewObjectValue creates a new ObjectValue
func NewObjectValue(instance ObjectInstance, className string) ObjectValue {
	return ObjectValue{
		Instance:  instance,
		ClassName: className,
	}
}
