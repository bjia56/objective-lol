package api

import (
	"time"

	"github.com/bjia56/objective-lol/pkg/types"
)

// ExecutionResult represents the result of executing Objective-LOL code
type ExecutionResult struct {
	// Return value from the execution (e.g., from MAIN function)
	Value interface{}
	
	// Raw Objective-LOL value (for advanced use)
	RawValue types.Value
	
	// Execution time
	Duration time.Duration
	
	// Output captured during execution (if configured)
	Output string
	
	// Any warnings or non-fatal issues
	Warnings []string
	
	// Debug information (if debug mode enabled)
	DebugInfo *DebugInfo
}

// DebugInfo contains debug information from execution
type DebugInfo struct {
	// Function calls made during execution
	FunctionCalls []FunctionCall
	
	// Memory usage statistics
	MemoryStats MemoryStats
	
	// Performance metrics
	PerformanceMetrics PerformanceMetrics
}

// FunctionCall represents a function call in the debug trace
type FunctionCall struct {
	Name      string
	Arguments []interface{}
	ReturnValue interface{}
	Duration  time.Duration
	Source    *SourceLocation
}

// MemoryStats contains memory usage information
type MemoryStats struct {
	PeakUsage    int64 // Peak memory usage in bytes
	CurrentUsage int64 // Current memory usage in bytes
	Allocations  int64 // Number of allocations
}

// PerformanceMetrics contains performance measurements
type PerformanceMetrics struct {
	ParseTime     time.Duration
	ExecutionTime time.Duration
	TotalTime     time.Duration
	
	// Instruction counts
	FunctionCalls int
	Operations    int
	Comparisons   int
}