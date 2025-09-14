package stdlib

import (
	"strings"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/interpreter"
)

// DefaultStdlibInitializers returns the standard library initializer map
// used by the Objective-LOL interpreter
func DefaultStdlibInitializers() map[string]interpreter.StdlibInitializer {
	return map[string]interpreter.StdlibInitializer{
		"CACHE":   RegisterCACHEInEnv,
		"FILE":    RegisterFILEInEnv,
		"HTTP":    RegisterHTTPInEnv,
		"IO":      RegisterIOInEnv,
		"MATH":    RegisterMATHInEnv,
		"PROCESS": RegisterPROCESSInEnv,
		"RANDOM":  RegisterRANDOMInEnv,
		"SOCKET":  RegisterSOCKETInEnv,
		"STDIO":   RegisterSTDIOInEnv,
		"STRING":  RegisterSTRINGInEnv,
		"SYSTEM":  RegisterSYSTEMInEnv,
		"TEST":    RegisterTESTInEnv,
		"THREAD":  RegisterTHREADInEnv,
		"TIME":    RegisterTIMEInEnv,
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
	StdlibDefinitionKindVariable
)

type StdlibDefinition struct {
	Name string
	Kind StdlibDefinitionKind
	Type string // For functions, this is the return type; for classes, this is always the class name
	Docs []string
}

// GetStdlibDefinitions extracts a list of standard library definitions (from the given initializers)
func GetStdlibDefinitions(fromInitializers interpreter.StdlibInitializer) []StdlibDefinition {
	definitions := []StdlibDefinition{}
	env := environment.NewEnvironment(nil)
	fromInitializers(env)
	for name, function := range env.GetAllFunctions() {
		if strings.HasPrefix(name, "stdlib:") {
			continue
		}
		definitions = append(definitions, StdlibDefinition{
			Name: name,
			Kind: StdlibDefinitionKindFunction,
			Type: function.ReturnType,
			Docs: function.Documentation,
		})
	}
	for name, class := range env.GetAllClasses() {
		if strings.HasPrefix(name, "stdlib:") {
			continue
		}
		definitions = append(definitions, StdlibDefinition{
			Name: name,
			Kind: StdlibDefinitionKindClass,
			Type: name,
			Docs: class.Documentation,
		})
		
		// Add class methods (public functions)
		for methodName, method := range class.PublicFunctions {
			definitions = append(definitions, StdlibDefinition{
				Name: name + "." + methodName,
				Kind: StdlibDefinitionKindFunction,
				Type: method.ReturnType,
				Docs: method.Documentation,
			})
		}
		
		// Add class member variables
		for varName, variable := range class.PublicVariables {
			definitions = append(definitions, StdlibDefinition{
				Name: name + "." + varName,
				Kind: StdlibDefinitionKindVariable,
				Type: variable.Type,
				Docs: variable.Documentation,
			})
		}
	}
	for name, variable := range env.GetAllVariables() {
		if strings.HasPrefix(name, "stdlib:") {
			continue
		}
		definitions = append(definitions, StdlibDefinition{
			Name: name,
			Kind: StdlibDefinitionKindVariable,
			Type: variable.Type,
			Docs: variable.Documentation,
		})
	}
	return definitions
}

// GetModuleCategories returns the ordered list of categories for a given module
func GetModuleCategories(moduleName string) []string {
	switch moduleName {
	case "MATH":
		return moduleMATHCategories
	case "STDIO":
		return moduleSTDIOCategories
	case "ARRAYS":
		return moduleArraysCategories
	case "MAPS":
		return moduleMapsCategories
	case "FILE":
		return moduleFileCategories
	case "HTTP":
		return moduleHTTPCategories
	case "IO":
		return moduleIOCategories
	case "SOCKET":
		return moduleSocketCategories
	case "PROCESS":
		return moduleProcessCategories
	case "SYSTEM":
		return moduleSystemCategories
	case "CACHE":
		return moduleCacheCategories
	case "RANDOM":
		return moduleRandomCategories
	case "TEST":
		return moduleTestCategories
	case "THREAD":
		return moduleThreadCategories
	default:
		return nil
	}
}
