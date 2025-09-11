package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/interpreter"
	"github.com/bjia56/objective-lol/pkg/parser"
	"github.com/bjia56/objective-lol/pkg/stdlib"
)

const (
	ForeignModuleNamespace = "foreign:anonymous"
)

// VM represents an Objective-LOL virtual machine instance
type VM struct {
	config      *VMConfig
	interpreter *interpreter.Interpreter
	isStarted   bool
}

// NewVM creates a new VM instance with the given config
func NewVM(config *VMConfig) (*VM, error) {
	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("VM configuration validation error: %v", err)
	}

	vm := &VM{
		config: config,
	}

	if err := vm.initialize(); err != nil {
		return nil, fmt.Errorf("VM initialization error: %v", err)
	}

	return vm, nil
}

// initialize sets up the VM interpreter and runtime
func (vm *VM) initialize() error {
	// Combine default stdlib with custom stdlib
	stdlibInit := stdlib.DefaultStdlibInitializers()
	for name, init := range vm.config.CustomStdlib {
		stdlibInit[name] = init
	}

	// Create interpreter
	vm.interpreter = interpreter.NewInterpreter(
		stdlibInit,
		stdlib.DefaultGlobalInitializers()...,
	)

	// Set working directory
	absWorkingDir, err := filepath.Abs(vm.config.WorkingDirectory)
	if err != nil {
		return NewConfigError("invalid working directory", err)
	}
	vm.interpreter.SetCurrentFile(filepath.Join(absWorkingDir, "main.olol"))

	// Configure I/O for stdlib
	if vm.config.Stdout != os.Stdout {
		stdlib.SetOutput(vm.config.Stdout)
	}
	if vm.config.Stdin != os.Stdin {
		stdlib.SetInput(vm.config.Stdin)
	}

	return nil
}

// GetCompatibilityShim returns a compatibility shim for the VM
func (vm *VM) GetCompatibilityShim() *VMCompatibilityShim {
	return &VMCompatibilityShim{vm: vm}
}

// Execute executes Objective-LOL code from a string
func (vm *VM) Execute(code string) (*ExecutionResult, error) {
	return vm.ExecuteWithContext(context.Background(), code)
}

// ExecuteWithContext executes code with a context for cancellation/timeout
func (vm *VM) ExecuteWithContext(ctx context.Context, code string) (*ExecutionResult, error) {
	if vm.isStarted {
		return nil, NewRuntimeError("VM has already been started; create a new VM instance to run code again", nil)
	}
	vm.isStarted = true

	startTime := time.Now()
	result := &ExecutionResult{}

	// Set up timeout if configured
	if vm.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, vm.config.Timeout)
		defer cancel()
	}

	// Parse the code
	program, err := vm.parseCode(code)
	if err != nil {
		return nil, err
	}

	// Execute with timeout handling
	var value environment.Value

	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(string(debug.Stack()))
				done <- NewRuntimeError(fmt.Sprintf("panic during execution: %v", r), nil)
			}
		}()
		var err error
		value, err = vm.interpreter.Interpret(ctx, program)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, wrapError(err, RuntimeErrorType, "execution failed")
		}
	case <-ctx.Done():
		return nil, NewTimeoutError(time.Since(startTime))
	}

	// Convert result value
	goValue, err := ToGoValue(value)
	if err != nil {
		return nil, wrapError(err, RuntimeErrorType, "could not convert result to Go value")
	}

	result.Value = goValue
	result.RawValue = value

	return result, nil
}

// parseCode parses the given code string into an AST
func (vm *VM) parseCode(code string) (*ast.ProgramNode, error) {
	lexer := parser.NewLexer(code)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Check for parsing errors
	if errors := p.Errors(); len(errors) > 0 {
		var errorMessages []string
		for _, err := range errors {
			errorMessages = append(errorMessages, fmt.Sprintf("%v", err))
		}
		return nil, NewCompileError(
			fmt.Sprintf("parsing failed with %d errors: %s",
				len(errors), strings.Join(errorMessages, "; ")),
			nil,
		)
	}

	return program, nil
}

func (vm *VM) NewObjectInstance(className string) (GoValue, error) {
	class, err := vm.interpreter.GetEnvironment().GetClass(strings.ToUpper(className))
	if err != nil {
		return WrapAny(nil), wrapError(err, RuntimeErrorType, fmt.Sprintf("class %s not found", className))
	}

	instance, err := class.NewObject(vm.interpreter.GetEnvironment())
	if err != nil {
		return WrapAny(nil), wrapError(err, RuntimeErrorType, "could not create object instance")
	}

	goVal, err := ToGoValue(instance)
	if err != nil {
		return WrapAny(nil), wrapError(err, RuntimeErrorType, "could not convert object instance to Go value")
	}

	return goVal, nil
}

// Call calls an Objective-LOL function with the given arguments
func (vm *VM) Call(functionName string, args []GoValue) (GoValue, error) {
	// Convert Go arguments to Objective-LOL values
	ololArgs, err := ConvertArguments(args)
	if err != nil {
		return WrapAny(nil), err
	}

	// Call function through interpreter
	result, err := vm.interpreter.CallFunction(strings.ToUpper(functionName), ololArgs)
	if err != nil {
		return WrapAny(nil), wrapError(err, RuntimeErrorType, "function call failed")
	}

	// Convert result back to Go value
	return ToGoValue(result)
}

// CallMethod calls a method on an Objective-LOL object
func (vm *VM) CallMethod(object GoValue, methodName string, args []GoValue) (GoValue, error) {
	// Convert object to environment.Value
	ololObject, err := FromGoValue(object)
	if err != nil {
		return WrapAny(nil), wrapError(err, RuntimeErrorType, "could not convert object to Objective-LOL value")
	}

	// Ensure it's an object instance
	objInstance, ok := ololObject.(*environment.ObjectInstance)
	if !ok {
		return WrapAny(nil), NewRuntimeError("provided value is not an Objective-LOL object", nil)
	}

	// Convert Go arguments to Objective-LOL values
	ololArgs, err := ConvertArguments(args)
	if err != nil {
		return WrapAny(nil), err
	}

	// Call method through interpreter
	result, err := vm.interpreter.CallMemberFunction(objInstance, strings.ToUpper(methodName), ololArgs)
	if err != nil {
		return WrapAny(nil), wrapError(err, RuntimeErrorType, "method call failed")
	}

	// Convert result back to Go value
	return ToGoValue(result)
}

// DefineVariable defines a global variable in the VM
func (vm *VM) DefineVariable(name string, value GoValue, constant bool) error {
	// Convert Go value to Objective-LOL value
	ololValue, err := FromGoValue(value)
	if err != nil {
		return err
	}

	// Define variable in the global environment
	return vm.interpreter.GetEnvironment().DefineVariable(strings.ToUpper(name), ololValue.Type(), ololValue, constant, nil)
}

// SetVariable sets a variable in the global environment
func (vm *VM) SetVariable(variableName string, value GoValue) error {
	// Convert Go value to Objective-LOL value
	ololValue, err := FromGoValue(value)
	if err != nil {
		return err
	}

	// Set the variable value in the global environment
	env := vm.interpreter.GetEnvironment()
	variableName = strings.ToUpper(variableName)

	// Check if variable already exists
	if _, err := env.GetVariable(variableName); err != nil {
		// Variable doesn't exist, create it
		return env.DefineVariable(variableName, ololValue.Type(), ololValue, false, nil)
	} else {
		// Variable exists, update its value
		return env.SetVariable(variableName, ololValue)
	}
}

// Get gets a variable from the global environment
func (vm *VM) GetVariable(variableName string) (GoValue, error) {
	// Get variable from global environment
	variable, err := vm.interpreter.GetEnvironment().GetVariable(strings.ToUpper(variableName))
	if err != nil {
		return WrapAny(nil), NewRuntimeError(
			fmt.Sprintf("variable %s not found", variableName),
			nil,
		)
	}

	// Convert to Go value
	return ToGoValue(variable.Value)
}

// DefineFunction defines a global function in the VM
func (vm *VM) DefineFunction(name string, argc int, function func(args []GoValue) (GoValue, error)) error {
	return vm.interpreter.GetEnvironment().DefineFunction(&environment.Function{
		Name:       strings.ToUpper(name),
		Parameters: make([]environment.Parameter, argc),
		NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
			// Convert environment.Value args to GoValue
			goArgs := make([]GoValue, len(args))
			for i, arg := range args {
				goVal, err := ToGoValue(arg)
				if err != nil {
					return nil, fmt.Errorf("error converting argument %d: %v", i, err)
				}
				goArgs[i] = goVal
			}

			// Call the provided Go function
			result, err := function(goArgs)
			if err != nil {
				return nil, err
			}

			// Convert result back to environment.Value
			envVal, err := FromGoValue(result)
			if err != nil {
				return nil, fmt.Errorf("error converting return value: %v", err)
			}
			return envVal, nil
		},
	})
}

type ClassVariable struct {
	Name   string
	Value  GoValue
	Locked bool
	Getter func(this GoValue) (GoValue, error)
	Setter func(this GoValue, value GoValue) error
}

type ClassMethod struct {
	Name     string
	Argc     int
	Function func(this GoValue, args []GoValue) (GoValue, error)
}

type UnknownFunctionHandler struct {
	Handler func(this GoValue, functionName string, fromContext string, args []GoValue) (GoValue, error)
}

type ClassDefinition struct {
	Name                   string
	PublicVariables        map[string]*ClassVariable
	PrivateVariables       map[string]*ClassVariable
	SharedVariables        map[string]*ClassVariable
	PublicMethods          map[string]*ClassMethod
	PrivateMethods         map[string]*ClassMethod
	UnknownFunctionHandler *UnknownFunctionHandler
}

func NewClassDefinition() *ClassDefinition {
	return &ClassDefinition{
		PublicVariables:  make(map[string]*ClassVariable),
		PrivateVariables: make(map[string]*ClassVariable),
		SharedVariables:  make(map[string]*ClassVariable),
		PublicMethods:    make(map[string]*ClassMethod),
		PrivateMethods:   make(map[string]*ClassMethod),
	}
}

// convertClassVariable converts a ClassVariable to a MemberVariable
func convertClassVariable(name string, classVar *ClassVariable, isPublic bool) (*environment.MemberVariable, error) {
	memberVar := &environment.MemberVariable{
		Variable: environment.Variable{
			Name:     strings.ToUpper(name),
			IsLocked: classVar.Locked,
			IsPublic: isPublic,
		},
	}

	// Set initial value if provided
	if classVar.Value.Get() != nil {
		ololValue, err := FromGoValue(classVar.Value)
		if err != nil {
			return nil, fmt.Errorf("error converting initial value for variable %s: %v", name, err)
		}
		memberVar.Value = ololValue
	} else {
		memberVar.Value = environment.NOTHIN
	}

	// Set custom getter/setter if provided
	if classVar.Getter != nil || classVar.Setter != nil {
		getter := classVar.Getter
		setter := classVar.Setter

		if getter != nil {
			memberVar.NativeGet = func(this *environment.ObjectInstance) (environment.Value, error) {
				thisGoValue, err := ToGoValue(this)
				if err != nil {
					return nil, err
				}
				result, err := getter(thisGoValue)
				if err != nil {
					return nil, err
				}
				return FromGoValue(result)
			}
		}

		if setter != nil {
			memberVar.NativeSet = func(this *environment.ObjectInstance, value environment.Value) error {
				thisGoValue, err := ToGoValue(this)
				if err != nil {
					return err
				}
				goValue, err := ToGoValue(value)
				if err != nil {
					return err
				}
				return setter(thisGoValue, goValue)
			}
		}
	}

	return memberVar, nil
}

// convertClassMethod converts a ClassMethod to a Function
func convertClassMethod(name string, method *ClassMethod) *environment.Function {
	return &environment.Function{
		Name:       strings.ToUpper(name),
		Parameters: make([]environment.Parameter, method.Argc),
		NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
			// Convert 'this' and args to Go values
			thisGoValue, err := ToGoValue(this)
			if err != nil {
				return nil, fmt.Errorf("error converting 'this' object: %v", err)
			}

			goArgs := make([]GoValue, len(args))
			for i, arg := range args {
				goVal, err := ToGoValue(arg)
				if err != nil {
					return nil, fmt.Errorf("error converting argument %d: %v", i, err)
				}
				goArgs[i] = goVal
			}

			// Call the method
			result, err := method.Function(thisGoValue, goArgs)
			if err != nil {
				return nil, err
			}

			// Convert result back to environment.Value
			return FromGoValue(result)
		},
	}
}

func (vm *VM) DefineClass(classDef *ClassDefinition) error {
	// Create the Objective-LOL class
	class := environment.NewClass(strings.ToUpper(classDef.Name), ForeignModuleNamespace, nil)

	// Convert and add public variables
	for name, classVar := range classDef.PublicVariables {
		memberVar, err := convertClassVariable(name, classVar, true)
		if err != nil {
			return err
		}
		class.PublicVariables[strings.ToUpper(name)] = memberVar
	}

	// Convert and add private variables
	for name, classVar := range classDef.PrivateVariables {
		memberVar, err := convertClassVariable(name, classVar, false)
		if err != nil {
			return err
		}
		class.PrivateVariables[strings.ToUpper(name)] = memberVar
	}

	// Convert and add shared variables
	for name, classVar := range classDef.SharedVariables {
		memberVar, err := convertClassVariable(name, classVar, true)
		if err != nil {
			return err
		}
		class.SharedVariables[strings.ToUpper(name)] = memberVar
	}

	// Convert and add public methods
	for name, method := range classDef.PublicMethods {
		function := convertClassMethod(name, method)
		class.PublicFunctions[strings.ToUpper(name)] = function
	}

	// Convert and add private methods
	for name, method := range classDef.PrivateMethods {
		function := convertClassMethod(name, method)
		class.PrivateFunctions[strings.ToUpper(name)] = function
	}

	// Set unknown function handler if provided
	if classDef.UnknownFunctionHandler != nil {
		handler := classDef.UnknownFunctionHandler.Handler
		class.UnknownFunctionHandler = func(fnName string, fromContext string) (*environment.Function, error) {
			return &environment.Function{
				Name:       strings.ToUpper(fnName),
				Parameters: []environment.Parameter{}, // Variadic parameters can be handled in the native impl
				IsVarargs:  true,
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					// Convert 'this' and args to Go values
					thisGoValue, err := ToGoValue(this)
					if err != nil {
						return nil, fmt.Errorf("error converting 'this' object: %v", err)
					}

					goArgs := make([]GoValue, len(args))
					for i, arg := range args {
						goVal, err := ToGoValue(arg)
						if err != nil {
							return nil, fmt.Errorf("error converting argument %d: %v", i, err)
						}
						goArgs[i] = goVal
					}

					// Call the unknown function handler
					result, err := handler(thisGoValue, fnName, fromContext, goArgs)
					if err != nil {
						return nil, err
					}

					// Convert result back to environment.Value
					return FromGoValue(result)
				},
			}, nil
		}
	}

	// Define the class in the environment
	return vm.interpreter.GetEnvironment().DefineClass(class)
}
