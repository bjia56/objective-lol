package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvironment(t *testing.T) {
	// Test creating environment without parent
	env := NewEnvironment(nil)
	assert.NotNil(t, env)
	assert.Nil(t, env.parent)
	assert.NotNil(t, env.variables)
	assert.NotNil(t, env.functions)
	assert.NotNil(t, env.classes)

	// Test creating environment with parent
	parent := NewEnvironment(nil)
	child := NewEnvironment(parent)
	assert.NotNil(t, child)
	assert.Equal(t, parent, child.parent)
}

func TestEnvironmentVariableOperations(t *testing.T) {
	env := NewEnvironment(nil)

	// Test defining a variable
	err := env.DefineVariable("x", "INTEGR", IntegerValue(42), false)
	assert.NoError(t, err)

	// Test getting the variable
	variable, err := env.GetVariable("x")
	assert.NoError(t, err)
	assert.Equal(t, "x", variable.Name)
	assert.Equal(t, "INTEGR", variable.Type)
	assert.Equal(t, IntegerValue(42), variable.Value)
	assert.False(t, variable.IsLocked)

	// Test setting the variable
	err = env.SetVariable("x", IntegerValue(100))
	assert.NoError(t, err)

	variable, err = env.GetVariable("x")
	assert.NoError(t, err)
	assert.Equal(t, IntegerValue(100), variable.Value)
}

func TestEnvironmentVariableScoping(t *testing.T) {
	parent := NewEnvironment(nil)
	child := NewEnvironment(parent)

	// Define variable in parent
	err := parent.DefineVariable("x", "INTEGR", IntegerValue(42), false)
	assert.NoError(t, err)

	// Child should be able to access parent variable
	variable, err := child.GetVariable("x")
	assert.NoError(t, err)
	assert.Equal(t, IntegerValue(42), variable.Value)

	// Define variable with same name in child (shadowing)
	err = child.DefineVariable("x", "INTEGR", IntegerValue(100), false)
	assert.NoError(t, err)

	// Child should see its own variable
	variable, err = child.GetVariable("x")
	assert.NoError(t, err)
	assert.Equal(t, IntegerValue(100), variable.Value)

	// Parent should still see its own variable
	variable, err = parent.GetVariable("x")
	assert.NoError(t, err)
	assert.Equal(t, IntegerValue(42), variable.Value)
}

func TestEnvironmentVariableErrors(t *testing.T) {
	env := NewEnvironment(nil)

	// Test duplicate variable definition
	err := env.DefineVariable("x", "INTEGR", IntegerValue(42), false)
	assert.NoError(t, err)

	err = env.DefineVariable("x", "STRIN", StringValue("hello"), false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already defined")

	// Test getting undefined variable
	_, err = env.GetVariable("undefined")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "undefined variable")

	// Test setting undefined variable
	err = env.SetVariable("undefined", IntegerValue(42))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "undefined variable")
}

func TestEnvironmentLockedVariables(t *testing.T) {
	env := NewEnvironment(nil)

	// Define locked variable
	err := env.DefineVariable("constant", "INTEGR", IntegerValue(42), true)
	assert.NoError(t, err)

	// Should not be able to modify locked variable
	err = env.SetVariable("constant", IntegerValue(100))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "locked")
}

func TestEnvironmentTypeCasting(t *testing.T) {
	env := NewEnvironment(nil)

	// Define variable with type casting
	err := env.DefineVariable("x", "INTEGR", StringValue("42"), false)
	assert.NoError(t, err)

	// Should have cast the string to integer
	variable, err := env.GetVariable("x")
	assert.NoError(t, err)
	assert.Equal(t, IntegerValue(42), variable.Value)
	assert.Equal(t, "INTEGR", variable.Type)

	// Test invalid cast
	err = env.DefineVariable("y", "INTEGR", StringValue("hello"), false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot initialize variable")
}

func TestEnvironmentFunctionOperations(t *testing.T) {
	env := NewEnvironment(nil)

	// Define a function
	function := &Function{
		Name:       "add",
		ReturnType: "INTEGR",
		Parameters: []Parameter{
			{Name: "x", Type: "INTEGR"},
			{Name: "y", Type: "INTEGR"},
		},
		Body: nil, // Would normally hold AST nodes
	}

	err := env.DefineFunction(function)
	assert.NoError(t, err)

	// Get the function
	retrieved, err := env.GetFunction("add")
	assert.NoError(t, err)
	assert.Equal(t, function, retrieved)

	// Test duplicate function definition
	err = env.DefineFunction(function)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already defined")

	// Test getting undefined function
	_, err = env.GetFunction("undefined")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "undefined function")
}

func TestEnvironmentClassOperations(t *testing.T) {
	env := NewEnvironment(nil)

	// Define a class
	class := &Class{
		Name:             "Person",
		ParentClasses:    []string{},
		MRO:              []string{"Person"},
		PublicVariables:  make(map[string]*Variable),
		PrivateVariables: make(map[string]*Variable),
		PublicFunctions:  make(map[string]*Function),
		PrivateFunctions: make(map[string]*Function),
		SharedVariables:  make(map[string]*Variable),
		SharedFunctions:  make(map[string]*Function),
	}

	err := env.DefineClass(class)
	assert.NoError(t, err)

	// Get the class
	retrieved, err := env.GetClass("Person")
	assert.NoError(t, err)
	assert.Equal(t, class, retrieved)

	// Test duplicate class definition
	err = env.DefineClass(class)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already defined")

	// Test getting undefined class
	_, err = env.GetClass("Undefined")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "undefined class")
}

func TestRuntimeEnvironment(t *testing.T) {
	runtime := NewRuntimeEnvironment()

	assert.NotNil(t, runtime)
	assert.NotNil(t, runtime.GlobalEnv)

	// Just test that the runtime environment is created properly
	// Built-in like BUKKIT may be registered elsewhere in the interpreter
	_, err := runtime.GlobalEnv.GetFunction("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "undefined function")
}

func TestObjectInstance(t *testing.T) {
	env := NewEnvironment(nil)

	// Create a simple class for testing
	class := &Class{
		Name:             "TestClass",
		ParentClasses:    []string{},
		MRO:              []string{"TestClass"},
		PublicVariables:  make(map[string]*Variable),
		PrivateVariables: make(map[string]*Variable),
		PublicFunctions:  make(map[string]*Function),
		PrivateFunctions: make(map[string]*Function),
		SharedVariables:  make(map[string]*Variable),
		SharedFunctions:  make(map[string]*Function),
	}

	// Add a public variable to the class
	class.PublicVariables["name"] = &Variable{
		Name:     "name",
		Type:     "STRIN",
		Value:    StringValue("test"),
		IsLocked: false,
		IsPublic: true,
	}

	// Define the class in environment
	err := env.DefineClass(class)
	assert.NoError(t, err)

	// Create an object instance
	instance, err := env.NewObjectInstance("TestClass")
	assert.NoError(t, err)
	assert.Equal(t, "TestClass", instance.Class.Name)
	assert.NotNil(t, instance.Variables)

	// Check that public variables are copied
	assert.Contains(t, instance.Variables, "name")
	assert.Equal(t, StringValue("test"), instance.Variables["name"].Value)

	// Value interface
	assert.Equal(t, "TestClass", instance.Type())
	assert.Equal(t, "<TestClass object>", instance.String())
	assert.Equal(t, YEZ, instance.ToBool()) // Objects are always truthy
	assert.False(t, instance.IsNothing())
	assert.Equal(t, instance, instance.Copy()) // Objects are reference types
}

func TestObjectInstanceMemberAccess(t *testing.T) {
	env := NewEnvironment(nil)

	class := &Class{
		Name:             "TestClass",
		ParentClasses:    []string{},
		MRO:              []string{"TestClass"},
		PublicVariables:  make(map[string]*Variable),
		PrivateVariables: make(map[string]*Variable),
		PublicFunctions:  make(map[string]*Function),
		PrivateFunctions: make(map[string]*Function),
		SharedVariables:  make(map[string]*Variable),
		SharedFunctions:  make(map[string]*Function),
	}

	class.PublicVariables["name"] = &Variable{
		Name:     "name",
		Type:     "STRIN",
		Value:    StringValue("initial"),
		IsLocked: false,
		IsPublic: true,
	}

	err := env.DefineClass(class)

	instance, err := env.NewObjectInstance("TestClass")
	assert.NoError(t, err)

	// Test getting member using the actual method names
	variable, err := instance.GetMemberVariable("name", "TestClass")
	assert.NoError(t, err)
	assert.Equal(t, StringValue("initial"), variable.Value)

	// Test setting member
	err = instance.SetMemberVariable("name", StringValue("updated"), "TestClass")
	assert.NoError(t, err)

	variable, err = instance.GetMemberVariable("name", "TestClass")
	assert.NoError(t, err)
	assert.Equal(t, StringValue("updated"), variable.Value)

	// Test accessing undefined member
	_, err = instance.GetMemberVariable("undefined", "TestClass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "undefined member")
}
