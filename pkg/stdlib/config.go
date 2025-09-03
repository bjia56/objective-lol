package stdlib

import (
	"github.com/bjia56/objective-lol/pkg/interpreter"
)

// DefaultStdlibInitializers returns the standard library initializer map
// used by the Objective-LOL interpreter
func DefaultStdlibInitializers() map[string]interpreter.StdlibInitializer {
	return map[string]interpreter.StdlibInitializer{
		"IO":     RegisterIOInEnv,
		"MATH":   RegisterMATHInEnv,
		"STDIO":  RegisterSTDIOInEnv,
		"STRING": RegisterSTRINGInEnv,
		"TEST":   RegisterTESTInEnv,
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
