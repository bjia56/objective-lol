package environment

import (
	"fmt"
)

var constructedObjects = make(map[string]*ObjectInstance)

func LookupObject(id string) (*ObjectInstance, error) {
	if obj, ok := constructedObjects[id]; ok {
		return obj, nil
	}
	return nil, &NotFound{Message: fmt.Sprintf("no object found with id %s", id)}
}

type Interpreter interface {
	CallFunction(function string, args []Value) (Value, error)
	CallMemberFunction(object *ObjectInstance, function string, args []Value) (Value, error)
	Fork() Interpreter
}

// Variable represents a variable with its type information and mutability
type Variable struct {
	Documentation []string
	Name          string
	Type          string
	Value         Value
	IsLocked      bool
	IsPublic      bool // Track if this variable is public
}

// MemberVariable represents a member variable in an object instance
type MemberVariable struct {
	Variable

	NativeGet func(this *ObjectInstance) (Value, error)
	NativeSet func(this *ObjectInstance, value Value) error
}

func (v *MemberVariable) Get(this *ObjectInstance) (Value, error) {
	if v.NativeGet != nil {
		return v.NativeGet(this)
	}
	return v.Variable.Value, nil
}

func (v *MemberVariable) Set(this *ObjectInstance, value Value) {
	if v.NativeSet != nil {
		v.NativeSet(this, value)
	} else {
		v.Variable.Value = value
	}
}

// Environment represents a lexical scope for variables and functions
type Environment struct {
	parent    *Environment
	variables map[string]*Variable
	functions map[string]*Function
	classes   map[string]*Class
}

// Function represents a user-defined or native function
type Function struct {
	Documentation []string
	Name          string
	ReturnType    string
	Parameters    []Parameter
	Body          interface{} // Will hold AST nodes
	IsShared      *bool       // nil for global functions, true/false for class methods
	IsVarargs     bool        // true if function accepts variable number of arguments
	NativeImpl    func(interpreter Interpreter, this *ObjectInstance, args []Value) (Value, error)
}

// Parameter represents a function parameter
type Parameter struct {
	Name string
	Type string
}

// NewEnvironment creates a new environment with an optional parent
func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent:    parent,
		variables: make(map[string]*Variable),
		functions: make(map[string]*Function),
		classes:   make(map[string]*Class),
	}
}

// DefineVariable defines a new variable in the current scope
func (e *Environment) DefineVariable(name, varType string, value Value, isLocked bool, docs []string) error {
	// Check if variable already exists in current scope
	if _, exists := e.variables[name]; exists {
		return &AlreadyExists{Message: fmt.Sprintf("variable '%s' already defined in current scope", name)}
	}

	var castedValue Value
	var actualType string

	// Handle dynamic typing for empty type specification
	if varType == "" {
		// Dynamic variable - infer type from value if provided, otherwise use NOTHIN
		if value == nil || value.IsNothing() {
			castedValue = NOTHIN
		} else {
			castedValue = value
		}
	} else {
		// Static typed variable - cast value to the specified type
		var err error
		castedValue, err = value.Cast(varType)
		if err != nil {
			return fmt.Errorf("cannot initialize variable '%s': %v", name, err)
		}
		actualType = varType
	}

	e.variables[name] = &Variable{
		Documentation: docs,
		Name:          name,
		Type:          actualType,
		Value:         castedValue,
		IsLocked:      isLocked,
		IsPublic:      true, // Regular variables are public by default
	}

	return nil
}

// GetVariable retrieves a variable from the current scope or parent scopes
func (e *Environment) GetVariable(name string) (*Variable, error) {
	if variable, exists := e.variables[name]; exists {
		return variable, nil
	}

	if e.parent != nil {
		return e.parent.GetVariable(name)
	}

	return nil, &NotFound{Message: fmt.Sprintf("undefined variable '%s'", name)}
}

// SetVariable sets the value of an existing variable
func (e *Environment) SetVariable(name string, value Value) error {
	variable, err := e.GetVariable(name)
	if err != nil {
		return err
	}

	if variable.IsLocked {
		return fmt.Errorf("cannot assign to locked variable '%s'", name)
	}

	castedValue, err := value.Cast(variable.Type)
	if err != nil {
		return fmt.Errorf("cannot assign to variable '%s': %v", name, err)
	}
	variable.Value = castedValue

	return nil
}

// DefineFunction defines a new function in the current scope
func (e *Environment) DefineFunction(function *Function) error {
	// Check if function already exists in current scope
	if _, exists := e.functions[function.Name]; exists {
		return &AlreadyExists{Message: fmt.Sprintf("function '%s' already defined in current scope", function.Name)}
	}

	e.functions[function.Name] = function
	return nil
}

// GetFunction retrieves a function from the current scope or parent scopes
func (e *Environment) GetFunction(name string) (*Function, error) {
	if function, exists := e.functions[name]; exists {
		return function, nil
	}

	if e.parent != nil {
		return e.parent.GetFunction(name)
	}

	return nil, &NotFound{Message: fmt.Sprintf("undefined function '%s'", name)}
}

// DefineClass defines a new class in the current scope
func (e *Environment) DefineClass(class *Class) error {
	// Store by qualified name (primary key for type safety)
	if _, exists := e.classes[class.QualifiedName]; exists {
		return &AlreadyExists{Message: fmt.Sprintf("class '%s' already defined in current scope", class.QualifiedName)}
	}
	e.classes[class.QualifiedName] = class

	// Also store by simple name for user code compatibility
	// This allows lookup by simple names like "READER" while maintaining qualified internal storage
	if existing, exists := e.classes[class.Name]; exists && existing.QualifiedName != class.QualifiedName {
		return &AlreadyExists{Message: fmt.Sprintf("class %s (%s) redeclared as %s (%s) in current scope", existing.Name, existing.QualifiedName, class.Name, class.QualifiedName)}
	}
	e.classes[class.Name] = class

	return nil
}

// GetClass retrieves a class from the current scope or parent scopes
// Supports both qualified names (e.g., "stdlib:IO.READER") and simple names (e.g., "READER")
func (e *Environment) GetClass(name string) (*Class, error) {
	if class, exists := e.classes[name]; exists {
		return class, nil
	}

	if e.parent != nil {
		return e.parent.GetClass(name)
	}

	return nil, &NotFound{Message: fmt.Sprintf("undefined class '%s'", name)}
}

func (e *Environment) NewObjectInstance(className string) (*ObjectInstance, error) {
	class, err := e.GetClass(className)
	if err != nil {
		return nil, err
	}

	instance, err := class.NewObject(e)
	if err != nil {
		return nil, err
	}

	constructedObjects[instance.ID()] = instance

	return instance, nil
}

// RuntimeEnvironment manages the global runtime state
type RuntimeEnvironment struct {
	GlobalEnv *Environment
	ExecDir   string
}

// NewRuntimeEnvironment creates a new runtime environment
func NewRuntimeEnvironment() *RuntimeEnvironment {
	return &RuntimeEnvironment{
		GlobalEnv: NewEnvironment(nil),
	}
}

// NewLocalEnv creates a new local environment with the global environment as parent
func (rt *RuntimeEnvironment) NewLocalEnv() *Environment {
	return NewEnvironment(rt.GlobalEnv)
}

// GetAllFunctions returns a copy of all functions in the environment (current scope only)
func (e *Environment) GetAllFunctions() map[string]*Function {
	result := make(map[string]*Function)
	for name, function := range e.functions {
		result[name] = function
	}
	return result
}

// GetAllClasses returns a copy of all classes in the environment (current scope only)
func (e *Environment) GetAllClasses() map[string]*Class {
	result := make(map[string]*Class)
	for name, class := range e.classes {
		result[name] = class
	}
	return result
}

// GetAllVariables returns a copy of all variables in the environment (current scope only)
func (e *Environment) GetAllVariables() map[string]*Variable {
	result := make(map[string]*Variable)
	for name, variable := range e.variables {
		result[name] = variable
	}
	return result
}

// InitializeInstanceVariablesWithMRO initializes variables using Method Resolution Order
func (e *Environment) InitializeInstanceVariablesWithMRO(instance *ObjectInstance) {
	// Initialize variables following MRO order (reverse for proper inheritance)
	for i := len(instance.Class.MRO) - 1; i >= 0; i-- {
		className := instance.Class.MRO[i]
		class, err := e.GetClass(className)
		if err != nil {
			continue // Skip invalid classes
		}

		// Initialize public variables
		for name, variable := range class.PublicVariables {
			var val Value
			if variable.Value != nil {
				val = variable.Value.Copy()
			}
			instance.Variables[name] = &MemberVariable{
				Variable: Variable{
					Name:     variable.Name,
					Type:     variable.Type,
					Value:    val,
					IsLocked: variable.IsLocked,
					IsPublic: true,
				},
				NativeGet: variable.NativeGet,
				NativeSet: variable.NativeSet,
			}
		}

		// Initialize private variables
		for name, variable := range class.PrivateVariables {
			var val Value
			if variable.Value != nil {
				val = variable.Value.Copy()
			}
			instance.Variables[name] = &MemberVariable{
				Variable: Variable{
					Name:     variable.Name,
					Type:     variable.Type,
					Value:    val,
					IsLocked: variable.IsLocked,
					IsPublic: false,
				},
				NativeGet: variable.NativeGet,
				NativeSet: variable.NativeSet,
			}
		}
	}
}
