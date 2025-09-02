package interpreter

import (
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// FunctionContext provides native functions with access to interpreter capabilities
type FunctionContext struct {
	interpreter *Interpreter
	environment *environment.Environment
}

// CallFunction calls a script-defined function by name with the given arguments
func (ctx *FunctionContext) CallFunction(name string, args []types.Value) (types.Value, error) {
	function, err := ctx.environment.GetFunction(name)
	if err != nil {
		return types.NOTHIN, err
	}
	return ctx.interpreter.callFunction(function, args)
}

// GetVariable retrieves a variable value from the current environment
func (ctx *FunctionContext) GetVariable(name string) (types.Value, error) {
	variable, err := ctx.environment.GetVariable(name)
	if err != nil {
		return types.NOTHIN, err
	}
	return variable.Value, nil
}

// SetVariable sets a variable value in the current environment
func (ctx *FunctionContext) SetVariable(name string, value types.Value) error {
	return ctx.environment.SetVariable(name, value)
}

// DefineVariable defines a new variable in the current environment
func (ctx *FunctionContext) DefineVariable(name, varType string, value types.Value, isLocked bool) error {
	return ctx.environment.DefineVariable(name, varType, value, isLocked)
}

// GetClass retrieves a class definition from the current environment
func (ctx *FunctionContext) GetClass(name string) (*environment.Class, error) {
	return ctx.environment.GetClass(name)
}

// NewObjectInstance creates a new instance of the specified class
func (ctx *FunctionContext) NewObjectInstance(className string) (interface{}, error) {
	return ctx.environment.NewObjectInstance(className)
}