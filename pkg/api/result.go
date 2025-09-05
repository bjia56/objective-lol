package api

import "github.com/bjia56/objective-lol/pkg/environment"

// ExecutionResult represents the result of executing Objective-LOL code
type ExecutionResult struct {
	// Return value from the execution (e.g., from MAIN function)
	Value interface{}

	// Raw Objective-LOL value (for advanced use)
	RawValue environment.Value

	// Output captured during execution (if configured)
	Output string
}
