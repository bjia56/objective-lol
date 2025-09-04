package environment

import (
	"fmt"

	"github.com/bjia56/objective-lol/pkg/types"
)

// Variable represents a variable with its type information and mutability
type Variable struct {
	Name     string
	Type     string
	Value    types.Value
	IsLocked bool
	IsPublic bool // Track if this variable is public
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
	Name        string
	ReturnType  string
	Parameters  []Parameter
	Body        interface{} // Will hold AST nodes
	IsShared    *bool       // nil for global functions, true/false for class methods
	NativeImpl  func(ctx interface{}, this *ObjectInstance, args []types.Value) (types.Value, error)
}

// Parameter represents a function parameter
type Parameter struct {
	Name string
	Type string
}

// Class represents an Objective-LOL class definition
type Class struct {
	Name             string                  // Display name: "READER"
	QualifiedName    string                  // Internal: "stdlib:IO.READER"
	ModulePath       string                  // Internal: "stdlib:IO"
	ParentClasses    []string                // Internal: qualified parent names (multiple inheritance)
	MRO              []string                // Method Resolution Order (C3 linearization)
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
func (e *Environment) DefineVariable(name, varType string, value types.Value, isLocked bool) error {
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
		Name:     name,
		Type:     varType,
		Value:    castedValue,
		IsLocked: isLocked,
		IsPublic: true, // Regular variables are public by default
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
func (e *Environment) SetVariable(name string, value types.Value) error {
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
		// Only warn about name collisions, don't fail - qualified names prevent actual conflicts
		// In a real implementation, we might want to handle import scoping here
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
	var qualifiedName string
	if modulePath != "" {
		qualifiedName = fmt.Sprintf("%s.%s", modulePath, name)
	} else {
		qualifiedName = name // Fallback for legacy/local classes
	}
	
	return &Class{
		Name:             name,
		QualifiedName:    qualifiedName,
		ModulePath:       modulePath,
		ParentClasses:    parentClasses, // Support multiple parents
		MRO:              []string{},     // Method Resolution Order (computed later)
		PublicVariables:  make(map[string]*Variable),
		PrivateVariables: make(map[string]*Variable),
		PublicFunctions:  make(map[string]*Function),
		PrivateFunctions: make(map[string]*Function),
		SharedVariables:  make(map[string]*Variable),
		SharedFunctions:  make(map[string]*Function),
	}
}

// NewClassLegacy creates a new class definition with single parent (backwards compatibility)
func NewClassLegacy(name, modulePath, parentClass string) *Class {
	var parentClasses []string
	if parentClass != "" {
		parentClasses = []string{parentClass}
	}
	return NewClass(name, modulePath, parentClasses)
}

// ObjectInstance represents an instance of a class
type ObjectInstance struct {
	Class           *Class
	MRO             []string             // Method Resolution Order (stored for efficiency)
	Variables       map[string]*Variable
	SharedVariables map[string]*Variable // Reference to class shared variables
	NativeData      any                  // For native classes, stores internal data
}

// NewObjectInstance creates a new instance of the specified class
func (e *Environment) NewObjectInstance(className string) (types.ObjectInstance, error) {
	class, err := e.GetClass(className)
	if err != nil {
		return nil, err
	}

	// Compute or retrieve cached MRO for this class
	mro, err := e.computeOrGetMRO(className)
	if err != nil {
		return nil, fmt.Errorf("failed to compute method resolution order for class %s: %v", className, err)
	}

	instance := &ObjectInstance{
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
func (e *Environment) computeOrGetMRO(className string) ([]string, error) {
	class, err := e.GetClass(className)
	if err != nil {
		return nil, err
	}

	// If MRO already computed, return it
	if len(class.MRO) > 0 {
		return class.MRO, nil
	}

	// Compute MRO using C3 linearization
	mro, err := e.computeC3Linearization(className)
	if err != nil {
		return nil, err
	}

	// Cache the result
	class.MRO = mro
	return mro, nil
}

// computeC3Linearization implements the C3 linearization algorithm
func (e *Environment) computeC3Linearization(className string) ([]string, error) {
	class, err := e.GetClass(className)
	if err != nil {
		return nil, err
	}

	// Base case: no parents
	if len(class.ParentClasses) == 0 {
		return []string{className}, nil
	}

	// Compute linearizations for all parent classes
	parentLinearizations := make([][]string, 0, len(class.ParentClasses))
	for _, parent := range class.ParentClasses {
		parentMRO, err := e.computeC3Linearization(parent)
		if err != nil {
			return nil, err
		}
		parentLinearizations = append(parentLinearizations, parentMRO)
	}

	// Add the list of direct parents
	parentLinearizations = append(parentLinearizations, class.ParentClasses)

	// Merge all linearizations using C3 algorithm
	merged, err := e.mergeLinearizations(parentLinearizations)
	if err != nil {
		return nil, fmt.Errorf("multiple inheritance conflict in class %s: %v", className, err)
	}

	// Result is current class + merged linearizations
	result := []string{className}
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
			if _, exists := instance.Variables[name]; !exists { // Don't override existing
				instance.Variables[name] = &Variable{
					Name:     variable.Name,
					Type:     variable.Type,
					Value:    variable.Value.Copy(),
					IsLocked: variable.IsLocked,
					IsPublic: true,
				}
			}
		}

		// Initialize private variables (only from same class)
		if className == instance.Class.QualifiedName {
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
}

// initializeInstanceVariables recursively initializes variables from class hierarchy (legacy method)
func (e *Environment) initializeInstanceVariables(instance *ObjectInstance, class *Class) {
	// First initialize parent class variables
	for _, parentClassName := range class.ParentClasses {
		if parentClass, err := e.GetClass(parentClassName); err == nil {
			e.initializeInstanceVariables(instance, parentClass)
		}
	}

	// Then initialize current class variables (may override parent variables)
	for name, variable := range class.PublicVariables {
		instance.Variables[name] = &Variable{
			Name:     variable.Name,
			Type:     variable.Type,
			Value:    variable.Value.Copy(),
			IsLocked: variable.IsLocked,
			IsPublic: true, // Public variables
		}
	}

	for name, variable := range class.PrivateVariables {
		instance.Variables[name] = &Variable{
			Name:     variable.Name,
			Type:     variable.Type,
			Value:    variable.Value.Copy(),
			IsLocked: variable.IsLocked,
			IsPublic: false, // Private variables
		}
	}
}

// GetMemberVariable retrieves a member variable from the object instance
func (obj *ObjectInstance) GetMemberVariable(name string, fromContext string) (*Variable, error) {
	// Check instance variables first
	if variable, exists := obj.Variables[name]; exists {
		// Check visibility using the variable's IsPublic flag
		if variable.IsPublic || fromContext == obj.Class.Name {
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
func (obj *ObjectInstance) SetMemberVariable(name string, value types.Value, fromContext string) error {
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

// GetMemberFunction retrieves a member function from the object's class
func (obj *ObjectInstance) GetMemberFunction(name string, fromContext string, env *Environment) (*Function, error) {
	return obj.Class.getMemberFunction(name, fromContext, env)
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
	} else {
		// Fallback for classes without computed MRO
		for _, parentClassName := range c.ParentClasses {
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
