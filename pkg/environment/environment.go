package environment

import (
	"fmt"
	"slices"
)

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
	NativeImpl    func(interpreter Interpreter, this *ObjectInstance, args []Value) (Value, error)
}

// Parameter represents a function parameter
type Parameter struct {
	Name string
	Type string
}

// Class represents an Objective-LOL class definition
type Class struct {
	Documentation    []string
	Name             string   // Display name: "READER"
	QualifiedName    string   // Internal: "stdlib:IO.READER"
	ModulePath       string   // Internal: "stdlib:IO"
	ParentClasses    []string // Internal: qualified parent names (multiple inheritance)
	MRO              []string // Method Resolution Order (C3 linearization)
	PublicVariables  map[string]*Variable
	PrivateVariables map[string]*Variable
	PublicFunctions  map[string]*Function
	PrivateFunctions map[string]*Function
	SharedVariables  map[string]*Variable
	SharedFunctions  map[string]*Function
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
		return fmt.Errorf("variable '%s' already defined in current scope", name)
	}

	// Cast value to the specified type
	castedValue, err := value.Cast(varType)
	if err != nil {
		return fmt.Errorf("cannot initialize variable '%s': %v", name, err)
	}

	e.variables[name] = &Variable{
		Documentation: docs,
		Name:          name,
		Type:          varType,
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

	return nil, fmt.Errorf("undefined variable '%s'", name)
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

	// Cast value to the variable's type
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
		return fmt.Errorf("function '%s' already defined in current scope", function.Name)
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

	return nil, fmt.Errorf("undefined function '%s'", name)
}

// DefineClass defines a new class in the current scope
func (e *Environment) DefineClass(class *Class) error {
	// Store by qualified name (primary key for type safety)
	if _, exists := e.classes[class.QualifiedName]; exists {
		return fmt.Errorf("class '%s' already defined in current scope", class.QualifiedName)
	}
	e.classes[class.QualifiedName] = class

	// Also store by simple name for user code compatibility
	// This allows lookup by simple names like "READER" while maintaining qualified internal storage
	if existing, exists := e.classes[class.Name]; exists && existing.QualifiedName != class.QualifiedName {
		return fmt.Errorf("class %s (%s) redeclared as %s (%s) in current scope", existing.Name, existing.QualifiedName, class.Name, class.QualifiedName)
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

	return nil, fmt.Errorf("undefined class '%s'", name)
}

// NewClass creates a new class definition with support for multiple inheritance
func NewClass(name, modulePath string, parentClasses []string) *Class {
	qualifiedName := fmt.Sprintf("%s.%s", modulePath, name)

	return &Class{
		Name:             name,
		QualifiedName:    qualifiedName,
		ModulePath:       modulePath,
		ParentClasses:    parentClasses, // Support multiple parents
		MRO:              []string{},    // Method Resolution Order (computed later)
		PublicVariables:  make(map[string]*Variable),
		PrivateVariables: make(map[string]*Variable),
		PublicFunctions:  make(map[string]*Function),
		PrivateFunctions: make(map[string]*Function),
		SharedVariables:  make(map[string]*Variable),
		SharedFunctions:  make(map[string]*Function),
	}
}

// ObjectInstance represents an instance of a class
type ObjectInstance struct {
	Environment     *Environment // Environment in which the instance was created
	Class           *Class
	MRO             []string // Method Resolution Order (stored for efficiency)
	Variables       map[string]*Variable
	SharedVariables map[string]*Variable // Reference to class shared variables
	NativeData      any                  // For native classes, stores internal data
}

// NewObjectInstance creates a new instance of the specified class
func (e *Environment) NewObjectInstance(className string) (*ObjectInstance, error) {
	class, err := e.GetClass(className)
	if err != nil {
		return nil, err
	}

	// Compute or retrieve cached MRO for this class
	mro, err := e.computeOrGetMRO(class)
	if err != nil {
		return nil, fmt.Errorf("failed to compute method resolution order for class %s: %v", className, err)
	}

	instance := &ObjectInstance{
		Environment:     e,
		Class:           class,
		MRO:             mro,
		Variables:       make(map[string]*Variable),
		SharedVariables: class.SharedVariables,
	}

	// Initialize instance variables using MRO
	e.initializeInstanceVariablesWithMRO(instance)

	return instance, nil
}

// computeOrGetMRO computes or retrieves cached Method Resolution Order for a class
func (e *Environment) computeOrGetMRO(class *Class) ([]string, error) {
	// If MRO already computed, return it
	if len(class.MRO) > 0 {
		return class.MRO, nil
	}

	// Compute MRO using C3 linearization
	mro, err := e.computeC3Linearization(class)
	if err != nil {
		return nil, err
	}

	// Cache the result
	class.MRO = mro
	return mro, nil
}

// computeC3Linearization implements the C3 linearization algorithm
func (e *Environment) computeC3Linearization(class *Class) ([]string, error) {
	// Base case: no parents
	if len(class.ParentClasses) == 0 {
		return []string{class.QualifiedName}, nil
	}

	// Compute linearizations for all parent classes
	parentLinearizations := make([][]string, 0, len(class.ParentClasses))
	for _, parent := range class.ParentClasses {
		parentClass, err := e.GetClass(parent)
		if err != nil {
			return nil, err
		}

		parentMRO, err := e.computeOrGetMRO(parentClass)
		if err != nil {
			return nil, err
		}
		parentLinearizations = append(parentLinearizations, parentMRO)
	}

	// Merge all linearizations using C3 algorithm
	merged, err := e.mergeLinearizations(parentLinearizations)
	if err != nil {
		return nil, fmt.Errorf("multiple inheritance conflict in class %s: %v", class.QualifiedName, err)
	}

	// Result is current class + merged linearizations
	result := []string{class.QualifiedName}
	result = append(result, merged...)
	return result, nil
}

// mergeLinearizations implements the C3 merge algorithm
func (e *Environment) mergeLinearizations(linearizations [][]string) ([]string, error) {
	result := []string{}

	// Create working copies of all linearizations
	working := make([][]string, len(linearizations))
	for i, lin := range linearizations {
		working[i] = make([]string, len(lin))
		copy(working[i], lin)
	}

	for {
		// Remove empty linearizations
		nonEmpty := [][]string{}
		for _, lin := range working {
			if len(lin) > 0 {
				nonEmpty = append(nonEmpty, lin)
			}
		}
		working = nonEmpty

		// If all linearizations are empty, we're done
		if len(working) == 0 {
			break
		}

		// Find a good head (appears first in some linearization and not in the tail of any)
		var goodHead string
		found := false

		for _, lin := range working {
			if len(lin) == 0 {
				continue
			}
			candidate := lin[0]
			isGood := true

			// Check if candidate appears in the tail of any linearization
			for _, otherLin := range working {
				for i := 1; i < len(otherLin); i++ {
					if otherLin[i] == candidate {
						isGood = false
						break
					}
				}
				if !isGood {
					break
				}
			}

			if isGood {
				goodHead = candidate
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("cannot create a consistent method resolution order")
		}

		// Add the good head to result
		result = append(result, goodHead)

		// Remove the good head from all linearizations where it appears first
		for i := range working {
			if len(working[i]) > 0 && working[i][0] == goodHead {
				working[i] = working[i][1:]
			}
		}
	}

	return result, nil
}

// initializeInstanceVariablesWithMRO initializes variables using Method Resolution Order
func (e *Environment) initializeInstanceVariablesWithMRO(instance *ObjectInstance) {
	// Initialize variables following MRO order (reverse for proper inheritance)
	for i := len(instance.MRO) - 1; i >= 0; i-- {
		className := instance.MRO[i]
		class, err := e.GetClass(className)
		if err != nil {
			continue // Skip invalid classes
		}

		// Initialize public variables
		for name, variable := range class.PublicVariables {
			instance.Variables[name] = &Variable{
				Name:     variable.Name,
				Type:     variable.Type,
				Value:    variable.Value.Copy(),
				IsLocked: variable.IsLocked,
				IsPublic: true,
			}
		}

		// Initialize private variables
		for name, variable := range class.PrivateVariables {
			instance.Variables[name] = &Variable{
				Name:     variable.Name,
				Type:     variable.Type,
				Value:    variable.Value.Copy(),
				IsLocked: variable.IsLocked,
				IsPublic: false,
			}
		}
	}
}

// GetMemberVariable retrieves a member variable from the object instance
func (obj *ObjectInstance) GetMemberVariable(name string, fromContext string) (*Variable, error) {
	// Check instance variables first
	if variable, exists := obj.Variables[name]; exists {
		// Check visibility using the variable's IsPublic flag
		if variable.IsPublic || fromContext == obj.Class.QualifiedName {
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

// SetMemberVariable sets a member variable in the object instance
func (obj *ObjectInstance) SetMemberVariable(name string, value Value, fromContext string) error {
	variable, err := obj.GetMemberVariable(name, fromContext)
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
func (obj *ObjectInstance) GetMemberFunction(name string, fromContext string) (*Function, error) {
	return obj.Class.getMemberFunction(name, fromContext, obj.Environment)
}

// GetHierarchy returns the class hierarchy as a slice of class names (MRO-based)
func (obj *ObjectInstance) GetHierarchy() []string {
	return obj.MRO
}

// getMemberFunction is a helper that recursively searches for a function in the class hierarchy
func (c *Class) getMemberFunction(name string, fromContext string, env *Environment) (*Function, error) {
	// Check public functions
	if function, exists := c.PublicFunctions[name]; exists {
		return function, nil
	}

	// Check private functions (only accessible from same class context)
	if function, exists := c.PrivateFunctions[name]; exists {
		if fromContext == c.Name {
			return function, nil
		}
		return nil, fmt.Errorf("function '%s' is private", name)
	}

	// Check shared functions
	if function, exists := c.SharedFunctions[name]; exists {
		return function, nil
	}

	// Check parent classes using MRO if available
	if len(c.MRO) > 1 { // Skip self (first element)
		for _, parentClassName := range c.MRO[1:] {
			if parentClass, err := env.GetClass(parentClassName); err == nil {
				if function, err := parentClass.getMemberFunction(name, fromContext, env); err == nil {
					return function, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("undefined member function '%s'", name)
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
