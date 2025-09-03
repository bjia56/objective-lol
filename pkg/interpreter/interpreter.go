package interpreter

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/modules"
	"github.com/bjia56/objective-lol/pkg/runtime"
	"github.com/bjia56/objective-lol/pkg/types"
)

// Interpreter implements the tree-walking interpreter for Objective-LOL
type Interpreter struct {
	// Global state
	runtime           *environment.RuntimeEnvironment
	moduleResolver    *modules.ModuleResolver      // For resolving and caching module imports
	stdlibInitializer map[string]StdlibInitializer // For loading standard library modules
	ctx               context.Context              // For cancellation and timeout support

	// Interpreter-local state
	environment       *environment.Environment    // Current environment/scope
	currentClass      string                      // For tracking visibility context
	currentObject     *environment.ObjectInstance // For tracking current object instance in method calls
	currentFile       string                      // For tracking current file being processed (for relative imports)
	currentModulePath string                      // For tracking current module context for qualified class names
}

type StdlibInitializer func(*environment.Environment, ...string) error

// NewInterpreter creates a new interpreter instance
func NewInterpreter(stdlib map[string]StdlibInitializer, globals ...StdlibInitializer) *Interpreter {
	runtime := environment.NewRuntimeEnvironment()
	// Use current working directory as base for module resolution
	workingDir, _ := filepath.Abs(".")

	// Register global types
	for _, init := range globals {
		init(runtime.GlobalEnv)
	}

	interpreter := &Interpreter{
		runtime:           runtime,
		moduleResolver:    modules.NewModuleResolver(workingDir),
		stdlibInitializer: stdlib,
		ctx:               context.Background(),
		environment:       runtime.NewLocalEnv(), // module scope
	}

	return interpreter
}

// SetCurrentFile sets the current file being processed for relative import resolution
func (i *Interpreter) SetCurrentFile(filename string) {
	i.currentFile = filename
	// Set module path for qualified class names
	if filename != "" {
		absPath, _ := filepath.Abs(filename)
		i.currentModulePath = fmt.Sprintf("file:%s", absPath)
	}
}

// ForkGlobal creates a new Interpreter from the current Interpreter with only global state
func (i *Interpreter) ForkGlobal() *Interpreter {
	return &Interpreter{
		runtime:           i.runtime,
		moduleResolver:    i.moduleResolver,
		stdlibInitializer: i.stdlibInitializer,
		ctx:               i.ctx,
		environment:       i.runtime.NewLocalEnv(),
	}
}

// ForkAll creates a new Interpreter from the current Interpreter with all state
func (i *Interpreter) ForkAll() *Interpreter {
	return &Interpreter{
		runtime:           i.runtime,
		moduleResolver:    i.moduleResolver,
		stdlibInitializer: i.stdlibInitializer,
		ctx:               i.ctx,
		environment:       i.environment,
		currentClass:      i.currentClass,
		currentObject:     i.currentObject,
		currentFile:       i.currentFile,
		currentModulePath: i.currentModulePath,
	}
}

// Interpret executes the given AST with context support for cancellation
func (i *Interpreter) Interpret(ctx context.Context, program *ast.ProgramNode) (types.Value, error) {
	oldCtx := i.ctx
	i.ctx = ctx
	defer func() { i.ctx = oldCtx }()

	return program.Accept(i)
}

// checkContext checks if the context has been cancelled and returns an error if so
func (i *Interpreter) checkContext() error {
	select {
	case <-i.ctx.Done():
		return i.ctx.Err()
	default:
		return nil
	}
}

// resolveClassName resolves class names to qualified names for type safety
func (i *Interpreter) resolveClassName(name string) string {
	// If already qualified (contains module separator), return as-is
	if strings.Contains(name, ":") || strings.Contains(name, ".") {
		return name
	}
	
	// Check if it's an imported class in current scope
	if class, err := i.environment.GetClass(name); err == nil {
		return class.QualifiedName
	}
	
	// Assume it's in current module
	if i.currentModulePath != "" {
		return fmt.Sprintf("%s.%s", i.currentModulePath, name)
	}
	
	return name // Fallback for legacy/local classes
}

// VisitProgram executes the entire program
func (i *Interpreter) VisitProgram(node *ast.ProgramNode) (types.Value, error) {
	// Pass 0: Process import statements first
	for _, decl := range node.Declarations {
		// Check for context cancellation
		if err := i.checkContext(); err != nil {
			return types.NOTHIN, err
		}

		switch n := decl.(type) {
		case *ast.ImportStatementNode:
			if _, err := i.VisitImportStatement(n); err != nil {
				return types.NOTHIN, err
			}
		}
	}

	// Pass 1: declare all functions and classes (now imports are available)
	for _, decl := range node.Declarations {
		// Check for context cancellation
		if err := i.checkContext(); err != nil {
			return types.NOTHIN, err
		}

		switch n := decl.(type) {
		case *ast.FunctionDeclarationNode:
			if _, err := i.VisitFunctionDeclaration(n); err != nil {
				return types.NOTHIN, err
			}
		case *ast.ClassDeclarationNode:
			if _, err := i.VisitClassDeclaration(n); err != nil {
				return types.NOTHIN, err
			}
		}
	}

	// Pass 2: execute variable declarations and other statements
	for _, decl := range node.Declarations {
		// Check for context cancellation
		if err := i.checkContext(); err != nil {
			return types.NOTHIN, err
		}

		switch n := decl.(type) {
		case *ast.VariableDeclarationNode:
			if _, err := i.VisitVariableDeclaration(n); err != nil {
				return types.NOTHIN, err
			}
		}
	}

	// Look for and execute MAIN function
	if mainFunc, err := i.environment.GetFunction("MAIN"); err == nil {
		return i.callFunction(mainFunc, []types.Value{})
	}

	return types.NOTHIN, nil
}

// VisitImportStatement handles module import statements
func (i *Interpreter) VisitImportStatement(node *ast.ImportStatementNode) (types.Value, error) {
	if node.IsFileImport {
		// Handle file module import
		return i.handleFileImport(node.ModuleName, node.Declarations)
	} else {
		// Handle built-in module import
		return i.handleBuiltinImport(node.ModuleName, node.Declarations)
	}
}

// handleBuiltinImport handles imports of built-in modules (STDIO, MATH, TIME)
func (i *Interpreter) handleBuiltinImport(moduleName string, declarations []string) (types.Value, error) {
	moduleName = strings.ToUpper(moduleName)

	// Load the requested built-in module into the current environment scope
	if init, ok := i.stdlibInitializer[moduleName]; ok {
		err := init(i.environment, declarations...)
		if err != nil {
			return types.NOTHIN, fmt.Errorf("failed to initialize module %s: %v", moduleName, err)
		}
	} else {
		return types.NOTHIN, fmt.Errorf("unknown built-in module: %s", moduleName)
	}

	return types.NOTHIN, nil
}

// handleFileImport handles imports of .olol file modules
func (i *Interpreter) handleFileImport(filePath string, declarations []string) (types.Value, error) {
	// Check for context cancellation
	if err := i.checkContext(); err != nil {
		return types.NOTHIN, err
	}

	// Determine the directory of the current file for relative path resolution
	var importingDir string
	if i.currentFile != "" {
		importingDir = filepath.Dir(i.currentFile)
	}

	// Load the module using the resolver with context-aware path resolution
	moduleAST, resolvedPath, err := i.moduleResolver.LoadModuleFromWithPath(filePath, importingDir)
	if err != nil {
		return types.NOTHIN, fmt.Errorf("failed to load module %s: %v", filePath, err)
	}

	var moduleEnv *environment.Environment

	// Check for circular imports during execution
	if i.moduleResolver.IsInExecutingStack(resolvedPath) {
		return types.NOTHIN, fmt.Errorf("circular import detected during execution: %s", filePath)
	}

	// Check if we have a cached environment for this module
	if cachedEnv, exists := i.moduleResolver.GetCachedEnvironment(resolvedPath); exists {
		moduleEnv = cachedEnv
	} else {
		// Add to execution stack to detect circular imports
		i.moduleResolver.AddToExecutingStack(resolvedPath)
		defer i.moduleResolver.RemoveFromExecutingStack()

		// Create a forked interpreter for module execution with isolated context
		moduleInterpreter := i.ForkGlobal()
		moduleInterpreter.SetCurrentFile(resolvedPath)

		moduleEnv = moduleInterpreter.environment

		// Execute the module using the forked interpreter
		// First pass: declare functions and classes
		for _, decl := range moduleAST.Declarations {
			switch n := decl.(type) {
			case *ast.FunctionDeclarationNode:
				if _, err := moduleInterpreter.VisitFunctionDeclaration(n); err != nil {
					return types.NOTHIN, fmt.Errorf("error declaring function in module %s: %v", filePath, err)
				}
			case *ast.ClassDeclarationNode:
				if _, err := moduleInterpreter.VisitClassDeclaration(n); err != nil {
					return types.NOTHIN, fmt.Errorf("error declaring class in module %s: %v", filePath, err)
				}
			}
		}

		// Second pass: execute variable declarations and import statements
		for _, decl := range moduleAST.Declarations {
			switch n := decl.(type) {
			case *ast.VariableDeclarationNode:
				if _, err := moduleInterpreter.VisitVariableDeclaration(n); err != nil {
					return types.NOTHIN, fmt.Errorf("error declaring variable in module %s: %v", filePath, err)
				}
			case *ast.ImportStatementNode:
				if _, err := moduleInterpreter.VisitImportStatement(n); err != nil {
					return types.NOTHIN, fmt.Errorf("error processing import in module %s: %v", filePath, err)
				}
			}
		}

		// Cache the executed environment
		i.moduleResolver.CacheEnvironment(resolvedPath, moduleEnv)
	}

	// Import the requested declarations from module environment to current environment
	if len(declarations) > 0 {
		// Import only specified declarations
		return i.importSelectiveDeclarations(moduleEnv, declarations, filePath)
	} else {
		// Import all public declarations
		return i.importAllDeclarations(moduleEnv, filePath)
	}
}

// importSelectiveDeclarations imports only the specified declarations from module environment
func (i *Interpreter) importSelectiveDeclarations(moduleEnv *environment.Environment, declarations []string, filePath string) (types.Value, error) {
	for _, declName := range declarations {
		declName = strings.ToUpper(declName)
		if strings.HasPrefix(declName, "_") {
			return types.NOTHIN, fmt.Errorf("importing private declaration %s from module %s is not allowed", declName, filePath)
		}

		// Try to import as function
		if function, err := moduleEnv.GetFunction(declName); err == nil {
			if err := i.environment.DefineFunction(function); err != nil {
				return types.NOTHIN, fmt.Errorf("failed to import function %s from module %s: %v", declName, filePath, err)
			}
			continue
		}

		// Try to import as class
		if class, err := moduleEnv.GetClass(declName); err == nil {
			if err := i.environment.DefineClass(class); err != nil {
				return types.NOTHIN, fmt.Errorf("failed to import class %s from module %s: %v", declName, filePath, err)
			}
			continue
		}

		// Try to import as variable
		if variable, err := moduleEnv.GetVariable(declName); err == nil {
			if err := i.environment.DefineVariable(declName, variable.Type, variable.Value, variable.IsLocked); err != nil {
				return types.NOTHIN, fmt.Errorf("failed to import variable %s from module %s: %v", declName, filePath, err)
			}
			continue
		}

		// Declaration not found
		return types.NOTHIN, fmt.Errorf("declaration %s not found in module %s", declName, filePath)
	}

	return types.NOTHIN, nil
}

// importAllDeclarations imports all public declarations from module environment
func (i *Interpreter) importAllDeclarations(moduleEnv *environment.Environment, filePath string) (types.Value, error) {
	// Import all functions (assuming all are public for now)
	if functions := moduleEnv.GetAllFunctions(); len(functions) > 0 {
		for name, function := range functions {
			// Skip private functions (starting with _)
			if strings.HasPrefix(name, "_") {
				continue
			}
			if err := i.environment.DefineFunction(function); err != nil {
				return types.NOTHIN, fmt.Errorf("failed to import function %s from module %s: %v", name, filePath, err)
			}
		}
	}

	// Import all classes (assuming all are public for now)
	if classes := moduleEnv.GetAllClasses(); len(classes) > 0 {
		for name, class := range classes {
			// Skip private classes (starting with _)
			if strings.HasPrefix(name, "_") {
				continue
			}
			if err := i.environment.DefineClass(class); err != nil {
				return types.NOTHIN, fmt.Errorf("failed to import class %s from module %s: %v", name, filePath, err)
			}
		}
	}

	// Import all variables (assuming all are public for now)
	if variables := moduleEnv.GetAllVariables(); len(variables) > 0 {
		for name, variable := range variables {
			// Skip private variables (starting with _)
			if strings.HasPrefix(name, "_") {
				continue
			}
			if err := i.environment.DefineVariable(name, variable.Type, variable.Value, variable.IsLocked); err != nil {
				return types.NOTHIN, fmt.Errorf("failed to import variable %s from module %s: %v", name, filePath, err)
			}
		}
	}

	return types.NOTHIN, nil
}

// VisitVariableDeclaration handles variable declarations
func (i *Interpreter) VisitVariableDeclaration(node *ast.VariableDeclarationNode) (types.Value, error) {
	var value types.Value = types.NOTHIN
	var err error

	if node.Value != nil {
		value, err = node.Value.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}
	}

	err = i.environment.DefineVariable(
		strings.ToUpper(node.Name),
		strings.ToUpper(node.Type),
		value,
		node.IsLocked,
	)

	return types.NOTHIN, err
}

// VisitFunctionDeclaration handles function declarations
func (i *Interpreter) VisitFunctionDeclaration(node *ast.FunctionDeclarationNode) (types.Value, error) {
	function := &environment.Function{
		Name:        strings.ToUpper(node.Name),
		ReturnType:  strings.ToUpper(node.ReturnType),
		Parameters:  node.Parameters,
		Body:        node.Body,
		IsShared:    node.IsShared,
		ParentClass: i.currentClass,
	}

	// Normalize parameter names and types
	for i := range function.Parameters {
		function.Parameters[i].Name = strings.ToUpper(function.Parameters[i].Name)
		function.Parameters[i].Type = strings.ToUpper(function.Parameters[i].Type)
	}

	return types.NOTHIN, i.environment.DefineFunction(function)
}

// VisitClassDeclaration handles class declarations
func (i *Interpreter) VisitClassDeclaration(node *ast.ClassDeclarationNode) (types.Value, error) {
	// Determine qualified parent class name
	var qualifiedParent string
	if node.ParentClass != "" {
		qualifiedParent = i.resolveClassName(strings.ToUpper(node.ParentClass))
	}
	
	class := environment.NewClass(
		strings.ToUpper(node.Name),
		i.currentModulePath,  // Use current module context
		qualifiedParent,
	)

	// Save current context
	oldClass := i.currentClass
	oldEnv := i.environment

	// Set class context for member processing
	i.currentClass = class.Name
	classEnv := environment.NewEnvironment(oldEnv)
	i.environment = classEnv

	// Process class members
	for _, member := range node.Members {
		// Check for context cancellation
		if err := i.checkContext(); err != nil {
			return types.NOTHIN, err
		}

		if member.IsVariable {
			// Handle member variables
			var value types.Value = types.NOTHIN
			var err error

			if member.Variable.Value != nil {
				value, err = member.Variable.Value.Accept(i)
				if err != nil {
					return types.NOTHIN, err
				}
			}

			variable := &environment.Variable{
				Name:     strings.ToUpper(member.Variable.Name),
				Type:     strings.ToUpper(member.Variable.Type),
				Value:    value,
				IsLocked: member.Variable.IsLocked,
				IsPublic: member.IsPublic, // Use visibility from member declaration
			}

			// Add to appropriate collection based on visibility and sharing
			if member.IsShared {
				class.SharedVariables[variable.Name] = variable
			} else if member.IsPublic {
				class.PublicVariables[variable.Name] = variable
			} else {
				class.PrivateVariables[variable.Name] = variable
			}

		} else {
			// Handle member functions
			function := &environment.Function{
				Name:        strings.ToUpper(member.Function.Name),
				ReturnType:  strings.ToUpper(member.Function.ReturnType),
				Parameters:  member.Function.Parameters,
				Body:        member.Function.Body,
				IsShared:    &member.IsShared,
				ParentClass: class.Name,
			}

			// Normalize parameter names and types
			for i := range function.Parameters {
				function.Parameters[i].Name = strings.ToUpper(function.Parameters[i].Name)
				function.Parameters[i].Type = strings.ToUpper(function.Parameters[i].Type)
			}

			// Add to appropriate collection based on visibility and sharing
			if member.IsShared {
				class.SharedFunctions[function.Name] = function
			} else if member.IsPublic {
				class.PublicFunctions[function.Name] = function
			} else {
				class.PrivateFunctions[function.Name] = function
			}
		}
	}

	// Restore context
	i.currentClass = oldClass
	i.environment = oldEnv

	return types.NOTHIN, i.environment.DefineClass(class)
}

// VisitAssignment handles variable assignments
func (i *Interpreter) VisitAssignment(node *ast.AssignmentNode) (types.Value, error) {
	value, err := node.Value.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	switch target := node.Target.(type) {
	case *ast.IdentifierNode:
		// Simple assignment - check if it's a member variable in method context
		name := strings.ToUpper(target.Name)

		// If we're in a method call and the identifier refers to a member variable, treat as member assignment
		if i.currentObject != nil {
			if _, err := i.currentObject.GetMemberVariable(name, i.currentClass); err == nil {
				return value, i.currentObject.SetMemberVariable(name, value, i.currentClass)
			}
		}

		// Otherwise, treat as local variable assignment
		return value, i.environment.SetVariable(name, value)

	case *ast.MemberAccessNode:
		// Member assignment (obj IN member ITZ value)
		objectValue, err := target.Object.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}

		if objVal, ok := objectValue.(types.ObjectValue); ok {
			if obj, ok := objVal.Instance.(*environment.ObjectInstance); ok {
				return value, obj.SetMemberVariable(strings.ToUpper(target.Member), value, i.currentClass)
			}
		}
		return types.NOTHIN, fmt.Errorf("cannot assign to member of non-object")

	default:
		return types.NOTHIN, fmt.Errorf("invalid assignment target")
	}
}

// VisitIfStatement handles if statements
func (i *Interpreter) VisitIfStatement(node *ast.IfStatementNode) (types.Value, error) {
	condition, err := node.Condition.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	if condition.ToBool() == types.YEZ {
		// Create new environment for then block
		oldEnv := i.environment
		thenEnv := environment.NewEnvironment(oldEnv)
		i.environment = thenEnv

		// Execute then block
		result, err := node.ThenBlock.Accept(i)

		// Restore environment
		i.environment = oldEnv
		return result, err
	} else if node.ElseBlock != nil {
		// Create new environment for else block
		oldEnv := i.environment
		elseEnv := environment.NewEnvironment(oldEnv)
		i.environment = elseEnv

		// Execute else block
		result, err := node.ElseBlock.Accept(i)

		// Restore environment
		i.environment = oldEnv
		return result, err
	}

	return types.NOTHIN, nil
}

// VisitWhileStatement handles while loops
func (i *Interpreter) VisitWhileStatement(node *ast.WhileStatementNode) (types.Value, error) {
	for {
		// Check for context cancellation
		if err := i.checkContext(); err != nil {
			return types.NOTHIN, err
		}

		condition, err := node.Condition.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}

		if condition.ToBool() != types.YEZ {
			break
		}

		// Create new environment for loop body
		oldEnv := i.environment
		loopEnv := environment.NewEnvironment(oldEnv)
		i.environment = loopEnv

		// Execute loop body
		_, err = node.Body.Accept(i)

		// Restore environment
		i.environment = oldEnv

		if err != nil {
			return types.NOTHIN, err
		}
	}

	return types.NOTHIN, nil
}

// VisitReturnStatement handles return statements
func (i *Interpreter) VisitReturnStatement(node *ast.ReturnStatementNode) (types.Value, error) {
	var value types.Value = types.NOTHIN
	var err error

	if node.Value != nil {
		value, err = node.Value.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}
	}

	return types.NOTHIN, runtime.ReturnValue{Value: value}
}

// VisitFunctionCall handles function calls
func (i *Interpreter) VisitFunctionCall(node *ast.FunctionCallNode) (types.Value, error) {
	// Check for context cancellation
	if err := i.checkContext(); err != nil {
		return types.NOTHIN, err
	}

	// Evaluate arguments
	args := make([]types.Value, len(node.Arguments))
	for j, arg := range node.Arguments {
		val, err := arg.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}
		args[j] = val
	}

	switch funcNode := node.Function.(type) {
	case *ast.IdentifierNode:
		// Global function call
		functionName := strings.ToUpper(funcNode.Name)

		// Check local environment first
		if function, err := i.environment.GetFunction(functionName); err == nil {
			return i.callFunction(function, args)
		}

		return types.NOTHIN, fmt.Errorf("undefined function '%s'", functionName)

	case *ast.MemberAccessNode:
		// Member function call (obj DO func WIT args)
		objectValue, err := funcNode.Object.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}

		if objVal, ok := objectValue.(types.ObjectValue); ok {
			if obj, ok := objVal.Instance.(*environment.ObjectInstance); ok {
				function, err := obj.GetMemberFunction(strings.ToUpper(funcNode.Member), i.currentClass, i.environment)
				if err != nil {
					return types.NOTHIN, err
				}
				return i.callMemberFunction(function, obj, args)
			}
		}
		return types.NOTHIN, fmt.Errorf("cannot call method on non-object")

	default:
		return types.NOTHIN, fmt.Errorf("invalid function call target")
	}
}

// callFunction executes a function with the given arguments
func (i *Interpreter) callFunction(function *environment.Function, args []types.Value) (types.Value, error) {
	// Check argument count
	if len(args) != len(function.Parameters) {
		return types.NOTHIN, fmt.Errorf("function '%s' expects %d arguments, got %d",
			function.Name, len(function.Parameters), len(args))
	}

	// Handle native functions
	if function.NativeImpl != nil {
		argsCasted := make([]types.Value, len(args))
		for i, arg := range args {
			casted, err := arg.Cast(function.Parameters[i].Type)
			if err != nil {
				return types.NOTHIN, fmt.Errorf("cannot cast function argument %s to %s: %v",
					function.Parameters[i].Name, function.Parameters[i].Type, err)
			}
			argsCasted[i] = casted
		}

		// Create function context
		ctx := NewFunctionContext(i, i.environment)

		return function.NativeImpl(ctx, i.currentObject, argsCasted)
	}

	// Create new environment for function execution
	funcEnv := environment.NewEnvironment(i.environment)

	// Bind parameters
	for j, param := range function.Parameters {
		err := funcEnv.DefineVariable(param.Name, param.Type, args[j], false)
		if err != nil {
			return types.NOTHIN, err
		}
	}

	// Save and set environment
	oldEnv := i.environment
	i.environment = funcEnv

	// Execute function body
	var result types.Value = types.NOTHIN
	if body, ok := function.Body.(*ast.StatementBlockNode); ok {
		_, err := body.Accept(i)
		if err != nil {
			if runtime.IsReturnValue(err) {
				result = runtime.GetReturnValue(err)
			} else {
				i.environment = oldEnv
				return types.NOTHIN, err
			}
		}
	}

	// Restore environment
	i.environment = oldEnv

	// Cast result to return type if specified
	if function.ReturnType != "" {
		castedResult, err := result.Cast(function.ReturnType)
		if err != nil {
			return types.NOTHIN, fmt.Errorf("cannot cast return value to %s: %v", function.ReturnType, err)
		}
		result = castedResult
	}

	return result, nil
}

// callMemberFunction executes a member function
func (i *Interpreter) callMemberFunction(function *environment.Function, obj *environment.ObjectInstance, args []types.Value) (types.Value, error) {
	memberInterpreter := i.ForkAll()
	memberInterpreter.currentClass = obj.Class.Name
	memberInterpreter.currentObject = obj
	return memberInterpreter.callFunction(function, args)
}

// VisitMemberAccess handles member access
func (i *Interpreter) VisitMemberAccess(node *ast.MemberAccessNode) (types.Value, error) {
	objectValue, err := node.Object.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	if objVal, ok := objectValue.(types.ObjectValue); ok {
		if obj, ok := objVal.Instance.(*environment.ObjectInstance); ok {
			variable, err := obj.GetMemberVariable(strings.ToUpper(node.Member), i.currentClass)
			if err != nil {
				return types.NOTHIN, err
			}
			return variable.Value, nil
		}
	}

	return types.NOTHIN, fmt.Errorf("cannot access member of non-object")
}

// VisitBinaryOp handles binary operations
func (i *Interpreter) VisitBinaryOp(node *ast.BinaryOpNode) (types.Value, error) {
	left, err := node.Left.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	right, err := node.Right.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	switch strings.ToUpper(node.Operator) {
	case "MOAR": // Addition
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.Add(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for addition")

	case "LES": // Subtraction
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.Subtract(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for subtraction")

	case "TIEMZ": // Multiplication
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.Multiply(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for multiplication")

	case "DIVIDEZ": // Division
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				// Check for division by zero
				if (rightNum == types.IntegerValue(0)) || (rightNum == types.DoubleValue(0.0)) {
					return types.NOTHIN, runtime.Exception{Message: "Division by zero"}
				}
				return leftNum.Divide(rightNum), nil
			}
		}
		return types.NOTHIN, runtime.Exception{Message: "Operands must be numbers for division"}

	case "BIGGR THAN": // Greater than
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.GreaterThan(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for comparison")

	case "SMALLR THAN": // Less than
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.LessThan(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for comparison")

	case "SAEM AS": // Equal to
		return left.EqualTo(right)

	case "AN": // Logical AND
		return types.BoolValue(left.ToBool() == types.YEZ && right.ToBool() == types.YEZ), nil

	case "OR": // Logical OR
		return types.BoolValue(left.ToBool() == types.YEZ || right.ToBool() == types.YEZ), nil

	default:
		return types.NOTHIN, fmt.Errorf("unknown binary operator: %s", node.Operator)
	}
}

// VisitUnaryOp handles unary operations
func (i *Interpreter) VisitUnaryOp(node *ast.UnaryOpNode) (types.Value, error) {
	operand, err := node.Operand.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	switch strings.ToUpper(node.Operator) {
	case "NOT":
		return types.BoolValue(operand.ToBool() == types.NO), nil
	default:
		return types.NOTHIN, fmt.Errorf("unknown unary operator: %s", node.Operator)
	}
}

// VisitCast handles type casting
func (i *Interpreter) VisitCast(node *ast.CastNode) (types.Value, error) {
	value, err := node.Expression.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	result, castErr := value.Cast(strings.ToUpper(node.TargetType))
	if castErr != nil {
		return types.NOTHIN, runtime.Exception{Message: castErr.Error()}
	}
	return result, nil
}

// VisitLiteral handles literal values
func (i *Interpreter) VisitLiteral(node *ast.LiteralNode) (types.Value, error) {
	return node.Value, nil
}

// VisitIdentifier handles identifier references
func (i *Interpreter) VisitIdentifier(node *ast.IdentifierNode) (types.Value, error) {
	name := strings.ToUpper(node.Name)

	// If we're in a method call and the identifier might be a member variable, check current object first
	if i.currentObject != nil {
		if memberVar, err := i.currentObject.GetMemberVariable(name, i.currentClass); err == nil {
			return memberVar.Value, nil
		}
	}

	// Check local variables
	if variable, err := i.environment.GetVariable(name); err == nil {
		return variable.Value, nil
	}

	// Check if it's a zero-argument function call
	if function, err := i.environment.GetFunction(name); err == nil {
		return i.callFunction(function, []types.Value{})
	}

	return types.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Undefined variable or function '%s'", name)}
}

// VisitObjectInstantiation handles object creation
func (i *Interpreter) VisitObjectInstantiation(node *ast.ObjectInstantiationNode) (types.Value, error) {
	var env *environment.Environment = i.environment

	instance, err := env.NewObjectInstance(strings.ToUpper(node.ClassName))
	if err != nil {
		return types.NOTHIN, err
	}

	objectValue := types.NewObjectValue(instance, strings.ToUpper(node.ClassName))

	// Check for constructor call
	if obj, ok := instance.(*environment.ObjectInstance); ok {
		// Look for a constructor method with the same name as the class
		constructorName := strings.ToUpper(node.ClassName)
		function, err := obj.GetMemberFunction(constructorName, i.currentClass, i.environment)

		// Only proceed if constructor exists
		if err == nil {
			// Evaluate constructor arguments
			args := make([]types.Value, len(node.ConstructorArgs))
			for j, arg := range node.ConstructorArgs {
				val, err := arg.Accept(i)
				if err != nil {
					return types.NOTHIN, err
				}
				args[j] = val
			}

			// Call the constructor method with arguments
			_, err = i.callMemberFunction(function, obj, args)
			if err != nil {
				return types.NOTHIN, fmt.Errorf("constructor call failed: %v", err)
			}
		} else if len(node.ConstructorArgs) > 0 {
			return types.NOTHIN, fmt.Errorf("default constructor '%s' expects %d arguments, got %d", constructorName, len(function.Parameters), len(node.ConstructorArgs))
		}
	}

	return objectValue, nil
}

// VisitStatementBlock handles statement blocks
func (i *Interpreter) VisitStatementBlock(node *ast.StatementBlockNode) (types.Value, error) {
	var result types.Value = types.NOTHIN

	for _, stmt := range node.Statements {
		// Check for context cancellation
		if err := i.checkContext(); err != nil {
			return types.NOTHIN, err
		}
		val, err := stmt.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}
		result = val
	}

	return result, nil
}

// VisitTryStatement handles try-catch-finally blocks
func (i *Interpreter) VisitTryStatement(node *ast.TryStatementNode) (types.Value, error) {
	// Check for context cancellation
	if err := i.checkContext(); err != nil {
		return types.NOTHIN, err
	}

	var result types.Value = types.NOTHIN
	var tryErr error

	// Execute try block with its own scope
	oldEnv := i.environment
	tryEnv := environment.NewEnvironment(oldEnv)
	i.environment = tryEnv

	result, tryErr = node.TryBody.Accept(i)

	// Restore environment after try block
	i.environment = oldEnv

	// If an exception occurred, handle it with the catch block
	if tryErr != nil && runtime.IsException(tryErr) {
		// Create new environment for catch block with exception variable
		catchEnv := environment.NewEnvironment(i.environment)
		oldCatchEnv := i.environment
		i.environment = catchEnv

		// Bind the exception message to the catch variable
		exceptionMsg := runtime.GetExceptionMessage(tryErr)
		i.environment.DefineVariable(node.CatchVar, "STRIN", types.StringValue(exceptionMsg), false)

		// Execute catch block
		result, tryErr = node.CatchBody.Accept(i)

		// Restore environment
		i.environment = oldCatchEnv
	}

	// Execute finally block if present (always runs) with its own scope
	if node.FinallyBody != nil {
		oldFinallyEnv := i.environment
		finallyEnv := environment.NewEnvironment(oldFinallyEnv)
		i.environment = finallyEnv

		// Finally block runs regardless of exceptions, but doesn't override result
		_, finallyErr := node.FinallyBody.Accept(i)

		// Restore environment
		i.environment = oldFinallyEnv

		// If finally block throws an exception, it takes precedence
		if finallyErr != nil {
			return types.NOTHIN, finallyErr
		}
	}

	return result, tryErr
}

// VisitThrowStatement handles throw statements
func (i *Interpreter) VisitThrowStatement(node *ast.ThrowStatementNode) (types.Value, error) {
	// Evaluate the expression to get the error message
	msgValue, err := node.Expression.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	// Cast to string if needed
	strValue, err := msgValue.Cast("STRIN")
	if err != nil {
		return types.NOTHIN, fmt.Errorf("exception message must be a string, got %s", msgValue.Type())
	}

	// Create and throw the exception
	exceptionMsg := strValue.String()
	return types.NOTHIN, runtime.Exception{Message: exceptionMsg}
}

// GetRuntime returns the runtime environment
func (i *Interpreter) GetRuntime() *environment.RuntimeEnvironment {
	return i.runtime
}

// GetEnvironment returns the current environment
func (i *Interpreter) GetEnvironment() *environment.Environment {
	return i.environment
}

// CallFunction calls a function with the provided arguments (for API use)
func (i *Interpreter) CallFunction(function *environment.Function, args []types.Value) (types.Value, error) {
	return i.callFunction(function, args)
}
