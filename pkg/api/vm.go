package api

import (
	"bytes"
	"context"
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
	"github.com/bjia56/objective-lol/pkg/stdlib"
)

// VM represents an Objective-LOL virtual machine instance
type VM struct {
	config      *VMConfig
	interpreter *interpreter.Interpreter
	mutex       sync.RWMutex // For thread safety

	// State tracking
	isInitialized bool
	currentFile   string

	// Output capture
	outputBuffer *bytes.Buffer
	originalOut  io.Writer
}

// NewVM creates a new VM instance with the given options
func NewVM(options ...VMOption) *VM {
	config := DefaultConfig()

	// Apply options
	for _, option := range options {
		if err := option(config); err != nil {
			// In a real implementation, you might want to handle this differently
			panic(fmt.Sprintf("VM configuration error: %v", err))
		}
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		panic(fmt.Sprintf("VM configuration validation error: %v", err))
	}

	vm := &VM{
		config:        config,
		outputBuffer:  &bytes.Buffer{},
		originalOut:   config.Stdout,
		isInitialized: false,
	}

	if err := vm.initialize(); err != nil {
		panic(fmt.Sprintf("VM initialization error: %v", err))
	}

	return vm
}

// initialize sets up the VM interpreter and runtime
func (vm *VM) initialize() error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	return vm.initializeUnsafe()
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

// ExecuteFile executes Objective-LOL code from a file
func (vm *VM) ExecuteFile(filename string) (*ExecutionResult, error) {
	return vm.ExecuteFileWithContext(context.Background(), filename)
}

// ExecuteFileWithContext executes a file with context
func (vm *VM) ExecuteFileWithContext(ctx context.Context, filename string) (*ExecutionResult, error) {
	// Validate file extension
	if !strings.HasSuffix(strings.ToLower(filename), ".olol") {
		return nil, NewConfigError("file must have .olol extension", nil)
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, NewConfigError(fmt.Sprintf("file not found: %s", filename), err)
	}

	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, NewConfigError(fmt.Sprintf("error reading file: %s", filename), err)
	}

	// Update current file for relative imports
	absFilename, err := filepath.Abs(filename)
	if err != nil {
		return nil, NewConfigError("could not resolve absolute path", err)
	}

	vm.mutex.Lock()
	oldFile := vm.currentFile
	vm.currentFile = absFilename
	vm.interpreter.SetCurrentFile(absFilename)
	vm.mutex.Unlock()

	defer func() {
		vm.mutex.Lock()
		vm.currentFile = oldFile
		vm.interpreter.SetCurrentFile(oldFile)
		vm.mutex.Unlock()
	}()

	return vm.ExecuteWithContext(ctx, string(content))
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
func (vm *VM) Call(functionName string, args ...interface{}) (interface{}, error) {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	// Convert Go arguments to Objective-LOL values
	ololArgs, err := ConvertArguments(args)
	if err != nil {
		return nil, err
	}

	// Get function from environment
	function, err := vm.interpreter.GetEnvironment().GetFunction(strings.ToUpper(functionName))
	if err != nil {
		return nil, NewRuntimeError(
			fmt.Sprintf("function %s not found", functionName),
			nil,
		)
	}

	// Call function through interpreter
	result, err := vm.interpreter.CallFunction(function, ololArgs)
	if err != nil {
		return nil, wrapError(err, RuntimeErrorType, "function call failed")
	}

	// Convert result back to Go value
	return ToGoValue(result)
}

// Set sets a variable in the global environment
func (vm *VM) Set(variableName string, value interface{}) error {
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
		return env.DefineVariable(variableName, ololValue.Type(), ololValue, false)
	} else {
		// Variable exists, update its value
		return env.SetVariable(variableName, ololValue)
	}
}

// Get gets a variable from the global environment
func (vm *VM) Get(variableName string) (interface{}, error) {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	// Get variable from global environment
	variable, err := vm.interpreter.GetEnvironment().GetVariable(strings.ToUpper(variableName))
	if err != nil {
		return nil, NewRuntimeError(
			fmt.Sprintf("variable %s not found", variableName),
			nil,
		)
	}

	// Convert to Go value
	return ToGoValue(variable.Value)
}

// Reset resets the VM to its initial state
func (vm *VM) Reset() error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	vm.isInitialized = false
	return vm.initializeUnsafe()
}

// initializeUnsafe sets up the VM interpreter and runtime without acquiring the lock
func (vm *VM) initializeUnsafe() error {
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

// GetConfig returns a copy of the current configuration
func (vm *VM) GetConfig() VMConfig {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	// Return a copy to prevent external modification
	configCopy := *vm.config
	return configCopy
}

// GetEnvironment returns the current global environment (for advanced use)
func (vm *VM) GetEnvironment() *environment.Environment {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	return vm.interpreter.GetEnvironment()
}
