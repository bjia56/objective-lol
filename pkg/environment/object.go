package environment

import (
	"fmt"
	"slices"
)

// GenericObject is an interface for all objects.
type GenericObject interface {
	Value
	GetQualifiedClassName() string
	GetMemberVariable(name string, callingContext string) (Value, error)
	SetMemberVariable(name string, value Value, callingContext string) error
	GetMemberFunction(name string, callingContext string) (*Function, error)
}

// ObjectInstance represents a native Objective-LOL instance of a class
type ObjectInstance struct {
	Environment     *Environment // Environment in which the instance was created
	Class           *Class
	MRO             []string // Method Resolution Order (stored for efficiency)
	Variables       map[string]*Variable
	SharedVariables map[string]*Variable // Reference to class shared variables
	NativeData      any                  // For native classes, stores internal data
}

// GetQualifiedClassName returns the class name of the object instance
func (obj *ObjectInstance) GetQualifiedClassName() string {
	return obj.Class.QualifiedName
}

// getMemberVariable retrieves a member variable from the object instance
func (obj *ObjectInstance) getMemberVariable(name string, callingContext string) (*Variable, error) {
	// Check instance variables first
	if variable, exists := obj.Variables[name]; exists {
		// Check visibility using the variable's IsPublic flag
		if variable.IsPublic || callingContext == obj.Class.QualifiedName {
			return variable, nil
		}
		return nil, fmt.Errorf("variable '%s' is private", name)
	}

	// Check shared variables
	if variable, exists := obj.SharedVariables[name]; exists {
		return variable, nil
	}

	return nil, fmt.Errorf("undefined member variable '%s'", name)
}

// GetMemberVariable retrieves the value of a member variable from the object instance
func (obj *ObjectInstance) GetMemberVariable(name string, callingContext string) (Value, error) {
	variable, err := obj.getMemberVariable(name, callingContext)
	if err != nil {
		return nil, err
	}
	return variable.Value, nil
}

// SetMemberVariable sets a member variable in the object instance
func (obj *ObjectInstance) SetMemberVariable(name string, value Value, callingContext string) error {
	variable, err := obj.getMemberVariable(name, callingContext)
	if err != nil {
		return err
	}

	if variable.IsLocked {
		return fmt.Errorf("cannot assign to locked variable '%s'", name)
	}

	// Cast value to the variable's type
	castedValue, err := value.Cast(variable.Type)
	if err != nil {
		return fmt.Errorf("cannot assign to variable '%s': %v", name, err)
	}

	variable.Value = castedValue
	return nil
}

func (o *ObjectInstance) Type() string      { return o.Class.Name }
func (o *ObjectInstance) String() string    { return fmt.Sprintf("<%s object>", o.Class.Name) }
func (o *ObjectInstance) Copy() Value       { return o }   // Objects are reference types
func (o *ObjectInstance) ToBool() BoolValue { return YEZ } // Objects are always truthy
func (o *ObjectInstance) IsNothing() bool   { return o == nil }

func (o *ObjectInstance) Cast(targetType string) (Value, error) {
	if targetType == "" {
		return o, nil
	}

	switch targetType {
	case "STRIN":
		return StringValue(o.String()), nil
	case "BOOL":
		return o.ToBool(), nil
	case "INTEGR", "DUBBLE", "NOTHIN":
		return nil, fmt.Errorf("cannot cast %s to %s", o.Class.Name, targetType)
	}

	if o.Environment == nil {
		return nil, fmt.Errorf("cannot cast %s to type %s", o.Class.Name, targetType)
	}

	targetCls, err := o.Environment.GetClass(targetType)
	if err != nil {
		return nil, fmt.Errorf("cannot cast %s to unknown type %s", o.Class.Name, targetType)
	}

	if slices.Contains(o.MRO, targetCls.QualifiedName) {
		return o, nil
	}

	return nil, fmt.Errorf("cannot cast %s (%s) to type %s (%s)", o.Class.Name, o.Class.QualifiedName, targetType, targetCls.QualifiedName)
}

func (o *ObjectInstance) EqualTo(other Value) (BoolValue, error) {
	if other.IsNothing() {
		return BoolValue(o.IsNothing()), nil
	}

	if otherObj, ok := other.(*ObjectInstance); ok {
		// Simple reference equality for now
		return BoolValue(o == otherObj), nil
	}

	return NO, nil
}

// GetMemberFunction retrieves a member function from the object's class
func (obj *ObjectInstance) GetMemberFunction(name string, callingContext string) (*Function, error) {
	return obj.Class.getMemberFunction(name, callingContext, obj.Environment)
}
