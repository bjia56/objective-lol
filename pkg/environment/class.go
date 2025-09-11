package environment

import "fmt"

// Class represents a native Objective-LOL class definition
type Class struct {
	Documentation    []string
	Name             string   // Display name: "READER"
	QualifiedName    string   // Internal: "stdlib:IO.READER"
	ModulePath       string   // Internal: "stdlib:IO"
	ParentClasses    []string // Internal: qualified parent names (multiple inheritance)
	MRO              []string // Method Resolution Order (C3 linearization)
	PublicVariables  map[string]*MemberVariable
	PrivateVariables map[string]*MemberVariable
	PublicFunctions  map[string]*Function
	PrivateFunctions map[string]*Function
	SharedVariables  map[string]*MemberVariable
	SharedFunctions  map[string]*Function

	UnknownFunctionHandler func(fnName string, fromContext string) (*Function, error) // Handler for unknown functions (if any)
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
		PublicVariables:  make(map[string]*MemberVariable),
		PrivateVariables: make(map[string]*MemberVariable),
		PublicFunctions:  make(map[string]*Function),
		PrivateFunctions: make(map[string]*Function),
		SharedVariables:  make(map[string]*MemberVariable),
		SharedFunctions:  make(map[string]*Function),
	}
}

// GetMemberFunction is a helper that recursively searches for a function in the class hierarchy
func (c *Class) GetMemberFunction(name string, fromContext string, env *Environment) (*Function, error) {
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
				if function, err := parentClass.GetMemberFunction(name, fromContext, env); err == nil {
					return function, nil
				}
			}
		}
	}

	return nil, &NotFound{Message: fmt.Sprintf("undefined member function '%s'", name)}
}

func (c *Class) CheckUnknownFunctionHandler(fnName, fromContext string) (*Function, error) {
	if c.UnknownFunctionHandler != nil {
		return c.UnknownFunctionHandler(fnName, fromContext)
	}
	return nil, &NotFound{Message: fmt.Sprintf("undefined member function '%s'", fnName)}
}

// NewObject creates a new instance of the specified class
func (c *Class) NewObject(env *Environment) (*ObjectInstance, error) {
	// Compute or retrieve cached MRO for this class
	_, err := c.computeOrGetMRO(env)
	if err != nil {
		return nil, fmt.Errorf("failed to compute method resolution order for class %s: %v", c.Name, err)
	}

	instance := &ObjectInstance{
		Environment:     env,
		Class:           c,
		Variables:       make(map[string]*MemberVariable),
		SharedVariables: c.SharedVariables,
	}

	// Initialize instance variables using MRO
	env.InitializeInstanceVariablesWithMRO(instance)

	return instance, nil
}

// computeOrGetMRO computes or retrieves cached Method Resolution Order for a class
func (c *Class) computeOrGetMRO(env *Environment) ([]string, error) {
	// If MRO already computed, return it
	if len(c.MRO) > 0 {
		return c.MRO, nil
	}

	// Compute MRO using C3 linearization
	mro, err := c.computeC3Linearization(env)
	if err != nil {
		return nil, err
	}

	// Cache the result
	c.MRO = mro
	return mro, nil
}

// computeC3Linearization implements the C3 linearization algorithm
func (c *Class) computeC3Linearization(env *Environment) ([]string, error) {
	// Base case: no parents
	if len(c.ParentClasses) == 0 {
		return []string{c.QualifiedName}, nil
	}

	// Compute linearizations for all parent classes
	parentLinearizations := make([][]string, 0, len(c.ParentClasses))
	for _, parent := range c.ParentClasses {
		parentClass, err := env.GetClass(parent)
		if err != nil {
			return nil, err
		}

		// Objective-LOL parent class
		parentMRO, err := parentClass.computeOrGetMRO(env)
		if err != nil {
			return nil, err
		}
		parentLinearizations = append(parentLinearizations, parentMRO)
	}

	// Merge all linearizations using C3 algorithm
	merged, err := mergeLinearizations(parentLinearizations)
	if err != nil {
		return nil, fmt.Errorf("multiple inheritance conflict in class %s: %v", c.QualifiedName, err)
	}

	// Result is current class + merged linearizations
	result := []string{c.QualifiedName}
	result = append(result, merged...)
	return result, nil
}

// mergeLinearizations implements the C3 merge algorithm
func mergeLinearizations(linearizations [][]string) ([]string, error) {
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
