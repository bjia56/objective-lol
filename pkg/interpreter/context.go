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

// NewFunctionContext creates a new FunctionContext with the given interpreter and environment
func NewFunctionContext(interp *Interpreter, env *environment.Environment) *FunctionContext {
	return &FunctionContext{
		interpreter: interp,
		environment: env,
	}
}

// Fork creates a new FunctionContext with a forked interpreter and environment
func (ctx *FunctionContext) Fork() *FunctionContext {
	return NewFunctionContext(ctx.interpreter.ForkAll(), ctx.environment)
}

// CallMethod calls a method on an object instance with the given arguments
func (ctx *FunctionContext) CallMethod(instance *environment.ObjectInstance, methodName string, fromContext string, args []types.Value) (types.Value, error) {
	method, err := instance.GetMemberFunction(methodName, fromContext, ctx.environment)
	if err != nil {
		return types.NOTHIN, err
	}

	return ctx.interpreter.callMemberFunction(method, instance, args)
}
