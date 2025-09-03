package types

import (
	"fmt"
	"slices"
	"strings"
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
	
	// Get qualified hierarchy for type safety
	hierarchy := o.Instance.GetHierarchy()
	
	// Try exact qualified match first (for type safety)
	if slices.Contains(hierarchy, targetType) {
		return o, nil
	}
	
	// For backward compatibility with unqualified names in user code,
	// check if targetType matches any class name in hierarchy
	for _, qualifiedName := range hierarchy {
		if strings.HasSuffix(qualifiedName, "."+targetType) {
			return o, nil
		}
	}
	
	// Standard type casts
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
