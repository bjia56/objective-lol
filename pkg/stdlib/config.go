package stdlib

import (
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/interpreter"
)

// DefaultStdlibInitializers returns the standard library initializer map
// used by the Objective-LOL interpreter
func DefaultStdlibInitializers() map[string]interpreter.StdlibInitializer {
	return map[string]interpreter.StdlibInitializer{
		"FILE":   RegisterFILEInEnv,
		"IO":     RegisterIOInEnv,
		"MATH":   RegisterMATHInEnv,
		"RANDOM": RegisterRANDOMInEnv,
		"STDIO":  RegisterSTDIOInEnv,
		"STRING": RegisterSTRINGInEnv,
		"TEST":   RegisterTESTInEnv,
		"THREAD": RegisterTHREADInEnv,
		"TIME":   RegisterTIMEInEnv,
	}
}

// DefaultGlobalInitializers returns the global initializers that are
// automatically registered in every interpreter instance
func DefaultGlobalInitializers() []interpreter.StdlibInitializer {
	return []interpreter.StdlibInitializer{
		RegisterArraysInEnv,
		RegisterMapsInEnv,
	}
}

type StdlibDefinitionKind int

const (
	StdlibDefinitionKindFunction StdlibDefinitionKind = iota
	StdlibDefinitionKindClass
)

type StdlibDefinition struct {
	Name string
	Kind StdlibDefinitionKind
	Type string // For functions, this is the return type; for classes, this is always the class name
}

// GetStdlibDefinitions extracts a list of standard library definitions (from the given initializers)
func GetStdlibDefinitions(fromInitializers interpreter.StdlibInitializer) []StdlibDefinition {
	definitions := []StdlibDefinition{}
	env := environment.NewEnvironment(nil)
	fromInitializers(env)
	for name, function := range env.GetAllFunctions() {
		definitions = append(definitions, StdlibDefinition{
			Name: name,
			Kind: StdlibDefinitionKindFunction,
			Type: function.ReturnType,
		})
	}
	for name := range env.GetAllClasses() {
		definitions = append(definitions, StdlibDefinition{
			Name: name,
			Kind: StdlibDefinitionKindClass,
			Type: name,
		})
	}
	return definitions
}
