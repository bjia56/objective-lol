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

// CallMethod calls a method on an object instance with the given arguments
func (ctx *FunctionContext) CallMethod(instance *environment.ObjectInstance, methodName string, fromContext string, args []types.Value) (types.Value, error) {
	method, err := instance.GetMemberFunction(methodName, fromContext, ctx.environment)
	if err != nil {
		return types.NOTHIN, err
	}

	return ctx.interpreter.callMemberFunction(method, instance, args)
}
