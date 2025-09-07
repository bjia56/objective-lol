package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/interpreter"
	"github.com/bjia56/objective-lol/pkg/parser"
	"github.com/bjia56/objective-lol/pkg/runtime"
	"github.com/bjia56/objective-lol/pkg/stdlib"
)

// VM represents an Objective-LOL virtual machine instance
type VM struct {
	config      *VMConfig
	interpreter *interpreter.Interpreter
	mutex       sync.RWMutex // For thread safety

	// State tracking
	isInitialized bool

	// Output capture
	outputBuffer *bytes.Buffer
	originalOut  io.Writer
}

// NewVM creates a new VM instance with the given config
func NewVM(config *VMConfig) (*VM, error) {
	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("VM configuration validation error: %v", err)
	}

	vm := &VM{
		config:        config,
		outputBuffer:  &bytes.Buffer{},
		originalOut:   config.Stdout,
		isInitialized: false,
	}

	if err := vm.initialize(); err != nil {
		return nil, fmt.Errorf("VM initialization error: %v", err)
	}

	return vm, nil
}

// initialize sets up the VM interpreter and runtime
func (vm *VM) initialize() error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	if vm.isInitialized {
		return nil
	}

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

	vm.isInitialized = true
	return nil
}

// Execute executes Objective-LOL code from a string
func (vm *VM) Execute(code string) (*ExecutionResult, error) {
	return vm.ExecuteWithContext(context.Background(), code)
}

// ExecuteWithContext executes code with a context for cancellation/timeout
func (vm *VM) ExecuteWithContext(ctx context.Context, code string) (*ExecutionResult, error) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

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

	// Set up output capture if needed
	var outputBuf bytes.Buffer
	if vm.config.Stdout != vm.originalOut {
		stdlib.SetOutput(&outputBuf)
		defer stdlib.SetOutput(vm.originalOut)
	}

	// Execute with timeout handling
	var value environment.Value

	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
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
	result.Output = outputBuf.String()

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

// Call calls an Objective-LOL function with the given arguments
func (vm *VM) Call(functionName string, args []GoValue) (GoValue, error) {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

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
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

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
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

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
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

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
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

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

// defineFunctionUnsafe defines a global function in the VM without holding the mutex
func (vm *VM) defineFunctionUnsafe(name string, argc int, function func(args []GoValue) (GoValue, error)) error {
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

// DefineFunction defines a global function in the VM
func (vm *VM) DefineFunction(name string, argc int, function func(args []GoValue) (GoValue, error)) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	return vm.defineFunctionUnsafe(name, argc, function)
}

// DefineFunctionMaxCompat defines a global function with maximum compatibility,
// wrapping arguments and return values as JSON strings.
// An optional id cookie is passed back to the function to identify it.
// jsonArgs is a JSON array string of the arguments.
// The function should return a JSON object string with "result" and "error" fields.
func (vm *VM) DefineFunctionMaxCompat(id, name string, argc int, function func(id, jsonArgs string) string) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	fn := func(args []GoValue) (GoValue, error) {
		if len(args) != argc {
			return WrapAny(nil), fmt.Errorf("expected %d arguments, got %d", argc, len(args))
		}

		// Convert args to JSON array string
		jsonBytes, err := json.Marshal(args)
		if err != nil {
			return WrapAny(nil), fmt.Errorf("error marshaling arguments to JSON: %v", err)
		}

		// Call the provided function
		jsonResult := function(id, string(jsonBytes))

		// Parse JSON result
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(jsonResult), &result); err != nil {
			return WrapAny(nil), fmt.Errorf("error unmarshaling result from JSON: %v", err)
		}
		if errVal, ok := result["error"]; ok && errVal != nil {
			return WrapAny(nil), runtime.Exception{Message: fmt.Sprintf("%v", errVal)}
		}
		if resultVal, ok := result["result"]; ok {
			return WrapAny(resultVal), nil
		}
		return WrapAny(nil), nil
	}

	return vm.defineFunctionUnsafe(name, argc, fn)
}
